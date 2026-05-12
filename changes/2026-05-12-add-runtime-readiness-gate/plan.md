# Plan

## Files To Create

- `changes/2026-05-12-add-runtime-readiness-gate/`
- `conversations/2026-05-12-runtime-readiness-gate.md`

## Files To Edit

- `CONSTITUTION.md`
- `CONSTITUTION.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `decisions/ADR-0002-productive-bootstrapping.md`

## Generated Artifacts

- None.

## Handwritten Logic

- None. This is a governance clarification.

## Tests

- `node tools/vibit check memory`
- `node tools/vibit check all --json`
- Secret scan
- Diff whitespace check

## Verification Commands

- `node tools/vibit check memory`
- `node tools/vibit check all --json | node -e '<JSON.parse assertion>'`
- `rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" . .git/config`
- `git diff --check`

## Rollback Or Migration Notes

Revert the readiness clarification if it causes agents to delay implementation without improving architecture decisions.
