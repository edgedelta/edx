---
name: ed-logs
description: Log management - search logs with CQL, log volume graphs, schema discovery.
metadata:
  version: "1.0.0"
  author: edgedelta
  repository: https://github.com/edgedelta/agent-skills
  tags: edgedelta,logs,logging,search,cql
  globs: "**/edgedelta*.yaml,**/*log*"
  alwaysApply: "false"
---

# Edge Delta Logs

Search and aggregate logs with CQL (Common Query Language).

## Prerequisites

The `edx` CLI must be installed and authenticated. See the **ed-edx** skill.

## Command Execution Order (Token-Efficient)

1. Check context first (prior outputs, conversation, saved values).
2. If field names/values are unknown, discover them first:
   `edx facets keys --scope log`, then `edx facets options --scope log --facet <field>`.
3. If still ambiguous, ask the user to confirm.
4. Then run the search.

## Search Logs

```bash
# Basic error search
edx logs search -q 'severity_text:"ERROR"' --lookback 1h

# Scoped to a service, more results
edx logs search -q 'service.name:"api" AND severity_text:"ERROR"' --lookback 1h --limit 100

# Absolute time window (incident windows from PagerDuty etc.)
edx logs search -q 'error' --from 2026-06-12T00:00:00.000Z --to 2026-06-12T01:00:00.000Z

# Human-readable table
edx logs search -q 'error' --output table --columns timestamp,severity_text,service.name,body
```

### CQL Syntax

| Query | Meaning |
|-------|---------|
| `error timeout` | Full-text search (bare words) |
| `severity_text:"ERROR"` | Field equals |
| `service.name:("api" OR "web")` | Multiple values |
| `-severity_text:"DEBUG"` | Negation |
| `@response.code > 400` | Numeric attribute comparison |
| `@Record.errorCode:"AccessDenied"` | Attribute (structured field) equals - note the leading `@` |
| `a AND b`, `a OR b` | Boolean operators |

**Not supported**: regular expressions (`/pattern/`), `=`/`!=` operators,
wildcards mid-string.

Common fields: `service.name`, `severity_text`, `host.name`, `ed.tag`,
`k8s.namespace.name`, `k8s.pod.name`, `body`.

### Attributes vs the body

Full-text search (bare words) only matches `body`. Structured logs often keep
the interesting value in an **attribute**, not the body - e.g. a CloudTrail
record has `body: "AssumeRole"` (just the event name) while the failure lives in
the `Record.errorCode` attribute. Query attributes with a leading `@`
(`@Record.errorCode:"AccessDenied"`); the same field without the `@` matches
nothing. So if a full-text search returns zero, the value is probably in an
attribute - inspect one raw record to find the real field names, then filter:

```bash
edx logs search -q 'service.name:"<svc>"' --limit 1 --output raw   # read the attributes map
edx logs search -q 'service.name:"<svc>" AND @Record.errorCode:"AccessDenied"'
```

## Log Volume Graphs

Aggregate log counts over time (find spikes, compare services):

```bash
# Error volume per service
edx logs graph -q 'severity_text:"ERROR"' --group-by service.name --lookback 6h

# Total volume for one service
edx logs graph -q 'service.name:"api"' --lookback 24h
```

## Pagination

Responses include cursors. Continue a search:

```bash
edx logs search -q 'error' --limit 100 --cursor "<next_cursor from previous response>"
```

## Empty Results Checklist

1. Widen the time range: `--lookback 24h`.
2. Verify field values exist: `edx facets options --scope log --facet service.name`.
3. Remove filters one at a time.
4. Remember full-text needs bare words, field filters need exact quoted values.
5. Full-text found nothing? The value may be in an attribute, not `body`.
   Inspect a raw record (`--limit 1 --output raw`) and query it with `@field:"value"`.

## References

- CQL works the same across logs, patterns and events scopes.
- Metric and trace scopes do NOT support full-text search.
