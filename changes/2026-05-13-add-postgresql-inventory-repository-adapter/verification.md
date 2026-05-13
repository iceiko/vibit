# Verification

Verified:

- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check migrations --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change add-postgresql-inventory-repository-adapter --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan for GitHub token patterns in tracked and untracked committable files.

Not verified:

- None.

Not applicable:

- PostgreSQL live integration tests; no local disposable PostgreSQL test environment standard exists yet.
- PostgreSQL migration apply/rollback execution; migration apply tooling is not implemented yet.
- Runtime PostgreSQL process wiring; this adapter is not wired into startup yet.
- Protobuf generation; no protocol schema changes.
