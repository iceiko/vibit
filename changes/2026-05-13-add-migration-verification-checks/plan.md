# Plan

## Files To Create

- `changes/2026-05-13-add-migration-verification-checks/*`

## Files To Edit

- `tools/vibit`
- `rules/check-rules.json`
- `.arch/conventions.yaml`
- `.arch/runtime.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `docs/schema-validation.md`
- `docs/schema-validation.zh-CN.md`
- `docs/postgresql-persistence-boundary.md`
- `docs/postgresql-persistence-boundary.zh-CN.md`
- `.arch/work-items.yaml`

## Generated Artifacts

- None.

## Handwritten Logic

Add a focused Node.js CLI check for PostgreSQL migration source files.

## Tests

- No Go tests are required because runtime behavior is unchanged.
- CLI verification is covered by direct `node tools/vibit check migrations --json` and aggregate `check all`.

## Verification Commands

```bash
node tools/vibit check migrations --json
node tools/vibit check schemas --json
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change add-migration-verification-checks --json
node tools/vibit check all --json
cd runtime && go test ./...
cd runtime && go vet ./...
git diff --check
```

## Rollback Or Migration Notes

This change adds source validation only. Revert the tooling and documentation changes if the check shape needs to be replaced.
