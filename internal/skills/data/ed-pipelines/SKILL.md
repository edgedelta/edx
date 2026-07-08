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
edx pipelines agents <conf-id>        # agents running this pipeline
edx pipelines status <conf-id>        # running / suspended
edx fleet agents                      # all agents org-wide
edx fleet deployments                 # rollout status across pipelines
edx health problems                   # components currently failing
```

## Config Change Workflow

```bash
# 1. Fetch the current config (content field holds the YAML)
edx pipelines get <conf-id> > pipeline.json
jq -r .content pipeline.json > pipeline.yaml

# 2. Edit pipeline.yaml (to develop the transform itself, see the ed-pipeline-tuning skill)

# 3. Dry-run the change on sample logs before saving (offline, no deploy)
edx pipelines test ottl <conf-id> --file samples.jsonl --statements '<ottl>'

# 4. Validate before saving
edx pipelines validate --file pipeline.yaml

# 5. Save a new version with a meaningful description
edx pipelines save <conf-id> --file pipeline.yaml -d "mask PII in checkout logs"

# 6. Deploy the version from the save response
edx pipelines deploy <conf-id> <version> --yes

# 7. Watch the rollout
edx fleet deployments <conf-id>
```

## Investigating Config Changes (Who Broke It?)

```bash
edx pipelines history <conf-id> --output table --columns timestamp,author,status,description
```

Correlate the deploy timestamps with the incident start. Roll back by
deploying the last good version:

```bash
edx pipelines deploy <conf-id> <last-good-version> --yes
```

## Live Capture (Debug Data In-Flight)

Sample real data before/after pipeline nodes on the running agents - the
fastest way to verify a processor or find where data is dropped:

```bash
# 1. Start a capture (all nodes, or scope with --nodes)
edx capture start <conf-id> --duration 2m --nodes mask_pii --max-items 50
# response contains the task "id"

# 2. Poll agent pickup status
edx capture status <task-id>

# 3. Fetch the captured before/after samples
edx capture results <conf-id>
```

Each node reports `before` and `after` arrays: compare them to verify
transformations, filters and routing.

## Troubleshooting

| Problem | Fix |
|---------|-----|
| Save rejected | Run `edx pipelines validate --file` for the error detail |
| Agents not picking up deploy | `edx fleet deployments <conf-id>`; check agent connectivity |
| Capture returns nothing | Confirm node names (`edx pipelines get`), extend `--duration`, ensure traffic flows |
| Data missing downstream | Live-capture the node chain to find where items drop |
