# PostgreSQL Persistence Boundary Standard

Status: Draft v0.2
Last updated: 2026-05-14
Scope: PostgreSQL repository, transaction, migration, and event-recording boundaries for the first durable Go runtime

This standard defines the implementation boundary that must exist before vibit adds PostgreSQL-backed module state.

Use this standard together with `.arch/runtime.yaml`, `.arch/dependencies.yaml`, `modules/inventory/module.yaml`, `modules/player/module.yaml`, `ADR-0011`, `ADR-0013`, `ADR-0014`, `ADR-0020`, `ADR-0021`, and `ADR-0022`.

## 1. Purpose

Durable persistence is where long-lived server projects often lose architectural clarity.

The main risk is not that a PostgreSQL repository cannot be written. The main risk is that agents hide transactions, SQL ownership, migration semantics, permission checks, event recording, or cross-module data access inside the most convenient package.

This standard prevents that drift by defining ownership before persistent inventory behavior is implemented and before later player account persistence adapters are added.

## 2. Layer Ownership

### Domain Module

Owner:

```text
runtime/internal/modules/<module>/
```

Responsibilities:

- Own module repository interfaces.
- Own module invariants and domain validation.
- Own domain request, response, event, permission, and error semantics.
- Depend only on vibit-owned interfaces and standard-library packages.

Must not:

- Import `github.com/jackc/pgx/v5`.
- Import `github.com/pressly/goose/v3`.
- Import S3 SDKs or MinIO clients.
- Construct SQL strings.
- Manage database transactions directly.
- Reach into another module's tables.

### PostgreSQL Adapter

Owner:

```text
runtime/internal/platform/persistence/postgres/
```

Responsibilities:

- Own `github.com/jackc/pgx/v5`.
- Implement module repository interfaces using PostgreSQL.
- Convert PostgreSQL rows into module-owned runtime structs.
- Convert PostgreSQL errors into stable module or application errors at adapter boundaries.
- Use transaction handles supplied by the application transaction boundary.
- Own PostgreSQL runtime configuration parsing and pgx connection-pool construction.
- Own the pgx-backed transaction runner adapter for the vibit-owned unit-of-work boundary.
- Keep SQL behavior explicit and covered by focused tests.

Must not:

- Enforce module permissions.
- Replace module invariant checks with hidden database-only logic.
- Import WebSocket or Protobuf platform packages.
- Expose `pgx` types to domain modules.

### Migration Tooling

Owner:

```text
runtime/internal/platform/migrations/
```

Responsibilities:

- Own `github.com/pressly/goose/v3` invocation and validation.
- Provide migration status, apply, and validation helpers. Rollback helpers remain future explicit work until a change requires them.
- Keep migration execution operationally explicit.

Must not:

- Hide domain business behavior in Go migrations.
- Allow migrations to be generated or applied without source files under `runtime/migrations/postgres/`.

The first explicit migration runner lives under:

```text
runtime/internal/platform/migrations/postgres.go
```

It exposes a vibit-owned `PostgresRunner` that accepts a caller-supplied `*sql.DB` and migration source filesystem or directory, configures goose for PostgreSQL, disables the global Go migration registry, lists known migration sources, reports structured migration status, and applies pending SQL-first migrations. It intentionally does not own database connection construction, does not close caller-owned database handles, and does not run during normal `cmd/vibit-server` startup.

The PostgreSQL platform owner may provide narrow helpers that adapt a package-owned pool into `*sql.DB` for migration tooling, such as:

```text
runtime/internal/platform/persistence/postgres/sql_db.go
```

Those helpers must keep pgx stdlib imports inside the PostgreSQL platform owner package. They do not move migration execution ownership out of `runtime/internal/platform/migrations/`.

### Migration Sources

Owner:

```text
runtime/migrations/postgres/
```

Responsibilities:

- Store SQL-first PostgreSQL migration source files.
- Use deterministic sequence numbers.
- Include both `-- +goose Up` and `-- +goose Down` sections.
- Declare module-owned tables, indexes, constraints, and foreign-key-like ownership references explicitly.

