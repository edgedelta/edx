package cli

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/edgedelta/edx/internal/config"
	"github.com/edgedelta/edx/internal/oauth"
)

func newAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage Edge Delta credentials",
		Long: `Manage Edge Delta API credentials stored in ~/.config/edx/config.yaml.

Each profile targets an environment (prod, staging or local), which selects
the main API host and the AI Teammate service hosts together.

Authenticate either with a static API token (--token) or interactively via
OAuth in your browser (--oauth). OAuth tokens are refreshed automatically.

Credentials can also be supplied via environment variables, which take
precedence over the config file:
  ED_API_TOKEN   API token (created under Admin > API Tokens)
  ED_ORG_ID      organization ID
  ED_ENV         environment: prod, staging or local`,
	}
	cmd.AddCommand(newAuthLoginCmd(), newAuthStatusCmd(), newAuthLogoutCmd())
	return cmd
}

func newAuthLoginCmd() *cobra.Command {
	var token, orgID, env string
	var useOAuth, setDefault bool
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Save credentials to a profile (API token or OAuth)",
		Example: `  edx auth login --token 00000000-0000-0000-0000-000000000000 --org-id <org-id>
  edx auth login --oauth --org-id <org-id>
  edx auth login --profile staging --env staging --oauth --org-id ...`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if orgID == "" {
				return fmt.Errorf("--org-id is required")
			}
			// --env on this command takes precedence over the persistent --env;
			// fall back to the persistent flag so `--env staging login` works too.
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

			if useOAuth {
				if token != "" {
					return fmt.Errorf("--oauth and --token are mutually exclusive")
				}
				fmt.Fprintf(os.Stderr, "Opening your browser to log in to %s …\n", eps.API)
				toks, err := oauth.Login(cmd.Context(), eps.API, oauth.LoginOptions{
					OpenBrowser: true,
					Prompt: func(u string) {
						fmt.Fprintf(os.Stderr, "If your browser did not open, visit:\n  %s\n", u)
					},
				})
				if err != nil {
					return fmt.Errorf("oauth login failed: %w", err)
				}
				if err := config.SaveOAuthTokens(name, env, orgID, toks.ClientID, toks.AccessToken, toks.RefreshToken, toks.Expiry); err != nil {
					return err
				}
				if setDefault {
					if err := setDefaultProfile(name); err != nil {
						return err
					}
				}
				path, _ := config.Path()
				fmt.Fprintf(os.Stderr, "Logged in via OAuth — saved profile %q (env: %s) to %s\n", name, env, path)
				return nil
			}

			if token == "" {
				return fmt.Errorf("provide --token, or use --oauth to log in via browser")
			}
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			cfg.Profiles[name] = &config.Profile{Env: env, OrgID: orgID, AuthMethod: config.AuthMethodToken, APIToken: token}
			if cfg.DefaultProfile == "" || setDefault {
				cfg.DefaultProfile = name
			}
			if err := cfg.Save(); err != nil {
				return err
			}
			path, _ := config.Path()
			fmt.Fprintf(os.Stderr, "Saved profile %q (env: %s) to %s\n", name, env, path)
			return nil
		},
	}
	cmd.Flags().StringVar(&token, "token", "", "Edge Delta API token (token auth)")
	cmd.Flags().BoolVar(&useOAuth, "oauth", false, "log in via OAuth in your browser instead of a token")
	cmd.Flags().StringVar(&orgID, "org-id", "", "Edge Delta organization ID (required)")
	cmd.Flags().StringVar(&env, "env", "", "environment for this profile: prod, staging or local (default prod)")
	cmd.Flags().BoolVar(&setDefault, "set-default", false, "make this profile the default")
	return cmd
}

// setDefaultProfile marks name as the config default.
func setDefaultProfile(name string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	cfg.DefaultProfile = name
	return cfg.Save()
}

func newAuthStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the active profile and verify the token against the API",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := config.Resolve(flagProfile, flagEnv, flagOrg, flagToken)
			if err != nil {
				return err
			}
			cred := maskToken(r.APIToken)
			if r.UsesOAuth() {
				cred = maskToken(r.OAuthAccessToken) + " (auto-refreshed)"
			}
			fmt.Fprintf(os.Stderr,
				"Profile:   %s\nEnv:       %s\nAuth:      %s\nAPI URL:   %s\nChat URL:  %s\nAgent URL: %s\nOrg ID:    %s\nCredential:%s\n",
				r.Profile, r.Env, r.AuthMethod, r.APIURL, r.ChatURL, r.AgentURL, r.OrgID, cred)

			c, err := newClient()
			if err != nil {
				return err
			}
			// Cheap authenticated call to verify token + org pairing.
			q := url.Values{}
			q.Set("scope", "log")
			if _, err := c.Get(cmdContext(cmd), "/facet_keys", q); err != nil {
				return fmt.Errorf("token verification failed: %w", err)
			}
			fmt.Fprintln(os.Stderr, "Status:  OK (token accepted)")
			return nil
		},
	}
}

func newAuthLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Remove the stored credentials for the active profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			name := flagProfile
			if name == "" {
				name = cfg.DefaultProfile
			}
			if name == "" {
				name = "default"
			}
			if _, ok := cfg.Profiles[name]; !ok {
				return fmt.Errorf("profile %q not found", name)
			}
			delete(cfg.Profiles, name)
			if cfg.DefaultProfile == name {
				cfg.DefaultProfile = ""
			}
			if err := cfg.Save(); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Removed profile %q\n", name)
			return nil
		},
	}
}

func maskToken(t string) string {
	if t == "" {
		return "(not set)"
	}
	if len(t) <= 8 {
		return "****"
	}
	return t[:4] + "..." + t[len(t)-4:]
}
