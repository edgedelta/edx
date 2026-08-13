package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// runCQLValidate runs `cql validate` with the given args and stdin, returning stdout,
// stderr and the error the command exited with.
func runCQLValidate(t *testing.T, stdin string, args ...string) (string, string, error) {
	t.Helper()
	cmd := newCQLCmd()
	var out, errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	cmd.SetIn(strings.NewReader(stdin))
	cmd.SetArgs(append([]string{"validate"}, args...))

	// warnf writes to the process's stderr rather than the command's, so capture it too.
	restore := captureStderr(t, &errOut)
	err := cmd.Execute()
	restore()
	return out.String(), errOut.String(), err
}

// captureStderr redirects os.Stderr into buf until the returned function is called.
func captureStderr(t *testing.T, buf *bytes.Buffer) func() {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	saved := os.Stderr
	os.Stderr = w

	// The drainer must not write into buf itself: the command under test
	// writes to the same buffer, and pipe EOF is not a happens-before edge
	// the race detector recognizes. Hand the bytes back over a channel and
	// let the caller's goroutine do the write.
	captured := make(chan []byte, 1)
	go func() {
		var chunk bytes.Buffer
		_, _ = chunk.ReadFrom(r)
		captured <- chunk.Bytes()
	}()

	return func() {
		os.Stderr = saved
		_ = w.Close()
		buf.Write(<-captured)
		_ = r.Close()
	}
}

func TestCQLValidateAcceptsValidQueries(t *testing.T) {
	cases := []struct{ dataType, query string }{
		{"metric", "sum:ed.host.cpu{*} by {host.name}.rollup(60)"},
		{"log", "{error AND ed.tag:$fleet} by {host.name}"},
		{"event", "{severity_text:ALERT event.type:log_threshold}"},
		{"pattern", "ed.tag:$fleet"},
		{"trace", "{service.name:checkout}"},
		{"formula", "timeshift(q1, 3600) / q2"},
	}

	for _, c := range cases {
		t.Run(c.dataType, func(t *testing.T) {
			stdout, _, err := runCQLValidate(t, "", "--type", c.dataType, c.query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(stdout, "valid") {
				t.Errorf("stdout = %q, want it to confirm validity", stdout)
			}
		})
	}
}

func TestCQLValidateReportsSyntaxErrorsWithACaret(t *testing.T) {
	stdout, stderr, err := runCQLValidate(t, "", "--type", "metric", "sum:foo{*}.rollup(abc)")
	if err == nil {
		t.Fatalf("expected a non-nil error for an invalid query; stdout=%q", stdout)
	}
	if !strings.Contains(err.Error(), "1 of 1 query failed") {
		t.Errorf("error = %q, want it to count the failures", err)
	}
	if !strings.Contains(stderr, "column 18") {
		t.Errorf("stderr = %q, want the column", stderr)
	}
	// The caret sits under `abc`, which starts at column 18.
	if !strings.Contains(stderr, strings.Repeat(" ", 18)+"^") {
		t.Errorf("stderr = %q, want a caret aligned to column 18", stderr)
	}
}

// The grammars are not interchangeable, so --type has to actually select one.
func TestCQLValidateEnforcesTheDataTypesGrammar(t *testing.T) {
	if _, _, err := runCQLValidate(t, "", "--type", "log", "sum:foo{*}"); err == nil {
		t.Error("log grammar accepted a metric query")
	}
	if _, _, err := runCQLValidate(t, "", "--type", "metric", "{just.a:filter}"); err == nil {
		t.Error("metric grammar accepted a log query")
	}
}

func TestCQLValidateRejectsUnknownType(t *testing.T) {
	_, _, err := runCQLValidate(t, "", "--type", "sql", "SELECT 1")
	if err == nil {
		t.Fatal("expected an error for an unknown --type")
	}
	if !strings.Contains(err.Error(), "sql") {
		t.Errorf("error = %q, want it to name the rejected type", err)
	}
	// The message should tell the user what is accepted instead.
	for _, want := range []string{"log", "metric", "formula"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error = %q, want it to list %q", err, want)
		}
	}
}

func TestCQLValidateReadsQueriesFromAFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "queries.txt")
	content := strings.Join([]string{
		"# comment lines and blanks are skipped",
		"",
		"sum:a{*}.rollup(60)",
		"avg:b{*} by {host.name}",
		"   ",
		"count_unique(x):c{*}",
	}, "\n")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	stdout, _, err := runCQLValidate(t, "", "--type", "metric", "--file", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stdout, "3 queries are valid") {
		t.Errorf("stdout = %q, want a count of 3 (comments and blanks skipped)", stdout)
	}
}

func TestCQLValidateReadsQueriesFromStdin(t *testing.T) {
	stdin := "sum:a{*}\nsum:b{*}.rollup(nope)\n"
	_, stderr, err := runCQLValidate(t, stdin, "--type", "metric", "--file", "-")
	if err == nil {
		t.Fatal("expected an error: one of the two queries is invalid")
	}
	if !strings.Contains(err.Error(), "1 of 2 queries failed") {
		t.Errorf("error = %q, want it to report 1 of 2", err)
	}
	if !strings.Contains(stderr, "nope") {
		t.Errorf("stderr = %q, want it to show the offending query", stderr)
	}
}

// Arguments and --file are additive, so a caller can pin one extra query onto a list.
func TestCQLValidateCombinesArgumentsAndFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "queries.txt")
	if err := os.WriteFile(path, []byte("sum:a{*}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	stdout, _, err := runCQLValidate(t, "", "--type", "metric", "--file", path, "avg:b{*}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stdout, "2 queries are valid") {
		t.Errorf("stdout = %q, want both queries counted", stdout)
	}
}

func TestCQLValidateQuietPrintsNothingOnSuccess(t *testing.T) {
	stdout, _, err := runCQLValidate(t, "", "--type", "metric", "--quiet", "sum:a{*}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stdout != "" {
		t.Errorf("stdout = %q, want it empty under --quiet", stdout)
	}
}

func TestCQLValidateRequiresAQuery(t *testing.T) {
	_, _, err := runCQLValidate(t, "", "--type", "metric")
	if err == nil {
		t.Fatal("expected an error when no queries are given")
	}
	if !strings.Contains(err.Error(), "no queries") {
		t.Errorf("error = %q, want it to say there is nothing to validate", err)
	}
}

// An empty query is valid everywhere else (a missing filter means match-all), so passing
// only whitespace has to read as "nothing given" rather than silently passing.
func TestCQLValidateTreatsBlankArgumentsAsNoQuery(t *testing.T) {
	_, _, err := runCQLValidate(t, "", "--type", "metric", "   ")
	if err == nil || !strings.Contains(err.Error(), "no queries") {
		t.Errorf("error = %v, want it to report no queries", err)
	}
}
