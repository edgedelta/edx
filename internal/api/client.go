// Package api implements the Edge Delta REST API client used by edx.
package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Service selects which Edge Delta host a request targets. The main API and
// the AI Teammate services (chat, agent) live on different hosts, so a request
// must name the service it belongs to.
type Service int

const (
	// ServiceAPI is the main Edge Delta API (api.edgedelta.com).
	ServiceAPI Service = iota
	// ServiceChat is the AI Teammate chat service (chat.ai.edgedelta.com):
	// issues, threads, channels, messages.
	ServiceChat
	// ServiceAgent is the AI Teammate agent service (agent.ai.edgedelta.com):
	// teammates (agents) and their versions.
	ServiceAgent
)

// Client talks to the Edge Delta services. Authentication uses the
// X-ED-API-Token header, which every service accepts for org-scoped,
// token-authenticated requests.
type Client struct {
	// BaseURL is the main API host; ChatURL and AgentURL are the AI Teammate
	// service hosts. All share the /v1/orgs/{org_id} path prefix.
	BaseURL  string
	ChatURL  string
	AgentURL string

	OrgID    string
	APIToken string

	HTTP      *http.Client
	UserAgent string

	// MaxRetries controls retry attempts for transient failures (429/5xx).
	MaxRetries int
}

// New builds a client with sane transport defaults. apiURL, chatURL and
// agentURL are the base URLs for the main API and the AI Teammate services.
func New(apiURL, chatURL, agentURL, orgID, token string, timeout time.Duration) *Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          64,
		MaxIdleConnsPerHost:   16,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{MinVersion: tls.VersionTLS12},
	}
	return &Client{
		BaseURL:    strings.TrimRight(apiURL, "/"),
		ChatURL:    strings.TrimRight(chatURL, "/"),
		AgentURL:   strings.TrimRight(agentURL, "/"),
		OrgID:      orgID,
		APIToken:   token,
		HTTP:       &http.Client{Transport: transport, Timeout: timeout},
		UserAgent:  "edx",
		MaxRetries: 3,
	}
}

// baseFor returns the base URL for a service.
func (c *Client) baseFor(svc Service) string {
	switch svc {
	case ServiceChat:
		return c.ChatURL
	case ServiceAgent:
		return c.AgentURL
	default:
		return c.BaseURL
	}
}

// Error is a non-2xx API response.
type Error struct {
	Status int
	Method string
	URL    string
	Body   string
}

func (e *Error) Error() string {
	body := strings.TrimSpace(e.Body)
	if len(body) > 2048 {
		body = body[:2048] + "..."
	}
	return fmt.Sprintf("API error %d on %s %s: %s", e.Status, e.Method, e.URL, body)
}

// OrgPath prefixes p with /v1/orgs/{org_id}.
func (c *Client) OrgPath(p string) string {
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return "/v1/orgs/" + c.OrgID + p
}

// Do performs an HTTP request against the main API and returns the response
// body. path must be absolute (e.g. /v1/orgs/<id>/pipelines). Transient errors
// (429, 502, 503, 504, connection failures) are retried with backoff.
func (c *Client) Do(ctx context.Context, method, path string, query url.Values, body []byte) ([]byte, error) {
	return c.DoOn(ctx, ServiceAPI, method, path, query, body)
}

// DoOn performs an HTTP request against the named service and returns the
// response body. path must be absolute (e.g. /v1/orgs/<id>/issues).
func (c *Client) DoOn(ctx context.Context, svc Service, method, path string, query url.Values, body []byte) ([]byte, error) {
	u := c.baseFor(svc) + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var lastErr error
	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := time.Duration(1<<uint(attempt-1)) * 500 * time.Millisecond
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		var reader io.Reader
		if body != nil {
			reader = bytes.NewReader(body)
		}
		req, err := http.NewRequestWithContext(ctx, method, u, reader)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", c.UserAgent)
		req.Header.Set("X-ED-API-Token", c.APIToken)

		resp, err := c.HTTP.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		respBody, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			lastErr = readErr
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return respBody, nil
		}

		apiErr := &Error{Status: resp.StatusCode, Method: method, URL: u, Body: string(respBody)}
		switch resp.StatusCode {
		case http.StatusTooManyRequests, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
			lastErr = apiErr
			continue
		default:
			return nil, apiErr
		}
	}
	return nil, fmt.Errorf("request failed after %d attempts: %w", c.MaxRetries+1, lastErr)
}

// Get performs a GET against an org-relative path.
func (c *Client) Get(ctx context.Context, orgRelPath string, query url.Values) ([]byte, error) {
	return c.Do(ctx, http.MethodGet, c.OrgPath(orgRelPath), query, nil)
}

// Post performs a POST against an org-relative path.
func (c *Client) Post(ctx context.Context, orgRelPath string, query url.Values, body []byte) ([]byte, error) {
	return c.Do(ctx, http.MethodPost, c.OrgPath(orgRelPath), query, body)
}

// Put performs a PUT against an org-relative path.
func (c *Client) Put(ctx context.Context, orgRelPath string, query url.Values, body []byte) ([]byte, error) {
	return c.Do(ctx, http.MethodPut, c.OrgPath(orgRelPath), query, body)
}

// Delete performs a DELETE against an org-relative path.
func (c *Client) Delete(ctx context.Context, orgRelPath string, query url.Values, body []byte) ([]byte, error) {
	return c.Do(ctx, http.MethodDelete, c.OrgPath(orgRelPath), query, body)
}

// GetFrom performs a GET against an org-relative path on the named service.
func (c *Client) GetFrom(ctx context.Context, svc Service, orgRelPath string, query url.Values) ([]byte, error) {
	return c.DoOn(ctx, svc, http.MethodGet, c.OrgPath(orgRelPath), query, nil)
}

// PostFrom performs a POST against an org-relative path on the named service.
func (c *Client) PostFrom(ctx context.Context, svc Service, orgRelPath string, query url.Values, body []byte) ([]byte, error) {
	return c.DoOn(ctx, svc, http.MethodPost, c.OrgPath(orgRelPath), query, body)
}

// PutFrom performs a PUT against an org-relative path on the named service.
func (c *Client) PutFrom(ctx context.Context, svc Service, orgRelPath string, query url.Values, body []byte) ([]byte, error) {
	return c.DoOn(ctx, svc, http.MethodPut, c.OrgPath(orgRelPath), query, body)
}

// DeleteFrom performs a DELETE against an org-relative path on the named service.
func (c *Client) DeleteFrom(ctx context.Context, svc Service, orgRelPath string, query url.Values, body []byte) ([]byte, error) {
	return c.DoOn(ctx, svc, http.MethodDelete, c.OrgPath(orgRelPath), query, body)
}
