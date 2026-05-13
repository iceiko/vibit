# Impact

## Affected Modules

- `runtime`: adds PostgreSQL platform configuration and transaction execution plumbing.

## Module Ownership Impact

No domain ownership changes.

The implementation stays within the accepted ownership model:

- `runtime/internal/platform/persistence/postgres/` owns `pgx`, PostgreSQL config parsing, connection pool construction, and pgx-backed transaction execution.
- `runtime/internal/platform/tx/` continues to own the driver-neutral `Runner` and `UnitOfWork` interfaces.
- `runtime/internal/app/` continues to orchestrate transaction use through `TransactionalDispatcher` and does not import PostgreSQL adapters.
- Domain modules continue to depend on module-owned repository interfaces, not pgx.

## Public Contract Impact

No command, query, event, permission, error, Protobuf, or WebSocket contract changes.

## Data And Migration Impact

No migration changes.

The first inventory migration remains source-checked only; live apply/rollback tooling is planned for a later work item.

## Test Impact

Add focused Go tests for:

- PostgreSQL config parsing.
- Environment config behavior.
- Transaction commit on success.
- Transaction rollback on handler error.
- Rollback after commit failure.
- Nil dependency validation.

These tests do not require a live PostgreSQL server.

## Documentation Impact

Update architecture manifests and persistence guidance to record that configuration and the pgx-backed runner now exist, while process wiring and live integration tests remain deferred.

## Compatibility Risks

The main risk is leaking pgx into the application or domain layer. The design avoids that by keeping pgx-owned executable handles inside the PostgreSQL platform package and returning only `tx.UnitOfWork` to application orchestration.
