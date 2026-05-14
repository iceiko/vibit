# Credential Record Schema Boundary

Status: Draft v0.1
Last updated: 2026-05-14
Scope: Credential record schema boundary for the first `device_credential_login` posture
Depends on: `docs/credential-token-session-schema-gates.md`
Canonical decision: `ADR-0032`

The paired Simplified Chinese translation is `docs/credential-record-schema-boundary.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This standard ratifies the credential record boundary required before vibit can implement the selected first login method:

```yaml
login_method: device_credential_login
credential_kind: high_entropy_installation_credential
default_durable_target: PostgreSQL
schema_boundary_status: ratified_no_schema_added
```

This document defines the future record semantics. It does not add SQL migration source, tables, repository interfaces, PostgreSQL adapters, runtime lookup, login handlers, token behavior, generated authentication shapes, Protobuf messages, WebSocket proof carriers, WebSocket handshake authentication, authentication dependencies, or production authentication behavior.

## 2. Required Reading

Read this standard together with:

- `docs/credential-token-session-schema-gates.md`
- `docs/first-login-method-set.md`
- `docs/token-lifecycle-storage-implications.md`
- `docs/authentication-contract-error-permission-surfaces.md`
- `docs/selected-login-token-boundary-checks.md`
- `docs/credential-storage-external-identity-linking-boundaries.md`
- `docs/postgresql-persistence-boundary.md`
- `docs/player-identity-session-boundary.md`
- `.arch/runtime.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `ADR-0025`
- `ADR-0029`
- `ADR-0030`
- `ADR-0031`
- `ADR-0032`

Reference reading:

- Nakama authentication concepts: `https://heroiclabs.com/docs/nakama/concepts/authentication/`
- Nakama session concepts: `https://heroiclabs.com/docs/nakama/concepts/session/`
- Pitaya features and session vocabulary: `https://pitaya.readthedocs.io/en/stable/features.html`
- Pitaya API and handler vocabulary: `https://pitaya.readthedocs.io/en/latest/API.html`

Nakama and Pitaya remain references for capability coverage and vocabulary. They do not govern vibit's credential schema, public API, session model, or generated boundaries.

## 3. Boundary Status

The credential record boundary is ratified as a schema boundary only:

```yaml
credential_record_schema_boundary:
  status: ratified_no_schema_added
  default_durable_target: PostgreSQL
  future_logical_table: authentication_device_credentials
  owner: runtime.authentication
  migration_source_added: false
  repository_interface_added: false
  postgres_adapter_added: false
  runtime_lookup_added: false
  authentication_implemented: false
```

The future logical table name is ratified so agents have a stable target during migration planning. No table exists until a later migration work item creates SQL source under `runtime/migrations/postgres/`.

## 4. Ownership

The future credential record belongs to the runtime authentication boundary:

```text
runtime.authentication
```

Future implementation ownership:

- Semantic login contract owner: `contracts/runtime/authentication/`
- Runtime validation and identity handoff owner: `runtime/internal/app`
- Future repository interface owner: a later ratified authentication repository boundary
- Future PostgreSQL adapter owner: `runtime/internal/platform/persistence/postgres/`
- Future migration source owner: `runtime/migrations/postgres/`

The player module owns player account lifecycle state. It does not own credential verifier records.

Player account lifecycle tables remain:

```text
player_accounts
player_account_events
```

They must stay credential-free, token-free, provider-subject-free, runtime-session-free, WebSocket-state-free, and request-validation-free.

## 5. Logical Record

The future credential record represents one durable high-entropy installation credential bound to one player account.

Future logical record:

```yaml
record: authentication_device_credential
fields:
  credential_record_id: log_safe_identifier
  player_id: player_account_reference
  credential_kind: device_credential_login
  credential_status: active | disabled | revoked | replaced
  credential_lookup_digest: secret_adjacent_index_material
  credential_verifier_digest: secret_verifier_material
  verifier_algorithm: versioned_non_plaintext_verifier
  verifier_version: integer
  verifier_key_id: secret_key_reference_not_secret_value
  client_instance_id_digest: optional_privacy_sensitive_correlation_material
  created_at: timestamp
  updated_at: timestamp
  last_verified_at: nullable_timestamp
  disabled_at: nullable_timestamp
  disabled_reason: nullable_catalog_value
  revoked_at: nullable_timestamp
  revoked_reason: nullable_catalog_value
  replaced_at: nullable_timestamp
  replaced_by_credential_record_id: nullable_log_safe_identifier
```

