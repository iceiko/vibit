# Storage Objects Repository Boundary

状态：Accepted v0.1
最后更新：2026-05-22
范围：在 PostgreSQL `storage_objects` migration source 之后，为未来 storage-neutral storage objects repository 定义 gate-only boundary
依赖：`docs/storage-objects-behavior-gate.md`、`docs/storage-objects-persistence-schema-gate.md`、`decisions/ADR-0111-storage-objects-migration-source.md`、`docs/reference-game-server-alignment.md`
Canonical decision：`ADR-0112`

英文原文是 `docs/storage-objects-repository-boundary.md`。英文文件是架构、契约和治理的权威版本；本文是面向中文读者的对应翻译。

本文定义 storage objects repository boundary。它是 gate artifact。它不添加 Go repository interfaces、PostgreSQL storage adapters、runtime behavior、protocol routes、Protobuf source、generated output、dependencies、migrations、automatic startup migration behavior、broad operations/admin behavior、authentication/session behavior changes、hosted deployments、release artifacts、public announcements、paid promotion、broad product module expansion、large object/blob storage、S3-compatible object storage 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Storage objects repository boundary 记录如下：

```yaml
storage_objects_repository_boundary: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0204
decision: ADR-0112
check_rule: runtime.storage_objects_repository_boundary
source_migration_source_decision: ADR-0111
source_migration_source: runtime/migrations/postgres/000006_create_storage_objects.sql
future_repository_owner_candidate: runtime/internal/modules/storage
future_repository_interface_candidate: runtime/internal/modules/storage.Repository
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
storage_objects_logical_table: storage_objects
repository_boundary_gate_only: true
repository_interface_added: false
postgresql_adapter_added: false
runtime_behavior_added: false
authentication_session_behavior_changed: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
dependency_added: false
migration_added: false
broad_operations_admin_behavior_added: false
product_module_expansion_added: false
hosted_deployment_added: false
additional_release_artifacts_authorized: false
public_announcements_beyond_github_release_authorized: false
paid_promotion_authorized: false
large_object_blob_storage_added: false
s3_compatible_object_storage_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_repository_interface_work_item: W-0205
future_repository_interface_direction: storage_objects_repository_interface_implementation
```

## 2. 目的

`W-0203` 已添加 `storage_objects` 的 PostgreSQL migration source。下一条有价值的边界，是未来 implementation code 可以使用的 storage-neutral repository seam，并避免 SQL detail、transport detail 或 protocol assumption 泄漏出去。

这个 boundary 为 prototype-ready storage objects 准备：

- repository ownership；
- candidate value types；
- create/read/list/update/delete vocabulary；
- optimistic conflict 和 version handoff posture；
- redaction 和 error posture；
- PostgreSQL adapter expectations；
- future implementation work 的 stop conditions。

它仍然不是 runtime feature。任何 route、handler、adapter 或 protocol message 在后续 bounded work item 明确授权前，都不得使用 storage objects。

## 3. 所有权

未来 repository 由 storage module 拥有：

```yaml
future_repository_owner_candidate: runtime/internal/modules/storage
future_repository_interface_candidate: runtime/internal/modules/storage.Repository
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
storage_objects_table_owner: runtime.storage
application_layer_owns_request_identity: true
postgresql_adapter_owns_sql_mapping: true
websocket_transport_owns_storage_objects: false
protocol_adapter_owns_storage_objects: false
authentication_module_owns_storage_objects: false
player_module_owns_storage_objects: false
```

规则：

- Future repository interface 必须是 storage-neutral、module-facing。
- Interface 不得提到 PostgreSQL、pgx、SQL rows、goose migrations、prepared statements、connection pools 或 database transaction 实现细节。
- PostgreSQL adapter 未来可以在 `runtime/internal/platform/persistence/postgres` 下实现该 interface，但必须等待单独 adapter gate。
- Application 或 handler code 未来可以通过 module/application boundary 调用 storage object behavior，而不是直接通过 SQL 或 transport state。
- Authentication 和 session code 提供 validated request identity；它们不拥有 storage object rows。
- Player account storage 拥有 player lifecycle state，不拥有 general player-owned storage objects。
- WebSocket transport 拥有 connection plumbing，不拥有 storage objects。
- Protocol adapters 拥有 wire conversion，不拥有 repository behavior。

## 4. 候选值类型

后续 implementation gate 可以重命名或缩小这些 shape，但第一版 repository interface implementation 应从这组词汇开始：

```yaml
candidate_value_types:
  - StorageObject
  - StorageObjectOwner
  - StorageObjectIdentity
  - StorageObjectValue
  - StorageObjectVersion
  - StorageObjectStatus
  - CreateStorageObjectInput
  - GetStorageObjectInput
  - ListStorageObjectsInput
  - UpdateStorageObjectInput
  - DeleteStorageObjectInput
  - StorageObjectConflict
  - StorageObjectRepositoryError
```

第一版字段词汇：

