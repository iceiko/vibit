# Token Verifier Record Schema Boundary

Status: Draft v0.1
Last updated: 2026-05-14
Scope: Token verifier record schema boundary for the first opaque access-token posture
Depends on: `docs/credential-token-session-schema-gates.md`, `docs/token-lifecycle-storage-implications.md`
Canonical decision: `ADR-0033`

The paired Simplified Chinese translation is `docs/token-verifier-record-schema-boundary.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This standard ratifies the token verifier record boundary required before vibit can implement the selected first access-token posture:

```yaml
access_token_format: opaque_high_entropy_token
token_issuance_carrier: login_command_response_token
request_proof_carrier: explicit_request_proof_payload
default_durable_target: PostgreSQL
schema_boundary_status: ratified_no_schema_added
```

This document defines future record semantics. It does not add SQL migration source, tables, repository interfaces, PostgreSQL adapters, runtime token validation, token issuance, logout behavior, refresh behavior, cleanup jobs, generated authentication shapes, Protobuf messages, WebSocket proof carriers, WebSocket handshake authentication, authentication dependencies, or production authentication behavior.

## 2. Required Reading

Read this standard together with:

- `docs/credential-token-session-schema-gates.md`
- `docs/token-lifecycle-storage-implications.md`
- `docs/first-token-format-proof-carrier-posture.md`
- `docs/authentication-contract-error-permission-surfaces.md`
- `docs/selected-login-token-boundary-checks.md`
- `docs/credential-record-schema-boundary.md`
- `docs/postgresql-persistence-boundary.md`
- `docs/player-identity-session-boundary.md`
- `.arch/runtime.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `ADR-0026`
- `ADR-0027`
- `ADR-0029`
- `ADR-0030`
- `ADR-0031`
- `ADR-0032`
- `ADR-0033`

Reference reading:

- Nakama authentication concepts: `https://heroiclabs.com/docs/nakama/concepts/authentication/`
- Nakama session concepts: `https://heroiclabs.com/docs/nakama/concepts/session/`
- Pitaya features and session vocabulary: `https://pitaya.readthedocs.io/en/stable/features.html`
- Pitaya API and handler vocabulary: `https://pitaya.readthedocs.io/en/latest/API.html`

Nakama and Pitaya remain references for capability coverage, session vocabulary, token lifecycle pressure, and handler context vocabulary. They do not govern vibit's token schema, public API, session model, transport behavior, or generated boundaries.

## 3. Boundary Status

The token verifier record boundary is ratified as a schema boundary only:

```yaml
token_verifier_record_schema_boundary:
  status: ratified_no_schema_added
  default_durable_target: PostgreSQL
  future_logical_table: authentication_access_tokens
  owner: runtime.authentication
  migration_source_added: false
  repository_interface_added: false
  postgres_adapter_added: false
  runtime_validation_added: false
  token_issuance_added: false
  logout_added: false
  cleanup_added: false
  authentication_implemented: false
```

The future logical table name is ratified so agents have a stable target during migration planning. No table exists until a later migration work item creates SQL source under `runtime/migrations/postgres/`.

## 4. Ownership

The future token verifier record belongs to the runtime authentication boundary:

```text
runtime.authentication
```

Future implementation ownership:

- Semantic token contract owner: `contracts/runtime/authentication/`
- Runtime validation and identity handoff owner: `runtime/internal/app`
- Future repository interface owner: a later ratified authentication repository boundary
- Future PostgreSQL adapter owner: `runtime/internal/platform/persistence/postgres/`
- Future migration source owner: `runtime/migrations/postgres/`
- Future cleanup owner: a later ratified authentication maintenance or token storage boundary

The player module owns player account lifecycle state. It does not own access-token verifier records.

Player account lifecycle tables remain:

```text
player_accounts
player_account_events
```

They must stay credential-free, token-free, provider-subject-free, runtime-session-free, WebSocket-state-free, and request-validation-free.

## 5. Logical Record

The future token verifier record represents one server-side verifier for one opaque access token.

