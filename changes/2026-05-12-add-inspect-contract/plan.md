# Plan

## Files To Create

- `conversations/2026-05-12-inspect-contract.md`

## Files To Edit

- `tools/vibit`
- `schema/inspect-output.schema.json`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `changes/2026-05-12-add-inspect-contract/*`

## Generated Artifacts

- None.

## Handwritten Logic

- Parse `inspect contract --module <module> --type <type> --id <id>` arguments.
- Normalize singular and plural contract type names.
- Read `.arch/contracts.yaml` through the existing lightweight registry entry parser.
- Read the referenced contract source file when present.
- Read `modules/<module>/module.yaml` when present.
- Print a `contract_inspection` JSON object.
- Return a non-zero exit code with a clear message when the requested contract is not registered.

## Tests

- Existing command contract inspection JSON parse check.
- Existing error catalog inspection JSON parse check.
- Missing contract failure check.
- Help text check.
- Contract consistency check.
- Schema check.
- Full repository check.
- Secret scan.

## Verification Commands

- `node tools/vibit inspect contract --module inventory --type command --id GrantItem`
- `node tools/vibit inspect contract --module inventory --type command --id GrantItem | node -e '<JSON.parse assertion>'`
- `node tools/vibit inspect contract --module inventory --type error --id inventory_errors | node -e '<JSON.parse assertion>'`
- `node tools/vibit inspect contract --module inventory --type command --id MissingCommand`
- `node tools/vibit --help | rg "inspect contract"`
- `node tools/vibit check contracts --json`
- `node tools/vibit check schemas`
- `node tools/vibit check all --json`
- `rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" --hidden --glob '!.git/**' --glob '!.vibit.local.env' .`
- `git diff --check`

## Rollback Or Migration Notes

Remove the inspect command if contract intake is replaced by a broader structured contract index.
