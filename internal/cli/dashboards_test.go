package cli

import (
	"encoding/json"
	"testing"
)

// The point of generating resource_accesses is that an author never writes it, so what
// matters is what ends up in the request body.
func TestFillResourceAccessesDerivesTheAllowlist(t *testing.T) {
	body := []byte(`{
      "dashboard_name": "X",
      "definition": {
        "version": 4,
        "widgets": [
          {"id": "root", "type": "grid", "grid": "72px / 1fr"},
          {"id": 1, "type": "viz", "visuals": [
            {"id": "A", "dataSource": {"type": "metric", "params": {"query": "sum:a{*}"}}}
          ]}
        ],
        "variables": [
          {"id": 1, "type": "facet-option", "key": "k", "label": "L",
           "params": {"facet": "host.name", "scope": "log"}}
        ]
      }
    }`)

	filled, err := fillResourceAccesses(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got struct {
		DashboardName    string `json:"dashboard_name"`
		ResourceAccesses []struct {
			Domain    string  `json:"domain"`
			Query     *string `json:"query"`
			Scope     *string `json:"scope"`
			FacetPath *string `json:"facet_path"`
		} `json:"resource_accesses"`
		Definition map[string]any `json:"definition"`
	}
	if err := json.Unmarshal(filled, &got); err != nil {
		t.Fatal(err)
	}

	if got.DashboardName != "X" {
		t.Errorf("dashboard_name = %q, want it preserved", got.DashboardName)
	}
	if got.Definition == nil {
		t.Error("definition was dropped")
	}
	if len(got.ResourceAccesses) != 2 {
		t.Fatalf("got %d entries, want 2: %s", len(got.ResourceAccesses), filled)
	}
	if got.ResourceAccesses[0].Domain != "metric" || *got.ResourceAccesses[0].Query != "sum:a{*}" {
		t.Errorf("entry 0 = %+v, want the metric query", got.ResourceAccesses[0])
	}
	if a := got.ResourceAccesses[1]; a.Domain != "facet_option" || *a.Scope != "log" || *a.FacetPath != "host.name" {
		t.Errorf("entry 1 = %+v, want the facet-option variable", a)
	}
}

// Whatever was in the body is replaced, so a stale allowlist copied from
// `edx dashboards get` cannot outlive the definition it described.
func TestFillResourceAccessesReplacesExistingEntries(t *testing.T) {
	body := []byte(`{
      "dashboard_name": "X",
      "resource_accesses": [{"domain": "metric", "query": "stale"}],
      "definition": {"widgets": [
        {"id": 1, "type": "viz", "visuals": [
          {"id": "A", "dataSource": {"type": "log", "params": {"query": "{a:b}"}}}
        ]}
      ]}
    }`)

	filled, err := fillResourceAccesses(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got struct {
		ResourceAccesses []struct {
			Domain string `json:"domain"`
			Query  string `json:"query"`
		} `json:"resource_accesses"`
	}
	if err := json.Unmarshal(filled, &got); err != nil {
		t.Fatal(err)
	}
	if len(got.ResourceAccesses) != 1 {
		t.Fatalf("got %d entries, want 1: %s", len(got.ResourceAccesses), filled)
	}
	if got.ResourceAccesses[0].Query == "stale" {
		t.Error("the stale entry survived")
	}
	if got.ResourceAccesses[0].Domain != "log" {
		t.Errorf("domain = %q, want log", got.ResourceAccesses[0].Domain)
	}
}

// A definition with nothing to allow must send [] rather than null.
func TestFillResourceAccessesEmitsAnEmptyListNotNull(t *testing.T) {
	body := []byte(`{"dashboard_name": "X", "definition": {"widgets": [
      {"id": "root", "type": "grid", "grid": "72px / 1fr"}
    ]}}`)

	filled, err := fillResourceAccesses(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(filled, &raw); err != nil {
		t.Fatal(err)
	}
	if string(raw["resource_accesses"]) != "[]" {
		t.Errorf("resource_accesses = %s, want []", raw["resource_accesses"])
	}
}

// `validate` accepts a bare definition; create and update always send a full body. A body
// with no "definition" key is left exactly as it was.
func TestFillResourceAccessesPassesThroughABareDefinition(t *testing.T) {
	body := []byte(`{"version":4,"widgets":[]}`)
	filled, err := fillResourceAccesses(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(filled) != string(body) {
		t.Errorf("body = %s, want it untouched", filled)
	}
}

func TestFillResourceAccessesRejectsInvalidJSON(t *testing.T) {
	if _, err := fillResourceAccesses([]byte(`{"definition":`)); err == nil {
		t.Error("expected an error for malformed JSON")
	}
}
