# Pitaya-Aligned Cluster Observability Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-14
Scope: Gate-only boundary for using Pitaya-aligned cluster observability vocabulary after the Pitaya service dispatch source-first sequence
Depends on: `decisions/ADR-0220-pitaya-aligned-cluster-event-bus-source-first-map.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0239`

The paired Simplified Chinese translation is `docs/pitaya-aligned-cluster-observability-boundary-gate.zh-CN.md`. The English file is authoritative.

This document defines cluster observability vocabulary only. It does not implement cluster observability behavior, service registry behavior, service selector behavior, heartbeat/liveness behavior, route targeting behavior, remote timeout or retry behavior, distributed session ownership behavior, distributed presence fanout behavior, cross-node error mapping behavior, cluster observability behavior, runtime behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

```yaml
pitaya_aligned_cluster_observability_boundary_gate: defined
completed_work_item: W-0331
decision: ADR-0239
check_rule: runtime.pitaya_aligned_cluster_observability_boundary_gate
source_pitaya_service_dispatch_map_decision: ADR-0220
standard: docs/pitaya-aligned-cluster-observability-boundary-gate.md
translation: docs/pitaya-aligned-cluster-observability-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: cluster_observability_vocabulary_boundary_defined_for_future_distributed_operations_planning
implementation_scope: gate_only_cluster_observability_vocabulary
future_implementation_work_item: W-0332
future_implementation_direction: implement_pitaya_aligned_cluster_observability_source_first_map
allowed_cluster_observability_vocabulary:
  - cluster_observability_boundary
  - local_operations_inspection_source
  - cluster_metrics_deferral
  - distributed_trace_deferral
  - node_health_view_deferral
  - redacted_cluster_diagnostics
runtime_behavior_added: false
node_identity_behavior_added: false
service_registry_behavior_added: false
service_selector_behavior_added: false
heartbeat_liveness_behavior_added: false
route_targeting_behavior_added: false
remote_timeout_retry_behavior_added: false
distributed_session_ownership_behavior_added: false
distributed_presence_fanout_behavior_added: false
cross_node_error_mapping_behavior_added: false
cluster_observability_behavior_added: false
service_export_behavior_added: false
remote_call_dispatch_behavior_added: false
server_to_server_rpc_behavior_added: false
frontend_message_forwarding_behavior_added: false
backend_service_route_ownership_behavior_added: false
cluster_event_bus_behavior_added: false
service_discovery_behavior_added: false
frontend_backend_role_runtime_behavior_added: false
distributed_runtime_implementation_added: false
runtime_endpoint_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
migration_added: false
dependency_added: false
persistence_added: false
startup_wiring_added: false
authentication_session_behavior_changed: false
hosted_deployment_added: false
sdk_added: false
release_artifact_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Purpose

The purpose is to keep cluster observability visible as a future Pitaya-class distributed operations concern without authorizing implementation. Agents may use this vocabulary for planning, inspection, ADRs, change specs, and future source-first maps only.

## 3. Vocabulary

- `cluster_observability_boundary`: planning vocabulary for cluster observability without implementing cluster observability behavior.
- `local_operations_inspection_source`: planning vocabulary for cluster observability without implementing cluster observability behavior.
- `cluster_metrics_deferral`: planning vocabulary for cluster observability without implementing cluster observability behavior.
- `distributed_trace_deferral`: planning vocabulary for cluster observability without implementing cluster observability behavior.
- `node_health_view_deferral`: planning vocabulary for cluster observability without implementing cluster observability behavior.
- `redacted_cluster_diagnostics`: planning vocabulary for cluster observability without implementing cluster observability behavior.

Forbidden use:

- Do not introduce concrete public API, package, route, method, wire, handler, dashboard, metrics, trace, admin, console, or inspector compatibility names from Pitaya or Nakama.
- Do not use this vocabulary as permission to add cluster observability behavior, service registry or selector behavior, heartbeat/liveness behavior, route targeting behavior, remote timeout/retry behavior, distributed session ownership behavior, distributed presence fanout behavior, cross-node error mapping behavior, cluster observability behavior, runtime behavior, protocol messages, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, or distributed runtime behavior.
- Do not classify raw tokens, credentials, lookup digests, verifier digests, verifier keys, DSNs, headers, cookies, query strings, subprotocol values, remote addresses, database payloads, local secret file contents, route payloads, event payloads, or concrete transport metadata as log-safe in this gate.

## 4. Current Mapping Candidate

```yaml
current_source_first_cluster_observability_mapping:
  local_operations_inspection:
    current: node tools/vibit inspect operations --json
    future_vocabulary: local_operations_inspection_source
    status: source_first_operations_inspection_only
  cluster_metrics:
    current: deferred
    future_vocabulary: cluster_metrics_deferral
    status: no_cluster_metrics_behavior
  distributed_trace:
    current: deferred
    future_vocabulary: distributed_trace_deferral
    status: no_distributed_trace_behavior
  node_health_view:
    current: deferred
    future_vocabulary: node_health_view_deferral
    status: no_node_health_view_behavior
  redacted_cluster_diagnostics:
    current: existing redaction posture in docs and checks
    future_vocabulary: redacted_cluster_diagnostics
    status: source_first_redaction_inspection_only
```

## 5. Ownership

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-cluster-observability-boundary-gate.md
  - .arch/reference.yaml
  - .arch/runtime.yaml
source_first_map_candidate_owner:
  - tools/vibit
runtime_behavior_owner: deferred
protocol_owner: unchanged
persistence_owner: unchanged
distributed_runtime_owner: deferred
```

Rules:

- Documentation and manifests may define cluster observability vocabulary and current source-first mapping.
- `tools/vibit` may emit a source-first cluster observability map because W-0332 is separately recorded.
- Existing runtime, protocol, route, repository, persistence, and transport behavior remain unchanged by this gate.

## 6. Reference Mapping

Nakama remains the primary product reference for broad game backend product capability pressure. Pitaya remains an architecture vocabulary reference for distributed runtime, components, handlers, services, sessions, routes, remote calls, service discovery, groups, lifecycle hooks, and operational concerns.

This gate maps those references into vibit-owned vocabulary only. It does not create direct compatibility, public API parity, runtime behavior, or distributed behavior.

## 7. Stop Conditions

Stop and require a later bounded work item before adding:

- cluster observability behavior;
- service registry or selector behavior;
- heartbeat/liveness behavior;
- route targeting behavior;
- remote timeout or retry behavior;
- distributed session ownership behavior;
- distributed presence fanout behavior;
- cross-node error mapping behavior;
- cluster observability behavior;
- runtime behavior;
- protocol messages or routes;
- Protobuf source;
- generated output;
- repository interfaces;
- PostgreSQL adapters;
- migrations;
- dependencies;
- hosted deployment;
- SDK publication;
- release artifacts;
- distributed runtime behavior;
- direct Nakama/Pitaya API compatibility;

## 8. Verification

Repository verification:

```text
node tools/vibit inspect rule runtime.pitaya_aligned_cluster_observability_boundary_gate
node tools/vibit check change define-pitaya-aligned-cluster-observability-boundary-gate --json
node tools/vibit check runtime --json
node tools/vibit check work --json
```
