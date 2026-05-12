# Plan

## Files To Create

- `modules/inventory/generated/contracts/GetInventory.generated.ts`
- `modules/inventory/queries/GetInventory.ts`
- `modules/inventory/tests/GetInventory.test.ts`
- `conversations/2026-05-12-implement-get-inventory-query.md`

## Files To Edit

- `tools/vibit`
- `contracts/inventory/queries/GetInventory.yaml`
- `modules/inventory/module.yaml`
- `modules/inventory/AGENTS.md`
- `modules/inventory/AGENTS.zh-CN.md`
- `modules/inventory/policies/inventoryPermissionPolicy.ts`
- `changes/2026-05-12-implement-get-inventory-query/*`

## Generated Artifacts

- `modules/inventory/generated/contracts/GetInventory.generated.ts`

Generate this file with:

```bash
node tools/vibit generate contract --module inventory --type query --id GetInventory
```

Do not hand-edit generated output.

## Handwritten Logic

- Extend contract generation to support query contracts.
- Implement `GetInventory` query handler against the generated contract shape.
- Reuse the inventory repository interface.
- Extend the permission policy for `inventory_read`.
- Add focused runtime tests using Node.js built-in test runner.

## Tests

- Direct Node runtime query test.
- `check runtime` text mode.
- `check runtime --json`.
- `check generated`.
- `check contracts --json`.
- `check all --json`.
- Secret scan.
- Diff check.

## Verification Commands

- `node tools/vibit generate contract --module inventory --type query --id GetInventory`
- `node tools/vibit inspect contract --module inventory --type query --id GetInventory`
- `node tools/vibit check generated`
- `node --experimental-strip-types --test modules/inventory/tests/GetInventory.test.ts`
- `node tools/vibit check runtime`
- `node tools/vibit check runtime --json`
- `node tools/vibit check contracts --json`
- `node tools/vibit check all --json`
- `rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" --hidden --glob '!.git/**' --glob '!.vibit.local.env' .`
- `git diff --check`

## Rollback Or Migration Notes

Remove the handwritten query handler, remove the generated `GetInventory` shape by regenerating module declarations accordingly, and restore `GetInventory` `runtime_status` if this read-side proof slice is replaced.
