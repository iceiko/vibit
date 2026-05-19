# ADR-0059: PostgreSQL Session Persistence Schema Gate

Status: Accepted
Date: 2026-05-17
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-first-message-connection-binding-implementation/`
- `changes/2026-05-17-define-postgres-session-persistence-schema-gate/`

Related conversations:

- `conversations/2026-05-17-postgres-session-persistence-schema-gate.md`

Related artifacts:

- `docs/postgres-session-persistence-schema-gate.md`
- `docs/postgres-session-persistence-schema-gate.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`W-0128` implemented first-message `runtime.authentication.BindConnection` as a bounded system route. The implementation validates opaque access-token proof through the existing application validator and uses server-observed connection metadata, but it intentionally does not create durable sessions, connection registries, route-policy bound identity, logout/revocation active-connection behavior, or reconnect/epoch policy.

The work queue reached `M-057/W-0129`, a confirmation gate after first-message connection binding. The maintainer asked the agent to continue according to the agent's recommendation. The agent recommended defining a PostgreSQL session persistence schema gate next.

## Decision

Select:

```text
define_postgres_session_persistence_schema_gate
```

Define a gate-only standard for the future first durable runtime session schema.

The future first session persistence target is PostgreSQL, with `runtime_sessions` as the logical table candidate and `runtime/migrations/postgres/000005_create_runtime_sessions.sql` as the future migration source candidate.

This ADR does not add SQL migration source, create tables, add repository interfaces, add PostgreSQL adapters, change runtime session validation behavior, change route policy, change the Protobuf envelope, add WebSocket handshake authentication, parse transport credential carriers, add logout/revocation active-connection behavior, add reconnect/epoch behavior, add dependencies, or adopt direct Nakama/Pitaya public API compatibility.

## Alternatives Considered

- Implement the `runtime_sessions` migration immediately.
- Add a durable connection registry before a session table.
- Implement logout/revocation active-connection behavior first.
- Implement reconnect/epoch behavior first.
- Make `BindConnection` identity satisfy ordinary protected route policy before a session schema exists.
- Use a Redis-like session store as the first durable session target.
- Adopt Nakama session APIs or Pitaya session APIs directly.

## Rationale

Nakama demonstrates that sessions are a core lifecycle concept for authenticated gameplay and realtime socket use. Pitaya demonstrates that session-like context should remain separate from transport acceptors and handlers. vibit should adapt these concepts by making the first durable session schema explicit before adding runtime behavior.

PostgreSQL is the right first durable target because it is already the accepted authoritative store and can keep player/account/token/session relationships inspectable and testable. A schema gate before migration prevents session persistence from silently absorbing connection registries, logout disconnect behavior, reconnect, route policy, or broader realtime features.

## Agent Reasoning Summary

After BindConnection, the highest leverage next boundary is not more realtime behavior. It is the durable session schema that later work can build on. Without that schema gate, future logout, reconnect, bound route identity, presence, rooms, and match runtime work would lack a clear lifecycle foundation.

## Decision Weights

```yaml
decision_weights:
  durable_lifecycle_foundation: high
  nakama_pitaya_alignment: high
  agent_readability: high
  migration_safety: high
  route_policy_safety: high
  immediate_feature_delivery: low
  dependency_expansion: low
confidence: high
```

## Consequences

- `docs/postgres-session-persistence-schema-gate.md` becomes the standard for the future first session schema implementation queue.
- `runtime.postgres_session_persistence_schema_gate` becomes the repository check rule.
- `M-057/W-0129` is closed with the selected direction.
- `M-058/W-0130` is completed as a gate-only milestone.
- The work queue blocks again before adding a session migration source or choosing another major direction.
- Existing runtime behavior is not changed by this ADR.

## Reversal Conditions

Revisit this decision if:

- A future operations requirement makes Redis-like session storage the first durable target.
- PostgreSQL session persistence becomes unsuitable for realtime lifecycle pressure before a first schema exists.
- A future ADR adopts direct Nakama or Pitaya public API compatibility.
- Handshake-level authentication is selected before a durable session schema.
- Route-policy bound identity is selected before durable session semantics are needed.

## Follow-Up

- Add the `runtime_sessions` SQL migration only through a later bounded work item.
- Define a session repository boundary before adding PostgreSQL adapter behavior.
- Define session validation runtime behavior before setting `RequestIdentity.SessionValidated` true.
- Define logout/revocation active-connection behavior before logout can invalidate bound WebSocket connections.
- Define reconnect and connection epoch behavior before duplicate connection replacement or resume behavior.
