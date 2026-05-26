# ADR-0144: Friends Relationship PostgreSQL Adapter Gate

Status: Accepted
Date: 2026-05-26
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-friends-relationship-postgresql-adapter-gate/`

Related conversations:

- `conversations/2026-05-26-friends-relationship-postgresql-adapter-gate.md`

Related artifacts:

- `docs/friends-relationship-postgresql-adapter-gate.md`
- `docs/friends-relationship-postgresql-adapter-gate.zh-CN.md`
- `runtime/internal/modules/friends/repository.go`
- `runtime/migrations/postgres/000007_create_friend_relationships.sql`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/friends/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-163/W-0235` implemented the storage-neutral friends relationship repository interface under `runtime/internal/modules/friends`. The next bounded step is to define the PostgreSQL adapter gate before any adapter code, SQL execution, runtime behavior, or protocol surface is added.

The existing `friend_relationships` migration source already defines the current-state table, and the repository interface already defines friend request, read, list, lifecycle transition, and block mutation vocabulary. A separate adapter gate keeps SQL mapping, transaction handoff, error mapping, redaction, and tests explicit before implementation.

Nakama remains the product capability reference for friends relationships as a core social graph primitive. Pitaya remains deferred as a future distributed architecture reference. vibit adapts the capability through explicit adapter boundaries and checkable manifests, not direct public API compatibility.

## Decision

Accept `docs/friends-relationship-postgresql-adapter-gate.md` as the gate for the future friends relationship PostgreSQL adapter.

The gate records:

- future adapter owner `runtime/internal/platform/persistence/postgres`;
- future source candidate `runtime/internal/platform/persistence/postgres/friend_relationship_repository.go`;
- future test candidate `runtime/internal/platform/persistence/postgres/friend_relationship_repository_test.go`;
- constructor candidate `NewFriendRelationshipRepositoryForUnitOfWork`;
- repository interface source `runtime/internal/modules/friends.Repository`;
- SQL mapping posture for `friend_relationships`;
- constructor and caller-supplied executor expectations;
- unit-of-work and transaction handoff expectations;
- redacted driver-error, affected-row, and conflict mapping;
- focused adapter implementation test expectations;
- stop conditions before implementation, runtime behavior, protocol routes, generated output, dependencies, migration changes, event/audit tables, hosted deployment, release artifacts, public announcements, paid promotion, broader social features, distributed runtime, or direct compatibility.

This ADR does not add PostgreSQL friends adapters, SQL execution behavior, unit-of-work factory wiring, runtime friend handlers, protocol routes, Protobuf sources, generated output, dependencies, migration changes, authentication/session behavior changes, request identity validation, event/audit tables, chat, groups, parties, matchmaking, match runtime, SDK publication, hosted deployments, release artifacts, public announcements, paid promotion, distributed runtime, Pitaya-style architecture, or direct Nakama/Pitaya API compatibility.

Open the next bounded work item:

```text
M-165/W-0237 Implement friends relationship PostgreSQL adapter
```

## Alternatives Considered

- Implement the PostgreSQL adapter immediately after the repository interface.
- Reuse the storage objects adapter implementation shape without a friends-specific gate.
- Put SQL execution under `runtime/internal/modules/friends`.
- Add runtime friend request/list/status behavior together with adapter implementation.
- Add Protobuf messages or public routes immediately.
- Add event/audit tables before current-state adapter semantics.
- Copy external Nakama or Pitaya public API compatibility.

## Rationale

The friends relationship repository has conflict and redaction pressure that should be explicit before implementation: canonical pair uniqueness, lifecycle transitions, pending/friends/blocked conflicts, actor-specific block columns, optimistic version checks, affected-row interpretation, foreign-key outcomes, and private social graph leakage.

A gate-only ADR keeps the next implementation slice bounded and makes repository checks able to reject accidental SQL, protocol, runtime, generated output, or broad social feature behavior before that work item is authorized.

## Agent Reasoning Summary

The safest continuation from `W-0235` is a platform adapter gate. It gives future implementation a precise owner, constructor posture, SQL mapping checklist, conflict mapping posture, and test list while preserving the separation between friends module vocabulary, PostgreSQL adapter behavior, runtime routing, protocol shape, and product-scope expansion.

## Decision Weights

```yaml
decision_weights:
  nakama_product_alignment: high
  ai_native_requirement_test_workflow_alignment: high
  boundary_clarity: high
  agent_readability: high
  implementation_risk: low
  adapter_risk: contained_by_next_work_item
  protocol_risk: deferred
  dependency_expansion: low
confidence: high
```

## Consequences

- `docs/friends-relationship-postgresql-adapter-gate.md` and `.zh-CN.md` exist.
- `runtime.friends_relationship_postgresql_adapter_gate` becomes the repository check rule for this slice.
- `M-164/W-0236` is completed.
- `M-165/W-0237 Implement friends relationship PostgreSQL adapter` becomes the next-ready work item.
- Existing runtime behavior, protocol behavior, migrations, dependencies, generated output, and SQL execution behavior remain unchanged by this ADR.

## Reversal Conditions

Revisit this decision if:

- `runtime/internal/modules/friends.Repository` changes before adapter implementation;
- the `friend_relationships` migration source changes materially;
- the project selects a different first friends relationship persistence engine;
- transaction ownership moves away from caller-supplied unit-of-work boundaries;
- event/audit table requirements become mandatory before adapter work;
- direct Nakama or Pitaya public API compatibility becomes an explicit project goal.

## Follow-Up

- Implement the friends relationship PostgreSQL adapter only after this gate is accepted.
- Keep runtime behavior, permissions, protocol routes, generated output, and local proof behind later gates.
- Keep chat, groups, parties, matchmaking, match runtime, SDKs, hosted surfaces, distributed runtime, event/audit tables, and direct compatibility behind future gates.