Future logical record:

```yaml
record: authentication_access_token
fields:
  token_record_id: log_safe_identifier
  token_kind: access_token
  token_status: active | expired | revoked
  actor_kind: player
  player_id: player_account_reference
  credential_record_id: credential_record_reference
  token_lookup_digest: secret_adjacent_index_material
  token_verifier_digest: secret_verifier_material
  verifier_algorithm: versioned_non_plaintext_verifier
  verifier_version: integer
  verifier_key_id: secret_key_reference_not_secret_value
  audience: route_or_gameplay_audience_catalog
  issued_at: timestamp
  expires_at: timestamp
  revoked_at: nullable_timestamp
  revoked_reason: nullable_catalog_value
  replaced_by_token_record_id: nullable_log_safe_identifier
  last_validated_at: nullable_timestamp
  last_failed_validation_at: nullable_timestamp
  cleanup_after: nullable_timestamp
  created_at: timestamp
  updated_at: timestamp
```

The future SQL migration may refine exact PostgreSQL types, constraints, and index names, but it must preserve these semantics unless a later ADR supersedes this boundary.

No generic `metadata` column is ratified for the first token verifier record. Additional metadata requires a future schema decision because arbitrary JSON fields are a common place for agents to hide raw tokens, claims, transport state, provider payloads, device fingerprints, or private request details.

## 6. Identifier Rules

`token_record_id` is the only token identifier safe for normal logs, change specs, ADRs, tests, and conversation logs.

Rules:

- `token_record_id` must be globally unique.
- `token_record_id` must not be derived from raw token material.
- `token_record_id` must not reveal `player_id`, `credential_record_id`, session metadata, route name, IP address, user agent, provider subject, or token issuance time.
- `token_lookup_digest` is not log-safe.
- `token_verifier_digest` is not log-safe.
- Raw access-token text is never log-safe.
- `player_id` is not secret proof and must not satisfy authentication by itself.
- `credential_record_id` is log-safe as a credential identifier, but it must not be treated as token proof.

## 7. Verifier Semantics

The raw access token is never stored.

The first verifier posture is:

```yaml
raw_token_storage: forbidden
token_lookup_digest_required: true
token_verifier_digest_required: true
verifier_algorithm_versioned: true
plaintext_comparison: forbidden
minimum_entropy_bits: 256
jwt_or_claim_parsing_required: false
signing_dependency_required: false
redis_like_store_required: false
```

The token is expected to be high entropy. Therefore this boundary does not ratify JWT, JWK, OAuth, OIDC, signing, key-management, Redis-like token/session stores, provider SDKs, bcrypt, or Argon2 dependencies.

Future implementation must still define the exact verifier algorithm before code exists. The implementation may use standard-library cryptographic primitives if a later implementation work item ratifies the exact algorithm, pepper/secret configuration boundary, digest format, and comparison behavior.

## 8. Lifecycle States

Ratified states:

| State | Meaning | Request validation allowed |
| --- | --- | --- |
| `active` | Token verifier may authenticate a request if it is not expired and the linked actor and credential gates pass. | Yes |
| `expired` | Token verifier exceeded its validity window. | No |
| `revoked` | Token verifier was explicitly invalidated before or after expiration. | No |

State rules:

- `expired` and `revoked` are terminal for the first posture.
- `active` tokens become invalid after `expires_at` even if `token_status` has not yet been materialized as `expired`.
- A future validator may treat expiration as computed state rather than eagerly updating every expired row.
- Revocation must take effect before production-sensitive domain dispatch.
- A disabled or deleted player account blocks token validation even if the token record is still active.
- A disabled, revoked, or replaced credential blocks token validation if the validator checks credential state as part of the authorized implementation.

## 9. Actor And Credential Relationship

The first token verifier relationship is:

```yaml
one_token_record_represents_one_access_token: true
one_access_token_belongs_to_one_actor: true
first_actor_kind: player
player_id_mutable: false
credential_record_id_required_for_first_posture: true
token_can_move_between_players: false
token_can_move_between_credentials: false
refresh_token_storage: forbidden_for_first_posture
runtime_session_binding: deferred
websocket_connection_binding: deferred
```

