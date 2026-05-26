# Friends Relationship PostgreSQL Adapter Gate

状态：Accepted v0.1
最后更新：2026-05-26
范围：为未来实现 `runtime/internal/modules/friends.Repository` 的 PostgreSQL adapter 定义 gate-only boundary
依赖：`runtime/internal/modules/friends/repository.go`、`docs/friends-relationship-repository-boundary.md`、`runtime/migrations/postgres/000007_create_friend_relationships.sql`、`docs/postgresql-persistence-boundary.md`、`docs/reference-game-server-alignment.md`
Canonical decision：`ADR-0144`

配套英文源文档是 `docs/friends-relationship-postgresql-adapter-gate.md`。英文文件是 authoritative。

本文定义 friends relationship PostgreSQL adapter gate。它是 gate artifact，不添加 PostgreSQL adapter implementation、SQL execution behavior、unit-of-work factory wiring、runtime friendship behavior、protocol routes、Protobuf source、generated output、dependencies、migrations、automatic startup migration behavior、event/audit tables、chat rooms、groups、parties、broadcast fanout、matchmaking、match runtime、operations/admin behavior、SDK publication、generated client libraries、hosted deployments、release artifacts、public announcements、paid promotion、Pitaya-style distributed architecture 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Friends relationship PostgreSQL adapter gate 记录如下：

```yaml
friends_relationship_postgresql_adapter_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0236
decision: ADR-0144
check_rule: runtime.friends_relationship_postgresql_adapter_gate
source_repository_interface_decision: ADR-0143
repository_interface: runtime/internal/modules/friends.Repository
repository_interface_source: runtime/internal/modules/friends/repository.go
repository_tests: runtime/internal/modules/friends/repository_test.go
source_migration_source_decision: ADR-0141
source_migration_source: runtime/migrations/postgres/000007_create_friend_relationships.sql
friend_relationships_logical_table: friend_relationships
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
future_adapter_source_candidate: runtime/internal/platform/persistence/postgres/friend_relationship_repository.go
future_adapter_tests_candidate: runtime/internal/platform/persistence/postgres/friend_relationship_repository_test.go
future_constructor_candidate: NewFriendRelationshipRepositoryForUnitOfWork
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
event_audit_table_added: false
chat_added: false
groups_added: false
parties_added: false
matchmaking_added: false
match_runtime_added: false
sdk_added: false
generated_client_library_added: false
hosted_deployment_added: false
release_artifact_added: false
distributed_runtime_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_adapter_implementation_work_item: W-0237
future_adapter_implementation_direction: implement_friends_relationship_postgresql_adapter
```

## 2. Purpose

`W-0235` 已实现 storage-neutral `runtime/internal/modules/friends.Repository` interface。下一步有价值的边界，是先定义 platform adapter gate，让后续实现能把 interface 映射到既有 PostgreSQL `friend_relationships` table。

本 gate 在编写任何 adapter SQL 之前记录未来 implementation shape：

- adapter ownership；
- constructor 与 executor handoff expectation；
- transaction 与 unit-of-work boundary；
- friend request、read、list、lifecycle transition、block mutation 方法的 SQL mapping posture；
- conflict、affected-row 与 driver-error mapping；
- timestamp、relationship version、canonical pair 与 block-column mapping；
- focused test expectations；
- 保持 runtime、protocol、generated output、broad social features 和 direct compatibility 不进入 adapter slice 的 stop conditions。

这不是 implementation。未来 adapter source path 只用于让后续 agents 按已接受的 boundary 校验实现。

## 3. Ownership

未来 adapter owner 是：

```yaml
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
repository_interface_owner: runtime/internal/modules/friends
sql_mapping_owner: runtime/internal/platform/persistence/postgres
transaction_owner: caller_supplied_unit_of_work
application_layer_owns_request_identity: true
friends_module_owns_repository_vocabulary: true
player_module_owns_player_lifecycle: true
websocket_transport_owns_friends_relationships: false
protocol_adapter_owns_friends_relationships: false
authentication_module_owns_friends_relationships: false
storage_module_owns_friends_relationships: false
```

规则：

