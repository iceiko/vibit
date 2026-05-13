# Request

## Original Request

```text
继续推进
```

## Current Planning Context

The previous runtime slice added the first inventory repository, policy, handler, and application dispatcher boundary. It intentionally left Protobuf-to-domain payload mapping for a later protocol adapter bridge.

The next useful slice is to connect the existing generated Protobuf inventory payloads to the handwritten inventory runtime request structs without letting generated Protobuf types leak into `runtime/internal/app/` or `runtime/internal/modules/inventory/`.

## Clarified Requirement

Add a narrow inventory Protobuf/domain bridge under the Protobuf protocol adapter package.

This change should:

- Convert generated inventory request payloads decoded from a protocol envelope into handwritten inventory runtime request structs.
- Convert inventory application result payloads and emitted events back into generated Protobuf payloads for protocol envelope construction.
- Keep inventory domain logic independent from generated Protobuf packages.
- Keep application dispatch independent from generated Protobuf packages.
- Prove the bridge through focused tests that include dispatcher integration.

## User-Visible Outcome

Future agents can follow a concrete runtime path:

```text
Protobuf envelope -> app.RouteRequest -> inventory runtime payload -> app.Dispatcher -> inventory handler -> application result -> Protobuf envelope
```

This advances the first complete backend slice while still deferring WebSocket transport and PostgreSQL persistence.

## Non-Goals

- Do not add WebSocket transport.
- Do not add PostgreSQL persistence.
- Do not add transaction or unit-of-work wiring.
- Do not add new third-party dependencies.
- Do not modify generated Protobuf files.
- Do not add or change public command, query, event, error, or permission contracts.
- Do not add authentication or player session validation shortcuts.

## Unknowns

- Exact transport error envelope behavior remains deferred until WebSocket transport work.
- Exact durable event publication and outbox behavior remains deferred.
- Generated route registration remains planned but not implemented in this change.

## Acceptance Criteria

- [ ] Inventory generated request payloads map to handwritten inventory runtime request structs.
- [ ] Inventory runtime response payloads map to generated Protobuf response payloads.
- [ ] Inventory runtime events map to generated Protobuf event payloads.
- [ ] `runtime/internal/modules/inventory/` still does not import generated Protobuf packages or Protobuf runtime packages.
- [ ] `runtime/internal/app/` still does not import generated Protobuf packages or Protobuf runtime packages.
- [ ] Tests cover `GrantItem`, `GetInventory`, response mapping, event mapping, and dispatcher integration.
- [ ] Runtime and repository checks pass.
