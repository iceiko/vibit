# Impact Analysis

## Affected Modules

- `inventory`

## Module Ownership Impact

No module ownership moves.

This change implements the first handwritten inventory runtime logic inside the declared `runtime/internal/modules/inventory/` boundary. The inventory module still owns inventory records, inventory items, item grant behavior, inventory reads, capacity rules, permissions, and inventory invariants.

The module still does not own player accounts, item catalog entries, currency balances, rewards, quests, match sessions, transport, protocol framing, PostgreSQL drivers, or object-storage clients.

## Public Contract Impact

No public command, query, event, error, or permission contract changes are planned.

The implementation uses existing contracts:

- `GrantItem`
- `GetInventory`
- `ItemGranted`
- `INVALID_ITEM_QUANTITY`
- `INVENTORY_CAPACITY_EXCEEDED`
- `INVENTORY_PERMISSION_DENIED`
- `inventory_grant_item`
- `inventory_read`

## Runtime Impact

Adds pure Go domain runtime behavior under `runtime/internal/modules/inventory/`.

Adds vibit-owned interfaces for:

- Inventory repository behavior.
- Inventory permission policy.
- Inventory capacity policy.
- Event identity and time generation for deterministic tests.

Adds an inventory dispatcher registration function that registers `inventory.GrantItem` and `inventory.GetInventory` handlers with `runtime/internal/app.Dispatcher`.

This change intentionally keeps Protobuf payload mapping outside the inventory module. Protocol adapters or generated bridges should map wire payloads to inventory request structs in a later change.

## Data And Migration Impact

No database migration is added.

The repository interface is defined before a PostgreSQL adapter exists. A later persistence change can implement the interface under `runtime/internal/platform/persistence/postgres/` without changing domain behavior.

## Test Impact

Adds Go tests for inventory runtime behavior and application dispatcher integration.

Expected coverage:

- Valid item grant mutates repository state.
- Invalid quantity fails before mutation.
- Capacity rejection fails before mutation.
- Permission rejection fails before mutation.
- Successful grant emits exactly one `ItemGranted` event.
- `GetInventory` reads without mutating state.
- Dispatcher integration preserves app route metadata.

## Documentation Impact

Updates module and runtime agent guides to reflect that inventory runtime boundary work has started.

## Compatibility Risks

The main risk is hardening handwritten request/response structs before generated Go contracts exist.

The mitigation is to keep these structs inside the inventory module runtime boundary, name them after existing contracts, and avoid claiming that they are generated public contract output. A later generator can replace or adapt this shape through an explicit generated-contract change.
