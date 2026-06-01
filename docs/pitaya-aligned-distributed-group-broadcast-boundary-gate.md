# Pitaya-Aligned Distributed Group And Broadcast Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-01
Scope: Gate-only boundary for using Pitaya-aligned distributed group and room broadcast vocabulary after the service discovery source-first map
Depends on: `decisions/ADR-0161-pitaya-aligned-service-discovery-source-first-map.md`, `docs/pitaya-aligned-service-discovery-boundary-gate.md`, `docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md`, `docs/first-server-push-realtime-messaging-gate.md`, `docs/realtime-protocol-websocket-outbound-delivery-gate.md`, `docs/game-protocol.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0162`

The paired Simplified Chinese translation is `docs/pitaya-aligned-distributed-group-broadcast-boundary-gate.zh-CN.md`. The English file is authoritative.

This document defines a distributed group and broadcast vocabulary gate only. It does not implement distributed groups, room broadcast fanout, delivery guarantees, stream subscriptions, group membership registries, groups, parties, chat rooms, matchmaking, match runtime, service discovery implementation, service registries, service selectors, node identity, server-to-server RPC, remote calls, frontend/backend server roles, distributed runtime behavior, cluster-safe session routing, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The Pitaya-aligned distributed group and broadcast boundary gate record is:

