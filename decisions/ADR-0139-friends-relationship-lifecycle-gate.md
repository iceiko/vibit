# ADR-0139: Friends Relationship Lifecycle Gate

Status: Accepted
Date: 2026-05-24
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-24-define-friends-relationship-lifecycle-gate/`

Related conversations:

- `conversations/2026-05-24-friends-relationship-lifecycle-gate.md`

Related artifacts:

- `docs/friends-relationship-lifecycle-gate.md`
- `docs/friends-relationship-lifecycle-gate.zh-CN.md`
- `docs/agent-native-feature-request-test-workflow.md`
- `docs/agent-native-feature-request-scaffolding.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0138` completed the scaffolded Nakama feature request intake pilot. It selected the `friends_groups_and_parties` capability family and opened `M-159/W-0231 Define friends relationship lifecycle gate`.

The maintainer's current product direction is Nakama-first. Pitaya remains deferred as a future distributed architecture reference. The product purpose is AI-native development and AI-native testing: a user requirement should become a bounded spec, acceptance criteria, test plan, implementation plan, verification, and durable project memory before behavior is added.

Friendship lifecycle is a core Nakama-class social graph capability, but vibit must adapt it into its own contract-first and server-authoritative model rather than copying external public APIs or compatibility surfaces.

## Decision

Define the friends relationship lifecycle gate in:

```text
docs/friends-relationship-lifecycle-gate.md
docs/friends-relationship-lifecycle-gate.zh-CN.md
```

The gate records the future semantic scope:

```text
request, accept, reject, remove, block, unblock, list, read relationship status
```

It records future command, query, event, error, permission, actor-relative state, invariant, redaction, and test-plan vocabulary without creating contract source files, runtime behavior, protocol routes, generated output, migrations, repository interfaces, adapters, dependencies, startup wiring, SDKs, hosted surfaces, distributed runtime, or direct Nakama/Pitaya API compatibility.

Every future command and query must require validated player identity. Metadata-only `player_id` or `session_id` is not proof.

The repository check rule is:

```text
runtime.friends_relationship_lifecycle_gate
```

Open the next bounded work item:

```text
M-160/W-0232 Define friends relationship persistence schema gate
```

## Alternatives Considered

- Implement friendship runtime behavior immediately.
- Add protocol routes and Protobuf messages immediately.
- Add a PostgreSQL migration immediately.
- Start with groups, parties, chat, matchmaking, or match runtime instead.
- Copy an external social API shape for compatibility.
- Reactivate Pitaya distributed architecture concerns for social graph work.

## Rationale

The selected step keeps the project aligned with both goals: Nakama-class capability coverage and agent-native maintainability. Friendship lifecycle is useful enough to move product capability forward, but it has enough privacy, identity, transition, persistence, and concurrency risk that a semantic gate should precede code.

Defining the lifecycle first gives future agents a stable target for acceptance criteria and tests. It also prevents the social graph from leaking into transport, protocol adapter, authentication, storage object, realtime delivery, chat, matchmaking, or match runtime boundaries.

The next step should be persistence schema gate, not runtime implementation, because relationship state is durable social graph state and future command behavior needs a stable table/state posture before migration, repository, and handler work.

## Agent Reasoning Summary

The agent continued from `W-0230` and treated `W-0231` as gate-only. It mapped the user's Nakama-first direction to friends/social graph capability planning, preserved Pitaya deferral, defined the lifecycle vocabulary and test expectations, and opened a schema gate as the next bounded continuation. No runtime, protocol, generated, migration, dependency, hosted, SDK, distributed, or direct compatibility scope was added.

## Decision Weights

```yaml
decision_weights:
  nakama_product_alignment: high
  ai_native_requirement_test_workflow_alignment: high
  contract_first_safety: high
  identity_and_privacy_risk_control: high
  future_testability: high
  persistence_readiness: medium
  runtime_scope_change: none
  protocol_scope_change: none
  migration_scope_change: none
  direct_api_compatibility: none
confidence: high
```

## Consequences

- `M-159/W-0231` is completed.
- `runtime.friends_relationship_lifecycle_gate` is registered.
- The friends relationship lifecycle semantic standard and Simplified Chinese translation exist.
- Future relationship commands, queries, events, errors, permissions, actor-relative states, invariants, redaction posture, and tests are defined as planning vocabulary.
- `M-160/W-0232 Define friends relationship persistence schema gate` becomes next-ready.
- Runtime behavior, protocol routes, Protobuf source, generated output, migrations, repository interfaces, PostgreSQL adapters, dependencies, startup wiring, authentication/session changes, groups, parties, chat, matchmaking, match runtime, operations/admin behavior, SDK publication, hosted surfaces, distributed runtime, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- the next persistence schema gate shows that player account identity or privacy requirements are not ready;
- alpha users need a different Nakama-class capability before social graph work continues;
- friendship lifecycle semantics prove too broad and need to split request/remove/block into separate standards;
- a later ADR authorizes a different social module ownership model;
- a later ADR explicitly adopts an external compatibility surface.

## Follow-Up

- Complete `W-0232`: define the friends relationship persistence schema gate.
- Keep runtime behavior, protocol routes, Protobuf source, generated output, migrations, repository interfaces, PostgreSQL adapters, dependencies, startup wiring, groups, parties, chat, matchmaking, match runtime, SDKs, hosted surfaces, distributed runtime, and direct compatibility behind later bounded work items.
