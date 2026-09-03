package cli

import (
	"net/url"

	"github.com/edgedelta/edx/internal/api"

	"github.com/spf13/cobra"
)

func newEventsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "events",
		Short: "Search events (anomalies, alerts, K8s events)",
		Long: `Search Edge Delta events: pattern anomalies, monitor alerts and
Kubernetes events.

Common queries:
  event.type:"pattern_anomaly"     log anomaly detections
  event.type:"metric_threshold"    metric alert triggers
  event.type:"log_threshold"       log alert triggers
  event.domain:("Monitor" OR "Monitor Alerts")   all monitor-triggered events
  event.domain:"K8s"               Kubernetes events

Domain values vary by org and environment; list the live set with
"edx facets options --scope event --facet event.domain".`,
	}
	cmd.AddCommand(newEventsSearchCmd())
	return cmd
}

func newEventsSearchCmd() *cobra.Command {
	var query string
	var all bool
	var tf timeFlags
	var pg pageFlags
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search events with CQL",
		Long: `Search events with CQL.

PAGINATION
  A page carries up to --limit items plus "next_cursor": non-empty means more
  results exist (re-run with --cursor <next_cursor> and the same query/time
  flags), "" means the result set is complete. A count from a single page is a
  lower bound unless next_cursor came back empty.

  For a full sweep (e.g. 30 days of monitor alerts) pass --all: edx follows
  next_cursor until it is empty and prints one combined response
  ({items, pages, total_items}), with per-page progress on stderr. --limit
  then sets the page size and defaults to 1000 instead of 20. Combine with
  --output-file for large sweeps. The same flags work on every
  cursor-paginated command (logs search, traces search, monitors list/states,
  rehydrations list).`,
		Example: `  edx events search --query 'event.type:"pattern_anomaly"' --lookback 6h
  edx events search --query 'event.domain:("Monitor" OR "Monitor Alerts")' --output table
  edx events search --query 'service.name:"api" AND event.type:"pattern_anomaly"'
  edx events search --query 'event.domain:("Monitor" OR "Monitor Alerts")' --lookback 720h --all --output-file events.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return pagedRequest{
				svc: api.ServiceAPI, path: "/events/search", itemsKey: "items",
				all: all, cursor: pg.cursor,
				allLimit: func() { pg.limit = 1000 },
				query: func() url.Values {
					q := url.Values{}
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
	cmd.Flags().StringVarP(&query, "query", "q", "", "CQL query (full-text search supported)")
	registerAllFlag(cmd, &all)
	tf.register(cmd, "1h")
	pg.register(cmd, 20)
	return cmd
}
