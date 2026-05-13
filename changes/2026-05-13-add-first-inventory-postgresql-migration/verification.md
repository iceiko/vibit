# Verification

Verified:

- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change add-first-inventory-postgresql-migration --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan for GitHub token patterns in tracked and untracked committable files.

Not verified:

- None.

Not applicable:

- PostgreSQL migration apply/rollback; migration verification tooling is planned for `W-0014`, and no disposable PostgreSQL test environment standard exists yet.
- PostgreSQL repository integration tests; no PostgreSQL adapter exists yet.
- Protobuf generation; no protocol schema changes.
