package cli

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/edgedelta/edx/internal/api"
)

// urlTimeFormat matches the backend's core.URLTimeFormat.
const urlTimeFormat = "2006-01-02T15:04:05.000Z"

const (
	// maxFollowFailures is how many consecutive failed polls --follow tolerates
	// before giving up. Transient blips (429/5xx, dropped connections) must not
	// end a tail, but a permanent failure - expired credentials, deleted
	// pipeline - has to surface as a non-zero exit instead of looking like an
	// idle pipeline.
	maxFollowFailures = 5
	// maxFollowSeen bounds the emitted-record memory of a long-running --follow.
	maxFollowSeen = 100000
)

func newCaptureCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "capture",
		Short: "Live capture: sample data flowing through pipeline nodes",
		Long: `Live capture samples real data before and after pipeline nodes on running
agents - the fastest way to debug a processor or verify data flow.

Workflow:
  1. edx capture start <conf-id> --duration 2m --nodes <node-name>
  2. edx capture status <task-id>        (poll until agents report)
  3. edx capture results <conf-id>       (fetch captured before/after samples)`,
	}
	cmd.AddCommand(
		newCaptureStartCmd(),
		newCaptureTaskCmd(),
		newCaptureStatusCmd(),
		newCaptureResultsCmd(),
	)
	return cmd
}

