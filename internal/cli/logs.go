package cli

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

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
"edx facets keys --scope log" and values with "edx facets options".`,
	}
	cmd.AddCommand(newLogsSearchCmd(), newLogsGraphCmd())
	return cmd
}

func newLogsSearchCmd() *cobra.Command {
	var query string
	var tf timeFlags
	var pg pageFlags
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search raw logs with CQL",
		Example: `  edx logs search --query 'severity_text:"ERROR"' --lookback 1h
  edx logs search --query 'service.name:"api" AND error' --from 2026-06-12T00:00:00.000Z --to 2026-06-12T01:00:00.000Z
  edx logs search --query 'k8s.namespace.name:"prod"' --limit 100 --output table --columns timestamp,severity_text,body`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			if query != "" {
				q.Set("query", query)
			}
			tf.apply(q)
			pg.apply(q)
			data, err := c.Get(cmdContext(cmd), "/logs/log_search/search", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&query, "query", "q", "", "CQL query (empty matches all logs)")
	tf.register(cmd, "1h")
	pg.register(cmd, 20)
	return cmd
}

func newLogsGraphCmd() *cobra.Command {
	var query string
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
  edx logs graph --query '*' --lookback 24h`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			cql := query
			if len(groupBy) > 0 {
				cql += fmt.Sprintf(" by {%s}", strings.Join(groupBy, ","))
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
	tf.register(cmd, "1h")
	pg.register(cmd, 0)
	return cmd
}
