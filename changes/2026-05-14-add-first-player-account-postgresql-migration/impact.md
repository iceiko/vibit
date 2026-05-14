# Impact Analysis

## Affected Modules

- `player`: gains the first durable PostgreSQL schema source for player-owned account lifecycle state.
- `runtime`: gains the second SQL migration source under the accepted migration root.

## Module Ownership Impact

The player module owns the new tables:

- `player_accounts`
- `player_account_events`

The schema stores player account lifecycle state only. It does not own authentication credentials, token storage, runtime sessions, WebSocket connection state, inventory state, or permission grants.

## Public Contract Impact

No public command, query, event, error, permission, Protobuf schema, or WebSocket protocol changes.

The migration stores fields already ratified by the player account semantic and wire contracts, plus persistence metadata required by `ADR-0022`.

## Data And Migration Impact

Adds the second PostgreSQL migration source file:

- `runtime/migrations/postgres/000002_create_player_account_state.sql`

The `Up` migration creates player account lifecycle and lifecycle event tables. The `Down` migration drops indexes and tables in reverse dependency order.

## Test Impact

No player account repository adapter exists yet, so there are no repository integration tests for player account behavior in this step.

Repository source checks should validate migration naming, goose markers, table references, forbidden persistence concerns, and manifest status.

## Documentation Impact

Update runtime and player manifests/guidance to record that the first player account migration source exists while repository interfaces, PostgreSQL adapters, runtime handlers, authentication, token behavior, credential storage, session persistence, Protobuf envelope changes, and WebSocket handshake changes remain deferred.

## Compatibility Risks

Low. This is the first player account migration and no player account runtime persistence path consumes it yet.

The main future risk is editing the migration after it is treated as applied in a shared environment. Future data changes should add new migrations instead.
