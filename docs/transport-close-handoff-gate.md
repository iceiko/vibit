# Transport Close Handoff Gate

Status: Draft v0.1
Last updated: 2026-05-18
Scope: Gate-only boundary for future application-to-WebSocket concrete close handoff
Depends on: `docs/websocket-close-policy-gate.md`, `decisions/ADR-0077-websocket-close-policy-single-process-implementation.md`, `docs/protocol-logout-route-gate.md`, `docs/nakama-pitaya-product-parity-roadmap.md`, `docs/runtime-protocol-adapter.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0080`

The paired Simplified Chinese translation is `docs/transport-close-handoff-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

vibit now has three separate lifecycle pieces:

- Active connection registry records server-observed connection state.
- Application close policy can resolve active bound records and produce redacted close intents.
- Protocol logout route can revoke the presented access-token record without closing sockets.

The missing piece is a narrow future handoff from application-owned close policy to concrete WebSocket socket close mechanics. Without a gate, future code could make WebSocket transport own authentication policy, make logout close sockets implicitly, close by client-supplied metadata, or hide session/reconnect decisions inside protocol handlers.

Nakama is the product reference for explicit lifecycle behavior across sessions, logout, realtime sockets, and server-directed disconnects. Pitaya is the architecture reference for separating acceptors, sessions, route handlers, groups, RPC, and kick/disconnect style connection management. vibit adapts those lessons into a gate that keeps policy in the application layer and allows transport to perform only a narrow concrete close action after a server-owned handoff.

This standard is gate-only.

```yaml
transport_close_handoff_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0172
decision: ADR-0080
check_rule: runtime.transport_close_handoff_gate
parity_phase: phase_2r_runtime_lifecycle_closure
application_policy_owner: runtime/internal/app/connection
active_connection_registry_owner: runtime/internal/app/connection
future_transport_handoff_owner: runtime/internal/platform/transport/ws
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
authentication_service_owner: runtime/internal/app/authentication
first_handoff_target: connection_id_and_epoch
server_observed_target_required: true
client_metadata_authority_allowed: false
transport_policy_ownership_allowed: false
transport_credential_parsing_allowed: false
transport_session_revocation_allowed: false
first_transport_action_candidate: close_socket
close_code_mapping_added: false
close_reason_text_added: false
protocol_close_message_added: false
logout_triggered_socket_close_added: false
runtime_session_revocation_added: false
reconnect_epoch_behavior_added: false
protocol_session_carrier_added: false
dependencies_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Ownership

The future handoff must preserve these ownership boundaries:

```yaml
application_policy_owner: runtime/internal/app/connection
active_connection_registry_owner: runtime/internal/app/connection
future_transport_handoff_owner: runtime/internal/platform/transport/ws
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
authentication_service_owner: runtime/internal/app/authentication
domain_module_owner: runtime/internal/modules/*
```

Rules:

- Application close policy owns the decision that a connection should be closed.
- Active connection registry owns server-observed connection records and invalidation/closed markers.
- WebSocket transport may own only concrete mechanics for closing an already accepted socket after application policy hands it a server-owned target.
- WebSocket transport must not parse credentials, validate tokens, select players, select runtime sessions, evaluate logout state, decide reconnect behavior, or choose player-facing text.
- Protocol adapters must not directly close sockets or create close targets from client payload metadata.
- Authentication service may revoke token records, but it must not call WebSocket transport or own close handoff.
- Domain modules must not import WebSocket transport or close concrete sockets directly.

## 3. Handoff Target

The first future handoff target must be:

```yaml
first_handoff_target: connection_id_and_epoch
connection_id_source: server_observed_websocket_accept_metadata
connection_epoch_source: server_observed_websocket_accept_metadata
requires_active_registry_record: true
client_supplied_connection_id_authority: false
client_supplied_epoch_authority: false
player_id_transport_authority: false
runtime_session_id_transport_authority: false
access_token_record_id_transport_authority: false
```

Rules:

- Future transport handoff must target concrete sockets by server-observed `connection_id` and `connection_epoch`.
- Application policy may resolve player, runtime session, or access-token record targets into concrete connection/epoch targets before handoff.
- Transport must not close by player id, runtime session id, access-token record id, route identity, request identity, envelope session metadata, headers, cookies, query strings, subprotocol values, or remote address.
- Connection epoch must prevent a stale close intent from closing a later socket that reused a connection id.
- If the concrete socket is not found or the epoch does not match, the failure must be redacted and policy-neutral.

## 4. Future Handoff Shape

A later implementation gate may choose exact Go types. The first vocabulary should be narrow:

```yaml
candidate_transport_close_request:
  connection_id: server_observed_connection_id
  connection_epoch: server_observed_connection_epoch
  reason_class: internal_redacted_close_reason_class
  public_visibility: silent_or_generic_disconnect_or_generic_reauth_required
  retryability: retryable_or_not_retryable_or_unknown
  requested_at: server_time

candidate_transport_close_result:
  connection_id: server_observed_connection_id
  connection_epoch: server_observed_connection_epoch
  outcome:
    - close_requested
    - socket_not_found
    - epoch_mismatch
    - already_closed
    - close_failed
  closed_at: server_time_optional
```

Rules:

- The request must not carry raw access-token material, raw credential material, lookup digests, verifier digests, verifier key ids, headers, cookies, query strings, subprotocol values, remote addresses, database errors, or full repository errors.
- `reason_class` is an internal semantic class, not a WebSocket close reason string.
- The first gate does not select WebSocket close codes, close reason text, protocol close messages, or player-facing system messages.
- The future implementation must define whether transport close errors are returned to the application policy, recorded as redacted outcomes, or both.
- The handoff must be testable without live network dependencies where possible.

## 5. Relationship To Existing Close Policy

