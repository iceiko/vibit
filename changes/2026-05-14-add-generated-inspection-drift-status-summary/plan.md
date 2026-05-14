# Plan

## Files To Create

- `changes/2026-05-14-add-generated-inspection-drift-status-summary/`

## Files To Edit

- `tools/vibit`
- `.arch/work-items.yaml`

## Generated Artifacts

- None.

## Handwritten Logic

Derive a compact generated contract shape status from missing and stale generated shape counts.

## Tests

Run unfiltered and filtered generated inspections plus generated output checks.

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit inspect generated --json`
- `node tools/vibit inspect generated --module inventory --json`
- `node tools/vibit inspect generated --type command --json`
- `node tools/vibit check generated --json`
- `node tools/vibit check change add-generated-inspection-drift-status-summary --json`
- `node tools/vibit check work --json`

## Rollback Or Migration Notes

Remove the additive status fields. No migration is involved.
