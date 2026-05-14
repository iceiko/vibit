# Plan

## Files To Create

- `changes/2026-05-14-add-stale-generated-contract-shape-detection/`

## Files To Edit

- `tools/vibit`
- `.arch/work-items.yaml`

## Generated Artifacts

- None.

## Handwritten Logic

Derive the expected generated contract shape path set from registered semantic contracts and fail any generated contract shape file outside that set.

## Tests

Run generated output checks and generated inspection.

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit check generated --json`
- `node tools/vibit inspect generated --json`
- `node tools/vibit check change add-stale-generated-contract-shape-detection --json`
- `node tools/vibit check work --json`

## Rollback Or Migration Notes

Remove the expected-path check. No data migration is involved.
