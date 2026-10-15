package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/edgedelta/edx/internal/api"
	"github.com/edgedelta/edx/internal/output"
)

func newMetricsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Discover and query metrics",
		Long: `Discover metric names and run aggregation queries.

Metric names must be EXACT (no wildcards). Agent workflow:

  1. edx metrics list --keyword <part>     find the exact name
  2. edx metrics describe <name>           kind, unit, labels (+ refs), cardinality,
                                           staleness and example queries
  3. edx metrics query ... --summary       per-series stats, trend, anomalies
     (add --baseline 1w to compare with the same window last week)
  4. edx metrics query ...                 raw points, only when you need them`,
	}
	cmd.AddCommand(newMetricsListCmd(), newMetricsDescribeCmd(), newMetricsQueryCmd())
	return cmd
}

func newMetricsListCmd() *cobra.Command {
	var keyword string
	var limit int
	var tf timeFlags
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List metric names that reported data in a time window",
		Long: `List metric names that reported at least one sample in the time window
(default: the last 24h). A metric that stopped reporting before the window
starts is not listed, so widen --lookback or pass --from/--to when you
investigate an older incident.`,
		Example: `  edx metrics list
  edx metrics list --keyword cpu
  edx metrics list --keyword http --from 2026-10-15T12:00:00.000Z --to 2026-10-15T18:00:00.000Z`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			q.Set("scope", "metric")
			q.Set("facet_path", "name")
			q.Set("limit", itoa(limit))
			tf.apply(q)
			data, err := c.Get(cmdContext(cmd), "/facet_options", q)
			if err != nil {
				return err
			}
			if keyword == "" {
				return printResult(data)
			}
			// Client-side keyword filter over the facet options.
			var facet struct {
				Options []map[string]any `json:"options"`
			}
			if err := json.Unmarshal(data, &facet); err != nil {
				return printResult(data)
			}
			filtered := make([]map[string]any, 0)
			kw := strings.ToLower(keyword)
			for _, opt := range facet.Options {
				if name, _ := opt["name"].(string); strings.Contains(strings.ToLower(name), kw) {
					filtered = append(filtered, opt)
					continue
				}
				if value, _ := opt["value"].(string); strings.Contains(strings.ToLower(value), kw) {
					filtered = append(filtered, opt)
				}
			}
			out, err := json.Marshal(map[string]any{"options": filtered})
			if err != nil {
				return err
			}
			return printResult(out)
		},
	}
	cmd.Flags().StringVar(&keyword, "keyword", "", "case-insensitive substring filter on metric names")
	cmd.Flags().IntVar(&limit, "limit", 1000, "maximum number of metric names to fetch")
	tf.register(cmd, "24h")
	return cmd
}

