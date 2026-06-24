---
name: ed-monitors
description: Monitors - create, manage, snooze and resolve Edge Delta monitors and alerts.
metadata:
  version: "1.0.0"
  author: edgedelta
  repository: https://github.com/edgedelta/agent-skills
  tags: edgedelta,monitors,alerts,alerting
  alwaysApply: "false"
---

# Edge Delta Monitors

Create, inspect, update and delete monitors; view triggered/resolved states.

## Prerequisites

The `edx` CLI must be installed and authenticated. See the **ed-edx** skill.

## Inspect

```bash
edx monitors list --output table
edx monitors get <monitor-id>          # full definition: query, thresholds, notifications
edx monitors states                    # current triggered/resolved states
edx monitors states -q 'monitor.status:"alert"'
```

## Create / Update

Monitor definitions are JSON. The reliable workflow is **clone-and-edit**:

```bash
# 1. Fetch an existing monitor of the same type as a template
edx monitors get <id> > monitor.json

# 2. Edit name, query, thresholds, notification targets

# 3. Create (or update in place)
edx monitors create --file monitor.json
edx monitors update <id> --file monitor.json
```

Monitor queries use the same CQL as search commands - validate the query
first by running it:

```bash
edx logs graph -q 'severity_text:"ERROR" AND service.name:"api"' --lookback 1h
```

If the graph returns sensible numbers, the monitor query will too.

## Delete

```bash
edx monitors delete <monitor-id> --yes
```

## Alert Triage

When an alert fires:

1. `edx monitors get <id>` - what condition fired? what query?
2. `edx events search -q 'event.domain:"Monitor Alerts"' --lookback 2h` -
   correlated alerts around the same time?
3. Run the monitor's underlying query yourself with a wider window to see
   the trend (`edx logs graph` / `edx metrics query`).
4. Pivot to the **ed-investigate** skill for full root-cause analysis.

## Good Practices

- Alert on symptoms (error rate, latency) not causes (CPU).
- Scope queries with `service.name`/`ed.tag` to keep alerts actionable.
- Bounded group-by keys only - high-cardinality grouping creates alert storms.
