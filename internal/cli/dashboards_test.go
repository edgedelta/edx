package cli

import (
	"encoding/json"
	"strings"
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

// A tag is the only thing --tag may change: the definition and everything else in the body
// has to arrive at the API the way the author wrote it.
func TestApplyTagsAddsWithoutDisturbingTheBody(t *testing.T) {
	body := []byte(`{"dashboard_name":"X","tags":["team-infra"],"definition":{"version":4,"widgets":[{"id":1000000}]}}`)

	tagged, err := applyTags(body, []string{"generated", "preview"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got struct {
		DashboardName string   `json:"dashboard_name"`
		Tags          []string `json:"tags"`
	}
	if err := json.Unmarshal(tagged, &got); err != nil {
		t.Fatal(err)
	}
	if got.DashboardName != "X" {
		t.Errorf("dashboard_name = %q, want it preserved", got.DashboardName)
	}
	want := []string{"team-infra", "generated", "preview"}
	if len(got.Tags) != len(want) {
		t.Fatalf("tags = %v, want %v", got.Tags, want)
	}
	for i := range want {
		if got.Tags[i] != want[i] {
			t.Errorf("tags = %v, want %v (existing tags first, in order)", got.Tags, want)
			break
		}
	}
	// A widget ID that round-tripped through float64 would come back as 1e+06 and the
	// dashboard would no longer match its own definition.
	if !strings.Contains(string(tagged), `"id":1000000`) {
		t.Errorf("numeric literals were rewritten: %s", tagged)
	}
}

func TestApplyTagsRemoves(t *testing.T) {
	body := []byte(`{"tags":["generated","preview","team-infra"]}`)

	tagged, err := applyTags(body, nil, []string{"preview"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got struct {
		Tags []string `json:"tags"`
	}
	if err := json.Unmarshal(tagged, &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Tags) != 2 || got.Tags[0] != "generated" || got.Tags[1] != "team-infra" {
		t.Errorf("tags = %v, want [generated team-infra]", got.Tags)
	}
}

// Removing the last tag has to drop the key, not leave an empty array behind that reads as
// "deliberately untagged".
func TestApplyTagsDropsTheKeyWhenNothingIsLeft(t *testing.T) {
	tagged, err := applyTags([]byte(`{"tags":["preview"]}`), nil, []string{"preview"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(tagged, &raw); err != nil {
		t.Fatal(err)
	}
	if _, ok := raw["tags"]; ok {
		t.Errorf("tags survived as %s, want the key gone", raw["tags"])
	}
}

// Re-running a create with the same flags must not accumulate duplicates.
func TestApplyTagsIsIdempotent(t *testing.T) {
	body := []byte(`{"tags":["generated"]}`)
	once, err := applyTags(body, []string{"generated"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	twice, err := applyTags(once, []string{"generated"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(once) != string(twice) {
		t.Errorf("not idempotent: %s then %s", once, twice)
	}
	var got struct {
		Tags []string `json:"tags"`
	}
	if err := json.Unmarshal(twice, &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Tags) != 1 {
		t.Errorf("tags = %v, want one entry", got.Tags)
	}
}

// Without --tag the body must be handed on byte for byte, so nothing is reformatted on a
// path that was not asked to change anything.
func TestApplyTagsWithNoChangesReturnsTheBodyUntouched(t *testing.T) {
	body := []byte(`{"dashboard_name":  "X"}`)
	same, err := applyTags(body, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(same) != string(body) {
		t.Errorf("body = %s, want it untouched", same)
	}
}

func TestApplyTagsTrimsAndIgnoresBlanks(t *testing.T) {
	tagged, err := applyTags([]byte(`{}`), []string{" generated ", "", "   "}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got struct {
		Tags []string `json:"tags"`
	}
	if err := json.Unmarshal(tagged, &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Tags) != 1 || got.Tags[0] != "generated" {
		t.Errorf("tags = %v, want [generated]", got.Tags)
	}
}

func TestApplyTagsRejectsInvalidJSON(t *testing.T) {
	if _, err := applyTags([]byte(`{"tags":`), []string{"generated"}, nil); err == nil {
		t.Error("expected an error for malformed JSON")
	}
}

// A tags array with junk in it must not crash the merge or leak the junk back.
func TestApplyTagsIgnoresNonStringTags(t *testing.T) {
	tagged, err := applyTags([]byte(`{"tags":["keep",7,null,{"a":1}]}`), []string{"generated"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got struct {
		Tags []string `json:"tags"`
	}
	if err := json.Unmarshal(tagged, &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Tags) != 2 || got.Tags[0] != "keep" || got.Tags[1] != "generated" {
		t.Errorf("tags = %v, want [keep generated]", got.Tags)
	}
}
