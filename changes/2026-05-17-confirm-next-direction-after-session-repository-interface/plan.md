# Plan

## Files To Create

- `changes/2026-05-17-confirm-next-direction-after-session-repository-interface/*`

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

The direction confirmation is verified through repository checks after the gate slice is recorded.

## Verification Commands

- `node tools/vibit check change confirm-next-direction-after-session-repository-interface --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check all --json`

## Rollback Or Migration Notes

Reversal means reopening `M-065/W-0137` and removing `M-066/W-0138` before any later PostgreSQL session adapter or runtime validation work depends on the gate.
