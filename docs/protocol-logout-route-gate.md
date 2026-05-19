# Protocol Logout Route Gate

Status: Draft v0.1
Last updated: 2026-05-18
Scope: Gate-only boundary for future client-facing `runtime.authentication.LogoutAccessToken` protocol route behavior
Depends on: `docs/nakama-pitaya-product-parity-roadmap.md`, `docs/logout-access-token-behavior-gate.md`, `decisions/ADR-0073-logout-access-token-behavior-implementation.md`, `docs/authentication-command-protocol-login-route-gate.md`, `docs/websocket-close-policy-gate.md`, `docs/runtime-protocol-adapter.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0079`

The paired Simplified Chinese translation is `docs/protocol-logout-route-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

vibit can now revoke a verified presented opaque access token through `authentication.Service.LogoutAccessToken`, but clients do not yet have a protocol route for requesting that behavior.

The Nakama/Pitaya product parity roadmap makes logout part of the near-term runtime lifecycle closure. Nakama's server runtime exposes session logout and session disconnect as distinct operations, which reinforces that token/session invalidation and socket disconnection are separate lifecycle actions. Pitaya's docs emphasize sessions, handler routing, acceptors, and kick/disconnect style behavior as separate framework surfaces. vibit adapts those lessons into a protocol logout route boundary that exposes token logout without silently closing sockets or moving token logic into transport.

This standard defines the next bounded protocol route implementation slice. It is gate-only for this work item.

```yaml
protocol_logout_route_gate: defined
implementation_authorized_by_this_standard: true
completed_gate_work_item: W-0169
future_implementation_work_item: W-0170
decision: ADR-0079
check_rule: runtime.protocol_logout_route_gate
public_logout_route: runtime.authentication.LogoutAccessToken
route_kind: command
route_protection_posture: explicit_service_validated_token_lifecycle_route
proof_carrier_posture: access_token_in_logout_request_payload
first_protocol_source: proto/vibit/authentication/v1/authentication.proto
first_generated_go_output: runtime/internal/generated/proto/vibit/authentication/v1/authentication.pb.go
request_message_candidate: vibit.authentication.v1.LogoutAccessTokenRequest
response_message_candidate: vibit.authentication.v1.LogoutAccessTokenResponse
application_handler_owner: runtime/internal/app/bootstrap
authentication_service_owner: runtime/internal/app/authentication
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
startup_owner: runtime/cmd/vibit-server
first_composed_runtime_store: postgres
memory_store_logout_route_status: unavailable_bootstrap
transaction_bypass_required: true
protobuf_envelope_change_status: unchanged
websocket_transport_credential_neutral: true
first_logout_scope: presented_access_token_only
already_revoked_token_public_behavior: invalid_token
expired_token_public_behavior: invalid_token
logout_closes_socket: false
runtime_session_revocation_added: false
active_connection_invalidation_added: false
concrete_socket_close_added: false
reconnect_epoch_behavior_added: false
protocol_session_carrier_added: false
dependencies_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Ownership

Future protocol logout route implementation must preserve these ownership boundaries:

```yaml
authentication_service_owner: runtime/internal/app/authentication
application_handler_owner: runtime/internal/app/bootstrap
route_policy_owner: runtime/internal/app
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
startup_owner: runtime/cmd/vibit-server
websocket_transport_owner: runtime/internal/platform/transport/ws
close_policy_owner: runtime/internal/app/connection
```

Rules:

- The authentication service remains the only owner of raw access-token proof validation, digest computation, verifier comparison, token lookup, token posture checks, and token revocation mutation.
- The future route handler may call only the existing `LogoutAccessToken` service method. It must not compute digests, compare verifiers, call repositories directly, open its own transaction, close sockets, or mutate runtime sessions.
- The Protobuf adapter may map request and response payloads only. It must not decide logout validity, route authorization, session revocation, active connection invalidation, or transport close behavior.
- WebSocket transport remains credential-neutral and policy-neutral. It must not parse logout proof from headers, cookies, query strings, bearer values, or subprotocols.
- The close policy and active connection registry remain separate from protocol logout route exposure.

