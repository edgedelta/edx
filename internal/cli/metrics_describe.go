package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/edgedelta/edx/internal/api"
	"github.com/edgedelta/edx/internal/output"
)

// metricLabel is one label (dimension) of a metric. Ref is the exact token to
// use in --filter and --group-by: resource labels by name, datapoint
// attributes with an "@" prefix.
type metricLabel struct {
	Name        string   `json:"name"`
	Kind        string   `json:"kind"`
	Ref         string   `json:"ref"`
	ValuesCount *uint64  `json:"values_count"`
	Samples     []string `json:"samples"`
}

type metricDescription struct {
	Name         string        `json:"name"`
	Found        bool          `json:"found"`
	Kind         string        `json:"kind"`
	Unit         string        `json:"unit"`
	Description  string        `json:"description"`
	Temporality  string        `json:"temporality,omitempty"`
	Monotonic    *bool         `json:"monotonic,omitempty"`
	PromSuffix   string        `json:"prom_suffix,omitempty"`
	SeriesCount  *uint64       `json:"series_count"`
	LatestSample string        `json:"latest_sample"`
	Age          string        `json:"age"`
	Window       string        `json:"window"`
	Labels       []metricLabel `json:"labels"`
	Examples     []string      `json:"examples"`
	Notes        []string      `json:"notes,omitempty"`
	Warnings     []string      `json:"warnings,omitempty"`
	Suggestions  []string      `json:"suggestions,omitempty"`
}

// rawMetricRow is the subset of a /metrics_v2/search row describe needs.
type rawMetricRow struct {
	TimeUnixNano int64          `json:"time_unix_nano"`
	Name         string         `json:"name"`
	Kind         string         `json:"kind"`
	Unit         string         `json:"unit"`
	Description  string         `json:"description"`
	Resource     map[string]any `json:"resource"`
	Attributes   map[string]any `json:"attributes"`
	Sum          *rawTemporal   `json:"sum"`
	Histogram    *rawTemporal   `json:"histogram"`
	ExpHistogram *rawTemporal   `json:"exponential_histogram"`
	Gauge        map[string]any `json:"gauge"`
}

type rawTemporal struct {
	AggregationTemporality string `json:"aggregation_temporality"`
	IsMonotonic            bool   `json:"is_monotonic"`
}

type cardinalityResp struct {
	LabelValuesCountTotal uint64 `json:"label_values_count_total"`
	LabelNamesCount       uint64 `json:"label_names_count"`
	TotalSubMetricCount   uint64 `json:"total_sub_metric_count"`
	Labels                []struct {
		LabelName        string `json:"label_name"`
		LabelKind        string `json:"label_kind"`
		LabelValuesCount uint64 `json:"label_values_count"`
		Cardinality      []struct {
			LabelValue  string `json:"label_value"`
			MetricCount int    `json:"metric_count"`
		} `json:"cardinality"`
	} `json:"labels"`
}

// nowFunc is the clock used for staleness and window math (overridden in tests).
var nowFunc = time.Now

