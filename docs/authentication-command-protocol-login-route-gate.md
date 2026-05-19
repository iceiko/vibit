# Authentication Command Protocol And Login Route Gate

Status: Draft v0.1
Last updated: 2026-05-17
Scope: Future Protobuf command payloads, application route registration, protocol bridge, and startup composition for the public `AuthenticateWithDeviceCredential` route
Depends on: `docs/authentication-contract-error-permission-surfaces.md`, `docs/device-credential-login-service-behavior-gate.md`, `docs/access-token-protocol-carrier-route-protection-gate.md`, `docs/runtime-authentication-startup-composition-gate.md`, `docs/runtime-protocol-adapter.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0055`

The paired Simplified Chinese translation is `docs/authentication-command-protocol-login-route-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

The runtime can now execute device credential login inside the application authentication service, validate opaque access tokens for protected routes, and wire route protection into the PostgreSQL runtime startup path.

The remaining client-facing gap is the public login command route. Clients need a protocol message and registered application route for:

```text
runtime.authentication.AuthenticateWithDeviceCredential
```

This gate defines the next bounded implementation slice. It may add only the public device-credential login command protocol messages, generated Protobuf output, protocol bridge behavior, application route handler registration, and PostgreSQL startup composition needed to expose the already-implemented service method.

It does not add session persistence, WebSocket handshake authentication, HTTP `Authorization` or Bearer carriers, cookies, query-string carriers, WebSocket subprotocol carriers, logout, refresh, cleanup, token rotation, repository interface changes, migration changes, new dependencies, or direct Nakama/Pitaya public API compatibility.

## 2. Core Rule

The authentication command protocol and login route gate is:

```yaml
authentication_command_protocol_login_route_gate: defined
implementation_authorized_by_this_standard: true
completed_gate_work_item: W-0119
future_implementation_work_item: W-0120
decision: ADR-0055
public_login_route: runtime.authentication.AuthenticateWithDeviceCredential
route_kind: command
route_public_policy_status: already_public_in_route_protection_policy
first_protocol_source: proto/vibit/authentication/v1/authentication.proto
first_generated_go_output: runtime/internal/generated/proto/vibit/authentication/v1/authentication.pb.go
application_handler_owner: runtime/internal/app/bootstrap
authentication_service_owner: runtime/internal/app/authentication
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
startup_owner: runtime/cmd/vibit-server
first_composed_runtime_store: postgres
memory_store_login_route_status: unavailable_bootstrap
websocket_transport_credential_neutral: true
protobuf_envelope_change_status: unchanged
session_persistence_status: deferred
websocket_handshake_authentication_status: deferred
```

Future implementation may expose only the existing device credential login service behavior through a public command route. It must not implement new authentication semantics or move authentication logic into WebSocket transport, Protobuf envelope metadata, domain modules, repositories, migrations, or generated files.

## 3. Selected Protocol Shape

The first authentication command protocol source is:

```text
proto/vibit/authentication/v1/authentication.proto
```

The first planned messages are:

```yaml
messages:
  AuthenticateWithDeviceCredentialRequest:
    fields:
      credential_proof: string
      requested_player_id: string
      client_instance_id: string
      account_creation_intent: AccountCreationIntent
  AuthenticateWithDeviceCredentialResponse:
    fields:
      authentication_status: string
      actor_kind: string
      player_id: string
      access_token: string
      token_type: string
      issued_at: string
      expires_at: string
      token_record_id: string
  AccountCreationIntent:
    values:
      - ACCOUNT_CREATION_INTENT_UNSPECIFIED
      - ACCOUNT_CREATION_INTENT_ALLOW_CREATE
      - ACCOUNT_CREATION_INTENT_AUTHENTICATE_EXISTING_ONLY
```

Rules:

- The existing `proto/vibit/protocol/v1/envelope.proto` must remain unchanged.
- The envelope route remains `kind=command`, `module=runtime.authentication`, `name=AuthenticateWithDeviceCredential`.
- The payload type is `vibit.authentication.v1.AuthenticateWithDeviceCredentialRequest`.
- The response payload type is `vibit.authentication.v1.AuthenticateWithDeviceCredentialResponse`.
- `credential_proof` and `access_token` are secret values and must not appear in errors, logs, events, or repository records.
- `access_token` is one-time client-visible response material from the existing service result; startup or protocol adapter code must not store it.
- `token_type` must reflect the service-local posture `opaque_access_token`.
- Time values must use RFC3339 or RFC3339Nano UTC text.

## 4. Route Registration Flow

Future implementation must preserve this layer sequence:

```yaml
login_route_flow:
  - websocket_transport_receives_binary_frame_without_reading_credentials
  - protobuf_adapter_decodes_existing_envelope
  - route_protector_allows_public_authentication_route_without_access_token
  - protobuf_adapter_decodes AuthenticateWithDeviceCredentialRequest
  - protocol_bridge_maps_request_to authentication.DeviceCredentialAuthenticationRequest
  - application_bootstrap_handler_calls authentication.Service.AuthenticateWithDeviceCredential
  - authentication_service_owns_unit_of_work_and token issuance
  - protocol_bridge_maps AuthenticationResult to AuthenticateWithDeviceCredentialResponse
  - protobuf_adapter_returns success or existing error envelope
