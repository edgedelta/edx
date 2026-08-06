package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func newLookupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lookup",
		Short: "Manage lookup tables (CSV enrichment tables)",
		Long: `Manage lookup tables: CSV files pipelines use to enrich telemetry.

A table's ID is its file name and must end in .csv or .csv.gz:

  edx lookup list                                  all tables with row counts
  edx lookup get users.csv                         one table's metadata
  edx lookup download users.csv --out users.csv    fetch the table data
  edx lookup create ./users.csv --description "user directory"
  edx lookup update users.csv --file ./new.csv     replace the table data
  edx lookup delete users.csv`,
	}
	cmd.AddCommand(
		newLookupListCmd(),
		newLookupGetCmd(),
		newLookupDownloadCmd(),
		newLookupCreateCmd(),
		newLookupUpdateCmd(),
		newLookupDeleteCmd(),
	)
	return cmd
}

// lookupTableName validates that name is a plausible lookup table ID. The
// backend derives table IDs from uploaded file names and requires this suffix.
func lookupTableName(name string) error {
	if !strings.HasSuffix(name, ".csv") && !strings.HasSuffix(name, ".csv.gz") {
		return fmt.Errorf("lookup table names end in .csv or .csv.gz, got %q", name)
	}
	return nil
}

func newLookupListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List lookup tables (metadata: name, size, row count, tags)",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/lookup_tables/metadata", nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

func newLookupGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <name>",
		Short: "Show one lookup table's metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/lookup_tables/"+url.PathEscape(args[0])+"/metadata", nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

func newLookupDownloadCmd() *cobra.Command {
	var out string
	cmd := &cobra.Command{
		Use:   "download <name>",
		Short: "Download a lookup table's CSV data",
		Example: `  edx lookup download users.csv                  # CSV to stdout
  edx lookup download users.csv --out users.csv  # CSV to file`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/lookup_tables/"+url.PathEscape(args[0]), nil)
			if err != nil {
				return err
			}
			var resp struct {
				Data []byte `json:"data"`
			}
			if err := json.Unmarshal(data, &resp); err != nil {
				return fmt.Errorf("unexpected response: %w", err)
			}
			if out == "" || out == "-" {
				_, err := os.Stdout.Write(resp.Data)
				return err
			}
			if err := os.WriteFile(out, resp.Data, 0o600); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Wrote %d bytes to %s\n", len(resp.Data), out)
			return nil
		},
	}
	cmd.Flags().StringVar(&out, "out", "", `output file ("-" or empty for stdout)`)
	return cmd
}

func newLookupCreateCmd() *cobra.Command {
	var name, description, tags string
	cmd := &cobra.Command{
		Use:   "create <file>",
		Short: "Create a lookup table from a CSV file",
		Long: `Create a lookup table by uploading a CSV file (or .csv.gz).

The table is named after the uploaded file; pass --name to override. The
first CSV row is treated as the header.`,
		Example: `  edx lookup create ./users.csv --description "user directory" --tags "auth,identity"
  edx lookup create ./export.csv --name users.csv
  cat users.csv | edx lookup create - --name users.csv`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				if args[0] == "-" {
					return fmt.Errorf("--name is required when reading from stdin")
				}
				name = filepath.Base(args[0])
			}
			if err := lookupTableName(name); err != nil {
				return err
			}
			file, err := readFileOrStdin(args[0])
			if err != nil {
				return err
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.DoMultipart(cmdContext(cmd), http.MethodPost, c.OrgPath("/lookup_tables"), nil,
				name, file, map[string]string{"description": description, "tags": tags})
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "table name (default: base name of <file>; must end in .csv or .csv.gz)")
	cmd.Flags().StringVar(&description, "description", "", "description of the table")
	cmd.Flags().StringVar(&tags, "tags", "", "tags for the table")
	return cmd
}

func newLookupUpdateCmd() *cobra.Command {
	var file, description, tags string
	cmd := &cobra.Command{
		Use:   "update <name>",
		Short: "Replace a lookup table's data (and optionally its metadata)",
		Long: `Replace an existing lookup table's data with a new CSV file. The local
file name does not matter: the upload targets the table named <name>.
Description and tags are kept unless new values are passed.`,
		Example: `  edx lookup update users.csv --file ./new-users.csv
  edx lookup update users.csv --file - --description "refreshed nightly" < new.csv`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if err := lookupTableName(name); err != nil {
				return err
			}
			body, err := readFileOrStdin(file)
			if err != nil {
				return err
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.DoMultipart(cmdContext(cmd), http.MethodPut, c.OrgPath("/lookup_tables/"+url.PathEscape(name)), nil,
				name, body, map[string]string{"description": description, "tags": tags})
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", `CSV file with the new data ("-" for stdin)`)
	_ = cmd.MarkFlagRequired("file")
	cmd.Flags().StringVar(&description, "description", "", "new description (default: keep current)")
	cmd.Flags().StringVar(&tags, "tags", "", "new tags (default: keep current)")
	return cmd
}

func newLookupDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a lookup table",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm(fmt.Sprintf("Delete lookup table %s?", args[0])) {
				return errAborted
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Delete(cmdContext(cmd), "/lookup_tables/"+url.PathEscape(args[0]), nil, nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}
