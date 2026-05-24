# Friends Relationship Persistence Schema Gate 中文版

状态：Accepted v0.1
最后更新：2026-05-24
范围：在添加 migration source、repository interface、adapter、runtime behavior、protocol route、generated output 或更广 social feature 前，定义未来 PostgreSQL friends relationship persistence schema 的 gate
依赖：`docs/friends-relationship-lifecycle-gate.md`、`docs/postgresql-persistence-boundary.md`、`docs/reference-game-server-alignment.md`、`docs/nakama-pitaya-product-parity-roadmap.md`
权威决策：`ADR-0140`

本文件是 `docs/friends-relationship-persistence-schema-gate.md` 的简体中文译本。英文版本是权威版本。

本文定义 friends relationship persistence schema gate。它是 gate artifact。它不添加 SQL migration source，不创建 `friend_relationships` table，不实现 friendship runtime behavior，不添加 protocol routes，不添加 Protobuf source 或 generated output，不添加 dependencies，不添加 repository interfaces，不添加 PostgreSQL adapters，不接入 startup，不改变 authentication/session behavior，不添加 delivery guarantees，不添加 stream subscriptions，不添加 chat rooms、groups、parties、broadcast fanout、matchmaking、match runtime、operations/admin behavior，不发布 SDK 或 generated client libraries，不创建 hosted deployments 或 release artifacts，不添加 Pitaya-style distributed architecture，也不添加 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

friends relationship persistence schema gate 记录如下：

```yaml
friends_relationship_persistence_schema_gate: defined
completed_work_item: W-0232
decision: ADR-0140
check_rule: runtime.friends_relationship_persistence_schema_gate
source_lifecycle_gate_decision: ADR-0139
source_lifecycle_gate_standard: docs/friends-relationship-lifecycle-gate.md
gate_standard: docs/friends-relationship-persistence-schema-gate.md
gate_standard_translation: docs/friends-relationship-persistence-schema-gate.zh-CN.md
selected_nakama_capability_family: friends_groups_and_parties
primary_product_reference: Nakama
pitaya_reference_status: deferred_future_architecture_reference
selected_first_friends_relationship_store: postgres
future_friends_relationships_logical_table: friend_relationships
future_friend_relationship_events_logical_table: deferred
future_migration_source_candidate: runtime/migrations/postgres/000007_create_friend_relationships.sql
future_repository_owner_candidate: runtime/internal/modules/friends
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
pair_identity_posture_recorded: true
relationship_state_representation_recorded: true
block_representation_recorded: true
index_uniqueness_posture_recorded: true
timestamp_posture_recorded: true
event_audit_posture_recorded: true
redaction_posture_recorded: true
future_repository_adapter_boundaries_recorded: true
future_migration_source_candidate_recorded: true
schema_gate_only: true
migration_source_added: false
friend_relationships_table_added: false
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
dependency_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_migration_work_item: W-0233
future_migration_direction: add_friends_relationship_migration_source
```

## 2. Product Intent

`ADR-0139` 已定义 Nakama-class social graph 的未来 friends relationship lifecycle。下一步是在添加 SQL 前，把未来 durable state 显式化。

本 gate 让未来 migration 可检查：

- table candidate 已知；
- pair identity 和 canonical ordering 在 SQL 前确定；
- relationship state 与 actor-relative public status 分开表示；
- block state 有独立表示，unblock 不会恢复旧 friendship；
- indexes 和 uniqueness 在 implementation 前规划；
- event/audit posture 在添加 event storage 前明确；
- redaction 和未来 repository/adapter boundaries 显式。

本 gate 保持保守。它准备下一项 migration-source-only slice，但不添加 migration file。

## 3. Selected Store And Table

第一版 friends relationship persistence target 是 PostgreSQL：

```yaml
selected_first_friends_relationship_store: postgres
future_friends_relationships_logical_table: friend_relationships
future_migration_source_candidate: runtime/migrations/postgres/000007_create_friend_relationships.sql
future_repository_boundary: separate_future_work_item
future_postgresql_adapter: separate_future_work_item
```

