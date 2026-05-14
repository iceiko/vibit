# Secret Configuration And Verifier Key Loading Boundary

Status: Draft v0.1
Last updated: 2026-05-15
Scope: Secret configuration ownership, verifier key separation, future key loading posture, key identifier handling, rotation expectations, development/test posture, production failure behavior, and redaction requirements for the first device-credential and opaque access-token posture
Depends on: `docs/token-credential-verifier-algorithm-redaction-boundary.md`, `docs/application-authentication-service-interface-boundary.md`, `docs/runtime-authentication-implementation-boundary.md`
Canonical decision: `ADR-0041`

The paired Simplified Chinese translation is `docs/secret-configuration-verifier-key-loading-boundary.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This standard defines the secret configuration and verifier key loading boundary that future application-owned authentication service code must follow.

It exists before secret loading, token material generation, credential material generation, verifier comparison, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository interface changes, migration schema changes, or production authentication behavior.

This is a boundary-only standard. It does not add Go code, imports, runtime services, environment parsing, secret loading, token generation, credential generation, digest comparison, repository methods, SQL migrations, Protobuf messages, WebSocket carriers, routes, KMS integration, secret-manager integration, or production authentication behavior.

## 2. Core Rule

The first secret configuration posture is:

```yaml
secret_configuration_boundary: defined
implementation_authorized_by_this_standard: false
planned_owner: runtime/internal/app
planned_optional_child_package: runtime/internal/app/authentication
first_local_key_source: process_environment_or_explicit_runtime_secret_input
external_kms_secret_manager_required_for_first_local_posture: false
external_secret_manager_dependency_adoption_record_required: true
verifier_key_set_required: true
verifier_key_id_required_for_stored_records: true
minimum_key_entropy_bits: 256
decoded_key_minimum_length_bytes: 32
key_text_encoding: base64url_unpadded_or_standard_base64_decoded_to_bytes
production_missing_key_behavior: fail_closed
development_default_production_keys: forbidden
automatic_persistent_random_key_generation: forbidden
runtime_secret_values_in_committed_artifacts: forbidden
```

The first local implementation may use process environment configuration or an explicit runtime secret input after a later code gate authorizes implementation. It does not require KMS or a cloud secret-manager dependency.

Production secret storage, external secret managers, KMS providers, cloud provider SDKs, operational rotation systems, or container orchestration secret integrations require a later dependency adoption record and operations boundary before implementation.

## 3. Ownership

Future secret configuration loading is application-owned:

```text
runtime/internal/app
```

Future code may use an application-owned child package such as:

```text
runtime/internal/app/authentication
```

Ownership rules:

- Application authentication code may load, validate, and hold verifier key material after a later implementation gate.
- `authentication.Repository` stores and retrieves already-computed digests, status fields, timestamps, and `verifier_key_id` values only.
- PostgreSQL adapters persist key identifiers and digest bytes; they do not load keys, decode secret configuration, compute HMACs, compare verifiers, rotate keys, or decide authentication outcomes.
- WebSocket transport and Protobuf protocol adapters do not read secret configuration or own key loading.
- Generated authentication contract shapes remain metadata-only and immutable.

Secret configuration must not be moved into transport, protocol, domain modules, generated output, migration files, SQL fixtures, or repository adapters for convenience.

## 4. Logical Key Set

The first verifier posture requires four distinct logical server-side keys:

```yaml
key_set_logical_keys:
  credential_lookup_key:
    purpose_label: vibit.credential.lookup.v1
    digest_class: credential_lookup_digest
    reuse_with_other_digest_classes: forbidden
  credential_verifier_key:
    purpose_label: vibit.credential.verifier.v1
    digest_class: credential_verifier_digest
    reuse_with_other_digest_classes: forbidden
  token_lookup_key:
    purpose_label: vibit.access_token.lookup.v1
    digest_class: token_lookup_digest
    reuse_with_other_digest_classes: forbidden
  token_verifier_key:
    purpose_label: vibit.access_token.verifier.v1
    digest_class: token_verifier_digest
    reuse_with_other_digest_classes: forbidden