```yaml
pitaya_aligned_distributed_group_broadcast_boundary_gate: defined
completed_work_item: W-0254
decision: ADR-0162
check_rule: runtime.pitaya_aligned_distributed_group_broadcast_boundary_gate
service_discovery_source_first_map_decision: ADR-0161
service_discovery_source_first_map_check_rule: runtime.pitaya_aligned_service_discovery_source_first_map
standard: docs/pitaya-aligned-distributed-group-broadcast-boundary-gate.md
translation: docs/pitaya-aligned-distributed-group-broadcast-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: distributed_group_broadcast_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_distributed_group_broadcast_vocabulary
future_implementation_work_item: W-0255
future_implementation_direction: implement_pitaya_aligned_distributed_group_broadcast_source_first_map
allowed_group_broadcast_vocabulary:
  - distributed_group
  - room_broadcast
  - broadcast_target
  - group_membership
  - broadcast_fanout
related_vocabulary:
  - target_scope
  - server_push_intent
  - route_handler
  - module_handler
  - frontend_server
  - backend_server
  - service_discovery
  - server_to_server_rpc
  - remote_call
  - cluster_safe_session_routing
current_single_process_group_broadcast_mapping:
  target_scope:
    current: protocol_envelope_target_metadata_only
    future_vocabulary: broadcast_target
    implementation_status: current_single_process_intent_only
  server_push_intent:
    current: application_owned_realtime_outbound_intent_single_process_delivery
    future_vocabulary: room_broadcast
    implementation_status: current_single_process_delivery_only
  distributed_group:
    current: no_distributed_group_current_single_process_target_scope_only
    future_vocabulary: distributed_group
    implementation_status: deferred_future_architecture_reference
  group_membership:
    current: no_group_membership_registry_or_subscription_state
    future_vocabulary: group_membership
    implementation_status: deferred_future_architecture_reference
  room_broadcast:
    current: no_room_broadcast_fanout_current_server_push_intent_only
    future_vocabulary: room_broadcast
    implementation_status: deferred_future_architecture_reference
  broadcast_target:
    current: target_scope_values_are_metadata_not_distributed_targets
    future_vocabulary: broadcast_target
    implementation_status: current_metadata_only
  broadcast_fanout:
    current: no_cluster_fanout_no_delivery_guarantee
    future_vocabulary: broadcast_fanout
    implementation_status: deferred_future_architecture_reference
distributed_group_implementation_added: false
distributed_groups_added: false
group_membership_registry_added: false
room_broadcast_fanout_added: false
broadcast_delivery_guarantee_added: false
stream_subscription_added: false
groups_parties_chat_runtime_behavior_added: false
match_runtime_behavior_added: false
service_discovery_implementation_added: false
service_registry_added: false
service_selector_added: false
node_registry_added: false
server_identity_added: false
server_to_server_rpc_implementation_added: false
remote_call_behavior_added: false
frontend_server_implementation_added: false
backend_server_implementation_added: false
frontend_backend_server_roles_added: false
distributed_runtime_implementation_added: false
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

`ADR-0161` made service discovery vocabulary inspectable through `node tools/vibit inspect pitaya-discovery --json`. That map intentionally left distributed groups and room broadcast as deferred follow-up vocabulary. The next high-risk Pitaya concept is broadcast: it can be mistaken for permission to add group membership state, room routing, fanout workers, delivery guarantees, stream subscriptions, or cluster-safe session routing.

This gate records the vocabulary and current mapping before any implementation. It maps the current single-process target-scope and server-push intent surfaces to future distributed group and broadcast concepts without changing runtime behavior.

## 3. Group And Broadcast Vocabulary

Allowed group and broadcast vocabulary:

- `distributed_group`: a future architecture-planning term for a cluster-aware grouping target. It is not a current group, party, chat room, match, or stream implementation.
- `room_broadcast`: a future broadcast vocabulary item. It is not a delivery guarantee, wire route, or fanout implementation in this slice.
- `broadcast_target`: a future targeting vocabulary item for fanout planning. Current `target_scope` values remain protocol metadata and do not identify distributed targets.
- `group_membership`: a future membership vocabulary item. It is not a registry, subscription table, persistence model, or runtime state in this slice.
- `broadcast_fanout`: a future delivery topology vocabulary item. It is not a worker, queue, retry policy, ordering policy, or cluster mechanism in this slice.

Related vocabulary:

- `target_scope`: current protocol target metadata defined by the game protocol and envelope posture.
- `server_push_intent`: current application-owned outbound realtime intent and single-process delivery posture.
- `route_handler` and `module_handler`: current in-process application/module behavior that may later produce broadcast intent.
- `frontend_server`, `backend_server`, `service_discovery`, `server_to_server_rpc`, `remote_call`, and `cluster_safe_session_routing`: prior Pitaya-aligned vocabulary families that remain deferred implementation concerns.

Forbidden vocabulary use:

- Do not introduce concrete public API, package, route, method, wire, group, room, channel, stream, or configuration compatibility names from Pitaya or Nakama.
- Do not use group or broadcast vocabulary as permission to add membership registries, room state, stream subscriptions, fanout workers, delivery guarantees, retries, ordering, durable offsets, queueing, cluster routing, service discovery, RPC, remote calls, protocol messages, generated output, persistence, or dependencies.
- Do not use future broadcast vocabulary to bypass module contracts, application dispatch boundaries, authentication/session validation gates, permission checks, generated output rules, redaction rules, or repository ownership.

## 4. Current Mapping

```yaml
current_single_process_group_broadcast_mapping:
  target_scope:
    current: Protobuf envelope target metadata and game-protocol target vocabulary
    future_vocabulary: broadcast_target
    status: current_single_process_intent_only
  server_push_intent:
    current: application-owned realtime outbound intent with single-process WebSocket delivery
    future_vocabulary: room_broadcast
    status: current_single_process_delivery_only
  distributed_group:
    current: no distributed group model; current target scope is metadata only
    future_vocabulary: distributed_group
    status: deferred_future_architecture_reference
  group_membership:
    current: no membership registry, subscription table, or distributed group state
    future_vocabulary: group_membership
    status: deferred_future_architecture_reference
  room_broadcast:
    current: no room broadcast fanout; current server push is narrow single-process delivery
    future_vocabulary: room_broadcast
    status: deferred_future_architecture_reference
  broadcast_target:
    current: target scope values are not distributed routing targets
    future_vocabulary: broadcast_target
    status: current_metadata_only
  broadcast_fanout:
    current: no fanout worker, queue, retry, ordering, or delivery guarantee
    future_vocabulary: broadcast_fanout
    status: deferred_future_architecture_reference
