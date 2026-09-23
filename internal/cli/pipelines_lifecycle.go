package cli

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

func newPipelinesCreateCmd() *cobra.Command {
	var file, tag, environment, subtype, description, runtime, runtimeVersion string
	cmd := &cobra.Command{
		Use:   "create --file pipeline.yaml --tag <tag> --environment <environment>",
		Short: "Create a new edge pipeline from YAML",
		Long:  "Create an edge pipeline and its initial deployed configuration. This does not install an agent. Check existing pipelines before creating to avoid duplicates.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(tag) == "" {
				return fmt.Errorf("tag must not be empty")
			}
			switch environment {
			case "Docker", "Linux", "Windows", "MacOS":
			case "Kubernetes":
				if subtype == "" {
					return fmt.Errorf("--fleet-subtype is required for Kubernetes")
				}
			default:
				return fmt.Errorf("unsupported edge environment %q", environment)
			}
			if runtime != "go" && runtime != "red" {
				return fmt.Errorf("runtime must be go or red")
			}
			if runtimeVersion != "" && runtime != "red" {
				return fmt.Errorf("runtime-version requires runtime red")
			}
			if subtype != "" && subtype != "Edge" && subtype != "Coordinator" && subtype != "Gateway" {
				return fmt.Errorf("unsupported fleet subtype %q", subtype)
			}
			content, err := readFileOrStdin(file)
			if err != nil {
				return err
			}
			if strings.TrimSpace(string(content)) == "" {
				return fmt.Errorf("pipeline YAML must not be empty")
			}
			body, err := json.Marshal(map[string]string{
				"content": string(content), "tag": tag, "environment": environment,
				"fleet_type": "Edge", "fleet_subtype": subtype, "description": description,
				"runtime": runtime, "runtime_version": runtimeVersion,
			})
			if err != nil {
				return err
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			// Creation has no idempotency key; a retry could create a second pipeline.
			c.MaxRetries = 0
			data, err := c.Post(cmdContext(cmd), "/confs", nil, body)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "pipeline YAML file (- for stdin)")
	cmd.Flags().StringVar(&tag, "tag", "", "pipeline tag")
	cmd.Flags().StringVar(&environment, "environment", "", "Docker, Kubernetes, Linux, Windows or MacOS")
	cmd.Flags().StringVar(&subtype, "fleet-subtype", "", "Edge, Coordinator or Gateway (required for Kubernetes)")
	cmd.Flags().StringVarP(&description, "description", "d", "", "pipeline description")
	cmd.Flags().StringVar(&runtime, "runtime", "go", "agent runtime: go or red")
	cmd.Flags().StringVar(&runtimeVersion, "runtime-version", "", "RED runtime release")
	for _, name := range []string{"file", "tag", "environment"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			cmd.PrintErrln(err)
		}
	}
	return cmd
}

func newPipelinesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <conf-id>",
		Short: "Delete a pipeline configuration (does not uninstall agents)",
		Long:  "Stop or uninstall agents before deleting their pipeline. Hosted fleets must be removed through their own lifecycle first.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm(fmt.Sprintf("Delete pipeline %s? This does not uninstall its agents.", args[0])) {
				return errAborted
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Delete(cmdContext(cmd), "/confs/"+url.PathEscape(args[0]), nil, nil)
			if err != nil {
				return err
			}
			if len(data) == 0 {
				data = []byte(`{"deleted":true}`)
			}
			return printResult(data)
		},
	}
}
