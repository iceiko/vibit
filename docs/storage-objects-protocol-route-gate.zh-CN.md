# Storage Objects Protocol Route Gate 中文版

状态：Accepted v0.1
最后更新：2026-05-22
范围：在 application runtime behavior 之后，为未来 client-facing storage objects protocol routes 定义 gate-only boundary
依赖：`docs/storage-objects-behavior-gate.md`、`docs/storage-objects-runtime-behavior-gate.md`、`decisions/ADR-0117-storage-objects-runtime-behavior-implementation.md`、`docs/runtime-protocol-adapter.md`、`docs/game-protocol.md`、`docs/generated-output.md`、`docs/reference-game-server-alignment.md`、`docs/nakama-pitaya-product-parity-roadmap.md`
Canonical decision：`ADR-0118`
说明：本文件是 `docs/storage-objects-protocol-route-gate.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本文定义 storage objects protocol route gate。它是一个 gate artifact。本文件不添加 protocol route implementation、Protobuf source、generated output、startup wiring、runtime handlers、repository interface changes、PostgreSQL adapter changes、migration changes、dependencies、authentication/session behavior changes、hosted deployments、release artifacts、public announcements、paid promotion、broad product module expansion、large object/blob storage、S3-compatible object storage 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Storage objects protocol route gate 记录如下：

```yaml
storage_objects_protocol_route_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0210
decision: ADR-0118
check_rule: runtime.storage_objects_protocol_route_gate
source_runtime_behavior_implementation_decision: ADR-0117
source_runtime_behavior_implementation: runtime/internal/app/storage/service.go
source_runtime_behavior_tests: runtime/internal/app/storage/service_test.go
source_runtime_behavior_gate_decision: ADR-0116
source_repository_interface_decision: ADR-0113
repository_interface: runtime/internal/modules/storage.Repository
future_protocol_source_candidate: proto/vibit/storage/v1/storage.proto
future_generated_go_output_candidate: runtime/internal/generated/proto/vibit/storage/v1/storage.pb.go
future_protocol_bridge_candidate: runtime/internal/platform/protocol/protobuf/storage_bridge.go
future_protocol_bridge_test_candidate: runtime/internal/platform/protocol/protobuf/storage_bridge_test.go
future_application_handler_candidate: runtime/internal/app/bootstrap/storage.go
future_application_handler_test_candidate: runtime/internal/app/bootstrap/storage_test.go
route_policy_requirement: request_token_required
authenticated_wrapper_required: true
request_identity_source: validated_authenticated_request_identity
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
client_supplied_owner_id_allowed: false
first_owner_kind: player
first_payload_package: vibit.storage.v1
protobuf_envelope_change_status: unchanged
websocket_transport_credential_neutral: true
protocol_route_gate_only: true
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
runtime_handler_added: false
startup_wiring_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
dependency_added: false
migration_added: false
authentication_session_behavior_changed: false
large_object_blob_storage_added: false
s3_compatible_object_storage_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_protocol_route_implementation_work_item: W-0211
future_protocol_route_implementation_direction: storage_objects_protocol_route_implementation
```

## 2. 目的

`W-0209` 已在 `runtime/internal/app/storage` 实现 application-owned storage object behavior。下一步有价值的边界不是直接写 route code 或生成 `.proto`，而是先记录未来 WebSocket/Protobuf 暴露方式如何调用该 service，并避免把 storage behavior 放进 transport、generated files 或 persistence adapters。

Nakama 提供产品面参考：client-facing storage objects 是常见 game backend 能力，包含 collection/key identity、owner scope、permissions、version conflict behavior、list/read/write/delete operations，以及 durable JSON game state。vibit 应覆盖这个能力类别。

Pitaya 提供架构面参考：acceptors、sessions、route handlers、serializers 和 backend services 应保持分离。vibit 通过保持 WebSocket transport credential-neutral、显式 Protobuf payload bridge、以及 application-owned route handlers 调用 application-owned storage services 来吸收这个经验。

本 gate 在实现前记录：

- candidate route names；
- candidate request/response message shapes；
- route protection 和 identity handoff posture；
- protocol adapter、application handler 和 startup ownership；
- generated-output expectations；
- error mapping 和 redaction expectations；
- Nakama/Pitaya reference mapping；
- 阻止本 slice 越界到 implementation 或 generated artifacts 的 stop conditions。

## 3. 未来 Route Surface

第一组 route family 只暴露 own-player storage object operations：

```yaml
candidate_routes:
  - kind: query
    module: storage
    name: GetOwnStorageObject
    route_id: storage.GetOwnStorageObject
    service_method: GetOwnStorageObject
  - kind: query
    module: storage
    name: ListOwnStorageObjects
    route_id: storage.ListOwnStorageObjects
    service_method: ListOwnStorageObjects
  - kind: command
    module: storage
    name: PutOwnStorageObject
    route_id: storage.PutOwnStorageObject
    service_method: PutOwnStorageObject
  - kind: command
    module: storage
    name: DeleteOwnStorageObject
    route_id: storage.DeleteOwnStorageObject
    service_method: DeleteOwnStorageObject
