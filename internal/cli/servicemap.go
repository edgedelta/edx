package cli

import (
	"net/url"

	"github.com/spf13/cobra"
)

func newServiceMapCmd() *cobra.Command {
	var tf timeFlags
	var ep extraParams
	cmd := &cobra.Command{
		Use:   "service-map",
		Short: "Get the service dependency map (from traces)",
		Example: `  edx service-map --lookback 1h
  edx service-map --output yaml`,
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
			data, err := c.Get(cmdContext(cmd), "/service_map", q)
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
