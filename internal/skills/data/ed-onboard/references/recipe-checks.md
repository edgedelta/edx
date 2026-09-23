# Recipe contracts and executable checks

These fixtures encode the mistakes found during the retained ECS/EKS/RDS runs. They
are a starting point for RDS Go v2.24.0 transformation testing, not a portable AWS
installer or proof of compatibility with every runtime. No onboarding CLI stages or
backend onboarding state are required.

The sanitized JSON pipelines in `assets/recipe-checks/` are valid YAML-compatible
representations derived from the retained RDS recipes. Adapt their database identity,
service identity, receiver/queue and runtime assumptions before deployment. Do not deploy
the example account/queue. Existing ECS/EKS references describe their tested prerequisites;
this fixture suite does not claim to execute ECS/EKS source collection.

From the skill directory:

```bash
# Offline structure checks and negative tests of the assertion checker itself.
python3 scripts/check_recipes.py --self-test

# Actual processor execution through an existing authorized edx profile/pipeline.
# This invokes only the dry-run endpoint, never save/deploy or a telemetry ingest API.
python3 scripts/check_recipes.py --run --edx /path/to/edx \
  --profile staging --conf-id "$CONF_ID" --evidence-dir /tmp/recipe-evidence
```

The second command submits fabricated, nonsecret fixtures to the processor preview API.
It requires that endpoint to support these sequence nodes. It stops at the first API
failure, preserves diagnostics and exits nonzero. If this endpoint returns an error, do not interpret a checker self-test or schema
validation as a
substitute for successful processor execution. Use native live capture in an authorized
isolated deployment when preview is unavailable, and keep that evidence separate.

Contracts cover:
- Average gauge values from sum/count; percent, byte, IOPS and second semantics. The
  deployed v2.24 extractor uses `1` for non-byte units; metric names/descriptions carry
  percent/rate/seconds meaning. An unknown metric emits nothing.
- Namespace and DB identity must both match. Aggregate records without an instance,
  another database, another namespace and zero count emit nothing. Valid zero stays zero.
- Exactly one matching metric per accepted input; extra outputs fail. This detects
  accidentally OR-ing independent conditions. An explicit post-extraction metric-only
  filter drops unmatched records: `keep_item: false` alone did not drop them in preview.
  Original raw metric logs are not retained.
- Engine batch expansion preserves both messages, log-group identity and each event's
  original millisecond timestamp, rather than the receiver's ingestion time.
- Directly attached source/destination sequences and meaningful nested processor names.
  Agent self-telemetry/internal statistics retain their established direct routing.

Each case in `assets/recipe-checks/cases.json` specifies an input and expected output
fields. Extra enrichment fields are allowed, but output count, type, signal name,
value, unit, source identity and timestamp must match. Output ordering is irrelevant.
Add source-specific positive/negative cases when changing a recipe; do not duplicate
an OTTL interpreter in Python to produce the supposed actual output.

For outputs obtained by a separate execution harness, normalize them to a JSON object
mapping every case ID to an array of actual item objects and run:

```bash
python3 scripts/check_recipes.py --actual actual.json
```

Record where those outputs came from. Copying expected values into actual.json tests
only the checker, not the transformation. Live telemetry still needs source and indexed
query verification as described in [verification](verification.md).
