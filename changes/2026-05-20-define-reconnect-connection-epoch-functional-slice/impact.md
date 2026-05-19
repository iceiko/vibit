# Impact

## Runtime

- Updates `runtime/internal/app/connection/registry.go`.
- Updates `runtime/internal/app/connection/registry_test.go`.
- Adds `ADR-0083` for the Tier 2 functional slice decision.
- Completes `M-103/W-0175` and opens the next lifecycle-closure work item.

## Application Connection Registry

The registry now treats connection epoch as a monotonic server-observed lifecycle dimension for a single connection id:

- Same active connection id and epoch is still rejected as `connection_already_open`.
- A higher epoch for the same connection id is accepted and marks earlier active records `superseded`.
- A lower or repeated epoch after a higher epoch exists is rejected as `connection_epoch_stale`.
- Superseded records retain `superseded_at` and `superseded_by_epoch`.
- `FindConnectionByID` can inspect superseded records.
- Active list methods continue to return only active bound records.

## Transport

WebSocket transport behavior is unchanged. Transport still owns accepted socket metadata and the concrete `RequestClose` handoff from W-0174, but it does not own reconnect policy.

## Authentication And Logout

Authentication behavior is unchanged:

- `LogoutAccessToken` still revokes only the verified presented token record.
- Authentication does not call WebSocket transport.
- Authentication does not own reconnect, epoch progression, socket close, or runtime session revocation.

## Protocol

- No Protobuf source is changed.
- No generated output is changed.
- The existing envelope is unchanged.
- No session carrier, reconnect route, resume route, close route, or system message is added.

## Nakama And Pitaya Reference Use

- Adopts Nakama's product lesson that realtime connection lifecycle needs explicit, testable semantics before presence and match runtime build on it.
- Adopts Pitaya's architectural lesson that session/connection lifecycle remains separate from route handlers and transport acceptors.
- Does not add direct Nakama or Pitaya API compatibility.
