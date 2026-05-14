# Plan

1. Add `generate contract-shapes <module|all>`.
2. Render generated Go files from registered non-runtime contract manifests.
3. Include generated/source/generator/contract traces in each file.
4. Extend `check generated` to require and reproduce contract shapes.
5. Generate contract shapes for inventory and player.
6. Run generated-output checks and Go tests.

## Files To Edit

- `tools/vibit`
- `docs/generated-output.md`
- `docs/generated-output.zh-CN.md`
- `.arch/work-items.yaml`

## Files To Create

- `runtime/internal/generated/contracts/inventory/commands/GrantItem.go`
- `runtime/internal/generated/contracts/inventory/queries/GetInventory.go`
- `runtime/internal/generated/contracts/inventory/events/ItemGranted.go`
- `runtime/internal/generated/contracts/inventory/errors/inventory_errors.go`
- `runtime/internal/generated/contracts/inventory/permissions/inventory_permissions.go`
- `runtime/internal/generated/contracts/player/commands/CreatePlayerAccount.go`
- `runtime/internal/generated/contracts/player/queries/GetPlayerAccount.go`
- `runtime/internal/generated/contracts/player/events/PlayerAccountCreated.go`
- `runtime/internal/generated/contracts/player/errors/player_account_errors.go`
- `runtime/internal/generated/contracts/player/permissions/player_account_permissions.go`

## Verification Commands

- `node tools/vibit generate contract-shapes all`
- `node tools/vibit inspect generated --json`
- `node tools/vibit check generated --json`
- `cd runtime && go test ./...`
