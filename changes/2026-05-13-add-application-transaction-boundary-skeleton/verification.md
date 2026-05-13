# Verification

Verified:

- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change add-application-transaction-boundary-skeleton --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- PostgreSQL migration apply/rollback; no migration files are added in this change.
- PostgreSQL repository integration tests; no PostgreSQL adapter exists yet.
- Protobuf generation; no protocol schema changes.
