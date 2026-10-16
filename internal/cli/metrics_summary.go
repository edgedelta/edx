package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

// This file turns a /graph metric response into a compact, agent-oriented
// summary: per-series statistics, trend, anomaly windows and (optionally) a
// baseline comparison, instead of every raw point.

// graphFormula is one formula's result in a /graph response. A formula that
// failed carries only Error.
type graphFormula struct {
	From    string          `json:"from"`
	To      string          `json:"to"`
	Window  json.RawMessage `json:"window"`
	Keys    []string        `json:"keys"`
	Records []graphRecord   `json:"records"`
	Error   string          `json:"error"`
}

type graphRecord struct {
	Values     []any        `json:"values"`
	Timeseries []graphPoint `json:"timeseries"`
	Aggregate  *struct {
		Value *float64 `json:"value"`
	} `json:"aggregate"`
}

type graphPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     *float64  `json:"value"`
}

// parseGraph decodes a /graph response keyed by formula name.
func parseGraph(data []byte) (map[string]graphFormula, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("unexpected /graph response: %w", err)
	}
	out := make(map[string]graphFormula, len(raw))
	for name, msg := range raw {
		var f graphFormula
		if err := json.Unmarshal(msg, &f); err != nil {
			continue // not a formula object (unknown envelope field)
		}
		out[name] = f
	}
	return out, nil
}

// windowDuration parses the response "window" (a Go duration string, or a
// number of nanoseconds from older endpoints).
func (f graphFormula) windowDuration() time.Duration {
	if len(f.Window) == 0 {
		return 0
	}
	var s string
	if json.Unmarshal(f.Window, &s) == nil {
		if d, err := time.ParseDuration(s); err == nil {
			return d
		}
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			return time.Duration(n) * time.Second
		}
		return 0
	}
	var n float64
	if json.Unmarshal(f.Window, &n) == nil {
		return time.Duration(n)
	}
	return 0
}

// seriesLabel renders a record's group labels as "k=v,k2=v2" ("*" when the
// query has no group-by).
func seriesLabel(keys []string, values []any) string {
	if len(keys) == 0 && len(values) == 0 {
		return "*"
	}
	parts := make([]string, 0, len(values))
	for i, v := range values {
		k := fmt.Sprintf("#%d", i)
		if i < len(keys) {
			k = keys[i]
		}
		s := ""
		if v != nil {
			s = fmt.Sprint(v)
		}
		if s == "" {
			s = "(empty)"
		}
		parts = append(parts, k+"="+s)
	}
	if len(parts) == 0 {
		return "(empty)"
	}
	return strings.Join(parts, ",")
}

// sumLikeAgg reports whether the aggregation's values add up across series and
// time (so the series total, and each series' share of it, is meaningful).
func sumLikeAgg(agg string) bool {
	switch strings.ToLower(agg) {
	case "sum", "count", "increase":
		return true
	}
	return false
}

type summaryOptions struct {
	Agg       string
	Query     string
	Rollup    int
	Top       int
	AnomalyK  float64
	MaxWindow int // anomaly windows kept per series
}

type anomalyWindow struct {
	From    string  `json:"from"`
	To      string  `json:"to"`
	Buckets int     `json:"buckets"`
	Peak    float64 `json:"peak"`
	Dir     string  `json:"dir"` // high | low
	score   float64
}

type baselineCmp struct {
	Value    *float64 `json:"value"`
	Points   int      `json:"points"`
	Delta    *float64 `json:"delta"`
	DeltaPct *float64 `json:"delta_pct"`
	Verdict  string   `json:"verdict"` // higher | lower | similar | insufficient_data | missing
}

type seriesSummary struct {
	Series    string          `json:"series"`
	Value     *float64        `json:"value"`
	SharePct  *float64        `json:"share_pct,omitempty"`
	Min       *float64        `json:"min"`
	Max       *float64        `json:"max"`
	Avg       *float64        `json:"avg"`
	P95       *float64        `json:"p95"`
	First     *float64        `json:"first"`
	Last      *float64        `json:"last"`
	Change    *float64        `json:"change"`
	ChangePct *float64        `json:"change_pct"`
	Trend     string          `json:"trend"` // rising | falling | flat | no_data
	Points    int             `json:"points"`
	Gaps      int             `json:"gaps"`
	Anomalies []anomalyWindow `json:"anomalies"`
	Baseline  *baselineCmp    `json:"baseline,omitempty"`

	rank float64
}

