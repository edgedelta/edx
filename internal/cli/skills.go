package cli

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/edgedelta/edx/internal/skills"
)

var (
	flagSkillProject bool
	flagSkillName    string
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
		Use:   "update [platform|all]",
		Short: "Re-install skills, overwriting the previously installed copies",
		Args:  cobra.MaximumNArgs(1),
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
