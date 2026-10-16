package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// metricsAPI is a path-routed fake of the metric endpoints. Each handler gets
// the request (and body) and returns status + body; missing routes 404.
type metricsAPI struct {
	routes map[string]func(r *http.Request, body string) (int, string)
	reqs   []logsRequest
}

func (m *metricsAPI) server(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		m.reqs = append(m.reqs, logsRequest{r.Method, r.URL.Path, r.URL.Query(), string(body)})
		prefix := "/v1/orgs/" + testOrg
		h, ok := m.routes[strings.TrimPrefix(r.URL.Path, prefix)]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		status, resp := h(r, string(body))
		w.WriteHeader(status)
		_, _ = w.Write([]byte(resp))
	}))
	t.Cleanup(srv.Close)
	useAPIEnv(t, srv.URL)
	return srv
}

func fixed(status int, body string) func(*http.Request, string) (int, string) {
	return func(*http.Request, string) (int, string) { return status, body }
}

func withNow(t *testing.T, now time.Time) {
	t.Helper()
	old := nowFunc
	nowFunc = func() time.Time { return now }
	t.Cleanup(func() { nowFunc = old })
}

// runEdxCapture runs edx and returns stdout, stderr and the command error.
func runEdxCapture(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	var errBuf bytes.Buffer
	restore := captureStderr(t, &errBuf)
	var err error
	out := captureStdout(t, func() { err = runEdx(t, args...) })
	restore()
	return out, errBuf.String(), err
}

var describeNow = time.Date(2026, 10, 15, 12, 0, 0, 0, time.UTC)

const describeSearchRows = `[
 {"time_unix_nano": 1760529570000000000, "name": "http.server.requests", "kind": "sum", "unit": "{request}",
  "description": "Count of HTTP requests", "resource": {"service.name": "checkout", "host.name": "web-1"},
  "attributes": {"http.route": "/pay"}, "sum": {"aggregation_temporality": "AGGREGATION_TEMPORALITY_CUMULATIVE", "is_monotonic": true, "value": 1042}},
 {"time_unix_nano": 1760529540000000000, "name": "http.server.requests", "kind": "sum", "unit": "{request}",
  "resource": {"service.name": "cart", "host.name": "web-2"}, "attributes": {"http.route": "/add"},
  "sum": {"aggregation_temporality": "AGGREGATION_TEMPORALITY_CUMULATIVE", "is_monotonic": true, "value": 77}}
]`

const describeCardinality = `{"label_values_count_total": 19, "label_names_count": 3, "total_sub_metric_count": 42,
 "labels": [
  {"label_name": "service.name", "label_kind": "resource", "label_values_count": 6, "metric_count": 42,
   "cardinality": [{"label_value": "checkout", "metric_count": 20}, {"label_value": "cart", "metric_count": 12}]},
  {"label_name": "host.name", "label_kind": "resource", "label_values_count": 3},
  {"label_name": "http.route", "label_kind": "attribute", "label_values_count": 10,
   "cardinality": [{"label_value": "/pay", "metric_count": 9}]}
 ]}`

