package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/edgedelta/edx/internal/api"
)

// fakeGetter returns a canned dashboard body per call, so the wait loop can be driven
// through a sequence of states.
type fakeGetter struct {
	states []*screenshotState
	calls  int
	err    error
}

func (f *fakeGetter) Get(ctx context.Context, path string, query url.Values) ([]byte, error) {
	if f.err != nil {
		return nil, f.err
	}
	i := f.calls
	f.calls++
	if i >= len(f.states) {
		i = len(f.states) - 1
	}
	return json.Marshal(map[string]any{
		"dashboard_id": "d1",
		"screenshot":   f.states[i],
	})
}

func pending(triggered string) *screenshotState {
	return &screenshotState{Status: screenshotPending, LastTriggeredAt: triggered}
}

func done(triggered, light string) *screenshotState {
	return &screenshotState{Status: screenshotOK, LastTriggeredAt: triggered, LightURL: light}
}

// The result is written over one field on the dashboard record, so "status is terminal" is
// not enough to prove the render we asked for has finished. These are the cases that decide
// whether edx hands back a fresh image or the previous one.
func TestScreenshotSettled(t *testing.T) {
	for name, tc := range map[string]struct {
		before, now *screenshotState
		want        bool
	}{
		"no screenshot yet":              {nil, nil, false},
		"first render still pending":     {nil, pending("t1"), false},
		"first render finished":          {nil, done("t1", "a.png"), true},
		"previous result, not restamped": {done("t1", "a.png"), done("t1", "a.png"), false},
		"previous result, restamped":     {done("t1", "a.png"), done("t2", "b.png"), true},
		"our render still pending":       {done("t1", "a.png"), pending("t2"), false},
		"attached to a run in flight":    {pending("t1"), done("t1", "a.png"), true},
		"errored render is settled":      {done("t1", "a.png"), &screenshotState{Status: screenshotError, LastTriggeredAt: "t2"}, true},
		"empty status is not a status":   {nil, &screenshotState{}, false},
	} {
		t.Run(name, func(t *testing.T) {
			if got := screenshotSettled(tc.before, tc.now); got != tc.want {
				t.Errorf("screenshotSettled(%+v, %+v) = %v, want %v", tc.before, tc.now, got, tc.want)
			}
		})
	}
}

func TestAwaitScreenshotWaitsForAFreshResult(t *testing.T) {
	before := done("t1", "old.png")
	g := &fakeGetter{states: []*screenshotState{
		before,        // the stale read that a naive loop would accept
		pending("t2"), // our render, now visible
		done("t2", "new.png"),
	}}

	state, err := awaitScreenshot(context.Background(), g, "d1", before, time.Second, time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.LightURL != "new.png" {
		t.Errorf("light_url = %q, want the fresh render", state.LightURL)
	}
	if g.calls != 3 {
		t.Errorf("polled %d times, want 3", g.calls)
	}
}

// A pending status that never resolves is the shape of a lost render, and the message has to
// say so rather than blame the caller's timeout.
func TestAwaitScreenshotTimesOutOnAStuckRender(t *testing.T) {
	g := &fakeGetter{states: []*screenshotState{pending("t2")}}

	_, err := awaitScreenshot(context.Background(), g, "d1", nil, 5*time.Millisecond, time.Millisecond)
	if err == nil {
		t.Fatal("expected a timeout error")
	}
	if !strings.Contains(err.Error(), screenshotPending) {
		t.Errorf("error %q does not report the last status", err)
	}
}

func TestAwaitScreenshotPropagatesAPIErrors(t *testing.T) {
	g := &fakeGetter{err: errors.New("boom")}
	if _, err := awaitScreenshot(context.Background(), g, "d1", nil, time.Second, time.Millisecond); err == nil {
		t.Error("expected the API error to surface")
	}
}

func TestAwaitScreenshotRespectsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	g := &fakeGetter{states: []*screenshotState{pending("t2")}}

	if _, err := awaitScreenshot(ctx, g, "d1", nil, time.Minute, time.Millisecond); !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
}

// A dashboard that has never been rendered has no screenshot object at all.
func TestFetchScreenshotStateHandlesAnAbsentScreenshot(t *testing.T) {
	g := &fakeGetter{states: []*screenshotState{nil}}
	state, err := fetchScreenshotState(context.Background(), g, "d1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state != nil {
		t.Errorf("state = %+v, want nil", state)
	}
	if statusOf(state) != "none" {
		t.Errorf("statusOf(nil) = %q, want none", statusOf(state))
	}
}

