# Impact

This change adds one SQL migration source:

- `runtime/migrations/postgres/000005_create_runtime_sessions.sql`

The migration creates the `runtime_sessions` table with:

- Server-generated `session_id`.
- Player actor identity fields.
- Session lifecycle status.
- Issued, expiration, last-seen, created, and updated timestamps.
- Optional revocation timestamp and reason.
- Optional access-token record linkage.
- Foreign key to `player_accounts(player_id)`.
- Optional foreign key to `authentication_access_tokens(token_record_id)`.

This is intentionally smaller than Nakama or Pitaya runtime session behavior. It borrows the durable lifecycle concept from Nakama and the session/transport separation from Pitaya, while keeping vibit's table, ownership, and verification rules native.

No Go code, generated code, Protobuf source, WebSocket behavior, repository interface, PostgreSQL adapter, dependency, route policy, logout/revocation behavior, reconnect behavior, or direct compatibility API is added.
