package cli

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// kgCapture records what the mock agent service received from a knowledge
// graph command, keeping the escaped path so tests can assert entity IDs
// containing "/" stay encoded as a single path segment.
type kgCapture struct {
	method string
	path   string // EscapedPath, e.g. /v1/orgs/x/knowledge-graph/entities/a%2Fb
	query  url.Values
	hits   int
}

// kgTestServer serves the agent host for knowledge graph reads, recording the
// request and replying with a {status,data,success} envelope.
func kgTestServer(t *testing.T, got *kgCapture, status int, respData string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got.hits++
		got.method = r.Method
		got.path = r.URL.EscapedPath()
		got.query = r.URL.Query()
		w.WriteHeader(status)
		if status >= 400 {
			_, _ = w.Write([]byte(`{"status":404,"error":"not found","success":false}`))
			return
		}
		_, _ = w.Write([]byte(`{"status":200,"data":` + respData + `,"success":true}`))
	}))
}

const kgBase = "/v1/orgs/" + testOrg + "/knowledge-graph"

func TestAIKnowledgeStats(t *testing.T) {
	var got kgCapture
	srv := kgTestServer(t, &got, http.StatusOK, `{"totalNodes":10,"totalEdges":4}`)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	if err := runEdx(t, "ai", "knowledge", "stats"); err != nil {
		t.Fatalf("stats: %v", err)
	}
	if got.method != http.MethodGet {
		t.Errorf("method = %q, want GET", got.method)
	}
	if want := kgBase + "/stats"; got.path != want {
		t.Errorf("path = %q, want %q", got.path, want)
	}
	if len(got.query) != 0 {
		t.Errorf("stats sent unexpected query params: %v", got.query)
	}
}

func TestAIKnowledgeTopologyFlags(t *testing.T) {
	var got kgCapture
	srv := kgTestServer(t, &got, http.StatusOK, `{"nodes":[],"edges":[]}`)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	if err := runEdx(t, "ai", "knowledge", "topology",
		"--limit", "100", "--namespaces", "topology,learned"); err != nil {
		t.Fatalf("topology: %v", err)
	}
	if want := kgBase + "/topology"; got.path != want {
		t.Errorf("path = %q, want %q", got.path, want)
	}
	if v := got.query.Get("limit"); v != "100" {
		t.Errorf("limit = %q, want 100", v)
	}
	if v := got.query.Get("namespaces"); v != "topology,learned" {
		t.Errorf("namespaces = %q, want topology,learned", v)
	}
}

func TestAIKnowledgeTopologyDefaultsOmitParams(t *testing.T) {
	var got kgCapture
	srv := kgTestServer(t, &got, http.StatusOK, `{"nodes":[],"edges":[]}`)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	if err := runEdx(t, "ai", "knowledge", "topology"); err != nil {
		t.Fatalf("topology: %v", err)
	}
	// Unset flags must be omitted so the server applies its own defaults.
	if len(got.query) != 0 {
		t.Errorf("default topology sent query params: %v", got.query)
	}
}

func TestAIKnowledgeSearch(t *testing.T) {
	var got kgCapture
	srv := kgTestServer(t, &got, http.StatusOK, `{"matches":[]}`)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	if err := runEdx(t, "ai", "knowledge", "search", "payment",
		"--types", "Service,Repo", "--min-confidence", "0.7",
		"--source", "github", "--limit", "50", "--cursor", "abc"); err != nil {
		t.Fatalf("search: %v", err)
	}
	if want := kgBase + "/search"; got.path != want {
		t.Errorf("path = %q, want %q", got.path, want)
	}
	for k, want := range map[string]string{
		"q":             "payment",
		"types":         "Service,Repo",
		"minConfidence": "0.7",
		"source":        "github",
		"limit":         "50",
		"cursor":        "abc",
	} {
		if v := got.query.Get(k); v != want {
			t.Errorf("query[%s] = %q, want %q", k, v, want)
		}
	}
}

func TestAIKnowledgeGetEscapesEntityID(t *testing.T) {
	var got kgCapture
	srv := kgTestServer(t, &got, http.StatusOK, `{"node":{}}`)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	// Entity IDs may contain "/" (e.g. ECS ARNs); the ID must stay one path
	// segment on the wire.
	id := testOrg + "::AwsResource::arn:aws:ecs/cluster/task"
	if err := runEdx(t, "ai", "knowledge", "get", id); err != nil {
		t.Fatalf("get: %v", err)
	}
	want := kgBase + "/entities/" + url.PathEscape(id)
	if got.path != want {
		t.Errorf("path = %q, want %q", got.path, want)
	}
}

func TestAIKnowledgeSubgraphHops(t *testing.T) {
	var got kgCapture
	srv := kgTestServer(t, &got, http.StatusOK, `{"root":{},"nodes":[]}`)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	id := testOrg + "::Service::api"
	if err := runEdx(t, "ai", "knowledge", "subgraph", id, "--hops", "3"); err != nil {
		t.Fatalf("subgraph: %v", err)
	}
	want := kgBase + "/entities/" + url.PathEscape(id) + "/subgraph"
	if got.path != want {
		t.Errorf("path = %q, want %q", got.path, want)
	}
	if v := got.query.Get("hops"); v != "3" {
		t.Errorf("hops = %q, want 3", v)
	}
}

func TestAIKnowledgeBlastRadiusMaxHops(t *testing.T) {
	var got kgCapture
	srv := kgTestServer(t, &got, http.StatusOK, `{"root":{},"affected":[]}`)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	id := testOrg + "::Service::api"
	if err := runEdx(t, "ai", "knowledge", "blast-radius", id, "--max-hops", "2"); err != nil {
		t.Fatalf("blast-radius: %v", err)
	}
	want := kgBase + "/entities/" + url.PathEscape(id) + "/blast-radius"
	if got.path != want {
		t.Errorf("path = %q, want %q", got.path, want)
	}
	if v := got.query.Get("maxHops"); v != "2" {
		t.Errorf("maxHops = %q, want 2", v)
	}
}

func TestAIKnowledgeCriticality(t *testing.T) {
	var got kgCapture
	srv := kgTestServer(t, &got, http.StatusOK, `{"basis":"dependency","scores":[]}`)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	if err := runEdx(t, "ai", "knowledge", "criticality",
		"--limit", "10", "--namespaces", "topology"); err != nil {
		t.Fatalf("criticality: %v", err)
	}
	if want := kgBase + "/criticality"; got.path != want {
		t.Errorf("path = %q, want %q", got.path, want)
	}
	if v := got.query.Get("limit"); v != "10" {
		t.Errorf("limit = %q, want 10", v)
	}
	if v := got.query.Get("namespaces"); v != "topology" {
		t.Errorf("namespaces = %q, want topology", v)
	}
}

func TestAIKnowledgeSurfacesServerError(t *testing.T) {
	var got kgCapture
	srv := kgTestServer(t, &got, http.StatusNotFound, "")
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	if err := runEdx(t, "ai", "knowledge", "get", "missing::Service::x"); err == nil {
		t.Fatal("expected error on 404 from the agent service")
	}
}
