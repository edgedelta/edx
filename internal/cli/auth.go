package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/edgedelta/edx/internal/api"
	"github.com/edgedelta/edx/internal/config"
	"github.com/edgedelta/edx/internal/oauth"
)

// cookieLifetime is how long an ed-admin-session cookie stays valid. The
// cookie is opaque, so this drives the estimated expiry stored at login time.
const cookieLifetime = 24 * time.Hour

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
	var useOAuth, setDefault, useCookie, useDevice, force bool
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Save credentials to a profile (OAuth by default, an API token, or a support cookie)",
		Example: `  edx auth login                                          # OAuth in your browser (default)
  edx auth login --profile staging --env staging          # OAuth against staging
  edx auth login --device                                 # no local browser (SSH/CI): enter a code in a browser
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
			if useDevice && (token != "" || useCookie) {
				return fmt.Errorf("--device cannot be combined with --token or --cookie")
			}

			// Device authorization flow (--device): for machines without a usable
			// browser (SSH, CI, headless). Yields the same OAuth token pair as the
			// default browser flow, so auto-refresh works the same way.
			if useDevice {
				return loginWithOAuth(cmd, eps, env, name, orgID, setDefault, true, "")
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
				cfg.Profiles[name] = newCookieProfile(env, orgID, cookie, time.Now())
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
			// With no local browser (SSH/headless), this auto-falls back to the
			// device flow.
			if token == "" {
				return loginWithOAuth(cmd, eps, env, name, orgID, setDefault, false, "")
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
	cmd.Flags().BoolVar(&useDevice, "device", false, "log in without a local browser using the device authorization grant (prints a code to enter in a browser on any device)")
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

	// Status is the offline health estimate ("ok (auto-refresh)",
	// "expires in 5h", "expired 3h ago", "unknown …", or "-" when the
	// method has no locally knowable expiry).
	Status string `json:"status"`
	// ExpiresAt is the RFC3339 credential expiry when one is known.
	ExpiresAt string `json:"expires_at,omitempty"`
	// Check is the live verification result ("ok" or "failed: …"),
	// populated only by `auth list --check`.
	Check string `json:"check,omitempty"`
}

// newCookieProfile builds a cookie-auth profile, stamping the estimated
// expiry so `auth list` can warn before the cookie dies.
func newCookieProfile(env, orgID, cookie string, now time.Time) *config.Profile {
	return &config.Profile{
		Env:           env,
		OrgID:         orgID,
		AuthMethod:    config.AuthMethodCookie,
		SessionCookie: cookie,
		CookieExpiry:  now.Add(cookieLifetime).UTC().Format(time.RFC3339),
	}
}

// profileStatus estimates a profile's credential health without touching the
// network. It returns the human status and, when known, the RFC3339 expiry.
func profileStatus(p *config.Profile, now time.Time) (status, expiresAt string) {
	switch p.AuthMethod {
	case config.AuthMethodOAuth:
		if p.OAuthRefreshToken != "" {
			// A stale access token refreshes transparently, so the pair
			// is healthy as long as the refresh token is accepted.
			return "ok (auto-refresh)", p.OAuthExpiry
		}
		return expiryStatus(p.OAuthExpiry, now, "unknown")
	case config.AuthMethodCookie:
		return expiryStatus(p.CookieExpiry, now, "unknown (re-login to track)")
	default: // static API tokens carry no local expiry
		return "-", ""
	}
}

// expiryStatus phrases an RFC3339 expiry relative to now, falling back to
// unknown when the timestamp is missing or malformed.
func expiryStatus(expiry string, now time.Time, unknown string) (string, string) {
	if expiry == "" {
		return unknown, ""
	}
	t, err := time.Parse(time.RFC3339, expiry)
	if err != nil {
		return unknown, ""
	}
	if t.Before(now) {
		return "expired " + humanDur(now.Sub(t)) + " ago", expiry
	}
	return "expires in " + humanDur(t.Sub(now)), expiry
}

// humanDur renders a duration at the coarsest useful grain: "<1m", "42m",
// "1h30m", "5h", "3d".
func humanDur(d time.Duration) string {
	if d < time.Minute {
		return "<1m"
	}
	d = d.Round(time.Minute)
	switch {
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 48*time.Hour:
		h, m := int(d.Hours()), int(d.Minutes())%60
		if m == 0 {
			return fmt.Sprintf("%dh", h)
		}
		return fmt.Sprintf("%dh%dm", h, m)
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
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

	now := time.Now()
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
		status, expiresAt := profileStatus(p, now)
		entries = append(entries, profileListEntry{
			Name:      name,
			Env:       env,
			OrgID:     p.OrgID,
			Auth:      auth,
			Default:   name == f.DefaultProfile,
			Status:    status,
			ExpiresAt: expiresAt,
		})
	}
	return entries
}

// formatProfileList renders the saved profiles as an aligned table. The default
// profile is prefixed with "* " and org IDs are shortened for readability.
func formatProfileList(f *config.File) string {
	return formatProfileEntries(profileListEntries(f))
}

// formatProfileEntries renders pre-computed rows (see formatProfileList).
// STATUS shows the live --check result when present, the offline estimate
// otherwise. Styling goes only on the last column so ANSI codes cannot skew
// tabwriter's alignment of the columns before it.
func formatProfileEntries(entries []profileListEntry) string {
	if len(entries) == 0 {
		return "No profiles yet. Run `edx auth login` to create one.\n"
	}
	var sb strings.Builder
	tw := tabwriter.NewWriter(&sb, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "  NAME\tENV\tORG\tAUTH\tSTATUS")
	for _, e := range entries {
		marker := "  "
		if e.Default {
			marker = "* "
		}
		org := "-"
		if e.OrgID != "" {
			org = shortID(e.OrgID)
		}
		fmt.Fprintf(tw, "%s%s\t%s\t%s\t%s\t%s\n", marker, e.Name, e.Env, org, e.Auth, statusCell(e))
	}
	_ = tw.Flush()
	return sb.String()
}

// statusCell picks and styles the STATUS column value for one row.
func statusCell(e profileListEntry) string {
	if e.Check != "" {
		if e.Check == "ok" {
			return okMark() + " ok"
		}
		return failMark() + " " + red(strings.TrimPrefix(e.Check, "failed: "))
	}
	switch {
	case strings.HasPrefix(e.Status, "expired"):
		return red(e.Status)
	case strings.HasPrefix(e.Status, "unknown"):
		return dim(e.Status)
	case strings.HasPrefix(e.Status, "expires in"):
		if t, err := time.Parse(time.RFC3339, e.ExpiresAt); err == nil && time.Until(t) < 2*time.Hour {
			return yellow(e.Status)
		}
	}
	return e.Status
}

// verifyAuth makes the cheapest authenticated API call that exercises the
// stored credentials. Shared by `auth status` and `auth list --check`.
func verifyAuth(ctx context.Context, c *api.Client) error {
	q := url.Values{}
	q.Set("scope", "log")
	_, err := c.Get(ctx, "/facet_keys", q)
	return err
}

// plainCheck is the machine-readable form of a live check result. Auth
// rejections are compacted (the request URL adds nothing in a status column);
// other errors keep their first line, truncated.
func plainCheck(err error) string {
	if err == nil {
		return "ok"
	}
	msg := strings.TrimSpace(err.Error())
	switch {
	case strings.Contains(msg, "API error 401"):
		return "failed: rejected (401 unauthorized — credentials expired or revoked)"
	case strings.Contains(msg, "API error 403"):
		return "failed: rejected (403 forbidden)"
	}
	if i := strings.IndexByte(msg, '\n'); i >= 0 {
		msg = msg[:i]
	}
	const maxMsg = 80
	if len(msg) > maxMsg {
		msg = msg[:maxMsg] + "…"
	}
	return "failed: " + msg
}

// checkLabel is the styled table form of a live check result.
func checkLabel(err error) string {
	return statusCell(profileListEntry{Check: plainCheck(err)})
}

func newAuthListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List saved profiles with estimated credential status (the default is marked with *)",
		Long: `List saved profiles. Always instant and offline.

