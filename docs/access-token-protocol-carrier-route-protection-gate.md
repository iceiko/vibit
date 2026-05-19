# Access Token Protocol Carrier And Route Protection Gate

Status: Draft v0.1
Last updated: 2026-05-17
Scope: Future access-token proof carrier exposure, request-level route protection, application validation handoff, public route policy, redaction, tests, and deferrals before protocol or route-protection implementation is added
Depends on: `docs/access-token-validation-service-behavior-gate.md`, `docs/runtime-protocol-adapter.md`, `docs/game-protocol.md`, `docs/session-persistence-websocket-handshake-decision-gates.md`
Canonical decision: `ADR-0053`

The paired Simplified Chinese translation is `docs/access-token-protocol-carrier-route-protection-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

The application authentication service can now issue and validate opaque access tokens inside a service-local boundary. The next risk is exposing those tokens to clients by placing bearer parsing, cookies, query parameters, WebSocket handshake authentication, Protobuf envelope changes, or route protection directly in transport, protocol, or domain packages.

This gate selects the next milestone direction after service-local authentication: define a request-level access-token carrier and route-protection boundary before implementation.

This is a gate-only standard. It does not add `.proto` sources, generated files, protocol adapter code, route-protection code, startup wiring, session persistence, WebSocket handshake authentication, repository changes, migrations, logout, refresh, cleanup, dependencies, or broader production authentication behavior.

## 2. Core Rule

The access-token protocol carrier and route-protection gate is:

```yaml
access_token_protocol_carrier_route_protection_gate: defined
implementation_authorized_by_this_standard: false
future_implementation_work_item: W-0114
completed_gate_work_item: W-0113
planned_owner: runtime/internal/app
planned_protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
planned_transport_owner: runtime/internal/platform/transport/ws
planned_authentication_service_owner: runtime/internal/app/authentication
validation_model: request_level_validation
first_proof_carrier_posture: protobuf_payload_wrapper
first_wrapper_message_candidate: vibit.authentication.v1.AuthenticatedRequest
protobuf_envelope_change_status: deferred
websocket_handshake_authentication_status: deferred
session_persistence_status: deferred
startup_wiring_status: deferred
route_policy_status: defined_not_implemented
```

Future implementation may only expose token proof to normal gameplay routes after a later work item explicitly authorizes the exact protocol source, generated output, adapter handoff, application route policy, and startup composition slice.

## 3. Selected Carrier Posture

The first route-protection carrier posture is an explicit Protobuf payload wrapper for protected requests.

Planned semantic shape:

```yaml
authenticated_request_payload_wrapper:
  message_candidate: vibit.authentication.v1.AuthenticatedRequest
  owner: proto/vibit/authentication/v1
  outer_envelope_kind: original_command_or_query
  outer_envelope_module: original_domain_module
  outer_envelope_name: original_domain_route_name
  outer_payload_type: vibit.authentication.v1.AuthenticatedRequest
  inner_payload_type: original_payload_type
  inner_payload: original_payload_bytes
  proof_field: access_token
  proof_encoding: base64url_unpadded_32_byte_opaque_access_token
```

Rules:

- The existing Protobuf envelope shape must not change for the first carrier posture.
- `Envelope.session.player_id`, `Envelope.session.session_id`, `Envelope.session.connection_id`, and `Envelope.session.connection_epoch` remain metadata only.
- The access token must not be placed in `Session`, `Target`, route fields, `request_id`, error details, logs, or connection metadata.
- The access token must not be parsed from HTTP `Authorization`, `Bearer` strings, cookies, query strings, or WebSocket subprotocols in this posture.
- The wrapper is a protocol payload carrier only. It does not validate proof, select token records, compare digests, decide player account state, or construct domain responses.
- The inner payload type must be decoded only after the wrapper is accepted structurally and proof validation succeeds.

## 4. Request-Level Validation Flow

The future implementation must preserve this layer sequence:

```yaml
request_level_validation_flow:
  - websocket_transport_receives_binary_frame_without_reading_authentication_proof
  - protobuf_adapter_decodes_existing_envelope
  - protobuf_adapter_recognizes_authenticated_payload_wrapper_for_protected_route
  - protobuf_adapter_extracts_access_token_text_as_secret
  - protobuf_adapter_keeps_inner_payload_bytes_undispatched
  - application_route_protection_policy_requires_authenticated_identity_for_protected_route
  - application_authentication_validator_calls_ValidateAccessToken
  - authentication_service_validates_token_and_returns_RequestIdentity
  - application_handoff_sets_validated_identity_with_SessionValidated_false
  - protocol_adapter_decodes_inner_payload_for_original_route
  - application_dispatch_calls_domain_handler_with_validated_identity
