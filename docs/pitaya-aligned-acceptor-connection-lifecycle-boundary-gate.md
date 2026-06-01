# Pitaya-Aligned Acceptor And Connection Lifecycle Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-01
Scope: Gate-only boundary for using Pitaya-aligned acceptor and connection lifecycle vocabulary after the serializer and message forwarding follow-up direction selection
Depends on: `decisions/ADR-0172-select-next-pitaya-aligned-direction-after-serializer-message-forwarding-map.md`, `decisions/ADR-0171-pitaya-aligned-serializer-message-forwarding-source-first-map.md`, `docs/pitaya-aligned-serializer-message-forwarding-boundary-gate.md`, `docs/runtime-protocol-adapter.md`, `docs/game-protocol.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0173`

The paired Simplified Chinese translation is `docs/pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.zh-CN.md`. The English file is authoritative.

This document defines an acceptor and connection lifecycle vocabulary gate only. It does not implement acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, session binding behavior, kick/disconnect behavior, serializer behavior, message forwarding behavior, route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, backend route targeting, cluster-safe session routing behavior, distributed session routing, distributed runtime behavior, distributed groups, room broadcast fanout, delivery guarantees, stream subscriptions, service discovery implementation, service registries, service selectors, node identity, server-to-server RPC, remote calls, frontend/backend server roles, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The Pitaya-aligned acceptor and connection lifecycle boundary gate record is:

```yaml
pitaya_aligned_acceptor_connection_lifecycle_boundary_gate: defined
completed_work_item: W-0265
decision: ADR-0173
check_rule: runtime.pitaya_aligned_acceptor_connection_lifecycle_boundary_gate
previous_direction_decision: ADR-0172
serializer_message_forwarding_source_first_map_decision: ADR-0171
serializer_message_forwarding_source_first_map_check_rule: runtime.pitaya_aligned_serializer_message_forwarding_source_first_map
standard: docs/pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.md
translation: docs/pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: acceptor_connection_lifecycle_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_acceptor_connection_lifecycle_vocabulary
future_implementation_work_item: W-0266
future_implementation_direction: implement_pitaya_aligned_acceptor_connection_lifecycle_source_first_map
allowed_acceptor_connection_lifecycle_vocabulary:
  - acceptor_boundary
  - websocket_acceptor
  - connection_id
  - connection_epoch
  - session_binding
  - active_connection_registry
  - close_handoff
  - presence_lifecycle_handoff
related_vocabulary:
  - first_message_binding
  - runtime_session
  - route_request
  - server_push_delivery
  - cluster_safe_session_routing
  - message_forwarding
current_single_process_acceptor_connection_mapping:
  websocket_acceptor:
    current: single_process_websocket_server_accept_loop
    future_vocabulary: acceptor_boundary
    implementation_status: no_tcp_acceptor_or_distributed_acceptor
  connection_identity:
    current: server_observed_connection_id
    future_vocabulary: connection_id
    implementation_status: local_process_metadata
  connection_epoch:
    current: server_observed_connection_epoch
    future_vocabulary: connection_epoch
    implementation_status: no_distributed_routing_epoch
  first_message_binding:
    current: authentication_bind_connection_route
    future_vocabulary: session_binding
    implementation_status: no_handshake_authentication_or_reconnect_binding
  active_connection_registry:
    current: application_owned_connection_registry
    future_vocabulary: active_connection_registry
    implementation_status: no_connection_owner_node_registry
  close_handoff:
    current: transport_close_to_application_policy
    future_vocabulary: close_handoff
    implementation_status: no_remote_disconnect_handoff
  presence_lifecycle:
    current: server_owned_presence_snapshot
    future_vocabulary: presence_lifecycle_handoff
    implementation_status: no_distributed_presence_lifecycle
acceptor_behavior_added: false
tcp_acceptor_added: false
websocket_acceptor_behavior_changed: false
connection_lifecycle_behavior_changed: false
session_binding_behavior_added: false
kick_disconnect_behavior_added: false
concrete_socket_close_behavior_changed: false
serializer_behavior_added: false
message_forwarding_behavior_added: false
route_handler_implementation_added: false
handler_routing_behavior_added: false
handler_pipeline_behavior_added: false
pipeline_middleware_behavior_added: false
backend_route_targeting_added: false
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
metrics_endpoint_added: false
tracing_pipeline_added: false
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

`ADR-0172` selected acceptor and connection lifecycle vocabulary as the next Pitaya-aligned direction after the serializer and message forwarding source-first map.

The risk is that agents may treat acceptor, session binding, kick, disconnect, or lifecycle terms as permission to add TCP acceptors, change WebSocket transport behavior, add handshake authentication, alter connection close semantics, add remote disconnect behavior, or introduce distributed session routing. This gate records vocabulary and mapping only. It keeps vibit's concrete single-process WebSocket accept loop, first-message binding, connection registry, close handoff, and presence lifecycle behavior unchanged.

## 3. Vocabulary

Allowed acceptor and connection lifecycle vocabulary:

- `acceptor_boundary`: future planning vocabulary for the boundary that accepts client connections. Current ownership remains the existing WebSocket transport.
- `websocket_acceptor`: future planning vocabulary for the first accepted transport family. It is not permission to change transport behavior in this slice.
- `connection_id`: future planning vocabulary for server-observed connection identity. Current ids remain local metadata.
- `connection_epoch`: future planning vocabulary for connection generation metadata. Current epochs are not distributed routing epochs.
- `session_binding`: future planning vocabulary for binding an authenticated runtime session to a connection. Current binding remains the existing first-message route behavior.
- `active_connection_registry`: future planning vocabulary for tracking active connections. Current state remains application-owned and single-process.
- `close_handoff`: future planning vocabulary for passing close facts from transport to application lifecycle policy. Current close behavior remains unchanged.
- `presence_lifecycle_handoff`: future planning vocabulary for connecting connection lifecycle facts to presence lifecycle. Current presence remains server-owned snapshot behavior.

Forbidden vocabulary use:

- Do not introduce concrete public API, package, route, method, wire, handler, acceptor, session, disconnect, registry, selector, or configuration compatibility names from Pitaya or Nakama.
- Do not use acceptor or connection lifecycle vocabulary as permission to add TCP acceptors, WebSocket behavior changes, session binding behavior, kick/disconnect behavior, remote connection handoff, reconnect routing, protocol messages, generated output, persistence, dependencies, topology, or distributed runtime behavior.
- Do not move session validation, first-message binding, close policy, presence lifecycle, or delivery behavior across transport, application, protocol, repository, or startup boundaries.

## 4. Current Mapping

```yaml
current_single_process_acceptor_connection_mapping:
  websocket_acceptor:
    current: single-process WebSocket server accept loop
    future_vocabulary: acceptor_boundary
    status: no_tcp_acceptor_or_distributed_acceptor
  connection_identity:
    current: server-observed connection id
    future_vocabulary: connection_id
    status: local_process_metadata
  connection_epoch:
    current: server-observed connection epoch
    future_vocabulary: connection_epoch
    status: no_distributed_routing_epoch
  first_message_binding:
    current: authentication BindConnection route
    future_vocabulary: session_binding
    status: no_handshake_authentication_or_reconnect_binding
  active_connection_registry:
    current: application-owned connection registry
    future_vocabulary: active_connection_registry
    status: no_connection_owner_node_registry
  close_handoff:
    current: transport close to application policy handoff
    future_vocabulary: close_handoff
    status: no_remote_disconnect_handoff
  presence_lifecycle:
    current: server-owned presence snapshot
    future_vocabulary: presence_lifecycle_handoff
    status: no_distributed_presence_lifecycle
