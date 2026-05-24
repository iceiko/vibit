# ADR-0140: Friends Relationship Persistence Schema Gate

Status: Accepted
Date: 2026-05-24
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-24-define-friends-relationship-persistence-schema-gate/`

Related conversations:

- `conversations/2026-05-24-friends-relationship-persistence-schema-gate.md`

Related artifacts:

- `docs/friends-relationship-persistence-schema-gate.md`
- `docs/friends-relationship-persistence-schema-gate.zh-CN.md`
- `docs/friends-relationship-lifecycle-gate.md`
- `docs/postgresql-persistence-boundary.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0139` completed the friends relationship lifecycle gate. It defined future request, accept, reject, remove, block, unblock, list, and status-read semantics for a Nakama-first social graph capability without adding runtime behavior, protocol routes, generated output, migrations, repositories, adapters, dependencies, or direct compatibility scope.

The next step is data-first: future relationship behavior needs a stable durable state posture before SQL, repositories, adapters, protocol, or runtime handlers exist.

## Decision

Define the friends relationship persistence schema gate in:

```text
docs/friends-relationship-persistence-schema-gate.md
docs/friends-relationship-persistence-schema-gate.zh-CN.md
```

The first persistence target is PostgreSQL. The future migration source candidate is:

```text
runtime/migrations/postgres/000007_create_friend_relationships.sql
```

The future logical current-state table is:

```text
friend_relationships
```

The gate records a pair-oriented schema posture with canonical unordered player pair identity, lifecycle state, actor-specific block timestamps, relationship versioning, timestamp checks, uniqueness/index expectations, event/audit deferral, redaction posture, and repository/adapter ownership candidates.

The repository check rule is:

```text
runtime.friends_relationship_persistence_schema_gate
```

This ADR is a schema gate only. It does not add the SQL migration source, create tables, add repositories or adapters, implement runtime behavior, expose protocol routes, add Protobuf source or generated output, add dependencies, wire startup, add event/audit tables, or create direct compatibility with Nakama or Pitaya.

Open the next bounded work item:

```text
M-161/W-0233 Add friends relationship migration source
```

## Alternatives Considered

- Add the SQL migration source in the same change.
- Add repository interfaces before a schema gate.
- Add runtime friendship behavior before persistence is explicit.
- Store actor-relative public states directly in the table.
- Add a relationship event/audit table in the first migration.
- Use a graph database or distributed social graph subsystem.
- Copy an external Nakama or Pitaya API or schema compatibility surface.

## Rationale

Pair-oriented persistence matches the lifecycle gate: relationship state is a fact about an unordered player pair, while public status is actor-relative and should be computed by future query behavior.

Using a single current-state table keeps the next migration-source slice narrow. Event/audit storage is valuable, but it should be separately authorized after current state, repository boundaries, and runtime event consistency expectations are clearer.

Block timestamps must be actor-specific because one player can block the other independently. Keeping block state separate from lifecycle state allows future behavior to derive `blocked_by_actor`, `blocked_actor`, and `mutual_blocked` without storing public actor-relative states as canonical database facts.

## Agent Reasoning Summary

The agent continued from `W-0231`, kept Nakama as the primary product capability reference, preserved Pitaya as deferred distributed architecture context, and defined a persistence schema gate only. It selected PostgreSQL, a future `friend_relationships` table, canonical pair identity, lifecycle state vocabulary, block timestamp representation, uniqueness/index posture, redaction, and future ownership boundaries while preserving all implementation deferrals.

## Decision Weights

```yaml
decision_weights:
  nakama_product_alignment: high
  ai_native_requirement_test_workflow_alignment: high
  contract_first_safety: high
  durable_social_graph_readiness: high
  privacy_and_redaction_control: high
  future_testability: high
  migration_scope_change: none
  runtime_scope_change: none
  protocol_scope_change: none
  direct_api_compatibility: none
confidence: high
```

## Consequences

- `M-160/W-0232` is completed.
- `runtime.friends_relationship_persistence_schema_gate` is registered.
- The friends relationship persistence schema standard and Simplified Chinese translation exist.
- The future migration source candidate is `runtime/migrations/postgres/000007_create_friend_relationships.sql`.
- The future current-state table candidate is `friend_relationships`.
- `M-161/W-0233 Add friends relationship migration source` becomes next-ready.
- Migration source creation, repository interfaces, PostgreSQL adapters, runtime behavior, protocol routes, Protobuf source, generated output, startup wiring, dependencies, event/audit tables, chat, groups, parties, matchmaking, match runtime, operations/admin behavior, SDKs, hosted surfaces, distributed runtime, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- the migration-source slice shows the pair-oriented current-state table cannot preserve lifecycle invariants;
- privacy requirements require event/audit retention before current-state rows;
- player id ordering is unsuitable for canonical pair identity;
- future runtime behavior needs a different block representation;
- a later ADR authorizes a different social module ownership model or external compatibility surface.

## Follow-Up

- Complete `W-0233`: add the friends relationship migration source.
- Keep repository interfaces, PostgreSQL adapters, runtime behavior, protocol routes, Protobuf source, generated output, startup wiring, dependencies, chat, groups, parties, matchmaking, match runtime, operations/admin behavior, SDKs, hosted surfaces, distributed runtime, event/audit tables, and direct compatibility behind later bounded work items.