func newMetricsDescribeCmd() *cobra.Command {
	var tf timeFlags
	var samples, maxLabels int
	cmd := &cobra.Command{
		Use:   "describe <name>",
		Short: "Everything needed to query a metric correctly, in one call",
		Long: `Describe a metric: kind (gauge/sum/histogram), unit, description, sum
temporality and monotonicity, Prometheus-style suffix, label keys split into
resource vs datapoint attribute with their cardinality and sample values,
series count, latest sample time (staleness) and ready-to-run example queries.

Each label carries "ref", the exact token to use in --filter and --group-by
(datapoint attributes take an "@" prefix).

Built from /metrics_v2/search (a few raw rows) and /cardinality; if one of them
fails the rest is still printed with a "warnings" field. When the metric is not
found, the closest metric names are suggested and the command exits non-zero.`,
		Example: `  edx metrics describe http.server.request.duration
  edx metrics describe ed.pipeline.bytes --lookback 24h`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			ctx := cmdContext(cmd)
			name := strings.TrimSpace(args[0])
			d := &metricDescription{Name: name, Labels: []metricLabel{}, Examples: []string{}}
			d.Window = tf.lookback
			if tf.from != "" {
				d.Window = tf.from + " -> " + tf.to
			}

			// 1. A few raw rows: kind, unit, description, temporality, latest timestamp.
			var rows []rawMetricRow
			sq := url.Values{}
			sq.Set("query", fmt.Sprintf("name:%q", name))
			sq.Set("limit", "20")
			tf.apply(sq)
			searchErr := error(nil)
			if data, err := c.Get(ctx, "/metrics_v2/search", sq); err != nil {
				searchErr = err
				d.Warnings = append(d.Warnings, "raw sample search failed: "+oneLine(err))
			} else if err := json.Unmarshal(data, &rows); err != nil {
				searchErr = err
				d.Warnings = append(d.Warnings, "raw sample search: unexpected response: "+oneLine(err))
			}

			// 2. Cardinality: label keys, value counts, top values, series count.
			var card *cardinalityResp
			cq := url.Values{}
			cq.Set("query", name+"{*}")
			tf.apply(cq)
			if data, err := c.Get(ctx, "/cardinality", cq); err != nil {
				d.Warnings = append(d.Warnings, "cardinality failed: "+oneLine(err))
			} else {
				var cr cardinalityResp
				if err := json.Unmarshal(data, &cr); err != nil {
					d.Warnings = append(d.Warnings, "cardinality: unexpected response: "+oneLine(err))
				} else {
					card = &cr
				}
			}

			applyRows(d, rows)
			applyCardinality(d, card, rows, samples, maxLabels)
			d.Found = len(rows) > 0 || (card != nil && (card.TotalSubMetricCount > 0 || len(card.Labels) > 0))

			// 3. Not found (or nothing could be fetched): check the name list for
			// an exact match outside the window, else suggest the closest names.
			if !d.Found {
				names, kinds, err := fetchMetricNames(cmd, c)
				switch {
				case err != nil:
					d.Warnings = append(d.Warnings, "metric name list failed: "+oneLine(err))
				case kinds[name] != "" || containsStr(names, name):
					d.Found = true
					if d.Kind == "" {
						d.Kind = kinds[name]
					}
					d.Notes = append(d.Notes, fmt.Sprintf("metric exists but has no samples in the window (%s); widen --lookback", d.Window))
				default:
					d.Suggestions = closestNames(name, names, 5)
				}
			}
			if !d.Found && searchErr != nil && card == nil && len(d.Suggestions) == 0 {
				return searchErr
			}
			if !d.Found {
				d.SeriesCount = nil
			}
			d.Examples = exampleQueries(d)

			if err := printDescription(d); err != nil {
				return err
			}
			if !d.Found {
				if len(d.Suggestions) > 0 {
					return fmt.Errorf("metric %q not found; closest: %s", name, strings.Join(d.Suggestions, ", "))
				}
				return fmt.Errorf("metric %q not found (metric names are exact and case-sensitive; see `edx metrics list --keyword <part>`)", name)
			}
			return nil
		},
	}
	tf.register(cmd, "1h")
	cmd.Flags().IntVar(&samples, "samples", 5, "sample values shown per label")
	cmd.Flags().IntVar(&maxLabels, "max-labels", 40, "maximum labels listed (lowest cardinality first within each kind)")
	return cmd
}

func printDescription(d *metricDescription) error {
	switch flagOutput {
	case "", "json":
		return writeOutput(func(w io.Writer) error { return writeCompactJSON(w, d) })
	}
	data, err := json.Marshal(d)
	if err != nil {
		return err
	}
	if flagOutput == "table" || flagOutput == "csv" {
		// Tabulate the labels; the header fields are in the json form.
		rows, _ := json.Marshal(map[string]any{"items": d.Labels})
		return printResultWith(rows, output.Options{Format: flagOutput, Columns: defaultCols([]string{"name", "kind", "ref", "values_count", "samples"})})
	}
	return printResult(data)
}