type baselineHeader struct {
	Offset   string   `json:"offset"`
	From     string   `json:"from"`
	To       string   `json:"to"`
	Total    *float64 `json:"total,omitempty"`
	DeltaPct *float64 `json:"total_delta_pct,omitempty"`
	Missing  int      `json:"series_missing"`
}

type metricSummary struct {
	Query           string          `json:"query"`
	From            string          `json:"from"`
	To              string          `json:"to"`
	Rollup          string          `json:"rollup"`
	Agg             string          `json:"agg"`
	ValueStat       string          `json:"value_stat"` // what "value" means per series: total | avg
	RankedBy        string          `json:"ranked_by"`  // share | max
	SeriesCount     int             `json:"series_count"`
	Shown           int             `json:"shown"`
	Omitted         int             `json:"omitted"`
	Total           *float64        `json:"total,omitempty"`
	TopContributors []string        `json:"top_contributors"`
	AnomalyK        float64         `json:"anomaly_k"`
	Baseline        *baselineHeader `json:"baseline,omitempty"`
	Notes           []string        `json:"notes,omitempty"`
	Series          []seriesSummary `json:"series"`
}

// sig rounds to 6 significant digits so summaries stay compact and stable.
func sig(v float64) float64 {
	if v == 0 || math.IsNaN(v) || math.IsInf(v, 0) {
		return v
	}
	r, err := strconv.ParseFloat(strconv.FormatFloat(v, 'g', 6, 64), 64)
	if err != nil {
		return v
	}
	return r
}

func fp(v float64) *float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return nil
	}
	r := sig(v)
	return &r
}

func ts(t time.Time) string { return t.UTC().Format(time.RFC3339) }

func median(sorted []float64) float64 {
	n := len(sorted)
	if n == 0 {
		return math.NaN()
	}
	if n%2 == 1 {
		return sorted[n/2]
	}
	return (sorted[n/2-1] + sorted[n/2]) / 2
}

// percentile uses the nearest-rank method on a sorted slice.
func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return math.NaN()
	}
	rank := int(math.Ceil(p / 100 * float64(len(sorted))))
	if rank < 1 {
		rank = 1
	}
	return sorted[rank-1]
}

// trendOf classifies a series by the significance of its least-squares slope:
// rising/falling only when the slope's t-statistic is >= 2 (~95%) AND the
// fitted change across the window is at least 5% of the series mean, so tiny
// but perfectly linear drifts don't read as trends.
func trendOf(xs, ys []float64) string {
	n := float64(len(ys))
	if len(ys) < 4 {
		return "flat"
	}
	var mx, my float64
	for i := range ys {
		mx += xs[i]
		my += ys[i]
	}
	mx /= n
	my /= n
	var sxx, sxy float64
	for i := range ys {
		sxx += (xs[i] - mx) * (xs[i] - mx)
		sxy += (xs[i] - mx) * (ys[i] - my)
	}
	if sxx == 0 {
		return "flat"
	}
	slope := sxy / sxx
	var sse float64
	for i := range ys {
		r := ys[i] - (my + slope*(xs[i]-mx))
		sse += r * r
	}
	fitted := math.Abs(slope * (xs[len(xs)-1] - xs[0]))
	scale := math.Abs(my)
	if scale == 0 {
		scale = math.Abs(slope) // all-zero mean: any real slope is relative
	}
	if fitted < 0.05*scale || slope == 0 {
		return "flat"
	}
	se := math.Sqrt(sse / (n - 2) / sxx)
	if se > 0 && math.Abs(slope/se) < 2 {
		return "flat"
	}
	if slope > 0 {
		return "rising"
	}
	return "falling"
}

