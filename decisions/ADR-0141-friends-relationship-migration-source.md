# ADR-0141: Friends Relationship Migration Source

Status: Accepted
Date: 2026-05-25
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-24-add-friends-relationship-migration-source/`

Related conversations:

- `conversations/2026-05-25-friends-relationship-migration-source.md`

Related artifacts:

- `runtime/migrations/postgres/000007_create_friend_relationships.sql`
- `docs/friends-relationship-persistence-schema-gate.md`
- `docs/friends-relationship-lifecycle-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0140` completed the friends relationship persistence schema gate. It selected PostgreSQL as the first friends relationship store, `friend_relationships` as the future logical current-state table, canonical unordered player pairs as the relationship identity, actor-specific block timestamps as the first block representation, and `runtime/migrations/postgres/000007_create_friend_relationships.sql` as the migration source candidate.

The work queue reached `M-161/W-0233 Add friends relationship migration source`. This slice is a migration-source-only step in the Nakama-first `friends_groups_and_parties` capability family. Pitaya remains deferred as a future distributed architecture reference.

## Decision

Add only the PostgreSQL migration source:

```text
runtime/migrations/postgres/000007_create_friend_relationships.sql
```

The migration creates one current-state table:

```text
friend_relationships
```

The table includes:

- `relationship_id` as the server-generated opaque record id;
- canonical pair columns `player_low_id` and `player_high_id`;
- foreign keys from both pair members to `player_accounts(player_id)`;
- `lifecycle_state` constrained to `pending`, `friends`, `rejected`, and `removed`;
- pair-member actor columns for request, response, and removal history;
- actor-specific block timestamps `blocked_by_low_at` and `blocked_by_high_at`;
- positive `relationship_version BIGINT` with default `1`;
- created, updated, and state-changed timestamps;
- rejected, removed, and block timestamp ordering checks;
- uniqueness for the canonical player pair;
- lookup indexes by pair member and lifecycle state plus an updated-at index.

This ADR does not add friends relationship repository interfaces, PostgreSQL adapter behavior, runtime friendship behavior, protocol routes, Protobuf source, generated output, dependencies, startup wiring, automatic migration application, authentication/session changes, delivery guarantees, stream subscriptions, chat rooms, groups, parties, broadcast fanout, matchmaking, match runtime, operations/admin behavior, SDK publication, generated client libraries, hosted deployments, release artifacts, public announcements, paid promotion, event/audit tables, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility.

Open the next bounded work item:

```text
M-162/W-0234 Define friends relationship repository boundary
```

## Alternatives Considered

- Add repository interfaces in the same slice.
- Add the PostgreSQL adapter in the same slice.
- Add runtime friend request, accept, reject, remove, block, unblock, list, or status behavior in the same slice.
- Add protocol routes or Protobuf source with the table.
- Add a friend relationship event/audit table in the first migration.
- Store actor-relative public relationship states directly.
- Copy external Nakama or Pitaya API or schema compatibility.
- Introduce distributed social graph routing before the local persistence shape is proven.

## Rationale

Nakama demonstrates that friends and social graph state are a core backend capability. vibit adapts that product need into a source-first, agent-native sequence: lifecycle gate, schema gate, migration source, repository boundary, adapter, runtime behavior, protocol, and proof.

The migration keeps durable state pair-oriented because a relationship is a fact about an unordered player pair. Public status remains actor-relative and must be derived by future query behavior. The block timestamps are actor-specific because either player can block the other independently.

Keeping the slice migration-only reduces risk. It gives future agents a stable SQL source and static checks without introducing business behavior, adapter mapping, protocol compatibility claims, or distributed architecture.

## Agent Reasoning Summary

The agent continued from `W-0232`, kept Nakama as the primary capability reference, and added the narrow SQL source described by the accepted persistence schema gate. It preserved the AI-native workflow by recording acceptance criteria, verification, decision memory, and check coverage before moving to repository behavior.

## Decision Weights

```yaml
decision_weights:
  nakama_product_alignment: high
  ai_native_requirement_test_workflow_alignment: high
  migration_safety: high
  agent_readability: high
  future_repository_testability: high
  runtime_behavior_risk: low
  dependency_expansion: low
  direct_api_compatibility: none
confidence: high
```

## Consequences

- `runtime/migrations/postgres/000007_create_friend_relationships.sql` exists.
- `runtime.friends_relationship_migration_source` becomes the repository check rule for this slice.
- `M-161/W-0233` is completed as a migration-source-only milestone.
- The work queue advances to `M-162/W-0234 Define friends relationship repository boundary`.
- Existing runtime behavior is not changed by this ADR.
- Friends relationships are not yet exposed through repository interfaces, PostgreSQL adapters, runtime services, protocol routes, generated output, SDKs, hosted surfaces, or direct external compatibility.

## Reversal Conditions

Revisit this decision if:

- a later repository or runtime gate proves the canonical pair table cannot preserve lifecycle invariants;
- a privacy or compliance requirement demands event/audit storage before repository behavior;
- player id ordering is unsuitable for canonical pair identity;
- future command behavior requires a different first block representation;
- a later ADR explicitly authorizes a different social module storage model or external compatibility target.

## Follow-Up

- Define the friends relationship repository boundary.
- Keep repository implementation, PostgreSQL adapter behavior, runtime friendship behavior, protocol routes, Protobuf source, generated output, startup wiring, dependencies, event/audit tables, groups, parties, chat, matchmaking, match runtime, SDKs, hosted surfaces, distributed runtime, and direct compatibility behind later bounded work items.
