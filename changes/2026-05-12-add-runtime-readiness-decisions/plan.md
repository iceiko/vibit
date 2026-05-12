# Plan

## Files To Create

- `.arch/runtime.yaml`
- `decisions/ADR-0003-first-reference-runtime-language.md`
- `decisions/ADR-0004-minimal-server-instance-model.md`
- `decisions/ADR-0005-contract-and-generation-boundary.md`
- `decisions/ADR-0006-first-runtime-proof-slice.md`
- `conversations/2026-05-12-runtime-readiness-decisions.md`
- `changes/2026-05-12-add-runtime-readiness-decisions/request.md`
- `changes/2026-05-12-add-runtime-readiness-decisions/spec.yaml`
- `changes/2026-05-12-add-runtime-readiness-decisions/impact.md`
- `changes/2026-05-12-add-runtime-readiness-decisions/plan.md`
- `changes/2026-05-12-add-runtime-readiness-decisions/checklist.md`
- `changes/2026-05-12-add-runtime-readiness-decisions/verification.md`

## Files To Edit

- `.arch/README.md`
- `.arch/README.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`

## Generated Artifacts

None.

## Handwritten Logic

None.

## Tests To Add Or Update

None.

This is a standards and architecture decision change.

## Verification Commands

- `node tools/vibit check memory`
- `node tools/vibit check all --json`
- Secret scan for obvious leaked tokens.
- `git diff --check`

## Rollback Notes

If the runtime direction changes, supersede the affected ADR and update `.arch/runtime.yaml`. Do not silently edit history after implementation has started.
