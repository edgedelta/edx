// Package skills installs the Edge Delta agent skills (ed-logs, ed-metrics,
// ed-monitors, ...) into coding assistants such as Claude Code, Cursor, Codex
// and opencode.
//
// The skills are authored in the top-level skills directory and embedded in the
// binary, so they always match the edx version that ships them and
// installation needs no network access.
package skills

import (
	"io/fs"

	skillsrc "github.com/edgedelta/edx/skills"
)

// Embedded returns the embedded skills filesystem, where each top-level entry
// is a skill such as "ed-logs".
func Embedded() fs.FS {
	return skillsrc.FS
}
