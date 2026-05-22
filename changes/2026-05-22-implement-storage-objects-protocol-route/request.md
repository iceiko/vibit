# Request

## Original Request

```text
继续推进。注意要贴合nakama pitaya
```

The maintainer later asked to continue with Nakama/Pitaya alignment and to commit and push the result.

## Clarified Requirement

Advance `W-0211 Implement storage objects protocol route` by exposing the already implemented own-player storage object application behavior through the bounded WebSocket/Protobuf route surface authorized by `ADR-0118`.

The slice must stay aligned with Nakama/Pitaya at the capability and layering level:

- Nakama alignment: durable player-owned storage objects with collection/key/value/version/list/read/write/delete semantics.
- Pitaya alignment: transport, session/identity, route, serializer/protocol adapter, handler, and backend behavior remain separated.
- Compatibility boundary: no direct Nakama or Pitaya public API compatibility is added.

## User-Visible Outcome

The PostgreSQL-backed runtime path can now register and dispatch the following protected storage object routes:

- `storage.GetOwnStorageObject`
- `storage.ListOwnStorageObjects`
- `storage.PutOwnStorageObject`
- `storage.DeleteOwnStorageObject`

The route family uses `vibit.storage.v1` Protobuf payloads and the existing authenticated request wrapper. Client payloads do not carry owner ids, player ids, session ids, access tokens, digest material, blob handles, or S3 fields.

## Non-Goals

- Changing storage object repository interfaces.
- Changing the PostgreSQL adapter.
- Changing migrations.
- Adding dependencies.
- Changing authentication/session behavior or route-protection semantics.
- Adding WebSocket handshake authentication.
- Adding public ACLs or cross-owner access.
- Adding admin search.
- Adding group/guild/party/room/match storage scopes.
- Adding batch writes.
- Adding JSON patch or merge semantics.
- Adding TTL or script hooks.
- Adding broad operations/admin behavior.
- Adding hosted deployments.
- Creating release binaries, packages, containers, checksums, provenance files, signing artifacts, install scripts, registry publications, or SDK packages.
- Executing public announcements beyond the GitHub release record.
- Running paid promotion.
- Adding large object/blob storage or S3-compatible object storage.
- Adding direct Nakama/Pitaya API compatibility.

## Acceptance Criteria

- [x] `proto/vibit/storage/v1/storage.proto` declares the own-player get/list/put/delete payload family.
- [x] `runtime/internal/generated/proto/vibit/storage/v1/storage.pb.go` is regenerated from the Protobuf source.
- [x] `runtime/internal/app/storage/routes.go` declares route keys for the storage object family.
- [x] `runtime/internal/platform/protocol/protobuf/storage_bridge.go` maps wire payloads to application storage requests and maps application results back to Protobuf responses.
- [x] `runtime/internal/app/bootstrap/storage.go` registers application route handlers and injects validated request identity into service requests.
- [x] PostgreSQL startup composition registers storage routes and bypasses the outer transaction wrapper for service-owned write routes.
- [x] Focused tests cover own-player get/list/put/delete mapping, optional expected versions, response mapping, redacted handler errors, route protection, and startup object-id shape.
- [x] `ADR-0119` records the implementation decision.
- [x] `runtime.storage_objects_protocol_route_implementation` check coverage exists.
- [x] `W-0211` is completed and `W-0212` is next-ready.
