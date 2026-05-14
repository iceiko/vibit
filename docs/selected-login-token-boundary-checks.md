# Selected Login And Token Boundary Checks

Status: Draft v0.2
Last updated: 2026-05-14
Scope: Repository checks for the selected first login method, opaque token posture, schema gates, dependency deferral, generated output deferral, protocol deferral, and runtime implementation deferral
Depends on: `docs/login-method-token-format-ratification.md`, `docs/token-lifecycle-storage-implications.md`, `docs/authentication-contract-error-permission-surfaces.md`, `docs/credential-token-session-schema-gates.md`
Canonical decision: `ADR-0030`

The paired Simplified Chinese translation is `docs/selected-login-token-boundary-checks.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This standard defines the narrow repository checks that protect vibit's first selected authentication posture before implementation begins:

```yaml
login_method: device_credential_login
access_token_format: opaque_high_entropy_token
token_issuance_carrier: login_command_response_token
request_proof_carrier: explicit_request_proof_payload
refresh_token: not_in_first_implementation
renewal_method: reauthenticate_with_device_credential_login
default_durable_target: PostgreSQL
implementation_authorized: persistence_adapter_checks_refined
```

The checks exist because the posture is now selected, but runtime behavior is still not authorized. Selection is not implementation.

This standard does not add runtime authentication, login handlers, token generation, token parsing, token validation, logout execution, refresh behavior, credential lookup behavior, session persistence, external identity linking, Protobuf messages, generated Go contract shapes, WebSocket route behavior, WebSocket handshake authentication, audit persistence, provider dependencies, signing dependencies, password-hashing dependencies, key-management dependencies, or Redis-like dependencies.

Authentication PostgreSQL adapter persistence files are the only selected exception to the broad credential/token vocabulary ban, and only inside the declared platform boundary:

```text
runtime/internal/platform/persistence/postgres/authentication_repository.go
runtime/internal/platform/persistence/postgres/authentication_repository_test.go
```

That exception permits persistence vocabulary for the ratified `authentication_device_credentials` and `authentication_access_tokens` tables. It does not permit token generation, token validation, verifier comparison, bearer parsing, transport behavior, Protobuf behavior, transaction ownership, generated authentication shapes, or major authentication dependencies.

## 2. Rule

The repository check rule is:

```text
runtime.selected_login_token_boundary
```

The default command is:

```bash
node tools/vibit check runtime --json
```

The rule must also run through:

```bash
node tools/vibit check all --json
```

The check is static and local. It must not require live PostgreSQL, Docker, Podman, cloud services, OAuth providers, OIDC providers, platform identity providers, Redis-like services, or network access.

## 3. What The Check Protects

The check protects these boundaries:

- The selected login method remains semantic-contract-only.
- Opaque access-token behavior remains semantic-contract-only.
- Refresh tokens remain unsupported in the first implementation posture.
- Credential and token verifier schema gates remain defined but not implemented.
- Runtime session persistence remains deferred.
- External identity linking remains deferred.
- Authentication generated output remains deferred.
- Authentication Protobuf source remains deferred.
- WebSocket transport remains credential-neutral.
- Current Protobuf `Session` metadata remains metadata-only and not proof.
- Player account lifecycle storage remains credential-free, token-free, external-identity-free, session-free, and WebSocket-state-free.
- New authentication dependencies remain blocked until adoption records and implementation gates authorize them.

## 4. Required Static Signals

The check requires the following static signals to remain present:

```yaml
selected_posture:
  login_method: device_credential_login
  token_format: opaque_high_entropy_token
  token_issuance_carrier: login_command_response_token
  request_proof_carrier: explicit_request_proof_payload
  access_token_ttl: 1h
  refresh_token: not_in_first_implementation
  renewal_method: reauthenticate_with_device_credential_login
  logout_scope_first_posture: presented_access_token
schema_gates:
  credential_record_schema_gate_status: migration_source_added
  credential_record_schema_boundary: docs/credential-record-schema-boundary.md
  token_verifier_record_schema_gate_status: migration_source_added
  token_verifier_record_schema_boundary: docs/token-verifier-record-schema-boundary.md
  credential_migration_source: runtime/migrations/postgres/000003_create_authentication_device_credentials.sql
  credential_migration_source_added: true
  token_verifier_migration_source: runtime/migrations/postgres/000004_create_authentication_access_tokens.sql
  token_verifier_migration_source_added: true
  external_identity_link_schema_gate_status: deferred_no_schema_added
  runtime_session_record_schema_gate_status: deferred_no_schema_added
implementation_status:
  runtime_authentication_implemented: false
  token_behavior_implemented: false
  credential_storage_implemented: false
  token_storage_schema_added: true
  credential_storage_schema_added: true
  external_identity_storage_schema_added: false
  session_storage_schema_added: false
  migration_sources_added: credential_and_token_verifier
  repository_interfaces_added: true
  repository_interface_source: runtime/internal/modules/authentication/repository.go
  postgres_adapters_added: true
  authentication_postgres_adapter_added: true
  authentication_postgresql_adapter_checks_refined: true
  runtime_lookup_added: false
  websocket_handshake_authentication_changed: false
  runtime_player_handlers_added: false
  websocket_routes_added: false
  protobuf_envelope_changed: false
  generated_contract_shapes_added: true
  authentication_generated_contract_shapes_added: true
