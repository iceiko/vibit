# Token And Credential Verifier Algorithm Redaction Boundary

Status: Draft v0.1
Last updated: 2026-05-15
Scope: Verifier algorithm posture, digest classification, key identifier treatment, constant-time comparison expectations, dependency posture, and redaction test expectations for the first device-credential and opaque access-token posture
Depends on: `docs/credential-record-schema-boundary.md`, `docs/token-verifier-record-schema-boundary.md`, `docs/application-authentication-service-interface-boundary.md`
Canonical decision: `ADR-0040`

The paired Simplified Chinese translation is `docs/token-credential-verifier-algorithm-redaction-boundary.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This standard defines the verifier algorithm and redaction boundary that future application-owned authentication service code must follow.

It exists before token material generation, verifier comparison, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository interface changes, migration schema changes, or production authentication behavior.

This is a boundary-only standard. It does not add Go code, imports, runtime services, token generation, digest comparison, key loading, configuration, repository methods, SQL migrations, Protobuf messages, WebSocket carriers, routes, or production authentication behavior.

## 2. Core Rule

The first verifier posture is:

```yaml
verifier_algorithm_family: vibit_hmac_sha256_v1
minimum_raw_token_entropy_bits: 256
minimum_raw_device_credential_entropy_bits: 256
token_text_encoding: url_safe_unpadded_base64_or_equivalent
credential_kind: high_entropy_installation_credential
credential_lookup_digest_required: true
credential_verifier_digest_required: true
token_lookup_digest_required: true
token_verifier_digest_required: true
raw_credential_storage: forbidden
raw_token_storage: forbidden
constant_time_verifier_comparison_required: true
password_hashing_dependency_required: false
external_cryptography_dependency_required: false
jwt_or_claim_parsing_required: false
oauth_oidc_dependency_required: false
kms_dependency_required_for_first_posture: false
implementation_authorized_by_this_standard: false
```

The first planned verifier algorithm can be implemented with the Go standard library after a later implementation work item authorizes code:

```yaml
future_allowed_standard_library_packages:
  - crypto/hmac
  - crypto/sha256
  - crypto/subtle
  - crypto/rand
  - encoding/base64
external_dependency_adoption_record_required_for_first_posture: false
```

`crypto/rand` and `encoding/base64` are named only because later token or credential material generation will need cryptographically secure randomness and safe text encoding. This standard does not authorize generation code.

## 3. Algorithm Identifiers

The first planned algorithm identifier is:

```yaml
verifier_algorithm: vibit_hmac_sha256_v1
verifier_version: 1
digest_size_bytes: 32
digest_storage_format: raw_32_byte_digest
test_fixture_text_encoding: base64url_unpadded_when_text_is_required
```

The algorithm identifier covers both lookup digests and verifier digests. Purpose labels and separate server-side keys distinguish each digest class.

Required purpose labels:

```yaml
purpose_labels:
  credential_lookup_digest: vibit.credential.lookup.v1
  credential_verifier_digest: vibit.credential.verifier.v1
  token_lookup_digest: vibit.access_token.lookup.v1
  token_verifier_digest: vibit.access_token.verifier.v1
```

Rules:

- Purpose labels are part of the digest input domain separation.
- Lookup and verifier digests must not reuse the same purpose label.
- Credential and token digests must not reuse the same purpose label.
- Future algorithm identifiers must be versioned before stored records use them.
- Stored records must retain `verifier_algorithm`, `verifier_version`, and `verifier_key_id`.

## 4. Digest Construction Posture

The planned digest construction is HMAC-SHA-256 over raw proof material with server-side secret keys and fixed purpose labels.

Future code must compute the digest over canonical byte input:

```yaml
digest_input:
  purpose_label: required
  separator: required_non_ambiguous_byte_separator
  raw_secret_material: credential_or_access_token_bytes
```

Planned digest classes:

```yaml
credential_lookup_digest:
  classification: secret_adjacent_index_material
  construction: HMAC-SHA-256(credential_lookup_key, purpose_label || separator || raw_credential)
  storage: authentication_device_credentials.credential_lookup_digest
  index_use: allowed
  log_safe: false
