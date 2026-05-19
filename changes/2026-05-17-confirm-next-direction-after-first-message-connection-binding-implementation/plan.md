# Plan

## Files To Create

- `changes/2026-05-17-confirm-next-direction-after-first-message-connection-binding-implementation/`
- `changes/2026-05-17-define-postgres-session-persistence-schema-gate/`

## Files To Edit

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`

## Generated Artifacts

- None.

## Handwritten Logic

- None. This is a direction-only change.

## Tests

- No Go tests required.
- Repository checks cover the work queue and manifests.

## Verification Commands

- `node tools/vibit check change confirm-next-direction-after-first-message-connection-binding-implementation --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check all --json`

## Rollback Or Migration Notes

No data migration is created. Reversal means reopening `M-057/W-0129` and selecting a different direction before the schema gate is implemented.
