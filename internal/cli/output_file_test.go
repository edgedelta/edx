package cli

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// captureStdout redirects os.Stdout for the duration of fn and returns what
// was written. printResult writes to os.Stdout directly, so cobra's SetOut is
// not enough to observe it.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = old }()
	fn()
	_ = w.Close()
	var sb strings.Builder
	buf := make([]byte, 64*1024)
	for {
		n, err := r.Read(buf)
		sb.Write(buf[:n])
		if err != nil {
			break
		}
	}
	return sb.String()
}

func TestOutputFileWritesFullResponse(t *testing.T) {
	// A response big enough that terminal / AI-harness stdout captures would
	// clip it; --output-file must land every byte in the file and keep stdout
	// clean.
	big := `{"matches":[` + strings.Repeat(`{"name":"`+strings.Repeat("x", 1024)+`"},`, 2047) +
		`{"name":"last"}]}`
	var got kgCapture
	srv := kgTestServer(t, &got, http.StatusOK, big)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	out := filepath.Join(t.TempDir(), "resp.json")
	stdout := captureStdout(t, func() {
		if err := runEdx(t, "ai", "knowledge", "search", "payment", "--output-file", out); err != nil {
			t.Errorf("search --output-file: %v", err)
		}
	})
	if stdout != "" {
		if len(stdout) > 200 {
			stdout = stdout[:200] + "..."
		}
		t.Errorf("stdout not empty with --output-file: %q", stdout)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}
	var v struct {
		Data struct {
			Matches []struct {
				Name string `json:"name"`
			} `json:"matches"`
		} `json:"data"`
	}
	// The whole point of --output-file is that the JSON is never syntactically
	// broken, so the file must parse as-is.
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatalf("output file is not valid JSON: %v", err)
	}
	if n := len(v.Data.Matches); n != 2048 {
		t.Errorf("matches in file = %d, want 2048", n)
	}
	if last := v.Data.Matches[len(v.Data.Matches)-1].Name; last != "last" {
		t.Errorf("last match = %q, want %q (file truncated?)", last, "last")
	}
}

func TestOutputFileHonorsFormat(t *testing.T) {
	var got kgCapture
	srv := kgTestServer(t, &got, http.StatusOK, `{"totalNodes":10}`)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	out := filepath.Join(t.TempDir(), "resp.yaml")
	if err := runEdx(t, "ai", "knowledge", "stats", "--output", "yaml", "--output-file", out); err != nil {
		t.Fatalf("stats --output yaml --output-file: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}
	if !strings.Contains(string(data), "totalNodes: 10") {
		t.Errorf("yaml output file missing rendered field, got: %q", string(data))
	}
}

func TestOutputFileUnwritablePathErrors(t *testing.T) {
	var got kgCapture
	srv := kgTestServer(t, &got, http.StatusOK, `{"totalNodes":10}`)
	defer srv.Close()
	useAgentEnv(t, srv.URL)

	out := filepath.Join(t.TempDir(), "no", "such", "dir", "resp.json")
	if err := runEdx(t, "ai", "knowledge", "stats", "--output-file", out); err == nil {
		t.Fatal("expected error for unwritable --output-file path")
	}
}
