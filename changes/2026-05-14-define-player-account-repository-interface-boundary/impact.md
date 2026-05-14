# Impact Analysis

## Affected Modules

- `player`
- `runtime`

## Module Ownership Impact

The player module now owns a storage-neutral runtime repository interface at `runtime/internal/modules/player/repository.go`.

Ownership does not move to PostgreSQL, transport, protocol, application dispatch, inventory, authentication, or runtime session packages.

## Public Contract Impact

No public command, query, event, error, permission, Protobuf envelope, or WebSocket contract changes.

The existing semantic contracts remain the source for `CreatePlayerAccount`, `GetPlayerAccount`, and `PlayerAccountCreated`. This change only creates the module-owned repository boundary that future adapters will implement.

## Data And Migration Impact

No new migration source and no data shape change.

The existing player account migration `runtime/migrations/postgres/000002_create_player_account_state.sql` remains unchanged.

## Test Impact

Added focused Go tests for:

- repository interface compile-time shape
- mutation required-field normalization
- UTC normalization of `occurred_at`
- initial account state guard
- account state validity

## Documentation Impact

Updated English and Simplified Chinese player/runtime guidance plus the PostgreSQL persistence boundary standard to record that the repository interface exists while adapter and handler implementation remains deferred.

## Compatibility Risks

Low. The new Go package is not wired into application dispatch, protocol bridges, transport routes, or persistence composition.
