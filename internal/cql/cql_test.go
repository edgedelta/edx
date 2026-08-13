package cql

import (
	"strings"
	"testing"
)

func TestValidateAcceptsWellFormedQueries(t *testing.T) {
	cases := []struct {
		dialect Dialect
		query   string
	}{
		// Variable references are ordinary TERMs, so dashboard templates parse unsubstituted.
		{DialectLog, "{ed.tag:$cluster AND k8s.namespace.name:$namespace} by {severity_text}"},
		{DialectLog, "count_unique():{ed.tag:$cluster} by {severity_text,event.type}"},
		{DialectLog, "{error AND ed.tag: $fleet} by {host.name}"},
		{DialectLog, `{service.name:"ed-agent-log" -ed.pipeline.node.type:"edge_rollup_service"}`},
		{DialectLog, "{(@alert.status:alert OR @alert.status:recovery) ed.tag:$fleet}"},
		{DialectLog, "ed.tag:$cluster"},
		{DialectLog, "*"},
		{DialectLog, "{a:b} by {c}.rollup(sum,d)"},
		{DialectMetric, "avg:ed.agent.cpu.milicores{ed.tag: $fleet} by {host.name}.rollup(60)"},
		{DialectMetric, "count_unique(k8s.pod.name):k8s.ksm.pod.status_phase.value{ed.tag:$cluster}.rollup(300)"},
		{DialectMetric, "sum:ed.pipeline.node.read_bytes{*} by {ed.tag}.rollup(60)"},
		{DialectMetric, "avg:m{*} by {a}.fill(zero,60).rollup(60)"},
		{DialectMetric, "$metric_name{*}"},
		{DialectFormula, "q1 + q2"},
		{DialectFormula, "(q1 + q2) * 100"},
		{DialectFormula, "timeshift(q1 + q2, 3600)"},
		{DialectFormula, "moving_average(q1, 5)"},
	}

	for _, c := range cases {
		errs, err := Validate(c.dialect, c.query)
		if err != nil {
			t.Fatalf("Validate(%s, %q): unexpected error: %v", c.dialect, c.query, err)
		}
		if len(errs) > 0 {
			t.Errorf("Validate(%s, %q) = %v, want no errors", c.dialect, c.query, errs)
		}
	}
}

func TestValidateRejectsMalformedQueries(t *testing.T) {
	cases := []struct {
		name    string
		dialect Dialect
		query   string
	}{
		{"unclosed brace", DialectLog, "{ed.tag:$fleet"},
		{"empty group by", DialectLog, "{a:b} by {"},
		{"dangling AND", DialectLog, "{a:b} AND"},
		{"unclosed paren", DialectLog, "{(a:b}"},
		{"rollup missing field", DialectLog, "{a:b} by {c}.rollup(sum)"},
		{"unclosed filter", DialectMetric, "avg:foo{"},
		{"non-numeric rollup", DialectMetric, "sum:foo{*}.rollup(abc)"},
		{"unknown fill method", DialectMetric, "sum:foo{*}.fill(bogus)"},
		// rollup takes a literal window; the frontend never substitutes a variable there.
		{"variable as rollup window", DialectMetric, "sum:foo{*}.rollup($window)"},
		{"timeshift missing shift", DialectFormula, "timeshift(q1, )"},
		{"unclosed function call", DialectFormula, "moving_average(q1"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			errs, err := Validate(c.dialect, c.query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(errs) == 0 {
				t.Fatalf("Validate(%s, %q) reported no errors, want at least one", c.dialect, c.query)
			}
			if errs[0].Message == "" {
				t.Error("error carries no message")
			}
		})
	}
}

// The dialects are not interchangeable: each rejects the other's shape. This is why
// QuerySites picks a grammar per data source type instead of using one tolerant grammar.
func TestDialectsAreNotInterchangeable(t *testing.T) {
	cases := []struct {
		dialect Dialect
		query   string
	}{
		{DialectLog, "sum:foo{*} by {host.name}.rollup(60)"},
		{DialectMetric, "{just.a:filter}"},
	}

	for _, c := range cases {
		errs, err := Validate(c.dialect, c.query)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(errs) == 0 {
			t.Errorf("%s dialect accepted %q, want rejection", c.dialect, c.query)
		}
	}
}

