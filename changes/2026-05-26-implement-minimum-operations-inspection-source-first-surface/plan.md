# Plan

## Files To Create

- `decisions/ADR-0153-minimum-operations-inspection-source-first-surface-implementation.md`
- `conversations/2026-05-26-minimum-operations-inspection-source-first-surface-implementation.md`
- `changes/2026-05-26-implement-minimum-operations-inspection-source-first-surface/`

## Files To Edit

- `tools/vibit`
- `rules/check-rules.json`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `modules/friends/module.yaml`
- `docs/runtime-runbook.md`
- `docs/runtime-runbook.zh-CN.md`
- `docs/alpha-acceptance-checklist.md`
- `docs/alpha-acceptance-checklist.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`

## Generated Artifacts

None.

## Handwritten Logic

Add `inspectOperations` support to `tools/vibit` and a repository check for the emitted JSON shape.

## TDD

1. Run `node tools/vibit inspect operations --json`.
2. Confirm it fails with `Unknown command`.
3. Implement the smallest command that emits the required source-first inspection record.
4. Add repository check coverage.

## Verification Commands

```bash
node -c tools/vibit
node tools/vibit inspect operations --json
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.minimum_operations_inspection_source_first_surface_implementation
node tools/vibit check change implement-minimum-operations-inspection-source-first-surface --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check schemas --json
node tools/vibit check memory --json
node tools/vibit check all --json
git diff --check
```

## Rollback Or Migration Notes

Rollback is source-only: remove the command, rule, ADR/change/memory artifacts, and restore W-0245 as next-ready. No runtime data migration is involved.
