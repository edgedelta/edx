package cli

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
	pflag "github.com/spf13/pflag"
)

// A flag shorthand that collides with a persistent one panics when cobra merges the
// flagsets, which happens on every invocation — including `edx --help`. Nothing else here
// builds the whole tree, so a new command's flags are otherwise unexercised until someone
// runs the binary.
func TestCommandTreeIsWellFormed(t *testing.T) {
	root := NewRootCmd()

	persistent := map[string]string{}
	root.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		if f.Shorthand != "" {
			persistent[f.Shorthand] = f.Name
		}
	})

	var walk func(cmd *cobra.Command, path string)
	walk = func(cmd *cobra.Command, path string) {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Shorthand == "" {
				return
			}
			if owner, taken := persistent[f.Shorthand]; taken && owner != f.Name {
				t.Errorf("%s: -%s is --%s on the local flag but --%s globally; cobra panics on this",
					path, f.Shorthand, f.Name, owner)
			}
		})
		if cmd.Runnable() && cmd.Short == "" {
			t.Errorf("%s: no Short description, so it is invisible in help", path)
		}
		for _, sub := range cmd.Commands() {
			walk(sub, path+" "+sub.Name())
		}
	}
	walk(root, "edx")

	// The help text is what actually triggered the panic, so render it for every command.
	var render func(cmd *cobra.Command, path string)
	render = func(cmd *cobra.Command, path string) {
		var buf strings.Builder
		cmd.SetOut(&buf)
		if err := cmd.Usage(); err != nil {
			t.Errorf("%s: rendering usage failed: %v", path, err)
		}
		for _, sub := range cmd.Commands() {
			render(sub, path+" "+sub.Name())
		}
	}
	render(root, "edx")
}
