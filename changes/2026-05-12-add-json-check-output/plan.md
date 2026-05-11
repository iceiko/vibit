# Plan

## Files To Create

- `conversations/2026-05-12-json-check-output.md`

## Files To Edit

- `tools/vibit`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `docs/schema-validation.md`
- `docs/schema-validation.zh-CN.md`
- `schema/inspect-output.schema.json`
- `changes/2026-05-12-add-json-check-output/*`

## Generated Artifacts

- None.

## Handwritten Logic

- Add `--json` argument parsing.
- Capture check results in memory.
- Preserve text output by default.
- Print JSON output when requested.
- Let `check all --json` include nested subcheck results.

## Tests

- Text `check all`
- JSON `check all`
- JSON `check schemas`
- JSON `check architecture`
- JSON `check change`
- JSON `check module`
- Secret scan

## Verification Commands

- `node tools/vibit check all`
- `node tools/vibit check all --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check architecture --json`
- `node tools/vibit check change add-json-check-output --json`
- `node tools/vibit check module inventory --json`
- `node -e 'JSON.parse(require("fs").readFileSync(0,"utf8"))'`
- `rg -n "ghp_[A-Za-z0-9]|github_pat_[A-Za-z0-9]" .`

## Rollback Or Migration Notes

Remove `--json` handling if the output contract proves premature, while preserving text checks.
