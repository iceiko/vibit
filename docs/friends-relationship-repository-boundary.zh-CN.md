# Friends Relationship Repository Boundary 中文版

状态：Accepted v0.1
最后更新：2026-05-25
范围：在 PostgreSQL `friend_relationships` migration source 之后，为未来 storage-neutral friends relationship repository 定义 gate-only boundary
依赖：`docs/friends-relationship-lifecycle-gate.md`、`docs/friends-relationship-persistence-schema-gate.md`、`decisions/ADR-0141-friends-relationship-migration-source.md`、`docs/reference-game-server-alignment.md`
权威决策：`ADR-0142`

说明：本文件是 `docs/friends-relationship-repository-boundary.md` 的简体中文译本。英文文件是权威版本。

本文定义 friends relationship repository boundary。它是 gate artifact。它不添加 Go repository interfaces、PostgreSQL adapter behavior、runtime friendship behavior、protocol routes、Protobuf source、generated output、dependencies、migrations、automatic startup migration behavior、event/audit tables、chat rooms、groups、parties、broadcast fanout、matchmaking、match runtime、operations/admin behavior、SDK publication、generated client libraries、hosted deployments、release artifacts、public announcements、paid promotion、Pitaya-style distributed architecture 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Friends relationship repository boundary 记录为：

```yaml
friends_relationship_repository_boundary: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0234
decision: ADR-0142
check_rule: runtime.friends_relationship_repository_boundary
source_migration_source_decision: ADR-0141
source_migration_source: runtime/migrations/postgres/000007_create_friend_relationships.sql
source_schema_gate_decision: ADR-0140
source_lifecycle_gate_decision: ADR-0139
future_repository_owner_candidate: runtime/internal/modules/friends
future_repository_interface_candidate: runtime/internal/modules/friends.Repository
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
friend_relationships_logical_table: friend_relationships
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
pitaya_distributed_architecture_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_repository_interface_work_item: W-0235
future_repository_interface_direction: implement_friends_relationship_repository_interface
```

## 2. 目的

`W-0233` 已添加 `friend_relationships` 的 PostgreSQL migration source。下一步有价值的边界，是先定义未来实现代码可以使用的 storage-neutral repository vocabulary，而不是把 SQL、transport 或 protocol 假设泄漏到模块里。

这个 boundary 为 Nakama-class friends relationship path 记录：

- repository ownership；
- candidate value types；
- lifecycle command 和 query vocabulary；
- pair identity 和 actor handoff rules；
- version、conflict 和 transaction handoff posture；
- redaction 和 error posture；
- PostgreSQL adapter expectations；
- 后续 implementation work 的 stop conditions。

这仍然不是 runtime feature。除非后续 bounded work item 明确授权，否则 handler、route、adapter、repository interface 或 protocol message 都不能使用 friends relationships。

## 3. Ownership

未来 repository 由 friends module 拥有：

```yaml
future_repository_owner_candidate: runtime/internal/modules/friends
future_repository_interface_candidate: runtime/internal/modules/friends.Repository
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
friend_relationships_table_owner: runtime.friends
application_layer_owns_request_identity: true
postgresql_adapter_owns_sql_mapping: true
websocket_transport_owns_friends_relationships: false
protocol_adapter_owns_friends_relationships: false
authentication_module_owns_friends_relationships: false
player_module_owns_friends_relationships: false
storage_module_owns_friends_relationships: false
```

规则：

- 未来 repository interface 必须是 storage-neutral 且面向 module。
- Interface 不得提到 PostgreSQL、pgx、SQL rows、goose migrations、prepared statements、connection pools、transaction runners 或 database driver errors。
- PostgreSQL adapter 以后可以在 `runtime/internal/platform/persistence/postgres` 下实现这个 interface，但必须等单独 adapter gate 授权。
- Application 或 handler code 以后应通过 module/application boundaries 调用 friends relationship behavior，而不是直接使用 SQL 或 transport state。
- Authentication 和 session code 提供 validated request identity，不拥有 friends relationship records。
- Player account storage 拥有 player lifecycle state，不拥有 social graph relationship state。
- Storage objects 拥有 player-owned small JSON game state，不拥有 friends relationships。
- WebSocket transport 只拥有 connection plumbing，不拥有 friendship state。
- Protocol adapters 只拥有 wire conversion，不拥有 repository behavior。

