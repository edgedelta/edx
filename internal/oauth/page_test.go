package oauth

import (
	"strings"
	"testing"
)

func TestResultPage(t *testing.T) {
	ok := resultPage(true, "You're signed in", "Close this tab.")
	for _, want := range []string{"You're signed in", "Close this tab.", "Edge Delta CLI", "<svg", "#16a34a", "Close tab", "window.close()"} {
		if !strings.Contains(ok, want) {
			t.Errorf("success page missing %q", want)
		}
	}
	if !strings.Contains(resultPage(false, "Sign-in failed", "x"), "#dc2626") {
		t.Error("error page should use the error accent color")
	}
}
