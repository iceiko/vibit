# Plan

## Files To Create

- `changes/2026-05-12-add-inspect-change/`
- `conversations/2026-05-12-inspect-change.md`

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

- Reuse change ID validation and change directory discovery.
- Read `spec.yaml` with existing lightweight scalar and list helpers.
- Return file presence for required change spec files.
- Return `exists: false` for missing change directories instead of requiring agents to infer that from prose.

## Tests

- Existing change inspection JSON parse check
- Missing change inspection JSON parse check
- `node tools/vibit check all`
- Secret scan
- Diff whitespace check

## Verification Commands

- `node tools/vibit inspect change add-inspect-change`
- `node tools/vibit inspect change add-inspect-change | node -e '<JSON.parse assertion>'`
- `node tools/vibit inspect change missing-change | node -e '<JSON.parse assertion>'`
- `node tools/vibit check all`
- `node tools/vibit check all --json`
- `rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" . .git/config`
- `git diff --check`

## Rollback Or Migration Notes

Remove the inspect command if change intake is replaced by a broader structured work-item API.
