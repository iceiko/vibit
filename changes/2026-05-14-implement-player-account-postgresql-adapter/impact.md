# Impact Analysis

## Affected Modules

- `player`
- `runtime`

## Module Ownership Impact

The player module keeps ownership of the storage-neutral `player.Repository` interface and account lifecycle structs.

The PostgreSQL platform package owns the new adapter implementation, SQL shape, pgx error mapping, and fake-executor tests. No domain package imports pgx.

## Public Contract Impact

No public commands, queries, events, errors, permissions, Protobuf messages, or WebSocket envelope behavior changed.

The adapter maps PostgreSQL missing-row, duplicate, and constraint failures into stable adapter sentinel errors. Runtime client-facing error mapping remains deferred until runtime player handlers are ratified.

## Data And Migration Impact

No migration source changed. The adapter uses the already-ratified `player_accounts` and `player_account_events` schema from `runtime/migrations/postgres/000002_create_player_account_state.sql`.

## Test Impact

Added focused fake-executor tests for:

- account insert SQL
- `PlayerAccountCreated` event insert SQL
- account lookup SQL
- no transaction-control SQL
- mutation normalization and UTC timestamps
- nullable lifecycle timestamp row mapping
- missing-row behavior
- duplicate and constraint error mapping
- no live PostgreSQL dependency by default
- `UnitOfWork.NewPlayerAccountRepository`

## Documentation Impact

Updated English and Simplified Chinese persistence and agent guidance to say the adapter is implemented while handlers, routes, authentication, credentials, tokens, external identities, and sessions remain deferred.

## Compatibility Risks

Runtime client behavior is unchanged because no handler or route consumes the adapter yet.

The main implementation risk is future agents mistaking adapter availability for runtime product availability. The manifests and `tools/vibit check runtime` now keep `player_account_persistence_implemented: false`, `runtime_player_handlers_added: false`, and `websocket_routes_added: false` while allowing `postgres_adapter_added: true`.
