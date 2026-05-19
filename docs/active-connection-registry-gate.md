# Active Connection Registry Gate

Status: Draft v0.1
Last updated: 2026-05-18
Scope: Gate-only boundary for future active WebSocket connection registry behavior
Depends on: `docs/logout-revocation-active-connection-gate.md`, `docs/logout-access-token-behavior-gate.md`, `decisions/ADR-0073-logout-access-token-behavior-implementation.md`, `docs/first-message-connection-binding-implementation-gate.md`, `docs/bound-identity-route-policy-gate.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0074`

The paired Simplified Chinese translation is `docs/active-connection-registry-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

`LogoutAccessToken` can now revoke the verified presented opaque access-token record, but that revocation does not close already-open WebSocket connections, revoke runtime sessions, or mark connection-bound identity unusable. That is intentional. The next missing boundary is the registry required before vibit can reason about active connections as server-owned runtime state.

Nakama makes the product pressure clear: authenticated session material, realtime sockets, logout, expiration, and revocation are related lifecycle concerns. Pitaya makes the layering pressure clear: acceptors, sessions, route handlers, and connection management are distinct surfaces, and handlers should not parse credentials or own transport lifecycle side effects.

vibit adapts those lessons by requiring an application-owned active connection registry boundary before any code can target open sockets by player, runtime session, token record, connection id, or epoch. This standard is gate-only.

```yaml
active_connection_registry_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0160
decision: ADR-0074
check_rule: runtime.active_connection_registry_gate
future_registry_owner: runtime/internal/app/connection
future_policy_owner: runtime/internal/app
future_transport_handoff_owner: runtime/internal/platform/transport/ws
future_protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
first_registry_posture: single_process_in_memory
registry_persistence_posture: non_durable_runtime_state
cluster_registry_posture: deferred
server_observed_connection_id_required: true
connection_epoch_required: true
metadata_only_targeting_allowed: false
bind_connection_integration_candidate: runtime.authentication.BindConnection
logout_revocation_integration_candidate: future_application_policy_only
active_connection_registry_added: false
active_connection_invalidation_added: false
websocket_close_policy_added: false
kick_disconnect_behavior_added: false
duplicate_connection_replacement_added: false
reconnect_epoch_behavior_added: false
runtime_session_revocation_added: false
protobuf_logout_route_added: false
protobuf_session_carrier_added: false
existing_protobuf_envelope_change_added: false
transport_credential_carrier_added: false
websocket_handshake_authentication_added: false
dependencies_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

This standard does not add a registry implementation.

## 2. Ownership

Future active connection registry behavior must stay application-owned:

```yaml
future_registry_owner: runtime/internal/app/connection
future_policy_owner: runtime/internal/app
future_connection_binding_caller: runtime/internal/platform/protocol/protobuf
future_transport_handoff_owner: runtime/internal/platform/transport/ws
authentication_service_owner: runtime/internal/app/authentication
session_repository_owner: runtime/internal/app/session
domain_handler_owner: runtime/internal/modules/*
```

Rules:

- The application layer owns active connection registry state and lifecycle policy.
- WebSocket transport may report server-observed connection open/close facts and may close concrete sockets only through a future narrow transport handoff.
- WebSocket transport must not own player identity, session identity, token validity, logout policy, route policy, or credential parsing.
- Protobuf adapters may call application-owned binding or registry methods only after a later implementation gate defines the exact handoff.
- Authentication service behavior may revoke tokens, but it must not directly close sockets or mutate registry state from token repository code.
- Domain modules must not read, write, close, kick, disconnect, or target WebSocket connections directly.

## 3. Future Registry Shape

A later implementation gate must choose exact Go types, but the first posture should model registry records with this vocabulary:

```yaml
candidate_registry_record:
  connection_id: server_observed_connection_id
  connection_epoch: server_observed_epoch
  state:
    - open_unbound
    - bound
    - closing
    - closed
  bound_actor:
    actor_kind: player
    player_id: validated_player_id
  runtime_session_id: optional_validated_runtime_session_id
  access_token_record_id: optional_server_token_record_id
  opened_at: server_clock_time
  bound_at: optional_server_clock_time
  last_seen_at: optional_server_clock_time
  closed_at: optional_server_clock_time
  close_reason_class: optional_redacted_internal_reason
```

Rules:

- A registry record represents server-observed connection state, not client proof.
- `connection_id` and `connection_epoch` are server-owned identifiers. They may help correlate transport events, but they are not authentication proof.
- Bound player, runtime session, or token linkage must come only from validated application identity, not from client-supplied metadata alone.
- The first registry posture is single-process, in-memory, and non-durable. Durable connection state, cross-node lookup, Redis-like stores, service discovery, and server-to-server RPC remain deferred.
- Registry state must not store raw access-token text, raw credential material, lookup digests, verifier digests, verifier key ids, Authorization headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or inner payload bytes.

## 4. Future Capabilities

A later implementation gate may define narrow capabilities such as:

```yaml
candidate_registry_capabilities:
  - RegisterOpenConnection
  - BindConnectionIdentity
  - MarkConnectionClosed
  - FindConnectionByID
  - ListConnectionsByPlayerID
  - ListConnectionsByRuntimeSessionID
  - ListConnectionsByAccessTokenRecordID
  - MarkConnectionInvalidated
```

