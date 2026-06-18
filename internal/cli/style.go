package cli

import (
	"os"
	"strings"
)

// colorEnabled is true when stderr is an interactive terminal and the user
// hasn't opted out via NO_COLOR. Computed once at startup.
var colorEnabled = func() bool {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}
	fi, err := os.Stderr.Stat()
	return err == nil && fi.Mode()&os.ModeCharDevice != 0
}()

func ansi(code, s string) string {
	if !colorEnabled {
		return s
	}
	return "\x1b[" + code + "m" + s + "\x1b[0m"
}

func dim(s string) string   { return ansi("2", s) }
func green(s string) string { return ansi("32", s) }
func okMark() string        { return green("✓") }

// shortID abbreviates a UUID/ULID for display, e.g. "0481a213…".
func shortID(id string) string {
	if len(id) <= 10 {
		return id
	}
	return id[:8] + "…"
}

// hostOnly returns the bare host of a URL for display (drops the scheme).
func hostOnly(rawURL string) string {
	s := strings.TrimPrefix(strings.TrimPrefix(rawURL, "https://"), "http://")
	return strings.TrimSuffix(s, "/")
}
