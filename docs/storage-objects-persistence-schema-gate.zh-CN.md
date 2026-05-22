# Storage Objects Persistence Schema Gate 中文版

状态：Accepted v0.1
最后更新：2026-05-22
范围：在添加 migration source 前定义第一版 PostgreSQL storage objects persistence schema 的 gate
依赖：`docs/storage-objects-behavior-gate.md`、`docs/postgresql-persistence-boundary.md`、`docs/reference-game-server-alignment.md`
权威决策：`ADR-0110`

本文件是 `docs/storage-objects-persistence-schema-gate.md` 的简体中文译本。英文版本是权威版本。

本文定义第一版 storage objects persistence schema gate。它是 gate artifact。它不添加 SQL migration source，不创建 `storage_objects` table，不实现 storage objects runtime behavior，不添加 protocol routes，不添加 Protobuf source 或 generated output，不添加 dependencies，不添加 repository interfaces，不添加 storage adapters，不扩展 operations/admin behavior，不添加 hosted deployments，不创建 release artifacts，不执行 public announcements，不运行 paid promotion，不改变 authentication/session behavior，不添加 large object/blob storage，不添加 S3-compatible object storage，不实现 broad product module，也不添加 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

storage objects persistence schema gate 记录如下：

```yaml
storage_objects_persistence_schema_gate: defined
completed_work_item: W-0202
decision: ADR-0110
check_rule: runtime.storage_objects_persistence_schema_gate
source_behavior_gate_decision: ADR-0109
source_behavior_gate_standard: docs/storage-objects-behavior-gate.md
gate_standard: docs/storage-objects-persistence-schema-gate.md
gate_standard_translation: docs/storage-objects-persistence-schema-gate.zh-CN.md
target_stage: prototype_ready_foundation
selected_first_storage_objects_store: postgres
future_storage_objects_logical_table: storage_objects
future_migration_source_candidate: runtime/migrations/postgres/000006_create_storage_objects.sql
future_repository_owner_candidate: runtime/internal/modules/storage
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
table_candidate_recorded: true
owner_identity_representation_recorded: true
collection_key_constraints_recorded: true
value_representation_recorded: true
version_representation_recorded: true
timestamp_posture_recorded: true
uniqueness_index_posture_recorded: true
redaction_posture_recorded: true
future_repository_adapter_boundaries_recorded: true
future_migration_source_candidate_recorded: true
schema_gate_only: true
migration_source_added: false
storage_objects_table_added: false
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
dependency_added: false
repository_interface_changed: false
storage_adapter_changed: false
broad_operations_admin_behavior_added: false
authentication_session_behavior_changed: false
product_module_expansion_added: false
hosted_deployment_added: false
additional_release_artifacts_authorized: false
public_announcements_beyond_github_release_authorized: false
paid_promotion_authorized: false
large_object_blob_storage_added: false
s3_compatible_object_storage_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_migration_work_item: W-0203
future_migration_direction: storage_objects_migration_source
```

## 2. Product Intent

`ADR-0109` 选择 player-owned small JSON storage objects 作为 inventory 之外第一版 general durable game-state behavior。下一步是在添加 SQL 前，把未来 persistence shape 明确下来。

本 gate 让未来 migration 可检查：

- table candidate 已知；
- ownership 和 identity columns 已选择；
- collection/key constraints 有边界；
- JSON value storage 明确是 small-object game state；
- optimistic versioning 有具体 storage posture；
- indexes 和 uniqueness 在 implementation 前规划；
- redaction 和未来 repository/adapter boundaries 是显式的。

本 gate 保持保守。它不添加 migration file，只准备下一项 migration-source-only slice。

## 3. Selected Store And Table

第一版 storage objects persistence target 是 PostgreSQL：

```yaml
selected_first_storage_objects_store: postgres
future_storage_objects_logical_table: storage_objects
future_migration_source_candidate: runtime/migrations/postgres/000006_create_storage_objects.sql
future_repository_boundary: separate_future_work_item
future_postgresql_adapter: separate_future_work_item
```

理由：

- PostgreSQL 是 vibit 第一版 accepted authoritative durable store。
- 当前 source alpha 已用 PostgreSQL 支撑 inventory、player accounts、authentication verifier records、token verifier records 和 runtime sessions。
- Storage objects 在 runtime behavior 存在前，需要 transactionally inspectable ownership 和 versioning。
- 先使用 PostgreSQL，避免在需求具体化前引入 document database 或 object/blob storage dependency。

## 4. Future `storage_objects` Table Candidate

未来第一版 migration 可以定义一个 logical table：