func TestMetricsDescribe(t *testing.T) {
	withNow(t, time.Unix(0, 1760529570000000000).Add(90*time.Second))
	api := &metricsAPI{routes: map[string]func(*http.Request, string) (int, string){
		"/metrics_v2/search": fixed(200, describeSearchRows),
		"/cardinality":       fixed(200, describeCardinality),
	}}
	api.server(t)

	out, _, err := runEdxCapture(t, "metrics", "describe", "http.server.requests")
	if err != nil {
		t.Fatalf("describe: %v", err)
	}
	t.Logf("describe output (%d bytes):\n%s", len(out), out)
	var d metricDescription
	if err := json.Unmarshal([]byte(out), &d); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, out)
	}
	if !d.Found || d.Kind != "sum" || d.Unit != "{request}" || d.Temporality != "cumulative" || d.Monotonic == nil || !*d.Monotonic {
		t.Errorf("metadata = found:%v kind:%q unit:%q temporality:%q monotonic:%v", d.Found, d.Kind, d.Unit, d.Temporality, d.Monotonic)
	}
	if d.SeriesCount == nil || *d.SeriesCount != 42 {
		t.Errorf("series_count = %v, want 42", d.SeriesCount)
	}
	if d.LatestSample != "2025-10-15T11:59:30Z" || d.Age != "1m30s" {
		t.Errorf("latest_sample/age = %q/%q", d.LatestSample, d.Age)
	}
	if len(d.Labels) != 3 || d.Labels[0].Kind != "resource" || d.Labels[2].Ref != "@http.route" {
		t.Fatalf("labels = %+v", d.Labels)
	}
	if d.Labels[1].Name != "service.name" || *d.Labels[1].ValuesCount != 6 || strings.Join(d.Labels[1].Samples, ",") != "checkout,cart" {
		t.Errorf("service.name label = %+v", d.Labels[1])
	}
	ex := strings.Join(d.Examples, "\n")
	for _, want := range []string{"--agg rate --group-by service.name --summary", "--baseline 1w", "--agg increase", `--filter '@http.route:"/pay"'`} {
		if !strings.Contains(ex, want) {
			t.Errorf("examples missing %q:\n%s", want, ex)
		}
	}
	if len(d.Warnings) != 0 {
		t.Errorf("warnings = %v", d.Warnings)
	}
	// Request shapes: search by exact name, cardinality scoped to the metric.
	for _, r := range api.reqs {
		switch {
		case strings.HasSuffix(r.path, "/metrics_v2/search"):
			if q := r.query.Get("query"); q != `name:"http.server.requests"` {
				t.Errorf("search query = %q", q)
			}
		case strings.HasSuffix(r.path, "/cardinality"):
			if q := r.query.Get("query"); q != `http.server.requests{*}` {
				t.Errorf("cardinality query = %q", q)
			}
			if lb := r.query.Get("lookback"); lb != "1h" {
				t.Errorf("cardinality lookback = %q", lb)
			}
		}
	}
	// The sample examples must be runnable edx invocations.
	root := NewRootCmd()
	for _, e := range d.Examples {
		inv := edxInvocations("`" + e + "`")
		if len(inv) != 1 {
			t.Fatalf("example not parsed: %q", e)
		}
		if p := checkInvocation(root, inv[0].words); len(p) > 0 {
			t.Errorf("example %q: %v", e, p)
		}
	}
}

func TestMetricsDescribePartialFailure(t *testing.T) {
	withNow(t, describeNow)
	api := &metricsAPI{routes: map[string]func(*http.Request, string) (int, string){
		"/metrics_v2/search": fixed(200, describeSearchRows),
		"/cardinality":       fixed(500, `{"error":"clickhouse timeout"}`),
	}}
	api.server(t)

	out, _, err := runEdxCapture(t, "metrics", "describe", "http.server.requests")
	if err != nil {
		t.Fatalf("describe must succeed on partial failure: %v", err)
	}
	t.Logf("partial describe output (%d bytes):\n%s", len(out), out)
	var d metricDescription
	if err := json.Unmarshal([]byte(out), &d); err != nil {
		t.Fatal(err)
	}
	if !d.Found || d.Kind != "sum" {
		t.Errorf("found/kind = %v/%q", d.Found, d.Kind)
	}
	if len(d.Warnings) != 1 || !strings.Contains(d.Warnings[0], "cardinality failed") {
		t.Errorf("warnings = %v", d.Warnings)
	}
	if d.SeriesCount != nil {
		t.Errorf("series_count = %v, want null when cardinality failed", *d.SeriesCount)
	}
	// Labels fall back to the raw rows' keys, values_count unknown.
	refs := map[string]metricLabel{}
	for _, l := range d.Labels {
		refs[l.Ref] = l
	}
	if l, ok := refs["@http.route"]; !ok || l.ValuesCount != nil || strings.Join(l.Samples, ",") != "/pay,/add" {
		t.Errorf("@http.route fallback label = %+v (ok=%v)", l, ok)
	}
	if _, ok := refs["service.name"]; !ok {
		t.Errorf("service.name missing from fallback labels: %+v", d.Labels)
	}
}

