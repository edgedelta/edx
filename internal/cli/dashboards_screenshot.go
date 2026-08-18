package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/edgedelta/edx/internal/api"
)

// Screenshot statuses, as written by the backend onto the dashboard record.
const (
	screenshotPending = "pending"
	screenshotOK      = "ok"
	screenshotError   = "error"
)

// The rendered image is a viewport-sized PNG or an A3 PDF, so this cap is far above any
// real result; it exists so a wrong URL cannot stream unbounded data into memory.
const maxScreenshotBytes = 64 << 20

// screenshotState is the "screenshot" object on a dashboard. Rendering is asynchronous:
// triggering sets status to pending, and the result is written back here.
type screenshotState struct {
	Status string `json:"status,omitempty"`
	// LightURL and DarkURL are the viewport-sized preview the dashboard list draws.
	LightURL string `json:"light_url,omitempty"`
	DarkURL  string `json:"dark_url,omitempty"`
	// FullLightURL and FullDarkURL hold exports — a full-page PNG or a PDF — which the
	// backend keeps apart so rendering one does not replace that preview.
	FullLightURL string `json:"full_light_url,omitempty"`
	FullDarkURL  string `json:"full_dark_url,omitempty"`
	// ReportedError says the dashboard rendered into its error state. The image is still
	// produced and the status is still ok; this says what the image shows.
	ReportedError   bool   `json:"reported_error,omitempty"`
	LastUpdatedAt   string `json:"last_updated_at,omitempty"`
	LastTriggeredAt string `json:"last_triggered_at,omitempty"`
}

// urlFor picks the render this invocation asked for. export covers a full-page PNG and a
// PDF, both of which land in the export slot rather than over the preview.
func (s *screenshotState) urlFor(theme string, export bool) string {
	if s == nil {
		return ""
	}
	if export {
		if theme == "dark" {
			return s.FullDarkURL
		}
		return s.FullLightURL
	}
	if theme == "dark" {
		return s.DarkURL
	}
	return s.LightURL
}

func statusOf(s *screenshotState) string {
	if s == nil || s.Status == "" {
		return "none"
	}
	return s.Status
}

// screenshotRequest is the trigger body. The share hash is minted server-side, so it is
// deliberately absent here.
type screenshotRequest struct {
	From        string            `json:"from,omitempty"`
	To          string            `json:"to,omitempty"`
	SavedView   string            `json:"saved_view,omitempty"`
	FacetParams map[string]string `json:"facet_params,omitempty"`
}

// dashboardGetter is the slice of the API client the wait loop needs, so the loop can be
// tested without a server.
type dashboardGetter interface {
	Get(ctx context.Context, orgRelPath string, query url.Values) ([]byte, error)
}

