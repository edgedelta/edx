// Package config manages edx CLI configuration: profiles stored in
// ~/.config/edx/config.yaml plus environment variable and flag overrides.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	DefaultAPIURL = "https://api.edgedelta.com"

	EnvAPIToken = "ED_API_TOKEN"
	EnvOrgID    = "ED_ORG_ID"
	EnvAPIURL   = "ED_API_URL"
	EnvProfile  = "EDX_PROFILE"
)

// Profile holds the credentials and endpoint for one Edge Delta organization.
type Profile struct {
	APIURL   string `yaml:"api_url,omitempty"`
	OrgID    string `yaml:"org_id,omitempty"`
	APIToken string `yaml:"api_token,omitempty"`
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
	APIURL   string
	OrgID    string
	APIToken string
}

// Resolve merges, in increasing precedence: profile from config file,
// environment variables, then explicit flag values (passed as overrides).
func Resolve(profileFlag, apiURLFlag, orgFlag, tokenFlag string) (*Resolved, error) {
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

	r := &Resolved{Profile: name, APIURL: DefaultAPIURL}
	if p, ok := cfg.Profiles[name]; ok {
		if p.APIURL != "" {
			r.APIURL = p.APIURL
		}
		r.OrgID = p.OrgID
		r.APIToken = p.APIToken
	} else if profileFlag != "" {
		return nil, fmt.Errorf("profile %q not found in config (run `edx auth login --profile %s`)", name, name)
	}

	if v := os.Getenv(EnvAPIURL); v != "" {
		r.APIURL = v
	}
	if v := os.Getenv(EnvOrgID); v != "" {
		r.OrgID = v
	}
	if v := os.Getenv(EnvAPIToken); v != "" {
		r.APIToken = v
	}

	if apiURLFlag != "" {
		r.APIURL = apiURLFlag
	}
	if orgFlag != "" {
		r.OrgID = orgFlag
	}
	if tokenFlag != "" {
		r.APIToken = tokenFlag
	}
	return r, nil
}
