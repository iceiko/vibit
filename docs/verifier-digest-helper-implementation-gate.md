# Verifier Digest Helper Implementation Gate

Status: Draft v0.1
Last updated: 2026-05-15
Scope: Verifier digest helper implementation gate, future helper file boundaries, canonical input construction, verifier key handoff, digest class mapping, redaction rules, tests, and deferrals before digest code is added
Depends on: `docs/verifier-digest-computation-comparison-boundary.md`, `docs/token-credential-verifier-algorithm-redaction-boundary.md`, `docs/token-credential-material-generation-implementation-gate.md`, `docs/authentication-service-implementation-readiness-gate.md`
Canonical decision: `ADR-0048`

The paired Simplified Chinese translation is `docs/verifier-digest-helper-implementation-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This gate defines the next bounded implementation slice for lookup digest and verifier digest computation helpers.

The repository already defines the verifier digest computation and comparison boundary, the verifier algorithm and redaction boundary, and the material generation helpers. The next risk is letting a future agent compute digests in service orchestration, repository adapters, transport handlers, test fixtures, or protocol code, or mixing digest computation with verifier comparison and authentication behavior.

This is an implementation-gate standard. It does not add Go code, imports, HMAC computation, digest helpers, verifier comparison, authentication service behavior, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, repository methods, SQL migrations, startup wiring, authentication dependencies, external cryptography services, KMS, cloud secret-manager integrations, or production authentication behavior.

## 2. Core Rule

The verifier digest helper implementation gate is:

```yaml
verifier_digest_helper_implementation_gate: defined
implementation_authorized_by_this_standard: false
future_implementation_work_item: W-0103
planned_owner: runtime/internal/app
planned_package: runtime/internal/app/authentication
future_slice: lookup_and_verifier_digest_computation_helpers
future_source: runtime/internal/app/authentication/verifier_digest.go
future_tests: runtime/internal/app/authentication/verifier_digest_test.go
verifier_algorithm_family: vibit_hmac_sha256_v1
verifier_version: 1
digest_size_bytes: 32
canonical_input_version: vibit.auth.verifier.input.v1
canonical_input_encoding: ascii_header_length_prefixed_purpose_label_length_prefixed_raw_material
hmac_hash: crypto/sha256
constant_time_comparison_primitive: crypto/hmac.Equal
verifier_key_handoff: VerifierKeySet
raw_material_handoff: GeneratedSecretMaterial_or_raw_bytes
digest_output_shape: raw_32_byte_digest
digest_copying_required: true
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

The future implementation must be a helper-only slice. It may build canonical digest input, compute HMAC-SHA-256 lookup digests, compute HMAC-SHA-256 verifier digests, and return copied digest byte slices. It must not compare verifier digests, write repositories, choose accounts, issue login responses, validate tokens, parse bearer proofs, or touch protocol carriers.

## 3. Future Helper Shape

Future implementation ownership:

```text
runtime/internal/app/authentication
```

Allowed future files after the implementation work item authorizes code:

```text
runtime/internal/app/authentication/verifier_digest.go
runtime/internal/app/authentication/verifier_digest_test.go
```

Preferred future API shape:

```yaml
future_types:
  DigestClass:
    values:
      - credential_lookup
      - credential_verifier
      - token_lookup
      - token_verifier
  ComputedDigest:
    owns: copied_digest_bytes_and_class
    methods:
      - Class() DigestClass
      - Bytes() []byte
    constraints:
      - Bytes returns a copy.
      - Error text and string formatting must not expose digest bytes, raw material, or key values.

future_constants:
  CanonicalInputVersion: "vibit.auth.verifier.input.v1"
  PurposeLabelCredentialLookup: "vibit.credential.lookup.v1"
  PurposeLabelCredentialVerifier: "vibit.credential.verifier.v1"
  PurposeLabelTokenLookup: "vibit.access_token.lookup.v1"
  PurposeLabelTokenVerifier: "vibit.access_token.verifier.v1"

future_functions:
  ComputeCredentialLookupDigest:
    signature: "func ComputeCredentialLookupDigest(keySet VerifierKeySet, rawMaterial []byte) (ComputedDigest, error)"
    behavior: build_canonical_input_and_hmac_with_credential_lookup_key
  ComputeCredentialVerifierDigest:
    signature: "func ComputeCredentialVerifierDigest(keySet VerifierKeySet, rawMaterial []byte) (ComputedDigest, error)"
    behavior: build_canonical_input_and_hmac_with_credential_verifier_key
  ComputeTokenLookupDigest:
    signature: "func ComputeTokenLookupDigest(keySet VerifierKeySet, rawMaterial []byte) (ComputedDigest, error)"
    behavior: build_canonical_input_and_hmac_with_token_lookup_key
  ComputeTokenVerifierDigest:
    signature: "func ComputeTokenVerifierDigest(keySet VerifierKeySet, rawMaterial []byte) (ComputedDigest, error)"
    behavior: build_canonical_input_and_hmac_with_token_verifier_key
```

The `VerifierKeySet` handoff is required. The helper must accept an already-validated key set; it must not load keys, parse environment variables, or choose key sets. The raw material handoff is the decoded raw bytes from `GeneratedSecretMaterial.RawBytes()` or equivalent; it must not accept encoded text, player id, session id, route name, or metadata.

## 4. Canonical Input Construction

The future helper must build a deterministic canonical byte input before HMAC computation.

