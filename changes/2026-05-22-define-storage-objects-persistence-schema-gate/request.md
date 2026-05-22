# Request

## Original Request

The maintainer asked to continue advancing vibit after the storage objects behavior gate was completed.

## Clarified Requirement

Complete `W-0202 Define storage objects persistence schema gate` by defining the first storage objects persistence schema posture before adding migration source, repository interfaces, storage adapters, protocol, or runtime behavior.

## User-Visible Outcome

Maintainers, contributors, and agents can now see:

- the storage objects persistence schema gate standard;
- the future `storage_objects` table candidate;
- owner identity representation;
- collection/key constraints;
- JSONB value representation;
- BIGINT version representation;
- timestamp and soft-delete posture;
- uniqueness/index posture;
- redaction posture;
- future migration source candidate;
- future repository and adapter ownership candidates;
- that `W-0203 Add storage objects migration source` is the next work item.

## Non-Goals

- Do not add SQL migration source.
- Do not create the `storage_objects` table.
- Do not implement storage objects runtime behavior.
- Do not add protocol routes.
- Do not add Protobuf source files or generated output.
- Do not add repository interfaces.
- Do not add storage adapters.
- Do not add dependencies.
- Do not change authentication/session behavior.
- Do not change route protection behavior.
- Do not broaden operations/admin behavior.
- Do not add hosted deployments or demos.
- Do not create release binaries, packages, containers, checksums, provenance files, signing artifacts, install scripts, registry publications, SDK packages, or additional release artifacts.
- Do not execute public announcements beyond the GitHub release record.
- Do not run paid promotion.
- Do not add large object/blob storage.
- Do not add S3-compatible object storage.
- Do not add broad product module implementation.
- Do not add direct Nakama/Pitaya API compatibility.

## Unknowns

- Exact SQL DDL remains deferred to `W-0203`.
- Repository interface and PostgreSQL adapter shapes remain deferred.
- Protocol and Protobuf shapes remain deferred.
- Runtime create/read/update/delete/conflict behavior remains deferred.

## Acceptance Criteria

- [x] The storage objects persistence schema gate is recorded in English and Simplified Chinese.
- [x] The gate records table candidate, owner identity representation, collection/key constraints, value representation, version representation, timestamps, uniqueness/index posture, redaction posture, migration source candidate, and future repository/adapter boundaries.
- [x] `ADR-0110` records the decision.
- [x] Repository checks cover the gate state.
- [x] The next work item advances to `W-0203 Add storage objects migration source`.
- [x] Runtime, protocol, generated output, migration source, dependency, repository, adapter, operations/admin, authentication/session, hosted deployment, release artifact, public announcement, paid promotion, broad product module, object/blob storage, S3-compatible object storage, and direct compatibility deferrals remain preserved.
