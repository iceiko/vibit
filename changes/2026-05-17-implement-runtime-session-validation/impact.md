# Impact

## Runtime

Adds `PersistentSessionValidator` under `runtime/internal/app`. It validates a persisted runtime session by:

- Requiring an already validated player identity.
- Normalizing and checking the request session id.
- Looking up an active persisted session through `session.Repository.FindActiveSessionByID`.
- Requiring active session status, expiration after `observed_at`, and actor/player match.
- Returning a `RequestIdentity` with `SessionValidated = true` only after all checks pass.

The validator is not wired into startup, route policy, WebSocket transport, or Protobuf protocol behavior in this slice.

## Authentication

Authentication remains the owner of access-token proof validation. This validator consumes already validated application identity and does not parse tokens, compute digests, issue tokens, refresh tokens, revoke tokens, or read raw credential material.

## Protocol And Transport

No protocol or transport behavior changes. The WebSocket transport remains credential-neutral and the existing Protobuf envelope remains unchanged.

## Game Server Reference Alignment

The implementation adapts Nakama's durable session lifecycle pressure by validating active, unexpired, unrevealed session state before treating a session as authenticated. It adapts Pitaya's session/handler separation by keeping validation in the application layer and outside transport acceptors and protocol routing.

Direct Nakama/Pitaya API compatibility is not added.