func defaultCols(cols []string) []string {
	if len(flagColumns) > 0 {
		return flagColumns
	}
	return cols
}

func oneLine(err error) string {
	s := strings.Join(strings.Fields(err.Error()), " ")
	return truncateStr(s, 200)
}

func normTemporality(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimPrefix(s, "aggregation_temporality_")
	return s
}

func applyRows(d *metricDescription, rows []rawMetricRow) {
	var latest int64
	for _, r := range rows {
		if d.Kind == "" {
			d.Kind = strings.ToLower(r.Kind)
		}
		if d.Unit == "" {
			d.Unit = r.Unit
		}
		if d.Description == "" {
			d.Description = r.Description
		}
		for kind, t := range map[string]*rawTemporal{"sum": r.Sum, "histogram": r.Histogram, "exponential_histogram": r.ExpHistogram} {
			if t == nil {
				continue
			}
			if d.Kind == "" {
				d.Kind = kind
			}
			if d.Temporality == "" && t.AggregationTemporality != "" {
				d.Temporality = normTemporality(t.AggregationTemporality)
			}
			if kind == "sum" && d.Monotonic == nil {
				m := t.IsMonotonic
				d.Monotonic = &m
			}
		}
		if d.Kind == "" && r.Gauge != nil {
			d.Kind = "gauge"
		}
		if r.TimeUnixNano > latest {
			latest = r.TimeUnixNano
		}
	}
	if latest > 0 {
		t := time.Unix(0, latest).UTC()
		d.LatestSample = t.Format(time.RFC3339)
		age := nowFunc().Sub(t).Round(time.Second)
		if age < 0 {
			age = 0
		}
		d.Age = fmtAge(age)
	}
	for _, suf := range []string{"_total", "_bucket", "_sum", "_count", "_created"} {
		if strings.HasSuffix(d.Name, suf) {
			d.PromSuffix = suf
			break
		}
	}
}

// fmtAge renders a staleness duration, switching to days past 48h.
func fmtAge(d time.Duration) string {
	if d < 48*time.Hour {
		return d.String()
	}
	days := int(d / (24 * time.Hour))
	return fmt.Sprintf("%dd%dh", days, int((d%(24*time.Hour))/time.Hour))
}

func applyCardinality(d *metricDescription, card *cardinalityResp, rows []rawMetricRow, samples, maxLabels int) {
	if card != nil {
		n := card.TotalSubMetricCount
		d.SeriesCount = &n
		for _, l := range card.Labels {
			kind := strings.ToLower(l.LabelKind)
			if kind == "" {
				kind = "resource"
			}
			vc := l.LabelValuesCount
			ml := metricLabel{Name: l.LabelName, Kind: kind, Ref: labelRef(l.LabelName, kind), ValuesCount: &vc, Samples: []string{}}
			for _, v := range l.Cardinality {
				if len(ml.Samples) >= samples {
					break
				}
				ml.Samples = append(ml.Samples, v.LabelValue)
			}
			d.Labels = append(d.Labels, ml)
		}
	} else if len(rows) > 0 {
		// Fall back to the label keys seen in the raw rows (cardinality unknown).
		seen := map[string]*metricLabel{}
		add := func(kind string, m map[string]any) {
			for k, v := range m {
				id := kind + "\x00" + k
				ml, ok := seen[id]
				if !ok {
					ml = &metricLabel{Name: k, Kind: kind, Ref: labelRef(k, kind), Samples: []string{}}
					seen[id] = ml
				}
				s := fmt.Sprint(v)
				if len(ml.Samples) < samples && !containsStr(ml.Samples, s) {
					ml.Samples = append(ml.Samples, s)
				}
			}
		}
		for _, r := range rows {
			add("resource", r.Resource)
			add("attribute", r.Attributes)
		}
		for _, ml := range seen {
			d.Labels = append(d.Labels, *ml)
		}
		d.Notes = append(d.Notes, fmt.Sprintf("labels and samples come from %d raw rows; values_count unknown", len(rows)))
	}
	sort.SliceStable(d.Labels, func(a, b int) bool {
		la, lb := d.Labels[a], d.Labels[b]
		if la.Kind != lb.Kind {
			return la.Kind == "resource"
		}
		return la.Name < lb.Name
	})
	if maxLabels > 0 && len(d.Labels) > maxLabels {
		d.Notes = append(d.Notes, fmt.Sprintf("%d labels, showing %d; raise --max-labels", len(d.Labels), maxLabels))
		d.Labels = d.Labels[:maxLabels]
	}
}

