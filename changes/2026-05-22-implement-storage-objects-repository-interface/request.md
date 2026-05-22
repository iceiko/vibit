# Request

## Original Request

```text
使用superpowers继续推进。
```

## Clarified Requirement

Advance `W-0205 Implement storage-neutral storage objects repository interface` by implementing only the storage-neutral repository interface, value vocabulary, normalizers, redacted errors, module manifest, and focused tests under the authorized storage module owner.

## User-Visible Outcome

Maintainers and agents can now inspect and build against a concrete Go repository interface for player-owned small JSON storage objects before PostgreSQL adapter, runtime behavior, or protocol work begins.

## Non-Goals

- Adding PostgreSQL storage adapters.
- Adding SQL execution behavior.
- Adding unit-of-work factory wiring.
- Changing migrations.
- Adding storage objects runtime behavior.
- Adding protocol routes.
- Adding Protobuf source files.
- Adding generated output.
- Adding dependencies.
- Adding automatic startup migration behavior.
- Changing authentication/session behavior.
- Adding broad operations/admin behavior.
- Adding hosted deployments.
- Creating release binaries, packages, containers, checksums, provenance files, signing artifacts, install scripts, registry publications, or SDK packages.
- Executing public announcements beyond the GitHub release record.
- Running paid promotion.
- Adding large object/blob storage or S3-compatible object storage.
- Adding direct Nakama/Pitaya API compatibility.

## Unknowns

- PostgreSQL adapter SQL mapping remains deferred to `W-0206` and later implementation work.
- Runtime permission and identity derivation remain deferred.
- Protocol route shape and generated output remain deferred.
- Admin, public, cross-owner, group/guild, room, party, match, object/blob, and S3-compatible storage scopes remain deferred.

## Acceptance Criteria

- [x] `runtime/internal/modules/storage/repository.go` defines storage-neutral storage object repository value types and `Repository`.
- [x] `runtime/internal/modules/storage/repository_test.go` verifies neutrality, closed vocabulary, normalization, copying, bounds, redaction, and forbidden material.
- [x] `modules/storage/module.yaml` exists and records module ownership.
- [x] `modules/storage/AGENTS.md` and `modules/storage/AGENTS.zh-CN.md` exist.
- [x] `ADR-0113` records the decision.
- [x] `runtime.storage_objects_repository_interface_implementation` check coverage exists.
- [x] `W-0205` is completed and `W-0206` is next-ready.
- [x] No adapter, SQL execution, runtime behavior, migration change, Protobuf, generated output, dependency, hosted deployment, release artifact, public announcement, paid promotion, blob/S3 storage, or direct compatibility is added by this slice.
