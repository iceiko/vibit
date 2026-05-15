# Authentication Service Implementation Readiness Gate

Status: Draft v0.1
Last updated: 2026-05-15
Scope: Entry criteria, package ownership, file boundaries, test expectations, sequencing, reference mapping, and deferrals before the first application authentication service implementation work
Depends on: `docs/runtime-authentication-implementation-boundary.md`, `docs/application-authentication-service-interface-boundary.md`, `docs/token-credential-verifier-algorithm-redaction-boundary.md`, `docs/secret-configuration-verifier-key-loading-boundary.md`, `docs/token-credential-material-generation-boundary.md`, `docs/verifier-digest-computation-comparison-boundary.md`
Canonical decision: `ADR-0044`

The paired Simplified Chinese translation is `docs/authentication-service-implementation-readiness-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This gate defines what must be true before vibit starts the first application authentication service implementation code.

It exists because authentication is security-sensitive and cross-cutting. Without an implementation readiness gate, a future agent could begin writing service code while quietly choosing package ownership, secret loading, token generation, digest helper names, repository call shape, protocol carriers, failure behavior, or test posture.

This is a readiness-only standard. It does not add Go code, imports, runtime services, secret loading, token generation, credential generation, verifier digest computation, verifier comparison, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket carriers, routes, repository methods, SQL migrations, authentication dependencies, or production authentication behavior.

## 2. Core Rule

The first authentication service implementation readiness posture is:

```yaml
authentication_service_implementation_readiness_gate: defined
implementation_authorized_by_this_standard: false
planned_owner: runtime/internal/app
planned_package_candidate: runtime/internal/app/authentication
first_code_slice_must_be_separately_authorized: true
service_code_status: deferred
secret_loading_code_status: deferred
material_generation_code_status: deferred
digest_helper_code_status: deferred
verifier_comparison_code_status: deferred
login_execution_status: deferred
token_validation_status: deferred
logout_execution_status: deferred
cleanup_execution_status: deferred
protobuf_authentication_messages_status: deferred
websocket_proof_carrier_status: deferred
repository_interface_change_status: deferred
migration_schema_change_status: deferred
major_dependency_status: deferred
```

This gate is complete only when agents can inspect a single standard and know:

- Which prior boundaries must be read.
- Which packages may be edited by the first implementation slice.
- Which behaviors remain separate gates.
- Which tests are required before behavior is accepted.
- Which Nakama and Pitaya capabilities are being adapted.
- Which public/protocol surfaces remain out of scope.

## 3. Required Prior Boundaries

Before implementation code starts, the implementing agent must read and preserve:

```yaml
required_boundaries:
  runtime_authentication_implementation_boundary: docs/runtime-authentication-implementation-boundary.md
  application_authentication_service_interface_boundary: docs/application-authentication-service-interface-boundary.md
  verifier_algorithm_redaction_boundary: docs/token-credential-verifier-algorithm-redaction-boundary.md
  secret_configuration_key_loading_boundary: docs/secret-configuration-verifier-key-loading-boundary.md
  material_generation_boundary: docs/token-credential-material-generation-boundary.md
  verifier_digest_computation_comparison_boundary: docs/verifier-digest-computation-comparison-boundary.md
