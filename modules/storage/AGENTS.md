# storage Module Agent Guide

Status: Draft v0.1

## When To Use This Module

Use this module for storage-object repository vocabulary and storage-neutral value types for player-owned small JSON game state.

The current implemented slice is intentionally narrow:

- `runtime/internal/modules/storage.Repository`
- `StorageObject`, owner, identity, value, status, version, create/read/list/update/delete input and result types
- optimistic conflict classes and redacted repository errors
- normalization helpers and focused Go tests

`M-133 Storage Objects Repository Interface Implementation` is completed by `W-0205`. The check rule is `runtime.storage_objects_repository_interface_implementation`.

`M-134 Storage Objects PostgreSQL Adapter Gate` is completed by `W-0206`. The check rule is `runtime.storage_objects_postgresql_adapter_gate`.

`M-135 Storage Objects PostgreSQL Adapter Implementation` is completed by `W-0207`. The check rule is `runtime.storage_objects_postgresql_adapter_implementation`.

`M-136 Storage Objects Runtime Behavior Gate` is completed by `W-0208`. The check rule is `runtime.storage_objects_runtime_behavior_gate`.

`M-137 Storage Objects Runtime Behavior Implementation` is completed by `W-0209`. The check rule is `runtime.storage_objects_runtime_behavior_implementation`. The application service lives at `runtime/internal/app/storage/service.go` with focused tests at `runtime/internal/app/storage/service_test.go`.

`M-138 Storage Objects Protocol Route Gate` is completed by `W-0210`. The check rule is `runtime.storage_objects_protocol_route_gate`. The gate lives at `docs/storage-objects-protocol-route-gate.md` and records future own-player storage routes, `vibit.storage.v1` message posture, protected-route identity handoff, and Nakama/Pitaya reference mapping without direct compatibility.

`M-139 Storage Objects Protocol Route Implementation` is completed by `W-0211`. The check rule is `runtime.storage_objects_protocol_route_implementation`. The implementation adds `vibit.storage.v1` Protobuf source and generated output, route keys, protocol bridge mapping, bootstrap handlers, startup registration, and focused tests for `storage.GetOwnStorageObject`, `storage.ListOwnStorageObjects`, `storage.PutOwnStorageObject`, and `storage.DeleteOwnStorageObject`.

`M-140 Storage Objects Protocol Route Local Proof` is completed by `W-0212`. The check rule is `runtime.storage_objects_protocol_route_local_proof`. The proof adds `TestStorageObjectsProtocolRouteLocalAlphaFlow`, extends the local authenticated E2E fixture with test-only storage route registration and a test-only storage repository, and updates `examples/local-alpha-request-loop.sh` to demonstrate authenticated own-player storage object put/get/list/delete through the existing WebSocket/Protobuf route flow.

Storage module storage-object work is complete through local proof, the first realtime runtime slice is implemented outside this module, the realtime protocol/WebSocket outbound delivery gate and implementation are completed outside this module, the agent-native feature request and test workflow is defined outside this module, the Nakama-aligned workflow pilot selected presence/status work outside this module, the presence/status local proof hardening slice is completed outside this module, authenticated gameplay failure-path verification is completed outside this module, the next Nakama prototype-ready capability selection chose developer experience outside this module, the local alpha example client path gate and implementation are completed outside this module, the feature request scaffolding gate and implementation are completed outside this module, the scaffolded Nakama feature request intake pilot selected friends relationship lifecycle outside this module, and the friends relationship lifecycle gate, persistence schema gate, migration source, repository boundary, repository interface implementation, and PostgreSQL adapter gate are completed outside this module. The repository next work item is `W-0237 Implement friends relationship PostgreSQL adapter`, outside this module. Do not use the storage module as owner for realtime protocol, WebSocket outbound delivery, repository-wide feature workflow behavior, presence/status proof hardening, authenticated gameplay failure-path verification, next Nakama capability selection, local alpha example client path gating, example path implementation, feature request scaffolding, scaffolded Nakama intake, friends relationship lifecycle, friends relationship persistence, friends relationship migration source, friends relationship repository boundary, friends relationship repository interface, friends relationship PostgreSQL adapter gate, or friends relationship PostgreSQL adapter implementation; preserve storage service, repository interface, PostgreSQL adapter, migration, public ACL, admin search, group/guild/party/room/match storage scope, batch write, JSON patch, merge, TTL, script hook, large object/blob storage, S3-compatible object storage, production memory storage behavior, and direct Nakama/Pitaya API compatibility deferrals unless a later bounded storage work item authorizes them.