- 未来 adapter 可以在 PostgreSQL platform package 下实现 `friends.Repository`。
- adapter 不得把 SQL 移入 `runtime/internal/modules/friends`。
- adapter 不拥有 request authentication、player identity validation、route policy、protocol parsing、WebSocket state、chat rooms、groups、parties、matchmaking、match runtime 或 distributed topology。
- adapter 必须接收已 normalize 的 repository input，或在 SQL mapping 前调用 friends module normalizers。
- adapter 必须返回 friends module value types 和 typed repository errors，而不是 driver-specific errors。
- adapter 可以把 player-account foreign-key outcome 映射为 storage conflict，但不能成为 player account lifecycle owner。

## 4. Future Constructor And Executor Handoff

第一版 adapter implementation 应遵循既有 PostgreSQL adapter pattern：

```yaml
future_constructor_candidate: NewFriendRelationshipRepositoryForUnitOfWork
future_repository_interface: runtime/internal/modules/friends.Repository
executor_source: caller_supplied
transaction_control_sql_allowed: false
unit_of_work_handoff_required: true
connection_pool_owned_by_adapter: false
context_required: true
```

规则：

- constructor 应接收由 unit-of-work boundary 提供的 executor 或 query interface，而不是直接拥有 pool。
- adapter 不得发出 `BEGIN`、`COMMIT` 或 `ROLLBACK`；transaction ownership 留给 unit-of-work runner。
- adapter 不得创建 automatic startup migrations。
- 如果既有 PostgreSQL platform dependencies 已覆盖实现需要，adapter 不得新增 dependency。
- 任何必要的 dependency change 都必须进入单独的 dependency-adoption decision。

## 5. SQL Mapping Posture

未来 adapter 可以把 repository methods 映射到 `friend_relationships` table：

```yaml
logical_table: friend_relationships
primary_key_column: relationship_id
pair_columns:
  - player_low_id
  - player_high_id
lifecycle_column: lifecycle_state
actor_columns:
  - requested_by_player_id
  - responded_by_player_id
  - removed_by_player_id
block_columns:
  - blocked_by_low_at
  - blocked_by_high_at
version_column: relationship_version
created_at_column: created_at
updated_at_column: updated_at
state_changed_at_column: state_changed_at
rejected_at_column: rejected_at
removed_at_column: removed_at
pair_unique_index: friend_relationships_pair_uq
player_low_state_index: friend_relationships_player_low_state_idx
player_high_state_index: friend_relationships_player_high_state_idx
updated_at_index: friend_relationships_updated_at_idx
```

Method posture：

- `CreateOrUpdateFriendRequest` 应为 canonical pair 插入 pending relationship，或仅按后续 implementation decision 更新 ended relationship。它必须保留 repository validation，尽量在 SQL 前拒绝 self-relationship，以正数 relationship version 起步，并把 active friendship、duplicate pending、blocked 和 invalid transition outcome 映射为 typed friends conflicts。
- `GetRelationshipByPair` 应按 canonical pair select，并返回 normalized `FriendRelationship`。
- `ListRelationshipsForPlayer` 必须 player-scoped、status-filtered、deterministically ordered、受 repository limit 约束，并为 pagination 做准备。它不能变成任意 social graph search 或 admin inspection。
- `AcceptFriendRequest`、`RejectFriendRequest` 和 `RemoveFriend` 应通过 expected lifecycle state 和 optional expected-version check 更新既有 row。Affected-row count 必须区分 not found、stale version 和 invalid transition，同时不泄漏 private social graph details。
- `SetPlayerBlock` 和 `ClearPlayerBlock` 应更新由 canonical pair member role 推导出的 actor-specific block timestamp column。Unblock 不得隐式恢复 friendship。

规则：

- SQL text 必须留在 PostgreSQL adapter package 内。
- Private social graph data 默认不 log-safe，不能出现在默认 error text 中。
- Driver-specific constraint names 可以在内部用于 error mapping，但 public module errors 必须保持 friends-module neutral。
- Update operation 必须检查 affected-row count。
- 除非后续 event/audit 或 retention decision 明确授权，adapter 不得 hard-delete relationship history。

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

- Application services 或 runtime composition 以后可以通过 explicit unit-of-work boundary 获取 adapter。
- 本 gate 不添加该 factory 或 composition。
- Adapter methods 必须使用 caller context。
- Adapter 不得通过返回 successful friends relationship result 来隐藏 transaction failure。
- Adapter 不得执行 route policy、session validation、access-token validation 或 WebSocket close behavior。

## 7. Error Mapping

未来 adapter 必须把 PostgreSQL details 收敛为 friends module errors：

