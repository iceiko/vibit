# Impact Analysis

## Nakama Product Capability Impact

This change supports the Nakama `friends_groups_and_parties` capability family by defining the future durable schema posture for player friendship relationships before implementation.

It adopts the product need for durable social graph state and adapts it into vibit's contract-first, agent-native workflow. It does not copy a Nakama API, route, schema, or compatibility surface.

## Pitaya Impact

Pitaya remains deferred. This change does not add distributed topology, frontend/backend split, RPC, service discovery, groups, cluster routing, distributed sessions, or distributed social graph routing.

## Affected Modules

- `runtime`: records the gate, next work item, and check rule.
- `reference`: records Nakama-first capability alignment and Pitaya deferral.
- `repository_workflow`: records the next migration-source-only continuation.
- `friends`: future module ownership candidate only; no module code is added.

## Module Ownership Impact

The future repository owner candidate is `runtime/internal/modules/friends`. The future PostgreSQL adapter owner is `runtime/internal/platform/persistence/postgres`.

No module directory, repository interface, adapter, runtime handler, protocol handler, or generated source is added by this gate.

## Public Contract Impact

No command, query, event, error, permission, route, protocol payload, generated output, or public API changes are added.

The future command/query/event/error vocabulary remains the semantic vocabulary recorded by `ADR-0139`.

## Data And Migration Impact

No migration source is added.

The gate records the future table candidate `friend_relationships`, canonical pair identity (`player_low_id + player_high_id`), lifecycle state, actor-specific block timestamps, version and timestamp posture, uniqueness/index expectations, event/audit deferral, redaction posture, and future migration source candidate `runtime/migrations/postgres/000007_create_friend_relationships.sql`.

## Test Impact

Behavior tests are not applicable because no runtime behavior changes.

Future migration-source checks should cover table name, columns, pair ordering, self-target prevention, lifecycle state vocabulary, block columns, relationship version, timestamps, pair uniqueness, list-query indexes, forbidden columns, and no runtime/protocol/generated/repository/adapter scope.

Repository checks added or updated for this gate include `runtime.friends_relationship_persistence_schema_gate`.

## Documentation And Memory Impact

This change adds:

- `docs/friends-relationship-persistence-schema-gate.md`
- `docs/friends-relationship-persistence-schema-gate.zh-CN.md`
- `decisions/ADR-0140-friends-relationship-persistence-schema-gate.md`
- `conversations/2026-05-24-friends-relationship-persistence-schema-gate.md`
- filled change artifacts under `changes/2026-05-24-define-friends-relationship-persistence-schema-gate/`

It also updates architecture manifests, README/AGENTS continuation pointers, `tools/vibit`, and `rules/check-rules.json`.

## Compatibility Risks

No API, event, data, protocol, SDK, hosted, distributed runtime, or direct Nakama/Pitaya compatibility surface changes are added.

The main risk is premature schema overreach. This is mitigated by keeping event/audit tables, runtime behavior, repository interfaces, PostgreSQL adapters, protocol routes, generated output, and direct compatibility deferred.
