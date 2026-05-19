# First Message Connection Binding Implementation Gate

Status: Draft v0.1
Last updated: 2026-05-17
Scope: Future bounded implementation slice for `runtime.authentication.BindConnection`
Depends on: `docs/first-message-connection-binding-gate.md`, `docs/session-persistence-websocket-handshake-ratification.md`, `docs/access-token-protocol-carrier-route-protection-gate.md`, `docs/authentication-command-protocol-login-route-gate.md`, `docs/runtime-authentication-startup-composition-gate.md`, `docs/runtime-protocol-adapter.md`, `docs/generated-output.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0058`

The paired Simplified Chinese translation is `docs/first-message-connection-binding-implementation-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

`ADR-0057` selected first-message protocol/application binding as the future connection-level identity posture. The selected route is:

```text
runtime.authentication.BindConnection
```

This gate defines what a later bounded implementation slice may add, and what it still must not add.

The implementation gate is:

```yaml
first_message_connection_binding_implementation_gate: defined
implementation_authorized_by_this_standard: false
completed_gate_work_item: W-0126
future_implementation_work_item: W-0128
decision: ADR-0058
check_rule: runtime.first_message_connection_binding_implementation_gate
future_route: runtime.authentication.BindConnection
future_route_kind: system
future_protocol_source: proto/vibit/authentication/v1/authentication.proto
future_generated_go_output: runtime/internal/generated/proto/vibit/authentication/v1/authentication.pb.go
future_protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
future_application_binding_owner: runtime/internal/app
future_startup_owner: runtime/cmd/vibit-server
future_request: vibit.authentication.v1.BindConnectionRequest
future_response: vibit.authentication.v1.BindConnectionResponse
future_status_enum: vibit.authentication.v1.ConnectionBindingStatus
first_composed_runtime_store: postgres
memory_store_binding_status: unavailable_bootstrap
websocket_transport_credential_neutral: true
websocket_handshake_authentication_added: false
existing_protobuf_envelope_change_added: false
session_persistence_added: false
route_policy_bound_identity_added_by_this_gate: false
```

This is a gate-only standard. It does not implement connection binding.

## 2. Future Implementation Boundary

A later implementation slice may add only these behavior families:

```yaml
future_allowed_implementation:
  protocol_source:
    - extend proto/vibit/authentication/v1/authentication.proto with BindConnection messages
  generated_output:
    - regenerate runtime/internal/generated/proto/vibit/authentication/v1/authentication.pb.go through Buf
  protocol_adapter:
    - recognize runtime.authentication.BindConnection as a system route
    - decode BindConnectionRequest payload
    - map public binding result to BindConnectionResponse
    - map public binding errors to error envelopes
  application_binding:
    - validate access-token proof through existing application authentication service boundary
    - create normalized RequestIdentity for the server-observed connection id
    - keep an in-memory process-local connection binding registry if explicitly implemented in the slice
  startup_composition:
    - wire the binder only when the PostgreSQL authentication service is composed
  tests:
    - protocol adapter tests
    - application binder tests
    - startup composition tests
    - WebSocket transport neutrality regression tests
```

The future implementation must not move authentication into transport, generated code, repositories, domain modules, or Protobuf envelope metadata.

## 3. Future Protocol Shape

The future implementation may extend:

```text
proto/vibit/authentication/v1/authentication.proto
```

Planned messages:

```yaml
BindConnectionRequest:
  fields:
    access_token:
      type: string
      secret: true
      semantics: opaque Base64URL unpadded 32-byte access-token proof
    client_instance_id:
      type: string
      secret: false
      semantics: optional client installation or runtime instance hint

BindConnectionResponse:
  fields:
    binding_status:
      type: ConnectionBindingStatus
      semantics: bound or rejected
    actor_kind:
      type: string
      semantics: player when bound
    player_id:
      type: string
      semantics: validated player id when bound
    connection_id:
      type: string
      semantics: server-observed connection id
    connection_epoch:
      type: uint64
      semantics: server-observed connection epoch
    session_validated:
      type: bool
      semantics: false until durable session persistence is implemented
    bound_at:
      type: string
      semantics: server timestamp in RFC3339Nano UTC text

ConnectionBindingStatus:
  values:
    - CONNECTION_BINDING_STATUS_UNSPECIFIED
    - CONNECTION_BINDING_STATUS_BOUND
    - CONNECTION_BINDING_STATUS_REJECTED