```

Rules:

- Domain handlers must receive a normalized `RequestIdentity`; they must not receive or parse access-token proof.
- The route-protection layer must deny protected routes before domain dispatch when proof is missing, malformed, invalid, unavailable, or not structurally carried in the selected wrapper.
- Metadata-only identity must never satisfy protected route policy.
- `SessionValidated` remains false until a later session persistence gate is selected and implemented.
- WebSocket transport remains credential-neutral.
- Protocol adapter code may extract the wrapper field only as a narrow proof handoff to application validation.

## 5. Route Policy Gate

Before route protection implementation, future work must declare an application-owned route policy.

First policy requirements:

```yaml
route_policy:
  owner: runtime/internal/app
  public_routes:
    - runtime.authentication.AuthenticateWithDeviceCredential
  protected_route_default: authenticated_player_required
  protected_identity_requirement:
    identity_status: validated
    actor_kind: player
    player_id_validated: true
    session_validated: false_allowed_until_session_persistence
  domain_module_token_parsing: forbidden
```

Rules:

- Public routes must be explicit. There is no implicit public gameplay route default.
- Protected routes require `IdentityValidationValidated`, `ActorKindPlayer`, and `PlayerIDValidated`.
- Route policy may allow `SessionValidated: false` only because access-token validation is request proof and session persistence is deferred.
- Inventory permission policies may continue to enforce domain permissions, but they must not replace route-level authentication policy.
- The route policy must be testable without a live PostgreSQL server.

## 6. Error Mapping

Future implementation must map authentication validation failures before domain dispatch.

First public posture:

```yaml
route_protection_error_mapping:
  missing_wrapper_or_missing_token: AUTHENTICATION_TOKEN_MISSING
  malformed_wrapper_or_malformed_token: AUTHENTICATION_TOKEN_MALFORMED
  invalid_token_family: AUTHENTICATION_TOKEN_INVALID
  validation_dependency_unavailable: AUTHENTICATION_TOKEN_STORE_UNAVAILABLE
  protected_route_without_validated_identity: AUTHENTICATION_TOKEN_INVALID
```

Rules:

- Public protocol errors must not reveal token lookup hit or miss, token lifecycle status, verifier key id, audience mismatch, player existence, player account state, or verifier mismatch.
- Route-not-authorized failures must not include raw token text, inner payload bytes, lookup digests, verifier digests, HMAC input/output, verifier keys, or full concrete key ids.
- Error mapping to Protobuf error envelopes must reuse the existing application-error boundary or a later explicit mapping gate.

## 7. Required Future Artifacts

The later implementation slice must define or update these artifacts before route protection becomes active:

```yaml
required_future_artifacts:
  protocol_source: proto/vibit/authentication/v1/authenticated_request.proto
  generated_go_proto_output: runtime/internal/generated/proto/vibit/authentication/v1/authenticated_request.pb.go
  protocol_adapter_tests: runtime/internal/platform/protocol/protobuf/*authentication*_test.go
  application_route_policy_source: runtime/internal/app/*route*_auth*.go
  application_route_policy_tests: runtime/internal/app/*route*_auth*_test.go
  authentication_validator_adapter_source: runtime/internal/app/authentication/*validator*.go
  process_startup_wiring: deferred_until_separate_startup_work_item
```

This gate does not create those artifacts.

## 8. Required Tests

Future implementation must add focused tests for:

- Missing wrapper on protected routes.
- Missing access token in wrapper.
- Malformed access token in wrapper.
- Invalid token public failure collapse.
- Store unavailable public failure mapping.
- Metadata-only identity rejected for protected routes.
- Valid token produces validated player identity before domain dispatch.
- `SessionValidated` remains false after access-token validation.
- Public authentication route remains explicit and does not require token proof.
- WebSocket transport does not parse `Authorization`, cookies, query strings, subprotocols, or bearer values.
- Protobuf envelope fields remain metadata only.
- Raw access-token text and inner payload bytes are not leaked in errors.

## 9. Deferrals

This gate does not authorize:

- Adding `.proto` files.
- Generating Protobuf output.
- Changing `proto/vibit/protocol/v1/envelope.proto`.
- Changing WebSocket handshake authentication.
- Parsing HTTP headers, cookies, query strings, bearer strings, or subprotocols.
- Wiring authentication service into process startup.
- Adding session persistence.
- Adding runtime session tables or migrations.
- Changing authentication repository interfaces.
- Changing PostgreSQL adapters.
- Implementing logout, refresh, cleanup, token rotation, or token validation audit mutation.
- Adding dependencies.
- Declaring direct Nakama or Pitaya API compatibility.

## 10. Reference Mapping

Nakama guides the capability need for token validation before gameplay requests. vibit adapts that through request-level validation and explicit route policy rather than direct API compatibility.

Pitaya guides the separation between frontend connection handling and backend route handler context. vibit adapts that by keeping WebSocket transport credential-neutral and passing validated identity to application dispatch before domain handlers run.

## 11. Verification

The repository check rule for this gate is:

```text
runtime.access_token_protocol_carrier_route_protection_gate
```
