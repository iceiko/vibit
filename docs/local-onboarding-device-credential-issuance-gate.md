# Local Onboarding Device Credential Issuance Gate

Status: Draft v0.1
Last updated: 2026-05-21
Scope: Gate-only boundary for the first local developer onboarding flow that creates a player account and issues a server-generated device credential
Depends on: `docs/v0.1-alpha-goal.md`, `docs/device-credential-login-service-behavior-gate.md`, `docs/token-credential-material-generation-implementation-gate.md`, `docs/verifier-digest-helper-implementation-gate.md`, `docs/verifier-digest-comparison-helper-gate.md`, `docs/credential-record-schema-boundary.md`, `docs/postgresql-persistence-boundary.md`
Canonical decision: `ADR-0089`

The paired Simplified Chinese translation is `docs/local-onboarding-device-credential-issuance-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

The alpha path now has login, access-token validation, runtime session metadata, first-message connection binding, protected inventory routes, a presence lifecycle primitive, a protected presence query, and logout. A new local developer still lacks the first credential needed to enter that path.

This gate defines the future local onboarding/device credential issuance boundary before any implementation is added.

This is a gate-only standard. It does not implement onboarding, generate or display credentials, create player accounts through a new flow, write credential records through a new flow, expose a public protocol route, change Protobuf sources, change generated output, change migrations, add dependencies, publish a release, add production signup, add external identity providers, add password login, add account recovery, add multi-device linking, or adopt direct Nakama/Pitaya API compatibility.

## 2. Core Rule

The local onboarding device credential issuance gate is:

```yaml
local_onboarding_device_credential_issuance_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0181
future_implementation_work_item: W-0182
decision: ADR-0089
check_rule: runtime.local_onboarding_device_credential_issuance_gate
first_surface_candidate: local_developer_onboarding_application_service
surface_visibility: local_only_not_public_signup
future_service_owner: runtime/internal/app/authentication
future_source: runtime/internal/app/authentication/service.go
future_tests: runtime/internal/app/authentication/service_test.go
future_service_method_candidate: OnboardLocalPlayerWithDeviceCredential
player_account_repository_method: CreatePlayerAccount
authentication_repository_method: StoreCredential
credential_material_helper: GenerateDeviceCredentialMaterial
credential_lookup_digest_helper: ComputeCredentialLookupDigest
credential_verifier_digest_helper: ComputeCredentialVerifierDigest
login_route_account_creation_behavior_changed: false
public_protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
production_signup_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

Future implementation may only add a local developer onboarding application service after a later work item explicitly authorizes code. It must not turn the existing public `runtime.authentication.AuthenticateWithDeviceCredential` login route into account creation behavior.

## 3. Future Surface

The first future surface candidate is application-local:

```yaml
local_onboarding_surface:
  owner: runtime/internal/app/authentication
  visibility: local_developer_only
  public_protocol_route: forbidden_by_this_gate
  websocket_route: forbidden_by_this_gate
  http_route: forbidden_by_this_gate
  startup_auto_onboarding: forbidden_by_this_gate
  production_signup: forbidden_by_this_gate
```

Rules:

- The future method may exist as an application service method that tests and later local tooling can call.
- It is not a public game client signup route.
- It must not be reachable through WebSocket Protobuf routing unless a later protocol gate explicitly adds a route.
- It must not create player accounts automatically during process startup.
- It must not interpret `AccountCreationIntent` on the existing login route as authorization to create accounts.
- It must not accept client-generated device credential material in the first posture.

The local-only posture is intentional. It gives alpha developers a controlled way to obtain a first credential without committing vibit to production signup, abuse controls, identity provider linking, recovery, account merge, or multi-device behavior.

## 4. Future Dependency Shape

The future implementation should extend application service dependencies only as needed for local onboarding.

Required future dependency categories:

```yaml
future_service_dependencies:
  unit_of_work_runner: already_present
  verifier_key_set: already_present
  device_credential_entropy_reader: required
  clock: already_present
  player_id_generator: required
  player_account_event_id_generator: required
  credential_record_id_generator: required
```

Rules:

