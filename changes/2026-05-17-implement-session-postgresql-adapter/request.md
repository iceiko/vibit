# Request

Implement the PostgreSQL adapter for `runtime/internal/app/session.Repository` after the session PostgreSQL adapter gate.

## Scope

Allowed:

- Add `runtime/internal/platform/persistence/postgres/session_repository.go`.
- Add `runtime/internal/platform/persistence/postgres/session_repository_test.go`.
- Add `UnitOfWork.NewSessionRepository()`.
- Add focused fake-executor tests for SQL shape, row mapping, error mapping, and unit-of-work factory wiring.
- Update manifests, ADR, conversation memory, AGENTS guides, rules, and checks.

Forbidden:

- Runtime session creation at login or `BindConnection`.
- Runtime session validation.
- Setting `RequestIdentity.SessionValidated` true.
- WebSocket handshake authentication or transport credential carriers.
- Protobuf session messages or envelope changes.
- Route-policy use of persisted session or bound identity.
- Logout-triggered active connection invalidation.
- Reconnect, resume, duplicate replacement, durable epoch behavior, cleanup jobs, dependencies, memory durable session behavior, or direct Nakama/Pitaya API compatibility.
