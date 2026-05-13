# PostgreSQL Persistence Boundary Standard

Status: Draft v0.1  
Last updated: 2026-05-13  
Scope: PostgreSQL repository, transaction, migration, and event-recording boundaries for the first durable Go runtime

This standard defines the implementation boundary that must exist before vibit adds PostgreSQL-backed module state.

Use this standard together with `.arch/runtime.yaml`, `.arch/dependencies.yaml`, `modules/inventory/module.yaml`, `ADR-0011`, `ADR-0013`, `ADR-0014`, and `ADR-0020`.

## 1. Purpose

Durable persistence is where long-lived server projects often lose architectural clarity.

The main risk is not that a PostgreSQL repository cannot be written. The main risk is that agents hide transactions, SQL ownership, migration semantics, permission checks, event recording, or cross-module data access inside the most convenient package.

This standard prevents that drift by defining ownership before persistent inventory behavior is implemented.

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
- Provide migration status, apply, rollback, and validation helpers when tooling is implemented.
- Keep migration execution operationally explicit.

Must not:

- Hide domain business behavior in Go migrations.
- Allow migrations to be generated or applied without source files under `runtime/migrations/postgres/`.

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

Must not:

- Let repositories silently open independent write transactions for command flows.
- Let domain modules access transaction handles directly.
- Treat a WebSocket connection as a transaction or session authority.

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

## 5. Migration Rules

SQL migrations are source artifacts.

Rules:

- Migrations must be reviewed as contract-bearing persistence changes.
- Migrations must not be edited after they are treated as applied in a shared environment; add a new migration instead.
- Every migration that creates module-owned state must name the owning module.
- Every migration that enforces an invariant should map back to a module invariant or persistence boundary rule.
- Destructive migrations require a change spec with rollback and data compatibility notes.
- Migration validation must be recorded before a persistent repository is considered verified.

## 6. Test And Verification Rules

Persistence work should add tests at the level that owns the behavior:

- Domain tests cover invariants and permission behavior without PostgreSQL.
- PostgreSQL adapter tests cover row mapping, atomic mutations, constraint handling, and concurrency-sensitive behavior.
- Migration tests or checks cover SQL file naming, `goose` markers, and apply/rollback validation.
- Runtime wiring tests cover configuration and composition, not SQL behavior.

Until a disposable PostgreSQL test environment exists, PostgreSQL integration tests may be gated by an explicit environment variable. Agents must record when those tests were skipped and why.

The current PostgreSQL adapter has focused fake-executor tests for SQL shape and transaction-bound behavior. Live repository integration tests are not mandatory until vibit defines a disposable PostgreSQL test environment standard.

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

It validates SQL migration naming, `goose` Up/Down markers, absence of unapproved Go migrations, owning-module traces, architecture manifest references, and the first inventory table references. It does not apply or roll back migrations against PostgreSQL yet.

Future persistence verification should add:

```bash
goose status
goose up
goose down
PostgreSQL repository integration tests
```

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
- Use database constraints as the only visible source of business rules.
- Add object storage to inventory persistence without a concrete large-object use case and dependency adoption record.
- Make MinIO mandatory for the durable inventory runtime.
