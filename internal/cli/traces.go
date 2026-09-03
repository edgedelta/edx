package cli

import (
	"net/url"

	"github.com/edgedelta/edx/internal/api"

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
	var includeChildren, all bool
	var tf timeFlags
	var pg pageFlags
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search spans with CQL",
		Example: `  edx traces search --query 'status.code:"ERROR"' --lookback 1h
  edx traces search --query 'service.name:"checkout" AND span.kind:"server"' --include-children
  edx traces search --query 'trace_id:"abc123"'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return pagedRequest{
				svc: api.ServiceAPI, path: "/traces", itemsKey: "items",
				all: all, cursor: pg.cursor,
				allLimit: func() { pg.limit = 1000 },
				query: func() url.Values {
					q := url.Values{}
					if query != "" {
						q.Set("query", query)
					}
					tf.apply(q)
					pg.apply(q)
					if includeChildren {
						q.Set("include_child_spans", "true")
					}
					return q
				},
			}.run(cmd)
		},
	}
	cmd.Flags().StringVarP(&query, "query", "q", "", `CQL filter (field:"value" syntax required)`)
	cmd.Flags().BoolVar(&includeChildren, "include-children", false, "include child spans of matched spans for full trace context")
	registerAllFlag(cmd, &all)
	tf.register(cmd, "1h")
	pg.register(cmd, 20)
	return cmd
}
