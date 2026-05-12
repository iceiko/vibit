# Plan

## Files To Create

- `modules/inventory/generated/contracts/GrantItem.generated.ts`
- `conversations/2026-05-12-generate-inventory-contract-shape.md`

## Files To Edit

- `tools/vibit`
- `rules/check-rules.json`
- `modules/inventory/module.yaml`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `changes/2026-05-12-generate-inventory-contract-shape/*`

## Generated Artifacts

- `modules/inventory/generated/contracts/GrantItem.generated.ts`

## Handwritten Logic

- Add a narrow `generate contract` CLI command.
- Generate only the first command contract shape from existing contract source fields.
- Add a `check generated` command that verifies declared generated files and trace markers.
- Include `check generated` in `check all`.

## Tests

- Inspect the `GrantItem` contract before generation.
- Generate the contract file.
- Check generated files in text and JSON modes.
- Run contract consistency checks.
- Run full repository checks.
- Confirm generated file trace markers.
- Run secret scan.

## Verification Commands

- `node tools/vibit inspect contract --module inventory --type command --id GrantItem`
- `node tools/vibit generate contract --module inventory --type command --id GrantItem`
- `node tools/vibit check generated`
- `node tools/vibit check generated --json`
- `node tools/vibit check contracts --json`
- `node tools/vibit check all --json`
- `rg -n "@generated|Source: contracts/inventory/commands/GrantItem.yaml|Generator: tools/vibit generate contract" modules/inventory/generated/contracts/GrantItem.generated.ts`
- `rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" --hidden --glob '!.git/**' --glob '!.vibit.local.env' .`
- `git diff --check`

## Rollback Or Migration Notes

Remove the generated file, remove its module manifest declaration, and remove the generator/check commands if the first generated shape direction is replaced before runtime implementation begins.
