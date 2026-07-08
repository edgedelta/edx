package cli

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

// newPipelinesTestCmd groups the dry-run "test" subcommands. Each one runs a
// processor against sample data on the backend WITHOUT saving or deploying, so a
// transform can be developed and verified before it ever touches live traffic.
func newPipelinesTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Dry-run a processor against sample data (no save, no deploy)",
		Long: `Dry-run a processor against sample logs on the backend without saving or
deploying it - the offline complement to "edx pipelines capture".

  test ottl   apply OTTL statements to sample OTEL log items
  test node   run a full processor node of any type (log_to_metric, log_to_pattern, ...)
  test grok   match a grok pattern against raw log lines
  test regex  match a regex against raw log lines

Iterate here until the output is right, then put the statements/node into the
config, "edx pipelines save", and "edx pipelines deploy".`,
	}
	cmd.AddCommand(
		newPipelinesTestOTTLCmd(),
		newPipelinesTestNodeCmd(),
		newPipelinesTestGrokCmd(),
		newPipelinesTestRegexCmd(),
	)
	return cmd
}

// readLines reads a file (or stdin for "-") and returns its non-empty lines.
func readLines(path string) ([]string, error) {
	b, err := readFileOrStdin(path)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0)
	for _, ln := range strings.Split(string(b), "\n") {
		ln = strings.TrimRight(ln, "\r")
		if strings.TrimSpace(ln) == "" {
			continue
		}
		out = append(out, ln)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no input lines read from %q", path)
	}
	return out, nil
}

// postTest marshals payload and POSTs it to /pipelines/<conf-id>/test/<kind>.
func postTest(cmd *cobra.Command, confID, kind string, payload map[string]any) error {
	c, err := newClient()
	if err != nil {
		return err
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	data, err := c.Post(cmdContext(cmd), "/pipelines/"+url.PathEscape(confID)+"/test/"+kind, nil, body)
	if err != nil {
		return err
	}
	return printResult(data)
}

func newPipelinesTestOTTLCmd() *cobra.Command {
	var file, statements, statementsFile string
	cmd := &cobra.Command{
		Use:   "ottl <conf-id> --file items.jsonl --statements '<ottl>'",
		Short: "Dry-run OTTL statements against sample OTEL log items",
		Long: `Apply OTTL statements to sample OTEL log items and print the transformed items.
Each line of --file is a JSON-encoded OTEL log record, e.g.
  {"body":"...","attributes":{},"resource":{},"timestamp":1}`,
		Example: `  edx pipelines test ottl <conf-id> --file items.jsonl --statements 'set(severity_text, "INFO")'
  edx pipelines test ottl <conf-id> --file items.jsonl --statements-file transform.ottl`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stmts, err := resolveInline(statements, statementsFile, "--statements", "--statements-file")
			if err != nil {
				return err
			}
			items, err := readLines(file)
			if err != nil {
				return err
			}
			return postTest(cmd, args[0], "ottl", map[string]any{"items": items, "statements": stmts})
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", `sample OTEL log items, one JSON object per line ("-" for stdin)`)
	_ = cmd.MarkFlagRequired("file")
	cmd.Flags().StringVarP(&statements, "statements", "s", "", "OTTL statements (may be multiline)")
	cmd.Flags().StringVar(&statementsFile, "statements-file", "", `file with OTTL statements ("-" for stdin)`)
	return cmd
}

func newPipelinesTestNodeCmd() *cobra.Command {
	var file, nodeFile string
	cmd := &cobra.Command{
		Use:   "node <conf-id> --file items.jsonl --node node.json",
		Short: "Dry-run a full processor node (any type) against sample items",
		Long: `Run a processor node of any type (log_to_metric, log_to_pattern, extract_metric,
parse_json_attributes, grok, mask, ...) against sample OTEL log items and print
the items it emits per output path.

--node is a JSON file holding the node in graph form:
  {"id":"<name>","type":"<node_type>","configuration":{ ...node config... }}`,
		Example: `  edx pipelines test node <conf-id> --file items.jsonl --node l2m.json`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nodeBytes, err := readFileOrStdin(nodeFile)
			if err != nil {
				return err
			}
			var node json.RawMessage
			if err := json.Unmarshal(nodeBytes, &node); err != nil {
				return fmt.Errorf("--node must be a JSON object: %w", err)
			}
			items, err := readLines(file)
			if err != nil {
				return err
			}
			return postTest(cmd, args[0], "node", map[string]any{"items": items, "node": node})
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", `sample OTEL log items, one JSON object per line ("-" for stdin)`)
	_ = cmd.MarkFlagRequired("file")
	cmd.Flags().StringVar(&nodeFile, "node", "", `JSON file with the node in graph form {"id","type","configuration":{...}} ("-" for stdin)`)
	_ = cmd.MarkFlagRequired("node")
	return cmd
}

func newPipelinesTestGrokCmd() *cobra.Command {
	var file, pattern string
	cmd := &cobra.Command{
		Use:     "grok <conf-id> --file logs.txt --pattern '<grok>'",
		Short:   "Dry-run a grok pattern against raw log lines",
		Example: `  edx pipelines test grok <conf-id> --file logs.txt --pattern '%{COMMONAPACHELOG}'`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			logs, err := readLines(file)
			if err != nil {
				return err
			}
			return postTest(cmd, args[0], "grok", map[string]any{"logs": logs, "pattern": pattern})
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", `raw log lines, one per line ("-" for stdin)`)
	_ = cmd.MarkFlagRequired("file")
	cmd.Flags().StringVarP(&pattern, "pattern", "p", "", "grok pattern, e.g. %{COMMONAPACHELOG}")
	_ = cmd.MarkFlagRequired("pattern")
	return cmd
}

func newPipelinesTestRegexCmd() *cobra.Command {
	var file, pattern string
	cmd := &cobra.Command{
		Use:     "regex <conf-id> --file logs.txt --pattern '<regex>'",
		Short:   "Dry-run a regex against raw log lines",
		Example: `  edx pipelines test regex <conf-id> --file logs.txt --pattern '(?P<ip>\S+) .*'`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := readLines(file)
			if err != nil {
				return err
			}
			return postTest(cmd, args[0], "regex", map[string]any{"input": input, "regex": pattern})
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", `raw log lines, one per line ("-" for stdin)`)
	_ = cmd.MarkFlagRequired("file")
	cmd.Flags().StringVarP(&pattern, "pattern", "p", "", "regex pattern (RE2), named groups become fields")
	_ = cmd.MarkFlagRequired("pattern")
	return cmd
}

// resolveInline returns the inline value, or the file's contents when the
// *-file variant is used. Exactly one of the two must be provided.
func resolveInline(inline, file, inlineFlag, fileFlag string) (string, error) {
	switch {
	case inline != "" && file != "":
		return "", fmt.Errorf("use only one of %s or %s", inlineFlag, fileFlag)
	case file != "":
		b, err := readFileOrStdin(file)
		if err != nil {
			return "", err
		}
		return string(b), nil
	case inline != "":
		return inline, nil
	default:
		return "", fmt.Errorf("one of %s or %s is required", inlineFlag, fileFlag)
	}
}