func TestMetricsDescribeNotFoundSuggests(t *testing.T) {
	withNow(t, describeNow)
	api := &metricsAPI{routes: map[string]func(*http.Request, string) (int, string){
		"/metrics_v2/search": fixed(200, `[]`),
		"/cardinality":       fixed(200, `{"label_values_count_total":0,"label_names_count":0,"total_sub_metric_count":0,"labels":[]}`),
		"/facet_options": fixed(200, `{"name":"name","options":[
			{"name":"http.server.request.duration","count":10,"extraFields":{"kind":"histogram"}},
			{"name":"http.server.requests","count":10,"extraFields":{"kind":"sum"}},
			{"name":"system.cpu.utilization","count":5},
			{"name":"k8s.pod.memory.usage","count":5}]}`),
	}}
	api.server(t)

	out, _, err := runEdxCapture(t, "metrics", "describe", "http.server.request")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("err = %v, want not found", err)
	}
	t.Logf("not-found output (%d bytes):\n%s", len(out), out)
	var d metricDescription
	if err := json.Unmarshal([]byte(out), &d); err != nil {
		t.Fatal(err)
	}
	if d.Found || len(d.Suggestions) < 2 || !strings.HasPrefix(d.Suggestions[0], "http.server.request") {
		t.Errorf("found=%v suggestions=%v", d.Found, d.Suggestions)
	}
	if containsStr(d.Suggestions, "k8s.pod.memory.usage") {
		t.Errorf("unrelated name suggested: %v", d.Suggestions)
	}
}

func TestMetricsDescribeExistsOutsideWindow(t *testing.T) {
	api := &metricsAPI{routes: map[string]func(*http.Request, string) (int, string){
		"/metrics_v2/search": fixed(200, `[]`),
		"/cardinality":       fixed(500, `boom`),
		"/facet_options":     fixed(200, `{"options":[{"name":"ed.pipeline.bytes","extraFields":{"kind":"gauge"}}]}`),
	}}
	api.server(t)
	out, _, err := runEdxCapture(t, "metrics", "describe", "ed.pipeline.bytes")
	if err != nil {
		t.Fatalf("describe: %v", err)
	}
	var d metricDescription
	if err := json.Unmarshal([]byte(out), &d); err != nil {
		t.Fatal(err)
	}
	if !d.Found || d.Kind != "gauge" || len(d.Notes) == 0 || !strings.Contains(strings.Join(d.Examples, " "), "--agg avg") {
		t.Errorf("describe = %+v", d)
	}
}

// --- summary math ---

func pts(start time.Time, step time.Duration, vals ...float64) []graphPoint {
	out := make([]graphPoint, len(vals))
	for i, v := range vals {
		out[i] = graphPoint{Timestamp: start.Add(time.Duration(i) * step)}
		if !math.IsNaN(v) {
			vv := v
			out[i].Value = &vv
		}
	}
	return out
}

func TestAnomaliesMAD(t *testing.T) {
	start := describeNow
	// Noisy baseline around 100 with a 2-bucket spike and a single dip.
	series := pts(start, time.Minute, 100, 102, 98, 101, 99, 100, 180, 190, 101, 99, 100, 40, 102, 98)
	wins, total := anomaliesOf(series, 3.5, time.Minute, 3)
	if total != 2 || len(wins) != 2 {
		t.Fatalf("windows = %+v (total %d), want 2", wins, total)
	}
	if w := wins[0]; w.Dir != "high" || w.Buckets != 2 || w.Peak != 190 ||
		w.From != ts(start.Add(6*time.Minute)) || w.To != ts(start.Add(8*time.Minute)) {
		t.Errorf("spike window = %+v", w)
	}
	if w := wins[1]; w.Dir != "low" || w.Buckets != 1 || w.Peak != 40 {
		t.Errorf("dip window = %+v", w)
	}
	// A higher k suppresses the dip but keeps the big spike.
	if wins, _ := anomaliesOf(series, 30, time.Minute, 3); len(wins) != 1 || wins[0].Dir != "high" {
		t.Errorf("k=30 windows = %+v", wins)
	}
	// Constant series: no anomalies (MAD and MeanAD both 0).
	if wins, _ := anomaliesOf(pts(start, time.Minute, 5, 5, 5, 5, 5, 5), 3.5, time.Minute, 3); wins != nil {
		t.Errorf("constant series windows = %+v", wins)
	}
	// MAD == 0 but not constant: falls back to MeanAD and still flags the outlier.
	if wins, _ := anomaliesOf(pts(start, time.Minute, 5, 5, 5, 5, 5, 5, 5, 50), 3.5, time.Minute, 3); len(wins) != 1 || wins[0].Peak != 50 {
		t.Errorf("MAD=0 windows = %+v", wins)
	}
	// Too few points: no detection.
	if wins, _ := anomaliesOf(pts(start, time.Minute, 1, 100, 1), 3.5, time.Minute, 3); wins != nil {
		t.Errorf("short series windows = %+v", wins)
	}
	// A null bucket splits windows.
	if _, total := anomaliesOf(pts(start, time.Minute, 10, 11, 9, 10, 10, 90, math.NaN(), 95, 10, 11), 3.5, time.Minute, 3); total != 2 {
		t.Errorf("null-split total = %d, want 2", total)
	}
}

