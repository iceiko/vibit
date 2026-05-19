# First Message Connection Binding Gate

Status: Draft v0.1
Last updated: 2026-05-17
Scope: Future protocol/application connection binding posture after session persistence and WebSocket handshake ratification
Depends on: `docs/session-persistence-websocket-handshake-ratification.md`, `docs/access-token-protocol-carrier-route-protection-gate.md`, `docs/authentication-command-protocol-login-route-gate.md`, `docs/runtime-protocol-adapter.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0057`

The paired Simplified Chinese translation is `docs/first-message-connection-binding-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

vibit now has:

- A public device credential login route: `runtime.authentication.AuthenticateWithDeviceCredential`.
- Request-level opaque access-token validation through `vibit.authentication.v1.AuthenticatedRequest`.
- PostgreSQL startup composition for the authentication service and route protection.
- A ratified session and WebSocket handshake posture in `ADR-0056`.

The next connection-lifecycle question is how an already authenticated player becomes associated with an open WebSocket connection without moving credential parsing into the WebSocket transport or changing the existing Protobuf envelope.

This gate defines the future first-message connection binding posture:

```yaml
first_message_connection_binding_gate: defined
implementation_authorized_by_this_standard: false
completed_gate_work_item: W-0124
decision: ADR-0057
check_rule: runtime.first_message_connection_binding_gate
selected_binding_message_kind: system
selected_binding_route: runtime.authentication.BindConnection
selected_binding_payload_candidate: vibit.authentication.v1.BindConnectionRequest
selected_binding_response_candidate: vibit.authentication.v1.BindConnectionResponse
selected_proof_carrier: protobuf_system_payload_access_token
websocket_transport_credential_neutral: true
websocket_handshake_authentication_added: false
existing_protobuf_envelope_change_added: false
protocol_source_added: false
generated_output_added: false
connection_binding_registry_added: false
route_policy_bound_identity_added: false
session_persistence_added: false
```

This is a gate-only standard. It does not implement connection binding.

## 2. Selected Future Shape

The selected future route is:

```yaml
route:
  kind: system
  module: runtime.authentication
  name: BindConnection
semantic_route_key: runtime.authentication.BindConnection
payload_candidate: vibit.authentication.v1.BindConnectionRequest
response_candidate: vibit.authentication.v1.BindConnectionResponse
```

This message is a protocol/application binding message. It is not a domain command, and it must not be routed through inventory, player, or other gameplay modules.

Planned request shape:

```yaml
BindConnectionRequest:
  access_token: opaque Base64URL unpadded 32-byte access-token proof
  client_instance_id: optional client installation or runtime instance hint
```

Planned response shape:

```yaml
BindConnectionResponse:
  binding_status: bound | rejected
  actor_kind: player
  player_id: validated player id when bound
  connection_id: server-assigned connection id
  connection_epoch: server-assigned or server-observed connection epoch
  session_validated: false until durable session persistence is implemented
  bound_at: server timestamp when binding succeeds
```

Rules:

- The access token remains an opaque proof and must be validated by the existing application authentication service boundary.
- The WebSocket transport must not parse the access token from HTTP headers, cookies, query strings, bearer strings, or subprotocols.
- The existing `proto/vibit/protocol/v1/envelope.proto` shape must not change for first-message binding.
- The future bind route must be handled before normal domain dispatch.
- Binding success must produce normalized application identity, not trust metadata-only envelope fields.

## 3. Future Layer Ownership

The future implementation must preserve this layer split:

```yaml
websocket_transport:
  owns:
    - connection accept
    - server-assigned connection id
    - binary frame read and write
  must_not_own:
    - access-token parsing
    - authentication service calls
    - player account lookup
    - request identity construction

protobuf_protocol_adapter:
  owns:
    - existing envelope decode
    - recognition of runtime.authentication.BindConnection system route
    - bind request payload decode
    - application binder handoff
    - bind response or public error envelope encode
  must_not_own:
    - token verifier lookup
    - credential lookup
    - durable session storage
    - logout or reconnect policy

application_runtime:
  owns:
    - connection binding policy
    - access-token validation handoff
    - normalized identity creation
    - connection-bound identity registry if implemented
    - route policy use of bound identity if implemented
  must_not_own:
    - WebSocket handshake credential extraction
    - direct Nakama or Pitaya API compatibility
