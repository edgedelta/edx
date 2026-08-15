package cli

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

// mustParseQuery parses a raw query string, failing the test on error.
func mustParseQuery(t *testing.T, raw string) url.Values {
	t.Helper()
	q, err := url.ParseQuery(raw)
	if err != nil {
		t.Fatalf("bad query %q: %v", raw, err)
	}
	return q
}

// rehydrationReq records one request received by the mock API server.
type rehydrationReq struct {
	method   string
	path     string
	rawQuery string
	body     []byte
}

// rehydrationServer serves canned responses keyed by "METHOD /org-relative-path"
// and records every request it receives.
func rehydrationServer(t *testing.T, responses map[string]string, log *[]rehydrationReq) *httptest.Server {
	t.Helper()
	prefix := "/v1/orgs/" + testOrg
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		*log = append(*log, rehydrationReq{r.Method, r.URL.Path, r.URL.RawQuery, body})
		key := r.Method + " " + strings.TrimPrefix(r.URL.Path, prefix)
		resp, ok := responses[key]
		if !ok {
			t.Errorf("unexpected request: %s", key)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte(resp))
	}))
}

func TestRehydrationsListQueriesV2(t *testing.T) {
	var log []rehydrationReq
	srv := rehydrationServer(t, map[string]string{
		"GET /rehydration_v2": `{"rehydrations":[],"next_cursor":""}`,
	}, &log)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	if err := runEdx(t, "rehydrations", "list", "--query", `ed.tag:"prod"`, "--lookback", "30m", "--limit", "5"); err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(log) != 1 {
		t.Fatalf("requests = %d, want 1", len(log))
	}
	q := mustParseQuery(t, log[0].rawQuery)
	if q.Get("query") != `ed.tag:"prod"` || q.Get("lookback") != "30m" || q.Get("limit") != "5" {
		t.Errorf("query params = %q", log[0].rawQuery)
	}
}

func TestRehydrationsGet(t *testing.T) {
	var log []rehydrationReq
	srv := rehydrationServer(t, map[string]string{
		"GET /rehydration/rh-123": `{"rehydration":{"rehydration_id":"rh-123"},"percentage":42}`,
	}, &log)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	if err := runEdx(t, "rehydrations", "get", "rh-123"); err != nil {
		t.Fatalf("get: %v", err)
	}
	if len(log) != 1 || log[0].method != http.MethodGet {
		t.Fatalf("log = %+v", log)
	}
}

func TestRehydrationsValidateRequiresQuery(t *testing.T) {
	useAPIEnv(t, "http://127.0.0.1:0")
	err := runEdx(t, "rehydrations", "validate")
	if err == nil || !strings.Contains(err.Error(), `"query"`) {
		t.Fatalf("err = %v, want required-flag error for --query", err)
	}
}

func TestRehydrationsValidate(t *testing.T) {
	var log []rehydrationReq
	srv := rehydrationServer(t, map[string]string{
		"GET /rehydration/validate": `{"total_jobs":1,"potential_rehydrations":[{"source":"s3-archive","bucket":"b","destination":"d"}]}`,
	}, &log)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	err := runEdx(t, "rehydrations", "validate",
		"--query", `service.name:"api"`,
		"--from", "2026-08-15T08:00:00.000Z", "--to", "2026-08-15T09:00:00.000Z",
		"--source", "s3-archive", "--destination", "dest-1")
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	q := mustParseQuery(t, log[0].rawQuery)
	if q.Get("query") != `service.name:"api"` ||
		q.Get("from") != "2026-08-15T08:00:00.000Z" || q.Get("to") != "2026-08-15T09:00:00.000Z" ||
		q.Get("source") != "s3-archive" || q.Get("destination") != "dest-1" {
		t.Errorf("query params = %q", log[0].rawQuery)
	}
}

func TestRehydrationsAnalyze(t *testing.T) {
	var log []rehydrationReq
	srv := rehydrationServer(t, map[string]string{
		"POST /rehydration/analyze": `{"approximate_size":1000,"approximate_count":10,"potential_overlaps":[]}`,
	}, &log)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	err := runEdx(t, "rehydrations", "analyze",
		"--query", `service.name:"api"`,
		"--from", "2026-08-15T08:00:00.000Z", "--to", "2026-08-15T09:00:00.000Z",
		"--tag", "prod", "--source", "s3-archive", "--bucket", "my-bucket")
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	if log[0].method != http.MethodPost {
		t.Fatalf("method = %q, want POST", log[0].method)
	}
	assertJSONEqual(t, log[0].body, `{
		"cql_query": "service.name:\"api\"",
		"from": "2026-08-15T08:00:00.000Z",
		"to": "2026-08-15T09:00:00.000Z",
		"tag": "prod",
		"source_integration": "s3-archive",
		"bucket": "my-bucket"
	}`)
}