func labelRef(name, kind string) string {
	if kind == "attribute" && !strings.HasPrefix(name, "@") {
		return "@" + name
	}
	return name
}

// fetchMetricNames lists metric names (and their kinds, from extraFields).
func fetchMetricNames(cmd *cobra.Command, c *api.Client) ([]string, map[string]string, error) {
	q := url.Values{}
	q.Set("scope", "metric")
	q.Set("facet_path", "name")
	q.Set("limit", "1000")
	data, err := c.Get(cmdContext(cmd), "/facet_options", q)
	if err != nil {
		return nil, nil, err
	}
	var facet struct {
		Options []struct {
			Name        string            `json:"name"`
			Value       string            `json:"value"`
			ExtraFields map[string]string `json:"extraFields"`
		} `json:"options"`
	}
	if err := json.Unmarshal(data, &facet); err != nil {
		return nil, nil, err
	}
	names := make([]string, 0, len(facet.Options))
	kinds := map[string]string{}
	for _, o := range facet.Options {
		n := o.Name
		if n == "" {
			n = o.Value
		}
		if n == "" {
			continue
		}
		names = append(names, n)
		if k := o.ExtraFields["kind"]; k != "" {
			kinds[n] = strings.ToLower(k)
		}
	}
	return names, kinds, nil
}

func containsStr(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// closestNames ranks candidates by normalized edit distance, rewarding
// substring matches and shared name tokens (split on . _ -).
func closestNames(target string, names []string, n int) []string {
	t := strings.ToLower(target)
	tTokens := nameTokens(t)
	type scored struct {
		name  string
		score float64
	}
	var all []scored
	for _, name := range names {
		l := strings.ToLower(name)
		dist := float64(levenshtein(t, l)) / math.Max(float64(len(t)), float64(len(l)))
		score := dist - 0.5*jaccard(tTokens, nameTokens(l))
		if strings.Contains(l, t) || strings.Contains(t, l) {
			score -= 0.3
		}
		all = append(all, scored{name, score})
	}
	sort.SliceStable(all, func(a, b int) bool {
		if all[a].score != all[b].score {
			return all[a].score < all[b].score
		}
		return all[a].name < all[b].name
	})
	var out []string
	for i := 0; i < len(all) && len(out) < n; i++ {
		if all[i].score > 0.75 {
			break
		}
		out = append(out, all[i].name)
	}
	return out
}

func nameTokens(s string) map[string]bool {
	m := map[string]bool{}
	for _, tok := range strings.FieldsFunc(s, func(r rune) bool { return r == '.' || r == '_' || r == '-' || r == '/' }) {
		m[tok] = true
	}
	return m
}

func jaccard(a, b map[string]bool) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 0
	}
	inter := 0
	for k := range a {
		if b[k] {
			inter++
		}
	}
	return float64(inter) / float64(len(a)+len(b)-inter)
}