理由：

- PostgreSQL 是 vibit 第一版 accepted authoritative durable store。
- Friends relationship state 是 durable social graph state，必须可在事务中检查。
- 第一版实现不应在本地 schema 证明前引入 graph database、document database、cache dependency 或 distributed social graph subsystem。
- Nakama 是 product capability reference；vibit 使用自己的 schema 和 contract posture。

## 4. Future `friend_relationships` Table Candidate

未来第一版 migration 可以定义一个 logical current-state table：

```yaml
friend_relationships:
  primary_key_candidate:
    - relationship_id
  required_columns:
    - relationship_id
    - player_low_id
    - player_high_id
    - lifecycle_state
    - relationship_version
    - created_at
    - updated_at
    - state_changed_at
  nullable_columns:
    - requested_by_player_id
    - responded_by_player_id
    - removed_by_player_id
    - rejected_at
    - removed_at
    - blocked_by_low_at
    - blocked_by_high_at
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
    - channel_id
    - chat_room_id
    - group_id
    - party_id
    - match_id
    - pitaya_server_id
    - nakama_api_path
```

`relationship_id` 是 server-generated opaque record id，不是 relationship 的 public identity。逻辑身份是 unordered player pair：

```text
player_low_id + player_high_id
```

`player_low_id` 和 `player_high_id` 是 canonical pair members。未来 migration source 应防止 self-target 并强制 canonical ordering，例如使用等价于以下规则的 check：

```text
player_low_id < player_high_id
```

schema 是 pair-oriented。Public query output 仍由后续 application behavior 相对 actor 计算。

## 5. Pair Identity And Player References

第一版 pair identity posture 是：

```yaml
pair_identity: canonical_unordered_player_pair
player_low_id_source: canonicalized_player_pair_member
player_high_id_source: canonicalized_player_pair_member
player_fk_candidate: player_accounts(player_id)
self_relationship_allowed: false
client_supplied_actor_id_as_proof_allowed: false
metadata_only_player_id_allowed_as_proof: false
```

规则：

- relationship row 不得表示 player 自己指向自己。
- 两个 pair member 都应引用已存在的 player account records。
- 后续 runtime behavior 必须从 validated request identity 得到 actor，而不是信任 client-supplied actor ids。
- pair columns 只是 persistence identity，不是 authentication proof。
- table 不得存储 session ids 或 transport connection identifiers。

## 6. Relationship State Representation

第一版 lifecycle state candidate 是：

```yaml
lifecycle_state_column: lifecycle_state
lifecycle_state_type: TEXT
allowed_lifecycle_states:
  - pending
  - friends
  - rejected
  - removed
actor_relative_public_status_stored: false
```

规则：

- `pending` 记录等待响应的 request。
- `friends` 记录已接受的 relationship。
- `rejected` 在第一版 schema candidate 中把被拒绝的 request 记录为当前 pair state。
- `removed` 记录已结束关系或为 block posture 保留的 neutral row。
- `outgoing_request_pending`、`incoming_request_pending`、`blocked_by_actor`、`blocked_actor` 等 actor-relative public states 由未来 query behavior 计算，不作为 canonical database states 存储。
- Duplicate request idempotency 仍是未来 runtime behavior decision；schema 只保证最多有一个 current pair row。

第一版 schema candidate 把 `rejected` 和 `removed` 存为 current row states，而不是 audit-only facts。Retention、cleanup、hard-delete behavior 继续推迟到后续 gate。

## 7. Request And Response Actor Columns

第一版 request/response posture 是：

```yaml
requested_by_player_id: nullable_pair_member
responded_by_player_id: nullable_pair_member
removed_by_player_id: nullable_pair_member
```

规则：

