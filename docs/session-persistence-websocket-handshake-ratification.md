# Session Persistence And WebSocket Handshake Ratification

Status: Draft v0.1
Last updated: 2026-05-17
Scope: First ratified posture for request-level validation, future session persistence, future connection binding, WebSocket handshake authentication deferral, and session-related implementation gates after public login route exposure
Depends on: `docs/session-persistence-websocket-handshake-decision-gates.md`, `docs/access-token-protocol-carrier-route-protection-gate.md`, `docs/authentication-command-protocol-login-route-gate.md`, `docs/runtime-protocol-adapter.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0056`

The paired Simplified Chinese translation is `docs/session-persistence-websocket-handshake-ratification.zh-CN.md`. The English file is authoritative.

## 1. Purpose

vibit now has the minimum authenticate-then-gameplay loop:

- Public device credential login route: `runtime.authentication.AuthenticateWithDeviceCredential`.
- Opaque access-token validation in the application authentication service.
- Request-level route protection through `vibit.authentication.v1.AuthenticatedRequest`.
- PostgreSQL startup composition for authentication service and route protection.

The next game-server concern is session and connection lifecycle. Mature servers such as Nakama expose an authentication/session model and then connect realtime sockets with authenticated client state. Pitaya separates acceptors, sessions, route handlers, groups, and later cluster roles. vibit should learn from those systems, but still preserve its own agent-native boundaries.

This standard ratifies the first posture after public login:

```yaml
session_persistence_websocket_handshake_ratification: defined
current_validation_model: request_level_access_token_validation
current_proof_carrier: protobuf_authenticated_request_payload_wrapper
current_websocket_handshake_authentication: not_selected
current_websocket_transport_credential_neutral: true
current_protobuf_envelope_change: unchanged
future_connection_binding_preferred_gate: first_message_protocol_binding
future_session_store_preferred_first_durable_target: postgres
implementation_authorized_by_this_standard: false
completed_work_item: W-0122
decision: ADR-0056
check_rule: runtime.session_persistence_websocket_handshake_ratification
```

This is a ratification gate, not an implementation gate.

## 2. Selected Current Path

The current production-sensitive path remains request-level validation:

```yaml
login:
  route: runtime.authentication.AuthenticateWithDeviceCredential
  result: opaque_access_token
protected_request:
  carrier: vibit.authentication.v1.AuthenticatedRequest
  validation_owner: runtime/internal/app/authentication
  route_policy_owner: runtime/internal/app
  domain_handoff: normalized RequestIdentity
  session_validated: false
websocket_transport:
  credential_neutral: true
protobuf_envelope:
  changed: false
```

Rules:

- `RequestIdentity.SessionValidated` remains false until a later session persistence or binding implementation explicitly validates a session.
- Client-supplied envelope `Session.player_id`, `Session.session_id`, `Session.connection_id`, and `Session.connection_epoch` remain metadata only.
- Protected domain routes must rely on normalized request identity from route protection, not metadata-only envelope fields.
- The current proof carrier remains the Protobuf authenticated request payload wrapper, not WebSocket handshake metadata.

## 3. Ratified Deferral Of WebSocket Handshake Authentication

WebSocket handshake authentication is not selected as the next implementation path.

The WebSocket transport must continue to ignore credential carriers:

```yaml
forbidden_current_transport_carriers:
  - HTTP Authorization header
  - Bearer header value
  - Cookie access_token or session token
  - query-string access_token or session token
  - Sec-WebSocket-Protocol authentication token
```

Rules:

- The WebSocket transport must not own player account lookup.
- The WebSocket transport must not own credential lookup.
- The WebSocket transport must not own token verifier lookup.
- The WebSocket transport must not create authenticated `RequestIdentity`.
- Any future transport-level extraction requires a separate transport-auth boundary and must hand off to application-owned or authentication-owned validation contracts.

Rationale:

- Browser and non-browser clients have different practical constraints for headers, cookies, query strings, and subprotocols.
- Rejecting a connection before a Protobuf envelope exists changes error behavior and compatibility.
- Transport-owned authentication would make future non-WebSocket transports harder to keep consistent.

## 4. Future Connection Binding Preference

The preferred next connection-level gate is first-message protocol binding:

```yaml
future_connection_binding_preferred_gate: first_message_protocol_binding
candidate_owner_layers:
  protocol_adapter: decode system/authentication binding message
  application_runtime: validate proof and bind normalized identity to connection context
  websocket_transport: keep connection plumbing credential-neutral
```

This preference does not implement first-message binding. A later gate must define:

