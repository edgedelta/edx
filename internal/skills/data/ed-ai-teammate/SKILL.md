---
name: ed-ai-teammate
description: AI Teammate - manage connectors (PagerDuty, Slack, GitHub, ...) and view teammate activity.
metadata:
  version: "1.0.0"
  author: edgedelta
  repository: https://github.com/edgedelta/agent-skills
  tags: edgedelta,ai,teammate,connectors
  alwaysApply: "false"
---

# Edge Delta AI Teammate

The AI Teammate ingests signals from connected tools (PagerDuty, Slack,
GitHub, ...) and acts on them. This skill manages those **connectors** and
inspects teammate **activity**.

## Prerequisites

The `edx` CLI must be installed and authenticated. See the **ed-edx** skill.

## Inspect

```bash
edx ai connectors list             # configured connectors
edx ai connectors specs            # available connector types + required fields
edx ai connectors environments     # where connectors can run
edx ai activity --lookback 24h     # teammate activity metrics
```

## Add or Update a Connector

1. Find the connector type and its required fields:

```bash
edx ai connectors specs --output json | jq '.[] | select(.type=="pagerduty")'
```

2. Build the request JSON per the spec (type, name, credentials/settings).

3. Apply it:

```bash
edx ai connectors update --file connector.json
```

Connector data flows through an ingestion pipeline that Edge Delta provisions
automatically; check it with `edx pipelines list --keyword ai`.

## Remove a Connector

The delete request body identifies the connector (same shape as update):

```bash
edx ai connectors delete --file connector.json --yes
```

## Troubleshooting

| Problem | Fix |
|---------|-----|
| Connector not ingesting | `edx pipelines list --keyword ai` then `edx health problems` |
| Unknown required fields | `edx ai connectors specs` is the source of truth |
| Credential errors | Re-apply with `edx ai connectors update --file` and fresh secrets |
