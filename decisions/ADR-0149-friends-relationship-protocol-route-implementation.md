# ADR-0149: Friends Relationship Protocol Route Implementation

Status: Accepted
Date: 2026-05-26
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-friends-relationship-protocol-route/`

Related conversations:

- `conversations/2026-05-26-friends-relationship-protocol-route-implementation.md`

Related artifacts:

- `proto/vibit/friends/v1/friends.proto`
- `runtime/internal/generated/proto/vibit/friends/v1/friends.pb.go`
- `runtime/internal/app/friends/routes.go`
- `runtime/internal/app/bootstrap/friends.go`
- `runtime/internal/app/bootstrap/friends_test.go`
- `runtime/internal/platform/protocol/protobuf/friends_bridge.go`
- `runtime/internal/platform/protocol/protobuf/friends_bridge_test.go`
- `runtime/cmd/vibit-server/main.go`
- `runtime/cmd/vibit-server/main_test.go`
- `modules/friends/module.yaml`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `.arch/protocol.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-168/W-0240` accepted `ADR-0148` and defined the friends relationship protocol route gate. The gate authorized a later bounded implementation of protected friends relationship command/query routes over the existing WebSocket/Protobuf request flow.

The application service already exists under `runtime/internal/app/friends`. It owns validated actor identity derivation, metadata-only identity refusal, unit-of-work repository orchestration, actor-relative public status, and redacted public error mapping for send, accept, reject, remove, block, unblock, list, and status operations.

Nakama provides the capability pressure: friends, friend requests, blocks, and relationship status are core game backend social graph capabilities. Pitaya provides the layering pressure: route handling, session/identity context, serializer/protocol adapter, application service behavior, and persistence remain separate. vibit adapts both references without direct public API compatibility.

## Decision

Implement the protected friends relationship protocol route family authorized by `ADR-0148`.

The implementation adds:

- Protobuf source `proto/vibit/friends/v1/friends.proto`;
- generated Go Protobuf output `runtime/internal/generated/proto/vibit/friends/v1/friends.pb.go` through Buf;
- route keys under `runtime/internal/app/friends/routes.go`;
- routes `friends.SendFriendRequest`, `friends.AcceptFriendRequest`, `friends.RejectFriendRequest`, `friends.RemoveFriend`, `friends.BlockPlayer`, `friends.UnblockPlayer`, `friends.ListFriendRelationships`, and `friends.GetFriendRelationshipStatus`;
- protocol bridge mapping under `runtime/internal/platform/protocol/protobuf/friends_bridge.go`;
- payload registry integration for friends request and response messages;
- bootstrap handlers under `runtime/internal/app/bootstrap/friends.go`;
- PostgreSQL startup composition that constructs the friends service and registers route handlers;
- transaction bypass for friends command routes because the friends service owns its unit of work;
- focused tests for route registration, optional expected-version mapping, response mapping, validated request identity handoff, redacted handler errors, forbidden protocol fields, and startup relationship-id generation;
- repository check coverage through `runtime.friends_relationship_protocol_route_implementation`.

The Protobuf payloads intentionally omit client-supplied actor ids, session ids, access-token fields, credential material, lookup digests, verifier digests, SQL details, chat fields, group fields, party fields, matchmaking fields, and match runtime fields. Actor identity continues to come from the authenticated request identity passed into the application service.

This ADR does not add repository interface changes, PostgreSQL adapter changes, migration changes, dependency additions, authentication/session behavior changes, route-protection semantic changes, delivery guarantees, stream subscriptions, chat rooms, groups, parties, broadcast fanout, matchmaking, match runtime, operations/admin behavior, SDK publication, generated client libraries, hosted deployments, release artifacts, public announcements, paid promotion, event/audit tables, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Keep friends relationship behavior application-only until a broader social proof exists.
- Add route keys and handlers but defer Protobuf source and generated output.
- Let client payloads include actor ids.
- Put friends relationship route behavior inside WebSocket transport or PostgreSQL adapter packages.
- Combine this route implementation with a local end-to-end proof in the same slice.
- Add groups, parties, chat, presence subscriptions, or event/audit history alongside the first friends route family.

## Rationale

The route gate already selected the safe surface: validated-player friends relationship operations with target player ids and optional expected versions. Implementing that route family now makes the already-ratified friends service visible through the same WebSocket/Protobuf architecture used by existing protected gameplay routes without broadening social product scope.

Keeping actor identity out of the payload is the important safety property. The authenticated wrapper and validated request identity remain the only actor source. That preserves the earlier authentication/session boundaries and prevents metadata-only envelope fields from becoming proof.

Separating the route implementation from a later local proof keeps this slice focused on protocol shape, bridge mapping, handler registration, startup composition, and focused unit tests.

## Agent Reasoning Summary

The smallest product-useful continuation after the route gate is to wire the existing friends application service into the established route/protocol/bootstrap layers. This advances Nakama-class friends capability while preserving Pitaya-style layering and vibit's generated-output and identity-safety rules.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  reference_alignment: high
  identity_safety: high
  protocol_compatibility_control: high
  boundary_clarity: high
  implementation_risk: medium
  generated_output_risk: medium
confidence: high
```

## Consequences

- `proto/vibit/friends/v1/friends.proto` exists.
- `runtime/internal/generated/proto/vibit/friends/v1/friends.pb.go` exists and is generated through Buf.
- `runtime/internal/app/friends/routes.go` exposes the friends route keys.
- `runtime/internal/app/bootstrap/friends.go` registers application handlers.
- `runtime/internal/platform/protocol/protobuf/friends_bridge.go` maps route payloads.
- PostgreSQL startup composition registers the friends route family.
- `runtime.friends_relationship_protocol_route_implementation` becomes the repository check rule for this slice.
- `M-169/W-0241` is completed.
- The next bounded direction is `W-0242 Prove friends relationship protocol route in local alpha request flow`.
- Broader social features, event/audit history, distributed runtime, hosted deployment, release artifact expansion, SDK/client libraries, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- route policy for protected gameplay requests changes materially;
- request identity validation semantics change materially;
- friends relationship behavior changes its actor derivation model;
- generated-output standards change the Protobuf source or generated Go output path;
- direct Nakama or Pitaya public API compatibility becomes an explicit future goal through a separate ADR;
- the local proof reveals that the route surface cannot support the intended alpha request flow without changing protocol shape.

## Follow-Up

- Add a narrow local proof for friends relationship send/list/status/remove/block/unblock or a smaller representative protected-route path in `W-0242`.
- Preserve the completed Protobuf shape unless the proof finds a specific compatibility issue.
- Keep groups, parties, chat, stream subscriptions, event/audit history, matchmaking, match runtime, repository/adapter/migration changes, authentication/session changes, SDK/client libraries, hosted deployment, and direct compatibility behind later bounded work items.