```yaml
storage_object_record:
  object_id: server_generated_opaque_id
  owner_kind: player
  owner_id: validated_player_id_handoff
  collection: bounded_text_identifier
  object_key: bounded_text_identifier
  value_json: small_json_object_payload
  version: server_managed_bigint_revision
  created_at: server_timestamp
  updated_at: server_timestamp
  deleted_at: nullable_server_timestamp
```

规则：

- `owner_kind` 必须保持 closed vocabulary。第一版 allowed owner kind 是 `player`。
- `owner_id` 是 input identity handoff，不是 proof。未来 runtime behavior 在调用 repository 前必须从 validated request identity 推导它。
- `collection` 和 `object_key` 必须保持 bounded identifiers；它们不是 path names、file names、bucket names 或 SQL fragments。
- `value_json` 默认不是 log-safe。
- `version` 由 server 管理，不能被接受为 client-authoritative state。
- `deleted_at` 只表达 soft-delete storage state；public delete semantics 留给后续 behavior decision。

## 5. 候选 Repository 能力

第一版 storage-neutral capability family 是：

```yaml
candidate_repository_capabilities:
  - CreateStorageObject
  - GetStorageObject
  - ListStorageObjects
  - UpdateStorageObject
  - DeleteStorageObject
```

能力规则：

- `CreateStorageObject` 未来可以为 already validated owner identity 以及 server-validated collection/key/value 创建 row。
- `GetStorageObject` 是 storage lookup。它不得 authenticate users、validate access tokens 或 create request identity。
- `ListStorageObjects` 必须 collection-scoped 且 pagination-ready。没有后续 gate 时，它不得变成 arbitrary owner search 或 admin inspection。
- `UpdateStorageObject` 必须支持 future expected-version handoff，并返回 typed conflict result，而不是泄漏 storage driver errors。
- `DeleteStorageObject` 在 caller 提供 expected version 时必须 version-aware，并且必须保留 future behavior gate 定义的 not-found 和 owner-mismatch leakage rules。
- 所有 methods 必须返回 typed module-owned records 和 errors，而不是 raw SQL rows 或 database driver errors。

Future repository interface 可以选择更短的名称，但必须保留 create、read、list、update、delete 和 conflict handling 的语义分离。

## 6. Version And Conflict Handoff

Repository boundary 为 optimistic concurrency 做准备，但不实现行为：

```yaml
version_storage: BIGINT
initial_create_version: 1
version_owner: server
client_authoritative_version_allowed: false
expected_version_handoff: future_behavior_or_interface_gate
conflict_public_shape: deferred_to_protocol_gate
```

候选 conflict classes：

```yaml
candidate_conflict_classes:
  - object_already_exists
  - object_not_found
  - version_mismatch
  - invalid_expected_version
  - owner_scope_mismatch
  - deleted_object
  - storage_unavailable
```

规则：

- Repository methods 未来可以区分 internal typed conflicts，但 public protocol error mapping 仍然 deferred。
- Version equality 不是 authentication proof。
- Stale expected version 不得被隐藏折叠成 successful write。
- Owner mismatch 不得泄漏另一个 player 的 object existence。
- PostgreSQL adapter 必须把 unique-index 和 affected-row outcomes 映射为 typed repository conflicts，不暴露 driver error text。

## 7. Redaction And Logging

Storage objects 是 game state，不是 authentication state，但仍然需要 redaction discipline。

```yaml
value_log_safe: false
collection_log_safe: conditional_after_validation
object_key_log_safe: conditional_after_validation
object_id_log_safe: conditional_after_validation
owner_id_log_safe: false_by_default
forbidden_repository_material:
  - raw_access_token
  - raw_credential
  - credential_lookup_digest
  - credential_verifier_digest
  - token_lookup_digest
  - token_verifier_digest
  - verifier_key
  - websocket_connection_id
  - websocket_subprotocol
  - remote_address
  - blob_bytes
  - file_path
  - s3_bucket
  - s3_object_key
```

规则：

- Repository errors 必须 redacted 且 typed。
- Raw values 和 storage driver errors 默认不得记录到日志。
- `value_json` 可能包含 player state 和 user-adjacent data；除非后续 redaction policy 收窄，否则未来 behavior 必须把它视为 non-log-safe。
- Repository 不得存储或返回 authentication material、token material、verifier digests、transport metadata、file/blob payloads 或 S3-compatible object references。

## 8. PostgreSQL Adapter Expectations

未来 PostgreSQL adapter 可以映射到：

```yaml
logical_table: storage_objects
active_identity_unique_index: storage_objects_active_identity_uq
owner_collection_index: storage_objects_owner_collection_idx
updated_at_index: storage_objects_updated_at_idx
soft_delete_column: deleted_at
version_column: version
```

Adapter expectations：

