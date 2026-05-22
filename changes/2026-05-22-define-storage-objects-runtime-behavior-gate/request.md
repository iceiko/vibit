# Request

Define the storage objects runtime behavior gate after the PostgreSQL adapter.

## Source

The maintainer requested: "继续推进。"

The current next-ready work item was `W-0208 Define storage objects runtime behavior gate`.

## Acceptance Criteria

- Define the storage objects runtime behavior gate.
- Record owner identity derivation, request identity requirements, permissions, validation, conflict semantics, route-policy expectations, service/application ownership, unit-of-work handoff, and stop conditions.
- Explicitly refuse to treat metadata-only `player_id` or `session_id` as authenticated proof.
- Register the check rule `runtime.storage_objects_runtime_behavior_gate`.
- Mark `M-136/W-0208` completed and select `M-137/W-0209 Implement storage objects runtime behavior` as next-ready.
- Preserve runtime behavior implementation, runtime handlers, protocol routes, Protobuf/generated output, repository interface changes, PostgreSQL adapter changes, dependencies, migration changes, operations/admin behavior, authentication/session changes, hosted deployment, release artifacts, public announcements, paid promotion, broad product module expansion, large object/blob storage, S3-compatible object storage, and direct Nakama/Pitaya API compatibility deferrals.
