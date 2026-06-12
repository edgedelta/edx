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

func TestPrintUnknownFormat(t *testing.T) {
	if err := Print(&bytes.Buffer{}, []byte(`{}`), Options{Format: "bogus"}); err == nil {
		t.Error("expected error for unknown format")
	}
}
