# ADR-0061: Session Repository Boundary

Status: Accepted
Date: 2026-05-17
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-runtime-sessions-migration-source/`
- `changes/2026-05-17-define-session-repository-boundary/`

Related conversations:

- `conversations/2026-05-17-session-repository-boundary.md`

Related artifacts:

- `docs/session-repository-boundary.md`
- `docs/session-repository-boundary.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`W-0132` added the PostgreSQL `runtime_sessions` migration source only. The table now exists as a SQL source, but the repository boundary, adapter behavior, session validation behavior, WebSocket behavior, logout/revocation behavior, reconnect behavior, and route-policy use of session identity remain deferred.

The work queue reached `M-061/W-0133`, a confirmation gate after the migration-source slice. The maintainer asked the agent to recommend the next ten steps and continue, while continuing to use Nakama and Pitaya as important references.

## Decision

Select:

```text
define_session_repository_boundary
```

Define a gate-only session repository boundary in:

```text
docs/session-repository-boundary.md
docs/session-repository-boundary.zh-CN.md
```

The future repository owner candidate is:

```text
runtime/internal/app/session
```

The future PostgreSQL adapter owner remains:

```text
runtime/internal/platform/persistence/postgres
```

This ADR does not add Go session repository interfaces, PostgreSQL session adapters, runtime session creation, runtime session validation, runtime session revocation execution, cleanup jobs, WebSocket handshake authentication, transport credential carriers, route-policy use of sessions or bound identity, Protobuf envelope changes, generated output, logout-triggered active WebSocket invalidation, reconnect/epoch behavior, dependencies, memory durable session behavior, or direct Nakama/Pitaya public API compatibility.

## Alternatives Considered

- Implement the session repository interface immediately.
- Implement the PostgreSQL session adapter immediately.
- Start runtime session validation before repository boundaries are documented.
- Attach session creation to login or `BindConnection` in the same slice.
- Make persisted session identity satisfy ordinary protected route policy.
- Start logout/revocation active-connection behavior before repository and validation boundaries.
- Adopt Nakama or Pitaya session APIs directly.

## Rationale

Nakama shows that session lifecycle lookup, expiration, refresh, logout, and operational management become central quickly in game servers. Pitaya shows that session context should be available to handler logic without making acceptors and transport own durable session behavior.

The right vibit adaptation is to define a storage-neutral application repository boundary before adding a Go interface or adapter. That gives later agents a clear owner, method vocabulary, forbidden data flows, and deferral list.

## Agent Reasoning Summary

After adding `runtime_sessions`, the smallest useful next step is not runtime behavior. It is the repository boundary that will keep future behavior from mixing session storage, token validation, WebSocket connection state, and route policy into one hidden subsystem.

## Decision Weights

```yaml
decision_weights:
  nakama_pitaya_alignment: high
  repository_safety: high
  agent_readability: high
  behavior_change_risk: low
  dependency_expansion: low
confidence: high
```

## Consequences

- `docs/session-repository-boundary.md` and `docs/session-repository-boundary.zh-CN.md` exist.
- `runtime.session_repository_boundary` becomes the repository check rule for this slice.
- `M-061/W-0133` is closed with `define_session_repository_boundary`.
- `M-062/W-0134` is completed as a gate-only milestone.
- The work queue blocks again at `M-063/W-0135` before adapters, runtime validation, route policy, logout/revocation, reconnect, operations, or broader game backend work.
- Existing runtime behavior is not changed by this ADR.

## Reversal Conditions

Revisit this decision if:

- A future ADR replaces `runtime_sessions` with a different first session persistence model.
- A future runtime session implementation needs an owner outside `runtime/internal/app/session`.
- The first production posture requires session validation before any repository boundary.
- A future ADR chooses direct Nakama or Pitaya public API compatibility.

## Follow-Up

- Define the PostgreSQL session adapter gate.
- Define the concrete session repository interface implementation slice.
- Define runtime session validation behavior before setting `RequestIdentity.SessionValidated` true.
- Define logout/revocation active-connection behavior before logout can close or invalidate bound WebSocket connections.
- Define reconnect and connection epoch behavior before duplicate connection replacement or resume behavior.
