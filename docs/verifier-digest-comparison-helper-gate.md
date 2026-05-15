# Verifier Digest Comparison Helper Gate

Status: Draft v0.1
Last updated: 2026-05-15
Scope: Verifier digest comparison helper implementation gate, future helper file boundaries, constant-time primitive posture, input redaction, failure posture, tests, and deferrals before comparison code is added
Depends on: `docs/verifier-digest-computation-comparison-boundary.md`, `docs/verifier-digest-helper-implementation-gate.md`, `docs/token-credential-verifier-algorithm-redaction-boundary.md`, `docs/authentication-service-implementation-readiness-gate.md`
Canonical decision: `ADR-0049`

The paired Simplified Chinese translation is `docs/verifier-digest-comparison-helper-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This gate defines the next bounded implementation slice for verifier digest comparison helpers.

The repository already has helper-only verifier digest computation under `runtime/internal/app/authentication`. The next risk is letting a future agent compare verifier material with non-constant-time primitives, compare the wrong material, hide authentication service behavior inside a helper, or move comparison into repositories, protocol adapters, transport handlers, or generated output.

This is an implementation-gate standard. It does not add Go code, imports, verifier comparison, authentication service behavior, login execution, token validation, logout execution, refresh behavior, cleanup jobs, Protobuf messages, WebSocket proof carriers, repository methods, SQL migrations, startup wiring, authentication dependencies, external cryptography services, KMS, cloud secret-manager integrations, or production authentication behavior.

## 2. Core Rule

The verifier digest comparison helper gate is:

```yaml
verifier_digest_comparison_helper_gate: defined
implementation_authorized_by_this_standard: false
future_implementation_work_item: W-0105
planned_owner: runtime/internal/app
planned_package: runtime/internal/app/authentication
future_slice: constant_time_verifier_digest_comparison_helpers
future_source: runtime/internal/app/authentication/verifier_comparison.go
future_tests: runtime/internal/app/authentication/verifier_comparison_test.go
verifier_algorithm_family: vibit_hmac_sha256_v1
verifier_version: 1
digest_size_bytes: 32
input_shape: computed_verifier_digest_and_stored_verifier_digest_bytes
computed_digest_handoff: ComputedDigest
stored_digest_handoff: copied_repository_digest_bytes
preferred_constant_time_primitive: crypto/hmac.Equal
acceptable_constant_time_primitive: crypto/subtle.ConstantTimeCompare
comparison_result_shape: redacted_match_or_mismatch_result
raw_material_comparison: forbidden
lookup_digest_comparison_for_authentication: forbidden
database_only_verifier_comparison: forbidden
service_behavior_status: deferred
authentication_service_behavior_status: deferred
protobuf_authentication_messages_status: deferred
websocket_proof_carrier_status: deferred
repository_interface_change_status: deferred
migration_schema_change_status: deferred
startup_wiring_status: deferred
authentication_dependencies_status: deferred
production_authentication_behavior_status: deferred
```

The future implementation must be a helper-only slice. It may compare a computed verifier digest to stored verifier digest bytes with a constant-time primitive and return a redacted comparison result. It must not compute digests, load keys, select records, parse proofs, issue login responses, validate tokens, revoke tokens, refresh tokens, call repositories, inspect lifecycle state, or touch protocol carriers.

## 3. Future Helper Shape

Future implementation ownership:

```text
runtime/internal/app/authentication
```

Allowed future files after the implementation work item authorizes code:

```text
runtime/internal/app/authentication/verifier_comparison.go
runtime/internal/app/authentication/verifier_comparison_test.go
```

`verifier_digest.go` remains computation-only. Comparison helpers must live in the separate comparison file so later agents can reason about computation and comparison independently.

Preferred future API shape:

```yaml
future_types:
  VerifierComparisonResult:
    owns: digest_class_and_match_status_only
    methods:
      - Class() DigestClass
      - Matched() bool
    constraints:
      - Result text and string formatting must not expose digest bytes, raw material, key ids, account ids, token ids, or lookup miss details.
      - A false match is not a public authentication error by itself; service behavior maps it later.

