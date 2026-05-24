# Impact

## Nakama Product Capability Impact

This change advances the Nakama-first `friends_groups_and_parties` capability family by defining the future friends relationship lifecycle semantics before implementation. It adopts the product need for friend request, accept, reject, remove, block, unblock, list, and relationship-status behavior, but adapts it into vibit's contract-first and server-authoritative model.

It does not copy external public routes, runtime API names, payloads, data models, SDK behavior, or compatibility guarantees.

## Pitaya Impact

Pitaya remains deferred as a future distributed architecture reference. This change does not introduce frontend/backend roles, RPC, service discovery, cluster routing, distributed sessions, group broadcast, or server-to-server messaging.

## Public Contract Impact

The gate records future semantic vocabulary only:

- commands: `SendFriendRequest`, `AcceptFriendRequest`, `RejectFriendRequest`, `RemoveFriend`, `BlockPlayer`, `UnblockPlayer`;
- queries: `ListFriendRelationships`, `GetFriendRelationshipStatus`;
- events: `FriendRequestCreated`, `FriendRequestAccepted`, `FriendRequestRejected`, `FriendRemoved`, `PlayerBlocked`, `PlayerUnblocked`;
- errors: `FRIENDSHIP_INVALID_TARGET`, `FRIENDSHIP_SELF_TARGET_FORBIDDEN`, `FRIENDSHIP_DUPLICATE_REQUEST`, `FRIENDSHIP_BLOCKED_RELATIONSHIP`, `FRIENDSHIP_INVALID_TRANSITION`, `FRIENDSHIP_RELATIONSHIP_NOT_FOUND`, `FRIENDSHIP_TARGET_NOT_FOUND`, `FRIENDSHIP_METADATA_IDENTITY_NOT_AUTHENTICATED`;
- permission: `validated_player_identity`.

No contract source files, protocol routes, Protobuf source, generated output, runtime handlers, or public API behavior are added.

## Data And Migration Impact

No migration, table, repository interface, PostgreSQL adapter, SQL execution behavior, or durable state is added. The change opens `W-0232 Define friends relationship persistence schema gate` to define persistence shape before migration source.

## Test Impact

The change records future positive, negative, permission/authentication, persistence/transaction, protocol, redaction, concurrency, and E2E test expectations. It adds repository checks for the gate but no runtime tests because no runtime behavior changes.

## Documentation And Memory Impact

The change adds the friends relationship lifecycle gate standard, paired Simplified Chinese translation, ADR, conversation log, filled change artifacts, manifest progression, roadmap updates, AGENTS guidance, `tools/vibit` check coverage, and rule catalog registration.

## Compatibility Risks

Breaking API risk is none because no public API or wire payload is added. Direct Nakama/Pitaya compatibility risk is controlled by explicit non-authorization and checker coverage. Runtime, persistence, SDK, hosted, release, and distributed runtime risks are deferred to future bounded work items.
