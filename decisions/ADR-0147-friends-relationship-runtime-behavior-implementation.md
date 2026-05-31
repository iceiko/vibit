# ADR-0147: Friends Relationship Runtime Behavior Implementation

Status: Accepted
Date: 2026-05-26
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-friends-relationship-runtime-behavior/`

Related conversations:

- `conversations/2026-05-26-friends-relationship-runtime-behavior-implementation.md`

Related artifacts:

- `runtime/internal/app/friends/service.go`
- `runtime/internal/app/friends/service_test.go`
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

`M-166/W-0238` defined the friends relationship runtime behavior gate after the friends relationship repository interface and PostgreSQL adapter were in place. The gate authorized a later bounded application-owned runtime behavior implementation under `runtime/internal/app/friends`.

The friends module owns storage-neutral value types, lifecycle/status vocabulary, and `runtime/internal/modules/friends.Repository`. The PostgreSQL adapter owns SQL mapping and unit-of-work repository construction. Runtime behavior must therefore compose the existing pieces through application dependencies and the unit-of-work boundary without importing persistence, protocol, transport, generated Protobuf, migration, chat, groups, parties, matchmaking, match runtime, SDK, hosted, or distributed runtime packages.

Nakama motivates durable friends relationships as a core social graph primitive. Pitaya reinforces separating route/session context, handlers, and persistence. vibit adapts those references through validated request identity and application-owned behavior, not direct public API compatibility.

## Decision

Implement friends relationship runtime behavior under:

```text
runtime/internal/app/friends
```

The implementation adds:

- `Service` and `NewService`;
- `SendFriendRequest`;
- `AcceptFriendRequest`;
- `RejectFriendRequest`;
- `RemoveFriend`;
- `BlockPlayer`;
- `UnblockPlayer`;
- `ListFriendRelationships`;
- `GetFriendRelationshipStatus`;
- request and result vocabulary for application callers;
- stable redacted public error codes for later handlers to map;
- validated-player actor derivation from `app.RequestIdentity`;
- metadata-only `player_id` and `session_id` refusal before id generation, unit-of-work access, or repository mutation;
- server-owned relationship id generation for send request behavior;
- incoming-only accept/reject behavior for pending requests;
- actor-relative public status computation for pending, friends, rejected, removed, blocked-by-actor, blocked-actor, and mutual-blocked states;
- successful `none` status mapping when a relationship status lookup is not found;
- unit-of-work friends repository handoff through `NewFriendRelationshipRepository`;
- redacted conflict mapping for invalid request, unauthenticated, forbidden, target not found, relationship not found, duplicate request, already friends, blocked relationship, invalid transition, version mismatch, and unavailable cases;
- focused fake-repository tests.

This ADR does not add runtime protocol handlers, route registration, startup wiring, Protobuf sources, generated output, repository interface changes, PostgreSQL adapter changes, dependencies, migration changes, authentication/session behavior changes, event/audit tables, chat rooms, groups, parties, broadcast fanout, matchmaking, match runtime, operations/admin behavior, SDK publication, generated client libraries, hosted deployments, release artifacts, public announcements, paid promotion, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Add WebSocket/Protobuf routes together with the service.
- Put the service under `runtime/internal/modules/friends`.
- Let client payloads provide actor ids.
- Treat metadata-only envelope/session fields as sufficient actor proof.
- Expose repository or SQL details directly to future callers.
- Implement event/audit tables together with runtime behavior.
- Copy Nakama friends API semantics directly.

## Rationale

The useful next step after the PostgreSQL adapter is an application service that future protocol slices can call. This keeps identity enforcement, validation, unit-of-work orchestration, public failure classes, and actor-relative status shaping in the application layer while preserving route and generated-output decisions for later gates.

Rejecting metadata-only identity before repository access prevents a client from turning unauthenticated envelope or session metadata into social graph authority. Computing public status relative to the requesting actor avoids exposing another player's private social graph framing. Fake-repository tests are sufficient for this slice because persistence mapping is already covered by the repository interface and PostgreSQL adapter tests.

## Agent Reasoning Summary

The safest continuation from `W-0238` was an application service only. It makes the friends relationship lifecycle usable inside the runtime while keeping protocol exposure, generated files, migrations, dependencies, event/audit tables, and broader social features behind later work items.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  identity_safety: high
  privacy_safety: high
  boundary_clarity: high
  agent_readability: high
  implementation_risk: medium
  protocol_risk: deferred
  dependency_expansion: low
confidence: high
```

## Consequences

- `runtime/internal/app/friends/service.go` exists.
- `runtime/internal/app/friends/service_test.go` exists.
- `runtime.friends_relationship_runtime_behavior_implementation` becomes the repository check rule for this slice.
- `M-167/W-0239` is completed.
- The next bounded direction is a friends relationship protocol route gate.
- Protocol routes, Protobuf source/generated output, startup wiring, migrations, dependencies, authentication/session semantics, event/audit tables, chat, groups, parties, matchmaking, match runtime, SDKs, hosted deployment, release artifacts, distributed runtime, and direct compatibility remain unchanged by this ADR.

## Reversal Conditions

Revisit this decision if:

- request identity validation semantics change materially;
- the friends repository interface changes materially;
- route policy selects a different first protected-route posture for friends relationships;
- friends relationships expand beyond player-to-player social graph behavior before protocol exposure;
- event/audit tables become mandatory before route exposure;
- direct Nakama or Pitaya public API compatibility becomes an explicit project goal.

## Follow-Up

- Define friends relationship protocol routes and generated output in a later gate before exposing this service over WebSocket/Protobuf.
- Preserve metadata-only identity refusal in future handlers.
- Keep chat, groups, parties, broadcast fanout, matchmaking, match runtime social context, event/audit tables, SDKs, hosted surfaces, distributed runtime, and direct compatibility behind later gates.
