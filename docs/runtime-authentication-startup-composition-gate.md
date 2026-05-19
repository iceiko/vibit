# Runtime Authentication Startup Composition Gate

Status: Draft v0.1
Last updated: 2026-05-17
Scope: Runtime process startup composition for the existing application authentication service, verifier key environment loader, request-level route protection, and Protobuf frame handler
Depends on: `docs/environment-verifier-key-loader-gate.md`, `docs/access-token-validation-service-behavior-gate.md`, `docs/access-token-protocol-carrier-route-protection-gate.md`, `docs/runtime-protocol-adapter.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0054`

The paired Simplified Chinese translation is `docs/runtime-authentication-startup-composition-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

The runtime can now validate opaque access-token proof at the application service layer, and the Protobuf adapter can enforce request-level route protection when a `RouteProtector` is injected.

The remaining gap is process startup composition: `runtime/cmd/vibit-server` must wire the existing authentication service, verifier key configuration, PostgreSQL unit-of-work runner, route access-token validator, route protector, and Protobuf frame handler together.

This is a startup-composition gate. It does not add new authentication semantics. It does not add WebSocket handshake authentication, session persistence, authentication command Protobuf messages, login route registration, repository interface changes, PostgreSQL adapter changes, migrations, dependencies, logout, refresh, cleanup, token rotation, token validation audit mutation, or broader production authentication behavior.

## 2. Core Rule

The runtime authentication startup composition gate is:

```yaml
runtime_authentication_startup_composition_gate: defined
implementation_authorized_by_this_standard: true
completed_gate_work_item: W-0116
future_implementation_work_item: W-0117
decision: ADR-0054
startup_owner: runtime/cmd/vibit-server
application_authentication_owner: runtime/internal/app/authentication
route_policy_owner: runtime/internal/app
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
transport_owner: runtime/internal/platform/transport/ws
first_composed_runtime_store: postgres
memory_store_authentication_status: bootstrap_metadata_only
websocket_transport_credential_neutral: true
protobuf_envelope_change_status: unchanged
session_persistence_status: deferred
websocket_handshake_authentication_status: deferred
```

Future implementation may wire existing authentication behavior into startup only through composition. It must not move authentication logic into WebSocket transport, Protobuf envelope metadata, domain handlers, repositories, or generated files.

## 3. Selected Startup Path

The first startup composition path is limited to:

```yaml
runtime_store: postgres
selector: VIBIT_RUNTIME_STORE=postgres
startup_file: runtime/cmd/vibit-server/main.go
required_store_capabilities:
  - postgres.UnitOfWork.NewAuthenticationRepository
  - postgres.UnitOfWork.NewPlayerAccountRepository
  - postgres.NewPoolRunner
```

Rules:

- PostgreSQL startup must fail closed when required authentication verifier key configuration is missing or invalid.
- The default in-memory runtime store remains a bootstrap path and must not pretend to have durable authentication repository capability.
- Startup must not apply migrations automatically.
- Startup must not create credential or token records by itself.
- Startup must not register authentication command routes unless a later protocol/route work item authorizes them.

## 4. Composition Flow

Future implementation must follow this layer sequence:

```yaml
startup_composition_flow:
  - read_runtime_store_selection
  - open_postgres_pool_for_postgres_store
  - build_postgres_unit_of_work_runner
  - load_verifier_key_set_from_explicit_environment_lookup
  - create_authentication_service_with_existing_dependencies
  - create_route_access_token_validator_from_service
  - create_application_route_protector
  - inject_route_protector_into_protobuf_frame_handler
  - mount_websocket_transport_with_opaque_frame_handler
```

Rules:

- The WebSocket transport continues to receive and emit opaque frames only.
- The Protobuf frame handler continues to unwrap the already-ratified `vibit.authentication.v1.AuthenticatedRequest` payload wrapper.
- The application route protector remains the owner of protected-route authentication policy.
- Domain modules receive `RequestIdentity`; they do not receive raw access-token proof.
- `SessionValidated` remains false until session persistence is ratified.

## 5. Startup Dependencies

Future implementation may use only already-ratified helpers and standard library packages:

```yaml
authentication_service_dependencies:
  unit_of_work_runner: postgres.NewPoolRunner(pool)
  verifier_key_set: authentication.LoadVerifierKeySetFromEnvironment(lookup)
  access_token_random: crypto/rand.Reader
  clock: startup_owned_system_clock
  token_record_id_generator: startup_owned_standard_library_generator
  access_token_lifetime:
    default: 1h
    optional_environment: VIBIT_AUTH_ACCESS_TOKEN_TTL
  token_audience:
    default: vibit_gameplay_runtime_requests
    optional_environment: VIBIT_AUTH_TOKEN_AUDIENCE
```

Rules:

- `VIBIT_AUTH_ACCESS_TOKEN_TTL` must parse as a positive Go duration when present.
- `VIBIT_AUTH_TOKEN_AUDIENCE` must trim whitespace and fall back to `vibit_gameplay_runtime_requests` when absent or empty.
- Startup-generated token record ids must be non-secret identifiers.
- Raw verifier keys, access tokens, credential proof, lookup digests, verifier digests, and full concrete verifier key ids must not be logged or included in startup errors.
- No external UUID, ULID, KSUID, JWT, OAuth, OIDC, password-hashing, Redis-like, KMS, cloud secret-manager, or session-store dependency is authorized by this gate.

## 6. Nakama And Pitaya Reference Mapping

Nakama guides the capability need: clients authenticate to obtain a session/token before using server and realtime features. vibit adapts that by composing token validation into request-level route protection for the PostgreSQL runtime path.

Pitaya guides the architectural separation: acceptors/connections, sessions, routing, and handlers are distinct concerns. vibit adapts that by keeping WebSocket transport credential-neutral and injecting validated identity before domain handlers run.

Neither reference overrides vibit's boundaries. This gate does not adopt direct Nakama or Pitaya public API compatibility.

## 7. Required Implementation Tests

The implementation slice must add or update focused tests under:

```text
runtime/cmd/vibit-server/main_test.go
```

Required test classes:

```yaml
required_tests:
  memory_startup_remains_bootstrap_without_route_protector
  explicit_route_protector_can_be_injected_into_frame_handler
  postgres_auth_startup_requires_verifier_key_configuration
  auth_startup_accepts_default_lifetime_and_audience
  auth_startup_accepts_configured_lifetime_and_audience
  auth_startup_rejects_invalid_lifetime
  token_record_id_generator_returns_non_secret_stable_shape
  startup_errors_do_not_include_verifier_key_values
```

Live PostgreSQL startup verification remains optional and must not be required by default repository checks.

## 8. Deferrals

This gate does not authorize:

- WebSocket handshake authentication.
- Session persistence.
- Runtime session tables or migrations.
- Authentication command Protobuf messages beyond the existing payload wrapper.
- Login route registration.
- HTTP `Authorization`, Bearer, cookie, query string, or WebSocket subprotocol credential carriers.
- Repository interface changes.
- PostgreSQL adapter changes.
- Migration changes.
- Automatic startup migrations.
- Logout, refresh, cleanup, token rotation, token validation audit mutation, or previous-token revocation.
- External dependencies.
- Direct Nakama or Pitaya API compatibility.

## 9. Verification

The repository check rule for this gate is:

```text
runtime.authentication_startup_composition_gate
```
