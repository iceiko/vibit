# Plan

## Files To Create

- `changes/2026-05-14-close-authentication-postgresql-adapter-implementation-milestone/`

## Files To Edit

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`

## Generated Artifacts

- None.

## Handwritten Logic

- None.

## Tests

- No runtime tests are added because this is a workflow and standards closeout only.

## Verification Commands

- `node tools/vibit inspect next --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check module authentication --json`
- `node tools/vibit check change close-authentication-postgresql-adapter-implementation-milestone --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

This change has no database migration or runtime rollback path. Reverting the closeout would restore `M-015` as active and `W-0085` as next ready.
