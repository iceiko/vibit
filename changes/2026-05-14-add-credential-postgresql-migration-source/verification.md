# Verification

Verified:
- `node tools/vibit inspect next --json`
- `node tools/vibit check migrations --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check architecture --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check change add-credential-postgresql-migration-source --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan excluding `.git`, `.vibit.local.env`, and `node_modules`

Not verified:
- Live PostgreSQL apply/rollback for `runtime/migrations/postgres/000003_create_authentication_device_credentials.sql`; this work item adds source only and no disposable DSN was required.

Not applicable:
- Runtime authentication tests, because no credential lookup, login, token issuance, token validation, logout, refresh, cleanup, handler, route, or authentication behavior was added.
- Authentication repository or adapter tests, because no repository interface or PostgreSQL adapter was added.
- Protobuf generation, because no Protobuf source or generated Protobuf output changed.
- WebSocket behavior tests, because no WebSocket handshake, proof carrier, route, or protocol behavior changed.
