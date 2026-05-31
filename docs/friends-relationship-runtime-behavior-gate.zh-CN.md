# Friends Relationship Runtime Behavior Gate

状态：Accepted v0.1
最后更新：2026-05-26
范围：PostgreSQL adapter 之后、未来 application-owned friends relationship runtime behavior 的 gate-only boundary
依赖：`docs/friends-relationship-lifecycle-gate.md`、`docs/friends-relationship-repository-boundary.md`、`docs/friends-relationship-postgresql-adapter-gate.md`、`runtime/internal/modules/friends/repository.go`、`runtime/internal/platform/persistence/postgres/friend_relationship_repository.go`、`docs/runtime-protocol-adapter.md`、`docs/bound-identity-route-policy-gate.md`、`docs/runtime-session-validation-gate.md`、`docs/reference-game-server-alignment.md`
Canonical decision：`ADR-0146`

配套英文原文是 `docs/friends-relationship-runtime-behavior-gate.md`。英文文件是权威版本。

本文定义 friends relationship runtime behavior gate。这是一个 gate artifact。它不添加 runtime behavior implementation、runtime handlers、startup wiring、protocol routes、Protobuf source、generated output、repository interface 变更、PostgreSQL adapter 变更、migration 变更、dependencies、authentication/session behavior 变更、delivery guarantees、stream subscriptions、chat rooms、groups、parties、broadcast fanout、matchmaking、match runtime、operations/admin behavior、SDK publication、generated client libraries、hosted deployments、release artifacts、public announcements、paid promotion、event/audit tables、Pitaya-style distributed architecture 或 direct Nakama/Pitaya API compatibility。

## 1. 核心规则

friends relationship runtime behavior gate 记录为：

```yaml
friends_relationship_runtime_behavior_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0238
decision: ADR-0146
check_rule: runtime.friends_relationship_runtime_behavior_gate
source_postgresql_adapter_implementation_decision: ADR-0145
source_postgresql_adapter: runtime/internal/platform/persistence/postgres/friend_relationship_repository.go
source_repository_interface_decision: ADR-0143
repository_interface: runtime/internal/modules/friends.Repository
repository_interface_source: runtime/internal/modules/friends/repository.go
future_runtime_owner_candidate: runtime/internal/app
future_friends_application_package_candidate: runtime/internal/app/friends
future_runtime_service_source_candidate: runtime/internal/app/friends/service.go
future_runtime_service_test_candidate: runtime/internal/app/friends/service_test.go
request_identity_required: true
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
actor_identity_source: validated_request_identity_player_id
first_actor_kind: player
actor_relative_public_status_required: true
route_policy_requirement: request_token_required
service_application_owner: runtime/internal/app
repository_handoff: unit_of_work_friend_relationship_repository_factory
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
event_audit_table_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_runtime_behavior_implementation_work_item: W-0239
future_runtime_behavior_implementation_direction: implement_friends_relationship_runtime_behavior
```

## 2. 目的

`W-0237` 已为 `runtime/internal/modules/friends.Repository` 实现 PostgreSQL adapter。下一个有用边界不是 route 或 protocol 变更，而是 runtime behavior gate：先定义未来 application code 如何把已验证的 player request 转换成 friends repository operations。

本 gate 在实现前记录未来行为形状：

- service 的 application ownership；
- actor identity 从 validated request identity 派生；
- actor-relative status 和 list behavior；
- permission 和 route-policy posture；
- validation 和 conflict mapping expectations；
- unit-of-work 和 repository handoff；
- redaction rules；
- test expectations；
- stop conditions，确保 protocol、generated output、authentication/session 变更、event/audit tables 和更宽 social features 不进入本 slice。

Nakama 提供 durable friends relationships 作为核心 social graph capability 的产品压力。Pitaya 提醒我们保持 handlers、sessions、route context 和 persistence responsibilities 分离。vibit 通过明确的 application-owned behavior 和 checks 来吸收这些参考，而不是追求 direct public API compatibility。

## 3. Ownership

未来 runtime behavior 由 application 层拥有：

