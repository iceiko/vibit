# Plan

## Files To Create

- `conversations/2026-05-12-inspect-rule.md`

## Files To Edit

- `tools/vibit`
- `schema/inspect-output.schema.json`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `changes/2026-05-12-add-inspect-rule/*`

## Generated Artifacts

- None.

## Handwritten Logic

- Load `rules/check-rules.json`.
- Find a rule by `rule_id`.
- Print a `rule_inspection` JSON object.
- Return a non-zero exit code with a clear message when the rule is unknown.

## Tests

- Existing rule inspection JSON parse check
- Missing rule failure check
- `node tools/vibit check all`
- Secret scan

## Verification Commands

- `node tools/vibit inspect rule check.subcheck`
- `node tools/vibit inspect rule check.subcheck | node -e '<JSON.parse assertion>'`
- `node tools/vibit inspect rule missing.rule`
- `node tools/vibit check all`
- `node tools/vibit check all --json`
- `rg -n "ghp_[A-Za-z0-9]|github_pat_[A-Za-z0-9]" .`
- `git diff --check`

## Rollback Or Migration Notes

Remove the inspect command if rule catalog access is replaced by a broader discovery API.
