# Device Credential Login Service Behavior Gate

Status: Draft v0.1
Last updated: 2026-05-16
Scope: Future device credential login service behavior sequence, dependency shape, repository handoff, token issuance posture, public failure collapse, redaction, tests, and deferrals before real login execution is added
Depends on: `docs/authentication-service-behavior-implementation-gate.md`, `docs/application-authentication-service-interface-boundary.md`, `docs/token-credential-material-generation-implementation-gate.md`, `docs/verifier-digest-helper-implementation-gate.md`, `docs/verifier-digest-comparison-helper-gate.md`
Canonical decision: `ADR-0051`

The paired Simplified Chinese translation is `docs/device-credential-login-service-behavior-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

The authentication service skeleton exists and fails closed. The next risk is implementing real device credential login by simply filling the skeleton method with ad hoc repository calls, token generation, public error mapping, or transport assumptions.

This gate defines the exact future login behavior before code is allowed to execute it.

This is a gate-only standard. It does not implement device credential login, issue access tokens, validate access tokens, change service method signatures, expose protocol carriers, change repositories, change migrations, wire startup, add dependencies, or add production authentication behavior.

## 2. Core Rule

The device credential login service behavior gate is:

```yaml
device_credential_login_service_behavior_gate: defined
implementation_authorized_by_this_standard: false
future_implementation_work_item: W-0109
completed_gate_work_item: W-0108
planned_owner: runtime/internal/app
planned_package: runtime/internal/app/authentication
planned_source: runtime/internal/app/authentication/service.go
planned_tests: runtime/internal/app/authentication/service_test.go
service_method: AuthenticateWithDeviceCredential
login_method: device_credential_login
credential_kind: device_credential_login
token_kind: access_token
token_type: opaque_access_token
proof_carrier_status: already_decoded_service_request_only
protocol_carrier_status: deferred
startup_wiring_status: deferred
repository_interface_change_status: deferred
migration_schema_change_status: deferred
authentication_dependencies_status: deferred
production_authentication_behavior_status: deferred
```

Future implementation may only execute the login flow after a later work item explicitly authorizes it.

## 3. Future Dependency Shape

The future implementation may extend the existing `ServiceDependencies` shape only inside the application service boundary.

Required future dependency categories:

```yaml
future_service_dependencies:
  unit_of_work_runner: already_present
  verifier_key_set: required
  access_token_entropy_reader: required
  clock: required
  token_record_id_generator: required
  access_token_lifetime: required_positive_duration
  token_audience: required_non_empty_string
```

Rules:

- `UnitOfWorkRunner` remains the only transaction entry point.
- Future login behavior must use a local capability interface on the provided `tx.UnitOfWork` to obtain `NewAuthenticationRepository()` and `NewPlayerAccountRepository()`.
- The global `tx.UnitOfWork` interface must not be expanded only for this login slice.
- Token record id generation must be injected. This gate does not choose UUID, ULID, KSUID, database-generated ids, or an external id package.
- Access-token lifetime must be configured through service dependencies and must be positive.
- No production default token lifetime is ratified by this gate.
- No startup configuration is wired by this gate.

## 4. Proof Input Shape

The future login method receives proof through `DeviceCredentialAuthenticationRequest.CredentialProof`.

First proof rules:

```yaml
device_credential_proof:
  source: already_decoded_service_request
  text_encoding: base64url_unpadded
  encoded_length_chars: 43
  raw_length_bytes: 32
  raw_entropy_floor_bits: 256
  bearer_prefix: forbidden
  raw_device_identifier: forbidden
  client_generated_low_entropy_value: forbidden
