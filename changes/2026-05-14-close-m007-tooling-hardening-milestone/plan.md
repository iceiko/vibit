# Plan

## Files To Create

- `changes/2026-05-14-close-m007-tooling-hardening-milestone/`

## Files To Edit

- `.arch/work-items.yaml`

## Generated Artifacts

- None.

## Handwritten Logic

None. This is a milestone closeout and workflow-state change.

## Tests

No tests are added.

## Verification Commands

- `node tools/vibit inspect work --json`
- `node tools/vibit inspect next --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change close-m007-tooling-hardening-milestone --json`
- `node tools/vibit check agent-tooling --json`
- `node tools/vibit check generated --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

Rollback by restoring M-007 to active, W-0048 to next_ready, and removing the M-008 confirmation gate. No data migration is involved.