Naming pattern:

```text
runtime/migrations/postgres/000001_create_inventory_state.sql
```

Go migrations are allowed only when a change spec explains why SQL is insufficient. They must not contain domain business logic.

### Transaction Boundary

Owner:

```text
runtime/internal/platform/tx/
```

Orchestrator:

```text
runtime/internal/app/
```

Responsibilities:

- State-changing commands run inside an application-owned unit of work.
- Repository mutations and durable event recording happen in the same unit of work.
- Query handlers do not mutate state and do not require a write transaction by default.
- Transaction retry policy is platform-owned and must be explicit before use.
- Event publication outside the transaction remains deferred until an event delivery or outbox standard exists.

The first Go skeleton for this boundary is:

```text
runtime/internal/platform/tx.Runner
runtime/internal/platform/tx.UnitOfWork
runtime/internal/app.TransactionalDispatcher
```

`TransactionalDispatcher` opens a unit of work for command routes and passes query routes through without a write unit of work by default. The skeleton is intentionally driver-neutral; it must not expose PostgreSQL driver handles to application handlers or domain modules.

The first pgx-backed implementation of this boundary lives under:

```text
runtime/internal/platform/persistence/postgres/runner.go
```

It implements `runtime/internal/platform/tx.Runner`, begins a `pgx` transaction inside the PostgreSQL platform package, commits when the application callback succeeds, rolls back when the callback fails, and attempts rollback after commit failure. The unit of work it passes to the callback still satisfies the vibit-owned `tx.UnitOfWork` interface. PostgreSQL-specific repository factories, such as `UnitOfWork.NewInventoryRepository`, must remain package-owned helpers and must not force application or domain packages to import `pgx`.

Must not:

- Let repositories silently open independent write transactions for command flows.
- Let domain modules access transaction handles directly.
- Treat a WebSocket connection as a transaction or session authority.

### PostgreSQL Runtime Configuration

Owner:

```text
runtime/internal/platform/persistence/postgres/
```

The first PostgreSQL configuration source is:

```text
runtime/internal/platform/persistence/postgres/config.go
```

It reads explicit process input through these environment variables:

```text
VIBIT_POSTGRES_DSN
VIBIT_POSTGRES_MAX_CONNS
VIBIT_POSTGRES_MIN_CONNS
```

Rules:

- `VIBIT_POSTGRES_DSN` is required when opening a PostgreSQL pool.
- Pool size settings are optional and must be non-negative integers.
- Configuration parsing may build a `pgxpool.Config`, but normal unit tests must not require a live PostgreSQL server.
- Connection strings and credentials must come from environment or explicit runtime input and must not be stored in tracked files.
- Process startup must keep PostgreSQL optional. `VIBIT_RUNTIME_STORE=postgres` explicitly selects the PostgreSQL-backed inventory composition path; the default remains in-memory startup.
- `*sql.DB` adaptation for migration tooling, when needed, must stay inside the PostgreSQL platform owner package.

## 3. Inventory Persistence Boundary

The first durable module state is inventory.

The inventory module owns:

```text
inventory account row
inventory item quantity row
inventory item grant event/audit row
```

The first persistent schema should model:

- One inventory account row per `player_id`.
- One inventory item row per `(player_id, item_id)`.
- One durable grant record per emitted `ItemGranted` event.

The persistent grant flow must be atomic:

```text
application command dispatch
-> open unit of work
-> lock inventory account for player_id
-> read current inventory
-> enforce permission and capacity policies
-> upsert item quantity
-> record ItemGranted grant row
-> commit unit of work
-> return application result
```

The lock is required because inventory capacity depends on the player's current set of items. The first PostgreSQL implementation should use an explicit inventory account row lock, not an implicit best-effort read.

The module-owned Go boundary for this flow is:

```text
Repository.LockInventoryForMutation(ctx, player_id) -> MutationLock
MutationLock.GetInventory(ctx, player_id)
MutationLock.GrantItem(ctx, GrantItemMutation)
MutationLock.Release()
```

