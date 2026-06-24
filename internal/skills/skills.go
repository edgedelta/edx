package skills

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Skill is one installable skill: a directory containing a SKILL.md.
type Skill struct {
	Name        string // directory name, e.g. "ed-logs"
	Description string // from the SKILL.md frontmatter `description:` field
}

// List returns every skill in fsys, sorted by name. A skill is any top-level
// directory that contains a SKILL.md file.
func List(fsys fs.FS) ([]Skill, error) {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, err
	}
	var skills []Skill
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		b, err := fs.ReadFile(fsys, path(e.Name(), "SKILL.md"))
		if err != nil {
			continue // directory without a SKILL.md is not a skill
		}
		skills = append(skills, Skill{
			Name:        e.Name(),
			Description: frontmatterDescription(b),
		})
	}
	sort.Slice(skills, func(i, j int) bool { return skills[i].Name < skills[j].Name })
	return skills, nil
}

// Names returns the sorted skill names in fsys.
func Names(fsys fs.FS) ([]string, error) {
	skills, err := List(fsys)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(skills))
	for i, s := range skills {
		names[i] = s.Name
	}
	return names, nil
}

// Read returns the SKILL.md contents of a single skill.
func Read(fsys fs.FS, name string) ([]byte, error) {
	if !exists(fsys, name) {
		return nil, fmt.Errorf("unknown skill %q", name)
	}
	return fs.ReadFile(fsys, path(name, "SKILL.md"))
}

// Install copies the skill's entire directory tree from fsys into
// destRoot/<name>/, creating directories as needed and overwriting existing
// files. It returns the number of files written.
func Install(fsys fs.FS, name, destRoot string) (int, error) {
	if !exists(fsys, name) {
		return 0, fmt.Errorf("unknown skill %q", name)
	}
	written := 0
	err := fs.WalkDir(fsys, name, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		dest := filepath.Join(destRoot, filepath.FromSlash(p))
		if d.IsDir() {
			return os.MkdirAll(dest, 0o755)
		}
		b, err := fs.ReadFile(fsys, p)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dest, b, 0o644); err != nil {
			return err
		}
		written++
		return nil
	})
	return written, err
}

// exists reports whether name is a skill (has a SKILL.md) in fsys.
func exists(fsys fs.FS, name string) bool {
	if name == "" || strings.ContainsRune(name, '/') {
		return false
	}
	_, err := fs.Stat(fsys, path(name, "SKILL.md"))
	return err == nil
}

// frontmatterDescription extracts the `description:` value from a SKILL.md's
// leading YAML frontmatter block (between the first two `---` lines). It
// returns an empty string when absent.
func frontmatterDescription(b []byte) string {
	sc := bufio.NewScanner(bytes.NewReader(b))
	inFrontmatter := false
	for sc.Scan() {
		line := sc.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			break // end of frontmatter
		}
		if !inFrontmatter {
			// No frontmatter started on the first non-empty line: give up.
			if trimmed != "" {
				return ""
			}
			continue
		}
		// Only top-level keys (no indentation) so we skip nested metadata.
		if line != strings.TrimLeft(line, " \t") {
			continue
		}
		if rest, ok := strings.CutPrefix(trimmed, "description:"); ok {
			return strings.Trim(strings.TrimSpace(rest), `"'`)
		}
	}
	return ""
}

// path joins slash-separated elements for use with io/fs (always forward slash).
func path(elem ...string) string {
	return strings.Join(elem, "/")
}
