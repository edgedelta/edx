package cli

import (
	"fmt"
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

// edgeDeltaBanner is the "EDGE DELTA" wordmark (figlet, Standard font).
const edgeDeltaBanner = ` _____ ____   ____ _____   ____  _____ _   _____  _
| ____|  _ \ / ___| ____| |  _ \| ____| | |_   _|/ \
|  _| | | | | |  _|  _|   | | | |  _| | |   | | / _ \
| |___| |_| | |_| | |___  | |_| | |___| |___| |/ ___ \
|_____|____/ \____|_____| |____/|_____|_____|_/_/   \_\`

const brandGreen = "\x1b[38;2;0;218;99m" // #00DA63 (Edge Delta logo green)

// renderBanner returns the wordmark, optionally painted Edge Delta green.
func renderBanner(color bool) string {
	if color {
		return brandGreen + edgeDeltaBanner + "\x1b[0m"
	}
	return edgeDeltaBanner
}

// fileIsTTY reports whether f is an interactive terminal.
func fileIsTTY(f *os.File) bool {
	fi, err := f.Stat()
	return err == nil && fi.Mode()&os.ModeCharDevice != 0
}

// writeBanner prints the wordmark to f when f is an interactive terminal and
// the user hasn't opted out via EDX_NO_BANNER. It stays silent for pipes and
// redirects so scripted output is never polluted.
func writeBanner(f *os.File) {
	if os.Getenv("EDX_NO_BANNER") != "" || !fileIsTTY(f) {
		return
	}
	color := os.Getenv("NO_COLOR") == "" && os.Getenv("TERM") != "dumb"
	fmt.Fprintln(f, renderBanner(color))
}
