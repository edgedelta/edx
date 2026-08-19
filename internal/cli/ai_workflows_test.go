package cli

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/edgedelta/edx/internal/config"
)

const (
	testWorkflowID  = "01KZP23T5Q1VCK2KB9TY250BB8"
	testExecutionID = "0198a7c0-0000-7000-8000-000000000000"
)

// workflowTestServer serves the workflow host, recording the request and
// replying with the AI services' {status,data,success} envelope.
func workflowTestServer(t *testing.T, got *capturedReq, respData string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got.hits++
		got.method = r.Method
		got.path = r.URL.Path
		got.query = r.URL.Query().Encode()
		got.body, _ = io.ReadAll(r.Body)
		_, _ = w.Write([]byte(`{"status":200,"data":` + respData + `,"success":true}`))
	}))
}

// useWorkflowEnv points edx at the given workflow host with token auth so
// workflow requests reach the mock server hermetically.
func useWorkflowEnv(t *testing.T, workflowURL string) {
	t.Helper()
	t.Setenv("EDX_CONFIG", filepath.Join(t.TempDir(), "config.yaml"))
	clearEnv(t)
	t.Setenv(config.EnvAPIToken, "tok-test")
	t.Setenv(config.EnvOrgID, testOrg)
	t.Setenv(config.EnvWorkflowURL, workflowURL)
}

func TestAIWorkflowsList(t *testing.T) {
	var got capturedReq
	srv := workflowTestServer(t, &got, `{"workflows":[{"workflowId":"`+testWorkflowID+`"}]}`)
	defer srv.Close()
	useWorkflowEnv(t, srv.URL)

	if err := runEdx(t, "ai", "workflows", "list", "--limit", "5", "--cursor", "abc", "--order", "descending"); err != nil {
		t.Fatalf("list: %v", err)
	}
	if got.method != http.MethodGet {
		t.Errorf("method = %q, want GET", got.method)
	}
	if want := "/v1/orgs/" + testOrg + "/workflows"; got.path != want {
		t.Errorf("path = %q, want %q", got.path, want)
	}
	if want := "cursor=abc&limit=5&order=descending"; got.query != want {
		t.Errorf("query = %q, want %q", got.query, want)
	}
}

func TestAIWorkflowsGet(t *testing.T) {
	var got capturedReq
	srv := workflowTestServer(t, &got, `{"workflowId":"`+testWorkflowID+`"}`)
	defer srv.Close()
	useWorkflowEnv(t, srv.URL)

	if err := runEdx(t, "ai", "workflows", "get", testWorkflowID); err != nil {
		t.Fatalf("get: %v", err)
	}
	if want := "/v1/orgs/" + testOrg + "/workflows/" + testWorkflowID; got.path != want {
		t.Errorf("path = %q, want %q", got.path, want)
	}
}

func TestAIWorkflowsRunsList(t *testing.T) {
	var got capturedReq
	srv := workflowTestServer(t, &got, `{"executions":[],"totalCount":0}`)
	defer srv.Close()
	useWorkflowEnv(t, srv.URL)

	if err := runEdx(t, "ai", "workflows", "runs", "list", testWorkflowID, "--limit", "10"); err != nil {
		t.Fatalf("runs list: %v", err)
	}
	if got.method != http.MethodGet {
		t.Errorf("method = %q, want GET", got.method)
	}
	if want := "/v1/orgs/" + testOrg + "/workflows/" + testWorkflowID + "/executions"; got.path != want {
		t.Errorf("path = %q, want %q", got.path, want)
	}
	if want := "limit=10"; got.query != want {
		t.Errorf("query = %q, want %q", got.query, want)
	}
}

func TestAIWorkflowsRunsGet(t *testing.T) {
	var got capturedReq
	srv := workflowTestServer(t, &got, `{"executionId":"`+testExecutionID+`"}`)
	defer srv.Close()
	useWorkflowEnv(t, srv.URL)

	if err := runEdx(t, "ai", "workflows", "runs", "get", testWorkflowID, testExecutionID); err != nil {
		t.Fatalf("runs get: %v", err)
	}
	if want := "/v1/orgs/" + testOrg + "/workflows/" + testWorkflowID + "/executions/" + testExecutionID; got.path != want {
		t.Errorf("path = %q, want %q", got.path, want)
	}
}

