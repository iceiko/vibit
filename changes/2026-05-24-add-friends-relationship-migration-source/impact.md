# Impact

This change adds one SQL migration source:

- `runtime/migrations/postgres/000007_create_friend_relationships.sql`

The migration creates the `friend_relationships` table with:

- Server-generated opaque `relationship_id`.
- Canonical unordered player pair identity through `player_low_id` and `player_high_id`.
- Pair-member foreign keys to `player_accounts(player_id)`.
- `lifecycle_state` constrained to `pending`, `friends`, `rejected`, and `removed`.
- Pair-member request, response, and removal actor columns.
- Actor-specific block timestamps through `blocked_by_low_at` and `blocked_by_high_at`.
- Positive server-managed `relationship_version BIGINT`, defaulting to `1`.
- Created, updated, state-changed, rejected, removed, and block timestamps.
- Canonical pair uniqueness.
- Pair-member/lifecycle-state indexes and an updated-at index.

This is intentionally smaller than full Nakama-style friends, groups, parties, chat, matchmaking, or social graph behavior. It creates only the durable schema foundation and keeps repository, adapter, protocol, runtime behavior, and event/audit storage behind later bounded work items.

No Go code, generated code, Protobuf source, WebSocket behavior, repository interface, PostgreSQL adapter, dependency, route policy, event/audit table, group, party, chat room, matchmaking, match runtime, hosted deployment, release artifact, public announcement, paid promotion, Pitaya-style distributed runtime, or direct compatibility API is added.
