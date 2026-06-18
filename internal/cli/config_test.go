package cli

import (
	"strings"
	"testing"

	"github.com/edgedelta/edx/internal/config"
)

func TestFormatConfigPath(t *testing.T) {
	if got := formatConfigPath("/x/config.yaml", true); got != "/x/config.yaml (exists)" {
		t.Errorf("exists case: got %q", got)
	}
	if got := formatConfigPath("/x/config.yaml", false); got != "/x/config.yaml (not found)" {
		t.Errorf("missing case: got %q", got)
	}
}

func TestConfigViewFromResolveWithOverrides(t *testing.T) {
	t.Setenv("EDX_CONFIG", t.TempDir()+"/config.yaml") // no file -> empty config
	for _, e := range []string{
		config.EnvAPIToken, config.EnvProfile, config.EnvAPIURL,
		config.EnvChatURL, config.EnvAgentURL,
	} {
		t.Setenv(e, "")
	}
	t.Setenv(config.EnvEnv, config.EnvStaging)
	t.Setenv(config.EnvOrgID, "org-from-env")

	r, err := config.Resolve("", "", "", "")
	if err != nil {
		t.Fatalf("resolve must succeed without credentials: %v", err)
	}
	v := newConfigView(r)
	if v.Env != config.EnvStaging {
		t.Errorf("env override not reflected: %+v", v)
	}
	if v.OrgID != "org-from-env" {
		t.Errorf("org override not reflected: %+v", v)
	}
	if v.ChatURL != "https://chat.ai.staging.edgedelta.com" {
		t.Errorf("env should drive hosts: %+v", v)
	}
	if v.APIToken != "(not set)" {
		t.Errorf("unset token should be %q, got %q", "(not set)", v.APIToken)
	}
}

func TestNewConfigViewMasksSecrets(t *testing.T) {
	r := &config.Resolved{
		Profile:           "default",
		Env:               "prod",
		APIURL:            "https://api.edgedelta.com",
		ChatURL:           "https://chat.ai.edgedelta.com",
		AgentURL:          "https://agent.ai.edgedelta.com",
		OrgID:             "org-123",
		AuthMethod:        "oauth",
		APIToken:          "edx_secrettoken_abcdef1234",
		OAuthClientID:     "client-abcdef",
		OAuthAccessToken:  "eyJ0aGlzaXNhdG9rZW4Ab3k",
		OAuthRefreshToken: "refreshtok-abcdef",
		OAuthExpiry:       "2026-06-18T20:00:00Z",
	}
	v := newConfigView(r)

	// Non-secret fields copied verbatim.
	if v.Env != "prod" || v.OrgID != "org-123" || v.AgentURL != "https://agent.ai.edgedelta.com" {
		t.Errorf("non-secret fields not copied: %+v", v)
	}
	if v.OAuthExpiry != "2026-06-18T20:00:00Z" {
		t.Errorf("expiry must not be masked: %q", v.OAuthExpiry)
	}
	// API token must be masked and not leak the raw value.
	if v.APIToken == "edx_secrettoken_abcdef1234" {
		t.Errorf("api_token raw value leaked: got %q", v.APIToken)
	}
	if !strings.Contains(v.APIToken, "...") {
		t.Errorf("api_token should be masked, got %q", v.APIToken)
	}
	// Secret values must never appear in full.
	for _, raw := range []string{"eyJ0aGlzaXNhdG9rZW4Ab3k", "refreshtok-abcdef", "client-abcdef", "edx_secrettoken_abcdef1234"} {
		if strings.Contains(v.OAuthAccessToken+v.OAuthRefreshToken+v.OAuthClientID+v.APIToken, raw) {
			t.Errorf("raw secret %q leaked into view: %+v", raw, v)
		}
	}
	if !strings.Contains(v.OAuthAccessToken, "...") {
		t.Errorf("access token should be masked, got %q", v.OAuthAccessToken)
	}
}
