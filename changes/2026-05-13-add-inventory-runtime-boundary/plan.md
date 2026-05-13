# Plan

## Files To Create

- `runtime/internal/modules/inventory/inventory.go`
- `runtime/internal/modules/inventory/inventory_test.go`
- `changes/2026-05-13-add-inventory-runtime-boundary/`

## Files To Edit

- `modules/inventory/module.yaml`
- `modules/inventory/AGENTS.md`
- `modules/inventory/AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`

## Generated Artifacts

None.

## Handwritten Logic

- Inventory request/response/event structs scoped to handwritten runtime behavior.
- Repository and policy interfaces owned by inventory.
- `GrantItemHandler`.
- `GetInventoryHandler`.
- `RegisterRoutes` for application dispatcher integration.
- In-memory test doubles inside test files only.

## Tests

- Unit tests for command, query, policies, repository mutation, and event emission.
- Dispatcher integration test for registered inventory command/query routes.

## Verification Commands

- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check module inventory --json`
- `node tools/vibit check change add-inventory-runtime-boundary --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan for GitHub tokens in tracked and unignored files

## Rollback Or Migration Notes

Rollback can remove the new inventory runtime files and revert the module/runtime guide updates. No public contract, wire schema, or database migration is involved.