`MutationLock` is a locked repository view for one inventory aggregate. It is not a transaction owner. `Release` only releases the aggregate lock or adapter-local resource; it must not commit or roll back the application-owned unit of work.

`GrantItemMutation` must carry the durable grant record metadata needed by PostgreSQL adapters:

```text
event_id
occurred_at
player_id
item_id
quantity
reason
```

The domain handler must create this metadata before calling `MutationLock.GrantItem`, so the adapter can record the item quantity change and the `inventory_item_grants` row with the same application-owned executor.

## 3.1 Player Account Persistence Schema Boundary

Player account persistence is ratified as account lifecycle state only.

The player module owns planned PostgreSQL state for:

```text
player account lifecycle row
player account lifecycle event/audit row
```

The first player account persistent schema must model:

- One `player_accounts` row per stable `player_id`.
- One `player_account_events` row per durable player account lifecycle fact.
- `PlayerAccountCreated` as the first lifecycle event type that must be recordable.

The planned `player_accounts` table owns these columns:

```text
player_id
display_name
account_state
created_at
updated_at
disabled_at
deleted_at
```

Column rules:

- `player_id` is the stable primary key and must be non-blank text.
- `display_name` must be non-blank text. Uniqueness remains deferred.
- `account_state` must be constrained to `active`, `disabled`, or `deleted`.
- `created_at` and `updated_at` are required timestamps.
- `disabled_at` is nullable and may be set only for disabled or deleted accounts.
- `deleted_at` is nullable and may be set only for deleted accounts.
- The first migration should use explicit check constraints for non-blank text and lifecycle state values.

The planned `player_account_events` table owns these columns:

```text
event_id
event_type
occurred_at
player_id
requested_by
account_state
display_name
metadata
recorded_at
```

Column rules:

- `event_id` is the primary key and must be non-blank text.
- `event_type` must be non-blank text. The first required value is `PlayerAccountCreated`.
- `occurred_at` records the domain event time.
- `player_id` references `player_accounts(player_id)` and must be non-blank text.
- `requested_by` must be non-blank text but is not authenticated proof by itself.
- `account_state` records the lifecycle state after the event.
- `display_name` records the display name after the event when relevant.
- `metadata` may use `JSONB NOT NULL DEFAULT '{}'::jsonb` for structured non-secret event metadata.
- `recorded_at` records when PostgreSQL stored the row.

The first player account migration should include indexes for:

```text
player_accounts(account_state)
player_account_events(player_id, occurred_at)
player_account_events(event_type, occurred_at)
```

The first player account schema must not store:

- authentication credentials
- password hashes
- authentication provider names or provider subject IDs
- external identity links
- access tokens
- refresh tokens
- token signing metadata
- runtime session rows
- WebSocket connection state
- request identity validation results
- inventory state
- permission grants

Those concerns require separate standards, contracts, migrations, and verification before implementation.

The first persistent player account creation flow, once runtime implementation is separately ratified, should be atomic:

```text
application command dispatch
-> open unit of work
-> enforce request validation and account creation permission
-> insert player_accounts row
-> insert player_account_events row for PlayerAccountCreated
-> commit unit of work
-> return application result
```

The player account repository boundary is module-owned and storage-neutral:

```text
Repository.CreatePlayerAccount(ctx, CreatePlayerAccountMutation)
Repository.GetPlayerAccount(ctx, player_id)
```

The first source for this boundary is:

```text
runtime/internal/modules/player/repository.go
```

It defines the account lifecycle view, the `Repository` interface, and the `CreatePlayerAccountMutation` fields required by the durable event record. It must remain free of PostgreSQL, WebSocket, Protobuf, authentication, token, credential, and session dependencies.

`CreatePlayerAccountMutation` must carry the durable event metadata needed by PostgreSQL adapters:

```text
event_id
occurred_at
player_id
display_name
account_state
requested_by
```