func newCaptureStartCmd() *cobra.Command {
	var duration time.Duration
	var nodes []string
	var interval string
	var maxItems int
	cmd := &cobra.Command{
		Use:   "start <conf-id>",
		Short: "Start a live capture task on a pipeline",
		Example: `  edx capture start <conf-id> --duration 2m
  edx capture start <conf-id> --duration 5m --nodes mask_pii,route_errors --max-items 50`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if maxItems < 1 || maxItems > 100 {
				return fmt.Errorf("--max-items must be between 1 and 100 (got %d)", maxItems)
			}
			if interval != "" {
				d, err := time.ParseDuration(interval)
				if err != nil {
					return fmt.Errorf("invalid --interval %q: %v", interval, err)
				}
				if d < time.Second {
					return fmt.Errorf("--interval must be at least 1s (got %s)", interval)
				}
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			task := map[string]any{
				"expires_at": time.Now().UTC().Add(duration).Format(urlTimeFormat),
			}
			if len(nodes) > 0 {
				task["nodes_to_capture"] = nodes
			}
			if interval != "" {
				task["polling_interval"] = interval
			}
			if maxItems > 0 {
				task["max_items"] = maxItems
			}
			body, err := json.Marshal(task)
			if err != nil {
				return err
			}
			data, err := c.Post(cmdContext(cmd), "/pipelines/"+url.PathEscape(args[0])+"/capture/task", nil, body)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().DurationVar(&duration, "duration", 2*time.Minute, "how long the capture stays active")
	cmd.Flags().StringSliceVar(&nodes, "nodes", nil, "pipeline node names to capture (default: all)")
	cmd.Flags().StringVar(&interval, "interval", "", "agent polling interval (Go duration, min 1s)")
	cmd.Flags().IntVar(&maxItems, "max-items", 20, "max items to capture per node (max 100)")
	return cmd
}

func newCaptureTaskCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "task <conf-id>",
		Short: "Show the active capture task for a pipeline",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/pipelines/"+url.PathEscape(args[0])+"/capture/task", nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

func newCaptureStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status <task-id>",
		Short: "Show per-agent status for a capture task",
		Long: `Shows which agents picked up the capture task and their reporting status.
The task ID is in the response of "edx capture start" (field "id").`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/pipelines/capture/task/%s/status", url.PathEscape(args[0]))
			data, err := c.Get(cmdContext(cmd), path, nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

func newCaptureResultsCmd() *cobra.Command {
	var ep extraParams
	var follow bool
	var sinceNow bool
	var interval time.Duration
	cmd := &cobra.Command{
		Use:   "results <conf-id>",
		Short: "Fetch captured samples (before/after each node)",
		Long: `Fetch captured samples. Each captured node reports "before" and "after"
arrays; the items are JSON-encoded strings, so decode them with jq's fromjson:

  edx capture results <conf-id> | jq '[.[].nodes[].after[]] | map(fromjson)'

FOLLOW MODE

With --follow, edx polls every --interval and prints one line of compact JSON
per newly captured item - already decoded, ready to pipe into jq, grep or an
agent's log monitor:

  {"timestamp":...,"source":"<agent>","node":"<node>","phase":"before|after","item":{...}}

Use --output raw to print the captured item alone, without the envelope.

The first poll prints everything currently held for the pipeline, which can be
several hundred lines for a wide capture. Pass --since-now for tail -f
semantics: the backlog is swallowed and only items captured after the tail
starts are printed. Either way you can narrow the stream with
--param nodes=<a,b> or --param limit_per_source=1.

Only --output json (the default) and raw are supported while following.

--follow tails the capture results; it does not keep the capture running. An
"edx capture start" task must be active or the stream stays silent, so re-arm
it as needed:

  edx capture start <conf-id> --duration 10m --nodes mask_pii
  edx capture results <conf-id> --follow --param nodes=mask_pii`,
		Example: `  edx capture results <conf-id>
  edx capture results <conf-id> --follow
  edx capture results <conf-id> --follow --since-now
  edx capture results <conf-id> --follow --interval 10s --param nodes=mask_pii
  edx capture results <conf-id> --follow --output raw | jq -r .body`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if err := ep.apply(q); err != nil {
				return err
			}
			if follow {
				if interval < time.Second {
					return fmt.Errorf("--interval must be at least 1s (got %s)", interval)
				}
				switch flagOutput {
				case "", "json", "raw":
				default:
					return fmt.Errorf("--follow supports --output json or raw, not %q", flagOutput)
				}
				return followCapture(cmd, args[0], q, interval, sinceNow)
			}
			if sinceNow {
				return errors.New("--since-now only applies with --follow")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), capturePath(args[0]), q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	ep.register(cmd)
	cmd.Flags().BoolVar(&follow, "follow", false, "stream newly captured items as NDJSON until interrupted")
	cmd.Flags().BoolVar(&sinceNow, "since-now", false, "with --follow, skip the existing backlog and print only items captured from now on")
	cmd.Flags().DurationVar(&interval, "interval", 5*time.Second, "poll interval while following (min 1s)")
	return cmd
}

func capturePath(confID string) string {
	return "/pipelines/" + url.PathEscape(confID) + "/capture"
}

// captureRecord is one captured item, flattened out of the nested
// CapturePayload response into a single line of NDJSON.
type captureRecord struct {
	Timestamp int64  `json:"timestamp"`
	Source    string `json:"source"`
	Node      string `json:"node"`
	Phase     string `json:"phase"`
	Item      any    `json:"item"`

	// key dedups across polls; raw is the item as the API returned it (a
	// JSON-encoded string). Unexported, so neither is serialized.
	key string
	raw string
}

// capturePayload is the subset of core.CapturePayload that --follow reads. The
// full response is []CapturePayload - one payload per agent upload.
type capturePayload struct {
	Timestamp int64  `json:"timestamp"`
	Source    string `json:"source"`
	Nodes     map[string]struct {
		Before []string `json:"before"`
		After  []string `json:"after"`
	} `json:"nodes"`
}

// flattenCapture turns a capture results response into one record per captured
// item, in a stable order so repeated polls emit new items in the same
// sequence. Items that are not valid JSON are passed through as strings.
func flattenCapture(data []byte) ([]captureRecord, error) {
	var payloads []capturePayload
	if err := json.Unmarshal(data, &payloads); err != nil {
		return nil, fmt.Errorf("unexpected capture response shape: %w", err)
	}
	var recs []captureRecord
	for _, p := range payloads {
		nodes := make([]string, 0, len(p.Nodes))
		for name := range p.Nodes {
			nodes = append(nodes, name)
		}
		sort.Strings(nodes)
		for _, name := range nodes {
			nd := p.Nodes[name]
			for _, phase := range []string{"before", "after"} {
				items := nd.Before
				if phase == "after" {
					items = nd.After
				}
				for i, item := range items {
					recs = append(recs, captureRecord{
						Timestamp: p.Timestamp,
						Source:    p.Source,
						Node:      name,
						Phase:     phase,
						Item:      decodeCaptureItem(item),
						key:       captureKey(p.Timestamp, p.Source, name, phase, i, item),
						raw:       item,
					})
				}
			}
		}
	}
	sort.SliceStable(recs, func(i, j int) bool { return recs[i].Timestamp < recs[j].Timestamp })
	return recs, nil
}

// decodeCaptureItem parses a captured item, which the API delivers as a
// JSON-encoded string. Anything that does not parse is kept verbatim rather
// than dropped.
func decodeCaptureItem(item string) any {
	var v any
	if err := json.Unmarshal([]byte(item), &v); err != nil {
		return item
	}
	return v
}

// captureKey identifies a captured item across polls. The item's own bytes are
// part of the key, so an agent re-uploading a batch with an item changed at the
// same index still surfaces it.
func captureKey(ts int64, source, node, phase string, idx int, item string) string {
	h := sha256.New()
	for _, part := range []string{strconv.FormatInt(ts, 10), source, node, phase, strconv.Itoa(idx), item} {
		h.Write([]byte(part))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}

// seenSet remembers which records have been emitted. It is bounded: a tail left
// running for days drops its oldest keys rather than growing without limit,
// which at worst re-emits an item the backend still holds long after it first
// appeared.
type seenSet struct {
	keys  map[string]struct{}
	order []string
	max   int
}

func newSeenSet(max int) *seenSet {
	return &seenSet{keys: make(map[string]struct{}), max: max}
}

// add records key and reports whether it is new.
func (s *seenSet) add(key string) bool {
	if _, ok := s.keys[key]; ok {
		return false
	}
	s.keys[key] = struct{}{}
	s.order = append(s.order, key)
	if len(s.order) > s.max {
		drop := len(s.order) / 2
		for _, k := range s.order[:drop] {
			delete(s.keys, k)
		}
		s.order = append([]string(nil), s.order[drop:]...)
	}
	return true
}

// captureTail writes each captured item it has not written before as one line
// on out, and puts operational notes on notes.
type captureTail struct {
	seen     *seenSet
	enc      *json.Encoder
	out      io.Writer
	notes    io.Writer
	rawItems bool
	// skipBacklog stays true until the first poll that actually reaches the API,
	// so a 404 or a retried failure at startup does not consume the skip.
	skipBacklog bool
}

func newCaptureTail(out, notes io.Writer, rawItems, sinceNow bool) *captureTail {
	enc := json.NewEncoder(out)
	enc.SetEscapeHTML(false)
	return &captureTail{
		seen:        newSeenSet(maxFollowSeen),
		enc:         enc,
		out:         out,
		notes:       notes,
		rawItems:    rawItems,
		skipBacklog: sinceNow,
	}
}

// emit writes the records of one poll that have not been written before. The
// first call under --since-now records the backlog as seen without writing it.
func (t *captureTail) emit(recs []captureRecord) error {
	if t.skipBacklog {
		for _, r := range recs {
			t.seen.add(r.key)
		}
		if len(recs) > 0 {
			fmt.Fprintf(t.notes, "skipped %d already-captured item(s); following new ones\n", len(recs))
		}
		// Consumed even on an empty response: anything arriving from the next
		// poll on is genuinely new, not backlog.
		t.skipBacklog = false
		return nil
	}
	for _, r := range recs {
		if !t.seen.add(r.key) {
			continue
		}
		if t.rawItems {
			if _, err := fmt.Fprintln(t.out, r.raw); err != nil {
				return err
			}
			continue
		}
		if err := t.enc.Encode(r); err != nil {
			return err
		}
	}
	return nil
}

// followCapture polls the capture results and prints each newly captured item
// as one line, until the context is cancelled or polling fails persistently.
// When sinceNow is set, whatever the first successful poll returns is recorded
// as already seen and not printed, giving tail -f semantics.
func followCapture(cmd *cobra.Command, confID string, q url.Values, interval time.Duration, sinceNow bool) error {
	c, err := newClient()
	if err != nil {
		return err
	}
	ctx := cmdContext(cmd)
	warnIfNoCaptureTask(ctx, c, confID)

	path := capturePath(confID)
	tail := newCaptureTail(os.Stdout, os.Stderr, flagOutput == "raw", sinceNow)
	failures := 0

	for {
		data, err := c.Get(ctx, path, q)
		switch {
		case err != nil && ctx.Err() != nil:
			return nil
		case isNotFound(err):
			// No capture data for this pipeline yet - a normal state while
			// waiting for the first agent upload, not a failure.
			failures = 0
		case err != nil:
			failures++
			fmt.Fprintf(os.Stderr, "capture poll failed (%d/%d): %v\n", failures, maxFollowFailures, err)
			if failures >= maxFollowFailures {
				return fmt.Errorf("giving up after %d consecutive failed polls: %w", failures, err)
			}
		default:
			failures = 0
			recs, err := flattenCapture(data)
			if err != nil {
				return err
			}
			if err := tail.emit(recs); err != nil {
				return err
			}
		}

		select {
		case <-time.After(interval):
		case <-ctx.Done():
			return nil
		}
	}
}

// warnIfNoCaptureTask notes on stderr when nothing is arming the capture, so a
// silent --follow is not mistaken for an idle pipeline. Best effort: a failed
// check is not worth ending the tail over.
func warnIfNoCaptureTask(ctx context.Context, c *api.Client, confID string) {
	data, err := c.Get(ctx, "/pipelines/"+url.PathEscape(confID)+"/capture/task", nil)
	if err != nil {
		if isNotFound(err) {
			fmt.Fprintf(os.Stderr, "warning: no active capture task for %s; run `edx capture start %s --duration 10m`\n", confID, confID)
		}
		return
	}
	var task struct {
		ID        string `json:"id"`
		ExpiresAt string `json:"expires_at"`
	}
	if err := json.Unmarshal(data, &task); err != nil || task.ID == "" {
		fmt.Fprintf(os.Stderr, "warning: no active capture task for %s; run `edx capture start %s --duration 10m`\n", confID, confID)
		return
	}
	if expiry, err := parseCaptureTime(task.ExpiresAt); err == nil && time.Now().After(expiry) {
		fmt.Fprintf(os.Stderr, "warning: capture task %s expired at %s; re-arm with `edx capture start %s --duration 10m`\n", task.ID, task.ExpiresAt, confID)
	}
}

func parseCaptureTime(s string) (time.Time, error) {
	if t, err := time.Parse(urlTimeFormat, s); err == nil {
		return t, nil
	}
	return time.Parse(time.RFC3339, s)
}

func isNotFound(err error) bool {
	var apiErr *api.Error
	return errors.As(err, &apiErr) && apiErr.Status == http.StatusNotFound
}
