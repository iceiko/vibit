# Pitaya-Aligned Distributed Runtime Vocabulary Reactivation Gate

Status: Accepted v0.1
Last updated: 2026-05-31
Scope: Gate-only boundary for reactivating Pitaya-aligned distributed runtime vocabulary after source-first operations inspection
Depends on: `decisions/ADR-0153-minimum-operations-inspection-source-first-surface-implementation.md`, `docs/reference-game-server-alignment.md`, `docs/nakama-pitaya-product-parity-roadmap.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0154`

The paired Simplified Chinese translation is `docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.zh-CN.md`. The English file is authoritative.

This document defines a vocabulary gate only. It does not implement distributed runtime behavior, frontend/backend server roles, server-to-server RPC, remote calls, service discovery, distributed groups, broadcast fanout, cluster-safe session routing, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The Pitaya-aligned distributed runtime vocabulary reactivation gate record is:

```yaml
pitaya_aligned_distributed_runtime_vocabulary_reactivation_gate: defined
completed_work_item: W-0246
decision: ADR-0154
check_rule: runtime.pitaya_aligned_distributed_runtime_vocabulary_reactivation_gate
source_operations_inspection_decision: ADR-0153
source_operations_inspection_check_rule: runtime.minimum_operations_inspection_source_first_surface_implementation
standard: docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md
translation: docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: vocabulary_reactivated_for_future_architecture_planning
implementation_scope: gate_only_architecture_vocabulary
future_implementation_work_item: W-0247
future_implementation_direction: implement_pitaya_aligned_distributed_runtime_vocabulary_source_first_map
allowed_vocabulary:
  - acceptor
  - frontend_server
  - backend_server
  - route_handler
  - session_binding
  - server_to_server_rpc
  - remote_call
  - service_discovery
  - distributed_group
  - room_broadcast
  - cluster_safe_session_routing
current_single_process_mapping:
  websocket_tcp_acceptors: current_websocket_acceptor_single_process_only
  session_binding: current_first_message_connection_binding_single_process_only
  route_handler_model: current_application_dispatch_and_protocol_bridge
  frontend_backend_server_roles: deferred_future_architecture_reference
  rpc_and_remote_calls: deferred_future_architecture_reference
  groups_rooms_broadcast: deferred_future_architecture_reference
  cluster_service_discovery: deferred_future_architecture_reference
distributed_runtime_implementation_added: false
frontend_backend_server_roles_added: false
server_to_server_rpc_added: false
service_discovery_added: false
distributed_groups_added: false
cluster_safe_session_routing_added: false
runtime_behavior_added: false
runtime_endpoint_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
migration_added: false
dependency_added: false
authentication_session_behavior_changed: false
hosted_deployment_added: false
sdk_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Purpose

`ADR-0153` made Pitaya architecture pressure visible in the source-first operations inspection output. That was useful but still embedded in an operations surface. The next safe step is to define a standalone vocabulary gate so future agents can talk about distributed runtime concepts without accidentally implementing them.

The gate reactivates Pitaya as vocabulary for future architecture planning only. It does not make Pitaya the primary product capability driver. Nakama remains the primary product reference for near-term capability breadth. Pitaya vocabulary is allowed when the discussion is about transport acceptors, session binding, route handling, frontend/backend roles, RPC/remotes, service discovery, distributed groups, broadcast, and cluster-safe routing.

## 3. Vocabulary

Allowed vocabulary:

- `acceptor`: a future abstraction for client connection acceptors such as WebSocket or TCP. Current vibit state remains the single-process WebSocket acceptor.
- `frontend_server`: a future role that owns client-facing acceptors and session ingress. Current vibit state has no frontend/backend role split.
- `backend_server`: a future role that owns backend services or domain handlers behind a frontend server. Current vibit state is a modular monolith.
- `route_handler`: the current application dispatch plus protocol bridge maps closest to Pitaya handler routing, but vibit keeps route contracts and module ownership.
- `session_binding`: the current first-message connection binding is the single-process predecessor of any future cluster-safe binding semantics.
- `server_to_server_rpc`: a future server-to-server call family. It must not bypass module contracts if introduced later.
- `remote_call`: a future distributed call vocabulary item, distinct from public client protocol routes.
- `service_discovery`: a future cluster membership and service lookup concern.
- `distributed_group`: a future cluster-aware group, room, party, stream, or match broadcast target concern.
- `room_broadcast`: a future broadcast vocabulary item. It does not authorize delivery guarantees, fanout code, or protocol messages.
- `cluster_safe_session_routing`: future routing semantics that keep player/session/connection ownership coherent across nodes.

Forbidden vocabulary use:

- Do not name concrete public APIs after Pitaya APIs.
- Do not add Pitaya package, namespace, route, method, or wire compatibility markers.
- Do not use vocabulary reactivation as permission to add implementation code.

## 4. Current Mapping

```yaml
current_single_process_mapping:
  websocket_tcp_acceptors:
    current: runtime/cmd/vibit-server WebSocket acceptor
    future_vocabulary: acceptor
    implementation_status: current_websocket_acceptor_single_process_only
  session_binding:
    current: first-message connection binding
    future_vocabulary: session_binding and cluster_safe_session_routing
    implementation_status: current_first_message_connection_binding_single_process_only
  route_handler_model:
    current: Protobuf bridge plus application dispatch and module handlers
    future_vocabulary: route_handler
    implementation_status: current_application_dispatch_and_protocol_bridge
  frontend_backend_server_roles:
    current: none
    future_vocabulary: frontend_server and backend_server
    implementation_status: deferred_future_architecture_reference
  rpc_and_remote_calls:
    current: none
    future_vocabulary: server_to_server_rpc and remote_call
    implementation_status: deferred_future_architecture_reference
  groups_rooms_broadcast:
    current: target scope vocabulary exists; distributed groups do not
    future_vocabulary: distributed_group and room_broadcast
    implementation_status: deferred_future_architecture_reference
  cluster_service_discovery:
    current: none
    future_vocabulary: service_discovery
    implementation_status: deferred_future_architecture_reference