STATUS is estimated from stored expiry data: OAuth profiles refresh
themselves, cookie profiles expire ~24h after login, and static API tokens
have no locally knowable expiry. Use "edx auth status" to verify every
profile live against its API.

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

// runProfileChecks verifies every entry concurrently via check and records
// the results. Checks are independent per profile; concurrent token refreshes
// are safe because config.SaveOAuthTokens serializes its writes.
func runProfileChecks(ctx context.Context, entries []profileListEntry, check func(ctx context.Context, name string) error) {
	var wg sync.WaitGroup
	for i := range entries {
		wg.Add(1)
		go func(e *profileListEntry) {
			defer wg.Done()
			e.Check = plainCheck(check(ctx, e.Name))
		}(&entries[i])
	}
	wg.Wait()
}

// checkProfileAuth verifies one saved profile's credentials against its API.
// It resolves the profile by name only — the global --env/--org/--token
// overrides are ignored so each row reflects what is actually stored.
func checkProfileAuth(ctx context.Context, name string) error {
	r, err := config.Resolve(name, "", "", "")
	if err != nil {
		return err
	}
	c, err := clientFromResolved(r)
	if err != nil {
		return err
	}
	return verifyAuth(ctx, c)
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

// loginWithOAuth runs the interactive OAuth login and saves the resulting tokens
// to the named profile. It uses the loopback browser flow by default and the
// device authorization flow when forceDevice is set or no local browser is
// available (SSH/headless) — both yield the same refreshable token pair. When
// email is non-empty (signup), it always uses the device flow and asks the
// server to email a passwordless magic link. Shared by "auth login" and "signup".
func loginWithOAuth(cmd *cobra.Command, eps config.Endpoints, env, name, orgID string, setDefault, forceDevice bool, email string) error {
	useDevice := forceDevice || email != ""
	if !useDevice && !oauth.BrowserAvailable() {
		useDevice = true
		fmt.Fprintln(os.Stderr, dim("No local browser detected; using device authorization."))
	}

	var toks oauth.Tokens
	var err error
	switch {
	case email != "":
		// Signup: the server emails a magic link; the user clicks it, then
		// confirms the code. No local browser is opened.
		toks, err = oauth.DeviceLogin(cmd.Context(), eps.API, oauth.DeviceLoginOptions{
			Email: email,
			Prompt: func(userCode, verificationURI, _ string, _ bool) {
				fmt.Fprintf(os.Stderr, "\n  We emailed a sign-in link to %s.\n  Open it, then confirm this code:\n\n    %s\n\n", email, userCode)
				fmt.Fprintf(os.Stderr, dim("  No email? Visit %s and enter the code.\n"), verificationURI)
				fmt.Fprintln(os.Stderr, dim("  Waiting for you to confirm…"))
			},
		})
	case useDevice:
		fmt.Fprintf(os.Stderr, "Signing in to %s via device authorization…\n", hostOnly(eps.API))
		toks, err = oauth.DeviceLogin(cmd.Context(), eps.API, oauth.DeviceLoginOptions{
			OpenBrowser: true,
			Prompt: func(userCode, verificationURI, _ string, opened bool) {
				fmt.Fprintf(os.Stderr, "\n  To finish signing in, open:\n    %s\n  and enter the code:\n    %s\n\n", verificationURI, userCode)
				if opened {
					fmt.Fprintln(os.Stderr, dim("  We opened your browser to the verification page."))
				}
				fmt.Fprintln(os.Stderr, dim("  Waiting for you to approve…"))
			},
		})
	default:
		fmt.Fprintf(os.Stderr, "Signing in to %s via your browser…\n", hostOnly(eps.API))
		toks, err = oauth.Login(cmd.Context(), eps.API, oauth.LoginOptions{
			OpenBrowser: true,
			Prompt: func(u string, opened bool) {
				if !opened {
					fmt.Fprintf(os.Stderr, "  Couldn't open a browser automatically — open this URL to continue:\n  %s\n", u)
				}
				fmt.Fprintln(os.Stderr, dim("  Waiting for you to authorize in the browser…"))
			},
		})
	}
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	// The org is carried in the access token; derive it so the user need not
	// pass --org-id. An explicit --org-id still overrides.
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
		Use:   "status [profile]",
		Short: "Live-check every saved profile against its API (or one profile in detail)",
		Long: `Verify saved credentials against their APIs.

Without arguments, every profile is checked concurrently (one cheap
authenticated call each) and the results are shown as a table, followed by
the active profile's endpoints. Exits non-zero if any profile fails, so it
can gate scripts.

With a profile name (or --profile), shows that profile's resolved endpoints
and credential in detail and verifies just that one.`,
		Example: `  edx auth status            # check all profiles
  edx auth status staging    # inspect and check one profile`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profile := flagProfile
			if len(args) == 1 {
				profile = args[0]
			}
			if profile != "" {
				return authStatusOne(cmd, profile)
			}
			return authStatusAll(cmd)
		},
	}
}

