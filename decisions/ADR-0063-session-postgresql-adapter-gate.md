# ADR-0063: Session PostgreSQL Adapter Gate

Status: Accepted
Date: 2026-05-17
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-session-repository-interface/`
- `changes/2026-05-17-define-session-postgresql-adapter-gate/`

Related conversations:

- `conversations/2026-05-17-session-postgresql-adapter-gate.md`

Related artifacts:

- `docs/session-postgresql-adapter-gate.md`
- `docs/session-postgresql-adapter-gate.zh-CN.md`
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

`W-0136` implemented the storage-neutral runtime session repository interface. The repository now has both a PostgreSQL `runtime_sessions` migration source and an application-owned `runtime/internal/app/session.Repository` interface, but no PostgreSQL adapter, no unit-of-work factory wiring, and no runtime session behavior.

The work queue reached `M-065/W-0137`, a confirmation gate. The maintainer asked the agent to recommend the next ten steps and continue, while continuing to use Nakama and Pitaya as important references.

## Decision

Select:

```text
define_session_postgresql_adapter_gate
```

Define a gate-only PostgreSQL adapter standard in:

```text
docs/session-postgresql-adapter-gate.md
docs/session-postgresql-adapter-gate.zh-CN.md
```

The future adapter owner remains:

```text
runtime/internal/platform/persistence/postgres
```

The storage-neutral repository owner remains:

```text
runtime/internal/app/session
```

This ADR does not add PostgreSQL adapter code, SQL execution behavior, unit-of-work factory wiring, runtime session creation, runtime session validation, runtime session revocation execution, cleanup jobs, WebSocket handshake authentication, transport credential carriers, route-policy use of sessions or bound identity, Protobuf session messages, existing envelope changes, generated output, logout-triggered active WebSocket invalidation, reconnect/epoch behavior, dependencies, memory durable session behavior, or direct Nakama/Pitaya public API compatibility.

## Alternatives Considered

- Implement the PostgreSQL session adapter immediately.
- Start runtime session validation before the adapter gate.
- Attach session creation to login or `BindConnection` in the same slice.
- Make persisted session identity satisfy ordinary protected route policy.
- Add logout/revocation active-connection invalidation before adapter and validation semantics are defined.
- Adopt Nakama or Pitaya public APIs directly.
- Move session persistence into WebSocket transport or protocol adapter code.

## Rationale

Nakama shows that session storage quickly needs durable lifecycle lookup, expiration, revocation, logout, and operational listing. Pitaya shows that session context belongs in handler-facing application context, not in transport acceptors or routing plumbing.

The right vibit adaptation is to define the PostgreSQL adapter gate before implementing the adapter. The gate documents transaction handoff, SQL shape, error mapping, test expectations, redaction rules, and non-goals so later agents do not accidentally mix persistence, token proof validation, request identity construction, route policy, and WebSocket connection behavior.

## Agent Reasoning Summary

After the storage-neutral interface, the next conservative step is not runtime validation. A validation gate needs a reliable persistence adapter boundary first. Defining the adapter gate keeps the next implementation narrow and makes the later runtime validation work easier to review.

## Decision Weights

```yaml
decision_weights:
  nakama_pitaya_alignment: high
  adapter_safety: high
  agent_readability: high
  behavior_change_risk: low
  dependency_expansion: low
confidence: high
```

## Consequences

- `docs/session-postgresql-adapter-gate.md` and `docs/session-postgresql-adapter-gate.zh-CN.md` exist.
- `runtime.session_postgresql_adapter_gate` becomes the repository check rule for this slice.
- `M-065/W-0137` is closed with `define_session_postgresql_adapter_gate`.
- `M-066/W-0138` is completed as a gate-only milestone.
- The work queue blocks again at `M-067/W-0139` before adapter implementation, runtime validation, route policy, logout/revocation, reconnect, operations, memory durable session behavior, or broader game backend work.
- Existing runtime behavior is not changed by this ADR.

## Reversal Conditions

Revisit this decision if:

- A future ADR replaces PostgreSQL as the first durable runtime session store.
- A future adapter implementation cannot satisfy `runtime/internal/app/session.Repository` without interface changes.
- The first production posture requires a non-SQL session store before the PostgreSQL adapter.
- A future ADR chooses direct Nakama or Pitaya public API compatibility.

## Follow-Up

- Implement the PostgreSQL session adapter in a bounded slice.
- Define runtime session validation behavior before setting `RequestIdentity.SessionValidated` true.
- Define logout/revocation active-connection behavior before logout can close or invalidate bound WebSocket connections.
- Define reconnect and connection epoch behavior before duplicate connection replacement or resume behavior.
