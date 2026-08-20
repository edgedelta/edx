package cli

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/edgedelta/edx/internal/config"
	"github.com/edgedelta/edx/internal/skills"
)

var (
	flagSkillProject bool
	flagSkillName    string
)

const (
	// envNoSkillsOffer disables the one-time first-run skills offer.
	envNoSkillsOffer = "EDX_NO_SKILLS_OFFER"

	// skillsOfferFile marks that the first-run offer was answered, so it is
	// never repeated. Lives in the state dir next to the update-check cache.
	skillsOfferFile = "skills-offer.json"
)

func newSkillsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skills",
		Short: "Install Edge Delta agent skills into your coding assistant",
		Long: fmt.Sprintf(`Install the Edge Delta agent skills into a coding assistant.

The skills (ed-logs, ed-metrics, ed-monitors, ...) teach an AI agent how to
drive edx. They are embedded in this binary, so they always match this edx
version and install without network access.

PLATFORMS
  %s

  With no argument the platform is auto-detected: from the environment when
  edx is launched by the agent, otherwise from the assistants installed on
  this machine. Pass one explicitly, or "all" to install everywhere.

EXAMPLES
  edx skills list
  edx skills install                      # auto-detect the running assistant
  edx skills install claude
  edx skills install all
  edx skills install claude --project     # into ./.claude/skills instead of $HOME
  edx skills install claude --name ed-monitors
  edx skills show ed-logs`, strings.Join(skills.PlatformNames(), ", ")),
	}
	cmd.AddCommand(newSkillsListCmd(), newSkillsInstallCmd(), newSkillsUpdateCmd(), newSkillsShowCmd())
	return cmd
}

func newSkillsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List the available skills",
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := skills.List(skills.Embedded())
			if err != nil {
				return err
			}
			data, err := json.Marshal(list)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

func newSkillsShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <skill>",
		Short: "Print a skill's SKILL.md",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := skills.Read(skills.Embedded(), args[0])
			if err != nil {
				return err
			}
			_, err = os.Stdout.Write(b)
			return err
		},
	}
}

func newSkillsInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install [platform|all]",
		Short: "Install skills for a platform (auto-detected if omitted)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSkillsInstall(args, false)
		},
	}
	addSkillInstallFlags(cmd)
	return cmd
}

func newSkillsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update [platform|all]",
		Aliases: []string{"upgrade"},
		Short:   "Re-install skills, overwriting the previously installed copies",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSkillsInstall(args, true)
		},
	}
	addSkillInstallFlags(cmd)
	return cmd
}

func addSkillInstallFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&flagSkillProject, "project", false, "install into the current project instead of the user-global directory")
	cmd.Flags().StringVar(&flagSkillName, "name", "", "install only the named skill (default: all skills)")
}

// runSkillsInstall resolves the target platform(s) and copies the embedded
// skills out. update only changes the wording shown to the user.
func runSkillsInstall(args []string, update bool) error {
	fsys := skills.Embedded()

	plats, err := resolvePlatforms(args)
	if err != nil {
		return err
	}

	names, err := selectSkillNames(fsys)
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot determine home directory: %w", err)
	}

	// Summarize every destination, then confirm once.
	verb := "Install"
	if update {
		verb = "Update"
	}
	roots := make([]string, len(plats))
	for i, p := range plats {
		roots[i] = p.SkillsRoot(home, flagSkillProject)
		fmt.Fprintf(os.Stderr, "  %-9s -> %s\n", p.Name, roots[i])
	}
	if !confirm(fmt.Sprintf("%s %d skill(s) into %d location(s)?", verb, len(names), len(plats))) {
		return nil
	}

	for i, p := range plats {
		total := 0
		for _, n := range names {
			written, err := skills.Install(fsys, n, roots[i])
			if err != nil {
				return fmt.Errorf("install %s for %s: %w", n, p.Name, err)
			}
			total += written
		}
		fmt.Fprintf(os.Stdout, "%s: %d skill(s), %d file(s) -> %s\n", p.Name, len(names), total, roots[i])
	}
	return nil
}

// resolvePlatforms turns the optional positional arg into the platforms to
// install for: "all", a named platform, or environment auto-detection.
func resolvePlatforms(args []string) ([]skills.Platform, error) {
	if len(args) == 1 {
		if args[0] == "all" {
			return skills.Platforms, nil
		}
		p, err := skills.PlatformByName(args[0])
		if err != nil {
			return nil, err
		}
		return []skills.Platform{p}, nil
	}
	// When edx is launched by the agent itself, the environment names it.
	if p, ok := skills.Detect(os.Getenv); ok {
		fmt.Fprintf(os.Stderr, "detected %s from the environment\n", p.Name)
		return []skills.Platform{p}, nil
	}
	// Otherwise (a human running this from a normal terminal), fall back to the
	// assistants actually installed on this machine.
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot determine home directory: %w", err)
	}
	if installed := skills.Installed(home, dirExists); len(installed) > 0 {
		names := make([]string, len(installed))
		for i, p := range installed {
			names[i] = p.Name
		}
		fmt.Fprintf(os.Stderr, "found installed: %s\n", strings.Join(names, ", "))
		return installed, nil
	}
	return nil, fmt.Errorf("could not detect a coding assistant; specify one: edx skills install <%s|all>", strings.Join(skills.PlatformNames(), "|"))
}

