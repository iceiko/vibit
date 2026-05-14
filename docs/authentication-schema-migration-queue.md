# Authentication Schema Migration Queue

Status: Draft v0.2
Last updated: 2026-05-14
Scope: Completed authentication schema migration, repository, adapter-boundary, redaction, and verification queue after credential and token verifier schema boundaries
Depends on: `docs/credential-record-schema-boundary.md`, `docs/token-verifier-record-schema-boundary.md`, `docs/postgresql-persistence-boundary.md`, `docs/postgresql-verification-environment.md`
Canonical decision: `ADR-0034`

The paired Simplified Chinese translation is `docs/authentication-schema-migration-queue.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This standard defines the next bounded work queue after both required authentication schema boundaries are ratified:

```yaml
credential_record_schema_boundary: migration_source_added
credential_migration_source_added: true
token_verifier_record_schema_boundary: migration_source_added
token_verifier_migration_source_added: true
default_durable_target: PostgreSQL
implementation_authorized_now: false
milestone_status: completed
next_milestone: M-015
```

The goal is to make the future authentication storage path deterministic without jumping directly from schema ratification into broad authentication implementation.

This document records migration order, repository-interface gates, PostgreSQL adapter gates, redaction checks, live PostgreSQL verification expectations, and milestone closeout. The `M-014` queue has now completed migration sources, static checks, the storage-neutral repository interface, and the PostgreSQL adapter boundary. It does not authorize runtime credential lookup, token issuance, token validation, logout execution, refresh, cleanup jobs, generated authentication shapes, Protobuf messages, WebSocket proof carriers, WebSocket handshake authentication, authentication dependencies, or production authentication behavior.

## 2. Required Reading

Read this standard together with:

- `docs/credential-record-schema-boundary.md`
- `docs/token-verifier-record-schema-boundary.md`
- `docs/credential-token-session-schema-gates.md`
- `docs/selected-login-token-boundary-checks.md`
- `docs/token-lifecycle-storage-implications.md`
- `docs/postgresql-persistence-boundary.md`
- `docs/postgresql-verification-environment.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `ADR-0029`
- `ADR-0032`
- `ADR-0033`
- `ADR-0034`

Reference reading:

- Nakama authentication concepts: `https://heroiclabs.com/docs/nakama/concepts/authentication/`
- Nakama session concepts: `https://heroiclabs.com/docs/nakama/concepts/session/`
- Pitaya features and session vocabulary: `https://pitaya.readthedocs.io/en/stable/features.html`
- Pitaya API and handler vocabulary: `https://pitaya.readthedocs.io/en/latest/API.html`

Nakama and Pitaya remain references for game-backend capability coverage and vocabulary only. They do not govern vibit's migration order, repository shape, public API, or transport behavior.

## 3. Queue Rule

Authentication persistence must move through traceable queue steps:

```yaml
schema_boundary:
  credential: completed_by_W_0074
  token_verifier: completed_by_W_0075
migration_queue_planning:
  completed_by: W_0076
migration_sources:
  credential: completed_by_W_0077
  token_verifier: completed_by_W_0078
schema_static_checks:
  completed_by_W_0079
repository_interfaces:
  completed_by_W_0080
postgres_adapters:
  boundary_defined_by_W_0081
redaction_and_live_verification:
  separate_future_work_after_adapter_implementation
milestone_closeout:
  completed_by_W_0082
runtime_authentication:
  blocked_until_later_milestone
```

Agents must not combine migration source creation, repository interface design, adapter implementation, runtime validation, protocol behavior, and authentication behavior in one broad change.

## 4. Planned Work Items

The planned M-014 queue after this planning step is:

