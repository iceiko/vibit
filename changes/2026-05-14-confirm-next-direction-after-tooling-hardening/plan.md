# Plan

## Files To Create

- `changes/2026-05-14-confirm-next-direction-after-tooling-hardening/`
- `conversations/2026-05-14-player-account-postgresql-persistence-direction.md`

## Files To Edit

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`

## Generated Artifacts

None.

## Handwritten Logic

None. This is workflow and architecture state only.

## Verification Commands

- `node tools/vibit check work --json`
- `node tools/vibit inspect next --json`
- `node tools/vibit check change confirm-next-direction-after-tooling-hardening --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

Rollback by restoring `M-008` and `W-0049` to the blocked confirmation gate and removing `M-009`. No data migration is involved.
