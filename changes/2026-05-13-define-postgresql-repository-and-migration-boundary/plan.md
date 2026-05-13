# Plan

## Files To Create

- `docs/postgresql-persistence-boundary.md`
- `docs/postgresql-persistence-boundary.zh-CN.md`
- `decisions/ADR-0020-postgresql-persistence-boundary.md`
- `changes/2026-05-13-define-postgresql-repository-and-migration-boundary/*`

## Files To Edit

- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/work-items.yaml`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/inventory/module.yaml`
- `modules/inventory/AGENTS.md`
- `modules/inventory/AGENTS.zh-CN.md`
- `docs/runtime-runbook.md`
- `docs/runtime-runbook.zh-CN.md`

## Generated Artifacts

- None.

## Handwritten Logic

- None. This change intentionally does not add runtime persistence implementation code.

## Tests

- No Go tests are required for this documentation and standards change.

## Verification Commands

```bash
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change define-postgresql-repository-and-migration-boundary --json
node tools/vibit check all --json
git diff --check
```

## Rollback Or Migration Notes

Rollback is a normal documentation and manifest revert. No persisted data is involved.