func TestTrend(t *testing.T) {
	idx := func(n int) []float64 {
		xs := make([]float64, n)
		for i := range xs {
			xs[i] = float64(i)
		}
		return xs
	}
	for _, tc := range []struct {
		name string
		ys   []float64
		want string
	}{
		{"rising", []float64{10, 12, 11, 14, 15, 17, 16, 19, 21, 22}, "rising"},
		{"falling", []float64{50, 47, 45, 44, 40, 38, 37, 33, 31, 30}, "falling"},
		{"noisy flat", []float64{100, 120, 90, 110, 95, 105, 115, 85, 100, 102}, "flat"},
		{"tiny drift", []float64{1000, 1000.1, 1000.2, 1000.3, 1000.4, 1000.5}, "flat"},
		{"too short", []float64{1, 5, 9}, "flat"},
	} {
		if got := trendOf(idx(len(tc.ys)), tc.ys); got != tc.want {
			t.Errorf("%s: trend = %q, want %q", tc.name, got, tc.want)
		}
	}
}

// graphJSON builds a /graph response for formula A with one series per label.
func graphJSON(start time.Time, step time.Duration, key string, series map[string][]float64) string {
	type rec struct {
		Values     []string     `json:"values"`
		Timeseries []graphPoint `json:"timeseries"`
	}
	var recs []rec
	for label, vals := range series {
		recs = append(recs, rec{Values: []string{label}, Timeseries: pts(start, step, vals...)})
	}
	b, _ := json.Marshal(map[string]any{"A": map[string]any{
		"from": ts(start), "to": ts(start.Add(time.Hour)), "window": step.String(),
		"keys": []string{key}, "records": recs,
	}})
	return string(b)
}

