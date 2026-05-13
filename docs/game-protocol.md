# Game Protocol Framework Standard

Status: Draft v0.1  
Last updated: 2026-05-13  
Scope: First gameplay/client protocol framework for vibit  
Canonical decision: `ADR-0015`

This standard defines how vibit's first game-facing protocol should be shaped before runtime transport handlers or generated clients are created.

The goal is not to clone an existing game backend. The goal is to absorb mature game-server patterns into an agent-native protocol surface that agents can inspect, extend, and verify.

## Problem

A game backend protocol is not only an RPC layer.

Real games need long-lived connections, authenticated sessions, player-scoped requests, server push, error semantics, reconnect behavior, room or match targeting, authoritative input handling, and future allocation or matchmaking flows.

If vibit starts with a generic request/response protocol, later agents may hide game behavior in route strings, transport handlers, or module internals. That would weaken the project's core promise: agents should understand and modify the server from explicit contracts and manifests.

## External References

Mature systems suggest useful patterns:

- Nakama separates ordinary APIs from realtime socket use and models presence, streams, parties, matchmaker flows, and server-authoritative multiplayer.
- Colyseus centers multiplayer behavior around rooms, server-owned state, synchronized state changes, and client messages inside a joined room.
- Pitaya and Pomelo show the value of connector or frontend network ownership, route-based dispatch, connection sessions, server push, handlers, and backend remotes.
- Skynet shows a service/message-dispatch style that is useful for server architecture and concurrency boundaries, but it should not define vibit's first client wire protocol.
- Agones and Open Match show that game server allocation and match construction are separate lifecycle concerns. They matter later, but they should not be forced into the first in-game WebSocket envelope.

vibit should borrow vocabulary and proven boundaries, not copy any framework wholesale.

Reference reading:

- Nakama Multiplayer Engine: `https://heroiclabs.com/docs/nakama/concepts/multiplayer/`
- Colyseus State Synchronization: `https://docs.colyseus.io/state`
- Pitaya Documentation: `https://pitaya.readthedocs.io/`
- Pomelo repository: `https://github.com/NetEase/pomelo`
- Skynet repository: `https://github.com/cloudwu/skynet`
- Agones GameServerAllocation: `https://agones.dev/site/docs/reference/gameserverallocation/`
- Open Match Match Function: `https://openmatch.dev/site/docs/guides/matchmaker/matchfunction/`

## Protocol Model

vibit's first game protocol is a WebSocket-framed Protobuf envelope.

The envelope is owned by platform protocol code, not by domain modules. Domain modules own semantic commands, queries, events, invariants, permissions, and errors through vibit manifests and contract sources.

The first protocol surface should use one gameplay WebSocket endpoint, with the initial endpoint path planned as:

```text
/v1/ws
```

The endpoint path may be finalized when the first transport implementation begins. Changing it after clients exist is a compatibility-sensitive protocol change.

## Envelope Responsibilities

The envelope must make these concepts explicit:

- Protocol version.
- Message kind.
- Request correlation.
- Module identity.
- Operation name.
- Game target scope.
- Session metadata.
- Payload encoding identity.
- Error mapping.

The envelope must not contain domain business logic.

The first envelope Protobuf source is:

```text
proto/vibit/protocol/v1/envelope.proto
```

The envelope uses a `bytes` payload plus explicit `payload_type` metadata. This keeps routing and payload identity inspectable without requiring `google.protobuf.Any` in the first protocol source.

The initial message kinds are:

```text
command
query
event
error
ack
heartbeat
system
input
state
```

Only `command`, `query`, `event`, `error`, and `system` need to be used by the first inventory slice. `ack`, `heartbeat`, `input`, and `state` are reserved because mature game protocols need them, but implementing them should wait for concrete requirements.

## Routing

Routing must be agent-readable.

The first route identity is structured fields:

```text
kind
module
name
```

The semantic route key may be rendered as:

```text
<module>.<name>
```

Examples:

```text
inventory.GrantItem
inventory.GetInventory
inventory.ItemGranted
```

Do not bury business meaning in an opaque path string. Transport handlers decode the envelope and hand the route to application dispatch. Domain modules do not parse WebSocket frames or route strings.

## Session Model

The protocol model distinguishes transport connection from logical session and player identity:

- `connection_id` is transport-local and may change after reconnect.
- `session_id` is the authenticated logical session when authentication exists.
- `player_id` is the domain identity of the player when available.
- `connection_epoch` or an equivalent reconnect version is reserved for future reconnect rules.

Until the player/auth module exists, inventory protocol work may treat authenticated identity as planned context rather than inventing an authentication shortcut.

## Target Model

Game messages often target more than a single RPC handler.

The envelope should reserve target metadata with these scope values:

```text
global
player
party
room
match
stream
system
```

The first inventory slice should use `player` scope for player-owned inventory operations.

