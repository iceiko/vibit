# Credential, Token, And Session Schema Gates

Status: Draft v0.1
Last updated: 2026-05-14
Scope: Future schema gates for credential records, external identity links, token verifier records, runtime session records, audit persistence, and player account lifecycle separation
Depends on: `docs/token-lifecycle-storage-implications.md`, `docs/authentication-contract-error-permission-surfaces.md`
Canonical decision: `ADR-0029`

The paired Simplified Chinese translation is `docs/credential-token-session-schema-gates.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This standard defines the schema gates that must exist before vibit implements the selected first authentication posture:

```yaml
login_method: device_credential_login
access_token_format: opaque_high_entropy_token
token_issuance_carrier: login_command_response_token
request_proof_carrier: explicit_request_proof_payload
refresh_token: not_in_first_implementation
renewal_method: reauthenticate_with_device_credential_login
default_durable_target: PostgreSQL
```

This document is a gate definition, not a database schema.

It does not add tables, migrations, repository interfaces, PostgreSQL adapters, runtime lookup code, token behavior, credential behavior, session persistence, audit persistence, Protobuf fields, WebSocket handshake behavior, runtime player handlers, or WebSocket routes.

## 2. Schema Gate Rule

A schema gate is a required planning artifact that must be satisfied before a future change may add persistent schema for a security-sensitive concern.

Each future persistent authentication concern must pass these gates in order:

```yaml
schema_ratification: required_before_migration
migration_source: separate_future_change
repository_interface: separate_future_change
postgres_adapter: separate_future_change
live_verification: separate_future_change_or_explicit_deferral
runtime_wiring: separate_future_implementation_milestone
```

Rules:

- A schema gate may name required record semantics.
- A schema gate may name forbidden shortcuts.
- A schema gate may name required future decisions.
- A schema gate must not create a migration.
- A schema gate must not imply that implementation is authorized.
- A schema gate must preserve the current player account lifecycle tables.

## 3. Gate Matrix

The W-0071 gate status is:

```yaml
credential_record_schema_gate:
  required_for_first_posture: true
  status: ratified_no_schema_added
  boundary: docs/credential-record-schema-boundary.md
  decision: ADR-0032
token_verifier_record_schema_gate:
  required_for_first_posture: true
  status: ratified_no_schema_added
  boundary: docs/token-verifier-record-schema-boundary.md
  decision: ADR-0033
external_identity_link_schema_gate:
  required_for_first_posture: false
  status: deferred_no_schema_added
runtime_session_record_schema_gate:
  required_for_first_posture: false
  status: deferred_no_schema_added
audit_persistence_schema_gate:
  required_before_durable_authentication_audit: true
  status: defined_no_schema_added
player_account_lifecycle_schema:
  status: preserved
  credential_columns_added: false
  token_columns_added: false
  external_identity_columns_added: false
  session_columns_added: false
```

The first implementation may not begin until the required credential and token verifier schema gates are turned into explicit schema ratification work, followed by migrations, repositories, adapters, tests, and runtime wiring.

W-0074 has ratified the credential record schema boundary without adding schema:

```yaml
credential_record_schema_boundary:
  status: migration_source_added
  standard: docs/credential-record-schema-boundary.md
  decision: ADR-0032
  future_logical_table: authentication_device_credentials
  migration_source: runtime/migrations/postgres/000003_create_authentication_device_credentials.sql
  migration_added_now: true
  repository_added_now: false
  runtime_lookup_added_now: false
```

W-0075 ratified the token verifier record schema boundary, and W-0078 has added its migration source:

```yaml
token_verifier_record_schema_boundary:
  status: migration_source_added
  standard: docs/token-verifier-record-schema-boundary.md
  decision: ADR-0033
  future_logical_table: authentication_access_tokens
  migration_source: runtime/migrations/postgres/000004_create_authentication_access_tokens.sql
  migration_added_now: true
  repository_added_now: false
  runtime_validation_added_now: false
  token_issuance_added_now: false
  logout_added_now: false
  cleanup_added_now: false
