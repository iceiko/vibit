# Impact Analysis

## Affected Modules

- `runtime`
- `storage`
- `workflow`
- `reference`

## Module Ownership Impact

The storage module continues to own storage-neutral value types and repository vocabulary through `runtime/internal/modules/storage`.

The storage application service remains under `runtime/internal/app/storage` and owns runtime behavior. This change adds only the protocol-facing route key, bootstrap handler, and Protobuf bridge needed to expose that service.

Protocol ownership remains under `runtime/internal/platform/protocol/protobuf`, generated Protobuf output remains under `runtime/internal/generated/proto`, and process startup wiring remains under `runtime/cmd/vibit-server`.

## Public Contract Impact

Adds a client-visible Protobuf route family:

- query `storage.GetOwnStorageObject`
- query `storage.ListOwnStorageObjects`
- command `storage.PutOwnStorageObject`
- command `storage.DeleteOwnStorageObject`

The payload package is `vibit.storage.v1`.

The payloads intentionally omit owner, player, session, token, digest, blob, and S3 fields. Owner derivation remains server-side through validated request identity.

## Generated Output Impact

`buf generate` produces:

- `runtime/internal/generated/proto/vibit/storage/v1/storage.pb.go`

The generated file contains the `protoc-gen-go` generated-code marker and traces back to `vibit/storage/v1/storage.proto`.

## Data And Migration Impact

No migrations are added or changed.

The storage object PostgreSQL adapter and existing `storage_objects` table are reused unchanged through the application service and unit-of-work boundary.

## Runtime Impact

The PostgreSQL runtime path now constructs the storage application service, registers the storage route handlers, and bypasses the outer `TransactionalDispatcher` for `PutOwnStorageObject` and `DeleteOwnStorageObject` because the storage service already runs its own unit of work.

The memory runtime path remains unchanged and does not gain a storage object backing store.

## Authentication And Session Impact

The routes use the existing protected-route request-token policy through the authenticated request wrapper. This change does not add handshake authentication, new token carriers, session persistence changes, bound-session policy changes, or metadata-only identity trust.

## Nakama And Pitaya Alignment

Nakama informs the capability: durable player-owned storage objects with collection/key/value/version/list/read/write/delete semantics.

Pitaya informs the layering: route/session context, serializer/protocol adapter, handler, and backend behavior remain separate.

The implementation keeps vibit's route names, payload shapes, identity handoff, generated-output standard, and bounded workflow. It does not copy Nakama or Pitaya public APIs.

## Test Impact

Adds or updates tests for:

- route registration and identity handoff for all four storage routes;
- malformed payload and service-error redaction;
- Protobuf request/response bridge mapping;
- optional expected-version preservation;
- forbidden proto fields;
- protected-route authenticated-wrapper enforcement;
- valid authenticated wrapper dispatch to a storage route;
- random storage object id shape.

## Compatibility Risks

The main compatibility risk is introducing a first storage object wire shape before semantic contract source files exist. The risk is bounded by the route gate, ADR, Protobuf source trace, focused tests, and check-rule coverage. A future W-0212 local proof remains responsible for exercising the route family through the local alpha request flow.