```

Rules:

- Lookup keys and verifier keys must be separated.
- Credential keys and token keys must be separated.
- A compromise of one logical key must not directly reveal every digest class.
- A key set is selected by an internal `verifier_key_id`.
- New credential and token verifier records must store the `verifier_key_id` for the key set used to compute their digests.
- Future verifier code must reject incomplete key sets instead of silently reusing a key for multiple digest classes.

## 5. Planned Local Configuration Surface

The first local implementation may use these planned environment variable names after a later code gate:

```yaml
planned_environment_variables:
  verifier_key_set_id: VIBIT_AUTH_VERIFIER_KEY_SET_ID
  credential_lookup_key: VIBIT_AUTH_CREDENTIAL_LOOKUP_KEY
  credential_verifier_key: VIBIT_AUTH_CREDENTIAL_VERIFIER_KEY
  token_lookup_key: VIBIT_AUTH_TOKEN_LOOKUP_KEY
  token_verifier_key: VIBIT_AUTH_TOKEN_VERIFIER_KEY
```

These names are a planned configuration contract, not an implementation. This standard does not authorize `os.Getenv`, environment parsing, config structs, CLI flags, process wiring, or runtime key loading.

Rules for future implementation:

- Environment variable values are secret material, except the key set id which is secret-adjacent and not log-safe by default.
- Key values must decode to at least 32 bytes.
- Key values must be generated from cryptographically secure randomness for production.
- Key values must not be human-readable passwords, passphrases, repeated-pattern strings, all-zero values, user identifiers, device identifiers, provider subjects, or copied tokens.
- Missing, malformed, too-short, duplicate, or incomplete production key configuration must fail closed.
- Local ignored files may be used by developers to populate environment variables, but no committed file may contain production-like values.
- This repository's ignored `.vibit.local.env` convention is for local machine configuration only and must not be cited as a committed secret source.

## 6. Key Identifier Rules

`verifier_key_id` identifies the logical verifier key set used for a stored credential or token verifier record. It is not the secret key value.

Classification:

```yaml
verifier_key_id:
  secret_value: false
  public_api_field: false
  database_record_field: true
  internal_selection_input: true
  log_safe_by_default: false
  public_error_safe: false
  documentation_example_safe: placeholder_only
  conversation_log_safe: placeholder_only
  change_spec_example_safe: placeholder_only
```

Rules:

- `verifier_key_id` must not contain key bytes, encoded key values, credentials, tokens, account identifiers, tenant identifiers, cloud secret paths, provider secret names, access keys, environment variable values, hostnames, deployment names, or operator names.
- Public errors must not disclose a key id, key-set miss, key rotation state, key decoding error, or key length failure.
- Logs, traces, metrics labels, audit-safe facts, ADRs, change specs, docs, and conversation logs may use placeholders such as `<verifier-key-id>` only.
- A future operations standard may define a short redacted fingerprint format, but no fingerprint format is ratified here.

## 7. Rotation Posture

Future key rotation must be key-set based.

```yaml
rotation_model:
  active_key_set: required
  accepted_previous_key_sets: allowed
  new_writes_use_active_key_set: required
  stored_records_keep_verifier_key_id: required
  verification_may_select_key_by_record_key_id: required
  automatic_rotation_implemented_by_this_standard: false
