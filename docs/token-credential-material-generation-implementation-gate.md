# Token And Credential Material Generation Implementation Gate

Status: Draft v0.1
Last updated: 2026-05-15
Scope: Token and credential material generation helper implementation gate, future helper file boundaries, entropy source, encoding posture, redaction rules, tests, and deferrals before generation code is added
Depends on: `docs/token-credential-material-generation-boundary.md`, `docs/verifier-digest-computation-comparison-boundary.md`, `docs/authentication-service-implementation-readiness-gate.md`, `docs/environment-verifier-key-loader-gate.md`
Canonical decision: `ADR-0047`

The paired Simplified Chinese translation is `docs/token-credential-material-generation-implementation-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This gate defines the next bounded implementation slice for raw device credential and opaque access-token material generation helpers.

The repository already defines the material generation boundary. The next risk is letting a future agent generate secrets in service orchestration, repository adapters, transport handlers, test fixtures, or protocol code, or mixing generation with verifier digest computation and authentication behavior.

This is an implementation-gate standard. It does not add Go code, imports, token generation, credential generation, verifier digest computation, verifier comparison, authentication service behavior, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, repository methods, SQL migrations, startup wiring, authentication dependencies, external randomness services, KMS, cloud secret-manager integrations, or production authentication behavior.

## 2. Core Rule

The token and credential material generation implementation gate is:

```yaml
token_credential_material_generation_implementation_gate: defined
implementation_authorized_by_this_standard: false
future_implementation_work_item: W-0101
planned_owner: runtime/internal/app
planned_package: runtime/internal/app/authentication
future_slice: raw_device_credential_and_access_token_material_helpers
future_source: runtime/internal/app/authentication/material_generation.go
future_tests: runtime/internal/app/authentication/material_generation_test.go
production_entropy_source: crypto/rand.Reader
testable_entropy_handoff: io.Reader
random_read_primitive: io.ReadFull
raw_material_size_bytes: 32
minimum_entropy_bits: 256
text_encoding: base64.RawURLEncoding
encoded_text_shape: url_safe_unpadded_base64
encoded_text_length_chars: 43
raw_material_copying_required: true
all_zero_material_fails_closed: true
repeated_single_byte_material_fails_closed: true
one_time_client_visible_presentation_required: true
raw_material_storage: forbidden
raw_material_repository_handoff: forbidden
digest_computation_status: deferred
verifier_comparison_status: deferred
authentication_service_behavior_status: deferred
protobuf_authentication_messages_status: deferred
websocket_proof_carrier_status: deferred
repository_interface_change_status: deferred
migration_schema_change_status: deferred
startup_wiring_status: deferred
authentication_dependencies_status: deferred
production_authentication_behavior_status: deferred
```

The future implementation must be a helper-only slice. It may create random raw bytes, reject malformed or weak generated material from a supplied reader, encode the material for one-time presentation, and return a small value object for later application-owned digest helpers. It must not compute lookup digests, compute verifier digests, compare verifiers, write repositories, choose accounts, issue login responses, validate tokens, parse bearer proofs, or touch protocol carriers.

## 3. Future Helper Shape

Future implementation ownership:

```text
runtime/internal/app/authentication
```

Allowed future files after the implementation work item authorizes code:

```text
runtime/internal/app/authentication/material_generation.go
runtime/internal/app/authentication/material_generation_test.go
```

Preferred future API shape:

```yaml
future_types:
  GeneratedSecretMaterial:
    owns: copied_raw_bytes_and_encoded_text
    methods:
      - Kind() MaterialKind
      - RawBytes() []byte
      - Text() string
    constraints:
      - RawBytes returns a copy.
      - Text returns URL-safe unpadded Base64 presentation text.
      - Error text and string formatting must not expose raw bytes or text.
  MaterialKind:
    values:
      - device_credential
      - access_token

future_functions:
  GenerateDeviceCredentialMaterial:
    signature: "func GenerateDeviceCredentialMaterial(random io.Reader) (GeneratedSecretMaterial, error)"
    behavior: read_32_random_bytes_validate_encode_as_device_credential_material
  GenerateAccessTokenMaterial:
    signature: "func GenerateAccessTokenMaterial(random io.Reader) (GeneratedSecretMaterial, error)"
    behavior: read_32_random_bytes_validate_encode_as_access_token_material
```

The explicit `io.Reader` handoff is required for testability. Future production service code may pass `crypto/rand.Reader` after its own service behavior gate authorizes the call path. The helper implementation itself may import `crypto/rand` only if it keeps the explicit reader seam and does not wire process startup or service behavior.

## 4. Entropy And Encoding

The future helper must generate canonical raw material before encoding.

```yaml
raw_material:
  bytes: 32
  entropy_bits: at_least_256
  production_source: crypto/rand.Reader
  read_primitive: io.ReadFull
  nil_reader: fail_closed
  short_read: fail_closed
  read_error: fail_closed
  all_zero_bytes: fail_closed
  repeated_single_byte_bytes: fail_closed

