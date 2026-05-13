# Request

## Original Request

Continue ten work items unless maintainer confirmation is required.

## Clarified Requirement

Advance `W-0007` by adding the first WebSocket Protobuf application composition adapter. The adapter must decode WebSocket frame payload bytes as Protobuf envelopes, bridge command and query payloads into application route requests, dispatch through `app.Dispatcher`, and encode application results or application errors back into Protobuf envelope bytes.

## User-Visible Outcome

Agents and maintainers can exercise the first in-process request loop without mounting the `/v1/ws` endpoint yet:

```text
binary frame bytes -> Protobuf envelope -> application dispatch -> Protobuf envelope bytes
```

## Non-Goals

- Do not change the protocol envelope shape.
- Do not add authentication or session validation.
- Do not mount `/v1/ws`.
- Do not add persistence, PostgreSQL transactions, or MinIO integration.
- Do not move generated Protobuf imports into application or domain packages.

## Unknowns

- None for this bounded step.

## Acceptance Criteria

- [x] Decode binary frame payloads as `vibit.protocol.v1.Envelope` outside the WebSocket transport package.
- [x] Convert command and query envelopes to application route requests through the existing Protobuf/domain bridge.
- [x] Dispatch requests through `app.Dispatcher`.
- [x] Encode successful application payloads back to Protobuf envelope bytes.
- [x] Encode application errors back to Protobuf error envelope bytes.
- [x] Preserve transport, protocol, application, domain, and generated output boundaries.
- [x] Cover the composition with Go tests.