- The application unit-of-work remains the only transaction entry point.
- The service must obtain `NewPlayerAccountRepository()` and `NewAuthenticationRepository()` from the unit-of-work capability.
- The global `tx.UnitOfWork` interface must not be expanded only for this slice.
- Identifier generation must be injected. This gate does not select UUID, ULID, KSUID, database-generated ids, or an external id dependency.
- The device credential entropy reader must be explicit so tests do not depend on nondeterministic process state.
- Production defaults for id generation, display names, local operator identity, or credential lifetime are not ratified by this gate.

## 5. Future Request And Result Shape

Candidate future request:

```yaml
LocalOnboardingDeviceCredentialIssuanceRequest:
  display_name: required_non_secret_text
  requested_by: optional_local_operator_label
```

Candidate future result:

```yaml
LocalOnboardingDeviceCredentialIssuanceResult:
  status: created
  player_id: generated_player_id
  credential_record_id: generated_credential_record_id
  device_credential: one_time_raw_credential_text
  created_at: server_time
```

Rules:

- The first posture should generate `player_id`, player account event id, and credential record id on the server side.
- The first posture should not allow caller-supplied `player_id` as identity proof.
- Display name is not proof and must not be embedded into credential material or digests.
- The raw device credential text may appear only in the successful result after the unit of work commits.
- The raw device credential bytes and text must never be stored.
- The result must not include an access token. Login remains a separate step through the existing device credential login route.

## 6. Required Future Sequence

When `W-0182` or a later work item authorizes behavior, the future onboarding method must execute this sequence:

```yaml
local_onboarding_sequence:
  - reject_invalid_local_request_before_unit_of_work
  - generate_device_credential_material_with_explicit_entropy_reader
  - compute_credential_lookup_digest_with_active_VerifierKeySet
  - compute_credential_verifier_digest_with_active_VerifierKeySet
  - enter_application_unit_of_work
  - obtain_player_account_repository_from_unit_of_work_capability
  - obtain_authentication_repository_from_unit_of_work_capability
  - generate_player_account_identifiers_with_injected_generators
  - create_active_player_account
  - create_active_device_credential_record_with_digest_only_storage
  - exit_unit_of_work_successfully
  - return_raw_device_credential_text_once_after_commit
```

Rules:

- No repository call may happen before basic request validation.
- Credential material generation must use `GenerateDeviceCredentialMaterial`.
- Credential digest computation must use `ComputeCredentialLookupDigest` and `ComputeCredentialVerifierDigest`.
- The credential record must use `credential_kind=device_credential_login`.
- The credential record must use `verifier_algorithm=vibit_hmac_sha256_v1` and `verifier_version=1`.
- `verifier_key_id` must come from the active `VerifierKeySet.KeySetID()`.
- Player account creation and credential record storage must commit or roll back together.
- If player account creation, credential record storage, dependency lookup, id generation, digest computation, or unit-of-work commit fails, the method must not return raw credential text as success.

## 7. Repository Handoff

Future implementation must use existing repository interfaces.

```yaml
repository_handoff:
  transaction_boundary: UnitOfWorkRunner.WithinUnitOfWork
  player_repository_source: unit_of_work.NewPlayerAccountRepository
  authentication_repository_source: unit_of_work.NewAuthenticationRepository
  player_mutation_method: CreatePlayerAccount
  credential_store_method: StoreCredential
  direct_postgres_import: forbidden
  repository_interface_change: forbidden_by_this_gate
```

The authentication repository receives only digest and metadata fields:

```yaml
credential_store_allowed_fields:
  - credential_record_id
  - player_id
  - credential_kind
  - credential_lookup_digest
  - credential_verifier_digest
  - verifier_algorithm
  - verifier_version
  - verifier_key_id
  - occurred_at
  - requested_by
```

Forbidden repository inputs:

- Raw device credential text.
- Raw device credential bytes.
- Encoded credential material.
- Verifier key bytes.
- Encoded verifier key values.
- Access-token material.
- Provider subjects, passwords, OAuth claims, OIDC claims, or account recovery data.

## 8. Relationship To Existing Login

This gate does not change device credential login.

