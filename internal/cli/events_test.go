package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// eventsAPIServer serves /events/search, capturing each request's query and
// replying from pages keyed by the incoming cursor ("" for the first page).
func eventsAPIServer(t *testing.T, gotQueries *[]url.Values, pages map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		*gotQueries = append(*gotQueries, q)
		resp, ok := pages[q.Get("cursor")]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, `{"error":"unknown cursor %q"}`, q.Get("cursor"))
			return
		}
		_, _ = w.Write([]byte(resp))
	}))
}

func eventItems(names ...string) string {
	out := ""
	for i, n := range names {
		if i > 0 {
			out += ","
		}
		out += `{"timestamp":1,"body":"` + n + `"}`
	}
	return out
}

func TestEventsSearchSendsQueryParams(t *testing.T) {
	var got []url.Values
	srv := eventsAPIServer(t, &got, map[string]string{
		"": `{"query_id":"x","items":[` + eventItems("a") + `],"next_cursor":""}`,
	})
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	if err := runEdx(t, "events", "search",
		"--query", `event.domain:"Monitor Alerts"`, "--lookback", "720h", "--limit", "50"); err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("requests = %d, want 1", len(got))
	}
	for k, want := range map[string]string{
		"query":    `event.domain:"Monitor Alerts"`,
		"lookback": "720h",
		"limit":    "50",
		"order":    "desc",
	} {
		if v := got[0].Get(k); v != want {
			t.Errorf("query[%s] = %q, want %q", k, v, want)
		}
	}
}

func TestEventsSearchAllFollowsCursorToCompletion(t *testing.T) {
	var got []url.Values
	srv := eventsAPIServer(t, &got, map[string]string{
		"":   `{"query_id":"x","items":[` + eventItems("a", "b") + `],"next_cursor":"c1"}`,
		"c1": `{"query_id":"x","items":[` + eventItems("c", "d") + `],"next_cursor":"c2"}`,
		"c2": `{"query_id":"x","items":[` + eventItems("e") + `],"next_cursor":""}`,
	})
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	stdout := captureStdout(t, func() {
		if err := runEdx(t, "events", "search", "--query", "x", "--all"); err != nil {
			t.Errorf("search --all: %v", err)
		}
	})

	if len(got) != 3 {
		t.Fatalf("requests = %d, want 3 (one per page)", len(got))
	}
	// --all without an explicit --limit uses the server's default page size,
	// not the CLI's one-page default of 20.
	if v := got[0].Get("limit"); v != "1000" {
		t.Errorf("page 1 limit = %q, want 1000", v)
	}
	if v := got[1].Get("cursor"); v != "c1" {
		t.Errorf("page 2 cursor = %q, want c1", v)
	}
	if v := got[2].Get("cursor"); v != "c2" {
		t.Errorf("page 3 cursor = %q, want c2", v)
	}

	var merged struct {
		Items      []map[string]any `json:"items"`
		Pages      int              `json:"pages"`
		TotalItems int              `json:"total_items"`
		NextCursor string           `json:"next_cursor"`
	}
	if err := json.Unmarshal([]byte(stdout), &merged); err != nil {
		t.Fatalf("merged output is not valid JSON: %v\n%s", err, stdout)
	}
	if merged.Pages != 3 || merged.TotalItems != 5 || len(merged.Items) != 5 {
		t.Errorf("merged = %d pages / %d total / %d items, want 3/5/5",
			merged.Pages, merged.TotalItems, len(merged.Items))
	}
	if merged.NextCursor != "" {
		t.Errorf("merged next_cursor = %q, want empty (sweep is complete)", merged.NextCursor)
	}
	if last, _ := merged.Items[4]["body"].(string); last != "e" {
		t.Errorf("last item body = %q, want %q", last, "e")
	}
}

func TestEventsSearchAllRespectsExplicitLimit(t *testing.T) {
	var got []url.Values
	srv := eventsAPIServer(t, &got, map[string]string{
		"": `{"query_id":"x","items":[` + eventItems("a") + `],"next_cursor":""}`,
	})
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	if err := runEdx(t, "events", "search", "--all", "--limit", "200"); err != nil {
		t.Fatalf("search --all --limit: %v", err)
	}
	if v := got[0].Get("limit"); v != "200" {
		t.Errorf("limit = %q, want 200", v)
	}
}

