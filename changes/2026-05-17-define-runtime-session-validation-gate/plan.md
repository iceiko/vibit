# Plan

## Files To Create

- `docs/runtime-session-validation-gate.md`
- `docs/runtime-session-validation-gate.zh-CN.md`
- `decisions/ADR-0065-runtime-session-validation-gate.md`
- `conversations/2026-05-17-runtime-session-validation-gate.md`
- `changes/2026-05-17-define-runtime-session-validation-gate/`

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

No runtime logic.

## Tests

No Go tests are required for this gate-only change.

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit check change define-runtime-session-validation-gate --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check all --json`

## Rollback Or Migration Notes

No migration rollback is needed. Reversal means replacing the gate with a new ADR and manifest updates.
