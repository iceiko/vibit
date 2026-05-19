# ADR-0064: Session PostgreSQL Adapter Implementation

Status: Accepted
Date: 2026-05-17
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-session-postgresql-adapter-gate/`
- `changes/2026-05-17-implement-session-postgresql-adapter/`

Related conversations:

- `conversations/2026-05-17-session-postgresql-adapter-implementation.md`

Related artifacts:

- `runtime/internal/platform/persistence/postgres/session_repository.go`
- `runtime/internal/platform/persistence/postgres/session_repository_test.go`
- `runtime/internal/platform/persistence/postgres/runner.go`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`W-0138` defined the session PostgreSQL adapter gate after the storage-neutral runtime session repository interface. The repository now has a `runtime_sessions` migration source, `runtime/internal/app/session.Repository`, and a ratified adapter boundary.

The work queue reached `M-067/W-0139`, a confirmation gate. The maintainer asked the agent to recommend the next ten steps and continue, with Nakama and Pitaya as important references.

## Decision

Select:

```text
implement_session_postgresql_adapter
```

Implement the bounded PostgreSQL adapter in:

```text
runtime/internal/platform/persistence/postgres/session_repository.go
runtime/internal/platform/persistence/postgres/session_repository_test.go
```

Expose the adapter through the PostgreSQL unit-of-work factory:

```text
UnitOfWork.NewSessionRepository() (session.Repository, error)
```

The adapter implements only `runtime/internal/app/session.Repository` methods and maps them to the existing `runtime_sessions` table. It remains persistence-only and transaction-bound.

This ADR does not add runtime session creation at login or `BindConnection`, runtime session validation, `RequestIdentity.SessionValidated = true`, WebSocket handshake authentication, transport credential carriers, Protobuf session messages, existing envelope changes, generated output, route-policy use of persisted session or bound identity, logout-triggered active connection invalidation, reconnect/resume/duplicate replacement/epoch behavior, cleanup jobs, dependencies, memory durable session behavior, or direct Nakama/Pitaya public API compatibility.

## Alternatives Considered

- Define runtime session validation before the adapter exists.
- Attach persisted session creation directly to login or connection binding in the same slice.
- Make persisted session identity satisfy ordinary protected route policy immediately.
- Add logout/revocation active-connection invalidation in the persistence adapter slice.
- Add a connection registry table with the session adapter.
- Copy Nakama or Pitaya session APIs directly.

## Rationale

Nakama shows that game sessions quickly need durable lifecycle operations: lookup, expiration, revocation, logout support, and management-ready listing. Pitaya shows that acceptors, session context, and route handlers should stay separated, so durable session persistence must not live in WebSocket transport or protocol routing.

The right vibit adaptation is a PostgreSQL adapter behind the existing storage-neutral repository interface. It gives later runtime validation work a durable persistence dependency without mixing proof validation, request identity construction, route policy, socket lifecycle, reconnect, or social/realtime behavior into the adapter.

## Agent Reasoning Summary

The agent treated Nakama and Pitaya as design pressure, not compatibility targets. Nakama supports prioritizing a durable session record adapter because authentication, session listing, expiration, and revocation become operational requirements quickly in real game backends. Pitaya supports keeping persistence below handler routing and transport acceptor concerns, so the adapter remains behind the application repository interface and unit-of-work boundary.

The selected implementation is the smallest code slice that unlocks later runtime session validation work without prematurely coupling login, WebSocket binding, route authorization, reconnect, or logout behavior to storage. The repository adapter can be verified with deterministic SQL-shape tests and error mapping tests before any runtime path depends on persisted sessions.

## Decision Weights

- Durable session lifecycle readiness: high.
- Separation between transport, protocol routing, application behavior, and persistence: high.
- Compatibility with existing repository and unit-of-work boundaries: high.
- Avoiding premature runtime authentication/session semantics: high.
- Direct Nakama/Pitaya API compatibility: low.
- Live database integration in this slice: low.

## Consequences

- `SessionRepository` implements `session.Repository`.
- The adapter queries and mutates only `runtime_sessions`.
- Unit-of-work code can create a session repository inside PostgreSQL transactions.
- Focused fake-executor tests cover SQL shape, row mapping, error mapping, nullable fields, UTC normalization, no live PostgreSQL requirement, and no repository-owned transaction control.
- `runtime.session_postgresql_adapter_implementation` becomes the repository check rule for this slice.
- The work queue blocks again after implementation at `M-069/W-0141`.

## Reversal Conditions

This decision should be revisited if the runtime session repository contract changes, if `runtime_sessions` is replaced by another first durable session store, if unit-of-work ownership moves out of the PostgreSQL persistence package, if later validation requires additional atomic operations that cannot be added without changing the adapter boundary, or if a future explicitly ratified clustering/session model requires a different storage abstraction.

## Follow-Up

- Define a runtime session validation behavior gate before any code sets `RequestIdentity.SessionValidated` true.
- Define login/session-creation composition before login creates durable session rows.
- Define logout/revocation active-connection behavior before logout can disconnect or invalidate bound WebSocket connections.
- Define reconnect and connection epoch behavior before duplicate connection replacement or resume behavior.
