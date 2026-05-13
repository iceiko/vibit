# Verification

Verified:

- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `node tools/vibit check migrations --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change add-postgresql-migration-apply-and-status-tooling --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan for GitHub token patterns in tracked changes.

Not verified:

- Live PostgreSQL migration apply/status execution. It remains deferred until `W-0019` defines the disposable PostgreSQL verification environment.

Not applicable:

- Runtime startup migration execution; this work item intentionally avoids startup side effects.
