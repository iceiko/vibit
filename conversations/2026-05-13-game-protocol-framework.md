# Conversation: Game Protocol Framework

Date: 2026-05-13  
Participants: Maintainer, Agent  
Related changes:

- `changes/2026-05-13-define-game-protocol-framework/`

Related artifacts:

- `docs/game-protocol.md`
- `docs/game-protocol.zh-CN.md`
- `.arch/protocol.yaml`
- `decisions/ADR-0015-game-protocol-framework.md`

## Context

The project had already ratified Go, WebSocket, Protobuf, PostgreSQL, a modular monolith first runtime model, and `inventory` as the first proof module. The next step was to define protocol shape before creating `.proto` files or runtime protocol handlers.

## Maintainer Narrative

The maintainer emphasized that protocol framework design must be grounded in real game application scenarios and should reference mature server-side game frameworks and engines:

```text
我们要考虑游戏的实际应用场景来进行协议框架的制定，你可以参考一下其他成熟游戏框架引擎的服务器端的实现方式和封装方式。
```

## Agent Response Summary

The agent treated this as a protocol architecture decision rather than immediate Protobuf or Go runtime implementation.

The agent reviewed mature game-server patterns:

- Nakama: realtime socket surface, presence, streams, parties, matchmaker flows, and server-authoritative multiplayer.
- Colyseus: room lifecycle, server-owned state, state synchronization, and room-scoped client messages.
- Pitaya and Pomelo: connector or frontend network ownership, route dispatch, sessions, server push, handlers, and backend remotes.
- Skynet: service/message-dispatch architecture and concurrency boundaries.
- Agones and Open Match: allocation and matchmaking lifecycle outside the first in-game protocol.

The agent recommended adopting a game-aware WebSocket-framed Protobuf envelope with explicit route, session, target, authority, error, and compatibility semantics, while deferring room state sync, matchmaking, allocation, reconnect replay, presence, and streams until concrete modules require them.

## Decisions

- Add `ADR-0015` for the game protocol framework.
- Add `.arch/protocol.yaml` as the machine-readable protocol manifest.
- Add `docs/game-protocol.md` and `docs/game-protocol.zh-CN.md` as the human-readable standard.
- Keep the first inventory slice player-scoped and avoid implementing broad multiplayer concepts prematurely.
- Do not create `.proto` files or Go runtime protocol handlers in this change.

## Artifacts

- `docs/game-protocol.md`
- `docs/game-protocol.zh-CN.md`
- `.arch/protocol.yaml`
- `decisions/ADR-0015-game-protocol-framework.md`
- `changes/2026-05-13-define-game-protocol-framework/`

## Open Questions

- Exact Protobuf envelope field names and field numbers.
- Whether the first endpoint path `/v1/ws` needs revision before transport implementation begins.
- Authentication/session handshake behavior before a player or auth module exists.
- Reconnect replay semantics.
- Room, match, input, state synchronization, presence, and stream standards.

## Follow-Up

- Extend protocol checks to validate `.arch/protocol.yaml`.
- Define the first envelope `.proto` after this standard is in place.
- Define reserved field policy before broad Protobuf generation.

## Redaction Notes

No secrets, tokens, account details, or private operational data were included.