func TestRehydrationsAnalyzeOmitsEmptyFields(t *testing.T) {
	var log []rehydrationReq
	srv := rehydrationServer(t, map[string]string{
		"POST /rehydration/analyze": `{"approximate_size":0,"approximate_count":0}`,
	}, &log)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	err := runEdx(t, "rehydrations", "analyze",
		"--query", `service.name:"api"`, "--tag", "prod",
		"--from", "2026-08-15T08:00:00.000Z", "--to", "2026-08-15T09:00:00.000Z")
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	// source_integration and bucket were not given, so the body must not
	// carry them at all — the API rejects an empty source_integration.
	assertJSONEqual(t, log[0].body, `{
		"cql_query": "service.name:\"api\"",
		"from": "2026-08-15T08:00:00.000Z",
		"to": "2026-08-15T09:00:00.000Z",
		"tag": "prod"
	}`)
}

func TestRehydrationsAnalyzeFromFile(t *testing.T) {
	var log []rehydrationReq
	srv := rehydrationServer(t, map[string]string{
		"POST /rehydration/analyze": `{"approximate_size":0,"approximate_count":0}`,
	}, &log)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	body := `{"from":"2026-08-15T08:00:00.000Z","to":"2026-08-15T09:00:00.000Z","archive_source":{"tag":"prod","type":"gcs_output"}}`
	file := writeTemp(t, body)
	if err := runEdx(t, "rehydrations", "analyze", "--file", file); err != nil {
		t.Fatalf("analyze --file: %v", err)
	}
	assertJSONEqual(t, log[0].body, body)
}

const validateTwoJobs = `{
	"total_jobs": 2,
	"potential_rehydrations": [
		{
			"source": "s3-archive",
			"bucket": "bucket-a",
			"destination": "dest-1",
			"efficiency_level": "fast",
			"archive_source": {"tag":"prod","node_name":"arch","type":"s3","bucket":"bucket-a"}
		},
		{
			"source": "gcs-archive",
			"bucket": "bucket-b",
			"destination": "dest-2",
			"filter_error_message": "filter not applicable"
		}
	]
}`

func TestRehydrationsCreateBatch(t *testing.T) {
	var log []rehydrationReq
	srv := rehydrationServer(t, map[string]string{
		"GET /rehydration/validate":           validateTwoJobs,
		"POST /rehydration/rehydration_batch": `[{"rehydration_id":"rh-new","status":"created"}]`,
	}, &log)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	err := runEdx(t, "rehydrations", "create",
		"--query", `service.name:"api"`,
		"--from", "2026-08-15T08:00:00.000Z", "--to", "2026-08-15T09:00:00.000Z",
		"--yes")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if len(log) != 2 {
		t.Fatalf("requests = %d, want 2 (validate + batch)", len(log))
	}
	post := log[1]
	if post.method != http.MethodPost || !strings.HasSuffix(post.path, "/rehydration/rehydration_batch") {
		t.Fatalf("second request = %s %s", post.method, post.path)
	}
	if q := mustParseQuery(t, post.rawQuery); q.Get("query") != `service.name:"api"` {
		t.Errorf("batch query param = %q", post.rawQuery)
	}
	// Only the entry without filter_error_message is submitted.
	assertJSONEqual(t, post.body, `[{
		"source_integration": "s3-archive",
		"bucket": "bucket-a",
		"destination": "dest-1",
		"archive_source": {"tag":"prod","node_name":"arch","type":"s3","bucket":"bucket-a"},
		"cql_query": "service.name:\"api\"",
		"from": "2026-08-15T08:00:00.000Z",
		"to": "2026-08-15T09:00:00.000Z",
		"exclude_overlap": true
	}]`)
}

func TestRehydrationsCreateDryRun(t *testing.T) {
	var log []rehydrationReq
	srv := rehydrationServer(t, map[string]string{
		"GET /rehydration/validate": validateTwoJobs,
	}, &log)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	err := runEdx(t, "rehydrations", "create",
		"--query", `service.name:"api"`,
		"--from", "2026-08-15T08:00:00.000Z", "--to", "2026-08-15T09:00:00.000Z",
		"--dry-run")
	if err != nil {
		t.Fatalf("create --dry-run: %v", err)
	}
	if len(log) != 1 {
		t.Fatalf("requests = %d, want 1 (validate only)", len(log))
	}
}

func TestRehydrationsCreateAbortsWithoutYes(t *testing.T) {
	var log []rehydrationReq
	srv := rehydrationServer(t, map[string]string{
		"GET /rehydration/validate": validateTwoJobs,
	}, &log)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	feedStdin(t, "n\n", func() {
		err := runEdx(t, "rehydrations", "create",
			"--query", `service.name:"api"`,
			"--from", "2026-08-15T08:00:00.000Z", "--to", "2026-08-15T09:00:00.000Z")
		if err != errAborted {
			t.Fatalf("err = %v, want errAborted", err)
		}
	})
	if len(log) != 1 {
		t.Errorf("declined create sent %d request(s), want 1 (validate only)", len(log))
	}
}

