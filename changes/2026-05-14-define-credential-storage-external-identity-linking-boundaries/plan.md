# Plan

## Files To Create

- `docs/credential-storage-external-identity-linking-boundaries.md`
- `docs/credential-storage-external-identity-linking-boundaries.zh-CN.md`
- `changes/2026-05-14-define-credential-storage-external-identity-linking-boundaries/request.md`
- `changes/2026-05-14-define-credential-storage-external-identity-linking-boundaries/spec.yaml`
- `changes/2026-05-14-define-credential-storage-external-identity-linking-boundaries/impact.md`
- `changes/2026-05-14-define-credential-storage-external-identity-linking-boundaries/plan.md`
- `changes/2026-05-14-define-credential-storage-external-identity-linking-boundaries/checklist.md`
- `changes/2026-05-14-define-credential-storage-external-identity-linking-boundaries/verification.md`

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
- `modules/player/module.yaml`

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
- `node tools/vibit check module player --json`
- `node tools/vibit check change define-credential-storage-external-identity-linking-boundaries --json`
- `node tools/vibit check all --json`
- `node tools/vibit inspect next --json`
- `git diff --check`

## Rollback Notes

This is a documentation and manifest boundary update. Rollback would remove the new standard, change spec, manifest references, guide references, module manifest markers, and work queue state update.
