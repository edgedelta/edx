package cli

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/edgedelta/edx/internal/config"
)

// capturedReq records what the mock agent service received.
type capturedReq struct {
	method string
	path   string
	body   []byte
	hits   int
}

// agentTestServer serves the agent host, recording the request and replying
// with a {status,data,success} envelope wrapping the given agent object.
func agentTestServer(t *testing.T, got *capturedReq, status int, respAgent string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got.hits++
		got.method = r.Method
		got.path = r.URL.Path
		got.body, _ = io.ReadAll(r.Body)
		w.WriteHeader(status)
		if status >= 400 {
			_, _ = w.Write([]byte(`{"status":` + strconv.Itoa(status) + `,"error":"boom","success":false}`))
			return
		}
		_, _ = w.Write([]byte(`{"status":200,"data":` + respAgent + `,"success":true}`))
	}))
}

// useAgentEnv points edx at the given agent host with token auth so agent
// requests reach the mock server hermetically.
func useAgentEnv(t *testing.T, agentURL string) {
	t.Helper()
	t.Setenv("EDX_CONFIG", filepath.Join(t.TempDir(), "config.yaml"))
	clearEnv(t)
	t.Setenv(config.EnvAPIToken, "tok-test")
	t.Setenv(config.EnvOrgID, testOrg)
	t.Setenv(config.EnvAgentURL, agentURL)
}

// writeTemp writes content to a temp file and returns its path.
func writeTemp(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "agent.json")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

const testAgentID = "01HXAGENT0000000000000000"

func TestAIAgentsUpdatePutsToAgentService(t *testing.T) {
	var got capturedReq
	srv := agentTestServer(t, &got, http.StatusOK, `{"id":"`+testAgentID+`","name":"Ollie"}`)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	file := writeTemp(t, `{"name":"Ollie","masterPrompt":"be helpful"}`)
	if err := runEdx(t, "ai", "agents", "update", testAgentID, "--file", file, "--yes"); err != nil {
		t.Fatalf("update: %v", err)
	}

	if got.method != http.MethodPut {
		t.Errorf("method = %q, want PUT", got.method)
	}
	wantPath := "/v1/orgs/" + testOrg + "/agents/" + testAgentID
	if got.path != wantPath {
		t.Errorf("path = %q, want %q", got.path, wantPath)
	}
	assertJSONEqual(t, got.body, `{"name":"Ollie","masterPrompt":"be helpful"}`)
}

func TestAIAgentsUpdateUnwrapsGetEnvelope(t *testing.T) {
	var got capturedReq
	srv := agentTestServer(t, &got, http.StatusOK, `{"id":"`+testAgentID+`"}`)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	// A file saved directly from "edx ai agents get" — the full envelope.
	file := writeTemp(t, `{"status":200,"success":true,"data":{"id":"`+testAgentID+`","name":"Ollie","model":"claude-opus-4-8"}}`)
	if err := runEdx(t, "ai", "agents", "update", testAgentID, "--file", file, "--yes"); err != nil {
		t.Fatalf("update: %v", err)
	}

	// The server must receive only the inner agent object, not the envelope.
	assertJSONEqual(t, got.body, `{"id":"`+testAgentID+`","name":"Ollie","model":"claude-opus-4-8"}`)
}

func TestAIAgentsUpdateBareBodyPassthrough(t *testing.T) {
	var got capturedReq
	srv := agentTestServer(t, &got, http.StatusOK, `{}`)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	// A hand-written bare body with no "data" key is sent unchanged.
	file := writeTemp(t, `{"model":"claude-opus-4-8"}`)
	if err := runEdx(t, "ai", "agents", "update", testAgentID, "--file", file, "--yes"); err != nil {
		t.Fatalf("update: %v", err)
	}
	assertJSONEqual(t, got.body, `{"model":"claude-opus-4-8"}`)
}

func TestAIAgentsUpdateStdin(t *testing.T) {
	var got capturedReq
	srv := agentTestServer(t, &got, http.StatusOK, `{}`)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	feedStdin(t, `{"data":{"name":"FromStdin"}}`, func() {
		if err := runEdx(t, "ai", "agents", "update", testAgentID, "--file", "-", "--yes"); err != nil {
			t.Fatalf("update: %v", err)
		}
	})
	assertJSONEqual(t, got.body, `{"name":"FromStdin"}`)
}

func TestAIAgentsUpdateAbortsWithoutYes(t *testing.T) {
	var got capturedReq
	srv := agentTestServer(t, &got, http.StatusOK, `{}`)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	file := writeTemp(t, `{"name":"Ollie"}`)
	// Decline the confirmation prompt.
	feedStdin(t, "n\n", func() {
		err := runEdx(t, "ai", "agents", "update", testAgentID, "--file", file)
		if err != errAborted {
			t.Fatalf("err = %v, want errAborted", err)
		}
	})
	if got.hits != 0 {
		t.Errorf("declined update still sent %d request(s)", got.hits)
	}
}

func TestAIAgentsUpdateSurfacesServerError(t *testing.T) {
	var got capturedReq
	srv := agentTestServer(t, &got, http.StatusNotFound, "")
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	file := writeTemp(t, `{"name":"Ollie"}`)
	err := runEdx(t, "ai", "agents", "update", testAgentID, "--file", file, "--yes")
	if err == nil {
		t.Fatal("expected error on 404 from the agent service")
	}
}

// rwCapture records a GET (backfill) followed by a PUT for prompt-mode tests.
type rwCapture struct {
	getHits int
	putHits int
	putBody []byte
}