The PostgreSQL adapter must use an executor supplied by the application-owned unit of work. It must not open hidden write transactions, parse credentials, validate tokens, bind sessions, decode Protobuf payloads, or enforce WebSocket behavior.

The first player account migration source is:

```text
runtime/migrations/postgres/000002_create_player_account_state.sql
```

It uses the next deterministic SQL migration number after existing migration files and includes:

```text
-- +goose Up
-- Module: player
-- +goose Down
```

Runtime player account handlers, WebSocket route wiring, authentication, token behavior, credential storage, external identity linking, and session persistence remain deferred until separately ratified.

The first credential migration source is:

```text
runtime/migrations/postgres/000003_create_authentication_device_credentials.sql
```

It is owned by `runtime.authentication`, creates only `authentication_device_credentials`, references `player_accounts(player_id)` without changing player lifecycle tables, and does not authorize token storage, repository interfaces, PostgreSQL adapters, runtime authentication, Protobuf messages, or WebSocket behavior.

The first token verifier migration source is:

```text
runtime/migrations/postgres/000004_create_authentication_access_tokens.sql
```

It is owned by `runtime.authentication`, creates only `authentication_access_tokens`, references `player_accounts(player_id)` and `authentication_device_credentials(credential_record_id)` without changing those tables, and does not authorize repository interfaces, PostgreSQL adapters, runtime authentication, Protobuf messages, or WebSocket behavior.

## 3.2 Player Account PostgreSQL Adapter Boundary

The first player account PostgreSQL adapter is implemented, but it remains a persistence adapter only. It does not authorize runtime handlers, WebSocket routes, authentication, token behavior, credential storage, external identity linking, or session persistence.

The adapter source path is:

```text
runtime/internal/platform/persistence/postgres/player_account_repository.go
```

The focused test path is:

```text
runtime/internal/platform/persistence/postgres/player_account_repository_test.go
```

The adapter is constructed with:

```text
NewPlayerAccountRepositoryForUnitOfWork(executor)
```

The constructor returns an implementation of:

```text
player.Repository
```

The executor must be supplied by the application-owned unit of work, normally from a transaction-bound handle such as `pgx.Tx`. The adapter uses the same small pgx-shaped executor interface as the inventory adapter for testability. It must not call `BEGIN`, `COMMIT`, or `ROLLBACK`; transaction lifetime belongs to the application-owned unit of work.

The adapter must not:

- Call `BEGIN`, `COMMIT`, or `ROLLBACK`.
- Open its own hidden write transaction.
- Open a PostgreSQL pool or read PostgreSQL configuration.
- Apply migrations or inspect migration status.
- Parse authentication credentials.
- Validate tokens.
- Bind or persist runtime sessions.
- Decode Protobuf payloads.
- Know WebSocket handshake or connection behavior.
- Enforce permissions.
- Import transport, protocol, application bootstrap, authentication, credential, token, session, inventory, S3, or MinIO packages.

The first allowed SQL operation scope is intentionally narrow:

- `CreatePlayerAccount` normalizes `player.CreatePlayerAccountMutation`.
- `CreatePlayerAccount` inserts one `player_accounts` row.
- `CreatePlayerAccount` inserts one `player_account_events` row for `PlayerAccountCreated` in the same caller-supplied executor path.
- `GetPlayerAccount` normalizes `player_id`.
- `GetPlayerAccount` reads the current account lifecycle row from `player_accounts`.
- Neither method reads or writes credentials, tokens, external identity links, runtime sessions, WebSocket state, inventory state, or permission grants.

Error mapping expectations:

- Missing rows must map to a stable player account not-found error path before runtime handlers expose the error to clients.
- Duplicate `player_id` or duplicate `event_id` constraint violations must map to stable duplicate/conflict error paths.
- Check constraint violations must map to stable invariant or validation error paths.
- Unexpected PostgreSQL errors may be wrapped with adapter context, but `pgx` or `pgconn` types must not leak into the module-owned repository interface.

The first adapter implementation includes focused tests for:

- Fake-executor SQL shape tests for account creation and account lookup.
- A no-transaction-control test proving the adapter does not issue `BEGIN`, `COMMIT`, or `ROLLBACK`.
- Mutation normalization and UTC timestamp tests.
- Row mapping tests for nullable lifecycle timestamps.
- Duplicate account, duplicate event, check constraint, and missing-row error mapping tests.
- Import-boundary tests through `node tools/vibit check runtime`.
- No live PostgreSQL dependency in default repository checks.
- Optional live integration coverage only through `VIBIT_POSTGRES_TEST_DSN`.

The PostgreSQL unit-of-work helper constructs this repository from a transaction executor through:

```text
runtime/internal/platform/persistence/postgres.UnitOfWork.NewPlayerAccountRepository
```

That helper must stay in the PostgreSQL platform package. It must not change the module-owned `player.Repository` interface, force application or domain packages to import `pgx`, or imply runtime player account handlers, WebSocket routes, authentication, credentials, tokens, external identity links, or session persistence.

## 3.3 Authentication PostgreSQL Adapter Boundary

The first authentication PostgreSQL adapter is implemented. It is a persistence adapter only. It does not authorize runtime authentication, token issuance, token validation, logout execution, refresh behavior, cleanup jobs, handlers, WebSocket routes, Protobuf messages, generated authentication shapes, authentication dependencies, or production authentication behavior.

The adapter source path is:

```text
runtime/internal/platform/persistence/postgres/authentication_repository.go
```

The focused test path is:

```text
runtime/internal/platform/persistence/postgres/authentication_repository_test.go
```

The adapter constructor is:

```text
NewAuthenticationRepositoryForUnitOfWork(executor)
```

The constructor returns an implementation of:

```text
authentication.Repository
```

The executor must be supplied by the application-owned unit of work, normally from a transaction-bound handle such as `pgx.Tx`. The adapter may use a small pgx-shaped executor interface for testability. It must not call `BEGIN`, `COMMIT`, or `ROLLBACK`; transaction lifetime belongs to the application-owned unit of work.

The adapter must not:

- Call `BEGIN`, `COMMIT`, or `ROLLBACK`.
- Open its own hidden write transaction.
- Open a PostgreSQL pool or read PostgreSQL configuration.
- Apply migrations or inspect migration status.
- Generate credential material or tokens.
- Compare credential or token verifiers.
- Parse bearer tokens.
- Validate access tokens.
- Execute login, logout, refresh, or cleanup behavior.
- Decode Protobuf payloads.
- Know WebSocket handshake or connection behavior.
- Enforce permissions.
- Import transport, protocol, application bootstrap, player, inventory, S3, or MinIO packages.

The first allowed SQL operation scope is intentionally limited to persistence operations for the already-ratified repository interface:

- `StoreCredential` may insert `authentication_device_credentials` rows.
- `FindCredentialByLookupDigest` may read one `authentication_device_credentials` row by `credential_lookup_digest`.
- `StoreToken` may insert `authentication_access_tokens` rows.
- `FindTokenByLookupDigest` may read one `authentication_access_tokens` row by `token_lookup_digest`.
- `RevokeCredential` may update credential terminal-state columns on `authentication_device_credentials`.
- `RevokeToken` may update token terminal-state and cleanup columns on `authentication_access_tokens`.
- `ListTokensEligibleForCleanup` may read token records whose `cleanup_after` is due.
- The adapter may reference `player_accounts(player_id)` through already-ratified foreign keys, but it must not read or write player account lifecycle columns.
- The adapter must not read or write external identity, runtime session, WebSocket state, inventory state, or audit persistence tables.

Error mapping expectations:

- Missing credential lookup rows must map to a stable credential not-found error path before runtime handlers expose the error to clients.
- Missing token lookup rows must map to a stable token not-found error path before runtime handlers expose the error to clients.
- Duplicate credential or token digest constraint violations must map to stable duplicate/conflict error paths.
- Foreign-key violations for `player_id`, `credential_record_id`, or replacement links must map to stable validation or invariant error paths.
- Check constraint violations must map to stable validation or invariant error paths.
- Unexpected PostgreSQL errors may be wrapped with adapter context, but `pgx` or `pgconn` types must not leak into the module-owned repository interface.

