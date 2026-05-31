# Conversation: Friends Relationship Protocol Route Gate

Date: 2026-05-26
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-friends-relationship-protocol-route-gate/`

Related artifacts:

- `docs/friends-relationship-protocol-route-gate.md`
- `docs/friends-relationship-protocol-route-gate.zh-CN.md`
- `decisions/ADR-0148-friends-relationship-protocol-route-gate.md`
- `changes/2026-05-26-define-friends-relationship-protocol-route-gate/`
- `conversations/2026-05-26-friends-relationship-protocol-route-gate.md`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/friends/module.yaml`
- `modules/friends/AGENTS.md`
- `modules/friends/AGENTS.zh-CN.md`

## Context

`M-167/W-0239` implemented application-owned friends relationship runtime behavior under `runtime/internal/app/friends`. The service covered friend request send, accept, reject, remove, block, unblock, list, and status behavior with validated request identity handoff, metadata-only identity refusal, unit-of-work repository access, and redacted public errors.

The next-ready work item was `W-0240 Define friends relationship protocol route gate`.

## Maintainer Narrative

The maintainer asked to continue the repository's next-ready work item.

```text
继续
```

The standing product direction remained Nakama-first friends/groups/parties capability coverage with Pitaya deferred as a future distributed architecture reference.

## Agent Response Summary

The agent advanced one bounded work item and defined a protocol route gate only. The slice recorded the future protected friends relationship route family, candidate Protobuf package and generated-output posture, request identity handoff, public error mapping, service handoff, Nakama/Pitaya reference mapping, verification expectations, and stop conditions.

The future route family is:

- `friends.SendFriendRequest`
- `friends.AcceptFriendRequest`
- `friends.RejectFriendRequest`
- `friends.RemoveFriend`
- `friends.BlockPlayer`
- `friends.UnblockPlayer`
- `friends.ListFriendRelationships`
- `friends.GetFriendRelationshipStatus`

The gate keeps actor identity derived only from validated authenticated request identity. Metadata-only `player_id` and `session_id` remain insufficient proof, and client payloads must not provide actor ids as proof.

## Decisions

- Complete `M-168/W-0240`.
- Accept `ADR-0148`.
- Add `runtime.friends_relationship_protocol_route_gate`.
- Record `docs/friends-relationship-protocol-route-gate.md` and `docs/friends-relationship-protocol-route-gate.zh-CN.md`.
- Preserve route implementation, Protobuf source, generated output, startup wiring, repository interface changes, adapter changes, dependencies, migrations, authentication/session behavior changes, event/audit tables, broad social features, hosted surfaces, SDKs, distributed runtime, and direct Nakama/Pitaya API compatibility as deferred concerns.
- Select `M-169/W-0241 Implement friends relationship protocol route` as the next-ready work item.

## Nakama And Pitaya Reference Basis

Nakama guided the route family pressure: friend request, accept, reject, remove, block, unblock, list, and status operations are core social graph access patterns.

Pitaya guided the layering pressure: route handlers, request/session identity, protocol adapter mapping, application behavior, and persistence must remain separated.

vibit adapted those references into protected request-token route planning with validated authenticated request identity, no metadata-only proof, and no direct public API compatibility.

## Artifacts

- `docs/friends-relationship-protocol-route-gate.md`
- `docs/friends-relationship-protocol-route-gate.zh-CN.md`
- `decisions/ADR-0148-friends-relationship-protocol-route-gate.md`
- `changes/2026-05-26-define-friends-relationship-protocol-route-gate/`
- `conversations/2026-05-26-friends-relationship-protocol-route-gate.md`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/friends/module.yaml`
- `modules/friends/AGENTS.md`
- `modules/friends/AGENTS.zh-CN.md`

## Open Questions

- Protobuf request and response field details remain deferred to `W-0241`.
- Protocol bridge implementation and bootstrap route registration remain deferred.
- Route-specific generated output remains deferred.
- Event/audit tables remain deferred.
- Chat rooms, groups, parties, broadcast fanout, matchmaking, match runtime social context, SDK publication, hosted surfaces, distributed runtime, and direct compatibility remain deferred.

## Follow-Up

The next-ready work item is:

```text
M-169/W-0241 Implement friends relationship protocol route
```

That next slice may implement only the protected friends relationship route family authorized by `ADR-0148`.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, or raw private social graph data from a real user are recorded in this conversation log.