```yaml
future_runtime_owner_candidate: runtime/internal/app
future_friends_application_package_candidate: runtime/internal/app/friends
future_runtime_service_source_candidate: runtime/internal/app/friends/service.go
future_runtime_service_test_candidate: runtime/internal/app/friends/service_test.go
repository_interface_owner: runtime/internal/modules/friends
postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
route_policy_owner: runtime/internal/app
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
websocket_transport_owner: runtime/internal/platform/transport/ws
player_account_owner: runtime/internal/app/player
```

规则：

- 未来 service behavior 可以放在 `runtime/internal/app/friends`，或由 implementation slice 批准的等价 application-owned package。
- service 只能通过 application 或 unit-of-work dependencies 调用 `runtime/internal/modules/friends.Repository`。
- service 只有在 implementation slice 授权 dependency handoff 后，才可以通过已有 application-owned player repository capability 检查 target player account 是否存在和状态。
- service 不得导入 PostgreSQL adapter packages、SQL row types、migration packages、WebSocket transport packages、generated Protobuf packages、chat packages、group packages、party packages、matchmaking packages、match runtime packages、SDK packages 或 distributed runtime packages。
- friends module 继续拥有 storage-neutral value types、normalizers、lifecycle/status vocabulary 和 repository error vocabulary。
- PostgreSQL adapter 仍然只负责 persistence，不派生 request identity、route policy 或 public protocol errors。
- transport 和 protocol adapters 不拥有 friends relationship permission 或 business behavior。

## 4. Request Identity And Actor Derivation

第一版 runtime behavior posture 是 player-to-player：

```yaml
first_actor_kind: player
actor_identity_source: validated_request_identity_player_id
request_identity_required: true
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
client_supplied_actor_id_allowed_as_proof: false
client_supplied_target_player_id_allowed: true
self_relationship_allowed: false
```

规则：

- 未来 friends operation 必须从 validated `app.RequestIdentity` 派生 actor。
- `RequestIdentity.Status` 必须是 `validated`。
- `RequestIdentity.ActorKind` 必须是 `player`。
- `RequestIdentity.PlayerIDValidated` 必须为 true。
- `RequestIdentity.PlayerID` 必须非空，并且在同时存在 actor identity 时必须一致。
- envelope/session metadata 中的 metadata-only `player_id` 绝不能满足本 gate。
- 单独 persisted `session_id` 绝不能成为 proof。
- client payload 可以为 player-to-player operations 指定 target player，但不能把 client-supplied actor 当成 proof。
- self-targeting 应尽可能在 repository mutation 前失败。

本 gate 不改变 `RequestIdentity`、access-token validation、bound connection identity、durable runtime session validation 或 WebSocket handshake behavior。它只记录未来 friends behavior 在 repository access 前必须要求的 identity 条件。

## 5. Future Runtime Behavior Shape

未来第一版实现可以暴露一个 application service，候选操作如下：

```yaml
candidate_operations:
  - send_friend_request
  - accept_friend_request
  - reject_friend_request
  - remove_friend
  - block_player
  - unblock_player
  - list_friend_relationships
  - get_friend_relationship_status
```

推荐的第一版 posture：

- `send_friend_request` 从 validated actor 向 target player 创建 outgoing pending relationship。
- `accept_friend_request` 只接受 incoming pending relationship。
- `reject_friend_request` 只拒绝 incoming pending relationship。
- `remove_friend` 按 repository lifecycle rules 移除已有 friendship 或 pending relationship，但不得移除 actor-specific block。
- `block_player` 设置 actor-specific block state，并在 blocked 状态阻止普通 friend operations。
- `unblock_player` 只清除 actor 自己的 block state，不得恢复之前的 friendship。
- `list_friend_relationships` 列出 actor-scoped relationships，并计算 actor-relative public statuses。
- `get_friend_relationship_status` 读取一个 actor-target relationship，并计算 actor-relative public status。

规则：

