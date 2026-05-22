# Storage Objects PostgreSQL Adapter Gate

状态：Accepted v0.1
最后更新：2026-05-22
范围：未来实现 `runtime/internal/modules/storage.Repository` 的 PostgreSQL adapter 的 gate-only 边界
依赖：`runtime/internal/modules/storage/repository.go`、`docs/storage-objects-repository-boundary.md`、`runtime/migrations/postgres/000006_create_storage_objects.sql`、`docs/postgresql-persistence-boundary.md`、`docs/reference-game-server-alignment.md`
规范决策：`ADR-0114`

配对英文源文件是 `docs/storage-objects-postgresql-adapter-gate.md`。英文文件是权威版本。

本文定义 storage objects PostgreSQL adapter gate。它是 gate artifact，不添加 PostgreSQL adapter implementation、SQL execution behavior、unit-of-work factory wiring、runtime behavior、protocol routes、Protobuf source、generated output、dependencies、migrations、automatic startup migration behavior、broad operations/admin behavior、authentication/session behavior changes、hosted deployments、release artifacts、public announcements、paid promotion、broad product module expansion、large object/blob storage、S3-compatible object storage 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Storage objects PostgreSQL adapter gate 记录为：

```yaml
storage_objects_postgresql_adapter_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0206
decision: ADR-0114
check_rule: runtime.storage_objects_postgresql_adapter_gate
source_repository_interface_decision: ADR-0113
repository_interface: runtime/internal/modules/storage.Repository
repository_interface_source: runtime/internal/modules/storage/repository.go
source_migration_source_decision: ADR-0111
source_migration_source: runtime/migrations/postgres/000006_create_storage_objects.sql
storage_objects_logical_table: storage_objects
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
future_adapter_source_candidate: runtime/internal/platform/persistence/postgres/storage_object_repository.go
future_adapter_tests_candidate: runtime/internal/platform/persistence/postgres/storage_object_repository_test.go
unit_of_work_handoff_required: true
transaction_owner: caller_supplied_unit_of_work
sql_mapping_owner: postgresql_adapter
adapter_gate_only: true
postgresql_adapter_added: false
sql_execution_added: false
unit_of_work_factory_added: false
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
future_adapter_implementation_work_item: W-0207
future_adapter_implementation_direction: storage_objects_postgresql_adapter_implementation
```

## 2. Purpose

`W-0205` 已实现 storage-neutral `runtime/internal/modules/storage.Repository` interface。下一步需要先定义 platform adapter gate，再让后续实现把该 interface 映射到已有 PostgreSQL `storage_objects` table。

该 gate 在写 SQL 之前记录未来实现形状：

- adapter ownership；
- constructor 与 executor handoff expectations；
- transaction 与 unit-of-work boundaries；
- create/read/list/update/delete 的 SQL mapping posture；
- conflict 与 driver-error mapping；
- timestamp、JSON、version、soft-delete mapping；
- focused test expectations；
- stop conditions，避免 runtime、protocol、generated output 和 product expansion 泄入 adapter slice。

这不是实现。未来 adapter source path 只是命名，便于后续 work 按已接受边界验证。

## 3. Ownership

未来 adapter owner 是：

```yaml
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
repository_interface_owner: runtime/internal/modules/storage
sql_mapping_owner: runtime/internal/platform/persistence/postgres
transaction_owner: caller_supplied_unit_of_work
application_layer_owns_request_identity: true
storage_module_owns_repository_vocabulary: true
websocket_transport_owns_storage_objects: false
protocol_adapter_owns_storage_objects: false
authentication_module_owns_storage_objects: false
player_module_owns_storage_objects: false
```

规则：

- Adapter 后续可以在 PostgreSQL platform package 下实现 `storage.Repository`。
- Adapter 不得把 SQL 移入 `runtime/internal/modules/storage`。
- Adapter 不拥有 request authentication、player identity validation、route policy、protocol parsing 或 WebSocket state。
- Adapter 应接收已 normalize 的 repository input，或在 SQL mapping 前调用 storage module normalizers。
- Adapter 必须返回 storage module value types 和 typed repository errors，而不是 driver-specific errors。

## 4. Future Constructor And Executor Handoff

第一版 adapter implementation 应遵循已有 PostgreSQL adapter 模式：

```yaml
future_constructor_candidate: NewStorageObjectRepository
future_repository_interface: runtime/internal/modules/storage.Repository
executor_source: caller_supplied
transaction_control_sql_allowed: false
unit_of_work_handoff_required: true
connection_pool_owned_by_adapter: false
context_required: true
```

规则：

- Constructor 应接收既有 executor 或 query interface，而不是直接拥有 pool。
- Adapter 不得执行 `BEGIN`、`COMMIT` 或 `ROLLBACK`；transaction ownership 属于 unit-of-work runner。
- Adapter 不得添加 automatic startup migrations。
- 如果已有 PostgreSQL platform dependency 足够，adapter 不得新增依赖。
- 任何依赖变化都必须走单独 dependency-adoption decision。

## 5. SQL Mapping Posture

未来 adapter 可以把 repository methods 映射到 `storage_objects` table：

```yaml
logical_table: storage_objects
primary_key_column: object_id
owner_columns:
  - owner_kind
  - owner_id
identity_columns:
  - collection
  - object_key
value_column: value_json
version_column: version
created_at_column: created_at
updated_at_column: updated_at
soft_delete_column: deleted_at
active_identity_unique_index: storage_objects_active_identity_uq
owner_collection_index: storage_objects_owner_collection_idx
updated_at_index: storage_objects_updated_at_idx
```

Method posture：

