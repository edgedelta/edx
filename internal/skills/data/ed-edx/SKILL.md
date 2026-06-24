---
name: ed-edx
description: Primary Edge Delta CLI - edx commands, authentication, output formats and conventions.
metadata:
  version: "1.0.0"
  author: edgedelta
  repository: https://github.com/edgedelta/agent-skills
  tags: edgedelta,cli,edx,auth,setup
  alwaysApply: "false"
---

# Edge Delta CLI (edx)

`edx` is the canonical way for agents to interact with Edge Delta: Pipeline
(fleet management, configs, live capture), Observability (logs, patterns,
metrics, traces, events, monitors) and AI Teammate (connectors, activity).

## Install

```bash
go install github.com/edgedelta/edx@latest
# verify
edx version
```

## Auth

```bash
edx auth login --token <api-token> --org-id <org-id>   # saved to ~/.config/edx/config.yaml
edx auth status                                        # verifies the token against the API
```

Environment variables override the config file: `ED_API_TOKEN`, `ED_ORG_ID`,
`ED_API_URL`. Multiple orgs: `edx auth login --profile <name> ...` then
`edx --profile <name> ...`.

If commands fail with 401, the token is invalid or does not match the org ID.
Re-run `edx auth login` with a token from Admin > API Tokens in the web app.

## Command Map

| Domain | Commands |
|--------|----------|
| Logs | `edx logs search`, `edx logs graph` |
| Patterns | `edx patterns list`, `edx patterns samples` |
| Metrics | `edx metrics list`, `edx metrics query` |
| Traces | `edx traces search`, `edx service-map` |
| Events | `edx events search` |
| Monitors | `edx monitors list/get/create/update/delete/states` |
| Pipelines | `edx pipelines list/get/history/save/deploy/validate/agents/status` |
| Fleet | `edx fleet agents`, `edx fleet deployments` |
| Live capture | `edx capture start/task/status/results` |
| Health | `edx health components`, `edx health problems` |
| Dashboards | `edx dashboards list/get` |
| Schema | `edx facets keys/options/list` |
| AI Teammate | `edx ai connectors ...`, `edx ai activity` |
| Ingestion | `edx ingest endpoints`, `edx ingest token` |
| Raw API | `edx api <METHOD> <path>` |

## Conventions

- **Output**: pretty JSON by default - parse it directly or pipe to `jq`.
  `--output table --columns a,b,c` for human summaries.
- **Time ranges**: `--lookback 15m|1h|24h` (Go durations) or
  `--from/--to` ISO 8601 (`2006-01-02T15:04:05.000Z`).
- **Pagination**: search responses include cursors; pass `--cursor` to continue.
- **Confirmations**: destructive commands (`deploy`, `delete`) prompt; add
  `--yes` in non-interactive contexts.
- **Limits**: default search limit is 20; raise with `--limit` (max 1000).

## Escape Hatch

For endpoints without a dedicated command:

```bash
edx api GET /v1/orgs/{org}/users          # {org} auto-substituted
edx api GET /tokens                       # org-relative shorthand
edx api POST /monitors --data @monitor.json
edx api GET /rehydration --param limit=10
```

## Troubleshooting

| Problem | Fix |
|---------|-----|
| 401 Invalid authentication token | `edx auth login` with a valid token + matching org ID |
| 404 on a path | Check the org ID; try `edx api` with the full /v1 path |
| Empty search results | Widen `--lookback`; verify fields with `edx facets options` |
| Timeout on large queries | Narrow the time range or add `--timeout 120s` |
