# ADR-0062: Session Repository Interface Implementation

Status: Accepted
Date: 2026-05-17
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-session-repository-boundary/`
- `changes/2026-05-17-implement-session-repository-interface/`

Related conversations:

- `conversations/2026-05-17-session-repository-interface-implementation.md`

Related artifacts:

- `runtime/internal/app/session/repository.go`
- `runtime/internal/app/session/repository_test.go`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`W-0134` defined the storage-neutral session repository boundary after the `runtime_sessions` migration source. The repository still had no Go interface for future session lifecycle behavior.

The work queue reached `M-063/W-0135`, a confirmation gate. The maintainer asked the agent to recommend the next ten steps and continue, while continuing to reference Nakama and Pitaya.

## Decision

Select:

```text
implement_session_repository_interface
```

Add a storage-neutral application package:

```text
runtime/internal/app/session
```

The package defines:

- `SessionStatus` closed vocabulary: `active`, `expired`, `revoked`.
- First actor kind vocabulary: `player`.
- `RuntimeSession` lifecycle record.
- Repository capability methods for creation, lookup, active lookup, last-seen update, expiration, revocation, and bounded player active-session listing.
- Normalization helpers for records, mutations, and queries.
- Focused unit tests.

This ADR does not add PostgreSQL adapters, SQL queries, unit-of-work factory methods, runtime session creation, runtime session validation, runtime session revocation execution, cleanup jobs, WebSocket handshake authentication, transport credential carriers, route-policy use of session or bound identity, Protobuf session messages, existing envelope changes, generated output, logout-triggered active WebSocket invalidation, reconnect/epoch behavior, dependencies, memory durable session behavior, or direct Nakama/Pitaya public API compatibility.

## Alternatives Considered

- Define a PostgreSQL session adapter gate first.
- Implement a PostgreSQL session adapter immediately.
- Add runtime session validation before a repository interface exists.
- Attach session creation to login or `BindConnection` in the same slice.
- Make persisted session identity satisfy ordinary protected route policy.
- Add logout/revocation active-connection invalidation.
- Adopt Nakama or Pitaya public APIs directly.

## Rationale

Nakama demonstrates that authenticated sessions become first-class lifecycle objects: lookup, expiration, revocation, logout, and management pressure appear quickly once realtime and social features grow. Pitaya demonstrates that handler-facing session context should remain separate from acceptors, transport, and routing.

The smallest useful vibit adaptation is therefore a storage-neutral application repository interface. It gives later agents a concrete typed boundary for session lifecycle data while preserving transport neutrality and keeping SQL, validation, route policy, and connection behavior behind later gates.

## Agent Reasoning Summary

The prior boundary already chose the owner candidate and capability vocabulary. Implementing the Go interface now reduces future ambiguity without changing runtime behavior. It also prevents later adapter or validation work from inventing ad hoc session shapes under pressure.

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

- `runtime/internal/app/session/repository.go` exists.
- `runtime/internal/app/session/repository_test.go` exists.
- `runtime.session_repository_interface_implementation` becomes the repository check rule for this slice.
- `M-063/W-0135` is closed with `implement_session_repository_interface`.
- `M-064/W-0136` is completed as an interface-only implementation milestone.
- The work queue blocks again at `M-065/W-0137` before adapters, runtime validation, route policy, logout/revocation, reconnect, operations, or broader game backend work.
- Existing runtime behavior is not changed by this ADR.

## Reversal Conditions

Revisit this decision if:

- A future ADR replaces `runtime_sessions` with a different first session persistence model.
- A future runtime session implementation needs an owner outside `runtime/internal/app/session`.
- Runtime validation requires a smaller or differently named capability set before adapter implementation.
- A future ADR chooses direct Nakama or Pitaya public API compatibility.

## Follow-Up

- Define the PostgreSQL session adapter gate.
- Define runtime session validation behavior before setting `RequestIdentity.SessionValidated` true.
- Define logout/revocation active-connection behavior before logout can close or invalidate bound WebSocket connections.
- Define reconnect and connection epoch behavior before duplicate connection replacement or resume behavior.
