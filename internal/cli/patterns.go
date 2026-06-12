package cli

import (
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
)

func newPatternsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "patterns",
		Short: "Log patterns (clustered message signatures)",
		Long: `Edge Delta clusters similar log messages into patterns. Each pattern carries
count, proportion, sentiment (positive/negative/neutral) and delta (change vs
an earlier window) - useful for spotting new or surging error signatures.`,
	}
	cmd.AddCommand(newPatternsListCmd(), newPatternsSamplesCmd())
	return cmd
}

func newPatternsListCmd() *cobra.Command {
	var query, offset string
	var limit int
	var summary, negative bool
	var tf timeFlags
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List top log patterns with statistics",
		Example: `  edx patterns list --lookback 1h --summary
  edx patterns list --query 'service.name:"api"' --negative --limit 20
  edx patterns list --offset 24h   # compute delta vs the same window 24h earlier`,
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
			if summary {
				q.Set("summary", "true")
			}
			if negative {
				q.Set("negative", "true")
			}
			if limit > 0 {
				q.Set("limit", itoa(limit))
			}
			if offset != "" {
				q.Set("offset", offset)
			}
			data, err := c.Get(cmdContext(cmd), "/clustering/stats", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&query, "query", "q", "", "CQL filter query")
	cmd.Flags().BoolVar(&summary, "summary", false, "return ~50 interesting clusters (top anomalies, top/bottom delta and count)")
	cmd.Flags().BoolVar(&negative, "negative", false, "only patterns with negative sentiment")
	cmd.Flags().IntVar(&limit, "limit", 20, "maximum number of clusters")
	cmd.Flags().StringVar(&offset, "offset", "", "comma-separated delta offsets in Go duration format (e.g. 24h)")
	tf.register(cmd, "1h")
	return cmd
}

func newPatternsSamplesCmd() *cobra.Command {
	var query string
	var limit int
	var tf timeFlags
	var ep extraParams
	cmd := &cobra.Command{
		Use:   "samples",
		Short: "Fetch raw log samples for a pattern",
		Example: `  edx patterns samples --query 'service.name:"api"' --param pattern='error connecting to *'`,
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
			if limit > 0 {
				q.Set("limit", itoa(limit))
			}
			if err := ep.apply(q); err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/clustering/samples", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&query, "query", "q", "", "CQL filter query")
	cmd.Flags().IntVar(&limit, "limit", 20, "maximum number of samples")
	tf.register(cmd, "1h")
	ep.register(cmd)
	return cmd
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