future_functions:
  CompareCredentialVerifierDigest:
    signature: "func CompareCredentialVerifierDigest(computed ComputedDigest, storedDigest []byte) (VerifierComparisonResult, error)"
    expected_computed_class: credential_verifier
    behavior: constant_time_compare_computed_verifier_digest_to_stored_credential_verifier_digest
  CompareTokenVerifierDigest:
    signature: "func CompareTokenVerifierDigest(computed ComputedDigest, storedDigest []byte) (VerifierComparisonResult, error)"
    expected_computed_class: token_verifier
    behavior: constant_time_compare_computed_verifier_digest_to_stored_token_verifier_digest
```

An unexported shared helper is acceptable if it keeps the exported functions class-specific and prevents accidental lookup-digest comparison.

## 4. Input Boundary

Allowed inputs:

- A `ComputedDigest` produced by the W-0103 digest helpers.
- Stored verifier digest bytes returned from repository records.

Forbidden inputs:

- Raw device credential material.
- Raw access-token material.
- Encoded credential or token text.
- Lookup digest bytes.
- `credential_lookup_digest`.
- `token_lookup_digest`.
- `verifier_key_id`.
- Player account id.
- Token record id.
- Credential record id.
- Provider subject.
- Session id.
- WebSocket connection metadata.
- Route name or protocol metadata.

The comparison helper must compare only verifier digest bytes. Lookup digest equality remains record selection only and is not proof of authentication.

## 5. Constant-Time Comparison Posture

Preferred primitive:

```yaml
preferred_go_comparison: crypto/hmac.Equal
```

Acceptable primitive:

```yaml
acceptable_go_comparison: crypto/subtle.ConstantTimeCompare
```

Forbidden for verifier digest comparison:

- `bytes.Equal`
- `reflect.DeepEqual`
- `==` on strings or arrays derived from digest bytes
- byte-slice to string conversion for equality
- map lookup equality
- database-only equality
- SQL comparison of verifier digest bytes as the final proof
- comparing encoded digest text
- comparing raw credential or token material

Length validation must fail closed. The first posture requires 32-byte verifier digests. The helper may reject missing or malformed digest input with a redacted error, but it must not disclose which side of an authentication attempt failed in a future public response.

## 6. Failure Posture

The helper must preserve a small internal distinction for tests while supporting public failure collapse later:

```yaml
internal_failure_classes:
  verifier_digest_mismatch: redacted
  missing_computed_digest: redacted
  wrong_computed_digest_class: redacted
  missing_stored_digest: redacted
  malformed_stored_digest: redacted
  invalid_comparison_input: redacted
future_public_failure_class: invalid_authentication_proof
```

Rules:

- Mismatch must not expose whether lookup succeeded.
- Missing stored verifier digest must not expose record existence or corruption through public service behavior.
- Wrong digest class must fail closed.
- Malformed stored digest must fail closed.
- Error text must not include digest bytes, raw material, key ids, account ids, credential ids, token ids, lookup values, or candidate counts.
- The helper must not decide whether a credential is active, expired, revoked, usable for a player, or bound to a session.

## 7. Package And File Boundary

Allowed future implementation area:

```text
runtime/internal/app/authentication/verifier_comparison.go
runtime/internal/app/authentication/verifier_comparison_test.go
```

Forbidden write areas for the future verifier comparison helper slice:

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

The future implementation must not wire comparison into login execution, account creation, token issuance, token validation, startup, WebSocket transport, Protobuf protocol, PostgreSQL persistence, generated contract output, migrations, or domain repository code.

## 8. Error And Redaction Requirements

Allowed error classification examples:

```yaml
error_classes:
  verifier_digest_mismatch
  missing_computed_digest
  wrong_computed_digest_class
  missing_stored_digest
  malformed_stored_digest
  invalid_verifier_comparison_input
