package skills

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func testFS() fstest.MapFS {
	return fstest.MapFS{
		"ed-logs/SKILL.md":         {Data: []byte("---\nname: ed-logs\ndescription: Search logs with CQL.\nmetadata:\n  description: nested should be ignored\n---\n\n# Logs\n")},
		"ed-metrics/SKILL.md":      {Data: []byte("---\nname: ed-metrics\ndescription: \"Aggregate metrics.\"\n---\nbody\n")},
		"ed-logs/reference/cql.md": {Data: []byte("cql reference")},
		"assets/logo.png":          {Data: []byte("not a skill")}, // no SKILL.md -> skipped
	}
}

func TestListReturnsSkillsSortedWithDescriptions(t *testing.T) {
	got, err := List(testFS())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 skills (assets dir has no SKILL.md), got %d: %+v", len(got), got)
	}
	if got[0].Name != "ed-logs" || got[1].Name != "ed-metrics" {
		t.Errorf("not sorted by name: %+v", got)
	}
	if got[0].Description != "Search logs with CQL." {
		t.Errorf("description parse: got %q", got[0].Description)
	}
	// Quoted description must be unquoted.
	if got[1].Description != "Aggregate metrics." {
		t.Errorf("quoted description: got %q", got[1].Description)
	}
}

func TestReadUnknownSkillErrors(t *testing.T) {
	if _, err := Read(testFS(), "ed-nope"); err == nil {
		t.Fatal("expected error for unknown skill")
	}
	if _, err := Read(testFS(), "ed-logs"); err != nil {
		t.Fatalf("Read known skill: %v", err)
	}
}

func TestInstallCopiesWholeTree(t *testing.T) {
	dest := t.TempDir()
	n, err := Install(testFS(), "ed-logs", dest)
	if err != nil {
		t.Fatalf("Install: %v", err)
	}
	if n != 2 { // SKILL.md + reference/cql.md
		t.Errorf("want 2 files written, got %d", n)
	}
	if b, err := os.ReadFile(filepath.Join(dest, "ed-logs", "SKILL.md")); err != nil || len(b) == 0 {
		t.Errorf("SKILL.md not installed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "ed-logs", "reference", "cql.md")); err != nil {
		t.Errorf("reference file not installed: %v", err)
	}
}

func TestInstallUnknownSkillErrors(t *testing.T) {
	if _, err := Install(testFS(), "ed-nope", t.TempDir()); err == nil {
		t.Fatal("expected error installing unknown skill")
	}
}

func TestInstallRejectsPathTraversal(t *testing.T) {
	if _, err := Install(testFS(), "../etc", t.TempDir()); err == nil {
		t.Fatal("expected error for traversal name")
	}
}

func TestEmbeddedSkillsArePresent(t *testing.T) {
	// Guards the vendored copy: `make sync-skills` must have run and the core
	// skills must be embedded.
	names, err := Names(Embedded())
	if err != nil {
		t.Fatalf("Names(Embedded): %v", err)
	}
	want := map[string]bool{"ed-edx": false, "ed-logs": false, "ed-monitors": false}
	for _, n := range names {
		if _, ok := want[n]; ok {
			want[n] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("embedded skills missing %q (run `make sync-skills`); have %v", name, names)
		}
	}
}
