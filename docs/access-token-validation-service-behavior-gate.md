# Access Token Validation Service Behavior Gate

Status: Draft v0.1
Last updated: 2026-05-16
Scope: Future access-token validation service behavior sequence, proof input shape, repository handoff, token lifecycle checks, request identity handoff, public failure collapse, redaction, tests, and deferrals before real validation execution is added
Depends on: `docs/authentication-service-behavior-implementation-gate.md`, `docs/application-authentication-service-interface-boundary.md`, `docs/verifier-digest-helper-implementation-gate.md`, `docs/verifier-digest-comparison-helper-gate.md`, `docs/device-credential-login-service-behavior-gate.md`
Canonical decision: `ADR-0052`

The paired Simplified Chinese translation is `docs/access-token-validation-service-behavior-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

Device credential login can now issue opaque access tokens through the application authentication service. The next risk is validating those tokens by adding ad hoc bearer parsing, repository calls, route protection, session persistence, or public error disclosure directly into runtime dispatch or transport code.

This gate defines the exact future access-token validation behavior before code is allowed to execute it.

This is a gate-only standard. It does not implement access-token validation, change service method signatures, expose protocol carriers, change repositories, change migrations, wire startup, add session persistence, add dependencies, or add production authentication behavior.

## 2. Core Rule

The access-token validation service behavior gate is:

```yaml
access_token_validation_service_behavior_gate: defined
implementation_authorized_by_this_standard: false
future_implementation_work_item: W-0111
completed_gate_work_item: W-0110
planned_owner: runtime/internal/app
planned_package: runtime/internal/app/authentication
planned_source: runtime/internal/app/authentication/service.go
planned_tests: runtime/internal/app/authentication/service_test.go
service_method: ValidateAccessToken
token_kind: access_token
token_type: opaque_access_token
proof_carrier_status: already_decoded_service_request_only
protocol_carrier_status: deferred
startup_wiring_status: deferred
repository_interface_change_status: deferred
migration_schema_change_status: deferred
session_persistence_status: deferred
authentication_dependencies_status: deferred
production_authentication_behavior_status: deferred
```

Future implementation may only execute token validation after a later work item explicitly authorizes it.

## 3. Future Dependency Shape

The future implementation should reuse the existing `ServiceDependencies` shape created for the login service slice.

Required future dependency categories:

```yaml
future_service_dependencies:
  unit_of_work_runner: already_present
  verifier_key_set: already_present
  clock: already_present
  token_audience: already_present
  access_token_entropy_reader: not_used_by_validation
  token_record_id_generator: not_used_by_validation
  access_token_lifetime: not_used_by_validation
```

Rules:

- `UnitOfWorkRunner` remains the only transaction entry point.
- Future validation behavior must use a local capability interface on the provided `tx.UnitOfWork` to obtain `NewAuthenticationRepository()` and `NewPlayerAccountRepository()`.
- The global `tx.UnitOfWork` interface must not be expanded only for this validation slice.
- Validation must not generate new token material or token record ids.
- The injected clock is used only to check token time windows.
- The configured token audience must match the stored token audience.
- No startup configuration is wired by this gate.

## 4. Proof Input Shape

The future validation method receives proof through `AccessTokenValidationRequest.AccessToken`.

First proof rules:

```yaml
access_token_proof:
  source: already_decoded_service_request
  text_encoding: base64url_unpadded
  encoded_length_chars: 43
  raw_length_bytes: 32
  raw_entropy_floor_bits: 256
  bearer_prefix: forbidden
  authorization_header_parsing: forbidden
  cookie_parsing: forbidden
  query_string_parsing: forbidden
  session_id_as_token: forbidden
