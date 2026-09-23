# Coverage and resumable handoff

Keep a small, durable local file in the task workspace; do not rely only on conversation
memory or temporary-directory state. Match the customer's existing IaC workflow. This
is an inventory/evidence artifact, not a required backend object or fixed execution plan.

Record:
- AWS account/profile/regions and Edge Delta org/profile separately; deployment owner
  (Terraform, CloudFormation, GitOps or manual) and selected runtime/chart versions.
- Per resource and requested signal: selected native source/path, direct-path capability
  gaps, authorized intermediary/cost exceptions, status and remaining work. RDS engine
  logs, native database metrics and application/client traces are separate entries.
- Immediately after creation: exact resource IDs, pipeline/environment IDs, dependencies,
  ownership (created versus preexisting), operation result and any partial failure.
- Verification query/time range and nonsecret evidence paths; distinguish schema checks,
  processor execution, live capture and indexed queries. Record invalid test data too.
- Retention instruction, recurring infrastructure/intermediary costs, previous revisions,
  and eventual rollback/cleanup order. Never put secrets or deploy credentials here.

On resume, reconcile the inventory with live state before retrying a create operation.
An API timeout can leave a created resource; list/inspect before trying again. Do not
identify ownership by a name prefix alone. Remove only owned resources when authorized;
when asked to retain, keep agents and senders running and supply inspection instructions.
