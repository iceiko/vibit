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

Runtime schema files and generated files are not created yet. Before implementing the first runtime slice, declare schemas for `GrantItem`, `GetInventory`, `ItemGranted`, inventory errors, and inventory permissions.

## Forbidden Shortcuts

- Do not bypass boundaries declared in `module.yaml`.
- Do not directly modify data owned by another module.
- Do not add unregistered public commands, queries, events, or permissions.
- Do not put inventory business rules in HTTP or transport handlers.
- Do not hand-edit generated files. If generated output is wrong, change the source contract, template, or generator.
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

If runtime test infrastructure does not exist yet, record the tests as not available instead of removing them from the manifest.
