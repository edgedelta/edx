package cli

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type logsRequest struct {
	method, path string
	query        url.Values
	body         string
}

// logsAPIServer records each request's method, path, query and body, and
// replies with resp.
func logsAPIServer(t *testing.T, got *[]logsRequest, resp string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		*got = append(*got, logsRequest{r.Method, r.URL.Path, r.URL.Query(), string(body)})
		_, _ = w.Write([]byte(resp))
	}))
}

func TestLogsSearchScope(t *testing.T) {
	for _, tc := range []struct {
		args      []string
		wantScope string
	}{
		{nil, "log"},
		{[]string{"--scope", "audit"}, "audit"},
	} {
		var got []logsRequest
		srv := logsAPIServer(t, &got, `{"items":[],"next_cursor":""}`)
		useAPIEnv(t, srv.URL)

		args := append([]string{"logs", "search", "--query", `method:"DELETE"`}, tc.args...)
		if err := runEdx(t, args...); err != nil {
			t.Fatalf("search %v: %v", tc.args, err)
		}
		srv.Close()
		if len(got) != 1 {
			t.Fatalf("requests = %d, want 1", len(got))
		}
		if !strings.HasSuffix(got[0].path, "/logs/log_search/search") {
			t.Errorf("path = %q", got[0].path)
		}
		if s := got[0].query.Get("scope"); s != tc.wantScope {
			t.Errorf("scope = %q, want %q", s, tc.wantScope)
		}
		if q := got[0].query.Get("query"); q != `method:"DELETE"` {
			t.Errorf("query = %q", q)
		}
	}
}

func TestLogsGraphAuditUsesLogSearchGraph(t *testing.T) {
	var got []logsRequest
	srv := logsAPIServer(t, &got, `{"records":[]}`)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	if err := runEdx(t, "logs", "graph", "--scope", "audit",
		"--query", `status.code:"500"`, "--group-by", "user.email", "--lookback", "24h"); err != nil {
		t.Fatalf("graph: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("requests = %d, want 1", len(got))
	}
	r := got[0]
	if r.method != http.MethodGet || !strings.HasSuffix(r.path, "/logs/log_search/graph") {
		t.Errorf("request = %s %s, want GET .../logs/log_search/graph", r.method, r.path)
	}
	for k, want := range map[string]string{
		"scope":    "audit",
		"query":    `status.code:"500" by {user.email}`,
		"lookback": "24h",
	} {
		if v := r.query.Get(k); v != want {
			t.Errorf("query[%s] = %q, want %q", k, v, want)
		}
	}
}

func TestLogsGraphLogScopeUsesFormulaGraph(t *testing.T) {
	var got []logsRequest
	srv := logsAPIServer(t, &got, `{"records":[]}`)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	if err := runEdx(t, "logs", "graph", "--query", "*"); err != nil {
		t.Fatalf("graph: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("requests = %d, want 1", len(got))
	}
	r := got[0]
	if r.method != http.MethodPost || !strings.HasSuffix(r.path, "/graph") {
		t.Errorf("request = %s %s, want POST .../graph", r.method, r.path)
	}
	var payload struct {
		Queries map[string]struct{ Scope string } `json:"queries"`
	}
	if err := json.Unmarshal([]byte(r.body), &payload); err != nil {
		t.Fatalf("body: %v", err)
	}
	if s := payload.Queries["Q1"].Scope; s != "log" {
		t.Errorf("Q1 scope = %q, want log", s)
	}
}

func TestLogsRejectsUnknownScope(t *testing.T) {
	for _, sub := range []string{"search", "graph"} {
		err := runEdx(t, "logs", sub, "--scope", "metric")
		if err == nil || !strings.Contains(err.Error(), "must be log or audit") {
			t.Errorf("logs %s --scope metric: err = %v, want scope error", sub, err)
		}
	}
}
