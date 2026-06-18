// Package cli defines the edx command tree.
package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/edgedelta/edx/internal/api"
	"github.com/edgedelta/edx/internal/config"
	"github.com/edgedelta/edx/internal/output"
)

var (
	flagProfile string
	flagEnv     string
	flagOrg     string
	flagToken   string
	flagOutput  string
	flagColumns []string
	flagTimeout time.Duration
	flagYes     bool
)

// Execute runs the root command and exits with an appropriate status code.
func Execute() {
	root := NewRootCmd()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

// NewRootCmd builds the full edx command tree.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "edx",
		Short: "Edge Delta CLI",
		Long: `edx is the Edge Delta command line interface.

It covers the three Edge Delta products:
  Pipeline       fleet management, pipeline configs, deployments, live capture
  Observability  logs, patterns, metrics, traces, events, monitors, dashboards
  AI Teammate    issues, threads, channels, teammates (agents), connectors

ENVIRONMENTS
  A profile targets an environment (prod, staging or local). The environment
  selects every service host at once: the main API plus the AI Teammate
  services, which live on their own hosts (chat.ai.edgedelta.com,
  agent.ai.edgedelta.com) rather than under api.edgedelta.com. Switch with
  "--profile <name>", "--env staging" or the ED_ENV variable.

AUTHENTICATION
  Run "edx auth login --token <api-token> --org-id <org-id>" once, or set the
  ED_API_TOKEN and ED_ORG_ID environment variables. Add "--env staging" to log
  in against staging. Tokens are created in the Edge Delta web app under
  Admin > API Tokens.

OUTPUT
  Responses print as pretty JSON by default. Use --output table|csv|yaml|raw,
  and --columns to pick table/csv columns by dot-path.

EXAMPLES
  edx logs search --query 'service.name:"api" AND severity_text:"ERROR"' --lookback 1h
  edx metrics query --name http.request.duration --agg avg --group-by service.name
  edx pipelines list --output table
  edx capture start <pipeline-id> --duration 2m --nodes my_node
  edx ai issues list --output table
  edx api GET /v1/orgs/{org}/dashboards`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	pf := root.PersistentFlags()
	pf.StringVar(&flagProfile, "profile", "", "config profile to use (default from config file or EDX_PROFILE)")
	pf.StringVar(&flagEnv, "env", "", "environment: prod, staging, or local (overrides profile and ED_ENV)")
	pf.StringVar(&flagOrg, "org", "", "Edge Delta organization ID (overrides profile and ED_ORG_ID)")
	pf.StringVar(&flagToken, "token", "", "API token (overrides profile and ED_API_TOKEN)")
	pf.StringVarP(&flagOutput, "output", "o", "json", "output format: json, yaml, table, csv, raw")
	pf.StringSliceVar(&flagColumns, "columns", nil, "columns (dot-paths) for table/csv output")
	pf.DurationVar(&flagTimeout, "timeout", 60*time.Second, "HTTP request timeout")
	pf.BoolVarP(&flagYes, "yes", "y", false, "skip confirmation prompts")

	root.AddCommand(
		newAuthCmd(),
		newLogsCmd(),
		newPatternsCmd(),
		newMetricsCmd(),
		newTracesCmd(),
		newEventsCmd(),
		newMonitorsCmd(),
		newPipelinesCmd(),
		newFleetCmd(),
		newCaptureCmd(),
		newDashboardsCmd(),
		newFacetsCmd(),
		newServiceMapCmd(),
		newAICmd(),
		newHealthCmd(),
		newIngestCmd(),
		newAPICmd(),
		newVersionCmd(),
	)

	// Show the Edge Delta wordmark above the root help (and bare `edx`).
	base := root.HelpFunc()
	root.SetHelpFunc(func(c *cobra.Command, args []string) {
		if c == root {
			writeBanner(os.Stdout)
		}
		base(c, args)
	})
	return root
}