credential_verifier_digest:
  classification: secret_verifier_material
  construction: HMAC-SHA-256(credential_verifier_key, purpose_label || separator || raw_credential)
  storage: authentication_device_credentials.credential_verifier_digest
  constant_time_comparison_required: true
  log_safe: false
token_lookup_digest:
  classification: secret_adjacent_index_material
  construction: HMAC-SHA-256(token_lookup_key, purpose_label || separator || raw_access_token)
  storage: authentication_access_tokens.token_lookup_digest
  index_use: allowed
  log_safe: false
token_verifier_digest:
  classification: secret_verifier_material
  construction: HMAC-SHA-256(token_verifier_key, purpose_label || separator || raw_access_token)
  storage: authentication_access_tokens.token_verifier_digest
  constant_time_comparison_required: true
  log_safe: false
```

The exact byte separator and helper function names remain for a later code implementation gate. The implementation must make the byte construction deterministic, unambiguous, covered by tests, and independent of transport or Protobuf message shape.

Lookup digest equality through a database index is not proof by itself. Future validation must still perform status checks, expiration or revocation checks where relevant, and a constant-time verifier digest comparison.

## 5. Entropy And Encoding

Raw access-token material must contain at least 256 bits of entropy.

Raw device credential material for the first `device_credential_login` posture must also contain at least 256 bits of entropy. It must be an installation credential generated for authentication, not a raw operating-system device identifier, advertising identifier, hardware serial number, account email, player name, or provider subject.

Token and credential text presentation rules:

```yaml
minimum_entropy_bits: 256
first_text_encoding: url_safe_unpadded_base64_or_equivalent
case_sensitive: true
allowed_in_url_query: false
allowed_in_route_name: false
allowed_in_session_metadata: false
allowed_in_logs: false
```

Future text encoding must avoid whitespace, control characters, path separators, query delimiters, and visually ambiguous formatting. The encoding is not a claim container and must not embed player, credential, route, permission, timestamp, or account lifecycle data.

## 6. Key Identifier And Secret Configuration

`verifier_key_id` identifies the server-side verifier key set used for a stored record. It is not the secret key value.

Classification:

```yaml
verifier_key_id:
  secret_value: false
  public_api_field: false
  log_safe_by_default: false
  allowed_in_database_record: true
  allowed_in_internal_rotation_plan: true
  allowed_in_public_error: false
  allowed_in_conversation_log: false
  allowed_in_change_spec_examples: false
```

Rules:

- Server-side verifier keys, peppers, and raw secret configuration values are secret material.
- `verifier_key_id` must not contain the key value, environment variable value, credential, token, cloud secret path that reveals tenancy, or provider secret.
- Future key loading, environment variables, rotation, fallback, and operational secret storage require a separate implementation/configuration gate.
- KMS or external secret-management integration is not required for the first local posture. Adding one later requires a dependency adoption record and an operations boundary.

## 7. Constant-Time Comparison

Future verifier comparison must compare verifier digests using constant-time equality.

Acceptable future Go primitives:

```yaml
preferred_go_comparison: crypto/hmac.Equal
acceptable_go_comparison: crypto/subtle.ConstantTimeCompare
plaintext_comparison: forbidden
bytes_equal_for_verifier_digest: forbidden
string_equal_for_verifier_digest: forbidden
```

Rules:

- Compute verifier digests before comparison.
- Compare verifier digest bytes, not raw credential or token strings.
- Do not compare verifier digests with `==`, `bytes.Equal`, string conversion, map lookup, or database-only equality.
- A missing lookup record must produce the same public failure class as an invalid verifier.
- Public failures must not disclose whether lookup digest, verifier digest, player account, credential record, token record, key id, algorithm version, or expiration state caused the failure unless the corresponding semantic error class explicitly allows it.

Timing equalization for missing records, disabled accounts, expired tokens, and revoked tokens may require later implementation tests. This standard requires constant-time digest comparison for verifier material, but it does not implement validation behavior.

## 8. Redaction Requirements

Forbidden in logs, traces, metrics labels, public errors, panic output, audit-safe facts, test snapshots, fixtures, ADRs, change specs, documentation examples, and conversation logs:

- Raw credential proof.
- Raw access-token text.
- Raw device identifiers when used as proof or proof input.
- Credential lookup digest.
- Credential verifier digest.
- Token lookup digest.
- Token verifier digest.
- Server-side HMAC keys, peppers, seed material, and secret configuration values.
- Full `verifier_key_id` values unless a later operations standard declares a specific safe representation.
- Provider secrets, OAuth tokens, JWTs, JWKs, passwords, or password-like credentials.

Allowed with care:

- `credential_record_id`
- `token_record_id`
- `player_id` when the surrounding artifact already allows player identifiers
- lifecycle state names
- registered error codes
- non-secret reason catalog values
- a future short redacted fingerprint only after a fingerprint standard exists

Digest values are not safe just because they are not raw tokens. They remain secret-adjacent or verifier material and must stay out of public artifacts.

## 9. Redaction Test Expectations

Future implementation work that handles credential proof, access-token proof, verifier inputs, verifier digests, or application authentication errors must add focused redaction tests.

Minimum test expectations:

```yaml
redaction_tests:
  raw_credential_absent_from_public_error: required
  raw_access_token_absent_from_public_error: required
  raw_credential_absent_from_logs: required_when_logging_path_exists
  raw_access_token_absent_from_logs: required_when_logging_path_exists
  raw_credential_absent_from_audit_safe_facts: required_when_audit_fact_path_exists
  raw_access_token_absent_from_audit_safe_facts: required_when_audit_fact_path_exists
  verifier_digest_absent_from_public_error: required
  verifier_key_value_absent_from_all_outputs: required
  key_identifier_absent_from_public_error: required
  registered_error_code_present: required
