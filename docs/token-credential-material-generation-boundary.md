# Token And Credential Material Generation Boundary

Status: Draft v0.1
Last updated: 2026-05-15
Scope: Future generation ownership, entropy enforcement, text encoding, one-time presentation, non-storage, repository handoff, redaction, dependency posture, and test expectations for high-entropy device credential and opaque access-token material
Depends on: `docs/token-credential-verifier-algorithm-redaction-boundary.md`, `docs/secret-configuration-verifier-key-loading-boundary.md`, `docs/token-lifecycle-storage-implications.md`, `docs/first-login-method-set.md`
Canonical decision: `ADR-0042`

The paired Simplified Chinese translation is `docs/token-credential-material-generation-boundary.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This standard defines how future application-owned authentication code may generate raw credential material and raw access-token material.

It exists before token generation code, credential generation code, secret loading, verifier digest computation, verifier comparison, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository interface changes, migration schema changes, or production authentication behavior.

This is a boundary-only standard. It does not add Go code, imports, runtime services, token generation, credential generation, verifier digest computation, verifier comparison, secret loading, repository methods, SQL migrations, Protobuf messages, WebSocket carriers, routes, or production authentication behavior.

## 2. Core Rule

The first material generation posture is:

```yaml
material_generation_boundary: defined
implementation_authorized_by_this_standard: false
planned_owner: runtime/internal/app
planned_optional_child_package: runtime/internal/app/authentication
first_device_credential_source: server_issued_application_generated
first_access_token_source: server_issued_application_generated
minimum_raw_device_credential_entropy_bits: 256
minimum_raw_access_token_entropy_bits: 256
raw_material_size_bytes: 32
text_encoding: url_safe_unpadded_base64_or_equivalent
one_time_client_visible_presentation_required: true
raw_credential_storage: forbidden
raw_token_storage: forbidden
raw_material_in_repository: forbidden
raw_material_in_transport_logs: forbidden
external_randomness_dependency_required_for_first_posture: false
future_allowed_standard_library_packages:
  - crypto/rand
  - encoding/base64
```

Future first-posture generation helpers may use Go standard library `crypto/rand` and `encoding/base64` after a later code gate authorizes implementation. No external randomness, cryptography, JWT, OAuth, OIDC, provider, KMS, cloud secret-manager, or operations dependency is required for this generation posture.

## 3. Ownership

Future material generation is application-owned:

```text
runtime/internal/app
```

Future helper code may live in an application-owned child package such as:

```text
runtime/internal/app/authentication
```

Ownership rules:

- Application authentication code owns raw token and credential material generation after a later implementation gate.
- `authentication.Repository` never generates, accepts for storage, or returns raw credential or raw token material.
- PostgreSQL adapters never generate, encode, log, persist, or return raw credential or raw token material.
- Protobuf adapters and WebSocket transports may carry already-decoded proof fields after later protocol gates, but they do not generate secrets.
- Generated authentication contract shapes remain metadata-only and immutable.

Material generation must not be placed in transport handlers, protocol adapters, domain modules, repositories, generated output, migrations, SQL fixtures, or tests that are not explicitly scoped to future generation helpers.

## 4. Device Credential Material

The first device credential posture is server-issued and application-generated.

```yaml
device_credential_material:
  source: server_issued_application_generated
  credential_kind: high_entropy_installation_credential
  minimum_entropy_bits: 256
  raw_size_bytes: 32
  text_encoding: url_safe_unpadded_base64_or_equivalent
  one_time_client_visible_presentation: required
  raw_server_storage: forbidden
  public_metadata_source: forbidden
```

Rules:

- The credential is secret proof material, not a raw operating-system device ID, advertising ID, hardware serial number, account email, player name, provider subject, player id, session id, connection id, or metadata-only value.
- The server may present the credential to the client only through a future explicitly authorized authentication response or credential bootstrap response.
- The client is responsible for storing the one-time presented credential after that future response.
- Future server code must compute the credential lookup and verifier digests before repository storage, but exact digest helper implementation remains a later gate.
- Future code must store only digest material, metadata, lifecycle state, timestamps, and `verifier_key_id`.
- The first posture does not accept client-generated installation credentials. A later boundary may add client-generated credential enrollment if it defines entropy, replay, collision, proof-carrier, and abuse controls.

This standard does not define the account creation policy, credential rotation command, recovery flow, or Protobuf response shape.

## 5. Access Token Material

The first access-token posture is server-issued and application-generated.

```yaml
access_token_material:
  source: server_issued_application_generated
  token_format: opaque_high_entropy_token
  minimum_entropy_bits: 256
  raw_size_bytes: 32
  text_encoding: url_safe_unpadded_base64_or_equivalent
  one_time_client_visible_presentation: required
  raw_server_storage: forbidden
  claim_container: false
