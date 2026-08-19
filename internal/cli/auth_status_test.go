package cli

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/edgedelta/edx/internal/api"
	"github.com/edgedelta/edx/internal/config"
)

func TestProfileStatus_TokenHasNoLocalExpiry(t *testing.T) {
	now := time.Date(2026, 8, 13, 12, 0, 0, 0, time.UTC)
	status, expiresAt := profileStatus(&config.Profile{AuthMethod: config.AuthMethodToken}, now)
	if status != "-" || expiresAt != "" {
		t.Errorf("token profile: got status %q expiresAt %q, want %q and empty", status, expiresAt, "-")
	}
	// Empty auth method defaults to token.
	status, _ = profileStatus(&config.Profile{}, now)
	if status != "-" {
		t.Errorf("empty auth method should behave like token, got %q", status)
	}
}

func TestProfileStatus_OAuthWithRefreshTokenIsOK(t *testing.T) {
	now := time.Date(2026, 8, 13, 12, 0, 0, 0, time.UTC)
	p := &config.Profile{
		AuthMethod:        config.AuthMethodOAuth,
		OAuthRefreshToken: "rt",
		OAuthExpiry:       now.Add(-1 * time.Hour).Format(time.RFC3339), // stale access token is fine: it refreshes
	}
	status, expiresAt := profileStatus(p, now)
	if status != "ok (auto-refresh)" {
		t.Errorf("oauth with refresh token: got %q, want %q", status, "ok (auto-refresh)")
	}
	if expiresAt != p.OAuthExpiry {
		t.Errorf("expiresAt should carry the access-token expiry, got %q", expiresAt)
	}
}

func TestProfileStatus_OAuthWithoutRefreshToken(t *testing.T) {
	now := time.Date(2026, 8, 13, 12, 0, 0, 0, time.UTC)

	expired := &config.Profile{
		AuthMethod:  config.AuthMethodOAuth,
		OAuthExpiry: now.Add(-3 * time.Hour).Format(time.RFC3339),
	}
	status, _ := profileStatus(expired, now)
	if status != "expired 3h ago" {
		t.Errorf("expired oauth: got %q, want %q", status, "expired 3h ago")
	}

	valid := &config.Profile{
		AuthMethod:  config.AuthMethodOAuth,
		OAuthExpiry: now.Add(2 * time.Hour).Format(time.RFC3339),
	}
	status, _ = profileStatus(valid, now)
	if status != "expires in 2h" {
		t.Errorf("valid oauth without refresh: got %q, want %q", status, "expires in 2h")
	}

	noExpiry := &config.Profile{AuthMethod: config.AuthMethodOAuth}
	status, expiresAt := profileStatus(noExpiry, now)
	if status != "unknown" || expiresAt != "" {
		t.Errorf("oauth with no expiry data: got %q/%q, want unknown and empty", status, expiresAt)
	}
}

func TestProfileStatus_Cookie(t *testing.T) {
	now := time.Date(2026, 8, 13, 12, 0, 0, 0, time.UTC)

	// Cookie saved before expiry tracking existed: no way to know.
	legacy := &config.Profile{AuthMethod: config.AuthMethodCookie, SessionCookie: "c"}
	status, expiresAt := profileStatus(legacy, now)
	if status != "unknown (re-login to track)" || expiresAt != "" {
		t.Errorf("legacy cookie: got %q/%q, want unknown (re-login to track) and empty", status, expiresAt)
	}

	fresh := &config.Profile{
		AuthMethod:    config.AuthMethodCookie,
		SessionCookie: "c",
		CookieExpiry:  now.Add(5 * time.Hour).Format(time.RFC3339),
	}
	status, expiresAt = profileStatus(fresh, now)
	if status != "expires in 5h" {
		t.Errorf("fresh cookie: got %q, want %q", status, "expires in 5h")
	}
	if expiresAt != fresh.CookieExpiry {
		t.Errorf("expiresAt should carry the cookie expiry, got %q", expiresAt)
	}

	dead := &config.Profile{
		AuthMethod:    config.AuthMethodCookie,
		SessionCookie: "c",
		CookieExpiry:  now.Add(-3 * time.Hour).Format(time.RFC3339),
	}
	status, _ = profileStatus(dead, now)
	if status != "expired 3h ago" {
		t.Errorf("dead cookie: got %q, want %q", status, "expired 3h ago")
	}
}

func TestHumanDur(t *testing.T) {
	cases := []struct {
		d    time.Duration
		want string
	}{
		{30 * time.Second, "<1m"},
		{42 * time.Minute, "42m"},
		{90 * time.Minute, "1h30m"},
		{5 * time.Hour, "5h"},
		{72 * time.Hour, "3d"},
	}
	for _, c := range cases {
		if got := humanDur(c.d); got != c.want {
			t.Errorf("humanDur(%v) = %q, want %q", c.d, got, c.want)
		}
	}
}

func TestNewCookieProfile_RecordsExpiry(t *testing.T) {
	now := time.Date(2026, 8, 13, 12, 0, 0, 0, time.UTC)
	p := newCookieProfile("prod", "org-1", "cookie-value", now)
	if p.AuthMethod != config.AuthMethodCookie || p.SessionCookie != "cookie-value" || p.OrgID != "org-1" || p.Env != "prod" {
		t.Fatalf("unexpected profile: %+v", p)
	}
	want := now.Add(cookieLifetime).UTC().Format(time.RFC3339)
	if p.CookieExpiry != want {
		t.Errorf("CookieExpiry = %q, want %q", p.CookieExpiry, want)
	}
}

