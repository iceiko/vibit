# Impact

This change closes the direction-confirmation gate after `ADR-0060` and selects the session repository boundary path.

No runtime behavior changes are introduced by this confirmation change.

The next work item may add:

- `docs/session-repository-boundary.md`
- `docs/session-repository-boundary.zh-CN.md`
- `decisions/ADR-0061-session-repository-boundary.md`

The next work item must not add:

- Go session repository code.
- PostgreSQL session adapters.
- Runtime session creation or validation.
- Route-policy use of persisted session or bound identity.
- WebSocket handshake authentication.
- Transport credential carriers.
- Protobuf envelope changes.
- Logout/revocation active-connection behavior.
- Reconnect/epoch behavior.
- Dependencies.
- Direct Nakama/Pitaya API compatibility.

Nakama informs the repository lifecycle capability pressure. Pitaya informs the separation between transport and session context.