```

Test fixtures should use obvious sentinel values such as:

```text
vibit-test-raw-credential-do-not-log
vibit-test-raw-access-token-do-not-log
vibit-test-verifier-key-do-not-log
```

These sentinel values are synthetic test material only. They must not become examples for production credential or token shape.

## 10. Dependency Posture

The first high-entropy verifier posture does not require an external dependency.

Accepted for later implementation after an explicit code gate:

```yaml
go_standard_library:
  hmac: crypto/hmac
  sha256: crypto/sha256
  constant_time: crypto/hmac.Equal_or_crypto/subtle
  randomness: crypto/rand
  base64url: encoding/base64
```

Deferred and not required for the first posture:

- bcrypt.
- Argon2 or Argon2id.
- JWT, JWK, or signing libraries.
- OAuth or OIDC provider SDKs.
- KMS or cloud secret-manager SDKs.
- Redis-like token/session stores.
- Password policy libraries.

If a future login method accepts password-like or low-entropy human input, it must define a separate credential boundary and dependency adoption record. HMAC-SHA-256 over low-entropy passwords is not ratified by this standard.

## 11. Ownership

Future verifier algorithm code is application-owned:

```text
runtime/internal/app
```

The code may later live in an application-owned child package such as:

```text
runtime/internal/app/authentication
```

Ownership rules:

- Application authentication service code owns token generation orchestration and verifier comparison after later gates authorize them.
- `authentication.Repository` stores and retrieves already-computed digest material; it does not compute or compare verifiers.
- PostgreSQL adapters persist digest bytes and key identifiers; they do not make authentication decisions.
- Protobuf adapters and WebSocket transports carry already-decoded proof fields after later protocol gates; they do not compute digests or compare verifiers.
- Generated authentication contract shapes remain metadata-only and immutable.

## 12. Non-Goals

This standard does not:

- Add token generation.
- Add credential generation.
- Add verifier comparison.
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
- Add secret configuration loading.
- Add KMS, OAuth, OIDC, JWT, bcrypt, or Argon2 dependencies.

## 13. Verification Path

The repository check rule for this boundary is:

```text
runtime.token_credential_verifier_algorithm_redaction_boundary
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

No runtime Go authentication behavior is verified by this standard because no behavior is added.

## 14. Follow-Up Gates

Recommended follow-up gates:

- Secret configuration and verifier key loading boundary.
- Token and credential material generation implementation gate.
- Verifier digest comparison implementation gate.
- Application authentication service implementation gate.
- Protobuf authentication message gate.
- WebSocket request proof carrier gate.
- Authentication redaction test implementation gate.
