package dashboards

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// oraclePath holds the output of the frontend's own fillDashboardResources() run over every
// shipped dashboard. Refresh it with `make sync-resource-accesses-oracle`, which runs
// `bun gen:resource-accesses` in the monorepo's web/.
const oraclePath = "testdata/resource-accesses.json"

// ResourceAccesses has to agree with the UI, or a dashboard created by edx and one created
// in the browser would grant different access when shared. This diffs against the real
// implementation's output rather than against expectations I wrote by hand.
func TestResourceAccessesMatchesTheFrontend(t *testing.T) {
	raw, err := os.ReadFile(oraclePath)
	if err != nil {
		t.Fatal(err)
	}
	var oracle map[string][]ResourceAccess
	if err := json.Unmarshal(raw, &oracle); err != nil {
		t.Fatal(err)
	}
	if len(oracle) == 0 {
		t.Fatal("oracle is empty")
	}

	var compared int
	for id, want := range oracle {
		t.Run(id, func(t *testing.T) {
			definition := loadFixture(t, id+".json")
			got := ResourceAccesses(definition)

			if len(got) != len(want) {
				t.Fatalf("got %d entries, want %d\ngot:  %s\nwant: %s",
					len(got), len(want), render(got), render(want))
			}

			for i := range want {
				// Order matters: this asserts the same derivation order as the UI, not
				// just the same set.
				if got[i].Domain != want[i].Domain {
					t.Errorf("entry %d domain = %q, want %q", i, got[i].Domain, want[i].Domain)
					continue
				}
				if !samePtr(got[i].Scope, want[i].Scope) {
					t.Errorf("entry %d (%s) scope = %s, want %s", i, want[i].Domain, show(got[i].Scope), show(want[i].Scope))
				}
				if !samePtr(got[i].FacetPath, want[i].FacetPath) {
					t.Errorf("entry %d (%s) facet_path = %s, want %s", i, want[i].Domain, show(got[i].FacetPath), show(want[i].FacetPath))
				}

				// An event entry's query is the one documented difference: the UI conjoins
				// an event.domain precondition. Assert the shape of that difference rather
				// than skipping it, so this fails if the UI stops rewriting (making the
				// simplification unnecessary) or if edx starts.
				if want[i].Domain == "event" {
					assertEventQueryDifference(t, i, got[i].Query, want[i].Query)
					continue
				}
				if !samePtr(got[i].Query, want[i].Query) {
					t.Errorf("entry %d (%s) query = %s, want %s", i, want[i].Domain, show(got[i].Query), show(want[i].Query))
				}
			}
			compared += len(want)
		})
	}
	t.Logf("compared %d entries across %d dashboards", compared, len(oracle))
}

// assertEventQueryDifference bounds the known event divergence. The frontend conjoins an
// event.domain precondition when the data source names an eventDomain, and leaves the
// query alone when it does not — so the two must agree exactly in the second case, and
// differ only by that precondition in the first.
func assertEventQueryDifference(t *testing.T, i int, got, want *string) {
	t.Helper()
	if got == nil || want == nil {
		t.Errorf("entry %d (event): got %s, want %s", i, show(got), show(want))
		return
	}
	if *got == *want {
		// No eventDomain on the data source, so the frontend rewrote nothing.
		return
	}
	if !strings.Contains(*want, "event.domain:") {
		t.Errorf("entry %d (event): queries differ but the frontend's %q carries no "+
			"event.domain precondition, so this is an unexplained divergence from %q",
			i, *want, *got)
	}
}

// Widget and variable types that resolve no resources must contribute nothing, so a
// dashboard of only layout and text produces an empty allowlist.
func TestResourceAccessesIgnoresNonQueryingTypes(t *testing.T) {
	definition := map[string]any{
		"version": float64(4),
		"widgets": []any{
			map[string]any{"id": "root", "type": "grid", "grid": "72px / 1fr"},
			map[string]any{"id": float64(1), "type": "markdown", "params": map[string]any{"content": "hi"}},
			map[string]any{"id": float64(2), "type": "variable-control", "variableId": float64(1)},
			map[string]any{"id": float64(3), "type": "viz", "resultType": "empty",
				"visualizer": map[string]any{"type": "empty"},
				"visuals": []any{
					// Neither an empty source nor a formula queries anything.
					map[string]any{"id": "A", "dataSource": map[string]any{"type": "empty"}},
					map[string]any{"id": "W", "dataSource": map[string]any{
						"type": "formula", "params": map[string]any{"formula": "A + 1"},
					}},
				}},
		},
		"variables": []any{
			map[string]any{"id": float64(1), "type": "string", "key": "s", "label": "S"},
			map[string]any{"id": float64(2), "type": "duration", "key": "d", "label": "D"},
			map[string]any{"id": float64(3), "type": "query", "key": "q", "label": "Q", "value": "a:b"},
			// metric-name shares the facet-option editor but resolves no resources.
			map[string]any{"id": float64(4), "type": "metric-name", "key": "m", "label": "M"},
		},
	}

	if got := ResourceAccesses(definition); len(got) != 0 {
		t.Errorf("got %s, want no entries", render(got))
	}
}

