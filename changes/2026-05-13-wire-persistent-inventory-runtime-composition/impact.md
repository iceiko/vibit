# Impact

## Architecture Impact

This change introduces the first persistent runtime composition path.

The application layer receives a narrow inventory repository provider abstraction so command handlers can resolve a transaction-bound repository from the current unit of work. This keeps the command orchestration point in `runtime/internal/app` while keeping PostgreSQL driver details inside `runtime/internal/platform/persistence/postgres`.

Process startup gains an explicit store selector:

```text
VIBIT_RUNTIME_STORE=memory
VIBIT_RUNTIME_STORE=postgres
```

The default remains `memory`.

## Module Impact

The inventory module's command and query contracts do not change.

The durable inventory implementation remains owned by the inventory module and the PostgreSQL platform adapter:

- Domain behavior remains in `runtime/internal/modules/inventory`.
- Persistent adapter behavior remains in `runtime/internal/platform/persistence/postgres`.
- Composition behavior is owned by `runtime/internal/app/bootstrap` and `runtime/cmd/vibit-server`.

## Data Impact

No new migration is required.

The existing migration source `runtime/migrations/postgres/000001_create_inventory_state.sql` remains the required schema for live PostgreSQL use. This change does not apply migrations automatically.

## Compatibility Impact

The WebSocket Protobuf protocol remains unchanged.

Default local runtime behavior remains compatible with the previous in-memory startup path.

## Agent Impact

Future agents get a clearer persistent runtime path:

- Start with in-memory composition for local behavior and most tests.
- Use explicit PostgreSQL composition when persistence is the subject of the change.
- Do not place business logic in process startup.
- Do not import pgx outside the PostgreSQL platform owner package.