```

Allowed in error text:

- Error classes.
- Digest class names.
- Non-secret numeric expectations such as `32` bytes.

Forbidden in errors, logs, test snapshots, docs, ADRs, change specs, and conversation logs:

- Raw credential material or bytes.
- Raw access-token material or bytes.
- Encoded generated material.
- Lookup digest bytes.
- Verifier digest bytes.
- HMAC input bytes.
- HMAC output bytes.
- Stored verifier digest bytes.
- Verifier key values.
- Encoded verifier key values.
- Full concrete `verifier_key_id` values.
- Account ids, credential record ids, or token record ids from real data.
- Candidate key-set counts.
- Lookup hit or miss details.

## 9. Required Tests For The Future Helper

Future implementation must add focused unit tests under:

```text
runtime/internal/app/authentication/verifier_comparison_test.go
```

Minimum test cases:

```yaml
required_tests:
  credential_verifier_digest_match_returns_matched
  credential_verifier_digest_mismatch_returns_not_matched
  token_verifier_digest_match_returns_matched
  token_verifier_digest_mismatch_returns_not_matched
  comparison_uses_crypto_hmac_equal
  comparison_rejects_lookup_digest_classes
  comparison_rejects_wrong_verifier_digest_class
  comparison_rejects_missing_computed_digest
  comparison_rejects_missing_stored_digest
  comparison_rejects_malformed_stored_digest_length
  comparison_does_not_compare_raw_material
  comparison_does_not_call_repositories
  comparison_result_does_not_expose_digest_bytes
  comparison_errors_are_redacted
  comparison_helpers_do_not_implement_authentication_service_behavior
```

Tests may use deterministic synthetic digest bytes only inside tests. They must not record production-like secrets, real credential material, real access tokens, real key ids, account ids, credential record ids, token record ids, or lookup digests.

## 10. Dependency Posture

No new external dependency is allowed by this gate.

Allowed Go standard library packages after the future implementation work item authorizes code:

```yaml
future_standard_library_imports_allowed:
  - crypto/hmac
  - crypto/subtle
  - errors
  - fmt
```

The first helper implementation must not add JWT, JWK, OAuth, OIDC, provider SDKs, password-hashing dependencies, Redis-like stores, KMS SDKs, cloud secret-manager SDKs, operations libraries, or external cryptography services.

## 11. Nakama And Pitaya Mapping

Nakama capability reference:

- Server-side credential and token validation needs a final verifier comparison step after lookup and digest computation. vibit adopts the capability need, not Nakama's public API shape.

Pitaya capability reference:

- Realtime route handlers should receive identity context after validation. vibit keeps verifier digest comparison in application-owned helpers, not in frontend acceptors, route handlers, or session binding code.

This gate maps those references into a narrow helper slice: comparison only, service behavior later.

## 12. Non-Goals

This gate does not:

- Add verifier comparison code.
- Add digest computation code.
- Add token generation orchestration.
- Add credential generation orchestration.
- Implement authentication service behavior.
- Execute login.
- Validate access tokens.
- Execute logout.
- Add refresh behavior.
- Add cleanup jobs.
- Add Protobuf authentication messages.
- Add WebSocket proof carriers.
- Change WebSocket handshake behavior.
- Change `authentication.Repository`.
- Change PostgreSQL migration schemas.
- Add startup wiring.
- Add authentication dependencies.
- Add production authentication behavior.

## 13. Verification Path

The repository check rule for this gate is:

```text
runtime.verifier_digest_comparison_helper_gate
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

No runtime Go verifier digest comparison behavior is verified by this gate because no comparison behavior is added.

## 14. Follow-Up Gates

Recommended follow-up gates:

- Implement verifier digest comparison helpers.
- Implement application authentication service behavior.
- Add Protobuf authentication messages.
- Add WebSocket request proof carriers.
