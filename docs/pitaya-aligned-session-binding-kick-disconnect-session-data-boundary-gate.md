# Pitaya-Aligned Session Binding, Kick/Disconnect, And Session Data Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-01
Scope: Gate-only boundary for using Pitaya-aligned session binding, kick/disconnect, and session data vocabulary after the acceptor and connection lifecycle follow-up direction selection
Depends on: `decisions/ADR-0175-select-next-pitaya-aligned-direction-after-acceptor-connection-lifecycle-map.md`, `decisions/ADR-0174-pitaya-aligned-acceptor-connection-lifecycle-source-first-map.md`, `docs/pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.md`, `docs/runtime-protocol-adapter.md`, `docs/game-protocol.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0176`

The paired Simplified Chinese translation is `docs/pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate.zh-CN.md`. The English file is authoritative.

This document defines a session binding, kick/disconnect, and session data vocabulary gate only. It does not implement session binding behavior, kick/disconnect behavior, session data behavior, session data persistence, acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, backend route targeting, cluster-safe session routing behavior, distributed session routing, distributed runtime behavior, distributed groups, room broadcast fanout, delivery guarantees, stream subscriptions, service discovery implementation, service registries, service selectors, node identity, server-to-server RPC, remote calls, frontend/backend server roles, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The Pitaya-aligned session binding, kick/disconnect, and session data boundary gate record is:

```yaml
pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate: defined
completed_work_item: W-0268
decision: ADR-0176
check_rule: runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate
previous_direction_decision: ADR-0175
acceptor_connection_lifecycle_source_first_map_decision: ADR-0174
acceptor_connection_lifecycle_source_first_map_check_rule: runtime.pitaya_aligned_acceptor_connection_lifecycle_source_first_map
standard: docs/pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate.md
translation: docs/pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: session_binding_kick_disconnect_session_data_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_session_binding_kick_disconnect_session_data_vocabulary
future_implementation_work_item: W-0269
future_implementation_direction: implement_pitaya_aligned_session_binding_kick_disconnect_session_data_source_first_map
allowed_session_binding_kick_disconnect_session_data_vocabulary:
  - session_binding_boundary
  - connection_bound_session
  - session_data
  - session_data_scope
  - server_initiated_disconnect
  - server_initiated_kick
  - session_unbind
  - session_close_reason
  - connection_session_handoff
  - presence_session_handoff
related_vocabulary:
  - first_message_binding
  - runtime_session
  - active_connection_registry
  - logout_service_behavior
  - websocket_close_handoff
  - presence_lifecycle
  - cluster_safe_session_routing
current_single_process_session_binding_kick_disconnect_session_data_mapping:
  first_message_binding:
    current: authentication_bind_connection_route
    future_vocabulary: session_binding_boundary
    implementation_status: no_handshake_authentication_or_reconnect_binding_change
  runtime_session_validation:
    current: request_level_access_token_validation_and_request_identity_handoff
    future_vocabulary: connection_bound_session
    implementation_status: no_session_persistence_or_every_request_policy_change
  session_metadata:
    current: request_identity_and_connection_metadata
    future_vocabulary: session_data
    implementation_status: no_session_data_store_or_public_api
  session_data_scope:
    current: no_general_session_data_scope
    future_vocabulary: session_data_scope
    implementation_status: planning_vocabulary_only
  active_connection_registry:
    current: application_owned_connection_registry
    future_vocabulary: connection_session_handoff
    implementation_status: no_cluster_safe_session_location_registry
  logout_disconnect_handoff:
    current: logout_service_revokes_token_and_transport_close_policy_remains_unchanged
    future_vocabulary: server_initiated_disconnect
    implementation_status: no_server_initiated_disconnect_behavior
  kick_policy:
    current: no_kick_policy_or_route
    future_vocabulary: server_initiated_kick
    implementation_status: planning_vocabulary_only
  session_unbind:
    current: close_handoff_and_connection_registry_cleanup
    future_vocabulary: session_unbind
    implementation_status: no_remote_unbind_or_reconnect_routing
  close_reason:
    current: existing_transport_close_reason_mapping
    future_vocabulary: session_close_reason
    implementation_status: no_close_policy_change
  presence_lifecycle:
    current: server_owned_presence_snapshot
    future_vocabulary: presence_session_handoff
    implementation_status: no_distributed_presence_session_handoff
session_binding_behavior_added: false
kick_disconnect_behavior_added: false
session_data_behavior_added: false
session_data_persistence_added: false
acceptor_behavior_added: false
tcp_acceptor_added: false
websocket_acceptor_behavior_changed: false
connection_lifecycle_behavior_changed: false
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

`ADR-0175` selected session binding, kick/disconnect, and session data vocabulary as the next Pitaya-aligned direction after the acceptor and connection lifecycle source-first map.

The risk is that agents may treat session, kick, disconnect, or session data terms as permission to change first-message binding, introduce socket kick routes, add general session storage, change logout or close behavior, add reconnect routing, or begin distributed session routing. This gate records vocabulary and mapping only. It keeps vibit's concrete single-process WebSocket runtime, request-level authentication, active connection registry, logout behavior, close handoff, and presence lifecycle behavior unchanged.

## 3. Vocabulary

Allowed session binding, kick/disconnect, and session data vocabulary:

- `session_binding_boundary`: future planning vocabulary for binding authenticated runtime identity to a connection. Current binding remains the existing first-message route.
- `connection_bound_session`: future planning vocabulary for a validated session associated with an active connection. Current request identity remains request-level validation and metadata handoff.
- `session_data`: future planning vocabulary for server-owned session metadata. It is not a general persistent data store in this slice.
- `session_data_scope`: future planning vocabulary for eventual limits on session data ownership. No concrete scope is added here.
- `server_initiated_disconnect`: future planning vocabulary for server-originated connection close intent. Current close behavior is unchanged.
- `server_initiated_kick`: future planning vocabulary for policy-driven removal from a session or connection. No kick policy or route is added here.
- `session_unbind`: future planning vocabulary for detaching a session from a connection. Current cleanup remains existing close handoff and registry cleanup.
- `session_close_reason`: future planning vocabulary for classifying server-visible close reasons. Current close reason mapping is unchanged.
- `connection_session_handoff`: future planning vocabulary for passing connection/session association facts inside the application boundary. No distributed owner registry is added here.
- `presence_session_handoff`: future planning vocabulary for connecting session lifecycle facts to presence lifecycle. Current presence remains server-owned snapshot behavior.

Forbidden vocabulary use:

- Do not introduce concrete public API, package, route, method, wire, handler, session, disconnect, kick, registry, selector, or configuration compatibility names from Pitaya or Nakama.
- Do not use session binding, kick/disconnect, or session data vocabulary as permission to add handshake authentication changes, reconnect routing, kick routes, disconnect routes, general session data persistence, protocol messages, generated output, persistence, dependencies, topology, or distributed runtime behavior.
- Do not move authentication validation, request identity, first-message binding, logout, close policy, connection registry, presence lifecycle, or delivery behavior across transport, application, protocol, repository, or startup boundaries.

## 4. Current Mapping

```yaml
current_single_process_session_binding_kick_disconnect_session_data_mapping:
  first_message_binding:
    current: authentication BindConnection route
    future_vocabulary: session_binding_boundary
    status: no_handshake_authentication_or_reconnect_binding_change
  runtime_session_validation:
    current: request-level access token validation and request identity handoff
    future_vocabulary: connection_bound_session
    status: no_session_persistence_or_every_request_policy_change
  session_metadata:
    current: request identity and connection metadata
    future_vocabulary: session_data
    status: no_session_data_store_or_public_api
  active_connection_registry:
    current: application-owned active connection registry
    future_vocabulary: connection_session_handoff
    status: no_cluster_safe_session_location_registry
  logout_disconnect_handoff:
    current: logout service revokes token; transport close policy remains unchanged
    future_vocabulary: server_initiated_disconnect
    status: no_server_initiated_disconnect_behavior
  kick_policy:
    current: no kick policy or route
    future_vocabulary: server_initiated_kick
    status: planning_vocabulary_only
  session_unbind:
    current: close handoff and connection registry cleanup
    future_vocabulary: session_unbind
    status: no_remote_unbind_or_reconnect_routing
  close_reason:
    current: existing transport close reason mapping
    future_vocabulary: session_close_reason
    status: no_close_policy_change
  presence_lifecycle:
    current: server-owned presence snapshot
    future_vocabulary: presence_session_handoff
    status: no_distributed_presence_session_handoff
