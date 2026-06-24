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

// Installed returns the platforms whose config directory exists under home,
// i.e. the assistants actually set up on this machine. dirExists reports
// whether a directory is present (os-backed in production, faked in tests).
// This is the fallback when environment detection fails — the common case when
// a human runs `edx skills install` from a normal terminal rather than from
// inside the agent.
func Installed(home string, dirExists func(string) bool) []Platform {
	var out []Platform
	for _, p := range Platforms {
		if dirExists(filepath.Join(home, p.markerDir())) {
			out = append(out, p)
		}
	}
	return out
}

// markerDir is the assistant's config directory (the parent of its skills
// directory), used to tell whether the assistant is installed.
func (p Platform) markerDir() string {
	return filepath.Dir(filepath.FromSlash(p.globalRel))
}

// SkillsRoot returns the directory skills install into. When project is true it
// is relative to the current directory; otherwise it is under home.
func (p Platform) SkillsRoot(home string, project bool) string {
	if project {
		return filepath.FromSlash(p.projectRel)
	}
	return filepath.Join(home, filepath.FromSlash(p.globalRel))
}