```yaml
existing_login_route: runtime.authentication.AuthenticateWithDeviceCredential
existing_login_route_account_creation_behavior_changed: false
existing_login_route_proof_required: true
login_route_returns_access_token: already_implemented
local_onboarding_returns_access_token: false
```

Rules:

- `AuthenticateWithDeviceCredential` continues to require a presented device credential proof.
- `AccountCreationIntent` on the existing login request remains non-creating until a later route behavior decision explicitly changes it.
- Local onboarding creates the credential. Login authenticates with that credential.
- Local onboarding must not bypass future login proof validation by issuing an access token directly.

## 9. Redaction Requirements

Future implementation must never place these values in errors, logs, traces, metrics labels, docs examples, test snapshots, ADRs, change specs, conversation logs, or public responses except the one-time successful local onboarding result carrier:

- Raw device credential text.
- Raw device credential bytes.
- Encoded generated credential material.
- Credential lookup digest bytes.
- Credential verifier digest bytes.
- HMAC input or output bytes.
- Verifier key bytes.
- Encoded verifier key values.
- Full concrete `verifier_key_id` values.
- Database connection strings or credentials.
- Credential lookup hit or miss details.
- Player account conflict details that disclose private state.

Allowed:

- Non-secret generated `player_id` and `credential_record_id` in the successful result.
- Registered or local-only redacted error classes.
- Redacted placeholders such as `<device-credential>`, `<credential-lookup-digest>`, and `<verifier-key-id>`.

Generated credential material is not safe just because it is local. One-time presentation means one successful local result carrier, not "safe to log once."

## 10. Future Tests

The future implementation must add focused tests under:

```text
runtime/internal/app/authentication/service_test.go
```

Minimum test classes:

```yaml
required_tests:
  onboarding_rejects_invalid_request_without_unit_of_work
  onboarding_generates_device_credential_with_explicit_reader
  onboarding_computes_credential_lookup_and_verifier_digests_before_storage
  onboarding_uses_player_repository_from_unit_of_work_only
  onboarding_uses_authentication_repository_from_unit_of_work_only
  onboarding_creates_active_player_account_before_credential_record
  onboarding_stores_credential_digests_only
  onboarding_returns_raw_device_credential_only_after_commit
  onboarding_does_not_return_credential_when_player_creation_fails
  onboarding_does_not_return_credential_when_credential_storage_fails
  onboarding_does_not_return_credential_when_commit_fails
  onboarding_does_not_issue_access_token_or_runtime_session
  onboarding_errors_do_not_leak_raw_credential_digest_or_key_material
  existing_login_route_still_does_not_create_accounts
```

The implementation should not require live PostgreSQL by default. Repository behavior can be covered with fakes or existing adapter tests unless a later work item explicitly requests live verification.

## 11. Deferrals

This gate preserves these deferrals:

- Runtime onboarding implementation.
- Public signup or production account creation.
- Public WebSocket, HTTP, or CLI surface selection.
- Protobuf request/response messages.
- Generated Go output.
- Migration schema changes.
- Repository interface changes.
- New dependencies.
- External identity providers, OAuth, OIDC, provider SDKs, password login, password hashing, account recovery, account merge, and multi-device linking.
- Credential rotation or replacement behavior.
- Access-token issuance from onboarding.
- Runtime session creation from onboarding.
- Memory-store durable authentication behavior.
- Release publishing.
- Direct Nakama/Pitaya public API compatibility.

## 12. Nakama And Pitaya Mapping

Nakama shows that usable game backends need a path to create or authenticate a player and obtain secret session/token material before gameplay requests. vibit adapts the capability need but keeps first local credential issuance separate from production signup and direct Nakama API compatibility.

Pitaya reinforces that transport acceptors, route handlers, and identity context should stay separated. vibit adapts that by keeping local onboarding in application service orchestration, not WebSocket transport or Protobuf envelope metadata.

## 13. Verification

The repository check rule for this gate is:

```text
runtime.local_onboarding_device_credential_issuance_gate
```

Recommended verification for this gate-only change:

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check change define-local-onboarding-device-credential-issuance-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Runtime Go tests and Buf generation are not required for this gate-only change because it does not add Go runtime behavior, Protobuf sources, generated output, migrations, dependencies, or release artifacts.