```

Rules:

- The proof is server-issued high-entropy device credential material.
- It is not a raw operating-system device id, advertising id, hardware serial number, user name, email address, provider subject, session id, or transport metadata.
- Missing, whitespace-only, padded, wrongly sized, non-Base64URL, or non-32-byte decoded proof must fail before any unit-of-work or repository call.
- Future implementation should decode the text into raw material only long enough to compute digests, then keep it out of logs, public errors, test snapshots, conversation logs, and docs examples.
- The application service must not parse an HTTP `Authorization` header or `Bearer` string for this login method. Protocol carriers are deferred.

## 5. Required Login Sequence

When W-0109 or a later work item authorizes behavior, `AuthenticateWithDeviceCredential` must execute this sequence:

```yaml
device_credential_login_sequence:
  - reject_missing_or_malformed_device_credential_proof_before_unit_of_work
  - decode_device_credential_proof_text_to_raw_32_byte_material
  - compute_credential_lookup_digest_with_active_VerifierKeySet
  - enter_application_unit_of_work
  - obtain_authentication_repository_from_unit_of_work_capability
  - obtain_player_account_repository_from_unit_of_work_capability
  - find_credential_by_lookup_digest
  - collapse_lookup_miss_to_public_invalid_credential
  - require_credential_kind_device_credential_login
  - require_credential_status_active
  - require_supported_verifier_algorithm_vibit_hmac_sha256_v1
  - require_supported_verifier_version_1
  - require_active_verifier_key_set_id_match_for_first_posture
  - compute_credential_verifier_digest
  - compare_credential_verifier_digest_with_CompareCredentialVerifierDigest
  - collapse_verifier_mismatch_to_public_invalid_credential
  - get_player_account_by_credential_player_id
  - require_player_account_active
  - generate_access_token_material_with_explicit_entropy_reader
  - compute_token_lookup_digest_with_active_VerifierKeySet
  - compute_token_verifier_digest_with_active_VerifierKeySet
  - create_token_record_id_with_injected_generator
  - store_access_token_record_through_authentication_repository
  - exit_unit_of_work_successfully
  - return_raw_access_token_text_once_after_unit_of_work_success
```

Rules:

- No repository call may happen before missing or malformed proof rejection.
- No verifier digest comparison may happen before credential lookup succeeds.
- No access-token material may be generated before credential proof and player account state are accepted.
- No raw access token may be returned if token storage, transaction commit, or any dependency fails.
- The raw access token text may be returned only once through `AuthenticationResult.AccessToken` after the unit of work succeeds.
- Lookup digest equality alone is not proof.
- Database equality alone is not final proof.
- Player account activity is required for the first login behavior.
- Credential disabled, revoked, replaced, wrong algorithm, wrong version, wrong key id, player missing, player disabled, player deleted, lookup miss, and verifier mismatch collapse to the same public invalid-credential family unless a later disclosure decision changes this.

## 6. Repository Handoff

Future login behavior must use existing repository interfaces. This gate does not authorize repository interface changes.

Required future repository access:

```yaml
repository_handoff:
  transaction_boundary: UnitOfWorkRunner.WithinUnitOfWork
  authentication_repository_source: unit_of_work.NewAuthenticationRepository
  player_repository_source: unit_of_work.NewPlayerAccountRepository
  authentication_lookup_method: FindCredentialByLookupDigest
  token_store_method: StoreToken
  player_lookup_method: GetPlayerAccount
  direct_postgres_import: forbidden
  repository_interface_change: forbidden_by_this_gate
```

The service must not import PostgreSQL adapters, pgx, goose, WebSocket transport packages, Protobuf adapter packages, generated Protobuf packages, or generated contract-shape packages for login behavior.

The authentication module remains storage-neutral. It stores and reads already-computed digest records only. It must not generate access tokens, compute digests, compare verifiers, decide proof validity, collapse public failures, or construct application responses.

## 7. Token Issuance Posture

First future login behavior may issue only opaque access tokens.

```yaml
token_issuance:
  token_kind: access_token
  token_type: opaque_access_token
  actor_kind: player
  token_material_generation: GenerateAccessTokenMaterial
  token_lookup_digest_helper: ComputeTokenLookupDigest
  token_verifier_digest_helper: ComputeTokenVerifierDigest
  token_record_store_method: StoreToken
  refresh_token_issued: false
  jwt_or_signed_claim_token: forbidden
  previous_token_rotation: deferred
  cleanup_job: deferred