## 3. Route And Carrier Posture

The first logout route is:

```yaml
route:
  kind: command
  module: runtime.authentication
  name: LogoutAccessToken
  semantic_contract: contracts/runtime/authentication/commands/LogoutAccessToken.yaml
```

The first route-protection posture is:

```yaml
route_protection_posture: explicit_service_validated_token_lifecycle_route
authenticated_request_wrapper_required: false
access_token_source: LogoutAccessTokenRequest.access_token
service_validation_required: true
```

Rationale:

- Logout needs the exact presented access token so the service can revoke the same token record after verifier comparison.
- The route is not a normal gameplay protected route. It is a token lifecycle route whose payload carries the proof being revoked.
- Marking the route as explicit and service-validated prevents the route protector from consuming the token before the logout service can revoke it.
- The route must still be explicitly registered; there is no implicit public authentication route family.

## 4. Protocol Shape

The future implementation may extend the existing authentication Protobuf source:

```text
proto/vibit/authentication/v1/authentication.proto
```

Planned request:

```yaml
LogoutAccessTokenRequest:
  access_token:
    type: string
    secret: true
    source: presented_opaque_access_token
  logout_reason:
    type: string
    required: false
    public_safe: true
```

Planned response:

```yaml
LogoutAccessTokenResponse:
  logout_status:
    type: string
    values:
      - revoked
      - rejected
  revoked:
    type: bool
  logout_scope:
    type: string
    first_value: presented_access_token
  revoked_at:
    type: string
    format: rfc3339nano_utc
  token_record_id:
    type: string
    visibility: audit_safe
```

Rules:

- The existing `proto/vibit/protocol/v1/envelope.proto` must remain unchanged.
- `access_token` is secret input and must never appear in errors, logs, events, close intents, connection records, or test names.
- `logout_reason` must not be treated as trusted server policy input and must not carry secrets.
- `logout_scope` must expose the first posture as `presented_access_token`, even if internal service vocabulary uses a narrower implementation enum.
- Time values must use RFC3339Nano UTC text.
- `token_record_id`, if exposed, is audit-safe only. It is not proof, not a session id, and not a connection target.

## 5. Future Route Flow

A future implementation must preserve this sequence:

```yaml
logout_route_flow:
  - websocket_transport_receives_binary_frame_without_reading_credentials
  - protobuf_adapter_decodes_existing_envelope
  - route_policy_allows_explicit_service_validated_logout_route_without_authenticated_wrapper
  - protobuf_adapter_decodes LogoutAccessTokenRequest
  - protocol_bridge_maps_request_to authentication.LogoutAccessTokenRequest
  - application_bootstrap_handler_calls authentication.Service.LogoutAccessToken
  - authentication_service_owns_unit_of_work_and token revocation
  - protocol_bridge_maps LogoutAccessTokenResult to LogoutAccessTokenResponse
  - protobuf_adapter_returns success or existing error envelope
```

Rules:

- The authentication route must bypass the outer `TransactionalDispatcher` unit of work because the authentication service owns its own unit-of-work boundary.
- A failed logout must not return `revoked: true`, `revoked_at`, or a successful status.
- A successful logout must be returned only after the service reports successful unit-of-work commit.
- The future route must not decode or dispatch an `AuthenticatedRequest` wrapper for logout.

## 6. Public Error Mapping

Future protocol behavior must map service public errors without leaking proof details:

```yaml
service_public_errors:
  AUTHENTICATION_TOKEN_MISSING: application_error_same_code
  AUTHENTICATION_TOKEN_MALFORMED: application_error_same_code
  AUTHENTICATION_TOKEN_INVALID: application_error_same_code
  AUTHENTICATION_TOKEN_STORE_UNAVAILABLE: application_error_same_code
  AUTHENTICATION_NOT_IMPLEMENTED: application_error_same_code
```

Rules:

