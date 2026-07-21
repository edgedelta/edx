---
name: ed-ai-teammate
description: AI Teammate - manage connectors (PagerDuty, Slack, GitHub, ...), update teammates (agents) and view teammate activity.
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
edx ai agents list                 # AI Teammates (agents); alias: edx ai teammates
edx ai agents get <agent-id>       # a single teammate's full definition
```

## Update a Teammate (Agent)

`edx ai agents update` prompts for confirmation unless `--yes`.

### Just the prompts (the common case)

Use the `--*-prompt` flags — each takes an inline string or `@file` (`@-` for
stdin). The command reads the current teammate, backfills whichever prompt you
did not pass (the service requires both `masterPrompt` and `userPrompt` on every
update), and sends only the prompts. Model, temperature, tools and everything
else are left untouched, so you never deal with model-tuning validation:

```bash
edx ai agents update <id> --master-prompt @master.md --user-prompt @user.md
edx ai agents update <id> --master-prompt "You are a concise SRE assistant."
```

### Any field (clone-and-edit)

Like `edx monitors update`: fetch, edit the JSON, apply it back.

```bash
edx ai agents get <agent-id> > agent.json      # edit fields under "data":
                                               # model, toolConfigurations,
                                               # connectors, priority, ...
edx ai agents update <agent-id> --file agent.json
```

The whole `get` envelope is accepted and unwrapped automatically. Only the
fields you include are changed; a field set to `null` is cleared. Use
`-f`/`--file -` to read from stdin. Note: re-sending the full object re-validates
model-tuning fields (`model`, `modelTemperature`) even if untouched — to change
only prompts prefer the flags above, or send just the prompt fields:

```bash
edx ai agents get <id> | jq '.data | {masterPrompt,userPrompt,toolingPrompt}' \
  | edx ai agents update <id> --file - --yes
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
| Teammate update rejected (4xx) | Re-fetch with `edx ai agents get <id>`, edit only the `data` fields, re-apply |