func newDashboardsScreenshotCmd() *cobra.Command {
	var (
		out       string
		theme     string
		format    string
		savedView string
		from, to  string
		facets    []string
		fullPage  bool
		wait      time.Duration
		interval  time.Duration
	)
	cmd := &cobra.Command{
		Use: "screenshot <dashboard-id>",
		// "pdf" and "png" set the format as well as naming the command, so the obvious
		// guess works. "export" is what the UI calls this, "render" is what it does.
		Aliases: []string{"pdf", "png", "export", "render"},
		Short:   "Render a saved dashboard to a PNG or PDF",
		Long: `Render a saved dashboard with real data and write the image to a file.

This is the check that offline validation cannot do. "validate" proves a definition
matches the schema and that its queries parse; a screenshot proves the dashboard
actually renders and that its widgets return data. Run it after validating, and
look at the image.

Rendering happens in a headless browser server-side, so it is slow — 30s is normal
and a minute is not unusual. Run this as a background task rather than blocking on
it.

This waits up to --wait (default 3m) for the render to finish. Running out of time
is reported apart from every other failure, because it says nothing about the
dashboard: exit code 3, and a JSON object carrying "timed_out": true and the last
status seen. Re-run the same command to pick the render up. A dashboard that did
render but failed is exit code 1, as usual.

PNG or PDF:

  png  the whole dashboard, 1123px wide (A3 at 96dpi) and as tall as it needs to
       be, light or dark. This is the one to read back: anything that can open an
       image can open it. --full-page=false renders a single 1123x1584 viewport
       instead, which is what the dashboard list draws as its preview image.
  pdf  the whole dashboard, paginated across A3 pages, light only. Opening it
       needs a PDF renderer, so this is the copy to hand to a person or attach to
       a ticket rather than the one to check your own work with.

"edx dashboards pdf <id>" and "edx dashboards png <id>" are the same command with
the format already chosen.

A dashboard that fails to load a widget is still rendered: the image shows which
one broke, and a note goes to stderr. Only a dashboard that never settles at all
counts as a failed render.

Authorization comes from the dashboard's stored resource_accesses, which edx
derives on create and update. Widgets that render empty while the same query
returns data elsewhere point at that allowlist rather than at the query.

A dashboard must be saved before it can be rendered, so the loop is:
  validate -> create (tagged "generated" and "preview") -> screenshot -> look ->
  update -> screenshot -> untag "preview"`,
		Example: `  edx dashboards screenshot <id> --out shot.png
  edx dashboards screenshot <id> --from 2026-08-14T00:00:00.000Z --to 2026-08-14T06:00:00.000Z
  edx dashboards screenshot <id> --theme dark --facet env=prod
  edx dashboards pdf <id>                      # whole dashboard, to send to someone`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			format, err := resolveFormat(cmd.CalledAs(), format, cmd.Flags().Changed("format"), out)
			if err != nil {
				return err
			}
			// The PDF is rendered once, as a document, and carries nothing in the dark
			// field, so a "dark PDF" would silently be the light one.
			if format == "pdf" && theme == "dark" {
				return errors.New("a PDF is rendered once as a document; --theme dark is only available for png")
			}
			if theme != "light" && theme != "dark" {
				return fmt.Errorf("invalid --theme %q: expected light or dark", theme)
			}
			if to != "" && from == "" {
				return errors.New("--to needs --from")
			}
			facetParams, parseErr := parseFacetParams(facets)
			if parseErr != nil {
				return parseErr
			}
			if out == "" {
				out = screenshotFilename(id, theme, format)
			}

			c, err := newClient()
			if err != nil {
				return err
			}
			ctx := cmdContext(cmd)

			// Read the current state first: the result is written in place, so this is what
			// tells a fresh render from the one already on the record.
			before, err := fetchScreenshotState(ctx, c, id)
			if err != nil {
				return err
			}

			body, err := json.Marshal(screenshotRequest{
				From: from, To: to, SavedView: savedView, FacetParams: facetParams,
			})
			if err != nil {
				return err
			}
			// A PDF is a whole-document render already, so full_page only means anything
			// for a PNG. Both land in the backend's export slot, away from the preview.
			export := format == "pdf" || fullPage
			query := url.Values{"screenshot_type": {format}}
			if fullPage && format == "png" {
				query.Set("full_page", "true")
			}
			if _, err := c.Post(ctx, "/dashboards/"+url.PathEscape(id)+"/screenshot/trigger", query, body); err != nil {
				if !isScreenshotInProgress(err) {
					return err
				}
				// Another render is already running. Attaching to it is what the caller
				// wanted anyway, and starting a second one is not possible.
				warnf("a screenshot is already being rendered for this dashboard; waiting for that one")
			}

			warnf("rendering %s (this usually takes 30-60s)...", id)
			state, err := awaitScreenshot(ctx, c, id, before, wait, interval)
			if err != nil {
				// A timeout is not a verdict on the dashboard, so say which one this was
				// in a form a script can read rather than only in the error text.
				var stillRendering *renderPendingError
				if errors.As(err, &stillRendering) && out != "-" {
					if result, mErr := json.Marshal(map[string]any{
						"dashboard_id": id,
						"status":       statusOf(stillRendering.state),
						"timed_out":    true,
						"waited":       wait.String(),
					}); mErr == nil {
						_ = printResult(result)
					}
				}
				return err
			}
			if state.Status == screenshotError {
				return fmt.Errorf("the dashboard did not finish rendering (screenshot status %q). "+
					"The renderer waits for the dashboard to report itself ready, so this usually "+
					"means a widget failed rather than that the screenshot service broke: check the "+
					"queries against live data with `edx logs`/`edx metrics`", state.Status)
			}
			imageURL := state.urlFor(theme, export)
			if imageURL == "" {
				return fmt.Errorf("the render reported %q but produced no %s image", state.Status, theme)
			}

			image, err := downloadScreenshot(ctx, &http.Client{Timeout: flagTimeout}, imageURL)
			if err != nil {
				return err
			}
			if state.ReportedError {
				// Not a failure any more: the render succeeded and the image shows which
				// widget broke. Say so, because the file looks normal otherwise.
				warnf("the dashboard rendered with errors — the image shows which widget failed")
			}

			// stdout is the image itself here, so the summary cannot go there.
			if out == "-" {
				if _, err := cmd.OutOrStdout().Write(image); err != nil {
					return err
				}
				warnf("wrote %d bytes to stdout", len(image))
				return nil
			}
			if err := os.WriteFile(out, image, 0o600); err != nil {
				return err
			}
			warnf("wrote %s (%d KiB)", out, len(image)/1024)

			result, err := json.Marshal(map[string]any{
				"path":            out,
				"status":          state.Status,
				"theme":           theme,
				"format":          format,
				"bytes":           len(image),
				"full_page":       fullPage && format == "png",
				"reported_error":  state.ReportedError,
				"image_url":       imageURL,
				"last_updated_at": state.LastUpdatedAt,
			})
			if err != nil {
				return err
			}
			return printResult(result)
		},
	}
	cmd.Flags().StringVar(&out, "out", "", `write the image here ("-" for stdout; default "<dashboard-id>.png")`)
	cmd.Flags().StringVar(&theme, "theme", "light", "light or dark (png only)")
	cmd.Flags().StringVar(&format, "format", "png", "png or pdf")
	cmd.Flags().StringVar(&from, "from", "", "start time, ISO 8601 (2006-01-02T15:04:05.000Z); without it the dashboard's own default range is used")
	cmd.Flags().StringVar(&to, "to", "", "end time, ISO 8601")
	cmd.Flags().StringVar(&savedView, "saved-view", "", "render a saved view by name")
	cmd.Flags().StringArrayVar(&facets, "facet", nil, "set a dashboard variable as key=value, by its key (repeatable)")
	cmd.Flags().BoolVar(&fullPage, "full-page", true, "capture the whole dashboard; --full-page=false renders one viewport, which is also what refreshes the dashboard list's preview image")
	cmd.Flags().DurationVar(&wait, "wait", 3*time.Minute, "how long to wait for the render to finish")
	cmd.Flags().DurationVar(&interval, "poll-interval", 3*time.Second, "how often to check whether the render finished")
	return cmd
}

