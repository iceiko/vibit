# Plan

## Files To Create

- `changes/2026-05-17-confirm-next-direction-after-session-postgresql-adapter-implementation/`

## Files To Edit

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`

## Generated Artifacts

- None.

## Handwritten Logic

No runtime logic.

## Tests

No Go tests are required for this direction-selection step.

## Verification Commands

- `node tools/vibit check change confirm-next-direction-after-session-postgresql-adapter-implementation --json`
- `node tools/vibit check work --json`
- `node tools/vibit check all --json`

## Rollback Or Migration Notes

No data rollback is needed. Reversal means reopening the confirmation gate and selecting a different direction.