| Work item | Title | Scope |
| --- | --- | --- |
| `W-0077` | Add credential PostgreSQL migration source | Completed. Creates only the ratified `authentication_device_credentials` SQL source. |
| `W-0078` | Add token verifier PostgreSQL migration source | Completed. Creates only the ratified `authentication_access_tokens` SQL source after the credential migration exists. |
| `W-0079` | Add authentication migration static checks | Completed. Hardened repository checks for authentication migration naming, ownership, forbidden raw secret columns, and player lifecycle table separation. |
| `W-0080` | Define authentication repository interface boundary | Completed. Defines storage-neutral interfaces and mutations for credential and token verifier storage without adapters or runtime behavior. |
| `W-0081` | Define authentication PostgreSQL adapter boundary | Completed. Reserved adapter source paths, transaction expectations, SQL operation scope, and focused tests without implementing the adapter. |
| `W-0082` | Close credential and token verifier schema ratification milestone | Completed. Reviewed M-014 against exit criteria and opened `M-015` as the bounded authentication PostgreSQL adapter implementation gate. |

This queue intentionally does not include runtime login, token issuance, token validation, logout execution, cleanup jobs, Protobuf, WebSocket, or dependency implementation work.

## 5. Migration Order

The first authentication migration order is:

```text
000003_create_authentication_device_credentials.sql
000004_create_authentication_access_tokens.sql
```

Rationale:

- Credential records must exist before token verifier records can reference `credential_record_id`.
- Player account lifecycle tables already exist in `000002_create_player_account_state.sql`.
- Token verifier records depend on the credential schema boundary for rotation and previous-token revocation linkage.
- Keeping credential and token verifier migrations separate makes each SQL source smaller and easier for agents to validate.

Future migration work may adjust exact sequence numbers only if existing migration files make these numbers unavailable. If numbers change, the change spec must explain the actual sequence and keep credential migration before token verifier migration.

## 6. Credential Migration Gate

`W-0077` has added only one SQL-first migration source:

```text
runtime/migrations/postgres/000003_create_authentication_device_credentials.sql
```

Allowed scope:

- Create `authentication_device_credentials`.
- Add indexes and constraints required by `docs/credential-record-schema-boundary.md`.
- Add `-- +goose Up`, `-- +goose Down`, and `-- Module: runtime.authentication`.
- Preserve player account lifecycle tables without alteration.

Forbidden scope:

- Token verifier table creation.
- Repository interfaces or adapters.
- Runtime credential lookup.
- Login behavior.
- Password hashing, OAuth, OIDC, provider SDK, JWT, key-management, or Redis-like dependencies.
- Generic metadata columns that can hide credentials, tokens, provider payloads, or transport state.

## 7. Token Verifier Migration Gate

`W-0078` has added only one SQL-first migration source:

```text
runtime/migrations/postgres/000004_create_authentication_access_tokens.sql
```

Allowed scope:

- Create `authentication_access_tokens`.
- Add indexes and constraints required by `docs/token-verifier-record-schema-boundary.md`.
- Add `-- +goose Up`, `-- +goose Down`, and `-- Module: runtime.authentication`.
- Preserve player account lifecycle tables without alteration.

Forbidden scope:

- Credential table changes unless required to preserve migration ordering and explicitly justified.
- Repository interfaces or adapters.
- Runtime token issuance, validation, logout, refresh, or cleanup.
- JWT, signing, OAuth, OIDC, provider SDK, key-management, password-hashing, or Redis-like dependencies.
- Generic metadata columns that can hide raw tokens, claims, transport state, provider payloads, or private request details.

## 8. Static Check Gate

`W-0079` has updated repository checks so authentication migrations are checkable before adapters exist.

Expected checks:

- Authentication migration files use deterministic sequence numbers.
- Authentication migration files include goose Up and Down sections.
- Authentication migration files include `-- Module: runtime.authentication`.
- `authentication_device_credentials` and `authentication_access_tokens` are the only authentication tables introduced by M-014 migration work.
- Raw credential and raw token columns are forbidden.
- `player_accounts` and `player_account_events` remain credential-free and token-free.
- JSON metadata columns remain absent unless a later ADR explicitly ratifies them.
- Repository-relative JSON check output uses forward slashes on every platform.

This check gate must remain static and local by default.

## 9. Repository Interface Gate

`W-0080` has defined storage-neutral interfaces after the migration sources and static checks exist.

Expected ownership:

```text
runtime/internal/modules/authentication/
```

Expected interface families:

