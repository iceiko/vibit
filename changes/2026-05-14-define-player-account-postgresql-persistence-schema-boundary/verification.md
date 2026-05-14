# Verification

Verified:

- `node tools/vibit check work --json`
- `node tools/vibit inspect next --json`
- `node tools/vibit check change define-player-account-postgresql-persistence-schema-boundary --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check contracts --json`
- `node tools/vibit check migrations --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- Runtime Go tests were covered indirectly by `node tools/vibit check runtime --json`; no Go runtime behavior was added.
- Live PostgreSQL verification is not required because no SQL migration source, repository adapter, or runtime persistence implementation is added.