// compareDataSource is not resolved by the UI's viz manifest, so it must not add an entry.
func TestResourceAccessesSkipsCompareDataSource(t *testing.T) {
	definition := map[string]any{
		"widgets": []any{map[string]any{
			"id": float64(1), "type": "viz",
			"visuals": []any{map[string]any{
				"id":                "A",
				"dataSource":        map[string]any{"type": "metric", "params": map[string]any{"query": "sum:a{*}"}},
				"compareDataSource": map[string]any{"params": map[string]any{"query": "sum:b{*}"}},
			}},
		}},
	}

	got := ResourceAccesses(definition)
	if len(got) != 1 {
		t.Fatalf("got %s, want exactly the primary data source", render(got))
	}
	if *got[0].Query != "sum:a{*}" {
		t.Errorf("query = %q, want the primary query", *got[0].Query)
	}
}

// The same query in two widgets yields two entries in the UI, so it must here too.
func TestResourceAccessesKeepsDuplicates(t *testing.T) {
	widget := func(id float64) any {
		return map[string]any{
			"id": id, "type": "viz",
			"visuals": []any{map[string]any{
				"id":         "A",
				"dataSource": map[string]any{"type": "metric", "params": map[string]any{"query": "sum:a{*}"}},
			}},
		}
	}
	got := ResourceAccesses(map[string]any{"widgets": []any{widget(1), widget(2)}})
	if len(got) != 2 {
		t.Errorf("got %d entries, want 2 (duplicates kept): %s", len(got), render(got))
	}
}

// A facet-option variable without params falls back to the UI's defaults, and a widget
// without a query still emits an empty one rather than omitting the key.
func TestResourceAccessesAppliesFrontendDefaults(t *testing.T) {
	got := ResourceAccesses(map[string]any{
		"widgets": []any{map[string]any{
			"id": float64(1), "type": "viz",
			"visuals": []any{map[string]any{"id": "A", "dataSource": map[string]any{"type": "log"}}},
		}},
		"variables": []any{
			map[string]any{"id": float64(1), "type": "facet-option", "key": "k", "label": "L"},
			map[string]any{"id": float64(2), "type": "facet", "key": "f", "label": "F"},
		},
	})

	if len(got) != 3 {
		t.Fatalf("got %s, want 3 entries", render(got))
	}
	if got[0].Query == nil || *got[0].Query != "" {
		t.Errorf("log entry query = %s, want an empty string, not an absent key", show(got[0].Query))
	}
	if got[1].Scope == nil || *got[1].Scope != "metric" || got[1].FacetPath == nil || *got[1].FacetPath != "ed.tag" {
		t.Errorf("facet-option defaults = scope %s facet_path %s, want metric / ed.tag",
			show(got[1].Scope), show(got[1].FacetPath))
	}
	// A `facet` variable with no scope must not invent one.
	if got[2].Scope != nil {
		t.Errorf("facet entry scope = %s, want it absent", show(got[2].Scope))
	}
	if got[2].Query != nil {
		t.Errorf("facet entry query = %s, want it absent", show(got[2].Query))
	}
}

func TestResourceAccessesToleratesJunk(t *testing.T) {
	for name, input := range map[string]any{
		"nil":              nil,
		"not an object":    []any{1, 2},
		"no widgets":       map[string]any{"version": float64(4)},
		"widgets not list": map[string]any{"widgets": "nope"},
		"widget not map":   map[string]any{"widgets": []any{"nope", float64(1)}},
		"visuals not list": map[string]any{"widgets": []any{map[string]any{"type": "viz", "visuals": "nope"}}},
		"params not map": map[string]any{"widgets": []any{map[string]any{"type": "viz",
			"visuals": []any{map[string]any{"dataSource": map[string]any{"type": "log", "params": "nope"}}}}}},
	} {
		t.Run(name, func(t *testing.T) {
			// Must not panic; content is unimportant since the input is malformed.
			_ = ResourceAccesses(input)
		})
	}
}

func samePtr(a, b *string) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

func show(s *string) string {
	if s == nil {
		return "<absent>"
	}
	return `"` + *s + `"`
}

func render(accesses []ResourceAccess) string {
	b, err := json.Marshal(accesses)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

// fixtureExists keeps the oracle and the definition fixtures from silently drifting apart.
func TestEveryOracleEntryHasAFixture(t *testing.T) {
	raw, err := os.ReadFile(oraclePath)
	if err != nil {
		t.Fatal(err)
	}
	var oracle map[string][]ResourceAccess
	if err := json.Unmarshal(raw, &oracle); err != nil {
		t.Fatal(err)
	}
	for id := range oracle {
		if _, err := os.Stat(filepath.Join("testdata", "definitions", id+".json")); err != nil {
			t.Errorf("oracle has %q but there is no matching definition fixture", id)
		}
	}
}
