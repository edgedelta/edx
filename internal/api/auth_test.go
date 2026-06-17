package api

import (
	"context"
	"net/http"
	"testing"
)

type fakeTokenSource struct {
	tok string
	err error
}

func (f fakeTokenSource) Token(context.Context) (string, error) { return f.tok, f.err }

func newReq() *http.Request {
	r, _ := http.NewRequest(http.MethodGet, "http://example/test", nil)
	return r
}

func TestAuthTokenMode(t *testing.T) {
	a := &Auth{APIToken: "tok123"}
	req := newReq()
	if err := a.apply(context.Background(), req, ServiceAPI, "/v1/orgs/o/pipelines"); err != nil {
		t.Fatal(err)
	}
	if got := req.Header.Get("X-ED-API-Token"); got != "tok123" {
		t.Errorf("X-ED-API-Token = %q", got)
	}
	if req.Header.Get("Authorization") != "" {
		t.Error("token mode should not set Authorization")
	}
}

func TestAuthOAuthMainAPI(t *testing.T) {
	a := &Auth{OAuth: fakeTokenSource{tok: "jwt"}, APIToken: "tok", APIDomain: "api.edgedelta.com"}
	req := newReq()
	if err := a.apply(context.Background(), req, ServiceAPI, "/v1/orgs/o/pipelines"); err != nil {
		t.Fatal(err)
	}
	if got := req.Header.Get("Authorization"); got != "Bearer jwt" {
		t.Errorf("Authorization = %q", got)
	}
	if req.Header.Get("X-ED-API-Domain") != "" {
		t.Error("main API must NOT receive X-ED-API-Domain")
	}
	if req.Header.Get("X-ED-API-Token") != "" {
		t.Error("oauth mode should not set X-ED-API-Token")
	}
}

func TestAuthOAuthAIServiceSetsDomain(t *testing.T) {
	a := &Auth{OAuth: fakeTokenSource{tok: "jwt"}, APIDomain: "api.staging.edgedelta.com"}
	req := newReq()
	if err := a.apply(context.Background(), req, ServiceChat, "/v1/orgs/o/issues"); err != nil {
		t.Fatal(err)
	}
	if req.Header.Get("Authorization") != "Bearer jwt" {
		t.Error("expected Bearer auth on chat service")
	}
	if got := req.Header.Get("X-ED-API-Domain"); got != "api.staging.edgedelta.com" {
		t.Errorf("X-ED-API-Domain = %q (AI services require it)", got)
	}
}

func TestAuthOAuthUnsupportedFallsBackToToken(t *testing.T) {
	orig := oauthUnsupported
	oauthUnsupported = []endpointRule{{svc: ServiceChat, contains: "/issues"}}
	defer func() { oauthUnsupported = orig }()

	a := &Auth{OAuth: fakeTokenSource{tok: "jwt"}, APIToken: "tok", APIDomain: "d"}
	req := newReq()
	if err := a.apply(context.Background(), req, ServiceChat, "/v1/orgs/o/issues"); err != nil {
		t.Fatal(err)
	}
	if req.Header.Get("X-ED-API-Token") != "tok" {
		t.Error("oauth-unsupported endpoint should fall back to API token")
	}
	if req.Header.Get("Authorization") != "" {
		t.Error("oauth-unsupported endpoint must not use Bearer")
	}

	// A different chat path still uses OAuth.
	req2 := newReq()
	_ = a.apply(context.Background(), req2, ServiceChat, "/v1/orgs/o/channels")
	if req2.Header.Get("Authorization") != "Bearer jwt" {
		t.Error("supported endpoint should still use OAuth")
	}
}

func TestAuthOAuthUnsupportedNoTokenErrors(t *testing.T) {
	orig := oauthUnsupported
	oauthUnsupported = []endpointRule{{svc: ServiceChat, contains: "/issues"}}
	defer func() { oauthUnsupported = orig }()

	a := &Auth{OAuth: fakeTokenSource{tok: "jwt"}} // no APIToken fallback
	if err := a.apply(context.Background(), newReq(), ServiceChat, "/v1/orgs/o/issues"); err == nil {
		t.Error("expected error when oauth-unsupported endpoint has no token fallback")
	}
}

func TestAuthNoCredentials(t *testing.T) {
	a := &Auth{}
	if err := a.apply(context.Background(), newReq(), ServiceAPI, "/v1/orgs/o/pipelines"); err == nil {
		t.Error("expected error with no credentials")
	}
}
