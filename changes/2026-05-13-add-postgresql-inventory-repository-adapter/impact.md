# Impact Analysis

## Affected Modules

- `runtime`: gains a PostgreSQL persistence adapter package and the accepted `pgx/v5` module dependency.
- `inventory`: its mutation repository contract now carries the durable grant record fields needed by a transaction-bound adapter.

## Module Ownership Impact

No data ownership changes.

Inventory still owns:

- `inventory_accounts`
- `inventory_items`
- `inventory_item_grants`

The PostgreSQL adapter implements inventory-owned repository interfaces but does not own inventory business rules, permissions, or capacity policy.

## Public Contract Impact

No WebSocket or Protobuf wire contract changes.

The internal Go repository mutation contract changes so `GrantItemMutation` carries:

- `event_id`
- `occurred_at`
- `reason`

This lets durable adapters record the `ItemGranted` grant row in the same unit of work as the item quantity mutation.

## Dependency Impact

`github.com/jackc/pgx/v5` becomes a direct Go module dependency inside `runtime/go.mod`.

The selected version is `v5.7.6` because it remains compatible with the current Go 1.23 baseline. The `go` directive is normalized to `1.23.0` because `pgx/v5@v5.7.6` declares `go 1.23.0`.

No new foundational dependency is introduced; `pgx/v5` was already accepted by `ADR-0013`.

## Data And Migration Impact

No new migration is required.

The adapter targets the existing migration source:

```text
runtime/migrations/postgres/000001_create_inventory_state.sql
```

## Test Impact

Focused package tests cover the adapter without requiring a live database by using a fake pgx-shaped executor.

The tests verify:

- inventory row mapping
- account row lock SQL
- absence of `BEGIN`, `COMMIT`, and `ROLLBACK` inside repository methods
- item quantity upsert SQL
- grant record insertion SQL
- required grant record fields
- `Release` does not commit or roll back the unit of work

## Documentation Impact

Update runtime and inventory guidance, architecture manifests, and the PostgreSQL persistence boundary standard to record the adapter, its transaction-bound construction rule, and the current live integration-test gap.

## Compatibility Risks

Medium.

The runtime still defaults to in-memory inventory wiring, so external behavior does not change yet. The internal mutation contract changed to support durable grant recording, so future adapters must provide event metadata before mutating state.
