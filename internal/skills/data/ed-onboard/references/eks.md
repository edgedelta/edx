# EKS

Inspect node capacity before choosing a deployment. For EC2-backed Linux nodes, use
the official Edge Delta Helm chart and pipeline deployment instructions. Reuse an
existing release where appropriate; record prior values/chart version for rollback.
Collect workload files with the node agent and receive OTLP from instrumented apps.
Infrastructure metrics do not prove application metric or trace collection.

Fargate pods cannot run a DaemonSet. Use an Edge Delta sidecar with shared application
files or direct application OTLP delivery where supported. Mixed clusters require
separate coverage for Fargate pods; a healthy node DaemonSet is not full coverage.
Treat control-plane logs separately from workload logs: identify AWS delivery constraints
rather than promising that a pod agent can bypass them.

Verify the actual cluster/resource identity in each requested signal. Do not create an
EKS cluster merely to onboard an existing ECS workload. Test EKS recipes separately
before labeling them live-validated.

References:
- https://docs.edgedelta.com/agent-installation/
- https://docs.aws.amazon.com/eks/latest/userguide/fargate.html

## Isolated validation fixture

Use a disposable namespace and restrict `kubernetes_input.include` to its test pod.
For example, `k8s.namespace.name=^onboarding-test$,k8s.pod.name=^telemetry-fixture$`.
A node agent collects container stdout; application metrics/traces still require
instrumentation. Expose an internal OTLP service using the official chart's `ports`
setting, with `port: 4318`, `protocol: TCP`, and `exposeInHost: false`.

Inspect the chart's rendered resources. For a small direct-path test, disable optional
coordinator, compactor, tracer, rollUp and forceReinstallApplications components when
not used by its pipeline. Supply the agent credential through a Kubernetes Secret.
Do not confuse the chart's secret-creation flag with referencing an existing secret;
its `secrets` environment-variable list can reference an existing `ED_API_KEY` secret.

An existing cluster returning Forbidden is not permission to alter its access controls.
For an authorized disposable-cluster test, create a separate cluster and scope any
access entry to that cluster. Remove namespace/Helm resources before deleting its
node group and cluster, and retain unrelated clusters unchanged.

## Native collection requirements

Use explicit include rules. The checked runtime's Kubernetes and metric filters reject
all records when the include list is empty, even though the configuration validates.
For an isolated test cluster, `kubernetes_input.include: ["k8s.namespace.name=.*"]`
collects all pod namespaces; scope this deliberately in a customer cluster. For direct
kubelet/cAdvisor metrics, set `ed_k8s_metrics_input.include` to `cadvisor=.*` and
`kubelet=.*`. Exclude node_exporter and kube_state_metrics when their exporters are not
installed; do not install additional collectors solely to satisfy default scrape targets.

Grant node-metrics read permissions and namespace-scoped lease permissions for the
agent's leader election. Chart 2.24.0 does not render nonResourceURLs through
clusterRoleRules; inspect rendered RBAC rather than assuming arbitrary rules survive.
For native host system metrics, mount host /proc and /sys read-only and set HOST_PROC
and HOST_SYS accordingly. Verify actual source records/metric points after deployment.

## Tested baseline and adaptation

The retained run used Go agent v2.24.0, Helm chart 2.24.0 and EC2-backed EKS 1.34.
It verified native stdout, kubelet/cAdvisor and host metrics, plus instrumented application
OTLP. This is a tested baseline, not a version requirement or a claim about every cluster.
Check runtime/chart compatibility, rendered RBAC, node architecture and capacity,
workload identity, source selection and egress before reusing the recipe.
