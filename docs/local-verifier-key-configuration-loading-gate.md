# Local Verifier Key Configuration Loading Gate

Status: Draft v0.1
Last updated: 2026-05-15
Scope: First local verifier key configuration loading implementation gate, explicit input posture, validation rules, redaction rules, package boundaries, test requirements, and deferrals before secret-loading code is added
Depends on: `docs/secret-configuration-verifier-key-loading-boundary.md`, `docs/verifier-digest-computation-comparison-boundary.md`, `docs/authentication-service-implementation-readiness-gate.md`
Canonical decision: `ADR-0045`

The paired Simplified Chinese translation is `docs/local-verifier-key-configuration-loading-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This gate defines the first bounded implementation slice for local verifier key configuration loading before adding code.

It exists because key configuration is security-sensitive, but the first implementation still needs to stay small enough for agents to implement and verify. The future implementation should first prove that a complete in-memory verifier key set can be validated, redacted, and tested without coupling the core validator to process environment parsing, KMS providers, cloud secret managers, protocol carriers, repositories, or login behavior.

This is an implementation-gate standard. It does not add Go code, imports, secret loading, environment parsing, token generation, credential generation, verifier digest computation, verifier comparison, authentication service behavior, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, repository methods, SQL migrations, external secret-manager integrations, authentication dependencies, or production authentication behavior.

## 2. Core Rule

The first local verifier key configuration loading gate is:

```yaml
local_verifier_key_configuration_loading_gate: defined
implementation_authorized_by_this_standard: false
planned_owner: runtime/internal/app
planned_package: runtime/internal/app/authentication
first_implementation_slice: explicit_in_memory_verifier_key_set_validation
first_environment_loader_status: deferred_to_follow_up_gate
process_environment_parsing_status: deferred
base64_text_decoding_status: deferred
external_kms_secret_manager_status: deferred
token_generation_status: deferred
credential_generation_status: deferred
digest_helper_status: deferred
verifier_comparison_status: deferred
authentication_service_behavior_status: deferred
login_execution_status: deferred
token_validation_status: deferred
protobuf_authentication_messages_status: deferred
websocket_proof_carrier_status: deferred
repository_interface_change_status: deferred
migration_schema_change_status: deferred
authentication_dependencies_status: deferred
```

The first implementation should validate explicit in-memory input only. Environment variable loading is deliberately sequenced after that validator exists.

The reason is simple: the validator is the invariant-bearing core, while environment parsing is an adapter. Agents should be able to test the invariant-bearing core without a process environment, local files, shell quoting, or deployment assumptions.

## 3. Package And File Boundary

Future implementation ownership:

```text
runtime/internal/app/authentication
```

Allowed first-slice files after a later implementation work item authorizes code:

```text
runtime/internal/app/authentication/verifier_key_config.go
runtime/internal/app/authentication/verifier_key_config_test.go
```

Allowed helper file only if the first implementation remains small:

```text
runtime/internal/app/authentication/errors.go
```

Forbidden first-slice write areas:

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

The first implementation must not wire configuration into process startup. It should define a small application-owned validation primitive that later loaders can call.

## 4. First Input Shape

The first code slice should accept explicit in-memory data with already-decoded bytes:

```yaml
VerifierKeySetConfig:
  key_set_id: string
  credential_lookup_key: []byte
  credential_verifier_key: []byte
  token_lookup_key: []byte
  token_verifier_key: []byte
```

The accepted output should be an immutable validated key-set value:

```yaml
VerifierKeySet:
  key_set_id: internal_non_log_safe_identifier
  credential_lookup_key: private_bytes
  credential_verifier_key: private_bytes
  token_lookup_key: private_bytes
  token_verifier_key: private_bytes
```

Rules:

- Input byte slices must be copied before storage.
- Output accessors must not expose mutable internal slices.
- The key set id is required but not log-safe by default.
- The key set id is not a key, not a credential, not a token, and not a public API value.
- The value must not implement `fmt.Stringer` in a way that exposes secrets.
- Error values must not include key bytes, encoded key values, full concrete key ids, environment variable values, or deployment identifiers.

## 5. Validation Rules

The first validator must reject:

```yaml
validation_failures:
  missing_key_set_id: fail_closed
  missing_credential_lookup_key: fail_closed
  missing_credential_verifier_key: fail_closed
  missing_token_lookup_key: fail_closed
  missing_token_verifier_key: fail_closed
  decoded_key_shorter_than_32_bytes: fail_closed
  duplicate_logical_key_bytes: fail_closed
  all_zero_key_bytes: fail_closed
  obvious_repeated_single_byte_key: fail_closed
