# Impact Analysis

## Affected Modules

- `runtime`: The next implementation slice will be application-owned under `runtime/internal/app`.
- `authentication`: Authentication remains token validation service owner, not route policy owner.

## Module Ownership Impact

This direction choice does not change code ownership by itself. It selects a route-policy implementation slice that must remain under `runtime/internal/app`.

## Public Contract Impact

No public command, query, event, permission, Protobuf, WebSocket, or database contract changes are made by this direction choice.

## Reference Alignment

Nakama informs the need for authenticated session material to become usable for gameplay access only through explicit lifecycle-aware policy. Pitaya informs the separation between acceptors, session context, and route handlers.

## Compatibility Risks

The risk is over-expanding the implementation into transport, protocol, logout, reconnect, or direct external API compatibility. The chosen slice is deliberately scoped to application route policy vocabulary and checks.
