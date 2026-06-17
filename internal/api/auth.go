package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// TokenSource yields a currently-valid OAuth access token, refreshing it
// transparently when it has expired. Implemented by the CLI layer (which can
// persist refreshed tokens); kept as an interface here to avoid a config
// dependency in the api package.
type TokenSource interface {
	Token(ctx context.Context) (string, error)
}

// Auth decides, per request, which credential to attach. It is the single
// place that knows about the two auth schemes and which endpoints support
// each — effectively an auth middleware in front of every request.
type Auth struct {
	// APIToken is the static token (X-ED-API-Token). May be empty when the
	// profile authenticates purely via OAuth.
	APIToken string

	// OAuth, when non-nil, is the profile's OAuth credential. Endpoints that
	// support OAuth get an "Authorization: Bearer <jwt>" header instead of the
	// API token.
	OAuth TokenSource

	// APIDomain is the bare host of the main API (e.g. "api.edgedelta.com").
	// The AI services require it as X-ED-API-Domain to validate a Bearer JWT's
	// issuer/audience; the main API ignores it.
	APIDomain string
}

// endpointRule matches a request by service and a substring of its path.
type endpointRule struct {
	svc      Service
	contains string
}

// oauthUnsupported lists endpoints that do NOT yet accept an OAuth Bearer
// token and must fall back to API-token auth. This is the "add them to the
// middleware" hook: as the backend wires OAuth onto more routes, delete the
// matching entry here.
//
// Matching is by service + a substring of the request path (which always
// includes the /v1/orgs/<id> prefix), e.g. {ServiceChat, "/issues"}.
var oauthUnsupported = []endpointRule{
	// NOTE: populate with the AI endpoints that reject OAuth in your
	// deployment. Left empty by default — every endpoint is assumed
	// OAuth-capable until proven otherwise. Example:
	//   {svc: ServiceAgent, contains: "/agents"},
}

// supportsOAuth reports whether (svc, path) may use an OAuth Bearer token.
func supportsOAuth(svc Service, path string) bool {
	for _, r := range oauthUnsupported {
		if r.svc == svc && strings.Contains(path, r.contains) {
			return false
		}
	}
	return true
}

// apply sets the appropriate auth header on req for a request to svc at path.
func (a *Auth) apply(ctx context.Context, req *http.Request, svc Service, path string) error {
	if a.OAuth != nil && supportsOAuth(svc, path) {
		tok, err := a.OAuth.Token(ctx)
		if err != nil {
			return fmt.Errorf("oauth: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+tok)
		// AI services validate the JWT against the API domain; the main API
		// does not need it.
		if svc != ServiceAPI && a.APIDomain != "" {
			req.Header.Set("X-ED-API-Domain", a.APIDomain)
		}
		return nil
	}

	if a.APIToken == "" {
		if a.OAuth != nil {
			return fmt.Errorf("this endpoint does not support OAuth yet and no API token is configured for fallback; run `edx auth login --token <token>` to add one")
		}
		return fmt.Errorf("no credentials configured")
	}
	req.Header.Set("X-ED-API-Token", a.APIToken)
	return nil
}