func TestEventsSearchAllAbortsOnRepeatedCursor(t *testing.T) {
	var got []url.Values
	srv := eventsAPIServer(t, &got, map[string]string{
		"":   `{"query_id":"x","items":[],"next_cursor":"c1"}`,
		"c1": `{"query_id":"x","items":[],"next_cursor":"c1"}`,
	})
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	if err := runEdx(t, "events", "search", "--all"); err == nil {
		t.Fatal("expected an error when the server repeats a cursor")
	}
	if len(got) > 2 {
		t.Errorf("requests = %d, want the loop to stop at 2", len(got))
	}
}

// quickRetries shrinks the --all retry backoff so failure tests stay fast.
func quickRetries(t *testing.T) {
	t.Helper()
	old := allRetryBaseDelay
	allRetryBaseDelay = time.Millisecond
	t.Cleanup(func() { allRetryBaseDelay = old })
}

func TestEventsSearchAllMidSweepErrorEmitsNothing(t *testing.T) {
	quickRetries(t)
	var got []url.Values
	srv := eventsAPIServer(t, &got, map[string]string{
		"": `{"query_id":"x","items":[` + eventItems("a") + `],"next_cursor":"gone"}`,
		// cursor "gone" is not in the map -> 400 from the mock server.
	})
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	stdout := captureStdout(t, func() {
		err := runEdx(t, "events", "search", "--all")
		if err == nil {
			t.Error("expected an error when a later page fails")
		} else if !strings.Contains(err.Error(), "--all --cursor 'gone'") {
			t.Errorf("error does not name the resume cursor: %v", err)
		}
	})
	// A partial sweep must never be printed: broken-looking-complete JSON is
	// exactly the failure mode --all exists to prevent.
	if stdout != "" {
		t.Errorf("stdout not empty after mid-sweep failure: %q", stdout)
	}
}

// A page that fails transiently (one 500) must be retried, not abort the
// sweep: this is what a live 107-page staging sweep dies on otherwise.
func TestEventsSearchAllRetriesTransientPageFailure(t *testing.T) {
	quickRetries(t)
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&hits, 1)
		switch {
		case r.URL.Query().Get("cursor") == "":
			_, _ = w.Write([]byte(`{"query_id":"x","items":[` + eventItems("a") + `],"next_cursor":"c1"}`))
		case n == 2: // first attempt at page 2 fails like staging under load
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"Error":"Failed to query event search"}`))
		default:
			_, _ = w.Write([]byte(`{"query_id":"x","items":[` + eventItems("b") + `],"next_cursor":""}`))
		}
	}))
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	stdout := captureStdout(t, func() {
		if err := runEdx(t, "events", "search", "--all"); err != nil {
			t.Errorf("search --all with one transient 500: %v", err)
		}
	})
	var merged struct {
		TotalItems int `json:"total_items"`
		Pages      int `json:"pages"`
	}
	if err := json.Unmarshal([]byte(stdout), &merged); err != nil {
		t.Fatalf("merged output is not valid JSON: %v", err)
	}
	if merged.Pages != 2 || merged.TotalItems != 2 {
		t.Errorf("merged = %d pages / %d items, want 2/2", merged.Pages, merged.TotalItems)
	}
	if hits != 3 {
		t.Errorf("requests = %d, want 3 (page1, failed page2, retried page2)", hits)
	}
}

func TestEventsSearchSinglePagePassesCursorThrough(t *testing.T) {
	var got []url.Values
	srv := eventsAPIServer(t, &got, map[string]string{
		"resume": `{"query_id":"x","items":[` + eventItems("z") + `],"next_cursor":""}`,
	})
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	if err := runEdx(t, "events", "search", "--cursor", "resume"); err != nil {
		t.Fatalf("search --cursor: %v", err)
	}
	if v := got[0].Get("cursor"); v != "resume" {
		t.Errorf("cursor = %q, want %q", v, "resume")
	}
}
