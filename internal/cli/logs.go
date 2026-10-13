package cli

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/edgedelta/edx/internal/api"

	"github.com/spf13/cobra"
)

func newLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Search and aggregate logs",
		Long: `Search and aggregate logs using CQL (Common Query Language).

CQL SYNTAX
  field equals       service.name:"api"
  multiple values    service.name:("api" OR "web")
  negation           -severity_text:"DEBUG"
  full-text search   error timeout      (bare words, logs/patterns/events only)
  numeric            @response.code > 400

Regular expressions are NOT supported. Verify field names with
"edx facets keys --scope log" and values with "edx facets options".

AUDIT LOGS
  --scope audit searches the org's audit trail instead: API requests
  (method, path, status.code, user.email, auth.method) and AI teammate
  activity (ed.audit.resource.type:"ai"). Discover its fields with
  "edx facets keys --scope audit".`,
	}
	cmd.AddCommand(newLogsSearchCmd(), newLogsGraphCmd())
	return cmd
}

// logScopeHelp describes --scope for the log commands. The log search
// endpoints serve two tables: the org's logs and its audit trail.
const logScopeHelp = `data scope: log or audit (API requests and AI teammate activity)`

func validateLogScope(scope string) error {
	if scope != "log" && scope != "audit" {
		return fmt.Errorf("invalid --scope %q: must be log or audit", scope)
	}
	return nil
}

func newLogsSearchCmd() *cobra.Command {
	var query, scope string
	var all bool
	var tf timeFlags
	var pg pageFlags
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search raw logs with CQL",
		Example: `  edx logs search --query 'severity_text:"ERROR"' --lookback 1h
  edx logs search --query 'service.name:"api" AND error' --from 2026-06-12T00:00:00.000Z --to 2026-06-12T01:00:00.000Z
  edx logs search --query 'k8s.namespace.name:"prod"' --limit 100 --output table --columns timestamp,severity_text,body
  edx logs search --scope audit --query 'method:"DELETE"' --lookback 24h`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateLogScope(scope); err != nil {
				return err
			}
			return pagedRequest{
				svc: api.ServiceAPI, path: "/logs/log_search/search", itemsKey: "items",
				all: all, cursor: pg.cursor,
				allLimit: func() { pg.limit = 1000 },
				query: func() url.Values {
					q := url.Values{}
					q.Set("scope", scope)
					if query != "" {
						q.Set("query", query)
					}
					tf.apply(q)
					pg.apply(q)
					return q
				},
			}.run(cmd)
		},
	}
	cmd.Flags().StringVarP(&query, "query", "q", "", "CQL query (empty matches all logs)")
	cmd.Flags().StringVar(&scope, "scope", "log", logScopeHelp)
	registerAllFlag(cmd, &all)
	tf.register(cmd, "1h")
	pg.register(cmd, 20)
	return cmd
}

func newLogsGraphCmd() *cobra.Command {
	var query, scope string
	var groupBy []string
	var tf timeFlags
	var pg pageFlags
	cmd := &cobra.Command{
		Use:   "graph",
		Short: "Get a log count timeseries for a CQL filter",
		Long: `Render a timeseries of log counts matching a CQL filter, optionally grouped
by fields. Use this to find spikes and compare services over time.`,
		Example: `  edx logs graph --query 'severity_text:"ERROR"' --group-by service.name --lookback 6h
  edx logs graph --query 'service.name:"api"'
  edx logs graph --query '*' --lookback 24h
  edx logs graph --scope audit --query '*' --group-by user.email --lookback 24h`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateLogScope(scope); err != nil {
				return err
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			cql := query
			if len(groupBy) > 0 {
				cql += fmt.Sprintf(" by {%s}", strings.Join(groupBy, ","))
			}
			if scope == "audit" {
				// The formula /graph endpoint has no audit evaluator; audit
				// timeseries are served by the log search graph endpoint.
				q := url.Values{}
				q.Set("scope", scope)
				q.Set("query", cql)
				tf.apply(q)
				pg.apply(q)
				data, err := c.Get(cmdContext(cmd), "/logs/log_search/graph", q)
				if err != nil {
					return err
				}
				return printResult(data)
			}
			payload := map[string]any{
				"queries":  map[string]any{"Q1": map[string]any{"scope": "log", "query": cql}},
				"formulas": map[string]any{"R1": map[string]any{"formula": "Q1"}},
			}
			body, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			q := url.Values{}
			tf.apply(q)
			pg.apply(q)
			data, err := c.Post(cmdContext(cmd), "/graph", q, body)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&query, "query", "q", "*", `CQL filter query ("*" for all logs)`)
	cmd.Flags().StringSliceVar(&groupBy, "group-by", nil, "group-by fields (e.g. service.name)")
	cmd.Flags().StringVar(&scope, "scope", "log", logScopeHelp)
	tf.register(cmd, "1h")
	pg.register(cmd, 0)
	return cmd
}
