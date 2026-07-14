package cli

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func newMonitorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitors",
		Short: "Manage monitors and view their states",
		Long: `Create, inspect, update and delete monitors, and view monitor states
(triggered/resolved) for alerting.`,
	}
	cmd.AddCommand(
		newMonitorsListCmd(),
		newMonitorsGetCmd(),
		newMonitorsCreateCmd(),
		newMonitorsUpdateCmd(),
		newMonitorsDeleteCmd(),
		newMonitorsStatesCmd(),
		newMonitorsEvaluateCmd(),
	)
	return cmd
}

func newMonitorsListCmd() *cobra.Command {
	var query string
	var pg pageFlags
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List monitors",
		Example: `  edx monitors list --output table --columns id,name,type,enabled`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			if query != "" {
				q.Set("query", query)
			}
			pg.apply(q)
			data, err := c.Get(cmdContext(cmd), "/monitors", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&query, "query", "q", "", "CQL filter expression")
	pg.register(cmd, 0)
	return cmd
}

func newMonitorsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <monitor-id>",
		Short: "Get a monitor by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/monitors/"+url.PathEscape(args[0]), nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

func newMonitorsCreateCmd() *cobra.Command {
	var file string
	cmd := &cobra.Command{
		Use:   "create --file monitor.json",
		Short: "Create a monitor from a JSON definition",
		Long: `Create a monitor. The definition is the JSON monitor body; pass a file
path or "-" to read from stdin. Tip: fetch an existing monitor with
"edx monitors get" to use as a starting template.`,
		Example: `  edx monitors create --file monitor.json
  edx monitors get <id> | jq '.name = "copy"' | edx monitors create --file -`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := readFileOrStdin(file)
			if err != nil {
				return err
			}
			data, err := c.Post(cmdContext(cmd), "/monitors", nil, body)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", `monitor JSON file ("-" for stdin)`)
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newMonitorsUpdateCmd() *cobra.Command {
	var file string
	cmd := &cobra.Command{
		Use:   "update <monitor-id> --file monitor.json",
		Short: "Update a monitor from a JSON definition",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := readFileOrStdin(file)
			if err != nil {
				return err
			}
			data, err := c.Put(cmdContext(cmd), "/monitors/"+url.PathEscape(args[0]), nil, body)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", `monitor JSON file ("-" for stdin)`)
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newMonitorsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <monitor-id>",
		Short: "Delete a monitor",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm(fmt.Sprintf("Delete monitor %s?", args[0])) {
				return errAborted
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Delete(cmdContext(cmd), "/monitors/"+url.PathEscape(args[0]), nil, nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

func newMonitorsStatesCmd() *cobra.Command {
	var query string
	var pg pageFlags
	cmd := &cobra.Command{
		Use:   "states",
		Short: "List monitor states (triggered / resolved)",
		Example: `  edx monitors states --output table
  edx monitors states --query 'monitor.status:"alert"'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			if query != "" {
				q.Set("query", query)
			}
			pg.apply(q)
			data, err := c.Get(cmdContext(cmd), "/monitor_states", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&query, "query", "q", "", "CQL filter expression")
	pg.register(cmd, 100)
	return cmd
}

func newMonitorsEvaluateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "evaluate <monitor-id>",
		Short: "Run a metric monitor's query now and show value vs thresholds (dry-run)",
		Long: `Evaluate a metric_threshold monitor against current data without changing its
state. Useful to confirm the monitor's query returns data and would (or would
not) fire — handy when a freshly created monitor still shows "No Data" because
its scheduled evaluation hasn't run yet.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			ctx := cmdContext(cmd)
			raw, err := c.Get(ctx, "/monitors/"+url.PathEscape(args[0]), nil)
			if err != nil {
				return err
			}
			var mon struct {
				Name               string  `json:"name"`
				Type               string  `json:"type"`
				EvaluationWindow   int     `json:"evaluation_window"`
				EvaluationFunction string  `json:"evaluation_function"`
				AlertThreshold     float64 `json:"alert_threshold"`
				WarningThreshold   float64 `json:"warning_threshold"`
				ThresholdType      string  `json:"threshold_type"`
				FormulaQuery       struct {
					Formula string                     `json:"formula"`
					Queries map[string]json.RawMessage `json:"queries"`
				} `json:"formula_query"`
			}
			if err := json.Unmarshal(raw, &mon); err != nil {
				return fmt.Errorf("could not parse monitor: %w", err)
			}
			if mon.Type != "metric_threshold" {
				return fmt.Errorf("evaluate supports metric_threshold monitors; %q is %q", args[0], mon.Type)
			}
			if len(mon.FormulaQuery.Queries) == 0 {
				return fmt.Errorf("monitor has no formula_query.queries to evaluate")
			}
			window := mon.EvaluationWindow
			if window <= 0 {
				window = 3600
			}
			formula := mon.FormulaQuery.Formula
			if formula == "" {
				formula = "A"
			}
			payload := map[string]any{
				"queries":  mon.FormulaQuery.Queries,
				"formulas": map[string]any{"A": map[string]any{"formula": formula}},
			}
			body, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			q := url.Values{}
			q.Set("lookback", fmt.Sprintf("%ds", window))
			q.Set("graph_type", "table")
			gdata, err := c.Post(ctx, "/graph", q, body)
			if err != nil {
				return err
			}
			value, records := sumGraphValue(gdata)
			verdict := classifyThreshold(value, mon.AlertThreshold, mon.WarningThreshold, mon.ThresholdType)
			result := map[string]any{
				"monitor":                   mon.Name,
				"type":                      mon.Type,
				"evaluation_window_seconds": window,
				"evaluation_function":       mon.EvaluationFunction,
				"value":                     value,
				"records":                   records,
				"warning_threshold":         mon.WarningThreshold,
				"alert_threshold":           mon.AlertThreshold,
				"threshold_type":            mon.ThresholdType,
				"verdict":                   verdict,
			}
			out, err := json.Marshal(result)
			if err != nil {
				return err
			}
			return printResult(out)
		},
	}
}

// sumGraphValue sums the aggregate values across all records of a /graph
// response and returns the total and the number of records (groups) seen.
func sumGraphValue(data []byte) (float64, int) {
	var g map[string]struct {
		Records []struct {
			Aggregate struct {
				Value float64 `json:"value"`
			} `json:"aggregate"`
		} `json:"records"`
	}
	if err := json.Unmarshal(data, &g); err != nil {
		return 0, 0
	}
	var total float64
	var n int
	for _, series := range g {
		for _, r := range series.Records {
			total += r.Aggregate.Value
			n++
		}
	}
	return total, n
}

// classifyThreshold reports ALERT / WARNING / OK for a value against a monitor's
// thresholds. threshold_type "below" inverts the comparison.
func classifyThreshold(value, alert, warn float64, thresholdType string) string {
	if thresholdType == "below" {
		switch {
		case value < alert:
			return "ALERT"
		case value < warn:
			return "WARNING"
		default:
			return "OK"
		}
	}
	switch {
	case value > alert:
		return "ALERT"
	case value > warn:
		return "WARNING"
	default:
		return "OK"
	}
}