func TestFormatProfileList_StatusColumn(t *testing.T) {
	now := time.Now()
	f := &config.File{
		Profiles: map[string]*config.Profile{
			"oauthy": {Env: config.EnvProd, OrgID: "2d6be233-f7bb-4fe1-90a5-28a95c86ec9c", AuthMethod: config.AuthMethodOAuth, OAuthRefreshToken: "rt"},
			"cookied": {
				Env: config.EnvProd, OrgID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
				AuthMethod: config.AuthMethodCookie, SessionCookie: "c",
				CookieExpiry: now.Add(-2 * time.Hour).UTC().Format(time.RFC3339),
			},
		},
	}
	out := formatProfileList(f)
	if !strings.Contains(out, "STATUS") {
		t.Fatalf("table should have a STATUS header:\n%s", out)
	}
	if !strings.Contains(out, "ok (auto-refresh)") {
		t.Errorf("oauth row should show ok (auto-refresh):\n%s", out)
	}
	if !strings.Contains(out, "expired 2h ago") {
		t.Errorf("cookie row should show expired 2h ago:\n%s", out)
	}
}

func TestProfileListEntries_StatusJSON(t *testing.T) {
	now := time.Now()
	expiry := now.Add(5 * time.Hour).UTC().Format(time.RFC3339)
	f := &config.File{
		Profiles: map[string]*config.Profile{
			"c": {Env: config.EnvProd, OrgID: "o", AuthMethod: config.AuthMethodCookie, SessionCookie: "s", CookieExpiry: expiry},
		},
	}
	data, err := json.Marshal(profileListEntries(f))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"status":"expires in 5h"`) || !strings.Contains(string(data), `"expires_at":"`+expiry+`"`) {
		t.Errorf("JSON should carry status and expires_at:\n%s", data)
	}
}

func TestVerifyAuth(t *testing.T) {
	ok := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[]`))
	}))
	defer ok.Close()
	c := api.New(ok.URL, ok.URL, ok.URL, ok.URL, "org", &api.Auth{APIToken: "t"}, 2*time.Second)
	if err := verifyAuth(context.Background(), c); err != nil {
		t.Errorf("200 response should verify, got %v", err)
	}

	unauth := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
	defer unauth.Close()
	c = api.New(unauth.URL, unauth.URL, unauth.URL, unauth.URL, "org", &api.Auth{APIToken: "t"}, 2*time.Second)
	if err := verifyAuth(context.Background(), c); err == nil {
		t.Error("401 response should fail verification")
	}
}

func TestCheckLabel(t *testing.T) {
	if got := checkLabel(nil); !strings.Contains(got, "ok") {
		t.Errorf("nil error should label ok, got %q", got)
	}
	got := checkLabel(context.DeadlineExceeded)
	if !strings.Contains(got, "context deadline exceeded") {
		t.Errorf("error label should carry the cause, got %q", got)
	}
}

func TestRunProfileChecks_FillsEveryEntry(t *testing.T) {
	entries := []profileListEntry{{Name: "a"}, {Name: "b"}, {Name: "c"}}
	runProfileChecks(context.Background(), entries, func(_ context.Context, name string) error {
		if name == "b" {
			return errors.New("boom")
		}
		return nil
	})
	if entries[0].Check != "ok" || entries[2].Check != "ok" {
		t.Errorf("healthy profiles should check ok: %+v", entries)
	}
	if entries[1].Check != "failed: boom" {
		t.Errorf("failing profile should carry the cause, got %q", entries[1].Check)
	}
}

func TestAuthList_HasNoCheckFlag(t *testing.T) {
	if newAuthListCmd().Flags().Lookup("check") != nil {
		t.Error("live checks belong to `auth status`; list must stay offline-only")
	}
}

func TestAuthStatus_AcceptsAtMostOneProfileArg(t *testing.T) {
	cmd := newAuthStatusCmd()
	if cmd.Args == nil {
		t.Fatal("status should declare an Args validator (MaximumNArgs(1))")
	}
	if err := cmd.Args(cmd, nil); err != nil {
		t.Errorf("no args should be accepted: %v", err)
	}
	if err := cmd.Args(cmd, []string{"demo"}); err != nil {
		t.Errorf("one profile arg should be accepted: %v", err)
	}
	if err := cmd.Args(cmd, []string{"a", "b"}); err == nil {
		t.Error("two args should be rejected")
	}
}

func TestPlainCheck_CompactsAuthRejections(t *testing.T) {
	err := errors.New("API error 401 on GET https://api.edgedelta.com/v1/orgs/e2c56e6c-9463-48de-9569-2/facet_keys: unauthorized")
	if got := plainCheck(err); got != "failed: rejected (401 unauthorized — credentials expired or revoked)" {
		t.Errorf("401 should compact to a rejection notice, got %q", got)
	}
	err = errors.New("API error 403 on GET https://api.edgedelta.com/v1/orgs/x/facet_keys: forbidden")
	if got := plainCheck(err); got != "failed: rejected (403 forbidden)" {
		t.Errorf("403 should compact to a rejection notice, got %q", got)
	}
	// Other errors pass through (first line, truncated).
	if got := plainCheck(errors.New("dial tcp: connection refused")); got != "failed: dial tcp: connection refused" {
		t.Errorf("non-auth error should pass through, got %q", got)
	}
}