// The image lives on a storage host that has nothing to do with the API, so the org's
// credentials must not travel with the request.
func TestDownloadScreenshotSendsNoCredentials(t *testing.T) {
	var gotAuth, gotCookie string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotCookie = r.Header.Get("Cookie")
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write([]byte("\x89PNG\r\n\x1a\n"))
	}))
	defer srv.Close()

	image, err := downloadScreenshot(context.Background(), srv.Client(), srv.URL+"/shot.png")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(string(image), "\x89PNG") {
		t.Errorf("got %q, want the PNG bytes", image)
	}
	if gotAuth != "" || gotCookie != "" {
		t.Errorf("request carried credentials: Authorization=%q Cookie=%q", gotAuth, gotCookie)
	}
}

func TestDownloadScreenshotReportsHTTPFailures(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "<Error>AccessDenied</Error>")
	}))
	defer srv.Close()

	_, err := downloadScreenshot(context.Background(), srv.Client(), srv.URL+"/shot.png")
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "403") || !strings.Contains(err.Error(), "AccessDenied") {
		t.Errorf("error %q should carry the status and the body", err)
	}
}

// An empty object would be written out as a 0-byte PNG and read as a blank dashboard, which
// is a worse answer than an error.
func TestDownloadScreenshotRejectsAnEmptyBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	if _, err := downloadScreenshot(context.Background(), srv.Client(), srv.URL+"/shot.png"); err == nil {
		t.Error("expected an error for an empty body")
	}
}

func TestScreenshotFilename(t *testing.T) {
	for name, tc := range map[string]struct {
		id, theme, format string
		want              string
	}{
		"plain":                 {"abc123", "light", "png", "abc123.png"},
		"dark is distinguished": {"abc123", "dark", "png", "abc123-dark.png"},
		"pdf":                   {"abc123", "light", "pdf", "abc123.pdf"},
		// A dashboard ID is server-assigned, but it lands in a path, so it must not be
		// able to steer the write somewhere else.
		"no path traversal": {"../../etc/passwd", "light", "png", "passwd.png"},
		"no separators":     {"a/b\\c", "light", "png", "c.png"},
		"only traversal":    {"../..", "light", "png", "dashboard.png"},
		"empty":             {"  ", "light", "png", "dashboard.png"},
	} {
		t.Run(name, func(t *testing.T) {
			if got := screenshotFilename(tc.id, tc.theme, tc.format); got != tc.want {
				t.Errorf("screenshotFilename(%q, %q, %q) = %q, want %q", tc.id, tc.theme, tc.format, got, tc.want)
			}
		})
	}
}

func TestParseFacetParams(t *testing.T) {
	got, err := parseFacetParams([]string{"env=prod", "service=api", "empty="})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]string{"env": "prod", "service": "api", "empty": ""}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s = %q, want %q", k, got[k], v)
		}
	}

	if _, err := parseFacetParams([]string{"nope"}); err == nil {
		t.Error("expected an error for a value without =")
	}
	if _, err := parseFacetParams([]string{"=prod"}); err == nil {
		t.Error("expected an error for an empty key")
	}
	if got, err := parseFacetParams(nil); err != nil || got != nil {
		t.Errorf("parseFacetParams(nil) = %v, %v, want nil, nil", got, err)
	}
}

// A refused trigger has to be told apart from a real failure: it means a render is already
// running, which is a reason to wait rather than to give up.
func TestIsScreenshotInProgress(t *testing.T) {
	inProgress := &api.Error{Status: http.StatusBadRequest, Body: "Not triggered screenshot because it is already in progress"}
	if !isScreenshotInProgress(inProgress) {
		t.Error("the in-progress response was not recognised")
	}
	if isScreenshotInProgress(&api.Error{Status: http.StatusBadRequest, Body: "org id parameter is required"}) {
		t.Error("an unrelated 400 was treated as an in-flight render")
	}
	if isScreenshotInProgress(&api.Error{Status: http.StatusNotFound, Body: "already in progress"}) {
		t.Error("a 404 was treated as an in-flight render")
	}
	if isScreenshotInProgress(errors.New("boom")) {
		t.Error("a non-API error was treated as an in-flight render")
	}
}

