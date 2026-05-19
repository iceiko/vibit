# Impact

## Runtime

- Updates `runtime/internal/app/connection/registry.go`.
- Updates `runtime/internal/platform/transport/ws/server.go`.
- Updates `runtime/cmd/vibit-server/main.go`.
- Adds startup composition adapters:
  - `runtime/cmd/vibit-server/connection_lifecycle.go`
  - `runtime/cmd/vibit-server/connection_binding_registry.go`
- Adds focused tests for presence snapshots, WebSocket lifecycle observation, and startup registry binding.
- Adds `ADR-0085` for the Tier 2 presence lifecycle functional slice decision.

## Presence

The first presence lifecycle behavior is derived from server-owned connection state:

- WebSocket accepted connections can notify a credential-neutral lifecycle observer on open and close.
- PostgreSQL runtime startup wires that observer to the existing in-memory active connection registry.
- Successful first-message connection binding records validated player identity in the registry.
- `connection.InMemoryRegistry.PresenceForPlayer` returns an online/offline player presence snapshot from active bound records.

## Protocol

No Protobuf source, generated output, envelope field, or client-visible presence route is added in this slice. The next bounded step may expose the registry-backed snapshot as a protected protocol query.

## Application And Authentication

Authentication service behavior is unchanged. The startup binding adapter composes the existing `app.ConnectionBinder` with the existing connection registry after validation succeeds.

The current `RequestIdentity` does not carry access-token record id or validated runtime session id from access-token validation. This slice therefore does not claim full token-record or session-presence linkage in startup wiring.

## Transport

WebSocket transport remains credential-neutral. It owns only server-observed connection lifecycle observation and does not parse headers, cookies, query strings, subprotocols, bearer values, access tokens, session ids, player ids, or route proofs.

## Nakama And Pitaya Reference Use

- Adopts Nakama's product lesson that presence depends on explicit authenticated socket lifecycle state.
- Adopts Pitaya's architecture lesson that acceptor lifecycle, session/binding context, and handler behavior should remain separated.
- Does not add direct Nakama or Pitaya API compatibility.

