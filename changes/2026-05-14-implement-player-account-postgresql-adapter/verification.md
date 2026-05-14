# Verification

Verified:

- `cd runtime && go test ./internal/platform/persistence/postgres`
- `cd runtime && go test ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check module player --json`
- `node tools/vibit check contracts --json`
- `node tools/vibit check migrations --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change implement-player-account-postgresql-adapter --json`
- `node tools/vibit check all --json`
- `node tools/vibit inspect next --json`
- `git diff --check`

Not verified:

- Live PostgreSQL execution for the player account adapter. `VIBIT_POSTGRES_TEST_DSN` is optional and default repository checks must not require a running PostgreSQL server.

Not applicable:

- Migration application or rollback, because this change does not add or modify migration sources.
- Runtime WebSocket request-loop verification for player account commands, because runtime player handlers and WebSocket routes remain deferred.
