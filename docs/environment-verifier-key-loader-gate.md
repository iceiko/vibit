# Environment Verifier Key Loader Gate

Status: Draft v0.1
Last updated: 2026-05-15
Scope: Process environment verifier key loader gate, environment variable contract, decoding posture, validation handoff, redaction rules, file boundaries, tests, and deferrals before loader code is added
Depends on: `docs/secret-configuration-verifier-key-loading-boundary.md`, `docs/local-verifier-key-configuration-loading-gate.md`
Canonical decision: `ADR-0046`

The paired Simplified Chinese translation is `docs/environment-verifier-key-loader-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This gate defines the future process environment verifier key loader before adding loader code.

The explicit in-memory verifier key set validator now exists under `runtime/internal/app/authentication`. The next risk is letting a future agent duplicate validation rules while adding environment parsing, or leak verifier key values through decoder errors, logs, test snapshots, change specs, or conversation logs.

This is an implementation-gate standard. It does not add Go code, imports, process environment parsing, Base64 text decoding, local secret file reading, `.env` behavior, CLI flag behavior, startup wiring, token generation, credential generation, verifier digest computation, verifier comparison, authentication service behavior, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, repository methods, SQL migrations, external secret-manager integrations, authentication dependencies, or production authentication behavior.

## 2. Core Rule

The environment verifier key loader gate is:

```yaml
environment_verifier_key_loader_gate: defined
implementation_authorized_by_this_standard: false
planned_owner: runtime/internal/app
planned_package: runtime/internal/app/authentication
future_loader_slice: process_environment_verifier_key_loader
required_validator_handoff: NewVerifierKeySet
explicit_in_memory_validator_required: true
explicit_in_memory_validator_source: runtime/internal/app/authentication/verifier_key_config.go
future_loader_source: runtime/internal/app/authentication/verifier_key_env.go
future_loader_tests: runtime/internal/app/authentication/verifier_key_env_test.go
input_source: process_environment
environment_variable_contract_declared: true
base64_text_decoding_posture_declared: true
startup_wiring_status: deferred
local_secret_file_status: deferred
dotenv_status: deferred
external_kms_secret_manager_status: deferred
token_generation_status: deferred
credential_generation_status: deferred
digest_helper_status: deferred
verifier_comparison_status: deferred
authentication_service_behavior_status: deferred
protobuf_authentication_messages_status: deferred
websocket_proof_carrier_status: deferred
repository_interface_change_status: deferred
migration_schema_change_status: deferred
authentication_dependencies_status: deferred
```

The future loader must remain an adapter. It may collect text values from a process environment source, decode key text into bytes, and hand the result to the explicit in-memory validator. It must not become a second validator, a startup composition point, a service implementation, or an authentication behavior entry point.

## 3. Environment Variable Contract

The future process environment contract is:

```yaml
environment_variables:
  VIBIT_AUTH_VERIFIER_KEY_SET_ID:
    required: true
    value_kind: key_set_identifier
    decoding: trim_space_string
    log_safe_value: false
  VIBIT_AUTH_CREDENTIAL_LOOKUP_KEY:
    required: true
    value_kind: verifier_key_bytes
    decoding: base64_text_to_bytes
    log_safe_value: false
  VIBIT_AUTH_CREDENTIAL_VERIFIER_KEY:
    required: true
    value_kind: verifier_key_bytes
    decoding: base64_text_to_bytes
    log_safe_value: false
  VIBIT_AUTH_TOKEN_LOOKUP_KEY:
    required: true
    value_kind: verifier_key_bytes
    decoding: base64_text_to_bytes
    log_safe_value: false
  VIBIT_AUTH_TOKEN_VERIFIER_KEY:
    required: true
    value_kind: verifier_key_bytes
    decoding: base64_text_to_bytes
    log_safe_value: false
```

Rules:

- All five variables are required for the first process environment posture.
- Environment variable names are log-safe.
- Environment variable values are not log-safe.
- The full concrete key set id value is not log-safe by default.
- Missing values must fail closed.
- Empty or whitespace-only key set id must fail closed through the explicit validator.
- Empty key text must fail closed before or through the explicit validator.
- Decoded key bytes must be passed to `NewVerifierKeySet`.
- The loader must not weaken validator requirements for missing, short, duplicate, all-zero, or repeated-byte keys.

## 4. Decoding Posture

The future loader may use only the Go standard library for decoding.

Allowed future imports after the implementation work item authorizes code:

```yaml
future_standard_library_imports_allowed:
  - encoding/base64
  - os
  - strings
  - errors
  - fmt
