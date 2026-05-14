# Plan

## Files To Create

- `decisions/ADR-0022-player-account-postgresql-schema-boundary.md`
- `changes/2026-05-14-define-player-account-postgresql-persistence-schema-boundary/`

## Files To Edit

- `docs/postgresql-persistence-boundary.md`
- `docs/postgresql-persistence-boundary.zh-CN.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/work-items.yaml`
- `modules/player/module.yaml`
- `modules/player/AGENTS.md`
- `modules/player/AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Generated Artifacts

None.

## Handwritten Logic

No runtime Go logic is added.

The CLI check logic is updated so player account migration files are blocked until schema ratification exists and constrained after it exists.

## Verification Commands

- `node tools/vibit check work --json`
- `node tools/vibit inspect next --json`
- `node tools/vibit check change define-player-account-postgresql-persistence-schema-boundary --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check contracts --json`
- `node tools/vibit check migrations --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

Rollback by restoring player account schema status to deferred and removing `ADR-0022`. No database migration exists yet.