```

The first implementation slice must not reinterpret those boundaries. If the implementation needs to change one, it must open a separate standard or ADR before code changes.

## 4. Package And File Ownership

The first authentication service implementation should stay application-owned.

Allowed future write area after a later implementation work item:

```text
runtime/internal/app/authentication/
```

Allowed application integration points after a later implementation work item:

```text
runtime/internal/app/
runtime/internal/app/bootstrap/
```

Allowed test area:

```text
runtime/internal/app/authentication/*_test.go
runtime/internal/app/*_test.go
```

Forbidden first-slice write areas unless a later work item explicitly names them:

- `runtime/internal/platform/transport/ws/`
- `runtime/internal/platform/protocol/protobuf/`
- `runtime/internal/generated/`
- `runtime/internal/modules/authentication/`
- `runtime/internal/platform/persistence/postgres/`
- `runtime/migrations/postgres/`
- `proto/`
- `contracts/runtime/authentication/`

The first implementation slice may use existing generated metadata and existing repository interfaces, but it must not hand-edit generated files or change public contracts as a side effect.

## 5. First Implementation Queue

The readiness gate recommends this implementation order:

```yaml
recommended_queue:
  - define_local_verifier_key_configuration_code_gate
  - implement_application_secret_configuration_loader
  - implement_token_credential_material_generation_helpers
  - implement_verifier_digest_helpers_and_comparison
  - implement_authentication_service_in_memory_unit_tests
  - implement_device_credential_authentication_service_flow
  - implement_access_token_validation_service_flow
  - define_protocol_authentication_message_gate
  - define_websocket_request_proof_carrier_gate
```

This queue is guidance, not implicit authorization. Each item still needs a bounded work item, change spec, tests, verification, and documentation update.

## 6. Minimum First Code Slice

The smallest acceptable first code slice should prove one narrow behavior without opening protocol or transport surfaces.

Recommended first code slice:

```yaml
first_code_slice:
  name: local_verifier_key_configuration_loader
  owner: runtime/internal/app/authentication
  behavior: load_explicit_in_memory_or_environment_supplied_key_set
  production_behavior: fail_closed_for_invalid_config
  external_dependencies: none
  protocol_changes: none
  repository_changes: none
  migration_changes: none
```

Why this comes first:

- Token and credential generation cannot be verified without key material.
- Digest helpers should receive validated logical keys instead of parsing configuration.
- Service code should not invent secret-loading behavior inside login logic.

This standard does not authorize the code slice. It names the recommended next implementation gate.

## 7. Service Behavior Entry Criteria

Before login or token validation behavior is implemented, these must be true:

```yaml
service_behavior_entry_criteria:
  explicit_key_configuration_loader_exists: required
  material_generation_helpers_exist: required_for_token_or_credential_issuance
  digest_helpers_exist: required
  constant_time_verifier_comparison_exists: required
  redaction_tests_exist: required
  repository_usage_through_unit_of_work: required
  generated_authentication_shapes_remain_metadata_only: required
  protobuf_authentication_messages_defined: required_before_wire_exposure
  websocket_proof_carrier_defined: required_before_realtime_exposure
```

Login execution must not be exposed through WebSocket or Protobuf until the protocol carrier gates exist.

Token validation must not become production-sensitive domain authorization until explicit request proof carrier behavior exists and request identity handoff is tested.

## 8. Required Tests For First Behavior

Future authentication implementation work must add tests proportional to the behavior it introduces.

Minimum test classes:

```yaml
required_test_classes:
  configuration_tests:
    - missing_key_fails_closed
    - malformed_key_fails_closed
    - duplicate_logical_keys_fail_closed
    - secret_values_absent_from_errors
  generation_tests:
    - generated_material_has_32_raw_bytes
    - generated_material_round_trips_through_text_encoding
    - raw_material_not_stored
  digest_tests:
    - canonical_input_is_stable
    - purpose_labels_are_separated
    - digest_output_is_32_bytes
    - lookup_digest_is_not_authentication_proof
  comparison_tests:
    - verifier_comparison_uses_constant_time_primitive
    - mismatch_and_missing_record_share_public_failure
  service_tests:
    - repository_is_used_through_unit_of_work
    - request_identity_is_populated_only_after_valid_proof
    - public_errors_are_registered_and_redacted
```

Tests must not require a running PostgreSQL server unless explicitly marked live and opt-in through the existing PostgreSQL verification standard.

## 9. Redaction And Observability

Future service code must preserve redaction by construction.

Forbidden in logs, public errors, traces, metrics labels, audit-safe facts, tests snapshots, ADRs, change specs, documentation examples, and conversation logs:

- Raw credentials.
- Raw access tokens.
- Encoded generated material.
- Lookup digests.
- Verifier digests.
- Verifier key values.
- Encoded verifier key values.
- Full concrete `verifier_key_id` values.
- HMAC input bytes.
- HMAC output bytes.
- Environment variable values.

Allowed:

- Registered error codes.
- Non-secret record identifiers.
- Lifecycle state names when explicitly safe.
- Placeholders such as `<raw-access-token>` and `<verifier-key-id>` when documenting redaction rules.

Observability must be useful without exposing proof material. A later observability standard may define redacted fingerprints, but none are ratified here.

## 10. Nakama And Pitaya Mapping

Nakama capabilities to adapt:

- Device/custom authentication entry point.
- Session token issuance.
- Session token validation.
- Token expiration and revocation.
- User/account status checks.

Pitaya capabilities to adapt:

- Frontend transport acceptor separated from backend handler logic.
- Route handler context receiving validated identity.
- Session binding as application/runtime context, not proof by itself.

vibit implementation rule:

- Transport accepts frames.
- Protocol decodes messages.
- Application authentication validates proof.
- Application request identity carries validated actor context.
- Domain modules consume request identity and do not parse proof.

Do not copy Nakama or Pitaya public APIs directly. Use them as capability coverage and vocabulary references.

## 11. Deferrals

This readiness gate does not authorize:

- Authentication service code.
- Secret loading code.
- Token generation.
- Credential generation.
- Verifier digest computation.
- Verifier comparison.
- Login execution.
- Access-token validation.
- Logout execution.
- Refresh behavior.
- Cleanup jobs.
- Protobuf authentication messages.
- WebSocket proof carriers.
- WebSocket handshake authentication.
- Authentication dependencies.
- `authentication.Repository` changes.
- PostgreSQL migration schema changes.
- Production authentication behavior.

## 12. Verification Path

The repository check rule for this gate is:

```text
runtime.authentication_service_implementation_readiness_gate
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

No runtime Go authentication behavior is verified by this standard because no behavior is added.

## 13. Completion Criteria

This gate is complete when:

- The readiness standard and translation exist.
- ADR-0044 records the readiness decision.
- Manifests and agent guides reference the gate.
- Repository checks enforce the readiness markers.
- `W-0095` is completed.
- The next ready work item is a bounded implementation or preparation gate.
