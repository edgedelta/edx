package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags "-X github.com/edgedelta/edx/internal/cli.Version=...".
var Version = "dev"

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the edx version",
		Run: func(cmd *cobra.Command, args []string) {
			writeBanner(os.Stdout)
			fmt.Printf("edx %s\n", Version)
		},
	}
}