```

Rules:

- The access token is opaque high-entropy material issued by the service login flow.
- It is not a JWT, signed claim token, session id, WebSocket connection id, device id, player id, route field, or transport metadata.
- Missing, whitespace-only, padded, wrongly sized, non-Base64URL, or non-32-byte decoded proof must fail before any unit-of-work or repository call.
- Future implementation should decode the text into raw material only long enough to compute digests, then keep it out of logs, public errors, test snapshots, conversation logs, and docs examples.
- The application service must not parse an HTTP `Authorization` header, `Bearer` string, cookie, query parameter, WebSocket handshake field, or Protobuf authentication carrier. Protocol carriers are deferred.

## 5. Required Validation Sequence

When W-0111 or a later work item authorizes behavior, `ValidateAccessToken` must execute this sequence:

```yaml
access_token_validation_sequence:
  - reject_missing_or_malformed_access_token_before_unit_of_work
  - decode_access_token_text_to_raw_32_byte_material
  - compute_token_lookup_digest_with_active_VerifierKeySet
  - enter_application_unit_of_work
  - obtain_authentication_repository_from_unit_of_work_capability
  - obtain_player_account_repository_from_unit_of_work_capability
  - find_token_by_lookup_digest
  - collapse_lookup_miss_to_public_invalid_token
  - require_token_kind_access_token
  - require_token_status_active
  - require_supported_verifier_algorithm_vibit_hmac_sha256_v1
  - require_supported_verifier_version_1
  - require_active_verifier_key_set_id_match_for_first_posture
  - require_configured_token_audience_match
  - require_token_issued_at_not_in_future_beyond_clock_tolerance
  - require_token_not_expired_by_injected_clock
  - compute_token_verifier_digest
  - compare_token_verifier_digest_with_CompareTokenVerifierDigest
  - collapse_verifier_mismatch_to_public_invalid_token
  - get_player_account_by_token_player_id
  - require_player_account_active
  - construct_validated_player_RequestIdentity
  - exit_unit_of_work_successfully
  - return_validated_identity_after_unit_of_work_success
```

Rules:

- No repository call may happen before missing or malformed proof rejection.
- No verifier digest comparison may happen before token lookup succeeds.
- No request identity may be marked validated before token verifier comparison, token lifecycle checks, audience check, and player account state checks succeed.
- Lookup digest equality alone is not proof.
- Database equality alone is not final proof.
- Player account activity is required for the first validation behavior.
- The first posture checks only the active `VerifierKeySet.KeySetID()`. Previous-key validation remains deferred.
- Token disabled, expired, revoked, replaced, wrong kind, wrong algorithm, wrong version, wrong key id, wrong audience, future-issued token, player missing, player disabled, player deleted, lookup miss, and verifier mismatch collapse to the same public invalid-token family unless a later disclosure decision changes this.

## 6. Repository Handoff

Future validation behavior must use existing repository interfaces. This gate does not authorize repository interface changes.

Required future repository access:

```yaml
repository_handoff:
  transaction_boundary: UnitOfWorkRunner.WithinUnitOfWork
  authentication_repository_source: unit_of_work.NewAuthenticationRepository
  player_repository_source: unit_of_work.NewPlayerAccountRepository
  authentication_lookup_method: FindTokenByLookupDigest
  token_mutation_method: deferred
  player_lookup_method: GetPlayerAccount
  direct_postgres_import: forbidden
  repository_interface_change: forbidden_by_this_gate
```

The service must not import PostgreSQL adapters, pgx, goose, WebSocket transport packages, Protobuf adapter packages, generated Protobuf packages, or generated contract-shape packages for validation behavior.

The authentication module remains storage-neutral. It stores and reads already-computed digest records only. It must not parse access-token carriers, compute digests, compare verifiers, decide proof validity, collapse public failures, or construct application responses.

Updating `LastValidatedAt` or `LastFailedValidationAt` is deferred unless a later work item authorizes a repository mutation shape and storage behavior.

## 7. Request Identity Handoff

Successful validation must create application-owned request identity before production-sensitive domain dispatch.

```yaml
request_identity_handoff:
  owner: runtime/internal/app
  source: AccessTokenValidationResult
  target_type: RequestIdentity
  success_validation_status: validated
  success_proof_status: valid
  identity_status: validated
  actor_kind: player
  actor_id: token.player_id
  player_id: token.player_id
  player_id_validated: true
  session_validated: false_until_session_persistence_gate
  session_id_source: request_metadata_only_if_present
  connection_id_source: request_metadata_only_if_present
  connection_epoch_source: request_metadata_only_if_present
  metadata_only_allowed_as_proof: false
