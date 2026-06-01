# Pitaya-Aligned Cluster-Safe Session Routing Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-01
Scope: Gate-only boundary for using Pitaya-aligned cluster-safe session routing vocabulary after the distributed group and broadcast source-first map
Depends on: `decisions/ADR-0163-pitaya-aligned-distributed-group-broadcast-source-first-map.md`, `docs/pitaya-aligned-distributed-group-broadcast-boundary-gate.md`, `docs/first-message-connection-binding-gate.md`, `docs/active-connection-registry-gate.md`, `docs/runtime-session-validation-gate.md`, `docs/logout-revocation-active-connection-gate.md`, `docs/bound-identity-route-policy-gate.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0164`

The paired Simplified Chinese translation is `docs/pitaya-aligned-cluster-safe-session-routing-boundary-gate.zh-CN.md`. The English file is authoritative.

This document defines a cluster-safe session routing vocabulary gate only. It does not implement cluster-safe session routing, session location registries, connection owner node registries, routing epoch behavior, session route target behavior, remote connection handoff, reconnect routing, distributed session routing, distributed runtime behavior, distributed groups, group membership registries, room broadcast fanout, delivery guarantees, stream subscriptions, service discovery implementation, service registries, service selectors, node identity, server-to-server RPC, remote calls, frontend/backend server roles, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The Pitaya-aligned cluster-safe session routing boundary gate record is:

