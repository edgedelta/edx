package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/edgedelta/edx/internal/cql"
)

func newCQLCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cql",
		Aliases: []string{"mql", "edcql", "edmql"},
		Short:   "Check Edge Delta query syntax offline",
		Long: `Check the syntax of Edge Delta queries (CQL/MQL) without a backend.

Queries are parsed with the same ANTLR grammars the backend parses them with, so
what is accepted here is what the backend accepts. The grammars differ by data
type — a metric query and a log query are not interchangeable — so --type
selects one.

This is a syntax check only. Whether a facet or metric name exists is resolved
by the backend against your organization's data and cannot be checked offline.`,
	}
	cmd.AddCommand(newCQLValidateCmd())
	return cmd
}

func newCQLValidateCmd() *cobra.Command {
	var dataType string
	var file string
	var quiet bool

	cmd := &cobra.Command{
		Use:   "validate [query...] --type log",
		Short: "Validate query syntax",
		Long: fmt.Sprintf(`Validate one or more queries against the grammar for a data type.

Pass queries as arguments, or use --file to read them one per line ("-" for
stdin), which is the form to use for a saved list of queries. Blank lines and
lines starting with "#" are skipped.

--type accepts: %s
The log grammar covers log, event, pattern and trace queries, which all share it.

Exits non-zero if any query fails to parse. Errors carry the column and a caret
under the offending token.`, strings.Join(cql.DataTypes(), ", ")),
		Example: `  edx cql validate --type metric 'sum:ed.host.cpu{*} by {host.name}.rollup(60)'
  edx cql validate --type log '{error AND ed.tag:$fleet} by {host.name}'
  edx cql validate --type formula 'timeshift(q1, 3600) / q2'
  edx cql validate --type metric --file queries.txt
  edx dashboards get <id> \
    | jq -r '.definition.widgets[]?.visuals[]?.dataSource | select(.type=="metric") | .params.query // empty' \
    | edx cql validate --type metric --file -`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dialect, ok := cql.DialectForDataType(dataType)
			if !ok {
				return fmt.Errorf("unknown --type %q; expected one of: %s",
					dataType, strings.Join(cql.DataTypes(), ", "))
			}

			queries, err := collectQueries(args, file, cmd.InOrStdin())
			if err != nil {
				return err
			}
			if len(queries) == 0 {
				return fmt.Errorf("no queries to validate; pass them as arguments or via --file")
			}

			out := cmd.OutOrStdout()
			var failed int
			for _, query := range queries {
				errs, err := cql.Validate(dialect, query)
				if err != nil {
					return err
				}
				if len(errs) == 0 {
					if !quiet {
						fmt.Fprintf(out, "ok  %s\n", query)
					}
					continue
				}
				failed++
				// Reported through warnf like the dashboard issues, so failures go to
				// stderr and stay visible when stdout is piped.
				for _, e := range errs {
					warnf("invalid %s query %s", dialect, e)
					if !quiet {
						fmt.Fprintln(cmd.ErrOrStderr(), indent(e.Annotate(query), "    "))
					}
				}
			}

			if failed > 0 {
				return fmt.Errorf("%d of %d quer%s failed to parse", failed, len(queries),
					plural(len(queries), "y", "ies"))
			}
			if quiet {
				return nil
			}
			fmt.Fprintf(out, "\n%d quer%s valid\n", len(queries), plural(len(queries), "y is", "ies are"))
			return nil
		},
	}

	cmd.Flags().StringVarP(&dataType, "type", "t", "", fmt.Sprintf("data type the query targets: %s", strings.Join(cql.DataTypes(), ", ")))
	cmd.Flags().StringVarP(&file, "file", "f", "", `read queries from a file, one per line ("-" for stdin)`)
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "print only failures")
	_ = cmd.MarkFlagRequired("type")
	return cmd
}

// collectQueries gathers queries from positional arguments and, when file is set, from a
// file or stdin one per line.
func collectQueries(args []string, file string, stdin io.Reader) ([]string, error) {
	queries := make([]string, 0, len(args))
	for _, a := range args {
		if strings.TrimSpace(a) != "" {
			queries = append(queries, a)
		}
	}

	if file == "" {
		return queries, nil
	}

	var raw []byte
	var err error
	if file == "-" {
		raw, err = io.ReadAll(stdin)
	} else {
		raw, err = readFileOrStdin(file)
	}
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(raw))
	// Queries can be long; the default 64KB token limit is not enough for a generated one.
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		queries = append(queries, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read queries: %w", err)
	}
	return queries, nil
}

func indent(s, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}
