package oauth

// This file implements the client side of the OAuth 2.0 Device Authorization
// Grant (RFC 8628), used by `edx auth login --device` to sign in on machines
// without a usable browser (SSH sessions, CI, headless boxes) where the loopback
// redirect of the standard authorization-code flow (oauth.go) cannot work.
//
// Unlike the loopback flow, nothing is delivered back to this process
// out-of-band: we display a short user_code, the human confirms it in a browser
// anywhere, and we poll the token endpoint until that approval completes. The
// user_code confirmation is what binds the approving human to this CLI instance.

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"
)

// BrowserAvailable reports whether a local browser can likely be opened for the
// loopback login flow. It returns false for remote shells (SSH) and headless
// Linux/BSD (no display server) — where the loopback redirect can't reach this
// machine and the CLI should use the device flow instead. macOS and Windows are
// assumed to have a usable local browser.
func BrowserAvailable() bool {
	// A remote shell: even if an opener exists, the browser it launches (and the
	// loopback redirect) would be on the wrong machine.
	if os.Getenv("SSH_CONNECTION") != "" || os.Getenv("SSH_TTY") != "" {
		return false
	}
	switch runtime.GOOS {
	case "darwin", "windows":
		return true
	default:
		return os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != ""
	}
}

// deviceCodeGrantType is the RFC 8628 grant_type for polling the token endpoint.
const deviceCodeGrantType = "urn:ietf:params:oauth:grant-type:device_code"

// DeviceLoginOptions tweaks the device flow.
type DeviceLoginOptions struct {
	// OpenBrowser opens the verification URI automatically (default true).
	OpenBrowser bool
	// Prompt receives the details to show the user: the user_code to enter, the
	// verification URI, the pre-filled verification URI, and whether a browser
	// was opened automatically.
	Prompt func(userCode, verificationURI, verificationURIComplete string, opened bool)
	// HTTPClient overrides the client used for the register/device/token calls.
	HTTPClient *http.Client
	// Timeout bounds how long to wait for approval. When zero, the server's
	// expires_in is used (falling back to 15m).
	Timeout time.Duration
}

// deviceAuthResponse is the RFC 8628 §3.2 device authorization response.
type deviceAuthResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int64  `json:"expires_in"`
	Interval                int64  `json:"interval"`
}

// DeviceLogin runs the full device authorization flow against apiBase and
// returns the minted tokens once the user approves in a browser.
func DeviceLogin(ctx context.Context, apiBase string, opts DeviceLoginOptions) (Tokens, error) {
	apiBase = strings.TrimRight(apiBase, "/")
	hc := opts.HTTPClient
	if hc == nil {
		hc = &http.Client{Timeout: 30 * time.Second}
	}

	clientID, err := registerDeviceClient(ctx, hc, apiBase)
	if err != nil {
		return Tokens{}, fmt.Errorf("register oauth client: %w", err)
	}

	verifier, challenge, err := pkce()
	if err != nil {
		return Tokens{}, err
	}

	dar, err := requestDeviceAuthorization(ctx, hc, apiBase, clientID, challenge)
	if err != nil {
		return Tokens{}, err
	}

	opened := false
	if opts.OpenBrowser {
		target := dar.VerificationURIComplete
		if target == "" {
			target = dar.VerificationURI
		}
		opened = openBrowser(target) == nil
	}
	if opts.Prompt != nil {
		opts.Prompt(dar.UserCode, dar.VerificationURI, dar.VerificationURIComplete, opened)
	}

	interval := time.Duration(dar.Interval) * time.Second
	if interval <= 0 {
		interval = 5 * time.Second
	}

	timeout := opts.Timeout
	if timeout <= 0 {
		if dar.ExpiresIn > 0 {
			timeout = time.Duration(dar.ExpiresIn) * time.Second
		} else {
			timeout = 15 * time.Minute
		}
	}
	deadline := time.Now().Add(timeout)

	for {
		select {
		case <-ctx.Done():
			return Tokens{}, ctx.Err()
		case <-time.After(interval):
		}

		if time.Now().After(deadline) {
			return Tokens{}, fmt.Errorf("timed out waiting for you to approve the device login")
		}

		toks, errCode, err := pollDeviceToken(ctx, hc, apiBase, clientID, dar.DeviceCode, verifier)
		if err != nil {
			return Tokens{}, err
		}
		switch errCode {
		case "":
			return toks, nil
		case "authorization_pending":
			// Not approved yet; keep polling.
		case "slow_down":
			// RFC 8628 §3.5: back off by 5s and keep polling.
			interval += 5 * time.Second
		case "access_denied":
			return Tokens{}, fmt.Errorf("the login was denied in the browser")
		case "expired_token":
			return Tokens{}, fmt.Errorf("the device code expired before you approved; run the command again")
		default:
			return Tokens{}, fmt.Errorf("device authorization failed: %s", errCode)
		}
	}
}

func registerDeviceClient(ctx context.Context, hc *http.Client, apiBase string) (string, error) {
	reqBody, _ := json.Marshal(map[string]any{
		"client_name": "edx CLI",
		// Device clients never use a redirect, but registration requires at least
		// one redirect URI; a loopback placeholder is stored and never used.
		"redirect_uris":              []string{"http://127.0.0.1/callback"},
		"grant_types":                []string{deviceCodeGrantType, "refresh_token"},
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

func requestDeviceAuthorization(ctx context.Context, hc *http.Client, apiBase, clientID, challenge string) (*deviceAuthResponse, error) {
	form := url.Values{
		"client_id":             {clientID},
		"scope":                 {scope},
		"code_challenge":        {challenge},
		"code_challenge_method": {"S256"},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiBase+"/oauth/device", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("device authorization returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var out deviceAuthResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("decode device authorization response: %w", err)
	}
	if out.DeviceCode == "" || out.UserCode == "" {
		return nil, fmt.Errorf("device authorization response was missing device_code/user_code")
	}
	return &out, nil
}

// pollDeviceToken makes one token-endpoint poll. On success it returns the
// tokens with an empty error code. On an OAuth error response it returns the
// error code (e.g. "authorization_pending", "slow_down") with a nil error, so
// the caller can decide whether to keep polling. A non-nil error means the
// request itself failed (transport/decoding).
func pollDeviceToken(ctx context.Context, hc *http.Client, apiBase, clientID, deviceCode, verifier string) (Tokens, string, error) {
	form := url.Values{
		"grant_type":    {deviceCodeGrantType},
		"device_code":   {deviceCode},
		"client_id":     {clientID},
		"code_verifier": {verifier},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiBase+"/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return Tokens{}, "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := hc.Do(req)
	if err != nil {
		return Tokens{}, "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var oauthErr struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		if err := json.Unmarshal(body, &oauthErr); err != nil || oauthErr.Error == "" {
			return Tokens{}, "", fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
		}
		return Tokens{}, oauthErr.Error, nil
	}

	var out struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return Tokens{}, "", fmt.Errorf("decode token response: %w", err)
	}
	if out.AccessToken == "" {
		return Tokens{}, "", fmt.Errorf("token response had no access_token")
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
	}, "", nil
}