```

## 5. Ownership

Vocabulary ownership:

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md
  - .arch/reference.yaml
  - .arch/runtime.yaml
source_first_map_candidate_owner:
  - tools/vibit
runtime_behavior_owner: unchanged
protocol_owner: unchanged
persistence_owner: unchanged
module_owner: unchanged
```

Rules:

- Documentation and manifests may define vocabulary and mapping.
- `tools/vibit` may later emit a source-first vocabulary map if a follow-up implementation work item authorizes it.
- Runtime, protocol, repository, persistence, generated output, startup wiring, and dependencies remain unchanged by this gate.
- No game module owns distributed runtime vocabulary by default.

## 6. Nakama And Pitaya Mapping

Nakama remains the primary product reference:

- The current product path still prioritizes common backend capability coverage, local alpha clarity, route proof, and prototype usefulness.
- This gate does not select a new social, realtime, matchmaking, match runtime, SDK, or operations implementation slice.

Pitaya is reactivated only as future architecture vocabulary:

- Adopted as vocabulary: acceptors, session binding, route handlers, frontend/backend roles, server-to-server RPC, remotes, service discovery, groups, broadcast, and cluster routing.
- Adapted to vibit: vocabulary must stay behind contract-first, module-owned, repository-checkable boundaries.
- Rejected for now: cluster runtime implementation, direct Pitaya API compatibility, public route naming compatibility, package namespace compatibility, and any bypass of vibit module contracts.

## 7. Future Implementation Work

Open:

```text
M-175/W-0247 Implement Pitaya-aligned distributed runtime vocabulary source-first map
```

The future work item may:

- add a source-first `tools/vibit inspect pitaya-vocabulary` command or equivalent repository inspection;
- summarize the vocabulary and current single-process mapping defined by this gate;
- update runbooks and acceptance docs to point to the vocabulary map;
- add repository checks that verify the map remains gate-only and redacted.

The future work item must not:

- add distributed runtime implementation;
- add frontend/backend server role implementation;
- add server-to-server RPC or remote call behavior;
- add service discovery;
- add distributed groups, broadcast fanout, or delivery guarantees;
- add cluster-safe session routing behavior;
- add runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## 8. Verification Expectations

This gate should verify:

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_distributed_runtime_vocabulary_reactivation_gate
node tools/vibit check change define-pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

This gate does not require Go tests by itself because it adds no Go runtime behavior.

## 9. Stop Conditions

Stop and create a separate gate if the work requires:

- process topology changes;
- new goroutines, listeners, network protocols, or server roles;
- RPC/remoting behavior;
- service discovery dependencies;
- distributed group or broadcast behavior;
- cluster-safe session routing behavior;
- protocol or Protobuf changes;
- repository, adapter, migration, or generated-output changes;
- public API compatibility with Pitaya or Nakama;
- any sensitive runtime state exposure.