// dirExists reports whether path is an existing directory.
func dirExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

// --- first-run offer ---------------------------------------------------------

// maybeOfferSkillsInstall offers, once, to install the embedded agent skills
// when a coding assistant is present but no skills are. Package managers can't
// do this at install time (Homebrew's post_install runs sandboxed with a fake
// HOME), so edx offers it the first time a human runs it — which also covers
// go install and curl installs. Safe for automation: it only fires when both
// stdin and stderr are TTYs, so agent, CI and piped callers never see it.
func maybeOfferSkillsInstall(cmd *cobra.Command) {
	if os.Getenv(envNoSkillsOffer) != "" {
		return
	}
	if !fileIsTTY(os.Stderr) || !fileIsTTY(os.Stdin) {
		return
	}
	if !skillsOfferApplies(cmd.CommandPath()) {
		return
	}
	if skillsOfferAnswered() {
		return
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	// Skills already installed somewhere: nothing to offer.
	if len(skillsRefreshTargets(home, dirExists)) > 0 {
		return
	}
	detected := skills.Installed(home, dirExists)
	if len(detected) == 0 {
		// No assistant on this machine yet. Stay silent and leave the offer
		// open: installing one later makes a future run ask.
		return
	}

	names := make([]string, len(detected))
	for i, p := range detected {
		names[i] = p.Name
	}
	fmt.Fprintf(os.Stderr, "edx ships agent skills that teach %s how to drive it.\n", strings.Join(names, ", "))
	accepted := confirmDefaultYes("Install them now? (later: edx skills install)")
	markSkillsOfferAnswered(accepted)
	if !accepted {
		return
	}
	installSkillsQuietly(home, detected)
}

// skillsOfferApplies reports whether the first-run offer may precede this
// command. Skills commands manage skills explicitly, and meta commands must
// stay noise-free (matching maybeNotifyUpdate's skip list).
func skillsOfferApplies(commandPath string) bool {
	fields := strings.Fields(commandPath)
	if len(fields) < 2 {
		return false // bare `edx` just prints help
	}
	switch fields[1] {
	case "skills", "update", "version", "help", "completion", "__complete", "__completeNoDesc":
		return false
	}
	return true
}

// installSkillsQuietly installs every embedded skill for the given platforms,
// printing one summary line per platform. Failures warn with the manual
// command instead of failing the user's actual command.
func installSkillsQuietly(home string, plats []skills.Platform) {
	fsys := skills.Embedded()
	names, err := skills.Names(fsys)
	if err != nil {
		return
	}
	for _, p := range plats {
		root := p.SkillsRoot(home, false)
		total := 0
		for _, n := range names {
			written, err := skills.Install(fsys, n, root)
			if err != nil {
				warnf("failed to install %s skills (%v) — run `edx skills install %s`", p.Name, err, p.Name)
				total = -1
				break
			}
			total += written
		}
		if total >= 0 {
			fmt.Fprintf(os.Stderr, "%s installed %d %s skill(s) (%d files) -> %s\n", okMark(), len(names), p.Name, total, root)
		}
	}
}

// skillsOfferAnswered reports whether the first-run offer was already made.
func skillsOfferAnswered() bool {
	dir, err := config.StateDir()
	if err != nil {
		return false
	}
	_, err = os.Stat(filepath.Join(dir, skillsOfferFile))
	return err == nil
}

// markSkillsOfferAnswered records the answer so the offer never repeats.
// Declining is remembered too: `edx skills install` stays the way back in.
func markSkillsOfferAnswered(accepted bool) {
	dir, err := config.StateDir()
	if err != nil {
		return
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return
	}
	data, err := json.Marshal(struct {
		AnsweredAt time.Time `json:"answered_at"`
		Accepted   bool      `json:"accepted"`
	}{time.Now().UTC(), accepted})
	if err != nil {
		return
	}
	_ = os.WriteFile(filepath.Join(dir, skillsOfferFile), data, 0o600)
}

// selectSkillNames returns the skills to install, honoring --name.
func selectSkillNames(fsys fs.FS) ([]string, error) {
	all, err := skills.Names(fsys)
	if err != nil {
		return nil, err
	}
	if flagSkillName == "" {
		return all, nil
	}
	for _, n := range all {
		if n == flagSkillName {
			return []string{flagSkillName}, nil
		}
	}
	return nil, fmt.Errorf("unknown skill %q (available: %s)", flagSkillName, strings.Join(all, ", "))
}
