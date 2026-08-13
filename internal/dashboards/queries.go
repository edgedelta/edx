package dashboards

import (
	"fmt"
	"strconv"

	"github.com/edgedelta/edx/internal/cql"
)

// QuerySite is one query embedded in a dashboard definition.
type QuerySite struct {
	// Path is a JSON Pointer to the query string itself, so an error can point at the
	// field the author wrote, e.g. "/widgets/2/visuals/0/dataSource/params/query".
	Path    string
	Dialect cql.Dialect
	Query   string
}

// dataSourceDialects maps a data source's `type` to the grammar the backend parses its
// query with. The log-family types all share EDCqlLogParser: the backend reaches them
// through pkg/antlrcql/log (chcommon for logs and events, tracerepo for traces,
// cluster/clickhouse for patterns).
//
// A type that is absent here has no query to check: "empty" carries none, and "formula"
// is handled separately because its query lives in params.formula.
var dataSourceDialects = map[string]cql.Dialect{
	"log":     cql.DialectLog,
	"event":   cql.DialectLog,
	"pattern": cql.DialectLog,
	"trace":   cql.DialectLog,
	"metric":  cql.DialectMetric,
}

// QuerySites returns every query in a definition, in document order.
//
// The walk is shape-aware rather than a blind search for "query" keys, so each site gets
// the right dialect and an exact path. TestQuerySitesFindsEveryQueryInFixtures guards
// against it falling behind the schema.
func QuerySites(definition any) []QuerySite {
	root, ok := definition.(map[string]any)
	if !ok {
		return nil
	}

	var sites []QuerySite

	for i, widget := range array(root["widgets"]) {
		w, ok := widget.(map[string]any)
		if !ok {
			continue
		}
		widgetPath := fmt.Sprintf("/widgets/%d", i)

		for j, visual := range array(w["visuals"]) {
			v, ok := visual.(map[string]any)
			if !ok {
				continue
			}
			visualPath := widgetPath + "/visuals/" + strconv.Itoa(j)

			primary, _ := v["dataSource"].(map[string]any)
			sites = append(sites, dataSourceSites(visualPath+"/dataSource", primary, "")...)

			// compareDataSource is a Partial<DataSource> the frontend merges over the
			// primary with `{...dataSource, ...compareDataSource}`, so it inherits the
			// type when it omits one. The spread is shallow: a compare `params` replaces
			// the primary's outright, and no `params` means no query of its own.
			if compare, ok := v["compareDataSource"].(map[string]any); ok {
				inherited, _ := primary["type"].(string)
				sites = append(sites, dataSourceSites(visualPath+"/compareDataSource", compare, inherited)...)
			}
		}
	}

	for i, variable := range array(root["variables"]) {
		v, ok := variable.(map[string]any)
		if !ok {
			continue
		}
		path := "/variables/" + strconv.Itoa(i)
		varType, _ := v["type"].(string)

		// A `query` variable holds a filter expression that gets substituted into the
		// widgets' queries.
		if varType == "query" {
			if value, ok := v["value"].(string); ok {
				sites = append(sites, QuerySite{Path: path + "/value", Dialect: cql.DialectLog, Query: value})
			}
		}

		// `facet-option` and `metric-name` narrow their option list with a filter, which
		// the backend parses through monitorvisitor — EDCqlLogParser again.
		if params, ok := v["params"].(map[string]any); ok {
			if query, ok := params["query"].(string); ok {
				sites = append(sites, QuerySite{Path: path + "/params/query", Dialect: cql.DialectLog, Query: query})
			}
		}
	}

	return sites
}

// dataSourceSites reads the query out of one data source. inheritedType applies when the
// source omits its own `type`, which only a compareDataSource does.
func dataSourceSites(path string, source map[string]any, inheritedType string) []QuerySite {
	if source == nil {
		return nil
	}

	sourceType, ok := source["type"].(string)
	if !ok {
		sourceType = inheritedType
	}

	params, ok := source["params"].(map[string]any)
	if !ok {
		return nil
	}

	if sourceType == "formula" {
		formula, ok := params["formula"].(string)
		if !ok {
			return nil
		}
		return []QuerySite{{Path: path + "/params/formula", Dialect: cql.DialectFormula, Query: formula}}
	}

	dialect, known := dataSourceDialects[sourceType]
	if !known {
		return nil
	}
	query, ok := params["query"].(string)
	if !ok {
		return nil
	}
	return []QuerySite{{Path: path + "/params/query", Dialect: dialect, Query: query}}
}

// validateQueries syntax-checks every query in a definition.
func validateQueries(definition any) ([]Issue, error) {
	var out []Issue
	for _, site := range QuerySites(definition) {
		errs, err := cql.Validate(site.Dialect, site.Query)
		if err != nil {
			return nil, err
		}
		for _, e := range errs {
			out = append(out, Issue{
				Path:    site.Path,
				Message: fmt.Sprintf("invalid %s query %s", site.Dialect, e),
			})
		}
	}
	return out, nil
}

func array(v any) []any {
	items, _ := v.([]any)
	return items
}
