package cli

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/edgedelta/edx/internal/config"
)

const (
	ghOwner = "edgedelta"
	ghRepo  = "edx"

	// updateCheckInterval bounds how often the passive startup notice hits the
	// network; between checks it reads the on-disk cache.
	updateCheckInterval = 24 * time.Hour
	updateCacheFile     = "update-check.json"

	// envNoUpdateCheck disables the passive "update available" startup notice.
	envNoUpdateCheck = "EDX_NO_UPDATE_CHECK"

	// maxDownloadBytes caps release downloads defensively.
	maxDownloadBytes = 200 << 20
)

// ghRelease is the subset of the GitHub Releases API response edx needs.
type ghRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
	Assets  []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

// newUpdateCmd builds `edx update`, the explicit self-upgrade command. The
// passive startup notice (maybeNotifyUpdate) only ever points users here; the
// interactive prompt and the actual install live in this command so automation
// never blocks on them.
func newUpdateCmd() *cobra.Command {
	var checkOnly bool
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update edx to the latest release",
		Long: `Update edx to the latest published release.

edx checks GitHub for the newest release and, if you are behind, replaces the
running binary in place after verifying the download against the release
checksums. Homebrew installs are upgraded with "brew upgrade" instead, so the
package manager stays in charge of its own files.

Use --check to only report whether an update is available, without installing.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(cmd.Context(), checkOnly)
		},
	}
	cmd.Flags().BoolVar(&checkOnly, "check", false, "only check for a newer version; do not install")
	return cmd
}

func runUpdate(ctx context.Context, checkOnly bool) error {
	rel, err := fetchLatestRelease(ctx, 15*time.Second)
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}
	writeUpdateCache(rel.TagName)

	if Version == "dev" {
		notef("this is a dev build; the latest release is %s", rel.TagName)
		notef("install a release with `brew upgrade edx` or `go install github.com/edgedelta/edx@latest`")
		return nil
	}
	if !isNewer(rel.TagName, Version) {
		fmt.Fprintf(os.Stderr, "%s edx is up to date (%s)\n", okMark(), Version)
		return nil
	}

	newVer := strings.TrimPrefix(rel.TagName, "v")
	fmt.Fprintf(os.Stderr, "A new version is available: %s (you have %s)\n", newVer, Version)
	if checkOnly {
		notef("run `edx update` to install it")
		return nil
	}

	// Homebrew installs must go through brew, not have their symlinked binary
	// swapped out from under the package manager.
	if _, ok := homebrewInstall(); ok {
		return updateViaHomebrew(ctx)
	}

	if !confirm(fmt.Sprintf("Update edx to %s now?", newVer)) {
		return fmt.Errorf("update cancelled")
	}
	return selfReplace(ctx, rel)
}

// maybeNotifyUpdate prints a one-line "update available" hint to stderr when a
// newer release exists. It is deliberately unobtrusive and safe for automation:
// it never prompts, never writes to stdout, and returns immediately unless
// stderr is an interactive terminal — so AI-driven, piped and CI invocations
// never see it and are never delayed. The version lookup is cached for
// updateCheckInterval; only the first interactive run per interval touches the
// network, with a tight timeout.
func maybeNotifyUpdate(cmd *cobra.Command) {
	if Version == "dev" || os.Getenv(envNoUpdateCheck) != "" {
		return
	}
	if !fileIsTTY(os.Stderr) {
		return
	}
	switch cmd.Name() {
	case "update", "version", "help", "completion", "__complete", "__completeNoDesc":
		return
	}

	latest, ok := cachedLatestVersion()
	if !ok {
		return
	}
	if isNewer(latest, Version) {
		warnf("edx %s is available (you have %s) — run `edx update`", strings.TrimPrefix(latest, "v"), Version)
	}
}

// cachedLatestVersion returns the latest release tag, using the on-disk cache
// when it is younger than updateCheckInterval and otherwise doing a short,
// bounded network check that refreshes the cache. Any error yields ok=false so
// the caller stays silent.
func cachedLatestVersion() (string, bool) {
	if c, ok := readUpdateCache(); ok && time.Since(c.CheckedAt) < updateCheckInterval {
		return c.Latest, true
	}
	rel, err := fetchLatestRelease(context.Background(), 2*time.Second)
	if err != nil {
		return "", false
	}
	writeUpdateCache(rel.TagName)
	return rel.TagName, true
}

// fetchLatestRelease queries the GitHub Releases API for the newest release,
// bounded by timeout.
func fetchLatestRelease(ctx context.Context, timeout time.Duration) (*ghRelease, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", ghOwner, ghRepo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "edx/"+Version)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %s", resp.Status)
	}
	var r ghRelease
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&r); err != nil {
		return nil, err
	}
	if r.TagName == "" {
		return nil, fmt.Errorf("latest release has no tag_name")
	}
	return &r, nil
}

// --- version comparison -----------------------------------------------------

// parseSemver parses "v1.2.3" (leading "v" and any pre-release/build suffix are
// optional) into numeric major, minor and patch. It returns ok=false for
// non-release strings like "dev" or a git-describe value with a commit suffix,
// which must not be compared.
func parseSemver(s string) (major, minor, patch int, ok bool) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "v")
	if i := strings.IndexAny(s, "-+"); i >= 0 {
		s = s[:i]
	}
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return 0, 0, 0, false
	}
	var nums [3]int
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return 0, 0, 0, false
		}
		nums[i] = n
	}
	return nums[0], nums[1], nums[2], true
}

// isNewer reports whether latest is a strictly higher release than current.
// It returns false unless both are clean release versions.
func isNewer(latest, current string) bool {
	lMa, lMi, lPa, ok1 := parseSemver(latest)
	cMa, cMi, cPa, ok2 := parseSemver(current)
	if !ok1 || !ok2 {
		return false
	}
	switch {
	case lMa != cMa:
		return lMa > cMa
	case lMi != cMi:
		return lMi > cMi
	default:
		return lPa > cPa
	}
}

// --- on-disk cache ----------------------------------------------------------

type updateCache struct {
	CheckedAt time.Time `json:"checked_at"`
	Latest    string    `json:"latest"`
}

func updateCachePath() (string, error) {
	dir, err := config.StateDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, updateCacheFile), nil
}

func readUpdateCache() (updateCache, bool) {
	var c updateCache
	p, err := updateCachePath()
	if err != nil {
		return c, false
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return c, false
	}
	if json.Unmarshal(data, &c) != nil || c.Latest == "" {
		return c, false
	}
	return c, true
}

func writeUpdateCache(latest string) {
	p, err := updateCachePath()
	if err != nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return
	}
	data, err := json.Marshal(updateCache{CheckedAt: time.Now().UTC(), Latest: latest})
	if err != nil {
		return
	}
	_ = os.WriteFile(p, data, 0o600)
}

// --- Homebrew path ----------------------------------------------------------

// homebrewInstall reports whether the running edx binary was installed by
// Homebrew, by checking whether its real (symlink-resolved) path lives under a
// Cellar directory.
func homebrewInstall() (string, bool) {
	exe, err := os.Executable()
	if err != nil {
		return "", false
	}
	real, err := filepath.EvalSymlinks(exe)
	if err != nil {
		real = exe
	}
	sep := string(filepath.Separator)
	if strings.Contains(real, sep+"Cellar"+sep) {
		return real, true
	}
	return "", false
}

func updateViaHomebrew(ctx context.Context) error {
	brew, err := exec.LookPath("brew")
	if err != nil {
		notef("this is a Homebrew install; upgrade it with:")
		fmt.Fprintln(os.Stderr, "  brew update && brew upgrade edx")
		return nil
	}
	if !confirm("Run `brew upgrade edx` now?") {
		return fmt.Errorf("update cancelled")
	}
	for _, args := range [][]string{{"update"}, {"upgrade", "edx"}} {
		c := exec.CommandContext(ctx, brew, args...)
		c.Stdout, c.Stderr, c.Stdin = os.Stderr, os.Stderr, os.Stdin
		if err := c.Run(); err != nil {
			return fmt.Errorf("brew %s failed: %w", strings.Join(args, " "), err)
		}
	}
	fmt.Fprintf(os.Stderr, "%s edx upgraded via Homebrew\n", okMark())
	return nil
}

// --- self-replacement -------------------------------------------------------

// selfReplace downloads the release archive for the current platform, verifies
// it against the published checksums, extracts the edx binary, and atomically
// swaps it over the running executable.
func selfReplace(ctx context.Context, rel *ghRelease) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot locate current binary: %w", err)
	}
	if real, err := filepath.EvalSymlinks(exe); err == nil {
		exe = real
	}

	assetName := archiveName(rel.TagName)
	assetURL := findAsset(rel, assetName)
	if assetURL == "" {
		return fmt.Errorf("no release asset %q for %s/%s", assetName, runtime.GOOS, runtime.GOARCH)
	}
	sumURL := findAsset(rel, "checksums.txt")
	if sumURL == "" {
		return fmt.Errorf("release %s is missing checksums.txt", rel.TagName)
	}

	notef("downloading %s", assetName)
	archive, err := download(ctx, assetURL)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	sums, err := download(ctx, sumURL)
	if err != nil {
		return fmt.Errorf("checksum download failed: %w", err)
	}

	want := checksumFor(string(sums), assetName)
	if want == "" {
		return fmt.Errorf("no checksum for %s in checksums.txt", assetName)
	}
	got := fmt.Sprintf("%x", sha256.Sum256(archive))
	if got != want {
		return fmt.Errorf("checksum mismatch for %s: got %s, want %s", assetName, got, want)
	}

	bin, err := extractBinary(archive, "edx")
	if err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}

	// Stage in the same directory so the final rename is atomic on one
	// filesystem, then swap it over the running binary.
	dir := filepath.Dir(exe)
	tmp, err := os.CreateTemp(dir, ".edx-update-*")
	if err != nil {
		return fmt.Errorf("cannot write to %s (try `sudo edx update` or `brew upgrade edx`): %w", dir, err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err := tmp.Write(bin); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpName, 0o755); err != nil {
		return err
	}
	if err := os.Rename(tmpName, exe); err != nil {
		return fmt.Errorf("cannot replace %s (try `sudo edx update` or `brew upgrade edx`): %w", exe, err)
	}

	fmt.Fprintf(os.Stderr, "%s updated edx to %s — re-run your command to use it\n", okMark(), strings.TrimPrefix(rel.TagName, "v"))
	return nil
}

// archiveName builds the release asset name for the current platform, matching
// the GoReleaser name_template: edx_<version>_<Os>_<Arch>.tar.gz, with the OS
// title-cased and amd64 rendered as x86_64.
func archiveName(tag string) string {
	ver := strings.TrimPrefix(tag, "v")
	osName := map[string]string{"darwin": "Darwin", "linux": "Linux"}[runtime.GOOS]
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x86_64"
	}
	return fmt.Sprintf("edx_%s_%s_%s.tar.gz", ver, osName, arch)
}

func findAsset(rel *ghRelease, name string) string {
	for _, a := range rel.Assets {
		if a.Name == name {
			return a.URL
		}
	}
	return ""
}

func download(ctx context.Context, url string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "edx/"+Version)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	return io.ReadAll(io.LimitReader(resp.Body, maxDownloadBytes))
}

// checksumFor returns the hex SHA256 recorded for name in a checksums.txt body
// whose lines are "<hex>  <filename>".
func checksumFor(sums, name string) string {
	for _, line := range strings.Split(sums, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == name {
			return fields[0]
		}
	}
	return ""
}

// extractBinary returns the file named name from a gzipped tarball.
func extractBinary(targz []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(targz))
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil, fmt.Errorf("%s not found in archive", name)
		}
		if err != nil {
			return nil, err
		}
		if hdr.Typeflag == tar.TypeReg && filepath.Base(hdr.Name) == name {
			return io.ReadAll(io.LimitReader(tr, maxDownloadBytes))
		}
	}
}
