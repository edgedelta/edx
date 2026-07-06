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

func TestAuthCookieMainAPI(t *testing.T) {
	a := &Auth{SessionCookie: "sess-abc"}
	req := newReq()
	if err := a.apply(context.Background(), req, ServiceAPI, "/v1/orgs/o/facet_keys"); err != nil {
		t.Fatal(err)
	}
	if got := req.Header.Get("Cookie"); got != "ed-admin-session=sess-abc" {
		t.Errorf("Cookie = %q", got)
	}
	if req.Header.Get("Authorization") != "" || req.Header.Get("X-ED-API-Token") != "" {
		t.Error("cookie mode should not set Authorization or X-ED-API-Token")
	}
}

func TestAuthCookieAIServiceUsesJWT(t *testing.T) {
	// In cookie mode the AI hosts can't take a cookie, so apply exchanges it for
	// a Bearer JWT (via the CookieJWT source) and adds X-ED-API-Domain.
	a := &Auth{SessionCookie: "sess-abc", CookieJWT: fakeTokenSource{tok: "jwt"}, APIDomain: "api.edgedelta.com"}
	for _, svc := range []Service{ServiceChat, ServiceAgent} {
		req := newReq()
		if err := a.apply(context.Background(), req, svc, "/v1/orgs/o/issues"); err != nil {
			t.Fatalf("service %v: %v", svc, err)
		}
		if got := req.Header.Get("Authorization"); got != "Bearer jwt" {
			t.Errorf("service %v Authorization = %q", svc, got)
		}
		if got := req.Header.Get("X-ED-API-Domain"); got != "api.edgedelta.com" {
			t.Errorf("service %v X-ED-API-Domain = %q", svc, got)
		}
		if req.Header.Get("Cookie") != "" {
			t.Errorf("service %v must NOT receive the cookie", svc)
		}
	}
}

func TestAuthCookieAIServiceNoSourceErrors(t *testing.T) {
	a := &Auth{SessionCookie: "sess-abc"} // no CookieJWT configured
	if err := a.apply(context.Background(), newReq(), ServiceChat, "/v1/orgs/o/issues"); err == nil {
		t.Error("expected error when AI service has no cookie-JWT source")
	}
}

func TestAuthNoCredentials(t *testing.T) {
	a := &Auth{}
	if err := a.apply(context.Background(), newReq(), ServiceAPI, "/v1/orgs/o/pipelines"); err == nil {
		t.Error("expected error with no credentials")
	}
}