func TestMetricsQuerySummaryTopAndBaseline(t *testing.T) {
	withNow(t, describeNow)
	start := describeNow.Add(-time.Hour)
	cur := map[string][]float64{}
	base := map[string][]float64{}
	for i := 0; i < 12; i++ {
		vals := make([]float64, 60)
		bvals := make([]float64, 60)
		for j := range vals {
			vals[j] = float64(100*(12-i)) + float64(j%3)
			bvals[j] = float64(100*(12-i)) + float64(j%3)
		}
		name := fmt.Sprintf("svc-%02d", i)
		switch i {
		case 0: // doubled vs last week, with a spike
			for j := range vals {
				vals[j] *= 2
			}
			vals[30], vals[31] = 9000, 9500
		case 1: // thin baseline: only 2 points last week
			for j := range bvals {
				if j > 1 {
					bvals[j] = math.NaN()
				}
			}
		}
		cur[name] = vals
		if i != 2 { // svc-02 is new this week
			base[name] = bvals
		}
	}
	curJSON := graphJSON(start, time.Minute, "service.name", cur)
	baseJSON := graphJSON(start.Add(-7*24*time.Hour), time.Minute, "service.name", base)
	api := &metricsAPI{routes: map[string]func(*http.Request, string) (int, string){
		"/graph": func(r *http.Request, _ string) (int, string) {
			if strings.HasPrefix(r.URL.Query().Get("from"), "2026-10-08") {
				return http.StatusMultiStatus, baseJSON
			}
			return http.StatusMultiStatus, curJSON
		},
		"/facet_keys": fixed(200, `[{"key":"service.name"}]`),
	}}
	api.server(t)

	out, _, err := runEdxCapture(t, "metrics", "query", "--name", "http.server.requests", "--agg", "sum",
		"--group-by", "service.name", "--summary", "--baseline", "1w", "--top", "5")
	if err != nil {
		t.Fatalf("query --summary: %v", err)
	}
	t.Logf("summary output (%d bytes):\n%s", len(out), out)

	var graphs []logsRequest
	for _, r := range api.reqs {
		if strings.HasSuffix(r.path, "/graph") {
			graphs = append(graphs, r)
		}
	}
	if len(graphs) != 2 {
		t.Fatalf("graph requests = %d, want 2", len(graphs))
	}
	if f, to := graphs[0].query.Get("from"), graphs[0].query.Get("to"); f != "2026-10-15T11:00:00.000Z" || to != "2026-10-15T12:00:00.000Z" {
		t.Errorf("current window = %s..%s", f, to)
	}
	if f, to := graphs[1].query.Get("from"), graphs[1].query.Get("to"); f != "2026-10-08T11:00:00.000Z" || to != "2026-10-08T12:00:00.000Z" {
		t.Errorf("baseline window = %s..%s", f, to)
	}
	if graphs[0].body != graphs[1].body || !strings.Contains(graphs[0].body, `sum:http.server.requests{*} by {service.name}`) {
		t.Errorf("bodies differ or wrong CQL: %s / %s", graphs[0].body, graphs[1].body)
	}

	var sm metricSummary
	if err := json.Unmarshal([]byte(out), &sm); err != nil {
		t.Fatalf("summary is not JSON: %v", err)
	}
	if sm.SeriesCount != 12 || sm.Shown != 5 || sm.Omitted != 7 || len(sm.Series) != 5 {
		t.Errorf("counts = series %d shown %d omitted %d len %d", sm.SeriesCount, sm.Shown, sm.Omitted, len(sm.Series))
	}
	if sm.ValueStat != "total" || sm.RankedBy != "share" || sm.Rollup != "1m0s" || len(sm.TopContributors) != 3 {
		t.Errorf("header = %+v", sm)
	}
	s0 := sm.Series[0]
	if s0.Series != "service.name=svc-00" || s0.SharePct == nil || *s0.SharePct < 20 {
		t.Errorf("top series = %+v", s0)
	}
	if len(s0.Anomalies) != 1 || s0.Anomalies[0].Dir != "high" || s0.Anomalies[0].Buckets != 2 {
		t.Errorf("svc-00 anomalies = %+v", s0.Anomalies)
	}
	if b := s0.Baseline; b == nil || b.Verdict != "higher" || b.DeltaPct == nil || *b.DeltaPct < 100 || b.Points != 60 {
		t.Errorf("svc-00 baseline = %+v", b)
	}
	byName := map[string]seriesSummary{}
	for _, s := range sm.Series {
		byName[s.Series] = s
	}
	if b := byName["service.name=svc-01"].Baseline; b == nil || b.Verdict != "insufficient_data" || b.Points != 2 {
		t.Errorf("svc-01 baseline = %+v", b)
	}
	if b := byName["service.name=svc-02"].Baseline; b == nil || b.Verdict != "missing" {
		t.Errorf("svc-02 baseline = %+v", b)
	}
	if b := byName["service.name=svc-03"].Baseline; b == nil || b.Verdict != "similar" || b.DeltaPct == nil || *b.DeltaPct != 0 {
		t.Errorf("svc-03 baseline = %+v", b)
	}
	if sm.Baseline == nil || sm.Baseline.From != "2026-10-08T11:00:00Z" || sm.Baseline.Offset != "1w" || sm.Baseline.Missing != 1 {
		t.Errorf("baseline header = %+v", sm.Baseline)
	}
	if len(out) > 8000 {
		t.Errorf("summary for 5 series is %d bytes, want < 8000", len(out))
	}
}

func TestMetricsQuerySummaryTenSeriesFitsBudget(t *testing.T) {
	withNow(t, describeNow)
	start := describeNow.Add(-time.Hour)
	cur := map[string][]float64{}
	for i := 0; i < 25; i++ {
		vals := make([]float64, 60)
		for j := range vals {
			vals[j] = 1234.5678*float64(i+1) + float64(j*j%17)
		}
		vals[10+i] = 1e6 // one anomaly window per series
		cur[fmt.Sprintf("k8s-namespace-%02d/checkout-deployment-%02d", i, i)] = vals
	}
	api := &metricsAPI{routes: map[string]func(*http.Request, string) (int, string){
		"/graph":      fixed(http.StatusMultiStatus, graphJSON(start, time.Minute, "k8s.deployment.name", cur)),
		"/facet_keys": fixed(200, `[{"key":"k8s.deployment.name"}]`),
	}}
	api.server(t)
	out, _, err := runEdxCapture(t, "metrics", "query", "--name", "k8s.container.cpu.usage", "--agg", "avg",
		"--group-by", "k8s.deployment.name", "--summary")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("10-series avg summary: %d bytes", len(out))
	var sm metricSummary
	if err := json.Unmarshal([]byte(out), &sm); err != nil {
		t.Fatal(err)
	}
	if sm.Shown != 10 || sm.Omitted != 15 || sm.RankedBy != "max" || sm.Total != nil {
		t.Errorf("header = shown %d omitted %d ranked %s total %v", sm.Shown, sm.Omitted, sm.RankedBy, sm.Total)
	}
	if len(out) > 6000 {
		t.Errorf("10-series summary is %d bytes, want well under 8000", len(out))
	}

	// The table form renders too.
	out, _, err = runEdxCapture(t, "metrics", "query", "--name", "k8s.container.cpu.usage", "--agg", "avg",
		"--group-by", "k8s.deployment.name", "--summary", "-o", "table", "--top", "3")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("table summary (%d bytes):\n%s", len(out), out)
	if !strings.Contains(out, "SERIES") || !strings.Contains(out, "series:   25 (showing 3, 22 omitted)") {
		t.Errorf("table summary:\n%s", out)
	}
}