// Exports and previews live in separate fields on the record, so reading the wrong one
// would hand back a stale image from an unrelated render rather than fail visibly.
func TestScreenshotStateURLFor(t *testing.T) {
	s := &screenshotState{
		LightURL: "preview-l.png", DarkURL: "preview-d.png",
		FullLightURL: "full-l.png", FullDarkURL: "full-d.png",
	}
	for _, tc := range []struct {
		theme  string
		export bool
		want   string
	}{
		{"light", false, "preview-l.png"},
		{"dark", false, "preview-d.png"},
		{"light", true, "full-l.png"},
		{"dark", true, "full-d.png"},
	} {
		if got := s.urlFor(tc.theme, tc.export); got != tc.want {
			t.Errorf("urlFor(%q, export=%v) = %q, want %q", tc.theme, tc.export, got, tc.want)
		}
	}

	// PDF renders produce no dark variant, so this must be visibly empty rather than
	// falling back to the light one.
	pdf := &screenshotState{FullLightURL: "report.pdf"}
	if pdf.urlFor("dark", true) != "" {
		t.Errorf("urlFor(dark) = %q, want empty", pdf.urlFor("dark", true))
	}
	var nilState *screenshotState
	if nilState.urlFor("light", true) != "" {
		t.Error("urlFor on a nil state should be empty")
	}
}

// The format can be stated three ways, and a PDF written into a file called .png is a
// mistake nobody catches until they try to open it.
func TestResolveFormat(t *testing.T) {
	for name, tc := range map[string]struct {
		calledAs   string
		format     string
		formatSet  bool
		out        string
		want       string
		wantErrHas string
	}{
		"default":                    {"screenshot", "png", false, "", "png", ""},
		"explicit flag":              {"screenshot", "pdf", true, "", "pdf", ""},
		"pdf alias":                  {"pdf", "png", false, "", "pdf", ""},
		"png alias":                  {"png", "png", false, "", "png", ""},
		"other aliases stay neutral": {"render", "png", false, "", "png", ""},
		"from the out extension":     {"screenshot", "png", false, "report.pdf", "pdf", ""},
		"extension is case blind":    {"screenshot", "png", false, "REPORT.PDF", "pdf", ""},
		"unrelated extension":        {"screenshot", "png", false, "shot.jpeg", "png", ""},
		"stdout has no extension":    {"screenshot", "pdf", true, "-", "pdf", ""},
		"sources agreeing is fine":   {"pdf", "pdf", true, "report.pdf", "pdf", ""},

		"alias contradicts the extension": {"pdf", "png", false, "shot.png", "", "disagree"},
		"flag contradicts the alias":      {"pdf", "png", true, "", "", "disagree"},
		"flag contradicts the extension":  {"screenshot", "pdf", true, "shot.png", "", "disagree"},
		"invalid format":                  {"screenshot", "gif", true, "", "", "invalid --format"},
	} {
		t.Run(name, func(t *testing.T) {
			got, err := resolveFormat(tc.calledAs, tc.format, tc.formatSet, tc.out)
			if tc.wantErrHas != "" {
				if err == nil {
					t.Fatalf("resolveFormat(...) = %q, want an error containing %q", got, tc.wantErrHas)
				}
				if !strings.Contains(err.Error(), tc.wantErrHas) {
					t.Errorf("error %q does not contain %q", err, tc.wantErrHas)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("resolveFormat(%q, %q, %v, %q) = %q, want %q",
					tc.calledAs, tc.format, tc.formatSet, tc.out, got, tc.want)
			}
		})
	}
}

// The aliases are the whole point of the shorthand, so they have to actually route.
func TestScreenshotAliasesResolve(t *testing.T) {
	for _, alias := range []string{"screenshot", "pdf", "png", "export", "render"} {
		root := NewRootCmd()
		cmd, _, err := root.Find([]string{"dashboards", alias})
		if err != nil {
			t.Errorf("%q did not resolve: %v", alias, err)
			continue
		}
		if cmd.Name() != "screenshot" {
			t.Errorf("%q resolved to %q, want screenshot", alias, cmd.Name())
		}
	}
}

// A timeout says nothing about the dashboard, so a caller running this unattended has to be
// able to tell it apart from a render that actually failed — without matching on text.
func TestRenderPendingErrorIsDistinguishable(t *testing.T) {
	g := &fakeGetter{states: []*screenshotState{pending("t2")}}
	_, err := awaitScreenshot(context.Background(), g, "d1", nil, 5*time.Millisecond, time.Millisecond)
	if err == nil {
		t.Fatal("expected a timeout error")
	}

	var stillRendering *renderPendingError
	if !errors.As(err, &stillRendering) {
		t.Fatalf("err is %T, want *renderPendingError", err)
	}
	if got := exitCodeFor(err); got != exitRenderPending {
		t.Errorf("exit code = %d, want %d", got, exitRenderPending)
	}
	if statusOf(stillRendering.state) != screenshotPending {
		t.Errorf("the error dropped the last status seen: %+v", stillRendering.state)
	}

	// A render that finished and failed is an ordinary failure, and must not borrow the
	// retry-me exit code.
	if got := exitCodeFor(errors.New("the dashboard did not finish rendering")); got != 1 {
		t.Errorf("a plain error exits %d, want 1", got)
	}
}
