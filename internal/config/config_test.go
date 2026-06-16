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
	t.Setenv(EnvEnv, "")
	t.Setenv(EnvProfile, "")
	t.Setenv(EnvAPIURL, "")
	t.Setenv(EnvChatURL, "")
	t.Setenv(EnvAgentURL, "")

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
	if r.OrgID != "org-from-file" || r.APIToken != "token-from-file" {
		t.Errorf("unexpected resolution from file: %+v", r)
	}
	if r.Env != EnvProd || r.APIURL != "https://api.edgedelta.com" {
		t.Errorf("expected default env prod, got %+v", r)
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

func TestResolveEnvironmentEndpoints(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	t.Setenv("EDX_CONFIG", cfgPath)
	t.Setenv(EnvAPIToken, "")
	t.Setenv(EnvOrgID, "")
	t.Setenv(EnvEnv, "")
	t.Setenv(EnvProfile, "")
	t.Setenv(EnvAPIURL, "")
	t.Setenv(EnvChatURL, "")
	t.Setenv(EnvAgentURL, "")

	content := `default_profile: prod
profiles:
  prod:
    org_id: o
    api_token: t
  staging:
    env: staging
    org_id: o
    api_token: t
`
	if err := os.WriteFile(cfgPath, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	// Default (prod) profile resolves to prod hosts for every service.
	r, err := Resolve("prod", "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if r.ChatURL != "https://chat.ai.edgedelta.com" || r.AgentURL != "https://agent.ai.edgedelta.com" {
		t.Errorf("prod AI hosts wrong: %+v", r)
	}

	// staging profile (env: staging) moves all hosts to staging.
	r, err = Resolve("staging", "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if r.Env != EnvStaging || r.APIURL != "https://api.staging.edgedelta.com" ||
		r.ChatURL != "https://chat.ai.staging.edgedelta.com" || r.AgentURL != "https://agent.ai.staging.edgedelta.com" {
		t.Errorf("staging endpoints wrong: %+v", r)
	}

	// --env flag overrides the profile's environment.
	r, err = Resolve("prod", "staging", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if r.ChatURL != "https://chat.ai.staging.edgedelta.com" {
		t.Errorf("--env should override profile env, got %+v", r)
	}

	// Unknown env is an error.
	if _, err := Resolve("prod", "bogus", "", ""); err == nil {
		t.Error("expected error for unknown environment")
	}
}

func TestResolveHostOverrides(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	t.Setenv("EDX_CONFIG", cfgPath)
	t.Setenv(EnvAPIToken, "")
	t.Setenv(EnvOrgID, "")
	t.Setenv(EnvEnv, "")
	t.Setenv(EnvProfile, "")
	t.Setenv(EnvAPIURL, "")
	t.Setenv(EnvChatURL, "")
	t.Setenv(EnvAgentURL, "")

	// Override only the chat host; the env (staging) still drives the rest.
	t.Setenv(EnvChatURL, "http://localhost:9999")
	r, err := Resolve("", "staging", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if r.ChatURL != "http://localhost:9999" {
		t.Errorf("ED_CHAT_URL should override chat host, got %q", r.ChatURL)
	}
	if r.APIURL != "https://api.staging.edgedelta.com" || r.AgentURL != "https://agent.ai.staging.edgedelta.com" {
		t.Errorf("override should not disturb other hosts: %+v", r)
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
			"default": {Env: EnvStaging, OrgID: "o", APIToken: "t"},
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
	if p == nil || p.OrgID != "o" || p.APIToken != "t" || p.Env != EnvStaging {
		t.Errorf("roundtrip mismatch: %+v", loaded)
	}
}