func newMetricsQueryCmd() *cobra.Command {
	var name, agg, filter, graphType, baseline string
	var groupBy []string
	var rollup, top int
	var summary bool
	var anomalyK float64
	var tf timeFlags
	var pg pageFlags
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Aggregate a metric as a timeseries or table",
		Long: `Aggregate a metric. The CQL built under the hood is:

  <agg>:<name>{<filter>} by {<group-by>}.rollup(<seconds>)

The filter uses CQL field syntax (service.name:"api", regex service.name:~"api-.*");
full-text search is not supported for metrics. Use "*" for no filter. Datapoint
attributes are referenced with an "@" prefix (@http.route); run
"edx metrics describe <name>" to see every label's exact ref.

Only indexed metric dimensions and "@" datapoint attributes can be grouped on.
Grouping by an unindexed custom dimension returns a single empty group rather
than an error — edx warns when a --group-by key is not indexed.

Aggregations: sum, avg, min, max, count, median, first, p25/p50/p75/p90/p95/p99
(percentiles over histograms and plain series), rate and increase (cumulative
counters: per-second rate / per-bucket increase). They are passed through to
the API as-is.

Output: by default the raw response keyed by query id ({"A": {"records": [...]}});
-o table/csv prints one row per series (labels, aggregate, points).

--summary (recommended for agents) replaces raw points with per-series
statistics: value (total for sum/count/increase, else avg), share of total,
min/avg/max/p95/first/last, change, trend (rising/falling/flat by slope
significance), gaps, and anomaly windows (runs of buckets whose modified
z-score 0.6745*|x-median|/MAD exceeds --anomaly-k). Series are ranked by share
of total (sum-like aggs) or by max, and --top keeps the first N.

--baseline <offset> (1w = same window last week, 1d = yesterday) runs the same
query over the shifted window and adds baseline value, delta, delta % and a
verdict per series (higher/lower/similar, or insufficient_data when either
window has too few points — never read a ratio off a thin baseline). Prefer 1w:
it survives daily and weekly cycles.`,
		Example: `  edx metrics query --name http.request.duration --agg avg --group-by service.name
  edx metrics query --name system.cpu.usage --agg max --filter 'host.name:"web-1"' --lookback 24h
  edx metrics query --name http.requests --agg sum --rollup 300 --graph-type table
  edx metrics query --name http.server.requests --agg rate --group-by service.name --summary
  edx metrics query --name http.server.requests --agg sum --filter 'service.name:~"checkout.*"' --summary --baseline 1w --top 5`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if top < 0 {
				return fmt.Errorf("--top must be >= 0")
			}
			if anomalyK <= 0 {
				return fmt.Errorf("--anomaly-k must be > 0")
			}
			var baseOffset time.Duration
			if baseline != "" {
				if !summary {
					return fmt.Errorf("--baseline requires --summary")
				}
				d, err := parseLookback(baseline)
				if err != nil || d <= 0 {
					return fmt.Errorf("invalid --baseline %q: use a duration such as 1w, 1d or 6h", baseline)
				}
				baseOffset = d
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			ctx := cmdContext(cmd)
			cql := buildMetricCQL(agg, name, filter, groupBy, rollup)
			if len(groupBy) > 0 {
				warnUnindexedGroupBy(ctx, c, groupBy)
			}
			body, err := metricGraphBody(cql)
			if err != nil {
				return err
			}
			q := url.Values{}
			var from, to time.Time
			if baseOffset > 0 {
				// Pin both windows to absolute bounds so the baseline is an
				// exact shift of the window under test.
				from, to, err = resolveWindow(tf)
				if err != nil {
					return err
				}
				q.Set("from", isoMillis(from))
				q.Set("to", isoMillis(to))
			} else {
				tf.apply(q)
			}
			pg.apply(q)
			if summary {
				q.Set("graph_type", "timeseries")
			} else if graphType != "" {
				q.Set("graph_type", graphType)
			}
			data, err := c.Post(ctx, "/graph", q, body)
			if err != nil {
				return err
			}
			if !summary {
				if flagOutput == "table" || flagOutput == "csv" {
					if rows, cols, ok := graphRows(data); ok {
						return printResultWith(rows, output.Options{Format: flagOutput, Columns: defaultCols(cols)})
					}
				}
				return printResult(data)
			}

			res, err := parseGraph(data)
			if err != nil {
				return err
			}
			f, ok := res["A"]
			if !ok {
				return fmt.Errorf("no result for query A in /graph response")
			}
			if f.Error != "" {
				return fmt.Errorf("query failed: %s", f.Error)
			}
			opts := summaryOptions{Agg: agg, Query: cql, Rollup: rollup, Top: top, AnomalyK: anomalyK, MaxWindow: 3}
			sm := summarizeGraph(f, opts)
			if baseOffset > 0 {
				bFrom, bTo := from.Add(-baseOffset), to.Add(-baseOffset)
				bq := url.Values{}
				bq.Set("from", isoMillis(bFrom))
				bq.Set("to", isoMillis(bTo))
				pg.apply(bq)
				bq.Set("graph_type", "timeseries")
				bdata, err := c.Post(ctx, "/graph", bq, body)
				var bf graphFormula
				if err == nil {
					var bres map[string]graphFormula
					if bres, err = parseGraph(bdata); err == nil {
						bf = bres["A"]
						if bf.Error != "" {
							err = fmt.Errorf("%s", bf.Error)
						}
					}
				}
				if err != nil {
					sm.Notes = append(sm.Notes, "baseline query failed: "+oneLine(err))
				} else {
					applyBaseline(sm, bf, opts, baseline, bFrom, bTo)
				}
			}
			return printSummary(sm)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "exact metric name (required; find it with \"edx metrics list\", inspect it with \"edx metrics describe\")")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&agg, "agg", "sum", "aggregation: sum, avg, min, max, count, median, first, p25/p50/p75/p90/p95/p99, rate, increase")
	cmd.Flags().StringVar(&filter, "filter", "*", `CQL filter (e.g. 'service.name:"api"', regex 'service.name:~"api-.*"', attribute '@http.route:"/login"'); "*" for none`)
	cmd.Flags().StringSliceVar(&groupBy, "group-by", nil, "group-by keys (e.g. service.name,@http.route)")
	cmd.Flags().IntVar(&rollup, "rollup", 0, "rollup period in seconds (default: auto from time range)")
	cmd.Flags().StringVar(&graphType, "graph-type", "timeseries", "graph type: timeseries or table (ignored with --summary)")
	cmd.Flags().BoolVar(&summary, "summary", false, "per-series stats, trend and anomaly windows instead of raw points (compact; best for agents)")
	cmd.Flags().StringVar(&baseline, "baseline", "", "with --summary: compare against the window shifted back by this offset (1w recommended, 1d, 6h)")
	cmd.Flags().IntVar(&top, "top", 10, "with --summary: keep the top N series (0 = all); the rest are counted as omitted")
	cmd.Flags().Float64Var(&anomalyK, "anomaly-k", 3.5, "with --summary: modified z-score threshold for anomaly windows")
	tf.register(cmd, "1h")
	pg.register(cmd, 0)
	return cmd
}

