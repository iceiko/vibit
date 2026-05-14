# Plan

## Files To Create

- `runtime/migrations/postgres/000003_create_authentication_device_credentials.sql`
- `changes/2026-05-14-add-credential-postgresql-migration-source/*`
- `conversations/2026-05-14-credential-postgresql-migration-source.md`

## Files To Edit

- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `docs/credential-token-session-schema-gates.md`
- `docs/credential-token-session-schema-gates.zh-CN.md`
- `docs/selected-login-token-boundary-checks.md`
- `docs/selected-login-token-boundary-checks.zh-CN.md`
- `docs/credential-record-schema-boundary.md`
- `docs/credential-record-schema-boundary.zh-CN.md`
- `docs/authentication-schema-migration-queue.md`
- `docs/authentication-schema-migration-queue.zh-CN.md`
- `docs/postgresql-persistence-boundary.md`
- `docs/postgresql-persistence-boundary.zh-CN.md`
- `tools/vibit`

## Generated Artifacts

- None.

## Handwritten Logic

- None. This change adds SQL migration source and repository check updates only.

## Tests

- No Go runtime behavior tests are required by the migration source itself.
- Existing runtime checks should still run and preserve the no-runtime-authentication, no-token-validation, no-handler, no-WebSocket-authentication boundary.
- Live PostgreSQL apply/rollback is not required for this source-only migration step unless a disposable DSN is explicitly supplied.

## Verification Commands

```bash
node tools/vibit check migrations --json
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change add-credential-postgresql-migration-source --json
node tools/vibit check all --json
git diff --check
```

## Rollback Or Migration Notes

Before the migration is applied in any shared environment, rollback is a normal source revert.

After shared application, do not edit this migration in place. Add a new migration instead.
