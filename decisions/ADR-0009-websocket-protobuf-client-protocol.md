# ADR-0009: WebSocket And Protobuf Client Protocol

Status: Accepted
Date: 2026-05-12
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-12-ratify-go-websocket-protobuf-runtime/`

Related conversations:

- `conversations/2026-05-12-go-websocket-protobuf-direction.md`

Related artifacts:

- `.arch/runtime.yaml`
- `.arch/contracts.yaml`
- `contracts/`
- `decisions/ADR-0004-minimal-server-instance-model.md`
- `decisions/ADR-0007-yaml-contract-source-format.md`
- `decisions/ADR-0008-go-server-runtime-language.md`

## Context

vibit's first domain is game-server architecture. The maintainer clarified that the first client protocol should be WebSocket and that Protobuf should be used for protocol messages.

The project already has YAML contract source files for business-level commands, queries, events, errors, and permissions. The maintainer selected the option where vibit manifests remain the source for business semantics and Protobuf becomes the wire schema for client/server protocol messages.

## Decision

WebSocket is the first gameplay/client protocol for vibit.

Protobuf is the first wire message format for client/server protocol messages.

vibit manifests and contract source files remain the source of truth for business semantics: command/query/event identity, ownership, permissions, errors, invariants, and module boundaries.

Protobuf files are the source of truth for wire-level message shape and compatibility. Tooling must check alignment between vibit manifests and Protobuf schemas before the protocol surface grows.

HTTP may be introduced later for health checks, admin APIs, development tooling, gateway integration, or observability. HTTP is not the first gameplay/client API unless a future ADR supersedes this decision.

## Alternatives Considered

- HTTP first with WebSocket later.
- WebSocket first with JSON messages.
- WebSocket first with Protobuf messages.
- Protobuf as the only source of truth for both business semantics and wire schema.
- Delay protocol selection until after the first runtime implementation.

## Rationale

Game servers usually need long-lived sessions, server push, bidirectional messages, and low-latency interaction patterns. WebSocket is a practical first protocol for that shape while keeping the first server instance simpler than a distributed actor or cluster runtime.

Protobuf gives compact, typed, versionable wire messages and has strong Go tooling. It is a better long-term protocol surface than ad hoc JSON once the framework is expected to support real client/server communication.

At the same time, Protobuf alone should not own vibit's business semantics. Module ownership, permissions, invariants, generated-file policy, and agent-readable intent need to remain explicit in vibit manifests and contracts. This keeps the architecture legible to agents and avoids burying server design in transport schemas.

## Agent Reasoning Summary

WebSocket plus Protobuf matches the expected realtime game-server path. Keeping vibit manifests as the semantic source preserves the agent-native architecture goal, while Protobuf owns the wire contract where it is strongest.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  human_ergonomics: medium
  implementation_cost: medium
  reversibility: medium
  long_term_maintainability: high
confidence: high
```

## Consequences

- First runtime protocol work should target WebSocket.
- First protocol message generation should target Protobuf.
- Transport adapters must not contain domain behavior.
- Domain modules must not depend directly on third-party WebSocket or Protobuf runtime libraries.
- Tooling must eventually verify that Protobuf messages match registered vibit commands, queries, events, errors, and permissions.
- HTTP references in existing docs should be treated as later operational surface, not gameplay protocol direction.

## Reversal Conditions

Revisit this decision if the first real client/server slice shows that WebSocket or Protobuf creates unacceptable agent ambiguity, weakens contract verification, or blocks essential client compatibility.

## Follow-Up

- Define the first `.proto` layout and package naming standard.
- Choose Protobuf tooling through the dependency adoption process.
- Choose the first WebSocket library through the dependency adoption process.
- Add an alignment check between vibit contract manifests and `.proto` files before broad protocol generation.
