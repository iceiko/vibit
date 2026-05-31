# ADR-0150: Friends Relationship Protocol Route Local Proof

Status: Accepted
Date: 2026-05-31
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-prove-friends-relationship-protocol-route-local-flow/`

Related conversations:

- `conversations/2026-05-26-friends-relationship-protocol-route-local-proof.md`

Related artifacts:

- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`
- `examples/local-alpha-example-client.sh`
- `examples/local-alpha-request-loop.sh`
- `examples/README.md`
- `examples/README.zh-CN.md`
- `examples/local-alpha-client/README.md`
- `examples/local-alpha-client/README.zh-CN.md`
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

`M-169/W-0241` implemented the friends relationship protocol route family. It added `vibit.friends.v1` Protobuf payloads, route keys, protocol bridge mapping, application bootstrap handlers, startup registration, and focused route tests for `friends.SendFriendRequest`, `friends.AcceptFriendRequest`, `friends.RejectFriendRequest`, `friends.RemoveFriend`, `friends.BlockPlayer`, `friends.UnblockPlayer`, `friends.ListFriendRelationships`, and `friends.GetFriendRelationshipStatus`.

The next useful step is not new friends behavior. The next useful step is to prove that the route family works inside the same local alpha request flow that already proves onboarding, device credential login, first-message connection binding, protected inventory, protected presence, storage objects, logout, and revoked-token rejection.

Nakama provides the capability reference: player-to-player friend request, acceptance, rejection, listing, status, removal, blocking, and unblocking are common social graph primitives for a game backend. Pitaya provides the layering reference: transport, session metadata, route handling, serializer/protocol adapter, application handler, service behavior, and repository handoff should remain separate. vibit adapts both references without direct public API compatibility.

## Decision

Add a proof-only local alpha E2E slice for the completed friends relationship route family.

The proof:

- extends the existing authenticated gameplay E2E fixture with friends route registration;
- adds deterministic multi-player fixture support for the local proof;
- adds a test-only friends relationship repository as an in-memory implementation of the existing `runtime/internal/modules/friends.Repository` interface;
- adds `TestFriendsRelationshipProtocolRouteLocalAlphaFlow`;
- exercises local onboarding and login for two players, first-message connection binding for both players, authenticated `friends.SendFriendRequest`, authenticated `friends.GetFriendRelationshipStatus`, authenticated `friends.AcceptFriendRequest`, authenticated `friends.ListFriendRelationships`, authenticated `friends.RemoveFriend`, authenticated `friends.BlockPlayer`, authenticated `friends.UnblockPlayer`, a second authenticated `friends.SendFriendRequest`, and authenticated `friends.RejectFriendRequest`;
- checks error and envelope redaction against access-token and one-time device-credential text;
- updates the local alpha example scripts and documentation with the exact proof command, route family, request-flow shape, redaction expectations, and Nakama/Pitaya alignment notes;
- registers `runtime.friends_relationship_protocol_route_local_proof` as the check rule for this slice.

This ADR does not add protocol messages, protocol routes, generated output, friends service behavior, repository interface changes, PostgreSQL adapter changes, migration changes, dependency additions, authentication/session behavior changes, route-protection semantic changes, hosted deployments, release artifacts, public announcements, paid promotion, event/audit tables, groups, parties, chat, matchmaking, match runtime, stream subscriptions, broadcast fanout, delivery guarantees, SDK publication, generated client libraries, production memory friends behavior, or direct Nakama/Pitaya API compatibility.

No direct Nakama or Pitaya API compatibility is added by this proof.

## Alternatives Considered

- Prove only one friends route instead of the full send/status/accept/list/remove/block/unblock/reject family.
- Add a standalone example client before proving route behavior through Go E2E tests.
- Use the PostgreSQL adapter for this proof by requiring a live database.
- Add production in-memory friends startup behavior.
- Combine the local proof with new social features such as groups, parties, chat, matchmaking, or match runtime.

## Rationale

The existing local alpha request-loop script is already the most reliable source-first proof path for new contributors. Extending it with the friends route family demonstrates the newly implemented Nakama-class social graph capability without creating a new runtime surface.

Keeping the repository implementation test-only avoids creating an accidental second production friends backend. The real production friends persistence path remains the existing PostgreSQL adapter and startup composition.

The proof is valuable because it validates the Pitaya-style separation already chosen by vibit: transport carries bytes, Protobuf adapts payloads, route protection validates request proof, application handlers inject validated identity, the friends service owns behavior, and the repository boundary owns persistence handoff.

## Agent Reasoning Summary

The smallest product-useful continuation after route implementation is an end-to-end local alpha proof. It makes the feature demonstrable to external contributors while preserving vibit's bounded workflow and avoiding premature expansion into groups, parties, chat, matchmaking, match runtime, SDKs, hosted demos, or direct compatibility.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  reference_alignment: high
  identity_safety: high
  protocol_compatibility_control: high
  boundary_clarity: high
  implementation_risk: low
  generated_output_risk: none
confidence: high
```

## Consequences

- `TestFriendsRelationshipProtocolRouteLocalAlphaFlow` proves the friends relationship route family over the existing `FrameHandler` path.
- `examples/local-alpha-example-client.sh` and `examples/local-alpha-request-loop.sh` include the friends relationship route proof.
- Examples documentation records the friends proof.
- `runtime.friends_relationship_protocol_route_local_proof` becomes the repository check rule for this slice.
- `M-170/W-0242` is completed.
- The next bounded direction is selecting the next Nakama prototype-ready capability after friends route proof.
- Broader social features, direct compatibility, production memory friends behavior, deployment, release artifact expansion, and larger product modules remain deferred.

## Reversal Conditions

Revisit this decision if:

- the friends route family changes protocol shape;
- route protection stops using request-level authenticated wrappers for friends routes;
- friends service actor derivation changes away from validated request identity;
- the local proof needs a live PostgreSQL database to catch meaningful behavior;
- direct Nakama or Pitaya public API compatibility becomes an explicit future goal through a separate ADR;
- production memory friends behavior becomes an explicit product requirement.

## Follow-Up

- Select the next Nakama prototype-ready capability after friends relationship route proof.
- Keep event/audit tables, groups, parties, chat, matchmaking, match runtime, stream subscriptions, broadcast fanout, delivery guarantees, SDK publication, hosted deployment, repository/adapter/migration changes, authentication/session changes, production memory friends behavior, and direct compatibility behind later bounded work items.
