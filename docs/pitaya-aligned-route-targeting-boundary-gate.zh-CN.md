# Pitaya-Aligned Route Targeting Boundary Gate 中文版

Status: Accepted v0.1
Last updated: 2026-06-14
Scope: Gate-only boundary for using Pitaya-aligned route targeting vocabulary after the Pitaya service dispatch source-first sequence
Depends on: `decisions/ADR-0220-pitaya-aligned-cluster-event-bus-source-first-map.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0229`

说明：本文件是 `docs/pitaya-aligned-route-targeting-boundary-gate.md` 的简体中文译本。英文版本是权威版本。

This document defines route targeting vocabulary only. It does not implement route targeting behavior, service registry behavior, service selector behavior, heartbeat/liveness behavior, route targeting behavior, remote timeout or retry behavior, distributed session ownership behavior, distributed presence fanout behavior, cross-node error mapping behavior, cluster observability behavior, runtime behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

```yaml
pitaya_aligned_route_targeting_boundary_gate: defined
completed_work_item: W-0321
decision: ADR-0229
check_rule: runtime.pitaya_aligned_route_targeting_boundary_gate
source_pitaya_service_dispatch_map_decision: ADR-0220
standard: docs/pitaya-aligned-route-targeting-boundary-gate.md
translation: docs/pitaya-aligned-route-targeting-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: route_targeting_vocabulary_boundary_defined_for_future_distributed_operations_planning
implementation_scope: gate_only_route_targeting_vocabulary
future_implementation_work_item: W-0322
future_implementation_direction: implement_pitaya_aligned_route_targeting_source_first_map
allowed_route_targeting_vocabulary:
  - route_targeting_boundary
  - local_route_registry_source
  - frontend_forwarding_handoff
  - backend_service_target_deferral
  - session_aware_targeting_deferral
  - distributed_route_targeting_deferral
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

The purpose is to keep route targeting visible as a future Pitaya-class distributed operations concern without authorizing implementation. Agents may use this vocabulary for planning, inspection, ADRs, change specs, and future source-first maps only.

## 3. Vocabulary

- `route_targeting_boundary`: planning vocabulary for route targeting without implementing route targeting behavior.
- `local_route_registry_source`: planning vocabulary for route targeting without implementing route targeting behavior.
- `frontend_forwarding_handoff`: planning vocabulary for route targeting without implementing route targeting behavior.
- `backend_service_target_deferral`: planning vocabulary for route targeting without implementing route targeting behavior.
- `session_aware_targeting_deferral`: planning vocabulary for route targeting without implementing route targeting behavior.
- `distributed_route_targeting_deferral`: planning vocabulary for route targeting without implementing route targeting behavior.

Forbidden use:

- Do not introduce concrete public API, package, route, method, wire, handler, dashboard, metrics, trace, admin, console, or inspector compatibility names from Pitaya or Nakama.
- Do not use this vocabulary as permission to add route targeting behavior, service registry or selector behavior, heartbeat/liveness behavior, route targeting behavior, remote timeout/retry behavior, distributed session ownership behavior, distributed presence fanout behavior, cross-node error mapping behavior, cluster observability behavior, runtime behavior, protocol messages, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, or distributed runtime behavior.
- Do not classify raw tokens, credentials, lookup digests, verifier digests, verifier keys, DSNs, headers, cookies, query strings, subprotocol values, remote addresses, database payloads, local secret file contents, route payloads, event payloads, or concrete transport metadata as log-safe in this gate.

## 4. Current Mapping Candidate

```yaml
current_source_first_route_targeting_mapping:
  local_route_registry:
    current: runtime/internal/app route key and bootstrap sources
    future_vocabulary: local_route_registry_source
    status: single_process_route_registry_only
  frontend_forwarding:
    current: node tools/vibit inspect pitaya-frontend-forwarding --json
    future_vocabulary: frontend_forwarding_handoff
    status: no_frontend_forwarding_behavior
  backend_service_targets:
    current: deferred
    future_vocabulary: backend_service_target_deferral
    status: no_backend_service_targeting_behavior
  session_aware_targeting:
    current: deferred
    future_vocabulary: session_aware_targeting_deferral
    status: no_session_aware_targeting_behavior
  distributed_route_targeting:
    current: deferred
    future_vocabulary: distributed_route_targeting_deferral
    status: no_distributed_runtime_behavior
```

## 5. Ownership

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-route-targeting-boundary-gate.md
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

- Documentation and manifests may define route targeting vocabulary and current source-first mapping.
- `tools/vibit` may emit a source-first route targeting map because W-0322 is separately recorded.
- Existing runtime, protocol, route, repository, persistence, and transport behavior remain unchanged by this gate.

## 6. Reference Mapping

Nakama remains the primary product reference for broad game backend product capability pressure. Pitaya remains an architecture vocabulary reference for distributed runtime, components, handlers, services, sessions, routes, remote calls, service discovery, groups, lifecycle hooks, and operational concerns.

This gate maps those references into vibit-owned vocabulary only. It does not create direct compatibility, public API parity, runtime behavior, or distributed behavior.

## 7. Stop Conditions

Stop and require a later bounded work item before adding:

- route targeting behavior;
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
node tools/vibit inspect rule runtime.pitaya_aligned_route_targeting_boundary_gate
node tools/vibit check change define-pitaya-aligned-route-targeting-boundary-gate --json
node tools/vibit check runtime --json
node tools/vibit check work --json
```
