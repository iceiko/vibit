# Plan

## Files To Create

- `runtime/internal/platform/migrations/postgres.go`
- `runtime/internal/platform/migrations/postgres_test.go`

## Files To Edit

- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/work-items.yaml`
- `docs/postgresql-persistence-boundary.md`
- `docs/postgresql-persistence-boundary.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `runtime/go.mod`
- This change spec.

## Generated Artifacts

- None.

## Handwritten Logic

Add migration status and apply helpers that validate options, configure goose for PostgreSQL, and operate on caller-supplied database handles and migration directories.

## Tests

Add unit tests that validate required options and verify that live database execution is explicitly skipped when `VIBIT_POSTGRES_MIGRATION_TEST_DSN` is not present.

## Verification Commands

- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `node tools/vibit check migrations --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change add-postgresql-migration-apply-and-status-tooling --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

Remove the new migration package files, restore manifest and guidance statuses to source-check-only, and move `W-0018` back to `next_ready`.
