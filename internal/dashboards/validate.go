// Package dashboards validates dashboard definitions against the frontend's own
// Dashboard type, deterministically and offline.
//
// dashboard-v4.schema.json is generated from web/src/modules/dashboards/versions/v4
// in the edgedelta monorepo, where the TypeScript type is the source of truth; a
// type-level drift guard there fails the web typecheck if the two diverge. Refresh
// the vendored copy with `make sync-dashboard-schema`.
package dashboards

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// printer renders jsonschema's error kinds; LocalizedString panics on a nil printer.
var printer = message.NewPrinter(language.English)

//go:embed dashboard-v4.schema.json
var schemaJSON []byte

// SupportedVersion is the definition version the embedded schema describes.
const SupportedVersion = 4

var (
	compileOnce sync.Once
	compiled    *jsonschema.Schema
	compileErr  error
)

func schema() (*jsonschema.Schema, error) {
	compileOnce.Do(func() {
		doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(schemaJSON))
		if err != nil {
			compileErr = fmt.Errorf("parse embedded dashboard schema: %w", err)
			return
		}
		c := jsonschema.NewCompiler()
		if err := c.AddResource("dashboard-v4.schema.json", doc); err != nil {
			compileErr = fmt.Errorf("load embedded dashboard schema: %w", err)
			return
		}
		compiled, compileErr = c.Compile("dashboard-v4.schema.json")
	})
	return compiled, compileErr
}

// Issue is a single problem found in a definition.
type Issue struct {
	// Path is a JSON Pointer into the definition, e.g. "/widgets/2/visualizer/type".
	Path    string
	Message string
}

func (i Issue) String() string {
	if i.Path == "" {
		return i.Message
	}
	return i.Path + ": " + i.Message
}

// ErrUnsupportedVersion is returned when the definition is not version 4. Other
// versions are migrated in the UI and have no schema here, so validating them
// against the v4 schema would produce misleading errors.
var ErrUnsupportedVersion = errors.New("unsupported definition version")

// ValidateDefinition checks a bare dashboard *definition* object (the value of the
// "definition" key) against the v4 schema. A nil error means the definition
// satisfies the frontend's Dashboard type.
func ValidateDefinition(definition any) ([]Issue, error) {
	sch, err := schema()
	if err != nil {
		return nil, err
	}

	if obj, ok := definition.(map[string]any); ok {
		if v, present := obj["version"]; present {
			// JSON numbers decode as float64.
			if n, isNum := v.(float64); isNum && int(n) != SupportedVersion {
				return nil, fmt.Errorf("%w: got %d, edx can only validate version %d",
					ErrUnsupportedVersion, int(n), SupportedVersion)
			}
		}
	}

	err = sch.Validate(definition)
	if err == nil {
		return nil, nil
	}
	var ve *jsonschema.ValidationError
	if !errors.As(err, &ve) {
		return nil, err
	}
	return issues(ve), nil
}

// issues flattens a validation error tree into the deepest causes, which carry the
// actionable message, then dedupes and sorts them by path.
//
// Every widget, visualizer, variable and data-source variant is a branch of one
// anyOf keyed on `type`, so one bad field fails every branch: a mistyped
// `displayOptions` key also reports "missing property 'grid'" from the grid branch
// and five more like it. resolveVariant discards the branches whose discriminant did
// not match, which leaves the errors that describe the variant the author meant.
func issues(ve *jsonschema.ValidationError) []Issue {
	seen := map[string]bool{}
	var out []Issue

	add := func(e *jsonschema.ValidationError) {
		path := "/" + strings.Join(e.InstanceLocation, "/")
		if path == "/" {
			path = ""
		}
		iss := Issue{Path: path, Message: e.ErrorKind.LocalizedString(printer)}
		if key := iss.String(); !seen[key] {
			seen[key] = true
			out = append(out, iss)
		}
	}

	var walk func(e *jsonschema.ValidationError)
	walk = func(e *jsonschema.ValidationError) {
		if len(e.Causes) == 0 {
			add(e)
			return
		}
		for _, c := range resolveVariant(e) {
			walk(c)
		}
	}
	walk(ve)

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Path != out[j].Path {
			return out[i].Path < out[j].Path
		}
		return out[i].Message < out[j].Message
	})
	return out
}

// resolveVariant narrows the causes of a discriminated anyOf/oneOf to the branch the
// instance was actually trying to be.
//
// A branch that reports a const/enum failure on the node's own `type` property is the
// wrong variant, so its other complaints are noise. If some branches match the
// discriminant, only those are returned. If none do, `type` itself is the problem and
// the discriminant failures — which enumerate the accepted values — are returned.
// Non-union nodes and unions without a usable discriminant pass through untouched.
func resolveVariant(e *jsonschema.ValidationError) []*jsonschema.ValidationError {
	switch e.ErrorKind.(type) {
	case *kind.AnyOf, *kind.OneOf:
	default:
		return e.Causes
	}

	discriminant := append(append([]string{}, e.InstanceLocation...), "type")

	var matched, discriminantFailures []*jsonschema.ValidationError
	for _, branch := range e.Causes {
		if failures := discriminantFailuresIn(branch, discriminant); len(failures) > 0 {
			discriminantFailures = append(discriminantFailures, failures...)
			continue
		}
		matched = append(matched, branch)
	}

	if len(matched) > 0 {
		return matched
	}
	if len(discriminantFailures) > 0 {
		return discriminantFailures
	}
	return e.Causes
}

// discriminantFailuresIn finds const/enum failures on exactly the given instance path
// within a branch, i.e. "this branch wanted a different `type`".
func discriminantFailuresIn(e *jsonschema.ValidationError, path []string) []*jsonschema.ValidationError {
	var found []*jsonschema.ValidationError

	var walk func(n *jsonschema.ValidationError)
	walk = func(n *jsonschema.ValidationError) {
		switch n.ErrorKind.(type) {
		case *kind.Const, *kind.Enum:
			if samePath(n.InstanceLocation, path) {
				found = append(found, n)
				return
			}
		}
		for _, c := range n.Causes {
			walk(c)
		}
	}
	walk(e)
	return found
}

func samePath(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
