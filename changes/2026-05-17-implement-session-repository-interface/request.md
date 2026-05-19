# Request

## Original Request

The maintainer asked the agent to recommend the next ten steps and continue them, with Nakama and Pitaya as important references.

## Clarified Requirement

Implement the selected storage-neutral runtime session repository interface under `runtime/internal/app/session`, with lifecycle value types and normalization tests only.

## User-Visible Outcome

Agents now have a concrete Go interface and value vocabulary for future runtime session lifecycle storage, inspired by Nakama's first-class session lifecycle and Pitaya's separation between session context and transport.

## Non-Goals

- No PostgreSQL session adapter or SQL query implementation.
- No unit-of-work factory wiring for session repositories.
- No runtime session creation during login or `BindConnection`.
- No runtime session validation and no `RequestIdentity.SessionValidated = true`.
- No logout execution, refresh, cleanup jobs, logout-triggered active WebSocket invalidation, reconnect/resume/duplicate replacement, WebSocket handshake authentication, transport credential carrier, Protobuf session message, existing envelope change, generated output, dependency addition, memory durable session behavior, or direct Nakama/Pitaya API compatibility.

## Unknowns

- The first session PostgreSQL adapter gate remains a future choice.
- Runtime session validation behavior remains a future choice.
- Permission and pagination semantics for player self-inspection versus admin session listing remain future choices.

## Acceptance Criteria

- [x] `runtime/internal/app/session/repository.go` defines storage-neutral session lifecycle value types and `Repository`.
- [x] `runtime/internal/app/session/repository_test.go` verifies closed status vocabulary, normalization, UTC times, listing bounds, and absence of secret material fields.
- [x] No PostgreSQL adapter, WebSocket transport, Protobuf, generated output, startup, or runtime validation behavior is added.
- [x] `runtime.session_repository_interface_implementation` is registered as the repository check rule.
