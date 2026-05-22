# Storage Objects Runtime Behavior Gate

状态：Accepted v0.1
最后更新：2026-05-22
范围：PostgreSQL adapter 之后，未来 application-owned storage objects runtime behavior 的 gate-only 边界
依赖：`docs/storage-objects-behavior-gate.md`、`docs/storage-objects-repository-boundary.md`、`docs/storage-objects-postgresql-adapter-gate.md`、`runtime/internal/modules/storage/repository.go`、`runtime/internal/platform/persistence/postgres/storage_object_repository.go`、`docs/runtime-protocol-adapter.md`、`docs/bound-identity-route-policy-gate.md`、`docs/runtime-session-validation-gate.md`、`docs/reference-game-server-alignment.md`
规范决策：`ADR-0116`

配对英文源文件是 `docs/storage-objects-runtime-behavior-gate.md`。英文文件是权威版本。

本文定义 storage objects runtime behavior gate。它是 gate artifact，不添加 runtime behavior implementation、runtime handlers、startup wiring、protocol routes、Protobuf source、generated output、repository interface changes、PostgreSQL adapter changes、migration changes、dependencies、authentication/session behavior changes、hosted deployments、release artifacts、public announcements、paid promotion、broad product module expansion、large object/blob storage、S3-compatible object storage 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Storage objects runtime behavior gate 记录为：

```yaml
storage_objects_runtime_behavior_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0208
decision: ADR-0116
check_rule: runtime.storage_objects_runtime_behavior_gate
source_adapter_implementation_decision: ADR-0115
source_adapter: runtime/internal/platform/persistence/postgres/storage_object_repository.go
source_repository_interface_decision: ADR-0113
repository_interface: runtime/internal/modules/storage.Repository
repository_interface_source: runtime/internal/modules/storage/repository.go
future_runtime_owner_candidate: runtime/internal/app
future_storage_application_package_candidate: runtime/internal/app/storage
future_runtime_service_source_candidate: runtime/internal/app/storage/service.go
future_runtime_service_test_candidate: runtime/internal/app/storage/service_test.go
request_identity_required: true
metadata_only_player_id_allowed_as_proof: false
first_owner_kind: player
owner_id_source: validated_request_identity_player_id
route_policy_requirement: request_token_required
service_application_owner: runtime/internal/app
repository_handoff: unit_of_work_storage_repository_factory
unit_of_work_handoff_required: true
runtime_behavior_gate_only: true
runtime_behavior_added: false
runtime_handlers_added: false
startup_wiring_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
dependency_added: false
migration_added: false
authentication_session_behavior_changed: false
large_object_blob_storage_added: false
s3_compatible_object_storage_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_runtime_behavior_implementation_work_item: W-0209
future_runtime_behavior_implementation_direction: storage_objects_runtime_behavior_implementation
```

## 2. Purpose

`W-0207` 已实现 `runtime/internal/modules/storage.Repository` 的 PostgreSQL adapter。下一步不应直接添加 route 或 protocol change，而是先定义 runtime behavior gate，说明 application code 后续如何把 validated route request 转成 storage repository operations。

该 gate 在实现前记录未来 behavior shape：

- application ownership；
- 从 validated request identity 派生 owner identity；
- permission 与 route-policy posture；
- validation 与 conflict mapping expectations；
- unit-of-work 与 repository handoff；
- redaction rules；
- test expectations；
- stop conditions，避免 protocol、generated output、authentication/session changes 和更广 storage product scope 泄入本 slice。

Nakama 提供 durable storage-object gameplay capability 的产品压力。Pitaya 强调 handlers、route policy 与 persistence responsibilities 的分离。vibit 把这些参考适配成显式 application-owned behavior 和 checks，而不是 direct public API compatibility。

## 3. Ownership

未来 runtime behavior 是 application-owned：

