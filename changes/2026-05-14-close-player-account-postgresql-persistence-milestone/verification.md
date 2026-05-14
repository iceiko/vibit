# Verification

Verified:

- `cd runtime && go test ./internal/platform/persistence/postgres`
- `cd runtime && go test ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check module player --json`
- `node tools/vibit check contracts --json`
- `node tools/vibit check migrations --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change close-player-account-postgresql-persistence-milestone --json`
- `node tools/vibit check all --json`
- `node tools/vibit inspect next --json`
- `git diff --check`

Not verified:

- Live PostgreSQL execution for the player account adapter. Default repository checks do not require a running PostgreSQL server, and `VIBIT_POSTGRES_TEST_DSN` was not required for this milestone closeout.

Not applicable:

- New runtime WebSocket request-loop verification for player account commands, because runtime player account handlers and WebSocket routes remain deferred.
- New migration apply or rollback verification, because this change does not add or modify migration sources.
