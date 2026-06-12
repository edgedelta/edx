package cli

import (
	"net/url"

	"github.com/spf13/cobra"
)

func newFleetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fleet",
		Short: "Fleet-wide agent and deployment visibility",
		Long: `Organization-wide view of agents and pipeline deployments. For a single
pipeline's agents use "edx pipelines agents <conf-id>".`,
	}
	cmd.AddCommand(newFleetAgentsCmd(), newFleetDeploymentsCmd())
	return cmd
}

func newFleetAgentsCmd() *cobra.Command {
	var ep extraParams
	cmd := &cobra.Command{
		Use:   "agents",
		Short: "List agents across all pipelines",
		Example: `  edx fleet agents --output table
  edx fleet agents --param keyword=prod`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			if err := ep.apply(q); err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/pipelines/agents_v2", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	ep.register(cmd)
	return cmd
}

func newFleetDeploymentsCmd() *cobra.Command {
	var ep extraParams
	cmd := &cobra.Command{
		Use:   "deployments [conf-id]",
		Short: "Show pipeline deployment status across the fleet",
		Long: `Without arguments, lists deployment status for all pipelines. With a conf
ID, shows deployment detail for that pipeline (agent versions, config rollout).`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			if err := ep.apply(q); err != nil {
				return err
			}
			path := "/pipelines/deployments_v2"
			if len(args) == 1 {
				path += "/" + url.PathEscape(args[0])
			}
			data, err := c.Get(cmdContext(cmd), path, q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	ep.register(cmd)
	return cmd
}