```

Rules:

- The route must be explicit. No implicit public route family is created.
- Application route registration belongs in `runtime/internal/app/bootstrap` or an equivalent application-composition package, not in the authentication module.
- The handler calls only the existing `AuthenticateWithDeviceCredential` service method.
- The handler must not compute digests, compare verifiers, generate tokens, call repositories directly, open transactions, or parse transport metadata as proof.
- Because the authentication service owns its own unit-of-work boundary, the application transaction wrapper must not open an outer inventory-style command transaction for the public authentication route.
- The memory runtime path may keep the login route unavailable because it has no durable authentication repository capability.

## 5. Error Mapping

Future implementation must map service public errors through the existing application error envelope boundary.

First mapping:

```yaml
service_public_errors:
  AUTHENTICATION_PROOF_MISSING: application_error_same_code
  AUTHENTICATION_PROOF_MALFORMED: application_error_same_code
  AUTHENTICATION_CREDENTIAL_INVALID: application_error_same_code
  AUTHENTICATION_CREDENTIAL_STORE_UNAVAILABLE: application_error_same_code
  AUTHENTICATION_TOKEN_STORE_UNAVAILABLE: application_error_same_code
  AUTHENTICATION_NOT_IMPLEMENTED: application_error_same_code
```

Rules:

- Public errors must not disclose credential lookup hit or miss, player account existence, verifier key id, digest class, HMAC input/output, or token record internals beyond audit-safe ids already in success output.
- Error messages must not include `credential_proof`, `access_token`, lookup digests, verifier digests, raw HMAC input, verifier keys, or full concrete key ids.
- A failed login must not return `access_token` or `token_record_id`.

## 6. Nakama And Pitaya Reference Mapping

Nakama guides the capability sequence: clients authenticate first and receive session/token material before using normal gameplay or realtime features. vibit adapts that by exposing the public login command before session persistence or handshake authentication.

Pitaya guides layering: acceptors/connections, session context, routes, and handlers stay separated. vibit adapts that by keeping WebSocket transport credential-neutral, using the existing envelope route, bridging Protobuf payloads in the protocol adapter, and invoking application-owned handlers.

Neither reference overrides vibit's boundaries. This gate does not adopt direct Nakama or Pitaya public API compatibility.

## 7. Required Implementation Tests

The future implementation slice must add or update focused tests for:

```yaml
required_tests:
  proto_source_and_generated_output_exist
  login_route_is_registered_only_when_authentication_service_is_composed
  login_route_is_public_and_does_not_require_access_token_wrapper
  login_route_bypasses_outer_transactional_dispatcher_unit_of_work
  login_request_maps_to_service_request_without_treating_metadata_as_proof
  login_success_maps_service_result_to_response_payload
  login_failure_maps_public_service_error_to_error_envelope
  login_errors_do_not_leak_credential_proof_or_access_token
  protected_routes_still_require_authenticated_wrapper
  websocket_transport_remains_credential_neutral
  existing_protobuf_envelope_remains_unchanged
```

Live PostgreSQL login verification remains optional and must not be required by default repository checks.

## 8. Deferrals

This gate does not authorize:

- Session persistence.
- Runtime session tables or migrations.
- WebSocket handshake authentication.
- HTTP `Authorization`, Bearer, cookie, query string, or WebSocket subprotocol credential carriers.
- Logout, refresh, cleanup, token rotation, token validation audit mutation, or previous-token revocation.
- Repository interface changes.
- PostgreSQL adapter changes.
- Migration changes.
- Automatic startup migrations.
- New external dependencies.
- Memory-store durable authentication behavior.
- Direct Nakama or Pitaya API compatibility.

## 9. Verification

The repository check rule for this gate is:

```text
runtime.authentication_command_protocol_login_route_gate
```
