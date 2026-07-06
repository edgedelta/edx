package cli

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/edgedelta/edx/internal/config"
)

// clearEnv resets every edx env var so tests are hermetic.
func clearEnv(t *testing.T) {
	t.Helper()
	for _, e := range []string{
		config.EnvAPIToken, config.EnvOrgID, config.EnvEnv, config.EnvProfile,
		config.EnvAPIURL, config.EnvChatURL, config.EnvAgentURL,
	} {
		t.Setenv(e, "")
	}
}

// runEdx builds a fresh root command and runs it with args, returning its error.
func runEdx(t *testing.T, args ...string) error {
	t.Helper()
	root := NewRootCmd()
	root.SetArgs(args)
	root.SetOut(os.Stderr)
	return root.Execute()
}

// feedStdin replaces os.Stdin with a pipe carrying s for the duration of fn.
// readCookie treats a non-terminal stdin as piped input and reads it to EOF.
func feedStdin(t *testing.T, s string, fn func()) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = old }()
	go func() { _, _ = w.WriteString(s); _ = w.Close() }()
	fn()
}

const testOrg = "00000000-0000-0000-0000-000000000000"

// cookieAPIServer accepts the probe only when the ed-admin-session cookie matches.
func cookieAPIServer(t *testing.T, cookie string, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/v1/orgs/"+testOrg+"/facet_keys") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if got := r.Header.Get("Cookie"); got != "ed-admin-session="+cookie {
			t.Errorf("probe Cookie header = %q, want ed-admin-session=%s", got, cookie)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`{}`))
	}))
}

func TestCleanCookie(t *testing.T) {
	cases := map[string]string{
		"  abc123  \n":                     "abc123",
		"\x1b[200~tok-en\x1b[201~":         "tok-en",
		"\x1b[200~  spaced  \x1b[201~\r\n": "spaced",
		"plain":                            "plain",
	}
	for in, want := range cases {
		if got := cleanCookie(in); got != want {
			t.Errorf("cleanCookie(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestAuthLoginCookieSavesProfile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("EDX_CONFIG", filepath.Join(dir, "config.yaml"))
	clearEnv(t)

	srv := cookieAPIServer(t, "cookie-xyz", http.StatusOK)
	defer srv.Close()
	t.Setenv(config.EnvAPIURL, srv.URL)

	feedStdin(t, "cookie-xyz\n", func() {
		if err := runEdx(t, "auth", "login", "--org-id", testOrg, "--cookie"); err != nil {
			t.Fatalf("auth login --cookie: %v", err)
		}
	})

	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	p := cfg.Profiles["default"]
	if p == nil {
		t.Fatal("default profile not created")
	}
	if p.AuthMethod != config.AuthMethodCookie || p.SessionCookie != "cookie-xyz" || p.OrgID != testOrg {
		t.Errorf("cookie profile fields wrong: %+v", p)
	}
	if cfg.DefaultProfile != "default" {
		t.Errorf("cookie login should set default profile, got %q", cfg.DefaultProfile)
	}
}

func TestAuthLoginCookieRequiresOrgID(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("EDX_CONFIG", filepath.Join(dir, "config.yaml"))
	clearEnv(t)

	feedStdin(t, "cookie-xyz\n", func() {
		if err := runEdx(t, "auth", "login", "--cookie"); err == nil {
			t.Fatal("expected error: --cookie without --org-id")
		}
	})
}

func TestAuthLoginCookieRejectsBadCookie(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("EDX_CONFIG", filepath.Join(dir, "config.yaml"))
	clearEnv(t)

	srv := cookieAPIServer(t, "bad", http.StatusUnauthorized)
	defer srv.Close()
	t.Setenv(config.EnvAPIURL, srv.URL)

	feedStdin(t, "bad\n", func() {
		if err := runEdx(t, "auth", "login", "--org-id", testOrg, "--cookie"); err == nil {
			t.Fatal("expected error when the cookie is rejected")
		}
	})

	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := cfg.Profiles["default"]; ok {
		t.Error("profile must not be saved when verification fails")
	}
}
