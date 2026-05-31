# ADR-0146: Friends Relationship Runtime Behavior Gate

Status: Accepted
Date: 2026-05-26
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-friends-relationship-runtime-behavior-gate/`

Related conversations:

- `conversations/2026-05-26-friends-relationship-runtime-behavior-gate.md`

Related artifacts:

- `docs/friends-relationship-runtime-behavior-gate.md`
- `docs/friends-relationship-runtime-behavior-gate.zh-CN.md`
- `runtime/internal/modules/friends/repository.go`
- `runtime/internal/platform/persistence/postgres/friend_relationship_repository.go`
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

`M-165/W-0237` implemented the friends relationship PostgreSQL adapter under `runtime/internal/platform/persistence/postgres`. The repository now has a storage-neutral friends repository interface, a PostgreSQL `friend_relationships` migration source, and a platform adapter with unit-of-work repository handoff.

The next bounded step is to define how runtime behavior may later use those persistence pieces. That behavior must derive the actor from validated request identity, preserve metadata-only identity protections, compute public relationship state relative to the actor, keep permission and route policy application-owned, and avoid protocol or generated-output changes until a later protocol slice.

Nakama keeps friends relationships as a core social graph primitive. Pitaya reinforces separating route/session context, handlers, and persistence. vibit adapts those lessons by defining an application-owned runtime behavior gate, not by copying public APIs.

## Decision

Accept `docs/friends-relationship-runtime-behavior-gate.md` as the gate for future friends relationship runtime behavior.

The gate records:

- future application owner `runtime/internal/app`;
- future package candidate `runtime/internal/app/friends`;
- future service source and test candidates;
- validated request identity as the only first-posture actor source;
- explicit refusal to treat metadata-only `player_id` or `session_id` as proof;
- first actor kind `player`;
- actor-relative public relationship status requirements;
- route-policy expectation `request_token_required`;
- permission, validation, conflict, redaction, and unit-of-work handoff expectations;
- future runtime behavior implementation direction `implement_friends_relationship_runtime_behavior`;
- stop conditions before runtime implementation, handlers, startup wiring, protocol routes, Protobuf sources, generated output, repository/adapter changes, dependencies, migrations, authentication/session changes, event/audit tables, hosted deployment, release artifacts, public announcements, paid promotion, broader social feature expansion, or direct compatibility.

This ADR does not add runtime behavior implementation, runtime handlers, startup wiring, protocol routes, Protobuf sources, generated output, repository interface changes, PostgreSQL adapter changes, dependencies, migration changes, authentication/session behavior changes, event/audit tables, hosted deployments, release artifacts, public announcements, paid promotion, groups, parties, chat, matchmaking, match runtime, SDKs, distributed runtime, or direct Nakama/Pitaya API compatibility.

Open the next bounded work item:

```text
M-167/W-0239 Implement friends relationship runtime behavior
```

## Alternatives Considered

- Implement runtime friends behavior immediately after the PostgreSQL adapter.
- Add protocol routes and Protobuf messages together with behavior.
- Let client payloads supply actor ids.
- Treat envelope metadata `player_id` as sufficient actor proof.
- Put runtime behavior inside the friends domain module or PostgreSQL adapter package.
- Add direct Nakama-compatible friends routes.

## Rationale

Friends relationship runtime behavior has identity, privacy, and conflict-mapping risk. The first posture must not let clients select actors, turn metadata-only player ids into proof, leak private relationship graph state, or hide social behavior inside transport, protocol, or persistence layers. A gate-only ADR keeps the next implementation slice focused on application behavior and makes route/protocol expansion a separate explicit decision.

Separating this gate from implementation also lets checks reject accidental runtime handlers, Protobuf shapes, generated output, event/audit tables, or broad social feature behavior before those surfaces are ratified.

## Agent Reasoning Summary

The safest continuation from the PostgreSQL adapter is an application behavior gate. It defines how future services should derive actor identity, require protected routes, validate inputs, compute actor-relative status, map conflicts, and use unit-of-work repository handoff while preserving protocol, authentication/session, event/audit, and broad product deferrals.

## Decision Weights

```yaml
decision_weights:
  nakama_product_alignment: high
  identity_safety: high
  privacy_safety: high
  boundary_clarity: high
  agent_readability: high
  implementation_risk: low
  protocol_risk: deferred
  dependency_expansion: low
confidence: high
```

## Consequences

- `docs/friends-relationship-runtime-behavior-gate.md` and `.zh-CN.md` exist.
- `runtime.friends_relationship_runtime_behavior_gate` becomes the repository check rule for this slice.
- `M-166/W-0238` is completed.
- `M-167/W-0239 Implement friends relationship runtime behavior` becomes the next-ready work item.
- Runtime behavior implementation, protocol behavior, generated output, migrations, dependencies, event/audit tables, and authentication/session behavior remain unchanged by this ADR.

## Reversal Conditions

Revisit this decision if:

- request identity validation semantics change materially;
- route policy selects a different first protected-route posture for friends relationships;
- friends relationships expand beyond player-to-player social graph behavior before implementation;
- the repository interface or PostgreSQL adapter changes before runtime behavior implementation;
- event/audit table requirements become mandatory before runtime behavior;
- direct Nakama or Pitaya public API compatibility becomes an explicit project goal.

## Follow-Up

- Implement friends relationship runtime behavior only after this gate is accepted.
- Keep protocol routes, Protobuf messages, generated output, startup wiring, event/audit tables, groups, parties, chat, matchmaking, match runtime, SDKs, hosted surfaces, distributed runtime, and direct compatibility behind later gates.
- Preserve metadata-only identity refusal and actor-relative status behavior in the implementation tests.