- Runtime behavior 必须使用 server-derived actor identity。
- Runtime behavior 必须尽可能在 repository calls 前验证 target player id、self-targeting、expected version、list status filters、pagination cursors 和 list limit。
- Runtime behavior 必须相对 requesting actor 计算 public status。
- Runtime behavior 第一版不得暴露其他玩家的 private social graph、任意 social graph search、admin inspection、group/guild relationships、party memberships、chat rooms、match rooms、matchmaking filters、server script hooks 或 direct external API compatibility。
- 除非后续 protocol gate 授权，runtime behavior 不得添加 public protocol routes 或 generated output。

## 6. Candidate Application Service Shape

第一版 implementation slice 应定义一个小的 application-owned service boundary。候选输入和输出：

```yaml
candidate_request_fields:
  - request_identity
  - target_player_id
  - expected_relationship_version
  - status_filter
  - list_limit
  - after_relationship_cursor
candidate_result_fields:
  - relationship
  - relationships
  - actor_relative_public_status
  - next_relationship_cursor
  - relationship_version
  - public_error_code
```

规则：

- service 应接收 already-normalized application identity，而不是 raw tokens、cookies、headers、WebSocket subprotocol values 或 envelope proof carriers。
- service 应在 repository handoff 前调用 friends module normalizers。
- service 应让 actor ids、target player ids、relationship ids、relationship state、block state 和 relationship versions 远离默认 errors 和 logs。
- service 应暴露稳定的 public error codes 或 classes，供后续 runtime handlers 映射。
- gate slice 不应添加 route registration、Protobuf conversion、startup composition 或 command/query dispatch wiring。

## 7. Validation Rules

未来 runtime behavior 必须在 persistence 前执行 validation：

```yaml
validation_required:
  - validated_player_identity
  - target_player_id_non_empty_and_normalized
  - self_target_forbidden
  - target_player_lookup_or_repository_conflict_handled
  - actor_must_be_pair_member_for_existing_relationship
  - expected_relationship_version_positive_when_present
  - list_limit_bounded
  - pagination_cursor_bounded
  - status_filter_known
```

规则：

- target player id validation 应复用已有 player identity rules 或 friends module normalizers。
- self-targeting 应尽可能在 repository mutation 前失败。
- implementation tests 必须明确 missing expected version behavior。
- invalid input 应尽可能在 repository mutation 前失败。
- target player not found、relationship not found、blocked relationship 和 privacy-sensitive cases 不得泄露隐藏 social graph state。
- repository unavailable errors 必须保持 redacted。

## 8. Permission And Route Policy Posture

第一版 route-policy posture 保守：

```yaml
route_policy_requirement: request_token_required
public_friends_routes_allowed: false
bound_connection_required_by_this_gate: false
session_validated_required_by_this_gate: false
bound_session_required_by_this_gate: false
```

后续 public contracts 的候选 permission families：

- send friend request；
- accept incoming friend request；
- reject incoming friend request；
- remove friend；
- block player；
- unblock player；
- list own friend relationships；
- read own actor-relative relationship status。

规则：

- Friends relationship routes 必须是 protected routes。
- 除非后续 route-policy ADR 改变 named routes，第一版 posture 应使用已有 `request_token_required` route protection family。
- public routes 不得读取或修改 friends relationships。
- bound connection identity 和 durable session validation 可以保留给未来 route families，但本 gate 不要求它们，也不改变普通 protected route behavior。
- metadata-only identity 必须 fail closed。

## 9. Conflict And Error Mapping

未来 runtime behavior 必须把 friends repository errors 映射成稳定 application classes：

```yaml
candidate_public_error_classes:
  - FRIENDSHIP_INVALID_REQUEST
  - FRIENDSHIP_UNAUTHENTICATED
  - FRIENDSHIP_FORBIDDEN
  - FRIENDSHIP_TARGET_NOT_FOUND
  - FRIENDSHIP_RELATIONSHIP_NOT_FOUND
  - FRIENDSHIP_DUPLICATE_REQUEST
  - FRIENDSHIP_ALREADY_FRIENDS
  - FRIENDSHIP_BLOCKED_RELATIONSHIP
  - FRIENDSHIP_INVALID_TRANSITION
  - FRIENDSHIP_VERSION_MISMATCH
  - FRIENDSHIP_UNAVAILABLE
```