```yaml
repository_error_owner: runtime/internal/modules/friends
driver_error_public_leakage_allowed: false
constraint_name_public_leakage_allowed: false
private_social_graph_public_leakage_allowed: false
conflict_classes:
  - relationship_not_found
  - target_player_not_found
  - self_relationship_forbidden
  - duplicate_pending_request
  - already_friends
  - blocked_relationship
  - invalid_transition
  - version_mismatch
  - stale_relationship_version
  - pair_identity_conflict
  - storage_unavailable
```

Mapping expectations：

- Pair uniqueness conflicts 应根据 current row state 和 repository operation 映射到 duplicate request、already friends、blocked relationship、invalid transition 或 pair identity conflict。
- Pair members 或 actor columns 的 foreign-key failures 应映射到 target-player-not-found 或 invalid input，不暴露 database constraint names。
- 带 expected version 的 no affected row 应映射到 relationship-not-found、version-mismatch、stale-version 或 invalid-transition，不泄漏 hidden relationship history。
- Malformed input 应尽可能在 SQL execution 前被拒绝。
- Unknown driver 或 executor failures 应映射为 storage-unavailable style errors，并使用 redacted reasons。
- Raw SQL、DSNs、credentials、token material、verifier digests、player ids、private relationship state 和 driver stack details 不得出现在 public error strings 中。

## 8. Test Expectations

后续 implementation slice 应添加 focused PostgreSQL adapter tests。在 live database verification 不可用前，可以使用 fake-executor 或 query-capture tests。

Implementation slice 必须覆盖的 test families：

```yaml
future_tests:
  - constructor_requires_executor
  - create_or_update_friend_request_maps_insert_update_and_conflicts
  - get_relationship_selects_by_canonical_pair
  - list_relationships_is_player_scoped_status_filtered_and_ordered
  - accept_reject_remove_check_expected_version_and_transition_state
  - block_unblock_update_actor_specific_block_columns
  - rows_scan_through_friends_normalizers
  - driver_errors_are_redacted
  - transaction_control_sql_is_absent
  - default_tests_do_not_require_live_postgresql
```

规则：

- Tests 不得要求 protocol routes。
- Tests 不得要求 WebSocket transport。
- Tests 不得打印 raw credentials、tokens、verifier keys、digests、DSNs、query strings、authorization headers、cookies、player ids 或 private relationship state。
- Live PostgreSQL verification 可以在后续 implementation slice 授权时添加；如果不可用，必须显式记录。

## 9. Relationship To Runtime, Protocol, And Authentication

本 gate 不改变 runtime 或 protocol behavior：

```yaml
runtime_friends_handlers_added: false
friends_protocol_routes_added: false
protobuf_friends_messages_added: false
generated_friends_output_added: false
authentication_session_behavior_changed: false
request_identity_handoff_changed: false
```

规则：

- Runtime friendship behavior 继续 deferred。
- Protocol routes 和 generated friends contract shapes 继续 deferred。
- Request identity validation 继续由 authentication/session boundaries 拥有。
- Adapter 不得解析 bearer tokens、cookies、WebSocket subprotocols、envelope metadata 或 transport connection identifiers。

## 10. Reference Alignment

Nakama 提供 durable friends relationship social graph behavior 的 product capability pressure。Pitaya 继续作为 deferred future architecture reference，用于分布式 runtime concern。vibit 只将这些 reference 用于 capability planning：

- 本 gate 不添加 direct Nakama 或 Pitaya API compatibility；
- 本 gate 不添加 public friends route；
- 本 gate 不添加 server runtime hook、group、party、chat、matchmaking、match runtime 或 admin surface；
- 未来 adapter 是 module/application behavior 之下的 platform persistence detail。

## 11. Stop Conditions

做以下任何事情前必须停止，并打开后续 bounded work item：

- 实现 `runtime/internal/platform/persistence/postgres/friend_relationship_repository.go`；
- 添加 friends relationships 的 SQL execution behavior；
- 添加 unit-of-work factory wiring 或 startup composition；
- 添加 runtime friend request/list/status behavior；
- 添加 protocol routes、Protobuf sources、generated output 或 generated clients；
- 修改 authentication/session behavior 或 request identity validation；
- 添加 event/audit tables；
- 添加 chat、groups、parties、broadcast fanout、matchmaking、match runtime、SDK、hosted、release 或 distributed runtime scope；
- 添加 direct Nakama 或 Pitaya public API compatibility。

