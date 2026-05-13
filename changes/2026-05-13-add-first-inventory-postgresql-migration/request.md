# Request

## Original Request

Continue ten work items unless maintainer confirmation is required.

## Clarified Requirement

Advance `W-0013` by adding the first SQL-first PostgreSQL migration for inventory durable state.

## User-Visible Outcome

Maintainers and agents can inspect the first durable inventory schema source under `runtime/migrations/postgres/`.

The migration creates:

- `inventory_accounts`: one aggregate row per `player_id`, used as the command-safe mutation lock target.
- `inventory_items`: current item quantities per `(player_id, item_id)`.
- `inventory_item_grants`: one durable grant record per emitted `ItemGranted` event.

## Non-Goals

- Do not add Go migrations.
- Do not add object storage or MinIO.
- Do not change durable inventory data ownership.
- Do not implement the PostgreSQL repository adapter yet.
- Do not add migration tooling invocation yet.
- Do not add outbox or event delivery storage yet.
- Do not add player-account or item-catalog foreign key ownership.

## Unknowns

- Migration apply/rollback against a disposable PostgreSQL database remains deferred until migration verification tooling exists.
- Concrete PostgreSQL repository SQL remains deferred to the adapter work item.

## Acceptance Criteria

- [x] Add a deterministic SQL-first goose migration under `runtime/migrations/postgres/`.
- [x] Include both `-- +goose Up` and `-- +goose Down` sections.
- [x] Create inventory account, item quantity, and item grant record tables.
- [x] Preserve inventory ownership and reference external `player_id` and `item_id` without claiming player or catalog ownership.
- [x] Update manifests and documentation to record the migration source.