encoded_material:
  encoding: base64.RawURLEncoding
  alphabet: url_safe_base64
  padding: forbidden
  expected_length_chars_for_32_bytes: 43
  whitespace: forbidden
  control_characters: forbidden
  path_separators: forbidden
  query_delimiters: forbidden
  claims_or_metadata: forbidden
```

Rules:

- Encoding is presentation only. It must not embed player id, credential id, token id, session id, timestamp, route name, permissions, account state, provider subject, or claims.
- Decoding the encoded text must recover the generated raw bytes exactly.
- Device credential material and access-token material share byte shape but must carry distinct `MaterialKind` values.
- The helper must not make generated values stable in non-test builds.
- The helper must not silently retry forever if randomness fails.

## 5. Package And File Boundary

Allowed future implementation area:

```text
runtime/internal/app/authentication/material_generation.go
runtime/internal/app/authentication/material_generation_test.go
```

Forbidden write areas for the future material generation helper slice:

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

The future implementation must not wire generation into login execution, account creation, token issuance, startup, WebSocket transport, Protobuf protocol, PostgreSQL persistence, generated contract output, migrations, or domain repository code.

## 6. Error And Redaction Requirements

The future helper should expose typed or sentinel errors that are useful to tests without disclosing generated material.

Allowed error classification examples:

```yaml
error_classes:
  missing_random_source
  random_read_failed
  invalid_generated_material
```

Allowed in error text:

- Error classes.
- Material kind names.
- Non-secret numeric expectations such as `32` bytes.

Forbidden in errors, logs, test snapshots, docs, ADRs, change specs, and conversation logs:

- Raw device credential text.
- Raw device credential bytes.
- Raw access-token text.
- Raw access-token bytes.
- Encoded generated material.
- Token or credential prefixes.
- Randomness seeds.
- Lookup digests.
- Verifier digests.
- Verifier keys.
- Encoded verifier keys.
- Full concrete `verifier_key_id` values.

Generated material is not safe just because it is short-lived. One-time presentation means exactly one future client-visible delivery path, not "safe to log once."

## 7. Required Tests For The Future Helper

Future implementation must add focused unit tests under:

```text
runtime/internal/app/authentication/material_generation_test.go
```

Minimum test cases:

```yaml
required_tests:
  device_credential_material_uses_32_random_bytes
  access_token_material_uses_32_random_bytes
  encoded_material_is_base64url_unpadded
  encoded_material_length_is_43_characters
  encoded_material_round_trips_to_raw_bytes
  generated_material_kind_is_preserved
  raw_bytes_are_copied_on_return
  nil_random_source_fails_closed
  random_source_error_fails_closed
  short_random_read_fails_closed
  all_zero_generated_material_fails_closed
  repeated_single_byte_generated_material_fails_closed
  generated_values_are_not_constant_with_progressing_source
  errors_do_not_include_raw_or_encoded_material
  helper_does_not_compute_digests_or_compare_verifiers
```

Tests may use deterministic readers only as explicit test seams. They must not introduce committed production-like secrets, stable generated tokens, stable generated credentials, logs containing material, repository fixtures containing material, or protocol fixtures containing material.

## 8. Dependency Posture

No new external dependency is allowed by this gate.

Allowed Go standard library packages after the future implementation work item authorizes code:

```yaml
future_standard_library_imports_allowed:
  - crypto/rand
  - encoding/base64
  - errors
  - fmt
  - io
```

The first helper implementation must not add JWT, JWK, OAuth, OIDC, provider SDKs, password-hashing dependencies, Redis-like stores, KMS SDKs, cloud secret-manager SDKs, operations libraries, or external randomness services.

## 9. Nakama And Pitaya Mapping

Nakama capability reference:

- Server-side account authentication commonly needs server-issued secret material and session token issuance.
- vibit adopts the capability need, not the exact implementation shape.

Pitaya capability reference:

- Realtime route handlers should receive identity context after framework/application validation.
- vibit keeps secret generation in application-owned helpers, not in transport acceptors or route dispatch.

This gate maps those references into a narrow helper slice: raw material creation first, digest helpers and authentication behavior later.

## 10. Non-Goals

This gate does not:

- Add token generation code.
- Add credential generation code.
- Compute lookup digests.
- Compute verifier digests.
- Compare verifiers.
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

## 11. Verification Path

The repository check rule for this gate is:

```text
runtime.token_credential_material_generation_implementation_gate
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

No runtime Go token or credential material generation behavior is verified by this gate because no generation behavior is added.

## 12. Follow-Up Gates

Recommended follow-up gates:

- Implement token and credential material generation helpers.
- Implement verifier digest computation helpers.
- Implement verifier digest comparison helpers.
- Implement application authentication service behavior.
- Add Protobuf authentication messages.
- Add WebSocket request proof carriers.
