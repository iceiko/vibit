# Request

The maintainer asked the agent to recommend the next ten steps and continue, with Nakama and Pitaya as important references.

## Selected Direction

Select `implement_session_postgresql_adapter` after the session PostgreSQL adapter gate.

## Rationale

The repository already has:

- A `runtime_sessions` PostgreSQL migration source.
- A storage-neutral `runtime/internal/app/session.Repository` interface.
- A gate-only PostgreSQL adapter standard.

The next bounded step is therefore the platform persistence adapter and unit-of-work factory. This follows Nakama's lesson that game session lifecycle records need durable lookup, expiration, revocation, and listing, while following Pitaya's lesson that transport acceptors and route handlers should not own durable session persistence.

## Non-Goals

- Runtime session creation at login or connection binding.
- Runtime session validation.
- Setting `RequestIdentity.SessionValidated` true.
- WebSocket handshake authentication.
- Transport credential carriers.
- Protobuf session messages or envelope changes.
- Route-policy use of persisted session or bound identity.
- Logout-triggered active connection invalidation.
- Reconnect, resume, duplicate replacement, or durable epoch behavior.
- Cleanup jobs, dependencies, memory durable session behavior, or direct Nakama/Pitaya API compatibility.
