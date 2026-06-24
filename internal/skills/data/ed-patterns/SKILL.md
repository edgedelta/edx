---
name: ed-patterns
description: Log patterns - clustered message signatures with counts, deltas and sentiment for anomaly hunting.
metadata:
  version: "1.0.0"
  author: edgedelta
  repository: https://github.com/edgedelta/agent-skills
  tags: edgedelta,patterns,clustering,anomaly,logs
  alwaysApply: "false"
---

# Edge Delta Log Patterns

Edge Delta clusters similar log messages into **patterns**. Each pattern
carries `count`, `proportion`, `sentiment` (positive/negative/neutral) and
`delta` (change vs an earlier window). Patterns are the fastest way to answer
"what's new or surging in the logs?" without reading raw lines.

## Prerequisites

The `edx` CLI must be installed and authenticated. See the **ed-edx** skill.

## Quick Start

```bash
# ~50 interesting clusters: top anomalies, biggest delta up/down, top counts
edx patterns list --summary --lookback 1h

# Negative-sentiment patterns only (errors, failures, timeouts)
edx patterns list --negative --lookback 1h

# Scoped to one service
edx patterns list -q 'service.name:"api"' --negative --limit 20
```

## Comparing Windows (Delta)

The `delta` stat compares against an earlier window. Default offset equals the
lookback. To compare against the same window yesterday:

```bash
edx patterns list --lookback 1h --offset 24h
```

A large positive delta on a negative-sentiment pattern is a strong incident
signal: a new or surging error signature.

## Getting Raw Samples

Once a suspicious pattern is found, fetch raw log lines behind it:

```bash
edx patterns samples -q 'service.name:"api"' --param pattern='<pattern text>'
```

Or search logs directly using distinctive tokens from the pattern:

```bash
edx logs search -q '"connection refused"' --lookback 1h
```

## Interpretation Guide

| Signal | Meaning |
|--------|---------|
| New pattern, high count | New failure mode or new deploy behavior |
| delta >> 0, negative sentiment | Surging error - investigate first |
| delta << 0 on normal traffic patterns | Possible traffic drop / outage |
| High proportion shift | Behavior change even if totals look flat |

## Empty Results Checklist

1. Widen `--lookback`.
2. Verify filters: `edx facets options --scope pattern --facet service.name`.
3. Drop `--negative` (the issue may be neutral-sentiment).
