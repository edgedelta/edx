package skills

import (
	"path/filepath"
	"testing"
)

func TestPlatformByName(t *testing.T) {
	if _, err := PlatformByName("claude"); err != nil {
		t.Errorf("claude should be known: %v", err)
	}
	if _, err := PlatformByName("emacs"); err == nil {
		t.Error("expected error for unknown platform")
	}
}

func TestSkillsRootGlobalAndProject(t *testing.T) {
	claude, _ := PlatformByName("claude")
	if got, want := claude.SkillsRoot("/home/x", false), filepath.FromSlash("/home/x/.claude/skills"); got != want {
		t.Errorf("global: got %q want %q", got, want)
	}
	if got, want := claude.SkillsRoot("/home/x", true), filepath.FromSlash(".claude/skills"); got != want {
		t.Errorf("project: got %q want %q", got, want)
	}

	oc, _ := PlatformByName("opencode")
	if got, want := oc.SkillsRoot("/home/x", false), filepath.FromSlash("/home/x/.config/opencode/skills"); got != want {
		t.Errorf("opencode global: got %q want %q", got, want)
	}
}

func TestInstalledDetectsByConfigDir(t *testing.T) {
	// Only ~/.claude and ~/.config/opencode exist on this fake machine.
	present := map[string]bool{
		filepath.FromSlash("/home/x/.claude"):          true,
		filepath.FromSlash("/home/x/.config/opencode"): true,
	}
	got := Installed("/home/x", func(p string) bool { return present[p] })

	names := map[string]bool{}
	for _, p := range got {
		names[p.Name] = true
	}
	if !names["claude"] || !names["opencode"] {
		t.Errorf("expected claude and opencode detected, got %v", names)
	}
	if names["cursor"] || names["codex"] {
		t.Errorf("detected an assistant whose dir does not exist: %v", names)
	}
}

func TestInstalledNoneWhenNothingPresent(t *testing.T) {
	if got := Installed("/home/x", func(string) bool { return false }); len(got) != 0 {
		t.Errorf("expected none installed, got %v", got)
	}
}

func TestDetect(t *testing.T) {
	env := map[string]string{"CURSOR_AGENT": "1"}
	getenv := func(k string) string { return env[k] }
	p, ok := Detect(getenv)
	if !ok || p.Name != "cursor" {
		t.Errorf("expected cursor, got %q ok=%v", p.Name, ok)
	}

	if _, ok := Detect(func(string) string { return "" }); ok {
		t.Error("expected no detection in empty env")
	}
}
