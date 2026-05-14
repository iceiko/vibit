# Plan

## Files To Create

- `runtime/internal/platform/persistence/postgres/authentication_repository.go`
- `runtime/internal/platform/persistence/postgres/authentication_repository_test.go`
- `changes/2026-05-14-implement-authentication-postgresql-adapter/`

## Files To Edit

- `runtime/internal/platform/persistence/postgres/runner.go`
- `runtime/internal/platform/persistence/postgres/inventory_repository_test.go`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `docs/postgresql-persistence-boundary.md`
- `docs/postgresql-persistence-boundary.zh-CN.md`
- `docs/selected-login-token-boundary-checks.md`
- `docs/selected-login-token-boundary-checks.zh-CN.md`
- `tools/vibit`

## Generated Artifacts

- None.

## Handwritten Logic

Implement a PostgreSQL adapter that:

- implements `authentication.Repository`
- normalizes mutations and lookup digests through authentication module helpers
- inserts and reads `authentication_device_credentials`
- inserts, reads, revokes, and lists cleanup-eligible `authentication_access_tokens`
- maps pgx no-row, duplicate, foreign-key, and check-constraint failures to stable adapter sentinel errors
- normalizes timestamp output to UTC
- keeps transaction ownership outside the repository

## Tests

Add fake-executor unit tests under the PostgreSQL platform package.

No live PostgreSQL test is required by default.

## Verification Commands

- `cd runtime && go test ./internal/platform/persistence/postgres`
- `cd runtime && go test ./...`
- `node -c tools/vibit`
- `node tools/vibit check runtime --json`
- `node tools/vibit check module authentication --json`
- `node tools/vibit check contracts --json`
- `node tools/vibit check migrations --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change implement-authentication-postgresql-adapter --json`
- `node tools/vibit inspect next --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

No database migration rollback is needed because no migration source changes.

If the adapter behavior is wrong, fix the adapter or its tests. Do not change `authentication.Repository` or the ratified migration schemas without a new change spec and maintainer confirmation.
