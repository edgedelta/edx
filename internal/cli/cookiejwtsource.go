package cli

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/edgedelta/edx/internal/oauth"
)

// cookieJWTSource implements api.TokenSource for a support (cookie) session. It
// exchanges the ed-admin-session cookie for a short-lived Bearer JWT via
// oauth.GetJWTFromCookie, caching it and re-minting when it nears expiry. This
// is what lets a cookie session reach the AI Teammate hosts, which accept a
// Bearer JWT rather than a cookie. There is no refresh token — a fresh JWT is
// obtained by calling the endpoint again with the (still-valid) cookie.
type cookieJWTSource struct {
	apiBase string
	cookie  string
	hc      *http.Client

	mu     sync.Mutex
	token  string
	expiry time.Time
}

func newCookieJWTSource(apiBase, cookie string, timeout time.Duration) *cookieJWTSource {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &cookieJWTSource{
		apiBase: apiBase,
		cookie:  cookie,
		hc:      &http.Client{Timeout: timeout},
	}
}

func (s *cookieJWTSource) Token(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.token != "" && time.Now().Before(s.expiry.Add(-60*time.Second)) {
		return s.token, nil
	}
	tok, exp, err := oauth.GetJWTFromCookie(ctx, s.apiBase, s.cookie, s.hc)
	if err != nil {
		return "", err
	}
	s.token = tok
	s.expiry = exp
	return s.token, nil
}
