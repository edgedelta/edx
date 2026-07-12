package cli

import (
	"strings"
	"testing"

	"github.com/edgedelta/edx/internal/config"
)

func TestFormatProfileList_Empty(t *testing.T) {
	out := formatProfileList(&config.File{})
	if !strings.Contains(out, "No profiles") {
		t.Errorf("empty list should mention no profiles, got %q", out)
	}
}

func TestFormatProfileList_MarksDefaultAndColumns(t *testing.T) {
	f := &config.File{
		DefaultProfile: "staging",
		Profiles: map[string]*config.Profile{
			"prod":    {Env: config.EnvProd, OrgID: "2d6be233-f7bb-4fe1-90a5-28a95c86ec9c", AuthMethod: config.AuthMethodOAuth},
			"staging": {Env: config.EnvStaging, OrgID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", AuthMethod: config.AuthMethodToken},
		},
	}
	out := formatProfileList(f)

	// Profiles are listed alphabetically, so find each line.
	var prodLine, stagingLine string
	for _, ln := range strings.Split(out, "\n") {
		if strings.Contains(ln, "prod") && !strings.Contains(ln, "NAME") {
			prodLine = ln
		}
		if strings.Contains(ln, "staging") {
			stagingLine = ln
		}
	}
	if strings.HasPrefix(strings.TrimSpace(prodLine), "*") {
		t.Errorf("prod is not default, should not be marked: %q", prodLine)
	}
	if !strings.HasPrefix(strings.TrimSpace(stagingLine), "*") {
		t.Errorf("staging is default, should be marked with *: %q", stagingLine)
	}
	if !strings.Contains(stagingLine, "token") || !strings.Contains(prodLine, "oauth") {
		t.Errorf("auth method missing from rows:\n%s", out)
	}
	// Org is shortened, not printed in full.
	if strings.Contains(out, "2d6be233-f7bb-4fe1-90a5-28a95c86ec9c") {
		t.Errorf("org id should be shortened, got full id:\n%s", out)
	}
}

func TestLoginWouldClobber(t *testing.T) {
	cfg := &config.File{
		Profiles: map[string]*config.Profile{"default": {Env: config.EnvProd, OrgID: "o"}},
	}

	// Implicit "default" that already exists → clobber (block).
	if !loginWouldClobber(cfg, "default", false /*explicit*/, false /*force*/) {
		t.Error("implicit re-login onto existing default should be blocked")
	}
	// --force overrides.
	if loginWouldClobber(cfg, "default", false, true) {
		t.Error("--force should allow overwrite")
	}
	// Explicit --profile means the user chose the name knowingly → allow.
	if loginWouldClobber(cfg, "default", true, false) {
		t.Error("explicit --profile should allow overwrite")
	}
	// New name that does not exist → never a clobber.
	if loginWouldClobber(cfg, "brand-new", false, false) {
		t.Error("saving a new profile is not a clobber")
	}
}