```

Minimum key length:

```yaml
decoded_key_minimum_length_bytes: 32
minimum_key_entropy_bits: 256
```

The first validator cannot prove entropy for caller-supplied bytes, but it must reject common invalid shapes that indicate misuse. Production key generation requirements remain defined by the secret configuration boundary.

## 6. Error And Redaction Requirements

The future first implementation should expose typed or sentinel errors that are useful to tests without disclosing secret material.

Allowed error classification examples:

```yaml
error_classes:
  missing_key_set_id
  missing_required_key
  key_too_short
  duplicate_logical_key
  weak_repeated_key
```

Forbidden in errors, logs, test snapshots, docs, ADRs, change specs, and conversation logs:

- Verifier key bytes.
- Encoded verifier key values.
- Environment variable values.
- Full concrete `verifier_key_id` values.
- Credentials.
- Access tokens.
- Lookup digests.
- Verifier digests.
- HMAC inputs.
- HMAC outputs.

Tests may use obvious sentinel strings or short byte slices to prove redaction, but they must not commit production-like keys.

## 7. Environment Loading Sequence

Environment loading is not part of the first implementation slice.

The later environment loader should call the explicit validator instead of duplicating validation rules. That later gate may authorize:

```yaml
future_environment_loader:
  variable_names:
    - VIBIT_AUTH_VERIFIER_KEY_SET_ID
    - VIBIT_AUTH_CREDENTIAL_LOOKUP_KEY
    - VIBIT_AUTH_CREDENTIAL_VERIFIER_KEY
    - VIBIT_AUTH_TOKEN_LOOKUP_KEY
    - VIBIT_AUTH_TOKEN_VERIFIER_KEY
  decoding: base64url_unpadded_or_standard_base64_to_bytes
  input_source: process_environment
  redaction_required: true
```

This gate does not authorize `os.Getenv`, `os.LookupEnv`, `.env` parsing, local file reading, CLI flags, startup wiring, or environment-driven runtime authentication.

## 8. Required Tests For The Future First Code Slice

The future implementation must add focused unit tests under:

```text
runtime/internal/app/authentication/verifier_key_config_test.go
```

Minimum test cases:

```yaml
required_tests:
  accepts_complete_distinct_32_byte_key_set
  copies_input_key_material
  accessors_do_not_expose_mutable_internal_slices
  missing_key_set_id_fails_closed
  missing_each_logical_key_fails_closed
  short_each_logical_key_fails_closed
  duplicate_logical_keys_fail_closed
  all_zero_key_fails_closed
  repeated_single_byte_key_fails_closed
  errors_do_not_include_secret_bytes
  errors_do_not_include_full_key_set_id
```

Tests must not require PostgreSQL, MinIO, WebSocket transport, Protobuf generation, process environment variables, KMS, cloud SDKs, or external services.

## 9. Dependency Posture

No new external dependency is allowed by this gate.

The future first implementation should use only the Go standard library. A major external KMS, cloud secret-manager, operations, cryptography, password-hashing, JWT, OAuth, OIDC, provider, or dotenv dependency requires a separate adoption record and operations boundary.

## 10. Nakama And Pitaya Mapping

Nakama capability reference:

- Account authentication and session token validation depend on trustworthy server-side secret material.
- Secret handling should remain server-owned and redacted.

Pitaya capability reference:

- Handler identity context should receive validated identity after proof validation.
- Transport/session state should not own verifier key material.

vibit adaptation:

- Verifier key configuration is application-owned.
- Key validation precedes digest helpers and service behavior.
- Transport, protocol, repositories, generated code, and domain modules do not parse or hold verifier keys.

## 11. Deferrals

This gate does not authorize:

- Go code by itself.
- Environment variable parsing.
- Base64 text decoding.
- `.env` parsing.
- Startup wiring.
- Token generation.
- Credential generation.
- Verifier digest computation.
- Verifier comparison.
- Authentication service behavior.
- Login execution.
- Access-token validation.
- Logout execution.
- Refresh execution.
- Cleanup jobs.
- Protobuf authentication messages.
- WebSocket proof carriers.
- WebSocket handshake authentication.
- Repository interface changes.
- PostgreSQL adapter changes.
- SQL migration changes.
- External KMS or secret-manager integrations.
- Authentication dependencies.
- Production authentication behavior.

## 12. Verification

The repository check rule for this gate is:

```text
runtime.local_verifier_key_configuration_loading_gate
```

The check should verify:

- This standard, translation, and ADR exist.
- Architecture manifests and agent guides reference this gate.
- Required markers identify explicit in-memory validation as the first implementation slice.
- Runtime code has not implemented environment parsing, secret loading, token generation, digest helpers, verifier comparison, authentication behavior, protocol carriers, repository changes, migrations, KMS, secret-manager integration, or new authentication dependencies from this gate.

Future code verification must include `go test ./...` from `runtime/` and `node tools/vibit check runtime --json`.
