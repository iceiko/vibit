# Plan

## Files To Create

- `runtime/internal/platform/persistence/postgres/inventory_repository.go`
- `runtime/internal/platform/persistence/postgres/inventory_repository_test.go`
- `changes/2026-05-13-add-postgresql-inventory-repository-adapter/*`

## Files To Edit

- `runtime/go.mod`
- `runtime/go.sum`
- `runtime/internal/modules/inventory/inventory.go`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `modules/inventory/module.yaml`
- `modules/inventory/AGENTS.md`
- `modules/inventory/AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `docs/postgresql-persistence-boundary.md`
- `docs/postgresql-persistence-boundary.zh-CN.md`
- `.arch/work-items.yaml`

## Generated Artifacts

- None.

## Handwritten Logic

Add the PostgreSQL inventory repository adapter inside the platform persistence boundary.

Extend the inventory mutation contract so durable adapters receive the emitted event metadata before they persist the item grant record.

## Tests

Add focused Go tests for the PostgreSQL adapter package without requiring a live PostgreSQL instance.

Live PostgreSQL integration tests remain deferred until a local disposable database standard exists.

## Verification Commands

```bash
cd runtime && go test ./...
cd runtime && go vet ./...
node tools/vibit check runtime --json
node tools/vibit check migrations --json
node tools/vibit check work --json
node tools/vibit check change add-postgresql-inventory-repository-adapter --json
node tools/vibit check all --json
git diff --check
```

## Rollback Or Migration Notes

No schema rollback is required because no migration is added.

If this adapter shape is replaced, remove the package under `runtime/internal/platform/persistence/postgres/` and restore the inventory mutation contract to the previous in-memory-only shape.