```

规则：

- Route names 必须保持 vibit-native，不复制 Nakama route paths 或 Pitaya route naming。
- `GetOwnStorageObject` 和 `ListOwnStorageObjects` 是 queries。
- `PutOwnStorageObject` 和 `DeleteOwnStorageObject` 是 commands。
- 第一组 route family 仅面向 validated player owner，不暴露任意 owner ids。
- Public ACLs、admin search、group/guild/party/room/match scopes、batch writes、JSON patch、merge semantics、TTL、script hooks 和 large object/blob storage 继续 deferred。
- 未来 route implementation 必须显式注册 routes。不得添加 catch-all storage route 或 reflective handler。

## 4. Protocol Shape

第一版 storage object protocol source candidate 是：

```text
proto/vibit/storage/v1/storage.proto
```

第一版 generated output candidate 是：

```text
runtime/internal/generated/proto/vibit/storage/v1/storage.pb.go
```

第一版 Protobuf package candidate 是：

```text
vibit.storage.v1
```

Candidate messages：

```yaml
messages:
  StorageObject:
    fields:
      collection: string
      key: string
      value_json: string
      version: int64
      created_at: string
      updated_at: string
  GetOwnStorageObjectRequest:
    fields:
      collection: string
      key: string
  GetOwnStorageObjectResponse:
    fields:
      object: StorageObject
  ListOwnStorageObjectsRequest:
    fields:
      collection: string
      limit: int32
      after_key: string
  ListOwnStorageObjectsResponse:
    fields:
      objects: repeated StorageObject
      next_key: string
  PutOwnStorageObjectRequest:
    fields:
      collection: string
      key: string
      value_json: string
      expected_version: int64
  PutOwnStorageObjectResponse:
    fields:
      object: StorageObject
      version: int64
  DeleteOwnStorageObjectRequest:
    fields:
      collection: string
      key: string
      expected_version: int64
  DeleteOwnStorageObjectResponse:
    fields:
      deleted: bool
      version: int64