// robustScores returns each value's modified z-score, 0.6745·(x−median)/MAD
// (Iglewicz & Hoaglin). When MAD is 0 (over half the points identical) it
// falls back to the mean absolute deviation, (x−median)/(1.253314·MeanAD).
// It returns nil when the series is constant.
func robustScores(vals []float64) []float64 {
	sorted := append([]float64(nil), vals...)
	sort.Float64s(sorted)
	med := median(sorted)
	dev := make([]float64, len(vals))
	var meanAD float64
	for i, v := range vals {
		dev[i] = math.Abs(v - med)
		meanAD += dev[i]
	}
	meanAD /= float64(len(vals))
	sort.Float64s(dev)
	mad := median(dev)
	scores := make([]float64, len(vals))
	switch {
	case mad > 0:
		for i, v := range vals {
			scores[i] = 0.6745 * (v - med) / mad
		}
	case meanAD > 0:
		for i, v := range vals {
			scores[i] = (v - med) / (1.253314 * meanAD)
		}
	default:
		return nil
	}
	return scores
}

// anomaliesOf groups contiguous buckets whose modified z-score exceeds k (same
// direction, no null bucket in between) into windows, keeping the maxKeep
// strongest in time order. It returns the windows and the total found.
func anomaliesOf(pts []graphPoint, k float64, step time.Duration, maxKeep int) ([]anomalyWindow, int) {
	var vals []float64
	var idx []int
	for i, p := range pts {
		if p.Value != nil && !math.IsNaN(*p.Value) {
			vals = append(vals, *p.Value)
			idx = append(idx, i)
		}
	}
	if len(vals) < 5 {
		return nil, 0
	}
	scores := robustScores(vals)
	if scores == nil {
		return nil, 0
	}
	var wins []anomalyWindow
	var cur *anomalyWindow
	lastIdx := -2
	for j, s := range scores {
		dir := ""
		if s > k {
			dir = "high"
		} else if s < -k {
			dir = "low"
		}
		if dir == "" || cur == nil || cur.Dir != dir || idx[j] != lastIdx+1 {
			if cur != nil {
				wins = append(wins, *cur)
				cur = nil
			}
		}
		if dir != "" {
			if cur == nil {
				cur = &anomalyWindow{From: ts(pts[idx[j]].Timestamp), Dir: dir}
			}
			cur.Buckets++
			cur.To = ts(pts[idx[j]].Timestamp.Add(step))
			if math.Abs(s) > math.Abs(cur.score) {
				cur.score = s
				cur.Peak = sig(vals[j])
			}
		}
		lastIdx = idx[j]
	}
	if cur != nil {
		wins = append(wins, *cur)
	}
	total := len(wins)
	if maxKeep > 0 && len(wins) > maxKeep {
		sort.SliceStable(wins, func(a, b int) bool { return math.Abs(wins[a].score) > math.Abs(wins[b].score) })
		wins = wins[:maxKeep]
		sort.SliceStable(wins, func(a, b int) bool { return wins[a].From < wins[b].From })
	}
	return wins, total
}

// inferStep returns the bucket width: the response window when known, else
// the smallest gap between consecutive points.
func inferStep(window time.Duration, recs []graphRecord) time.Duration {
	if window > 0 {
		return window
	}
	var step time.Duration
	for _, r := range recs {
		for i := 1; i < len(r.Timeseries); i++ {
			d := r.Timeseries[i].Timestamp.Sub(r.Timeseries[i-1].Timestamp)
			if d > 0 && (step == 0 || d < step) {
				step = d
			}
		}
	}
	return step
}

