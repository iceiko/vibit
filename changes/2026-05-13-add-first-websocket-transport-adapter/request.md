# Request

## Original Request

```text
继续
```

## Current Planning Context

The runtime now has:

- Protobuf envelope decoding and encoding.
- Inventory Protobuf/domain payload bridging.
- Application dispatch.
- Inventory runtime handlers.
- Application error envelope mapping.

The next useful runtime slice is the first WebSocket transport adapter. The adapter should stay narrow: own WebSocket connections and binary frames, then delegate protocol and application behavior to injected handlers.

## Clarified Requirement

Add the first WebSocket transport adapter under `runtime/internal/platform/transport/ws/`.

This change should:

- Introduce the accepted `github.com/coder/websocket` dependency.
- Add a small transport server/handler that accepts WebSocket connections.
- Read binary frames from a connection.
- Pass frame bytes and connection metadata to an injected frame handler.
- Write returned binary frames back to the connection.
- Reject non-binary client messages.
- Keep business behavior out of transport code.

## User-Visible Outcome

Future runtime composition can bind:

```text
WebSocket frame -> Protobuf adapter -> app.Dispatcher -> domain handler -> Protobuf adapter -> WebSocket frame
```

without putting protocol or domain behavior into the WebSocket package.

## Non-Goals

- Do not implement a full server command in `cmd/vibit-server`.
- Do not add authentication, session validation, or player ownership.
- Do not add PostgreSQL persistence.
- Do not add transaction or unit-of-work wiring.
- Do not implement heartbeat, ack, reconnect, or close-code policy beyond the minimal transport loop.
- Do not implement room, match, party, stream, or presence lifecycle behavior.
- Do not put Protobuf, application dispatch, or inventory business logic inside the WebSocket transport package.

## Unknowns

- Final `/v1/ws` route mounting belongs to process wiring and remains deferred unless needed by tests.
- Connection/session id assignment policy remains minimal until the auth/session standard exists.
- Transport-level error close policy may need refinement after the first composed server loop exists.

## Acceptance Criteria

- [x] WebSocket transport package imports `github.com/coder/websocket` only inside the allowed package boundary.
- [x] Transport accepts a WebSocket request and processes binary frames.
- [x] Transport rejects text frames without invoking business handlers.
- [x] Transport frame handler receives connection metadata and copied frame bytes.
- [x] Transport writes returned binary frames back to the client.
- [x] Transport package does not import generated Protobuf, application dispatch, or domain module packages.
- [x] Runtime and repository checks pass.
