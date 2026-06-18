package cli

import (
	"strings"
	"testing"
)

func TestRenderBanner(t *testing.T) {
	plain := renderBanner(false)
	if strings.Count(plain, "\n") != 4 { // 5 lines
		t.Errorf("expected 5-line banner, got:\n%s", plain)
	}
	if strings.Contains(plain, "\x1b[") {
		t.Error("plain banner must not contain ANSI escapes")
	}
	colored := renderBanner(true)
	if !strings.HasPrefix(colored, "\x1b[38;2;10;120;230m") || !strings.HasSuffix(colored, "\x1b[0m") {
		t.Error("colored banner should be wrapped in the brand-blue ANSI escape")
	}
}
