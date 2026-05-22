# Request

## Original Request

```text
继续推进。
```

## Clarified Requirement

Advance the next-ready work item `W-0209 Implement storage objects runtime behavior` by adding application-owned storage object runtime behavior under the boundary ratified by `ADR-0116`, without adding protocol routes, Protobuf source, generated output, startup wiring, repository interface changes, PostgreSQL adapter changes, migrations, dependencies, authentication/session behavior changes, hosted deployments, release artifacts, public announcements, paid promotion, blob/S3 storage, or direct Nakama/Pitaya API compatibility.

## User-Visible Outcome

Maintainers and agents can use `runtime/internal/app/storage.Service` as the first application behavior surface for player-owned small JSON storage objects. The service supports own-object get, list, put, and delete behavior, derives ownership from validated player request identity, rejects metadata-only identity before repository access, maps conflicts to redacted public classes, and uses unit-of-work storage repository handoff.

## Non-Goals

- Add WebSocket, HTTP, or other protocol routes.
- Add Protobuf source files or generated output.
- Register runtime handlers or startup composition.
- Change the storage repository interface.
- Change the PostgreSQL adapter.
- Change migrations or automatic startup migration behavior.
- Add dependencies.
- Change authentication/session validation behavior or route-protection semantics.
- Add public ACLs, admin search, group/guild/party/room/match storage scopes, batch writes, JSON patch, merge semantics, TTL, or script hooks.
- Add large object/blob storage or S3-compatible object storage.
- Add hosted deployments, release artifacts, public announcements, paid promotion, or direct Nakama/Pitaya API compatibility.

## Unknowns

- Exact protocol route names and Protobuf request/response shapes remain deferred.
- Startup wiring and route registration remain deferred.
- Future permission names remain deferred until public contracts are ratified.

## Acceptance Criteria

- [x] Application-owned service source exists at `runtime/internal/app/storage/service.go`.
- [x] Focused service tests exist at `runtime/internal/app/storage/service_test.go`.
- [x] Metadata-only identity is rejected before unit-of-work or repository access.
- [x] Validated player identity is used to derive `storage.StorageObjectOwner`.
- [x] Get, list, put, and delete service methods call the storage repository through unit-of-work handoff.
- [x] Input validation and value-size bounds happen before repository mutation when possible.
- [x] Repository conflicts map to stable redacted public error codes.
- [x] Protocol, generated output, repository interface, PostgreSQL adapter, dependency, migration, authentication/session, hosted, release, blob/S3, and direct compatibility deferrals are preserved.