```

Candidate future files are:

```yaml
protocol_source: proto/vibit/authentication/v1/authentication.proto
generated_go_output: runtime/internal/generated/proto/vibit/authentication/v1/authentication.pb.go
protocol_adapter_source: runtime/internal/platform/protocol/protobuf/connection_binding*.go
application_binding_source: runtime/internal/app/connection_binding*.go
startup_composition_owner: runtime/cmd/vibit-server
```

Generated Go Protobuf output must be produced by Buf and must not be edited by hand.

## 4. Binding Flow

The planned future flow is:

```yaml
connection_binding_flow:
  - websocket_transport_accepts_connection_without_reading_credentials
  - websocket_transport_assigns_connection_id
  - client_may_call_public_login_route_if_it_needs_an_access_token
  - client_sends_runtime.authentication.BindConnection_system_message
  - protobuf_adapter_decodes_bind_payload
  - application_binding_policy_validates_access_token
  - authentication_service_returns_validated_player_identity
  - application_runtime_binds_identity_to_connection_id
  - bind_response_confirms_bound_identity_without_returning_secret_material
```

Anonymous pre-bind messages:

```yaml
allowed_before_binding:
  - runtime.authentication.AuthenticateWithDeviceCredential
  - runtime.authentication.BindConnection
  - heartbeat_or_ack_frames_after_a_later_protocol_rule_allows_them
forbidden_before_binding:
  - bound_identity_only_gameplay_routes
  - future_presence_room_party_match_membership_routes
  - future_group_broadcast_membership_routes
```

Rules:

- Public login remains usable before connection binding because a new client may not have an access token yet.
- First-message binding means the first identity-binding message for a connection, not necessarily the first frame ever sent on that connection.
- Metadata-only envelope session fields never satisfy binding.
- A successful bind must use the server-observed connection id, not a client-supplied connection id.

## 5. Route Policy Relationship

The current route-protection behavior remains request-level access-token validation through `vibit.authentication.v1.AuthenticatedRequest`.

Future binding implementation may add connection-bound identity as a proof source only if the implementation explicitly updates application route policy. The planned order is:

```yaml
future_protected_route_identity_sources:
  first: request_level_access_token_wrapper
  second: bound_connection_identity_after_explicit_route_policy_update
```

Rules:

- Existing protected routes must not silently become public.
- A bound connection identity must still be represented as normalized `RequestIdentity`.
- `RequestIdentity.SessionValidated` remains false until a separate durable session persistence implementation validates a session.
- Domain handlers must not query connection binding registries directly.
- Inventory, player, authentication repository interfaces, and PostgreSQL adapters must not be changed by this gate.

## 6. Failure Behavior

Future implementation must collapse public failures:

```yaml
future_public_errors:
  missing_bind_proof: CONNECTION_BINDING_TOKEN_MISSING
  malformed_bind_payload_or_token: CONNECTION_BINDING_TOKEN_MALFORMED
  invalid_or_expired_or_revoked_token: CONNECTION_BINDING_TOKEN_INVALID
  validation_dependency_unavailable: CONNECTION_BINDING_UNAVAILABLE
  protected_bound_route_without_binding: CONNECTION_BINDING_REQUIRED
```

Rules:

- Public errors must not reveal token lookup hit or miss, token lifecycle status, verifier key id, player account state, verifier mismatch, or internal binding registry state.
- Invalid binding proof must not bind any identity to the connection.
- Rebinding an already bound connection, duplicate player connections, active-connection invalidation, and kick/replace behavior require later gates.
- This gate selects no bind timeout and no automatic close-on-bind-failure behavior. Later implementation may return a protocol error without closing the WebSocket.
- Transport-level close policy for repeated failures is an operations or abuse-control concern and remains deferred.

## 7. Session Persistence And Reconnect Deferrals

First-message connection binding is not durable session persistence.

Rules:

- No `runtime_sessions` table is authorized.
- No session repository is authorized.
- No session migration is authorized.
- No Redis-like or external session-store dependency is selected.
- No reconnect resume behavior is authorized.
- No connection epoch replacement behavior is authorized.
- No logout-triggered active connection invalidation is authorized.
- No token rotation or refresh behavior is authorized.

Future durable sessions need a PostgreSQL-first schema/repository/migration gate. Future reconnect and duplicate connection behavior need a reconnect/epoch gate.

## 8. Nakama And Pitaya Reference Mapping

Nakama reference mapping:

```yaml
adopted_concepts:
  - authenticate_before_realtime_socket_features
  - socket_connection_associated_with_authenticated_state
  - session_and_socket_lifecycle_are_related_but_not_identical
