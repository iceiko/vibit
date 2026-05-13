# Impact Analysis

## Affected Modules

- `inventory`: gains the first durable PostgreSQL schema source for inventory-owned state.
- `runtime`: gains the first SQL migration source under the accepted migration root.

## Module Ownership Impact

Inventory owns the new tables:

- `inventory_accounts`
- `inventory_items`
- `inventory_item_grants`

The schema references `player_id` and `item_id` as external stable identifiers. It does not claim ownership of player accounts or the item catalog.

## Public Contract Impact

No public command, query, event, error, permission, Protobuf schema, or WebSocket protocol changes.

The migration stores fields already present in the inventory command and event contracts.

## Data And Migration Impact

Adds the first PostgreSQL migration source file:

- `runtime/migrations/postgres/000001_create_inventory_state.sql`

The `Up` migration creates inventory aggregate, quantity, and grant-record tables. The `Down` migration drops them in reverse dependency order.

## Test Impact

No PostgreSQL adapter or migration runner exists yet, so this change cannot run apply/rollback verification.

Repository checks and Go tests should still pass. Migration validation checks are planned for `W-0014`.

## Documentation Impact

Update runtime and inventory manifests/guidance to record that the first migration source exists while PostgreSQL adapter and apply/rollback verification remain pending.

## Compatibility Risks

Low. This is the first migration and no shared environment is defined yet.

The main future risk is editing the migration after it is treated as applied in a shared environment. Future data changes should add new migrations instead.