The first posture requires `credential_record_id` linkage so later implementation can revoke previous active tokens for the same credential during successful login rotation. This linkage is not a substitute for token proof and does not authorize credential lookup, token validation, or logout behavior by itself.

Future service or admin tokens require separate actor-kind decisions because their lifetime, audience, permission, storage, and audit rules differ from player access tokens.

## 10. Expiration

The first access-token TTL remains:

```yaml
access_token_ttl: 1h
```

Future token verifier records must preserve:

```yaml
issued_at_required: true
expires_at_required: true
expires_after_issued_at: true
expired_proof_failure_class: required
expired_record_retention: required
```

An expired token must fail as `expired_proof`, not as missing proof or malformed proof. Expired records may be retained for replay analysis, logout idempotency, abuse investigation, and audit correlation.

## 11. Revocation, Logout, And Rotation

The first logout scope remains:

```yaml
logout_scope_first_posture: presented_access_token
```

Revocation semantics:

- Logout revokes the presented access token only.
- Admin revocation remains deferred to a future permission surface.
- Logout-all-devices remains deferred.
- Forced revocation for a disabled or deleted account remains deferred to account policy and audit work.
- Revoked tokens must fail distinctly from malformed, missing, invalid, and expired proof.
- Revocation reason must be retained as a non-secret catalog value.

Rotation semantics:

```yaml
rotation_on_successful_login: required_when_implementation_exists
previous_active_token_for_same_credential: revoke_when_repository_supports_it
replaced_by_token_record_id: optional_lineage_field
automatic_background_rotation: deferred
refresh_token_rotation: deferred
```

Successful login should issue a new token and revoke previous active access tokens for the same credential once schema, repository, and implementation gates authorize that behavior. This does not require runtime session persistence.

## 12. Retention And Cleanup

Cleanup is required before production token storage is enabled, but no cleanup job is added by this standard.

First retention posture:

```yaml
active_records_retained_until_expiration_or_revocation: true
expired_record_default_retention_recommendation: 7d
revoked_record_default_retention_recommendation: 7d
cleanup_target: expired_and_revoked_token_verifier_records
cleanup_trigger_first_posture: explicit_maintenance_command_or_scheduled_runtime_job_deferred
cleanup_owner: future_authentication_or_token_storage_boundary
```

The future migration may use `cleanup_after` to make retention explicit. If it does not, the repository or cleanup boundary must define how records become eligible for deletion.

Cleanup must not delete active tokens. Cleanup must not be hidden inside request validation. Cleanup must be idempotent, concurrency-aware, and auditable before production use.

## 13. Replay-Sensitive Failure Classes

Future validation must preserve these failure classes:

```yaml
failure_classes:
  - missing_proof
  - malformed_proof
  - unsupported_proof
  - invalid_proof
  - expired_proof
  - revoked_proof
  - actor_disabled
  - validator_unavailable
```

Replay-sensitive rules:

- A malformed token must not trigger lookup behavior that leaks whether a valid token exists.
- An invalid token must not reveal whether the lookup digest matched but the verifier digest failed.
- Expired and revoked tokens may produce distinct stable failure classes after proof format is accepted.
- Failure responses must not disclose whether a player account, credential record, or token record exists beyond the ratified error surface.
- Rate limiting and abuse controls remain future work, but schema and error design must not prevent them.

## 14. Uniqueness And Index Rules

Future migration work must preserve these uniqueness semantics:

```yaml
unique:
  - token_record_id
  - token_lookup_digest
indexes:
  - player_id
  - credential_record_id
  - token_status
  - expires_at
  - cleanup_after
foreign_key_like_relationships:
  - player_id references player account lifecycle identity
  - credential_record_id references authentication_device_credentials
```