规则：

- missing 或 unvalidated request identity 必须在 repository access 前映射到已有 protected-route authentication posture。
- self-targeting 和 malformed inputs 可以是 public invalid-request classes。
- duplicate pending request、already friends、blocked relationship、invalid transition 和 version mismatch 只有在不泄露 actor 自身关系之外的 private graph state 时，才可以成为 public conflict classes。
- target-player-not-found 和 relationship-not-found leakage rules 必须保守；public output 可以把 privacy-sensitive cases 折叠成 not-found 或 forbidden。
- stored private relationship state、block details、actor ids、target ids、relationship ids、SQL details、driver errors、DSNs、credentials、token material、verifier digests 和 route proof carriers 不得泄露。
- repository `storage_unavailable` errors 必须映射成 unavailable class，且不暴露 platform internals。
- Runtime behavior 不得添加超过现有 application route-protection classes 的 authentication/token/session failure detail。

## 10. Unit-Of-Work And Repository Handoff

未来 runtime behavior 应使用已有 application transaction boundary：

```yaml
unit_of_work_handoff_required: true
repository_handoff: unit_of_work_friend_relationship_repository_factory
repository_interface: runtime/internal/modules/friends.Repository
postgresql_adapter: runtime/internal/platform/persistence/postgres.FriendRelationshipRepository
service_starts_transactions: false
repository_starts_transactions: false
```

规则：

- state-changing operations 必须通过 application-owned unit-of-work 或等价 transaction boundary 执行。
- service 应从 unit-of-work 获取 `friends.Repository`，而不是直接构造 PostgreSQL adapter。
- service 不得发出 `BEGIN`、`COMMIT` 或 `ROLLBACK` SQL。
- service 不得导入 PostgreSQL-specific packages。
- read-only operations 可以直接使用 repository dependencies，也可以通过 read-only unit-of-work；具体形状由后续 implementation gate 批准。
- 当 target player lookup 和 relationship mutation 都需要时，必须在未来 unit-of-work 中确定顺序。
- Repository errors 必须在 rollback 或 failed unit-of-work outcome 后映射，且不泄露 platform internals。

## 11. Actor-Relative Status Rules

Public status 是 actor-relative：

```yaml
actor_relative_public_states:
  - none
  - outgoing_request_pending
  - incoming_request_pending
  - friends
  - blocked_by_actor
  - blocked_actor
  - mutual_blocked
  - removed
  - rejected
```

规则：

- persisted pair-oriented state 必须先相对 requesting actor 转换，再进入 public output。
- pending requests 必须区分 actor 的 outgoing 和 incoming。
- actor-specific block columns 必须区分 actor blocked target、target blocked actor 和 mutual block。
- `unblock_player` 不得恢复 friendship 或 pending request state。
- 除非后续 public contract 明确授权，removed 和 rejected states 不得暴露其他玩家的 private relationship history。
- list output 必须保持 player-scoped，不得泄露其他玩家之间的 relationships。

## 12. Redaction And Logging

Friends relationship runtime data 是 private social graph data。

```yaml
private_social_graph_log_safe: false
actor_player_id_log_safe: false_by_default
target_player_id_log_safe: false_by_default
relationship_id_log_safe: conditional_after_validation
relationship_state_log_safe: conditional_after_validation
relationship_version_log_safe: conditional_after_validation
forbidden_runtime_log_material:
  - raw_access_token
  - raw_credential
  - credential_lookup_digest
  - credential_verifier_digest
  - token_lookup_digest
  - token_verifier_digest
  - verifier_key
  - verifier_key_id
  - authorization_header
  - cookie
  - query_string_token
  - websocket_subprotocol
  - websocket_connection_id
  - remote_address
  - sql_text
  - database_dsn
  - private_relationship_graph
  - chat_room_id
  - group_id
  - party_id
  - match_id
  - pitaya_server_id
```

规则：

