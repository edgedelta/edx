package cli

import (
	"path/filepath"
	"testing"
)

func TestSkillsOfferApplies(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{"edx logs search", true},
		{"edx auth login", true},
		{"edx pipelines list", true},
		{"edx", false},                // bare edx just prints help
		{"edx skills install", false}, // managing skills explicitly
		{"edx skills list", false},
		{"edx update", false}, // covers the `upgrade` alias too (same Name)
		{"edx version", false},
		{"edx help", false},
		{"edx completion bash", false},
		{"edx __complete logs", false},
	}
	for _, c := range cases {
		if got := skillsOfferApplies(c.path); got != c.want {
			t.Errorf("skillsOfferApplies(%q) = %v, want %v", c.path, got, c.want)
		}
	}
}

func TestSkillsOfferMarker(t *testing.T) {
	// StateDir derives from the config path, which EDX_CONFIG overrides.
	t.Setenv("EDX_CONFIG", filepath.Join(t.TempDir(), "config.yaml"))

	if skillsOfferAnswered() {
		t.Fatal("offer reported answered before any answer")
	}
	markSkillsOfferAnswered(false)
	if !skillsOfferAnswered() {
		t.Fatal("offer not reported answered after marking")
	}
	// Declining must be remembered exactly like accepting: never re-ask.
	markSkillsOfferAnswered(true)
	if !skillsOfferAnswered() {
		t.Fatal("marker lost after re-marking")
	}
}
