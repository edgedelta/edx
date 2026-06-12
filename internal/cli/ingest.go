package cli

import (
	"net/url"

	"github.com/spf13/cobra"
)

func newIngestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ingest",
		Short: "Ingestion endpoints and tokens",
		Long: `Discover the organization's HTTPS ingestion endpoints (for OTLP and raw
data) and mint ingestion tokens for sending telemetry.`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "endpoints",
		Short: "Show HTTPS ingestion endpoints, sample payloads and test commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/ingestion_endpoints", nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	})

	var confID, nodeName string
	token := &cobra.Command{
		Use:   "token",
		Short: "Get or create an ingestion token",
		Example: `  edx ingest token --conf-id <pipeline-id> --node-name my_otlp_input`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			q.Set("conf_id", confID)
			q.Set("node_name", nodeName)
			data, err := c.Get(cmdContext(cmd), "/ingestion_token", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	token.Flags().StringVar(&confID, "conf-id", "", "pipeline conf ID the token is bound to (required)")
	_ = token.MarkFlagRequired("conf-id")
	token.Flags().StringVar(&nodeName, "node-name", "", "input node name in the pipeline (required)")
	_ = token.MarkFlagRequired("node-name")
	cmd.AddCommand(token)

	return cmd
}