```

These signals are not a substitute for schemas or tests. They are tripwires that force later implementation work to explicitly change the architecture state.

`credential_storage_schema_added: true` and `token_storage_schema_added: true` mean only that the SQL migration sources exist. `repository_interfaces_added: true` means only that a storage-neutral interface boundary exists. `postgres_adapters_added: true` and `authentication_postgres_adapter_added: true` mean only that the bounded PostgreSQL persistence adapter exists under the platform package. `authentication_postgresql_adapter_checks_refined: true` means only that checks distinguish persistence-adapter vocabulary from runtime authentication behavior. `authentication_generated_contract_shapes_added: true` means only that source-traced metadata-only Go contract shapes exist under `runtime/internal/generated/contracts/runtime/authentication/`. These signals do not authorize runtime credential lookup behavior, login, token behavior, Protobuf messages, WebSocket behavior, or authentication dependencies.

## 5. Forbidden Shortcuts

Until a later bounded work item authorizes implementation, agents must not add:

- `AuthenticateWithDeviceCredential` runtime behavior.
- `ValidateAccessToken` runtime behavior.
- `LogoutAccessToken` runtime behavior.
- `RefreshAccessToken` runtime behavior.
- `AuthService`, `Authenticator`, `TokenValidator`, `TokenIssuer`, or `TokenVerifier` runtime implementation. The approved storage-neutral `Repository` interface under `runtime/internal/modules/authentication/repository.go` and the separately gated PostgreSQL adapter under `runtime/internal/platform/persistence/postgres/authentication_repository.go` are the only current exceptions.
- Token or credential random generation code.
- Bearer-token parsing or acceptance in runtime code.
- WebSocket `Authorization`, `Bearer`, `Cookie`, `Sec-WebSocket-Protocol`, or handshake header authentication behavior.
- Runtime authentication Protobuf source under `proto/vibit/runtime/`.
- Generated Go authentication contract shapes under `runtime/internal/generated/contracts/runtime/authentication/`.
- Additional credential, token, refresh-token, runtime-session, external-identity, provider-subject, or authentication-audit migrations beyond the ratified credential and token verifier sources.
- Player account lifecycle table changes that store credential, token, external identity, session, request validation, or WebSocket state.
- JWT, OAuth, OIDC, password-hashing, provider SDK, key-management, or Redis-like dependencies for authentication.

## 6. Allowed Current Artifacts

The following artifacts are allowed because they are semantic or gate artifacts, not runtime implementation:

- `contracts/runtime/authentication/**`
- `docs/authentication-contract-error-permission-surfaces.md`
- `docs/credential-token-session-schema-gates.md`
- `docs/credential-record-schema-boundary.md`
- `docs/token-verifier-record-schema-boundary.md`
- `docs/token-lifecycle-storage-implications.md`
- `docs/first-login-method-set.md`
- `docs/first-token-format-proof-carrier-posture.md`
- `runtime/internal/modules/authentication/repository.go`
- `runtime/internal/platform/persistence/postgres/authentication_repository.go` only after `W-0084` implements the declared persistence adapter.
- `runtime/internal/platform/persistence/postgres/authentication_repository_test.go` only for persistence-adapter tests.
- Architecture manifest markers that state implementation is deferred.
- Agent-facing guides that explain the boundary.

Generated Go contract shapes are allowed only as metadata-only source-traced output under `runtime/internal/generated/contracts/runtime/authentication/`. Protobuf wire messages remain deferred. Runtime handlers remain deferred. Runtime token validation remains deferred.

## 7. Machine-Readable Output

Every emitted check item must include:

```yaml
rule_id: runtime.selected_login_token_boundary
artifact: <repo-relative-forward-slash-path>
```

Repository-relative paths in JSON output must use forward slashes on all platforms, including Windows.

Agents should use:

```bash
node tools/vibit inspect rule runtime.selected_login_token_boundary --json
```

when they need rule metadata, and:

```bash
node tools/vibit check runtime --json
```

when they need actionable results.

## 8. Relationship To Nakama And Pitaya

Nakama and Pitaya remain active references for game backend capability and Go game-server vocabulary.

This check does not copy their authentication API shapes. It protects vibit's agent-native path from prematurely importing reference-framework behavior into transport, Protobuf, player persistence, or domain modules.

Reference alignment:

```yaml
nakama:
  device_style_login: adapted_as_capability_reference
  session_token_refresh_logout_vocabulary: adapted_as_lifecycle_reference
  direct_api_compatibility: rejected_for_now
pitaya:
  handler_session_context: adapted_as_request_identity_vocabulary
  session_binding: deferred
  transport_owned_authentication: rejected
```

## 9. Verification Path

For changes that touch this boundary, run:

```bash
node tools/vibit check runtime --json
node tools/vibit check architecture --json
node tools/vibit check agent-tooling --json
node tools/vibit check memory --json
node tools/vibit check work --json
node tools/vibit check all --json
```

If a change spec exists, also run:

```bash
node tools/vibit check change <change-id> --json
```

Live PostgreSQL verification is not required by this check because W-0072 does not add schema, migrations, repositories, adapters, or runtime behavior.

## 10. Future Migration Path

A future implementation milestone may intentionally change this boundary. That work must:

- Create or update a change spec.
- Update the English standard and Simplified Chinese translation.
- Update the related ADR or create a new ADR.
- Update `.arch/runtime.yaml`, `.arch/conventions.yaml`, `.arch/contracts.yaml`, `.arch/reference.yaml`, and relevant module manifests.
- Ratify schema before migrations.
- Ratify repository interfaces before adapters.
- Ratify generated output before generating shapes.
- Ratify Protobuf impact before adding wire messages.
- Ratify WebSocket carrier behavior before changing handshake or route behavior.
- Add focused tests.
- Update `runtime.selected_login_token_boundary` so it blocks only unapproved shortcuts and permits the newly approved slice.
