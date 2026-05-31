# ADR-0145: Friends Relationship PostgreSQL Adapter Implementation

Status: Accepted
Date: 2026-05-26
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-friends-relationship-postgresql-adapter/`

Related conversations:

- `conversations/2026-05-26-friends-relationship-postgresql-adapter-implementation.md`

Related artifacts:

- `runtime/internal/platform/persistence/postgres/friend_relationship_repository.go`
- `runtime/internal/platform/persistence/postgres/friend_relationship_repository_test.go`
- `runtime/internal/platform/persistence/postgres/runner.go`
- `modules/friends/module.yaml`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-164/W-0236` defined the friends relationship PostgreSQL adapter gate after the storage-neutral friends repository interface. The gate authorized a later bounded implementation under `runtime/internal/platform/persistence/postgres` and kept runtime friendship behavior, protocol routes, generated output, dependencies, migration changes, event/audit tables, broad social features, hosted surfaces, SDKs, distributed runtime, and direct Nakama/Pitaya API compatibility deferred.

The existing `friend_relationships` migration source defines the current-state relationship table. The friends module owns `runtime/internal/modules/friends.Repository`, canonical pair value types, lifecycle/status vocabularies, conflict classes, normalizers, and redacted repository errors.

Nakama keeps friends relationships as a core social graph primitive. Pitaya reinforces keeping persistence below handlers, routes, RPC, and cluster behavior. vibit adapts those lessons through a platform adapter implementing a storage-neutral repository interface, not through direct public API compatibility.

## Decision

Implement the friends relationship PostgreSQL adapter under:

```text
runtime/internal/platform/persistence/postgres
```

The implementation adds:

- `FriendRelationshipRepository`;
- `NewFriendRelationshipRepositoryForUnitOfWork`;
- `UnitOfWork.NewFriendRelationshipRepository`;
- create/update request mapping to the existing `friend_relationships` table;
- pair lookup and player-scoped list mapping;
- pending/friends/ended/blocked status filtering;
- one-row-overflow pagination using canonical pair tokens;
- accept, reject, and remove lifecycle transitions with optional expected-version checks;
- actor-specific block and unblock column updates;
- server-side `relationship_version` increments;
- row scanning through friends module normalizers;
- redacted mapping to friends module repository errors;
- focused fake-executor tests.

This ADR does not add runtime friendship handlers, application service behavior, protocol routes, Protobuf sources, generated output, dependency changes, migration changes, automatic startup migration behavior, authentication/session behavior changes, event/audit tables, chat, groups, parties, matchmaking, match runtime, operations/admin behavior, SDK publication, hosted deployments, release artifacts, public announcements, paid promotion, distributed runtime, Pitaya-style architecture, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Wait for a live PostgreSQL integration environment before implementing the adapter.
- Implement runtime friendship behavior together with persistence.
- Add protocol routes and Protobuf messages together with persistence.
- Put SQL execution under `runtime/internal/modules/friends`.
- Add event/audit tables before current-state adapter semantics.
- Add a new PostgreSQL dependency or direct Nakama-compatible API.

## Rationale

The repository interface, migration source, and adapter gate already exist. Implementing the adapter now gives future friends relationship runtime behavior a real persistence target while preserving later gates for validated request identity, actor-relative public status, permission, route policy, protocol, generated output, and broader social feature decisions.

Fake-executor tests follow the existing PostgreSQL adapter pattern in this repository. They verify query shape, argument normalization, row scanning, redaction, conflict handling, expected-version checks, block-column selection, and absence of transaction-control SQL without making default checks depend on a live database.

## Agent Reasoning Summary

The safest continuation from `W-0236` was a platform adapter slice only. It provides durable friends relationship persistence behind `friends.Repository`, keeps SQL out of the friends module, and leaves runtime and protocol behavior to later gates.

## Decision Weights

```yaml
decision_weights:
  nakama_product_alignment: high
  boundary_clarity: high
  agent_readability: high
  implementation_risk: medium
  protocol_risk: deferred
  dependency_expansion: low
confidence: high
```

## Consequences

- `runtime/internal/platform/persistence/postgres/friend_relationship_repository.go` exists.
- `runtime/internal/platform/persistence/postgres/friend_relationship_repository_test.go` exists.
- `UnitOfWork.NewFriendRelationshipRepository` exists as an adapter handoff helper.
- `runtime.friends_relationship_postgresql_adapter_implementation` becomes the repository check rule for this slice.
- `M-165/W-0237` is completed.
- `M-166/W-0238 Define friends relationship runtime behavior gate` becomes the next-ready work item.
- Runtime behavior, protocol routes, generated output, migrations, dependencies, event/audit tables, and authentication/session semantics remain unchanged by this ADR.

## Reversal Conditions

Revisit this decision if:

- the friends repository interface changes materially;
- the `friend_relationships` migration source changes materially;
- a later ADR selects a different pagination or conflict model;
- the project requires live PostgreSQL adapter tests as default verification;
- event/audit table requirements become mandatory before adapter use;
- direct Nakama or Pitaya public API compatibility becomes an explicit project goal.

## Follow-Up

- Define the friends relationship runtime behavior gate before adding application services or handlers.
- Define protocol routes and Protobuf messages only after runtime behavior is ratified.
- Keep chat, groups, parties, matchmaking, match runtime, SDKs, hosted surfaces, distributed runtime, event/audit tables, and direct compatibility behind future gates.
