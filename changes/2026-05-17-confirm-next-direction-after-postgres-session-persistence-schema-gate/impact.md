# Impact

This change closes the direction-confirmation gate after `ADR-0059` and selects the migration-source-only path.

No runtime behavior changes are introduced by this confirmation change.

The next work item may add:

- `runtime/migrations/postgres/000005_create_runtime_sessions.sql`

The next work item must not add:

- Session repositories.
- PostgreSQL session adapters.
- Runtime session validation.
- Route-policy use of persisted session or bound identity.
- WebSocket handshake authentication.
- Transport credential carriers.
- Protobuf envelope changes.
- Logout/revocation active-connection behavior.
- Reconnect/epoch behavior.
- Dependencies.
- Direct Nakama/Pitaya API compatibility.

Nakama informs the lifecycle-record priority. Pitaya informs the separation between transport and session context.