// resolveFormat picks the output format from, in precedence order: an explicit --format,
// the alias the command was invoked as, and the extension of --out.
//
// Sources that disagree are an error rather than a silent choice: "edx dashboards pdf <id>
// --out shot.png" writing a PDF into a .png file is the kind of thing nobody notices until
// something tries to open it.
func resolveFormat(calledAs, format string, formatSet bool, out string) (string, error) {
	type source struct{ what, format string }
	var sources []source

	if formatSet {
		if format != "png" && format != "pdf" {
			return "", fmt.Errorf("invalid --format %q: expected png or pdf", format)
		}
		sources = append(sources, source{"--format " + format, format})
	}
	if calledAs == "png" || calledAs == "pdf" {
		sources = append(sources, source{`"dashboards ` + calledAs + `"`, calledAs})
	}
	if ext := strings.ToLower(strings.TrimPrefix(path.Ext(out), ".")); ext == "png" || ext == "pdf" {
		sources = append(sources, source{"--out " + out, ext})
	}

	if len(sources) == 0 {
		return "png", nil
	}
	for _, s := range sources[1:] {
		if s.format != sources[0].format {
			return "", fmt.Errorf("%s and %s disagree about the format; drop one", sources[0].what, s.what)
		}
	}
	return sources[0].format, nil
}

// parseFacetParams turns repeated key=value flags into the trigger's facet_params. The keys
// are dashboard variable keys, which the share URL carries as facet-<key>.
func parseFacetParams(facets []string) (map[string]string, error) {
	if len(facets) == 0 {
		return nil, nil
	}
	out := make(map[string]string, len(facets))
	for _, f := range facets {
		k, v, ok := strings.Cut(f, "=")
		if !ok || strings.TrimSpace(k) == "" {
			return nil, fmt.Errorf("invalid --facet %q: expected key=value", f)
		}
		out[k] = v
	}
	return out, nil
}

