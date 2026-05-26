# ADR-0143: Friends Relationship Repository Interface Implementation

Status: Accepted
Date: 2026-05-25
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-25-implement-friends-relationship-repository-interface/`

Related conversations:

- `conversations/2026-05-26-friends-relationship-repository-interface-implementation.md`

Related artifacts:

- `runtime/internal/modules/friends/repository.go`
- `runtime/internal/modules/friends/repository_test.go`
- `modules/friends/module.yaml`
- `modules/friends/AGENTS.md`
- `modules/friends/AGENTS.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-162/W-0234` defined the friends relationship repository boundary after the PostgreSQL `friend_relationships` migration source. The next bounded step was to turn that boundary into storage-neutral Go repository vocabulary without adding adapter behavior, runtime friendship behavior, protocol routes, Protobuf source, generated output, dependencies, migration changes, event/audit tables, or direct external compatibility.

Nakama remains the product capability reference for friends relationships as a core game/backend social graph primitive. Pitaya remains deferred as a future distributed architecture reference. vibit adapts the capability through agent-native workflow: requirement, spec, tests, implementation, verification, and durable memory.

## Decision

Implement the storage-neutral repository interface under:

```text
runtime/internal/modules/friends
```

The package defines:

- `runtime/internal/modules/friends.Repository`;
- `FriendRelationship`, `FriendRelationshipPair`, `FriendRelationshipActor`, `FriendRelationshipBlockState`, and `FriendRelationshipVersion`;
- closed first-posture vocabulary for lifecycle states `pending`, `friends`, `rejected`, and `removed`;
- closed public status vocabulary for `pending`, `friends`, `blocked`, and `ended`;
- request, accept, reject, remove, block, unblock, pair lookup, and player-scoped list input/result types;
- conflict classes including `version_mismatch`;
- redacted repository error types;
- normalization helpers for records, list results, pair identity, actors, block state, and repository inputs;
- focused tests for storage neutrality, closed vocabularies, canonical pair handling, self rejection, actor validation, version validation and copying, list-result copying, redaction, and forbidden material.

Add the first friends module manifest and module AGENTS guides so future agents can discover ownership before adding adapters or runtime behavior.

This ADR does not add PostgreSQL friends adapters, SQL execution behavior, unit-of-work factory wiring, runtime friend handlers, protocol routes, Protobuf sources, generated output, dependencies, migration changes, authentication/session behavior changes, request identity validation, event/audit tables, chat, groups, parties, matchmaking, match runtime, SDK publication, hosted deployments, release artifacts, public announcements, paid promotion, distributed runtime, Pitaya-style architecture, or direct Nakama/Pitaya API compatibility.

Open the next bounded work item:

```text
M-164/W-0236 Define friends relationship PostgreSQL adapter gate
```

## Alternatives Considered

- Implement the PostgreSQL adapter in the same slice.
- Add runtime friend request/list/status behavior with the repository interface.
- Add Protobuf messages or public routes immediately.
- Place the interface under `runtime/internal/app`.
- Treat actor ids supplied by clients as proof.
- Add event/audit tables before current-state adapter semantics.
- Copy external Nakama or Pitaya public API compatibility.

## Rationale

The repository boundary already selected the owner candidate and capability vocabulary. Implementing only the storage-neutral interface now reduces future adapter ambiguity while keeping SQL, behavior, protocol, and event/audit concerns behind later gates.

Putting the interface in `runtime/internal/modules/friends` makes friends relationships a first-class domain module without making it own player accounts, authentication, sessions, transport behavior, chat, groups, parties, matchmaking, match runtime, or distributed topology.

The result serves the user's AI-native product goal: a future friendship feature request can be mapped into explicit specs, tests, implementation, and checks instead of relying on hidden architectural memory.

## Agent Reasoning Summary

The safest continuation from `W-0234` was an interface-only code slice. It gives later PostgreSQL adapter work a stable typed contract, adds tests for redaction and storage neutrality, and preserves stop conditions that keep protocol/runtime behavior from leaking into this package.

## Decision Weights

```yaml
decision_weights:
  nakama_product_alignment: high
  ai_native_requirement_test_workflow_alignment: high
  boundary_clarity: high
  agent_readability: high
  implementation_risk: low
  adapter_risk: deferred
  protocol_risk: low
  dependency_expansion: low
confidence: high
```

## Consequences

- `runtime/internal/modules/friends/repository.go` exists.
- `runtime/internal/modules/friends/repository_test.go` exists.
- `modules/friends/module.yaml` and paired module guides exist.
- `runtime.friends_relationship_repository_interface_implementation` becomes the repository check rule for this slice.
- `M-163/W-0235` is completed.
- `M-164/W-0236 Define friends relationship PostgreSQL adapter gate` becomes the next-ready work item.
- Existing runtime behavior, protocol behavior, migrations, generated output, and dependencies are unchanged by this ADR.

## Reversal Conditions

Revisit this decision if:

- friends relationships stop being a module-owned social graph capability;
- the first adapter needs a different repository owner;
- the `friend_relationships` table shape changes in a way that invalidates the value vocabulary;
- a later ADR selects a materially different conflict or leakage model;
- event/audit table requirements become mandatory before adapter work;
- direct Nakama or Pitaya public API compatibility becomes an explicit project goal.

## Follow-Up

- Define the friends relationship PostgreSQL adapter gate.
- Implement the adapter only after the gate is accepted.
- Define runtime behavior, permissions, protocol routes, generated output, and local proof only after repository and adapter boundaries are accepted.
- Keep chat, groups, parties, matchmaking, match runtime, SDKs, hosted surfaces, distributed runtime, event/audit tables, and direct compatibility behind future gates.
