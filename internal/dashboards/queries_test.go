package dashboards

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/edgedelta/edx/internal/cql"
)

// The shipped dashboards are the ground truth for false positives: every query in them
// is one the product runs, so all of them must parse. They are full of unsubstituted
// `$variable` references, which is the case worth pinning down — it is the reason
// validating a definition at rest works at all.
func TestFixtureQueriesAllParse(t *testing.T) {
	var checked int
	for name, definition := range loadFixtures(t) {
		t.Run(name, func(t *testing.T) {
			for _, site := range QuerySites(definition) {
				errs, err := cql.Validate(site.Dialect, site.Query)
				if err != nil {
					t.Fatalf("%s: %v", site.Path, err)
				}
				if len(errs) > 0 {
					t.Errorf("%s: %s query %q failed to parse: %v", site.Path, site.Dialect, site.Query, errs)
				}
				checked++
			}
		})
	}
	if checked == 0 {
		t.Fatal("no queries found in the fixtures; the walk or the fixtures are broken")
	}
	t.Logf("parsed %d queries across the shipped dashboards", checked)
}

// QuerySites walks the definition's known shape, so a new query field in the schema would
// go unchecked. This finds query-ish keys generically and fails if the walk missed one.
func TestQuerySitesFindsEveryQueryInFixtures(t *testing.T) {
	for name, definition := range loadFixtures(t) {
		found := map[string]bool{}
		for _, site := range QuerySites(definition) {
			found[site.Path] = true
		}

		for _, path := range scanForQueryKeys(definition, "") {
			if !found[path] {
				t.Errorf("%s: QuerySites missed %s; add it to the walk in queries.go", name, path)
			}
		}
	}
}

// scanForQueryKeys finds every string-valued "query" or "formula" key anywhere in the
// definition, as a cross-check on the shape-aware walk.
func scanForQueryKeys(node any, path string) []string {
	var out []string
	switch n := node.(type) {
	case map[string]any:
		for key, value := range n {
			child := path + "/" + key
			if s, ok := value.(string); ok && (key == "query" || key == "formula") && s != "" {
				out = append(out, child)
				continue
			}
			out = append(out, scanForQueryKeys(value, child)...)
		}
	case []any:
		for i, value := range n {
			out = append(out, scanForQueryKeys(value, path+"/"+strconv.Itoa(i))...)
		}
	}
	sort.Strings(out)
	return out
}