```yaml
pitaya_aligned_cluster_safe_session_routing_boundary_gate: defined
completed_work_item: W-0256
decision: ADR-0164
check_rule: runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate
distributed_group_broadcast_source_first_map_decision: ADR-0163
distributed_group_broadcast_source_first_map_check_rule: runtime.pitaya_aligned_distributed_group_broadcast_source_first_map
standard: docs/pitaya-aligned-cluster-safe-session-routing-boundary-gate.md
translation: docs/pitaya-aligned-cluster-safe-session-routing-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: cluster_safe_session_routing_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_cluster_safe_session_routing_vocabulary
future_implementation_work_item: W-0257
future_implementation_direction: implement_pitaya_aligned_cluster_safe_session_routing_source_first_map
allowed_cluster_safe_session_routing_vocabulary:
  - cluster_safe_session_routing
  - session_location
  - connection_owner_node
  - routing_epoch
  - session_route_target
  - connection_handoff
  - reconnect_route
related_vocabulary:
  - connection_id
  - connection_epoch
  - first_message_connection_binding
  - active_connection_registry
  - runtime_session
  - bound_connection_identity
  - request_token_identity
  - session_validated_identity
  - single_process_connection_binding
  - frontend_server
  - backend_server
  - service_discovery
  - server_to_server_rpc
  - remote_call
  - distributed_group
  - room_broadcast
current_single_process_session_routing_mapping:
  connection_id:
    current: server_observed_connection_id_epoch
    future_vocabulary: connection_owner_node
    implementation_status: current_single_process_connection_binding
  connection_epoch:
    current: server_observed_connection_id_epoch
    future_vocabulary: routing_epoch
    implementation_status: current_single_process_connection_binding
  first_message_connection_binding:
    current: current_single_process_connection_binding
    future_vocabulary: session_location
    implementation_status: active_connection_registry_single_process
  active_connection_registry:
    current: active_connection_registry_single_process
    future_vocabulary: session_location
    implementation_status: no_cross_node_session_location
  runtime_session:
    current: metadata_only_session_id_not_routing_proof
    future_vocabulary: session_route_target
    implementation_status: no_cluster_route_target
  bound_connection_identity:
    current: current_single_process_connection_binding
    future_vocabulary: session_location
    implementation_status: no_distributed_session_routing
  request_token_identity:
    current: request_level_identity_not_cluster_route
    future_vocabulary: session_route_target
    implementation_status: no_cluster_route_target
  session_validated_identity:
    current: request_validation_status_not_cluster_route
    future_vocabulary: session_route_target
    implementation_status: no_cross_node_session_location
  connection_handoff:
    current: no_remote_connection_handoff
    future_vocabulary: connection_handoff
    implementation_status: deferred_future_architecture_reference
  reconnect_route:
    current: reconnect_epoch_local_only_not_cluster_routing
    future_vocabulary: reconnect_route
    implementation_status: deferred_future_architecture_reference
cluster_safe_session_routing_added: false
session_location_registry_added: false
connection_owner_node_registry_added: false
routing_epoch_behavior_added: false
session_route_target_added: false
remote_connection_handoff_added: false
reconnect_route_added: false
distributed_session_routing_added: false
distributed_group_implementation_added: false
distributed_groups_added: false
group_membership_registry_added: false
room_broadcast_fanout_added: false
broadcast_delivery_guarantee_added: false
stream_subscription_added: false
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

`ADR-0163` made distributed group and broadcast vocabulary inspectable through `node tools/vibit inspect pitaya-groups --json`. That map intentionally left cluster-safe session routing as the next deferred Pitaya-aligned planning family. Session routing is high risk because it can be mistaken for permission to add cross-node connection lookup, remote connection handoff, reconnect routing, service discovery, RPC, or transport-level authentication behavior.

This gate records vocabulary and current mapping before any implementation. It maps the current single-process connection binding, active connection registry, runtime session validation, and request-level identity surfaces to future session routing concepts without changing runtime behavior.

## 3. Session Routing Vocabulary

Allowed cluster-safe session routing vocabulary:

- `cluster_safe_session_routing`: future planning vocabulary for routing a validated session or connection target across nodes. It is not behavior in this slice.
- `session_location`: future planning vocabulary for where a validated runtime session or bound connection currently lives. It is not a registry, cache, table, or service discovery record in this slice.
- `connection_owner_node`: future planning vocabulary for the node that owns an open connection. The current runtime has no node ownership registry.
- `routing_epoch`: future planning vocabulary for avoiding stale routing decisions. Current `connection_epoch` remains single-process connection lifecycle metadata.
- `session_route_target`: future planning vocabulary for an application-selected route target. Current session metadata is not a cluster route target.
- `connection_handoff`: future planning vocabulary for moving or delegating connection handling. It is not a remote call, socket migration, or close policy in this slice.
- `reconnect_route`: future planning vocabulary for directing reconnect behavior. Current reconnect and epoch behavior remains local and non-cluster.

Related vocabulary:

- `connection_id` and `connection_epoch`: server-observed connection lifecycle metadata in the current single-process runtime.
- `first_message_connection_binding`: current application/protocol binding posture for associating an authenticated identity with a WebSocket connection.
- `active_connection_registry`: current single-process runtime state vocabulary for server-observed open connections.
- `runtime_session`, `bound_connection_identity`, `request_token_identity`, and `session_validated_identity`: current authentication/session validation concepts; none is a cluster route target in this slice.
- `frontend_server`, `backend_server`, `service_discovery`, `server_to_server_rpc`, `remote_call`, `distributed_group`, and `room_broadcast`: prior Pitaya-aligned vocabulary families that remain deferred implementation concerns.

Forbidden vocabulary use:

- Do not introduce concrete public API, package, route, method, wire, session, channel, registry, selector, handoff, or configuration compatibility names from Pitaya or Nakama.
- Do not use session routing vocabulary as permission to add cross-node registries, session-location tables, node identity, service discovery, RPC, remote calls, transport carriers, protocol messages, generated output, persistence, dependencies, routing caches, handoff workers, reconnect routers, or cluster runtime behavior.
- Do not treat metadata-only `session_id`, client-supplied connection metadata, target-scope metadata, or transport metadata as routing proof.
- Do not bypass application dispatch, authenticated request validation, bound identity route policy, generated output rules, redaction rules, or module ownership.

## 4. Current Mapping

```yaml
current_single_process_session_routing_mapping:
  connection_id:
    current: server-observed connection id and epoch metadata
    future_vocabulary: connection_owner_node
    status: current_single_process_connection_binding
  connection_epoch:
    current: server-observed connection epoch metadata
    future_vocabulary: routing_epoch
    status: current_single_process_connection_binding
  first_message_connection_binding:
    current: application-owned first-message bind posture for one process
    future_vocabulary: session_location
    status: active_connection_registry_single_process
  active_connection_registry:
    current: single-process active connection registry vocabulary
    future_vocabulary: session_location
    status: no_cross_node_session_location
  runtime_session:
    current: session validation metadata; metadata-only session id is not proof
    future_vocabulary: session_route_target
    status: no_cluster_route_target
  bound_connection_identity:
    current: one-process bound identity vocabulary
    future_vocabulary: session_location
    status: no_distributed_session_routing
  request_token_identity:
    current: request-level token identity
    future_vocabulary: session_route_target
    status: no_cluster_route_target
  session_validated_identity:
    current: validation status vocabulary, not a route record
    future_vocabulary: session_route_target
    status: no_cross_node_session_location
  connection_handoff:
    current: no remote connection handoff
    future_vocabulary: connection_handoff
    status: deferred_future_architecture_reference
  reconnect_route:
    current: reconnect epoch behavior is local, not cluster routing
    future_vocabulary: reconnect_route
    status: deferred_future_architecture_reference
```

## 5. Ownership

Cluster-safe session routing vocabulary ownership:

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-cluster-safe-session-routing-boundary-gate.md
  - .arch/reference.yaml
  - .arch/runtime.yaml
source_first_map_candidate_owner:
  - tools/vibit
runtime_behavior_owner: unchanged
transport_owner: unchanged
protocol_owner: unchanged
persistence_owner: unchanged
module_owner: unchanged
dependency_owner: unchanged
```