Rules:

- Registration must be driven by server-observed transport lifecycle events.
- Binding must require validated application identity from existing access-token and optional runtime-session validation paths.
- Listing by player, session, or token record id is only a targeting primitive. It must not itself decide close policy.
- Invalidating a registry record and closing a WebSocket are separate future actions.
- Any future public route or admin surface that targets connections must pass through application policy and redacted authorization checks.

## 5. Relationship To Logout And Revocation

This gate exists because logout/revocation needs a safe target model, but it does not implement active invalidation:

```yaml
logout_access_token_behavior_changed: false
token_revocation_behavior_changed: false
runtime_session_revocation_added: false
active_connection_invalidation_added: false
websocket_close_policy_added: false
logout_close_socket_default: not_selected_by_this_gate
logout_registry_lookup_default: not_selected_by_this_gate
```

Rules:

- `LogoutAccessToken` continues to revoke only the verified presented access-token record.
- Existing request-level access-token validation rejects revoked material on later protected requests when validation runs.
- A future policy must decide whether token revocation should look up active connections by access-token record id, runtime session id, player id, or connection id.
- A future policy must decide whether active connection invalidation failure fails logout, succeeds with a redacted warning, or is retried asynchronously.
- A future policy must decide whether a bound connection remains open but fails protected routes, receives a system message, closes immediately, or waits for reconnect/epoch handling.

## 6. Relationship To WebSocket And Protocol

This gate does not change WebSocket or Protobuf behavior:

```yaml
websocket_transport_credential_neutral: true
websocket_close_policy_added: false
transport_credential_carrier_added: false
websocket_handshake_authentication_added: false
protobuf_logout_route_added: false
protobuf_session_carrier_added: false
existing_protobuf_envelope_change_added: false
generated_output_added: false
```

Rules:

- WebSocket transport must not parse Authorization headers, bearer values, cookies, query-string tokens, session tokens, or `Sec-WebSocket-Protocol` authentication material for this gate.
- No close code, close reason, kick/disconnect message, logout route, session carrier, reconnect token, resume token, or connection epoch protocol field is authorized here.
- Generated Go Protobuf output must not change for this gate.
- The existing first-message `BindConnection` route remains the only selected connection identity route. This gate does not add registry-backed route-policy use of bound identity.

## 7. Future Test Expectations

A later implementation gate must require focused tests for:

- Registering an open connection from server-observed connection id and epoch.
- Rejecting duplicate active records for the same connection id and epoch unless a replacement policy is explicitly selected.
- Binding only after validated player identity is available.
- Rejecting metadata-only player id, session id, token record id, or client-supplied connection metadata as targeting proof.
- Marking close state from server-observed transport lifecycle events.
- Listing active connections by player id, runtime session id, and access-token record id without leaking raw proof material.
- Preserving WebSocket transport credential neutrality.
- Keeping invalidation policy and concrete socket close behavior outside the registry unless separately authorized.
- Redacting raw token text, lookup digests, verifier digests, verifier key ids, remote addresses, headers, cookies, query strings, and subprotocol values from errors and logs.

## 8. Nakama And Pitaya Reference Mapping

Nakama reference mapping:

- Adopt the lesson that realtime sockets and authenticated session material need coordinated lifecycle semantics.
- Adapt that into a vibit-owned registry that tracks server-observed active connections and validated identity linkage.
- Defer direct Nakama session APIs, JWT/session token shape, realtime socket compatibility, dashboard operations, and cluster session routing.

Pitaya reference mapping:

- Adopt the separation between acceptors, sessions, handlers, and connection management.
- Adapt Pitaya-style session/context separation into application-owned registry state plus narrow transport lifecycle handoff.
- Defer frontend/backend cluster routing, distributed kick/disconnect, groups/rooms integration, and server-to-server RPC invalidation.

## 9. Non-Goals

This gate does not authorize:

- Go runtime registry implementation.
- WebSocket socket close, kick, disconnect, or close reason behavior.
- Runtime session revocation.
- Logout-all, admin revocation, refresh, cleanup jobs, or token rotation.
- Reconnect, resume, duplicate replacement, or durable epoch behavior.
- Protocol logout routes, protocol session carriers, protocol close messages, or existing envelope changes.
- WebSocket handshake authentication or transport credential carriers.
- PostgreSQL, Redis-like, distributed, or durable active connection storage.
- Major new dependencies.
- Direct Nakama or Pitaya public API compatibility.

## 10. Next Gate

After this gate, the work queue must stop at a new confirmation point. The recommended next choices are:

```yaml
candidate_next_directions:
  - implement_active_connection_registry_single_process
  - define_websocket_close_policy_gate
  - define_protocol_logout_route_gate
  - define_reconnect_connection_epoch_gate
  - define_protocol_session_carrier_gate
  - strengthen_operations_observability_and_admin_tooling
  - expand_core_game_backend_modules_after_nakama_pitaya_review
```

The conservative recommendation is `implement_active_connection_registry_single_process`, but it must be selected explicitly.
