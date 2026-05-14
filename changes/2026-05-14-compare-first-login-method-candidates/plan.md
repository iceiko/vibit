# Plan

## Files To Create

- `docs/first-login-method-candidates.md`
- `docs/first-login-method-candidates.zh-CN.md`
- `conversations/2026-05-14-first-login-method-candidate-comparison.md`

## Files To Edit

- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`

## Generated Artifacts

None.

## Handwritten Logic

None.

## Tests To Add Or Update

None. This is a comparison and planning change only.

## Verification Commands

- `node tools/vibit inspect next --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check architecture --json`
- `node tools/vibit check agent-tooling --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check change compare-first-login-method-candidates --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Notes

Revert the comparison document, conversation log, manifest updates, and work queue transition. No runtime or data rollback is needed.