- `CreateStorageObject` 应在 validation 后插入 active rows，保留 `object_id`，从 positive version 开始，并把 active identity conflicts 映射为 `object_already_exists`。
- `GetStorageObject` 只返回 requested owner、collection、key 对应的 active rows。
- `ListStorageObjects` 必须 owner and collection scoped，使用 deterministic key ordering，遵守 bounded limits，并在第一版中使用 `AfterObjectKey` 作为足够不透明的 pagination cursor。
- `UpdateStorageObject` 应更新 active rows，server-side increment version，保留 ownership，并在 supplied expected version 时执行 expected-version check。
- `DeleteStorageObject` 应 soft-delete active rows，在 supplied expected version 时 version-aware，并保留 owner-mismatch leakage rules。

规则：

- SQL text 必须留在 PostgreSQL adapter package。
- `value_json` 默认不是 log-safe，不能放入默认 error text。
- Driver-specific constraint names 可以用于内部 error mapping，但 public module errors 必须保持 storage-neutral。
- Update 和 delete operations 必须检查 affected-row counts。

## 6. Transaction And Unit-Of-Work Boundary

未来 adapter 参与既有 transaction handoff：

```yaml
unit_of_work_handoff_required: true
adapter_starts_transactions: false
adapter_commits_transactions: false
adapter_rolls_back_transactions: false
adapter_safe_for_existing_runner: true
```

规则：

- Application services 或 runtime composition 后续可通过显式 unit-of-work boundary 获取 adapter。
- 本 gate 不添加该 factory 或 composition。
- Adapter methods 必须使用 caller 的 context。
- Adapter 不得把 transaction failures 隐藏成成功的 storage results。

## 7. Error Mapping

未来 adapter 必须把 PostgreSQL details 折叠为 storage module errors：

```yaml
repository_error_owner: runtime/internal/modules/storage
driver_error_public_leakage_allowed: false
constraint_name_public_leakage_allowed: false
value_json_public_leakage_allowed: false
conflict_classes:
  - object_already_exists
  - object_not_found
  - version_mismatch
  - invalid_expected_version
  - owner_scope_mismatch
  - deleted_object
  - storage_unavailable
```

Mapping expectations：

- Unique active identity conflicts 应映射为 `object_already_exists`。
- 有 expected version 但没有 affected row 时，应映射为 typed not-found 或 version conflict，且不得泄露其他 owner 的 object existence。
- Malformed input 应尽可能在 SQL execution 之前拒绝。
- Unknown driver 或 executor failures 应映射为带 redacted reasons 的 storage-unavailable 类错误。
- Raw SQL、DSNs、credentials、token material、verifier digests、`value_json` 和 driver stack details 不得出现在 public error strings。

## 8. Test Expectations

后续 implementation slice 应添加 focused PostgreSQL adapter tests。在 live database verification 不可用前，可以使用 fake-executor 或 query-capture tests。

Implementation slice 的 required test families：

```yaml
future_tests:
  - constructor_requires_executor
  - create_maps_insert_and_unique_conflict
  - get_maps_active_lookup_and_not_found
  - list_is_owner_collection_scoped_and_ordered
  - update_checks_expected_version_and_increments_version
  - delete_soft_deletes_and_checks_expected_version
  - json_value_is_copied_and_redacted
  - timestamps_are_utc_normalized
  - driver_errors_are_redacted
  - transaction_control_sql_is_absent
```

规则：

- Tests 不得依赖 protocol routes。
- Tests 不得依赖 WebSocket transport。
- Tests 不得打印 raw values、DSNs、credentials、tokens 或 digests。
- Live PostgreSQL verification 可在 implementation slice 授权后添加；如果不可用，必须明确记录。

## 9. Relationship To Runtime, Protocol, And Authentication

该 gate 不改变 runtime 或 protocol behavior：

```yaml
runtime_storage_handlers_added: false
storage_protocol_routes_added: false
protobuf_storage_messages_added: false
generated_storage_output_added: false
authentication_session_behavior_changed: false
request_identity_handoff_changed: false
```

规则：

- Runtime handlers 仍然 deferred。
- Protocol routes 和 generated storage contract shapes 仍然 deferred。
- Request identity validation 仍由 authentication/session boundaries 拥有。
- Adapter 不得解析 bearer tokens、cookies、WebSocket subprotocols 或 envelope metadata。

## 10. Reference Alignment

Nakama 提供 durable storage-object-like game state capability。Pitaya 强调把 persistence 放在 handlers 与 runtime routing 下层。vibit 只把它们用于 capability planning：

- 本 gate 不添加 direct Nakama 或 Pitaya API compatibility；
- 本 gate 不添加 public storage route；
- 本 gate 不添加 server runtime hook 或 admin surface。

## 11. Verification

该 gate 后需要运行：

```text
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.storage_objects_postgresql_adapter_gate
node tools/vibit check change define-storage-objects-postgresql-adapter-gate --json
node tools/vibit check module storage --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
cd runtime && go test ./...
git diff --check
```

本 gate 不需要 live PostgreSQL adapter verification，因为没有添加 adapter implementation 或 SQL execution。

## 12. Stop Conditions

添加下列内容前必须停止并创建后续 bounded work item：

- `runtime/internal/platform/persistence/postgres/storage_object_repository.go`；
- PostgreSQL storage adapter implementation；
- SQL execution behavior；
- unit-of-work factory wiring；
- runtime storage object handlers；
- protocol routes；
- Protobuf source 或 generated output；
- dependency changes；
- migration changes；
- automatic startup migration behavior；
- authentication/session behavior changes；
- route-protection changes；
- broad operations/admin behavior；
- hosted deployments；
- release binaries、packages、containers、checksums、provenance files、signing artifacts、install scripts、registry publication 或 SDK packages；
- public announcements beyond the GitHub release record；
- paid promotion；
- large object/blob storage；
- S3-compatible object storage；
- direct Nakama/Pitaya API compatibility。
