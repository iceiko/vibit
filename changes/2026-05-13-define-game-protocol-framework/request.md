# Request

Date: 2026-05-13
Change ID: `define-game-protocol-framework`
Type: standard

## Maintainer Request

The maintainer asked:

```text
我们要考虑游戏的实际应用场景来进行协议框架的制定，你可以参考一下其他成熟游戏框架引擎的服务器端的实现方式和封装方式。
```

## Clarified Requirement

Define vibit's first game protocol framework before creating `.proto` files or runtime protocol handlers.

The standard should account for real game backend needs:

- Long-lived client connections.
- Session and player identity.
- Server push.
- Request/response correlation.
- Player, room, match, stream, and system targets.
- Server-authoritative commands, events, and future inputs.
- Compatibility and generated protocol discipline.
- Clear transport, protocol, application, and domain boundaries.

## User-Visible Outcome

Future agents should be able to inspect a protocol standard and machine-readable protocol manifest before writing Protobuf or runtime transport code.

## Non-Goals

- Do not create `.proto` files.
- Do not generate Protobuf output.
- Do not add Go runtime protocol handlers.
- Do not implement authentication, reconnect replay, rooms, matches, streams, presence, matchmaking, or allocation.
- Do not add new external dependencies.
- Do not change inventory command, query, event, error, or permission contracts.

## Acceptance Criteria

- [ ] Add a game protocol ADR.
- [ ] Add a machine-readable `.arch/protocol.yaml`.
- [ ] Add English and Simplified Chinese protocol standard documents.
- [ ] Update repository manifests and guides so agents know the protocol standard exists.
- [ ] Extend protocol verification to check the protocol manifest.
- [ ] Record conversation memory.
- [ ] Record verification results exactly.
