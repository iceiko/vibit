# Impact

This change adds a gate-only standard for the future runtime session repository boundary.

Files added:

- `docs/session-repository-boundary.md`
- `docs/session-repository-boundary.zh-CN.md`
- `decisions/ADR-0061-session-repository-boundary.md`
- `conversations/2026-05-17-session-repository-boundary.md`

No Go runtime behavior changes are introduced.

Future work gains:

- A storage-neutral repository owner candidate.
- Candidate repository capability vocabulary.
- Data boundary rules for `runtime_sessions`.
- Clear separation from authentication token validation.
- Clear separation from WebSocket transport and Protobuf protocol behavior.

The change explicitly does not add:

- Go repository interfaces.
- PostgreSQL session adapters.
- Runtime session creation, validation, revocation, or cleanup.
- Route-policy use of session identity.
- WebSocket handshake authentication.
- Transport credential carriers.
- Protobuf session messages or envelope changes.
- Logout/revocation active-connection behavior.
- Reconnect/epoch behavior.
- Dependencies.
- Direct Nakama/Pitaya API compatibility.

Nakama informs the need for queryable session lifecycle capabilities. Pitaya informs the boundary between session context and transport acceptors.
