// Package config manages edx CLI configuration: profiles stored in
// ~/.config/edx/config.yaml plus environment variable and flag overrides.
//
// A profile names an Edge Delta environment ("prod", "staging" or "local").
// The environment resolves to a *set* of service base URLs — the main API
// plus the AI Teammate services (chat, agent), which live on their own hosts
// (e.g. chat.ai.edgedelta.com) rather than under api.edgedelta.com. This is
// why a single api-url is not enough: switching environments must move every
// service endpoint together.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Edge Delta environments. Each maps to a full set of service endpoints.
const (
	EnvProd    = "prod"
	EnvStaging = "staging"
	EnvLocal   = "local"

	DefaultEnv = EnvProd
)

// Authentication methods a profile can use.
const (
	// AuthMethodToken sends a static API token (X-ED-API-Token). Default.
	AuthMethodToken = "token"
	// AuthMethodOAuth sends a Bearer JWT obtained via the OAuth login flow,
	// refreshing it when it expires.
	AuthMethodOAuth = "oauth"
)

const (
	EnvAPIToken = "ED_API_TOKEN"
	EnvOrgID    = "ED_ORG_ID"
	EnvEnv      = "ED_ENV"
	EnvProfile  = "EDX_PROFILE"

	// Per-service host overrides. These are an undocumented escape hatch for
	// pointing at a non-standard deployment (a branch deploy, an odd local
	// port); when set they override only that service's host for the resolved
	// environment. Deliberately env-var only — no flags or profile fields — so
	// the common path never sees a URL.
	EnvAPIURL   = "ED_API_URL"
	EnvChatURL  = "ED_CHAT_URL"
	EnvAgentURL = "ED_AGENT_URL"
)

// Endpoints is the full set of service base URLs for one environment.
// The AI Teammate features (issues, threads, channels) are served by Chat;
// teammates (agents) by Agent. Both are distinct hosts from the main API.
type Endpoints struct {
	API   string
	Chat  string
	Agent string
}

// envEndpoints mirrors backend/core/ai.go's service endpoint maps so edx
// targets the same hosts the product uses.
var envEndpoints = map[string]Endpoints{
	EnvProd: {
		API:   "https://api.edgedelta.com",
		Chat:  "https://chat.ai.edgedelta.com",
		Agent: "https://agent.ai.edgedelta.com",
	},
	EnvStaging: {
		API:   "https://api.staging.edgedelta.com",
		Chat:  "https://chat.ai.staging.edgedelta.com",
		Agent: "https://agent.ai.staging.edgedelta.com",
	},
	EnvLocal: {
		API:   "http://localhost:4444",
		Chat:  "http://localhost:3001",
		Agent: "http://localhost:3002",
	},
}

// KnownEnvs returns the supported environment names, prod first.
func KnownEnvs() []string {
	return []string{EnvProd, EnvStaging, EnvLocal}
}

// EndpointsForEnv returns the service endpoints for a named environment.
func EndpointsForEnv(env string) (Endpoints, bool) {
	e, ok := envEndpoints[env]
	return e, ok
}

// Profile holds the credentials and environment for one Edge Delta org.
// The environment is the only knob for endpoints — there are no per-service
// URL overrides; switch hosts by switching env.
type Profile struct {
	// Env selects the environment ("prod", "staging", "local") and thus the
	// full set of service endpoints.
	Env   string `yaml:"env,omitempty"`
	OrgID string `yaml:"org_id,omitempty"`

	// AuthMethod is "token" (default) or "oauth".
	AuthMethod string `yaml:"auth_method,omitempty"`

	// APIToken is used when AuthMethod is "token" (or as a fallback for
	// endpoints that do not yet accept OAuth — see api.oauthUnsupported).
	APIToken string `yaml:"api_token,omitempty"`

	// OAuth* fields hold the credentials minted by `edx auth login --oauth`.
	// The access token is refreshed automatically using the refresh token.
	OAuthClientID     string `yaml:"oauth_client_id,omitempty"`
	OAuthAccessToken  string `yaml:"oauth_access_token,omitempty"`
	OAuthRefreshToken string `yaml:"oauth_refresh_token,omitempty"`
	OAuthExpiry       string `yaml:"oauth_expiry,omitempty"` // RFC3339
}

// File is the on-disk configuration document.
type File struct {
	DefaultProfile string              `yaml:"default_profile,omitempty"`
	Profiles       map[string]*Profile `yaml:"profiles,omitempty"`
}

// Path returns the config file location, honoring EDX_CONFIG override.
func Path() (string, error) {
	if p := os.Getenv("EDX_CONFIG"); p != "" {
		return p, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "edx", "config.yaml"), nil
}

