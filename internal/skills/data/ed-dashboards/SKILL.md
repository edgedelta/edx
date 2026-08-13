---
name: ed-dashboards
description: Dashboards - create, update, inspect and validate metric dashboards from the CLI.
metadata:
  version: "1.1.0"
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
Offline validation needs a build that ships the `cql` command; check with
`edx cql validate --help` and upgrade via `edx update` if it is missing.

## Inspect

```bash
edx dashboards list --output table --columns dashboard_id,dashboard_name,creator
edx dashboards get <dashboard-id>          # full definition - use as a template
```

## Build by Validating, Not by Guessing

**Do not write a dashboard JSON and hope it renders.** Validation is offline,
instant and needs no API call, so run it after every edit. A dashboard can
**save fine via the API yet fail to render** in the UI ("Dashboard could not be
found"), and validation is what catches that before you pollute the backend.

The loop:

```bash
# 1. Check the query on its own, before it goes anywhere near a dashboard.
edx cql validate --type metric 'sum:service.tokens{*} by {model}.rollup(60)'

# 2. Draft the definition, then validate after every edit. Repeat until clean.
edx dashboards validate --file dashboard.json

# 3. Only once it validates, create it.
edx dashboards create --file dashboard.json
```

Iterate on step 2 until it prints `definition is valid`. Errors are JSON
Pointers into your file, so fix exactly what is named and re-run:

```
$ edx dashboards validate --file dashboard.json
! /widgets/1/displayOptions: additional properties 'titel' not allowed
! /widgets/1/visualizer/type: value must be 'raw-table'
! /widgets/1/visualizer/type: value must be one of 'empty', 'json', 'table', 'list',
  'geomap', 'bignumber', 'gauge', 'pie', 'donut', 'column', 'radar', 'sunburst', ...
Error: dashboard definition failed validation (3 issue(s)); fix, or pass --skip-validation to override
```

Fix one reported path at a time and re-run rather than rewriting the file — the
error lists the accepted values, so you rarely have to guess.

One path can be reported more than once, as `visualizer/type` is above. Widget
shapes are variants, and some (like `raw-table`) accept a different set of
sibling fields, so each variant reports its own accepted values. Treat the
**union** of those messages as the valid set.

`validate` takes either a full dashboard body or a bare `definition` object, so
you can validate a draft before you have wrapped it in one.

`create` and `update` run the same checks automatically and refuse to submit a
bad definition. Validating first is still worth it: you iterate without network
round-trips.

### What validation does and does not catch

| Checked offline | Not checked |
| --- | --- |
| Definition structure, against the schema generated from the UI's own types | Whether a facet or metric name **exists** (needs the backend) |
| Unknown/misspelled keys, wrong enum values, wrong types | Whether a query returns any **data** |
| Query **syntax**, per data type, with the backend's own grammars | Whether a group-by dimension is **indexed** (see Gotchas) |

So a definition that validates is well-formed, not necessarily populated. After
creating it, confirm the queries actually return data with `edx metrics query`
or `edx logs search`.

Two things are reported as **warnings** and do **not** fail the command, so read
the output rather than only the exit code:

- `resource_accesses is empty` - see below, this one breaks rendering.
- `unsupported definition version` - see below, this means you got **no**
  schema validation at all.

## Use `version: 4`

Author new dashboards with `definition.version: 4`. It is the current schema and
**the only version `edx` can validate**. Older versions still render (the UI
migrates them) but validation skips them entirely:

```
$ edx dashboards validate --file old.json
! unsupported definition version: got 3, edx can only validate version 4; skipping schema validation
definition is valid          # <- this says nothing was checked, not that it is correct
```

That output is a trap: exit code 0 with no schema check performed. If you see
it, bump the definition to version 4.

When you use `edx dashboards get <id>` as a template, check its `version` — an
older dashboard hands you an older schema.

## Missing `resource_accesses` Breaks Rendering

The UI resolves a dashboard through `resource_accesses`; it needs one
`{domain, query}` entry per widget query. An empty `resource_accesses` renders
blank or errors, and this is a **warning**, not an error — the command still
exits 0. Mirror every widget query there.

## Anatomy of a Metric Dashboard

```json
{
  "dashboard_name": "Service Usage",
  "description": "Tokens and cost",
  "definition": {
    "version": 4,
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
  placed with `position.area.{column,columnSpan,row,rowSpan}` (1-indexed). Add
  more `72px` row tracks as the dashboard grows.
- **`resultType`**: `aggregate` for a single value (bignumber), `timeseries` for
  a trend (line/area/bar), `raw` for a raw-table, `empty` for none.
- **`visualizer.type`**: `bignumber` | `line` | `area` | `bar` | `pie` | `table`
  | `gauge` | ... — validation lists every accepted value if you get it wrong.
- **`visuals[].id`**: `A`-`F` for queries, `W`-`Z` for formulas.
- **`resource_accesses`**: mirror every widget query, one entry each.

## Validating Queries On Their Own

Query syntax differs by data type and the grammars are **not** interchangeable —
a metric query is a syntax error in a log widget and vice versa. Check a query
before embedding it:

```bash
edx cql validate --type metric 'sum:ed.host.cpu{*} by {host.name}.rollup(60)'
edx cql validate --type log '{error AND ed.tag:$fleet} by {host.name}'
edx cql validate --type formula 'timeshift(q1, 3600) / q2'
```

`--type` takes the data type the widget targets: `log`, `event`, `pattern`,
`trace`, `metric` or `formula`. Errors point at a column with a caret:

```
$ edx cql validate --type metric 'sum:foo{*}.rollup(abc)'
! invalid metric query at column 18: mismatched input 'abc' expecting NUMBER
    sum:foo{*}.rollup(abc)
                      ^
```

To check every query in an existing dashboard at once:

```bash
edx dashboards get <id> \
  | jq -r '.definition.widgets[]?.visuals[]?.dataSource | select(.type=="metric") | .params.query // empty' \
  | edx cql validate --type metric --file -
```

`$variable` references are valid in a dashboard query and parse unsubstituted,
so you can validate a template as-is.

## Create / Update / Delete

```bash
edx dashboards create --file dashboard.json
edx dashboards update <dashboard-id> --file dashboard.json
edx dashboards delete <dashboard-id> --yes
```

Fastest authoring path: `edx dashboards get <id>` of a working dashboard, swap
the widget queries and name, bump `version` to 4 if it is older, then validate
and `create`.

## Gotchas

- Group-by breakdowns only work on **indexed** metric dimensions (see
  **ed-metrics**). `by {model}` on a non-indexed OTLP attribute collapses to a
  single series - break those down with logs or a dedicated `log_to_metric`.
  Validation cannot detect this; it is a data-modelling issue, not a syntax one.
- `--skip-validation` exists as an escape hatch, but a validation failure almost
  always means the dashboard won't render. Fix the definition instead; reach for
  the flag only when you have confirmed the validator is wrong.
- `position.targetId` is a widget reference, and validation checks its type but
  not that the target exists. Every shipped dashboard points grid-positioned
  widgets at the root grid widget (`"root"`), so follow that and confirm the
  target id is spelled the same in both places yourself.