```

Rules:

- `StoreTokenMutation` must receive digest bytes only, never raw token material.
- `VerifierAlgorithm` must be `vibit_hmac_sha256_v1`.
- `VerifierVersion` must be `1` for the first implementation.
- `VerifierKeyID` must come from the active `VerifierKeySet.KeySetID()`.
- `IssuedAt` and `ExpiresAt` must come from the injected clock and positive access-token lifetime.
- `TokenRecordID` must come from the injected token record id generator.
- `Audience` must come from service configuration.
- Refresh-token behavior remains unsupported.
- Token rotation or revocation of previous tokens for the same credential remains deferred unless a later work item authorizes the exact behavior.

## 8. Public Error Collapse

Future behavior must preserve redacted internal classes while collapsing public proof failures.

```yaml
public_error_collapse:
  missing_proof: AUTHENTICATION_PROOF_MISSING
  malformed_proof: AUTHENTICATION_PROOF_MALFORMED
  lookup_miss: AUTHENTICATION_CREDENTIAL_INVALID
  wrong_credential_kind: AUTHENTICATION_CREDENTIAL_INVALID
  inactive_credential: AUTHENTICATION_CREDENTIAL_INVALID
  unsupported_algorithm: AUTHENTICATION_CREDENTIAL_INVALID
  unsupported_version: AUTHENTICATION_CREDENTIAL_INVALID
  verifier_key_id_mismatch: AUTHENTICATION_CREDENTIAL_INVALID
  verifier_mismatch: AUTHENTICATION_CREDENTIAL_INVALID
  player_missing_or_inactive: AUTHENTICATION_CREDENTIAL_INVALID
  repository_or_unit_of_work_unavailable: AUTHENTICATION_CREDENTIAL_STORE_UNAVAILABLE
  token_generation_or_storage_unavailable: AUTHENTICATION_CREDENTIAL_STORE_UNAVAILABLE
```

The public surface must not reveal whether a credential record exists, whether the player exists, which verifier key id is stored, whether a token store failed after proof acceptance, or whether a verifier mismatch happened.

## 9. Redaction Requirements

Future behavior must never place these values in errors, logs, docs examples, test failure messages, conversation logs, ADRs, change specs, or public responses:

- Raw device credential text.
- Raw device credential bytes.
- Raw access-token text except the one-time successful `AuthenticationResult.AccessToken` carrier.
- Raw access-token bytes.
- Lookup digest bytes.
- Verifier digest bytes.
- HMAC input or output bytes.
- Verifier key bytes.
- Full concrete `verifier_key_id` values.
- Credential lookup hit or miss details.
- Player lookup hit or miss details.
- Token store internals for a failed proof.

Allowed:

- Registered public error codes.
- Redacted placeholders such as `<device-credential-proof>` and `<access-token>`.
- Non-secret record ids only after proof succeeds and only in authorized result fields or internal tests.

## 10. Required Tests

Future implementation must add or update focused tests in `runtime/internal/app/authentication/service_test.go`.

Required test classes:

```yaml
required_tests:
  login_rejects_missing_proof_without_unit_of_work
  login_rejects_malformed_proof_without_unit_of_work
  login_computes_lookup_digest_before_repository_lookup
  login_uses_authentication_repository_from_unit_of_work_only
  login_uses_player_repository_from_unit_of_work_only
  login_collapses_lookup_miss_to_invalid_credential
  login_rejects_inactive_or_wrong_kind_credential
  login_rejects_wrong_algorithm_version_or_key_id
  login_compares_verifier_digest_before_token_generation
  login_collapses_verifier_mismatch_to_invalid_credential
  login_requires_active_player_account
  login_generates_access_token_only_after_proof_and_player_acceptance
  login_stores_token_digest_only
  login_does_not_return_access_token_when_store_or_commit_fails
  login_returns_access_token_once_after_unit_of_work_success
  login_errors_do_not_leak_raw_proof_or_token_material
  login_does_not_validate_access_tokens_or_touch_protocol_carriers
```

Normal tests must use fake unit-of-work and repository implementations. They must not require a live PostgreSQL server.

## 11. Nakama And Pitaya Mapping

Nakama capability mapping:

- Server-side account authentication is adapted as application-owned proof validation.
- Session/access-token issuance is adapted as opaque access-token generation and digest storage.
- Token expiration and revocation remain token validation behavior, not login behavior.

Pitaya capability mapping:

- Realtime handlers should receive identity context after validation.
- Frontend and backend handler separation maps to vibit's transport/protocol/application separation.
- Login proof validation remains outside route/domain handlers.

These references guide capability coverage only. They do not authorize copying Nakama or Pitaya public APIs.

## 12. Deferrals

This gate does not authorize:

- Device credential login execution.
- Access-token validation execution.
- Logout execution.
- Refresh execution.
- Token cleanup jobs.
- Protocol authentication messages.
- WebSocket proof carriers.
- WebSocket handshake authentication.
- Startup wiring.
- Repository interface changes.
- PostgreSQL adapter changes.
- SQL migration changes.
- Generated file changes.
- External cryptography, JWT, OAuth, OIDC, provider, KMS, cloud secret-manager, password-hashing, Redis-like, queue, or session-store dependencies.
- Production authentication behavior.

## 13. Verification

The repository check rule for this gate is:

```text
runtime.device_credential_login_service_behavior_gate
```
