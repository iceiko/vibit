# Impact

## Runtime

- Updates `runtime/internal/platform/protocol/protobuf/inventory_bridge.go`.
- Updates `runtime/internal/platform/protocol/protobuf/authentication_bridge.go`.
- Updates focused Protobuf adapter tests.
- Adds `ADR-0084` for the Tier 2 functional slice decision.
- Completes `M-104/W-0176` and opens the next lifecycle-closure work item.

## Protocol

The existing `vibit.protocol.v1.Envelope.Session` metadata is reused. No `.proto` file is changed and no generated output is changed.

Successful `runtime.authentication.AuthenticateWithDeviceCredential` responses can now put these server-owned values into the response envelope session metadata:

- `session_id`, from the runtime session created during login.
- `player_id`, from the authenticated player result.
- `connection_id` and `connection_epoch`, from existing server-observed frame/request metadata when present.

## Application And Authentication

Authentication service behavior is unchanged. The service already returns application-owned `AuthenticationResult.SessionID` after successful session creation. This slice only makes the Protobuf response builder expose that value in the existing envelope session metadata.

Session id remains metadata for lifecycle correlation. It is not a proof by itself.

## Transport

WebSocket transport behavior is unchanged. Transport still owns accepted connection metadata and concrete close handoff only; it does not parse session proof, authenticate handshakes, or route durable sessions.

## Nakama And Pitaya Reference Use

- Adopts Nakama's product lesson that clients need explicit session metadata after authentication before richer realtime lifecycle behavior can be usable.
- Adopts Pitaya's architectural lesson that session context should be visible to routing/session layers without moving auth into the transport acceptor.
- Does not add direct Nakama or Pitaya API compatibility.
