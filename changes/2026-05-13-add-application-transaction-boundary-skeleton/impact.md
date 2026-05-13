# Impact Analysis

## Affected Modules

- `runtime`: adds the first transaction boundary interfaces and application dispatch orchestration.

## Module Ownership Impact

No domain data ownership changes.

The transaction boundary remains split as accepted:

- `runtime/internal/platform/tx/` owns unit-of-work interfaces and future platform execution.
- `runtime/internal/app/` owns command/query orchestration.

Domain modules still depend on repository and policy interfaces, not platform transaction handles.

## Public Contract Impact

No public command, query, event, error, permission, Protobuf schema, or WebSocket protocol changes.

The changed interfaces are internal Go runtime interfaces.

## Data And Migration Impact

No migration files are added.

This change prepares the boundary needed before SQL migrations and PostgreSQL repository adapters.

## Test Impact

Application tests should verify:

- Command dispatch runs inside an injected unit of work.
- Query dispatch does not open a write unit of work.
- Unit-of-work errors preserve application result metadata.

Existing runtime request-loop tests should continue to pass.

## Documentation Impact

Update runtime agent guidance, the runtime manifest, the PostgreSQL persistence boundary, and paired Simplified Chinese translations where public guidance changes.

## Compatibility Risks

Low. Existing `app.Dispatcher` behavior stays unchanged; transaction orchestration is added as a wrapper that can be introduced by persistent composition later.

The main architectural risk is allowing `runtime/internal/app/` to import arbitrary platform adapters. The runtime check should allow only the `platform/tx` boundary package.
