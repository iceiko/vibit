# Verification

Verified:

- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change add-postgresql-configuration-and-transaction-runner --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan for GitHub token patterns in tracked and untracked committable files.

Not verified:

- None.

Not applicable:

- Live PostgreSQL integration tests; this work item intentionally uses fake transaction handles until the disposable PostgreSQL environment standard exists.
- Migration apply/rollback; planned for `W-0018`.
- Runtime PostgreSQL process wiring; planned for `W-0020`.
