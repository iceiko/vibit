# Impact

This change adds the storage-neutral storage objects repository interface implementation:

- `runtime/internal/modules/storage/repository.go`
- `runtime/internal/modules/storage/repository_test.go`
- `modules/storage/module.yaml`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`
- `decisions/ADR-0113-storage-objects-repository-interface-implementation.md`
- `runtime.storage_objects_repository_interface_implementation`

The interface records:

- Repository owner `runtime/internal/modules/storage`.
- Interface `runtime/internal/modules/storage.Repository`.
- Module-owned storage object value types.
- Player owner first posture.
- Active/deleted object status vocabulary.
- JSON object value normalization.
- Create/read/list/update/delete repository capabilities.
- Version and optimistic conflict handoff.
- Redacted repository errors.
- Future PostgreSQL adapter owner `runtime/internal/platform/persistence/postgres`.

No PostgreSQL adapter code, SQL execution behavior, migration change, runtime handler, protocol route, Protobuf source, generated output, dependency, hosted deployment, release artifact, public announcement, paid promotion, large object/blob storage, S3-compatible object storage, or direct Nakama/Pitaya compatibility is added.
