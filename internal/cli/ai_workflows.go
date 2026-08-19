package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/edgedelta/edx/internal/api"
)

func newAIWorkflowsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflows",
		Short: "Inspect and run AI Team workflows",
		Long: `Inspect AI Team workflows and their runs, and trigger manual runs.

Workflows are node graphs (Start, Task, Action, If/Else, Transform, Wait)
built in the Edge Delta web app under AI Team > Workflows. Served by the
workflow service (workflow.ai.edgedelta.com).

Run stats (total runs, success rate, avg runtime) are derived from the
ed.oncall.ai.workflow_execution.count and .duration metrics; query them with
"edx metrics query".`,
		Example: `  edx ai workflows list --output table --columns workflowId,displayName,status
  edx ai workflows get <workflow-id>
  edx ai workflows run <workflow-id> --input '{"alert":"cpu high"}'
  edx ai workflows runs list <workflow-id>
  edx ai workflows runs steps <workflow-id> <execution-id>`,
	}
	cmd.AddCommand(
		newAIWorkflowsListCmd(),
		newAIWorkflowsGetCmd(),
		newAIWorkflowsRunCmd(),
		newAIWorkflowsRunsCmd(),
	)
	return cmd
}

// wfGet performs a GET against a path on the workflow service.
func wfGet(cmd *cobra.Command, path string, q url.Values) error {
	c, err := newClient()
	if err != nil {
		return err
	}
	data, err := c.GetFrom(cmdContext(cmd), api.ServiceWorkflow, path, q)
	if err != nil {
		return err
	}
	return printResult(data)
}

func newAIWorkflowsListCmd() *cobra.Command {
	var page aiPageFlags
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List workflows for the organization",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			page.apply(q)
			return wfGet(cmd, "/workflows", q)
		},
	}
	page.register(cmd)
	return cmd
}

func newAIWorkflowsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <workflow-id>",
		Short: "Get a single workflow by ID (the node graph is in \"content\")",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return wfGet(cmd, "/workflows/"+url.PathEscape(args[0]), nil)
		},
	}
}

func newAIWorkflowsRunsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runs",
		Short: "Run history: executions and their steps",
		Long: `Inspect workflow runs (executions). "list" takes a workflow ID; "get" and
"steps" take the workflow ID plus an execution ID from a previous "list" or
"run".`,
	}
	cmd.AddCommand(
		newAIWorkflowsRunsListCmd(),
		newAIWorkflowsRunsGetCmd(),
		newAIWorkflowsRunsStepsCmd(),
	)
	return cmd
}

func newAIWorkflowsRunsListCmd() *cobra.Command {
	var page aiPageFlags
	cmd := &cobra.Command{
		Use:   "list <workflow-id>",
		Short: "List runs of a workflow, newest first",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			page.apply(q)
			return wfGet(cmd, "/workflows/"+url.PathEscape(args[0])+"/executions", q)
		},
	}
	page.register(cmd)
	return cmd
}

func newAIWorkflowsRunsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <workflow-id> <execution-id>",
		Short: "Get a single run",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return wfGet(cmd, workflowExecutionPath(args[0], args[1], ""), nil)
		},
	}
}

func newAIWorkflowsRunsStepsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "steps <workflow-id> <execution-id>",
		Short: "List the recorded steps of a run",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return wfGet(cmd, workflowExecutionPath(args[0], args[1], "/steps"), nil)
		},
	}
}

// workflowExecutionPath builds a /workflows/{id}/executions/{id}... path.
func workflowExecutionPath(workflowID, executionID, suffix string) string {
	return "/workflows/" + url.PathEscape(workflowID) + "/executions/" + url.PathEscape(executionID) + suffix
}

