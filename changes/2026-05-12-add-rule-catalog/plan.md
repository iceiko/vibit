# Plan

## Files To Create

- `rules/README.md`
- `rules/README.zh-CN.md`
- `rules/check-rules.json`
- `schema/rule-catalog.schema.json`
- `conversations/2026-05-12-rule-catalog.md`

## Files To Edit

- `tools/vibit`
- `.arch/conventions.yaml`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `docs/schema-validation.md`
- `docs/schema-validation.zh-CN.md`
- `changes/2026-05-12-add-rule-catalog/*`

## Generated Artifacts

- None.

## Handwritten Logic

- Add rule catalog loading and validation to `check schemas`.
- Check that all rule IDs known by `tools/vibit` exist in `rules/check-rules.json`.
- Keep validation dependency-free.

## Tests

- `node tools/vibit check schemas`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all`
- `node tools/vibit check all --json`
- Secret scan

## Verification Commands

- `node tools/vibit check schemas`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all`
- `node tools/vibit check all --json`
- `rg -n "ghp_[A-Za-z0-9]|github_pat_[A-Za-z0-9]" .`
- `git diff --check`

## Rollback Or Migration Notes

Remove the catalog files and schema checks if rule ID handling is replaced by generated metadata later.