- Application errors 必须 redacted and typed。
- raw player relationship details、target identifiers、block details 和 storage driver errors 默认不得写入日志。
- Private social graph data 必须被视为 non-log-safe，除非后续 redaction policy 缩小该范围。
- Runtime behavior 不得存储或返回 authentication material、token material、verifier digests、transport metadata、chat room state、group state、party state、match state 或 distributed routing state。

## 13. Future Test Expectations

后续 implementation slice 应添加 focused application service tests：

```yaml
future_tests:
  - service_requires_validated_player_identity
  - service_rejects_metadata_only_identity
  - service_rejects_client_supplied_actor_as_proof
  - send_friend_request_validates_target_and_self_relationship
  - send_friend_request_maps_duplicate_already_friends_blocked_conflicts
  - accept_reject_require_incoming_pending_request
  - remove_friend_preserves_block_semantics
  - block_unblock_use_actor_specific_block_state
  - list_relationships_is_actor_scoped_and_status_filtered
  - get_status_computes_actor_relative_public_status
  - expected_version_mismatch_maps_to_public_conflict
  - repository_unavailable_is_redacted
  - state_changing_operations_use_unit_of_work
  - no_protocol_or_transport_dependency
```

规则：

- 除非后续 implementation slice 授权 live PostgreSQL integration，tests 应使用 fakes 或 in-memory stubs。
- Tests 必须验证 missing 或 metadata-only identity 会在 repository mutation 前失败。
- Tests 必须验证 pending、friends、actor block、target block、mutual block、removed 和 rejected states 的 actor-relative status conversion。
- Tests 必须验证 target player not found、relationship not found 和 privacy-sensitive failures 不泄露 hidden graph details。
- Tests 不得要求 protocol routes、WebSocket transport、Protobuf generated code 或 generated clients。
- Tests 不得输出 raw credentials、tokens、verifier keys、digests、DSNs、query strings、authorization headers、cookies、player ids 或 private relationship state。

## 14. Relationship To Runtime, Protocol, And Authentication

本 gate 不改变 runtime 或 protocol behavior：

```yaml
runtime_friends_service_added: false
runtime_friends_handlers_added: false
friends_protocol_routes_added: false
protobuf_friends_messages_added: false
generated_friends_output_added: false
authentication_session_behavior_changed: false
request_identity_handoff_changed: false
```

规则：

- Runtime friendship behavior implementation 仍然 deferred。
- Protocol routes 和 generated friends contract shapes 仍然 deferred。
- Request identity validation 继续由 authentication/session 和 route-policy boundaries 拥有。
- 未来 service 不得解析 bearer tokens、cookies、WebSocket subprotocols、envelope metadata 或 transport connection identifiers。
- 本 gate 不改变 access-token validation、runtime session validation、bound connection identity、route protection、WebSocket handshake behavior、logout behavior 或 reconnect behavior。

## 15. Reference Alignment

Nakama 提供 durable friends relationship social graph behavior 的产品能力压力。Pitaya 作为 distributed runtime concerns 的未来架构参考继续 deferred。vibit 只把这些参考用于 capability planning：

- 本 gate 不添加 direct Nakama 或 Pitaya API compatibility；
- 本 gate 不添加 public friends route；
- 本 gate 不添加 server runtime hook、group、party、chat、matchmaking、match runtime、admin surface 或 distributed routing behavior；
- 未来 runtime behavior 是位于 module repository vocabulary 之上、protocol route/public contract behavior 之下的 application service。

## 16. Stop Conditions

在执行下列任何事项前，停止并打开后续 bounded work item：

- 实现 `runtime/internal/app/friends/service.go`；
- 添加 runtime friend request/list/status handlers；
- 添加 command/query dispatch registration 或 startup composition；
- 添加 protocol routes、Protobuf sources、generated output 或 generated clients；
- 改变 authentication/session behavior 或 request identity validation；
- 改变 friends repository interface 或 PostgreSQL adapter；
- 添加或修改 migrations；
- 添加 event/audit tables；
- 添加 chat、groups、parties、broadcast fanout、matchmaking、match runtime、SDK、hosted、release 或 distributed runtime scope；
- 添加 direct Nakama 或 Pitaya public API compatibility。