// credentialSummary renders a masked one-line description of the credential
// a resolved profile would send.
func credentialSummary(r *config.Resolved) string {
	switch {
	case r.UsesOAuth():
		return maskToken(r.OAuthAccessToken) + " (auto-refreshed)"
	case r.UsesCookie():
		return maskToken(r.SessionCookie) + " (ed-admin-session cookie)"
	default:
		return maskToken(r.APIToken)
	}
}

// authStatusAll live-checks every saved profile concurrently and prints the
// roster, then the active profile's endpoints. Non-zero exit if any failed.
func authStatusAll(cmd *cobra.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	entries := profileListEntries(cfg)
	if len(entries) == 0 {
		fmt.Fprint(os.Stdout, formatProfileEntries(entries))
		return nil
	}
	notef("Verifying %d profile(s)…", len(entries))
	runProfileChecks(cmdContext(cmd), entries, checkProfileAuth)

	if cmd.Flags().Changed("output") {
		data, err := json.Marshal(entries)
		if err != nil {
			return err
		}
		if err := printResult(data); err != nil {
			return err
		}
	} else {
		fmt.Fprint(os.Stdout, formatProfileEntries(entries))
		if r, err := config.Resolve("", flagEnv, flagOrg, flagToken); err == nil {
			fmt.Fprintf(os.Stderr, "\nActive profile: %s (%s, org %s)\n  API URL:    %s\n  Credential: %s\n",
				r.Profile, r.Env, shortID(r.OrgID), r.APIURL, credentialSummary(r))
		}
	}

	failed := 0
	for _, e := range entries {
		if e.Check != "ok" {
			failed++
		}
	}
	if failed > 0 {
		return fmt.Errorf("%d of %d profile(s) failed verification", failed, len(entries))
	}
	return nil
}

// authStatusOne prints one profile's resolved endpoints and credential in
// detail and verifies it against the API.
func authStatusOne(cmd *cobra.Command, profile string) error {
	r, err := config.Resolve(profile, flagEnv, flagOrg, flagToken)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr,
		"Profile:   %s\nEnv:       %s\nAuth:      %s\nAPI URL:   %s\nChat URL:  %s\nAgent URL: %s\nOrg ID:    %s\nCredential:%s\n",
		r.Profile, r.Env, r.AuthMethod, r.APIURL, r.ChatURL, r.AgentURL, r.OrgID, credentialSummary(r))

	c, err := clientFromResolved(r)
	if err != nil {
		return err
	}
	// Cheap authenticated call to verify credential + org pairing.
	if err := verifyAuth(cmdContext(cmd), c); err != nil {
		return fmt.Errorf("token verification failed: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Status:    %s ok (credential accepted)\n", okMark())
	return nil
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
