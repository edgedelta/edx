package cli

import (
	"net/url"

	"github.com/spf13/cobra"
)

func newHealthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Pipeline component health",
		Long: `Health status of pipeline components (sources, processors, destinations)
as reported by the agents.`,
	}
	cmd.AddCommand(newHealthComponentsCmd(), newHealthProblemsCmd())
	return cmd
}

func newHealthComponentsCmd() *cobra.Command {
	var host, tag string
	var all bool
	var tf timeFlags
	var ep extraParams
	cmd := &cobra.Command{
		Use:   "components",
		Short: "Component health entries for one host",
		Long: `Component health for a specific host (agent). Find host names with
"edx facets options --scope log --facet host.name" or "edx pipelines agents".`,
		Example: `  edx health components --host web-1 --lookback 1h
  edx health components --host web-1 --tag prod --all`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			q.Set("host", host)
			if tag != "" {
				q.Set("tag", tag)
			}
			if all {
				q.Set("all", "true")
			}
			tf.apply(q)
			if err := ep.apply(q); err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/health/component", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVar(&host, "host", "", "host name to inspect (required)")
	_ = cmd.MarkFlagRequired("host")
	cmd.Flags().StringVar(&tag, "tag", "", "pipeline tag (required for some orgs)")
	cmd.Flags().BoolVar(&all, "all", false, "include only customer-visible components")
	tf.register(cmd, "1h")
	ep.register(cmd)
	return cmd
}

func newHealthProblemsCmd() *cobra.Command {
	var ep extraParams
	var tf timeFlags
	cmd := &cobra.Command{
		Use:   "problems",
		Short: "Components currently reporting problems",
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
			data, err := c.Get(cmdContext(cmd), "/health/problems", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	tf.register(cmd, "1h")
	ep.register(cmd)
	return cmd
}
