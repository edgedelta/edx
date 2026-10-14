# Edge Delta Agent Skills

Edge Delta skills for AI agents (Claude Code, Cursor, Copilot and other
SKILL.md-compatible tools), built on the [`edx` CLI](https://github.com/edgedelta/edx)
and Edge Delta's SDKs.

> The skills are authored in [`edgedelta/edx/skills`](https://github.com/edgedelta/edx/tree/main/skills)
> and published to [`edgedelta/agent-skills`](https://github.com/edgedelta/agent-skills)
> on every edx release. Open issues and pull requests against edx.

## See it in action

Ask your agent a root-cause question and the skills drive the investigation.
Here is what it looks like in Claude Code — surging error patterns found,
onset quantified, verdict delivered (**ed-investigate**):

![root-cause investigation in Claude Code](assets/skill-rca.gif)

And tracing a fleet regression back to the offending pipeline config change,
with the rollback command ready (**ed-pipelines**):

![config-change root cause in Claude Code](assets/skill-config-rca.gif)

The `edx` outputs in these clips are from real investigations against a live
org, replayed with scripted timing to keep the clips short (see
[`demo/`](demo/)). Re-record with `vhs demo/<name>.tape` — no credentials
appear on screen.

## Skills

| Skill | Description |
|-------|-------------|
| [ed-onboard](ed-onboard/SKILL.md) | Connect cloud, on-premises and hybrid telemetry sources to Edge Delta |
| [ed-edx](ed-edx/SKILL.md) | Primary CLI - all edx commands, auth, setup |
| [ed-logs](ed-logs/SKILL.md) | Search logs with CQL, log volume graphs |
| [ed-patterns](ed-patterns/SKILL.md) | Log patterns, anomaly and sentiment analysis |
| [ed-metrics](ed-metrics/SKILL.md) | Discover and aggregate metrics |
| [ed-traces](ed-traces/SKILL.md) | Distributed traces and the service map |
| [ed-events](ed-events/SKILL.md) | Events: anomalies, monitor alerts, K8s events |
| [ed-monitors](ed-monitors/SKILL.md) | Create, manage and resolve monitors |
| [ed-pipelines](ed-pipelines/SKILL.md) | Fleet management, config changes, deployments, live capture |
| [ed-rehydrations](ed-rehydrations/SKILL.md) | Replay archived data back through a pipeline |
| [ed-pipeline-tuning](ed-pipeline-tuning/SKILL.md) | Improve data quality: parse, structure, enrich, mask, roll up to metrics/patterns |
| [ed-dashboards](ed-dashboards/SKILL.md) | Create, update, validate and screenshot metric dashboards |
| [ed-investigate](ed-investigate/SKILL.md) | Cross-signal incident investigation workflow |
| [ed-ai-teammate](ed-ai-teammate/SKILL.md) | AI Teammate connectors and activity |
| [ed-browser-rum](ed-browser-rum/SKILL.md) | Browser RUM: page-load performance, Core Web Vitals and JS errors from real visitors |

## Setup edx

```bash
brew install edgedelta/tap/edx          # or: go install github.com/edgedelta/edx@latest

# Token auth (good for CI):
edx auth login --token <api-token> --org-id <org-id>
# Or OAuth (browser login; org is read from the token):
edx auth login --oauth

edx auth status
```

API tokens are created in the Edge Delta web app under **Admin → API Tokens**.
Environment variables `ED_API_TOKEN`, `ED_ORG_ID` and `ED_ENV` (`prod`,
`staging` or `local`) override the config file - useful in CI.

## Install Skills

```bash
npx skills add edgedelta/agent-skills \
  --skill ed-edx \
  --skill ed-onboard \
  --skill ed-logs \
  --skill ed-patterns \
  --skill ed-metrics \
  --skill ed-traces \
  --skill ed-monitors \
  --skill ed-pipelines \
  --skill ed-rehydrations \
  --skill ed-pipeline-tuning \
  --skill ed-dashboards \
  --skill ed-events \
  --skill ed-investigate \
  --skill ed-ai-teammate \
  --skill ed-browser-rum \
  --full-depth -y
```

Or copy the skill directories you need into your agent's skills folder
(e.g. `.claude/skills/`).

## Conventions Used Across Skills

- Discovery before queries: `edx facets keys` / `edx facets options` to find
  field names and values before writing CQL.
- Time ranges: `--lookback 1h` (Go durations) or `--from/--to` (ISO 8601).
- Default output is JSON; `--output table --columns ...` for summaries.
- Destructive operations (`deploy`, `delete`) require `--yes` when
  non-interactive.