The future SQL migration may refine exact PostgreSQL types, constraints, and index names, but it must preserve these semantics unless a later ADR supersedes this boundary.

No generic `metadata` column is ratified for the first credential record. Additional metadata requires a future schema decision because arbitrary JSON fields are a common place for agents to hide credentials, provider payloads, tokens, or transport state.

## 6. Identifier Rules

`credential_record_id` is the only credential identifier safe for normal logs, change specs, ADRs, tests, and conversation logs.

Rules:

- `credential_record_id` must be globally unique.
- `credential_record_id` must not be derived from raw credential material.
- `credential_record_id` must not reveal device model, platform account, provider subject, IP address, user agent, player name, or player display name.
- `credential_lookup_digest` is not log-safe.
- `credential_verifier_digest` is not log-safe.
- `client_instance_id_digest` is not log-safe by default.
- `player_id` is not secret proof and must not satisfy authentication by itself.

## 7. Verifier Semantics

The raw credential proof is never stored.

The first verifier posture is:

```yaml
raw_credential_storage: forbidden
raw_os_device_id_as_credential: forbidden
credential_lookup_digest_required: true
credential_verifier_digest_required: true
verifier_algorithm_versioned: true
plaintext_comparison: forbidden
password_hashing_dependency_required: false
external_provider_dependency_required: false
```

The credential is expected to be high entropy. Therefore this boundary does not ratify bcrypt, Argon2, OAuth, OIDC, JWT, provider SDKs, or key-management dependencies.

Future implementation must still define the exact verifier algorithm before code exists. The implementation may use standard-library cryptographic primitives if a later implementation work item ratifies the exact algorithm and secret configuration boundary. If a future design wants a password-like or low-entropy credential, it must use a separate credential boundary and dependency adoption record.

## 8. Lifecycle States

Ratified states:

| State | Meaning | Login allowed |
| --- | --- | --- |
| `active` | Credential may authenticate if the linked player account is also active and token issuance gates pass. | Yes |
| `disabled` | Credential is temporarily blocked by policy or operations. | No |
| `revoked` | Credential is permanently invalidated. | No |
| `replaced` | Credential was superseded by a rotated or replacement credential. | No |

State rules:

- `revoked` and `replaced` are terminal.
- `disabled` may become `active` only through a later explicitly authorized administrative or recovery flow.
- `active` login is still blocked when the linked player account is disabled or deleted.
- Credential status must not override player account lifecycle state.
- A disabled or deleted player account must not be silently re-enabled by credential login.

## 9. Player Relationship

The first credential record relationship is:

```yaml
one_credential_record_belongs_to_one_player: true
credential_player_id_mutable: false
credential_can_move_between_players: false
one_player_active_device_credentials_first_posture: at_most_one
historical_records_per_player_allowed: true
multi_device_linking: deferred
account_recovery: deferred
account_merge: deferred
```

The first posture allows historical credential records for the same player because rotation and revocation need an audit trail. It does not authorize multiple concurrent active device credentials for one player. Multi-device linking, account recovery, and account merge require later account-linking decisions.

When first login creates a player account, player account creation and credential record creation must be atomic in the same application-owned unit of work. When a login authenticates an existing account, credential verification and token verifier creation must preserve consistent failure behavior without exposing whether a credential, player account, or token record exists.

## 10. Uniqueness And Index Rules

Future migration work must preserve these uniqueness semantics:

```yaml
unique:
  - credential_record_id
  - credential_lookup_digest
conditional_unique:
  - at_most_one_active_device_credential_per_player_for_first_posture
foreign_key_like_relationship:
  - player_id references player account lifecycle identity
```

