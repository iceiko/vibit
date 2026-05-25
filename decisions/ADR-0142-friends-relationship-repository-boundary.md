# ADR-0142: Friends Relationship Repository Boundary

Status: Accepted
Date: 2026-05-25
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-25-define-friends-relationship-repository-boundary/`

Related conversations:

- `conversations/2026-05-25-friends-relationship-repository-boundary.md`

Related artifacts:

- `docs/friends-relationship-repository-boundary.md`
- `docs/friends-relationship-repository-boundary.zh-CN.md`
- `runtime/migrations/postgres/000007_create_friend_relationships.sql`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0141` added the PostgreSQL `friend_relationships` migration source as the first durable current-state table for a future Nakama-class friends relationship capability. The project has intentionally sequenced this path as lifecycle gate, persistence schema gate, migration source, repository boundary, repository interface, adapter gate, adapter implementation, runtime behavior, protocol route, and local proof.

The work queue reached `M-162/W-0234 Define friends relationship repository boundary`. This slice is a boundary-only step in the Nakama-first `friends_groups_and_parties` capability family. Pitaya remains deferred as a future distributed architecture reference.

## Decision

Accept `docs/friends-relationship-repository-boundary.md` as the storage-neutral repository boundary for future friends relationship lifecycle behavior.

The boundary records:

- future repository owner candidate `runtime/internal/modules/friends`;
- future interface candidate `runtime/internal/modules/friends.Repository`;
- future PostgreSQL adapter owner `runtime/internal/platform/persistence/postgres`;
- source migration `runtime/migrations/postgres/000007_create_friend_relationships.sql`;
- logical table `friend_relationships`;
- candidate value types for relationships, pairs, actors, lifecycle state, versions, block state, inputs, conflicts, and repository errors;
- candidate capabilities for request creation/update, pair lookup, player listing, accept, reject, remove, block, and unblock;
- canonical unordered player pair identity;
- request identity handoff from already validated application identity;
- optimistic version and typed conflict posture;
- private social graph redaction posture;
- PostgreSQL adapter expectations tied to the indexes created by `ADR-0141`;
- future implementation queue boundaries.

This ADR does not add repository interface implementation, PostgreSQL adapter behavior, runtime friendship behavior, protocol routes, Protobuf source, generated output, dependencies, migrations, startup wiring, automatic migration application, authentication/session changes, delivery guarantees, stream subscriptions, chat rooms, groups, parties, broadcast fanout, matchmaking, match runtime, operations/admin behavior, SDK publication, generated client libraries, hosted deployments, release artifacts, public announcements, paid promotion, event/audit tables, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility.

Open the next bounded work item:

```text
M-163/W-0235 Implement storage-neutral friends relationship repository interface
```

## Alternatives Considered

- Add the Go repository interface in the same slice.
- Add the PostgreSQL adapter in the same slice.
- Add runtime friend request, accept, reject, remove, block, unblock, list, or status behavior in the same slice.
- Add protocol routes or Protobuf source with the repository boundary.
- Add a friend relationship event/audit table before repository vocabulary.
- Treat actor ids supplied by clients as proof.
- Copy external Nakama or Pitaya API compatibility.
- Introduce distributed social graph routing before single-process repository semantics are proven.

## Rationale

Nakama demonstrates that friends and social graph state are a core game/backend capability. vibit adapts that product need into an agent-native sequence where each layer becomes explicit and checkable before the next layer depends on it.

The repository boundary is needed before interface implementation because future agents need a concise vocabulary for ownership, value types, conflicts, redaction, and adapter expectations. Defining that vocabulary now reduces the risk that W-0235 invents ad hoc SQL-shaped types or mixes identity proof with social graph state.

Keeping the slice boundary-only reduces risk. It lets static checks prove that no runtime behavior, adapter behavior, protocol shape, generated output, or direct external compatibility was added.

## Agent Reasoning Summary

The agent continued from `W-0233`, kept Nakama as the primary capability reference, and converted the migration source into a storage-neutral future repository plan. The change preserves the user's AI-native product purpose: a future user requirement should become a bounded spec, acceptance criteria, test plan, tests, implementation, verification, and durable memory.

## Decision Weights

```yaml
decision_weights:
  nakama_product_alignment: high
  ai_native_requirement_test_workflow_alignment: high
  repository_testability: high
  agent_readability: high
  privacy_and_redaction: high
  runtime_behavior_risk: low
  dependency_expansion: low
  direct_api_compatibility: none
confidence: high
```

## Consequences

- `docs/friends-relationship-repository-boundary.md` and its Simplified Chinese translation exist.
- `runtime.friends_relationship_repository_boundary` becomes the repository check rule for this slice.
- `M-162/W-0234` is completed as a repository-boundary-only milestone.
- The work queue advances to `M-163/W-0235 Implement storage-neutral friends relationship repository interface`.
- Existing runtime behavior is not changed by this ADR.
- Friends relationships are not yet exposed through repository interfaces, PostgreSQL adapters, runtime services, protocol routes, generated output, SDKs, hosted surfaces, or direct external compatibility.

## Reversal Conditions

Revisit this decision if:

- future repository interface implementation proves the candidate vocabulary cannot express the canonical pair lifecycle safely;
- privacy or compliance requirements demand a different redaction posture before interface implementation;
- a future runtime behavior gate selects a materially different conflict or leakage model;
- a future ADR explicitly authorizes a different social module storage model or external compatibility target.

## Follow-Up

- Implement the storage-neutral friends relationship repository interface as `W-0235`.
- Keep PostgreSQL adapter behavior, runtime friendship behavior, protocol routes, Protobuf source, generated output, startup wiring, dependencies, event/audit tables, groups, parties, chat, matchmaking, match runtime, SDKs, hosted surfaces, distributed runtime, and direct compatibility behind later bounded work items.
