# Plan

## Files To Create

- `changes/2026-05-12-add-check-memory/`
- `conversations/2026-05-12-check-memory.md`

## Files To Edit

- `tools/vibit`
- `rules/check-rules.json`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`

## Generated Artifacts

- None.

## Handwritten Logic

- Add `check memory`.
- Check conversation headings, dates, and required sections.
- Check ADR headings, status, dates, ID consistency, and required sections.
- Register new rule IDs used by the check output.
- Add `check memory` to `check all`.

## Tests

- `check memory` text output
- `check memory --json` parse and status assertion
- `check all` text output
- `check all --json` parse and subcheck assertion
- Secret scan
- Diff whitespace check

## Verification Commands

- `node tools/vibit check memory`
- `node tools/vibit check memory --json | node -e '<JSON.parse assertion>'`
- `node tools/vibit check all`
- `node tools/vibit check all --json | node -e '<JSON.parse assertion>'`
- `rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" . .git/config`
- `git diff --check`

## Rollback Or Migration Notes

Remove `check memory` from `check all` if the project memory standard changes before the check can be migrated.