```

## 5. Ownership

Group and broadcast vocabulary ownership:

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-distributed-group-broadcast-boundary-gate.md
  - .arch/reference.yaml
  - .arch/runtime.yaml
source_first_map_candidate_owner:
  - tools/vibit
runtime_behavior_owner: unchanged
protocol_owner: unchanged
persistence_owner: unchanged
module_owner: unchanged
dependency_owner: unchanged
```

Rules:

- Documentation and manifests may define distributed group and broadcast vocabulary and mapping.
- `tools/vibit` may later emit a source-first distributed group and broadcast map if a follow-up implementation work item authorizes it.
- Runtime, protocol, repository, persistence, generated output, startup wiring, dependencies, RPC, remote calls, service discovery, frontend/backend role behavior, and cluster-safe session routing remain unchanged by this gate.
- Domain modules do not gain group, party, chat, stream, match, or broadcast ownership by default. Module contracts remain the source of module behavior and data ownership.

## 6. Nakama And Pitaya Mapping

Nakama remains the primary product reference for near-term capability breadth. Pitaya remains an architecture vocabulary reference for acceptors, sessions, route handlers, frontend/backend roles, RPC/remotes, service discovery, groups, broadcast, and cluster routing.

Adopted as vocabulary:

- distributed group as future cluster-aware target vocabulary;
- room broadcast as future fanout vocabulary;
- broadcast target, group membership, and broadcast fanout as planning vocabulary;
- mapping current target-scope metadata and server-push intent to deferred distributed group and broadcast concepts.

Adapted to vibit:

- Any future group or broadcast model must preserve vibit module ownership, application dispatch boundaries, server-authoritative validation, generated output rules, redaction, and repository checks.
- Current single-process runtime and narrow realtime outbound delivery remain the concrete implementation.
- Any future distributed group or room broadcast implementation must be separately gated and verified before behavior exists.

Rejected for now:

- direct Pitaya or Nakama API compatibility;
- Pitaya or Nakama package, method, group, room, stream, or route naming compatibility;
- distributed group implementation;
- group membership registry or subscription state;
- room broadcast fanout;
- delivery guarantees, retries, ordering, durable offsets, or queueing;
- protocol messages or routes for groups or broadcast;
- service discovery, RPC, remote calls, frontend/backend process split, or cluster-safe session routing.

## 7. Future Implementation Work

Open:

```text
M-183/W-0255 Implement Pitaya-aligned distributed group and broadcast source-first map
```

The future work item may:

- add a source-first repository inspection map for distributed group and broadcast vocabulary;
- summarize current target-scope metadata, server-push intent, and single-process delivery mapping;
- update runbooks and acceptance docs to point to the distributed group and broadcast map;
- add repository checks that verify the map remains gate-only and redacted.

The future work item must not:

- add distributed group implementation;
- add room broadcast fanout;
- add delivery guarantees, retries, ordering, durable offsets, queueing, or backpressure behavior;
- add stream subscriptions, group membership registries, groups, parties, chat rooms, matchmaking, or match runtime behavior;
- add service discovery implementation, service registries, selectors, node identity, or topology behavior;
- add server-to-server RPC implementation;
- add remote call behavior;
- add frontend/backend server role implementation;
- add distributed runtime implementation;
- add cluster-safe session routing behavior;
- add runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## 8. Verification Expectations

This gate should verify:

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_distributed_group_broadcast_boundary_gate
node tools/vibit check change define-pitaya-aligned-distributed-group-broadcast-boundary-gate --json
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

- distributed group implementation;
- room broadcast fanout behavior;
- group membership registry, subscription, stream, chat, party, group, room, matchmaking, or match runtime behavior;
- delivery guarantees, retries, ordering, durable offsets, queueing, or backpressure behavior;
- service discovery implementation, registry, selector, membership, heartbeat, or node identity behavior;
- server-to-server RPC or remote call behavior;
- process topology changes;
- new goroutines, listeners, network protocols, or server roles;
- cluster-safe session routing behavior;
- protocol or Protobuf changes;
- repository, adapter, migration, dependency, or generated-output changes;
- public API compatibility with Pitaya or Nakama;
- any sensitive runtime state exposure.
