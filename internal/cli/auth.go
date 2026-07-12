package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/edgedelta/edx/internal/api"
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

"edx auth login" logs in via OAuth in your browser by default (tokens are
refreshed automatically). Pass --token <api-token> --org-id <org-id> to use a
static API token instead (handy for CI).

Credentials can also be supplied via environment variables, which take
precedence over the config file:
  ED_API_TOKEN   API token (created under Admin > API Tokens)
  ED_ORG_ID      organization ID
  ED_ENV         environment: prod, staging or local

You can keep several logins at once. Name each with --profile at login time,
list them with "edx auth list", and switch the default with "edx auth use
<name>". Per-command, --profile or the EDX_PROFILE env var override the default.`,
	}
	cmd.AddCommand(newAuthLoginCmd(), newAuthListCmd(), newAuthUseCmd(), newAuthStatusCmd(), newAuthLogoutCmd())
	return cmd
}

func newAuthLoginCmd() *cobra.Command {
	var token, orgID, env string
	var useOAuth, setDefault, useCookie, force bool
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Save credentials to a profile (OAuth by default, an API token, or a support cookie)",
		Example: `  edx auth login                                          # OAuth in your browser (default)
  edx auth login --profile staging --env staging          # OAuth against staging
  edx auth login --token 00000000-0000-0000-0000-000000000000 --org-id <org-id>
  edx auth login --org-id <org-id> --cookie               # paste an ed-admin-session cookie (support-org access)`,
		RunE: func(cmd *cobra.Command, args []string) error {
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

			cfg0, err := config.Load()
			if err != nil {
				return err
			}
			if loginWouldClobber(cfg0, name, cmd.Flags().Changed("profile"), force) {
				return fmt.Errorf("profile %q already exists; pass --profile <name> to save under a different name, or --force to overwrite it", name)
			}

			if token != "" && useOAuth {
				return fmt.Errorf("--token and --oauth are mutually exclusive")
			}
			if useCookie && (token != "" || useOAuth) {
				return fmt.Errorf("--cookie cannot be combined with --token or --oauth")
			}

			// Cookie auth: store a pasted ed-admin-session cookie. The cookie is
			// opaque, so --org-id is required (it cannot be derived like OAuth).
			if useCookie {
				if orgID == "" {
					return fmt.Errorf("--org-id is required with --cookie")
				}
				cookie, err := readCookie()
				if err != nil {
					return err
				}
				if cookie == "" {
					return fmt.Errorf("no cookie provided; paste it at the prompt or pipe it via stdin")
				}
				// Verify before saving so a stale/wrong cookie is never stored.
				apiURL := eps.API
				if v := os.Getenv(config.EnvAPIURL); v != "" {
					apiURL = v
				}
				c := api.New(apiURL, eps.Chat, eps.Agent, orgID, &api.Auth{SessionCookie: cookie}, flagTimeout)
				q := url.Values{}
				q.Set("scope", "log")
				if _, err := c.Get(cmdContext(cmd), "/facet_keys", q); err != nil {
					return fmt.Errorf("could not verify the cookie for org %s: %v\n"+
						"check that it is current (cookies expire ~24h) and that the org has support access enabled", shortID(orgID), err)
				}
				cfg, err := config.Load()
				if err != nil {
					return err
				}
				cfg.Profiles[name] = &config.Profile{Env: env, OrgID: orgID, AuthMethod: config.AuthMethodCookie, SessionCookie: cookie}
				if cfg.DefaultProfile == "" || setDefault {
					cfg.DefaultProfile = name
				}
				if err := cfg.Save(); err != nil {
					return err
				}
				fmt.Fprintf(os.Stderr, "%s Signed in with a support cookie — profile %q (env: %s, org %s)\n", okMark(), name, env, shortID(orgID))
				fmt.Fprintln(os.Stderr, dim("  The cookie expires (~24h); re-run this command when it does."))
				return nil
			}

			// OAuth is the default; a static API token is used only with --token.
			if token == "" {
				fmt.Fprintf(os.Stderr, "Signing in to %s via your browser…\n", hostOnly(eps.API))
				toks, err := oauth.Login(cmd.Context(), eps.API, oauth.LoginOptions{
					OpenBrowser: true,
					Prompt: func(u string, opened bool) {
						if !opened {
							fmt.Fprintf(os.Stderr, "  Couldn't open a browser automatically — open this URL to continue:\n  %s\n", u)
						}
						fmt.Fprintln(os.Stderr, dim("  Waiting for you to authorize in the browser…"))
					},
				})
				if err != nil {
					return fmt.Errorf("oauth login failed: %w", err)
				}
				// The org is carried in the access token; derive it so the user
				// need not pass --org-id. An explicit --org-id still overrides.
				if orgID == "" {
					orgID = oauth.OrgIDFromToken(toks.AccessToken)
				}
				if orgID == "" {
					return fmt.Errorf("could not determine organization from the token; pass --org-id")
				}
				if err := config.SaveOAuthTokens(name, env, orgID, toks.ClientID, toks.AccessToken, toks.RefreshToken, toks.Expiry); err != nil {
					return err
				}
				if setDefault {
					if err := setDefaultProfile(name); err != nil {
						return err
					}
				}
				fmt.Fprintf(os.Stderr, "%s Signed in — profile %q (env: %s, org %s)\n", okMark(), name, env, shortID(orgID))
				return nil
			}

			if orgID == "" {
				return fmt.Errorf("--org-id is required with --token")
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
	cmd.Flags().StringVar(&token, "token", "", "use static API token auth instead of OAuth (requires --org-id)")
	cmd.Flags().BoolVar(&useOAuth, "oauth", false, "")
	_ = cmd.Flags().MarkHidden("oauth") // OAuth is the default now; flag kept as a no-op for back-compat
	cmd.Flags().StringVar(&orgID, "org-id", "", "Edge Delta organization ID (required with --token/--cookie; derived from the token for OAuth)")
	cmd.Flags().StringVar(&env, "env", "", "environment for this profile: prod, staging or local (default prod)")
	cmd.Flags().BoolVar(&setDefault, "set-default", false, "make this profile the default")
	cmd.Flags().BoolVar(&useCookie, "cookie", false, "authenticate with a pasted ed-admin-session cookie for support-org access (requires --org-id; prompts for the value)")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite the target profile if it already exists")
	return cmd
}

// loginWouldClobber reports whether saving credentials to profile `name` would
// silently overwrite an existing profile the user did not explicitly name — the
// footgun a forgotten --profile causes (a bare `auth login` targets "default").
// An explicit --profile or --force is the user opting in, so neither clobbers.
func loginWouldClobber(cfg *config.File, name string, explicitName, force bool) bool {
	if explicitName || force {
		return false
	}
	_, exists := cfg.Profiles[name]
	return exists
}

// profileListEntry is the JSON shape of one row from `auth list --json`.
type profileListEntry struct {
	Name    string `json:"name"`
	Env     string `json:"env"`
	OrgID   string `json:"org_id"`
	Auth    string `json:"auth_method"`
	Default bool   `json:"default"`
}

// profileListEntries returns the saved profiles sorted by name, with empty
// fields resolved to their effective defaults. Org IDs are returned in full
// (unlike the table, which shortens them) so machine consumers get the real id.
func profileListEntries(f *config.File) []profileListEntry {
	names := make([]string, 0, len(f.Profiles))
	for name := range f.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)

	entries := make([]profileListEntry, 0, len(names))
	for _, name := range names {
		p := f.Profiles[name]
		env := p.Env
		if env == "" {
			env = config.DefaultEnv
		}
		auth := p.AuthMethod
		if auth == "" {
			auth = config.AuthMethodToken
		}
		entries = append(entries, profileListEntry{
			Name:    name,
			Env:     env,
			OrgID:   p.OrgID,
			Auth:    auth,
			Default: name == f.DefaultProfile,
		})
	}
	return entries
}

// formatProfileList renders the saved profiles as an aligned table. The default
// profile is prefixed with "* " and org IDs are shortened for readability.
func formatProfileList(f *config.File) string {
	entries := profileListEntries(f)
	if len(entries) == 0 {
		return "No profiles yet. Run `edx auth login` to create one.\n"
	}
	var sb strings.Builder
	tw := tabwriter.NewWriter(&sb, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "  NAME\tENV\tORG\tAUTH")
	for _, e := range entries {
		marker := "  "
		if e.Default {
			marker = "* "
		}
		org := "-"
		if e.OrgID != "" {
			org = shortID(e.OrgID)
		}
		fmt.Fprintf(tw, "%s%s\t%s\t%s\t%s\n", marker, e.Name, e.Env, org, e.Auth)
	}
	_ = tw.Flush()
	return sb.String()
}

func newAuthListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List saved profiles (the default is marked with *)",
		Long: `List saved profiles.

Without -o, prints a human table with the default profile marked "*" and org
IDs shortened. With -o (json, yaml, table, csv, raw) the profiles are rendered
in that format with org IDs in full — e.g. "edx auth list -o json | jq".`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			// Default (no explicit -o) keeps the human table with the "*"
			// default marker and shortened org IDs. An explicit -o routes the
			// profiles through the shared renderer like every other command.
			if !cmd.Flags().Changed("output") {
				fmt.Fprint(os.Stdout, formatProfileList(cfg))
				return nil
			}
			data, err := json.Marshal(profileListEntries(cfg))
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

func newAuthUseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "use <profile>",
		Short: "Set the default profile used when --profile/EDX_PROFILE is not given",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if err := config.UseProfile(name); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "%s Now using profile %q\n", okMark(), name)
			return nil
		},
	}
}

