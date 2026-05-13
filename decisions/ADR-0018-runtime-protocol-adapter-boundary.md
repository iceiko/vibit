# ADR-0018: Runtime Protocol Adapter Boundary

Status: Accepted
Date: 2026-05-13
Decision Makers: Agent
Related changes:

- `changes/2026-05-13-define-runtime-protocol-adapter-boundary/`

Related conversations:

- `conversations/2026-05-13-runtime-protocol-adapter-boundary.md`

Related artifacts:

- `docs/runtime-protocol-adapter.md`
- `docs/runtime-protocol-adapter.zh-CN.md`
- `.arch/runtime.yaml`
- `.arch/protocol.yaml`
- `runtime/AGENTS.md`
- `tools/vibit`

## Context

The project has accepted Go, WebSocket, Protobuf, generated output rules, and the first protocol envelope. The local environment still lacks Go, Buf, and `protoc`, so generated output and Go runtime code should not be faked.

The next useful preparation is to define the runtime boundary between WebSocket transport, Protobuf protocol adaptation, application dispatch, generated code, and domain modules before implementation starts.

## Decision

Define `docs/runtime-protocol-adapter.md` as the runtime protocol adapter boundary standard, with a paired Simplified Chinese translation.

The first runtime flow is:

```text
websocket frame
-> transport frame
-> protocol envelope
-> route request
-> application dispatch
-> domain command or query handler
-> application result
-> protocol envelope
-> websocket frame
```

Layer ownership is:

- `runtime/internal/platform/transport/ws/` owns WebSocket connection and frame adaptation.
- `runtime/internal/platform/protocol/protobuf/` owns Protobuf envelope encode/decode and protocol-level conversion.
- `runtime/internal/app/` owns command/query dispatch and unit-of-work orchestration.
- `runtime/internal/modules/<module>/` owns handwritten domain behavior.
- `runtime/internal/generated/contracts/` and `runtime/internal/generated/proto/` own generated shapes only.

Extend repository checks so `node tools/vibit check runtime` verifies that the runtime protocol adapter boundary standard is wired into runtime and protocol manifests before implementation starts.

## Alternatives Considered

- Start writing Go runtime code immediately.
- Wait until generated Protobuf files exist before defining adapter boundaries.
- Put all decode, dispatch, and handler wiring into the WebSocket package.
- Let generated Protobuf packages become the application dispatch contract.

## Rationale

Starting implementation without the boundary would increase the chance that business behavior gets placed in transport or protocol conversion code.

Defining the handoff before code is cheaper than refactoring misplaced runtime responsibilities later. It also gives future agents a direct intake document for the first runtime adapter slice.

Keeping generated Protobuf packages out of the application dispatch contract preserves vibit's semantic contract source model. Protobuf owns wire shape; application dispatch owns route and handler invocation.

## Agent Reasoning Summary

Because toolchain-dependent generation cannot run locally, the highest-value next step is not to fake generated files or Go tests. It is to make the next implementation boundary explicit and checkable.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  human_ergonomics: medium
  implementation_cost: low
  reversibility: high
  long_term_maintainability: high
confidence: high
```

## Consequences

- Future runtime code must preserve the transport, protocol, application, domain, and generated output boundaries.
- `node tools/vibit check runtime` should fail if the boundary standard is disconnected from the runtime manifest or agent guides.
- The first implementation should introduce narrow handoff types before WebSocket wiring grows.
- Go, Buf, and Protobuf generation remain not verified until those tools are available.

## Reversal Conditions

Revisit this decision if early runtime implementation shows the handoff is too abstract, if generated dispatch becomes simpler and still preserves semantic contract ownership, or if a future runtime model changes the transport/protocol/application split.

## Follow-Up

- Add initial Go handoff types once Go tooling is available.
- Add protocol adapter tests before WebSocket transport wiring.
- Add import-boundary checks once Go source files exist.
