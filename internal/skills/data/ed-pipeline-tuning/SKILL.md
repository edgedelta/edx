---
name: ed-pipeline-tuning
description: Tune pipeline data quality - parse, structure, enrich, mask and roll logs up to metrics/patterns, with offline dry-run and live-capture verification.
metadata:
  version: "1.0.0"
  author: edgedelta
  repository: https://github.com/edgedelta/agent-skills
  tags: edgedelta,pipelines,data-quality,ottl,parsing,enrichment,masking,log-to-metric
  globs: "**/edgedelta*.yaml,**/pipeline*.yaml"
  alwaysApply: "false"
---

# Edge Delta Pipeline Tuning (Data Quality)

Improve the *content* of telemetry: parse unstructured logs into fields, set
severity and timestamps, mask PII, drop noise, and roll logs up into metrics and
patterns. **ed-pipelines** covers how to ship a change (save/deploy/rollback);
this skill covers *what* to change and how to express it, and how to prove it
worked.

## Prerequisites

`edx` installed and authenticated (see **ed-edx**). Requires edx >= 0.8.0 for
`edx pipelines test`.

## The golden rule

**A valid config is not a working transform.** `validate` only checks that the
YAML is well-formed; it does NOT prove your rule changed the data. Always:

1. **`edx pipelines test`** - dry-run the rule on sample logs (offline, no deploy)
2. **`edx capture`** - after deploy, confirm the transform on live data

Skipping these is the #1 cause of "I deployed it but nothing changed" (e.g. an
OTTL statement that references the wrong field passes validation and silently
no-ops).

## The tuning loop

```bash
edx pipelines get <pipeline-id> | jq -r .content > pipeline.yaml   # 1. current config
edx pipelines test ottl <pipeline-id> -f samples.jsonl -s '<ottl>' # 2. dry-run offline
# 3. edit pipeline.yaml to add the processor
edx pipelines validate -f pipeline.yaml                            # 4. validate
edx pipelines save <pipeline-id> -f pipeline.yaml -d "why"         # 5. save version
edx pipelines deploy <pipeline-id> <version> --yes                 # 6. deploy
edx capture start <pipeline-id> --duration 2m && edx capture results <pipeline-id>  # 7. verify live
```

`samples.jsonl` is one JSON-encoded OTEL log per line, e.g.
`{"body":"...","attributes":{},"resource":{},"timestamp":1}`.

## Processor toolbox

Processors live in a `sequence` multiprocessor node (`processors:` list) or, for
log->metric/pattern, as top-level nodes wired via `links`.

| Goal | Node / OTTL |
|------|-------------|
| Transform anything | `ottl_transform` (`statements:`) |
| Parse JSON body | `merge_maps(attributes, ParseJSON(body), "upsert")` |
| Parse key=value / logfmt | `merge_maps(attributes, ParseKeyValue(body), "upsert")` |
| Extract fields by regex | `merge_maps(attributes, ExtractPatterns(body, "(?P<f>...)"), "upsert")` |
| Parse a known format | `grok` node (`pattern: %{COMMONAPACHELOG}`) |
| Mask PII | `replace_pattern(body, "<regex>", "[REDACTED]")` or a `mask` node |
| Drop noise | `filter`, `delete_key(attributes, "k")` |
| Logs -> metric | `log_to_metric` node (`pattern`, `metric_name`, `dimension_groups`) |
| Logs -> patterns | `log_to_pattern` node (`num_of_clusters`, `field_path`) |

Dry-run non-OTTL nodes with `edx pipelines test node <id> -f samples.jsonl --node node.json`,
where node.json is `{"id":"...","type":"log_to_metric","configuration":{...}}`.

## Recipes (verified OTTL)

**Auto-detect format and parse** (guard each parser so it only runs on matching lines):
```
merge_maps(attributes, ParseJSON(body), "upsert") where IsMatch(body, "^\\s*[{]")
merge_maps(attributes, ExtractPatterns(body, "^(?P<client_ip>\\S+) \\S+ \\S+ \\[[^\\]]+\\] \"(?P<method>\\S+) (?P<path>\\S+) [^\"]+\" (?P<status>\\d+) (?P<bytes>\\d+)"), "upsert") where IsMatch(body, "HTTP/")
merge_maps(attributes, ParseKeyValue(body), "upsert") where not IsMatch(body, "^\\s*[{]") and IsMatch(body, "[a-zA-Z0-9_]+=")
```

