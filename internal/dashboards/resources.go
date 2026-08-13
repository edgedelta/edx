package dashboards

// resource_accesses generation.
//
// resource_accesses is the allowlist of what an anonymous viewer may query. The backend
// consults it only for dashboard-token auth — public share links and screenshot
// generation — comparing `domain`, `scope` and `facet_path`, never `query`.
//
// The UI derives it from the definition on every save (fillDashboardResources ->
// each widget's and variable's `resolveResources`), so it is not something an author
// writes. ResourceAccesses mirrors that derivation so edx can fill it in too, which keeps
// it out of the hands of whoever is authoring the definition.
//
// Fidelity is checked against the real frontend implementation: testdata/resource-accesses.json
// is the output of fillDashboardResources() run over the shipped dashboards, and
// TestResourceAccessesMatchesTheFrontend diffs against it.

// ResourceAccess is one entry of a dashboard's resource_accesses.
//
// The optional fields are pointers because the frontend distinguishes an empty query from
// an absent one: a log widget with no query emits `"query": ""`, while a `facet` variable
// emits no query key at all.
type ResourceAccess struct {
	Domain    string  `json:"domain"`
	Query     *string `json:"query,omitempty"`
	Scope     *string `json:"scope,omitempty"`
	FacetPath *string `json:"facet_path,omitempty"`
}

// edTag is the frontend's default facet for a facet-option variable (semconv1.EDTag).
const edTag = "ed.tag"

// defaultFacetScope is the frontend's default scope for a facet-option variable.
const defaultFacetScope = "metric"

// queryDomains are the data source types that contribute a `{domain, query}` entry, keyed
// by the domain name the entry carries. Every one of these uses its own type as the
// domain. `empty` and `formula` sources contribute nothing: the former has no query, and
// the latter reads other visuals rather than querying a data source.
var queryDomains = map[string]bool{
	"log":     true,
	"metric":  true,
	"event":   true,
	"pattern": true,
	"trace":   true,
}

// ResourceAccesses derives a definition's resource_accesses, mirroring what the UI writes
// when a dashboard is saved.
//
// Entries are returned in the UI's order — widgets first, in document order, then
// variables — and duplicates are kept rather than collapsed, because the same query used
// by two widgets produces two entries there too.
//
// One known difference: for an `event` data source the UI rewrites the query to conjoin an
// `event.domain` precondition, e.g. `{a:b} by {c}` becomes
// `{(a:b) AND (event.domain:"K8s")} by {c}`. This emits the widget's query unchanged.
// Authorization compares only `domain`, and the sole reader of `query` (the metrics
// inventory's usage panel) filters to `domain == "metric"`, so nothing depends on the
// difference — and the UI rewrites the whole array on the next save anyway.
func ResourceAccesses(definition any) []ResourceAccess {
	root, ok := definition.(map[string]any)
	if !ok {
		return nil
	}

	var out []ResourceAccess

	for _, w := range array(root["widgets"]) {
		widget, ok := w.(map[string]any)
		if !ok {
			continue
		}
		// Only viz widgets resolve resources; grid, markdown, tabs, variable-control and
		// empty contribute nothing.
		if kind, _ := widget["type"].(string); kind != "viz" {
			continue
		}
		for _, v := range array(widget["visuals"]) {
			visual, ok := v.(map[string]any)
			if !ok {
				continue
			}
			// compareDataSource is deliberately skipped: the UI's viz manifest resolves
			// only `visuals[].dataSource`.
			source, ok := visual["dataSource"].(map[string]any)
			if !ok {
				continue
			}
			sourceType, _ := source["type"].(string)
			if !queryDomains[sourceType] {
				continue
			}
			params, _ := source["params"].(map[string]any)
			query, _ := params["query"].(string)
			out = append(out, ResourceAccess{Domain: sourceType, Query: &query})
		}
	}

	for _, v := range array(root["variables"]) {
		variable, ok := v.(map[string]any)
		if !ok {
			continue
		}
		params, _ := variable["params"].(map[string]any)

		switch varType, _ := variable["type"].(string); varType {
		case "facet-option":
			// `metric-name` shares this editor but has no resolveResources, so it is
			// intentionally absent here.
			query, _ := params["query"].(string)
			scope := stringOr(params["scope"], defaultFacetScope)
			facet := stringOr(params["facet"], edTag)
			out = append(out, ResourceAccess{
				Domain:    "facet_option",
				Query:     &query,
				Scope:     &scope,
				FacetPath: &facet,
			})
		case "facet":
			// Carries only a scope, and no default: the UI passes params.scope straight
			// through, so an absent scope stays absent.
			access := ResourceAccess{Domain: "facet"}
			if scope, ok := params["scope"].(string); ok {
				access.Scope = &scope
			}
			out = append(out, access)
		}
	}

	return out
}

func stringOr(v any, fallback string) string {
	if s, ok := v.(string); ok && s != "" {
		return s
	}
	return fallback
}
