# Impact

## Runtime

- Adds `runtime/internal/app/connection/close_policy.go`.
- Adds `runtime/internal/app/connection/close_policy_test.go`.
- Adds an application-owned `ClosePolicy` that targets active bound registry records by connection id and epoch, player id, runtime session id, or access-token record id.
- Marks matched registry records invalidated and returns redacted `CloseIntent` values.
- Uses only the `mark_invalidated_only` transport action.

## Authentication

- Authentication remains token lifecycle owner only.
- `LogoutAccessToken` behavior is unchanged.
- Authentication service does not own the active connection registry or close policy.

## Transport And Protocol

- WebSocket transport remains credential-neutral and policy-neutral.
- No concrete socket close handoff is added.
- No close codes, close reason text, protocol close messages, logout routes, Protobuf changes, generated output, or WebSocket handshake changes are added.

## Nakama And Pitaya Reference Use

- Adopts Nakama's lesson that realtime socket lifecycle and token/session lifecycle need explicit server-owned semantics.
- Adopts Pitaya's lesson that acceptors, sessions, handlers, and connection management should remain separated.
- Does not add direct Nakama or Pitaya API compatibility.
