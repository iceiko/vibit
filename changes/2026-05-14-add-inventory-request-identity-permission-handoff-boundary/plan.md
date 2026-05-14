# Plan

## Files To Create

- None beyond this change spec.

## Files To Edit

- `runtime/internal/modules/inventory/inventory.go`
- `runtime/internal/modules/inventory/memory.go`
- `runtime/internal/modules/inventory/inventory_test.go`
- `runtime/internal/app/bootstrap/inventory_test.go`
- `docs/player-identity-session-boundary.md`
- `docs/player-identity-session-boundary.zh-CN.md`
- `modules/inventory/module.yaml`
- `modules/inventory/AGENTS.md`
- `modules/inventory/AGENTS.zh-CN.md`
- `.arch/runtime.yaml`
- `.arch/work-items.yaml`

## Generated Artifacts

- None.

## Handwritten Logic

Add a small permission context for inventory permission decisions and bridge route request identity into handler permission calls. Preserve static bootstrap policy as an explicit allow/deny fixture and provide an identity-aware guard policy that denies metadata-only identity for privileged grants.

## Tests

Add or update inventory and bootstrap tests for identity handoff and metadata-only denial/neutral behavior.

## Verification Commands

- `cd runtime && go test ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change add-inventory-request-identity-permission-handoff-boundary --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

Rollback is source-only. No persisted data or generated outputs are changed.