```yaml
future_runtime_owner_candidate: runtime/internal/app
future_storage_application_package_candidate: runtime/internal/app/storage
future_runtime_service_source_candidate: runtime/internal/app/storage/service.go
future_runtime_service_test_candidate: runtime/internal/app/storage/service_test.go
repository_interface_owner: runtime/internal/modules/storage
postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
route_policy_owner: runtime/internal/app
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
websocket_transport_owner: runtime/internal/platform/transport/ws
```

规则：

- 未来 service behavior 可以放在 `runtime/internal/app/storage`，或 implementation slice 明确 ratify 的等价 application-owned package。
- Service 只能通过 application 或 unit-of-work dependencies 调用 `runtime/internal/modules/storage.Repository`。
- Service 不得 import PostgreSQL adapter packages、SQL row types、migration packages、WebSocket transport packages、generated Protobuf packages、S3 SDKs 或 MinIO clients。
- Storage module 仍拥有 storage-neutral value types、validation helpers 和 repository error vocabulary。
- PostgreSQL adapter 仍是 persistence-only，不得派生 request identity、route policy 或 public protocol errors。
- Transport 与 protocol adapters 不拥有 storage object permission 或业务行为。

## 4. Request Identity And Owner Derivation

第一版 runtime behavior posture 是 player-owned：

```yaml
first_owner_kind: player
owner_id_source: validated_request_identity_player_id
request_identity_required: true
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
client_supplied_owner_id_allowed: false
owner_id_overrides_allowed: false
```

规则：

- 未来 storage operation 必须从 validated `app.RequestIdentity` 派生 `storage.StorageObjectOwner{Kind: player, ID: identity.PlayerID}`。
- `RequestIdentity.Status` 必须是 `validated`。
- `RequestIdentity.ActorKind` 必须是 `player`。
- `RequestIdentity.PlayerIDValidated` 必须为 true。
- `RequestIdentity.PlayerID` 必须非空，并且在 actor identity 同时存在时与其匹配。
- Envelope/session metadata 中 metadata-only `player_id` 永远不能满足该 gate。
- 持久化 `session_id` 本身永远不能成为 proof。
- 第一版 posture 中，client payload 不得选择另一个 owner id。

该 gate 不改变 `RequestIdentity`、access-token validation、bound connection identity、durable runtime session validation 或 WebSocket handshake behavior。它只记录未来 storage behavior 在 repository access 前必须要求的 identity 条件。

## 5. Future Runtime Behavior Shape

未来第一版实现可以暴露 application service 的这些 candidate operations：

```yaml
candidate_operations:
  - get_own_storage_object
  - list_own_storage_objects
  - put_own_storage_object
  - delete_own_storage_object
```

推荐第一版 posture：

- `get` 按 collection 与 key 读取 validated player 自己的一个 active object。
- `list` 按 collection 分页读取 validated player 自己的 active objects。
- `put` 为 validated player 创建或替换一个 object，并返回当前 version。
- `delete` soft-delete validated player 的一个 object，并返回稳定 success 或 conflict behavior。

规则：

- Runtime behavior 必须使用 server-derived owner identity。
- Runtime behavior 必须在 repository calls 前验证 collection、key、value shape、value size、expected version 和 list limit。
- 第一版实现不得暴露 cross-owner reads、cross-owner writes、public ACLs、admin bypass、group/guild scopes、room scopes、match scopes、batch writes、JSON patch、merge semantics、TTL 或 server script hooks。
- 除非后续 protocol gate 授权，runtime behavior 不得添加 public protocol routes 或 generated output。

## 6. Candidate Application Service Shape

第一版 implementation slice 应定义小型 application-owned service boundary。Candidate inputs 和 outputs：

```yaml
candidate_request_fields:
  - request_identity
  - collection
  - key
  - value_json
  - expected_version
  - list_limit
  - after_object_key
candidate_result_fields:
  - object
  - objects
  - next_object_key
  - version
  - public_error_code
```

规则：

