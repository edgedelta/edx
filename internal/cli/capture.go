package cli

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

// urlTimeFormat matches the backend's core.URLTimeFormat.
const urlTimeFormat = "2006-01-02T15:04:05.000Z"

func newCaptureCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "capture",
		Short: "Live capture: sample data flowing through pipeline nodes",
		Long: `Live capture samples real data before and after pipeline nodes on running
agents - the fastest way to debug a processor or verify data flow.

Workflow:
  1. edx capture start <conf-id> --duration 2m --nodes <node-name>
  2. edx capture status <task-id>        (poll until agents report)
  3. edx capture results <conf-id>       (fetch captured before/after samples)`,
	}
	cmd.AddCommand(
		newCaptureStartCmd(),
		newCaptureTaskCmd(),
		newCaptureStatusCmd(),
		newCaptureResultsCmd(),
	)
	return cmd
}

func newCaptureStartCmd() *cobra.Command {
	var duration time.Duration
	var nodes []string
	var interval string
	var maxItems int
	cmd := &cobra.Command{
		Use:   "start <conf-id>",
		Short: "Start a live capture task on a pipeline",
		Example: `  edx capture start <conf-id> --duration 2m
  edx capture start <conf-id> --duration 5m --nodes mask_pii,route_errors --max-items 50`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if maxItems < 1 || maxItems > 100 {
				return fmt.Errorf("--max-items must be between 1 and 100 (got %d)", maxItems)
			}
			if interval != "" {
				d, err := time.ParseDuration(interval)
				if err != nil {
					return fmt.Errorf("invalid --interval %q: %v", interval, err)
				}
				if d < time.Second {
					return fmt.Errorf("--interval must be at least 1s (got %s)", interval)
				}
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			task := map[string]any{
				"expires_at": time.Now().UTC().Add(duration).Format(urlTimeFormat),
			}
			if len(nodes) > 0 {
				task["nodes_to_capture"] = nodes
			}
			if interval != "" {
				task["polling_interval"] = interval
			}
			if maxItems > 0 {
				task["max_items"] = maxItems
			}
			body, err := json.Marshal(task)
			if err != nil {
				return err
			}
			data, err := c.Post(cmdContext(cmd), "/pipelines/"+url.PathEscape(args[0])+"/capture/task", nil, body)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().DurationVar(&duration, "duration", 2*time.Minute, "how long the capture stays active")
	cmd.Flags().StringSliceVar(&nodes, "nodes", nil, "pipeline node names to capture (default: all)")
	cmd.Flags().StringVar(&interval, "interval", "", "agent polling interval (Go duration, min 1s)")
	cmd.Flags().IntVar(&maxItems, "max-items", 20, "max items to capture per node (max 100)")
	return cmd
}

func newCaptureTaskCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "task <conf-id>",
		Short: "Show the active capture task for a pipeline",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/pipelines/"+url.PathEscape(args[0])+"/capture/task", nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

func newCaptureStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status <task-id>",
		Short: "Show per-agent status for a capture task",
		Long: `Shows which agents picked up the capture task and their reporting status.
The task ID is in the response of "edx capture start" (field "id").`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/pipelines/capture/task/%s/status", url.PathEscape(args[0]))
			data, err := c.Get(cmdContext(cmd), path, nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

func newCaptureResultsCmd() *cobra.Command {
	var ep extraParams
	cmd := &cobra.Command{
		Use:   "results <conf-id>",
		Short: "Fetch captured samples (before/after each node)",
		Long: `Fetch captured samples. Each captured node reports "before" and "after"
arrays; the items are JSON-encoded strings, so decode them with jq's fromjson:

  edx capture results <conf-id> | jq '[.[].nodes[].after[]] | map(fromjson)'`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			if err := ep.apply(q); err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/pipelines/"+url.PathEscape(args[0])+"/capture", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	ep.register(cmd)
	return cmd
}
