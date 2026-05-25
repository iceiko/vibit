# Impact Analysis

## Nakama Product Capability Impact

This change supports the `friends_groups_and_parties` capability family. It prepares the storage-neutral repository boundary for future friends relationship lifecycle behavior, which is a prerequisite for Nakama-class social graph features such as friend requests, friend lists, blocking, and actor-relative relationship status.

It does not add direct Nakama API compatibility.

## Pitaya Impact

Pitaya remains deferred as a future distributed architecture reference. This change does not add distributed topology, frontend/backend split, RPC, service discovery, groups, cluster routing, distributed sessions, or distributed social graph routing.

## Affected Modules

- `friends`: future module owner candidate only; no runtime module directory or Go source is created in this slice.
- `runtime`: manifests and static checks are updated.
- `workflow`: W-0234 is completed and W-0235 is opened.
- `reference`: Nakama-first roadmap state is updated.

## Module Ownership Impact

The future friends repository owner is recorded as `runtime/internal/modules/friends`, with future interface candidate `runtime/internal/modules/friends.Repository`. PostgreSQL adapter ownership is recorded as `runtime/internal/platform/persistence/postgres`.

The storage, player, authentication, protocol, WebSocket transport, and PostgreSQL platform boundaries are explicitly kept separate.

## Public Contract Impact

No public command, query, event, error, permission, route, protocol payload, or generated contract is added.

The boundary records candidate future repository vocabulary only:

- value types such as `FriendRelationship`, `FriendRelationshipPair`, `FriendRelationshipVersion`, and `FriendRelationshipRepositoryError`;
- capabilities such as `CreateOrUpdateFriendRequest`, `GetRelationshipByPair`, `ListRelationshipsForPlayer`, `AcceptFriendRequest`, `RejectFriendRequest`, `RemoveFriend`, `SetPlayerBlock`, and `ClearPlayerBlock`;
- conflict classes such as `relationship_not_found`, `duplicate_pending_request`, `blocked_relationship`, and `version_mismatch`.

## Data And Migration Impact

No migration is added or changed. The existing source remains:

```text
runtime/migrations/postgres/000007_create_friend_relationships.sql
```

The boundary records future adapter expectations for:

- `friend_relationships_pair_uq`;
- `friend_relationships_player_low_state_idx`;
- `friend_relationships_player_high_state_idx`;
- `friend_relationships_updated_at_idx`;
- `relationship_version`.

## Test Impact

Runtime behavior tests are not applicable because this change adds no runtime behavior.

Static checks cover:

- boundary docs, ADR, conversation, and change artifacts;
- W-0234 completed and W-0235 next-ready;
- manifests and guide references;
- absence of friends repository implementation before W-0235;
- absence of PostgreSQL friends adapter behavior;
- absence of friends protocol/protobuf/generated output;
- forbidden direct compatibility and token-like text patterns.

## Documentation And Memory Impact

Added:

- `docs/friends-relationship-repository-boundary.md`;
- `docs/friends-relationship-repository-boundary.zh-CN.md`;
- `decisions/ADR-0142-friends-relationship-repository-boundary.md`;
- `conversations/2026-05-25-friends-relationship-repository-boundary.md`.

Updated manifests, README files, alpha docs, product maturity docs, roadmap docs, AGENTS guides, `tools/vibit`, and `rules/check-rules.json`.

## Compatibility Risks

This change does not affect API, event, data, protocol, SDK, hosted, distributed runtime, or direct Nakama/Pitaya compatibility. The main risk is future over-implementation from the boundary; static checks and W-0235 scope are designed to prevent that.
