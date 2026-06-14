# Pitaya-Aligned Distributed Presence Fanout Boundary Gate 中文版

Status: Accepted v0.1
Last updated: 2026-06-14
Scope: Gate-only boundary for using Pitaya-aligned distributed presence fanout vocabulary after the Pitaya service dispatch source-first sequence
Depends on: `decisions/ADR-0220-pitaya-aligned-cluster-event-bus-source-first-map.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0235`

说明：本文件是 `docs/pitaya-aligned-distributed-presence-fanout-boundary-gate.md` 的简体中文译本。英文版本是权威版本。

This document defines distributed presence fanout vocabulary only. It does not implement distributed presence fanout behavior, service registry behavior, service selector behavior, heartbeat/liveness behavior, route targeting behavior, remote timeout or retry behavior, distributed session ownership behavior, distributed presence fanout behavior, cross-node error mapping behavior, cluster observability behavior, runtime behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

```yaml
pitaya_aligned_distributed_presence_fanout_boundary_gate: defined
completed_work_item: W-0327
decision: ADR-0235
check_rule: runtime.pitaya_aligned_distributed_presence_fanout_boundary_gate
source_pitaya_service_dispatch_map_decision: ADR-0220
standard: docs/pitaya-aligned-distributed-presence-fanout-boundary-gate.md
translation: docs/pitaya-aligned-distributed-presence-fanout-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: distributed_presence_fanout_vocabulary_boundary_defined_for_future_distributed_operations_planning
implementation_scope: gate_only_distributed_presence_fanout_vocabulary
future_implementation_work_item: W-0328
future_implementation_direction: implement_pitaya_aligned_distributed_presence_fanout_source_first_map
allowed_distributed_presence_fanout_vocabulary:
  - distributed_presence_fanout_boundary
  - local_presence_source
  - presence_event_handoff
  - node_fanout_deferral
  - subscription_registry_deferral
  - delivery_guarantee_deferral
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

The purpose is to keep distributed presence fanout visible as a future Pitaya-class distributed operations concern without authorizing implementation. Agents may use this vocabulary for planning, inspection, ADRs, change specs, and future source-first maps only.

## 3. Vocabulary

- `distributed_presence_fanout_boundary`: planning vocabulary for distributed presence fanout without implementing distributed presence fanout behavior.
- `local_presence_source`: planning vocabulary for distributed presence fanout without implementing distributed presence fanout behavior.
- `presence_event_handoff`: planning vocabulary for distributed presence fanout without implementing distributed presence fanout behavior.
- `node_fanout_deferral`: planning vocabulary for distributed presence fanout without implementing distributed presence fanout behavior.
- `subscription_registry_deferral`: planning vocabulary for distributed presence fanout without implementing distributed presence fanout behavior.
- `delivery_guarantee_deferral`: planning vocabulary for distributed presence fanout without implementing distributed presence fanout behavior.

Forbidden use:

- Do not introduce concrete public API, package, route, method, wire, handler, dashboard, metrics, trace, admin, console, or inspector compatibility names from Pitaya or Nakama.
- Do not use this vocabulary as permission to add distributed presence fanout behavior, service registry or selector behavior, heartbeat/liveness behavior, route targeting behavior, remote timeout/retry behavior, distributed session ownership behavior, distributed presence fanout behavior, cross-node error mapping behavior, cluster observability behavior, runtime behavior, protocol messages, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, or distributed runtime behavior.
- Do not classify raw tokens, credentials, lookup digests, verifier digests, verifier keys, DSNs, headers, cookies, query strings, subprotocol values, remote addresses, database payloads, local secret file contents, route payloads, event payloads, or concrete transport metadata as log-safe in this gate.

## 4. Current Mapping Candidate

```yaml
current_source_first_distributed_presence_fanout_mapping:
  local_presence_sources:
    current: runtime/internal/app presence lifecycle and outbound delivery sources
    future_vocabulary: local_presence_source
    status: single_process_presence_only
  presence_events:
    current: runtime/internal/platform/events and presence local proof sources
    future_vocabulary: presence_event_handoff
    status: local_event_source_inspection_only
  node_fanout:
    current: deferred
    future_vocabulary: node_fanout_deferral
    status: no_node_fanout_behavior
  subscription_registry:
    current: deferred
    future_vocabulary: subscription_registry_deferral
    status: no_subscription_registry_behavior
  delivery_guarantees:
    current: deferred
    future_vocabulary: delivery_guarantee_deferral
    status: no_delivery_guarantee_behavior
```

## 5. Ownership

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-distributed-presence-fanout-boundary-gate.md
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

- Documentation and manifests may define distributed presence fanout vocabulary and current source-first mapping.
- `tools/vibit` may emit a source-first distributed presence fanout map because W-0328 is separately recorded.
- Existing runtime, protocol, route, repository, persistence, and transport behavior remain unchanged by this gate.

## 6. Reference Mapping

Nakama remains the primary product reference for broad game backend product capability pressure. Pitaya remains an architecture vocabulary reference for distributed runtime, components, handlers, services, sessions, routes, remote calls, service discovery, groups, lifecycle hooks, and operational concerns.

This gate maps those references into vibit-owned vocabulary only. It does not create direct compatibility, public API parity, runtime behavior, or distributed behavior.

## 7. Stop Conditions

Stop and require a later bounded work item before adding:

- distributed presence fanout behavior;
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
node tools/vibit inspect rule runtime.pitaya_aligned_distributed_presence_fanout_boundary_gate
node tools/vibit check change define-pitaya-aligned-distributed-presence-fanout-boundary-gate --json
node tools/vibit check runtime --json
node tools/vibit check work --json
```