Current close policy emits `CloseIntent` values with:

```yaml
current_transport_action: mark_invalidated_only
concrete_socket_close_added: false
```

Future handoff work must not silently reinterpret existing `mark_invalidated_only` behavior as concrete close. A later implementation must explicitly choose how application policy produces close requests:

```yaml
future_options:
  - keep_mark_invalidated_only_and_add_separate_close_socket_action
  - allow_close_policy_to_emit_close_socket_after_gate
  - compose_invalidated_marker_then_transport_close_request
```

Rules:

- Registry invalidation and concrete socket close remain separate lifecycle facts.
- A transport close request must not be generated unless application policy selected a close action.
- Transport-observed peer close and application-directed close must remain distinguishable for registry updates.
- Failure to close a concrete socket must not make token revocation, logout, or session mutation appear to have succeeded or failed unless a later policy explicitly defines that coupling.

## 6. Relationship To Logout, Sessions, And Reconnect

This gate does not change logout, session, or reconnect behavior:

```yaml
logout_triggered_socket_close_added: false
runtime_session_revocation_added: false
session_revocation_close_added: false
duplicate_connection_replacement_added: false
reconnect_epoch_behavior_added: false
protocol_session_carrier_added: false
```

Rules:

- `LogoutAccessToken` continues to revoke only the verified presented access-token record.
- This gate does not decide whether logout should later close the current socket or any sockets linked to the token record.
- Runtime session revocation remains a separate future behavior.
- Duplicate connection replacement and reconnect/resume behavior remain separate future gates.
- Protocol session carriers remain separate from transport close handoff.

## 7. WebSocket And Protocol Boundaries

This gate keeps WebSocket transport credential-neutral and protocol-neutral:

```yaml
websocket_transport_credential_neutral: true
transport_credential_parsing_added: false
websocket_handshake_authentication_added: false
protocol_close_message_added: false
protobuf_source_added: false
generated_output_added: false
existing_protobuf_envelope_change_added: false
```

Rules:

- No `.proto` source or generated output is authorized by this gate.
- No existing Protobuf envelope fields may change for this gate.
- No close route, kick route, disconnect route, system close message, or admin disconnect protocol surface is authorized by this gate.
- WebSocket transport must not read Authorization headers, bearer values, cookies, query-string tokens, session tokens, or subprotocol authentication material.
- Future close reason text must remain redacted and separately ratified.

## 8. Required Future Tests

A later implementation gate must require focused tests for:

```yaml
required_tests:
  - application_policy_remains_owner_of_close_decision
  - transport_handoff_targets_connection_id_and_epoch_only
  - stale_epoch_does_not_close_new_socket
  - missing_socket_returns_redacted_not_found_outcome
  - transport_does_not_parse_credentials_or_tokens
  - protocol_adapter_does_not_close_sockets_directly
  - authentication_service_does_not_call_transport
  - domain_modules_do_not_import_websocket_transport
  - close_reason_class_does_not_leak_secrets_or_internal_ids
  - existing_protobuf_envelope_remains_unchanged
  - logout_route_behavior_remains_token_record_scoped
```

Live end-to-end socket tests may be useful later but are not required by this gate-only standard.

## 9. Nakama And Pitaya Reference Mapping

Nakama reference mapping:

- Adopt the product lesson that logout, session invalidation, realtime disconnect, and server-directed disconnect are explicit lifecycle concerns.
- Adapt that into vibit's application-owned close policy and narrow WebSocket close handoff.
- Defer Nakama session API compatibility, realtime socket compatibility, dashboard disconnect behavior, session logout-all behavior, and cluster session routing.

Pitaya reference mapping:

- Adopt the layering lesson that acceptors, sessions, handlers, groups, RPC, and kick/disconnect behavior should remain separate.
- Adapt that into a transport handoff that closes concrete sockets only after application policy selects a server-observed connection target.
- Defer Pitaya route naming compatibility, frontend/backend topology, distributed kick/disconnect, group broadcast integration, and RPC/session propagation.

## 10. Non-Goals

This gate does not authorize:

- Concrete WebSocket socket close implementation.
- Close code mapping.
- Close reason text.
- Player-facing kick/disconnect messages.
- Protocol close messages.
- Protobuf source or generated output changes.
- Logout-triggered socket close.
- Runtime session revocation.
- Duplicate connection replacement.
- Reconnect/resume/epoch behavior.
- Protocol session carriers.
- Presence, chat, groups, parties, matchmaking, match runtime, SDK, cluster, or distributed runtime behavior.
- New dependencies.
- Direct Nakama/Pitaya API compatibility.

## 11. Agent Rules

Agents must:

- Read this standard before adding concrete WebSocket close handoff code.
- Keep close policy decisions in `runtime/internal/app/connection`.
- Keep concrete socket mechanics in `runtime/internal/platform/transport/ws`.
- Use only server-owned connection id and epoch as the first transport target.
- Preserve redaction boundaries for tokens, credentials, digests, verifier keys, headers, cookies, query strings, subprotocol values, remote addresses, database errors, and internal repository errors.
- Keep generated Protobuf output unchanged unless a later protocol gate authorizes it.

Agents must not:

- Close sockets from authentication service behavior.
- Close sockets directly from protocol bridge code.
- Close sockets from domain modules.
- Treat client-supplied player id, session id, token record id, connection id, epoch, or envelope metadata as close authority.
- Add close codes, reason text, reconnect behavior, session carriers, operations/admin disconnect, dependencies, or direct Nakama/Pitaya compatibility in this gate.

## 12. Verification

Required repository checks for this gate:

```bash
node -c tools/vibit
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change define-transport-close-handoff-gate --json
node tools/vibit inspect next --json
```

Go runtime tests are not required for this gate-only standard because no Go runtime behavior is added.
