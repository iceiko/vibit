# Plan

## Files To Create

- `runtime/migrations/postgres/000002_create_player_account_state.sql`
- `changes/2026-05-14-add-first-player-account-postgresql-migration/*`

## Files To Edit

- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `modules/player/module.yaml`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `docs/postgresql-persistence-boundary.md`
- `docs/postgresql-persistence-boundary.zh-CN.md`
- `tools/vibit`
- `.arch/work-items.yaml`

## Generated Artifacts

- None.

## Handwritten Logic

- None. This change adds SQL migration source and verification metadata only.

## Tests

- No Go runtime behavior tests are required by the migration source itself.
- Existing runtime checks should still run and preserve the no-authentication, no-session, no-handler boundary.
- Live PostgreSQL apply/rollback is not required for this source-only migration step unless a disposable DSN is explicitly supplied.

## Verification Commands

```bash
node tools/vibit check migrations --json
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change add-first-player-account-postgresql-migration --json
node tools/vibit check all --json
git diff --check
```

## Rollback Or Migration Notes

Before the migration is applied in any shared environment, rollback is a normal source revert.

After shared application, do not edit this migration in place. Add a new migration instead.