- Credential lookup and mutation repository boundary.
- Token verifier creation, lookup, revocation, expiration query, and cleanup eligibility boundary.
- Unit-of-work expectations that keep credential, player account, and token mutations atomic when implementation later authorizes them.

Current approved source:

```text
runtime/internal/modules/authentication/repository.go
```

Forbidden scope:

- PostgreSQL adapter implementation.
- Runtime login handlers.
- Token issuance or validation.
- Logout behavior.
- Cleanup jobs.
- WebSocket routes or proof carriers.
- Protobuf messages.
- Authentication dependencies.

## 10. PostgreSQL Adapter Boundary Gate

`W-0081` has defined adapter boundaries after storage-neutral interfaces exist.

Expected ownership:

```text
runtime/internal/platform/persistence/postgres/
```

Expected boundary content:

- Reserved source and test paths.
- Constructor names.
- Caller-owned executor and transaction rules.
- SQL operation scope.
- Error mapping expectations.
- Fake-executor test expectations.
- Optional live PostgreSQL verification expectations.

Forbidden scope:

- Runtime authentication implementation.
- WebSocket route or handshake behavior.
- Protobuf behavior.
- Token generation or cryptographic verifier implementation.
- External provider or Redis-like dependencies.

## 11. Verification Expectations

Default verification remains local:

```bash
node tools/vibit check migrations --json
node tools/vibit check runtime --json
node tools/vibit check all --json
```

Live PostgreSQL verification remains opt-in through:

```text
VIBIT_POSTGRES_TEST_DSN
VIBIT_POSTGRES_TEST_ALLOW_DESTRUCTIVE=1
```

Future migration and adapter work should record one of:

```text
Verified: live PostgreSQL verification ran against VIBIT_POSTGRES_TEST_DSN
Not verified: live PostgreSQL verification skipped because VIBIT_POSTGRES_TEST_DSN was not set
Not applicable: no migration, table, repository, adapter, or persistence behavior changed
```

No W-0076 work requires live PostgreSQL verification because this planning step adds no schema or runtime behavior.

## 12. Redaction Expectations

Future authentication migration, repository, adapter, and test work must preserve redaction.

Forbidden in logs, errors, traces, fixtures, ADRs, change specs, and conversation logs:

- Raw credential proof.
- Raw access tokens.
- Credential lookup or verifier digests.
- Token lookup or verifier digests.
- Server-side verifier secrets or peppers.
- Authorization headers.
- Cookie contents.
- WebSocket subprotocol token carriers.
- URL query token carriers.
- Provider secrets.

Allowed with care:

- `credential_record_id`
- `token_record_id`
- `player_id`
- lifecycle state names
- non-secret reason catalog values
- stable failure class names

## 13. Non-Authorization

This standard records the completed `M-014` queue. It authorized only the bounded artifacts named by `W-0077` through `W-0082`: credential and token verifier migration sources, static checks, a storage-neutral repository interface, a PostgreSQL adapter boundary, and milestone closeout.

It still does not authorize:

- Additional authentication migration sources beyond the ratified credential and token verifier migrations.
- External identity, runtime session, audit, or refresh-token tables.
- PostgreSQL authentication adapter implementation outside the separately opened `M-015` gate.
- Runtime credential lookup behavior.
- Token generation, issuance, parsing, validation, refresh, rotation, replay handling, cleanup jobs, or production token behavior.
- Login handlers.
- Runtime player handlers.
- WebSocket routes.
- Protobuf messages or generated Protobuf output.
- WebSocket proof carriers.
- WebSocket handshake authentication.
- Password hashing, JWT, OAuth, OIDC, provider SDK, Redis-like, cryptography, key-management, or major authentication dependencies.
- Treating metadata-only `player_id`, `session_id`, `connection_id`, or `connection_epoch` as proof.

## 14. Follow-Up

Next work:

```text
W-0083 Refine authentication adapter implementation checks
```

It may refine repository checks so the declared PostgreSQL adapter boundary can be implemented later without weakening runtime authentication, protocol, WebSocket, generated-output, and dependency deferrals.