## 4. Candidate Value Types

后续 implementation gate 可以重命名或缩减这些 shape，但第一版 repository interface implementation 应从以下 vocabulary 开始：

```yaml
candidate_value_types:
  - FriendRelationship
  - FriendRelationshipPair
  - FriendRelationshipID
  - FriendRelationshipStatus
  - FriendRelationshipLifecycleState
  - FriendRelationshipActor
  - FriendRelationshipVersion
  - FriendRelationshipBlockState
  - SendFriendRequestInput
  - AcceptFriendRequestInput
  - RejectFriendRequestInput
  - RemoveFriendInput
  - BlockPlayerInput
  - UnblockPlayerInput
  - GetFriendRelationshipInput
  - ListFriendRelationshipsInput
  - FriendRelationshipConflict
  - FriendRelationshipRepositoryError
```

第一版 record vocabulary：

```yaml
friend_relationship_record:
  relationship_id: server_generated_opaque_id
  pair: canonical_unordered_player_pair
  player_low_id: canonical_pair_member
  player_high_id: canonical_pair_member
  lifecycle_state: pending_or_friends_or_rejected_or_removed
  requested_by_player_id: nullable_pair_member_actor
  responded_by_player_id: nullable_pair_member_actor
  removed_by_player_id: nullable_pair_member_actor
  blocked_by_low_at: nullable_server_timestamp
  blocked_by_high_at: nullable_server_timestamp
  relationship_version: server_managed_bigint_revision
  created_at: server_timestamp
  updated_at: server_timestamp
  state_changed_at: server_timestamp
  rejected_at: nullable_server_timestamp
  removed_at: nullable_server_timestamp
```

规则：

- Pair identity 必须是 canonical unordered player pair。
- `player_low_id` 和 `player_high_id` 是 persistence identity fields，不是 authentication proof。
- `requested_by_player_id`、`responded_by_player_id` 和 `removed_by_player_id` 是 normalized pair-member actors，不是 proof。
- 未来 runtime behavior 必须先从 validated request identity 派生 actor，再调用 repository。
- Public actor-relative states，例如 outgoing pending、incoming pending、blocked by actor 和 actor blocked target，后续再计算；它们不是 repository lifecycle states。
- `relationship_version` 由 server 管理，不能成为 client-authoritative state。
- Private social graph data 默认不是 log-safe。

## 5. Candidate Repository Capabilities

第一版 storage-neutral capability family 是：

```yaml
candidate_repository_capabilities:
  - CreateOrUpdateFriendRequest
  - GetRelationshipByPair
  - ListRelationshipsForPlayer
  - AcceptFriendRequest
  - RejectFriendRequest
  - RemoveFriend
  - SetPlayerBlock
  - ClearPlayerBlock
```

Capability 规则：

- `CreateOrUpdateFriendRequest` 只能为已经 validated 的 actor 和 normalized target pair 创建或更新 relationship row。
- `GetRelationshipByPair` 是 storage lookup。它不得 authenticate users、validate access tokens 或创建 request identity。
- `ListRelationshipsForPlayer` 必须 player-scoped 且 pagination-ready。没有后续 gate 时，不得变成任意 social graph search 或 admin inspection。
- `AcceptFriendRequest` 和 `RejectFriendRequest` 必须是 existing pending relationship 上的 lifecycle transitions。
- `RemoveFriend` 必须按后续 behavior rules 结束 active friendship 或 pending relationship；默认不得 hard-delete 可能需要审计的历史。
- `SetPlayerBlock` 和 `ClearPlayerBlock` 必须保留 actor-specific block semantics，且 unblock 后不得隐式恢复 friendship。
- 所有 methods 必须返回 typed module-owned records 和 errors，而不是 raw SQL rows 或 database driver errors。

