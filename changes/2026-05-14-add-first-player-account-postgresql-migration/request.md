# Request

## Original Request

按照你的专业建议进行推进。

## Clarified Requirement

Advance `W-0051` by adding the first SQL-first PostgreSQL migration source for the ratified player account lifecycle schema.

## User-Visible Outcome

Maintainers and agents can inspect the first player account persistence schema source under `runtime/migrations/postgres/`.

The migration creates:

- `player_accounts`: one lifecycle row per stable `player_id`.
- `player_account_events`: durable lifecycle event records, including support for the first required event type `PlayerAccountCreated`.

## Non-Goals

- Do not add authentication providers.
- Do not add token behavior, token storage, refresh tokens, or signing metadata.
- Do not add credential or password storage.
- Do not add runtime session persistence.
- Do not add WebSocket connection state or handshake authentication.
- Do not add Protobuf envelope changes.
- Do not add runtime player account handlers or WebSocket routes.
- Do not add player account repository interfaces or PostgreSQL adapters in this step.
- Do not copy Nakama or Pitaya public API shapes.

## Unknowns

- Live PostgreSQL apply/rollback for this migration is not required by this work item and depends on an explicit disposable DSN.
- Account lifecycle update and deletion event types remain future semantic-contract work.
- Display-name uniqueness remains deferred.

## Acceptance Criteria

- [x] Add deterministic SQL-first goose migration source `runtime/migrations/postgres/000002_create_player_account_state.sql`.
- [x] Include `-- +goose Up`, `-- Module: player`, and `-- +goose Down`.
- [x] Create only `player_accounts` and `player_account_events`.
- [x] Include explicit constraints and required indexes from `ADR-0022`.
- [x] Update module/runtime manifests and guidance to record that migration source exists while runtime implementation remains deferred.
- [x] Verify migration and runtime boundary checks without adding repository adapters, handlers, authentication, tokens, credentials, sessions, Protobuf envelope changes, or WebSocket handshake changes.