- SQL execution 属于 `runtime/internal/platform/persistence/postgres`。
- Unit-of-work 和 transaction handoff 必须遵循现有 platform transaction boundary。
- SQL 不得泄漏到 `runtime/internal/modules/storage`。
- Adapter 必须保留 `owner_kind + owner_id + collection + object_key` 的 active-object uniqueness。
- Updates 和 deletes 必须检查 affected rows，并在 expected version 被提供时 version-aware。
- Soft-deleted rows 不得被 ordinary active-object lookups 返回，除非 future admin/operations gate 定义该行为。
- Adapter gate 被接受后，adapter tests 应覆盖 unique conflicts、not found、stale version、soft delete、owner/collection listing、value copying、timestamp mapping 和 redacted errors。

## 9. 与 Runtime、Protocol 和 Authentication 的关系

这个 boundary 不改变 runtime 或 protocol behavior：

```yaml
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
authentication_session_behavior_changed: false
route_protection_semantics_changed: false
websocket_transport_credential_neutral: true
```

规则：

- Future runtime behavior 必须先从 validated request identity 推导 owner identity，再调用 repository。
- Repository 不得 parse access tokens、validate sessions、inspect WebSocket metadata 或 set request identity。
- Repository 不得暴露 Protobuf request/response types。
- Future protocol routes 在 protocol gate 定义 exact route names、message shapes、generated output 和 error mapping 前都只是 planning candidates。
- 这个 boundary 不授权 public read/write permissions、ACLs、admin search、cross-player access 或 direct Nakama/Pitaya API compatibility。

## 10. Nakama 和 Pitaya 参考映射

Nakama reference mapping：

```yaml
adopted_concepts:
  - storage_objects_are_core_durable_game_state
  - collection_key_identity_supports_prototype_state
  - version_conflicts_need_first_class_semantics
adapted_concepts:
  - repository_is_vibit_storage_neutral_module_boundary
  - owner_identity_comes_from_validated_request_identity
  - public_api_compatibility_is_not_a_goal
deferred_concepts:
  - public_or_permission_bit_storage_objects
  - admin_storage_object_search
  - multi_owner_and_group_storage
  - direct_client_api_compatibility
rejected_for_now:
  - direct_nakama_storage_api_compatibility
```

Pitaya reference mapping：

```yaml
adopted_concepts:
  - handlers_should_not_own_persistence_details
  - routing_and_storage_boundaries_should_remain_separate
adapted_concepts:
  - storage_repository_is_module_owned_not_transport_owned
  - durable_state_handoff_stays_below_protocol_adapter
deferred_concepts:
  - cluster_safe_storage_routing
  - distributed_cache_or_document_store
  - frontend_backend_rpc_storage_facade
rejected_for_now:
  - direct_pitaya_handler_or_rpc_compatibility
```

Nakama 和 Pitaya 用来指导 capability planning。它们不覆盖 vibit constitution、ADRs、manifests、generated boundaries 或 verification commands。

## 11. Future Implementation Queue

这个 boundary 之后，后续工作仍应拆分：

```yaml
future_work_items:
  storage_objects_repository_interface_implementation:
    may_add:
      - runtime/internal/modules/storage
      - storage-neutral repository types
      - focused unit tests
    must_not_add:
      - PostgreSQL adapter
      - SQL execution
      - protocol routes
      - runtime behavior
  storage_objects_postgresql_adapter_gate:
    may_define:
      - adapter ownership
      - transaction handoff
      - SQL query shape
      - adapter tests
  storage_objects_postgresql_adapter_implementation:
    may_add:
      - platform persistence adapter
      - unit-of-work factory handoff if separately authorized
  storage_objects_runtime_behavior_gate:
    may_define:
      - owner identity derivation
      - permissions
      - validation
      - conflict semantics
      - route-policy expectations
  storage_objects_protocol_gate:
    may_define:
      - Protobuf source
      - generated output
      - route names
      - error mapping
```

没有新 ADR 时，不要把这些合并成一个 broad storage subsystem slice。

## 12. Stop Conditions

出现以下任一情况时停止，并要求后续 bounded work item：

- 创建 `runtime/internal/modules/storage` source files；
- 添加 repository interface implementation；
- 添加 PostgreSQL adapter files 或 SQL execution；
- changing migrations；
- 添加 runtime handlers 或 application behavior；
- 添加 protocol routes、Protobuf sources 或 generated output；
- 添加 dependencies；
- changing authentication/session behavior 或 route protection；
- 添加 broad operations/admin behavior；
- 添加 hosted deployments；
- 创建 release binaries、packages、containers、checksums、provenance files、signing artifacts、install scripts、registry publications 或 SDK packages；
- 执行 GitHub release record 之外的 public announcements；
- 运行 paid promotion；
- 添加 large object/blob storage 或 S3-compatible object storage；
- 添加 direct Nakama/Pitaya API compatibility。

## 13. Verification

该 boundary 的 repository verification 是：

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit inspect rule runtime.storage_objects_repository_boundary
node tools/vibit check change define-storage-objects-repository-boundary --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
cd runtime && go test ./...
git diff --check
```

Repository check rule 是：

```yaml
runtime.storage_objects_repository_boundary
```
