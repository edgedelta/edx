package cli

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func newDashboardsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dashboards",
		Short: "Manage dashboards",
		Long: `List, inspect, create, update and delete dashboards.

Dashboards are defined by a JSON body with a "definition" object. Tip: fetch an
existing dashboard with "edx dashboards get <id>" to use as a starting template.
"create" and "update" validate the definition client-side to catch the common
mistakes that make a dashboard save via the API but fail to render in the UI.`,
	}
	cmd.AddCommand(
		newDashboardsListCmd(),
		newDashboardsGetCmd(),
		newDashboardsCreateCmd(),
		newDashboardsUpdateCmd(),
		newDashboardsDeleteCmd(),
	)
	return cmd
}

func newDashboardsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List dashboards",
		Example: `  edx dashboards list --output table --columns dashboard_id,dashboard_name,creator`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/dashboards", nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

func newDashboardsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <dashboard-id>",
		Short: "Get a dashboard definition",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/dashboards/"+url.PathEscape(args[0]), nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

func newDashboardsCreateCmd() *cobra.Command {
	var file string
	var skipValidation bool
	cmd := &cobra.Command{
		Use:   "create --file dashboard.json",
		Short: "Create a dashboard from a JSON definition",
		Example: `  edx dashboards create --file dashboard.json
  edx dashboards get <id> | jq '.dashboard_name="copy"' | edx dashboards create --file -`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := readFileOrStdin(file)
			if err != nil {
				return err
			}
			if err := checkDashboard(body, skipValidation); err != nil {
				return err
			}
			data, err := c.Post(cmdContext(cmd), "/dashboards", nil, body)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", `dashboard JSON file ("-" for stdin)`)
	_ = cmd.MarkFlagRequired("file")
	cmd.Flags().BoolVar(&skipValidation, "skip-validation", false, "skip client-side definition checks")
	return cmd
}

func newDashboardsUpdateCmd() *cobra.Command {
	var file string
	var skipValidation bool
	cmd := &cobra.Command{
		Use:   "update <dashboard-id> --file dashboard.json",
		Short: "Update a dashboard from a JSON definition",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := readFileOrStdin(file)
			if err != nil {
				return err
			}
			if err := checkDashboard(body, skipValidation); err != nil {
				return err
			}
			data, err := c.Put(cmdContext(cmd), "/dashboards/"+url.PathEscape(args[0]), nil, body)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", `dashboard JSON file ("-" for stdin)`)
	_ = cmd.MarkFlagRequired("file")
	cmd.Flags().BoolVar(&skipValidation, "skip-validation", false, "skip client-side definition checks")
	return cmd
}

func newDashboardsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <dashboard-id>",
		Short: "Delete a dashboard",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm(fmt.Sprintf("Delete dashboard %s?", args[0])) {
				return errAborted
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Delete(cmdContext(cmd), "/dashboards/"+url.PathEscape(args[0]), nil, nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

// checkDashboard runs client-side validation, printing warnings to stderr and
// returning an error for hard failures unless skip is set.
func checkDashboard(body []byte, skip bool) error {
	if skip {
		return nil
	}
	errs, warns := validateDashboard(body)
	for _, w := range warns {
		warnf("%s", w)
	}
	if len(errs) > 0 {
		for _, e := range errs {
			warnf("%s", e)
		}
		return fmt.Errorf("dashboard definition failed validation (%d issue(s)); fix, or pass --skip-validation to override", len(errs))
	}
	return nil
}

// validateDashboard catches the common definition mistakes that let a dashboard
// save via the API yet fail to render in the UI: a widget-schema/version
// mismatch, and missing resource_accesses.
func validateDashboard(body []byte) (errs, warns []string) {
	var d struct {
		Definition struct {
			Version int              `json:"version"`
			Widgets []map[string]any `json:"widgets"`
		} `json:"definition"`
		ResourceAccesses []any `json:"resource_accesses"`
	}
	if err := json.Unmarshal(body, &d); err != nil {
		return []string{fmt.Sprintf("invalid JSON: %v", err)}, nil
	}
	if len(d.Definition.Widgets) == 0 {
		warns = append(warns, "definition.widgets is empty — the dashboard will render blank")
	}
	usesVisuals := false
	vizCount := 0
	for _, w := range d.Definition.Widgets {
		if t, _ := w["type"].(string); t == "viz" {
			vizCount++
		}
		if _, ok := w["visuals"]; ok {
			usesVisuals = true
		}
	}
	if usesVisuals && d.Definition.Version != 3 {
		errs = append(errs, fmt.Sprintf(
			"widgets use the visuals[] schema, which requires definition.version 3 (got %d); the UI reports \"Dashboard could not be found\" on a version mismatch",
			d.Definition.Version))
	}
	if vizCount > 0 && len(d.ResourceAccesses) == 0 {
		warns = append(warns, "resource_accesses is empty; the UI may fail to resolve the dashboard. Add one {\"domain\":...,\"query\":...} entry per widget query.")
	}
	return errs, warns
}
