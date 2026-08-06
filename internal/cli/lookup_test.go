package cli

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/edgedelta/edx/internal/config"
)

// lookupCapture records what the mock API server received for one request.
type lookupCapture struct {
	method   string
	path     string
	filename string
	file     []byte
	form     map[string]string
	hits     int
}

// lookupAPIServer serves the main API host, capturing multipart uploads and
// replying with resp.
func lookupAPIServer(t *testing.T, got *lookupCapture, resp string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got.hits++
		got.method = r.Method
		got.path = r.URL.Path
		if err := r.ParseMultipartForm(1 << 20); err == nil {
			got.form = map[string]string{}
			for k, vs := range r.MultipartForm.Value {
				if len(vs) > 0 {
					got.form[k] = vs[0]
				}
			}
			if f, h, err := r.FormFile("data"); err == nil {
				got.filename = h.Filename
				got.file, _ = io.ReadAll(f)
				f.Close()
			}
		}
		_, _ = w.Write([]byte(resp))
	}))
}

// useAPIEnv points edx at the given main API host with token auth.
func useAPIEnv(t *testing.T, apiURL string) {
	t.Helper()
	t.Setenv("EDX_CONFIG", filepath.Join(t.TempDir(), "config.yaml"))
	clearEnv(t)
	t.Setenv(config.EnvAPIToken, "tok-test")
	t.Setenv(config.EnvOrgID, testOrg)
	t.Setenv(config.EnvAPIURL, apiURL)
}

func TestLookupListGetsMetadata(t *testing.T) {
	var got lookupCapture
	srv := lookupAPIServer(t, &got, `{"metadatas":[{"name":"users.csv","count":10}]}`)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	if err := runEdx(t, "lookup", "list"); err != nil {
		t.Fatalf("list: %v", err)
	}
	if got.method != http.MethodGet {
		t.Errorf("method = %q, want GET", got.method)
	}
	if want := "/v1/orgs/" + testOrg + "/lookup_tables/metadata"; got.path != want {
		t.Errorf("path = %q, want %q", got.path, want)
	}
}

func TestLookupGetFetchesTableMetadata(t *testing.T) {
	var got lookupCapture
	srv := lookupAPIServer(t, &got, `{"metadata":{"name":"users.csv"}}`)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	if err := runEdx(t, "lookup", "get", "users.csv"); err != nil {
		t.Fatalf("get: %v", err)
	}
	if want := "/v1/orgs/" + testOrg + "/lookup_tables/users.csv/metadata"; got.path != want {
		t.Errorf("path = %q, want %q", got.path, want)
	}
}

func TestLookupCreateUploadsMultipart(t *testing.T) {
	var got lookupCapture
	srv := lookupAPIServer(t, &got, `{"name":"users.csv"}`)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	csv := "id,name\n1,ada\n"
	dir := t.TempDir()
	file := filepath.Join(dir, "users.csv")
	if err := os.WriteFile(file, []byte(csv), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := runEdx(t, "lookup", "create", file, "--description", "user directory", "--tags", "auth"); err != nil {
		t.Fatalf("create: %v", err)
	}
	if got.method != http.MethodPost {
		t.Errorf("method = %q, want POST", got.method)
	}
	if want := "/v1/orgs/" + testOrg + "/lookup_tables"; got.path != want {
		t.Errorf("path = %q, want %q", got.path, want)
	}
	if got.filename != "users.csv" {
		t.Errorf("uploaded filename = %q, want users.csv", got.filename)
	}
	if string(got.file) != csv {
		t.Errorf("uploaded data = %q, want %q", got.file, csv)
	}
	if got.form["description"] != "user directory" || got.form["tags"] != "auth" {
		t.Errorf("form fields = %v", got.form)
	}
}

func TestLookupCreateRejectsBadName(t *testing.T) {
	useAPIEnv(t, "http://127.0.0.1:0")
	file := filepath.Join(t.TempDir(), "users.txt")
	if err := os.WriteFile(file, []byte("a,b\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := runEdx(t, "lookup", "create", file); err == nil {
		t.Fatal("create with .txt name should fail before any request")
	}
}

func TestLookupUpdateUploadsUnderTableName(t *testing.T) {
	var got lookupCapture
	srv := lookupAPIServer(t, &got, `{"name":"users.csv"}`)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	csv := "id,name\n2,grace\n"
	// Local file name differs from the table name on purpose: the backend
	// takes the table ID from the uploaded filename, so update must send
	// the file under the target table's name.
	file := filepath.Join(t.TempDir(), "fresh-export.csv")
	if err := os.WriteFile(file, []byte(csv), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := runEdx(t, "lookup", "update", "users.csv", "--file", file); err != nil {
		t.Fatalf("update: %v", err)
	}
	if got.method != http.MethodPut {
		t.Errorf("method = %q, want PUT", got.method)
	}
	if want := "/v1/orgs/" + testOrg + "/lookup_tables/users.csv"; got.path != want {
		t.Errorf("path = %q, want %q", got.path, want)
	}
	if got.filename != "users.csv" {
		t.Errorf("uploaded filename = %q, want users.csv (the table name, not the local file name)", got.filename)
	}
	if string(got.file) != csv {
		t.Errorf("uploaded data = %q, want %q", got.file, csv)
	}
}

func TestLookupDownloadWritesData(t *testing.T) {
	csv := "id,name\n1,ada\n"
	var got lookupCapture
	srv := lookupAPIServer(t, &got,
		`{"data":"`+base64.StdEncoding.EncodeToString([]byte(csv))+`","metadata":{"name":"users.csv"}}`)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	out := filepath.Join(t.TempDir(), "out.csv")
	if err := runEdx(t, "lookup", "download", "users.csv", "--out", out); err != nil {
		t.Fatalf("download: %v", err)
	}
	if want := "/v1/orgs/" + testOrg + "/lookup_tables/users.csv"; got.path != want {
		t.Errorf("path = %q, want %q", got.path, want)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != csv {
		t.Errorf("downloaded data = %q, want %q", data, csv)
	}
}

func TestLookupDeleteConfirmsAndDeletes(t *testing.T) {
	var got lookupCapture
	srv := lookupAPIServer(t, &got, ``)
	defer srv.Close()
	useAPIEnv(t, srv.URL)

	if err := runEdx(t, "lookup", "delete", "users.csv", "--yes"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if got.method != http.MethodDelete {
		t.Errorf("method = %q, want DELETE", got.method)
	}
	if want := "/v1/orgs/" + testOrg + "/lookup_tables/users.csv"; got.path != want {
		t.Errorf("path = %q, want %q", got.path, want)
	}

	// Declined confirmation must not send a request.
	hits := got.hits
	feedStdin(t, "n\n", func() {
		if err := runEdx(t, "lookup", "delete", "users.csv"); err == nil {
			t.Error("declined delete should return an error")
		}
	})
	if got.hits != hits {
		t.Errorf("declined delete still hit the API (%d -> %d)", hits, got.hits)
	}
}
