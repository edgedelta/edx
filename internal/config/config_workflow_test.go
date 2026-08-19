package config

import (
	"path/filepath"
	"testing"
)

// clearWorkflowEnv makes workflow-endpoint resolution hermetic.
func clearWorkflowEnv(t *testing.T) {
	t.Helper()
	t.Setenv("EDX_CONFIG", filepath.Join(t.TempDir(), "config.yaml"))
	for _, e := range []string{EnvAPIToken, EnvOrgID, EnvEnv, EnvProfile, EnvAPIURL, EnvChatURL, EnvAgentURL, EnvWorkflowURL} {
		t.Setenv(e, "")
	}
}

func TestResolveWorkflowEndpointPerEnv(t *testing.T) {
	clearWorkflowEnv(t)

	cases := map[string]string{
		EnvProd:    "https://workflow.ai.edgedelta.com",
		EnvStaging: "https://workflow.ai.staging.edgedelta.com",
		EnvLocal:   "http://localhost:3010",
	}
	for env, want := range cases {
		r, err := Resolve("", env, "", "")
		if err != nil {
			t.Fatalf("Resolve(%s): %v", env, err)
		}
		if r.WorkflowURL != want {
			t.Errorf("env %s: WorkflowURL = %q, want %q", env, r.WorkflowURL, want)
		}
	}
}

func TestResolveWorkflowHostOverride(t *testing.T) {
	clearWorkflowEnv(t)
	t.Setenv(EnvWorkflowURL, "http://localhost:8888")

	r, err := Resolve("", EnvStaging, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if r.WorkflowURL != "http://localhost:8888" {
		t.Errorf("ED_WORKFLOW_URL should override workflow host, got %q", r.WorkflowURL)
	}
	if r.AgentURL != "https://agent.ai.staging.edgedelta.com" {
		t.Errorf("other hosts must be untouched, got agent %q", r.AgentURL)
	}
}
