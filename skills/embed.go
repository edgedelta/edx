// Package skills holds the Edge Delta agent skills (ed-logs, ed-metrics,
// ed-monitors, ...). This directory is their source of truth: each edx release
// mirrors it, minus this file, to github.com/edgedelta/agent-skills, where
// `npx skills add edgedelta/agent-skills` installs from.
package skills

import "embed"

// FS holds the skill directories. `all:` ensures dotfiles and
// underscore-prefixed reference files inside a skill are embedded too.
//
//go:embed all:ed-*
var FS embed.FS
