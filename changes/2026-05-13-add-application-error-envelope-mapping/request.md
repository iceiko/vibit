# Request

## Original Request

```text
继续
```

## Current Planning Context

The inventory Protobuf/domain bridge now maps successful inventory request, response, and event payloads across the protocol adapter boundary.

Before adding WebSocket transport, the runtime still needs a stable way to turn application errors into protocol error envelopes. Without that, transport code would be forced to invent response error behavior at the frame-handling layer.

## Clarified Requirement

Add a narrow application error envelope mapping in the Protobuf protocol adapter.

This change should:

- Map `app.ApplicationResult` values with `ApplicationError` into Protobuf `Envelope` values with `MESSAGE_KIND_ERROR`.
- Preserve request correlation, target metadata, and session metadata.
- Copy stable error code and public message into `protocolv1.Error`.
- Mark the first application errors as non-retryable by default.
- Keep transport, domain, and application packages independent from generated Protobuf packages.

## User-Visible Outcome

Future WebSocket transport work can send successful responses and application error responses through the same protocol adapter boundary without placing error-format decisions inside the transport handler.

## Non-Goals

- Do not add WebSocket transport.
- Do not add PostgreSQL persistence.
- Do not modify `.proto` sources or generated Protobuf files.
- Do not add new public error codes.
- Do not add structured error details.
- Do not decide transport close codes or retry policy beyond the first non-retryable application error default.
- Do not add authentication or session validation shortcuts.

## Unknowns

- Transport-level error handling and WebSocket close behavior remain deferred.
- Structured `Error.details` payloads remain deferred until there is a concrete catalog-backed details schema.
- Retryability may later become catalog-driven when error catalogs are generated into runtime code.

## Acceptance Criteria

- [ ] `app.ApplicationResult` with `ApplicationError` builds a Protobuf error envelope.
- [ ] Error envelope uses `MESSAGE_KIND_ERROR`.
- [ ] Error envelope preserves `request_id`, route metadata, target, and session.
- [ ] Error envelope contains code, message, related request id, and retryable flag.
- [ ] Successful inventory response/event mapping remains unchanged.
- [ ] `runtime/internal/app/` and `runtime/internal/modules/inventory/` still do not import generated Protobuf packages.
- [ ] Runtime and repository checks pass.
