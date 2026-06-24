// Package skills installs the Edge Delta agent skills (ed-logs, ed-metrics,
// ed-monitors, ...) into coding assistants such as Claude Code, Cursor, Codex
// and opencode.
//
// The skills are authored in the github.com/edgedelta/agent-skills repository
// and vendored into ./data at build time (see `make sync-skills`). They are
// embedded in the binary so they always match the edx version that ships them
// and so installation needs no network access.
package skills

import (
	"embed"
	"io/fs"
)

// data holds the vendored skill directories. `all:` ensures dotfiles and
// underscore-prefixed reference files inside a skill are embedded too.
//
//go:embed all:data
var data embed.FS

// Embedded returns the embedded skills filesystem rooted at the skills
// directory (so each top-level entry is a skill such as "ed-logs").
func Embedded() fs.FS {
	sub, err := fs.Sub(data, "data")
	if err != nil {
		// Unreachable: "data" is embedded at compile time.
		panic(err)
	}
	return sub
}
