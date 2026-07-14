package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/edgedelta/edx/internal/api"
)

func newPipelinesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pipelines",
		Aliases: []string{"pipeline"},
		Short:   "Manage pipelines: configs, history, deployments",
		Long: `Manage Edge Delta pipelines (fleet configurations).

A pipeline is identified by its conf ID. Use "edx pipelines list" to discover
IDs, then inspect, change and deploy configurations:

  list > get > save (new version) > deploy <version>

Version history is kept server-side; "history" shows it and "deploy" can roll
forward or back to any version.`,
	}
	cmd.AddCommand(
		newPipelinesListCmd(),
		newPipelinesGetCmd(),
		newPipelinesHistoryCmd(),
		newPipelinesSaveCmd(),
		newPipelinesDeployCmd(),
		newPipelinesAgentsCmd(),
		newPipelinesStatusCmd(),
		newPipelinesValidateCmd(),
		newPipelinesDeployCommandCmd(),
		newPipelinesTestCmd(),
	)
	return cmd
}

func newPipelinesListCmd() *cobra.Command {
	var keyword string
	var includeSuspended bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List pipelines (fleets)",
		Example: `  edx pipelines list --output table --columns id,tag,fleet_type,environment,status,updated
  edx pipelines list --keyword prod`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/pipelines", nil)
			if err != nil {
				return err
			}
			if keyword == "" && includeSuspended {
				return printResult(data)
			}
			var pipelines []map[string]any
			if err := json.Unmarshal(data, &pipelines); err != nil {
				return printResult(data)
			}
			filtered := make([]map[string]any, 0, len(pipelines))
			for _, p := range pipelines {
				tag, _ := p["tag"].(string)
				status, _ := p["status"].(string)
				if keyword != "" && !containsFold(tag, keyword) {
					continue
				}
				if !includeSuspended && status == "suspended" {
					continue
				}
				filtered = append(filtered, p)
			}
			out, err := json.Marshal(filtered)
			if err != nil {
				return err
			}
			return printResult(out)
		},
	}
	cmd.Flags().StringVar(&keyword, "keyword", "", "filter by tag substring (case-insensitive)")
	cmd.Flags().BoolVar(&includeSuspended, "all", false, "include suspended pipelines")
	return cmd
}

func newPipelinesGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <conf-id>",
		Short: "Get a pipeline definition (including config content)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/pipelines/"+url.PathEscape(args[0]), nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	return cmd
}

func newPipelinesHistoryCmd() *cobra.Command {
	var ep extraParams
	cmd := &cobra.Command{
		Use:   "history <conf-id>",
		Short: "Show pipeline version history (who changed what, when)",
		Long: `Show the saved version history for a pipeline. The version identifier used
by "edx pipelines deploy" is the entry's epoch-millisecond timestamp, surfaced
here as the "version" column.`,
		Example: `  edx pipelines history <conf-id> --output table --columns version,description,creator,created`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			if err := ep.apply(q); err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/pipelines/"+url.PathEscape(args[0])+"/history", q)
			if err != nil {
				return err
			}
			return printResult(withVersionField(data))
		},
	}
	ep.register(cmd)
	return cmd
}

