package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/edgedelta/edx/internal/api"

	"github.com/spf13/cobra"
)

// Cursor pagination is one convention across the API's list/search endpoints:
// an array under a per-endpoint key plus a cursor that is non-empty while more
// results exist. pagedRequest is the shared tail for every such command: a
// single page by default (with a more-results hint on stderr), or a complete
// sweep with --all.

// registerAllFlag adds the --all sweep flag to a cursor-paginated command.
func registerAllFlag(cmd *cobra.Command, all *bool) {
	cmd.Flags().BoolVar(all, "all", false, "fetch every page (follow the response cursor until empty) and print one combined response")
}

// pagedRequest is one cursor-paginated list/search invocation.
type pagedRequest struct {
	svc      api.Service
	path     string
	itemsKey string // the response's array field: items, states, monitors, rehydrations, data, matches, ...
	all      bool
	cursor   string            // starting cursor from the --cursor flag
	allLimit func()            // optional: raise the page size for --all when the user did not set --limit
	query    func() url.Values // command-specific params; may also set cursor (overridden while sweeping)
}

// run executes the request: a single page by default (with a more-results hint
// on stderr when the response carries a cursor), or a complete sweep with
// --all.
func (r pagedRequest) run(cmd *cobra.Command) error {
	c, err := newClient()
	if err != nil {
		return err
	}
	if r.all && r.allLimit != nil && !cmd.Flags().Changed("limit") {
		r.allLimit()
	}
	buildQuery := func(cursor string) url.Values {
		q := r.query()
		if cursor != "" {
			q.Set("cursor", cursor)
		}
		return q
	}
	if !r.all {
		data, err := c.GetFrom(cmdContext(cmd), r.svc, r.path, buildQuery(""))
		if err != nil {
			return err
		}
		// Legibility for one-page reads: say on stderr when the count is a
		// lower bound, since table/csv output hides the cursor.
		if _, next, err := parsePage(data, r.itemsKey); err == nil && next != "" {
			warnf("more results exist: re-run with --cursor '%s', or use --all to fetch every page", next)
		}
		return printResult(data)
	}
	merged, err := r.fetchAll(cmdContext(cmd), c, buildQuery)
	if err != nil {
		return err
	}
	return printResult(merged)
}

// parsePage pulls the array and the next cursor out of a paginated envelope.
// It understands the API's three envelope families:
//
//	{<itemsKey>: [...], "next_cursor": "..."}           observability endpoints
//	{"data": [...], "nextCursor": "...", ...}           AI chat/workflow lists
//	{"data": {<itemsKey>: [...], "nextCursor": "..."}}  knowledge graph search
func parsePage(data []byte, itemsKey string) ([]json.RawMessage, string, error) {
	var env map[string]json.RawMessage
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, "", err
	}
	// The AI services wrap the payload in "data"; descend when the array is
	// not at the top level.
	if _, ok := env[itemsKey]; !ok || itemsKey == "data" {
		if raw, ok := env["data"]; ok && len(raw) > 0 && raw[0] == '{' {
			var inner map[string]json.RawMessage
			if err := json.Unmarshal(raw, &inner); err != nil {
				return nil, "", err
			}
			env = inner
		}
	}
	next, err := stringField(env, "next_cursor")
	if err != nil {
		return nil, "", err
	}
	if next == "" {
		if next, err = stringField(env, "nextCursor"); err != nil {
			return nil, "", err
		}
	}
	var items []json.RawMessage
	if raw, ok := env[itemsKey]; ok {
		if err := json.Unmarshal(raw, &items); err != nil {
			return nil, "", err
		}
	}
	return items, next, nil
}

// stringField reads an optional string (or null) field from a decoded envelope.
func stringField(env map[string]json.RawMessage, key string) (string, error) {
	raw, ok := env[key]
	if !ok || string(raw) == "null" {
		return "", nil
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return "", err
	}
	return s, nil
}

// allPageAttempts is how many times --all tries one page before giving up.
// The client already retries transport errors and 429/gateway codes; this
// covers plain 500s, which search endpoints return transiently under load and
// which would otherwise throw away every page fetched so far.
const allPageAttempts = 4

// allRetryBaseDelay is the first retry backoff for a failed --all page
// (doubles per attempt). A variable so tests can shrink it.
var allRetryBaseDelay = time.Second

// fetchAll follows the response cursor until the server reports the result
// set complete, returning a single combined envelope. It fails atomically: any
// mid-sweep error aborts with no output, so a partial sweep can never be
// mistaken for a complete one; the error names the cursor to resume from.
func (r pagedRequest) fetchAll(ctx context.Context, c *api.Client, buildQuery func(cursor string) url.Values) ([]byte, error) {
	items := []json.RawMessage{}
	pages := 0
	cursor := r.cursor
	for {
		data, err := r.fetchPage(ctx, c, buildQuery(cursor), pages+1)
		if err != nil {
			if cursor != "" {
				return nil, fmt.Errorf("%w (items fetched so far were discarded; resume from this point with --all --cursor '%s')", err, cursor)
			}
			return nil, err
		}
		pageItems, next, err := parsePage(data, r.itemsKey)
		if err != nil {
			return nil, fmt.Errorf("page %d: unexpected response shape: %w", pages+1, err)
		}
		items = append(items, pageItems...)
		pages++
		warnf("page %d: %d item(s), %d total", pages, len(pageItems), len(items))
		if next == "" {
			break
		}
		if next == cursor {
			return nil, fmt.Errorf("page %d: server returned the same cursor twice; aborting to avoid an infinite loop", pages)
		}
		cursor = next
	}
	return json.Marshal(map[string]any{
		r.itemsKey:    items,
		"pages":       pages,
		"total_items": len(items),
		"next_cursor": "",
	})
}

// fetchPage gets one page, retrying transient failures with backoff. The
// request is an idempotent GET, so retrying is always safe.
func (r pagedRequest) fetchPage(ctx context.Context, c *api.Client, q url.Values, page int) ([]byte, error) {
	var lastErr error
	for attempt := 1; attempt <= allPageAttempts; attempt++ {
		if attempt > 1 {
			delay := allRetryBaseDelay << uint(attempt-2)
			warnf("page %d attempt %d/%d failed (%v); retrying in %s", page, attempt-1, allPageAttempts, lastErr, delay)
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
		data, err := c.GetFrom(ctx, r.svc, r.path, q)
		if err == nil {
			return data, nil
		}
		lastErr = err
		// A 4xx is deterministic (bad cursor, bad query): retrying cannot help.
		var apiErr *api.Error
		if errors.As(err, &apiErr) && apiErr.Status >= 400 && apiErr.Status < 500 {
			return nil, fmt.Errorf("page %d: %w", page, err)
		}
	}
	return nil, fmt.Errorf("page %d failed after %d attempts: %w", page, allPageAttempts, lastErr)
}
