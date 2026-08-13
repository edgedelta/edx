package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/edgedelta/edx/internal/dashboards"
)

func newDashboardsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dashboards",
		Short: "Manage dashboards",
		Long: `List, inspect, create, update, validate and delete dashboards.

Dashboards are defined by a JSON body with a "definition" object. Tip: fetch an
existing dashboard with "edx dashboards get <id>" to use as a starting template.
"create" and "update" validate the definition against the frontend's own Dashboard
schema first, so a definition that would save via the API yet fail to render in the
UI is rejected here instead. Use "validate" to check a definition without saving.`,
	}
	cmd.AddCommand(
		newDashboardsListCmd(),
		newDashboardsGetCmd(),
		newDashboardsCreateCmd(),
		newDashboardsUpdateCmd(),
		newDashboardsValidateCmd(),
		newDashboardsOptionsCmd(),
		newDashboardsDeleteCmd(),
	)
	return cmd
}

func newDashboardsOptionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "options",
		Aliases: []string{"types"},
		Short:   "List the accepted widget, visualizer and data types",
		Long: `List every accepted value for a definition's typed fields: widget types,
visualizer types, data source types, variable types, position types and result
types.

Read out of the same embedded schema that "validate" enforces, so this is always
in step with what validation accepts — useful when picking a visualizer or
diagnosing a rejected "type".`,
		Example: `  edx dashboards options
  edx dashboards options --output yaml
  edx dashboards options | jq -r '.visualizerTypes[]'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := dashboards.SchemaOptions()
			if err != nil {
				return err
			}
			data, err := json.Marshal(opts)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

func newDashboardsValidateCmd() *cobra.Command {
	var file string
	cmd := &cobra.Command{
		Use:   "validate --file dashboard.json",
		Short: "Validate a dashboard definition without saving it",
		Long: `Validate a dashboard JSON body, or a bare definition object, against the
frontend's Dashboard schema, and check the syntax of the queries inside it. Exits
non-zero when the definition would not render.

Queries are parsed with the same ANTLR grammars the backend uses, choosing the
grammar from each data source's type — metric, log and formula queries have
different syntax. This is a syntax check only: telling whether a facet or metric
name exists needs a live backend.

Accepts either a full dashboard body ({"dashboard_name": ..., "definition": {...}})
or just the definition object ({"version": 4, "widgets": [...]}).`,
		Example: `  edx dashboards validate --file dashboard.json
  edx dashboards get <id> | edx dashboards validate --file -`,
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := readFileOrStdin(file)
			if err != nil {
				return err
			}
			if err := checkDashboard(body, false); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "definition is valid")
			return nil
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", `dashboard JSON file ("-" for stdin)`)
	_ = cmd.MarkFlagRequired("file")
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

// validateDashboard checks the definition against the frontend's Dashboard schema
// (see internal/dashboards), which is generated from the TypeScript type the UI
// itself renders from. Warnings cover things the schema cannot express.
func validateDashboard(body []byte) (errs, warns []string) {
	var d struct {
		Definition       map[string]any `json:"definition"`
		ResourceAccesses []any          `json:"resource_accesses"`
	}
	if err := json.Unmarshal(body, &d); err != nil {
		return []string{fmt.Sprintf("invalid JSON: %v", err)}, nil
	}

	definition := d.Definition
	if definition == nil {
		// Also accept a bare definition object, which is what "validate" is usually
		// handed while a definition is still being drafted.
		var bare map[string]any
		if err := json.Unmarshal(body, &bare); err != nil {
			return []string{fmt.Sprintf("invalid JSON: %v", err)}, nil
		}
		if _, hasWidgets := bare["widgets"]; !hasWidgets {
			return []string{`no "definition" object found, and the body is not a bare definition (no "widgets" key)`}, nil
		}
		definition = bare
	}

	widgets, _ := definition["widgets"].([]any)
	if len(widgets) == 0 {
		warns = append(warns, "definition.widgets is empty — the dashboard will render blank")
	}
	vizCount := 0
	for _, w := range widgets {
		m, ok := w.(map[string]any)
		if !ok {
			continue
		}
		if t, _ := m["type"].(string); t == "viz" {
			vizCount++
		}
	}
	if vizCount > 0 && len(d.ResourceAccesses) == 0 && d.Definition != nil {
		warns = append(warns, `resource_accesses is empty; the UI may fail to resolve the dashboard. Add one {"domain":...,"query":...} entry per widget query.`)
	}

	issues, err := dashboards.ValidateDefinition(definition)
	if err != nil {
		if errors.Is(err, dashboards.ErrUnsupportedVersion) {
			// v3 and older are migrated in the UI; edx only ships the v4 schema.
			warns = append(warns, fmt.Sprintf("%v; skipping schema validation", err))
			return errs, warns
		}
		return append(errs, err.Error()), warns
	}
	for _, iss := range issues {
		errs = append(errs, iss.String())
	}
	return errs, warns
}