func TestAIWorkflowsRunsSteps(t *testing.T) {
	var got capturedReq
	srv := workflowTestServer(t, &got, `{"steps":[]}`)
	defer srv.Close()
	useWorkflowEnv(t, srv.URL)

	if err := runEdx(t, "ai", "workflows", "runs", "steps", testWorkflowID, testExecutionID); err != nil {
		t.Fatalf("runs steps: %v", err)
	}
	if want := "/v1/orgs/" + testOrg + "/workflows/" + testWorkflowID + "/executions/" + testExecutionID + "/steps"; got.path != want {
		t.Errorf("path = %q, want %q", got.path, want)
	}
}

// sseWorkflowServer serves the manual-run endpoint, recording the request and
// streaming the given SSE events (already-formatted "event:/data:" blocks).
func sseWorkflowServer(t *testing.T, got *capturedReq, events []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got.hits++
		got.method = r.Method
		got.path = r.URL.Path
		got.body, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "text/event-stream")
		fl, _ := w.(http.Flusher)
		for _, ev := range events {
			_, _ = fmt.Fprint(w, ev)
			if fl != nil {
				fl.Flush()
			}
		}
	}))
}

func TestAIWorkflowsRunStreamsUntilDone(t *testing.T) {
	var got capturedReq
	srv := sseWorkflowServer(t, &got, []string{
		"event: message\ndata: {\"type\":\"execution-create\",\"execution\":{\"executionId\":\"" + testExecutionID + "\"}}\n\n",
		"event: message\ndata: {\"type\":\"step-result\",\"stepIndex\":0,\"stepResult\":{\"type\":\"workflow-complete\",\"done\":true}}\n\n",
		"event: done\ndata: \n\n",
	})
	defer srv.Close()
	useWorkflowEnv(t, srv.URL)

	if err := runEdx(t, "ai", "workflows", "run", testWorkflowID); err != nil {
		t.Fatalf("run: %v", err)
	}
	if got.method != http.MethodPost {
		t.Errorf("method = %q, want POST", got.method)
	}
	if want := "/v1/orgs/" + testOrg + "/workflows/" + testWorkflowID + "/executions"; got.path != want {
		t.Errorf("path = %q, want %q", got.path, want)
	}
	// The manual-run endpoint requires an "input" field; default is {}.
	assertJSONEqual(t, got.body, `{"input":{}}`)
}

func TestAIWorkflowsRunInputFlag(t *testing.T) {
	var got capturedReq
	srv := sseWorkflowServer(t, &got, []string{"event: done\ndata: \n\n"})
	defer srv.Close()
	useWorkflowEnv(t, srv.URL)

	if err := runEdx(t, "ai", "workflows", "run", testWorkflowID, "--input", `{"alert":"cpu high"}`); err != nil {
		t.Fatalf("run: %v", err)
	}
	assertJSONEqual(t, got.body, `{"input":{"alert":"cpu high"}}`)
}

func TestAIWorkflowsRunPlainStringInput(t *testing.T) {
	var got capturedReq
	srv := sseWorkflowServer(t, &got, []string{"event: done\ndata: \n\n"})
	defer srv.Close()
	useWorkflowEnv(t, srv.URL)

	// A non-JSON input value is sent as a JSON string.
	if err := runEdx(t, "ai", "workflows", "run", testWorkflowID, "--input", "investigate cpu"); err != nil {
		t.Fatalf("run: %v", err)
	}
	assertJSONEqual(t, got.body, `{"input":"investigate cpu"}`)
}

func TestAIWorkflowsRunErrorEventFails(t *testing.T) {
	var got capturedReq
	srv := sseWorkflowServer(t, &got, []string{
		"event: message\ndata: {\"type\":\"run-error\",\"message\":\"Workflow not active\"}\n\n",
		"event: done\ndata: \n\n",
	})
	defer srv.Close()
	useWorkflowEnv(t, srv.URL)

	if err := runEdx(t, "ai", "workflows", "run", testWorkflowID); err == nil {
		t.Fatal("run should fail when the stream reports run-error")
	}
}