```

Decoding policy:

```yaml
preferred_encoding: base64url_unpadded
compatibility_encoding: standard_base64_padded
invalid_base64_text: fail_closed
raw_unencoded_key_text: forbidden
hex_key_text: forbidden
json_key_blob: forbidden
partial_key_set: forbidden
```

The first future implementation should prefer URL-safe unpadded Base64 because it is practical for environment variables and shell configuration. It may also accept standard padded Base64 for operator ergonomics if tests cover both paths. It must not accept raw key text, hex text, JSON blobs, comma-delimited values, or partial key sets.

## 5. Package And File Boundary

Future implementation ownership:

```text
runtime/internal/app/authentication
```

Allowed future loader files after a later implementation work item authorizes code:

```text
runtime/internal/app/authentication/verifier_key_env.go
runtime/internal/app/authentication/verifier_key_env_test.go
```

The future loader may use the existing validator in:

```text
runtime/internal/app/authentication/verifier_key_config.go
```

Forbidden write areas for the environment loader implementation slice:

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

The future implementation must not wire the loader into process startup. Startup wiring is a separate composition decision.

## 6. Future Loader Shape

The preferred future shape is:

```yaml
EnvLookup:
  type: function
  signature: "func(name string) (string, bool)"

future_functions:
  LoadVerifierKeySetFromEnvironment:
    input: EnvLookup
    output: VerifierKeySet
    behavior: decode_required_environment_values_then_call_NewVerifierKeySet
  LoadVerifierKeySetFromProcessEnvironment:
    input: none
    output: VerifierKeySet
    behavior: small_os_LookupEnv_adapter_only
```

The testable loader should accept an explicit lookup function so tests can avoid mutating global process environment when that is simpler. A tiny process adapter may call `os.LookupEnv` after the implementation gate authorizes it, but it must still not wire anything into server startup.

## 7. Error And Redaction Requirements

The future loader should expose typed or sentinel errors that are useful to tests without disclosing secret material.

Allowed error classification examples:

```yaml
error_classes:
  missing_environment_variable
  invalid_environment_key_encoding
  invalid_environment_key_set
```

Allowed in public error text:

- Environment variable names.
- Error classes.
- Logical key purposes.

Forbidden in errors, logs, test snapshots, docs, ADRs, change specs, and conversation logs:

- Environment variable values.
- Encoded verifier key values.
- Decoded verifier key bytes.
- Full concrete key set id values.
- Credentials.
- Access tokens.
- Lookup digests.
- Verifier digests.
- HMAC inputs.
- HMAC outputs.
- Deployment-specific identifiers.

The future loader must wrap or map validation errors without adding secret values to their text.

## 8. Required Tests For The Future Loader

The future implementation must add focused unit tests under:

```text
runtime/internal/app/authentication/verifier_key_env_test.go
```

Minimum test cases:

```yaml
required_tests:
  accepts_complete_base64url_unpadded_environment_key_set
  accepts_complete_standard_base64_padded_environment_key_set
  missing_each_environment_variable_fails_closed
  invalid_each_encoded_key_fails_closed
  decoded_short_key_fails_through_validator
  duplicate_decoded_keys_fail_through_validator
  all_zero_decoded_key_fails_through_validator
  repeated_single_byte_decoded_key_fails_through_validator
  loader_calls_explicit_in_memory_validator
  errors_include_environment_variable_name_when_safe
  errors_do_not_include_environment_variable_values
  errors_do_not_include_full_key_set_id
  process_environment_adapter_is_small_and_unwired
```

Tests must not require PostgreSQL, MinIO, WebSocket transport, Protobuf generation, KMS, cloud SDKs, local files, `.env` files, or external services.

## 9. Dependency Posture

No new external dependency is allowed by this gate.

The future loader should use only the Go standard library. A major external KMS, cloud secret-manager, operations, cryptography, password-hashing, JWT, OAuth, OIDC, provider, or dotenv dependency requires a separate adoption record and operations boundary.

## 10. Nakama And Pitaya Mapping

Nakama capability reference:

- Server-side authentication secret material must be configured reliably before authentication token behavior can be production meaningful.
- Session token validation depends on stable server-owned verifier material.

Pitaya capability reference:

- Realtime route handlers should receive identity context after proof validation.
- Server/session infrastructure should not own verifier key material.

vibit adaptation:

- Environment configuration is an application-owned adapter.
- The invariant-bearing verifier key set validator remains the source of validation truth.
- Transport, protocol, repositories, generated code, and domain modules do not parse or hold verifier keys.
- Startup composition remains separate from loader implementation.

## 11. Deferrals

This gate does not authorize:

- Go code by itself.
- Process environment parsing in this change.
- Base64 text decoding in this change.
- Local secret file reading.
- `.env` parsing.
- CLI flag secret input.
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

This gate's repository check rule is:

```yaml
runtime.environment_verifier_key_loader_gate
```

Verification for this gate must prove:

- The English standard and Simplified Chinese translation exist.
- ADR-0046 exists.
- Architecture manifests and agent guides reference the gate.
- The environment variable contract is declared.
- The future loader must hand off to `NewVerifierKeySet`.
- No process environment parsing or Base64 decoding code is added in this gate.
- Deferrals remain visible in manifests.

## 13. Migration Path

1. Complete this gate as standards and manifest work only.
2. Implement the environment verifier key loader in a later code slice.
3. Keep startup wiring behind a separate composition gate.
4. Keep KMS, cloud secret-manager, `.env`, local secret files, CLI secret input, and production operations posture behind separate decisions.