// screenshotFilename is the default output path: the dashboard ID, so repeated renders of
// different dashboards do not overwrite each other.
//
// The ID is server-assigned but it ends up in a path, so it is reduced to a single bare
// filename in the working directory rather than trusted to be one.
func screenshotFilename(id, theme, format string) string {
	name := path.Base(strings.ReplaceAll(strings.TrimSpace(id), `\`, "/"))
	name = strings.ReplaceAll(name, ":", "-")
	name = strings.TrimLeft(name, ".")
	if name == "" {
		name = "dashboard"
	}
	if theme == "dark" {
		name += "-dark"
	}
	return name + "." + format
}

// isScreenshotInProgress reports whether a trigger was refused because a render is already
// running, which the backend answers with 400 rather than a conflict status.
func isScreenshotInProgress(err error) bool {
	var apiErr *api.Error
	if !errors.As(err, &apiErr) || apiErr.Status != http.StatusBadRequest {
		return false
	}
	return strings.Contains(apiErr.Body, "already in progress")
}

func fetchScreenshotState(ctx context.Context, g dashboardGetter, id string) (*screenshotState, error) {
	data, err := g.Get(ctx, "/dashboards/"+url.PathEscape(id), nil)
	if err != nil {
		return nil, err
	}
	var d struct {
		Screenshot *screenshotState `json:"screenshot"`
	}
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, fmt.Errorf("decode dashboard %q: %w", id, err)
	}
	return d.Screenshot, nil
}

// screenshotSettled reports whether state is a finished render of our own request rather
// than the result that was already on the record.
//
// The backend overwrites one screenshot object in place, so a terminal status alone proves
// nothing: a read that lands before the pending write is visible would see the previous
// run's "ok" and we would hand back a stale image. last_triggered_at is stamped when the
// render is requested and carried through to the result, so it identifies the run. When the
// record was already pending we are attaching to someone else's render and cannot expect a
// new stamp, so any terminal status counts.
func screenshotSettled(before, now *screenshotState) bool {
	if now == nil || now.Status == "" || now.Status == screenshotPending {
		return false
	}
	if before == nil || before.Status == screenshotPending {
		return true
	}
	return now.LastTriggeredAt != before.LastTriggeredAt
}

// renderPendingError says the wait ran out with the render still going. It is not the same
// outcome as a failed render — the dashboard may be perfectly fine — so it carries its own
// exit code, letting a background caller retry instead of reporting a broken dashboard.
type renderPendingError struct {
	wait  time.Duration
	state *screenshotState
}

func (e *renderPendingError) Error() string {
	return fmt.Sprintf("still rendering after %s (status %q); the server caps a render at 60s, "+
		"so a status still pending this late usually means the render was lost — check again "+
		"with this same command, or trigger a new one", e.wait, statusOf(e.state))
}

func (e *renderPendingError) ExitCode() int { return exitRenderPending }

func awaitScreenshot(ctx context.Context, g dashboardGetter, id string, before *screenshotState, wait, interval time.Duration) (*screenshotState, error) {
	if interval <= 0 {
		interval = time.Second
	}
	deadline := time.Now().Add(wait)
	for {
		state, err := fetchScreenshotState(ctx, g, id)
		if err != nil {
			return nil, err
		}
		if screenshotSettled(before, state) {
			return state, nil
		}
		if !time.Now().Before(deadline) {
			return nil, &renderPendingError{wait: wait, state: state}
		}
		select {
		case <-time.After(interval):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// downloadScreenshot fetches the rendered image.
//
// No credentials are sent: the URL is a public object on a storage host unrelated to the
// API, and attaching the org's token would hand it to a third party.
func downloadScreenshot(ctx context.Context, hc *http.Client, rawURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "edx")

	resp, err := hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download screenshot: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxScreenshotBytes))
	if err != nil {
		return nil, fmt.Errorf("download screenshot: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		snippet := strings.TrimSpace(string(body))
		if len(snippet) > 256 {
			snippet = snippet[:256] + "..."
		}
		return nil, fmt.Errorf("download screenshot: HTTP %d from %s: %s", resp.StatusCode, rawURL, snippet)
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("download screenshot: %s returned an empty body", rawURL)
	}
	return body, nil
}
