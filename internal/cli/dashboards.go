package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"

	"github.com/edgedelta/edx/internal/dashboards"
)

func newDashboardsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dashboards",
		Short: "Manage dashboards",
		Long: `List, inspect, create, update, validate, render and delete dashboards.

Dashboards are defined by a JSON body with a "definition" object. Tip: fetch an
existing dashboard with "edx dashboards get <id>" to use as a starting template.
"create" and "update" validate the definition against the frontend's own Dashboard
schema first, so a definition that would save via the API yet fail to render in the
UI is rejected here instead. Use "validate" to check a definition without saving.

Validation is offline and instant; it cannot tell you whether a dashboard looks
right or whether its widgets return data. "screenshot" answers that by rendering the
saved dashboard in a real browser, which is slow enough to be worth running as a
background task. Together they give the loop:

  validate  ->  create --tag generated --tag preview  ->  screenshot  ->  look
            ->  update  ->  screenshot  ->  tag --remove preview`,
	}
	cmd.AddCommand(
		newDashboardsListCmd(),
		newDashboardsGetCmd(),
		newDashboardsCreateCmd(),
		newDashboardsUpdateCmd(),
		newDashboardsValidateCmd(),
		newDashboardsScreenshotCmd(),
		newDashboardsTagCmd(),
		newDashboardsOptionsCmd(),
		newDashboardsDeleteCmd(),
	)
	return cmd
}