// newClient resolves configuration and returns an authenticated API client.
func newClient() (*api.Client, error) {
	r, err := config.Resolve(flagProfile, flagEnv, flagOrg, flagToken)
	if err != nil {
		return nil, err
	}
	if r.APIToken == "" && !r.UsesOAuth() {
		return nil, fmt.Errorf("no credentials configured: run `edx auth login --token <token> --org-id <org>` (or `--oauth`), or set %s", config.EnvAPIToken)
	}
	if r.OrgID == "" {
		return nil, fmt.Errorf("no organization ID configured: run `edx auth login` with --org-id, set %s, or pass --org", config.EnvOrgID)
	}
	auth := &api.Auth{APIToken: r.APIToken, APIDomain: hostOf(r.APIURL)}
	if r.UsesOAuth() {
		expiry, _ := time.Parse(time.RFC3339, r.OAuthExpiry)
		auth.OAuth = newOAuthTokenSource(r, expiry)
	}
	return api.New(r.APIURL, r.ChatURL, r.AgentURL, r.OrgID, auth, flagTimeout), nil
}

// hostOf returns the bare host (host:port) of a URL, used as X-ED-API-Domain.
func hostOf(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return u.Host
}

// printResult renders raw response bytes using the global output flags.
func printResult(data []byte) error {
	return output.Print(os.Stdout, data, output.Options{Format: flagOutput, Columns: flagColumns})
}

// confirm prompts unless --yes was given. Returns true when the user accepts.
func confirm(prompt string) bool {
	if flagYes {
		return true
	}
	fmt.Fprintf(os.Stderr, "%s [y/N]: ", prompt)
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.ToLower(strings.TrimSpace(line))
	return line == "y" || line == "yes"
}

// readFileOrStdin reads path, treating "-" as stdin.
func readFileOrStdin(path string) ([]byte, error) {
	if path == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}

// timeFlags are the shared time-range parameters for search endpoints.
type timeFlags struct {
	lookback string
	from     string
	to       string
}

func (t *timeFlags) register(cmd *cobra.Command, defaultLookback string) {
	cmd.Flags().StringVar(&t.lookback, "lookback", defaultLookback, "lookback window in Go duration format (e.g. 15m, 1h, 24h); ignored when --from is set")
	cmd.Flags().StringVar(&t.from, "from", "", "start time, ISO 8601 (2006-01-02T15:04:05.000Z)")
	cmd.Flags().StringVar(&t.to, "to", "", "end time, ISO 8601 (2006-01-02T15:04:05.000Z)")
}

func (t *timeFlags) apply(q url.Values) {
	if t.from != "" {
		q.Set("from", t.from)
		if t.to != "" {
			q.Set("to", t.to)
		}
		return
	}
	if t.lookback != "" {
		q.Set("lookback", t.lookback)
	}
	if t.to != "" {
		q.Set("to", t.to)
	}
}

// pageFlags are the shared pagination parameters for search endpoints.
type pageFlags struct {
	limit  int
	cursor string
	order  string
}

func (p *pageFlags) register(cmd *cobra.Command, defaultLimit int) {
	cmd.Flags().IntVar(&p.limit, "limit", defaultLimit, "maximum number of results")
	cmd.Flags().StringVar(&p.cursor, "cursor", "", "pagination cursor from a previous response")
	cmd.Flags().StringVar(&p.order, "order", "desc", "sort order: asc or desc")
}

func (p *pageFlags) apply(q url.Values) {
	if p.limit != 0 {
		q.Set("limit", fmt.Sprintf("%d", p.limit))
	}
	if p.cursor != "" {
		q.Set("cursor", p.cursor)
	}
	if p.order != "" {
		q.Set("order", p.order)
	}
}

// extraParams supports arbitrary --param key=value passthrough query params.
type extraParams struct {
	params []string
}

func (e *extraParams) register(cmd *cobra.Command) {
	cmd.Flags().StringArrayVar(&e.params, "param", nil, "extra query parameter as key=value (repeatable)")
}

func (e *extraParams) apply(q url.Values) error {
	for _, p := range e.params {
		k, v, ok := strings.Cut(p, "=")
		if !ok {
			return fmt.Errorf("invalid --param %q: expected key=value", p)
		}
		q.Add(k, v)
	}
	return nil
}

func cmdContext(cmd *cobra.Command) context.Context {
	return cmd.Context()
}
