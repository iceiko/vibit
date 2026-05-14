# Plan

## Files To Create

- `changes/2026-05-14-add-contract-inspection-filters/`

## Files To Edit

- `tools/vibit`
- `docs/agent-tooling.md`
- `docs/agent-tooling.zh-CN.md`
- `.arch/work-items.yaml`

## Generated Artifacts

- None.

## Handwritten Logic

Add compatible argument parsing and filtering for contract index inspection.

## Tests

Run focused CLI inspection checks and existing repository checks.

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit inspect contracts --type command --json`
- `node tools/vibit inspect contracts --status draft --json`
- `node tools/vibit inspect contracts --module inventory --type query --json`
- `node tools/vibit check agent-tooling --json`
- `node tools/vibit check change add-contract-inspection-filters --json`
- `node tools/vibit check work --json`

## Rollback Or Migration Notes

Remove the optional parser branches and documentation references. No data migration is involved.
