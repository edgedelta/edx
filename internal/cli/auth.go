package cli

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"

	"github.com/edgedelta/edx/internal/config"
)

func newAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage Edge Delta credentials",
		Long: `Manage Edge Delta API credentials stored in ~/.config/edx/config.yaml.

Credentials can also be supplied via environment variables, which take
precedence over the config file:
  ED_API_TOKEN   API token (created under Admin > API Tokens)
  ED_ORG_ID      organization ID
  ED_API_URL     API base URL (defaults to https://api.edgedelta.com)`,
	}
	cmd.AddCommand(newAuthLoginCmd(), newAuthStatusCmd(), newAuthLogoutCmd())
	return cmd
}

func newAuthLoginCmd() *cobra.Command {
	var token, orgID, apiURL string
	var setDefault bool
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Save an API token and organization ID to a profile",
		Example: `  edx auth login --token 00000000-0000-0000-0000-000000000000 --org-id <org-id>
  edx auth login --profile staging --api-url https://api.staging.edgedelta.com --token ... --org-id ...`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if token == "" {
				return fmt.Errorf("--token is required")
			}
			if orgID == "" {
				return fmt.Errorf("--org-id is required")
			}
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			name := flagProfile
			if name == "" {
				name = "default"
			}
			cfg.Profiles[name] = &config.Profile{APIURL: apiURL, OrgID: orgID, APIToken: token}
			if cfg.DefaultProfile == "" || setDefault {
				cfg.DefaultProfile = name
			}
			if err := cfg.Save(); err != nil {
				return err
			}
			path, _ := config.Path()
			fmt.Fprintf(os.Stderr, "Saved profile %q to %s\n", name, path)
			return nil
		},
	}
	cmd.Flags().StringVar(&token, "token", "", "Edge Delta API token (required)")
	cmd.Flags().StringVar(&orgID, "org-id", "", "Edge Delta organization ID (required)")
	cmd.Flags().StringVar(&apiURL, "api-url", "", "API base URL override for this profile")
	cmd.Flags().BoolVar(&setDefault, "set-default", false, "make this profile the default")
	return cmd
}

func newAuthStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the active profile and verify the token against the API",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := config.Resolve(flagProfile, flagAPIURL, flagOrg, flagToken)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Profile: %s\nAPI URL: %s\nOrg ID:  %s\nToken:   %s\n",
				r.Profile, r.APIURL, r.OrgID, maskToken(r.APIToken))

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
