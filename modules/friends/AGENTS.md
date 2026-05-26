# friends Module Agent Guide

Status: Draft v0.1

## When To Use This Module

Use this module for friends relationship social graph repository vocabulary and storage-neutral value types.

The current implemented slice is intentionally narrow:

- `runtime/internal/modules/friends.Repository`
- `FriendRelationship`, canonical unordered pair, actor, block state, lifecycle state, public status, and version value types
- request, accept, reject, remove, block, unblock, pair lookup, and player-scoped list input/result types
- optimistic conflict classes and redacted repository errors
- normalization helpers and focused Go tests

`M-163 Friends Relationship Repository Interface Implementation` is completed by `W-0235`. The check rule is `runtime.friends_relationship_repository_interface_implementation`.

`M-164 Friends Relationship PostgreSQL Adapter Gate` is completed by `W-0236`. The check rule is `runtime.friends_relationship_postgresql_adapter_gate`.

The repository next work item is `W-0237 Implement friends relationship PostgreSQL adapter`. Do not implement runtime friendship behavior, protocol routes, Protobuf source, generated output, event/audit tables, startup wiring beyond the bounded adapter handoff, groups, parties, chat, matchmaking, match runtime, SDK publication, hosted deployments, distributed runtime, or direct Nakama/Pitaya API compatibility in this module until a later bounded work item authorizes that scope.

## When Not To Use This Module

Do not use this module for:

- WebSocket, HTTP, Protobuf, or generated wire behavior.
- PostgreSQL adapter implementation or SQL execution.
- Runtime friend request, accept, reject, remove, block, unblock, list, or status behavior.
- Player account lifecycle.
- Authentication, token formats, credential storage, or session validation.
- Storage objects, inventory, presence, realtime messaging, chat, groups, parties, matchmaking, or match runtime.
- Event/audit tables for friendship history.
- Direct Nakama or Pitaya public API compatibility.

If a requirement needs one of those concepts, create or update the owning boundary instead of adding hidden ownership here.

## Extension Points

- Repository interface: `runtime/internal/modules/friends.Repository`
- Repository value types: `FriendRelationship`, `FriendRelationshipPair`, `FriendRelationshipActor`, `FriendRelationshipBlockState`, `FriendRelationshipVersion`
- Lifecycle vocabulary: `pending`, `friends`, `rejected`, `removed`
- Public status vocabulary: `pending`, `friends`, `blocked`, `ended`
- Normalizers: relationship records, list results, pair identity, actors, block state, and repository inputs
- Tests: `runtime/internal/modules/friends/repository_test.go`
- PostgreSQL adapter owner candidate: `runtime/internal/platform/persistence/postgres`
- PostgreSQL adapter gate: `docs/friends-relationship-postgresql-adapter-gate.md`
- PostgreSQL adapter implementation follow-up: `W-0237 Implement friends relationship PostgreSQL adapter`

The first public runtime commands and queries are still deferred. Future runtime behavior must derive actor identity from validated request identity before calling this repository interface; client-supplied player ids are not authentication proof.

## Forbidden Shortcuts

- Do not bypass boundaries declared in `module.yaml`.
- Do not add unregistered public commands, queries, events, errors, or permissions.
- Do not add PostgreSQL adapter code under this module.
- Do not import `pgx`, `database/sql`, WebSocket packages, generated Protobuf packages, SDK packages, or distributed runtime packages into this module.
- Do not execute SQL or mention concrete SQL statements in friends module source.
- Do not change migrations from this module.
- Do not wire runtime handlers, startup composition, route policy, protocol adapters, or transport behavior from this module.
- Do not add protocol routes, Protobuf sources, or generated output without a later bounded work item.
- Do not store raw credentials, raw tokens, verifier material, lookup digests, verifier digests, cookies, transport subprotocols, connection metadata, chat, group, party, matchmaking, match runtime, Pitaya server routing, or direct Nakama/Pitaya compatibility fields in friends value types.
- Do not treat pair member ids, actor ids, `player_id`, `session_id`, or transport metadata as authenticated proof.

## Required Tests

See `tests.required` in `module.yaml`.

For the current repository interface slice, tests must cover:

- Repository interface storage neutrality.
- Closed lifecycle state, public status, and conflict vocabularies.
- Canonical unordered pair normalization and self-relationship rejection.
- Returned record normalization.
- Returned list result normalization and slice copying.
- Send/block/unblock self-target rejection.
- Accept/reject/remove actor-in-pair requirements.
- Expected version validation and pointer copying.
- List pagination bounds.
- Redacted conflict and repository errors.
- Absence of secret, transport, protocol, distributed, chat, group, party, match, and direct compatibility fields.

Run `node tools/vibit check runtime` after changing friends runtime source. When Go is available, also run `cd runtime && go test ./...`.
