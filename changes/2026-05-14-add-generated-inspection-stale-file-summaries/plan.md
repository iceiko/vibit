# Plan

## Files To Create

- `changes/2026-05-14-add-generated-inspection-stale-file-summaries/`

## Files To Edit

- `tools/vibit`
- `.arch/work-items.yaml`

## Generated Artifacts

- None.

## Handwritten Logic

Reuse generated contract shape expected output paths to report stale generated files in inspection output.

## Tests

Run unfiltered and filtered generated inspections plus generated output checks.

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit inspect generated --json`
- `node tools/vibit inspect generated --module inventory --json`
- `node tools/vibit inspect generated --type command --json`
- `node tools/vibit check generated --json`
- `node tools/vibit check change add-generated-inspection-stale-file-summaries --json`
- `node tools/vibit check work --json`

## Rollback Or Migration Notes

Remove the additive stale summary fields. No migration is involved.
