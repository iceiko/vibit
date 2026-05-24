# 好友关系生命周期 Gate

状态：Accepted v0.1
最后更新：2026-05-24
范围：在 persistence、protocol、runtime behavior 或更大的 social features 之前，为未来 player friendship relationship lifecycle 定义语义 gate
依赖：`docs/agent-native-feature-request-test-workflow.md`、`docs/agent-native-feature-request-scaffolding.md`、`docs/reference-game-server-alignment.md`、`docs/nakama-pitaya-product-parity-roadmap.md`
Canonical decision: `ADR-0139`

配套英文源文件是 `docs/friends-relationship-lifecycle-gate.md`。英文文件是权威版本。

本文定义 friends relationship lifecycle semantic gate。它是 gate artifact。它不添加 runtime friendship behavior、protocol routes、Protobuf source、generated output、migrations、repository interfaces、PostgreSQL adapters、dependencies、startup wiring、authentication/session behavior changes、delivery guarantees、stream subscriptions、chat rooms、groups、parties、broadcast fanout、matchmaking、match runtime、operations/admin behavior、SDK publication、generated client libraries、hosted deployments、release artifacts、Pitaya-style distributed architecture，或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

好友关系生命周期 gate 记录如下：

```yaml
friends_relationship_lifecycle_gate: defined
completed_work_item: W-0231
decision: ADR-0139
check_rule: runtime.friends_relationship_lifecycle_gate
source_intake_decision: ADR-0138
source_scaffolding_decision: ADR-0137
source_workflow_decision: ADR-0128
gate_standard: docs/friends-relationship-lifecycle-gate.md
gate_standard_translation: docs/friends-relationship-lifecycle-gate.zh-CN.md
selected_nakama_capability_family: friends_groups_and_parties
primary_product_reference: Nakama
pitaya_reference_status: deferred_future_architecture_reference
ai_native_development_testing_goal: user_requirement_to_spec_tests_implementation_verification
semantic_gate_only: true
future_persistence_schema_gate_work_item: W-0232
future_persistence_schema_gate_direction: define_friends_relationship_persistence_schema_gate
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
dependency_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Product Intent

Friendship 是 Nakama-class game backends 的核心 social graph primitive。它后续可以支撑 groups、parties、chat targeting、invites、matchmaking filters、player discovery 和 match social context。vibit 采用的是产品能力，不采用外部 public API shape。

vibit 的姿态是：

- implementation 之前先定义 contract-first lifecycle semantics；
- state transitions 必须 server-authoritative；
- public relationship status 必须 actor-relative；
- 每个未来 command 和 query 都要求 validated player identity；
- 写代码之前先定义测试期望；
- private social graph data 默认不是 log-safe；
- 除非后续 ADR 明确授权，否则拒绝 direct external API compatibility。

Pitaya 仍然只是 future distributed architecture reference。本 gate 不能引入 distributed routing、frontend/backend roles、RPC、cluster groups、service discovery 或 server-to-server messaging。

## 3. Future Semantic Scope

未来好友关系生命周期必须覆盖：

```yaml
semantic_scope:
  - request
  - accept
  - reject
  - remove
  - block
  - unblock
  - list
  - read_relationship_status
```

生命周期是 player-to-player 并且 server-authoritative。未来 domain owner 是 social/friends capability boundary，而不是 WebSocket transport、protocol adapters、authentication、storage objects、inventory、realtime delivery、chat、matchmaking 或 match runtime。

## 4. Future Contract Vocabulary

未来 command vocabulary 是：

```yaml
commands:
  - SendFriendRequest
  - AcceptFriendRequest
  - RejectFriendRequest
  - RemoveFriend
  - BlockPlayer
  - UnblockPlayer
```

未来 query vocabulary 是：

```yaml
queries:
  - ListFriendRelationships
  - GetFriendRelationshipStatus
```

未来 event vocabulary 是：

```yaml
events:
  - FriendRequestCreated
  - FriendRequestAccepted
  - FriendRequestRejected
  - FriendRemoved
  - PlayerBlocked
  - PlayerUnblocked
```

未来 error vocabulary 是：

```yaml
errors:
  - FRIENDSHIP_INVALID_TARGET
  - FRIENDSHIP_SELF_TARGET_FORBIDDEN
  - FRIENDSHIP_DUPLICATE_REQUEST
  - FRIENDSHIP_BLOCKED_RELATIONSHIP
  - FRIENDSHIP_INVALID_TRANSITION
  - FRIENDSHIP_RELATIONSHIP_NOT_FOUND
  - FRIENDSHIP_TARGET_NOT_FOUND
  - FRIENDSHIP_METADATA_IDENTITY_NOT_AUTHENTICATED
