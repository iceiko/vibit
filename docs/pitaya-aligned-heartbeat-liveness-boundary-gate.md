# Pitaya-Aligned Heartbeat And Liveness Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-14
Scope: Gate-only boundary for using Pitaya-aligned heartbeat and liveness vocabulary after the Pitaya service dispatch source-first sequence
Depends on: `decisions/ADR-0220-pitaya-aligned-cluster-event-bus-source-first-map.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0227`

The paired Simplified Chinese translation is `docs/pitaya-aligned-heartbeat-liveness-boundary-gate.zh-CN.md`. The English file is authoritative.

This document defines heartbeat and liveness vocabulary only. It does not implement heartbeat and liveness behavior, service registry behavior, service selector behavior, heartbeat/liveness behavior, route targeting behavior, remote timeout or retry behavior, distributed session ownership behavior, distributed presence fanout behavior, cross-node error mapping behavior, cluster observability behavior, runtime behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

```yaml
pitaya_aligned_heartbeat_liveness_boundary_gate: defined
completed_work_item: W-0319
decision: ADR-0227
check_rule: runtime.pitaya_aligned_heartbeat_liveness_boundary_gate
source_pitaya_service_dispatch_map_decision: ADR-0220
standard: docs/pitaya-aligned-heartbeat-liveness-boundary-gate.md
translation: docs/pitaya-aligned-heartbeat-liveness-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: heartbeat_liveness_vocabulary_boundary_defined_for_future_distributed_operations_planning
implementation_scope: gate_only_heartbeat_liveness_vocabulary
future_implementation_work_item: W-0320
future_implementation_direction: implement_pitaya_aligned_heartbeat_liveness_source_first_map
allowed_heartbeat_liveness_vocabulary:
  - heartbeat_liveness_boundary
  - local_health_source
  - node_heartbeat_deferral
  - liveness_ttl_deferral
  - failure_detection_deferral
  - cluster_membership_deferral
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

The purpose is to keep heartbeat and liveness visible as a future Pitaya-class distributed operations concern without authorizing implementation. Agents may use this vocabulary for planning, inspection, ADRs, change specs, and future source-first maps only.

## 3. Vocabulary

- `heartbeat_liveness_boundary`: planning vocabulary for heartbeat and liveness without implementing heartbeat and liveness behavior.
- `local_health_source`: planning vocabulary for heartbeat and liveness without implementing heartbeat and liveness behavior.
- `node_heartbeat_deferral`: planning vocabulary for heartbeat and liveness without implementing heartbeat and liveness behavior.
- `liveness_ttl_deferral`: planning vocabulary for heartbeat and liveness without implementing heartbeat and liveness behavior.
- `failure_detection_deferral`: planning vocabulary for heartbeat and liveness without implementing heartbeat and liveness behavior.
- `cluster_membership_deferral`: planning vocabulary for heartbeat and liveness without implementing heartbeat and liveness behavior.

Forbidden use:

- Do not introduce concrete public API, package, route, method, wire, handler, dashboard, metrics, trace, admin, console, or inspector compatibility names from Pitaya or Nakama.
- Do not use this vocabulary as permission to add heartbeat and liveness behavior, service registry or selector behavior, heartbeat/liveness behavior, route targeting behavior, remote timeout/retry behavior, distributed session ownership behavior, distributed presence fanout behavior, cross-node error mapping behavior, cluster observability behavior, runtime behavior, protocol messages, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, or distributed runtime behavior.
- Do not classify raw tokens, credentials, lookup digests, verifier digests, verifier keys, DSNs, headers, cookies, query strings, subprotocol values, remote addresses, database payloads, local secret file contents, route payloads, event payloads, or concrete transport metadata as log-safe in this gate.

## 4. Current Mapping Candidate

```yaml
current_source_first_heartbeat_liveness_mapping:
  local_health_sources:
    current: runtime health readiness version config surface and runbook sources
    future_vocabulary: local_health_source
    status: local_health_inspection_only
  node_heartbeat:
    current: deferred
    future_vocabulary: node_heartbeat_deferral
    status: no_node_heartbeat_behavior
  liveness_ttl:
    current: deferred
    future_vocabulary: liveness_ttl_deferral
    status: no_liveness_ttl_behavior
  failure_detection:
    current: deferred
    future_vocabulary: failure_detection_deferral
    status: no_failure_detection_behavior
  cluster_membership:
    current: deferred
    future_vocabulary: cluster_membership_deferral
    status: no_cluster_membership_behavior
```

## 5. Ownership

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-heartbeat-liveness-boundary-gate.md
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

- Documentation and manifests may define heartbeat and liveness vocabulary and current source-first mapping.
- `tools/vibit` may emit a source-first heartbeat and liveness map because W-0320 is separately recorded.
- Existing runtime, protocol, route, repository, persistence, and transport behavior remain unchanged by this gate.

## 6. Reference Mapping

Nakama remains the primary product reference for broad game backend product capability pressure. Pitaya remains an architecture vocabulary reference for distributed runtime, components, handlers, services, sessions, routes, remote calls, service discovery, groups, lifecycle hooks, and operational concerns.

This gate maps those references into vibit-owned vocabulary only. It does not create direct compatibility, public API parity, runtime behavior, or distributed behavior.

## 7. Stop Conditions

Stop and require a later bounded work item before adding:

- heartbeat and liveness behavior;
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
node tools/vibit inspect rule runtime.pitaya_aligned_heartbeat_liveness_boundary_gate
node tools/vibit check change define-pitaya-aligned-heartbeat-liveness-boundary-gate --json
node tools/vibit check runtime --json
node tools/vibit check work --json
```