```yaml
canonical_digest_input:
  version_header_ascii: "vibit.auth.verifier.input.v1"
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

- The version header must be literal ASCII bytes. It must change if the canonical input shape changes.
- Purpose labels must be the registered constant for each digest class.
- Raw material must be the decoded raw secret material bytes, not normalized text, metadata, player id, session id, route name, or provider subject.
- The raw material length for the first posture is 32 bytes. The helper must reject raw material with length zero.
- Length prefixes use big-endian unsigned 16-bit integers and make the input unambiguous even if future raw material shape changes.
- Future tests must include deterministic fixture vectors that verify canonical byte construction by comparing the entire input byte sequence.

## 5. Digest Classes And Key Mapping

Each digest class uses a registered purpose label and the matching logical key from the verifier key set:

```yaml
credential_lookup_digest:
  purpose_label: vibit.credential.lookup.v1
  key_accessor: VerifierKeySet.CredentialLookupKey
  output_bytes: 32
credential_verifier_digest:
  purpose_label: vibit.credential.verifier.v1
  key_accessor: VerifierKeySet.CredentialVerifierKey
  output_bytes: 32
token_lookup_digest:
  purpose_label: vibit.access_token.lookup.v1
  key_accessor: VerifierKeySet.TokenLookupKey
  output_bytes: 32
token_verifier_digest:
  purpose_label: vibit.access_token.verifier.v1
  key_accessor: VerifierKeySet.TokenVerifierKey
  output_bytes: 32
```

Rules:

- Lookup digests and verifier digests must use different purpose labels.
- Credential digests and token digests must use different purpose labels.
- Each compute function must use exactly the matching key; passing a mismatched key must not be possible through the function signature.
- Digest bytes must not be truncated for storage, lookup, comparison, or logging.
- The helper must copy digest bytes on return.

## 6. Package And File Boundary

Allowed future implementation area:

```text
runtime/internal/app/authentication/verifier_digest.go
runtime/internal/app/authentication/verifier_digest_test.go
```

Forbidden write areas for the future verifier digest helper slice:

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

The future implementation must not wire digest computation into login execution, account creation, token issuance, token validation, startup, WebSocket transport, Protobuf protocol, PostgreSQL persistence, generated contract output, migrations, or domain repository code.

## 7. Error And Redaction Requirements

The future helper should expose typed or sentinel errors that are useful to tests without disclosing digest material or key values.

Allowed error classification examples:

```yaml
error_classes:
  missing_key_set
  missing_raw_material
  invalid_digest_computation
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
- Verifier key values.
- Encoded verifier key values.
- Full concrete `verifier_key_id` values.
- Candidate key-set counts.

Digest material is not safe just because it is deterministic. Digest bytes remain secret-adjacent and must stay out of public artifacts.

## 8. Required Tests For The Future Helper

Future implementation must add focused unit tests under:

```text
runtime/internal/app/authentication/verifier_digest_test.go
```

Minimum test cases:

```yaml
required_tests:
  canonical_input_is_deterministic
  canonical_input_uses_version_header
  canonical_input_null_separator_present
  canonical_input_length_prefixes_purpose_label
  canonical_input_length_prefixes_raw_material
  lookup_and_verifier_purpose_labels_differ
  credential_and_token_purpose_labels_differ
  digest_output_is_32_bytes
  credential_lookup_digest_uses_credential_lookup_key
  credential_verifier_digest_uses_credential_verifier_key
  token_lookup_digest_uses_token_lookup_key
  token_verifier_digest_uses_token_verifier_key
  different_keys_produce_different_digests
  different_raw_material_produces_different_digests
  digest_bytes_are_copied_on_return
  empty_raw_material_fails_closed
  errors_do_not_include_digest_or_key_material
  helper_does_not_compare_verifiers
```

Tests may use deterministic synthetic keys and raw material only inside tests. They must not become production defaults, documentation examples of real secret shape, or committed production-like secret values.

## 9. Dependency Posture

No new external dependency is allowed by this gate.

Allowed Go standard library packages after the future implementation work item authorizes code:

```yaml
future_standard_library_imports_allowed:
  - crypto/hmac
  - crypto/sha256
  - encoding/binary
  - errors
  - fmt
```

The first helper implementation must not add JWT, JWK, OAuth, OIDC, provider SDKs, password-hashing dependencies, Redis-like stores, KMS SDKs, cloud secret-manager SDKs, operations libraries, or external cryptography services.

## 10. Nakama And Pitaya Mapping

Nakama capability reference:

- Server-side authentication commonly needs HMAC-based credential and token verification. vibit adopts the capability need, not the exact implementation shape.

Pitaya capability reference:

- Realtime route handlers should receive identity context after framework/application validation. vibit keeps digest computation in application-owned helpers, not in transport acceptors or route dispatch.

This gate maps those references into a narrow helper slice: digest computation first, verifier comparison and authentication behavior later.

## 11. Non-Goals

This gate does not:

- Add verifier digest computation code.
- Add verifier comparison code.
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

## 12. Verification Path

The repository check rule for this gate is:

```text
runtime.verifier_digest_helper_implementation_gate
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

No runtime Go verifier digest computation behavior is verified by this gate because no computation behavior is added.

## 13. Follow-Up Gates

Recommended follow-up gates:

- Implement verifier digest computation helpers.
- Implement verifier digest comparison helpers.
- Implement application authentication service behavior.
- Add Protobuf authentication messages.
- Add WebSocket request proof carriers.
