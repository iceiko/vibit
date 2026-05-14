# Plan

## Files To Create

- `docs/session-persistence-websocket-handshake-decision-gates.md`
- `docs/session-persistence-websocket-handshake-decision-gates.zh-CN.md`
- `changes/2026-05-14-define-session-persistence-websocket-handshake-decision-gates/request.md`
- `changes/2026-05-14-define-session-persistence-websocket-handshake-decision-gates/spec.yaml`
- `changes/2026-05-14-define-session-persistence-websocket-handshake-decision-gates/impact.md`
- `changes/2026-05-14-define-session-persistence-websocket-handshake-decision-gates/plan.md`
- `changes/2026-05-14-define-session-persistence-websocket-handshake-decision-gates/checklist.md`
- `changes/2026-05-14-define-session-persistence-websocket-handshake-decision-gates/verification.md`

## Files To Edit

- `docs/authentication-token-session-validation.md`
- `docs/authentication-token-session-validation.zh-CN.md`
- `.arch/conventions.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`

## Generated Artifacts

None.

## Handwritten Logic

None.

## Tests To Add Or Update

No Go tests are required because no Go runtime code changes.

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit check runtime --json`
- `node tools/vibit check architecture --json`
- `node tools/vibit check work --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check change define-session-persistence-websocket-handshake-decision-gates --json`
- `node tools/vibit check all --json`
- `node tools/vibit inspect next --json`
- `git diff --check`

## Rollback Notes

This is a documentation and manifest boundary update. Rollback would remove the new standard, change spec, manifest references, guide references, and work queue state update.
