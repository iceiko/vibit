# ADR-0148: Friends Relationship Protocol Route Gate

Status: Accepted
Date: 2026-05-26
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-friends-relationship-protocol-route-gate/`

Related conversations:

- `conversations/2026-05-26-friends-relationship-protocol-route-gate.md`

Related artifacts:

- `docs/friends-relationship-protocol-route-gate.md`
- `docs/friends-relationship-protocol-route-gate.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/friends/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-167/W-0239` implemented application-owned friends relationship runtime behavior under `runtime/internal/app/friends`. The service sends, accepts, rejects, removes, blocks, unblocks, lists, and reads relationship status while deriving actor identity from validated `app.RequestIdentity`, rejecting metadata-only identity before id generation or repository access, using unit-of-work repository handoff, computing actor-relative public status, and mapping repository conflicts to redacted public errors.

The next bounded step is to define how this service may later become client-facing over the existing WebSocket/Protobuf protocol. That requires route names, message-shape posture, protected-route policy, protocol adapter ownership, generated-output expectations, public error mapping, redaction expectations, and stop conditions before `.proto` or route code exists.

Nakama is the reference for the product capability: friends, friend requests, blocks, and player-scoped social graph status are common game backend surfaces. Pitaya is the reference for layering: acceptors, sessions, route handlers, serializers, and backend services remain separate. vibit adapts both references without copying either public API.

## Decision

Accept `docs/friends-relationship-protocol-route-gate.md` as the gate for future friends relationship protocol routes.

The gate records:

- future route family `friends.SendFriendRequest`, `friends.AcceptFriendRequest`, `friends.RejectFriendRequest`, `friends.RemoveFriend`, `friends.BlockPlayer`, `friends.UnblockPlayer`, `friends.ListFriendRelationships`, and `friends.GetFriendRelationshipStatus`;
- command/query split for mutation and read routes;
- future Protobuf source candidate `proto/vibit/friends/v1/friends.proto`;
- future generated Go output candidate `runtime/internal/generated/proto/vibit/friends/v1/friends.pb.go`;
- candidate request/response messages and fields;
- unchanged envelope posture;
- `request_token_required` protected-route posture with authenticated wrapper requirement;
- validated authenticated request identity as the only actor source;
- refusal to treat metadata-only `player_id` or `session_id` as proof;
- no client-supplied actor id proof in the first posture;
- future protocol bridge and application handler ownership;
- generated-output and no-hand-edit expectations;
- public service error mapping for the `FRIENDSHIP_*` public error family;
- Nakama/Pitaya reference mapping;
- required future tests;
- stop conditions before route implementation, Protobuf source, generated output, startup wiring, repository/adapter/migration/dependency changes, authentication/session changes, delivery guarantees, stream subscriptions, chat, groups, parties, matchmaking, match runtime, event/audit tables, hosted deployment, release artifacts, public announcements, paid promotion, SDK work, distributed runtime, or direct compatibility.

This ADR does not add protocol route implementation, Protobuf source, generated output, startup wiring, runtime handlers, repository interface changes, PostgreSQL adapter changes, dependencies, migration changes, authentication/session behavior changes, delivery guarantees, stream subscriptions, chat rooms, groups, parties, broadcast fanout, matchmaking, match runtime, operations/admin behavior, SDK publication, generated client libraries, hosted deployments, release artifacts, public announcements, paid promotion, event/audit tables, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Add friends `.proto` messages and route handlers in the same slice.
- Copy Nakama friends route names and response models directly.
- Use a Pitaya-style opaque route string instead of vibit's `kind/module/name` envelope route identity.
- Allow client payloads to include actor ids.
- Treat envelope `player_id` or `session_id` metadata as sufficient actor proof.
- Put friends route behavior in WebSocket transport, generated bridge, or PostgreSQL adapter packages.

## Rationale

The service behavior exists, but exposing it over the client protocol creates compatibility, identity, and generated-output risk. A gate-only slice lets the project preserve useful Nakama-style friends capability while keeping Pitaya-style route/session/handler separation and vibit's own contract-first protocol discipline.

The future route family should be player-actor scoped because the service is already built around validated player identity. Broader social features such as groups, parties, chat targeting, subscriptions, event/audit history, matchmaking integration, match social context, and SDK/client libraries are product-useful, but they need separate contracts and checks.

## Agent Reasoning Summary

The safest continuation from `W-0239` is a protocol route gate. It records the future route and message shape, connects the work to Nakama/Pitaya capability mapping, and prevents agents from jumping directly into `.proto`, generated code, startup wiring, event/audit tables, or direct compatibility.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  reference_alignment: high
  identity_safety: high
  protocol_compatibility_control: high
  boundary_clarity: high
  implementation_risk: low
  generated_output_risk: deferred
confidence: high
```

## Consequences

- `docs/friends-relationship-protocol-route-gate.md` and `.zh-CN.md` exist.
- `runtime.friends_relationship_protocol_route_gate` becomes the repository check rule for this slice.
- `M-168/W-0240` is completed.
- The next bounded direction is friends relationship protocol route implementation.
- Protobuf source/generated output, route implementation, startup wiring, repository/adapter/migration/dependency changes, authentication/session semantics, delivery guarantees, event/audit tables, hosted deployment, release artifacts, SDK/client libraries, distributed runtime, and direct compatibility remain unchanged by this ADR.

## Reversal Conditions

Revisit this decision if:

- route policy changes the first protected-route posture for friends;
- request identity validation semantics change materially;
- friends runtime behavior changes its actor derivation model;
- the project chooses direct Nakama or Pitaya API compatibility through a future ADR;
- friends relationship capability expands to groups, parties, chat, subscriptions, matchmaking, match runtime, event/audit history, or SDK/client libraries before protocol exposure;
- generated-output standards change before the future implementation slice.

## Follow-Up

- Implement friends relationship protocol routes only after this gate is accepted.
- Keep Protobuf source/generated output, startup wiring, groups, parties, chat, stream subscriptions, event/audit tables, matchmaking, match runtime, SDK publication, generated client libraries, hosted deployment, and direct compatibility behind later bounded work items.
- Preserve metadata-only identity refusal in future route tests.
