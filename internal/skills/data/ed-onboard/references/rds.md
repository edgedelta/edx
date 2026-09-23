# RDS

Select collection per engine and signal. Inventory engine logs, database infrastructure
metrics and application/database-client traces separately. Application telemetry alone
is incomplete RDS onboarding. An Edge Delta agent cannot be installed on an RDS host
or read its private filesystem. When a direct path is unavailable, proactively propose
the supported exception and its cost; do not silently omit server signals. An existing
user instruction authorizing that exception is sufficient to proceed.

- Database logs: RDS provides DescribeDBLogFiles/DownloadDBLogFilePortion for supported
  engines/logs. Check whether the deployed Edge Delta runtime has a maintained input
  implementing this path. If absent, report a capability gap; do not substitute an
  ad-hoc polling loop and call it production onboarding. Reliable collection needs
  durable offsets, rotation handling, throttling, retries and duplicate control.
- Database metrics: check supported direct database receivers and required SQL grants.
  Their metrics may differ from AWS infrastructure metrics; do not claim equivalence.
- AWS metrics / Enhanced Monitoring: discover the requested AWS-specific signal.
  Enhanced Monitoring delivers into CloudWatch Logs. Explain this exception and obtain
  scope for that path rather than enabling CloudWatch as a blanket default.
- Traces: instrument the application making database calls. RDS log forwarding does
  not create distributed application traces.

Never enable verbose SQL logging, alter database parameters, or grant database access
merely during discovery. Describe required changes and their operational effects.
Mark RDS live verification separately from successful ECS tests.

References:
- https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/DownloadCompleteDBLogFile.html
- https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_Monitoring.OS.Enabling.html

## Native-source verification

No direct RDS log/SQL input was found in the checked Go agent v2.24.0. Do not use a
one-shot AWS API downloader or SQL-to-OTLP probe as evidence of RDS source support.
Verify collection using a maintained Edge Delta source/forwarder. The documented RDS
integration uses CloudWatch; if avoiding CloudWatch is required, report the server-side
capability gap instead of silently substituting an adapter. Obtain scope for any
CloudWatch exception and disclose its additional ingestion/storage path.

An instrumented application can still emit its own measured database-client spans and
request metrics. A node agent can collect that application's stdout. Label those as
application telemetry, separately from RDS engine logs and database infrastructure
metrics. They do not establish complete RDS onboarding.

## Supported CloudWatch exception

After the user authorizes this path:

- Engine logs: enable only the engine's supported RDS log export, then use a CloudWatch
  Logs subscription to the official Edge Delta Lambda Forwarder and an Edge Delta HTTP
  receiver. Prefer this push path. Scope Lambda invocation to the exact log group and
  account. The forwarder uses `ED_ENDPOINT`; do not assume arbitrary authorization
  headers are supported. Set explicit log retention for the test/customer requirement.
- Native AWS/RDS metrics: CloudWatch Metric Stream → Firehose → compressed S3 objects
  with S3 event notifications → SQS → Edge Delta `s3_input`. The agent reads objects
  itself. Use a scoped workload identity (for EKS, IRSA), not long-lived access keys.
  No additional collector or SQL-to-OTLP probe is required.
- Filter metric names at the stream and the selected DBInstanceIdentifier inside the
  Edge Delta pipeline. AWS streams include all dimension combinations for selected
  namespace/metric names: downstream filtering does not remove their AWS streaming cost.
- CloudWatch JSON metric streams contain newline-delimited individual records; do not
  assume a `records` wrapper. Verify the actual format through native live capture.
  CPU utilization, connections, free memory/storage, IOPS and latency are gauges:
  derive each interval's average from `value.sum / value.count`, guard count > 0,
  retain DB identity and use the source measurement timestamp. Do not sum gauge samples
  or treat IOPS as a cumulative counter. Each extract_metric rule should use one
  conjunctive condition for namespace, instance, metric name and count: separate entries
  in `conditions` are OR, not AND. Drop raw metric records after extraction unless
  keeping them as logs is explicitly useful. In a sequence, `keep_item: false` alone
  does not guarantee unmatched records are dropped: add an explicit post-extraction
  metric-only filter and test zero-count/unknown-metric cases. Agent v2.24 extract_metric accepts units `1`
  and `b`; document percent/seconds/rate semantics in descriptions when using `1`.
- Expand the forwarder's `logEvents` array, preserve event message, timestamp and log
  group/stream identity. In the Go agent's custom OTTL context, `timestamp` is epoch
  milliseconds; setting an arbitrary `time_unix_nano` field does not change event time.
- Attach `<source>_multiprocessor` and `<destination>_multiprocessor`; retain empty
  attachments and give every nested processor a descriptive `metadata.name`.
- Verify indexed engine records with their AWS log-group identity and nonempty native
  metric timeseries in edx. A running stream, success response or empty aggregate zero
  is insufficient. Account for AWS delivery buffering and Edge Delta indexing latency.
  Keep application/client traces separate from these server-side signals.

Standard RDS metrics already publish to CloudWatch. This exception adds streaming and
transport charges; exported engine logs also add CloudWatch Logs ingestion/storage.
Do not describe all of these costs as a flat doubling of the entire bill. Minimize
retention and metric selection according to user requirements, and retain deployed
resources when the user has requested review before cleanup.

References:
- https://docs.edgedelta.com/aws-lambda-forwarder/
- https://docs.edgedelta.com/rds-cloudwatch-metrics-ingestion/
- https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/CloudWatch-metric-streams-setup.html
- https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/monitoring-cloudwatch.html
