# WebSocket Close Policy Gate

Status: Draft v0.1
Last updated: 2026-05-18
Scope: Gate-only boundary for future WebSocket close, kick, and disconnect policy behavior
Depends on: `docs/active-connection-registry-gate.md`, `decisions/ADR-0075-active-connection-registry-single-process-implementation.md`, `docs/logout-revocation-active-connection-gate.md`, `docs/bound-identity-route-policy-gate.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0076`

The paired Simplified Chinese translation is `docs/websocket-close-policy-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

The active connection registry can now model server-observed connection state and validated identity linkage, but it deliberately cannot close sockets. That separation is required before vibit can decide what happens when logout, token revocation, session revocation, duplicate connection policy, reconnect policy, admin action, or operational drain targets an active WebSocket connection.

Nakama makes the product pressure clear: realtime sockets, authenticated sessions, logout, expiration, and single-socket style policies are lifecycle concerns that must behave predictably for players. Pitaya makes the layering pressure clear: acceptors, agents/sessions, handlers, and connection management are separate surfaces, and handlers should not close network connections directly as hidden business logic.

vibit adapts those lessons by requiring an application-owned WebSocket close policy boundary before any code can convert registry invalidation into concrete socket close behavior. This standard is gate-only.

```yaml
websocket_close_policy_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0164
decision: ADR-0076
check_rule: runtime.websocket_close_policy_gate
future_policy_owner: runtime/internal/app
future_registry_owner: runtime/internal/app/connection
future_transport_handoff_owner: runtime/internal/platform/transport/ws
future_protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
first_close_policy_posture: application_owned_policy_before_transport_handoff
transport_close_handoff_posture: deferred
registry_invalidation_to_close_default: not_selected_by_this_gate
logout_close_socket_default: not_selected_by_this_gate
kick_disconnect_behavior_added: false
websocket_close_implementation_added: false
close_code_mapping_added: false
close_reason_text_added: false
protocol_close_message_added: false
protobuf_logout_route_added: false
protobuf_session_carrier_added: false
reconnect_epoch_behavior_added: false
runtime_session_revocation_added: false
transport_credential_carrier_added: false
websocket_handshake_authentication_added: false
dependencies_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

This standard does not add close behavior.

## 2. Ownership

Future WebSocket close policy must stay application-owned:

```yaml
future_close_policy_owner: runtime/internal/app
future_registry_owner: runtime/internal/app/connection
future_transport_handoff_owner: runtime/internal/platform/transport/ws
future_protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
authentication_service_owner: runtime/internal/app/authentication
session_repository_owner: runtime/internal/app/session
domain_handler_owner: runtime/internal/modules/*
```

Rules:

- The application layer owns the policy decision that a connection should be closed, disconnected, kicked, drained, invalidated, or left open.
- WebSocket transport may own only a future narrow concrete close handoff after application policy has produced a redacted close intent.
- WebSocket transport must not parse credentials, decide identity, decide logout policy, decide session policy, decide route policy, or choose player-facing close reason text.
- Authentication service behavior may revoke tokens, but it must not directly close sockets or own transport close handoff.
- Active connection registry may mark records invalidated or closed, but registry listing and invalidation are not themselves close policy.
- Protobuf adapters may call future application close policy only after a later implementation gate defines exact route and payload behavior.
- Domain modules must not close, kick, disconnect, or target WebSocket connections directly.

## 3. Future Close Intent Vocabulary

A later implementation gate must choose exact Go types, but the first policy vocabulary should separate intent, target, transport action, and outcome:

```yaml
candidate_close_intent:
  target:
    - connection_id_and_epoch
    - player_id
    - runtime_session_id
    - access_token_record_id
  reason_class:
    - token_revoked
    - logout_presented_token
    - session_revoked
    - duplicate_connection_policy
    - server_shutdown_or_drain
    - policy_violation
    - administrative_action
    - protocol_error
    - idle_timeout
    - unknown_internal
  transport_action:
    - close_socket
    - mark_invalidated_only
    - send_system_message_then_close
    - defer_to_reconnect_policy
  retryability:
    - retryable
    - not_retryable
    - unknown
  public_visibility:
    - silent
    - generic_disconnect
    - generic_reauth_required
```

Rules:

- Close intents must be derived from server-owned registry records and validated application identity, not client metadata alone.
- Close reason classes are internal semantic classes. They are not raw WebSocket close reason text.
- Future player-visible close text, system messages, or protocol errors must be separately ratified.
- Registry invalidation and concrete socket close are separate actions until an implementation gate selects a policy.
- Close policy must define whether failure to close a socket fails the caller, succeeds with a redacted warning, records a retryable action, or is ignored.
- Close policy must define whether it targets by connection id/epoch, player id, runtime session id, access-token record id, or a composed policy.

## 4. Close Code And Reason Boundaries

This gate does not select WebSocket close codes or reason strings:

```yaml
close_code_mapping_added: false
close_reason_text_added: false
custom_close_codes_selected: false
websocket_status_code_dependency_added: false
player_visible_close_reason_added: false
```

Rules:

- No custom WebSocket close code is authorized by this gate.
- No close reason text, kick reason, disconnect reason, or player-facing system message is authorized by this gate.
- A future implementation must map internal reason classes to redacted close behavior without leaking raw token material, verifier material, session ids, internal repository ids, remote addresses, headers, cookies, query strings, subprotocol values, or database errors.
- Transport-level abnormal closes, peer disconnects, network failures, and application-directed closes must remain distinguishable in internal records, but public behavior must be redacted.

