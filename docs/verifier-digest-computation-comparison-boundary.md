# Verifier Digest Computation And Comparison Boundary

Status: Draft v0.1
Last updated: 2026-05-15
Scope: Future verifier digest computation ownership, canonical byte input construction, purpose-label use, logical key use, key-id selection, lookup digest handoff, verifier digest comparison, failure redaction, dependency posture, and test expectations for the first device-credential and opaque access-token posture
Depends on: `docs/token-credential-verifier-algorithm-redaction-boundary.md`, `docs/secret-configuration-verifier-key-loading-boundary.md`, `docs/token-credential-material-generation-boundary.md`
Canonical decision: `ADR-0043`

The paired Simplified Chinese translation is `docs/verifier-digest-computation-comparison-boundary.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This standard defines how future application-owned authentication code may compute lookup digests, compute verifier digests, and compare verifier digests.

It exists before verifier digest helper code, verifier comparison code, token generation code, credential generation code, secret loading, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository interface changes, migration schema changes, or production authentication behavior.

This is a boundary-only standard. It does not add Go code, imports, runtime services, HMAC helpers, digest helpers, comparison helpers, token generation, credential generation, secret loading, repository methods, SQL migrations, Protobuf messages, WebSocket carriers, routes, or production authentication behavior.

## 2. Core Rule

The first verifier digest computation and comparison posture is:

```yaml
verifier_digest_computation_comparison_boundary: defined
implementation_authorized_by_this_standard: false
planned_owner: runtime/internal/app
planned_optional_child_package: runtime/internal/app/authentication
verifier_algorithm_family: vibit_hmac_sha256_v1
verifier_version: 1
digest_size_bytes: 32
digest_output_shape: raw_32_byte_digest
canonical_input_version: vibit.auth.verifier.input.v1
canonical_input_encoding: ascii_header_length_prefixed_purpose_label_length_prefixed_raw_material
lookup_digest_database_equality_allowed_for_record_selection: true
lookup_digest_database_equality_sufficient_for_authentication: false
constant_time_verifier_comparison_required: true
missing_lookup_public_failure: invalid_authentication_proof
invalid_verifier_public_failure: invalid_authentication_proof
unknown_key_id_public_failure: invalid_authentication_proof
external_cryptography_dependency_required_for_first_posture: false
future_allowed_standard_library_packages:
  - crypto/hmac
  - crypto/sha256
  - crypto/subtle
```

Future first-posture digest helpers may use Go standard library `crypto/hmac`, `crypto/sha256`, and `crypto/subtle` after a later code gate authorizes implementation. No external cryptography, password-hashing, JWT, OAuth, OIDC, provider, KMS, cloud secret-manager, or operations dependency is required for this digest computation and comparison posture.

## 3. Ownership

Future verifier digest computation and comparison is application-owned:

```text
runtime/internal/app
```

Future helper code may live in an application-owned child package such as:

```text
runtime/internal/app/authentication
```

Ownership rules:

- Application authentication code owns digest computation and verifier comparison after a later implementation gate.
- `authentication.Repository` stores and retrieves already-computed digest material only.
- PostgreSQL adapters may use lookup digest equality for record selection, but they do not compute digests, compare verifier digests, choose keys, or decide authentication outcomes.
- Protobuf adapters and WebSocket transports may carry already-decoded proof fields after later protocol gates, but they do not compute digests or compare verifiers.
- Generated authentication contract shapes remain metadata-only and immutable.

Digest computation and comparison must not be placed in transport handlers, protocol adapters, domain modules, repositories, generated output, migrations, SQL fixtures, or tests that are not explicitly scoped to future application-owned digest helpers.

## 4. Canonical Digest Input

Future digest helpers must build a deterministic canonical byte input before HMAC computation.

The first planned canonical input is:

```yaml
canonical_digest_input:
  version_header_ascii: vibit.auth.verifier.input.v1
  header_separator: 0x00
  purpose_label_length: uint16_big_endian_byte_length
  purpose_label: ascii_bytes
  raw_material_length: uint16_big_endian_byte_length
  raw_material: generated_secret_material_bytes
