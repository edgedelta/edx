package cli

import (
	"net/url"

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
  event.domain:"Monitor Alerts"    all monitor-triggered events
  event.domain:"K8s"               Kubernetes events`,
	}
	cmd.AddCommand(newEventsSearchCmd())
	return cmd
}

func newEventsSearchCmd() *cobra.Command {
	var query string
	var tf timeFlags
	var pg pageFlags
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search events with CQL",
		Example: `  edx events search --query 'event.type:"pattern_anomaly"' --lookback 6h
  edx events search --query 'event.domain:"Monitor Alerts"' --output table
  edx events search --query 'service.name:"api" AND event.type:"pattern_anomaly"'`,
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
			data, err := c.Get(cmdContext(cmd), "/events/search", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&query, "query", "q", "", "CQL query (full-text search supported)")
	tf.register(cmd, "1h")
	pg.register(cmd, 20)
	return cmd
}
