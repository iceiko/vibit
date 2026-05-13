# inventory Module Agent Guide

Status: Draft v0.1

## When To Use This Module

Use this module for requirements that change or read player inventory state.

The first proof-slice capability is intentionally narrow:

- Grant an item to a player's inventory through `GrantItem`.
- Read a player's inventory through `GetInventory`.
- Publish `ItemGranted` after a successful grant.
- Enforce inventory permissions, positive quantities, capacity limits, and inventory item identity.

This module owns inventory records and inventory items. It may reference `player_id` and `item_id`, but it does not own player accounts or the item catalog.

## When Not To Use This Module

Do not use this module for:

- Player account lifecycle.
- Item catalog definitions or item balancing.
- Currency balances or purchases.
- Reward claim eligibility.
- Quest progress.
- Match or session lifecycle.

If a requirement needs one of those concepts, create or update the owning module contract instead of adding hidden ownership here.

## Extension Points

- Command handler: `GrantItem`
- Query handler: `GetInventory`
- Published event: `ItemGranted`
- Policies: inventory capacity and inventory permission checks
- Repository: inventory persistence behind the module boundary
- Tests: command, query, event, contract, invariant, and architecture tests

Runtime contract source files live under `contracts/inventory/` and are registered in `.arch/contracts.yaml`.

Generated contract shapes:

- Planned Go contract shapes for `GrantItem`, `GetInventory`, and `ItemGranted`
- Generated Protobuf wire schemas for the first WebSocket client/server inventory messages

Before implementing the first runtime slice, read:

- `contracts/inventory/commands/GrantItem.yaml`
- `contracts/inventory/queries/GetInventory.yaml`
- `contracts/inventory/events/ItemGranted.yaml`
- `contracts/inventory/errors/inventory_errors.yaml`
- `contracts/inventory/permissions/inventory_permissions.yaml`

The first handwritten runtime path now starts under `runtime/internal/modules/inventory/`:

- `GrantItemHandler`
- `GetInventoryHandler`
- Inventory repository interface behind the module boundary
- Inventory capacity and permission policy interfaces
- `RegisterRoutes` for `runtime/internal/app.Dispatcher`
- Tests for command, query, event, permission, capacity, and dispatcher integration behavior

The first inventory Protobuf/domain bridge lives under `runtime/internal/platform/protocol/protobuf/inventory_bridge.go`. It maps generated inventory wire payloads into inventory runtime request structs and maps inventory runtime results/events back into generated Protobuf payloads.

Do not import generated Protobuf types directly into this module. Protocol adapters or generated bridges should translate wire payloads into inventory runtime request structs.

PostgreSQL persistence work must follow `docs/postgresql-persistence-boundary.md`, `docs/postgresql-verification-environment.md`, and `ADR-0020`.

For the first durable implementation:

- Inventory repository interfaces remain owned by this module.
- PostgreSQL adapter code belongs under `runtime/internal/platform/persistence/postgres/`.
- SQL migrations belong under `runtime/migrations/postgres/`.
- The first migration source is `runtime/migrations/postgres/000001_create_inventory_state.sql`.
- `GrantItem` must call `LockInventoryForMutation` after request validation and permission checks, then use the returned `MutationLock` for the current inventory read and grant mutation.
- PostgreSQL adapters must implement that lock as an inventory account row lock for `player_id` inside the application-owned unit of work.
- `MutationLock.Release` releases the aggregate lock or adapter-local resource; it must not commit or roll back a transaction.
- The first PostgreSQL adapter is `runtime/internal/platform/persistence/postgres/inventory_repository.go`, with focused tests in `runtime/internal/platform/persistence/postgres/inventory_repository_test.go`.
- Durable grant behavior must record the item quantity change and the `ItemGranted` grant record inside the same application-owned unit of work. `GrantItemMutation` carries `event_id`, `occurred_at`, and `reason` for this purpose.
- Migration apply/status and repository integration verification require an explicit disposable PostgreSQL environment through `VIBIT_POSTGRES_TEST_DSN`.
- Live PostgreSQL repository integration tests are opt-in and must be skipped with an explicit verification note when `VIBIT_POSTGRES_TEST_DSN` is not set.

## Forbidden Shortcuts

- Do not bypass boundaries declared in `module.yaml`.
- Do not directly modify data owned by another module.
- Do not add unregistered public commands, queries, events, or permissions.
- Do not put inventory business rules in WebSocket, HTTP, Protobuf, or transport handlers.
- Do not make this module depend directly on third-party WebSocket or Protobuf libraries.
- Do not make this module depend directly on PostgreSQL drivers, S3 SDKs, or MinIO clients. Use vibit-owned repository and storage interfaces when persistence implementation begins.
- Do not hide transaction creation inside inventory repositories for command flows.
- Do not read current inventory for capacity-sensitive grants outside the mutation lock.
- Do not drop event metadata from `GrantItemMutation`; durable adapters need it to persist `inventory_item_grants` atomically with the quantity update.
- Do not hand-edit generated files. If generated output is wrong, change the source contract, template, or generator.
- Do not invent payload fields in implementation. Update the relevant contract source file first.
- Do not introduce a dependency on player, currency, reward, quest, or match modules without updating the manifest and change spec first.
- Do not use negative or zero item quantities in grant flows.

## Required Tests

See `tests.required` in `module.yaml`.

For the first runtime slice, tests should cover:

- `GrantItem` accepts a valid grant and records the item.
- `GrantItem` rejects invalid quantity.
- `GrantItem` rejects grants that exceed capacity.
- `GrantItem` emits `ItemGranted` exactly once for a successful grant.
- `GetInventory` returns inventory state without mutating it.
- Permission failures use `INVENTORY_PERMISSION_DENIED`.
- Architecture checks still pass.

Run `node tools/vibit check runtime` after changing inventory runtime behavior. When Go is available, also run `cd runtime && go test ./...`.

PostgreSQL is the first authoritative durable store when inventory persistence begins. S3-compatible object storage is not required for the first inventory slice unless a future contract introduces large object artifacts that inventory owns.
