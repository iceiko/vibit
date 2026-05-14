# Plan

## Files To Create

- `docs/login-method-token-format-ratification.md`
- `docs/login-method-token-format-ratification.zh-CN.md`
- `decisions/ADR-0024-login-method-token-format-ratification-boundary.md`
- `conversations/2026-05-14-login-method-token-format-ratification-standard.md`

## Files To Edit

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

None. This is a standards and manifest change only.

## Verification Commands

- `node tools/vibit inspect next --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check architecture --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check change define-login-method-token-format-ratification-standard --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Notes

Revert the standard, ADR, manifest references, and work queue changes. No runtime or data migration rollback is needed.