- Service 应接收已 normalize 的 application identity，而不是 raw tokens、cookies、headers、WebSocket subprotocol values 或 envelope proof carriers。
- Service 应在 repository handoff 前调用 storage module normalizers。
- Service 应避免把 raw JSON values 放入默认 errors 和 logs。
- Service 应暴露稳定 public error codes 或 classes，供后续 runtime handlers 映射。
- 本 gate slice 不添加 route registration、Protobuf conversion、startup composition 或 command/query dispatch wiring。

## 7. Validation Rules

未来 runtime behavior 必须在 persistence 前执行 validation：

```yaml
validation_required:
  - validated_player_identity
  - collection_non_empty_length_bounded
  - key_non_empty_length_bounded
  - value_top_level_json_object
  - value_size_bounded
  - expected_version_positive_when_present
  - list_limit_bounded
  - pagination_cursor_bounded
```

规则：

- Collection 和 key validation 应复用 storage module normalization rules，除非未来 contract 明确收紧 protocol-visible syntax。
- Value JSON 不是 log-safe，跨边界时必须保持 copied 或 immutable。
- Missing expected version behavior 必须在 implementation tests 中明确。
- Invalid input 应尽可能在 repository mutation 前失败。
- Repository unavailable errors 必须保持 redacted。

## 8. Permission And Route Policy Posture

第一版 route-policy posture 保守：

```yaml
route_policy_requirement: request_token_required
public_storage_routes_allowed: false
bound_connection_required_by_this_gate: false
session_validated_required_by_this_gate: false
bound_session_required_by_this_gate: false
```

后续 public contracts 的 candidate permission families：

- read own storage object；
- list own storage objects；
- write own storage object；
- delete own storage object。

规则：

- Storage object routes 必须是 protected routes。
- 第一版 posture 应使用现有 `request_token_required` route protection family，除非后续 route-policy ADR 改变 named routes。
- Public routes 不得读取或修改 storage objects。
- Bound connection identity 和 durable session validation 可保留给未来 route families，但该 gate 不要求它们，也不改变 ordinary protected route behavior。
- Metadata-only identity 必须 fail closed。

## 9. Conflict And Error Mapping

未来 runtime behavior 必须把 storage repository errors 映射为稳定 application classes：

```yaml
candidate_public_error_classes:
  - STORAGE_OBJECT_INVALID_REQUEST
  - STORAGE_OBJECT_NOT_FOUND
  - STORAGE_OBJECT_ALREADY_EXISTS
  - STORAGE_OBJECT_VERSION_MISMATCH
  - STORAGE_OBJECT_UNAVAILABLE
  - STORAGE_OBJECT_FORBIDDEN
```

规则：

- Not-found、owner mismatch 和 deleted-object cases 必须避免 cross-player existence leaks。
- Version mismatch 可以作为 public conflict class，但 stored values、raw JSON、SQL details、driver errors、DSNs、credentials、token material、verifier digests 和 route proof carriers 不得泄露。
- Repository `storage_unavailable` errors 必须映射为 retryable 或 unavailable class，且不暴露 platform internals。
- 当 request identity 未 validated 时，permission failure 必须在 repository access 前发生。
- Runtime behavior 不得在现有 application route-protection classes 之外添加 authentication/token/session failure details。

## 10. Unit-Of-Work And Repository Handoff

未来 runtime behavior 应使用既有 application transaction boundary：

```yaml
repository_handoff: unit_of_work_storage_repository_factory
unit_of_work_handoff_required: true
service_starts_transactions: false
service_commits_transactions: false
service_rolls_back_transactions: false
repository_factory_candidate: NewStorageObjectRepository
```

规则：

- Mutating operations 从 command dispatch 调用时，应运行在现有 command unit-of-work boundary 内。
- Service 应从 application dependency 或 implementation ratify 的 unit-of-work capability 获取 `storage.Repository`。
- Service 不得为了创建 repositories 而 import PostgreSQL packages。
- 如果 implementation 选择 read path，query behavior 可以使用显式 query repository provider，但 owner derivation 和 validation rules 不变。
- 该 gate 不改变 `TransactionalDispatcher` 或 startup wiring。

