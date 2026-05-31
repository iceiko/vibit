# Conversation: Friends Relationship Protocol Route Implementation

Date: 2026-05-26
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-friends-relationship-protocol-route/`

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
- `decisions/ADR-0149-friends-relationship-protocol-route-implementation.md`
- `changes/2026-05-26-implement-friends-relationship-protocol-route/`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `.arch/protocol.yaml`
- `modules/friends/module.yaml`
- `modules/friends/AGENTS.md`
- `modules/friends/AGENTS.zh-CN.md`

## Context

`M-168/W-0240` completed the friends relationship protocol route gate. It authorized a bounded route implementation for the friends relationship route family while preserving repository interface changes, PostgreSQL adapter changes, migrations, dependencies, authentication/session changes, broad social features, event/audit tables, hosted deployment, release artifact expansion, paid promotion, public announcements, and direct Nakama/Pitaya API compatibility as deferred concerns.

The next-ready work item was `W-0241 Implement friends relationship protocol route`.

## Maintainer Narrative

The maintainer asked:

```text
继续
```

English summary: continue the next ready work.

## Agent Response Summary

The agent advanced one bounded work item and implemented the protected friends relationship protocol route family.

The work added:

- `vibit.friends.v1` Protobuf request and response messages;
- generated Go Protobuf output through Buf;
- route keys for the eight friends relationship routes;
- protocol bridge mapping between Protobuf payloads and application friends payloads;
- bootstrap route handlers that inject validated request identity and map redacted friends errors;
- PostgreSQL startup registration for the friends route family;
- transaction bypass for friends command routes because the friends service owns its unit of work;
- focused route, bridge, and startup tests;
- ADR, change spec, manifest, check-rule, and continuation updates.

## Decisions

- Complete `M-169/W-0241`.
- Accept `ADR-0149`.
- Add `runtime.friends_relationship_protocol_route_implementation`.
- Keep actor identity out of Protobuf payloads and derive it from validated request identity.
- Keep metadata-only `player_id` and `session_id` as non-proof.
- Keep groups, parties, chat, subscriptions, event/audit tables, matchmaking, match runtime, SDKs, hosted deployment, distributed runtime, and direct compatibility deferred.
- Select `W-0242 Prove friends relationship protocol route in local alpha request flow` as the next bounded direction.

## Nakama And Pitaya Reference Basis

Nakama guided the capability pressure: friends, friend requests, blocks, and relationship status are common game backend social graph primitives.

Pitaya guided the layering pressure: route registration, authenticated request/session context, protocol serialization, application handlers, backend service behavior, and persistence remain separated.

vibit adapts those lessons into its own WebSocket/Protobuf route model and application-owned service boundary. This slice does not add direct Nakama or Pitaya public API compatibility.

## Artifacts

- `proto/vibit/friends/v1/friends.proto`
- `runtime/internal/generated/proto/vibit/friends/v1/friends.pb.go`
- `runtime/internal/app/friends/routes.go`
- `runtime/internal/app/bootstrap/friends.go`
- `runtime/internal/app/bootstrap/friends_test.go`
- `runtime/internal/platform/protocol/protobuf/friends_bridge.go`
- `runtime/internal/platform/protocol/protobuf/friends_bridge_test.go`
- `runtime/internal/platform/protocol/protobuf/inventory_bridge.go`
- `runtime/internal/platform/protocol/protobuf/payload_registry.go`
- `runtime/cmd/vibit-server/main.go`
- `runtime/cmd/vibit-server/main_test.go`
- `decisions/ADR-0149-friends-relationship-protocol-route-implementation.md`
- `changes/2026-05-26-implement-friends-relationship-protocol-route/`
- `conversations/2026-05-26-friends-relationship-protocol-route-implementation.md`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `.arch/protocol.yaml`
- `modules/friends/module.yaml`
- `modules/friends/AGENTS.md`
- `modules/friends/AGENTS.zh-CN.md`

## Open Questions

- The local alpha request flow still needs a focused proof for friends relationship protocol routes.
- Groups, parties, chat, stream subscriptions, event/audit history, matchmaking, match runtime, SDKs, and direct compatibility remain deferred.
- Live PostgreSQL verification is not required by this route implementation slice; the next proof may decide whether local flow coverage needs an opt-in persistent path.

## Follow-Up

- Implement `W-0242 Prove friends relationship protocol route in local alpha request flow`.
- Preserve the current friends route Protobuf shape unless the proof reveals a concrete issue.
- Keep Nakama/Pitaya alignment explicit as capability and layering guidance, not direct API compatibility.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, DSNs with credentials, or raw private social graph data from a real user are recorded in this conversation log.
