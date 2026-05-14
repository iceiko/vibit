# Request

## Original Request

Continue the next work item in the queue.

## Clarified Requirement

Advance `W-0054 Implement player account PostgreSQL adapter` by implementing only the player account PostgreSQL repository adapter inside the PostgreSQL platform package, adding focused fake-executor tests, updating architecture metadata, and preserving all runtime handler, WebSocket, authentication, credential, token, external identity, and session deferrals.

## User-Visible Outcome

Agents can now see a bounded player account PostgreSQL adapter implementation at `runtime/internal/platform/persistence/postgres/player_account_repository.go`, focused tests at `runtime/internal/platform/persistence/postgres/player_account_repository_test.go`, and a PostgreSQL package helper `UnitOfWork.NewPlayerAccountRepository`.

The adapter remains internal persistence infrastructure. It is not exposed through runtime player handlers or WebSocket routes.

## Non-Goals

- Do not add runtime player account command or query handlers.
- Do not add WebSocket routes.
- Do not add authentication, tokens, credentials, external identity linking, or session persistence.
- Do not change the ratified player account migration schema.
- Do not change the module-owned `player.Repository` interface.
- Do not make live PostgreSQL verification mandatory for default checks.

## Unknowns

- The later runtime handler and route composition milestone is not selected by this change.
- Production authentication, token, credential, external identity, and session persistence designs remain deferred.

## Acceptance Criteria

- [x] The adapter implements `player.Repository` in the PostgreSQL platform package.
- [x] `CreatePlayerAccount` inserts `player_accounts` and `player_account_events` through a caller-supplied executor.
- [x] `GetPlayerAccount` reads current lifecycle state from `player_accounts`.
- [x] The adapter does not call `BEGIN`, `COMMIT`, or `ROLLBACK`.
- [x] Focused fake-executor tests cover SQL shape, normalization, UTC and nullable timestamp mapping, missing rows, duplicate conflicts, constraint errors, and no-live-PostgreSQL default behavior.
- [x] Architecture metadata and bilingual guidance record that the adapter is implemented while runtime handlers and authentication remain deferred.