func buildMetricCQL(agg, name, filter string, groupBy []string, rollup int) string {
	if strings.TrimSpace(filter) == "" {
		filter = "*"
	}
	cql := fmt.Sprintf("%s:%s{%s}", agg, name, filter)
	if len(groupBy) > 0 {
		cql += fmt.Sprintf(" by {%s}", strings.Join(groupBy, ","))
	}
	if rollup > 0 {
		cql += fmt.Sprintf(".rollup(%d)", rollup)
	}
	return cql
}

func metricGraphBody(cql string) ([]byte, error) {
	return json.Marshal(map[string]any{
		"queries":  map[string]any{"A": map[string]any{"scope": "metric", "query": cql}},
		"formulas": map[string]any{"A": map[string]any{"formula": "A"}},
	})
}

// graphRows flattens a /graph response into one row per series — the group
// label columns, the window aggregate and the point count — for table/csv
// output. It returns the default column order alongside the rows.
func graphRows(data []byte) ([]byte, []string, bool) {
	res, err := parseGraph(data)
	if err != nil || len(res) == 0 {
		return nil, nil, false
	}
	names := make([]string, 0, len(res))
	for n := range res {
		names = append(names, n)
	}
	sort.Strings(names)
	var cols []string
	seen := map[string]bool{}
	addCol := func(c string) {
		if !seen[c] {
			seen[c] = true
			cols = append(cols, c)
		}
	}
	if len(names) > 1 {
		addCol("formula")
	}
	var items []map[string]any
	for _, n := range names {
		f := res[n]
		if f.Error != "" {
			addCol("error")
			items = append(items, map[string]any{"formula": n, "error": f.Error})
			continue
		}
		for _, k := range f.Keys {
			addCol(k)
		}
		for _, rec := range f.Records {
			row := map[string]any{}
			if len(names) > 1 {
				row["formula"] = n
			}
			for i, v := range rec.Values {
				k := fmt.Sprintf("#%d", i)
				if i < len(f.Keys) {
					k = f.Keys[i]
				}
				addCol(k)
				row[k] = v
			}
			points := 0
			var last any
			for _, p := range rec.Timeseries {
				if p.Value != nil {
					points++
					last = *p.Value
				}
			}
			if rec.Aggregate != nil && rec.Aggregate.Value != nil {
				row["aggregate"] = *rec.Aggregate.Value
			}
			row["points"] = points
			row["last"] = last
			items = append(items, row)
		}
	}
	for _, c := range []string{"aggregate", "points", "last"} {
		addCol(c)
	}
	// Label keys such as "service.name" stay flat top-level keys, which the
	// table renderer looks up verbatim.
	out, err := json.Marshal(map[string]any{"items": items})
	if err != nil {
		return nil, nil, false
	}
	return out, cols, true
}

