// Package output renders API responses as json, yaml, table or csv.
package output

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

// Options controls rendering.
type Options struct {
	// Format is one of: json (default), yaml, table, csv, raw.
	Format string
	// Columns optionally restricts/orders table and csv columns. Values are
	// dot-paths into each row object (e.g. "resource.host.name").
	Columns []string
}

// Print renders raw JSON response bytes to w in the requested format.
func Print(w io.Writer, data []byte, opts Options) error {
	switch opts.Format {
	case "", "json":
		return printJSON(w, data)
	case "raw":
		_, err := w.Write(append(data, '\n'))
		return err
	case "yaml":
		return printYAML(w, data)
	case "table":
		return printTabular(w, data, opts.Columns, "table")
	case "csv":
		return printTabular(w, data, opts.Columns, "csv")
	default:
		return fmt.Errorf("unknown output format %q (expected json, yaml, table, csv or raw)", opts.Format)
	}
}

func printJSON(w io.Writer, data []byte) error {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		// Not JSON (e.g. plain text or empty body): print as-is.
		_, werr := fmt.Fprintln(w, strings.TrimSpace(string(data)))
		return werr
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

func printYAML(w io.Writer, data []byte) error {
	v, err := decodeNumeric(data)
	if err != nil {
		// Not JSON (e.g. plain text or empty body): print as-is.
		_, werr := fmt.Fprintln(w, strings.TrimSpace(string(data)))
		return werr
	}
	out, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(out)
	return err
}

// decodeNumeric unmarshals JSON, decoding numbers as json.Number and then
// normalizing each to int64 (when integral) or float64. A plain unmarshal into
// any decodes every JSON number as float64, which makes yaml.Marshal render
// large integers — counts, epoch-millis timestamps — in lossy scientific
// notation (e.g. 1.135995e+06). Keeping integers as int64 avoids that.
func decodeNumeric(data []byte) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	return normalizeNumbers(v), nil
}

func normalizeNumbers(v any) any {
	switch t := v.(type) {
	case map[string]any:
		for k, val := range t {
			t[k] = normalizeNumbers(val)
		}
		return t
	case []any:
		for i, val := range t {
			t[i] = normalizeNumbers(val)
		}
		return t
	case json.Number:
		if i, err := t.Int64(); err == nil {
			return i
		}
		if f, err := t.Float64(); err == nil {
			return f
		}
		return t.String()
	default:
		return v
	}
}

// Rows extracts a list of row objects from an arbitrary API response. It
// understands the common Edge Delta envelope keys (items, stats, records, data).
func Rows(data []byte) ([]map[string]any, error) {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("response is not JSON: %w", err)
	}
	rows := extractRows(v, 0)
	if rows == nil {
		// Unwrap the AI services' {status, data: X} envelope. X is either a
		// thin collection wrapper ({"issues": [...]}, {"timeline": [...]}) whose
		// sole array we tabulate, or a single resource object — possibly with an
		// embedded sub-collection (e.g. a thread carrying its messages), which we
		// must still render as one row rather than descending into the embed.
		if m, ok := v.(map[string]any); ok {
			if inner, ok := m["data"].(map[string]any); ok {
				if arr := soleArrayValue(inner); arr != nil {
					return extractRows(arr, 0), nil
				}
				return []map[string]any{inner}, nil
			}
			return []map[string]any{m}, nil
		}
		return nil, fmt.Errorf("could not find a list of rows in the response; use --output json")
	}
	return rows, nil
}

// soleArrayValue returns the array when m is a thin collection wrapper — a
// single key whose value is an array (e.g. {"issues": [...]}). It returns nil
// for resource objects, so an embedded sub-collection is never mistaken for the
// payload.
func soleArrayValue(m map[string]any) []any {
	if len(m) != 1 {
		return nil
	}
	for _, v := range m {
		if arr, ok := v.([]any); ok {
			return arr
		}
	}
	return nil
}

func extractRows(v any, depth int) []map[string]any {
	if depth > 3 {
		return nil
	}
	switch t := v.(type) {
	case []any:
		rows := make([]map[string]any, 0, len(t))
		for _, item := range t {
			if m, ok := item.(map[string]any); ok {
				rows = append(rows, m)
			} else {
				rows = append(rows, map[string]any{"value": item})
			}
		}
		return rows
	case map[string]any:
		for _, key := range []string{"items", "stats", "records", "data", "agents", "deployments", "spans", "traces", "monitors", "dashboards", "options"} {
			if inner, ok := t[key]; ok {
				if rows := extractRows(inner, depth+1); rows != nil {
					return rows
				}
			}
		}
	}
	return nil
}

func printTabular(w io.Writer, data []byte, columns []string, kind string) error {
	rows, err := Rows(data)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		_, err := fmt.Fprintln(w, "(no results)")
		return err
	}

	flat := make([]map[string]string, len(rows))
	keySet := map[string]bool{}
	for i, row := range rows {
		f := map[string]string{}
		flatten("", row, f, 0)
		flat[i] = f
		for k := range f {
			keySet[k] = true
		}
	}

	cols := columns
	if len(cols) == 0 {
		cols = orderKeys(keySet)
	}

	if kind == "csv" {
		cw := csv.NewWriter(w)
		if err := cw.Write(cols); err != nil {
			return err
		}
		for _, f := range flat {
			rec := make([]string, len(cols))
			for j, c := range cols {
				rec[j] = f[c]
			}
			if err := cw.Write(rec); err != nil {
				return err
			}
		}
		cw.Flush()
		return cw.Error()
	}

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, strings.ToUpper(strings.Join(cols, "\t")))
	for _, f := range flat {
		vals := make([]string, len(cols))
		for j, c := range cols {
			vals[j] = truncate(f[c], 120)
		}
		fmt.Fprintln(tw, strings.Join(vals, "\t"))
	}
	return tw.Flush()
}

// flatten converts nested objects into dot-path keys with string values.
// Arrays and deeply nested objects are rendered as compact JSON.
func flatten(prefix string, v any, out map[string]string, depth int) {
	m, ok := v.(map[string]any)
	if !ok || depth > 2 {
		out[prefix] = compact(v)
		return
	}
	for k, val := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch val.(type) {
		case map[string]any:
			flatten(key, val, out, depth+1)
		default:
			out[key] = compact(val)
		}
	}
}

func compact(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return strings.ReplaceAll(strings.ReplaceAll(t, "\n", "\\n"), "\t", " ")
	case float64:
		if t == float64(int64(t)) {
			return fmt.Sprintf("%d", int64(t))
		}
		return fmt.Sprintf("%g", t)
	case bool:
		return fmt.Sprintf("%v", t)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(b)
	}
}

// orderKeys returns columns in a stable order, with well-known identity and
// time fields first so default tables read naturally.
func orderKeys(keys map[string]bool) []string {
	priority := []string{
		"id", "name", "tag", "title", "timestamp", "severity_text", "service.name",
		"status", "state", "type", "fleet_type", "environment", "version",
		"creator", "created", "updater", "updated", "body", "message", "pattern",
		"count", "host.name",
	}
	var cols []string
	seen := map[string]bool{}
	for _, p := range priority {
		if keys[p] {
			cols = append(cols, p)
			seen[p] = true
		}
	}
	var rest []string
	for k := range keys {
		if !seen[k] {
			rest = append(rest, k)
		}
	}
	sort.Strings(rest)
	// Cap implicit columns so wide objects stay readable; explicit --columns
	// bypasses this in printTabular.
	const maxImplicit = 12
	for _, k := range rest {
		if len(cols) >= maxImplicit {
			break
		}
		cols = append(cols, k)
	}
	return cols
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}
