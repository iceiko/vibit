# Runtime Protocol Adapter Boundary Standard

Status: Draft v0.1
Last updated: 2026-05-13
Scope: Runtime boundary between WebSocket transport, Protobuf protocol adaptation, application dispatch, and domain modules

This standard defines the first runtime handoff shape for vibit's WebSocket Protobuf server path.

Use this standard together with `.arch/runtime.yaml`, `.arch/protocol.yaml`, `docs/game-protocol.md`, `docs/generated-output.md`, `ADR-0014`, `ADR-0015`, `ADR-0016`, `ADR-0017`, and `ADR-0018`.

## 1. Purpose

vibit's first runtime path must stay understandable to agents.

The main risk is not that a WebSocket server cannot be written. The main risk is that agents place validation, routing, permission checks, session shortcuts, payload decoding, or domain behavior in whichever layer is most convenient at the time.

This standard prevents that drift by defining a narrow handoff between layers before runtime implementation starts.

## 2. First Runtime Flow

The first gameplay request flow is:

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

The first server-push flow is:

```text
domain event
-> application publication decision
-> protocol envelope
-> transport send
```

Server-push persistence and outbox behavior remain deferred until a separate event delivery decision exists.

## 3. Layer Ownership

### WebSocket Transport Adapter

Owner:

```text
runtime/internal/platform/transport/ws/
```

Responsibilities:

- Own `github.com/coder/websocket`.
- Accept and close WebSocket connections.
- Read and write binary frames.
- Enforce transport-level size, close, and lifecycle behavior when defined.
- Pass opaque frame bytes to the protocol adapter.
- Send encoded frame bytes returned by the protocol adapter.

Must not:

- Parse domain payloads.
- Dispatch directly to domain modules.
- Interpret command, query, event, permission, or invariant semantics.
- Construct module-specific Protobuf payloads.
- Mutate domain state.

### Protobuf Protocol Adapter

Owner:

```text
runtime/internal/platform/protocol/protobuf/
```

Responsibilities:

- Own Protobuf envelope encode and decode.
- Validate envelope-level shape.
- Convert `kind`, `module`, `name`, `request_id`, `target`, `session`, `payload_type`, and `payload` into an application route request.
- Decode generated Protobuf payloads only through generated Protobuf packages.
- Convert generated wire payloads into handwritten domain runtime payloads only through explicit protocol bridge functions.
- Convert application results and events back into generated wire payloads only through explicit protocol bridge functions.
- Map protocol errors into error envelopes.
- Preserve request correlation.

Must not:

- Open WebSocket connections.
- Own long-lived player sessions beyond envelope metadata interpretation.
- Perform domain permission decisions.
- Enforce domain invariants.
- Call domain repositories directly.
- Hide business behavior in envelope conversion.

### Application Dispatch

Owner:

```text
runtime/internal/app/
```

Responsibilities:

- Own command and query dispatch.
- Match route requests to registered application handlers.
- Create or invoke unit-of-work boundaries for state-changing commands.
- Call domain module handlers through vibit-owned interfaces.
- Map application results to protocol response metadata.
- Keep route registration explicit and generated when generation exists.

Must not:

- Parse WebSocket frames.
- Own Protobuf wire framing.
- Import `github.com/coder/websocket`.
- Import generated Protobuf packages unless a later adapter decision explicitly allows a narrow bridge.
- Hide module-specific business rules.

### Protocol-To-Domain Payload Bridges

Owner:

```text
runtime/internal/platform/protocol/protobuf/
```

The first active bridge is:

```text
runtime/internal/platform/protocol/protobuf/inventory_bridge.go
```

Responsibilities:

- Map decoded generated Protobuf payloads into handwritten module runtime request structs before application dispatch.
- Map application result payloads into generated Protobuf response payloads after application dispatch.
- Map application events into generated Protobuf event payloads when they are ready for protocol output.
- Keep field mapping explicit and test-covered.
- Preserve the original envelope metadata and request correlation.

Must not:

- Enforce domain permissions or invariants.
- Call repositories directly.
- Add authentication shortcuts.
- Change generated Protobuf output.
- Become a hidden place for business behavior.

Rules:

- Domain modules must not import generated Protobuf packages.
- `runtime/internal/app/` must not import generated Protobuf packages.
- A bridge belongs in the protocol adapter layer until a later generated bridge standard replaces the handwritten bridge.
- Unknown or not-yet-bridged routes may pass through unchanged only when their payload is already a protocol payload. Inventory routes must fail fast on mismatched payload types.

### Domain Module Runtime Logic

Owner:

```text
runtime/internal/modules/<module>/
```

Responsibilities:

- Implement handwritten command and query behavior.
- Enforce module invariants.
- Use vibit-owned repository and policy interfaces.
- Emit domain events as server facts.
- Stay aligned with module manifests and semantic contract sources.

Must not:

- Import WebSocket libraries.
- Own Protobuf framing.
- Parse envelopes.
- Depend on `google.golang.org/protobuf` directly.
- Reach into other modules' internals.

### Generated Code

Owners:

```text
runtime/internal/generated/contracts/
runtime/internal/generated/proto/
```

Responsibilities:

- Provide generated contract and wire shapes.
- Remain traceable to source contracts or `.proto` files.
- Remain immutable to non-system agents.

Must not:

- Contain handwritten runtime logic.
- Own transport, application dispatch, or domain behavior.

