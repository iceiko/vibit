# Plan

## Files To Create

- `docs/postgres-session-persistence-schema-gate.md`
- `docs/postgres-session-persistence-schema-gate.zh-CN.md`
- `decisions/ADR-0059-postgres-session-persistence-schema-gate.md`
- `changes/2026-05-17-define-postgres-session-persistence-schema-gate/`
- `conversations/2026-05-17-postgres-session-persistence-schema-gate.md`

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
- `rules/check-rules.json`
- `tools/vibit`

## Generated Artifacts

- None.

## Handwritten Logic

- No Go runtime logic.
- Add Node repository check logic only.

## Tests

- Repository checks only.

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit check runtime --json`
- `node tools/vibit check module authentication --json`
- `node tools/vibit check work --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check change define-postgres-session-persistence-schema-gate --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

No migration source is created. Reversal means removing or superseding the gate before any future `runtime_sessions` migration source is added.