```

## 5. Ownership

Acceptor and connection lifecycle vocabulary ownership:

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.md
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

- Documentation and manifests may define acceptor and connection lifecycle vocabulary and current mapping.
- `tools/vibit` may later emit a source-first acceptor and connection lifecycle map if a follow-up implementation work item authorizes it.
- Runtime, transport, protocol, repository, persistence, generated output, startup wiring, dependencies, service discovery, RPC, remote calls, frontend/backend role behavior, cluster-safe session routing, distributed group behavior, and room broadcast behavior remain unchanged by this gate.
- Domain modules do not gain acceptor, transport, session binding, disconnect, close policy, presence lifecycle, service discovery, RPC, or distributed runtime ownership by default.

## 6. Nakama And Pitaya Mapping

Nakama remains the primary product reference for near-term capability breadth. Pitaya remains an architecture vocabulary reference for acceptors, sessions, route handlers, frontend/backend roles, RPC/remotes, service discovery, groups, broadcast, cluster routing, handler pipelines, serializers, forwarding, and connection lifecycle.

Adopted as vocabulary:

- acceptor and connection lifecycle vocabulary for future architecture planning;
- current single-process WebSocket accept loop, connection metadata, binding, registry, close handoff, and presence lifecycle mapping;
- explicit deferral language for future source-first inspection work.

Not adopted:

- direct Nakama or Pitaya API compatibility;
- concrete TCP acceptors;
- concrete session binding behavior changes;
- concrete kick/disconnect behavior;
- distributed connection owner registries;
- distributed session routing, reconnect routing, or remote handoff behavior;
- metrics, tracing, dashboards, hosted surfaces, SDKs, or release artifacts.

## 7. Verification

Required repository checks:

```bash
node -c tools/vibit
node tools/vibit inspect rule runtime.pitaya_aligned_acceptor_connection_lifecycle_boundary_gate
node tools/vibit check change define-pitaya-aligned-acceptor-connection-lifecycle-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

## 8. Stop Conditions

Stop and ask before:

- adding or changing acceptor behavior;
- adding TCP acceptors;
- changing WebSocket behavior;
- changing connection lifecycle behavior;
- changing session binding behavior;
- adding kick/disconnect behavior;
- changing protocol messages or routes;
- adding generated output;
- changing repositories, PostgreSQL adapters, migrations, or dependencies;
- adding metrics endpoints, tracing pipelines, dashboards, hosted deployment, SDK publication, release artifacts, or direct Nakama/Pitaya API compatibility.