未来 repository interface 可以选择更短名称，但必须保留 request creation、relationship read/list、lifecycle transition、block mutation 和 conflict handling 之间的语义拆分。

## 6. Pair Identity And Request Identity Handoff

Repository boundary 只准备 canonical pair handling，不实现 behavior：

```yaml
pair_identity: canonical_unordered_player_pair
self_relationship_allowed: false
actor_identity_source: validated_request_identity_before_repository_call
client_supplied_actor_id_as_proof_allowed: false
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
actor_relative_public_status_stored: false
```

规则：

- Repository 可以接收 normalized actor ids 作为数据，但 actor ids 不是 authentication proof。
- 玩家永远不能与自己形成 relationship。
- Pair canonicalization 必须是 deterministic，且不依赖请求方向。
- 当后续 behavior gate 要求 leakage collapse 时，public errors 不应泄露隐藏 relationship history。
- Repository 不得把 transport metadata、WebSocket connection identifiers、cookies、headers、tokens 或 sessions 当作 identity proof。

## 7. Version And Conflict Handoff

Repository boundary 只准备 optimistic concurrency，不实现 behavior：

```yaml
version_storage: BIGINT
initial_create_version: 1
version_owner: server
client_authoritative_version_allowed: false
expected_version_handoff: future_behavior_or_interface_gate
conflict_public_shape: deferred_to_protocol_gate
```

Candidate conflict classes：

```yaml
candidate_conflict_classes:
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

规则：

- Repository methods 以后可以区分 internal typed conflicts，但 public protocol error mapping 仍然 deferred。
- Version equality 不是 authentication proof。
- Stale expected version 不得被折叠成隐藏的成功写入。
- Target-player-not-found 和 relationship-not-found 的泄露规则仍由未来 behavior decision 决定。
- PostgreSQL adapter 必须把 unique-index、affected-row 和 foreign-key outcomes 映射为 typed repository conflicts，不得暴露 driver error text。

## 8. Redaction And Logging

Friends relationship state 是 private social graph data。

```yaml
private_social_graph_log_safe: false
relationship_id_log_safe: conditional_after_validation
player_id_log_safe: false_by_default
lifecycle_state_log_safe: conditional_after_validation
version_log_safe: conditional_after_validation
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
  - authorization_header
  - cookie
  - query_string_token
  - chat_room_id
  - group_id
  - party_id
  - match_id
  - pitaya_server_id
```

规则：

- Repository errors 必须 redacted 且 typed。
- Raw player relationship details 和 storage driver errors 默认不得记录到日志。
- Private social graph data 默认按 non-log-safe 处理，除非后续 redaction policy 收窄。
- Repository 不得存储或返回 authentication material、token material、verifier digests、transport metadata、chat room state、group state、party state、match state 或 distributed routing state。

## 9. PostgreSQL Adapter Expectations

未来 PostgreSQL adapter 可以映射到：

```yaml
logical_table: friend_relationships
pair_unique_index: friend_relationships_pair_uq
player_low_state_index: friend_relationships_player_low_state_idx
player_high_state_index: friend_relationships_player_high_state_idx
updated_at_index: friend_relationships_updated_at_idx
version_column: relationship_version
```

Adapter expectations：

- SQL execution 属于 `runtime/internal/platform/persistence/postgres`。
- Unit-of-work 和 transaction handoff 必须遵循已有 platform transaction boundary。
- SQL 不得泄漏到 `runtime/internal/modules/friends`。
- Adapter 必须保留 `player_low_id + player_high_id` 的 canonical pair uniqueness。
- 如果 caller 提供 expected version，updates 必须检查 affected-row 且 version-aware。
- Actor columns 必须在存储前检查 pair-member，或映射为 typed conflicts。
- Adapter gate 接受之后，adapter tests 应覆盖 request creation、duplicate request、accept、reject、remove、block、unblock、pair lookup、player listing、stale version、foreign-key target absence、canonical pair normalization、timestamp mapping 和 redacted errors。

## 10. Relationship To Runtime, Protocol, And Authentication

本 boundary 不改变 runtime 或 protocol behavior：

```yaml
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
authentication_session_behavior_changed: false
request_identity_validation_added: false
```

规则：

- 未来 runtime behavior 必须从 already validated request identity 派生 actor identity。
- 这个 repository boundary 不授权 friendship request、accept、reject、remove、block、unblock、list 或 status runtime handlers。
- 这个 repository boundary 不授权任何 protocol routes、Protobuf messages、generated clients、SDKs 或 transport carriers。
- Repository 不 authenticate users、不 parse tokens、不 validate sessions、不 bind WebSocket connections，也不 enforce route policy。

## 11. Nakama And Pitaya Reference Mapping

Nakama reference mapping：

```yaml
adopted_concepts:
  - friends_relationships_are_core_social_graph_state
  - friend_request_accept_reject_remove_block_unblock_need_durable_state
  - list_and_status_queries_need_repository_ready_boundaries
