# Impact

## Affected Areas

- `runtime/internal/platform/protocol/protobuf/`: adds an integration test that exercises the Protobuf frame handler with PostgreSQL-backed inventory composition.
- `runtime/internal/platform/persistence/postgres/`: adds a package-owned helper for exposing `*sql.DB` to migration tooling without leaking pgx stdlib imports outside the PostgreSQL platform owner.
- `docs/postgresql-verification-environment.md` and `docs/runtime-runbook.md`: document the concrete opt-in live verification command.
- `.arch/runtime.yaml`, `.arch/conventions.yaml`, and `modules/inventory/module.yaml`: record the new verification path and current live execution status.

## Boundary Impact

The change keeps all third-party ownership rules intact:

- `pgx` stays under `runtime/internal/platform/persistence/postgres/`.
- `goose` stays under `runtime/internal/platform/migrations/`.
- Protobuf runtime and generated Protobuf imports stay under generated packages and the Protobuf protocol adapter.
- Application and domain packages continue to depend on vibit-owned interfaces.

## Runtime Impact

Default runtime startup and default tests remain unchanged. `VIBIT_RUNTIME_STORE=memory` remains the default. PostgreSQL persistence still requires explicit `VIBIT_RUNTIME_STORE=postgres` and `VIBIT_POSTGRES_DSN` for process startup. Migrations remain explicit and are not applied during normal startup.

## Verification Impact

The new integration test is part of `go test ./...`, but it skips unless `VIBIT_POSTGRES_TEST_DSN` is set. When the DSN is set and destructive cleanup is allowed, it verifies migration status, migration apply, and a persistent grant/read request loop against a real PostgreSQL database.