Room, match, party, and stream behavior is reserved. It must not be implemented as hidden ad hoc fields before the relevant module, contracts, and runtime lifecycle rules exist.

## Client-To-Server Messages

Client-to-server messages may represent:

- Commands: client intent to change server state.
- Queries: client request to read server state.
- Inputs: realtime authoritative gameplay input, reserved until a match or room module exists.
- System messages: protocol-level negotiation or lifecycle messages.

Client messages are requests, not facts. The server decides whether the command, query, or input is valid.

## Server-To-Client Messages

Server-to-client messages may represent:

- Command or query responses correlated to a `request_id`.
- Domain event pushes.
- Error envelopes.
- System/session messages.
- State snapshots or patches, reserved until a room or match state-sync standard exists.
- Acknowledgements and heartbeats, reserved until the transport lifecycle requires them above WebSocket ping/pong.

Domain events are server facts. Client code must not be allowed to publish facts directly into the domain model.

## Error Model

Protocol errors must map to registered vibit error catalogs when they refer to public module behavior.

The first wire error shape should include:

- Stable error code.
- Human-readable message.
- Related `request_id`, when available.
- Retryability signal.
- Optional structured details.

Transport errors, malformed envelopes, permission errors, invariant failures, and unknown routes must be distinguishable.

## Protobuf Sources

The first protocol source files are:

```text
proto/vibit/protocol/v1/envelope.proto
proto/vibit/inventory/v1/inventory.proto
```

The first Buf configuration files are:

```text
buf.yaml
buf.gen.yaml
```

Buf is accepted as the Protobuf orchestration tool by `ADR-0013`. Generated Go Protobuf output belongs under:

```text
runtime/internal/generated/proto/
```

Generated output must be produced from `.proto` sources. Agents must not create or hand-edit generated Go Protobuf files to work around source or generator problems.

Generated output rules are defined by `docs/generated-output.md` and `ADR-0017`.

Runtime transport, protocol adapter, application dispatch, generated code, and domain module handoff rules are defined by `docs/runtime-protocol-adapter.md` and `ADR-0018`.

## Compatibility Rules

The envelope Protobuf package is:

```text
vibit.protocol.v1
```

The first Protobuf package version remains:

```text
vibit.<module>.v1
```

Envelope versioning and module message versioning are related but separate:

- Envelope version governs framing, routing, target, session, and system semantics.
- Module package version governs module command, query, event, and payload wire shapes.

Before first `.proto` files are created, the protocol standard should define reserved fields and extension discipline. Once `.proto` files exist:

- Do not reuse removed Protobuf field numbers.
- Do not rename public field semantics casually.
- Prefer additive evolution.
- Breaking changes require a change spec and ADR.

## Agent Rules

Before adding or changing `.proto` files, runtime protocol handlers, or generated protocol code, agents must read:

- `docs/game-protocol.md`
- `.arch/protocol.yaml`
- `ADR-0009`
- `ADR-0014`
- `ADR-0015`
- `ADR-0017`
- `ADR-0018`
- `.arch/contracts.yaml`
- `docs/generated-output.md`
- `docs/runtime-protocol-adapter.md`
- The affected module manifest and contract files

Agents must preserve these boundaries:

- WebSocket transport belongs under `runtime/internal/platform/transport/ws/`.
- Protobuf envelope encode/decode belongs under `runtime/internal/platform/protocol/protobuf/`.
- Application dispatch belongs under `runtime/internal/app/`.
- Domain behavior belongs under `runtime/internal/modules/<module>/`.
- Generated Protobuf output belongs under `runtime/internal/generated/proto/`.
- Protobuf source files belong under `proto/vibit/<module>/v1/`.
- The protocol envelope source belongs under `proto/vibit/protocol/v1/`.
- Root Buf generation configuration belongs in `buf.yaml` and `buf.gen.yaml`.

## First Implementation Guidance

The first inventory runtime slice should prove the protocol framework without pretending to implement all multiplayer concerns.

Recommended first flow:

```text
WebSocket frame
-> Protobuf envelope
-> kind/module/name dispatch
-> generated command or query payload
-> application dispatcher
-> inventory handler
-> response envelope
```

Do not implement room state sync, matchmaking, allocation, reconnect replay, presence, or streams in the inventory slice unless a new change spec and ADR make them part of the target.

## Verification Direction

Current verification:

```bash
node tools/vibit check protocol
```

The protocol check should verify:

- `.arch/protocol.yaml` exists.
- The manifest references `ADR-0015`.
- The manifest records WebSocket, Protobuf, envelope routing, session identity, target scopes, message kinds, and server authority.
- `buf.yaml` and `buf.gen.yaml` contain the accepted source root, generation output, lint, and breaking-check settings.
- The first envelope `.proto` declares expected envelope, target, session, error, kind, and target-scope shape.
- Protobuf source files align with registered contracts once they exist.

Future checks should validate generated output traceability, reserved field policy, and runtime dispatch ownership.
