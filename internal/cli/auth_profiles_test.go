package cli

import (
	"encoding/json"
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

func TestProfileListEntries_JSON(t *testing.T) {
	f := &config.File{
		DefaultProfile: "staging",
		Profiles: map[string]*config.Profile{
			"prod":    {OrgID: "2d6be233-f7bb-4fe1-90a5-28a95c86ec9c", AuthMethod: config.AuthMethodOAuth}, // Env empty → resolves to default
			"staging": {Env: config.EnvStaging, OrgID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", AuthMethod: config.AuthMethodToken},
		},
	}
	entries := profileListEntries(f)

	// Sorted by name: prod, then staging.
	if len(entries) != 2 || entries[0].Name != "prod" || entries[1].Name != "staging" {
		t.Fatalf("entries not sorted by name: %+v", entries)
	}
	// Empty env resolves to the effective default; org IDs are full (not shortened).
	if entries[0].Env != config.DefaultEnv {
		t.Errorf("empty env should resolve to default, got %q", entries[0].Env)
	}
	if entries[0].OrgID != "2d6be233-f7bb-4fe1-90a5-28a95c86ec9c" {
		t.Errorf("JSON should carry the full org id, got %q", entries[0].OrgID)
	}
	// Only the default profile is flagged.
	if entries[0].Default {
		t.Error("prod is not the default")
	}
	if !entries[1].Default {
		t.Error("staging is the default and should be flagged")
	}

	// Marshals to valid JSON with snake_case keys.
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"auth_method": "oauth"`) || !strings.Contains(string(data), `"default": true`) {
		t.Errorf("unexpected JSON:\n%s", data)
	}

	// Empty config marshals to an empty array, not null.
	empty, err := json.Marshal(profileListEntries(&config.File{}))
	if err != nil {
		t.Fatal(err)
	}
	if string(empty) != "[]" {
		t.Errorf("empty profiles should marshal to [], got %s", empty)
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
