# Impact

## Runtime

Adds an application-owned presence query handler under `runtime/internal/app/presence` and registers it only in the PostgreSQL runtime composition where authentication, route protection, connection binding, and presence lifecycle wiring already exist.

The query is request-token protected by the existing route protection posture. It accepts an optional requested `player_id`, but allows only the authenticated player to query their own snapshot in this first slice.

## Protocol

Adds `proto/vibit/presence/v1/presence.proto` with `GetPlayerPresenceRequest`, `GetPlayerPresenceResponse`, and bounded connection metadata. Generated Go Protobuf output is produced by Buf under `runtime/internal/generated/proto/vibit/presence/v1/presence.pb.go`.

No existing envelope field changes. No authentication, inventory, player, or transport wire schema changes.

## Data

No migrations and no durable presence store. The query reads the existing single-process in-memory active connection registry.

## Boundaries

The slice does not add presence subscriptions, broadcasts, durable/distributed presence, chat/social/matchmaking/match runtime behavior, reconnect/resume tokens, logout-triggered socket close, runtime session revocation, operations/admin behavior, dependencies, SDK behavior, or direct Nakama/Pitaya API compatibility.
