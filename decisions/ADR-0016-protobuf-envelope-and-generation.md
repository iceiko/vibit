# ADR-0016: Protobuf Envelope And Generation Config

Status: Accepted
Date: 2026-05-13
Decision Makers: Agent
Related changes:

- `changes/2026-05-13-define-protobuf-envelope-and-generation/`

Related conversations:

- `conversations/2026-05-13-protobuf-envelope-and-generation.md`

Related artifacts:

- `.arch/protocol.yaml`
- `.arch/runtime.yaml`
- `.arch/contracts.yaml`
- `.arch/conventions.yaml`
- `docs/game-protocol.md`
- `docs/game-protocol.zh-CN.md`
- `proto/README.md`
- `proto/README.zh-CN.md`
- `buf.yaml`
- `buf.gen.yaml`
- `proto/vibit/protocol/v1/envelope.proto`
- `proto/vibit/inventory/v1/inventory.proto`
- `decisions/ADR-0013-first-go-runtime-dependencies.md`
- `decisions/ADR-0014-go-runtime-layout-and-boundaries.md`
- `decisions/ADR-0015-game-protocol-framework.md`

## Context

The project has accepted Go as the first runtime language, WebSocket as the first gameplay/client protocol, Protobuf as the first wire message format, and a game-aware WebSocket-framed Protobuf envelope as the first protocol framework.

Before writing WebSocket handlers or Go runtime business logic, vibit needs concrete Protobuf source files and generation configuration that agents can inspect and verify. Without those source files, future agents would have to infer envelope shape, payload mapping, Go package output, and inventory wire messages from prose.

The first protocol source change should stay narrow:

- Define the protocol envelope source.
- Define inventory command, query, and event wire messages from registered semantic contracts.
- Add Buf configuration for future generation.
- Extend repository checks to validate protocol source shape.
- Avoid committing generated Go Protobuf output until the local toolchain is available and generated-output trace rules are stable.

## Decision

Create the first protocol envelope source at:

```text
proto/vibit/protocol/v1/envelope.proto
```

Use package:

```text
vibit.protocol.v1
```

Use Go package:

```text
github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/protocol/v1;protocolv1
```

The envelope defines:

- `Envelope`
- `MessageKind`
- `Target`
- `TargetScope`
- `Session`
- `Error`

The envelope uses explicit route fields:

```text
kind
module
name
```

Use a `bytes payload` field with an explicit `payload_type` string instead of `google.protobuf.Any` for the first protocol source.

Create the first inventory Protobuf source at:

```text
proto/vibit/inventory/v1/inventory.proto
```

Use package:

```text
vibit.inventory.v1
```

Define messages derived from registered inventory contracts:

```text
GrantItemRequest
GrantItemResponse
GetInventoryRequest
GetInventoryResponse
InventoryItem
ItemGranted
```

Create root Buf configuration:

```text
buf.yaml
buf.gen.yaml
```

`buf.yaml` uses `proto/` as the module path with `STANDARD` linting and `FILE` breaking checks. `buf.gen.yaml` uses the remote Go Protobuf plugin and writes planned generated output under:

```text
runtime/internal/generated/proto
```

Do not generate or commit Go Protobuf output in this change. Generated output remains immutable and must be produced by the accepted generator path in a later change.

## Alternatives Considered

- Use `google.protobuf.Any` for payloads immediately.
- Encode each module message directly as a top-level WebSocket frame with no envelope.
- Delay Buf configuration until the first generated Go output is committed.
- Generate Go Protobuf output immediately.
- Keep only manifest-level protocol planning and postpone `.proto` files.

## Rationale

Using `bytes payload` plus `payload_type` keeps the envelope small and explicitly inspectable. Agents can reason about routing and payload identity from stable fields without needing `Any` semantics, type URL conventions, or unpacking behavior before protocol adapters exist.

Creating `.proto` files now converts the accepted protocol framework into source artifacts that can be checked. The inventory file is narrow enough to remain aligned with current registered contracts, while the envelope file reserves game protocol concerns without implementing room, match, presence, reconnect, input, or state synchronization features.

Adding Buf configuration now makes the future generation path explicit before generated output appears. That reduces the chance that an agent chooses a conflicting generator, output path, or package option later.

Skipping generated Go output is intentional because this change is about source, config, and verification. The local environment may not have Buf, `protoc`, or Go available, and generated-output trace conventions should be enforced before generated files are committed.

## Agent Reasoning Summary

The next durable step should turn protocol decisions into verifiable source files while keeping runtime behavior untouched. The envelope should be explicit enough for agents to route and validate messages, but not broad enough to invite premature multiplayer runtime implementation.

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

- `node tools/vibit check protocol` now validates Buf configuration, the envelope source, protocol enum values, Go package options, module messages, fields, and source traces.
- `.arch/protocol.yaml`, `.arch/runtime.yaml`, `.arch/contracts.yaml`, and `.arch/conventions.yaml` now point to the first protocol source and generation configuration.
- Future WebSocket and Protobuf runtime adapters must use the envelope source as the wire-schema authority.
- Domain modules must not own or parse the envelope.
- Generated Go Protobuf output must go under `runtime/internal/generated/proto/` and must not be handwritten.
- Any future breaking envelope or inventory wire-schema change requires a change spec and ADR.

## Reversal Conditions

Revisit this decision if:

- `bytes payload` plus `payload_type` creates avoidable ambiguity once Go protocol adapters are implemented.
- Buf remote plugin usage is unsuitable for offline, pinned, or reproducible local generation.
- The generated Go package path conflicts with Go module layout or release packaging.
- The first WebSocket implementation shows that envelope and module package versioning should be combined differently.
- A concrete client generator requires a different source layout before public clients exist.

## Follow-Up

- Add generated-output trace rules before committing generated Go Protobuf files.
- Run `buf lint` and `buf generate` once Buf and the Go Protobuf toolchain are available.
- Implement the Protobuf adapter only after generated output rules and dispatch boundaries are stable.
- Define reserved field-number policy before broad public protocol evolution.
- Keep room, match, input, state sync, presence, stream, and reconnect behavior deferred until separate modules and standards require them.
