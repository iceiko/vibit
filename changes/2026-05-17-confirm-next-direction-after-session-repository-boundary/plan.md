# Plan

## Files To Create

- `changes/2026-05-17-confirm-next-direction-after-session-repository-boundary/*`

## Files To Edit

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- AGENTS guides and module manifest references

## Generated Artifacts

- None.

## Handwritten Logic

No runtime logic is added by the direction confirmation itself.

## Tests

The direction confirmation is verified through repository checks after the implementation slice is recorded.

## Verification Commands

- `node tools/vibit check change confirm-next-direction-after-session-repository-boundary --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check all --json`

## Rollback Or Migration Notes

Reversal means reopening `M-063/W-0135` and removing `M-064/W-0136` before any later adapter or runtime validation work depends on the interface.