```

这些 vocabulary 只是语义规划。它不创建 contract source files、generated shapes、protocol payloads、routes、repositories 或 runtime handlers。

## 5. Identity And Permissions

每个未来 command 和 query 都要求：

```yaml
permission: validated_player_identity
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
actor_identity_source: validated_request_identity
```

规则：

- actor 来自 validated request identity，不来自 client-supplied actor id。
- `player_id` 和 `session_id` metadata 不是 authentication proof。
- 禁止 self-targeting。
- 在 authorization 或 privacy 要求 collapse 时，public failures 不能泄露 target existence 或 hidden relationship details。
- relationship records、target ids、actor ids、private statuses 和 conflict details 默认不是 log-safe。

## 6. Relationship State Model

未来 public relationship status 是 actor-relative：

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

未来 persisted canonical state 可以是 pair-oriented。如果采用 pair-oriented，public query output 仍必须相对请求 actor 计算。storage shape 有意推迟到 `W-0232`。

第一批 lifecycle invariants 是：

- 禁止 self-targeting。
- duplicate request behavior 必须在 implementation 之前明确为 forbidden 或显式 idempotent。
- accept 只能用于 incoming pending request。
- reject 只能用于 incoming pending request。
- remove 只能用于 existing friendship，不能移除 block。
- block 会移除或覆盖 pending 和 friend relationship，转为 actor 的 block。
- unblock 只移除 actor 的 block，不能自动恢复之前的 friendship。
- mutual block status 只在双方都有 active block 时表示。
- list 和 status queries 是 actor-relative，不能暴露其他玩家的 private graph。

## 7. Future Persistence Gate

下一项 bounded work item 是：

```text
W-0232 Define friends relationship persistence schema gate
```

该 follow-up 应在 migration source 出现前定义 table candidates、pair identity posture、block representation、indexes、uniqueness、event/audit posture、redaction 和 repository/adapter boundaries。

本 gate 有意不决定：

- exact table names；
- canonical pair key encoding；
- rejected 或 removed states 是 tombstones、audit-only facts，还是 current rows；
- duplicate request idempotency；
- concurrency conflict resolution；
- hard delete 或 retention policy；
- repository interface shape；
- PostgreSQL adapter SQL；
- protocol routes 或 payloads。

## 8. Future Test Expectations

未来 behavior tests 必须在 implementation 前规划。

Positive tests：

- send friend request；
- accept incoming request；
- reject incoming request；
- remove existing friend；
- block target player；
- unblock previously blocked player；
- list actor-relative relationships；
- read actor-relative relationship status。

Negative tests：

- self-targeting；
- duplicate request 或明确选定的 idempotency behavior；
- invalid transition；
- blocked relationship interaction；
- missing 或 unknown target；
- missing relationship；
- metadata-only identity。

Permission and authentication tests：

- 每个 command/query 都要求 validated player identity；
- client-supplied actor id 会被 ignored 或 rejected；
- metadata-only `player_id` 和 `session_id` 作为 proof 会被拒绝。

Persistence and transaction tests：

- schema 和 repository tests 必须在 `W-0232` 之后、migration/adapter/runtime implementation 之前定义；
- command transitions 必须 transactional；
- emitted events 和 state changes 必须在未来 unit-of-work boundary 内保持一致。

Failure and redaction tests：

- 在 privacy 要求 collapse 的地方，public errors 不泄露 private relationship graph details；
- logs 不暴露 raw credentials、tokens、verifier keys、digests、transport metadata 或 private social graph internals。

Concurrency tests：

- simultaneous request、accept、reject、remove、block 和 unblock conflicts 必须在 runtime implementation 之前有明确 expected outcomes。

Integration and end-to-end tests：

- 推迟到 protocol routes 和 runtime handlers 获得授权后。

## 9. Non-Authorization

本 gate 不授权：

- runtime friendship behavior；
- protocol routes；
- Protobuf source；
- generated output；
- migrations；
- repository interfaces；
- PostgreSQL adapters；
- dependencies；
- startup wiring；
- authentication/session behavior changes；
- delivery guarantees；
- stream subscriptions；
- chat rooms；
- groups；
- parties；
- broadcast fanout；
- matchmaking；
- match runtime；
- operations/admin behavior；
- SDK publication；
- generated client libraries；
- hosted deployments；
- release artifacts；
- Pitaya-style distributed architecture；
- direct Nakama/Pitaya API compatibility。

这些领域的任何未来工作都需要单独的 bounded work item 和 verification record。

## 10. Verification

本 gate 通过以下命令验证：

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.friends_relationship_lifecycle_gate
node tools/vibit check change define-friends-relationship-lifecycle-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```
