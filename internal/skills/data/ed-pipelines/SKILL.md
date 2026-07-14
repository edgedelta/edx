---
name: ed-pipelines
description: Pipelines - fleet management, config changes, version history, deployments and live capture.
metadata:
  version: "1.0.0"
  author: edgedelta
  repository: https://github.com/edgedelta/agent-skills
  tags: edgedelta,pipelines,fleet,agents,deploy,live-capture
  globs: "**/edgedelta*.yaml,**/pipeline*.yaml"
  alwaysApply: "false"
---

# Edge Delta Pipelines

Manage fleets of Edge Delta agents: pipeline configurations, version history,
deployments, agent health and live capture.

## Prerequisites

The `edx` CLI must be installed and authenticated. See the **ed-edx** skill.

## Concepts

- A **pipeline** (conf) is an agent configuration with server-side version
  history. A **fleet** is the set of agents running that pipeline.
- Changing a pipeline is two steps: `save` creates a new version, `deploy`
  rolls a version out to the fleet. Deploying an older version is the
  supported rollback.

## Fleet Visibility

```bash
edx pipelines list --output table --columns id,tag,fleet_type,environment,status,updated
edx pipelines agents <pipeline-id>        # agents running this pipeline
edx pipelines status <pipeline-id>        # running / suspended
edx fleet agents                      # all agents org-wide
edx fleet deployments                 # rollout status across pipelines
edx health problems                   # components currently failing
```

## Config Change Workflow

```bash
# 1. Fetch the current config (content field holds the YAML)
edx pipelines get <pipeline-id> > pipeline.json
jq -r .content pipeline.json > pipeline.yaml

# 2. Edit pipeline.yaml (to develop the transform itself, see the ed-pipeline-tuning skill)

# 3. Dry-run the change on sample logs before saving (offline, no deploy)
edx pipelines test ottl <pipeline-id> --file samples.jsonl --statements '<ottl>'

# 4. Validate before saving
edx pipelines validate --file pipeline.yaml

# 5. Save a new version with a meaningful description
#    (save prints the new version - the epoch-ms timestamp deploy expects)
edx pipelines save <pipeline-id> --file pipeline.yaml -d "mask PII in checkout logs"

# 6. Deploy. --latest deploys the newest saved version; --wait blocks until
#    agents check in after the rollout and fails if one never does.
edx pipelines deploy <pipeline-id> --latest --wait --yes

# 7. Watch the rollout
edx fleet deployments <pipeline-id>
```

The deploy **version** is a saved version's epoch-millisecond timestamp, shown as
the `version` column of `edx pipelines history`. Use `--latest` to skip the
lookup, or pass an explicit version to roll forward/back. (`--latest`/`--wait`
require `edx` >= 0.10.0.)

## OTLP and Other Push Sources

An OTLP source ingests over the network. One `otlp_input` gRPC node on a port
receives **all signals** (logs, metrics and traces) - a single node multiplexes
them, and the per-signal `data_type` field is deprecated:

```yaml
- name: my_otlp
  type: otlp_input
  port: 4317
  protocol: grpc        # or http
```

Point the sender (e.g. an OpenTelemetry SDK/collector, or Claude Code with
`CLAUDE_CODE_ENABLE_TELEMETRY=1`) at `http://<agent-host>:<port>`.

> **A new source whose port is already bound fails the ENTIRE pipeline graph, not
> just that node.** The agent stops, its heartbeat freezes, and no data flows.
> Deploy with `edx pipelines deploy --wait` so a frozen heartbeat is reported
> immediately; if the agent was already up, it must be restarted to bind the new
> port. Pick a free, non-ephemeral port (below the OS ephemeral range, ~49152+).

## Structuring Processors: Attach Per Source

Prefer a **per-source multiprocessor wired directly off each source** for
source-specific parsing/enrichment/PII, and reserve a shared multiprocessor for
genuinely common, cross-source logic. This is the pattern to reach for by
default - it keeps each source's transforms next to the source (the Visual
Builder renders a source's dedicated multiprocessor attached to it) instead of
funnelling everything through one node in the middle.

```yaml
links:
- from: my_otlp
  to: my_otlp_multiprocessor      # source-specific: attached to the source
- from: my_otlp_multiprocessor
  to: common_multiprocessor       # shared cross-source logic (optional)
- from: common_multiprocessor
  to: edgedelta

nodes:
- name: my_otlp_multiprocessor
  type: sequence
  processors:
  - type: ottl_transform
    data_types: [log, metric]     # OTTL runs on metrics too, not just logs
    statements: |-
      set(resource["telemetry.source"], "my_app")
```

Name it `<source>_multiprocessor` and link it 1:1 off the source. The
`processors` list lives on the `sequence` node (not the source node).

## Redacting PII (OTTL recipes)

Push sources (OTLP especially) put identifiers in **attributes and resource, not
`body`** - so field-level redaction beats body regexes. Scrub structured fields,
and keep body regexes as defense-in-depth:

```yaml
- type: ottl_transform
  data_types: [log, metric]
  statements: |-
    delete_key(attributes, "user.email")
    set(resource["host.ip"], "[REDACTED_IP]") where EDXCoalesce(resource["host.ip"],"")!=""
- type: ottl_transform
  data_types: [log]
  statements: |-
    replace_pattern(body, "[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}", "[REDACTED_EMAIL]")
```

Verify with live capture (below): compare `before`/`after` on the multiprocessor
and confirm the PII fields are gone.

## Investigating Config Changes (Who Broke It?)

```bash
edx pipelines history <pipeline-id> --output table --columns version,timestamp,author,status,description
```

Correlate the deploy timestamps with the incident start. Roll back by
deploying the last good version:

```bash
edx pipelines deploy <pipeline-id> <last-good-version> --yes
```

## Live Capture (Debug Data In-Flight)

Sample real data before/after pipeline nodes on the running agents - the
fastest way to verify a processor or find where data is dropped:

```bash
# 1. Start a capture (all nodes, or scope with --nodes)
edx capture start <pipeline-id> --duration 2m --nodes mask_pii --max-items 50
# response contains the task "id"

# 2. Poll agent pickup status
edx capture status <task-id>

# 3. Fetch the captured before/after samples
edx capture results <pipeline-id>
```

Each node reports `before` and `after` arrays: compare them to verify
transformations, filters and routing. Captured items are **JSON-encoded
strings** - decode with jq's `fromjson`:

```bash
edx capture results <pipeline-id> | jq '[.[].nodes[].after[]] | map(fromjson)'
```

`--max-items` is capped at 100 per node.

## Troubleshooting

| Problem | Fix |
|---------|-----|
| Save rejected | Run `edx pipelines validate --file` for the error detail |
| Agent wedged / heartbeat frozen after deploy | A new source's port is already bound (fails the whole graph). Deploy with `--wait`; check `edx health problems` and agent logs; free the port or restart the agent |
| Agents not picking up deploy | `edx fleet deployments <pipeline-id>`; check agent connectivity |
| Capture returns nothing | Confirm node names (`edx pipelines get`), extend `--duration`, ensure traffic flows |
| Data missing downstream | Live-capture the node chain to find where items drop |
