package oauth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestPKCE(t *testing.T) {
	v, c, err := pkce()
	if err != nil {
		t.Fatal(err)
	}
	// RFC 7636: verifier must be 43-128 chars.
	if len(v) < 43 || len(v) > 128 {
		t.Errorf("verifier length %d out of range", len(v))
	}
	sum := sha256.Sum256([]byte(v))
	want := base64.RawURLEncoding.EncodeToString(sum[:])
	if c != want {
		t.Errorf("challenge != S256(verifier): %q vs %q", c, want)
	}
	// No padding allowed in PKCE.
	for _, s := range []string{v, c} {
		if s[len(s)-1] == '=' {
			t.Errorf("unexpected base64 padding in %q", s)
		}
	}
}

func TestRefresh(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth/token" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("grant_type") != "refresh_token" {
			t.Errorf("grant_type = %q", r.Form.Get("grant_type"))
		}
		if r.Form.Get("refresh_token") != "rtok" {
			t.Errorf("refresh_token = %q", r.Form.Get("refresh_token"))
		}
		if r.Form.Get("client_id") != "cid" {
			t.Errorf("client_id = %q", r.Form.Get("client_id"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "newaccess",
			"refresh_token": "newrefresh",
			"token_type":    "Bearer",
			"expires_in":    3600,
		})
	}))
	defer srv.Close()

	toks, err := Refresh(context.Background(), srv.URL, "cid", "rtok", srv.Client())
	if err != nil {
		t.Fatal(err)
	}
	if toks.AccessToken != "newaccess" || toks.RefreshToken != "newrefresh" {
		t.Errorf("unexpected tokens: %+v", toks)
	}
	if d := time.Until(toks.Expiry); d < 50*time.Minute || d > 70*time.Minute {
		t.Errorf("expiry ~1h expected, got %s", d)
	}
}

func TestOrgIDFromToken(t *testing.T) {
	// Build a JWT-shaped token: header.payload.sig (sig unused).
	enc := func(v any) string {
		b, _ := json.Marshal(v)
		return base64.RawURLEncoding.EncodeToString(b)
	}
	mk := func(payload any) string {
		return enc(map[string]string{"alg": "RS256"}) + "." + enc(payload) + ".sig"
	}

	// Edge Delta shape: attr.organization_id is a string array.
	tok := mk(map[string]any{"sub": "u", "attr": map[string][]string{"organization_id": {"org-abc"}}})
	if got := OrgIDFromToken(tok); got != "org-abc" {
		t.Errorf("attr.organization_id: got %q", got)
	}

	// Fallback: top-level string claim.
	if got := OrgIDFromToken(mk(map[string]any{"organization_id": "org-top"})); got != "org-top" {
		t.Errorf("top-level organization_id: got %q", got)
	}

	// Missing claim and malformed token return "".
	if got := OrgIDFromToken(mk(map[string]any{"sub": "u"})); got != "" {
		t.Errorf("missing claim should be empty, got %q", got)
	}
	if got := OrgIDFromToken("not-a-jwt"); got != "" {
		t.Errorf("malformed token should be empty, got %q", got)
	}
}

func TestRefreshServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid_grant"}`))
	}))
	defer srv.Close()
	if _, err := Refresh(context.Background(), srv.URL, "cid", "bad", srv.Client()); err == nil {
		t.Error("expected error on 400 from token endpoint")
	}
}

func TestGetJWTFromCookie(t *testing.T) {
	// A JWT whose payload carries an exp claim one hour out.
	exp := time.Now().Add(time.Hour).Unix()
	payload := base64.RawURLEncoding.EncodeToString([]byte(
		`{"exp":` + itoa(exp) + `}`))
	jwt := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256"}`)) + "." + payload + ".sig"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/cookie_service/get_jwt_from_cookie" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if got := r.Header.Get("Cookie"); got != "ed-admin-session=sess-abc" {
			t.Errorf("Cookie header = %q", got)
		}
		_, _ = w.Write([]byte(`{"bearer_token":"` + jwt + `"}`))
	}))
	defer srv.Close()

	tok, expiry, err := GetJWTFromCookie(context.Background(), srv.URL, "sess-abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	if tok != jwt {
		t.Errorf("token = %q, want %q", tok, jwt)
	}
	if d := time.Until(expiry); d < 55*time.Minute || d > 65*time.Minute {
		t.Errorf("expiry parsed wrong: %v away", d)
	}
}

func itoa(n int64) string {
	return strconv.FormatInt(n, 10)
}
