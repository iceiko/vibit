# Plan

## Files To Create

- `changes/2026-05-14-add-contract-source-field-summaries/`

## Files To Edit

- `tools/vibit`
- `.arch/work-items.yaml`

## Generated Artifacts

- None.

## Handwritten Logic

Add small YAML path helpers for extracting required field lists and property names from existing contract source manifests.

## Tests

Run focused CLI inspections for commands, events, and a module-filtered query.

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit inspect contracts --type command --json`
- `node tools/vibit inspect contracts --type event --json`
- `node tools/vibit inspect contracts --module inventory --type query --json`
- `node tools/vibit check change add-contract-source-field-summaries --json`
- `node tools/vibit check work --json`

## Rollback Or Migration Notes

Remove the additive `fields` inspection output and helper functions. No migration is involved.
