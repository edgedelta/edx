package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"

	"github.com/edgedelta/edx/internal/api"
)

func newMetricsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Discover and query metrics",
		Long: `Discover metric names and run aggregation queries.

Metric names must be EXACT (no wildcards). Use "edx metrics list" first to
find the metric, then "edx metrics query" to aggregate it.`,
	}
	cmd.AddCommand(newMetricsListCmd(), newMetricsQueryCmd())
	return cmd
}

func newMetricsListCmd() *cobra.Command {
	var keyword string
	var limit int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available metric names",
		Example: `  edx metrics list
  edx metrics list --keyword cpu`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			q.Set("scope", "metric")
			q.Set("facet_path", "name")
			q.Set("limit", itoa(limit))
			data, err := c.Get(cmdContext(cmd), "/facet_options", q)
			if err != nil {
				return err
			}
			if keyword == "" {
				return printResult(data)
			}
			// Client-side keyword filter over the facet options.
			var facet struct {
				Options []map[string]any `json:"options"`
			}
			if err := json.Unmarshal(data, &facet); err != nil {
				return printResult(data)
			}
			filtered := make([]map[string]any, 0)
			kw := strings.ToLower(keyword)
			for _, opt := range facet.Options {
				if name, _ := opt["name"].(string); strings.Contains(strings.ToLower(name), kw) {
					filtered = append(filtered, opt)
					continue
				}
				if value, _ := opt["value"].(string); strings.Contains(strings.ToLower(value), kw) {
					filtered = append(filtered, opt)
				}
			}
			out, err := json.Marshal(map[string]any{"options": filtered})
			if err != nil {
				return err
			}
			return printResult(out)
		},
	}
	cmd.Flags().StringVar(&keyword, "keyword", "", "case-insensitive substring filter on metric names")
	cmd.Flags().IntVar(&limit, "limit", 1000, "maximum number of metric names to fetch")
	return cmd
}

func newMetricsQueryCmd() *cobra.Command {
	var name, agg, filter, graphType string
	var groupBy []string
	var rollup int
	var tf timeFlags
	var pg pageFlags
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Aggregate a metric as a timeseries or table",
		Long: `Aggregate a metric. The CQL built under the hood is:

  <agg>:<name>{<filter>} by {<group-by>}.rollup(<seconds>)

The filter uses CQL field syntax (service.name:"api"); full-text search is not
supported for metrics. Use "*" for no filter.

Only indexed metric dimensions can be grouped on (discover them with
"edx facets keys --scope metric"). Grouping by an OTLP datapoint attribute or a
custom log_to_metric field dimension that is not indexed returns a single empty
group rather than an error — edx warns when a --group-by key is not indexed.

The response is keyed by query id, e.g. {"A": {"records": [...]}}.`,
		Example: `  edx metrics query --name http.request.duration --agg avg --group-by service.name
  edx metrics query --name system.cpu.usage --agg max --filter 'host.name:"web-1"' --lookback 24h
  edx metrics query --name http.requests --agg sum --rollup 300 --graph-type table`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			cql := fmt.Sprintf("%s:%s{%s}", agg, name, filter)
			if len(groupBy) > 0 {
				warnUnindexedGroupBy(cmdContext(cmd), c, groupBy)
				cql += fmt.Sprintf(" by {%s}", strings.Join(groupBy, ","))
			}
			if rollup > 0 {
				cql += fmt.Sprintf(".rollup(%d)", rollup)
			}
			payload := map[string]any{
				"queries":  map[string]any{"A": map[string]any{"scope": "metric", "query": cql}},
				"formulas": map[string]any{"A": map[string]any{"formula": "A"}},
			}
			body, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			q := url.Values{}
			tf.apply(q)
			pg.apply(q)
			if graphType != "" {
				q.Set("graph_type", graphType)
			}
			data, err := c.Post(cmdContext(cmd), "/graph", q, body)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "exact metric name (required; discover with `edx metrics list`)")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&agg, "agg", "sum", "aggregation: sum, avg, min, max, count, median")
	cmd.Flags().StringVar(&filter, "filter", "*", `CQL filter (e.g. 'service.name:"api"'); "*" for none`)
	cmd.Flags().StringSliceVar(&groupBy, "group-by", nil, "group-by keys (e.g. service.name,host.name)")
	cmd.Flags().IntVar(&rollup, "rollup", 0, "rollup period in seconds (default: auto from time range)")
	cmd.Flags().StringVar(&graphType, "graph-type", "timeseries", "graph type: timeseries or table")
	tf.register(cmd, "1h")
	pg.register(cmd, 0)
	return cmd
}

// warnUnindexedGroupBy warns (best-effort) when a --group-by key is not an
// indexed metric dimension, since such a group-by silently returns one empty
// group instead of failing.
func warnUnindexedGroupBy(ctx context.Context, c *api.Client, groupBy []string) {
	q := url.Values{}
	q.Set("scope", "metric")
	data, err := c.Get(ctx, "/facet_keys", q)
	if err != nil {
		return // don't block the query on a discovery failure
	}
	var keys []struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &keys); err != nil {
		return
	}
	indexed := make(map[string]bool, len(keys))
	for _, k := range keys {
		if k.Key != "" {
			indexed[k.Key] = true
		}
		if k.Name != "" {
			indexed[k.Name] = true
		}
	}
	if len(indexed) == 0 {
		return
	}
	for _, g := range groupBy {
		if !indexed[strings.TrimSpace(g)] {
			warnf("group-by key %q is not an indexed metric dimension; results will be ungrouped (one empty group). See `edx facets keys --scope metric`.", g)
		}
	}
}
