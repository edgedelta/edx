package cli

import (
	"net/url"

	"github.com/spf13/cobra"
)

func newTracesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "traces",
		Short: "Search distributed traces (OTel spans)",
		Long: `Search distributed traces. Queries require CQL field syntax
(service.name:"api"); full-text search is NOT supported for traces.

Common fields: service.name, status.code, span.kind, trace_id, ed.tag`,
	}
	cmd.AddCommand(newTracesSearchCmd())
	return cmd
}

func newTracesSearchCmd() *cobra.Command {
	var query string
	var includeChildren bool
	var tf timeFlags
	var pg pageFlags
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search spans with CQL",
		Example: `  edx traces search --query 'status.code:"ERROR"' --lookback 1h
  edx traces search --query 'service.name:"checkout" AND span.kind:"server"' --include-children
  edx traces search --query 'trace_id:"abc123"'`,
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
			if includeChildren {
				q.Set("include_child_spans", "true")
			}
			data, err := c.Get(cmdContext(cmd), "/traces", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&query, "query", "q", "", `CQL filter (field:"value" syntax required)`)
	cmd.Flags().BoolVar(&includeChildren, "include-children", false, "include child spans of matched spans for full trace context")
	tf.register(cmd, "1h")
	pg.register(cmd, 20)
	return cmd
}
