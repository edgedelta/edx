package cli

import (
	"strings"
	"testing"

	"github.com/edgedelta/edx/internal/config"
)

func TestNewConfigViewMasksSecrets(t *testing.T) {
	r := &config.Resolved{
		Profile:           "default",
		Env:               "prod",
		APIURL:            "https://api.edgedelta.com",
		ChatURL:           "https://chat.ai.edgedelta.com",
		AgentURL:          "https://agent.ai.edgedelta.com",
		OrgID:             "org-123",
		AuthMethod:        "oauth",
		APIToken:          "",
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
	// Unset token -> (not set).
	if v.APIToken != "(not set)" {
		t.Errorf("empty api_token should be %q, got %q", "(not set)", v.APIToken)
	}
	// Secret values must never appear in full.
	for _, raw := range []string{"eyJ0aGlzaXNhdG9rZW4Ab3k", "refreshtok-abcdef", "client-abcdef"} {
		if strings.Contains(v.OAuthAccessToken+v.OAuthRefreshToken+v.OAuthClientID, raw) {
			t.Errorf("raw secret %q leaked into view: %+v", raw, v)
		}
	}
	if !strings.Contains(v.OAuthAccessToken, "...") {
		t.Errorf("access token should be masked, got %q", v.OAuthAccessToken)
	}
}
