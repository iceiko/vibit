# Plan

## Files To Create

- `runtime/migrations/postgres/000001_create_inventory_state.sql`
- `changes/2026-05-13-add-first-inventory-postgresql-migration/*`

## Files To Edit

- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `modules/inventory/module.yaml`
- `modules/inventory/AGENTS.md`
- `modules/inventory/AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `.arch/work-items.yaml`

## Generated Artifacts

- None.

## Handwritten Logic

- None. This change adds SQL migration source only.

## Tests

- No Go runtime behavior tests are required by the migration itself, but the existing Go suite should still pass.
- Migration apply/rollback verification remains deferred until migration tooling checks exist.

## Verification Commands

```bash
cd runtime && go test ./...
cd runtime && go vet ./...
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change add-first-inventory-postgresql-migration --json
node tools/vibit check all --json
git diff --check
```

## Rollback Or Migration Notes

Before the migration is applied in any shared environment, rollback is a normal source revert.

After shared application, do not edit this migration in place. Add a new migration instead.