func newAIWorkflowsRunCmd() *cobra.Command {
	var input string
	cmd := &cobra.Command{
		Use:   "run <workflow-id>",
		Short: "Trigger a manual run and stream its progress until it finishes",
		Long: `Trigger a manual run of a workflow. The workflow service streams progress as
events (execution created, per-node step results, completion); each event is
printed as one JSON line as it arrives. The command exits 0 when the run
completes and non-zero when the run reports an error. A run can take minutes;
cancel watching with Ctrl-C (the run itself keeps going server-side).`,
		Example: `  edx ai workflows run <workflow-id>
  edx ai workflows run <workflow-id> --input '{"alert":"cpu high"}'
  edx ai workflows run <workflow-id> --input @input.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := workflowRunBody(input)
			if err != nil {
				return err
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			stream, err := c.StreamFrom(cmdContext(cmd), api.ServiceWorkflow, http.MethodPost,
				"/workflows/"+url.PathEscape(args[0])+"/executions", nil, body)
			if err != nil {
				return err
			}
			defer stream.Close()
			return printWorkflowRunStream(os.Stdout, stream)
		},
	}
	cmd.Flags().StringVar(&input, "input", "{}", `run input: inline JSON or plain string, or "@file" / "@-" for stdin`)
	return cmd
}

// workflowRunBody builds the manual-run request body. The endpoint requires an
// "input" field; the value is inline JSON, a plain string (sent as a JSON
// string), or "@file" / "@-" content (same convention as --*-prompt flags).
func workflowRunBody(input string) ([]byte, error) {
	v, err := resolvePromptValue(input)
	if err != nil {
		return nil, err
	}
	raw := json.RawMessage(v)
	if !json.Valid(raw) {
		raw, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
	}
	return json.Marshal(map[string]json.RawMessage{"input": raw})
}

// workflowRunFailureTypes are stream event types that mean the run did not
// complete; stepResult sub-type "workflow-error" counts too.
var workflowRunFailureTypes = map[string]bool{
	"run-error":        true,
	"create-error":     true,
	"step-fetch-error": true,
}

// printWorkflowRunStream consumes the manual-run SSE stream, printing each
// event's data as one JSON line, and returns an error when the stream reported
// a failure. Events are "event:"/"data:" blocks separated by blank lines; the
// terminal event is named "done".
func printWorkflowRunStream(w io.Writer, stream io.Reader) error {
	sc := bufio.NewScanner(stream)
	// Step results can carry large node outputs; the default 64KB token limit
	// is not enough.
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)

	var runErr error
	event := ""
	var data strings.Builder
	dispatch := func() bool {
		defer func() { event = ""; data.Reset() }()
		if event == "done" {
			return true
		}
		payload := strings.TrimSpace(data.String())
		if payload == "" {
			return false
		}
		fmt.Fprintln(w, payload)
		var ev struct {
			Type       string `json:"type"`
			Message    string `json:"message"`
			StepResult struct {
				Type    string `json:"type"`
				Message string `json:"message"`
			} `json:"stepResult"`
		}
		if json.Unmarshal([]byte(payload), &ev) == nil {
			if workflowRunFailureTypes[ev.Type] {
				runErr = fmt.Errorf("workflow run failed: %s", ev.Message)
			}
			if ev.StepResult.Type == "workflow-error" {
				msg := ev.StepResult.Message
				if msg == "" {
					msg = "a node reported an error"
				}
				runErr = fmt.Errorf("workflow run failed: %s", msg)
			}
		}
		return false
	}

	for sc.Scan() {
		line := sc.Text()
		switch {
		case line == "":
			if dispatch() {
				return runErr
			}
		case strings.HasPrefix(line, "event:"):
			event = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
		case strings.HasPrefix(line, "data:"):
			data.WriteString(strings.TrimPrefix(strings.TrimPrefix(line, "data:"), " "))
		}
	}
	if err := sc.Err(); err != nil {
		return fmt.Errorf("reading run stream: %w", err)
	}
	// Stream ended without a "done" event; flush any trailing block.
	dispatch()
	return runErr
}