- 当 state 源自 request 时，`pending`、`friends` 和 `rejected` states 应有 `requested_by_player_id`。
- accepted 或 rejected states 可以有 `responded_by_player_id`。
- removed states 可以有 `removed_by_player_id`。
- 这些 columns 是 state history 和 conflict handling 使用的 pair-member references；它们不是 authentication proof。
- 当 privacy 需要 collapse 时，public errors 和 logs 不得暴露 hidden relationship history。

## 8. Block Representation

Block state 是 actor-specific，并独立于 lifecycle state：

```yaml
block_columns:
  - blocked_by_low_at
  - blocked_by_high_at
block_representation: per_pair_member_timestamp
mutual_block_representation: both_block_columns_present
unblock_restores_prior_friendship: false
```

规则：

- `blocked_by_low_at` 表示 `player_low_id` block 了 `player_high_id`。
- `blocked_by_high_at` 表示 `player_high_id` block 了 `player_low_id`。
- 两者都存在时，public actor-relative status 是 mutual block。
- 后续 behavior 中，block 必须 override pending 或 friends state。
- Unblock 清除 actor 自己的 block timestamp，且不得自动恢复之前的 friendship。
- block-only row 在未来 behavior gate 选择其他表示前，可以使用 `lifecycle_state: removed`。

## 9. Version, Timestamp, And Retention Posture

第一版 version 和 timestamp posture 是：

```yaml
relationship_version_column: relationship_version
relationship_version_type_candidate: BIGINT
initial_relationship_version_candidate: 1
created_at: TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at: TIMESTAMPTZ NOT NULL DEFAULT now()
state_changed_at: TIMESTAMPTZ NOT NULL DEFAULT now()
soft_delete_column: deferred
hard_delete_policy: deferred
```

规则：

- `relationship_version` 必须为正数，并由 server 管理。
- 后续 command behavior 应在成功 state mutation 后递增 version。
- 应强制 `updated_at >= created_at`。
- 应强制 `state_changed_at >= created_at`。
- `rejected_at`、`removed_at`、`blocked_by_low_at` 和 `blocked_by_high_at` 如果存在，不应早于 `created_at`。
- Cleanup、retention windows、hard delete 和 tombstone pruning 均推迟。

## 10. Uniqueness And Index Posture

第一版 uniqueness posture 是：

```yaml
logical_pair_unique_candidate:
  - player_low_id
  - player_high_id
```

未来 migration-source slice 推荐 indexes：

- `(player_low_id, player_high_id)` unique pair identity index；
- `(player_low_id, lifecycle_state)` lookup index；
- `(player_high_id, lifecycle_state)` lookup index；
- updated-at index，用于 future diagnostics 或 cleanup；
- block indexes 只有在 migration-source slice 能保持 narrow 且有明确理由时才添加。

本 gate 不授权 global player search、relationship recommendations、social graph traversal、analytics indexes、admin dashboards、chat targeting indexes、group 或 party indexes、matchmaking indexes，或 distributed graph routing。

## 11. Event And Audit Posture

第一版 migration source candidate 只应添加 current-state table：

```yaml
future_friend_relationship_events_logical_table: deferred
outbox_table_added_by_schema_gate: false
audit_table_added_by_schema_gate: false
domain_events_defined_by_lifecycle_gate: true
```

理由：

- lifecycle gate 已定义未来 domain events；
- current-state table 足以支撑第一版 migration source；
- event history、audit retention、outbox delivery 和 analytics 应拆到后续 bounded gate。

未来 runtime behavior 一旦被授权，仍必须在 unit-of-work boundary 内保持 emitted events 与 state changes 一致。

## 12. Ownership Boundaries

未来 friends relationship behavior 应拥有自己的 module boundary：

```yaml
future_repository_owner_candidate: runtime/internal/modules/friends
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
future_contract_owner_candidate: contracts/friends
future_proto_source_candidate: proto/vibit/friends/v1
friends_module_owns_friend_relationships: true
storage_module_owns_friend_relationships: false
player_module_owns_friend_relationships: false
authentication_module_owns_friend_relationships: false
websocket_transport_owns_friend_relationships: false
```

