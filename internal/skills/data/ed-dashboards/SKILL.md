---
name: ed-dashboards
description: Dashboards - create, update, inspect and validate metric dashboards from the CLI.
metadata:
  version: "1.0.0"
  author: edgedelta
  repository: https://github.com/edgedelta/agent-skills
  tags: edgedelta,dashboards,metrics,visualization
  alwaysApply: "false"
---

# Edge Delta Dashboards

List, inspect, create, update and delete dashboards.

## Prerequisites

The `edx` CLI must be installed and authenticated. See the **ed-edx** skill.
`create`/`update`/`delete` require `edx` >= 0.10.0 (`list`/`get` work on any version).

## Inspect

```bash
edx dashboards list --output table --columns dashboard_id,dashboard_name,creator
edx dashboards get <dashboard-id>          # full definition - use as a template
```

## The One Rule That Saves Hours

A dashboard can **save fine via the API yet fail to render** in the UI
("Dashboard could not be found"). Two causes:

1. **Schema / version mismatch.** Widgets that use the `visuals[]` array require
   `definition.version: 3`. A `visuals[]` widget saved with version 1 is accepted
   by the API but the UI cannot parse it.
2. **Missing `resource_accesses`.** The UI resolves a dashboard through
   `resource_accesses`; it must contain one `{domain, query}` entry per widget
   query. An empty `resource_accesses` renders blank or errors.

`edx dashboards create`/`update` validate both **client-side before** calling the
API, so you catch these without polluting the backend.

## Anatomy of a Metric Dashboard

```json
{
  "dashboard_name": "Service Usage",
  "description": "Tokens and cost",
  "definition": {
    "version": 3,
    "timeFilters": {"lookback": "1h"},
    "widgets": [
      {"id": "root", "type": "grid", "displayOptions": {"hideBackground": true},
       "grid": "72px 72px 72px / 1fr 1fr 1fr 1fr 1fr 1fr 1fr 1fr 1fr 1fr 1fr 1fr"},
      {"id": 1, "type": "viz", "displayOptions": {"title": "Total Tokens"},
       "position": {"area": {"column": 1, "columnSpan": 3, "row": 1, "rowSpan": 3},
                    "targetId": "root", "type": "grid"},
       "resultType": "aggregate",
       "visualizer": {"type": "bignumber"},
       "visuals": [{"id": "A",
                    "dataSource": {"type": "metric",
                                   "params": {"query": "sum:service.tokens{*}"}}}]}
    ]
  },
  "resource_accesses": [{"domain": "metric", "query": "sum:service.tokens{*}"}]
}
```

Key fields:

- **root grid** is 12 columns (`1fr` x12), rows sized `72px`. Each viz widget is
  placed with `position.area.{column,columnSpan,row,rowSpan}` (1-indexed).
- **`resultType`**: `aggregate` for a single value (bignumber), `timeseries` for
  a trend (line/area/bar).
- **`visualizer.type`**: `bignumber` | `line` | `table` | `bar` | `pie` | ...
- **metric query CQL**: `<agg>:<name>{<filter>} by {<dim>}.rollup(<secs>)` -
  identical to `edx metrics query` (see the **ed-metrics** skill).
- **`resource_accesses`**: mirror every widget query, one entry each.

## Create / Update / Delete

```bash
edx dashboards create --file dashboard.json
edx dashboards update <dashboard-id> --file dashboard.json
edx dashboards delete <dashboard-id> --yes
```

Fastest authoring path: `edx dashboards get <id>` of a working dashboard, swap
the widget queries and name, then `create`.

## Gotchas

- Group-by breakdowns only work on **indexed** metric dimensions (see
  **ed-metrics**). `by {model}` on a non-indexed OTLP attribute collapses to a
  single series - break those down with logs or a dedicated `log_to_metric`.
- Keep `definition.version` (3) in sync with the `visuals[]` widget schema.
- `--skip-validation` exists as an escape hatch, but a validation failure almost
  always means the dashboard won't render.
