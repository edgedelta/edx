package dashboards

import (
	"fmt"
	"sort"
)

// Options lists the accepted values for a definition's discriminated fields, read out of
// the embedded schema rather than hand-maintained, so they cannot drift from what
// validation enforces.
type Options struct {
	// WidgetTypes are the values for a widget's `type`.
	WidgetTypes []string `json:"widgetTypes"`
	// VisualizerTypes are the values for a viz widget's `visualizer.type`.
	VisualizerTypes []string `json:"visualizerTypes"`
	// DataTypes are the values for a `dataSource.type`.
	DataTypes []string `json:"dataTypes"`
	// VariableTypes are the values for a variable's `type`.
	VariableTypes []string `json:"variableTypes"`
	// PositionTypes are the values for a widget's `position.type`.
	PositionTypes []string `json:"positionTypes"`
	// ResultTypes are the values for a viz widget's `resultType`.
	ResultTypes []string `json:"resultTypes"`
	// Version is the definition version these options describe.
	Version int `json:"version"`
}

// SchemaOptions reads the accepted values for each discriminated field out of the
// embedded schema.
//
// Every list comes from the schema document, so regenerating the schema updates them with
// no code change here. A field whose values cannot be found is an error rather than an
// empty list: silently offering no options would read as "this field accepts nothing".
func SchemaOptions() (*Options, error) {
	doc, err := schemaDocument()
	if err != nil {
		return nil, err
	}

	opts := &Options{Version: SupportedVersion}
	for _, field := range []struct {
		name   string
		target *[]string
		read   func(map[string]any) ([]string, error)
	}{
		{"widget types", &opts.WidgetTypes, func(d map[string]any) ([]string, error) {
			return discriminantValues(d, "DashboardsV4.Widget")
		}},
		{"visualizer types", &opts.VisualizerTypes, func(d map[string]any) ([]string, error) {
			return discriminantValues(d, "DashboardsV4.VizOptions")
		}},
		{"data types", &opts.DataTypes, func(d map[string]any) ([]string, error) {
			return discriminantValues(d, "DashboardsV4.DataSource")
		}},
		{"variable types", &opts.VariableTypes, func(d map[string]any) ([]string, error) {
			return discriminantValues(d, "DashboardsV4.Variable")
		}},
		{"position types", &opts.PositionTypes, func(d map[string]any) ([]string, error) {
			return discriminantValues(d, "DashboardsV4.Position")
		}},
		{"result types", &opts.ResultTypes, func(d map[string]any) ([]string, error) {
			return enumValues(d, "DashboardsV4.ResultType")
		}},
	} {
		values, err := field.read(doc)
		if err != nil {
			return nil, fmt.Errorf("read %s from the embedded schema: %w", field.name, err)
		}
		*field.target = values
	}

	return opts, nil
}

// discriminantValues collects the `type` values a definition's variants accept.
//
// A definition is either a union of variants (anyOf/oneOf) or a single object, and a
// variant's `type` is a const, an enum, or a $ref to one. Variants can also be $refs to
// further unions, so this recurses.
func discriminantValues(doc map[string]any, name string) ([]string, error) {
	definition, err := definitionByName(doc, name)
	if err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	var walk func(node map[string]any, depth int) error
	walk = func(node map[string]any, depth int) error {
		// The schema is generated, not arbitrary input, but bound the recursion anyway so
		// a future cyclic $ref cannot hang the command.
		if depth > 20 {
			return fmt.Errorf("%s: $ref chain deeper than 20 levels", name)
		}

		if ref, ok := node["$ref"].(string); ok {
			target, err := definitionByRef(doc, ref)
			if err != nil {
				return err
			}
			return walk(target, depth+1)
		}

		for _, key := range []string{"anyOf", "oneOf"} {
			if branches, ok := node[key].([]any); ok {
				for _, branch := range branches {
					b, ok := branch.(map[string]any)
					if !ok {
						continue
					}
					if err := walk(b, depth+1); err != nil {
						return err
					}
				}
				return nil
			}
		}

		properties, ok := node["properties"].(map[string]any)
		if !ok {
			return nil
		}
		discriminant, ok := properties["type"].(map[string]any)
		if !ok {
			return nil
		}
		return collectValues(doc, discriminant, seen, depth)
	}

	if err := walk(definition, 0); err != nil {
		return nil, err
	}
	if len(seen) == 0 {
		return nil, fmt.Errorf("%s: no `type` values found", name)
	}
	return sorted(seen), nil
}

// collectValues reads a const, an enum, or a $ref to either, into seen.
func collectValues(doc, node map[string]any, seen map[string]bool, depth int) error {
	if depth > 20 {
		return fmt.Errorf("$ref chain deeper than 20 levels")
	}
	if ref, ok := node["$ref"].(string); ok {
		target, err := definitionByRef(doc, ref)
		if err != nil {
			return err
		}
		return collectValues(doc, target, seen, depth+1)
	}
	if constant, ok := node["const"].(string); ok {
		seen[constant] = true
		return nil
	}
	for _, v := range asStrings(node["enum"]) {
		seen[v] = true
	}
	return nil
}

// enumValues reads a definition that is a plain string enum, such as ResultType.
func enumValues(doc map[string]any, name string) ([]string, error) {
	definition, err := definitionByName(doc, name)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	if err := collectValues(doc, definition, seen, 0); err != nil {
		return nil, err
	}
	if len(seen) == 0 {
		return nil, fmt.Errorf("%s: not an enum", name)
	}
	return sorted(seen), nil
}

func definitionByName(doc map[string]any, name string) (map[string]any, error) {
	definitions, ok := doc["definitions"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("schema has no definitions object")
	}
	definition, ok := definitions[name].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("schema has no definition %q", name)
	}
	return definition, nil
}

// definitionByRef resolves a local "#/definitions/<name>" pointer.
func definitionByRef(doc map[string]any, ref string) (map[string]any, error) {
	const prefix = "#/definitions/"
	if len(ref) <= len(prefix) || ref[:len(prefix)] != prefix {
		return nil, fmt.Errorf("unsupported $ref %q", ref)
	}
	return definitionByName(doc, ref[len(prefix):])
}

func asStrings(v any) []string {
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

func sorted(set map[string]bool) []string {
	out := make([]string, 0, len(set))
	for v := range set {
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}