```

规则：

- 现有 `proto/vibit/protocol/v1/envelope.proto` 必须保持不变，除非后续 protocol ADR 明确修改 envelope semantics。
- `value_json` 不是 log-safe，不能出现在默认 errors、logs、route policy diagnostics 或 test names 中。
- 暴露时间值时应使用 RFC3339Nano UTC text。
- `expected_version` 如果使用 `0` 或字段缺省语义，必须在 implementation ADR 中明确 no-precondition posture。Service 已有 optional expected-version vocabulary；未来 Protobuf mapping 必须保留这个区别，不得发明 merge semantics。
- Protocol shape 不得包含 `owner_id`、`player_id`、`session_id`、raw access tokens、credential material、lookup digests、verifier digests、SQL details、blob bytes、S3 bucket names 或 transport metadata。

## 5. Route Protection And Identity Handoff

第一版 route-policy posture 是：

```yaml
route_policy_requirement: request_token_required
authenticated_wrapper_required: true
request_identity_source: validated_authenticated_request_identity
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
client_supplied_owner_id_allowed: false
```

规则：

- 未来 storage routes 必须是 protected gameplay routes。
- 未来 handlers 必须从现有 protected-route flow 接收 validated `app.RequestIdentity`。
- Envelope/session metadata 中 metadata-only `player_id` 或 `session_id` 绝不能成为 storage owner proof。
- Client payloads 不得选择 owner ids。
- Service 仍负责在 repository access 前拒绝 invalid identity。
- 本 gate 不改变 authentication、token validation、session persistence、first-message binding、WebSocket handshake behavior、bound-identity policy 或 route-protection semantics。

## 6. 未来 Route Flow

未来 implementation 必须保留该顺序：

```yaml
storage_route_flow:
  - websocket_transport_receives_binary_frame_without_reading_credentials
  - protobuf_adapter_decodes_existing_envelope
  - route_policy_requires_authenticated_request_wrapper
  - route_policy_validates_access_token_through_existing authentication/session behavior
  - protobuf_adapter_decodes storage request payload
  - protocol_bridge_maps_payload_to runtime/internal/app/storage request
  - application_handler_calls storage.Service own-object method
  - storage_service_derives_owner_from validated app.RequestIdentity
  - storage_service_uses unit-of-work NewStorageObjectRepository handoff
  - protocol_bridge_maps storage service result to storage response payload
  - protobuf_adapter_returns success or existing error envelope
```

规则：

- WebSocket transport 保持 credential-neutral 和 payload-neutral。
- Protobuf adapter 可以 decode/encode storage payloads，但不能拥有 permission decisions、repository calls 或 storage behavior。
- Application handler registration 属于 `runtime/internal/app/bootstrap` 或等价 application-composition package。
- Handler 只能调用 application storage service。它不得直接调用 repositories、打开 ad hoc SQL、导入 PostgreSQL adapters，或把 transport metadata 当 proof 解析。
- Normal query/command transaction wrapping 只能通过现有 application dispatch/transaction boundaries 使用。本 gate 不添加 startup composition。

## 7. Error Mapping

未来 protocol behavior 必须通过现有 application error envelopes 映射 service public errors：

```yaml
service_public_errors:
  STORAGE_OBJECT_INVALID_REQUEST: application_error_same_code
  STORAGE_OBJECT_NOT_FOUND: application_error_same_code
  STORAGE_OBJECT_ALREADY_EXISTS: application_error_same_code
  STORAGE_OBJECT_VERSION_MISMATCH: application_error_same_code
  STORAGE_OBJECT_UNAVAILABLE: application_error_same_code
  STORAGE_OBJECT_FORBIDDEN: application_error_same_code