```

Rules:

- The route kind is `system`.
- The route key is `runtime.authentication.BindConnection`.
- The request payload type is `vibit.authentication.v1.BindConnectionRequest`.
- The response payload type is `vibit.authentication.v1.BindConnectionResponse`.
- The existing `proto/vibit/protocol/v1/envelope.proto` must remain unchanged.
- Generated Go Protobuf output must be produced by Buf and must not be hand-edited.
- `access_token` is secret material and must not appear in errors, logs, events, debug details, or repository records.

## 4. Future Binding Flow

The future implementation must preserve this flow:

```yaml
bind_connection_flow:
  - websocket_transport_accepts_connection_without_reading_credentials
  - websocket_transport_assigns_or_observes_connection_id
  - client_may_call_public_login_route_if_access_token_is_needed
  - client_sends_runtime.authentication.BindConnection_system_message
  - protobuf_adapter_decodes_existing_envelope_and_bind_payload
  - application_binding_service_validates_access_token_through_authentication_service
  - application_binding_service_builds_validated_player_identity_for_connection
  - application_binding_service_records_process_local_binding_if_registry_is_in_slice
  - protocol_adapter_returns_bind_response_without_secret_material
```

Rules:

- Binding validates access-token proof through the existing application authentication service path.
- Binding must use the server-observed `connection_id` and `connection_epoch`, not client-supplied connection identity.
- Binding success may create connection-bound identity, but route policy may use that identity only if the implementation slice explicitly updates route policy.
- `RequestIdentity.SessionValidated` remains false until durable session persistence is separately implemented.
- Domain handlers must not query connection binding registries directly.

## 5. Future Application Ownership

The future application boundary should use app-owned types under:

```text
runtime/internal/app
```

Candidate files:

```yaml
application_binding_source:
  - runtime/internal/app/connection_binding.go
  - runtime/internal/app/connection_binding_test.go
optional_authentication_adapter_source:
  - runtime/internal/app/authentication/connection_binding_validator.go
  - runtime/internal/app/authentication/connection_binding_validator_test.go
```

Candidate application vocabulary:

```yaml
BindConnectionRequest:
  access_token: string
  route: RouteKey
  connection_id: string
  connection_epoch: uint64
  client_instance_id: string

BindConnectionResult:
  bound: bool
  identity: RequestIdentity
  binding_status: bound | rejected
  public_error_code: ErrorCode
  connection_id: string
  connection_epoch: uint64
  bound_at: time.Time
```

Rules:

- The binder may depend on an interface that validates access tokens; it must not import platform persistence or protocol packages.
- The binder must collapse public token failures into binding-specific public error codes.
- The binder must not generate tokens, revoke tokens, refresh tokens, mutate token audit state, or create session records.
- If a registry is added, it must be process-local only in the first implementation and must not claim durability.

## 6. Future Protocol Adapter Ownership

The future Protobuf adapter boundary should use files under:

```text
runtime/internal/platform/protocol/protobuf
```

Candidate files:

```yaml
protocol_adapter_source:
  - runtime/internal/platform/protocol/protobuf/connection_binding.go
  - runtime/internal/platform/protocol/protobuf/connection_binding_test.go
```

Rules:

- The protocol adapter may decode `BindConnectionRequest`.
- The protocol adapter may map binding results to `BindConnectionResponse`.
- The protocol adapter may map public binding errors to existing error envelopes.
- The protocol adapter must not validate token digests, query repositories, load verifier keys, create sessions, or decide token lifecycle state.
- The adapter must not change normal request-level `vibit.authentication.v1.AuthenticatedRequest` route protection behavior.

## 7. Startup Composition

Future startup composition is limited to:

```yaml
first_runtime_store: postgres
startup_owner: runtime/cmd/vibit-server
memory_store_binding_status: unavailable_bootstrap
```

Rules:

- PostgreSQL runtime startup may wire connection binding only when the existing authentication service and route protector are already composed.
- Memory runtime startup may keep `BindConnection` unavailable.
- Startup must not parse access tokens from process arguments, environment variables, HTTP headers, cookies, query strings, or WebSocket subprotocols.
- Startup must not apply migrations automatically.
- Startup must not register durable session stores unless a later session persistence gate authorizes them.

## 8. Error Mapping

Future implementation must use binding-specific public errors:

```yaml
future_public_errors:
  missing_bind_proof: CONNECTION_BINDING_TOKEN_MISSING
  malformed_bind_payload_or_token: CONNECTION_BINDING_TOKEN_MALFORMED
  invalid_or_expired_or_revoked_token: CONNECTION_BINDING_TOKEN_INVALID
  validation_dependency_unavailable: CONNECTION_BINDING_UNAVAILABLE
  protected_bound_route_without_binding: CONNECTION_BINDING_REQUIRED
```

Rules:

- Public errors must not disclose lookup hit or miss, verifier key id, verifier mismatch, token status, player account state, audience mismatch, connection registry state, or internal dependency class.
- A failed bind must not bind identity to the connection.
- A failed bind must not return token record ids, credential record ids, lookup digests, verifier digests, or raw proof.
- Rebinding an already bound connection, duplicate player connections, kick/replace behavior, and repeated failure close policy remain separate gates.

## 9. Nakama And Pitaya Reference Mapping

Nakama reference mapping:

```yaml
adopted_concepts:
  - authenticated_session_material_precedes_authenticated_realtime_socket_use
  - realtime_socket_lifecycle_can_have_authenticated_player_context
  - session_token_lifetime_and_socket_connection_lifetime_are_related_but_distinct
