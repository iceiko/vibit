# Authentication Service Behavior Implementation Gate

Status: Draft v0.1
Last updated: 2026-05-16
Scope: Future application authentication service behavior ownership, file boundaries, repository handoff, helper composition flow, public error collapse, redaction, tests, and deferrals before service behavior code is added
Depends on: `docs/runtime-authentication-implementation-boundary.md`, `docs/application-authentication-service-interface-boundary.md`, `docs/token-credential-verifier-algorithm-redaction-boundary.md`, `docs/secret-configuration-verifier-key-loading-boundary.md`, `docs/token-credential-material-generation-boundary.md`, `docs/verifier-digest-computation-comparison-boundary.md`, `docs/authentication-service-implementation-readiness-gate.md`, `docs/verifier-digest-comparison-helper-gate.md`
Canonical decision: `ADR-0050`

The paired Simplified Chinese translation is `docs/authentication-service-behavior-implementation-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This gate defines the future authentication service behavior slice after the verifier key, material generation, digest computation, and verifier comparison helpers exist.

The repository now has the small application-owned helper chain needed for the first selected posture:

```text
VerifierKeySet
-> raw credential/token material generation
-> lookup and verifier digest computation
-> constant-time verifier digest comparison
```

The next risk is letting a future agent wire those helpers directly into transport, Protobuf adapters, repositories, startup, or generated files, or implement login and token validation while inventing error mapping, repository call ordering, proof redaction, key selection, or request identity handoff.

This is an implementation-gate standard. It does not add Go service code, login execution, access-token validation, logout execution, refresh behavior, cleanup jobs, protocol carriers, repository interface changes, SQL migrations, startup wiring, external dependencies, or production authentication behavior.

## 2. Core Rule

The authentication service behavior implementation gate is:

```yaml
authentication_service_behavior_implementation_gate: defined
implementation_authorized_by_this_standard: false
future_implementation_work_item: W-0107
planned_owner: runtime/internal/app
planned_package: runtime/internal/app/authentication
future_slice: authentication_service_behavior_skeleton
future_source: runtime/internal/app/authentication/service.go
future_tests: runtime/internal/app/authentication/service_test.go
repository_handoff: application_unit_of_work
helper_composition_flow_defined: true
public_error_collapse_defined: true
request_identity_handoff_defined: true
service_behavior_status: gated
login_execution_status: deferred
token_validation_status: deferred
logout_execution_status: deferred
refresh_behavior_status: deferred
cleanup_execution_status: deferred
protocol_carrier_status: deferred
protobuf_authentication_messages_status: deferred
websocket_proof_carrier_status: deferred
repository_interface_change_status: deferred
migration_schema_change_status: deferred
startup_wiring_status: deferred
authentication_dependencies_status: deferred
production_authentication_behavior_status: deferred
```

The future service behavior must be application-owned. It may orchestrate existing helper outputs, repository interfaces, and request identity handoff only after a later bounded work item explicitly authorizes the exact slice.

## 3. Future Service Ownership

Future service behavior owner:

```text
runtime/internal/app/authentication
```

Allowed future files after the implementation work item authorizes code:

```text
runtime/internal/app/authentication/service.go
runtime/internal/app/authentication/service_test.go
```

The first future service code slice should be a skeleton only unless a later work item explicitly authorizes real login or token validation behavior. The skeleton may define typed dependency boundaries, request/result vocabulary, redacted internal error classes, and fail-closed `not implemented` behavior. It must not call repositories, issue tokens, validate tokens, revoke tokens, refresh tokens, clean up tokens, expose protocol carriers, or wire startup.

## 4. Repository Handoff

Authentication service behavior may use the storage-neutral authentication repository only through the application unit-of-work boundary.

Required future flow for state-changing behavior:

```text
application service method
-> application-owned unit-of-work runner
-> UnitOfWork.NewAuthenticationRepository(...)
-> authentication.Repository
-> persistence-only PostgreSQL adapter
```

Rules:

- The service owns orchestration.
- The repository owns storage-neutral record lookup and mutation only.
- The PostgreSQL adapter owns SQL persistence only.
- The service must not import PostgreSQL driver packages.
- The service must not bypass the unit-of-work boundary for state-changing authentication behavior.
- Repositories must not generate material, compute digests, compare verifiers, parse proof, map public failures, or construct `RequestIdentity`.
- PostgreSQL adapters must not decide authentication outcomes.

Repository interface changes remain deferred. If the service needs a repository method that does not exist, a later work item must update the module boundary and adapter tests before behavior code consumes it.

## 5. Helper Composition Flow

Future device credential login behavior, when separately authorized, must compose helpers in this order:

```yaml
device_credential_login_flow:
  input: already_decoded_service_request
  proof_shape: raw_device_credential_material
  steps:
    - reject_missing_or_malformed_proof_before_repository_lookup
    - compute_credential_lookup_digest_with_VerifierKeySet
    - find_credential_by_lookup_digest_through_unit_of_work_repository
    - collapse_lookup_miss_or_unusable_record_to_public_invalid_credential
    - compute_credential_verifier_digest_with_record_verifier_key_context
    - compare_credential_verifier_digest_with_CompareCredentialVerifierDigest
    - collapse_mismatch_to_public_invalid_credential
    - require_active_credential_and_allowed_player_account_state
    - generate_access_token_material_with_explicit_entropy_reader
    - compute_token_lookup_and_verifier_digests
    - store_token_record_through_unit_of_work_repository
    - return_raw_access_token_text_once_to_the_authorized_response_carrier
  output: redacted_authentication_result