```

The byte sequence is:

```text
ascii("vibit.auth.verifier.input.v1")
|| 0x00
|| uint16be(len(purpose_label))
|| ascii(purpose_label)
|| uint16be(len(raw_material))
|| raw_material
```

Rules:

- Purpose labels must be ASCII and must match the registered labels for the digest class.
- Raw material must be the decoded raw secret material bytes, not normalized text, metadata, player id, session id, route name, or provider subject.
- The raw material length for the first posture is 32 bytes.
- Length prefixes make the input unambiguous even if future raw material shape changes.
- The version header is part of the HMAC input and must change if the canonical input shape changes.
- Future tests must include deterministic fixture vectors for canonical byte construction.

This standard defines the input shape for future helpers. It does not implement those helpers.

## 5. Digest Classes

The first posture has four digest classes:

```yaml
credential_lookup_digest:
  purpose_label: vibit.credential.lookup.v1
  logical_key: credential_lookup_key
  output_bytes: 32
  storage_column: authentication_device_credentials.credential_lookup_digest
  database_equality_for_selection: allowed
  log_safe: false
credential_verifier_digest:
  purpose_label: vibit.credential.verifier.v1
  logical_key: credential_verifier_key
  output_bytes: 32
  storage_column: authentication_device_credentials.credential_verifier_digest
  constant_time_comparison_required: true
  log_safe: false
token_lookup_digest:
  purpose_label: vibit.access_token.lookup.v1
  logical_key: token_lookup_key
  output_bytes: 32
  storage_column: authentication_access_tokens.token_lookup_digest
  database_equality_for_selection: allowed
  log_safe: false
token_verifier_digest:
  purpose_label: vibit.access_token.verifier.v1
  logical_key: token_verifier_key
  output_bytes: 32
  storage_column: authentication_access_tokens.token_verifier_digest
  constant_time_comparison_required: true
  log_safe: false
```

Rules:

- Lookup digests and verifier digests must use different purpose labels.
- Credential digests and token digests must use different purpose labels.
- Each digest class must use the matching logical key from the selected verifier key set.
- Digest bytes must not be truncated for storage, lookup, comparison, logs, public errors, metrics labels, conversation logs, ADRs, or change specs.
- Any future digest prefix or fingerprint format requires a separate redaction/fingerprint standard.

## 6. Key Selection

New credential and token writes must use the active verifier key set selected by future secret configuration code.

Future validation of presented credential or access-token proof must consider accepted key sets because the opaque proof text does not carry `verifier_key_id`.

Planned validation posture:

```yaml
new_write_key_selection:
  key_set: active_key_set
  stored_record_keeps_verifier_key_id: true
validation_lookup_key_selection:
  candidate_key_sets: active_and_accepted_previous_key_sets
  compute_lookup_digest_per_candidate_key_set: required
  repository_lookup: by_lookup_digest
validation_verifier_key_selection:
  selected_by_stored_record_verifier_key_id: required
  unknown_key_id_public_failure: invalid_authentication_proof
  retired_key_id_public_failure: invalid_authentication_proof
```

Rules:

- Future validation may compute multiple lookup digests for accepted key sets.
- Repository lookup by lookup digest selects a candidate record only; it is not authentication proof.
- After a record is selected, future code must compute the matching verifier digest with the key set identified by the stored `verifier_key_id`.
- If the stored key id is unknown, retired, malformed, or unavailable, public behavior must use the same invalid proof failure class.
- No public error may reveal the active key set id, previous key set ids, candidate count, lookup miss, key miss, rotation state, or digest value.

This standard does not change the repository interface. If a later implementation needs batch lookup by multiple digest candidates, that work item must update the repository boundary explicitly.

## 7. Comparison Boundary

Future verifier comparison must compare computed verifier digest bytes against stored verifier digest bytes using constant-time equality.

Accepted future Go comparison primitives:

```yaml
preferred_go_comparison: crypto/hmac.Equal
acceptable_go_comparison: crypto/subtle.ConstantTimeCompare
forbidden_comparison:
  - bytes.Equal
  - string equality
  - byte_slice_string_conversion
  - reflect.DeepEqual
  - database_only_equality
  - map_lookup_equality