func TestRehydrationsCreateNoEligibleJobs(t *testing.T) {
	var log []rehydrationReq
	srv := rehydrationServer(t, map[string]string{
		"GET /rehydration/validate": `{"total_jobs":0,"errors":["no archive sources found"]}`,
	}, &log)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	err := runEdx(t, "rehydrations", "create",
		"--query", `service.name:"api"`,
		"--from", "2026-08-15T08:00:00.000Z", "--to", "2026-08-15T09:00:00.000Z",
		"--yes")
	if err == nil || !strings.Contains(err.Error(), "no archive sources found") {
		t.Fatalf("err = %v, want error mentioning server-side errors", err)
	}
	if len(log) != 1 {
		t.Errorf("requests = %d, want 1", len(log))
	}
}

func TestRehydrationsCancel(t *testing.T) {
	var log []rehydrationReq
	srv := rehydrationServer(t, map[string]string{
		"PUT /rehydration/rh-123": `{"rehydration_id":"rh-123","status":"cancelled"}`,
	}, &log)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	if err := runEdx(t, "rehydrations", "cancel", "rh-123", "--yes"); err != nil {
		t.Fatalf("cancel: %v", err)
	}
	if log[0].method != http.MethodPut {
		t.Fatalf("method = %q, want PUT", log[0].method)
	}
	assertJSONEqual(t, log[0].body, `{"status":"cancelled"}`)
}

func TestRehydrationsDelete(t *testing.T) {
	var log []rehydrationReq
	srv := rehydrationServer(t, map[string]string{
		"DELETE /rehydration/rh-123": `{}`,
	}, &log)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	if err := runEdx(t, "rehydrations", "delete", "rh-123", "--yes"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if log[0].method != http.MethodDelete {
		t.Fatalf("method = %q, want DELETE", log[0].method)
	}
}

func TestRehydrationsDeleteAbortsWithoutYes(t *testing.T) {
	var log []rehydrationReq
	srv := rehydrationServer(t, map[string]string{}, &log)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	feedStdin(t, "n\n", func() {
		if err := runEdx(t, "rehydrations", "delete", "rh-123"); err != errAborted {
			t.Fatalf("err = %v, want errAborted", err)
		}
	})
	if len(log) != 0 {
		t.Errorf("declined delete sent %d request(s)", len(log))
	}
}

func TestDescribePotentialRehydration(t *testing.T) {
	// Legacy integration-backed entry: top-level source/bucket are set.
	legacy := potentialRehydration{Source: "s3-archive", Bucket: "b", Destination: "d", EfficiencyLevel: "fast"}
	if got := describePotentialRehydration(legacy); got != "s3-archive (bucket b) -> d [efficiency: fast]" {
		t.Errorf("legacy = %q", got)
	}

	// Pipeline-archive entry: source/bucket are empty, fall back to archive_source.
	pipeline := potentialRehydration{
		Destination:   "d",
		ArchiveSource: []byte(`{"tag":"prod","node_name":"gcs-arch","bucket":"arch-bucket","type":"gcs_output"}`),
	}
	if got := describePotentialRehydration(pipeline); got != "prod/gcs-arch (bucket arch-bucket) -> d [efficiency: unknown]" {
		t.Errorf("pipeline = %q", got)
	}
}

func TestResolveTimeRange(t *testing.T) {
	now := time.Date(2026, 8, 15, 9, 15, 0, 0, time.UTC)

	// Explicit from/to pass through unchanged.
	from, to, err := resolveTimeRange(timeFlags{from: "2026-08-15T08:00:00.000Z", to: "2026-08-15T09:00:00.000Z"}, now)
	if err != nil || from != "2026-08-15T08:00:00.000Z" || to != "2026-08-15T09:00:00.000Z" {
		t.Errorf("explicit range = %q..%q, err %v", from, to, err)
	}

	// From without to defaults to now.
	_, to, err = resolveTimeRange(timeFlags{from: "2026-08-15T08:00:00.000Z"}, now)
	if err != nil || to != "2026-08-15T09:15:00.000Z" {
		t.Errorf("default to = %q, err %v", to, err)
	}

	// Lookback computes from/to around now.
	from, to, err = resolveTimeRange(timeFlags{lookback: "1h"}, now)
	if err != nil || from != "2026-08-15T08:15:00.000Z" || to != "2026-08-15T09:15:00.000Z" {
		t.Errorf("lookback range = %q..%q, err %v", from, to, err)
	}

	// Invalid lookback errors.
	if _, _, err = resolveTimeRange(timeFlags{lookback: "bogus"}, now); err == nil {
		t.Error("expected error for invalid lookback")
	}
}