The future SQL source may implement `player_id` as a real foreign key only if that does not break module ownership, migration order, or test isolation. Whether the relationship is a database foreign key or an application-enforced reference must be explicit in the migration work item.

## 11. Rotation And Replacement

Credential rotation creates a new credential record and marks the previous active credential as `replaced`.

Rules:

- Do not overwrite verifier material in place for rotation.
- Preserve lineage through `replaced_by_credential_record_id`.
- Rotation must happen in one unit of work when it changes both old and new records.
- If token verifier records exist, rotation impact on active access tokens must follow the token verifier schema boundary.
- Presented credential revocation must not require raw credential storage.

The first schema boundary prepares for rotation. It does not implement rotation commands or runtime behavior.

## 12. Redaction

Forbidden in logs, errors, traces, tests, fixtures, ADRs, change specs, and conversation logs:

- Raw credential proof.
- Raw operating-system device ID used as proof.
- Credential lookup digest.
- Credential verifier digest.
- Server-side verifier secrets or peppers.
- Full client instance identifiers.
- Raw access tokens.
- Token verifier hashes.
- Provider secrets.

Allowed with care:

- `credential_record_id`
- `player_id`
- lifecycle state names
- non-secret reason catalog values

Test fixtures must use synthetic values and must not contain real device identifiers, real credentials, real tokens, or copied production data.

## 13. Forbidden Shortcuts

Agents must not:

- Add `authentication_device_credentials` SQL migration source from this standard alone.
- Add credential repository interfaces from this standard alone.
- Add PostgreSQL credential adapters from this standard alone.
- Add runtime credential lookup or login handlers from this standard alone.
- Store credentials in `player_accounts` or `player_account_events`.
- Store provider subjects in credential records.
- Store access tokens, refresh tokens, runtime sessions, WebSocket connection state, or request validation results in credential records.
- Treat current Protobuf `Session` fields as credential proof.
- Put credential proof into WebSocket handshake headers, cookies, subprotocols, query parameters, routes, request IDs, player IDs, session IDs, or connection IDs.
- Add password hashing, OAuth, OIDC, provider SDK, JWT, key-management, Redis-like, or other major authentication dependencies from this standard alone.
- Copy Nakama or Pitaya public API shapes.

## 14. Required Future Gates

Before credential storage can be implemented, future work must complete these gates:

```yaml
credential_record_schema_boundary: completed_by_W_0074
token_verifier_record_schema_boundary: required_before_authentication_implementation
authentication_schema_migration_queue: required_before_migration
credential_migration_source: separate_future_work
credential_repository_interface: separate_future_work
credential_postgres_adapter: separate_future_work
redaction_tests: separate_future_or_implementation_work
runtime_authentication_wiring: separate_future_implementation_milestone
```

The migration gate must prove that the SQL source creates only ratified credential structures and does not modify player account lifecycle tables.

## 15. Reference Alignment

### Nakama

Nakama demonstrates that a mature game backend can support device-style authentication, automatic account creation posture, session tokens, refresh, logout, and multiple linked authentication methods.

vibit adapts only the low-friction device-login capability for the first login method. It rejects raw public device identifiers as credential proof, defers multiple linked authentication methods, and keeps session-token and refresh-token behavior outside the credential record.

### Pitaya

Pitaya demonstrates useful session and handler vocabulary: connection acceptors, sessions accessible in handler context, ID binding, frontend/backend session differences, and route-aware message handling.

vibit keeps those ideas separate from credential storage. Credential records are persistent authentication verifier records, not transport sessions, handler context state, or connection binding records.

## 16. Verification

Default verification for this standard:

```bash
node tools/vibit inspect next --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check architecture --json
node tools/vibit check memory --json
node tools/vibit check change define-credential-record-schema-boundary --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Live PostgreSQL verification is not required because this work does not add migration source, tables, repositories, adapters, or runtime behavior.

Go tests are not required because no Go runtime behavior changes.

## 17. Follow-Up

Next work:

```text
W-0075 Define token verifier record schema boundary
```

The token verifier record boundary must define opaque access-token verifier semantics, statuses, expiration, revocation, credential-token linkage, retention, cleanup, and replay-sensitive failure classes before migration planning begins.