```

Rules:

- Compare verifier digest bytes, not raw credential text, raw access-token text, encoded material, lookup digest bytes, key ids, or public identifiers.
- Database equality on lookup digest is allowed only for candidate record selection.
- A lookup hit must still pass lifecycle checks and constant-time verifier digest comparison before proof is accepted.
- A lookup miss, missing stored verifier digest, mismatched verifier digest, unknown key id, unsupported algorithm version, expired token, revoked token, disabled credential, or disabled account must not expose a more detailed public failure reason unless a later semantic error standard explicitly allows it.
- Future code should avoid obvious timing differences between missing records and invalid verifier digests; exact equalization strategy belongs to the later implementation gate.

## 8. Failure Redaction

Future public failures for credential and token verifier problems must collapse to a registered invalid-proof class unless a later semantic standard explicitly allows a more specific public class.

Required public failure posture:

```yaml
public_failure_class:
  missing_lookup_record: invalid_authentication_proof
  verifier_digest_mismatch: invalid_authentication_proof
  unknown_verifier_key_id: invalid_authentication_proof
  unsupported_verifier_algorithm: invalid_authentication_proof
  malformed_presented_proof: invalid_authentication_proof
  expired_or_revoked_token: invalid_authentication_proof
```

Forbidden in public errors, logs intended for client support, audit-safe facts, traces, metrics labels, ADRs, change specs, documentation examples, and conversation logs:

- Raw credential material.
- Raw access-token material.
- Encoded credential or token material.
- Lookup digest bytes.
- Verifier digest bytes.
- Verifier key values.
- Encoded verifier key values.
- Full concrete `verifier_key_id` values.
- Candidate key-set counts.
- HMAC input bytes.
- HMAC output bytes.
- A reason that distinguishes lookup miss from verifier mismatch.

Allowed with care:

- Registered error codes.
- Non-secret lifecycle state names when the surrounding artifact already permits them.
- Placeholders such as `<lookup-digest>` or `<verifier-key-id>` when documenting redaction rules.

## 9. Test Expectations

Future implementation work that computes lookup digests, computes verifier digests, or compares verifier digests must add focused tests.

Minimum expectations:

```yaml
digest_tests:
  canonical_input_is_deterministic: required
  canonical_input_uses_version_header: required
  canonical_input_length_prefixes_purpose_label: required
  canonical_input_length_prefixes_raw_material: required
  lookup_and_verifier_purpose_labels_differ: required
  credential_and_token_purpose_labels_differ: required
  digest_output_is_32_bytes: required
  lookup_digest_uses_lookup_key: required
  verifier_digest_uses_verifier_key: required
  lookup_digest_not_used_as_authentication_proof: required
  verifier_comparison_uses_constant_time_primitive: required
  missing_record_and_mismatch_share_public_failure: required
  raw_material_absent_from_outputs: required
  digest_material_absent_from_outputs: required
```

Test fixtures may use deterministic synthetic keys and raw material only inside tests. They must not become production defaults, documentation examples of real secret shape, or committed production-like secret values.

## 10. Dependency Posture

The first digest computation and comparison posture does not require an external dependency.

Accepted for later implementation after an explicit code gate:

```yaml
go_standard_library:
  hmac: crypto/hmac
  hash: crypto/sha256
  constant_time: crypto/subtle
external_dependency_adoption_record_required_for_first_posture: false
```

Deferred and not required for the first posture:

- External cryptography libraries.
- Password hashing dependencies.
- JWT, JWK, or signing libraries.
- OAuth or OIDC provider SDKs.
- KMS or cloud secret-manager SDKs.
- Redis-like token/session stores.

If a future login method accepts password-like or low-entropy material, it must define a separate credential boundary and dependency adoption record.

## 11. Non-Goals

This standard does not:

- Add verifier digest computation.
- Add verifier comparison.
- Add token generation.
- Add credential generation.
- Add secret loading.
- Add application authentication service code.
- Add login execution.
- Add access-token validation.
- Add logout execution.
- Add refresh behavior.
- Add cleanup jobs.
- Add Protobuf authentication messages.
- Add WebSocket proof carriers.
- Change the WebSocket handshake.
- Add authentication dependencies.
- Change `authentication.Repository`.
- Change PostgreSQL migration schemas.
- Add production authentication behavior.

## 12. Verification Path

The repository check rule for this boundary is:

```text
runtime.verifier_digest_computation_comparison_boundary
```

For changes that touch this boundary, run:

```bash
node tools/vibit check runtime --json
node tools/vibit check contracts --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check all --json
```

No runtime Go verifier digest computation or comparison behavior is verified by this standard because no behavior is added.

## 13. Follow-Up Gates

Recommended follow-up gates:

- Authentication service implementation readiness gate.
- Local verifier key configuration implementation gate.
- Token and credential material generation implementation gate.
- Verifier digest helper implementation gate.
- Authentication redaction test implementation gate.
- Protobuf authentication message gate.
- WebSocket request proof carrier gate.
