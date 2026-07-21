package cli

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"testing"
)

func TestParseSemver(t *testing.T) {
	cases := []struct {
		in                 string
		wantMa, wantMi, wp int
		wantOK             bool
	}{
		{"v1.2.3", 1, 2, 3, true},
		{"1.2.3", 1, 2, 3, true},
		{"v10.0.11", 10, 0, 11, true},
		{"v1.2.3-rc1", 1, 2, 3, true},          // pre-release suffix dropped
		{"v1.2.3-3-gabc123-dirty", 1, 2, 3, true}, // git-describe suffix dropped
		{"v1.2.3+meta", 1, 2, 3, true},
		{"dev", 0, 0, 0, false},
		{"", 0, 0, 0, false},
		{"1.2", 0, 0, 0, false},
		{"1.2.x", 0, 0, 0, false},
	}
	for _, c := range cases {
		ma, mi, p, ok := parseSemver(c.in)
		if ok != c.wantOK {
			t.Errorf("parseSemver(%q) ok=%v, want %v", c.in, ok, c.wantOK)
			continue
		}
		if ok && (ma != c.wantMa || mi != c.wantMi || p != c.wp) {
			t.Errorf("parseSemver(%q) = %d.%d.%d, want %d.%d.%d", c.in, ma, mi, p, c.wantMa, c.wantMi, c.wp)
		}
	}
}

func TestIsNewer(t *testing.T) {
	cases := []struct {
		latest, current string
		want            bool
	}{
		{"v1.2.4", "v1.2.3", true},
		{"v1.3.0", "v1.2.9", true},
		{"v2.0.0", "v1.9.9", true},
		{"v1.2.3", "v1.2.3", false}, // equal
		{"v1.2.2", "v1.2.3", false}, // older
		{"v1.2.4", "dev", false},    // dev never triggers an update
		{"garbage", "v1.2.3", false},
	}
	for _, c := range cases {
		if got := isNewer(c.latest, c.current); got != c.want {
			t.Errorf("isNewer(%q, %q) = %v, want %v", c.latest, c.current, got, c.want)
		}
	}
}

func TestArchiveNameMatchesGoReleaser(t *testing.T) {
	// The name must match the .goreleaser.yaml name_template exactly, or the
	// self-update download 404s.
	got := archiveName("v1.4.0")
	// Depends on the host GOOS/GOARCH, so just assert the invariant shape.
	if !bytes.HasPrefix([]byte(got), []byte("edx_1.4.0_")) {
		t.Errorf("archiveName leading segment wrong: %q", got)
	}
	if !bytes.HasSuffix([]byte(got), []byte(".tar.gz")) {
		t.Errorf("archiveName should end in .tar.gz: %q", got)
	}
	// amd64 must render as x86_64; arm64 stays arm64. Never the raw GOARCH "amd64".
	if bytes.Contains([]byte(got), []byte("amd64")) {
		t.Errorf("archiveName must render amd64 as x86_64: %q", got)
	}
}

func TestChecksumFor(t *testing.T) {
	body := "abc123  edx_1.4.0_Darwin_arm64.tar.gz\n" +
		"def456  edx_1.4.0_Linux_x86_64.tar.gz\n"
	if got := checksumFor(body, "edx_1.4.0_Darwin_arm64.tar.gz"); got != "abc123" {
		t.Errorf("checksumFor darwin = %q, want abc123", got)
	}
	if got := checksumFor(body, "edx_1.4.0_Linux_x86_64.tar.gz"); got != "def456" {
		t.Errorf("checksumFor linux = %q, want def456", got)
	}
	if got := checksumFor(body, "missing.tar.gz"); got != "" {
		t.Errorf("checksumFor(missing) = %q, want empty", got)
	}
}

func TestExtractBinary(t *testing.T) {
	want := []byte("#!/fake edx binary contents")
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	// A non-target file plus the target, to prove selection works.
	writeTar(t, tw, "README.md", []byte("readme"))
	writeTar(t, tw, "edx", want)
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}

	got, err := extractBinary(buf.Bytes(), "edx")
	if err != nil {
		t.Fatalf("extractBinary: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("extractBinary = %q, want %q", got, want)
	}

	if _, err := extractBinary(buf.Bytes(), "nope"); err == nil {
		t.Error("extractBinary should error when the named file is absent")
	}
}

func writeTar(t *testing.T, tw *tar.Writer, name string, data []byte) {
	t.Helper()
	if err := tw.WriteHeader(&tar.Header{
		Name:     name,
		Mode:     0o755,
		Size:     int64(len(data)),
		Typeflag: tar.TypeReg,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(data); err != nil {
		t.Fatal(err)
	}
}
