package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

const captureResponse = `[
  {
    "org_id": "o1", "config_id": "c1", "timestamp": 200, "source": "agent-b",
    "nodes": {
      "mask_pii": {
        "before": ["{\"body\":\"card 4111111111111111\"}"],
        "after":  ["{\"body\":\"card [REDACTED]\"}"]
      }
    }
  },
  {
    "org_id": "o1", "config_id": "c1", "timestamp": 100, "source": "agent-a",
    "nodes": {
      "route_errors": {"before": ["{\"body\":\"boom\"}"]},
      "mask_pii":     {"after":  ["not json"]}
    }
  }
]`

func TestFlattenCapture(t *testing.T) {
	recs, err := flattenCapture([]byte(captureResponse))
	if err != nil {
		t.Fatalf("flattenCapture: %v", err)
	}
	if len(recs) != 4 {
		t.Fatalf("got %d records, want 4", len(recs))
	}
	// Oldest payload first, then nodes sorted by name, before before after.
	want := []struct {
		ts    int64
		src   string
		node  string
		phase string
	}{
		{100, "agent-a", "mask_pii", "after"},
		{100, "agent-a", "route_errors", "before"},
		{200, "agent-b", "mask_pii", "before"},
		{200, "agent-b", "mask_pii", "after"},
	}
	for i, w := range want {
		got := recs[i]
		if got.Timestamp != w.ts || got.Source != w.src || got.Node != w.node || got.Phase != w.phase {
			t.Errorf("record %d = {%d %s %s %s}, want {%d %s %s %s}",
				i, got.Timestamp, got.Source, got.Node, got.Phase, w.ts, w.src, w.node, w.phase)
		}
	}
	// A JSON-encoded item is decoded; anything else is passed through verbatim.
	if m, ok := recs[3].Item.(map[string]any); !ok || m["body"] != "card [REDACTED]" {
		t.Errorf("item not decoded: %#v", recs[3].Item)
	}
	if s, ok := recs[0].Item.(string); !ok || s != "not json" {
		t.Errorf("non-JSON item not passed through: %#v", recs[0].Item)
	}
}

func TestFlattenCaptureEmptyAndBadShape(t *testing.T) {
	for _, in := range []string{`[]`, `null`} {
		recs, err := flattenCapture([]byte(in))
		if err != nil {
			t.Errorf("flattenCapture(%s): %v", in, err)
		}
		if len(recs) != 0 {
			t.Errorf("flattenCapture(%s) = %d records, want 0", in, len(recs))
		}
	}
	if _, err := flattenCapture([]byte(`{"error":"nope"}`)); err == nil {
		t.Error("flattenCapture on an object: want error, got nil")
	}
}

// The dedup key must be stable across identical polls (so a re-fetched batch is
// not re-emitted) and distinct for every axis that makes an item a new event.
func TestCaptureKey(t *testing.T) {
	base := captureKey(1, "a", "n", "before", 0, "x")
	if base != captureKey(1, "a", "n", "before", 0, "x") {
		t.Error("key is not stable for identical input")
	}
	others := map[string]string{
		"timestamp": captureKey(2, "a", "n", "before", 0, "x"),
		"source":    captureKey(1, "b", "n", "before", 0, "x"),
		"node":      captureKey(1, "a", "m", "before", 0, "x"),
		"phase":     captureKey(1, "a", "n", "after", 0, "x"),
		"index":     captureKey(1, "a", "n", "before", 1, "x"),
		"item":      captureKey(1, "a", "n", "before", 0, "y"),
	}
	for axis, k := range others {
		if k == base {
			t.Errorf("key collides when %s differs", axis)
		}
	}
	// Field boundaries must not be ambiguous: ("ab","c") != ("a","bc").
	if captureKey(1, "ab", "c", "before", 0, "x") == captureKey(1, "a", "bc", "before", 0, "x") {
		t.Error("key collides across a field boundary")
	}
}

// Two identical items in the same batch are two events, not one - only a
// re-fetch of the same batch is a duplicate.
func TestFlattenCaptureDuplicateItemsInBatch(t *testing.T) {
	const dup = `[{"timestamp":1,"source":"a","nodes":{"n":{"before":["{\"body\":\"same\"}","{\"body\":\"same\"}"]}}}]`
	recs, err := flattenCapture([]byte(dup))
	if err != nil {
		t.Fatalf("flattenCapture: %v", err)
	}
	if len(recs) != 2 {
		t.Fatalf("got %d records, want 2", len(recs))
	}
	seen := newSeenSet(maxFollowSeen)
	for i, r := range recs {
		if !seen.add(r.key) {
			t.Errorf("record %d treated as a duplicate within the batch", i)
		}
	}
	for _, r := range recs {
		if seen.add(r.key) {
			t.Error("re-fetched record not deduped")
		}
	}
}

