# Impact Analysis

## Affected Modules

- `inventory`: extends the module-owned repository boundary with a command-safe mutation lock and routes `GrantItem` through it.
- `runtime`: updates in-memory bootstrap support and persistence boundary guidance.

## Module Ownership Impact

Inventory repository interfaces remain owned by the inventory module. PostgreSQL adapter code will still belong under `runtime/internal/platform/persistence/postgres/` when implemented.

This change does not move transaction ownership into repositories. It only defines the module-owned repository view that must be bound to an application-owned unit of work once `W-0012` exists.

## Public Contract Impact

No public command, query, event, error, permission, Protobuf schema, or WebSocket protocol changes.

The changed interface is internal Go runtime shape under `runtime/internal/modules/inventory/`.

## Data And Migration Impact

No migration files are added.

The future PostgreSQL adapter must implement this boundary with an explicit inventory account row lock for `player_id`, after the migration creates the required account table.

## Test Impact

Inventory runtime tests should verify that successful `GrantItem` calls lock the aggregate before reading inventory and applying the grant, and that validation or permission failures do not acquire a mutation lock.

Existing request-loop tests should continue to pass with the in-memory repository.

## Documentation Impact

Update the PostgreSQL persistence boundary standard, runtime agent guide, inventory agent guide, and inventory module manifest with the explicit locked repository view.

Update paired Simplified Chinese translations for materially changed public guidance.

## Compatibility Risks

Low. The change affects internal runtime interfaces before any PostgreSQL adapter exists.

The main implementation risk is accidentally making the lock optional for persistent command flows. The docs and tests should keep the persistent path mandatory.
