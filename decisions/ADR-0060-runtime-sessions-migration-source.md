# ADR-0060: Runtime Sessions Migration Source

Status: Accepted
Date: 2026-05-17
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-postgres-session-persistence-schema-gate/`
- `changes/2026-05-17-implement-runtime-sessions-migration-source/`

Related conversations:

- `conversations/2026-05-17-postgres-session-persistence-migration-source.md`

Related artifacts:

- `runtime/migrations/postgres/000005_create_runtime_sessions.sql`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`W-0130` defined the PostgreSQL session persistence schema gate. The gate selected PostgreSQL as the first durable session target, `runtime_sessions` as the future logical table candidate, and `runtime/migrations/postgres/000005_create_runtime_sessions.sql` as the future migration source candidate.

The work queue reached `M-059/W-0131`, a confirmation gate after the schema gate. The maintainer asked the agent to recommend the next ten steps and then continue, while continuing to use Nakama and Pitaya as important references.

## Decision

Select:

```text
implement_runtime_sessions_migration_source
```

Add only the PostgreSQL migration source:

```text
runtime/migrations/postgres/000005_create_runtime_sessions.sql
```

The migration creates `runtime_sessions` with actor/player identity, session status, issued/expiration/last-seen timestamps, optional revocation fields, and optional `authentication_access_tokens(token_record_id)` linkage.

This ADR does not add session repository interfaces, PostgreSQL session adapters, runtime session creation, runtime session validation, runtime revocation execution, cleanup jobs, WebSocket handshake authentication, transport credential carriers, route-policy use of sessions or bound identity, Protobuf envelope changes, generated output, logout-triggered active WebSocket invalidation, reconnect/epoch behavior, dependencies, memory durable session behavior, or direct Nakama/Pitaya public API compatibility.

## Alternatives Considered

- Define the session repository boundary before adding the migration source.
- Implement the PostgreSQL session adapter immediately.
- Create a durable `runtime_session_connections` table together with `runtime_sessions`.
- Add runtime session creation to login or BindConnection in the same slice.
- Make first-message bound identity satisfy protected route policy after adding the table.
- Start logout/revocation active-connection behavior before repository and validation boundaries.
- Adopt Nakama session APIs or Pitaya session APIs directly.

## Rationale

Nakama demonstrates that authenticated sessions deserve a first-class lifecycle record. Pitaya reinforces that session-like context should remain separate from transport acceptors and handler plumbing. The right vibit adaptation is a small, inspectable PostgreSQL lifecycle table before repository or runtime behavior.

Keeping this slice migration-only makes the boundary clear for agents. A later repository boundary can define storage-neutral application behavior against a stable table shape. A later runtime validation gate can decide when `RequestIdentity.SessionValidated` may become true. A later connection lifecycle gate can decide whether logout closes active sockets or whether reconnect/epoch state is durable.

## Agent Reasoning Summary

After the schema gate, the smallest useful next step is the migration source itself. It turns the planned `runtime_sessions` shape into a concrete SQL artifact while avoiding the more dangerous behavior changes: route policy, WebSocket handshake authentication, session creation, logout disconnects, and reconnect.

## Decision Weights

```yaml
decision_weights:
  durable_lifecycle_foundation: high
  nakama_pitaya_alignment: high
  migration_safety: high
  agent_readability: high
  behavior_change_risk: low
  dependency_expansion: low
confidence: high
```

## Consequences

- `runtime/migrations/postgres/000005_create_runtime_sessions.sql` exists.
- `runtime.runtime_sessions_migration_source` becomes the repository check rule for this slice.
- `M-059/W-0131` is closed with `implement_runtime_sessions_migration_source`.
- `M-060/W-0132` is completed as a migration-source-only milestone.
- The work queue blocks again at `M-061/W-0133` before repositories, adapters, runtime validation, route policy, logout/revocation, reconnect, operations, or broader game backend work.
- Existing runtime behavior is not changed by this ADR.

## Reversal Conditions

Revisit this decision if:

- A future ADR replaces PostgreSQL as the first durable session target.
- The first runtime session implementation needs a different minimum lifecycle vocabulary before any repository code exists.
- A future ADR chooses a direct Nakama or Pitaya public API compatibility target.
- Connection lifecycle storage must be introduced before session lifecycle storage for a documented production reason.

## Follow-Up

- Define the storage-neutral session repository boundary.
- Define the PostgreSQL session adapter gate after the repository boundary.
- Define runtime session creation and validation behavior before setting `RequestIdentity.SessionValidated` true.
- Define logout/revocation active-connection behavior before logout can close or invalidate bound WebSocket connections.
- Define reconnect and connection epoch behavior before duplicate connection replacement or resume behavior.
