# Impact

## Runtime

- Adds `runtime/internal/platform/transport/ws/close_handoff.go`.
- Adds `runtime/internal/platform/transport/ws/close_handoff_test.go`.
- Updates `runtime/internal/platform/transport/ws/server.go` so accepted sockets are registered in a single-process transport-owned socket table and unregistered when the connection exits.
- Updates `runtime/internal/platform/transport/ws/server_test.go` with a live WebSocket close handoff check.

## Transport

- WebSocket transport now owns concrete socket close mechanics for an already accepted socket.
- The first target is server-observed `connection_id + connection_epoch`.
- The handoff returns redacted outcomes and does not expose database errors, credential material, player ids, session ids, token record ids, headers, cookies, query strings, subprotocol values, or remote addresses.
- The concrete first action uses `CloseNow`, deliberately avoiding close code and reason-text selection.

## Application Connection Policy

- `runtime/internal/app/connection` remains the owner of close decisions and registry lifecycle markers.
- The existing `ClosePolicy` still emits `mark_invalidated_only`; this change does not reinterpret it as an automatic socket close.

## Authentication

- Authentication remains token lifecycle owner only.
- `LogoutAccessToken` behavior is unchanged and does not close sockets.
- Authentication service does not import or call WebSocket transport.

## Protocol

- No Protobuf source is changed.
- No generated output is changed.
- The existing protocol envelope is unchanged.
- No protocol close message, session carrier, kick route, disconnect route, or logout-triggered close behavior is added.

## Nakama And Pitaya Reference Use

- Adopts Nakama's lifecycle lesson that realtime socket closure is explicit and separate from token logout.
- Adopts Pitaya's layering lesson that concrete connection mechanics remain separate from route handlers and application policy.
- Does not add direct Nakama or Pitaya API compatibility.
