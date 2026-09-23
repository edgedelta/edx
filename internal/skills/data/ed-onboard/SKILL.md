---
name: ed-onboard
description: Connect customer environments to Edge Delta, selecting direct collection paths, provisioning pipelines and agents, and verifying source-specific logs, metrics and traces. Use for first-time onboarding or adding telemetry sources.
metadata:
  version: "1.0.0"
  author: edgedelta
  repository: https://github.com/edgedelta/agent-skills
  tags: edgedelta,onboarding,aws,ecs,eks,rds
---

# Onboard telemetry to Edge Delta

Use `edx` for Edge Delta operations and the customer's cloud CLI or infrastructure
as code for their environment. See **ed-edx** for authentication and **ed-pipelines**
for existing pipeline changes. Do not confuse AWS profiles with Edge Delta profiles.

## Collection principles

- Minimize hops. Prefer source → Edge Delta agent → Edge Delta, or compatible
  source → Edge Delta when local collection/processing is unnecessary.
- Prefer the Edge Delta agent. Add another collector only for a verified capability
  gap; explain the requirement and cost before selecting that exception.
- Avoid CloudWatch or other paid intermediary ingestion/storage when direct collection
  is supported. Existing CloudWatch usage does not make forwarding the preferred path.
- Test the actual collection path. Use Edge Delta sources to pull external data or
  receive native pushes. Do not fetch source data in an ad-hoc script and resubmit it
  to OTLP as proof of source support. Application instrumentation may emit its own
  activity, but does not establish native infrastructure/database collection.
- Where agents can run, include applicable workload logs and host/container/Kubernetes
  metrics; synthetic OTLP traffic alone is not environment onboarding.
- Decide per signal. An installed agent cannot manufacture application traces, access
  an isolated container's stdout, or read a managed database's host filesystem.
- Preserve existing collection until the replacement is verified; report duplicate
  ingestion during migration and remove the old path only within authorized scope.

## Adaptive onboarding responsibilities

These are responsibilities, not mandatory command stages. Revisit them as discovery
changes the best collection path; edx provides operations, not an onboarding state machine.

1. Establish AWS account/profile/regions, Edge Delta profile/org, resources and requested
   signals from the user's scope. Read-only discovery may continue while ambiguities
   are resolved. Never infer a destination organization from an AWS profile name.
2. Discover resources, current telemetry paths, agent deployments and pipeline IDs.
   Record each region/service as inspected, empty or inaccessible. Access denied is
   not evidence of absence. Read [AWS discovery](references/aws.md).
3. Select the shortest supported path. Read only the applicable recipe:
   [ECS Fargate](references/ecs-fargate.md), [EKS](references/eks.md),
   [RDS](references/rds.md). For other sources, inspect agent source capabilities and
   deployment documentation before claiming support; record gaps explicitly. Build a
   per-resource coverage table: requested signal, native source, proposed path, exception,
   and evidence/status. Separate server signals from application/client telemetry.
   Minimizing hops does not justify silently omitting a signal: propose the supported
   exception and tradeoff, respecting any authorization already given.
4. Prepare concrete configuration and infrastructure changes. Record account, org,
   resource IDs, signals, paths, versions, permissions, network needs, rollout effects,
   ownership, and cleanup/rollback actions. Keep credentials out of plans and output.
   Keep a durable local inventory as described in [handoff](references/handoff.md);
   no backend onboarding record or fixed plan format is required.
   Reuse existing owned resources when appropriate; do not create duplicates on reruns.
   Respect Terraform/GitOps ownership by changing its source rather than introducing drift.
5. Execute within the user's authorization. Existing authorization to deploy/test is
   sufficient; do not ask again for every step. For unapproved mutations, present the
   concrete plan first. Track created IDs immediately, including partial failures.
6. Verify every requested signal with [verification](references/verification.md).
   Use [recipe checks](references/recipe-checks.md) when adapting parsing/extraction.
   Distinguish healthy deployment, accepted export, correct transformation and queryable
   source telemetry.
   Report unsupported, blocked and unverified signals separately from passing ones.
7. Honor retention instructions. If the user asks to keep the deployment, leave agents
   and workloads running and record IDs, access instructions, ongoing costs and the
   eventual cleanup procedure. Do not stop agents to force a verification flush.
   Otherwise, for temporary tests, stop senders, uninstall agents, remove owned cloud resources,
   then delete the temporary pipeline. Verify deletion and list residual resources.
   Never delete preexisting customer resources or broaden cleanup by name prefix alone.

## Pipeline authoring

Follow **ed-pipelines** for source-attached and destination-attached multiprocessors.
Every application source and destination gets a `type: sequence` node named exactly
`<node-name>_multiprocessor`, immediately after the source or before the destination,
even when empty. These attachments are distinct from standalone middle processors.
Keep the default direct path for agent self-telemetry and internal statistics.

Give every nested processor a meaningful display name in its JSON-encoded `metadata`
string, for example `metadata: '{"name":"Set service identity"}'`. Preserve other
metadata keys; do not leave Custom OTTL processors displayed as "Custom". Check these
conventions in generated templates and before validating or creating the pipeline.

## New edge pipelines

Requires an edx build exposing `pipelines create/delete`; inspect `--help` first.

```bash
edx --profile "$ED_PROFILE" pipelines list --keyword "$RUN_ID"
edx --profile "$ED_PROFILE" pipelines validate --file pipeline.yaml
edx --profile "$ED_PROFILE" pipelines create --file pipeline.yaml \
  --tag "$RUN_ID" --environment Docker
# Capture the returned id privately. The initial configuration is already deployed.
edx --profile "$ED_PROFILE" pipelines deploy-command "$CONF_ID"
```

Treat deployment commands as sensitive: they may contain agent credentials. Use the
server-provided image and environment settings for the selected org/runtime. A pipeline
ID is not an organization API token. Do not embed an admin API token in a workload.
Use the deploy credential through the customer's secret delivery mechanism.

For Kubernetes use `--environment Kubernetes --fleet-subtype Edge` (or the actual
Gateway/Coordinator role). Validate both the API call and its returned validation
result; a successful HTTP response alone does not prove valid configuration.

Stop/uninstall test agents before `edx pipelines delete "$CONF_ID" --yes`.
Deleting a pipeline does not remove cloud infrastructure or already ingested telemetry.