adapted_concepts:
  - vibit_uses_opaque_access_token_proof_not_nakama_jwt_session_compatibility
  - vibit_binds_connection_through_protocol_application_message_not_websocket_handshake
deferred_concepts:
  - refresh_token_model
  - single_socket_single_session_policy
  - server_side_session_disconnect
  - reconnect_restore_behavior
rejected_for_now:
  - direct_nakama_api_compatibility
```

Pitaya reference mapping:

```yaml
adopted_concepts:
  - acceptor_session_route_handler_separation
  - bind_user_identity_to_session_like_connection_context
  - groups_rooms_and_broadcast_need_bound_identity_later
adapted_concepts:
  - vibit_binding_is_application_owned_and_protocol_explicit
  - vibit_single_process_registry_precedes_cluster_session_broadcast
deferred_concepts:
  - frontend_backend_split
  - server_to_server_rpc
  - remote_session_bind_broadcast
  - direct_group_membership_api
rejected_for_now:
  - direct_pitaya_api_compatibility
```

Nakama and Pitaya guide capability planning. They do not override vibit's constitution, manifests, ADRs, generated boundaries, or verification commands.

## 9. Required Future Gates

Before first-message binding behavior exists, a later implementation gate must define:

- Exact Protobuf message fields and field numbers.
- Whether `authentication.proto` is extended or a new authentication protocol source is created.
- Generated output update path.
- Protocol adapter recognition and response mapping.
- Application connection binding registry shape.
- Startup composition of the binder.
- Route policy use of bound identity.
- Focused tests for missing, malformed, invalid, expired, revoked, and unavailable proof.
- Redaction rules for raw tokens, token digests, key ids, connection ids, and internal binding state.
- Close/error behavior for failed binding attempts.

Before durable or multi-connection behavior exists, separate future gates must define:

- PostgreSQL session schema.
- Logout, revocation, and active-connection invalidation.
- Reconnect and connection epoch behavior.
- Single-socket, duplicate connection, kick, and replacement behavior.
- Presence, room, party, match, and group membership attachment.
- Operations, metrics, and abuse-control posture.

## 10. Deferrals

This gate does not authorize:

- New `.proto` messages.
- Generated Protobuf output.
- Changes to `proto/vibit/protocol/v1/envelope.proto`.
- WebSocket handshake authentication.
- HTTP `Authorization`, Bearer, cookie, query-string, or subprotocol proof carriers.
- Connection-bound identity registry code.
- Route policy use of bound identity.
- Session persistence.
- Session tables, repositories, PostgreSQL adapters, migrations, or cleanup jobs.
- Logout, refresh, token rotation, token validation audit mutation, or active-connection invalidation.
- Reconnect, resume, connection epoch replacement, or duplicate-connection policy.
- New dependencies.
- Memory durable authentication behavior.
- Direct Nakama or Pitaya public API compatibility.

## 11. Verification

The repository check rule for this gate is:

```text
runtime.first_message_connection_binding_gate
```

The check must verify:

- The standard, translation, ADR, change specs, and conversation log exist.
- Runtime, conventions, contracts, reference, work-item, module, and AGENTS markers exist.
- WebSocket transport non-test Go files remain credential-neutral.
- The existing Protobuf envelope remains free of authentication proof and binding fields.
- No `BindConnection` Protobuf source or generated output has been added before an implementation gate.
- No session or connection-binding migration source is present.

## 12. References

- Nakama sockets: `https://heroiclabs.com/docs/nakama/concepts/sockets/`
- Nakama sessions: `https://heroiclabs.com/docs/nakama/concepts/session/`
- Nakama configuration session limits: `https://heroiclabs.com/docs/nakama/getting-started/configuration/`
- Pitaya session package documentation: `https://pkg.go.dev/github.com/topfreegames/pitaya/v3/pkg/session`
