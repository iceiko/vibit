# Plan

## Files To Create

- `.arch/work-items.yaml`
- `docs/workflow.md`
- `docs/workflow.zh-CN.md`
- `changes/2026-05-13-add-work-item-continuation-system/`

## Files To Edit

- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `.arch/README.md`
- `.arch/README.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`
- `schema/inspect-output.schema.json`

## Generated Artifacts

- None.

## Handwritten Logic

- Add `check work` to verify `.arch/work-items.yaml` and workflow docs.
- Add `inspect work` to return current milestone, next-ready work items, active items, and blocked items.
- Add work check to `check all`.

## Tests

- Use CLI checks as the test surface for the workflow metadata.
- Keep implementation dependency-free.

## Verification Commands

- `node tools/vibit check work --json`
- `node tools/vibit inspect work`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan for GitHub tokens in tracked and unignored files

## Rollback Or Migration Notes

Rollback can remove the work-item manifest, workflow docs, and CLI support. No runtime or data migration is involved.
