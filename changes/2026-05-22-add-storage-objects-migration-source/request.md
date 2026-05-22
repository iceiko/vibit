# Request

## Original Request

```text
继续推进
```

## Clarified Requirement

Advance `W-0203 Add storage objects migration source` by adding the first PostgreSQL migration source for `storage_objects`, following the accepted storage objects persistence schema gate.

## User-Visible Outcome

Maintainers and agents can inspect and apply a concrete SQL migration source for player-owned small JSON storage objects.

## Non-Goals

- Adding storage objects runtime behavior.
- Adding repository interfaces.
- Adding PostgreSQL storage adapters.
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

- Storage-neutral repository shape remains deferred.
- PostgreSQL adapter mapping remains deferred.
- Protocol route shape remains deferred.
- Runtime permission and optimistic conflict behavior remain deferred.
- Object/blob storage and S3-compatible storage remain separate future capability families.

## Acceptance Criteria

- [x] `runtime/migrations/postgres/000006_create_storage_objects.sql` exists.
- [x] The migration declares goose Up and Down markers.
- [x] The migration declares `-- Module: storage`.
- [x] The migration creates `storage_objects`.
- [x] Required owner, collection, key, value, version, timestamp, and deletion columns are present.
- [x] Active logical identity uniqueness is declared.
- [x] No raw secret, digest, transport metadata, blob, file path, or S3 object columns are added.
- [x] No Go repository, adapter, runtime behavior, Protobuf, generated output, dependency, hosted deployment, release artifact, public announcement, or paid promotion is added by this slice.