func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	prev := make([]int, len(rb)+1)
	cur := make([]int, len(rb)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(ra); i++ {
		cur[0] = i
		for j := 1; j <= len(rb); j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			cur[j] = min(prev[j]+1, cur[j-1]+1, prev[j-1]+cost)
		}
		prev, cur = cur, prev
	}
	return prev[len(rb)]
}

// exampleQueries suggests ready-to-run queries that fit the metric's kind.
func exampleQueries(d *metricDescription) []string {
	if !d.Found {
		return []string{}
	}
	base := "edx metrics query --name " + shellQuote(d.Name)
	groupBy := pickGroupBy(d.Labels)
	gb := ""
	if groupBy != "" {
		gb = " --group-by " + shellQuote(groupBy)
	}
	var aggs []string
	cumulativeCounter := d.Kind == "sum" && d.Monotonic != nil && *d.Monotonic && d.Temporality != "delta"
	switch {
	case d.PromSuffix == "_bucket":
		aggs = []string{"max"}
		d.Notes = append(d.Notes, "Prometheus histogram bucket series (one series per `le` bound): prefer the _count/_sum siblings, or p95 on a native histogram metric")
	case cumulativeCounter || (d.Kind == "" && d.PromSuffix == "_total"):
		aggs = []string{"rate", "increase"}
		d.Notes = append(d.Notes, "cumulative counter: use rate (per-second) or increase (per bucket); raw sum/avg of a cumulative counter only shows the running total")
	case d.Kind == "sum" && d.Temporality == "delta":
		aggs = []string{"sum"}
	case d.Kind == "sum":
		aggs = []string{"avg", "max"}
	case d.Kind == "histogram" || d.Kind == "exponential_histogram":
		aggs = []string{"p95", "avg"}
	case d.Kind == "gauge":
		aggs = []string{"avg", "max"}
	default:
		aggs = []string{"avg"}
	}
	out := []string{}
	for i, a := range aggs {
		q := base + " --agg " + a + gb + " --summary"
		if i == 0 {
			out = append(out, q, q+" --baseline 1w")
		} else {
			out = append(out, q)
		}
	}
	if f := pickFilter(d.Labels, groupBy); f != "" {
		out = append(out, base+" --agg "+aggs[0]+" --filter "+shellQuote(f)+" --summary")
	}
	return out
}

// pickGroupBy prefers a well-known service/host dimension with a small but
// non-trivial number of values, then any label with 2..100 values.
func pickGroupBy(labels []metricLabel) string {
	ok := func(l metricLabel) bool {
		return l.ValuesCount == nil && len(l.Samples) > 1 || l.ValuesCount != nil && *l.ValuesCount >= 2 && *l.ValuesCount <= 100
	}
	for _, pref := range []string{"service.name", "k8s.deployment.name", "k8s.namespace.name", "host.name", "k8s.pod.name"} {
		for _, l := range labels {
			if l.Name == pref && ok(l) {
				return l.Ref
			}
		}
	}
	best := ""
	var bestN uint64 = math.MaxUint64
	for _, l := range labels {
		if !ok(l) {
			continue
		}
		n := uint64(len(l.Samples))
		if l.ValuesCount != nil {
			n = *l.ValuesCount
		}
		if n < bestN {
			best, bestN = l.Ref, n
		}
	}
	return best
}

func pickFilter(labels []metricLabel, skip string) string {
	for _, l := range labels {
		if l.Ref != skip && len(l.Samples) > 0 && l.Samples[0] != "" {
			return l.Ref + ":" + strconv.Quote(l.Samples[0])
		}
	}
	for _, l := range labels {
		if len(l.Samples) > 0 && l.Samples[0] != "" {
			return l.Ref + ":" + strconv.Quote(l.Samples[0])
		}
	}
	return ""
}

// shellQuote single-quotes s when it holds shell metacharacters.
func shellQuote(s string) string {
	if s != "" && !strings.ContainsAny(s, " \t\"'$`\\|&;<>(){}*?[]!#~") {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
