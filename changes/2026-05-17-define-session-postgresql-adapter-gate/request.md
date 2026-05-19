# Request

## Original Request

The maintainer asked the agent to recommend the next ten steps and continue them, with Nakama and Pitaya as important references.

## Clarified Requirement

Define a gate-only PostgreSQL session adapter standard after the storage-neutral session repository interface implementation. The gate should prepare the future adapter implementation while keeping runtime behavior unchanged.

## User-Visible Outcome

Agents now have an explicit standard for how a future PostgreSQL adapter may implement `runtime/internal/app/session.Repository`, including ownership, SQL shape, transaction handoff, error mapping, tests, and deferrals.

## Non-Goals

- No PostgreSQL session adapter Go files.
- No SQL execution behavior for `runtime_sessions`.
- No unit-of-work factory wiring for `NewSessionRepository`.
- No runtime session creation during login or `BindConnection`.
- No runtime session validation and no `RequestIdentity.SessionValidated = true`.
- No logout execution, refresh, cleanup jobs, logout-triggered active WebSocket invalidation, reconnect/resume/duplicate replacement, WebSocket handshake authentication, transport credential carrier, Protobuf session message, existing envelope change, generated output, dependency addition, memory durable session behavior, or direct Nakama/Pitaya API compatibility.

## Unknowns

- The first concrete PostgreSQL session adapter implementation remains a future choice.
- Runtime session validation behavior remains a future choice.
- Whether update/revoke zero affected rows should map to not found or stale state remains an implementation-gate detail.
- Permission and pagination semantics for player self-inspection versus admin session listing remain future choices.

## Acceptance Criteria

- [x] `docs/session-postgresql-adapter-gate.md` exists.
- [x] `docs/session-postgresql-adapter-gate.zh-CN.md` exists.
- [x] `ADR-0063` records the adapter gate decision.
- [x] Future adapter ownership, SQL shape, transaction boundary, error mapping, redaction, and tests are documented.
- [x] No PostgreSQL adapter, runtime validation, WebSocket, Protobuf, generated output, dependency, or direct Nakama/Pitaya API behavior is added.
- [x] `runtime.session_postgresql_adapter_gate` is registered as the repository check rule.