规则：

- 未来 friends module 拥有 relationship lifecycle domain behavior。
- player module 拥有 player account lifecycle，不拥有 social graph transitions。
- authentication 拥有 credentials、tokens 和 sessions，不拥有 friendship state。
- storage objects 拥有 generic player-owned JSON objects，不拥有 social graph relationships。
- WebSocket transport 拥有 connection plumbing，不拥有 durable relationship records。
- 只有 repository boundaries 被授权后，PostgreSQL adapters 才能实现 friends repositories。

## 13. Redaction Posture

Friends relationship records 默认不 log-safe。

默认不 log-safe：

- relationship ids；
- 组成 social graph records 的 pair member ids；
- lifecycle state；
- request、response、removal、rejection 或 block actor ids；
- block timestamps；
- conflict details；
- 暴露 pair identity 或 private relationship history 的 database errors。

禁止的 secret 和 transport material：

- raw device credentials；
- raw access tokens；
- verifier keys；
- lookup 或 verifier digests；
- 带 credentials 的 PostgreSQL DSNs；
- headers、cookies、query strings、WebSocket subprotocol values、remote addresses，或具体 transport metadata。

未来 adapters 和 handlers 必须返回 redacted errors；当 privacy 需要 collapse 时，public failures 不得泄漏 hidden relationship details。

## 14. Future Migration Source Expectations

下一项 bounded work item 可以添加：

```text
runtime/migrations/postgres/000007_create_friend_relationships.sql
```

该 migration-source-only slice 可以为 `friend_relationships` 添加 SQL DDL、comments、indexes 和 migration checks。

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
- chat、groups、parties、matchmaking、match runtime 或 operations/admin behavior；
- Pitaya-style distributed architecture；
- direct Nakama/Pitaya API compatibility。

## 15. Verification Expectations

未来 migration-source slice 应验证：

- goose up/down markers；
- table name 和 required columns；
- pair ordering 和 self-target prevention checks；
- lifecycle state vocabulary checks；
- relationship version 和 timestamp checks；
- pair uniqueness 和 list-query indexes；
- 没有 forbidden secret、digest、transport、chat、group、party、match、Pitaya 或 Nakama compatibility columns；
- 没有 Go runtime behavior；
- migration boundary 的 repository checks。

后续 repository/adapter/runtime work 只有在对应 implementation gates 被接受后，才应添加 request、accept、reject、remove、block、unblock、list、status、privacy、redaction 和 concurrency behavior 的 focused tests。

## 16. Stop Conditions

做以下任何事项前必须停止并请求 maintainer 授权：

- 在本 gate 同一 change 中添加 SQL migration source；
- 创建 `friend_relationships`；
- 实现 friends relationship runtime behavior；
- 添加 protocol routes；
- 添加 Protobuf source files 或 generated output；
- 添加 repository interfaces；
- 添加 PostgreSQL 或其他 storage adapters；
- 添加 dependencies；
- 改变 authentication/session semantics；
- 改变 route protection semantics；
- 添加 chat rooms、groups、parties、matchmaking、match runtime、operations/admin behavior、social graph search、recommendation、analytics 或 distributed graph routing；
- 添加 hosted deployments 或 demos；
- 创建 release binaries、packages、containers、checksums、signing/provenance artifacts、install scripts、registry publications、SDK packages 或 additional release artifacts；
- 执行 GitHub release record 之外的 public announcements；
- 运行 paid promotion；
- 添加 direct Nakama/Pitaya API compatibility。

## 17. Next Work

下一项 bounded direction 是：

```text
W-0233 Add friends relationship migration source
```

该工作可以添加 `friend_relationships` 的第一版 SQL migration source 和对应 static checks，同时继续延后 repository interfaces、adapters、protocol、generated output 和 runtime behavior。
