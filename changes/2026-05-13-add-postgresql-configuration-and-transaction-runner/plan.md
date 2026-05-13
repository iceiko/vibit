# Plan

1. Add PostgreSQL config parsing under `runtime/internal/platform/persistence/postgres/`.
2. Add a pgx transaction runner under `runtime/internal/platform/persistence/postgres/`.
3. Add tests using fake transaction handles rather than a live database.
4. Update architecture manifests and persistence guidance.
5. Mark `W-0017` complete and move `W-0018` to `next_ready`.
6. Run Go tests, vet, runtime checks, work checks, change checks, all checks, and whitespace checks.

## Files To Create

- `runtime/internal/platform/persistence/postgres/config.go`
- `runtime/internal/platform/persistence/postgres/config_test.go`
- `runtime/internal/platform/persistence/postgres/runner.go`
- `runtime/internal/platform/persistence/postgres/runner_test.go`

## Files To Edit

- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/work-items.yaml`
- `docs/postgresql-persistence-boundary.md`
- `docs/postgresql-persistence-boundary.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- This change spec.

## Rollback Notes

Remove the new PostgreSQL config and runner files, restore manifest/documentation status to `adapter_created_not_wired`, and move `W-0017` back to `next_ready`.
