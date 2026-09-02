package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/edgedelta/edx/internal/api"
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
  event.domain:"Monitor Alerts"    all monitor-triggered events
  event.domain:"K8s"               Kubernetes events`,
	}
	cmd.AddCommand(newEventsSearchCmd())
	return cmd
}

// eventPage is the /events/search response envelope. next_cursor is always
// present: non-empty means another page exists, "" means this was the last.
type eventPage struct {
	QueryID    string            `json:"query_id"`
	Items      []json.RawMessage `json:"items"`
	NextCursor string            `json:"next_cursor"`
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
  then sets the page size and defaults to 1000 (the server's own default)
  instead of 20. Combine with --output-file for large sweeps.`,
		Example: `  edx events search --query 'event.type:"pattern_anomaly"' --lookback 6h
  edx events search --query 'event.domain:"Monitor Alerts"' --output table
  edx events search --query 'service.name:"api" AND event.type:"pattern_anomaly"'
  edx events search --query 'event.domain:"Monitor Alerts"' --lookback 720h --all --output-file events.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			buildQuery := func(cursor string) url.Values {
				q := url.Values{}
				if query != "" {
					q.Set("query", query)
				}
				tf.apply(q)
				pg.apply(q)
				if cursor != "" {
					q.Set("cursor", cursor)
				}
				return q
			}
			if !all {
				data, err := c.Get(cmdContext(cmd), "/events/search", buildQuery(""))
				if err != nil {
					return err
				}
				// Legibility for one-page reads: say on stderr when the count
				// is a lower bound, since table/csv output hides next_cursor.
				var page eventPage
				if json.Unmarshal(data, &page) == nil && page.NextCursor != "" {
					warnf("more results exist: re-run with --cursor '%s', or use --all to fetch every page", page.NextCursor)
				}
				return printResult(data)
			}
			// --all: --limit is the page size; default to the server's own
			// default page size rather than the CLI's one-page default of 20.
			if !cmd.Flags().Changed("limit") {
				pg.limit = 1000
			}
			merged, err := fetchAllEventPages(cmdContext(cmd), c, buildQuery, pg.cursor)
			if err != nil {
				return err
			}
			return printResult(merged)
		},
	}
	cmd.Flags().StringVarP(&query, "query", "q", "", "CQL query (full-text search supported)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch every page (follow next_cursor until empty) and print one combined response")
	tf.register(cmd, "1h")
	pg.register(cmd, 20)
	return cmd
}

// fetchAllEventPages follows next_cursor until the server reports the result
// set complete, returning a single combined envelope. It fails atomically: any
// mid-sweep error aborts with no output, so a partial sweep can never be
// mistaken for a complete one.
func fetchAllEventPages(ctx context.Context, c *api.Client, buildQuery func(cursor string) url.Values, cursor string) ([]byte, error) {
	var items []json.RawMessage
	pages := 0
	for {
		data, err := c.Get(ctx, "/events/search", buildQuery(cursor))
		if err != nil {
			return nil, fmt.Errorf("page %d: %w", pages+1, err)
		}
		var page eventPage
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("page %d: unexpected response shape: %w", pages+1, err)
		}
		items = append(items, page.Items...)
		pages++
		warnf("page %d: %d item(s), %d total", pages, len(page.Items), len(items))
		if page.NextCursor == "" {
			break
		}
		if page.NextCursor == cursor {
			return nil, fmt.Errorf("page %d: server returned the same cursor twice; aborting to avoid an infinite loop", pages)
		}
		cursor = page.NextCursor
	}
	if items == nil {
		items = []json.RawMessage{}
	}
	return json.Marshal(struct {
		Items      []json.RawMessage `json:"items"`
		Pages      int               `json:"pages"`
		TotalItems int               `json:"total_items"`
		NextCursor string            `json:"next_cursor"`
	}{Items: items, Pages: pages, TotalItems: len(items), NextCursor: ""})
}
