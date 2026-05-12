# ADR-0015: Game Protocol Framework

Status: Accepted  
Date: 2026-05-13  
Decision Makers: Maintainer, Agent  
Related changes:

- `changes/2026-05-13-define-game-protocol-framework/`

Related conversations:

- `conversations/2026-05-13-game-protocol-framework.md`

Related artifacts:

- `.arch/protocol.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `docs/game-protocol.md`
- `docs/game-protocol.zh-CN.md`
- `proto/README.md`
- `decisions/ADR-0009-websocket-protobuf-client-protocol.md`
- `decisions/ADR-0014-go-runtime-layout-and-boundaries.md`

## Context

The maintainer asked that vibit's protocol framework be shaped around real game application scenarios and that mature game frameworks and engines be used as references for server-side implementation and encapsulation patterns.

The project has already accepted Go as the first runtime language, WebSocket as the first gameplay/client protocol, Protobuf as the first wire format, and a modular monolith as the first server model. The next risk is starting with a protocol that is technically typed but too generic for games.

Mature game backends show repeated concerns:

- Long-lived realtime connections.
- Explicit session and player identity.
- Request/response correlation.
- Server push.
- Presence, streams, rooms, parties, matches, and matchmaker flows.
- Server-authoritative command or input validation.
- Clear separation between network connectors, routing, application handlers, domain logic, and distributed lifecycle concerns.

vibit needs to absorb those concerns without cloning Nakama, Colyseus, Pitaya, Pomelo, Skynet, Agones, or Open Match.

## Decision

vibit's first game protocol framework is a WebSocket-framed Protobuf envelope with explicit routing, session, target, authority, and compatibility semantics.

The first gameplay WebSocket endpoint is planned as:

```text
/v1/ws
```

The endpoint path remains planned until transport implementation begins. Changing it after clients exist is compatibility-sensitive and requires a change spec and ADR.

The envelope must expose structured route identity:

```text
kind
module
name
```

The semantic route key may be rendered as:

```text
<module>.<name>
```

The envelope must reserve or define:

- Protocol version.
- Request correlation.
- Session metadata.
- Target metadata.
- Payload encoding identity.
- Error mapping.
- Message kind.

The initial message kinds are:

```text
command
query
event
error
system
ack
heartbeat
input
state
```

The first inventory slice should use only `command`, `query`, `event`, `error`, and `system`. `ack`, `heartbeat`, `input`, and `state` are reserved until concrete room, match, realtime input, or lifecycle requirements exist.

The target model reserves these scopes:

```text
global
player
party
room
match
stream
system
```

The first inventory slice uses `player` scope.

The server is authoritative. Client-to-server commands, queries, and future inputs are requests, not facts. Server-pushed domain events are facts.

`.arch/protocol.yaml` is the machine-readable protocol framework manifest. `docs/game-protocol.md` is the human and agent-readable standard.

## Alternatives Considered

- Raw per-message Protobuf frames without an envelope.
- Generic RPC over WebSocket only.
- JSON envelope with Protobuf payload.
- Pomelo/Pitaya-style opaque route strings as the primary contract.
- Colyseus-style room and state synchronization as the first abstraction.
- Nakama-like broad realtime API surface immediately.
- Actor or service-message runtime protocol first.
- Game server allocation and matchmaking lifecycle first.

## Rationale

A raw per-message Protobuf protocol is compact, but it makes routing, target, session, and error semantics harder for agents to inspect before generated code exists.

A generic RPC protocol is easy to implement but underrepresents real game concerns. It risks treating server push, rooms, matches, presence, reconnects, and authoritative input as ad hoc extensions.

Opaque route strings are familiar in some game server frameworks, but vibit should prefer structured route fields because agents can inspect, validate, and generate them without parsing business meaning from strings.

Starting with room/state synchronization or a broad realtime API would copy too much from a specific framework before vibit has proven its contract-first agent workflow. The first inventory slice needs player-scoped command/query behavior, not a full multiplayer state engine.

Separating allocation and matchmaking lifecycle from the first in-game protocol keeps the initial runtime focused while preserving a future path for Agones/Open Match-style concerns.

## Agent Reasoning Summary

The protocol should be game-aware from the beginning, but the first implementation should stay narrow. The durable standard should reserve concepts that real game servers need while preventing agents from implementing them opportunistically inside transport code or route strings.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  human_ergonomics: high
  implementation_cost: medium
  reversibility: medium
  long_term_maintainability: high
confidence: high
```

## Consequences

- Future `.proto` envelope work must align with `.arch/protocol.yaml` and `docs/game-protocol.md`.
- `node tools/vibit check protocol` should verify the protocol manifest as well as manifest-to-Protobuf alignment.
- Transport handlers must not hide business behavior.
- Domain modules must not parse WebSocket frames, own Protobuf framing, or infer game targets from opaque route strings.
- The first inventory runtime slice can remain player-scoped while still fitting a future room, match, stream, or matchmaking model.
- Reconnect replay, presence, stream subscriptions, room state sync, matchmaker integration, and allocation remain deferred until separate changes define them.

## Reversal Conditions

Revisit this decision if:

- The first WebSocket implementation shows the envelope is too heavy for simple client usage.
- Generated Protobuf code becomes unclear because envelope and module package versioning are separated incorrectly.
- A real room or match module needs a different target model.
- Client compatibility requirements make `/v1/ws` unsuitable before clients exist.
- A future runtime distribution model requires a different dispatch boundary.

## Follow-Up

- Add the first envelope `.proto` only after the protocol framework standard is accepted.
- Extend `node tools/vibit check protocol` to validate envelope `.proto` shape once it exists.
- Define reserved Protobuf field numbers and extension policy before broad generation.
- Define authentication/session protocol behavior when a player or auth module is introduced.
- Define room, match, input, state synchronization, presence, reconnect, and stream standards only when concrete modules require them.