// readCookie reads an ed-admin-session cookie value: from piped stdin, or by
// prompting when stdin is an interactive terminal.
//
// The interactive read MUST use raw mode (term.MakeRaw, which clears ICANON) —
// NOT a normal line read and NOT term.ReadPassword (which keeps canonical mode
// on). Canonical mode caps a line at MAX_CANON (~1024 bytes), so a multi-KB
// cookie paste overflows the line buffer and the trailing Enter is never
// delivered, hanging the read. Raw mode has no such limit and keeps the secret
// off the screen.
func readCookie() (string, error) {
	if !fileIsTTY(os.Stdin) {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return cleanCookie(string(data)), nil
	}
	fmt.Fprint(os.Stderr, "Paste your ed-admin-session cookie value (input hidden), then press Enter: ")
	line, err := readLineRaw(os.Stdin)
	fmt.Fprintln(os.Stderr) // nothing was echoed; move to a fresh line.
	if err != nil {
		return "", err
	}
	return cleanCookie(line), nil
}

// readLineRaw reads a single line from an interactive terminal in raw mode,
// terminating at CR or LF. Ctrl-C cancels.
func readLineRaw(f *os.File) (string, error) {
	fd := int(f.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer func() { _ = term.Restore(fd, oldState) }()

	var sb strings.Builder
	buf := make([]byte, 1)
	for {
		n, rerr := f.Read(buf)
		if n > 0 {
			switch buf[0] {
			case '\r', '\n':
				return sb.String(), nil
			case 3: // Ctrl-C
				return "", fmt.Errorf("cancelled")
			default:
				sb.WriteByte(buf[0])
			}
		}
		if rerr != nil {
			if rerr == io.EOF {
				return sb.String(), nil
			}
			return "", rerr
		}
	}
}

// cleanCookie trims surrounding whitespace and strips the bracketed-paste
// escape markers (ESC[200~ … ESC[201~) a terminal may wrap a raw-mode paste in.
func cleanCookie(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "\x1b[200~")
	s = strings.TrimSuffix(s, "\x1b[201~")
	return strings.TrimSpace(s)
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
			if r.UsesCookie() {
				cred = maskToken(r.SessionCookie) + " (ed-admin-session cookie)"
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
