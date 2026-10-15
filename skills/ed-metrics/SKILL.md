---
name: ed-metrics
description: Metrics - discover and describe metrics, summarize series (stats, trend, anomalies, week-over-week baseline) and run aggregation queries.
metadata:
  version: "1.0.0"
  author: edgedelta
  repository: https://github.com/edgedelta/agent-skills
  tags: edgedelta,metrics,timeseries,aggregation
  alwaysApply: "false"
---

# Edge Delta Metrics

Discover, describe and aggregate metrics.

## Prerequisites

The `edx` CLI must be installed and authenticated. See the **ed-edx** skill.

## Agent Workflow (token-efficient)

1. **Find the exact name**: `edx metrics list --keyword <part>`. It only lists metrics that reported in the window (default last 24h); for an older incident pass `--from/--to` around it.
2. **Describe it first**: `edx metrics describe <name>` - one call returns kind,
   unit, temporality, labels with cardinality and sample values, staleness and
   ready-to-run example queries. Copy an example rather than guessing the
   aggregation or label syntax.
3. **Investigate with `--summary`**: `edx metrics query ... --summary` returns
   per-series stats, trend and anomaly windows instead of every point (a 10-series
   summary is a few KB; raw points can be 100x that). Add `--baseline 1w` for
   "is this abnormal?" questions.
4. **Raw points only when needed** (plotting, exact timestamps): plain
   `edx metrics query` without `--summary`; for big responses use `--output-file`.

## Rules That Prevent Failed Queries

1. **Metric names must be EXACT** - no wildcards, no regex, case-sensitive.
   Always discover the name first.
2. **No full-text search** in metric filters - `field:"value"` syntax only
   (regex: `field:~"re"`, negated `field!~"re"`).
3. Group-by keys must be real labels - use the `ref` values from
   `edx metrics describe <name>` (datapoint attributes take an `@` prefix,
   e.g. `@http.route`).
4. **Pick the aggregation by kind** (describe tells you the kind):
   - cumulative counter (monotonic sum, `_total`): `--agg rate` or `--agg increase`;
     sum/avg of a cumulative counter only shows the running total.
   - gauge: `--agg avg` / `--agg max`.
   - histogram: `--agg p95` (also p25/p50/p75/p90/p99), or `--agg avg`.
   - delta sum / count-style: `--agg sum`.

## Discover and Describe

```bash
edx metrics list                    # all metric names
edx metrics list --keyword cpu      # fuzzy filter
edx metrics list --keyword http --lookback 7d   # widen the window for older incidents
edx metrics describe http.server.request.duration
edx metrics describe ed.pipeline.bytes --lookback 24h
```

`describe` prints `found`, `kind`, `unit`, `temporality`, `monotonic`,
`prom_suffix`, `series_count`, `latest_sample` and `age` (staleness), `labels`
(`name`, `kind` resource/attribute, `ref` to use in queries, `values_count`,
`samples`) and `examples`. If one backend call fails it still prints what it
got plus `warnings`. An unknown name exits non-zero with `suggestions` (closest
metric names).

## Summarize a Metric (preferred for investigation)

```bash
# Per-service request rate: stats, trend, anomaly windows, top 10 series
edx metrics query --name http.server.requests --agg rate --group-by service.name --summary

# Compare with the same hour last week (baseline discipline)
edx metrics query --name http.server.requests --agg sum --group-by service.name \
  --summary --baseline 1w --top 5

# p95 latency for one route over 6h, stricter anomaly threshold
edx metrics query --name http.server.request.duration --agg p95 \
  --filter '@http.route:"/checkout"' --lookback 6h --summary --anomaly-k 5
```

Summary header: `query` (the CQL), `from`/`to`, `rollup`, `value_stat`
(`total` for sum/count/increase, else `avg`), `ranked_by` (`share` of total for
sum-like aggs, else `max`), `series_count`, `shown`, `omitted`, `total`,
`top_contributors`, and `baseline` window bounds when requested.