// recsOf builds the records for one poll of a single node's before items.
func recsOf(t *testing.T, ts int64, bodies ...string) []captureRecord {
	t.Helper()
	items := make([]string, len(bodies))
	for i, b := range bodies {
		items[i] = `{"body":"` + b + `"}`
	}
	blob, err := json.Marshal([]any{map[string]any{
		"timestamp": ts, "source": "agent-a",
		"nodes": map[string]any{"n": map[string]any{"before": items}},
	}})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	recs, err := flattenCapture(blob)
	if err != nil {
		t.Fatalf("flattenCapture: %v", err)
	}
	return recs
}

func TestCaptureTailEmitsNewItemsOnce(t *testing.T) {
	var out, notes bytes.Buffer
	tail := newCaptureTail(&out, &notes, false, false)

	if _, err := tail.emit(recsOf(t, 100, "a", "b")); err != nil {
		t.Fatalf("emit: %v", err)
	}
	// A re-fetched batch that grew: only the appended item is new.
	if _, err := tail.emit(recsOf(t, 100, "a", "b", "c")); err != nil {
		t.Fatalf("emit: %v", err)
	}
	lines := nonEmptyLines(out.String())
	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3:\n%s", len(lines), out.String())
	}
	for i, want := range []string{"a", "b", "c"} {
		if !strings.Contains(lines[i], `"body":"`+want+`"`) {
			t.Errorf("line %d = %s, want body %q", i, lines[i], want)
		}
	}
	if notes.Len() != 0 {
		t.Errorf("unexpected note without --since-now: %s", notes.String())
	}
}

func TestCaptureTailSinceNowSkipsBacklog(t *testing.T) {
	var out, notes bytes.Buffer
	tail := newCaptureTail(&out, &notes, false, true)

	if _, err := tail.emit(recsOf(t, 100, "old1", "old2")); err != nil {
		t.Fatalf("emit: %v", err)
	}
	if out.Len() != 0 {
		t.Fatalf("backlog was printed: %s", out.String())
	}
	if !strings.Contains(notes.String(), "skipped 2") {
		t.Errorf("note = %q, want it to report 2 skipped", notes.String())
	}
	// The backlog is still suppressed on re-fetch, and new items print.
	if _, err := tail.emit(recsOf(t, 100, "old1", "old2")); err != nil {
		t.Fatalf("emit: %v", err)
	}
	if _, err := tail.emit(recsOf(t, 300, "fresh")); err != nil {
		t.Fatalf("emit: %v", err)
	}
	lines := nonEmptyLines(out.String())
	if len(lines) != 1 || !strings.Contains(lines[0], `"body":"fresh"`) {
		t.Errorf("got %v, want only the fresh item", lines)
	}
}

// An empty first poll must still consume the skip - otherwise the next poll's
// items, which arrived after the tail started, would be swallowed as backlog.
func TestCaptureTailSinceNowEmptyFirstPoll(t *testing.T) {
	var out, notes bytes.Buffer
	tail := newCaptureTail(&out, &notes, false, true)

	if _, err := tail.emit(nil); err != nil {
		t.Fatalf("emit: %v", err)
	}
	if notes.Len() != 0 {
		t.Errorf("note printed for an empty backlog: %s", notes.String())
	}
	if _, err := tail.emit(recsOf(t, 100, "fresh")); err != nil {
		t.Fatalf("emit: %v", err)
	}
	if lines := nonEmptyLines(out.String()); len(lines) != 1 {
		t.Errorf("got %d lines, want the fresh item to print:\n%s", len(lines), out.String())
	}
}

// emit's count drives the heartbeat, so it has to match what was written.
func TestCaptureTailEmitReportsWrittenCount(t *testing.T) {
	var out, notes bytes.Buffer
	tail := newCaptureTail(&out, &notes, false, false)

	if n, err := tail.emit(recsOf(t, 100, "a", "b")); err != nil || n != 2 {
		t.Errorf("first emit = (%d, %v), want (2, nil)", n, err)
	}
	// Same batch plus one: only the appended item counts.
	if n, err := tail.emit(recsOf(t, 100, "a", "b", "c")); err != nil || n != 1 {
		t.Errorf("second emit = (%d, %v), want (1, nil)", n, err)
	}
	if n, err := tail.emit(recsOf(t, 100, "a", "b", "c")); err != nil || n != 0 {
		t.Errorf("re-fetch emit = (%d, %v), want (0, nil)", n, err)
	}
	// A skipped backlog writes nothing, so it must report nothing.
	skipping := newCaptureTail(&out, &notes, false, true)
	if n, err := skipping.emit(recsOf(t, 100, "old")); err != nil || n != 0 {
		t.Errorf("backlog emit = (%d, %v), want (0, nil)", n, err)
	}
}

