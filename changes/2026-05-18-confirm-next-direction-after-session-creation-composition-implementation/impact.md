# Impact Analysis

## Affected Modules

- `runtime`
- `authentication`

## Module Ownership Impact

This confirmation step changes planning state only. It does not move ownership. Future route-policy work remains application-owned under `runtime/internal/app`.

## Public Contract Impact

No command, query, event, permission, error, Protobuf, generated output, WebSocket, or database contract is changed by this confirmation step.

## Reference Alignment

Nakama informs the need to treat authenticated session material as relevant to gameplay access while keeping logout, refresh, session management, and active socket behavior as separate policy choices.

Pitaya informs the separation of acceptor, session context, and route handler execution. vibit adapts that by keeping route policy in the application layer and keeping transport credential-neutral.

## Risk

The risk is over-expanding route policy into logout, reconnect, protocol carriers, or social/realtime behavior. The selected direction intentionally creates a gate-only milestone first.