Per series: `series` (labels), `value`, `share_pct`, `min`, `max`, `avg`, `p95`,
`first`, `last`, `change`, `change_pct`, `trend` (rising/falling/flat - slope
significance, not eyeballing), `points`, `gaps` (null/missing buckets) and
`anomalies` (windows where the modified z-score vs the series median/MAD
exceeds `--anomaly-k`, default 3.5; `dir` high/low, `peak`).

With `--baseline <offset>` each series gets `baseline.value`, `delta`,
`delta_pct` and `verdict`:

| verdict | meaning |
|---------|---------|
| `higher` / `lower` | change of 10% or more, with enough data in both windows |
| `similar` | within 10% |
| `insufficient_data` | a window has fewer than 3 points (or the baseline has under half the current points; count aggs under 10) - do not read a ratio |
| `missing` | series did not exist in the baseline window |

Baseline rules (same as ed-baseline-compare): prefer `--baseline 1w` (survives
daily and weekly cycles), then `1d`; always quote absolute numbers with the
ratio ("2.2x, 72k -> 158k req"); state both windows' exact bounds (`from`/`to`
and `baseline.from`/`baseline.to`).

`--top N` (default 10, 0 = all) keeps the highest-ranked series and reports
how many were `omitted`. `-o table` prints a readable table; `-o csv` one row
per series.

## Raw Query

```bash
# Average request duration per service, last hour
edx metrics query --name http.request.duration --agg avg --group-by service.name

# Max CPU on one host over 24h with 5-minute rollups
edx metrics query --name system.cpu.usage --agg max \
  --filter 'host.name:"web-1"' --rollup 300 --lookback 24h

# Regex filter across a service family
edx metrics query --name http.server.requests --agg sum \
  --filter 'service.name:~"checkout-.*"' --group-by service.name

# Table instead of timeseries (current values, good for top-N)
edx metrics query --name http.requests --agg sum \
  --group-by service.name --graph-type table --output table
```

Aggregations: `sum`, `avg`, `min`, `max`, `count`, `median`, `first`,
`p25`/`p50`/`p75`/`p90`/`p95`/`p99`, `rate`, `increase`.
Filter: CQL field syntax (`field:"value"`, regex `field:~"re"`), `"*"` for none.

Under the hood this builds the CQL `agg:name{filter} by {keys}.rollup(secs)`
and returns records keyed by formula, e.g. `{"A": {"records": [...]}}`.
`--output table`/`csv` flattens it to one row per series: label columns,
`aggregate`, `points`, `last`.

## Groupable Dimensions (Important)

Resource dimensions are a **fixed indexed allowlist** (`service.name`, `host.ip`,
`ed.*`, `k8s.*`, ...) - discover them with `edx facets keys --scope metric` or
`edx metrics describe <name>`. Datapoint attributes are grouped with an `@`
prefix (`--group-by @http.route`). Grouping by anything else (e.g. a custom
`log_to_metric` field dimension that is not indexed) **silently returns one
empty group** rather than an error; `edx metrics query` warns when a
non-`@` `--group-by` key is not indexed.

## Common Investigations

| Question | Command |
|----------|---------|
| What is this metric and how do I query it? | `edx metrics describe <name>` |
| Is CPU/memory elevated vs normal? | `edx metrics query --name system.cpu.usage --agg avg --group-by host.name --summary --baseline 1w` |
| Which service has the most errors? | `edx metrics query --name <error-count-metric> --agg sum --group-by service.name --summary --top 5` |
| When did latency regress? | `edx metrics query --name http.request.duration --agg p95 --lookback 6h --summary` (read `anomalies[].from`) |
| Exact points for a chart | `edx metrics query --name <name> --agg avg --output-file points.json` |

## Empty Results Checklist

1. Verify the metric name exactly: `edx metrics describe <name>` (suggests close names).
2. Check `age` / `latest_sample` in describe - the metric may have stopped reporting.
3. Verify filter values: the `samples` in describe, or `edx facets options --scope metric --facet service.name`.
4. Widen `--lookback`.
5. Remove `--filter` (use `"*"`), then add filters back one at a time.
