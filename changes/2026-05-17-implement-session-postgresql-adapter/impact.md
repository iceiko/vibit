# Impact

## Runtime

Adds a PostgreSQL adapter for the storage-neutral runtime session repository and exposes it through the PostgreSQL unit-of-work factory.

The adapter maps the existing repository methods to `runtime_sessions` SQL:

- Create one runtime session row.
- Get one runtime session row by `session_id`.
- Find one active, unexpired runtime session by `session_id`.
- Update `last_seen_at`.
- Mark an active session expired.
- Revoke a session with `revoked_at` and `revocation_reason`.
- List active, unexpired sessions for one player with a bounded limit.

## Authentication

Authentication remains the owner of access-token proof validation. The adapter only stores and returns the opaque `access_token_record_id` linkage already present in `runtime_sessions`.

## Protocol And Transport

No protocol or transport behavior changes. The WebSocket transport remains credential-neutral and the existing Protobuf envelope remains unchanged.

## Game Server Reference Alignment

The implementation adapts Nakama's durable lifecycle session pressure and Pitaya's transport/session/handler separation. It does not copy either public API.