```yaml
storage_objects:
  primary_key_candidate:
    - object_id
  required_columns:
    - object_id
    - owner_kind
    - owner_id
    - collection
    - object_key
    - value_json
    - version
    - created_at
    - updated_at
  nullable_columns:
    - deleted_at
  forbidden_columns:
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

`object_id` 是 server-generated opaque record id。它不是 object 的 public identity。逻辑身份仍然是：

```text
owner_kind + owner_id + collection + object_key
```

第一版 table candidate 可以使用 `object_key`，而不是裸 `key` column，以减少 SQL 可读性和工具解析歧义。

## 5. Owner Identity Representation

第一版 schema posture 只支持 player-owned objects：

```yaml
owner_kind_first_value: player
owner_id_source: validated_request_identity_player_id
owner_player_fk_candidate: player_accounts(player_id)
owner_kind_check_candidate: owner_kind = 'player'
```

规则：

- `owner_kind` 在第一版 migration source 中必须是 closed vocabulary。
- `owner_id` 必须非空。
- 第一版 owner kind 是 `player`。
- 第一候选是外键到 `player_accounts(player_id)`。
- 表不得把 client-supplied owner id 当作 proof；后续 runtime behavior 必须从 validated identity 得出 owner id。

Deferred owner identities：

- global；
- group/guild；
- party；
- room；
- match；
- server shard；
- public catalog；
- admin-owned objects。

## 6. Collection And Key Constraints

未来 migration 应在 runtime behavior 依赖 collection/key 前约束这些字段。

第一候选：

```yaml
collection_column: collection
key_column: object_key
collection_type: TEXT
object_key_type: TEXT
collection_not_blank: true
object_key_not_blank: true
collection_max_length_candidate: 128
object_key_max_length_candidate: 256
path_semantics_allowed: false
case_sensitive: true
```

推荐 SQL-level checks：

- `length(btrim(collection)) > 0`；
- `length(collection) <= 128`；
- `length(btrim(object_key)) > 0`；
- `length(object_key) <= 256`。

如果 precise protocol identifier rules 先被 ratify，后续 implementation gate 或 migration-source slice 可以添加更严格的 ASCII-safe 或 pattern checks。

## 7. Value Representation

第一版 value representation candidate 是 PostgreSQL `JSONB`：

```yaml
value_column: value_json
value_type_candidate: JSONB
value_top_level_shape: object
value_binary_storage: false
value_log_safe: false
```

规则：

- value 必须保持 small durable game state，不是 large binary content。
- value 不得包含 raw credentials、raw tokens、verifier keys、DSNs、digest bytes、transport metadata 或其他 secrets。
- 当 PostgreSQL 支持时，未来 migration source 应包含 SQL check，要求 top-level value 是 JSON object。
- Maximum byte size 仍是 implementation-gate requirement；如果 migration-source slice ratify 了 maximum，可以选择 SQL check。

Deferred value behavior：

- arbitrary JSON scalar 或 array payloads；
- binary blobs；
- external file references；
- S3-compatible object references；
- encrypted-at-rest application payload envelope；
- document search indexes；
- JSON patch 或 merge behavior。

## 8. Version Representation

第一版 version representation candidate 是 server-managed integer revision：

```yaml
version_column: version
version_type_candidate: BIGINT
initial_version_candidate: 1
version_increment_policy: server_managed_on_successful_mutation
client_authoritative_version: false
```

规则：

- `version` 必须为正数。
- 成功 create 从 version `1` 开始，除非未来 migration-source decision 选择另一个正数起点。
- 成功 update 会递增 version。
- Delete version behavior 必须由未来 runtime behavior gate 或 implementation gate 定义。
- Public protocol 后续可以暴露 string form 或 opaque token，但第一版 stored representation candidate 是 numeric revision。

本 gate 不添加 compare-and-swap runtime behavior。它只定义可以支持 optimistic concurrency 的 persistence representation candidate。

## 9. Timestamp And Deletion Posture

第一版 timestamp posture 是：

```yaml
created_at: TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at: TIMESTAMPTZ NOT NULL DEFAULT now()
deleted_at: TIMESTAMPTZ NULL
soft_delete_candidate: true
hard_delete_candidate: deferred
```

规则：

- 应强制 `updated_at >= created_at`。
- 如果存在 `deleted_at`，应强制 `deleted_at >= created_at`。
- Soft delete 是第一版 schema candidate，因为它能支持未来 conflict 和 audit behavior，而不需要立即添加 event history。
- Runtime behavior 后续仍可选择 delete 是隐藏 rows、保留 tombstones，还是经过 cleanup gate 后 hard-delete rows。

本 schema 不添加 audit/event tables。

## 10. Uniqueness And Index Posture

第一版 uniqueness posture 是：

```yaml
logical_identity_unique_candidate:
  - owner_kind
  - owner_id
  - collection
  - object_key
active_unique_index_candidate:
  - owner_kind
  - owner_id
  - collection
  - object_key
  - deleted_at IS NULL
