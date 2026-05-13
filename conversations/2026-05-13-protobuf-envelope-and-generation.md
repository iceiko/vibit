# Conversation: Protobuf Envelope And Generation

Date: 2026-05-13
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-13-define-protobuf-envelope-and-generation/`

Related artifacts:

- `buf.yaml`
- `buf.gen.yaml`
- `proto/vibit/protocol/v1/envelope.proto`
- `proto/vibit/inventory/v1/inventory.proto`
- `decisions/ADR-0016-protobuf-envelope-and-generation.md`

## Context

After reviewing project progress and setting the next stage, the project moved into the Protocol-to-Runtime Bootstrap phase. The previous change defined the game protocol framework but intentionally did not create `.proto` files or runtime handlers.

The accepted direction before this conversation was:

- Go is the first server runtime language.
- WebSocket is the first gameplay/client protocol.
- Protobuf is the first wire format.
- PostgreSQL is the first authoritative durable store.
- The first server model is a modular monolith.
- The first proof module is `inventory`.
- Mature game-server patterns should inform protocol design without cloning a broad framework surface.

## Maintainer Narrative

The maintainer asked to continue development after the game protocol framework was defined:

```text
继续推进
```

Earlier maintainer intent still applied: decisions that can create major branches should be discussed, while professional sub-decisions inside ratified directions may be made by the agent.

## Agent Response Summary

The agent treated the next step as a bounded protocol-source and generation-configuration change, not runtime business implementation.

The agent added:

- Root Buf configuration.
- A protocol envelope `.proto`.
- An inventory `.proto` aligned with registered command, query, and event contracts.
- Tooling checks for Buf configuration, envelope shape, enum values, Go package options, source traces, module messages, and contract-aligned fields.
- Architecture, protocol, contract, agent-guide, README, and Protobuf README updates.
- `ADR-0016` to record the durable envelope and generation configuration decision.

The agent did not generate Go Protobuf output because local Buf, `protoc`, and Go toolchains were not available. Generated output remains planned and must not be handwritten.

## Decisions

- Create `proto/vibit/protocol/v1/envelope.proto` as the first protocol envelope source.
- Create `proto/vibit/inventory/v1/inventory.proto` as the first inventory wire-schema source.
- Use `bytes payload` plus explicit `payload_type` in the first envelope instead of `google.protobuf.Any`.
- Add `buf.yaml` and `buf.gen.yaml` now, before generated output is committed.
- Keep generated Go Protobuf output planned under `runtime/internal/generated/proto/`.
- Extend `node tools/vibit check protocol` so future agents can verify protocol source shape without relying on visual inspection.

## Artifacts

- `buf.yaml`
- `buf.gen.yaml`
- `proto/vibit/protocol/v1/envelope.proto`
- `proto/vibit/inventory/v1/inventory.proto`
- `decisions/ADR-0016-protobuf-envelope-and-generation.md`
- `changes/2026-05-13-define-protobuf-envelope-and-generation/`

## Open Questions

- Exact generated-output trace marker format for Go Protobuf files.
- Whether Buf remote plugin usage should be pinned or replaced by local plugin invocation for reproducible generation.
- Whether `payload_type` should use a fixed naming convention such as `vibit.inventory.v1.GrantItemRequest`.
- Reserved field-number policy before public client compatibility begins.
- When to implement the Go Protobuf adapter and WebSocket transport adapter.

## Follow-Up

- Add generated-output trace rules before committing generated Go Protobuf output.
- Run Buf lint/generation when the toolchain is available.
- Define the first runtime Protobuf adapter after generation rules are stable.
- Keep room, match, state sync, presence, streams, reconnect, and realtime input deferred until their own module standards exist.

## Redaction Notes

No secrets, tokens, account details, or private operational data were included.
