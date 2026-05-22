# Impact

This change adds one SQL migration source:

- `runtime/migrations/postgres/000006_create_storage_objects.sql`

The migration creates the `storage_objects` table with:

- Server-generated opaque `object_id`.
- Player-only owner posture through `owner_kind = 'player'`.
- `owner_id` linked to `player_accounts(player_id)`.
- Logical identity over `owner_kind + owner_id + collection + object_key`.
- Small JSON game-state payloads in `value_json JSONB`, constrained to top-level JSON objects.
- Positive server-managed `version BIGINT`, defaulting to `1`.
- Created, updated, and optional soft-delete timestamps.
- A unique active logical identity index.
- Owner/collection and updated-at lookup indexes.

This is intentionally smaller than full Nakama-style storage objects or any Pitaya-style handler/RPC integration. It creates only durable schema foundation and keeps repository, adapter, protocol, and runtime behavior behind later bounded work items.

No Go code, generated code, Protobuf source, WebSocket behavior, repository interface, PostgreSQL adapter, dependency, route policy, object/blob storage, S3-compatible object storage, hosted deployment, release artifact, public announcement, paid promotion, or direct compatibility API is added.
