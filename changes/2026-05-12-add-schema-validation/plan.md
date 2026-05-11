# Plan

## Files To Create

- `docs/schema-validation.md`
- `docs/schema-validation.zh-CN.md`
- `schema/module-manifest.schema.json`
- `schema/change-spec.schema.json`
- `schema/agent-decision-record.schema.json`
- `schema/inspect-output.schema.json`
- `conversations/2026-05-12-schema-validation-start.md`

## Files To Edit

- `tools/vibit`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `.arch/conventions.yaml`
- `changes/2026-05-12-add-schema-validation/*`

## Generated Artifacts

- None.

## Handwritten Logic

- Add `check schemas`.
- Add no-dependency critical-field validation for:
  - `.arch/conventions.yaml`
  - `.arch/modules.yaml`
  - `modules/<module>/module.yaml`
  - `changes/<change>/spec.yaml`
  - `decisions/ADR-*.md`
- Add schema file existence and JSON parse checks.

## Tests

- CLI schema check
- Aggregate check
- Inspect commands
- Secret scan

## Verification Commands

- `node tools/vibit check schemas`
- `node tools/vibit check all`
- `node tools/vibit inspect module inventory`
- `node tools/vibit inspect boundary --from inventory --to player`
- `rg -n "ghp_[A-Za-z0-9]|github_pat_[A-Za-z0-9]" .`

## Rollback Or Migration Notes

Schema validation can be relaxed or moved behind a warning mode if early strictness blocks valid design evolution.