```

规则：

- Not-found、owner mismatch 和 deleted-object cases 不得泄露 cross-player existence。
- Version mismatch 可以作为 conflict class 公开。
- Errors 不得包含 raw JSON values、除已验证 caller identity 外的 owner ids、raw token material、credential material、lookup digests、verifier digests、HMAC input/output、SQL strings、database errors、DSNs、headers、cookies、query strings、WebSocket subprotocol values、remote addresses、connection ids 或 session ids。
- Malformed payload 的 protocol adapter errors 必须与 service validation errors 区分，但不得暴露 payload contents。

## 8. Ownership

未来 implementation 必须保留这些 owner：

```yaml
storage_service_owner: runtime/internal/app/storage
application_handler_owner: runtime/internal/app/bootstrap
route_policy_owner: runtime/internal/app
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
storage_repository_interface_owner: runtime/internal/modules/storage
postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
websocket_transport_owner: runtime/internal/platform/transport/ws
startup_owner: runtime/cmd/vibit-server
generated_output_owner: runtime/internal/generated/proto/vibit/storage/v1
protobuf_source_owner: proto/vibit/storage/v1
```

规则：

- Storage runtime behavior 继续位于 `runtime/internal/app/storage`。
- Protocol bridge code 只能 map payload fields。
- Persistence code 仍只负责 storage adapter behavior。
- Generated output 必须由 `.proto` sources 产生，不得手工编辑。
- Startup wiring、route registration 和 generated output 留给后续 implementation work item。

## 9. Required Future Tests

未来 implementation slice 必须添加聚焦测试：

```yaml
required_tests:
  proto_source_and_generated_output_include_storage_messages
  storage_routes_are_registered_only_when_storage_service_is_composed
  storage_routes_are_protected_and_require_authenticated_wrapper
  storage_routes_do_not_accept_metadata_only_player_id_or_session_id
  storage_route_payloads_do_not_include_owner_id
  get_request_maps_to_GetOwnStorageObject
  list_request_maps_to_ListOwnStorageObjects
  put_request_maps_to_PutOwnStorageObject
  delete_request_maps_to_DeleteOwnStorageObject
  storage_success_maps_service_results_to_response_payloads
  storage_public_errors_map_to_error_envelopes
  storage_errors_do_not_leak_value_json_or_repository_details
  websocket_transport_remains_credential_neutral
  existing_protobuf_envelope_remains_unchanged
```

Live PostgreSQL verification 仍是 optional，默认 repository checks 不得要求它。

## 10. Nakama And Pitaya Reference Mapping

Nakama reference mapping：

- 采用 durable player-owned storage objects 是一等 game backend capability 的产品预期。
- 将 collection/key/value/version/list/read/write/delete semantics 适配进 vibit service 和 protocol model。
- 适配 optimistic version conflict semantics，但不复制 Nakama route paths、permission integers、data model names、server runtime APIs、JavaScript/Lua hook model 或 direct compatibility。
- Deferred public ACLs、cross-user reads、system-owned storage、admin search、batch writes、match/party/group scoped storage、TTL、script hooks 和 storage object API compatibility。

Pitaya reference mapping：

- 采用 acceptors、sessions、routes、serializers、handlers 和 backend behavior 的分离。
- 通过保持 `kind/module/name` route identity 显式，并让 application handlers 调用 application services 来适配 handler routing。
- 保持 WebSocket acceptor behavior credential-neutral 和 storage-neutral。
- Deferred Pitaya route naming compatibility、frontend/backend cluster routing、remote calls、groups integration、distributed push 和 RPC/session propagation。

## 11. Stop Conditions

遇到以下任一事项必须停止并打开后续 bounded work item：

- protocol route implementation；
- Protobuf source creation；
- generated output creation or editing；
- protocol bridge implementation；
- application route registration；
- startup wiring；
- repository interface changes；
- PostgreSQL adapter changes；
- migration changes；
- new dependencies；
- authentication/session behavior changes；
- route-protection semantic changes；
- public ACLs 或 cross-owner access；
- admin search；
- group/guild/party/room/match storage scopes；
- batch writes；
- JSON patch 或 merge semantics；
- TTL 或 script hooks；
- large object/blob storage；
- S3-compatible object storage；
- hosted deployments；
- release artifacts；
- public announcements；
- paid promotion；
- direct Nakama/Pitaya API compatibility。

## 12. Verification

本 gate 的 repository check rule 是：

```text
runtime.storage_objects_protocol_route_gate
```

本 gate 预期 verification：

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.storage_objects_protocol_route_gate
node tools/vibit check change define-storage-objects-protocol-route-gate --json
node tools/vibit check module storage --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

本 gate 不需要 Go tests，因为它不得添加 Go runtime behavior。