func printSummary(sm *metricSummary) error {
	switch flagOutput {
	case "", "json":
		return writeOutput(func(w io.Writer) error { return writeCompactJSON(w, sm) })
	case "table":
		return writeOutput(func(w io.Writer) error { return printSummaryTable(w, sm) })
	case "csv":
		rows, err := json.Marshal(map[string]any{"items": sm.Series})
		if err != nil {
			return err
		}
		return printResultWith(rows, output.Options{Format: "csv", Columns: flagColumns})
	}
	data, err := json.Marshal(sm)
	if err != nil {
		return err
	}
	return printResult(data)
}

// parseLookback parses a Go duration, also accepting whole days (7d) and
// weeks (1w).
func parseLookback(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	for suffix, unit := range map[string]time.Duration{"d": 24 * time.Hour, "w": 7 * 24 * time.Hour} {
		if n, ok := strings.CutSuffix(s, suffix); ok {
			v, err := strconv.Atoi(n)
			if err != nil {
				return 0, err
			}
			return time.Duration(v) * unit, nil
		}
	}
	return time.ParseDuration(s)
}

func parseISOTime(s string) (time.Time, error) {
	for _, layout := range []string{time.RFC3339Nano, "2006-01-02T15:04:05.000Z", "2006-01-02T15:04:05Z07:00", "2006-01-02T15:04Z"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid time %q: use ISO 8601 (2006-01-02T15:04:05.000Z)", s)
}

func isoMillis(t time.Time) string { return t.UTC().Format("2006-01-02T15:04:05.000Z") }

// resolveWindow turns --lookback/--from/--to into absolute bounds.
func resolveWindow(tf timeFlags) (time.Time, time.Time, error) {
	to := nowFunc().UTC().Truncate(time.Second)
	if tf.to != "" {
		t, err := parseISOTime(tf.to)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		to = t
	}
	if tf.from != "" {
		from, err := parseISOTime(tf.from)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		return from, to, nil
	}
	lb, err := parseLookback(tf.lookback)
	if err != nil || lb <= 0 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid --lookback %q", tf.lookback)
	}
	return to.Add(-lb), to, nil
}

// warnUnindexedGroupBy warns (best-effort) when a --group-by key is not an
// indexed metric dimension, since such a group-by silently returns one empty
// group instead of failing.
func warnUnindexedGroupBy(ctx context.Context, c *api.Client, groupBy []string) {
	needCheck := false
	for _, g := range groupBy {
		if !strings.HasPrefix(strings.TrimSpace(g), "@") {
			needCheck = true
		}
	}
	if !needCheck {
		return
	}
	q := url.Values{}
	q.Set("scope", "metric")
	data, err := c.Get(ctx, "/facet_keys", q)
	if err != nil {
		return // don't block the query on a discovery failure
	}
	var keys []struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &keys); err != nil {
		return
	}
	indexed := make(map[string]bool, len(keys))
	for _, k := range keys {
		if k.Key != "" {
			indexed[k.Key] = true
		}
		if k.Name != "" {
			indexed[k.Name] = true
		}
	}
	if len(indexed) == 0 {
		return
	}
	for _, g := range groupBy {
		g = strings.TrimSpace(g)
		// "@attr" group-bys address datapoint attributes, which are not in
		// the facet-key index but group correctly.
		if strings.HasPrefix(g, "@") {
			continue
		}
		if !indexed[g] {
			warnf("group-by key %q is not an indexed metric dimension; results will be ungrouped (one empty group). See `edx facets keys --scope metric`.", g)
		}
	}
}
