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

`M-165 Friends Relationship PostgreSQL Adapter Implementation` is completed by `W-0237`. The check rule is `runtime.friends_relationship_postgresql_adapter_implementation`.

`M-166 Friends Relationship Runtime Behavior Gate` is completed by `W-0238`. The check rule is `runtime.friends_relationship_runtime_behavior_gate`.

`M-167 Friends Relationship Runtime Behavior Implementation` is completed by `W-0239`. The check rule is `runtime.friends_relationship_runtime_behavior_implementation`.

`M-168 Friends Relationship Protocol Route Gate` is completed by `W-0240`. The check rule is `runtime.friends_relationship_protocol_route_gate`.

`M-169 Friends Relationship Protocol Route Implementation` is completed by `W-0241`. The check rule is `runtime.friends_relationship_protocol_route_implementation`.

`M-170 Friends Relationship Protocol Route Local Proof` is completed by `W-0242`. The check rule is `runtime.friends_relationship_protocol_route_local_proof`.

The repository next work item is `W-0249 Implement Pitaya-aligned frontend/backend role source-first map`. That work is outside this module and must stay gate-only; do not add protocol shape changes, repository interface changes, PostgreSQL adapter changes, migrations, dependencies, event/audit tables, groups, parties, chat, matchmaking, match runtime, SDK publication, hosted deployments, distributed runtime behavior, frontend/backend server role implementation, operations/admin implementation, or direct Nakama/Pitaya API compatibility in this module until a later bounded work item authorizes that scope.

## When Not To Use This Module

Do not use this module for:

- New WebSocket, HTTP, Protobuf, or generated wire behavior beyond the completed W-0241 route family.
- PostgreSQL adapter implementation or SQL execution under this module.
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
- PostgreSQL adapter implementation: `runtime/internal/platform/persistence/postgres/friend_relationship_repository.go`
- PostgreSQL adapter tests: `runtime/internal/platform/persistence/postgres/friend_relationship_repository_test.go`
- Runtime behavior gate: `docs/friends-relationship-runtime-behavior-gate.md`
- Runtime behavior implementation: `runtime/internal/app/friends/service.go`
- Runtime behavior tests: `runtime/internal/app/friends/service_test.go`
- Protocol route gate: `docs/friends-relationship-protocol-route-gate.md`
- Protocol route implementation: `proto/vibit/friends/v1/friends.proto`, `runtime/internal/generated/proto/vibit/friends/v1/friends.pb.go`, `runtime/internal/platform/protocol/protobuf/friends_bridge.go`, `runtime/internal/app/bootstrap/friends.go`, and `runtime/internal/app/friends/routes.go`
- Protocol route local proof: `W-0242 Prove friends relationship protocol route in local alpha request flow`
- Next Pitaya-aligned vocabulary follow-up: `W-0249 Implement Pitaya-aligned frontend/backend role source-first map`

The first public runtime commands and queries are still deferred. Future runtime behavior must derive actor identity from validated request identity before calling this repository interface; client-supplied player ids are not authentication proof.

## Forbidden Shortcuts

- Do not bypass boundaries declared in `module.yaml`.
- Do not add unregistered public commands, queries, events, errors, or permissions.
- Do not add PostgreSQL adapter code under this module.
- Do not import `pgx`, `database/sql`, WebSocket packages, generated Protobuf packages, SDK packages, or distributed runtime packages into this module.
- Do not execute SQL or mention concrete SQL statements in friends module source.
- Do not change migrations from this module.
- Do not wire runtime handlers, startup composition, route policy, protocol adapters, or transport behavior from this module.
- Do not add new protocol routes, Protobuf sources, or generated output beyond the W-0241 friends route family without a later bounded work item.
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