func newPipelinesSaveCmd() *cobra.Command {
	var file, description string
	cmd := &cobra.Command{
		Use:   "save <conf-id> --file pipeline.yaml",
		Short: "Save a new pipeline version (does not deploy)",
		Long: `Save a new version of the pipeline configuration. The file is the agent
config YAML content. Saving creates a version; use "edx pipelines deploy" to
roll it out to the fleet.`,
		Example: `  edx pipelines save <conf-id> --file pipeline.yaml --description "add k8s source"`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			content, err := readFileOrStdin(file)
			if err != nil {
				return err
			}
			payload := map[string]any{
				"description": description,
				"content":     string(content),
			}
			body, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			data, err := c.Post(cmdContext(cmd), "/pipelines/"+url.PathEscape(args[0])+"/save", nil, body)
			if err != nil {
				return err
			}
			// The save response doesn't cleanly carry the new version, so resolve
			// it from history and surface it for the deploy step.
			if v := latestVersion(cmdContext(cmd), c, args[0]); v != "" {
				notef("Saved version %s — deploy it with: edx pipelines deploy %s %s   (or --latest)", v, args[0], v)
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", `pipeline YAML file ("-" for stdin)`)
	_ = cmd.MarkFlagRequired("file")
	cmd.Flags().StringVarP(&description, "description", "d", "", "change description shown in history")
	return cmd
}

func newPipelinesDeployCmd() *cobra.Command {
	var latest bool
	var wait bool
	var waitTimeout time.Duration
	cmd := &cobra.Command{
		Use:   "deploy <conf-id> [version]",
		Short: "Deploy a saved pipeline version to the fleet",
		Long: `Deploy a specific saved version to all agents in the fleet. Deploying an
older version is the supported way to roll back.

The version is a saved version's epoch-millisecond timestamp (see the "version"
column of "edx pipelines history"). Use "--latest" to deploy the most recently
saved version without looking it up.

With "--wait", edx polls the fleet after deploying and reports whether agents
check in (heartbeat) after the rollout. A frozen heartbeat is the tell-tale
sign that an agent failed to apply the new config (for example, a new source
whose port is already bound fails the whole pipeline).`,
		Example: `  edx pipelines deploy <conf-id> 1783965150095 --yes
  edx pipelines deploy <conf-id> --latest --wait`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			confID := args[0]
			c, err := newClient()
			if err != nil {
				return err
			}
			ctx := cmdContext(cmd)

			version := ""
			if len(args) == 2 {
				version = args[1]
			}
			switch {
			case latest && version != "":
				return fmt.Errorf("pass either a version argument or --latest, not both")
			case latest:
				version = latestVersion(ctx, c, confID)
				if version == "" {
					return fmt.Errorf("could not resolve latest version from history for pipeline %s", confID)
				}
			case version == "":
				return fmt.Errorf("a version is required (see `edx pipelines history`), or pass --latest")
			}

			if !confirm(fmt.Sprintf("Deploy version %s of pipeline %s to the fleet?", version, confID)) {
				return errAborted
			}
			path := "/pipelines/" + url.PathEscape(confID) + "/deploy/" + url.PathEscape(version)
			data, err := c.Post(ctx, path, nil, nil)
			if err != nil {
				return err
			}
			if err := printResult(data); err != nil {
				return err
			}
			if wait {
				return waitForRollout(ctx, c, confID, waitTimeout)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&latest, "latest", false, "deploy the most recently saved version")
	cmd.Flags().BoolVar(&wait, "wait", false, "after deploying, wait for agents to check in and report rollout health")
	cmd.Flags().DurationVar(&waitTimeout, "wait-timeout", 120*time.Second, "how long to wait for agents to check in (with --wait)")
	return cmd
}

func newPipelinesAgentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agents <conf-id>",
		Short: "List agents running this pipeline",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/pipelines/"+url.PathEscape(args[0])+"/agents", nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	return cmd
}

func newPipelinesStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status <conf-id>",
		Short: "Get fleet status (running / suspended) for a pipeline",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/confs/"+url.PathEscape(args[0])+"/fleet/status", nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	return cmd
}

func newPipelinesValidateCmd() *cobra.Command {
	var file string
	cmd := &cobra.Command{
		Use:     "validate --file pipeline.yaml",
		Short:   "Validate an agent configuration without saving it",
		Example: `  edx pipelines validate --file pipeline.yaml`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			content, err := readFileOrStdin(file)
			if err != nil {
				return err
			}
			payload := map[string]any{"content": string(content)}
			body, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			data, err := c.Post(cmdContext(cmd), "/confs/validate", nil, body)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", `pipeline YAML file ("-" for stdin)`)
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newPipelinesDeployCommandCmd() *cobra.Command {
	var ep extraParams
	cmd := &cobra.Command{
		Use:   "deploy-command <conf-id>",
		Short: "Get the install/deploy command for this pipeline's agents",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			if err := ep.apply(q); err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/pipelines/"+url.PathEscape(args[0])+"/deploy_command", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	ep.register(cmd)
	return cmd
}

func containsFold(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// latestVersion resolves the newest saved version (the epoch-millisecond
// timestamp, as a string) from a pipeline's history. Returns "" on any error.
func latestVersion(ctx context.Context, c *api.Client, confID string) string {
	data, err := c.Get(ctx, "/pipelines/"+url.PathEscape(confID)+"/history", nil)
	if err != nil {
		return ""
	}
	var entries []struct {
		Timestamp int64 `json:"timestamp"`
	}
	if err := json.Unmarshal(data, &entries); err != nil {
		return ""
	}
	var newest int64
	for _, e := range entries {
		if e.Timestamp > newest {
			newest = e.Timestamp
		}
	}
	if newest == 0 {
		return ""
	}
	return strconv.FormatInt(newest, 10)
}

// withVersionField annotates each history entry with a "version" field equal to
// its timestamp, so `--columns version` renders the deploy identifier. The input
// is returned unchanged if it isn't a JSON array of objects.
func withVersionField(data []byte) []byte {
	var entries []map[string]any
	if err := json.Unmarshal(data, &entries); err != nil {
		return data
	}
	for _, e := range entries {
		if _, ok := e["version"]; ok {
			continue
		}
		ts, ok := e["timestamp"]
		if !ok {
			continue
		}
		if f, ok := ts.(float64); ok {
			e["version"] = strconv.FormatInt(int64(f), 10)
		} else {
			e["version"] = fmt.Sprintf("%v", ts)
		}
	}
	out, err := json.Marshal(entries)
	if err != nil {
		return data
	}
	return out
}

// agentBeat is the subset of a pipeline agent record used to observe rollout.
type agentBeat struct {
	Host          string `json:"Host"`
	Identifier    string `json:"Identifier"`
	AgentVersion  string `json:"AgentVersion"`
	LastHeartBeat string `json:"LastHeartBeat"`
	Active        bool   `json:"active"`
}

func (a agentBeat) key() string {
	if a.Identifier != "" {
		return a.Identifier
	}
	return a.Host
}

func (a agentBeat) display() string {
	if a.Host != "" {
		return a.Host
	}
	return a.Identifier
}

func fetchAgentBeats(ctx context.Context, c *api.Client, confID string) ([]agentBeat, error) {
	data, err := c.Get(ctx, "/pipelines/"+url.PathEscape(confID)+"/agents", nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Agents []agentBeat `json:"agents"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Agents, nil
}

func parseBeat(s string) time.Time {
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, urlTimeFormat} {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

// waitForRollout polls the fleet after a deploy and confirms every active agent
// checks in (advances its heartbeat) afterwards. A heartbeat that never advances
// signals an agent that failed to apply the new config.
func waitForRollout(ctx context.Context, c *api.Client, confID string, timeout time.Duration) error {
	notef("Waiting up to %s for agents to check in after deploy…", timeout)
	baseline, err := fetchAgentBeats(ctx, c, confID)
	if err != nil {
		warnf("could not read fleet status to confirm rollout: %v", err)
		return nil
	}
	base := map[string]time.Time{}
	var activeKeys []string
	for _, a := range baseline {
		if !a.Active {
			continue
		}
		base[a.key()] = parseBeat(a.LastHeartBeat)
		activeKeys = append(activeKeys, a.key())
	}
	if len(activeKeys) == 0 {
		warnf("no active agents on this pipeline to observe; skipping rollout check")
		return nil
	}

	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		if cur, err := fetchAgentBeats(ctx, c, confID); err == nil {
			advanced := 0
			for _, a := range cur {
				if !a.Active {
					continue
				}
				if b, ok := base[a.key()]; ok && parseBeat(a.LastHeartBeat).After(b) {
					advanced++
				}
			}
			if advanced >= len(activeKeys) {
				fmt.Fprintf(os.Stderr, "%s All %d agent(s) checked in after the deploy.\n", okMark(), len(activeKeys))
				reportHealthProblems(ctx, c)
				return nil
			}
		}
		if time.Now().After(deadline) {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}

	// Timed out: report which agents never checked in.
	cur, _ := fetchAgentBeats(ctx, c, confID)
	curByKey := map[string]agentBeat{}
	for _, a := range cur {
		curByKey[a.key()] = a
	}
	stale := 0
	for _, k := range activeKeys {
		a := curByKey[k]
		if !parseBeat(a.LastHeartBeat).After(base[k]) {
			stale++
			warnf("  agent %s has not checked in since deploy (last heartbeat %s)", a.display(), a.LastHeartBeat)
		}
	}
	warnf("A frozen heartbeat usually means the agent failed to apply the new config")
	warnf("(e.g. a new source whose port is already bound fails the whole pipeline).")
	warnf("Check `edx health problems` and the agent logs.")
	return fmt.Errorf("rollout not confirmed within %s: %d/%d agent(s) did not check in", timeout, stale, len(activeKeys))
}

// reportHealthProblems best-effort surfaces any currently-failing components.
// The endpoint can include already-resolved ("ok") entries, so count only the
// components whose status is not healthy.
func reportHealthProblems(ctx context.Context, c *api.Client) {
	data, err := c.Get(ctx, "/health/problems", nil)
	if err != nil {
		return
	}
	var resp struct {
		Items []struct {
			Status string `json:"status"`
		} `json:"items"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return
	}
	unhealthy := 0
	for _, it := range resp.Items {
		switch strings.ToLower(it.Status) {
		case "", "ok", "healthy", "green":
		default:
			unhealthy++
		}
	}
	if unhealthy > 0 {
		warnf("%d component problem(s) currently reported — see `edx health problems`", unhealthy)
	}
}