**Severity - explicit level tag wins over free-text keywords** (a `<Warning>` whose
message contains the word "error" must stay WARN):
```
set(severity_text, ToUpperCase(attributes["level"])) where EDXCoalesce(attributes["level"], "") != ""
set(severity_text, "ERROR") where EDXCoalesce(severity_text, "") == "" and IsMatch(body, "(?i)<(error|fault|critical)>")
set(severity_text, "WARN")  where EDXCoalesce(severity_text, "") == "" and IsMatch(body, "(?i)<warning>")
set(severity_text, "ERROR") where EDXCoalesce(severity_text, "") == "" and IsMatch(body, "(?i)(\\berror\\b|\\bfatal\\b|\\bpanic\\b)")
set(severity_text, "WARN")  where EDXCoalesce(severity_text, "") == "" and IsMatch(body, "(?i)(\\bwarn\\b|\\bwarning\\b)")
set(severity_text, "INFO")  where EDXCoalesce(severity_text, "") == ""
```

**Timestamp from an embedded string** (the field is `timestamp`; extract to
seconds precision first if the format has sub-seconds):
```
set(timestamp, UnixMilli(Time(attributes["ts"], "%Y-%m-%dT%H:%M:%SZ"))) where EDXCoalesce(attributes["ts"], "") != ""
```

**Mask PII (IPv4) + drop empty/redundant fields:**
```
replace_pattern(body, "\\b(?:[0-9]{1,3}\\.){3}[0-9]{1,3}\\b", "[REDACTED_IP]")
delete_key(attributes, "pid") where EDXCoalesce(attributes["pid"], "") == ""
```

**Logs -> metric** (count errors, dimensioned) as a top-level node:
```yaml
- name: error_count
  type: log_to_metric
  pattern: (?i)(error|fatal|panic)
  interval: 1m0s
  metric_name: error_count
  dimension_groups:
    - field_dimensions: [ item["resource"]["service.name"] ]
      enabled_stats: [ count ]
```

Notes: guards matter - `ParseJSON` on a body that starts with `{` but is
malformed errors the statement (skipped at runtime, so the log passes through
unparsed). OTTL string literals need doubled backslashes (`\\b` in YAML = `\b`
regex). `EDXCoalesce(x, "")` safely treats a missing/nil field as empty.

## Where to put processors (pre vs post — avoid the middle)

A pipeline has two natural processing stages. Put each processor in the one that
matches its scope; a processor that sits in neither is a smell.

```
input ──▶ [PRE: per-source MP] ──▶ (rare middle) ──▶ [POST: pre-output MP] ──▶ output
```

- **PRE — the per-source multiprocessor** (`<source>_multiprocessor`, right after
  each input). Everything **specific to that source**, because each source has one
  known shape: parse/structure it, set severity and timestamp from *its* format,
  drop noise (`filter`), sample low-value logs, and branch to `log_to_metric` /
  `log_to_pattern`. Also clean up the fields *this parser* created (e.g.
  `delete_key(attributes, "log_ts")` right after using it to set `timestamp`).

- **POST — the multiprocessor just before the output** (all sources converge).
  Only things **common to every source**: destination tagging (e.g. everything
  bound for Splunk gets a `resource` tag), org-wide PII masking, a severity
  fallback. This is the "output node" stage.

- **Middle — anything between pre and post — should be rare.** If a step is
  source-specific it belongs in PRE; if it's universal it belongs in POST. A
  source-specific rule stranded in the shared stage (e.g. macOS `<Level>` severity
  parsing applied to every source) is the most common mistake — move it into the
  source's own multiprocessor.

Rule of thumb: *does this apply to one source or all of them?* One → PRE. All →
POST. Neither cleanly → reconsider the design.

## Troubleshooting

| Problem | Fix |
|---------|-----|
| Deployed but data unchanged | You skipped dry-run/capture. `edx pipelines test ottl` the statement; check it references the right field. |
| `test` rejects the statement | Read the error - usually an undefined OTTL function or a bad regex escape (`\\` in YAML). |
| Non-OTTL node (log_to_metric, ...) won't dry-run in `test ottl` | Use `edx pipelines test node` with the node as `{"id","type","configuration":{...}}`. |
| Severity misclassified | Set from the explicit level tag/field *before* free-text keyword matching. |
| Metric/pattern shows nothing in capture | They emit on an interval; extend `--duration` past one `interval`/reporting window. |