func newDashboardsTagCmd() *cobra.Command {
	var add, remove []string
	cmd := &cobra.Command{
		Use:   "tag <dashboard-id>",
		Short: "Add or remove tags on a dashboard",
		Long: `Add or remove tags on an existing dashboard without touching its definition.

Tags are how a dashboard says where it came from and whether anyone has looked at
it. Two conventions worth following when generating dashboards:

  generated  written by a tool or an agent rather than by hand
  preview    not verified yet — carry it from creation until a rendered screenshot
             has been checked, then remove it

The dashboard is read and written back whole, so tags are the only thing that
changes. resource_accesses is left exactly as it was, since the definition is
unchanged.`,
		Example: `  edx dashboards tag <id> --add generated --add preview
  edx dashboards tag <id> --remove preview
  edx dashboards list | jq -r '.[] | select(.tags | index("preview")) | .dashboard_id'`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(add) == 0 && len(remove) == 0 {
				return errors.New("nothing to do: pass --add and/or --remove")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			ctx := cmdContext(cmd)
			path := "/dashboards/" + url.PathEscape(args[0])

			current, err := c.Get(ctx, path, nil)
			if err != nil {
				return err
			}
			updated, err := applyTags(current, add, remove)
			if err != nil {
				return err
			}
			data, err := c.Put(ctx, path, nil, updated)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringArrayVar(&add, "add", nil, "tag to add (repeatable)")
	cmd.Flags().StringArrayVar(&remove, "remove", nil, "tag to remove (repeatable)")
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
	var tags []string
	cmd := &cobra.Command{
		Use:   "create --file dashboard.json",
		Short: "Create a dashboard from a JSON definition",
		Long: `Create a dashboard from a JSON body, after running the same checks as
"validate".

resource_accesses is derived from the definition and added for you, replacing
anything in the file, exactly as the UI does when a dashboard is saved. It is the
allowlist public share links and screenshots are authorized against, and it is
fully determined by the widgets and variables, so there is no reason to write it
by hand.

Tag what you create. A dashboard written by a tool should say so with "generated",
and one nobody has looked at yet should carry "preview" until a rendered screenshot
has been checked (see "edx dashboards screenshot" and "edx dashboards tag").`,
		Example: `  edx dashboards create --file dashboard.json --tag generated --tag preview
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
			body, err = applyTags(body, tags, nil)
			if err != nil {
				return err
			}
			body, err = fillResourceAccesses(body)
			if err != nil {
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
	cmd.Flags().StringArrayVar(&tags, "tag", nil, `tag to add, on top of any in the file (repeatable); "generated" and "preview" are the conventions`)
	return cmd
}

func newDashboardsUpdateCmd() *cobra.Command {
	var file string
	var skipValidation bool
	var tags []string
	cmd := &cobra.Command{
		Use:   "update <dashboard-id> --file dashboard.json",
		Short: "Update a dashboard from a JSON definition",
		Long: `Update a dashboard from a JSON body, after running the same checks as
"validate".

resource_accesses is regenerated from the definition, so a body round-tripped
through "get" carries a fresh allowlist rather than the one it was fetched with.

The whole record is replaced, so a body that omits "tags" clears them: pass --tag
to keep them, or use "edx dashboards tag" to change tags on their own.`,
		Args: cobra.ExactArgs(1),
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
			body, err = applyTags(body, tags, nil)
			if err != nil {
				return err
			}
			body, err = fillResourceAccesses(body)
			if err != nil {
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
	cmd.Flags().StringArrayVar(&tags, "tag", nil, "tag to add, on top of any in the file (repeatable)")
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
// fillResourceAccesses derives resource_accesses from the definition and writes it into
// the request body, the same way the UI does when a dashboard is saved.
//
// It is the allowlist that public share links and screenshots are authorized against, and
// it is fully determined by the definition, so nobody authoring a dashboard should have to
// write it. Any value already in the body is replaced — again matching the UI, which
// regenerates the whole array on every save.
//
// A body without a "definition" object is passed through untouched: `validate` accepts a
// bare definition, but create and update always send a full body, and inventing a
// top-level key for something else's payload would be wrong.
func fillResourceAccesses(body []byte) ([]byte, error) {
	var envelope map[string]any
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	definition, ok := envelope["definition"]
	if !ok {
		return body, nil
	}

	accesses := dashboards.ResourceAccesses(definition)
	if accesses == nil {
		// Encode as [] rather than null: the field is a list, and a null would read as a
		// different thing from "nothing to allow".
		accesses = []dashboards.ResourceAccess{}
	}
	envelope["resource_accesses"] = accesses

	filled, err := json.Marshal(envelope)
	if err != nil {
		return nil, fmt.Errorf("re-encode dashboard body: %w", err)
	}
	return filled, nil
}

// applyTags rewrites the body's top-level "tags" array. A body with nothing to change is
// returned untouched, so the request carries exactly what the caller wrote.
func applyTags(body []byte, add, remove []string) ([]byte, error) {
	if len(add) == 0 && len(remove) == 0 {
		return body, nil
	}
	// UseNumber so a round-tripped definition keeps its numeric literals: widget IDs and
	// thresholds must come back the way they went in.
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()
	var envelope map[string]any
	if err := dec.Decode(&envelope); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	tags := mergeTags(stringList(envelope["tags"]), add, remove)
	if len(tags) == 0 {
		delete(envelope, "tags")
	} else {
		envelope["tags"] = tags
	}

	tagged, err := json.Marshal(envelope)
	if err != nil {
		return nil, fmt.Errorf("re-encode dashboard body: %w", err)
	}
	return tagged, nil
}

// mergeTags applies the additions and removals, keeping the existing order and never
// duplicating a tag.
func mergeTags(existing, add, remove []string) []string {
	drop := make(map[string]bool, len(remove))
	for _, t := range remove {
		drop[strings.TrimSpace(t)] = true
	}
	seen := make(map[string]bool, len(existing)+len(add))
	out := make([]string, 0, len(existing)+len(add))
	for _, t := range append(append([]string{}, existing...), add...) {
		t = strings.TrimSpace(t)
		if t == "" || drop[t] || seen[t] {
			continue
		}
		seen[t] = true
		out = append(out, t)
	}
	return out
}

// stringList reads a decoded JSON array of strings, ignoring anything else in it.
func stringList(v any) []string {
	items, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(items))
	for _, item := range items {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

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