## 4. First Handoff Types

The first runtime implementation should introduce vibit-owned handoff types before binding to transport or generated Protobuf packages.

The conceptual type names are:

```text
TransportFrame
ProtocolEnvelope
RouteRequest
ApplicationResult
OutboundMessage
```

The first Go runtime slices implement the application-owned `RouteRequest` and `ApplicationResult` concepts under `runtime/internal/app/`, an explicit application dispatcher for command and query routes, and the Protobuf-to-application conversion under `runtime/internal/platform/protocol/protobuf/`. The remaining concepts may use idiomatic Go names when they are implemented, but they must preserve the responsibilities.

Required concepts:

- `TransportFrame` carries frame bytes and connection metadata, but no domain semantics.
- `ProtocolEnvelope` represents decoded envelope metadata and payload bytes.
- `RouteRequest` carries `kind`, `module`, `name`, `request_id`, target, session, payload identity, and decoded command/query payload.
- `ApplicationResult` carries response payloads, emitted events, and application-level errors.
- `OutboundMessage` carries protocol-level output ready for transport encoding.

## 5. Routing Rules

Runtime routing must use the structured route fields:

```text
kind
module
name
```

The rendered route key may be:

```text
<module>.<name>
```

Route registration must be explicit. When generators exist, route registration should be generated from contracts and manifests. Until generators exist, handwritten route registration must be small, local to application dispatch, and covered by change specs.

The current handwritten application dispatcher only dispatches `command` and `query` route requests. `event`, `error`, `system`, `ack`, `heartbeat`, `input`, and `state` messages are not application-dispatchable until a later standard defines their lifecycle.

Transport handlers must not build ad hoc route strings from WebSocket paths or message text.

## 6. Error Boundaries

Errors should be mapped at the layer that owns the failure:

- Transport failures belong to the transport adapter.
- Malformed envelope and unknown payload type failures belong to the protocol adapter.
- Unknown route failures belong to application dispatch.
- Permission, invariant, and domain validation failures belong to domain modules and application policy boundaries.
- Internal failures should not leak implementation details into public error payloads.

Public module errors must map to registered error catalogs.

Application errors are encoded by the Protobuf protocol adapter as `MESSAGE_KIND_ERROR` envelopes. The first active mapper is:

```text
runtime/internal/platform/protocol/protobuf/error_envelope.go
```

Rules:

- Preserve `request_id`, route metadata, target metadata, and session metadata from the application result.
- Copy stable application error code and public message into `protocolv1.Error`.
- Set `Error.request_id` to the correlated request id.
- Leave `payload_type` and `payload` empty for error envelopes.
- Treat application errors as non-retryable by default until retryability is generated from registered error catalogs.
- Do not expose non-application internal errors through this mapper without a separate protocol error handling decision.

## 7. Session And Target Boundaries

The protocol adapter may parse session and target metadata, but it must not invent authentication shortcuts.

Until player or auth modules exist:

- `player_id` can be treated as planned context for inventory protocol shape.
- Authentication/session validation remains a deferred explicit module or platform decision.
- Transport connection identity must not be treated as a durable player identity.

Target scopes beyond `player` remain reserved until the relevant module and lifecycle standard exists.

## 8. Agent Rules

Agents must:

- Read this standard before adding WebSocket transport code, Protobuf runtime adapter code, application dispatch code, or domain runtime handlers.
- Keep transport, protocol, application, domain, and generated output responsibilities separate.
- Add or update checks when a new boundary rule can be machine-verified.
- Record unavailable Go, Buf, or Protobuf tooling instead of claiming generation or runtime tests ran.

Agents must not:

- Put business behavior in WebSocket handlers.
- Put domain permission or invariant logic in Protobuf envelope decoding.
- Let domain modules import WebSocket or Protobuf runtime libraries directly.
- Let generated output become a handwritten adapter layer.
- Add authentication shortcuts inside inventory protocol work.

## 9. Verification

Current verification:

```bash
node tools/vibit check runtime
node tools/vibit check protocol
node tools/vibit check generated
node tools/vibit check all
```

`check runtime` should verify that the boundary standard, runtime manifest, protocol manifest, runtime agent guide, and repository guide all point to the runtime protocol adapter boundary before Go runtime implementation starts.

Future verification should inspect Go imports and package contents to ensure:

- `github.com/coder/websocket` imports stay under `runtime/internal/platform/transport/ws/`.
- Protobuf runtime imports stay under generated Protobuf packages and protocol adapters.
- Domain modules do not import transport or Protobuf libraries directly.
- Application dispatch does not parse WebSocket frames.
- Application and domain packages do not import platform adapters or generated Protobuf packages directly.
- Generated output does not contain handwritten adapter code.
- Inventory Protobuf/domain bridge code remains under `runtime/internal/platform/protocol/protobuf/`.

## 10. Migration Path

Before runtime implementation:

1. Keep package directories as skeletons.
2. Define the runtime protocol adapter boundary.
3. Keep `.proto` sources and generated-output rules stable.
4. Add toolchain-dependent generated output only when Go, Buf, and Protobuf tooling are available.

When implementation begins:

1. Add narrow Go handoff types first.
2. Add protocol adapter tests before wiring WebSocket transport.
3. Add application dispatch tests before domain runtime behavior grows.
4. Keep the first inventory slice small and player-scoped.
