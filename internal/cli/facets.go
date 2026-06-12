package cli

import (
	"net/url"

	"github.com/spf13/cobra"
)

const scopeHelp = "data scope: log, metric, trace, pattern or event"

func newFacetsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "facets",
		Short: "Discover queryable fields and their values",
		Long: `Discover the schema of your data before writing CQL queries:

  edx facets keys --scope log                      all filterable field names
  edx facets options --scope log --facet service.name   values for one field
  edx facets list --scope log                      configured facets`,
	}
	cmd.AddCommand(newFacetsListCmd(), newFacetsOptionsCmd(), newFacetsKeysCmd())
	return cmd
}

func newFacetsListCmd() *cobra.Command {
	var scope string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List configured facets (builtin and user-defined)",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			if scope != "" {
				q.Set("scope", scope)
			}
			data, err := c.Get(cmdContext(cmd), "/facets", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVar(&scope, "scope", "log", scopeHelp)
	return cmd
}

func newFacetsOptionsCmd() *cobra.Command {
	var scope, facet string
	var limit int
	var tf timeFlags
	cmd := &cobra.Command{
		Use:   "options",
		Short: "List observed values for a field",
		Example: `  edx facets options --scope log --facet service.name
  edx facets options --scope metric --facet name --limit 500`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			q.Set("scope", scope)
			q.Set("facet_path", facet)
			if limit > 0 {
				q.Set("limit", itoa(limit))
			}
			tf.apply(q)
			data, err := c.Get(cmdContext(cmd), "/facet_options", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVar(&scope, "scope", "log", scopeHelp)
	cmd.Flags().StringVar(&facet, "facet", "", "field name, e.g. service.name (required)")
	_ = cmd.MarkFlagRequired("facet")
	cmd.Flags().IntVar(&limit, "limit", 50, "maximum number of values")
	tf.register(cmd, "")
	return cmd
}

func newFacetsKeysCmd() *cobra.Command {
	var scope string
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "List all filterable field names for a scope",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			q.Set("scope", scope)
			data, err := c.Get(cmdContext(cmd), "/facet_keys", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVar(&scope, "scope", "log", scopeHelp)
	return cmd
}