Rules:

- Documentation and manifests may define cluster-safe session routing vocabulary and mapping.
- `tools/vibit` may later emit a source-first cluster-safe session routing map if a follow-up implementation work item authorizes it.
- Runtime, transport, protocol, repository, persistence, generated output, startup wiring, dependencies, service discovery, RPC, remote calls, frontend/backend role behavior, distributed group behavior, and room broadcast behavior remain unchanged by this gate.
- Domain modules do not gain session routing, connection targeting, reconnect routing, handoff, service discovery, RPC, or transport ownership by default.

## 6. Nakama And Pitaya Mapping

Nakama remains the primary product reference for near-term capability breadth. Pitaya remains an architecture vocabulary reference for acceptors, sessions, route handlers, frontend/backend roles, RPC/remotes, service discovery, groups, broadcast, and cluster routing.

Adopted as vocabulary:

- cluster-safe session routing as future architecture-planning vocabulary;
- session location, connection owner node, routing epoch, session route target, connection handoff, and reconnect route as planning vocabulary;
- mapping current single-process connection binding, active connection registry, and request/session validation vocabulary to deferred cluster-routing concepts.

Adapted to vibit:

- Any future routing model must preserve vibit application-owned identity validation, bound identity route policy, module boundaries, generated output rules, redaction, and repository checks.
- Current single-process runtime, connection binding, active connection registry, and request-level validation remain the concrete implementation.
- Any future cluster-safe session routing implementation must be separately gated and verified before behavior exists.

Rejected for now:

- direct Pitaya or Nakama API compatibility;
- Pitaya or Nakama package, method, session, route, registry, selector, or handoff naming compatibility;
- cluster-safe session routing behavior;
- session location registry, connection owner node registry, routing epoch behavior, session route targets, remote connection handoff, reconnect routing, or distributed session routing;
- service discovery, RPC, remote calls, frontend/backend process split, distributed groups, room broadcast fanout, or distributed runtime behavior;
- protocol messages or routes, generated output, persistence, migrations, dependencies, hosted deployment, SDK publication, or release artifacts for routing.

## 7. Future Implementation Work

Open:

```text
M-185/W-0257 Implement Pitaya-aligned cluster-safe session routing source-first map
```

The future work item may:

- add a source-first repository inspection map for cluster-safe session routing vocabulary;
- summarize current connection id, connection epoch, first-message binding, active connection registry, runtime session validation, bound identity, and request identity mappings;
- update runbooks and acceptance docs to point to the cluster-safe session routing map;
- add repository checks that verify the map remains gate-only and redacted.

The future work item must not:

- add cluster-safe session routing behavior;
- add session location registries or connection owner node registries;
- add routing epoch behavior, session route targets, remote connection handoff, or reconnect routing;
- add service discovery implementation, service registries, selectors, node identity, or topology behavior;
- add server-to-server RPC implementation or remote call behavior;
- add frontend/backend server role implementation;
- add distributed runtime implementation;
- add distributed group implementation, group membership registries, stream subscriptions, room broadcast fanout, or delivery guarantees;
- add runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## 8. Verification Expectations

This gate should verify:

- `runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate` is registered.
- `ADR-0164`, this standard, the Simplified Chinese translation, change artifacts, and conversation memory exist.
- W-0256 is completed and W-0257 is next-ready.
- Current single-process connection/session mapping is recorded.
- Deferrals remain explicit for cluster-safe routing behavior, registries, handoff, distributed runtime, service discovery, RPC, remote calls, protocol, generated output, persistence, dependencies, hosted surfaces, SDKs, and direct compatibility.

Required commands:

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate
node tools/vibit check change define-pitaya-aligned-cluster-safe-session-routing-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

## 9. Stop Conditions

Stop and require a later bounded work item before adding:

- cluster-safe session routing behavior;
- session location registry or connection owner node registry behavior;
- routing epoch behavior, route target resolution, remote connection handoff, or reconnect route behavior;
- service discovery, service registry, service selector, node identity, server identity, RPC, remote calls, or frontend/backend roles;
- distributed runtime behavior, distributed groups, group membership registries, room broadcast fanout, delivery guarantees, or stream subscriptions;
- runtime endpoint behavior, protocol messages, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted surfaces, SDKs, or direct compatibility.

## 10. Non-Authorization

This is a boundary-only standard. It authorizes vocabulary, mapping, manifests, checks, ADRs, and memory only. It does not authorize runtime behavior, protocol behavior, generated output, persistence, service discovery, RPC, remote calls, distributed runtime behavior, hosted deployment, SDK publication, release execution, or direct Nakama/Pitaya API compatibility.
