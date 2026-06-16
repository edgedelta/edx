package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintJSON(t *testing.T) {
	var buf bytes.Buffer
	err := Print(&buf, []byte(`{"b":1,"a":"x"}`), Options{Format: "json"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"a": "x"`) {
		t.Errorf("expected pretty JSON, got %q", buf.String())
	}
}

func TestPrintTableFromItemsEnvelope(t *testing.T) {
	data := []byte(`{"items":[{"id":"1","name":"a","nested":{"k":"v"}},{"id":"2","name":"b"}],"next_cursor":"x"}`)
	var buf bytes.Buffer
	if err := Print(&buf, data, Options{Format: "table"}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{"ID", "NAME", "1", "2", "NESTED.K", "v"} {
		if !strings.Contains(out, want) {
			t.Errorf("table output missing %q:\n%s", want, out)
		}
	}
}

func TestPrintTableRootArray(t *testing.T) {
	data := []byte(`[{"tag":"prod","status":"running"}]`)
	var buf bytes.Buffer
	if err := Print(&buf, data, Options{Format: "table", Columns: []string{"tag", "status"}}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "prod") || !strings.Contains(out, "running") {
		t.Errorf("unexpected table output:\n%s", out)
	}
}

func TestPrintCSV(t *testing.T) {
	data := []byte(`[{"a":"1","b":"2"}]`)
	var buf bytes.Buffer
	if err := Print(&buf, data, Options{Format: "csv", Columns: []string{"a", "b"}}); err != nil {
		t.Fatal(err)
	}
	if got := buf.String(); got != "a,b\n1,2\n" {
		t.Errorf("unexpected csv output: %q", got)
	}
}

func TestPrintSingleObjectTable(t *testing.T) {
	data := []byte(`{"id":"abc","status":"ok"}`)
	var buf bytes.Buffer
	if err := Print(&buf, data, Options{Format: "table"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "abc") {
		t.Errorf("unexpected output:\n%s", buf.String())
	}
}

func TestPrintTableNestedCollectionEnvelope(t *testing.T) {
	// Chat service shape: the collection sits under data.issues, not data.
	data := []byte(`{"status":200,"data":{"issues":[{"issueId":"i1","severity":"high"},{"issueId":"i2","severity":"low"}]}}`)
	var buf bytes.Buffer
	if err := Print(&buf, data, Options{Format: "table", Columns: []string{"issueId", "severity"}}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{"i1", "high", "i2", "low"} {
		if !strings.Contains(out, want) {
			t.Errorf("nested-collection table missing %q:\n%s", want, out)
		}
	}
}

func TestPrintTableSingleResourceEnvelope(t *testing.T) {
	// Single-resource GET shape: the object sits under data.
	data := []byte(`{"status":200,"data":{"issueId":"i1","severity":"high","state":"open"}}`)
	var buf bytes.Buffer
	if err := Print(&buf, data, Options{Format: "table", Columns: []string{"issueId", "state"}}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "i1") || !strings.Contains(out, "open") {
		t.Errorf("single-resource envelope should unwrap data:\n%s", out)
	}
}

func TestPrintTableResourceWithEmbeddedCollection(t *testing.T) {
	// A single resource (thread) that embeds a sub-collection (messages) must
	// render as one row of the resource, NOT descend into the embed.
	data := []byte(`{"status":200,"data":{"id":"t1","state":"open","messages":[{"id":"m1"},{"id":"m2"}]}}`)
	var buf bytes.Buffer
	if err := Print(&buf, data, Options{Format: "table", Columns: []string{"id", "state"}}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "t1") || !strings.Contains(out, "open") {
		t.Errorf("should render the thread resource, got:\n%s", out)
	}
	if strings.Contains(out, "m1") || strings.Contains(out, "m2") {
		t.Errorf("should NOT descend into embedded messages, got:\n%s", out)
	}
}

func TestPrintUnknownFormat(t *testing.T) {
	if err := Print(&bytes.Buffer{}, []byte(`{}`), Options{Format: "bogus"}); err == nil {
		t.Error("expected error for unknown format")
	}
}