```

Rules:

- Future code must keep SessionValidated false until session persistence is ratified; it must not mark `SessionValidated` as true from access-token validation alone.
- If it uses application identity helpers, it must override or avoid any helper behavior that treats metadata-only session id as validated proof.
- `ConnectionID`, `ConnectionEpoch`, route fields, and session metadata may be copied as context, but they are not proof.
- Domain modules must consume `RequestIdentity`; they must not parse access tokens, compare verifier digests, select token records, or decide authentication proof validity.

## 8. Public Error Collapse

Future behavior must preserve redacted internal classes while collapsing public proof failures.

```yaml
public_error_collapse:
  missing_token: AUTHENTICATION_TOKEN_MISSING
  malformed_token: AUTHENTICATION_TOKEN_MALFORMED
  lookup_miss: AUTHENTICATION_TOKEN_INVALID
  wrong_token_kind: AUTHENTICATION_TOKEN_INVALID
  inactive_token: AUTHENTICATION_TOKEN_INVALID
  expired_token: AUTHENTICATION_TOKEN_INVALID
  revoked_token: AUTHENTICATION_TOKEN_INVALID
  unsupported_algorithm: AUTHENTICATION_TOKEN_INVALID
  unsupported_version: AUTHENTICATION_TOKEN_INVALID
  verifier_key_id_mismatch: AUTHENTICATION_TOKEN_INVALID
  audience_mismatch: AUTHENTICATION_TOKEN_INVALID
  future_issued_token: AUTHENTICATION_TOKEN_INVALID
  verifier_mismatch: AUTHENTICATION_TOKEN_INVALID
  player_missing_or_inactive: AUTHENTICATION_TOKEN_INVALID
  repository_or_unit_of_work_unavailable: AUTHENTICATION_TOKEN_STORE_UNAVAILABLE
```

`AUTHENTICATION_TOKEN_EXPIRED`, `AUTHENTICATION_TOKEN_REVOKED`, and `AUTHENTICATION_ACCOUNT_DISABLED` exist in the semantic catalog. The first validation behavior should not publicly distinguish them unless a later disclosure decision explicitly authorizes that posture.

The public surface must not reveal whether a token record exists, whether the token is expired or revoked, whether the player exists, which verifier key id is stored, whether the audience matched, or whether a verifier mismatch happened.

## 9. Redaction Requirements

Future behavior must never place these values in errors, logs, docs examples, test failure messages, conversation logs, ADRs, change specs, or public responses:

- Raw access-token text.
- Raw access-token bytes.
- Lookup digest bytes.
- Verifier digest bytes.
- HMAC input or output bytes.
- Verifier key bytes.
- Full concrete `verifier_key_id` values.
- Token lookup hit or miss details.
- Token lifecycle details for failed proof.
- Audience mismatch details for failed proof.
- Player lookup hit or miss details.
- Repository implementation details for a failed proof.

Allowed:

- Registered public error codes.
- Redacted placeholders such as `<access-token>` and `<verifier-key-id>`.
- Non-secret record ids only after proof succeeds and only in authorized result fields or internal tests.

## 10. Required Tests

Future implementation must add or update focused tests in `runtime/internal/app/authentication/service_test.go`.

Required test classes:

```yaml
required_tests:
  validation_rejects_missing_token_without_unit_of_work
  validation_rejects_malformed_token_without_unit_of_work
  validation_computes_lookup_digest_before_repository_lookup
  validation_uses_authentication_repository_from_unit_of_work_only
  validation_uses_player_repository_from_unit_of_work_only
  validation_collapses_lookup_miss_to_invalid_token
  validation_rejects_inactive_expired_or_revoked_token
  validation_rejects_wrong_kind_algorithm_version_key_or_audience
  validation_compares_verifier_digest_before_request_identity
  validation_collapses_verifier_mismatch_to_invalid_token
  validation_requires_active_player_account
  validation_returns_request_identity_only_after_unit_of_work_success
  validation_keeps_session_validated_false_without_session_persistence
  validation_errors_do_not_leak_raw_token_or_digest_material
  validation_does_not_touch_protocol_carriers_or_route_protection
```

Normal tests must use fake unit-of-work and repository implementations. They must not require a live PostgreSQL server.

## 11. Nakama And Pitaya Mapping

Nakama capability mapping:

- Session/access-token validation is adapted as application-owned proof validation.
- Token expiration and revocation are checked inside the service validation flow.
- Server-side account state remains authoritative for whether a token can produce player identity.

Pitaya capability mapping:

- Realtime handlers should receive identity context after validation.
- Frontend acceptors and backend handlers remain separated from proof validation.
- Session binding is runtime context, not proof by itself.

These references guide capability coverage only. They do not authorize copying Nakama or Pitaya public APIs.

## 12. Deferrals

This gate does not authorize:

- Access-token validation execution.
- WebSocket proof carriers.
- WebSocket handshake authentication.
- Route protection.
- Session persistence.
- Logout execution.
- Refresh execution.
- Token cleanup jobs.
- Token validation audit mutation.
- Protocol authentication messages.
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
runtime.access_token_validation_service_behavior_gate
```