```

W-0077 and W-0078 have added the credential and token verifier migration sources. Authentication migration static checks remain the next required step before repository interfaces, adapters, or runtime authentication behavior may be added.

## 4. Credential Record Gate

Credential records are required before `device_credential_login` can be implemented.

Gate status:

```yaml
required_for_first_posture: true
default_durable_target: PostgreSQL
schema_added_now: false
migration_added_now: false
repository_added_now: false
runtime_lookup_added_now: false
```

Future credential schema ratification must define:

- Credential record owner.
- Record lifecycle states.
- Account creation and account lookup relationship.
- Whether one player can have multiple device credential records.
- Whether one credential can ever move between players.
- Which identifier is safe to log.
- Which secret verifier is stored.
- Which raw credential material is forbidden.
- Which fields are unique.
- Which fields are mutable.
- Rotation, replacement, revocation, and disabled-credential behavior.
- How disabled or deleted player account states block login.
- Which operations must be atomic with player account lifecycle changes.
- Redaction rules for logs, errors, traces, tests, fixtures, and conversation logs.
- Abuse-control and retryability classes.

Required first posture constraints:

```yaml
credential_kind: device_credential_login
raw_os_device_id_as_credential: forbidden
raw_credential_storage: forbidden
credential_verifier_storage: required
credential_record_id_log_safe: required
player_account_lifecycle_tables_store_credentials: forbidden
```

This gate does not choose a password model, OAuth provider, OIDC issuer, platform identity provider, provider SDK, key-management dependency, password-hashing dependency, or external identity linking behavior.

## 5. Token Verifier Record Gate

Token verifier records are required before opaque access-token validation can be implemented.

Gate status:

```yaml
required_for_first_posture: true
default_durable_target: PostgreSQL
schema_added_now: false
migration_added_now: false
repository_added_now: false
runtime_validation_added_now: false
redis_like_store_selected: false
```

Future token verifier schema ratification must define:

- Token record owner.
- Non-plaintext verifier storage.
- Verifier algorithm and comparison rules.
- Token record identifier that is safe to log.
- Actor kind and actor identifier relationship.
- Player account relationship for player access tokens.
- Credential-token linkage required for rotation and presented-token logout.
- Audience and route eligibility semantics.
- Minimum statuses: `active`, `expired`, and `revoked`.
- `issued_at`, `expires_at`, and `revoked_at` semantics.
- Revocation reason retention.
- Replacement or rotation lineage.
- Expired and revoked record retention.
- Cleanup owner and cleanup trigger.
- Replay-sensitive failure classes.
- Redaction rules for logs, errors, traces, tests, fixtures, and conversation logs.

Required first posture constraints:

```yaml
raw_token_storage: forbidden
minimum_entropy_bits: 256
access_token_ttl: 1h
refresh_token_storage: forbidden_for_first_posture
session_token_vocabulary: deferred_until_session_persistence
token_record_id_log_safe: required
```

The first token verifier schema must be able to revoke the presented access token. It should also preserve enough credential linkage to support rotation of previous active tokens for the same credential when that implementation step is authorized.

## 6. External Identity Link Gate

External identity storage is not required for the first `device_credential_login` posture.

Gate status:

```yaml
required_for_first_posture: false
status: deferred
schema_added_now: false
migration_added_now: false
provider_dependency_added_now: false
```

Future external identity schema ratification must define:

- Provider namespace semantics.
- Provider subject semantics.
- Whether provider subjects are globally unique or provider-scoped.
- Whether one account can have multiple provider links.
- Whether one provider subject can map to more than one vibit account.
- Link, unlink, conflict, recovery, and merge behavior.
- Provider metadata retention and redaction.
- Which events are audit-only and which may be client-visible.
- Which provider dependencies are adopted and where they are allowed.

This gate keeps provider identity separate from credential storage, token verifier storage, runtime sessions, and player account lifecycle storage.

## 7. Runtime Session Record Gate

Runtime session persistence is not required for the first access-token posture.

Gate status:

```yaml
required_for_first_posture: false
status: deferred
schema_added_now: false
migration_added_now: false
session_store_selected_now: false
websocket_connection_binding_selected_now: false
```

Future runtime session schema ratification must define:

- Whether runtime sessions exist as durable records.
- Whether PostgreSQL remains sufficient or a Redis-like store is justified.
- Session identifier semantics.
- Session token semantics, if any.
- Access-token-to-session relationship, if any.
- WebSocket connection binding.
- Reconnect behavior.
- Connection epoch behavior.
- Expiration, revocation, and cleanup.
- Whether validation happens per request, on first message, during WebSocket handshake, or by a hybrid model.
- Whether Protobuf envelope fields change.

Until that gate is ratified, current Protobuf `Session` fields remain metadata-only and must not become proof.

## 8. Audit Persistence Gate

W-0070 defined semantic audit-oriented events. Durable audit persistence is not added by W-0071.

Future durable audit schema ratification must define:

- Whether authentication audit is stored in module tables, an event log, an outbox, operational logs, or another approved store.
- Which semantic events become durable rows.
- Which identifiers are safe to store.
- Which fields are explicitly forbidden.
- Retention and cleanup policy.
- Query and inspection requirements.
- Whether audit writes share a unit of work with credential or token verifier mutations.

Forbidden audit payload material:

- Raw credential values.
- Raw access-token values.
- Token verifier hashes.
- Password hashes.
- Provider secrets.
- Full provider payloads.
- WebSocket connection secrets.

## 9. Player Account Lifecycle Preservation

Current player account lifecycle storage remains:

```text
player_accounts
player_account_events
```

These tables remain lifecycle storage only.

Forbidden in player lifecycle storage:

- Credential columns or rows.
- Password hashes.
- External identity provider subjects.
- Access-token state.
- Refresh-token state.
- Runtime session state.
- WebSocket connection state.
- Request validation results.
- Raw authentication proof.

Future authentication work may reference player accounts through explicit repository boundaries, but it must not convert player lifecycle storage into authentication storage.

## 10. Future Work Separation

Future implementation must be split at least this way:

```yaml
credential_schema_ratification: separate_work
token_verifier_schema_ratification: separate_work
credential_and_token_migration_sources: completed_by_W_0077_and_W_0078
repository_interfaces: separate_work
postgres_adapters: separate_work
redaction_and_schema_tests: separate_work
runtime_authentication_wiring: separate_work
protocol_or_websocket_changes: separate_decision_if_needed
```

Agents must not combine schema ratification, migration creation, repository implementation, and runtime authentication behavior in a single broad change unless a future work item explicitly grants that scope.

## 11. Reference Alignment

### Nakama

Nakama remains the capability reference for accounts, authentication methods, session tokens, refresh, logout, expiry, and revocation vocabulary.

vibit adapts the capability coverage but does not copy Nakama's public API, token/session model, or storage shape. Refresh tokens and session token vocabulary remain deferred.

### Pitaya

Pitaya remains the vocabulary reference for sessions, handler context, connection binding, frontend/backend roles, and realtime server routing.

vibit adapts the separation between connection acceptors and application identity. It does not place credential or token validation inside WebSocket acceptors, and it does not copy Pitaya's public API shape.

## 12. Verification Path

Default verification for this standard:

```bash
node tools/vibit inspect next --json
node tools/vibit check runtime --json
node tools/vibit check architecture --json
node tools/vibit check agent-tooling --json
node tools/vibit check memory --json
node tools/vibit check work --json
node tools/vibit check change define-credential-token-session-schema-gates --json
node tools/vibit check all --json
git diff --check
```

Go runtime tests are not required for this gate-only change because no Go runtime behavior changes.

## 13. Non-Authorization

This standard does not authorize:

- Credential tables.
- Token tables.
- External identity tables.
- Runtime session tables.
- Audit tables.
- Migrations.
- Repository interfaces.
- PostgreSQL adapters.
- Runtime credential lookup.
- Token generation, parsing, validation, refresh, revocation, rotation, replay handling, cleanup, or storage.
- Login handlers.
- Runtime player handlers.
- WebSocket routes.
- Protobuf messages or generated Protobuf output.
- WebSocket handshake authentication.
- First system-message authentication.
- Password hashing, JWT, OAuth, OIDC, provider SDK, Redis-like, cryptography, key-management, or major authentication dependencies.
- Treating metadata-only `player_id`, `session_id`, `connection_id`, or `connection_epoch` as proof.

## 14. Follow-Up

Next work:

```text
W-0079 Add authentication migration static checks
```

W-0079 should harden local checks for the ratified authentication migration sources without adding repositories, adapters, runtime token validation, generated output, Protobuf messages, WebSocket behavior, or authentication dependencies.