// agentRMWServer serves GET /agents/{id} with the given current agent object and
// captures the subsequent PUT body — the read-modify-write path prompt mode uses.
func agentRMWServer(t *testing.T, current string, rw *rwCapture) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rw.getHits++
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":200,"data":` + current + `,"success":true}`))
		case http.MethodPut:
			rw.putHits++
			rw.putBody, _ = io.ReadAll(r.Body)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":200,"data":` + current + `,"success":true}`))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}

const currentAgent = `{"id":"` + testAgentID + `","name":"Ollie","masterPrompt":"OLD master","userPrompt":"OLD user","toolingPrompt":"OLD tooling","model":"gpt-5.2","modelTemperature":0.1}`

func TestAIAgentsUpdatePromptModeSendsOnlyPrompts(t *testing.T) {
	var rw rwCapture
	srv := agentRMWServer(t, currentAgent, &rw)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	if err := runEdx(t, "ai", "agents", "update", testAgentID, "--master-prompt", "NEW master", "--yes"); err != nil {
		t.Fatalf("update: %v", err)
	}
	if rw.getHits != 1 || rw.putHits != 1 {
		t.Fatalf("want 1 GET + 1 PUT, got GET=%d PUT=%d", rw.getHits, rw.putHits)
	}
	// masterPrompt is the new value; userPrompt is backfilled from current; no
	// model/temperature/tooling are sent (untouched by the partial merge).
	assertJSONEqual(t, rw.putBody, `{"masterPrompt":"NEW master","userPrompt":"OLD user"}`)
}

func TestAIAgentsUpdatePromptModeFromFile(t *testing.T) {
	var rw rwCapture
	srv := agentRMWServer(t, currentAgent, &rw)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	pf := writeTemp(t, "FROM FILE master")
	if err := runEdx(t, "ai", "agents", "update", testAgentID, "--master-prompt", "@"+pf, "--yes"); err != nil {
		t.Fatalf("update: %v", err)
	}
	assertJSONEqual(t, rw.putBody, `{"masterPrompt":"FROM FILE master","userPrompt":"OLD user"}`)
}

func TestAIAgentsUpdatePromptModeIncludesToolingWhenSet(t *testing.T) {
	var rw rwCapture
	srv := agentRMWServer(t, currentAgent, &rw)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	if err := runEdx(t, "ai", "agents", "update", testAgentID,
		"--user-prompt", "NEW user", "--tooling-prompt", "NEW tooling", "--yes"); err != nil {
		t.Fatalf("update: %v", err)
	}
	assertJSONEqual(t, rw.putBody, `{"masterPrompt":"OLD master","userPrompt":"NEW user","toolingPrompt":"NEW tooling"}`)
}

func TestAIAgentsUpdateFileAndPromptFlagConflict(t *testing.T) {
	var rw rwCapture
	srv := agentRMWServer(t, currentAgent, &rw)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	file := writeTemp(t, `{"name":"x"}`)
	err := runEdx(t, "ai", "agents", "update", testAgentID, "--file", file, "--master-prompt", "x", "--yes")
	if err == nil {
		t.Fatal("expected error combining --file with a prompt flag")
	}
	if rw.getHits != 0 || rw.putHits != 0 {
		t.Errorf("conflicting flags still hit the server: GET=%d PUT=%d", rw.getHits, rw.putHits)
	}
}

func TestAIAgentsUpdateRequiresFileOrPromptFlag(t *testing.T) {
	var rw rwCapture
	srv := agentRMWServer(t, currentAgent, &rw)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	if err := runEdx(t, "ai", "agents", "update", testAgentID, "--yes"); err == nil {
		t.Fatal("expected error when neither --file nor a prompt flag is given")
	}
}

func TestAIAgentsUpdatePromptModeMissingRequiredPrompt(t *testing.T) {
	var rw rwCapture
	// Current teammate has no userPrompt to backfill.
	srv := agentRMWServer(t, `{"id":"`+testAgentID+`","masterPrompt":"has master"}`, &rw)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	err := runEdx(t, "ai", "agents", "update", testAgentID, "--master-prompt", "new", "--yes")
	if err == nil {
		t.Fatal("expected error: userPrompt required but not backfillable")
	}
	if rw.putHits != 0 {
		t.Errorf("must not PUT when a required prompt is missing (PUT=%d)", rw.putHits)
	}
}

func TestResolvePromptValue(t *testing.T) {
	if got, err := resolvePromptValue("literal text"); err != nil || got != "literal text" {
		t.Errorf("inline: got %q err %v", got, err)
	}
	pf := writeTemp(t, "file contents")
	if got, err := resolvePromptValue("@" + pf); err != nil || got != "file contents" {
		t.Errorf("@file: got %q err %v", got, err)
	}
}

func TestUnwrapDataEnvelope(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"envelope", `{"status":200,"data":{"name":"x"},"success":true}`, `{"name":"x"}`},
		{"bare object", `{"name":"x"}`, `{"name":"x"}`},
		{"data is null", `{"data":null}`, `{"data":null}`},
		{"data is array", `{"data":[1,2]}`, `{"data":[1,2]}`},
		{"not json", `not json`, `not json`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := string(unwrapDataEnvelope([]byte(c.in))); got != c.want {
				t.Errorf("unwrapDataEnvelope(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

// assertJSONEqual compares two JSON documents for semantic equality.
func assertJSONEqual(t *testing.T, got []byte, want string) {
	t.Helper()
	var g, w any
	if err := json.Unmarshal(got, &g); err != nil {
		t.Fatalf("got is not JSON: %v (%s)", err, got)
	}
	if err := json.Unmarshal([]byte(want), &w); err != nil {
		t.Fatalf("want is not JSON: %v", err)
	}
	gb, _ := json.Marshal(g)
	wb, _ := json.Marshal(w)
	if string(gb) != string(wb) {
		t.Errorf("body = %s, want %s", gb, wb)
	}
}