// summarizeRecord computes the per-series statistics for one record.
func summarizeRecord(label string, rec graphRecord, agg string, k float64, step time.Duration, maxWin int) (seriesSummary, int) {
	s := seriesSummary{Series: label, Trend: "no_data", Anomalies: []anomalyWindow{}}
	pts := append([]graphPoint(nil), rec.Timeseries...)
	sort.SliceStable(pts, func(a, b int) bool { return pts[a].Timestamp.Before(pts[b].Timestamp) })

	var vals, xs []float64
	for i, p := range pts {
		if p.Value == nil || math.IsNaN(*p.Value) {
			s.Gaps++
			continue
		}
		vals = append(vals, *p.Value)
		xs = append(xs, float64(i))
		if i > 0 && step > 0 {
			if missing := int(math.Round(float64(p.Timestamp.Sub(pts[i-1].Timestamp))/float64(step))) - 1; missing > 0 {
				s.Gaps += missing
			}
		}
	}
	s.Points = len(vals)
	if len(vals) == 0 {
		if rec.Aggregate != nil && rec.Aggregate.Value != nil {
			s.Value = fp(*rec.Aggregate.Value)
			s.rank = *rec.Aggregate.Value
		}
		return s, 0
	}
	sorted := append([]float64(nil), vals...)
	sort.Float64s(sorted)
	var total float64
	for _, v := range vals {
		total += v
	}
	avg := total / float64(len(vals))
	first, last := vals[0], vals[len(vals)-1]
	s.Min, s.Max, s.Avg = fp(sorted[0]), fp(sorted[len(sorted)-1]), fp(avg)
	s.P95 = fp(percentile(sorted, 95))
	s.First, s.Last = fp(first), fp(last)
	s.Change = fp(last - first)
	if first != 0 {
		s.ChangePct = fp((last - first) / math.Abs(first) * 100)
	}
	s.Trend = trendOf(xs, vals)
	if sumLikeAgg(agg) {
		s.Value = fp(total)
		s.rank = total
	} else {
		s.Value = fp(avg)
		s.rank = sorted[len(sorted)-1]
	}
	wins, n := anomaliesOf(pts, k, step, maxWin)
	if wins != nil {
		s.Anomalies = wins
	}
	return s, n
}

// summarizeGraph builds the summary for one formula result.
func summarizeGraph(f graphFormula, o summaryOptions) *metricSummary {
	sumLike := sumLikeAgg(o.Agg)
	step := inferStep(f.windowDuration(), f.Records)
	sm := &metricSummary{
		Query:           o.Query,
		From:            f.From,
		To:              f.To,
		Agg:             o.Agg,
		AnomalyK:        o.AnomalyK,
		SeriesCount:     len(f.Records),
		TopContributors: []string{},
		Series:          []seriesSummary{},
	}
	switch {
	case step > 0:
		sm.Rollup = step.String()
	case o.Rollup > 0:
		sm.Rollup = (time.Duration(o.Rollup) * time.Second).String()
	default:
		sm.Rollup = "auto"
	}
	if sumLike {
		sm.ValueStat, sm.RankedBy = "total", "share"
	} else {
		sm.ValueStat, sm.RankedBy = "avg", "max"
	}
	var all []seriesSummary
	var grand float64
	for _, rec := range f.Records {
		s, nWin := summarizeRecord(seriesLabel(f.Keys, rec.Values), rec, o.Agg, o.AnomalyK, step, o.MaxWindow)
		if nWin > len(s.Anomalies) {
			sm.Notes = append(sm.Notes, fmt.Sprintf("%s: %d anomaly windows, showing the %d strongest", s.Series, nWin, len(s.Anomalies)))
		}
		if sumLike && s.Value != nil {
			grand += s.rank
		}
		all = append(all, s)
	}
	sort.SliceStable(all, func(a, b int) bool {
		if all[a].rank != all[b].rank {
			return all[a].rank > all[b].rank
		}
		return all[a].Series < all[b].Series
	})
	if sumLike {
		sm.Total = fp(grand)
		for i := range all {
			if grand != 0 && all[i].Value != nil {
				all[i].SharePct = fp(math.Round(all[i].rank/grand*1000) / 10)
			}
		}
	}
	for i := 0; i < len(all) && i < 3; i++ {
		if all[i].Points == 0 {
			break
		}
		if sumLike && all[i].SharePct != nil {
			sm.TopContributors = append(sm.TopContributors, fmt.Sprintf("%s (%s%%)", all[i].Series, fmtNum(*all[i].SharePct)))
		} else if all[i].Max != nil {
			sm.TopContributors = append(sm.TopContributors, fmt.Sprintf("%s (max %s)", all[i].Series, fmtNum(*all[i].Max)))
		}
	}
	shown := all
	if o.Top > 0 && len(all) > o.Top {
		shown = all[:o.Top]
	}
	sm.Shown = len(shown)
	sm.Omitted = len(all) - len(shown)
	if sm.Omitted > 0 {
		sm.Notes = append(sm.Notes, fmt.Sprintf("%d lower-ranked series omitted; raise --top or narrow --filter", sm.Omitted))
	}
	sm.Series = shown
	return sm
}

