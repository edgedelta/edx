package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/edgedelta/edx/internal/config"
	"github.com/spf13/cobra"
)

// configView is the masked, serializable form of a resolved configuration.
// Secret fields are run through maskToken so the output is safe to share.
type configView struct {
	Profile           string `json:"profile"`
	Env               string `json:"env"`
	AuthMethod        string `json:"auth_method"`
	APIURL            string `json:"api_url"`
	ChatURL           string `json:"chat_url"`
	AgentURL          string `json:"agent_url"`
	OrgID             string `json:"org_id"`
	APIToken          string `json:"api_token"`
	OAuthClientID     string `json:"oauth_client_id"`
	OAuthAccessToken  string `json:"oauth_access_token"`
	OAuthRefreshToken string `json:"oauth_refresh_token"`
	OAuthExpiry       string `json:"oauth_expiry"`
}

// newConfigView copies a Resolved into its display form, masking every secret.
// OAuthExpiry is not a secret and is shown verbatim.
func newConfigView(r *config.Resolved) configView {
	return configView{
		Profile:           r.Profile,
		Env:               r.Env,
		AuthMethod:        r.AuthMethod,
		APIURL:            r.APIURL,
		ChatURL:           r.ChatURL,
		AgentURL:          r.AgentURL,
		OrgID:             r.OrgID,
		APIToken:          maskToken(r.APIToken),
		OAuthClientID:     maskToken(r.OAuthClientID),
		OAuthAccessToken:  maskToken(r.OAuthAccessToken),
		OAuthRefreshToken: maskToken(r.OAuthRefreshToken),
		OAuthExpiry:       r.OAuthExpiry,
	}
}

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Inspect edx configuration",
		Long: `Inspect the edx configuration.

  config path   print the config file location
  config show   print the resolved configuration (secrets masked)`,
	}
	cmd.AddCommand(newConfigShowCmd(), newConfigPathCmd())
	return cmd
}

func newConfigShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Print the resolved configuration (secrets masked, no network call)",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := config.Resolve(flagProfile, flagEnv, flagOrg, flagToken)
			if err != nil {
				return err
			}
			data, err := json.Marshal(newConfigView(r))
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

// formatConfigPath renders the path line, noting whether the file exists.
func formatConfigPath(path string, exists bool) string {
	state := "not found"
	if exists {
		state = "exists"
	}
	return fmt.Sprintf("%s (%s)", path, state)
}

func newConfigPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print the config file path",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := config.Path()
			if err != nil {
				return err
			}
			_, statErr := os.Stat(path)
			fmt.Fprintln(os.Stdout, formatConfigPath(path, statErr == nil))
			return nil
		},
	}
}
