---
name: ed-events
description: Events - pattern anomalies, monitor alert triggers and Kubernetes events.
metadata:
  version: "1.0.0"
  author: edgedelta
  repository: https://github.com/edgedelta/agent-skills
  tags: edgedelta,events,anomalies,alerts,kubernetes
  alwaysApply: "false"
---

# Edge Delta Events

Events are the "what happened" stream: anomaly detections, monitor alert
triggers and Kubernetes events. Always check events early in an incident -
they often point straight at the cause.

## Prerequisites

The `edx` CLI must be installed and authenticated. See the **ed-edx** skill.

## Event Types

| Query | Meaning |
|-------|---------|
| `event.type:"pattern_anomaly"` | Log anomaly detections |
| `event.type:"metric_threshold"` | Metric alert triggers |
| `event.type:"log_threshold"` | Log alert triggers |
| `event.domain:"Monitor Alerts"` | All monitor-triggered events |
| `event.domain:"K8s"` | Kubernetes events (OOMKilled, BackOff, ...) |

Discover the live set: `edx facets options --scope event --facet event.type`
and `--facet event.domain`.

## Search Events

```bash
# All anomalies in the last 6 hours
edx events search -q 'event.type:"pattern_anomaly"' --lookback 6h

# Anomalies for one service
edx events search -q 'service.name:"api" AND event.type:"pattern_anomaly"'

# Everything monitors fired recently, as a table
edx events search -q 'event.domain:"Monitor Alerts"' --output table

# Kubernetes trouble
edx events search -q 'event.domain:"K8s" AND OOMKilled' --lookback 24h
```

Full-text search is supported in the event scope (bare words work).

## Incident Usage

1. Establish the incident window (from the page/alert).
2. `edx events search --from <start> --to <end>` - what fired in the window?
3. Pattern anomalies name the service and signature - pivot to
   `edx patterns list` / `edx logs search` for detail.
4. Monitor alerts carry the monitor ID - `edx monitors get <id>` for the
   query and thresholds behind it.