func TestMetricsQueryBaselineRequiresSummary(t *testing.T) {
	useAPIEnv(t, "http://127.0.0.1:1")
	err := runEdx(t, "metrics", "query", "--name", "x", "--baseline", "1w")
	if err == nil || !strings.Contains(err.Error(), "--baseline requires --summary") {
		t.Errorf("err = %v", err)
	}
	err = runEdx(t, "metrics", "query", "--name", "x", "--summary", "--baseline", "lastweek")
	if err == nil || !strings.Contains(err.Error(), "invalid --baseline") {
		t.Errorf("err = %v", err)
	}
}

func TestMetricsQueryTableFlattensSeries(t *testing.T) {
	resp := `{"A":{"from":"2026-10-15T11:00:00Z","to":"2026-10-15T12:00:00Z","window":"1m0s",
	 "keys":["service.name","host.name"],
	 "records":[
	  {"values":["checkout","web-1"],"timeseries":[{"timestamp":"2026-10-15T11:00:00Z","value":1},{"timestamp":"2026-10-15T11:01:00Z","value":null},{"timestamp":"2026-10-15T11:02:00Z","value":3}],"aggregate":{"value":4}},
	  {"values":["cart","web-2"],"timeseries":[{"timestamp":"2026-10-15T11:00:00Z","value":7}],"aggregate":{"value":7}}]}}`
	api := &metricsAPI{routes: map[string]func(*http.Request, string) (int, string){
		"/graph":      fixed(http.StatusMultiStatus, resp),
		"/facet_keys": fixed(200, `[{"key":"service.name"},{"key":"host.name"}]`),
	}}
	api.server(t)

	out, _, err := runEdxCapture(t, "metrics", "query", "--name", "m", "--group-by", "service.name,host.name", "-o", "table")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("table:\n%s", out)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("table has %d lines, want header + 2 series:\n%s", len(lines), out)
	}
	if f := strings.Fields(lines[0]); strings.Join(f, " ") != "SERVICE.NAME HOST.NAME AGGREGATE POINTS LAST" {
		t.Errorf("header = %q", lines[0])
	}
	if f := strings.Fields(lines[1]); strings.Join(f, " ") != "checkout web-1 4 2 3" {
		t.Errorf("row 1 = %q", lines[1])
	}

	out, _, err = runEdxCapture(t, "metrics", "query", "--name", "m", "-o", "csv")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(out, "service.name,host.name,aggregate,points,last\ncheckout,web-1,4,2,3\n") {
		t.Errorf("csv:\n%s", out)
	}

	// JSON output stays the raw response.
	out, _, err = runEdxCapture(t, "metrics", "query", "--name", "m")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"records"`) {
		t.Errorf("json output changed:\n%s", out)
	}
}

func TestMetricsQueryPassesRateAgg(t *testing.T) {
	api := &metricsAPI{routes: map[string]func(*http.Request, string) (int, string){
		"/graph": fixed(http.StatusMultiStatus, `{"A":{"records":[]}}`),
	}}
	api.server(t)
	for _, agg := range []string{"rate", "increase", "p95"} {
		if _, _, err := runEdxCapture(t, "metrics", "query", "--name", "http.server.requests", "--agg", agg,
			"--filter", `service.name:~"check.*"`); err != nil {
			t.Fatalf("--agg %s: %v", agg, err)
		}
		last := api.reqs[len(api.reqs)-1]
		if want := agg + `:http.server.requests{service.name:~\"check.*\"}`; !strings.Contains(last.body, want) {
			t.Errorf("--agg %s body = %s, want %s", agg, last.body, want)
		}
	}
}

