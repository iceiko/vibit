# Request

The maintainer asked:

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya
```

After recommending `implement_first_message_connection_binding`, this change implements the bounded `runtime.authentication.BindConnection` slice authorized by `docs/first-message-connection-binding-implementation-gate.md`.

## Scope

Implement first-message connection binding through a Protobuf `system` route:

```text
runtime.authentication.BindConnection
```

The client sends `vibit.authentication.v1.BindConnectionRequest` as the first binding message after WebSocket connection establishment. The Protobuf adapter decodes the message and hands it to application-owned connection binding logic. The binder validates the opaque access-token proof through the existing access-token validation boundary, then returns a connection-bound validated player identity for the server-observed connection metadata.

## User-Visible Outcome

Clients can bind a WebSocket connection to a validated player identity by sending a `BindConnection` system message with an opaque access token.

## Non-Goals

- Do not parse credentials in the WebSocket handshake, HTTP headers, cookies, query strings, or subprotocols.
- Do not change `proto/vibit/protocol/v1/envelope.proto`.
- Do not add session persistence, session tables, repositories, migrations, or cleanup jobs.
- Do not use connection-bound identity for ordinary protected routes yet; the existing request-level `AuthenticatedRequest` wrapper remains the protected-route path.
- Do not add logout-triggered active connection invalidation.
- Do not add reconnect, resume, duplicate replacement, or connection epoch policy beyond transport-observed metadata handoff.
- Do not add dependencies.
- Do not adopt direct Nakama or Pitaya public API compatibility.

## Acceptance Criteria

- [x] `BindConnectionRequest`, `BindConnectionResponse`, and `ConnectionBindingStatus` are added to authentication Protobuf source.
- [x] Generated Go Protobuf output is produced by Buf.
- [x] Application-owned connection binding validates token proof through the existing validator interface.
- [x] Protocol adapter handles the system route without normal command/query dispatch.
- [x] Successful response contains no access-token material.
- [x] Public binding errors are binding-specific and redacted.
- [x] PostgreSQL startup composition injects the binder when authentication service exists.
- [x] Memory bootstrap keeps binding unavailable.
- [x] WebSocket transport remains credential-neutral.
- [x] The work queue stops at a new next-direction confirmation gate.
