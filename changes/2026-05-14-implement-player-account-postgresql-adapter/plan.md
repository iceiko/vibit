# Plan

## Files To Create

- `runtime/internal/platform/persistence/postgres/player_account_repository.go`
- `runtime/internal/platform/persistence/postgres/player_account_repository_test.go`
- `changes/2026-05-14-implement-player-account-postgresql-adapter/`

## Files To Edit

- `runtime/internal/platform/persistence/postgres/runner.go`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `.arch/work-items.yaml`
- `modules/player/module.yaml`
- `modules/player/AGENTS.md`
- `modules/player/AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `docs/postgresql-persistence-boundary.md`
- `docs/postgresql-persistence-boundary.zh-CN.md`
- `tools/vibit`

## Generated Artifacts

- None.

## Handwritten Logic

Implement a PostgreSQL adapter that:

- implements `player.Repository`
- normalizes create mutations and lookup IDs through player module helpers
- inserts `player_accounts`
- inserts `player_account_events` for `PlayerAccountCreated`
- reads account lifecycle rows from `player_accounts`
- maps pgx no-row, duplicate, and constraint failures to stable adapter sentinel errors
- normalizes timestamp output to UTC
- keeps transaction ownership outside the repository

## Tests

Add fake-executor unit tests under the PostgreSQL platform package.

No live PostgreSQL test is required by default.

## Verification Commands

- `cd runtime && go test ./internal/platform/persistence/postgres`
- `cd runtime && go test ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check module player --json`
- `node tools/vibit check contracts --json`
- `node tools/vibit check migrations --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change implement-player-account-postgresql-adapter --json`
- `node tools/vibit check all --json`
- `node tools/vibit inspect next --json`
- `git diff --check`

## Rollback Or Migration Notes

No database migration rollback is needed because no migration source changes.

If the adapter behavior is wrong, fix the adapter or its tests. Do not change `player.Repository` or the ratified migration schema without a new change spec and maintainer confirmation.
