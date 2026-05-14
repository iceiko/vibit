# Plan

## Files To Create

- `changes/2026-05-14-close-durable-inventory-and-plan-player-session-boundary/request.md`
- `changes/2026-05-14-close-durable-inventory-and-plan-player-session-boundary/spec.yaml`
- `changes/2026-05-14-close-durable-inventory-and-plan-player-session-boundary/impact.md`
- `changes/2026-05-14-close-durable-inventory-and-plan-player-session-boundary/plan.md`
- `changes/2026-05-14-close-durable-inventory-and-plan-player-session-boundary/checklist.md`
- `changes/2026-05-14-close-durable-inventory-and-plan-player-session-boundary/verification.md`

## Files To Edit

- `.arch/work-items.yaml`
- `.arch/reference.yaml`
- `.arch/runtime.yaml`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`

## Generated Artifacts

- None yet.

## Handwritten Logic

No runtime logic changes.

## Tests

No new tests. This change is verified with repository checks.

## Verification Commands

- `node tools/vibit inspect work --json`
- `node tools/vibit check work --json`
- `node tools/vibit check architecture --json`
- `node tools/vibit check postgres-env --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check change close-durable-inventory-and-plan-player-session-boundary --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

Rollback by restoring the previous work queue and manifest statuses. No data migration is involved.
