package cli

import (
	"fmt"
	"io/fs"
	"regexp"
	"strings"
	"testing"

	skillsrc "github.com/edgedelta/edx/skills"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Agents run the skills' edx commands verbatim, so a renamed command or a
// removed flag in a skill is a broken instruction. Every `edx ...` invocation in
// a skill's code blocks and inline code must resolve against the real command
// tree: each subcommand must exist and each flag must be defined on it.
func TestSkillDocsReferenceRealCommands(t *testing.T) {
	root := NewRootCmd()
	checked := 0
	err := fs.WalkDir(skillsrc.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}
		raw, err := fs.ReadFile(skillsrc.FS, path)
		if err != nil {
			return err
		}
		for _, inv := range edxInvocations(string(raw)) {
			checked++
			for _, problem := range checkInvocation(root, inv.words) {
				t.Errorf("%s:%d: `%s`: %s", path, inv.line, inv.text, problem)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	// Guards the extractor itself: if it silently matched nothing, the test
	// would pass while checking nothing.
	if checked < 100 {
		t.Fatalf("only %d edx invocations found in the skills; the extractor is likely broken", checked)
	}
}

type edxInvocation struct {
	line  int
	text  string
	words []string // the words after "edx"
}

var (
	inlineCode = regexp.MustCompile("`([^`]+)`")
	// An invocation starts a line (optionally after a "$ " prompt) or follows
	// a shell separator. Prose that mentions edx, such as a quoted warning in
	// sample output or brew's edgedelta/tap/edx, is not an invocation.
	edxStart = regexp.MustCompile(`(?:^\s*(?:\$\s+)?|\|\s*|&&\s*|;\s*|\$\(\s*)edx\s`)
)

// edxInvocations extracts edx command lines from a Markdown document's fenced
// code blocks (joining backslash continuations) and inline code spans.
func edxInvocations(doc string) []edxInvocation {
	var out []edxInvocation
	add := func(line int, code string) {
		for _, loc := range edxStart.FindAllStringIndex(code, -1) {
			text := code[loc[1]-len("edx")-1:]
			words := shellWords(strings.TrimPrefix(text, "edx"))
			out = append(out, edxInvocation{line: line, text: text, words: words})
		}
	}
	inFence := false
	pending, pendingLine := "", 0
	for i, line := range strings.Split(doc, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inFence = !inFence
			continue
		}
		if inFence {
			if pending == "" {
				pendingLine = i + 1
			}
			if strings.HasSuffix(trimmed, `\`) {
				pending += strings.TrimSuffix(trimmed, `\`) + " "
				continue
			}
			add(pendingLine, pending+trimmed)
			pending = ""
			continue
		}
		for _, m := range inlineCode.FindAllStringSubmatch(line, -1) {
			add(i+1, m[1])
		}
	}
	return out
}

// shellWords splits a command line into words, honoring single and double
// quotes, and stops at the first pipe, redirect, separator or comment.
func shellWords(s string) []string {
	var words []string
	var cur strings.Builder
	inWord := false
	var quote rune
	flush := func() {
		if inWord {
			words = append(words, cur.String())
			cur.Reset()
			inWord = false
		}
	}
	for _, r := range s {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				cur.WriteRune(r)
			}
		case r == '\'' || r == '"':
			quote, inWord = r, true
		case r == ' ' || r == '\t':
			flush()
		case !inWord && strings.ContainsRune("|;&>#", r):
			flush()
			return words
		default:
			cur.WriteRune(r)
			inWord = true
		}
	}
	flush()
	return words
}

// checkInvocation resolves words against the command tree and returns what does
// not exist. Documentation shorthand is understood: "a/b" names alternative
// subcommands (each must exist), and "..." or a <placeholder> ends the check.
func checkInvocation(root *cobra.Command, words []string) []string {
	var problems []string
	var walk func(cmd *cobra.Command, path string, words []string)
	walk = func(cmd *cobra.Command, path string, words []string) {
		skipValue := false
		for i, w := range words {
			if skipValue {
				skipValue = false
				continue
			}
			if w == "..." || strings.HasPrefix(w, "<") {
				return
			}
			if strings.HasPrefix(w, "-") {
				f, p := checkFlag(cmd, w)
				if p != "" {
					problems = append(problems, fmt.Sprintf("%s: %s", path, p))
				}
				// "--flag value": the next word is the value, not a
				// subcommand or argument.
				skipValue = f != nil && f.Value.Type() != "bool" && !strings.Contains(w, "=")
				continue
			}
			if !cmd.HasSubCommands() {
				// A positional argument (or a flag's value): the rest are
				// flags and arguments of this command.
				continue
			}
			alts := strings.Split(w, "/")
			var found []*cobra.Command
			for _, a := range alts {
				sub := findSub(cmd, a)
				if sub == nil {
					if cmd.Runnable() && len(alts) == 1 {
						break // an argument to a runnable parent
					}
					problems = append(problems, fmt.Sprintf("%s has no subcommand %q", path, a))
					continue
				}
				found = append(found, sub)
			}
			if len(found) == 0 {
				continue
			}
			for _, sub := range found {
				walk(sub, path+" "+sub.Name(), words[i+1:])
			}
			return
		}
	}
	walk(root, "edx", words)
	return problems
}

func findSub(cmd *cobra.Command, name string) *cobra.Command {
	for _, sub := range cmd.Commands() {
		if sub.Name() == name || sub.HasAlias(name) {
			return sub
		}
	}
	return nil
}

// checkFlag looks w up among the flags cmd accepts, including persistent flags
// defined on it or its parents. It returns the flag when found, or a problem.
func checkFlag(cmd *cobra.Command, w string) (*pflag.Flag, string) {
	lookup := func(name string) *pflag.Flag {
		for _, fs := range []*pflag.FlagSet{cmd.Flags(), cmd.PersistentFlags(), cmd.InheritedFlags()} {
			if f := fs.Lookup(name); f != nil {
				return f
			}
		}
		return nil
	}
	if strings.HasPrefix(w, "--") {
		name, _, _ := strings.Cut(strings.TrimPrefix(w, "--"), "=")
		if name == "" || name == "help" || name == "version" {
			return nil, ""
		}
		if f := lookup(name); f != nil {
			return f, ""
		}
		return nil, fmt.Sprintf("unknown flag --%s", name)
	}
	short := strings.TrimPrefix(w, "-")
	if len(short) != 1 || short == "h" {
		return nil, "" // a value like -severity_text:"DEBUG", not a flag
	}
	for _, fs := range []*pflag.FlagSet{cmd.Flags(), cmd.PersistentFlags(), cmd.InheritedFlags()} {
		if f := fs.ShorthandLookup(short); f != nil {
			return f, ""
		}
	}
	return nil, fmt.Sprintf("unknown flag -%s", short)
}

func TestCheckInvocationCatchesDrift(t *testing.T) {
	root := NewRootCmd()
	for _, tc := range []struct {
		line    string
		wantBad bool
	}{
		{`edx logs search --scope audit -q 'method:"DELETE"' --lookback 24h`, false},
		{`edx --profile staging logs search -q error`, false},
		{`edx monitors list/get/delete`, false},
		{`edx api GET /tokens`, false},
		{`edx logs serch -q error`, true},                            // misspelled subcommand
		{`edx metrics query --name x -q 'service.name:"api"'`, true}, // flag that does not exist
		{`edx logs search --max-itemz 5`, true},                      // misspelled flag
		{`edx monitors list/frob`, true},                             // one bad alternative
	} {
		invs := edxInvocations("`" + tc.line + "`")
		if len(invs) != 1 {
			t.Fatalf("%q: extracted %d invocations, want 1", tc.line, len(invs))
		}
		problems := checkInvocation(root, invs[0].words)
		if bad := len(problems) > 0; bad != tc.wantBad {
			t.Errorf("%q: problems = %v, want bad=%v", tc.line, problems, tc.wantBad)
		}
	}
	// Sample output in a fenced block mentions edx without invoking it.
	if n := len(edxInvocations("```\n! got 3, edx can only validate version 4\n```")); n != 0 {
		t.Errorf("prose mentioning edx extracted %d invocations, want 0", n)
	}
}