```

Future access-token validation behavior, when separately authorized, must compose helpers in this order:

```yaml
access_token_validation_flow:
  input: already_decoded_explicit_request_proof_payload
  proof_shape: raw_access_token_material
  steps:
    - reject_missing_or_malformed_proof_before_repository_lookup
    - compute_token_lookup_digest_with_VerifierKeySet
    - find_token_by_lookup_digest_through_unit_of_work_repository
    - collapse_lookup_miss_or_unusable_record_to_public_invalid_token
    - compute_token_verifier_digest_with_record_verifier_key_context
    - compare_token_verifier_digest_with_CompareTokenVerifierDigest
    - collapse_mismatch_to_public_invalid_token
    - require_active_token_lifecycle_state_and_unexpired_window
    - convert_validated_actor_to_RequestIdentity
  output: application_owned_request_identity
```

This gate records the intended composition only. It does not authorize either flow to execute.

## 6. Public Error Collapse

The service behavior must keep internal proof distinctions useful for tests while collapsing public proof failures.

Internal failure classes may include:

```yaml
internal_failure_classes:
  missing_proof: redacted
  malformed_proof: redacted
  lookup_miss: redacted
  wrong_verifier_algorithm: redacted
  unknown_verifier_key_id: redacted
  unsupported_verifier_version: redacted
  verifier_digest_mismatch: redacted
  credential_not_active: redacted
  token_not_active: redacted
  token_expired: redacted
  token_revoked: redacted
  repository_unavailable: redacted_dependency
```

Required first public collapse:

```yaml
public_error_collapse:
  missing_device_credential_proof: AUTHENTICATION_PROOF_MISSING
  malformed_device_credential_proof: AUTHENTICATION_PROOF_MALFORMED
  invalid_device_credential_proof_family: AUTHENTICATION_CREDENTIAL_INVALID
  missing_access_token_proof: AUTHENTICATION_TOKEN_MISSING
  malformed_access_token_proof: AUTHENTICATION_TOKEN_MALFORMED
  invalid_access_token_proof_family: AUTHENTICATION_TOKEN_INVALID
  credential_store_unavailable: AUTHENTICATION_CREDENTIAL_STORE_UNAVAILABLE
  token_store_unavailable: AUTHENTICATION_TOKEN_STORE_UNAVAILABLE
  unsupported_refresh: AUTHENTICATION_REFRESH_NOT_SUPPORTED
  not_implemented: AUTHENTICATION_NOT_IMPLEMENTED
```

`AUTHENTICATION_TOKEN_EXPIRED`, `AUTHENTICATION_TOKEN_REVOKED`, and account-disabled public distinctions exist in the semantic catalog, but exposing those distinctions in first behavior remains a later explicit decision. A future work item may authorize more specific public mapping only after it confirms the disclosure posture.

## 7. Request Identity Handoff

Access-token validation must eventually convert proof into application-owned request identity before production-sensitive domain dispatch.

Target handoff:

```yaml
request_identity_handoff:
  owner: runtime/internal/app
  source: access_token_validation_result
  target_type: RequestIdentity
  success_status: authentication_proven
  actor_kind: player
  player_id_validated: true
  session_validated: false_until_session_persistence_gate
  metadata_only_allowed_as_proof: false