```

Rules:

- Access tokens are bearer secrets.
- Token text must not embed player id, credential id, token record id, session id, route name, timestamp, permission, account state, provider subject, or claims.
- Future server code must compute token lookup and verifier digests before repository storage, but exact digest helper implementation remains a later gate.
- Future code must store only digest material, token metadata, lifecycle state, timestamps, and `verifier_key_id`.
- Token text may be returned to the client only once in a future explicitly authorized authentication response.
- Token text must not be accepted in URL query parameters, route names, current Protobuf `Session` metadata fields, logs, traces, metrics labels, public errors, audit-safe facts, conversation logs, or change specs.

This standard does not define login execution, token issuance command behavior, token validation, logout, cleanup, refresh, Protobuf messages, or WebSocket proof carriers.

## 6. Encoding And Byte Shape

Future first-posture generation must produce canonical raw bytes before text encoding.

```yaml
raw_material:
  source: cryptographically_secure_random_bytes
  bytes: 32
  entropy_bits: at_least_256
  all_zero_value: forbidden
  repeated_pattern_value: forbidden
  human_readable_password_like_value: forbidden
text_material:
  encoding: url_safe_unpadded_base64_or_equivalent
  whitespace: forbidden
  control_characters: forbidden
  path_separators: forbidden
  query_delimiters: forbidden
  case_sensitive: true
```

Rules:

- Encoding is presentation only. It must not become a claim container.
- Decoding must recover the original generated bytes without lossy normalization.
- Future tests should cover entropy length, text alphabet, no padding if base64url unpadded is selected, and round-trip decoding.
- Future code must fail closed if a generation helper cannot obtain cryptographically secure random bytes.

## 7. Repository Handoff

Future application authentication code must keep the repository handoff non-raw.

Allowed repository inputs after later gates:

- `credential_lookup_digest`
- `credential_verifier_digest`
- `token_lookup_digest`
- `token_verifier_digest`
- `verifier_algorithm`
- `verifier_version`
- `verifier_key_id`
- lifecycle state
- timestamps
- stable record identifiers

Forbidden repository inputs:

- Raw credential text.
- Raw credential bytes.
- Raw access-token text.
- Raw access-token bytes.
- Token prefixes.
- Credential prefixes.
- Randomness seeds.
- Verifier keys.
- Encoded verifier keys.
- Secret-manager payloads.

The repository may persist already-computed digest bytes. It must not generate raw material, encode raw material, compute verifier digests, compare verifiers, or decide authentication outcomes.

## 8. Redaction Requirements

Forbidden in logs, traces, metrics labels, public errors, panic output, audit-safe facts, test snapshots, fixtures, ADRs, change specs, documentation examples, and conversation logs:

- Raw credential text.
- Raw credential bytes.
- Raw access-token text.
- Raw access-token bytes.
- Token or credential prefixes.
- Randomness seeds.
- Generated material before encoding.
- Encoded generated material.
- Lookup digests.
- Verifier digests.
- Verifier keys.
- Encoded verifier keys.
- Full concrete `verifier_key_id` values.

Allowed with care:

- Non-secret record identifiers.
- Lifecycle state names.
- Registered error codes.
- Placeholder names such as `<raw-access-token>` only when describing redaction rules.

Generated material is not safe just because it is short-lived. One-time presentation means one client-visible delivery, not "safe to log once."

## 9. Test Expectations

Future implementation work that generates raw credential or access-token material must add focused tests.

Minimum expectations:

```yaml
generation_tests:
  random_source_error_fails_closed: required
  raw_material_length_is_32_bytes: required
  encoded_material_round_trips_to_raw_bytes: required
  encoded_material_uses_allowed_text_alphabet: required
  generated_values_are_not_constant: required
  raw_material_absent_from_logs: required_when_logging_path_exists
  raw_material_absent_from_public_errors: required
  raw_material_absent_from_repository_records: required
  token_text_contains_no_claims: required
  credential_text_contains_no_metadata: required
```

Tests may use deterministic fakes only behind a test seam in future implementation code. They must not replace production randomness, leak fixture secrets, or make generated values stable in non-test builds.

## 10. Dependency Posture

The first material generation posture does not require an external dependency.

Accepted for later implementation after an explicit code gate:

```yaml
go_standard_library:
  randomness: crypto/rand
  base64url: encoding/base64
external_dependency_adoption_record_required_for_first_posture: false
```

Deferred and not required for the first posture:

- External randomness services.
- KMS or cloud secret-manager SDKs.
- JWT, JWK, or signing libraries.
- OAuth or OIDC provider SDKs.
- Password-hashing dependencies.
- Redis-like token/session stores.

If a future login method accepts password-like or low-entropy material, it must define a separate credential boundary and dependency adoption record.

## 11. Non-Goals

This standard does not:

- Add token generation.
- Add credential generation.
- Add secret loading.
- Add verifier digest computation.
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
- Add production authentication behavior.

## 12. Verification Path

The repository check rule for this boundary is:

```text
runtime.token_credential_material_generation_boundary
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

No runtime Go token or credential generation behavior is verified by this standard because no behavior is added.

## 13. Follow-Up Gates

Recommended follow-up gates:

- Verifier digest computation and constant-time comparison boundary.
- Application authentication service implementation gate.
- Authentication redaction test implementation gate.
- Protobuf authentication message gate.
- WebSocket request proof carrier gate.