// Load reads the config file. A missing file yields an empty config, not an error.
func Load() (*File, error) {
	path, err := Path()
	if err != nil {
		return nil, err
	}
	f := &File{Profiles: map[string]*Profile{}}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return f, nil
	}
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, f); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, err)
	}
	if f.Profiles == nil {
		f.Profiles = map[string]*Profile{}
	}
	return f, nil
}

// Save writes the config file with 0600 permissions (it contains tokens).
func (f *File) Save() error {
	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := yaml.Marshal(f)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// Resolved is the effective configuration after merging profile, env and flags.
type Resolved struct {
	Profile  string
	Env      string
	APIURL   string
	ChatURL  string
	AgentURL string
	OrgID    string

	AuthMethod string
	APIToken   string

	OAuthClientID     string
	OAuthAccessToken  string
	OAuthRefreshToken string
	OAuthExpiry       string
}

// UsesOAuth reports whether the resolved profile authenticates via OAuth.
func (r *Resolved) UsesOAuth() bool {
	return r.AuthMethod == AuthMethodOAuth && r.OAuthAccessToken != ""
}

// Resolve picks the environment and credentials. The environment
// (envFlag > ED_ENV > profile.Env > default) selects every service endpoint;
// org and token may still be overridden by env vars and flags.
func Resolve(profileFlag, envFlag, orgFlag, tokenFlag string) (*Resolved, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	name := profileFlag
	if name == "" {
		name = os.Getenv(EnvProfile)
	}
	if name == "" {
		name = cfg.DefaultProfile
	}
	if name == "" {
		name = "default"
	}

	var p *Profile
	if found, ok := cfg.Profiles[name]; ok {
		p = found
	} else if profileFlag != "" {
		return nil, fmt.Errorf("profile %q not found in config (run `edx auth login --profile %s`)", name, name)
	}

	env := envFlag
	if env == "" {
		env = os.Getenv(EnvEnv)
	}
	if env == "" && p != nil {
		env = p.Env
	}
	if env == "" {
		env = DefaultEnv
	}

	eps, ok := EndpointsForEnv(env)
	if !ok {
		return nil, fmt.Errorf("unknown environment %q (valid: %s)", env, strings.Join(KnownEnvs(), ", "))
	}

	r := &Resolved{
		Profile:  name,
		Env:      env,
		APIURL:   eps.API,
		ChatURL:  eps.Chat,
		AgentURL: eps.Agent,
	}
	if p != nil {
		r.OrgID = p.OrgID
		r.AuthMethod = p.AuthMethod
		r.APIToken = p.APIToken
		r.OAuthClientID = p.OAuthClientID
		r.OAuthAccessToken = p.OAuthAccessToken
		r.OAuthRefreshToken = p.OAuthRefreshToken
		r.OAuthExpiry = p.OAuthExpiry
	}
	if r.AuthMethod == "" {
		r.AuthMethod = AuthMethodToken
	}

	// Per-service host overrides (escape hatch): override a single service's
	// host without disturbing the others or the chosen environment.
	if v := os.Getenv(EnvAPIURL); v != "" {
		r.APIURL = v
	}
	if v := os.Getenv(EnvChatURL); v != "" {
		r.ChatURL = v
	}
	if v := os.Getenv(EnvAgentURL); v != "" {
		r.AgentURL = v
	}

	if v := os.Getenv(EnvOrgID); v != "" {
		r.OrgID = v
	}
	if v := os.Getenv(EnvAPIToken); v != "" {
		r.APIToken = v
	}
	if orgFlag != "" {
		r.OrgID = orgFlag
	}
	if tokenFlag != "" {
		r.APIToken = tokenFlag
	}
	return r, nil
}

// SaveOAuthTokens persists a refreshed OAuth credential set onto a profile,
// creating the profile if it does not exist. Used both by the login flow and
// by the auto-refresh path so a renewed access token survives across runs.
func SaveOAuthTokens(profileName, env, orgID, clientID, access, refresh string, expiry time.Time) error {
	cfg, err := Load()
	if err != nil {
		return err
	}
	p := cfg.Profiles[profileName]
	if p == nil {
		p = &Profile{}
		cfg.Profiles[profileName] = p
	}
	if env != "" {
		p.Env = env
	}
	if orgID != "" {
		p.OrgID = orgID
	}
	p.AuthMethod = AuthMethodOAuth
	p.OAuthClientID = clientID
	p.OAuthAccessToken = access
	if refresh != "" {
		p.OAuthRefreshToken = refresh
	}
	p.OAuthExpiry = expiry.UTC().Format(time.RFC3339)
	if cfg.DefaultProfile == "" {
		cfg.DefaultProfile = profileName
	}
	return cfg.Save()
}