```

Domain modules must receive `RequestIdentity`. They must not parse access tokens, compare verifier digests, select credential records, or decide authentication proof validity.

## 8. Redaction Requirements

Forbidden in errors, logs, traces, metrics labels, tests snapshots, ADRs, change specs, documentation examples, conversation logs, and public responses except the one-time authorized token response carrier:

- Raw device credential material.
- Raw access-token material.
- Encoded generated credential material.
- Encoded generated access-token material outside the one-time response carrier.
- Lookup digest bytes.
- Verifier digest bytes.
- HMAC input bytes.
- HMAC output bytes.
- Verifier key values.
- Encoded verifier key values.
- Full concrete `verifier_key_id` values.
- Credential lookup hit or miss details.
- Token lookup hit or miss details.
- Candidate key-set counts.
- Repository implementation details for a failed proof.

Allowed:

- Registered public error codes.
- Non-secret record ids when the flow has already proven the actor and the target carrier is authorized.
- Redacted placeholders such as `<raw-access-token>` and `<verifier-key-id>`.
- Internal test-only failure class names when no secret or existence detail is exposed publicly.

## 9. File Boundary

Allowed future implementation area:

```text
runtime/internal/app/authentication/service.go
runtime/internal/app/authentication/service_test.go
```

Allowed existing helper dependencies after a later service code work item:

```text
runtime/internal/app/authentication/verifier_key_config.go
runtime/internal/app/authentication/verifier_key_env.go
runtime/internal/app/authentication/material_generation.go
runtime/internal/app/authentication/verifier_digest.go
runtime/internal/app/authentication/verifier_comparison.go
```

Forbidden write areas for the first service behavior slice unless a later work item explicitly names them:

- `runtime/cmd/vibit-server/`
- `runtime/internal/app/bootstrap/`
- `runtime/internal/platform/transport/ws/`
- `runtime/internal/platform/protocol/protobuf/`
- `runtime/internal/platform/persistence/postgres/`
- `runtime/internal/platform/migrations/`
- `runtime/internal/modules/authentication/`
- `runtime/internal/generated/`
- `runtime/migrations/postgres/`
- `proto/`
- `contracts/runtime/authentication/`

The future service behavior must not hand-edit generated files or hide service behavior in transport, protocol, repository, migration, or startup code.

## 10. Required Tests For Future Service Behavior

Future implementation must add focused tests under:

```text
runtime/internal/app/authentication/service_test.go
```

Minimum test classes before real behavior is accepted:

```yaml
required_tests:
  service_skeleton_fails_closed_without_behavior_authorization
  service_dependencies_reject_nil_or_missing_unit_of_work
  service_does_not_log_or_return_raw_proof_material
  credential_login_composes_lookup_digest_before_repository_lookup
  credential_login_compares_verifier_digest_before_token_issuance
  credential_login_collapses_lookup_miss_and_mismatch_to_same_public_error
  access_token_validation_composes_lookup_digest_before_repository_lookup
  access_token_validation_compares_verifier_digest_before_request_identity
  access_token_validation_collapses_lookup_miss_and_mismatch_to_same_public_error
  request_identity_is_populated_only_after_valid_access_token_proof
  repository_is_used_only_through_unit_of_work_boundary
  protocol_and_transport_packages_do_not_import_service_behavior
```

Tests must use fakes for repository and unit-of-work behavior unless a later live integration work item explicitly opts into PostgreSQL. Normal service behavior tests must not require a running PostgreSQL server.

## 11. Nakama And Pitaya Mapping

Nakama capabilities to adapt:

- Server-side account authentication.
- Session/access-token issuance.
- Session/access-token validation.
- Token expiration and revocation checks.
- Account disabled checks.

Pitaya capabilities to adapt:

- Frontend acceptors stay separated from backend handler logic.
- Realtime handlers receive identity/session context after validation.
- Session binding is runtime context, not proof by itself.

vibit rule:

```text
transport accepts frames
-> protocol decodes messages
-> application authentication validates proof
-> application request identity carries validated actor context
-> domain modules consume request identity
```

Do not copy Nakama or Pitaya public APIs directly. Use them as capability and vocabulary references while preserving vibit's agent-native boundaries.

## 12. Deferrals

This gate does not authorize:

- Go authentication service behavior code.
- Login execution.
- Access-token validation execution.
- Logout execution.
- Refresh execution.
- Token cleanup jobs.
- Credential bootstrap or account creation policy.
- Session persistence.
- Protocol authentication messages.
- WebSocket proof carriers.
- WebSocket handshake authentication.
- Repository interface changes.
- PostgreSQL adapter changes.
- SQL migration changes.
- Startup wiring.
- External dependencies.
- KMS or cloud secret-manager integration.
- Redis-like token/session stores.
- Production authentication behavior.

## 13. Verification

The repository check rule for this gate is:

```text
runtime.authentication_service_behavior_implementation_gate
```

For changes that touch this gate, run:

```bash
node tools/vibit check runtime --json
node tools/vibit check contracts --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check all --json
```

The check must preserve that `service.go` and `service_test.go` do not exist until a later work item explicitly authorizes service code.
