# Verification

Verified:

- `cd runtime && go test ./internal/platform/persistence/postgres`
- `cd runtime && go test ./...`
- `node -c tools/vibit`
- `node tools/vibit check runtime --json`
- `node tools/vibit check module authentication --json`
- `node tools/vibit check contracts --json`
- `node tools/vibit check migrations --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change implement-authentication-postgresql-adapter --json`
- `node tools/vibit inspect next --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan before commit: `rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" . .git/config`

Pending for final repository state:

- None.

Not verified:

- Live PostgreSQL execution for the authentication adapter. `VIBIT_POSTGRES_TEST_DSN` is optional and default repository checks must not require a running PostgreSQL server.

Warnings:

- `node tools/vibit check runtime --json` reports one `runtime.identity_boundary` warning because the bounded PostgreSQL authentication adapter necessarily mentions credential vocabulary. The adapter remains inside the explicit ratified persistence boundary.

Not applicable:

- Migration application or rollback, because this change does not add or modify migration sources.
- Runtime authentication request-loop verification, because runtime authentication handlers, token validators, Protobuf messages, WebSocket routes, and generated authentication shapes remain deferred.
