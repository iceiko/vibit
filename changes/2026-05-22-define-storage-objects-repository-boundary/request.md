# Request

## Original Request

```text
继续推进
```

## Clarified Requirement

Advance `W-0204 Define storage objects repository boundary` by defining the storage-neutral storage objects repository boundary after the `storage_objects` migration source.

## User-Visible Outcome

Maintainers and agents can inspect a clear boundary standard before adding Go storage object repository interfaces, PostgreSQL storage adapters, protocol routes, or runtime behavior.

## Non-Goals

- Implementing repository interfaces.
- Creating `runtime/internal/modules/storage`.
- Adding PostgreSQL storage adapters.
- Adding SQL execution behavior.
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

- Exact Go repository type names remain deferred to `W-0205`.
- PostgreSQL adapter mapping remains deferred.
- Runtime permission and identity derivation remain deferred.
- Protocol route shape and generated output remain deferred.
- Admin, public, cross-owner, group/guild, room, party, match, object/blob, and S3-compatible storage scopes remain deferred.

## Acceptance Criteria

- [x] `docs/storage-objects-repository-boundary.md` exists.
- [x] `docs/storage-objects-repository-boundary.zh-CN.md` exists.
- [x] `ADR-0112` records the decision.
- [x] The boundary defines future repository ownership, candidate value types, create/read/list/update/delete vocabulary, version and conflict handoff, redaction posture, PostgreSQL adapter expectations, future queue, and stop conditions.
- [x] `runtime.storage_objects_repository_boundary` check coverage exists.
- [x] `W-0204` is completed and `W-0205` is next-ready.
- [x] No Go repository, adapter, runtime behavior, SQL execution, migration change, Protobuf, generated output, dependency, hosted deployment, release artifact, public announcement, or paid promotion is added by this slice.
