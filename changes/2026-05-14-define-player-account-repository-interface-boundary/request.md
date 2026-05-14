# Request

## Original Request

继续推进

## Clarified Requirement

Advance `W-0052` by defining the module-owned player account repository interface boundary after the first player account PostgreSQL migration source exists.

## User-Visible Outcome

Maintainers and agents can inspect the player account repository interface at `runtime/internal/modules/player/repository.go`.

The boundary declares:

- `Repository.CreatePlayerAccount(ctx, CreatePlayerAccountMutation)`
- `Repository.GetPlayerAccount(ctx, player_id)`
- Storage-neutral account lifecycle structs and mutation metadata for future PostgreSQL adapters.

## Non-Goals

- Do not add a PostgreSQL player account repository adapter.
- Do not add runtime player account command or query handlers.
- Do not add WebSocket routes or handshake behavior.
- Do not add authentication providers.
- Do not add token behavior, token storage, refresh tokens, or signing metadata.
- Do not add credential or password storage.
- Do not add runtime session persistence.
- Do not change the Protobuf envelope.
- Do not change the ratified player account migration schema or table ownership.
- Do not copy Nakama or Pitaya public API shapes.

## Unknowns

- PostgreSQL adapter implementation remains a future work item.
- Runtime player account handler and protocol bridge sequencing remains future work after the adapter boundary exists.
- Account lifecycle update and deletion repository methods remain future contract work.

## Acceptance Criteria

- [x] Add storage-neutral player account repository interface source under `runtime/internal/modules/player/`.
- [x] Add focused tests for repository interface shape and mutation normalization.
- [x] Update runtime identity checks to allow only this repository boundary under the player runtime module.
- [x] Update player/runtime manifests and guides so future PostgreSQL adapter work has a declared interface to implement.
- [x] Verify repository boundary checks without adding PostgreSQL adapter implementation, runtime player routes, authentication, tokens, credentials, session persistence, Protobuf envelope changes, or WebSocket handshake changes.
