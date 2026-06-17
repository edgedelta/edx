package cli

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/edgedelta/edx/internal/config"
	"github.com/edgedelta/edx/internal/oauth"
)

// oauthTokenSource implements api.TokenSource. It hands out the profile's OAuth
// access token, refreshing it via the refresh token when it is within a minute
// of expiry and persisting the renewed tokens back to the config file so they
// survive across invocations.
type oauthTokenSource struct {
	profile  string
	env      string
	orgID    string
	apiBase  string
	clientID string

	mu      sync.Mutex
	access  string
	refresh string
	expiry  time.Time
}

func newOAuthTokenSource(r *config.Resolved, expiry time.Time) *oauthTokenSource {
	return &oauthTokenSource{
		profile:  r.Profile,
		env:      r.Env,
		orgID:    r.OrgID,
		apiBase:  r.APIURL,
		clientID: r.OAuthClientID,
		access:   r.OAuthAccessToken,
		refresh:  r.OAuthRefreshToken,
		expiry:   expiry,
	}
}

func (s *oauthTokenSource) Token(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.access != "" && time.Now().Before(s.expiry.Add(-60*time.Second)) {
		return s.access, nil
	}
	if s.refresh == "" {
		return "", fmt.Errorf("access token expired and no refresh token available; run `edx auth login --oauth` again")
	}

	t, err := oauth.Refresh(ctx, s.apiBase, s.clientID, s.refresh, nil)
	if err != nil {
		return "", fmt.Errorf("token refresh failed (run `edx auth login --oauth` again): %w", err)
	}
	s.access = t.AccessToken
	if t.RefreshToken != "" {
		s.refresh = t.RefreshToken
	}
	s.expiry = t.Expiry
	if err := config.SaveOAuthTokens(s.profile, s.env, s.orgID, s.clientID, s.access, s.refresh, s.expiry); err != nil {
		// Non-fatal: the in-memory token still works for this run.
		fmt.Printf("warning: could not persist refreshed token: %v\n", err)
	}
	return s.access, nil
}