adapted_concepts:
  - vibit_uses_opaque_access_token_not_nakama_session_api_compatibility
  - vibit_binds_connection_through_existing_protobuf_system_route_not_transport_handshake
  - vibit_keeps_request_level_validation_as_current_protected_route_path
deferred_concepts:
  - refresh_token_flow
  - single_session_or_single_socket_enforcement
  - disconnect_on_session_revocation
  - reconnect_restore_behavior
rejected_for_now:
  - direct_nakama_api_compatibility
```

Pitaya reference mapping:

```yaml
adopted_concepts:
  - transport_acceptor_separated_from_session_binding
  - user_identity_can_bind_to_session_like_connection_context
  - route_handlers_should_receive_context_not_parse_transport_credentials
adapted_concepts:
  - vibit_connection_binding_is_application_owned_and_protocol_explicit
  - vibit_first_registry_is_process_local_before_cluster_session_broadcast
  - vibit_routes_remain_contract_first_instead_of_pitaya_route_string_compatibility
deferred_concepts:
  - frontend_backend_cluster_split
  - remote_session_binding_broadcast
  - groups_rooms_and_presence_attachment
  - server_to_server_rpc
rejected_for_now:
  - direct_pitaya_api_compatibility
```

Nakama and Pitaya guide capability planning. They do not override vibit's constitution, manifests, ADRs, generated boundaries, or verification commands.

## 10. Required Future Tests

The future implementation slice must add focused tests for:

```yaml
required_tests:
  protocol_source_includes_bind_connection_messages
  generated_output_is_buf_generated_and_traced_to_proto_source
  bind_connection_payload_decodes_to_application_request
  bind_connection_success_maps_to_response_without_secret_material
  bind_connection_missing_token_maps_to_CONNECTION_BINDING_TOKEN_MISSING
  bind_connection_malformed_token_maps_to_CONNECTION_BINDING_TOKEN_MALFORMED
  bind_connection_invalid_token_maps_to_CONNECTION_BINDING_TOKEN_INVALID
  bind_connection_unavailable_dependency_maps_to_CONNECTION_BINDING_UNAVAILABLE
  bind_connection_uses_server_observed_connection_id
  bind_connection_keeps_session_validated_false
  failed_bind_does_not_create_bound_identity
  public_login_route_remains_available_before_binding
  request_level_authenticated_wrapper_remains_current_protected_route_path
  websocket_transport_remains_credential_neutral
  existing_protobuf_envelope_remains_unchanged
```

Live PostgreSQL verification remains optional unless a later change explicitly requires an opt-in live check.

## 11. Deferrals

This gate does not authorize:

- Implementing `BindConnection`.
- Adding `.proto` messages.
- Running `buf generate`.
- Editing generated `.pb.go` files.
- Changing `proto/vibit/protocol/v1/envelope.proto`.
- WebSocket handshake authentication.
- HTTP `Authorization`, Bearer, cookie, query-string, or subprotocol proof carriers.
- Durable session persistence.
- Session tables, session repositories, session migrations, or session cleanup jobs.
- Route-policy use of bound identity without the future implementation slice explicitly changing route policy.
- Logout-triggered active-connection invalidation.
- Refresh behavior, token rotation, or token validation audit mutation.
- Reconnect, resume, duplicate connection replacement, or connection epoch policy.
- Presence, rooms, parties, match runtime, group membership, or broadcast behavior.
- New dependencies.
- Memory durable authentication behavior.
- Direct Nakama or Pitaya public API compatibility.

## 12. Verification

The repository check rule for this gate is:

```text
runtime.first_message_connection_binding_implementation_gate
```

The check must verify:

- The standard, translation, ADR, change specs, and conversation log exist.
- Runtime, conventions, contracts, reference, work-item, module, and AGENTS markers exist.
- The gate declares the future implementation artifacts and deferrals.
- WebSocket transport non-test Go files remain credential-neutral.
- The existing Protobuf envelope remains free of authentication proof and binding fields.
- No `BindConnection` Protobuf source or generated output has been added by this gate-only change.
- No session or connection-binding migration source is present.

## 13. References

- Nakama sockets: `https://heroiclabs.com/docs/nakama/concepts/sockets/`
- Nakama sessions: `https://heroiclabs.com/docs/nakama/concepts/session/`
- Nakama configuration session limits: `https://heroiclabs.com/docs/nakama/getting-started/configuration/`
- Pitaya session package documentation: `https://pkg.go.dev/github.com/topfreegames/pitaya/v3/pkg/session`
