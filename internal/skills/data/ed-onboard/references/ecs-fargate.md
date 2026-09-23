# ECS Fargate: Edge Delta sidecar

Preferred path: application → same-task Edge Delta agent → Edge Delta.
No FireLens, Fluent Bit or CloudWatch is required when the application can send OTLP
or write files to a shared volume.

1. Inspect the actual application telemetry output. For files, mount the same task
   volume into both containers and configure `file_input`, rotation and permissions.
   For OTLP, configure the agent receiver and point the application to localhost.
2. Use [the OTLP pipeline](../assets/otlp-pipeline.yaml) as a minimal starting point.
   One HTTP receiver handles logs, metrics and traces. Do not add one listener per signal
   on the same port. Confirm the selected runtime supports this configuration.
3. Resolve the image and agent environment through `edx pipelines deploy-command`.
   Pin a tested image version/digest. Inject credentials with ECS secrets, granting the
   execution role only access to that secret and required decryption permissions.
4. Add the agent to a new task definition revision. Allocate appropriate CPU/memory,
   configure outbound connectivity and bounded application export retries. The receiver
   is internal to the task; no inbound security-group rule is needed for localhost.
5. For existing services, deploy through their owner (Terraform/CloudFormation/etc.)
   and record the previous task definition revision. Adding a sidecar requires a rollout.
6. Verify each requested signal before retiring the previous collection path.

Fargate does not give a sidecar automatic access to another container's stdout/stderr.
If the application only emits stdout, evaluate shared-file or network logging changes.
Use another router only if those options cannot satisfy the requirement, and explain
that specific gap. Do not claim the Edge Delta image is a FireLens-compatible router
without verifying that capability.

## Temporary verification workload

`scripts/emit_otlp.py` is a standard-library-only application fixture, not a collector.
Run it beside the agent with `OTEL_EXPORTER_OTLP_ENDPOINT=http://127.0.0.1:4318`,
`ONBOARD_RUN_ID` set uniquely, and `ONBOARD_DURATION_SECONDS` bounded. It emits a log,
a gauge and a span with the same service name, plus a file log when `ONBOARD_LOG_FILE`
is configured. Copy it into your test application image or task command.

Do not add the awslogs driver to the test task: it would reintroduce CloudWatch costs.
Capture diagnostics through a controlled alternative if needed, without shipping
credentials. Stop the task, deregister/delete its task definition, remove its cluster,
secret, IAM role/policy and security group, then delete its pipeline. Keep a cleanup
ledger so interrupted tests can resume cleanup.

Reference: https://docs.edgedelta.com/ingest-from-ecs/

### Reproduce the disposable Fargate fixture

Render `python3 scripts/render_ecs_smoke.py --output /tmp/ecs-smoke.json` from the
skill directory. This creates a CloudFormation template, not live resources. Validate
it with `aws cloudformation validate-template`. Parameters are RunID, VpcID, AgentKey,
AgentImage (pin a digest), APIEndpoint and AIEndpoint from the deployment instructions.
AgentKey is a NoEcho parameter delivered through Secrets Manager. Use a private
parameter file, never a credential on the command line or in committed artifacts.

Create a uniquely named stack with CAPABILITY_IAM and ownership tags. The template
creates an ECS cluster, outbound-only security group, scoped secret execution role,
secret and task definition. After CREATE_COMPLETE, read its Cluster, SecurityGroup
and TaskDefinition outputs and run one FARGATE task in a subnet with working egress.
For a public test subnet set assignPublicIp=ENABLED; do not create a NAT gateway just
for this test. No ingress rule is needed. The application runs for at most 15 minutes.

Record the task ARN and every stack resource ID. Query the three signals as described
in verification.md. Stop the task and wait for STOPPED, then delete the stack and wait
for DELETE_COMPLETE. Delete any remaining inactive task definition revision, verify
secret/role/security group removal, and delete the test pipeline. CloudFormation
creation alone does not start the task or prove ingestion. Existing VPCs/subnets are
parameters and must not be deleted during cleanup.

For native file collection as well as OTLP, render with `--with-file-logs` and create
its pipeline from `assets/ecs-file-otlp-pipeline.yaml`. This mounts one task volume in
both containers; the Edge Delta mount is read-only. Verify a body containing
`onboarding file log` and `resource.ed.source.type=file_input` in the returned record.
File logs derive their service name from the file path, so the OTLP service-name filter
will not find them. Query their unique pipeline tag and marker instead.

## Tested baseline and native infrastructure coverage

The retained run used Go agent v2.24.0 with a shared-file input and application OTLP.
Native `http_pull_input` collected the Fargate metadata v4 `/stats` and `/task/stats`
endpoints from `ECS_CONTAINER_METADATA_URI_V4`; no router was needed. Agent-container
CPU/memory extraction is not proof of per-container metrics for the whole task. State
exactly which containers/signals have metric extraction versus only raw statistics.
Check application file rotation/permissions, shared mounts, metadata availability,
architecture, task resource budget and secret/egress permissions when adapting.

Cleanup examples above apply only when cleanup is authorized. A user instruction to
retain overrides disposable-fixture cleanup: keep the workload running and inventory it.
