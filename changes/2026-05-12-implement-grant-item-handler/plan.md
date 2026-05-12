# Plan

## Files To Create

- `modules/inventory/commands/GrantItem.ts`
- `modules/inventory/repositories/InMemoryInventoryRepository.ts`
- `modules/inventory/policies/inventoryCapacityPolicy.ts`
- `modules/inventory/policies/inventoryPermissionPolicy.ts`
- `modules/inventory/tests/GrantItem.test.ts`
- `conversations/2026-05-12-implement-grant-item-handler.md`

## Files To Edit

- `tools/vibit`
- `rules/check-rules.json`
- `modules/inventory/module.yaml`
- `modules/inventory/AGENTS.md`
- `modules/inventory/AGENTS.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `changes/2026-05-12-implement-grant-item-handler/*`

## Generated Artifacts

- None. Do not hand-edit `modules/inventory/generated/contracts/GrantItem.generated.ts`.

## Handwritten Logic

- Implement `GrantItem` command handler against the generated contract shape.
- Implement an in-memory inventory repository.
- Implement inventory capacity and permission policies.
- Add focused runtime tests using Node.js built-in test runner.
- Add `check runtime` to execute runtime tests.
- Include `check runtime` in `check all`.

## Tests

- Direct Node runtime test.
- `check runtime` text mode.
- `check runtime --json`.
- `check generated`.
- `check all --json`.
- Secret scan.
- Diff check.

## Verification Commands

- `node tools/vibit inspect contract --module inventory --type command --id GrantItem`
- `node tools/vibit check generated`
- `node --experimental-strip-types --test modules/inventory/tests/GrantItem.test.ts`
- `node tools/vibit check runtime`
- `node tools/vibit check runtime --json`
- `node tools/vibit check all --json`
- `rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" --hidden --glob '!.git/**' --glob '!.vibit.local.env' .`
- `git diff --check`

## Rollback Or Migration Notes

Remove the handwritten inventory files, remove the runtime check, and restore `proof_slice.runtime_implementation` if this first executable path is replaced before more runtime code depends on it.