The future SQL source may implement `player_id` and `credential_record_id` as real foreign keys only if that does not break module ownership, migration order, or test isolation. Whether each relationship is a database foreign key or an application-enforced reference must be explicit in the migration work item.

## 15. Redaction

Forbidden in logs, errors, traces, tests, fixtures, ADRs, change specs, and conversation logs:

- Raw access tokens.
- Token lookup digest.
- Token verifier digest.
- Server-side verifier secrets or peppers.
- Authorization header contents.
- Cookie contents.
- WebSocket subprotocol token carriers.
- URL query token carriers.
- Full request proof payloads containing token text.
- Credential lookup digest.
- Credential verifier digest.
- Provider secrets.

Allowed with care:

- `token_record_id`
- `credential_record_id`
- `player_id`
- lifecycle state names
- non-secret reason catalog values
- stable failure class names

Test fixtures must use synthetic values and must not contain real tokens, credentials, device identifiers, provider payloads, authorization headers, cookies, or copied production data.

## 16. Forbidden Shortcuts

Agents must not:

- Add `authentication_access_tokens` SQL migration source from this standard alone.
- Add token repository interfaces from this standard alone.
- Add PostgreSQL token adapters from this standard alone.
- Add runtime token issuance, validation, logout, refresh, or cleanup from this standard alone.
- Store tokens in `player_accounts` or `player_account_events`.
- Store token verifier records in credential records.
- Store raw access tokens, refresh tokens, runtime sessions, WebSocket connection state, or request validation results in token verifier records.
- Treat current Protobuf `Session` fields as access-token proof.
- Put token proof into WebSocket handshake headers, cookies, subprotocols, query parameters, routes, request IDs, player IDs, session IDs, or connection IDs.
- Add JWT, signing, OAuth, OIDC, provider SDK, key-management, Redis-like, password hashing, or other major authentication dependencies from this standard alone.
- Copy Nakama or Pitaya public API shapes.

## 17. Required Future Gates

Before token storage can be implemented, future work must complete these gates:

```yaml
credential_record_schema_boundary: completed_by_W_0074
token_verifier_record_schema_boundary: completed_by_W_0075
authentication_schema_migration_queue: required_before_migration
token_verifier_migration_source: separate_future_work
authentication_repository_interface: separate_future_work
token_postgres_adapter: separate_future_work
verifier_algorithm_and_secret_configuration: separate_future_work
redaction_tests: separate_future_or_implementation_work
runtime_token_validation_wiring: separate_future_implementation_milestone
logout_behavior: separate_future_implementation_milestone
cleanup_behavior: separate_future_maintenance_milestone
```

The migration gate must prove that SQL source creates only ratified token verifier structures and does not modify player account lifecycle tables.

## 18. Reference Alignment

### Nakama

Nakama demonstrates that mature game backends commonly provide access tokens or session tokens, expiration, refresh, revocation, logout, and multiple authentication entry points.

vibit adapts the lifecycle pressure but not the public API shape. The first vibit posture uses opaque high-entropy access tokens with server-side verifier storage, no refresh token, login-command issuance, and explicit request proof payloads.

### Pitaya

Pitaya demonstrates useful session and handler vocabulary around connection acceptors, sessions, ID binding, frontend/backend session differences, and route-aware message handling.

vibit keeps those ideas separate from token verifier storage. Token verifier records are persistent authentication proof verifier records, not transport sessions, handler context state, or WebSocket connection binding records.

## 19. Verification

Default verification for this standard:

```bash
node tools/vibit inspect next --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check architecture --json
node tools/vibit check memory --json
node tools/vibit check change define-token-verifier-record-schema-boundary --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Live PostgreSQL verification is not required because this work does not add migration source, tables, repositories, adapters, or runtime behavior.

Go tests are not required because no Go runtime behavior changes.

## 20. Follow-Up

Next work:

```text
W-0076 Plan authentication schema migration queue
```

The migration queue must plan credential and token verifier migration order, repository-interface gates, PostgreSQL adapter gates, redaction checks, and live verification expectations before any authentication migration source is added.
