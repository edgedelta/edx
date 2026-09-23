# AWS discovery

Always bind commands to the chosen `--profile` and `--region`. First resolve
`aws sts get-caller-identity`, then `aws ec2 describe-regions`. Inspect all selected
regions and retain failures rather than silently treating them as empty inventories.

- ECS: list clusters, services and tasks; describe task definitions, service capacity
  provider strategies and task launch types. FARGATE_SPOT also counts as Fargate.
  Inspect application log drivers, mounts, existing sidecars and task secrets.
- EKS: list/describe clusters, nodegroups and Fargate profiles; inspect Kubernetes
  nodes/workloads with the correct context and existing Edge Delta Helm releases.
- RDS: describe DB instances/clusters, engines, enabled log exports and Enhanced
  Monitoring. Do not turn on database logging or change parameter groups during discovery.
- Edge Delta: list pipelines and agents; identify matching resource ownership and
  destination profile before reusing a configuration.

Prefer existing egress. Do not expose OTLP publicly: same-task ECS traffic can use
localhost; Kubernetes receivers should remain internal unless deliberately designed
otherwise. Account for NAT/egress and agent compute even when CloudWatch is bypassed.

An empty environment is a valid discovery result. Create demo workloads only when
requested. For authorized tests use a unique run ID, tag owned resources, maintain
an exact ID ledger, and set bounded sender lifetimes. Reconcile create timeouts before
retrying, since the resource may already exist.
