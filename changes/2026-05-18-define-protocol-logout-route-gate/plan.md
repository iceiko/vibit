# Plan

## Files To Create

- `docs/protocol-logout-route-gate.md`
- `docs/protocol-logout-route-gate.zh-CN.md`
- `decisions/ADR-0079-protocol-logout-route-gate.md`
- `conversations/2026-05-18-protocol-logout-route-gate.md`
- `changes/2026-05-18-define-protocol-logout-route-gate/*`

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
- `tools/vibit`
- `rules/check-rules.json`

## Generated Artifacts

None.

## Handwritten Logic

Only repository check logic in `tools/vibit`; no Go runtime behavior.

## Tests

No Go tests. Use repository checks and JavaScript syntax check.

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit check schemas --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change define-protocol-logout-route-gate --json`
- `node tools/vibit check all --json`
- `git diff --check`
- `node tools/vibit inspect next --json`

## Rollback Or Migration Notes

This is a gate-only documentation and check-rule change. Reversal would remove ADR-0079, the protocol logout route gate standard, and associated manifest/check markers.
