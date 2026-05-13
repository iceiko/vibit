# Impact Analysis

## Affected Modules

- `runtime`: adds PostgreSQL migration execution helpers under the migration platform owner package.

## Module Ownership Impact

No domain ownership changes.

`runtime/internal/platform/migrations/` becomes the only package that imports `github.com/pressly/goose/v3`. SQL source ownership remains under `runtime/migrations/postgres/`.

## Public Contract Impact

No command, query, event, permission, error, Protobuf, or WebSocket contract changes.

## Data And Migration Impact

No new migration files are added.

This change makes existing SQL migration sources executable through a caller-supplied database handle, but it does not apply migrations during default tests and does not change the first inventory migration.

## Test Impact

Add focused Go tests for option validation and live PostgreSQL gating behavior. Live apply/status execution remains deferred until the disposable PostgreSQL verification environment exists.

## Documentation Impact

Update architecture manifests and persistence guidance to record that migration apply/status tooling exists, while live database verification remains deferred.

## Compatibility Risks

The main risk is creating hidden startup side effects. This change avoids that by exposing explicit helper functions only and leaving startup wiring untouched.