```

未来 migration-source slice 推荐 indexes：

- `(owner_kind, owner_id, collection, object_key)` 上的 active logical identity unique index，条件是 `deleted_at IS NULL`；
- `(owner_kind, owner_id, collection)` lookup index；
- future diagnostics 或 cleanup 使用的 updated-at index；
- deleted-at index 只有在后续 cleanup 或 tombstone queries 获得授权后才添加。

本 gate 不授权 global search、cross-owner search、JSONB GIN indexes、admin listing、analytics indexes 或 operations dashboards。

## 11. Ownership Boundaries

未来 storage object behavior 应拥有自己的 module boundary：

```yaml
future_repository_owner_candidate: runtime/internal/modules/storage
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
future_contract_owner_candidate: contracts/storage
future_proto_source_candidate: proto/vibit/storage/v1
storage_module_owns_storage_objects: true
inventory_module_owns_storage_objects: false
player_module_owns_storage_objects: false
authentication_module_owns_storage_objects: false
websocket_transport_owns_storage_objects: false
```

规则：

- 未来 storage module 拥有 storage object domain behavior。
- Player module 拥有 player account lifecycle，不拥有 storage object behavior。
- Inventory module 拥有 inventory state，不拥有 general storage objects。
- Authentication 拥有 credentials/tokens/sessions，不拥有 storage object data。
- WebSocket transport 拥有 connection plumbing，不拥有 durable object records。
- PostgreSQL adapters 只有在 repository boundaries 获得授权后才能实现 storage repositories。

## 12. Redaction Posture

Storage object values 默认不是 log-safe。

默认 not log-safe：

- `value_json`；
- validation errors 中的 raw object values；
- 与 collection/key 组合出现在 public diagnostics 中的 owner ids；
- collection 和 key，直到后续 redaction decision 允许特定 logs；
- rejected payloads 中的 JSON path fragments；
- 包含 values 的 database error details。

Forbidden secret material：

- raw device credentials；
- raw access tokens；
- verifier keys；
- lookup 或 verifier digests；
- 带 credentials 的 PostgreSQL DSNs；
- headers、cookies、query strings、WebSocket subprotocol values、remote addresses 或 concrete transport metadata。

未来 adapters 和 handlers 必须返回 redacted errors。

## 13. Future Migration Source Expectations

下一项 bounded work item 可以添加：

```text
runtime/migrations/postgres/000006_create_storage_objects.sql
```

该 migration-source-only slice 可以为 `storage_objects` 添加 SQL DDL、comments、indexes 和 migration checks。

它不得添加：

- Go repository interfaces；
- PostgreSQL adapter behavior；
- runtime handlers；
- protocol routes；
- Protobuf source files；
- generated output；
- startup wiring；
- automatic migration apply behavior；
- dependencies；
- large object/blob storage；
- S3-compatible object storage；
- direct Nakama/Pitaya API compatibility。

## 14. Verification Expectations

未来 migration-source slice 应验证：

- goose up/down markers；
- table name 和 required columns；
- owner、collection、key、value、version、timestamp 和 deletion checks；
- active logical identity uniqueness；
- 没有 forbidden secret、digest、transport、blob、file 或 S3 columns；
- 没有 Go runtime behavior；
- migration boundary 的 repository checks。

后续 repository/adapter/runtime work 只有在对应 implementation gates accepted 后，才应添加 create/read/update/delete/conflict/cross-owner behavior focused tests。

## 15. Stop Conditions

做以下任何事情前必须停止并请求 maintainer authorization：

- 在本 gate 同一个 change 中添加 SQL migration source；
- 实现 storage objects runtime behavior；
- 添加 protocol routes；
- 添加 Protobuf source files 或 generated output；
- 添加 repository interfaces；
- 添加 PostgreSQL 或其他 storage adapters；
- 添加 dependencies；
- 改变 authentication/session semantics；
- 改变 route protection semantics；
- 添加 cross-player、global、group、party、room、match、public、admin 或 ACL object scopes；
- 添加 large object/blob storage 或 S3-compatible object storage；
- 添加 JSONB search indexes 用于 product search behavior；
- 添加 server-side custom logic hooks；
- 添加 broad operations/admin behavior；
- 添加 hosted deployments 或 demos；
- 创建 release binaries、packages、containers、checksums、signing/provenance artifacts、install scripts、registry publications、SDK packages 或 additional release artifacts；
- 执行 GitHub release record 之外的 public announcements；
- 运行 paid promotion；
- 添加 direct Nakama/Pitaya API compatibility。

## 16. Next Work

下一项 bounded direction 是：

```text
W-0203 Add storage objects migration source
```

该工作可以添加 `storage_objects` 的第一版 SQL migration source 和对应 static checks，同时继续推迟 repository interfaces、adapters、protocol 和 runtime behavior。