adapted_concepts:
  - repository_is_vibit_storage_neutral_module_boundary
  - public_status_is_actor_relative_and_computed_later
  - schema_and_protocol_are_vibit_native_not_direct_api_compatibility
deferred_concepts:
  - groups
  - parties
  - chat_targeting
  - matchmaking_social_filters
  - match_runtime_social_context
rejected_for_now:
  - direct_nakama_session_or_friends_api_compatibility
```

Pitaya reference mapping：

```yaml
pitaya_reference_status: deferred_future_architecture_reference
deferred_concepts:
  - frontend_backend_cluster_social_graph_routing
  - distributed_group_membership
  - RPC_or_service_discovery_for_friendship_operations
  - distributed_session_social_context
rejected_for_now:
  - direct_pitaya_api_compatibility
```

Nakama 是近期 capability coverage 的 primary product reference。Pitaya 继续作为 future distributed architecture reference 延后。两者都不能覆盖 vibit 的 constitution、ADRs、manifests、generated boundaries 或 verification commands。

## 12. Future Implementation Queue

这个 boundary 之后，后续工作仍应拆开：

```yaml
future_work_items:
  friends_relationship_repository_interface_implementation:
    work_item: W-0235
    may_add:
      - runtime/internal/modules/friends
      - storage-neutral repository interface and value types
      - focused repository vocabulary tests
    must_not_add:
      - PostgreSQL adapter behavior
      - runtime friendship behavior
      - protocol routes
      - Protobuf source
      - generated output
  friends_relationship_postgresql_adapter_gate:
    may_define:
      - adapter ownership
      - transaction handoff
      - SQL mapping
      - adapter tests
  friends_relationship_runtime_behavior_gate:
    may_define:
      - actor derivation from validated request identity
      - lifecycle command behavior
      - public conflict and leakage policy
  friends_relationship_protocol_route_gate:
    may_define:
      - route names
      - payloads
      - generated output posture
```

不要在没有新 ADR 的情况下把这些合并成一个宽泛 social subsystem slice。

## 13. Stop Conditions

在做以下事情之前，必须停止并要求单独 bounded work item：

- 添加 `runtime/internal/modules/friends`；
- 添加 friends repository interface implementation；
- 添加 PostgreSQL friends adapter behavior；
- 添加 runtime friendship behavior；
- 添加 protocol routes、Protobuf source 或 generated output；
- 添加 migrations 或 event/audit tables；
- 添加 chat、groups、parties、matchmaking、match runtime、SDKs、hosted surfaces、release artifacts 或 distributed runtime；
- 添加 direct Nakama/Pitaya API compatibility。

## 14. Verification

Repository verification for this boundary 是：

```bash
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change define-friends-relationship-repository-boundary --json
node tools/vibit check all --json
```

Repository check rule 是：

```yaml
runtime.friends_relationship_repository_boundary
```
