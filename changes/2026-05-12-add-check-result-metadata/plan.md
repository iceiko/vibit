# Plan

## Files To Create

- `conversations/2026-05-12-check-result-metadata.md`

## Files To Edit

- `tools/vibit`
- `schema/inspect-output.schema.json`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `changes/2026-05-12-add-check-result-metadata/*`

## Generated Artifacts

- None.

## Handwritten Logic

- Add a shared result recording helper.
- Add `rule_id` and `artifact` to result entries.
- Pass explicit rule metadata from existing check helpers.
- Add aggregate subcheck result entries to `check all --json`.

## Tests

- Text `check all`
- JSON `check all`
- JSON metadata assertion across nested subchecks
- JSON `check schemas`
- JSON `check architecture`
- JSON `check change`
- JSON `check module`
- Secret scan

## Verification Commands

- `node tools/vibit check all`
- `node tools/vibit check all --json`
- `node tools/vibit check all --json | node -e '<metadata assertion>'`
- `node tools/vibit check schemas --json`
- `node tools/vibit check architecture --json`
- `node tools/vibit check change add-check-result-metadata --json`
- `node tools/vibit check module inventory --json`
- `rg -n "ghp_[A-Za-z0-9]|github_pat_[A-Za-z0-9]" .`

## Rollback Or Migration Notes

Remove the added result metadata fields and aggregate subcheck result entries if the result contract needs to be redesigned before stabilization.