- System or authentication message name.
- Protobuf payload shape.
- Timeout and failure behavior.
- Whether anonymous pre-bind messages are allowed.
- Connection context storage owner.
- Revalidation behavior.
- Logout, revocation, expiration, and reconnect behavior.
- How route protection uses bound identity versus per-request proof.
- How duplicate active connections are handled.
- How future room, party, match, presence, and group membership attach to binding.

Rules:

- First-message binding must not be invented as an ad hoc domain command.
- Binding state must not live in inventory, player, or other domain modules.
- Binding identity must still enter domain dispatch through normalized application request identity.

## 5. Future Session Persistence Posture

The preferred first durable target for future session persistence is PostgreSQL, because PostgreSQL is already the accepted first authoritative durable store for module state.

This preference does not add session persistence.

Future session persistence requires a separate schema gate that defines:

- Whether sessions are persisted at all.
- Session record ownership.
- Session id generation.
- Actor kind and actor id fields.
- Player id linkage.
- Issue, expiration, revocation, and last-seen semantics.
- Connection binding fields, if any.
- Token record linkage, if any.
- Logout and revocation transaction boundaries.
- Cleanup strategy.
- Migration source.
- Repository interface and PostgreSQL adapter boundaries.
- Live verification requirements.

Rules:

- No `runtime_sessions` table is authorized by this standard.
- No repository method is authorized by this standard.
- No PostgreSQL adapter behavior is authorized by this standard.
- No Redis-like or external session-store dependency is selected.
- Memory session storage is not selected as production behavior.

## 6. Nakama And Pitaya Reference Mapping

Nakama reference mapping:

```yaml
adopted_concepts:
  - authenticate_before_protected_gameplay
  - session_lifetime_dimensions
  - realtime_socket_associated_with_authenticated_state
adapted_concepts:
  - token_or_session_material_is_vibit_opaque_access_token_first
  - realtime_socket_binding_is_deferred_to_protocol_application_gate
deferred_concepts:
  - refresh_token_behavior
  - session_persistence_schema
  - socket_reconnect_restore_behavior
rejected_for_now:
  - direct_nakama_api_compatibility
```

Pitaya reference mapping:

```yaml
adopted_concepts:
  - transport_acceptor_separated_from_session_binding
  - route_handler_separated_from_connection_plumbing
  - groups_rooms_broadcast_require_later_binding_context
adapted_concepts:
  - session_binding_becomes_application_protocol_first_message_gate
deferred_concepts:
  - frontend_backend_split
  - server_to_server_rpc
  - cluster_service_discovery
rejected_for_now:
  - direct_pitaya_api_compatibility
```

Nakama and Pitaya guide capability planning. They do not override vibit's constitution, manifests, ADRs, generated boundaries, or verification commands.

## 7. Required Future Gates

Before implementing session or connection behavior, future work must create one of these gates:

```yaml
future_gates:
  first_message_connection_binding_gate:
    required_before:
      - system/authentication binding protobuf messages
      - connection-bound identity cache
      - route policy use of bound identity
      - reconnect or epoch behavior
  postgres_session_persistence_schema_gate:
    required_before:
      - session tables
      - session repository interfaces
      - PostgreSQL session adapters
      - session cleanup jobs
  websocket_handshake_authentication_gate:
    required_before:
      - Authorization/Bearer/cookie/query/subprotocol credential parsing
      - pre-envelope handshake rejection behavior
      - transport-level proof extraction
  logout_revocation_active_connection_gate:
    required_before:
      - logout execution
      - active session revocation
      - closing or invalidating bound WebSocket connections
  reconnect_epoch_gate:
    required_before:
      - connection_epoch semantics
      - duplicate connection replacement
      - session resume behavior
```

## 8. Deferrals

This standard does not authorize:

- Session persistence implementation.
- Session schema, migrations, repositories, or PostgreSQL adapters.
- WebSocket handshake authentication.
- HTTP `Authorization`, Bearer, cookie, query-string, or WebSocket subprotocol proof carriers.
- Existing Protobuf envelope changes.
- New Protobuf system messages.
- Generated Protobuf output.
- Connection-bound identity cache.
- Reconnect, resume, or connection epoch behavior.
- Logout, refresh, cleanup, token rotation, or token validation audit mutation.
- New dependencies.
- Memory durable authentication behavior.
- Direct Nakama or Pitaya public API compatibility.

## 9. Verification

The repository check rule for this ratification is:

```text
runtime.session_persistence_websocket_handshake_ratification
```

The check must verify:

- The standard, translation, ADR, change specs, and conversation log exist.
- Manifests and AGENTS guides record the ratified posture.
- The WebSocket transport remains credential-neutral.
- The existing Protobuf envelope has no token, credential, login, or handshake proof fields.
- No session migration source exists as part of this gate.
- No new generated Protobuf output is required by this gate.
