# Request

Define the gate-only boundary for future durable runtime session creation composition after runtime session validation implementation.

## Scope

Allowed:

- Add `docs/session-creation-composition-gate.md`.
- Add `docs/session-creation-composition-gate.zh-CN.md`.
- Add `ADR-0067`.
- Define future session creation composition ownership under `runtime/internal/app`.
- Define future login-time composition candidate using `AuthenticateWithDeviceCredential`.
- Define future use of `session.Repository.CreateRuntimeSession` through unit-of-work capabilities.
- Define session id, lifetime, token linkage, redaction, and future test expectations.
- Update manifests, ADR, conversation memory, AGENTS guides, rules, and checks.

Forbidden:

- Implement session creation.
- Modify authentication service behavior.
- Add session id generation code.
- Create sessions at login or `BindConnection`.
- Change runtime session validation behavior.
- Change route policy to use session or bound identity.
- Change WebSocket handshake authentication or transport credential carriers.
- Add Protobuf session messages, generated output, or envelope changes.
- Add logout/revocation active-connection invalidation.
- Add reconnect, resume, duplicate replacement, durable epoch behavior, cleanup jobs, dependencies, memory durable session behavior, or direct Nakama/Pitaya API compatibility.