func TestQuerySitesAssignsDialectsByDataSourceType(t *testing.T) {
	definition := map[string]any{
		"version": float64(4),
		"widgets": []any{
			map[string]any{
				"id":   float64(1),
				"type": "viz",
				"visuals": []any{
					map[string]any{"id": "A", "dataSource": map[string]any{
						"type": "log", "params": map[string]any{"query": "{a:b}"},
					}},
					map[string]any{"id": "B", "dataSource": map[string]any{
						"type": "metric", "params": map[string]any{"query": "sum:m{*}"},
					}},
					map[string]any{"id": "C", "dataSource": map[string]any{
						"type": "trace", "params": map[string]any{"query": "{c:d}"},
					}},
					map[string]any{"id": "D", "dataSource": map[string]any{
						"type": "formula", "params": map[string]any{"formula": "q1 + q2"},
					}},
					// No query to check: "empty" carries none.
					map[string]any{"id": "E", "dataSource": map[string]any{"type": "empty"}},
				},
			},
		},
		"variables": []any{
			map[string]any{"id": float64(1), "type": "query", "value": "ed.tag:$x"},
			map[string]any{"id": float64(2), "type": "facet-option", "params": map[string]any{"query": "ed.tag:$y"}},
		},
	}

	want := []QuerySite{
		{Path: "/widgets/0/visuals/0/dataSource/params/query", Dialect: cql.DialectLog, Query: "{a:b}"},
		{Path: "/widgets/0/visuals/1/dataSource/params/query", Dialect: cql.DialectMetric, Query: "sum:m{*}"},
		{Path: "/widgets/0/visuals/2/dataSource/params/query", Dialect: cql.DialectLog, Query: "{c:d}"},
		{Path: "/widgets/0/visuals/3/dataSource/params/formula", Dialect: cql.DialectFormula, Query: "q1 + q2"},
		{Path: "/variables/0/value", Dialect: cql.DialectLog, Query: "ed.tag:$x"},
		{Path: "/variables/1/params/query", Dialect: cql.DialectLog, Query: "ed.tag:$y"},
	}

	got := QuerySites(definition)
	if len(got) != len(want) {
		t.Fatalf("found %d sites, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("site %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// compareDataSource is a Partial<DataSource> merged over the primary, so it inherits the
// primary's type when it omits one — the dialect has to follow.
func TestQuerySitesInheritsCompareDataSourceType(t *testing.T) {
	visual := func(compare map[string]any) any {
		return map[string]any{
			"id":   float64(1),
			"type": "viz",
			"visuals": []any{map[string]any{
				"id":                "A",
				"dataSource":        map[string]any{"type": "metric", "params": map[string]any{"query": "sum:m{*}"}},
				"compareDataSource": compare,
			}},
		}
	}

	t.Run("inherits the type when it has its own params", func(t *testing.T) {
		got := QuerySites(map[string]any{"widgets": []any{visual(map[string]any{
			"params": map[string]any{"query": "avg:m{*}"},
		})}})
		if len(got) != 2 {
			t.Fatalf("found %d sites, want 2: %+v", len(got), got)
		}
		compare := got[1]
		if compare.Path != "/widgets/0/visuals/0/compareDataSource/params/query" {
			t.Errorf("Path = %q", compare.Path)
		}
		if compare.Dialect != cql.DialectMetric {
			t.Errorf("Dialect = %q, want metric inherited from dataSource", compare.Dialect)
		}
	})

	// The frontend's merge is a shallow spread, so a compare source without `params`
	// reuses the primary's query. Reporting it twice would double every error.
	t.Run("adds nothing when it only offsets the time range", func(t *testing.T) {
		got := QuerySites(map[string]any{"widgets": []any{visual(map[string]any{
			"timeFiltersOffset": "inherit",
		})}})
		if len(got) != 1 {
			t.Fatalf("found %d sites, want 1: %+v", len(got), got)
		}
	})
}

// breakFirstMetricQuery replaces the first metric query in a definition with a
// syntactically invalid one, and reports where it did so.
func breakFirstMetricQuery(t *testing.T, definition map[string]any, query string) string {
	t.Helper()
	for _, w := range array(definition["widgets"]) {
		widget, _ := w.(map[string]any)
		for _, v := range array(widget["visuals"]) {
			visual, _ := v.(map[string]any)
			source, _ := visual["dataSource"].(map[string]any)
			if kind, _ := source["type"].(string); kind != "metric" {
				continue
			}
			params, ok := source["params"].(map[string]any)
			if !ok {
				continue
			}
			params["query"] = query
			return query
		}
	}
	t.Fatal("no metric data source in the fixture to break")
	return ""
}

func TestValidateDefinitionReportsQueryErrors(t *testing.T) {
	definition := loadFixture(t, "fleet-resource-usage.json")
	breakFirstMetricQuery(t, definition, "avg:ed.agent.cpu.milicores{ed.tag: $fleet} by {host.name}.rollup(sixty)")

	issues, err := ValidateDefinition(definition)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("got %d issues, want exactly 1:\n%s", len(issues), format(issues))
	}
	if !strings.Contains(issues[0].Path, "/dataSource/params/query") {
		t.Errorf("Path = %q, want it to point at the query", issues[0].Path)
	}
	if !strings.Contains(issues[0].Message, "invalid metric query") {
		t.Errorf("Message = %q, want it to name the dialect", issues[0].Message)
	}
}

// A malformed query is not a schema violation — `query` is just a string — so the schema
// check has to stay silent about it and the query check has to catch it.
func TestSchemaAndQueryChecksAreIndependent(t *testing.T) {
	definition := loadFixture(t, "fleet-resource-usage.json")
	definition["unknownTopLevel"] = true
	breakFirstMetricQuery(t, definition, "sum:m{*}.rollup(nope)")

	issues, err := ValidateDefinition(definition)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var sawSchema, sawQuery bool
	for _, issue := range issues {
		if strings.Contains(issue.Message, "unknownTopLevel") {
			sawSchema = true
		}
		if strings.Contains(issue.Message, "invalid metric query") {
			sawQuery = true
		}
	}
	if !sawSchema {
		t.Errorf("schema violation not reported:\n%s", format(issues))
	}
	if !sawQuery {
		t.Errorf("query error not reported:\n%s", format(issues))
	}
}

func format(issues []Issue) string {
	var b strings.Builder
	for _, issue := range issues {
		b.WriteString("  " + issue.String() + "\n")
	}
	return b.String()
}

func loadFixtures(t *testing.T) map[string]map[string]any {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join("testdata", "definitions", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) == 0 {
		t.Fatal("no fixtures found")
	}
	out := map[string]map[string]any{}
	for _, path := range paths {
		out[filepath.Base(path)] = loadFixture(t, filepath.Base(path))
	}
	return out
}

func loadFixture(t *testing.T, name string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", "definitions", name))
	if err != nil {
		t.Fatal(err)
	}
	var definition map[string]any
	if err := json.Unmarshal(raw, &definition); err != nil {
		t.Fatal(err)
	}
	return definition
}
