# Plan

## Files To Create

- `conversations/2026-05-12-inspect-rules.md`

## Files To Edit

- `tools/vibit`
- `schema/inspect-output.schema.json`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `changes/2026-05-12-add-inspect-rules/*`

## Generated Artifacts

- None.

## Handwritten Logic

- Parse optional `--category <category>` for `inspect rules`.
- Load `rules/check-rules.json`.
- Print a `rules_inspection` JSON object with filter metadata and matching rules.

## Tests

- Full rules inspection JSON parse check
- Category-filtered inspection JSON parse check
- Category filtering assertion
- `node tools/vibit check all`
- Secret scan

## Verification Commands

- `node tools/vibit inspect rules`
- `node tools/vibit inspect rules --category check`
- `node tools/vibit inspect rules --category check | node -e '<filter assertion>'`
- `node tools/vibit check all`
- `node tools/vibit check all --json`
- `rg -n "ghp_[A-Za-z0-9]|github_pat_[A-Za-z0-9]" . .git/config`
- `git diff --check`

## Rollback Or Migration Notes

Remove the inspect command if rule discovery is replaced by a broader query API.
