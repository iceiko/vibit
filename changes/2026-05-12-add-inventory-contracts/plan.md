# Plan

## Files To Create

- `decisions/ADR-0007-yaml-contract-source-format.md`
- `.arch/contracts.yaml`
- `contracts/inventory/commands/GrantItem.yaml`
- `contracts/inventory/queries/GetInventory.yaml`
- `contracts/inventory/events/ItemGranted.yaml`
- `contracts/inventory/errors/inventory_errors.yaml`
- `contracts/inventory/permissions/inventory_permissions.yaml`
- `conversations/2026-05-12-inventory-contracts.md`
- `changes/2026-05-12-add-inventory-contracts/request.md`
- `changes/2026-05-12-add-inventory-contracts/spec.yaml`
- `changes/2026-05-12-add-inventory-contracts/impact.md`
- `changes/2026-05-12-add-inventory-contracts/plan.md`
- `changes/2026-05-12-add-inventory-contracts/checklist.md`
- `changes/2026-05-12-add-inventory-contracts/verification.md`

## Files To Edit

- `.arch/README.md`
- `.arch/README.zh-CN.md`
- `.arch/conventions.yaml`
- `modules/inventory/module.yaml`
- `modules/inventory/AGENTS.md`
- `modules/inventory/AGENTS.zh-CN.md`

## Generated Artifacts

None.

## Handwritten Logic

None.

## Tests To Add Or Update

None.

Runtime and generator tests are deferred until implementation exists.

## Verification Commands

- `node tools/vibit inspect module inventory`
- `node tools/vibit check module inventory`
- `node tools/vibit check memory`
- `node tools/vibit check all --json`
- Secret scan for GitHub token forms.
- `git diff --check`

## Rollback Notes

If the contract source format changes, supersede ADR-0007 and update `.arch/contracts.yaml` before generating runtime code.
