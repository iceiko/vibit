# Plan

## Files To Create

- `decisions/ADR-0002-productive-bootstrapping.md`
- `changes/2026-05-12-add-bootstrapping-control/`
- `conversations/2026-05-12-bootstrapping-control.md`

## Files To Edit

- `CONSTITUTION.md`
- `CONSTITUTION.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`

## Generated Artifacts

- None.

## Handwritten Logic

- None. This is a governance and standards change only.

## Tests

- `node tools/vibit check memory`
- `node tools/vibit check all`
- `node tools/vibit check all --json`
- Secret scan
- Diff whitespace check

## Verification Commands

- `node tools/vibit check memory`
- `node tools/vibit check all`
- `node tools/vibit check all --json | node -e '<JSON.parse assertion>'`
- `rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" . .git/config`
- `git diff --check`

## Rollback Or Migration Notes

Remove or revise ADR-0002 and the paired constitution/AGENTS sections if runtime work shows that the rule blocks necessary tooling rather than preventing drift.
