# edx — Edge Delta CLI

`edx` is the Edge Delta command line interface. It covers the three Edge Delta
products:

| Product | Commands |
|---------|----------|
| **Pipeline** | `pipelines`, `fleet`, `capture`, `health`, `ingest` |
| **Observability** | `logs`, `patterns`, `metrics`, `traces`, `events`, `monitors`, `dashboards`, `facets`, `service-map` |
| **AI Teammate** | `ai issues`, `ai threads`, `ai channels`, `ai agents`, `ai connectors`, `ai activity` |

Plus `api` as a raw escape hatch for any endpoint, and `auth` for credentials.

## See it in action

Search logs with CQL:

![edx logs search](assets/edx-logs.gif)

List pipelines across the fleet:

![edx pipelines list](assets/edx-pipelines.gif)

Discover and query metrics:

![edx metrics](assets/edx-metrics.gif)

Surface negative log patterns (clustered error signatures):

![edx patterns](assets/edx-patterns.gif)

The clips are generated with [vhs](https://github.com/charmbracelet/vhs) from
the tape scripts in [`demo/`](demo/) — re-record with `vhs demo/<name>.tape`
(requires an authenticated `edx`; no credentials appear on screen).

## Install

Homebrew (macOS and Linux, Intel and Apple Silicon):

```bash
brew install edgedelta/tap/edx
# or
brew tap edgedelta/tap
brew install edx
```

Go toolchain:

```bash
go install github.com/edgedelta/edx@latest
# or from a checkout:
make install
```

## Authenticate

**OAuth (default)** — `edx auth login` opens your browser to log in; the access
token is refreshed automatically and your organization is read from the token:

```bash
edx auth login
edx auth status
```

**API token** — for CI/automation, create one in the Edge Delta web app
(Admin → API Tokens) and pass `--token` (with `--org-id`):

```bash
edx auth login --token <api-token> --org-id <org-id>
```

OAuth sends a Bearer JWT instead of the API token. A few AI endpoints don't
accept OAuth yet; for those, edx automatically falls back to the API token if
the profile has one (the fallback list lives in `internal/api/auth.go` —
`oauthUnsupported`). Pass `--token` if you rely on those endpoints.

Or use environment variables (they override the config file):

```bash
export ED_API_TOKEN=...
export ED_ORG_ID=...
export ED_ENV=staging                          # prod (default), staging or local
```

Config lives in `~/.config/edx/config.yaml` (mode 0600). To find that path and
inspect the configuration edx actually resolves (profile + env vars + flags
merged, secrets masked, no network call):

```bash
edx config path                                # print the config file location
edx config show                                # resolved config, secrets masked
edx config show --output yaml
```

Unlike `edx auth status`, `config show` prints to stdout and makes no API call,
so it works offline and pipes cleanly into `jq`.

### Environments and profiles

A profile targets an **environment** — `prod` (default), `staging` or `local` —
and the environment selects every service host together. This matters because
the AI Teammate features live on their own hosts, not under `api.edgedelta.com`:

| Environment | Main API | Chat (issues, threads, channels) | Agent (teammates) |
|-------------|----------|----------------------------------|-------------------|
| `prod`      | `api.edgedelta.com` | `chat.ai.edgedelta.com` | `agent.ai.edgedelta.com` |
| `staging`   | `api.staging.edgedelta.com` | `chat.ai.staging.edgedelta.com` | `agent.ai.staging.edgedelta.com` |
| `local`     | `localhost:4444` | `localhost:3001` | `localhost:3002` |

Switch environments with a profile, the `--env` flag, or `ED_ENV` — never by
hand-editing a single base URL:

```bash
edx auth login --profile staging --env staging --token ... --org-id ...
edx --profile staging ai issues list
edx --env staging ai issues list          # one-off, no profile needed
```

A profile is just `env` + `org_id` + `api_token` — there are no per-service URL
knobs to set; switching the environment moves every host together.

Escape hatch (advanced): to point a single service at a non-standard deployment
(a branch deploy, an odd local port) without disturbing the rest, set
`ED_API_URL`, `ED_CHAT_URL` or `ED_AGENT_URL`. When set, each overrides only
that service's host for the chosen environment.

## Output

Every command prints pretty JSON by default. Other formats:

```bash
edx pipelines list --output table
edx pipelines list --output table --columns id,tag,fleet_type,status,updated
edx monitors list --output yaml
edx logs search -q 'error' --output csv --columns timestamp,severity_text,body
```

## Quick reference

```bash
# --- Observability ----------------------------------------------------------
edx logs search -q 'service.name:"api" AND severity_text:"ERROR"' --lookback 1h
edx logs graph -q 'severity_text:"ERROR"' --group-by service.name --lookback 6h
edx patterns list --summary --lookback 1h          # clustered log signatures
edx patterns list --negative --query 'service.name:"api"'
edx metrics list --keyword cpu
edx metrics query --name http.request.duration --agg avg --group-by service.name
edx traces search -q 'status.code:"ERROR"' --include-children
edx events search -q 'event.type:"pattern_anomaly"' --lookback 6h
edx monitors list --output table
edx monitors states                                 # triggered/resolved states
edx service-map --lookback 1h
edx dashboards list

# --- Schema discovery (do this before writing CQL) --------------------------
edx facets keys --scope log
edx facets options --scope log --facet service.name
edx facets options --scope metric --facet name      # all metric names

# --- Pipeline / fleet management --------------------------------------------
edx pipelines list --output table
edx pipelines get <conf-id>                         # includes config content
edx pipelines history <conf-id>                     # version history
edx pipelines validate --file pipeline.yaml
edx pipelines save <conf-id> --file pipeline.yaml -d "add k8s source"
edx pipelines deploy <conf-id> <version> --yes      # deploy or roll back
edx pipelines agents <conf-id>
edx pipelines status <conf-id>
edx fleet agents
edx fleet deployments
edx health problems

# --- Live capture (debug data flowing through nodes) -------------------------
edx capture start <conf-id> --duration 2m --nodes mask_pii
edx capture status <task-id>
edx capture results <conf-id>

# --- AI Teammate (issues/threads/channels on chat.*, teammates on agent.*) ---
edx ai issues list                                  # open issues
edx ai issues list --include-closed --with-threads
edx ai issues critical --limit 10                   # ranked by criticality
edx ai issues health                                # current health score
edx ai issues timeline --lookback 24h
edx ai issues get <issue-id> --with-threads
edx ai issues close <issue-id>
edx ai channels list
edx ai threads list --channel <channel-id>
edx ai threads messages --channel <channel-id> <thread-id>
edx ai agents list                                  # AI teammates (agents)
edx ai connectors list
edx ai connectors specs
edx ai connectors update --file connector.json
edx ai activity

# --- Config inspection -------------------------------------------------------
edx config path                                     # where config.yaml lives
edx config show                                     # resolved config (masked)

# --- Agent skills ------------------------------------------------------------
edx skills list                                     # available skills
edx skills install                                  # auto-detect assistant, install all
edx skills install claude                           # or: cursor | codex | opencode | all
edx skills install claude --name ed-monitors        # just one skill
edx skills show ed-logs                             # print a skill's SKILL.md

# --- Anything else -----------------------------------------------------------
edx api GET /v1/orgs/{org}/users
edx api GET /tokens                                 # org-relative shorthand
edx api POST /pipelines/<conf-id>/save --data @save.json
edx api GET /issues --service chat                  # hit an AI service host
```

## CQL in 30 seconds

```
service.name:"api"                       field equals
service.name:("api" OR "web")            multiple values
-severity_text:"DEBUG"                   negation
error timeout                            full-text (logs/patterns/events only)
@response.code > 400                     numeric comparison on attributes
```

Regular expressions are **not** supported. Metrics and traces require
`field:"value"` syntax (no full-text). Time ranges: `--lookback 15m|1h|24h`
(Go durations) or `--from/--to` in ISO 8601 (`2006-01-02T15:04:05.000Z`).

## Agent skills

`edx skills install` drops the Edge Delta agent skills (ed-logs, ed-metrics,
ed-monitors, ...) into your coding assistant so an AI agent knows how to drive
edx. The platform is auto-detected from the environment, or pass it explicitly
(`claude`, `cursor`, `codex`, `opencode`, or `all`); `--project` installs into
the current repo instead of your home directory.

The skills are authored in [edgedelta/agent-skills](https://github.com/edgedelta/agent-skills)
and embedded in the binary, so an installed skill always matches the edx version
that wrote it — no network access needed.

## Development

```bash
make build         # bin/edx
make test
make vet
make sync-skills   # refresh the embedded skills from ../agent-skills, then commit
```

The embedded skills live in `internal/skills/data` (a generated copy of the
`agent-skills` repo). After skills change upstream, run `make sync-skills`
(optionally `SKILLS_SRC=/path/to/agent-skills`) and commit the result.

## Related

- [edgedelta/agent-skills](https://github.com/edgedelta/agent-skills) — AI agent skills built on edx
- [edgedelta/edgedelta-mcp-server](https://github.com/edgedelta/edgedelta-mcp-server) — MCP server for Edge Delta
