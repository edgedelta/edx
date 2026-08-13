package config

import (
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// Concurrent refreshes (auth status checks every profile in parallel) must
// not lose each other's writes: a dropped update can discard a rotated
// refresh token, permanently killing that profile.
func TestSaveOAuthTokens_ConcurrentProfilesKeepAllUpdates(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("EDX_CONFIG", filepath.Join(dir, "config.yaml"))
	for _, e := range []string{EnvAPIToken, EnvOrgID, EnvEnv, EnvProfile, EnvAPIURL, EnvChatURL, EnvAgentURL} {
		t.Setenv(e, "")
	}

	const n = 8
	exp := time.Now().Add(time.Hour)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			name := fmt.Sprintf("p%d", i)
			if err := SaveOAuthTokens(name, EnvProd, "org", "cid", "access", "refresh-"+name, exp); err != nil {
				t.Error(err)
			}
		}(i)
	}
	wg.Wait()

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Profiles) != n {
		t.Fatalf("lost concurrent updates: %d of %d profiles survived", len(cfg.Profiles), n)
	}
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("p%d", i)
		if p := cfg.Profiles[name]; p == nil || p.OAuthRefreshToken != "refresh-"+name {
			t.Errorf("profile %s lost its refresh token: %+v", name, p)
		}
	}
}