```

Expected rotation phases:

1. Introduce a new complete key set.
2. Mark the new key set active for new credential and token verifier records.
3. Keep previous key sets available for the access-token TTL, credential replacement window, or other later ratified retention window.
4. Reissue, rotate, revoke, or replace records according to a later authentication behavior gate.
5. Retire old key sets only after no valid record requires them.

Rotation rules:

- Rotation must never require hand-editing stored digests.
- Rotation must not use database-only equality as verifier proof.
- Rotation failures must use the same public failure class as other invalid proof failures unless a later semantic error standard explicitly allows more detail.
- Key retirement must be observable internally without exposing raw keys or full key identifiers in public artifacts.

No rotation behavior is implemented by this standard.

## 8. Development And Test Posture

Development and test configuration must remain explicit.

Allowed after later implementation gates:

- Test-only deterministic fixture keys inside test code when they are clearly non-production and used only for repeatable unit tests.
- Local ignored files that set environment variables for a developer machine.
- Explicit in-memory test configuration supplied by tests.
- Redaction tests that use obvious sentinel strings to prove secret values do not leak.

Forbidden:

- Committed production-like key values.
- Default production keys.
- Automatic persistent random key generation.
- Silent fallback to a shared hard-coded development key.
- Using the key set id as a secret key.
- Using raw credential material, raw access-token material, provider tokens, user identifiers, device identifiers, or player identifiers as key material.

Automatic random keys may be allowed only for an explicitly ephemeral in-memory development mode after a later gate. They must not be used for durable local authentication state because restart would make stored verifier records unverifiable.

## 9. Failure And Fallback Behavior

Future production behavior must fail closed when required secret configuration is unavailable or invalid.

```yaml
failure_behavior:
  missing_required_key: fail_closed
  malformed_key_encoding: fail_closed
  decoded_key_too_short: fail_closed
  duplicate_logical_keys: fail_closed
  incomplete_key_set: fail_closed
  unknown_record_key_id: invalid_proof_public_failure
  retired_record_key_id: invalid_proof_public_failure
  public_error_discloses_secret_config_problem: forbidden
```

Rules:

- Startup may fail before serving requests when key configuration is invalid.
- Request-time validation must not disclose whether failure was caused by key configuration, missing records, invalid proof, expiration, revocation, account state, algorithm version, or key id unless a later semantic error standard explicitly allows it.
- Development behavior may use clearer local diagnostics only if diagnostics redact secret values and remain outside public client responses.
- Metrics labels must not include key values, encoded key values, full key ids, credentials, tokens, or cloud secret paths.

## 10. Redaction Requirements

Forbidden in logs, traces, metrics labels, public errors, panic output, audit-safe facts, test snapshots, fixtures, ADRs, change specs, documentation examples, and conversation logs:

- Verifier key values.
- Encoded verifier key values.
- Decoded verifier key bytes.
- Environment variable values.
- Secret-manager response payloads.
- Cloud secret paths that reveal tenancy or deployment structure.
- Provider credentials or provider secret names.
- Raw credential proof.
- Raw access-token text.
- Credential lookup digest.
- Credential verifier digest.
- Token lookup digest.
- Token verifier digest.
- Full concrete `verifier_key_id` values.

Allowed with care:

- Environment variable names.
- Placeholder key identifiers such as `<verifier-key-id>`.
- Non-secret configuration field names.
- Registered error codes.
- Internal lifecycle state names.

Secret configuration values are not safe because they are "only local." Agent-facing documents and conversation logs must treat local secrets as real secrets.

## 11. Dependency Posture

The first local secret configuration posture does not require an external dependency.

```yaml
external_kms_secret_manager_required_for_first_local_posture: false
process_environment_allowed_after_code_gate: true
explicit_runtime_secret_input_allowed_after_code_gate: true
dependency_adoption_record_required_for_external_secret_manager: true
operations_boundary_required_for_external_secret_manager: true
```

Deferred and not selected by this standard:

- KMS providers.
- Cloud secret-manager SDKs.
- Vault-like systems.
- Container orchestration secret APIs.
- Password-hashing dependencies.
- JWT, JWK, OAuth, or OIDC dependencies.
- Provider SDKs.

A future production deployment guide may require an external secret manager, but that decision belongs to a separate dependency and operations gate.

## 12. Non-Goals

This standard does not:

- Add secret loading.
- Add environment parsing.
- Add config structs.
- Add CLI flags.
- Add KMS or cloud secret-manager integration.
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
- Add production authentication behavior.

## 13. Verification Path

The repository check rule for this boundary is:

```text
runtime.secret_configuration_verifier_key_loading_boundary
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

- Token and credential material generation boundary.
- Verifier digest computation and constant-time comparison implementation boundary.
- Application authentication service implementation gate.
- Authentication redaction test implementation gate.
- Protobuf authentication message gate.
- WebSocket request proof carrier gate.
- Operations secret-management adoption gate.