## 11. Relationship To Protocol

该 gate 不添加 protocol behavior：

```yaml
storage_protocol_routes_added: false
protobuf_storage_messages_added: false
generated_storage_output_added: false
existing_envelope_changed: false
websocket_transport_changed: false
```

规则：

- 未来 protocol work 必须在单独 bounded work item 中定义 exact route names、module/name routing、request/response messages、generated output 和 error mapping。
- `docs/storage-objects-behavior-gate.md` 中的 candidate route names 仍只是 planning vocabulary。
- WebSocket transport 必须保持 credential-neutral。
- Protobuf adapters 不得派生 permissions 或 owner identity。

## 12. Relationship To Authentication And Session

该 gate 使用现有 validated request identity vocabulary，但不改变 authentication 或 session behavior：

```yaml
access_token_validation_changed: false
request_identity_shape_changed: false
session_validation_changed: false
websocket_handshake_authentication_changed: false
first_message_connection_binding_changed: false
metadata_only_player_id_allowed_as_proof: false
```

规则：

- Access-token validation 仍是当前 protected routes 的 proof path。
- Durable session validation 仍由独立边界拥有。
- First-message connection binding 不会因该 gate 授权 ordinary storage operations。
- WebSocket handshake authentication 仍然 deferred。
- Implementation slice 不得放松 metadata-only identity protections。

## 13. Test Expectations

后续 implementation slice 应添加 focused tests：

```yaml
future_tests:
  - metadata_only_identity_rejected_before_repository
  - validated_player_identity_derives_owner
  - client_owner_id_ignored_or_rejected
  - get_maps_repository_object_to_runtime_result
  - list_is_owner_collection_scoped_and_bounded
  - put_validates_json_object_and_version
  - delete_checks_expected_version_when_supplied
  - owner_mismatch_and_not_found_do_not_leak_existence
  - repository_errors_are_redacted
  - unit_of_work_repository_handoff_is_used_for_mutations
  - no_protocol_route_or_generated_output_required
```

规则：

- Tests 应尽量使用 fake repositories 或 fake unit-of-work providers。
- 默认 repository checks 不得要求 live PostgreSQL。
- Tests 不得打印 raw values、DSNs、credentials、access tokens、verifier material、lookup digests、verifier digests、cookies、headers、query strings 或 WebSocket subprotocol values。

## 14. Stop Conditions

添加以下内容前必须停止，并要求后续 bounded work item：

- runtime behavior implementation；
- runtime handlers 或 route registration；
- protocol routes；
- Protobuf source files；
- generated output；
- repository interface changes；
- PostgreSQL adapter changes；
- migration changes；
- new dependencies；
- startup wiring；
- authentication/session behavior changes；
- WebSocket handshake authentication；
- public ACLs、admin storage search、group/guild/party/room/match storage scopes、batch writes、JSON patch、merge semantics、TTL 或 script hooks；
- large object/blob storage；
- S3-compatible object storage；
- hosted deployments；
- release artifacts；
- GitHub release record 之外的 public announcements；
- paid promotion；
- direct Nakama/Pitaya API compatibility。

## 15. Verification

该 gate 通过 repository checks 验证：

```yaml
check_rule: runtime.storage_objects_runtime_behavior_gate
required_commands:
  - node -c tools/vibit
  - node tools/vibit inspect next --json
  - node tools/vibit inspect rule runtime.storage_objects_runtime_behavior_gate
  - node tools/vibit check change define-storage-objects-runtime-behavior-gate --json
  - node tools/vibit check module storage --json
  - node tools/vibit check work --json
  - node tools/vibit check runtime --json
  - node tools/vibit check memory --json
  - node tools/vibit check schemas --json
  - node tools/vibit check all --json
  - cd runtime && go test ./...
  - git diff --check
```

该 gate 不添加 runtime behavior 或 SQL execution，因此不需要 live PostgreSQL verification。
