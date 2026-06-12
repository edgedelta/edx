package cli

import (
	"net/url"

	"github.com/spf13/cobra"
)

func newDashboardsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dashboards",
		Short: "List and inspect dashboards",
	}
	cmd.AddCommand(
		&cobra.Command{
			Use:     "list",
			Short:   "List dashboards",
			Example: `  edx dashboards list --output table --columns id,name,creator`,
			RunE: func(cmd *cobra.Command, args []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				data, err := c.Get(cmdContext(cmd), "/dashboards", nil)
				if err != nil {
					return err
				}
				return printResult(data)
			},
		},
		&cobra.Command{
			Use:   "get <dashboard-id>",
			Short: "Get a dashboard definition",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				data, err := c.Get(cmdContext(cmd), "/dashboards/"+url.PathEscape(args[0]), nil)
				if err != nil {
					return err
				}
				return printResult(data)
			},
		},
	)
	return cmd
}