```

## 5. Ownership

Session binding, kick/disconnect, and session data vocabulary ownership:

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate.md
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

- Documentation and manifests may define session binding, kick/disconnect, and session data vocabulary and current mapping.
- `tools/vibit` may later emit a source-first session binding, kick/disconnect, and session data map if a follow-up implementation work item authorizes it.
- Runtime, transport, protocol, repository, persistence, generated output, startup wiring, dependencies, service discovery, RPC, remote calls, frontend/backend role behavior, cluster-safe session routing, distributed group behavior, room broadcast behavior, and session data persistence remain unchanged by this gate.
- Domain modules do not gain session binding, kick/disconnect, session data, close policy, presence lifecycle, service discovery, RPC, or distributed runtime ownership by default.

## 6. Nakama And Pitaya Mapping

Nakama remains the primary product reference for near-term capability breadth. Pitaya remains an architecture vocabulary reference for acceptors, session binding, session data, route handlers, frontend/backend roles, RPC/remotes, service discovery, groups, broadcast, cluster routing, handler pipelines, serializers, forwarding, and connection lifecycle.

This gate maps those references into vibit-owned vocabulary only. It does not create direct compatibility, public API parity, or runtime behavior.

## 7. Stop Conditions

Stop and require a later bounded work item before adding:

- Session binding behavior, kick/disconnect behavior, session data behavior, or session data persistence.
- WebSocket handshake authentication changes, reconnect routing, remote disconnect, remote kick, route targeting, connection owner registries, or cluster-safe session routing behavior.
- Protocol messages, Protobuf sources, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, metrics/tracing behavior, hosted surfaces, SDKs, release artifacts, or direct Nakama/Pitaya API compatibility.

## 8. Verification

The repository check rule for this boundary is:

```text
runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate
```

Required verification for this gate:

```sh
node -c tools/vibit
node tools/vibit inspect rule runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate
node tools/vibit inspect next --json
node tools/vibit check change define-pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```
