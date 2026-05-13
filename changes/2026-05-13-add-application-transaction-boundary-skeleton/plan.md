# Plan

## Files To Create

- `runtime/internal/platform/tx/uow.go`
- `runtime/internal/platform/tx/uow_test.go`
- `runtime/internal/app/transactional_dispatch.go`
- `runtime/internal/app/transactional_dispatch_test.go`
- `changes/2026-05-13-add-application-transaction-boundary-skeleton/*`

## Files To Edit

- `tools/vibit`
- `.arch/runtime.yaml`
- `docs/postgresql-persistence-boundary.md`
- `docs/postgresql-persistence-boundary.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `docs/runtime-protocol-adapter.md`
- `docs/runtime-protocol-adapter.zh-CN.md`
- `.arch/work-items.yaml`

## Generated Artifacts

- None.

## Handwritten Logic

- Add a small `tx.Runner` and `tx.UnitOfWork` abstraction.
- Add a no-op runner for non-persistent tests and bootstrap use.
- Add an application `TransactionalDispatcher` wrapper that opens a unit of work for command routes and bypasses it for query routes.
- Adjust runtime architecture checks to allow `runtime/internal/app/` to import `runtime/internal/platform/tx` only.

## Tests

- Add unit tests for the no-op transaction runner.
- Add application dispatch tests for command/query transaction behavior and failure propagation.

## Verification Commands

```bash
cd runtime && go test ./...
cd runtime && go vet ./...
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change add-application-transaction-boundary-skeleton --json
node tools/vibit check all --json
git diff --check
```

## Rollback Or Migration Notes

Rollback is a normal code and documentation revert. No persisted data or migrations are involved.