func TestWarnUnindexedGroupBySkipsAttributes(t *testing.T) {
	api := &metricsAPI{routes: map[string]func(*http.Request, string) (int, string){
		"/graph":      fixed(http.StatusMultiStatus, `{"A":{"records":[]}}`),
		"/facet_keys": fixed(200, `[{"key":"service.name","name":"Service"}]`),
	}}
	api.server(t)

	_, stderr, err := runEdxCapture(t, "metrics", "query", "--name", "m", "--group-by", "@http.route,service.name")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(stderr, "not an indexed") {
		t.Errorf("unexpected warning for @attr group-by:\n%s", stderr)
	}

	_, stderr, err = runEdxCapture(t, "metrics", "query", "--name", "m", "--group-by", "model")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stderr, `"model" is not an indexed`) {
		t.Errorf("missing warning for unindexed key:\n%s", stderr)
	}

	// An @-only group-by skips the facet_keys lookup entirely.
	before := len(api.reqs)
	if _, _, err := runEdxCapture(t, "metrics", "query", "--name", "m", "--group-by", "@http.route"); err != nil {
		t.Fatal(err)
	}
	for _, r := range api.reqs[before:] {
		if strings.HasSuffix(r.path, "/facet_keys") {
			t.Errorf("facet_keys fetched for @-only group-by")
		}
	}
}

func TestParseLookback(t *testing.T) {
	for in, want := range map[string]time.Duration{"1w": 7 * 24 * time.Hour, "2d": 48 * time.Hour, "90m": 90 * time.Minute} {
		if got, err := parseLookback(in); err != nil || got != want {
			t.Errorf("parseLookback(%q) = %v, %v", in, got, err)
		}
	}
	if _, err := parseLookback("xw"); err == nil {
		t.Error("parseLookback(xw) succeeded")
	}
}

func TestWriteCompactJSONRoundTrips(t *testing.T) {
	v := map[string]any{"a": 1, "list": []map[string]any{{"x": 1}, {"x": 2}}, "s": []string{"p", "q"}}
	var buf bytes.Buffer
	if err := writeCompactJSON(&buf, v); err != nil {
		t.Fatal(err)
	}
	var back map[string]any
	if err := json.Unmarshal(buf.Bytes(), &back); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, buf.String())
	}
	if !strings.Contains(buf.String(), "\n    {\"x\":1},\n    {\"x\":2}\n  ]") {
		t.Errorf("array of objects not one-per-line:\n%s", buf.String())
	}
}

func TestMetricsListSendsTimeWindow(t *testing.T) {
	facet := `{"options":[{"name":"qa_http.server.requests"},{"name":"system.cpu.usage"}]}`
	cases := []struct {
		args []string
		want map[string]string
	}{
		{[]string{"metrics", "list"}, map[string]string{"lookback": "24h"}},
		{[]string{"metrics", "list", "--lookback", "7h"}, map[string]string{"lookback": "7h"}},
		{[]string{"metrics", "list", "--from", "2026-10-15T12:00:00.000Z", "--to", "2026-10-15T18:00:00.000Z"},
			map[string]string{"from": "2026-10-15T12:00:00.000Z", "to": "2026-10-15T18:00:00.000Z", "lookback": ""}},
	}
	for _, tc := range cases {
		api := &metricsAPI{routes: map[string]func(*http.Request, string) (int, string){
			"/facet_options": fixed(200, facet),
		}}
		api.server(t)
		out, _, err := runEdxCapture(t, append(tc.args, "--keyword", "qa_")...)
		if err != nil {
			t.Fatalf("%v: %v", tc.args, err)
		}
		if !strings.Contains(out, "qa_http.server.requests") || strings.Contains(out, "system.cpu.usage") {
			t.Errorf("%v: keyword filter output = %s", tc.args, out)
		}
		if len(api.reqs) != 1 {
			t.Fatalf("%v: %d requests", tc.args, len(api.reqs))
		}
		q := api.reqs[0].query
		for k, v := range tc.want {
			if got := q.Get(k); got != v {
				t.Errorf("%v: %s = %q, want %q", tc.args, k, got, v)
			}
		}
	}
}
