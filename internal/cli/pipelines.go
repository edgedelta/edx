package cli

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
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
		Example: `  edx pipelines history <conf-id> --output table --columns version,description,creator,created`,
		Args: cobra.ExactArgs(1),
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
			return printResult(data)
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
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", `pipeline YAML file ("-" for stdin)`)
	_ = cmd.MarkFlagRequired("file")
	cmd.Flags().StringVarP(&description, "description", "d", "", "change description shown in history")
	return cmd
}

func newPipelinesDeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy <conf-id> <version>",
		Short: "Deploy a saved pipeline version to the fleet",
		Long: `Deploy a specific saved version to all agents in the fleet. Deploying an
older version is the supported way to roll back.`,
		Example: `  edx pipelines deploy <conf-id> 42 --yes`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm(fmt.Sprintf("Deploy version %s of pipeline %s to the fleet?", args[1], args[0])) {
				return errAborted
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			path := "/pipelines/" + url.PathEscape(args[0]) + "/deploy/" + url.PathEscape(args[1])
			data, err := c.Post(cmdContext(cmd), path, nil, nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
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
		Use:   "validate --file pipeline.yaml",
		Short: "Validate an agent configuration without saving it",
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
