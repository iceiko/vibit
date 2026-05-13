# Impact Analysis

## Affected Modules

- `runtime`: adds the PostgreSQL persistence boundary standard and records runtime ownership rules for adapters, migrations, transaction boundaries, and verification.
- `inventory`: records that the first durable inventory implementation must use PostgreSQL through vibit-owned repository interfaces and must protect capacity-sensitive grants with an inventory account row lock.

## Module Ownership Impact

No data ownership changes are implemented in this change.

The boundary confirms existing ownership:

- Inventory owns inventory records, inventory items, and item grant semantics.
- PostgreSQL adapters implement module repository interfaces.
- Migration source files live under `runtime/migrations/postgres/`.
- Transaction boundary interfaces live under `runtime/internal/platform/tx/`.

## Public Contract Impact

No public command, query, event, error, permission, Protobuf schema, or WebSocket protocol changes.

## Data And Migration Impact

No migration files are added.

The future migration direction is defined: SQL-first `goose` migrations with deterministic sequence numbers and both Up and Down sections.

## Test Impact

No runtime tests are added because this is a boundary and standards change.

Future persistence work must add domain tests, PostgreSQL adapter tests, migration checks, and integration verification when a disposable PostgreSQL environment exists.

## Documentation Impact

Adds the English and Simplified Chinese PostgreSQL persistence boundary standard, adds `ADR-0020`, and updates runtime, inventory, and work-queue guidance.

## Compatibility Risks

Low. This change defines implementation constraints before persistent runtime behavior exists.
