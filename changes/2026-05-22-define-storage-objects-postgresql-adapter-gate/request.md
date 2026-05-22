# Request

Define the storage objects PostgreSQL adapter gate after the storage-neutral repository interface.

## Source

The maintainer requested: "使用superpowers继续推进。"

The current next-ready work item was `W-0206 Define storage objects PostgreSQL adapter gate`.

## Acceptance Criteria

- Define the storage objects PostgreSQL adapter gate.
- Record future PostgreSQL adapter owner, constructor and unit-of-work expectations, SQL mapping posture, transaction handoff, error mapping, test expectations, and stop conditions.
- Register the check rule `runtime.storage_objects_postgresql_adapter_gate`.
- Mark `M-134/W-0206` completed and select `M-135/W-0207 Implement storage objects PostgreSQL adapter` as next-ready.
- Preserve PostgreSQL adapter implementation, SQL execution, runtime behavior, protocol routes, Protobuf/generated output, dependencies, migration changes, operations/admin behavior, authentication/session changes, hosted deployment, release artifacts, public announcements, paid promotion, broad product module expansion, large object/blob storage, S3-compatible object storage, and direct Nakama/Pitaya API compatibility deferrals.