## When Not To Use This Module

Do not use this module for:

- WebSocket, HTTP, Protobuf, or generated wire behavior.
- PostgreSQL adapter implementation or SQL execution.
- Player account lifecycle.
- Authentication, token formats, credential storage, or session validation.
- Inventory item grants or inventory capacity rules.
- Large object/blob storage.
- S3-compatible object storage.
- Direct Nakama or Pitaya public API compatibility.

If a requirement needs one of those concepts, create or update the owning boundary instead of adding hidden ownership here.

## Extension Points

- Repository interface: `runtime/internal/modules/storage.Repository`
- Repository value types: `StorageObject`, `StorageObjectOwner`, `StorageObjectIdentity`, `StorageObjectValue`, `StorageObjectVersion`
- Normalizers: create/read/list/update/delete inputs and returned records
- Tests: `runtime/internal/modules/storage/repository_test.go`
- PostgreSQL adapter owner: `runtime/internal/platform/persistence/postgres`
- PostgreSQL adapter source: `runtime/internal/platform/persistence/postgres/storage_object_repository.go`
- PostgreSQL adapter tests: `runtime/internal/platform/persistence/postgres/storage_object_repository_test.go`
- Runtime behavior owner candidate: `runtime/internal/app/storage`
- Runtime behavior service: `runtime/internal/app/storage/service.go`
- Runtime behavior tests: `runtime/internal/app/storage/service_test.go`
- Runtime behavior gate: `docs/storage-objects-runtime-behavior-gate.md`
- Protocol route gate: `docs/storage-objects-protocol-route-gate.md`
- Protocol route Protobuf source: `proto/vibit/storage/v1/storage.proto`
- Generated protocol output: `runtime/internal/generated/proto/vibit/storage/v1/storage.pb.go`
- Protocol bridge: `runtime/internal/platform/protocol/protobuf/storage_bridge.go`
- Application route handlers: `runtime/internal/app/bootstrap/storage.go`

The first public commands and queries are the own-player storage object route family. Public ACLs, cross-owner access, admin search, group/guild/party/room/match scopes, batch operations, JSON patch/merge, TTL, script hooks, and blob/S3 storage remain deferred.

## Forbidden Shortcuts

- Do not bypass boundaries declared in `module.yaml`.
- Do not add unregistered public commands, queries, events, errors, or permissions.
- Do not add PostgreSQL adapter code under this module.
- Do not import `pgx`, `database/sql`, WebSocket packages, generated Protobuf packages, S3 SDKs, or MinIO clients into this module.
- Do not execute SQL or mention concrete SQL statements in storage module source.
- Do not change migrations from this module.
- Do not wire new runtime handlers, startup composition, route policy, or transport behavior from this module.
- Do not add new protocol routes, Protobuf sources, or generated output without a later bounded work item.
- Do not store raw credentials, raw tokens, verifier material, lookup digests, verifier digests, cookies, transport subprotocols, connection metadata, blob buckets, S3 keys, or large object payloads in storage object value types.
- Do not treat `owner_id` as authenticated proof. Player identity and session validation remain owned by their own boundaries.

## Required Tests

See `tests.required` in `module.yaml`.

For the current repository interface slice, tests must cover:

- Repository interface storage neutrality.
- Closed owner-kind and object-status vocabulary.
- Returned record normalization.
- Create/get/list/update/delete input normalization.
- Top-level JSON object value validation and byte copying.
- Positive version and expected-version constraints.
- Pagination bounds.
- Redacted conflict and repository errors.
- Absence of secret, transport, blob, S3, and direct compatibility fields.

For the current runtime behavior slice, tests must cover:

- Metadata-only identity rejection before repository access.
- Validated player identity owner derivation.
- Own-object get/list/put/delete behavior.
- Value shape, value size, list pagination, and expected-version validation.
- Redacted conflict mapping.
- Unit-of-work storage repository handoff.

For the current protocol route slice, tests must cover:

- Protobuf route payload mapping for own-player get/list/put/delete.
- Optional expected-version preservation.
- Response mapping and RFC3339Nano timestamp formatting.
- Redacted handler errors.
- Validated request identity handoff.
- Protected-route authentication wrapper enforcement.

Run `node tools/vibit check runtime` after changing storage runtime source. When Go is available, also run `cd runtime && go test ./...`.