- Lookup miss, already revoked token, expired token, wrong audience, wrong kind, unsupported verifier metadata, unknown verifier key id, and verifier mismatch must remain collapsed to public invalid-token behavior.
- Error messages must not include raw token text, lookup digest, verifier digest, HMAC input, verifier key id, headers, cookies, query strings, subprotocol values, session ids, connection ids, remote addresses, or database errors.
- Memory runtime may keep durable logout route behavior unavailable until a later memory authentication posture is ratified.

## 7. Relationship To Socket Close And Sessions

This gate does not change socket or session lifecycle behavior:

```yaml
logout_closes_socket: false
close_policy_called_by_logout_route: false
active_connection_invalidation_added: false
runtime_session_revocation_added: false
bound_connection_identity_after_logout_policy: deferred
bound_session_identity_after_logout_policy: deferred
reconnect_epoch_behavior_added: false
protocol_session_carrier_added: false
```

Rules:

- Successful protocol logout revokes only the verified presented access-token record.
- Successful protocol logout must not implicitly close the WebSocket connection.
- Successful protocol logout must not implicitly revoke the linked runtime session.
- Successful protocol logout must not invalidate active connection registry records.
- Existing request-level token validation will reject the revoked token on later protected requests when validation runs.
- Whether a bound connection can continue using already-bound identity after logout remains a later bound-identity/session/reconnect policy question.

## 8. Required Future Tests

The future implementation slice must add focused tests for:

```yaml
required_tests:
  proto_source_and_generated_output_include_logout_messages
  logout_route_is_registered_only_when_authentication_service_is_composed
  logout_route_is_explicit_service_validated_token_lifecycle_route
  logout_route_bypasses_outer_transactional_dispatcher_unit_of_work
  logout_request_maps_access_token_to_service_request_without_logging_it
  logout_success_maps_service_result_to_response_payload
  logout_failure_maps_public_service_error_to_error_envelope
  logout_errors_do_not_leak_access_token
  websocket_transport_remains_credential_neutral
  existing_protobuf_envelope_remains_unchanged
  logout_route_does_not_call_close_policy_or_close_socket
  protected_gameplay_routes_still_require_authenticated_wrapper
```

Live PostgreSQL verification remains optional and must not be required by default repository checks.

## 9. Nakama And Pitaya Reference Mapping

Nakama reference mapping:

- Adopt the distinction between logout/token invalidation and explicit session disconnect.
- Adapt the product expectation that clients can request logout through a stable server API.
- Defer Nakama session token compatibility, refresh token compatibility, logout-all semantics, realtime socket disconnect compatibility, and dashboard/admin compatibility.

Pitaya reference mapping:

- Adopt the separation between client connection acceptors, sessions, routes, handlers, and connection management.
- Adapt the session/kick lifecycle lesson by keeping logout proof validation in application service behavior and socket close in a separate future policy/handoff.
- Defer Pitaya route naming compatibility, frontend/backend routing, distributed kick/disconnect, groups integration, and RPC/session propagation.

## 10. Non-Goals

This gate does not authorize:

- Adding Protobuf logout messages in this work item.
- Adding generated Go Protobuf output in this work item.
- Registering `runtime.authentication.LogoutAccessToken` in this work item.
- Changing the existing Protobuf envelope.
- Changing WebSocket handshake authentication.
- Parsing HTTP headers, bearer values, cookies, query strings, or WebSocket subprotocols.
- Closing sockets, choosing close codes, choosing close reason text, or sending protocol close messages.
- Revoking runtime sessions.
- Invalidating active connection registry records.
- Adding reconnect, resume, duplicate replacement, or connection epoch behavior.
- Adding protocol session carriers.
- Adding refresh, logout-all, admin revocation, cleanup jobs, presence, chat, social modules, matchmaking, match runtime, SDKs, cluster, RPC, service discovery, dependencies, or direct Nakama/Pitaya API compatibility.

## 11. Verification

The repository check rule for this gate is:

```text
runtime.protocol_logout_route_gate
```
