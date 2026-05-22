# Impact

This change adds a gate-only storage objects repository boundary:

- `docs/storage-objects-repository-boundary.md`
- `docs/storage-objects-repository-boundary.zh-CN.md`
- `decisions/ADR-0112-storage-objects-repository-boundary.md`
- `runtime.storage_objects_repository_boundary`

The boundary records:

- Future repository owner candidate `runtime/internal/modules/storage`.
- Future interface candidate `runtime/internal/modules/storage.Repository`.
- Future PostgreSQL adapter owner `runtime/internal/platform/persistence/postgres`.
- Logical table `storage_objects`.
- Candidate storage object value types.
- Candidate create/read/list/update/delete repository capabilities.
- Version and optimistic conflict handoff.
- Redaction expectations for values, owner identity, identifiers, errors, and forbidden material.
- PostgreSQL adapter expectations for later work.
- Stop conditions and the next repository interface implementation slice.

No Go source files, module directories, repository interfaces, PostgreSQL adapters, SQL execution behavior, migrations, runtime handlers, protocol routes, Protobuf sources, generated output, dependencies, hosted deployment, release artifacts, public announcements, paid promotion, large object/blob storage, S3-compatible object storage, or direct Nakama/Pitaya compatibility are added.
