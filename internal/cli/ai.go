package cli

import (
	"net/url"

	"github.com/spf13/cobra"
)

func newAICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai",
		Short: "AI Teammate: connectors and activity",
		Long: `Manage the AI Teammate product: data connectors that feed the teammate
(PagerDuty, Slack, GitHub, ...) and teammate activity metrics.`,
	}
	cmd.AddCommand(newAIConnectorsCmd(), newAIActivityCmd())
	return cmd
}

func newAIConnectorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connectors",
		Short: "Manage AI Teammate connectors",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List configured AI connectors",
			RunE: func(cmd *cobra.Command, args []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				data, err := c.Get(cmdContext(cmd), "/ai/connectors", nil)
				if err != nil {
					return err
				}
				return printResult(data)
			},
		},
		&cobra.Command{
			Use:   "specs",
			Short: "List available connector types and their configuration specs",
			RunE: func(cmd *cobra.Command, args []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				data, err := c.Get(cmdContext(cmd), "/ai/connectors/specs", nil)
				if err != nil {
					return err
				}
				return printResult(data)
			},
		},
		&cobra.Command{
			Use:   "environments",
			Short: "List environments available for connector deployment",
			RunE: func(cmd *cobra.Command, args []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				data, err := c.Get(cmdContext(cmd), "/ai/connectors/environments", nil)
				if err != nil {
					return err
				}
				return printResult(data)
			},
		},
		newAIConnectorsUpdateCmd(),
		newAIConnectorsDeleteCmd(),
	)
	return cmd
}

func newAIConnectorsUpdateCmd() *cobra.Command {
	var file string
	cmd := &cobra.Command{
		Use:   "update --file connector.json",
		Short: "Create or update a connector from a JSON request",
		Long: `Create or update a connector. The request body follows the connector spec;
use "edx ai connectors specs" to see required fields per connector type.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := readFileOrStdin(file)
			if err != nil {
				return err
			}
			data, err := c.Post(cmdContext(cmd), "/ai/connectors", nil, body)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", `connector request JSON file ("-" for stdin)`)
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newAIConnectorsDeleteCmd() *cobra.Command {
	var file string
	cmd := &cobra.Command{
		Use:   "delete --file connector.json",
		Short: "Delete a connector (request body identifies the connector)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm("Delete this AI connector?") {
				return errAborted
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := readFileOrStdin(file)
			if err != nil {
				return err
			}
			data, err := c.Delete(cmdContext(cmd), "/ai/connectors", nil, body)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", `connector request JSON file ("-" for stdin)`)
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newAIActivityCmd() *cobra.Command {
	var ep extraParams
	var tf timeFlags
	cmd := &cobra.Command{
		Use:   "activity",
		Short: "Show AI Teammate activity metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			tf.apply(q)
			if err := ep.apply(q); err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/ai/activity_metrics", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	tf.register(cmd, "24h")
	ep.register(cmd)
	return cmd
}
