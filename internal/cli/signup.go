package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/edgedelta/edx/internal/config"
)

// newSignupCmd creates an Edge Delta account and logs in. It is a thin,
// onboarding-flavored wrapper over the same OAuth transport as "auth login": it
// discovers the user's Git email as a starting point, then hands off to the
// browser (or a device code on headless machines) where the account is created
// and the login completes. The only long-lived credential stored is the API
// token the login issues.
func newSignupCmd() *cobra.Command {
	var env string
	var setDefault, force bool
	cmd := &cobra.Command{
		Use:   "signup",
		Short: "Create an Edge Delta account and log in",
		Long: `Create a new Edge Delta account and log in.

signup finds your Git email as a starting point, confirms it, then emails you a
one-click sign-in link. Open the link, confirm the code shown in your terminal,
and you're signed in. The login is saved as a profile, exactly like
"edx auth login" — existing accounts are simply logged in.`,
		Example: `  edx signup                                  # email a sign-in link, then confirm the code
  edx signup --profile staging --env staging  # sign up against staging under a named profile`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if env == "" {
				env = flagEnv
			}
			if env == "" {
				env = config.DefaultEnv
			}
			eps, ok := config.EndpointsForEnv(env)
			if !ok {
				return fmt.Errorf("unknown environment %q (valid: %s)", env, strings.Join(config.KnownEnvs(), ", "))
			}
			name := flagProfile
			if name == "" {
				name = "default"
			}

			cfg0, err := config.Load()
			if err != nil {
				return err
			}
			if loginWouldClobber(cfg0, name, cmd.Flags().Changed("profile"), force) {
				return fmt.Errorf("profile %q already exists; pass --profile <name> to save under a different name, or --force to overwrite it", name)
			}

			email, err := confirmSignupEmail(discoverGitEmail())
			if err != nil {
				return err
			}
			if email == "" {
				return fmt.Errorf("an email address is required to sign up")
			}
			fmt.Fprintf(os.Stderr, "Setting up Edge Delta for %s.\n", email)

			// The server emails a magic link; clicking it creates the account (if
			// new) and, after the code confirmation, issues the token signup stores.
			return loginWithOAuth(cmd, eps, env, name, "", setDefault, false, email)
		},
	}
	cmd.Flags().StringVar(&env, "env", "", "environment for this profile: prod, staging or local (default prod)")
	cmd.Flags().BoolVar(&setDefault, "set-default", false, "make this profile the default")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite the target profile if it already exists")
	return cmd
}

// discoverGitEmail returns the user's configured Git email, or "" if it is unset
// or a GitHub no-reply address (which can't receive mail, so it's a poor default).
func discoverGitEmail() string {
	out, err := exec.Command("git", "config", "--get", "user.email").Output()
	if err != nil {
		return ""
	}
	email := strings.TrimSpace(string(out))
	if email == "" || strings.HasSuffix(strings.ToLower(email), "noreply.github.com") {
		return ""
	}
	return email
}

// confirmSignupEmail shows the discovered email and lets the user accept it
// (Enter) or type a different one. With no discovered email it prompts for one.
// A non-interactive stdin, or an empty answer, keeps the discovered value (which
// may be ""); the browser step asks for the address either way, so this is only
// a convenience.
func confirmSignupEmail(discovered string) (string, error) {
	if !fileIsTTY(os.Stdin) {
		return discovered, nil
	}
	if discovered != "" {
		fmt.Fprintf(os.Stderr, "Found your Git email:\n\n    %s\n\nPress Enter to use it, or type another email: ", discovered)
	} else {
		fmt.Fprint(os.Stderr, "Enter your email (or press Enter to choose it in the browser): ")
	}
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	line = strings.TrimSpace(line)
	if line != "" {
		return line, nil
	}
	if err != nil {
		// EOF with nothing typed — fall back to the discovered value.
		return discovered, nil
	}
	return discovered, nil
}