// Baseline verdict thresholds (ed-baseline-compare): a comparison needs enough
// samples in BOTH windows, and a change must clear a relative threshold.
const (
	baselineMinPoints   = 3    // absolute floor of non-null buckets per window
	baselineMinCount    = 10.0 // absolute floor of the value for count aggregations
	baselineSimilarPct  = 10.0 // |delta %| below this reads as "similar"
	baselinePointsRatio = 0.5  // baseline needs >= half the current window's points
)

// applyBaseline adds per-series baseline comparisons from a time-shifted query.
func applyBaseline(sm *metricSummary, base graphFormula, o summaryOptions, offset string, bFrom, bTo time.Time) {
	step := inferStep(base.windowDuration(), base.Records)
	byLabel := map[string]seriesSummary{}
	var baseGrand float64
	for _, rec := range base.Records {
		s, _ := summarizeRecord(seriesLabel(base.Keys, rec.Values), rec, o.Agg, o.AnomalyK, step, 0)
		byLabel[s.Series] = s
		if s.Value != nil {
			baseGrand += s.rank
		}
	}
	hdr := &baselineHeader{Offset: offset, From: ts(bFrom), To: ts(bTo)}
	countAgg := strings.EqualFold(o.Agg, "count") || strings.EqualFold(o.Agg, "increase")
	for i := range sm.Series {
		cur := &sm.Series[i]
		b, ok := byLabel[cur.Series]
		if !ok || b.Points == 0 {
			hdr.Missing++
			cur.Baseline = &baselineCmp{Verdict: "missing"}
			continue
		}
		cmp := &baselineCmp{Value: b.Value, Points: b.Points}
		if cur.Value != nil && b.Value != nil {
			cmp.Delta = fp(*cur.Value - *b.Value)
			if *b.Value != 0 {
				cmp.DeltaPct = fp(math.Round((*cur.Value-*b.Value)/math.Abs(*b.Value)*1000) / 10)
			}
		}
		minPts := baselineMinPoints
		if r := int(math.Ceil(baselinePointsRatio * float64(cur.Points))); r > minPts {
			minPts = r
		}
		switch {
		case cur.Value == nil || b.Points < minPts || cur.Points < baselineMinPoints:
			cmp.Verdict = "insufficient_data"
		case countAgg && math.Max(math.Abs(*cur.Value), math.Abs(*b.Value)) < baselineMinCount:
			cmp.Verdict = "insufficient_data"
		case cmp.DeltaPct == nil && cmp.Delta != nil && *cmp.Delta > 0:
			cmp.Verdict = "higher"
		case cmp.DeltaPct == nil:
			cmp.Verdict = "similar"
		case math.Abs(*cmp.DeltaPct) < baselineSimilarPct:
			cmp.Verdict = "similar"
		case *cmp.DeltaPct > 0:
			cmp.Verdict = "higher"
		default:
			cmp.Verdict = "lower"
		}
		cur.Baseline = cmp
	}
	if sumLikeAgg(o.Agg) {
		hdr.Total = fp(baseGrand)
		if sm.Total != nil && baseGrand != 0 {
			hdr.DeltaPct = fp(math.Round((*sm.Total-baseGrand)/math.Abs(baseGrand)*1000) / 10)
		}
	}
	sm.Baseline = hdr
}

func fmtNum(v float64) string {
	if a := math.Abs(v); a == 0 || (a >= 1e-4 && a < 1e15) {
		return strconv.FormatFloat(sig(v), 'f', -1, 64)
	}
	return strconv.FormatFloat(v, 'g', 6, 64)
}

func fmtPtr(v *float64) string {
	if v == nil {
		return "-"
	}
	return fmtNum(*v)
}

