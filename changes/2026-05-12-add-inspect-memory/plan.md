# Plan

## Files To Create

- `changes/2026-05-12-add-inspect-memory/`
- `conversations/2026-05-12-inspect-memory.md`

## Files To Edit

- `tools/vibit`
- `schema/inspect-output.schema.json`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`

## Generated Artifacts

- None.

## Handwritten Logic

- Reuse change inspection data for change memory entries.
- Read conversation log headings, dates, related changes, and related artifacts.
- Read Agent Decision Record headings, status, date, related changes, related conversations, and related artifacts.
- Return counts and intake guidance as JSON.

## Tests

- Memory inspection JSON parse check
- Memory count assertion
- `node tools/vibit check all`
- Secret scan
- Diff whitespace check

## Verification Commands

- `node tools/vibit inspect memory`
- `node tools/vibit inspect memory | node -e '<JSON.parse assertion>'`
- `node tools/vibit --help | rg "inspect memory"`
- `node tools/vibit check all`
- `node tools/vibit check all --json`
- `rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" . .git/config`
- `git diff --check`

## Rollback Or Migration Notes

Remove the inspect command if project memory indexing moves into a broader artifact registry.
