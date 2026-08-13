package dashboards

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// loadDefinition reads a definition fixture. The fixtures under testdata/definitions
// are the dashboards shipped by the frontend, typed against DashboardsV4.Definition —
// so any rejection here is a false positive in the schema, not a bad fixture.
func loadDefinition(t *testing.T, name string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", "definitions", name))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var def map[string]any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	// Decode through the same path the CLI uses, so number handling matches.
	if err := json.Unmarshal(raw, &def); err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	return def
}

func TestValidateDefinitionAcceptsShippedDashboards(t *testing.T) {
	entries, err := os.ReadDir(filepath.Join("testdata", "definitions"))
	if err != nil {
		t.Fatalf("read fixtures: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("no fixtures found")
	}
	for _, e := range entries {
		t.Run(e.Name(), func(t *testing.T) {
			issues, err := ValidateDefinition(loadDefinition(t, e.Name()))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(issues) > 0 {
				for _, iss := range issues {
					t.Errorf("false positive: %s", iss)
				}
			}
		})
	}
}

func TestValidateDefinitionRejects(t *testing.T) {
	vizWidget := func(def map[string]any) map[string]any {
		for _, w := range def["widgets"].([]any) {
			m := w.(map[string]any)
			if m["type"] == "viz" {
				return m
			}
		}
		t.Fatal("no viz widget in fixture")
		return nil
	}

	cases := map[string]struct {
		mutate   func(def map[string]any)
		wantPath string
	}{
		"unknown viz type": {
			mutate:   func(d map[string]any) { vizWidget(d)["visualizer"].(map[string]any)["type"] = "lyne" },
			wantPath: "visualizer/type",
		},
		"excess property on a widget": {
			mutate:   func(d map[string]any) { vizWidget(d)["notAWidgetField"] = "x" },
			wantPath: "widgets/",
		},
		"excess property on displayOptions": {
			mutate: func(d map[string]any) {
				vizWidget(d)["displayOptions"].(map[string]any)["typo"] = true
			},
			wantPath: "displayOptions",
		},
		"excess property at the top level": {
			mutate:   func(d map[string]any) { d["unknownTopLevel"] = 1 },
			wantPath: "",
		},
		"missing widgets": {
			mutate:   func(d map[string]any) { delete(d, "widgets") },
			wantPath: "",
		},
		"widget without an id": {
			mutate:   func(d map[string]any) { delete(vizWidget(d), "id") },
			wantPath: "widgets/",
		},
		"unknown result type": {
			mutate:   func(d map[string]any) { vizWidget(d)["resultType"] = "timeserie" },
			wantPath: "widgets/",
		},
		"data source id outside A-F/W-Z": {
			mutate: func(d map[string]any) {
				vizWidget(d)["visuals"].([]any)[0].(map[string]any)["id"] = "Q"
			},
			wantPath: "visuals/0/id",
		},
		"unknown aggregation mode": {
			mutate: func(d map[string]any) {
				vizWidget(d)["visualizer"].(map[string]any)["aggregation"] = map[string]any{"mode": "total"}
			},
			wantPath: "aggregation/mode",
		},
		"unknown facet scope inside params": {
			mutate: func(d map[string]any) {
				d["variables"] = []any{map[string]any{
					"id": 1, "label": "L", "key": "k", "type": "query",
					"params": map[string]any{"scope": "logz"},
				}}
			},
			wantPath: "variables/0",
		},
		"markdown params content as a number": {
			mutate: func(d map[string]any) {
				d["widgets"] = append(d["widgets"].([]any), map[string]any{
					"id": 998, "type": "markdown", "params": map[string]any{"content": 42},
				})
			},
			wantPath: "widgets/",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			def := loadDefinition(t, "top-talkers.json")
			tc.mutate(def)

			issues, err := ValidateDefinition(def)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(issues) == 0 {
				t.Fatal("expected the definition to be rejected")
			}
			if tc.wantPath != "" {
				var found bool
				for _, iss := range issues {
					if strings.Contains(iss.Path, tc.wantPath) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("no issue mentioned %q; got:", tc.wantPath)
					for _, iss := range issues {
						t.Errorf("  %s", iss)
					}
				}
			}
		})
	}
}

// Without discriminant narrowing a single bad field reports against all six widget
// branches ("missing property 'grid'", "additional properties ... not allowed", ...).
func TestValidateDefinitionReportsOnlyTheIntendedVariant(t *testing.T) {
	def := loadDefinition(t, "top-talkers.json")
	var viz map[string]any
	for _, w := range def["widgets"].([]any) {
		if m := w.(map[string]any); m["type"] == "viz" {
			viz = m
			break
		}
	}
	viz["displayOptions"].(map[string]any)["tittle"] = "typo"

	issues, err := ValidateDefinition(def)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 1 {
		t.Errorf("want exactly 1 issue, got %d:", len(issues))
		for _, iss := range issues {
			t.Errorf("  %s", iss)
		}
		return
	}
	if !strings.Contains(issues[0].Path, "displayOptions") || !strings.Contains(issues[0].Message, "tittle") {
		t.Errorf("issue does not point at the typo: %s", issues[0])
	}
}

// The params escape hatch must keep accepting keys the schema does not know about.
func TestValidateDefinitionAllowsUnknownParamsKeys(t *testing.T) {
	def := loadDefinition(t, "top-talkers.json")
	def["widgets"] = append(def["widgets"].([]any), map[string]any{
		"id": 998, "type": "markdown",
		"params": map[string]any{"content": "x", "custom": map[string]any{"a": 1}},
	})
	def["variables"] = []any{map[string]any{
		"id": 1, "label": "L", "key": "k", "type": "query",
		"params": map[string]any{"scope": "log", "extra": true},
	}}

	issues, err := ValidateDefinition(def)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, iss := range issues {
		t.Errorf("false positive: %s", iss)
	}
}

func TestValidateDefinitionRejectsOtherVersions(t *testing.T) {
	def := loadDefinition(t, "top-talkers.json")
	def["version"] = float64(3)

	if _, err := ValidateDefinition(def); !errors.Is(err, ErrUnsupportedVersion) {
		t.Fatalf("want ErrUnsupportedVersion, got %v", err)
	}
}

func TestEmbeddedSchemaCompiles(t *testing.T) {
	if _, err := schema(); err != nil {
		t.Fatalf("embedded schema does not compile: %v", err)
	}
}
