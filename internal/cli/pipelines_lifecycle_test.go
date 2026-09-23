package cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestPipelineCreate(t *testing.T) {
	for _, tc := range []struct {
		name, environment, subtype string
		invalid                    bool
	}{
		{"docker", "Docker", "", false},
		{"kubernetes", "Kubernetes", "Edge", false},
		{"missing subtype", "Kubernetes", "", true},
		{"invalid environment", "ECS", "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got map[string]string
			calls := 0
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.Method != "POST" || r.URL.Path != "/v1/orgs/"+testOrg+"/confs" {
					t.Errorf("%s: got: %s %s, want: POST confs", tc.name, r.Method, r.URL.Path)
				}
				if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
					t.Errorf("%s: decode: %v", tc.name, err)
				}
				w.WriteHeader(http.StatusCreated)
				if _, err := w.Write([]byte(`{"id":"created"}`)); err != nil {
					t.Errorf("%s: write: %v", tc.name, err)
				}
			}))
			defer srv.Close()
			useAPIEnv(t, srv.URL)
			file := filepath.Join(t.TempDir(), "pipeline.yaml")
			content := "version: v3\nsettings:\n  tag: test\n"
			if err := os.WriteFile(file, []byte(content), 0600); err != nil {
				t.Fatalf("%s: write file: %v", tc.name, err)
			}
			args := []string{"pipelines", "create", "--file", file, "--tag", "test", "--environment", tc.environment}
			if tc.subtype != "" {
				args = append(args, "--fleet-subtype", tc.subtype)
			}
			err := runEdx(t, args...)
			if tc.invalid {
				if err == nil || calls != 0 {
					t.Fatalf("%s: got: err=%v calls=%d, want: validation error and no requests", tc.name, err, calls)
				}
				return
			}
			if err != nil {
				t.Fatalf("%s: create: %v", tc.name, err)
			}
			want := map[string]string{"content": content, "tag": "test", "environment": tc.environment, "fleet_type": "Edge", "fleet_subtype": tc.subtype, "description": "", "runtime": "go", "runtime_version": ""}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("%s: got: %#v, want: %#v", tc.name, got, want)
			}
			if calls != 1 {
				t.Errorf("%s: got: %d calls, want: 1", tc.name, calls)
			}
		})
	}
}

func TestPipelineDelete(t *testing.T) {
	for _, tc := range []struct {
		name    string
		yes     bool
		status  int
		wantErr bool
	}{
		{"confirmed", true, 200, false}, {"denied", true, 403, true}, {"unconfirmed", false, 200, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.Method != "DELETE" || r.URL.Path != "/v1/orgs/"+testOrg+"/confs/test-id" {
					t.Errorf("%s: got: %s %s, want: DELETE conf", tc.name, r.Method, r.URL.Path)
				}
				w.WriteHeader(tc.status)
			}))
			defer srv.Close()
			useAPIEnv(t, srv.URL)
			args := []string{"pipelines", "delete", "test-id"}
			if tc.yes {
				args = append(args, "--yes")
			}
			feedStdin(t, "no\n", func() {
				err := runEdx(t, args...)
				if (err != nil) != tc.wantErr {
					t.Errorf("%s: got: %v, want error: %v", tc.name, err, tc.wantErr)
				}
			})
			want := 0
			if tc.yes {
				want = 1
			}
			if calls != want {
				t.Errorf("%s: got: %d requests, want: %d", tc.name, calls, want)
			}
		})
	}
}

func TestPipelineCreateDoesNotRetry(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()
	useAPIEnv(t, srv.URL)
	file := filepath.Join(t.TempDir(), "pipeline.yaml")
	if err := os.WriteFile(file, []byte("version: v3"), 0600); err != nil {
		t.Fatalf("%s: write: %v", t.Name(), err)
	}
	err := runEdx(t, "pipelines", "create", "--file", file, "--tag", "test", "--environment", "Docker")
	if err == nil || calls != 1 {
		t.Errorf("%s: got: err=%v calls=%d, want: failure with one request", t.Name(), err, calls)
	}
}
