# Verification

Verified:

- `node tools/vibit check migrations --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change add-migration-verification-checks --json`
- `node tools/vibit check all --json`
- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `git diff --check`
- Secret scan for GitHub token patterns in tracked and untracked committable files.

Not verified:

- None.

Not applicable:

- PostgreSQL migration apply/rollback execution; this change validates migration source files only.
- PostgreSQL repository integration tests; no PostgreSQL adapter exists yet.
- Protobuf generation; no protocol schema changes.
