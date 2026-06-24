package skills

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Platform is a coding assistant that reads SKILL.md skills from a known
// directory. All supported platforms use the same SKILL.md format, so a skill
// installs identically everywhere — only the destination directory differs.
type Platform struct {
	Name string // canonical name used as the install argument

	// globalRel / projectRel are the skills directories relative to the user
	// home directory and to the current project root, respectively.
	globalRel  string
	projectRel string

	// detectEnv are environment variables whose presence indicates this
	// assistant is the one driving edx. Detection is best-effort; an explicit
	// platform argument is always authoritative.
	detectEnv []string
}

// Platforms are the assistants edx can install skills into.
var Platforms = []Platform{
	{Name: "claude", globalRel: ".claude/skills", projectRel: ".claude/skills", detectEnv: []string{"CLAUDECODE", "CLAUDE_CODE"}},
	{Name: "cursor", globalRel: ".cursor/skills", projectRel: ".cursor/skills", detectEnv: []string{"CURSOR_AGENT", "CURSOR_TRACE_ID"}},
	{Name: "codex", globalRel: ".codex/skills", projectRel: ".codex/skills", detectEnv: []string{"CODEX_SANDBOX", "CODEX_HOME"}},
	{Name: "opencode", globalRel: ".config/opencode/skills", projectRel: ".opencode/skills", detectEnv: []string{"OPENCODE"}},
}

// PlatformNames returns the supported platform names, for help and errors.
func PlatformNames() []string {
	names := make([]string, len(Platforms))
	for i, p := range Platforms {
		names[i] = p.Name
	}
	return names
}

// PlatformByName looks up a platform by its canonical name.
func PlatformByName(name string) (Platform, error) {
	for _, p := range Platforms {
		if p.Name == name {
			return p, nil
		}
	}
	return Platform{}, fmt.Errorf("unknown platform %q (supported: %s)", name, strings.Join(PlatformNames(), ", "))
}

// Detect returns the platform indicated by the environment, using getenv to
// read variables (os.Getenv in production, a fake in tests). ok is false when
// no platform can be determined.
func Detect(getenv func(string) string) (Platform, bool) {
	for _, p := range Platforms {
		for _, key := range p.detectEnv {
			if getenv(key) != "" {
				return p, true
			}
		}
	}
	return Platform{}, false
}

// SkillsRoot returns the directory skills install into. When project is true it
// is relative to the current directory; otherwise it is under home.
func (p Platform) SkillsRoot(home string, project bool) string {
	if project {
		return filepath.FromSlash(p.projectRel)
	}
	return filepath.Join(home, filepath.FromSlash(p.globalRel))
}
