package dashboards

import (
	"slices"
	"strings"
	"testing"
)

func TestSchemaOptionsReadsEveryField(t *testing.T) {
	opts, err := SchemaOptions()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fields := map[string][]string{
		"widgetTypes":     opts.WidgetTypes,
		"visualizerTypes": opts.VisualizerTypes,
		"dataTypes":       opts.DataTypes,
		"variableTypes":   opts.VariableTypes,
		"positionTypes":   opts.PositionTypes,
		"resultTypes":     opts.ResultTypes,
	}
	for name, values := range fields {
		if len(values) == 0 {
			t.Errorf("%s is empty; the schema shape it reads must have changed", name)
		}
		if !slices.IsSorted(values) {
			t.Errorf("%s is not sorted: %v", name, values)
		}
		for _, v := range values {
			if v == "" {
				t.Errorf("%s contains an empty value: %v", name, values)
			}
		}
	}

	if opts.Version != SupportedVersion {
		t.Errorf("Version = %d, want %d", opts.Version, SupportedVersion)
	}

	// Spot-check values that must be present, so a mis-resolved $ref that still yields a
	// non-empty list is caught. `raw-table` in particular lives in a separate variant
	// from the other visualizers and is only found by recursing into the union.
	for _, want := range []struct {
		field  string
		values []string
		expect string
	}{
		{"widgetTypes", opts.WidgetTypes, "viz"},
		{"widgetTypes", opts.WidgetTypes, "grid"},
		{"visualizerTypes", opts.VisualizerTypes, "bignumber"},
		{"visualizerTypes", opts.VisualizerTypes, "raw-table"},
		{"dataTypes", opts.DataTypes, "metric"},
		{"dataTypes", opts.DataTypes, "formula"},
		{"variableTypes", opts.VariableTypes, "facet-option"},
		{"positionTypes", opts.PositionTypes, "grid"},
		{"resultTypes", opts.ResultTypes, "timeseries"},
	} {
		if !slices.Contains(want.values, want.expect) {
			t.Errorf("%s is missing %q: %v", want.field, want.expect, want.values)
		}
	}
}

// The listed options are only useful if they are exactly what validation accepts. For
// each one, a minimal definition using it must validate, and a value not in the list must
// be rejected.
func TestListedWidgetAndVisualizerTypesAreWhatValidationAccepts(t *testing.T) {
	opts, err := SchemaOptions()
	if err != nil {
		t.Fatal(err)
	}

	for _, vizType := range opts.VisualizerTypes {
		issues, err := ValidateDefinition(vizDefinition(vizType))
		if err != nil {
			t.Fatalf("%s: %v", vizType, err)
		}
		// A listed visualizer may still require extra fields of its own; what must never
		// happen is a complaint about `visualizer/type` itself.
		for _, issue := range issues {
			if strings.HasSuffix(issue.Path, "/visualizer/type") {
				t.Errorf("listed visualizer %q rejected: %s", vizType, issue)
			}
		}
	}

	unlisted := ValidateOrFail(t, vizDefinition("bignumbers"))
	var sawTypeIssue bool
	for _, issue := range unlisted {
		if strings.HasSuffix(issue.Path, "/visualizer/type") {
			sawTypeIssue = true
		}
	}
	if !sawTypeIssue {
		t.Errorf("an unlisted visualizer type was accepted; issues: %v", unlisted)
	}
}

func TestListedDataTypesAreWhatValidationAccepts(t *testing.T) {
	opts, err := SchemaOptions()
	if err != nil {
		t.Fatal(err)
	}

	for _, dataType := range opts.DataTypes {
		definition := vizDefinition("bignumber")
		widgets, _ := definition["widgets"].([]any)
		widget, _ := widgets[1].(map[string]any)
		visuals, _ := widget["visuals"].([]any)
		visual, _ := visuals[0].(map[string]any)
		visual["dataSource"] = map[string]any{"type": dataType}

		for _, issue := range ValidateOrFail(t, definition) {
			if strings.HasSuffix(issue.Path, "/dataSource/type") {
				t.Errorf("listed data type %q rejected: %s", dataType, issue)
			}
		}
	}
}

// ValidateOrFail validates a definition and fails the test on a hard error, returning the
// issues for inspection.
func ValidateOrFail(t *testing.T, definition any) []Issue {
	t.Helper()
	issues, err := ValidateDefinition(definition)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return issues
}

// vizDefinition builds the smallest definition holding one viz widget of the given
// visualizer type.
func vizDefinition(vizType string) map[string]any {
	return map[string]any{
		"version":     float64(4),
		"timeFilters": map[string]any{"lookback": "1h"},
		"widgets": []any{
			map[string]any{
				"id":   "root",
				"type": "grid",
				"grid": "72px / 1fr",
			},
			map[string]any{
				"id":   float64(1),
				"type": "viz",
				"position": map[string]any{
					"type":     "grid",
					"targetId": "root",
					"area":     map[string]any{"column": float64(1), "columnSpan": float64(1), "row": float64(1), "rowSpan": float64(1)},
				},
				"resultType": "aggregate",
				"visualizer": map[string]any{"type": vizType},
				"visuals": []any{
					map[string]any{
						"id":         "A",
						"dataSource": map[string]any{"type": "metric", "params": map[string]any{"query": "sum:m{*}"}},
					},
				},
			},
		},
	}
}
