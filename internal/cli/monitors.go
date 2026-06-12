package cli

import (
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