func TestCaptureTailHeartbeat(t *testing.T) {
	var out, notes bytes.Buffer
	tail := newCaptureTail(&out, &notes, false, false)
	clock := tail.quietSince
	tail.now = func() time.Time { return clock }
	tail.heartbeat = time.Minute

	// Quiet, but not yet long enough to be worth mentioning.
	for i := 0; i < 5; i++ {
		clock = clock.Add(10 * time.Second)
		tail.tick(0)
	}
	if notes.Len() != 0 {
		t.Fatalf("note before the interval elapsed: %s", notes.String())
	}
	// Crossing the interval reports the silence and how much polling it covered.
	clock = clock.Add(10 * time.Second)
	tail.tick(0)
	got := notes.String()
	if !strings.Contains(got, "no new items in the last 1m0s") || !strings.Contains(got, "(6 polls)") {
		t.Errorf("note = %q, want it to report 1m0s and 6 polls", got)
	}
	// The window restarts, so the next note is another full interval away.
	notes.Reset()
	clock = clock.Add(59 * time.Second)
	tail.tick(0)
	if notes.Len() != 0 {
		t.Errorf("note repeated before a fresh interval elapsed: %s", notes.String())
	}
}

// Items flowing are their own evidence - a tail that is producing output should
// never also narrate.
func TestCaptureTailHeartbeatSilentWhileItemsFlow(t *testing.T) {
	var out, notes bytes.Buffer
	tail := newCaptureTail(&out, &notes, false, false)
	clock := tail.quietSince
	tail.now = func() time.Time { return clock }
	tail.heartbeat = time.Minute

	for i := 0; i < 10; i++ {
		clock = clock.Add(30 * time.Second)
		tail.tick(1)
	}
	if notes.Len() != 0 {
		t.Errorf("heartbeat fired while items were being written: %s", notes.String())
	}
	// One quiet poll right after a productive one is not yet a silence either.
	clock = clock.Add(30 * time.Second)
	tail.tick(0)
	if notes.Len() != 0 {
		t.Errorf("heartbeat fired only 30s after the last item: %s", notes.String())
	}
}

func TestCaptureTailRawOutput(t *testing.T) {
	var out, notes bytes.Buffer
	tail := newCaptureTail(&out, &notes, true, false)
	if _, err := tail.emit(recsOf(t, 100, "a")); err != nil {
		t.Fatalf("emit: %v", err)
	}
	got := strings.TrimSpace(out.String())
	if got != `{"body":"a"}` {
		t.Errorf("raw output = %q, want the item alone with no envelope", got)
	}
}

func nonEmptyLines(s string) []string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		if strings.TrimSpace(l) != "" {
			out = append(out, l)
		}
	}
	return out
}

func TestSeenSetEviction(t *testing.T) {
	s := newSeenSet(10)
	keys := make([]string, 12)
	for i := range keys {
		keys[i] = captureKey(int64(i), "a", "n", "before", 0, "x")
		if !s.add(keys[i]) {
			t.Fatalf("key %d reported as already seen", i)
		}
	}
	if len(s.keys) != len(s.order) {
		t.Errorf("map (%d) and order (%d) diverged", len(s.keys), len(s.order))
	}
	if len(s.order) > s.max {
		t.Errorf("order grew past max: %d > %d", len(s.order), s.max)
	}
	// The newest keys survive eviction; the oldest are dropped.
	if s.add(keys[11]) {
		t.Error("newest key was evicted")
	}
	if !s.add(keys[0]) {
		t.Error("oldest key was not evicted")
	}
}

// Each emitted record must be exactly one line, or a line-oriented consumer
// (jq, grep, a log monitor) sees one event as many.
func TestCaptureRecordEncodesToOneLine(t *testing.T) {
	recs, err := flattenCapture([]byte(`[{"timestamp":1,"source":"a","nodes":{"n":{"before":["{\"body\":\"line1\nline2\",\"html\":\"<b>&\"}"]}}}]`))
	if err != nil {
		t.Fatalf("flattenCapture: %v", err)
	}
	b, err := json.Marshal(recs[0])
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, c := range string(b) {
		if c == '\n' {
			t.Fatalf("encoded record contains a newline: %s", b)
		}
	}
	// The dedup key and raw item are internal and must not leak into output.
	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, k := range []string{"key", "raw"} {
		if _, ok := got[k]; ok {
			t.Errorf("internal field %q leaked into output: %s", k, b)
		}
	}
	for _, k := range []string{"timestamp", "source", "node", "phase", "item"} {
		if _, ok := got[k]; !ok {
			t.Errorf("missing field %q: %s", k, b)
		}
	}
}
