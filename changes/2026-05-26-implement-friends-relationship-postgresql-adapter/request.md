# Request

Implement the friends relationship PostgreSQL adapter after the adapter gate.

## Source

The maintainer requested: "继续"

The current next-ready work item was `W-0237 Implement friends relationship PostgreSQL adapter`.

## Acceptance Criteria

- Implement the friends relationship PostgreSQL adapter under `runtime/internal/platform/persistence/postgres`.
- Map `runtime/internal/modules/friends.Repository` to the existing `friend_relationships` table.
- Cover create/update request, pair lookup, player-scoped list, accept, reject, remove, block, and unblock behavior with focused fake-executor tests.
- Preserve canonical pair handling, relationship version increments, actor-specific block semantics, redacted error behavior, conflict mapping, and absence of transaction-control SQL.
- Add unit-of-work repository handoff without wiring runtime startup, handlers, routes, or protocol behavior.
- Register the check rule `runtime.friends_relationship_postgresql_adapter_implementation`.
- Mark `M-165/W-0237` completed and select `M-166/W-0238 Define friends relationship runtime behavior gate` as next-ready.
- Preserve runtime friendship behavior, protocol route, Protobuf/generated output, dependency, migration, event/audit table, operations/admin, authentication/session, hosted deployment, release artifact, public announcement, paid promotion, broad social feature, distributed runtime, and direct Nakama/Pitaya API compatibility deferrals.
