# Impact Analysis

## Affected Modules

- `player`
- `runtime`

## Module Ownership Impact

The player module keeps ownership of the storage-neutral repository interface at `runtime/internal/modules/player/repository.go`.

The future PostgreSQL implementation is explicitly owned by the PostgreSQL platform package:

- `runtime/internal/platform/persistence/postgres/player_account_repository.go`
- `runtime/internal/platform/persistence/postgres/player_account_repository_test.go`

No ownership moves to transport, protocol, application dispatch, authentication, runtime sessions, inventory, S3, or MinIO.

## Public Contract Impact

No public command, query, event, error, permission, Protobuf envelope, or WebSocket contract changes.

The existing semantic contracts remain the source for `CreatePlayerAccount`, `GetPlayerAccount`, and `PlayerAccountCreated`. This change only defines the adapter boundary for a future implementation of `player.Repository`.

## Data And Migration Impact

No new migration source and no data shape change.

The existing player account migration `runtime/migrations/postgres/000002_create_player_account_state.sql` remains unchanged.

## Test Impact

No new Go adapter tests were added because adapter implementation is intentionally deferred.

The next adapter implementation must add focused fake-executor tests for SQL shape, transaction-control absence, mutation normalization, row mapping, missing rows, duplicate/constraint error mapping, and default no-live-PostgreSQL behavior.

## Documentation Impact

Updated the PostgreSQL persistence boundary standard and paired Simplified Chinese translation. Updated player/runtime manifests and guides to record that the adapter boundary is defined while adapter implementation and runtime handlers remain deferred.

## Compatibility Risks

Low. This change adds no runtime behavior.

The main risk reduced by this change is future drift: an agent should no longer invent a player account adapter in a handler, protocol package, authentication package, or hidden transaction path.