## 5. Relationship To Active Connection Registry

The registry is the target model; it is not the close policy:

```yaml
registry_lookup_allowed_in_future_policy: true
registry_invalidation_to_close_default: not_selected_by_this_gate
registry_mark_closed_from_transport_lifecycle: future_handoff_only
registry_mark_invalidated_from_policy: future_policy_only
```

Rules:

- Future close policy may query the registry by connection id/epoch, player id, runtime session id, or access-token record id.
- Listing registry records must not decide whether a socket should close.
- Marking a registry record invalidated must not itself close a concrete socket unless a later implementation gate explicitly wires that action.
- Transport-observed peer close and application-directed close are different lifecycle facts and must not be collapsed.
- Duplicate connection replacement requires a separate reconnect/epoch or duplicate-policy gate.

## 6. Relationship To Logout And Session Revocation

This gate does not change logout or session behavior:

```yaml
logout_access_token_behavior_changed: false
token_revocation_behavior_changed: false
runtime_session_revocation_added: false
logout_close_socket_default: not_selected_by_this_gate
session_revocation_close_socket_default: not_selected_by_this_gate
admin_kick_default: not_selected_by_this_gate
```

Rules:

- `LogoutAccessToken` continues to revoke only the verified presented access-token record.
- Existing request-level access-token validation rejects revoked material on later protected requests when validation runs.
- A future policy must decide whether token revocation looks up active connections by access-token record id, runtime session id, player id, or connection id.
- A future policy must decide whether logout succeeds if socket close fails.
- A future policy must decide whether runtime session revocation closes active sockets, invalidates registry records only, or defers to reconnect/session validation policy.

## 7. Relationship To WebSocket And Protocol

This gate does not change WebSocket or Protobuf behavior:

```yaml
websocket_transport_credential_neutral: true
websocket_close_implementation_added: false
transport_close_handoff_added: false
transport_credential_carrier_added: false
websocket_handshake_authentication_added: false
protobuf_logout_route_added: false
protobuf_session_carrier_added: false
protocol_close_message_added: false
existing_protobuf_envelope_change_added: false
generated_output_added: false
```

Rules:

- WebSocket transport must not parse Authorization headers, bearer values, cookies, query-string tokens, session tokens, or `Sec-WebSocket-Protocol` authentication material for this gate.
- No close code, close reason, kick/disconnect message, logout route, session carrier, reconnect token, resume token, or connection epoch protocol field is authorized here.
- Generated Go Protobuf output must not change for this gate.
- A future close handoff must be narrow enough that transport closes concrete sockets without owning application policy.

## 8. Future Test Expectations

A later implementation gate must require focused tests for:

- Producing close intents from server-owned registry records only.
- Rejecting metadata-only player id, session id, token record id, or client-supplied connection metadata as close targets.
- Keeping close decision policy in the application layer.
- Keeping concrete socket close mechanics in a narrow transport handoff.
- Keeping registry invalidation separate from concrete socket close unless explicitly selected.
- Mapping internal reason classes to redacted public behavior.
- Preserving logout/token revocation behavior if socket close fails unless a future policy selects a different outcome.
- Preserving WebSocket transport credential neutrality.
- Avoiding generated Protobuf changes and protocol close messages unless separately authorized.

## 9. Nakama And Pitaya Reference Mapping

Nakama reference mapping:

- Adopt the lesson that realtime socket lifecycle, session expiration, logout, and single-socket style policies need explicit lifecycle semantics.
- Adapt that into vibit-owned close policy vocabulary that can later coordinate registry records, token revocation, runtime sessions, and reconnect behavior.
- Defer direct Nakama session APIs, JWT/session token shape, realtime socket compatibility, dashboard operations, cluster session routing, and exact close behavior.

Pitaya reference mapping:

- Adopt the separation between acceptors, agents/sessions, route handlers, and connection management.
- Adapt Pitaya-style connection management into application-owned close policy plus a future narrow transport handoff.
- Defer frontend/backend cluster routing, distributed kick/disconnect, groups/rooms integration, and server-to-server RPC invalidation.

## 10. Non-Goals

This gate does not authorize:

- Go runtime WebSocket close implementation.
- Transport close handoff code.
- Close codes, close reason strings, kick messages, disconnect messages, or protocol close messages.
- Logout-triggered socket close.
- Runtime session revocation-triggered socket close.
- Admin kick/disconnect behavior.
- Duplicate connection replacement.
- Reconnect, resume, durable epoch behavior, or reconnect token behavior.
- Protocol logout routes, protocol session carriers, protocol close messages, or existing envelope changes.
- WebSocket handshake authentication or transport credential carriers.
- PostgreSQL, Redis-like, distributed, or durable active connection storage.
- Major new dependencies.
- Direct Nakama or Pitaya public API compatibility.

## 11. Next Gate

After this gate, the work queue must stop at a new confirmation point. The recommended next choices are:

```yaml
candidate_next_directions:
  - implement_websocket_close_policy_single_process
  - define_protocol_logout_route_gate
  - define_reconnect_connection_epoch_gate
  - define_protocol_session_carrier_gate
  - strengthen_operations_observability_and_admin_tooling
  - expand_core_game_backend_modules_after_nakama_pitaya_review
```

The conservative recommendation is `implement_websocket_close_policy_single_process`, but it must be selected explicitly.
