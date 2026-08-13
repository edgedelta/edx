package dashboards

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// examplesDir holds the dashboard examples shipped with the ed-dashboards skill, vendored
// from the agent-skills repo by `make sync-skills`.
const examplesDir = "../skills/data/ed-dashboards/examples"

// Agents copy these examples as a starting point, so an example that does not validate is
// worse than no example at all: it teaches a shape the product rejects. This is the guard
// that keeps them true as the schema changes.
func TestSkillExamplesValidate(t *testing.T) {
	paths, err := filepath.Glob(filepath.Join(examplesDir, "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) == 0 {
		t.Fatalf("no examples found in %s; if they moved, update this test rather than deleting it", examplesDir)
	}

	for _, path := range paths {
		t.Run(filepath.Base(path), func(t *testing.T) {
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}

			// The examples are full dashboard bodies so they can be fed straight to
			// `edx dashboards create`, so reach through to the definition.
			var body struct {
				DashboardName    string         `json:"dashboard_name"`
				Tags             []string       `json:"tags"`
				Definition       map[string]any `json:"definition"`
				ResourceAccesses []struct {
					Domain string `json:"domain"`
					Query  string `json:"query"`
				} `json:"resource_accesses"`
			}
			if err := json.Unmarshal(raw, &body); err != nil {
				t.Fatalf("invalid JSON: %v", err)
			}
			if body.DashboardName == "" {
				t.Error("example has no dashboard_name, so it cannot be created as-is")
			}
			// The examples are copied verbatim as starting points, so they are where the
			// tag convention is actually taught: an example that ships untagged teaches
			// authors to leave generated dashboards anonymous.
			if !slices.Contains(body.Tags, "generated") {
				t.Errorf(`example tags = %v, want "generated" among them`, body.Tags)
			}
			if !slices.Contains(body.Tags, "preview") {
				t.Errorf(`example tags = %v, want "preview" among them: an example is `+
					"copied before anyone has looked at the result", body.Tags)
			}
			if body.Definition == nil {
				t.Fatal(`example has no "definition" object`)
			}

			issues, err := ValidateDefinition(body.Definition)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			for _, issue := range issues {
				t.Errorf("%s", issue)
			}

			// edx derives resource_accesses from the definition on create and update, so
			// an example must not carry one. Shipping a hand-written copy would teach
			// authors to maintain a field they should never touch, and it would go stale
			// the moment someone edited the widgets.
			if len(body.ResourceAccesses) > 0 {
				t.Errorf("example declares %d resource_accesses entries; edx generates them, "+
					"so remove the field", len(body.ResourceAccesses))
			}
		})
	}
}

// The skill lists every accepted widget, visualizer, data, variable, position and result
// type so an agent knows its options without a CLI call. Those lists are prose, so this
// asserts they still match the schema: adding a visualizer type to definition.ts and
// regenerating must not silently leave the skill a value short.
//
// `edx dashboards options` prints the same values from the same source, which is the fix
// when this fails — copy them across.
func TestSkillDocumentsEverySchemaOption(t *testing.T) {
	raw, err := os.ReadFile("../skills/data/ed-dashboards/SKILL.md")
	if err != nil {
		t.Fatal(err)
	}
	skill := string(raw)

	opts, err := SchemaOptions()
	if err != nil {
		t.Fatal(err)
	}

	for _, field := range []struct {
		name   string
		values []string
	}{
		{"widget type", opts.WidgetTypes},
		{"visualizer type", opts.VisualizerTypes},
		{"data type", opts.DataTypes},
		{"variable type", opts.VariableTypes},
		{"position type", opts.PositionTypes},
		{"result type", opts.ResultTypes},
	} {
		for _, value := range field.values {
			// The skill writes every option as inline code, which keeps this from
			// matching the value incidentally in prose.
			if !strings.Contains(skill, "`"+value+"`") {
				t.Errorf("ed-dashboards/SKILL.md does not document the %s %q; "+
					"run `edx dashboards options` and add it", field.name, value)
			}
		}
	}
}