// Every query field in a definition is optional and the backend reads a missing filter
// as match-all, so an empty string must not be an error.
func TestValidateAcceptsEmptyQuery(t *testing.T) {
	for _, d := range []Dialect{DialectLog, DialectMetric, DialectFormula} {
		errs, err := Validate(d, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(errs) > 0 {
			t.Errorf("Validate(%s, \"\") = %v, want no errors", d, errs)
		}
	}
}

// An unrecognised dialect must be an error, never a silent pass.
func TestValidateRejectsUnknownDialect(t *testing.T) {
	errs, err := Validate(Dialect("sql"), "SELECT 1")
	if err == nil {
		t.Fatalf("Validate with unknown dialect returned errs=%v, want an error", errs)
	}
	if !strings.Contains(err.Error(), "sql") {
		t.Errorf("error %q does not name the offending dialect", err)
	}
}

func TestDialectForDataType(t *testing.T) {
	// The log-family data types share one grammar; this mapping is what lets a dashboard
	// and `edx cql validate --type trace` agree.
	for _, dataType := range []string{"log", "event", "pattern", "trace"} {
		if d, ok := DialectForDataType(dataType); !ok || d != DialectLog {
			t.Errorf("DialectForDataType(%q) = (%q, %v), want (log, true)", dataType, d, ok)
		}
	}
	if d, ok := DialectForDataType("metric"); !ok || d != DialectMetric {
		t.Errorf("metric mapped to (%q, %v)", d, ok)
	}
	if d, ok := DialectForDataType("formula"); !ok || d != DialectFormula {
		t.Errorf("formula mapped to (%q, %v)", d, ok)
	}
	// A dashboard's "empty" data source carries no query, so it has no grammar.
	if _, ok := DialectForDataType("empty"); ok {
		t.Error(`DialectForDataType("empty") reported a dialect, want none`)
	}
	if _, ok := DialectForDataType(""); ok {
		t.Error("the empty data type reported a dialect, want none")
	}
}

func TestDataTypesListsEveryMapping(t *testing.T) {
	got := DataTypes()
	want := []string{"event", "formula", "log", "metric", "pattern", "trace"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("DataTypes() = %v, want %v (sorted)", got, want)
	}
	for _, dataType := range got {
		if _, ok := DialectForDataType(dataType); !ok {
			t.Errorf("DataTypes() lists %q but DialectForDataType rejects it", dataType)
		}
	}
}

func TestAnnotatePutsTheCaretUnderTheProblem(t *testing.T) {
	const query = "sum:foo{*}.rollup(abc)"
	errs, err := Validate(DialectMetric, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) != 1 {
		t.Fatalf("got %d errors, want 1", len(errs))
	}

	lines := strings.Split(errs[0].Annotate(query), "\n")
	if len(lines) != 2 {
		t.Fatalf("Annotate returned %d lines, want 2: %q", len(lines), lines)
	}
	if lines[0] != query {
		t.Errorf("first line = %q, want the query", lines[0])
	}
	caret := strings.IndexByte(lines[1], '^')
	if caret != errs[0].Column {
		t.Errorf("caret at %d, want column %d", caret, errs[0].Column)
	}
	// The caret should land on `abc`, the token the parser objected to.
	if !strings.HasPrefix(query[caret:], "abc") {
		t.Errorf("caret points at %q, want it at `abc`", query[caret:])
	}
}

// Columns are byte offsets, so a multi-byte character earlier in the query must not push
// the caret out of alignment.
func TestAnnotateAlignsWithMultiByteCharacters(t *testing.T) {
	const query = `{service.name:"café"} by {`
	errs, err := Validate(DialectLog, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) == 0 {
		t.Fatal("expected a syntax error for the unclosed brace")
	}

	annotated := errs[0].Annotate(query)
	lines := strings.Split(annotated, "\n")
	if len(lines) != 2 {
		t.Fatalf("Annotate returned %d lines, want 2", len(lines))
	}
	// The caret's rune offset must match the rune offset of the byte column, not the
	// byte count, or every character after "café" would be off by one.
	wantRunes := len([]rune(query[:errs[0].Column]))
	if gotRunes := len([]rune(lines[1])) - 1; gotRunes != wantRunes {
		t.Errorf("caret at rune %d, want %d\n%s", gotRunes, wantRunes, annotated)
	}
}

// An out-of-range column must not panic; it just yields the query unannotated.
func TestAnnotateToleratesAnOutOfRangeColumn(t *testing.T) {
	const query = "sum:a{*}"
	for _, column := range []int{-1, len(query) + 1, 9999} {
		e := SyntaxError{Column: column, Message: "synthetic"}
		if got := e.Annotate(query); got != query {
			t.Errorf("Annotate with column %d = %q, want the bare query", column, got)
		}
	}
}

func TestSyntaxErrorLocatesTheProblem(t *testing.T) {
	errs, err := Validate(DialectMetric, "sum:foo{*}.rollup(abc)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) != 1 {
		t.Fatalf("got %d errors, want 1: %v", len(errs), errs)
	}
	if errs[0].Column != 18 {
		t.Errorf("Column = %d, want 18 (the position of `abc`)", errs[0].Column)
	}
	if !strings.Contains(errs[0].String(), "column 18") {
		t.Errorf("String() = %q, want it to carry the column", errs[0].String())
	}
}
