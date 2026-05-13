# Plan

## Files To Create

- `docs/postgresql-verification-environment.md`
- `docs/postgresql-verification-environment.zh-CN.md`

## Files To Edit

- `tools/vibit`
- `rules/check-rules.json`
- `.arch/conventions.yaml`
- `.arch/runtime.yaml`
- `.arch/work-items.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `docs/runtime-runbook.md`
- `docs/runtime-runbook.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/inventory/AGENTS.md`
- `modules/inventory/AGENTS.zh-CN.md`
- `modules/inventory/module.yaml`
- This change spec.

## Generated Artifacts

- None.

## Handwritten Logic

Add a static `check postgres-env` CLI command that verifies environment-standard documents, manifest references, and guidance artifacts. The command must not connect to PostgreSQL or inspect local service managers.

## Tests

Use repository checks to validate the new command and guidance references.

## Verification Commands

- `node tools/vibit check postgres-env --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change define-disposable-postgresql-verification-environment --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

Remove the new standard documents, CLI command, rule metadata, manifest references, and guidance updates. Move `W-0019` back to `next_ready`.
