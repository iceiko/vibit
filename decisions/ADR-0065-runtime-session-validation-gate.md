# ADR-0065: Runtime Session Validation Gate

Status: Accepted
Date: 2026-05-17
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-session-postgresql-adapter-implementation/`
- `changes/2026-05-17-define-runtime-session-validation-gate/`

Related conversations:

- `conversations/2026-05-17-runtime-session-validation-gate.md`

Related artifacts:

- `docs/runtime-session-validation-gate.md`
- `docs/runtime-session-validation-gate.zh-CN.md`
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

`W-0140` implemented the PostgreSQL adapter for the storage-neutral runtime session repository. The repository now has a durable `runtime_sessions` table, a session repository interface, and a PostgreSQL adapter with unit-of-work wiring.

The work queue reached `M-069/W-0141`, a confirmation gate. The maintainer asked the agent to recommend the next ten steps and continue, while continuing to learn from Nakama and Pitaya.

## Decision

Select:

```text
define_runtime_session_validation_gate
```

Define a gate-only runtime session validation standard in:

```text
docs/runtime-session-validation-gate.md
docs/runtime-session-validation-gate.zh-CN.md
```

The future validation owner is:

```text
runtime/internal/app
```

The session record and repository owner remains:

```text
runtime/internal/app/session
```

The PostgreSQL adapter owner remains:

```text
runtime/internal/platform/persistence/postgres
```

This ADR does not add runtime session validation code, set `RequestIdentity.SessionValidated` true, create sessions at login or `BindConnection`, change route policy, change WebSocket handshake authentication, add transport credential carriers, add Protobuf session messages, change the existing Protobuf envelope, add generated output, add logout/revocation active-connection invalidation, add reconnect/epoch behavior, add cleanup jobs, add dependencies, add memory durable session behavior, or adopt direct Nakama/Pitaya public API compatibility.

## Alternatives Considered

- Implement runtime session validation immediately.
- Attach durable session creation to login or `BindConnection` before validation semantics are defined.
- Make first-message bound identity or persisted session identity satisfy ordinary protected routes immediately.
- Put session validation inside the PostgreSQL adapter.
- Put session validation inside WebSocket transport or Protobuf protocol adapters.
- Define logout/revocation active-connection invalidation before validation behavior.
- Copy Nakama or Pitaya session APIs directly.

## Rationale

Nakama shows that session validity is a real lifecycle concern: session records can expire, be revoked, and interact with logout and realtime socket behavior. Nakama also shows that token/session logout and active socket disconnect behavior should not be collapsed accidentally.

Pitaya shows that acceptors, session context, and handlers should remain separated. Handler-facing context can carry session state, but transport acceptors should not own durable validation or request identity construction.

The right vibit adaptation is to define an application-owned runtime session validation gate before adding code. The gate records the future validation order, identity handoff, error collapse, redaction expectations, tests, and deferrals so later agents cannot mistake a durable session row for authenticated proof.

## Agent Reasoning Summary

After the PostgreSQL adapter, the highest-leverage next step is to decide how persisted session records can become validated request identity. Implementing validation immediately would cross several unresolved boundaries: session creation composition, route-policy use of session identity, logout/revocation socket effects, reconnect/epoch behavior, and protocol carriers.

The gate gives the project a precise next implementation target while preserving vibit's separation between authentication proof validation, session lifecycle storage, request identity construction, transport connection plumbing, and protocol routing.

## Decision Weights

```yaml
decision_weights:
  validated_identity_safety: high
  nakama_pitaya_alignment: high
  separation_of_transport_protocol_app_persistence: high
  future_implementation_clarity: high
  immediate_feature_delivery: medium
  direct_nakama_pitaya_api_compatibility: low
confidence: high
```

## Consequences

- `docs/runtime-session-validation-gate.md` and `docs/runtime-session-validation-gate.zh-CN.md` exist.
- `runtime.runtime_session_validation_gate` becomes the repository check rule for this slice.
- `M-069/W-0141` is closed with `define_runtime_session_validation_gate`.
- `M-070/W-0142` is completed as a gate-only milestone.
- The work queue blocks again at the next confirmation gate before implementation, session creation composition, route-policy use of session or bound identity, logout/revocation active-connection behavior, reconnect/epoch behavior, operations, memory durable session behavior, or broader game backend expansion.
- Existing runtime behavior is not changed by this ADR.

## Reversal Conditions

Revisit this decision if:

- A future ADR chooses handshake-level authentication as the primary session proof path.
- A future ADR replaces PostgreSQL-backed runtime sessions with a different first durable session store.
- A future session creation composition decision requires different validation inputs.
- Route-policy requirements force a different identity handoff model.
- A later ADR adopts direct Nakama or Pitaya public API compatibility.

## Follow-Up

- Define or implement runtime session validation in a bounded slice before setting `RequestIdentity.SessionValidated` true.
- Define session creation composition before login or `BindConnection` creates durable runtime session rows.
- Define bound-identity route policy before protected routes can use session-validated or connection-bound identity.
- Define logout/revocation active-connection behavior before revocation closes or invalidates WebSocket connections.
- Define reconnect and connection epoch behavior before duplicate replacement or resume behavior.
