package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePrecedence(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	t.Setenv("EDX_CONFIG", cfgPath)
	t.Setenv(EnvAPIToken, "")
	t.Setenv(EnvOrgID, "")
	t.Setenv(EnvAPIURL, "")
	t.Setenv(EnvProfile, "")

	content := `default_profile: default
profiles:
  default:
    org_id: org-from-file
    api_token: token-from-file
`
	if err := os.WriteFile(cfgPath, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	r, err := Resolve("", "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if r.OrgID != "org-from-file" || r.APIToken != "token-from-file" || r.APIURL != DefaultAPIURL {
		t.Errorf("unexpected resolution from file: %+v", r)
	}

	t.Setenv(EnvOrgID, "org-from-env")
	r, err = Resolve("", "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if r.OrgID != "org-from-env" {
		t.Errorf("env should override file, got %+v", r)
	}

	r, err = Resolve("", "", "org-from-flag", "token-from-flag")
	if err != nil {
		t.Fatal(err)
	}
	if r.OrgID != "org-from-flag" || r.APIToken != "token-from-flag" {
		t.Errorf("flags should override env, got %+v", r)
	}
}

func TestResolveMissingProfileFails(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("EDX_CONFIG", filepath.Join(dir, "config.yaml"))
	if _, err := Resolve("nope", "", "", ""); err == nil {
		t.Error("expected error for explicitly requested missing profile")
	}
}

func TestSaveAndLoadRoundtrip(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	t.Setenv("EDX_CONFIG", cfgPath)

	f := &File{
		DefaultProfile: "default",
		Profiles: map[string]*Profile{
			"default": {OrgID: "o", APIToken: "t", APIURL: "https://example.com"},
		},
	}
	if err := f.Save(); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("config file should be 0600, got %v", info.Mode().Perm())
	}
	loaded, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	p := loaded.Profiles["default"]
	if p == nil || p.OrgID != "o" || p.APIToken != "t" {
		t.Errorf("roundtrip mismatch: %+v", loaded)
	}
}
