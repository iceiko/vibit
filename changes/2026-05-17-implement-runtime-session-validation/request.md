# Request

Implement the application-owned runtime session validator after the runtime session validation gate.

## Scope

Allowed:

- Add `runtime/internal/app/runtime_session_validator.go`.
- Add `runtime/internal/app/runtime_session_validator_test.go`.
- Depend only on the storage-neutral `runtime/internal/app/session.Repository`.
- Look up active persisted sessions through `FindActiveSessionByID`.
- Require an already-validated player identity before trusting the persisted session row.
- Set `RequestIdentity.SessionValidated = true` only after durable active-session validation succeeds.
- Collapse public invalid-session failures to a stable redacted reason.
- Add focused fake-repository tests.
- Update manifests, ADR, conversation memory, AGENTS guides, rules, and checks.

Forbidden:

- Runtime session creation at login or `BindConnection`.
- Route-policy use of session-validated or bound identity.
- WebSocket handshake authentication or transport credential carriers.
- Protobuf session messages or envelope changes.
- Logout-triggered active connection invalidation.
- Reconnect, resume, duplicate replacement, durable epoch behavior, cleanup jobs, dependencies, memory durable session behavior, or direct Nakama/Pitaya API compatibility.
