package oauth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
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

// TestBrowserAvailableSSH guards the headless-fallback signal: a remote shell
// must report no browser (so `auth login` falls back to the device flow),
// regardless of OS or a set DISPLAY, since the loopback redirect can't reach the
// remote machine.
func TestBrowserAvailableSSH(t *testing.T) {
	t.Setenv("SSH_CONNECTION", "1.2.3.4 51000 6.7.8.9 22")
	t.Setenv("DISPLAY", ":0")
	if BrowserAvailable() {
		t.Error("BrowserAvailable() = true under SSH_CONNECTION, want false")
	}
}

// TestDeviceLogin drives the full device flow against a fake authorization
// server: client registration, the device authorization request, and a token
// poll that returns authorization_pending once before approving. It checks the
// poll loop keeps going on pending, that PKCE binds the exchange (the verifier
// hashes to the challenge sent at the device request), and that the prompt gets
// the user_code + verification URI.
func TestDeviceLogin(t *testing.T) {
	var mu sync.Mutex
	var challenge string
	var tokenPolls int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth/client/register":
			var body struct {
				GrantTypes              []string `json:"grant_types"`
				TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
				RedirectURIs            []string `json:"redirect_uris"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode register body: %v", err)
			}
			if body.TokenEndpointAuthMethod != "none" {
				t.Errorf("token_endpoint_auth_method = %q, want none", body.TokenEndpointAuthMethod)
			}
			hasDeviceGrant := false
			for _, g := range body.GrantTypes {
				if g == deviceCodeGrantType {
					hasDeviceGrant = true
				}
			}
			if !hasDeviceGrant {
				t.Errorf("grant_types %v missing device_code grant", body.GrantTypes)
			}
			if len(body.RedirectURIs) == 0 {
				t.Errorf("redirect_uris must be non-empty")
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"client_id": "cid"})

		case "/oauth/device":
			if err := r.ParseForm(); err != nil {
				t.Fatal(err)
			}
			if r.Form.Get("client_id") != "cid" {
				t.Errorf("client_id = %q", r.Form.Get("client_id"))
			}
			if r.Form.Get("code_challenge_method") != "S256" {
				t.Errorf("code_challenge_method = %q", r.Form.Get("code_challenge_method"))
			}
			c := r.Form.Get("code_challenge")
			if c == "" {
				t.Errorf("code_challenge is empty")
			}
			mu.Lock()
			challenge = c
			mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"device_code":               "dcode",
				"user_code":                 "BCDF-GHJK",
				"verification_uri":          "https://app.example.com/device",
				"verification_uri_complete": "https://app.example.com/device?user_code=BCDF-GHJK",
				"expires_in":                900,
				"interval":                  1,
			})

		case "/oauth/token":
			if err := r.ParseForm(); err != nil {
				t.Fatal(err)
			}
			if r.Form.Get("grant_type") != deviceCodeGrantType {
				t.Errorf("grant_type = %q", r.Form.Get("grant_type"))
			}
			if r.Form.Get("device_code") != "dcode" {
				t.Errorf("device_code = %q", r.Form.Get("device_code"))
			}
			mu.Lock()
			wantChallenge := challenge
			tokenPolls++
			n := tokenPolls
			mu.Unlock()
			// PKCE: the verifier must hash (S256) to the challenge sent earlier.
			sum := sha256.Sum256([]byte(r.Form.Get("code_verifier")))
			if got := base64.RawURLEncoding.EncodeToString(sum[:]); got != wantChallenge {
				t.Errorf("code_verifier does not match challenge: %q vs %q", got, wantChallenge)
			}
			w.Header().Set("Content-Type", "application/json")
			if n == 1 {
				// First poll: user hasn't approved yet.
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]any{"error": "authorization_pending"})
				return
			}
			json.NewEncoder(w).Encode(map[string]any{
				"access_token":  "atok",
				"refresh_token": "rtok",
				"token_type":    "Bearer",
				"expires_in":    3600,
			})

		default:
			t.Errorf("unexpected path %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	var gotUserCode, gotURI string
	toks, err := DeviceLogin(context.Background(), srv.URL, DeviceLoginOptions{
		Timeout: 10 * time.Second,
		Prompt: func(userCode, verificationURI, _ string, _ bool) {
			gotUserCode = userCode
			gotURI = verificationURI
		},
	})
	if err != nil {
		t.Fatalf("DeviceLogin: %v", err)
	}
	if toks.AccessToken != "atok" || toks.RefreshToken != "rtok" {
		t.Errorf("tokens = %+v", toks)
	}
	if toks.ClientID != "cid" {
		t.Errorf("client_id = %q, want cid", toks.ClientID)
	}
	if toks.Expiry.Before(time.Now()) {
		t.Errorf("expiry not in the future: %v", toks.Expiry)
	}
	if gotUserCode != "BCDF-GHJK" {
		t.Errorf("prompt user_code = %q", gotUserCode)
	}
	if gotURI != "https://app.example.com/device" {
		t.Errorf("prompt verification_uri = %q", gotURI)
	}
	mu.Lock()
	polls := tokenPolls
	mu.Unlock()
	if polls != 2 {
		t.Errorf("token polls = %d, want 2 (pending then success)", polls)
	}
}