// writeCompactJSON writes v (a struct or map) as JSON with one top-level key
// per line and every array of objects one element per line: still valid,
// greppable JSON, but a fraction of the size of fully indented output.
func writeCompactJSON(w io.Writer, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if _, err := dec.Token(); err != nil { // '{'
		return err
	}
	var buf bytes.Buffer
	buf.WriteString("{")
	first := true
	for dec.More() {
		keyTok, err := dec.Token()
		if err != nil {
			return err
		}
		var val json.RawMessage
		if err := dec.Decode(&val); err != nil {
			return err
		}
		if !first {
			buf.WriteString(",")
		}
		first = false
		key, _ := json.Marshal(keyTok)
		fmt.Fprintf(&buf, "\n  %s: ", key)
		var elems []json.RawMessage
		if len(val) > 0 && val[0] == '[' && json.Unmarshal(val, &elems) == nil && len(elems) > 0 && elems[0][0] == '{' {
			buf.WriteString("[")
			for i, e := range elems {
				if i > 0 {
					buf.WriteString(",")
				}
				buf.WriteString("\n    ")
				buf.Write(e)
			}
			buf.WriteString("\n  ]")
			continue
		}
		buf.Write(val)
	}
	buf.WriteString("\n}\n")
	_, err = w.Write(buf.Bytes())
	return err
}

// printSummaryTable renders a summary for humans (-o table).
func printSummaryTable(w io.Writer, sm *metricSummary) error {
	fmt.Fprintf(w, "query:    %s\n", sm.Query)
	fmt.Fprintf(w, "window:   %s -> %s (rollup %s)\n", sm.From, sm.To, sm.Rollup)
	fmt.Fprintf(w, "series:   %d (showing %d, %d omitted); value = %s, ranked by %s\n", sm.SeriesCount, sm.Shown, sm.Omitted, sm.ValueStat, sm.RankedBy)
	if sm.Total != nil {
		fmt.Fprintf(w, "total:    %s\n", fmtNum(*sm.Total))
	}
	if len(sm.TopContributors) > 0 {
		fmt.Fprintf(w, "top:      %s\n", strings.Join(sm.TopContributors, "; "))
	}
	if b := sm.Baseline; b != nil {
		line := fmt.Sprintf("baseline: -%s %s -> %s", b.Offset, b.From, b.To)
		if b.Total != nil {
			line += fmt.Sprintf(", total %s", fmtNum(*b.Total))
			if b.DeltaPct != nil {
				line += fmt.Sprintf(" (%+g%%)", *b.DeltaPct)
			}
		}
		if b.Missing > 0 {
			line += fmt.Sprintf(", %d series missing", b.Missing)
		}
		fmt.Fprintln(w, line)
	}
	for _, n := range sm.Notes {
		fmt.Fprintf(w, "note:     %s\n", n)
	}
	fmt.Fprintln(w)
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	hdr := "SERIES\tVALUE\tSHARE%\tMIN\tAVG\tMAX\tP95\tLAST\tCHANGE%\tTREND\tPTS\tGAPS\tANOMALIES"
	if sm.Baseline != nil {
		hdr += "\tBASELINE\tDELTA%\tVERDICT"
	}
	fmt.Fprintln(tw, hdr)
	for _, s := range sm.Series {
		var an []string
		for _, a := range s.Anomalies {
			an = append(an, fmt.Sprintf("%s %s..%s peak %s", a.Dir, clock(a.From), clock(a.To), fmtNum(a.Peak)))
		}
		anom := "-"
		if len(an) > 0 {
			anom = strings.Join(an, "; ")
		}
		row := []string{truncateStr(s.Series, 60), fmtPtr(s.Value), fmtPtr(s.SharePct), fmtPtr(s.Min), fmtPtr(s.Avg), fmtPtr(s.Max),
			fmtPtr(s.P95), fmtPtr(s.Last), fmtPtr(s.ChangePct), s.Trend, strconv.Itoa(s.Points), strconv.Itoa(s.Gaps), anom}
		if sm.Baseline != nil {
			if b := s.Baseline; b != nil {
				row = append(row, fmtPtr(b.Value), fmtPtr(b.DeltaPct), b.Verdict)
			} else {
				row = append(row, "-", "-", "-")
			}
		}
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}
	return tw.Flush()
}

// clock shortens an RFC3339 timestamp to HH:MM for table cells.
func clock(s string) string {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return s
	}
	return t.UTC().Format("15:04")
}

func truncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}
