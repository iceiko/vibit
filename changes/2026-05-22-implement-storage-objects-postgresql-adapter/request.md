# Request

Implement the storage objects PostgreSQL adapter after the adapter gate.

## Source

The maintainer requested: "继续推进。"

The current next-ready work item was `W-0207 Implement storage objects PostgreSQL adapter`.

## Acceptance Criteria

- Implement the storage objects PostgreSQL adapter under `runtime/internal/platform/persistence/postgres`.
- Map `runtime/internal/modules/storage.Repository` to the existing `storage_objects` table.
- Cover create, get, list, update, and delete behavior with focused fake-executor tests.
- Preserve redacted error behavior, version-aware conflicts, active-row filtering, and soft-delete behavior.
- Add unit-of-work repository handoff without wiring runtime startup, handlers, routes, or protocol behavior.
- Register the check rule `runtime.storage_objects_postgresql_adapter_implementation`.
- Mark `M-135/W-0207` completed and select `M-136/W-0208 Define storage objects runtime behavior gate` as next-ready.
- Preserve runtime handler, protocol route, Protobuf/generated output, dependency, migration, operations/admin, authentication/session, hosted deployment, release artifact, public announcement, paid promotion, broad product module, large object/blob storage, S3-compatible object storage, and direct Nakama/Pitaya API compatibility deferrals.