The implementation includes focused tests for:

- Fake-executor SQL shape tests for credential create/lookup/revocation.
- Fake-executor SQL shape tests for token create/lookup/revocation/cleanup eligibility.
- A no-transaction-control test proving the adapter does not issue `BEGIN`, `COMMIT`, or `ROLLBACK`.
- Mutation normalization and UTC timestamp tests.
- Row mapping tests for nullable terminal-state and cleanup timestamps.
- Missing row, duplicate, foreign-key, and check-constraint error mapping tests.
- Import-boundary tests through `node tools/vibit check runtime`.
- No live PostgreSQL dependency in default repository checks.
- Optional live integration coverage only through `VIBIT_POSTGRES_TEST_DSN`.

The PostgreSQL unit-of-work helper is:

```text
runtime/internal/platform/persistence/postgres.UnitOfWork.NewAuthenticationRepository
```

That helper must stay in the PostgreSQL platform package. It must not change the module-owned `authentication.Repository` interface, force application or domain packages to import `pgx`, or imply runtime authentication handlers, WebSocket routes, Protobuf messages, generated authentication shapes, authentication dependencies, or production authentication behavior.

## 4. Repository Rules

Repository interfaces are module-owned. PostgreSQL adapters implement them.

Rules:

- A repository method must not perform permission checks.
- A repository method must not decode Protobuf payloads.
- A repository method must not know about WebSocket sessions.
- A repository method must preserve module-owned invariants that are also expressible as database constraints.
- A persistent command flow must use a transaction-bound repository.
- An in-memory repository may remain for tests and pre-persistence bootstrap, but it is not the authoritative durable store.
- Repositories receive transaction binding from application composition. They must not create their own hidden write transaction for a command flow.

For inventory, `GrantItem` must acquire `LockInventoryForMutation` after request validation and permission checks, and before reading current inventory for capacity enforcement. Capacity-sensitive reads and the grant mutation must go through the returned `MutationLock`.

The PostgreSQL adapter must implement the lock with an explicit inventory account row lock for `player_id` inside the application-owned unit of work. It must not replace this with an advisory lock or hidden repository-owned transaction without a superseding decision.

The first adapter source is:

```text
runtime/internal/platform/persistence/postgres/inventory_repository.go
```

It is constructed with:

```text
NewInventoryRepositoryForUnitOfWork(executor)
```

The executor must be supplied by application composition, normally from a transaction-bound handle such as `pgx.Tx`. The adapter may depend on a small pgx-shaped executor interface for testability, but it must not call `BEGIN`, `COMMIT`, or `ROLLBACK` itself.

The first PostgreSQL unit-of-work helper can construct this repository from a transaction executor through:

```text
runtime/internal/platform/persistence/postgres.UnitOfWork.NewInventoryRepository
```

This helper is intentionally owned by the PostgreSQL platform package. It does not change the module-owned repository interface and does not expose `pgx` to inventory domain code.

The first persistent runtime composition is explicit:

```text
VIBIT_RUNTIME_STORE=postgres
```

That path opens a PostgreSQL pool from `VIBIT_POSTGRES_DSN`, uses `runtime/internal/platform/persistence/postgres.Runner` for command unit-of-work ownership, creates command repositories from `postgres.UnitOfWork.NewInventoryRepository`, and uses a PostgreSQL inventory repository for query routes. The default server path remains `VIBIT_RUNTIME_STORE=memory`.

Normal server startup does not apply migrations. Migration execution stays an explicit operator or tooling action unless a future change spec authorizes automatic startup migration behavior.

## 5. Migration Rules

SQL migrations are source artifacts.

Rules:

- Migrations must be reviewed as contract-bearing persistence changes.
- Migrations must not be edited after they are treated as applied in a shared environment; add a new migration instead.
- Every migration that creates module-owned state must name the owning module.
- Every migration that enforces an invariant should map back to a module invariant or persistence boundary rule.
- Player account migrations may create only the lifecycle tables ratified in Section 3.1 until authentication, credential, token, external identity, or session persistence standards are separately ratified.
- Destructive migrations require a change spec with rollback and data compatibility notes.
- Migration validation must be recorded before a persistent repository is considered verified.
- Migration execution helpers must be invoked explicitly by an operator, tool, or approved runtime composition path. Normal server startup must not apply migrations unless a later change spec authorizes that behavior.

## 6. Test And Verification Rules

Persistence work should add tests at the level that owns the behavior:

- Domain tests cover invariants and permission behavior without PostgreSQL.
- PostgreSQL adapter tests cover row mapping, atomic mutations, constraint handling, and concurrency-sensitive behavior.
- PostgreSQL config tests cover environment parsing and pool config construction without opening a live connection.
- PostgreSQL transaction runner tests cover begin, commit, rollback, and dependency validation using fake pgx transactions when a disposable PostgreSQL environment is not yet defined.
- Migration tests or checks cover SQL file naming, `goose` markers, and apply/rollback validation.
- Runtime wiring tests cover configuration and composition, not SQL behavior.

Until a disposable PostgreSQL test environment exists, PostgreSQL integration tests may be gated by an explicit environment variable. Agents must record when those tests were skipped and why.

The current PostgreSQL adapter has focused fake-executor tests for SQL shape and transaction-bound behavior. The current runtime wiring tests cover store selection and application composition without opening a live PostgreSQL connection by default. Live repository integration tests are opt-in through the disposable PostgreSQL verification environment standard.

The current PostgreSQL configuration and transaction runner have focused tests under:

```text
runtime/internal/platform/persistence/postgres/config_test.go
runtime/internal/platform/persistence/postgres/runner_test.go
runtime/internal/app/bootstrap/inventory_test.go
runtime/cmd/vibit-server/main_test.go
```

The current PostgreSQL migration runner has focused tests under:

```text
runtime/internal/platform/migrations/postgres_test.go
```

These tests validate option handling, SQL source discovery, and error propagation without requiring a live PostgreSQL server.

Current repository verification remains:

```bash
cd runtime && go test ./...
cd runtime && go vet ./...
node tools/vibit check migrations
node tools/vibit check runtime
node tools/vibit check all
```

The current migration source check is:

```bash
node tools/vibit check migrations
```

It validates SQL migration naming, `goose` Up/Down markers, absence of unapproved Go migrations, owning-module traces, architecture manifest references, the first inventory table references, the player account lifecycle migration shape, the credential migration source shape, and the token verifier migration source shape. It does not apply or roll back migrations against PostgreSQL yet.

The current migration apply/status API is:

```text
runtime/internal/platform/migrations.NewPostgresRunner
runtime/internal/platform/migrations.PostgresRunner.Status
runtime/internal/platform/migrations.PostgresRunner.Apply
```

Live migration apply/status verification remains deferred until the disposable PostgreSQL verification environment is defined.

Future persistence verification should add PostgreSQL-backed migration status, apply, rollback, and repository integration checks against that disposable environment.

## 7. Agent Rules

Agents must:

- Read this standard before adding PostgreSQL repositories, migrations, transaction code, or persistent runtime wiring.
- Update module manifests before changing durable data ownership.
- Keep `pgx` inside `runtime/internal/platform/persistence/postgres/`.
- Keep `goose` inside `runtime/internal/platform/migrations/`.
- Keep SQL source files under `runtime/migrations/postgres/`.
- Record skipped PostgreSQL integration verification explicitly.

Agents must not:

- Add PostgreSQL imports to domain modules.
- Add migration side effects to command handlers.
- Add migration side effects to normal server startup without a change spec.
- Use database constraints as the only visible source of business rules.
- Add object storage to inventory persistence without a concrete large-object use case and dependency adoption record.
- Make MinIO mandatory for the durable inventory runtime.
