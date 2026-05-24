# storage Module Agent Guide 中文版

状态：Draft v0.1

说明：本文件是 `modules/storage/AGENTS.md` 的简体中文译本。英文版本是权威版本。

## 何时使用本模块

当需求涉及 player-owned small JSON game state 的 storage-object repository vocabulary 和 storage-neutral value types 时，使用本模块。

当前已实现的 slice 有意保持很窄：

- `runtime/internal/modules/storage.Repository`
- `StorageObject`、owner、identity、value、status、version，以及 create/read/list/update/delete input 和 result types
- optimistic conflict classes 和 redacted repository errors
- normalization helpers 和聚焦的 Go tests

`M-133 Storage Objects Repository Interface Implementation` 已由 `W-0205` 完成。检查规则是 `runtime.storage_objects_repository_interface_implementation`。

`M-134 Storage Objects PostgreSQL Adapter Gate` 已由 `W-0206` 完成。检查规则是 `runtime.storage_objects_postgresql_adapter_gate`。

`M-135 Storage Objects PostgreSQL Adapter Implementation` 已由 `W-0207` 完成。检查规则是 `runtime.storage_objects_postgresql_adapter_implementation`。

`M-136 Storage Objects Runtime Behavior Gate` 已由 `W-0208` 完成。检查规则是 `runtime.storage_objects_runtime_behavior_gate`。

`M-137 Storage Objects Runtime Behavior Implementation` 已由 `W-0209` 完成。检查规则是 `runtime.storage_objects_runtime_behavior_implementation`。Application service 位于 `runtime/internal/app/storage/service.go`，聚焦测试位于 `runtime/internal/app/storage/service_test.go`。

`M-138 Storage Objects Protocol Route Gate` 已由 `W-0210` 完成。检查规则是 `runtime.storage_objects_protocol_route_gate`。Gate 位于 `docs/storage-objects-protocol-route-gate.md`，记录未来 own-player storage routes、`vibit.storage.v1` message posture、protected-route identity handoff 和 Nakama/Pitaya reference mapping，但不添加 direct compatibility。

`M-139 Storage Objects Protocol Route Implementation` 已由 `W-0211` 完成。检查规则是 `runtime.storage_objects_protocol_route_implementation`。该实现添加了 `vibit.storage.v1` Protobuf source 和 generated output、route keys、protocol bridge mapping、bootstrap handlers、startup registration，以及针对 `storage.GetOwnStorageObject`、`storage.ListOwnStorageObjects`、`storage.PutOwnStorageObject` 和 `storage.DeleteOwnStorageObject` 的聚焦测试。

`M-140 Storage Objects Protocol Route Local Proof` 已由 `W-0212` 完成。检查规则是 `runtime.storage_objects_protocol_route_local_proof`。该 proof 添加了 `TestStorageObjectsProtocolRouteLocalAlphaFlow`，在 local authenticated E2E fixture 中加入 test-only storage route registration 和 test-only storage repository，并更新 `examples/local-alpha-request-loop.sh`，用于通过现有 WebSocket/Protobuf route flow 演示 authenticated own-player storage object put/get/list/delete。

Storage module 的 storage-object work 已完成到 local proof，第一版 realtime runtime slice 已在本模块外完成，realtime protocol/WebSocket outbound delivery gate 和 implementation 也已在本模块外完成，并且 outbound delivery 后的 next direction 已选择 AI-native feature request and test workflow。Repository 下一项 work item 是 `W-0220 Define agent-native feature request and test workflow`，不属于本模块。不要把 storage module 当作 realtime protocol、WebSocket outbound delivery 或 repository-wide feature workflow behavior 的 owner；除非后续 bounded storage work item 明确授权，继续保留 storage service、repository interface、PostgreSQL adapter、migration、public ACL、admin search、group/guild/party/room/match storage scope、batch write、JSON patch、merge、TTL、script hook、large object/blob storage、S3-compatible object storage、production memory storage behavior 和 direct Nakama/Pitaya API compatibility 的 deferrals。

## 何时不要使用本模块

不要把本模块用于：

- WebSocket、HTTP、Protobuf 或 generated wire behavior。
- PostgreSQL adapter implementation 或 SQL execution。
- Player account lifecycle。
- Authentication、token formats、credential storage 或 session validation。
- Inventory item grants 或 inventory capacity rules。
- Large object/blob storage。
- S3-compatible object storage。
- Direct Nakama 或 Pitaya public API compatibility。

如果需求需要这些概念，应创建或更新对应 owner boundary，不要在这里隐藏 ownership。

## Extension Points

- Repository interface: `runtime/internal/modules/storage.Repository`
- Repository value types: `StorageObject`, `StorageObjectOwner`, `StorageObjectIdentity`, `StorageObjectValue`, `StorageObjectVersion`
- Normalizers: create/read/list/update/delete inputs 和 returned records
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

第一批 public commands 和 queries 是 own-player storage object route family。Public ACLs、cross-owner access、admin search、group/guild/party/room/match scopes、batch operations、JSON patch/merge、TTL、script hooks 和 blob/S3 storage 继续 deferred。

## Forbidden Shortcuts

- 不要绕过 `module.yaml` 中声明的边界。
- 不要添加未登记的 public commands、queries、events、errors 或 permissions。
- 不要在本模块下添加 PostgreSQL adapter code。
- 不要把 `pgx`、`database/sql`、WebSocket packages、generated Protobuf packages、S3 SDKs 或 MinIO clients import 到本模块。
- 不要在 storage module source 中执行 SQL 或写具体 SQL statements。
- 不要从本模块修改 migrations。
- 不要从本模块 wire 新的 runtime handlers、startup composition、route policy 或 transport behavior。
- 没有后续 bounded work item 时，不要添加新的 protocol routes、Protobuf sources 或 generated output。
- 不要在 storage object value types 中存 raw credentials、raw tokens、verifier material、lookup digests、verifier digests、cookies、transport subprotocols、connection metadata、blob buckets、S3 keys 或 large object payloads。
- 不要把 `owner_id` 当成 authenticated proof。Player identity 和 session validation 仍由各自 boundary 拥有。

## Required Tests

参见 `module.yaml` 中的 `tests.required`。

当前 repository interface slice 的测试必须覆盖：

- Repository interface storage neutrality。
- Closed owner-kind 和 object-status vocabulary。
- Returned record normalization。
- Create/get/list/update/delete input normalization。
- Top-level JSON object value validation 和 byte copying。
- Positive version 和 expected-version constraints。
- Pagination bounds。
- Redacted conflict 和 repository errors。
- 不包含 secret、transport、blob、S3 和 direct compatibility fields。

当前 runtime behavior slice 的测试必须覆盖：

- Metadata-only identity 在 repository access 前被拒绝。
- Validated player identity owner derivation。
- Own-object get/list/put/delete behavior。
- Value shape、value size、list pagination 和 expected-version validation。
- Redacted conflict mapping。
- Unit-of-work storage repository handoff。

当前 protocol route slice 的测试必须覆盖：

- Own-player get/list/put/delete 的 Protobuf route payload mapping。
- Optional expected-version preservation。
- Response mapping 和 RFC3339Nano timestamp formatting。
- Redacted handler errors。
- Validated request identity handoff。
- Protected-route authentication wrapper enforcement。

修改 storage runtime source 后运行 `node tools/vibit check runtime`。Go 可用时，也运行 `cd runtime && go test ./...`。
