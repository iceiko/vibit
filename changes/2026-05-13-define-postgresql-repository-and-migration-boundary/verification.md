# Verification

Verified:

- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change define-postgresql-repository-and-migration-boundary --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- PostgreSQL migration apply/rollback; no migration files or disposable PostgreSQL environment exist yet.
- PostgreSQL repository integration tests; no repository adapter exists yet.

Not applicable:

- Go runtime behavior tests; this change does not modify Go runtime behavior.
- Protobuf generation; no protocol schema changes.
