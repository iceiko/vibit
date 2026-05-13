# Impact Analysis

## Affected Modules

- `inventory` remains the first proof slice, but no inventory business behavior is implemented in this change.
- `runtime/internal/app/` gains the first application dispatch skeleton.
- `tools/vibit` gains an additional runtime layer-boundary check.

## Module Ownership Impact

No module ownership changes.

Application dispatch owns route registration and command/query handler invocation. Domain modules will later own business behavior and invariants behind vibit-owned handler interfaces.

## Public Contract Impact

No command, query, event, error, or permission contract changes.

This change uses the existing `kind`, `module`, and `name` route model but does not alter public protocol schemas.

## Data And Migration Impact

No data ownership changes.

No migrations are added.

## Test Impact

Adds Go unit tests for:

- Route registration.
- Duplicate route rejection.
- Unknown route rejection.
- Unsupported message kind rejection.
- Request correlation metadata preservation.
- Handler error propagation.

## Documentation Impact

Updates runtime-facing documentation and translations to say that a dispatch skeleton now exists while WebSocket transport, PostgreSQL persistence, migrations, and inventory business logic remain unimplemented.

## Compatibility Risks

The main risk is committing too much handler shape before persistence and generated route registration exist.

That risk is bounded by keeping the dispatcher generic, platform-free, and small. The dispatcher registers explicit `RouteKey` values and calls vibit-owned handler interfaces without binding to generated Protobuf payload types.
