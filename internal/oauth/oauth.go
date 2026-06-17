// Package oauth implements the browser-based OAuth login the edx CLI uses as
// an alternative to a static API token. It performs an RFC 6749
// authorization-code flow with PKCE (RFC 7636) against the Edge Delta API:
//
//	POST /oauth/client/register   dynamically register a public CLI client
//	GET  /oauth/authorize         open in a browser; user logs in
//	POST /oauth/token             exchange the code (and later refresh) for a JWT
//
// The redirect target is a loopback HTTP server bound to 127.0.0.1 on an
// ephemeral port, the standard approach for native/CLI apps.
package oauth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// scope requests offline_access (so the server provisions a refresh token) and
// full_access (API access), plus the standard OIDC scopes.
const scope = "openid profile email offline_access full_access"

// Tokens is the result of a login or refresh.
type Tokens struct {
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
	ClientID     string
}

// LoginOptions tweaks the interactive flow.
type LoginOptions struct {
	// OpenBrowser opens the authorize URL automatically (default true).
	OpenBrowser bool
	// Prompt receives the authorize URL and a human message (default: stderr).
	Prompt func(url string)
	// HTTPClient overrides the client used for the register/token calls.
	HTTPClient *http.Client
	// Timeout bounds how long to wait for the browser callback (default 3m).
	Timeout time.Duration
}

// Login runs the full interactive flow against apiBase (e.g.
// https://api.staging.edgedelta.com) and returns the minted tokens.
func Login(ctx context.Context, apiBase string, opts LoginOptions) (Tokens, error) {
	apiBase = strings.TrimRight(apiBase, "/")
	hc := opts.HTTPClient
	if hc == nil {
		hc = &http.Client{Timeout: 30 * time.Second}
	}
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 3 * time.Minute
	}

	// Loopback server first, so we know the redirect URI before registering.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return Tokens{}, fmt.Errorf("start loopback listener: %w", err)
	}
	defer ln.Close()
	redirectURI := fmt.Sprintf("http://127.0.0.1:%d/callback", ln.Addr().(*net.TCPAddr).Port)

	clientID, err := registerClient(ctx, hc, apiBase, redirectURI)
	if err != nil {
		return Tokens{}, fmt.Errorf("register oauth client: %w", err)
	}

	verifier, challenge, err := pkce()
	if err != nil {
		return Tokens{}, err
	}
	state, err := randomURLSafe(24)
	if err != nil {
		return Tokens{}, err
	}

	authURL := authorizeURL(apiBase, clientID, redirectURI, state, challenge)

	resultCh := make(chan callbackResult, 1)
	srv := &http.Server{Handler: callbackHandler(state, resultCh)}
	go srv.Serve(ln)
	defer srv.Close()

	if opts.Prompt != nil {
		opts.Prompt(authURL)
	}
	if opts.OpenBrowser {
		_ = openBrowser(authURL)
	}

	select {
	case <-ctx.Done():
		return Tokens{}, ctx.Err()
	case <-time.After(timeout):
		return Tokens{}, fmt.Errorf("timed out after %s waiting for browser login", timeout)
	case res := <-resultCh:
		if res.err != nil {
			return Tokens{}, res.err
		}
		return exchangeCode(ctx, hc, apiBase, clientID, redirectURI, res.code, verifier)
	}
}

// Refresh exchanges a refresh token for a fresh access token.
func Refresh(ctx context.Context, apiBase, clientID, refreshToken string, hc *http.Client) (Tokens, error) {
	if hc == nil {
		hc = &http.Client{Timeout: 30 * time.Second}
	}
	form := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
		"client_id":     {clientID},
	}
	return postToken(ctx, hc, strings.TrimRight(apiBase, "/"), clientID, form)
}

// --- internals ---

type callbackResult struct {
	code string
	err  error
}

func callbackHandler(wantState string, ch chan<- callbackResult) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if e := q.Get("error"); e != "" {
			finish(w, false)
			ch <- callbackResult{err: fmt.Errorf("authorization failed: %s %s", e, q.Get("error_description"))}
			return
		}
		if q.Get("state") != wantState {
			finish(w, false)
			ch <- callbackResult{err: fmt.Errorf("state mismatch (possible CSRF); aborting")}
			return
		}
		code := q.Get("code")
		if code == "" {
			finish(w, false)
			ch <- callbackResult{err: fmt.Errorf("no authorization code in callback")}
			return
		}
		finish(w, true)
		ch <- callbackResult{code: code}
	})
	return mux
}

func finish(w http.ResponseWriter, ok bool) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	msg := "Login successful — you can close this tab and return to the terminal."
	if !ok {
		msg = "Login failed — return to the terminal for details."
	}
	fmt.Fprintf(w, "<!doctype html><meta charset=utf-8><title>edx</title>"+
		"<body style=\"font:16px -apple-system,sans-serif;display:grid;place-items:center;height:90vh;color:#1a1a1a\">"+
		"<div><h2>edx</h2><p>%s</p></div></body>", msg)
}

func registerClient(ctx context.Context, hc *http.Client, apiBase, redirectURI string) (string, error) {
	reqBody, _ := json.Marshal(map[string]any{
		"client_name":                "edx CLI",
		"redirect_uris":              []string{redirectURI},
		"grant_types":                []string{"authorization_code", "refresh_token"},
		"response_types":             []string{"code"},
		"token_endpoint_auth_method": "none", // public client (PKCE, no secret)
		"scope":                      scope,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiBase+"/oauth/client/register", strings.NewReader(string(reqBody)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("register returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var out struct {
		ClientID string `json:"client_id"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("decode register response: %w", err)
	}
	if out.ClientID == "" {
		return "", fmt.Errorf("register response had no client_id")
	}
	return out.ClientID, nil
}

func authorizeURL(apiBase, clientID, redirectURI, state, challenge string) string {
	q := url.Values{
		"response_type":         {"code"},
		"client_id":             {clientID},
		"redirect_uri":          {redirectURI},
		"scope":                 {scope},
		"state":                 {state},
		"code_challenge":        {challenge},
		"code_challenge_method": {"S256"},
	}
	return apiBase + "/oauth/authorize?" + q.Encode()
}

func exchangeCode(ctx context.Context, hc *http.Client, apiBase, clientID, redirectURI, code, verifier string) (Tokens, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"code_verifier": {verifier},
		"client_id":     {clientID},
		"redirect_uri":  {redirectURI},
	}
	return postToken(ctx, hc, apiBase, clientID, form)
}

func postToken(ctx context.Context, hc *http.Client, apiBase, clientID string, form url.Values) (Tokens, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiBase+"/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return Tokens{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := hc.Do(req)
	if err != nil {
		return Tokens{}, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Tokens{}, fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var out struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return Tokens{}, fmt.Errorf("decode token response: %w", err)
	}
	if out.AccessToken == "" {
		return Tokens{}, fmt.Errorf("token response had no access_token")
	}
	expiry := time.Now().Add(6 * time.Hour)
	if out.ExpiresIn > 0 {
		expiry = time.Now().Add(time.Duration(out.ExpiresIn) * time.Second)
	}
	return Tokens{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
		Expiry:       expiry,
		ClientID:     clientID,
	}, nil
}

// pkce returns a verifier and its S256 challenge (RFC 7636).
func pkce() (verifier, challenge string, err error) {
	verifier, err = randomURLSafe(64)
	if err != nil {
		return "", "", err
	}
	sum := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(sum[:])
	return verifier, challenge, nil
}

func randomURLSafe(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func openBrowser(u string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", u).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", u).Start()
	default:
		return exec.Command("xdg-open", u).Start()
	}
}
