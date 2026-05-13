# Plan

## Files To Create

- `changes/2026-05-13-add-inventory-repository-mutation-lock-boundary/*`

## Files To Edit

- `runtime/internal/modules/inventory/inventory.go`
- `runtime/internal/modules/inventory/memory.go`
- `runtime/internal/modules/inventory/inventory_test.go`
- `docs/postgresql-persistence-boundary.md`
- `docs/postgresql-persistence-boundary.zh-CN.md`
- `modules/inventory/module.yaml`
- `modules/inventory/AGENTS.md`
- `modules/inventory/AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `.arch/runtime.yaml`
- `.arch/work-items.yaml`

## Generated Artifacts

- None.

## Handwritten Logic

- Add inventory-owned interfaces for acquiring a mutation lock and using a locked mutation repository view.
- Update `GrantItem` to acquire the lock before reading current inventory and applying the grant.
- Update the in-memory repository to provide a real mutex-backed locked view.

## Tests

- Add focused inventory tests for lock ordering and lock avoidance on pre-mutation failures.
- Run the full Go runtime test suite to catch request-loop and bootstrap regressions.

## Verification Commands

```bash
cd runtime && go test ./...
cd runtime && go vet ./...
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change add-inventory-repository-mutation-lock-boundary --json
node tools/vibit check all --json
git diff --check
```

## Rollback Or Migration Notes

Rollback is a normal code and documentation revert. No persisted data or migrations are involved.
