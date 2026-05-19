# 运行时会话验证门

状态：Draft v0.1
最后更新：2026-05-17
范围：PostgreSQL 会话适配器之后，未来运行时会话验证行为的 gate-only 边界
依赖：`docs/session-postgresql-adapter-gate.md`、`decisions/ADR-0064-session-postgresql-adapter-implementation.md`、`docs/session-repository-boundary.md`、`docs/session-persistence-websocket-handshake-ratification.md`、`docs/reference-game-server-alignment.md`
规范决策：`ADR-0065`

配对英文源文件是 `docs/runtime-session-validation-gate.md`。英文文件是权威版本。

## 1. 目的

vibit 现在已经具备设计会话验证之前需要的 durable 基础：

- PostgreSQL `runtime_sessions` 迁移源。
- 存储中立的 `runtime/internal/app/session.Repository` 接口。
- 该接口的 PostgreSQL 适配器。
- 现有 request-level access-token 验证和 first-message connection binding，它们当前生产路径都保持 `SessionValidated` 为 false。

下一步不应该静默启用会话验证。下一步应该先定义未来验证 gate。成熟游戏服务器给这个边界施加了设计压力：

- Nakama 把已认证 session 视为带过期、刷新、登出和 socket 影响的生命周期对象。
- Nakama 也把 token/session logout 和 active socket disconnect 区分开。
- Pitaya 把 session context 暴露给 handler，但保持 acceptor、session binding、handler execution 分层。

vibit 应该吸收这些经验，把运行时会话验证保持为 application-owned 且显式的边界。本标准只定义 gate。

```yaml
runtime_session_validation_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0142
decision: ADR-0065
check_rule: runtime.runtime_session_validation_gate
future_validation_owner: runtime/internal/app
future_session_repository_owner: runtime/internal/app/session
postgresql_session_adapter_owner: runtime/internal/platform/persistence/postgres
future_validator_source_candidate: runtime/internal/app/runtime_session_validator.go
future_validator_test_candidate: runtime/internal/app/runtime_session_validator_test.go
runtime_session_validation_added: false
request_identity_session_validated_true_added: false
session_creation_composition_added: false
route_policy_session_identity_added: false
route_policy_bound_identity_added: false
websocket_handshake_authentication_added: false
transport_credential_carrier_added: false
protobuf_session_messages_added: false
existing_protobuf_envelope_change_added: false
logout_revocation_active_connection_added: false
reconnect_epoch_behavior_added: false
cleanup_jobs_added: false
dependencies_added: false
memory_durable_session_behavior_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

这是 gate-only 标准。它不添加运行时验证代码。

## 2. 所有权

未来运行时会话验证归 application 层所有：

```yaml
future_validation_owner: runtime/internal/app
session_record_owner: runtime/internal/app/session
postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
authentication_service_owner: runtime/internal/app/authentication
transport_owner: runtime/internal/platform/transport/ws
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
```

规则：

- 未来 validator 可以依赖 application-facing session repository 能力，但不能依赖 PostgreSQL 细节。
- 未来 validator 只能通过 application-owned 类型消费已经规范化的 request identity、access-token 验证结果或 connection binding 结果。
- 未来 validator 不能 import WebSocket transport、generated Protobuf、PostgreSQL adapter、pgx、SQL row 或 migration package。
- PostgreSQL session adapter 必须保持 persistence-only，不能创建 `RequestIdentity`。
- 本 gate 下 WebSocket transport 必须继续保持 credential-neutral。
- Protocol adapter 不能决定 session 是否有效。只有后续独立 protocol gate 才能允许其承载 proof 或 metadata。

## 3. 未来验证语义

后续实现 slice 可以定义包含这些候选输入的 validator：

```yaml
candidate_validation_input:
  - route_request_identity
  - route_request_session_metadata
  - server_observed_connection_id
  - server_observed_connection_epoch
  - observed_at
  - session_repository
```

候选未来验证顺序：

1. repository lookup 前拒绝缺失或格式非法的必要输入。
2. 信任 persisted session row 之前，必须已有 validated actor identity。
3. 通过 `session.Repository` 使用服务器拥有的 `session_id` 查找 runtime session。
4. 要求 `session_status = active`。
5. 要求 `expires_at > observed_at`。
6. 要求 session actor 匹配已经验证过的 actor identity。
7. player account identity handoff 必须继续保持规范化且由 application 层拥有。
8. 只有后续实现 gate 授权时，才可以更新 `last_seen_at`。
9. 只有所有被选择的检查通过后，才能返回 `RequestIdentity.SessionValidated = true`。
10. 对外失败输出折叠为稳定的 invalid-session 错误类。

规则：

- 仅 persisted `session_id` 不是 proof。
- `access_token_record_id` linkage 是私有 metadata。除非后续 ADR 明确定义组合方式，否则它不能替代 access-token proof validation。
- 客户端提供的 envelope `Session.session_id` 在被 durable state 和已经验证过的 actor identity 校验前仍然只是 metadata。
- session validation 不能隐式创建 session。
- session validation 不能刷新 access token、延长 session 过期时间、撤销 token 或登出用户。
- session validation 不能关闭 WebSocket 连接。

## 4. Request Identity 交接

未来验证可以设置：

```yaml
future_request_identity:
  status: validated
  actor_kind: player
  player_id_validated: true
  session_validated: true
```

只有后续实现 slice 可以在生产行为里把 `SessionValidated` 设为 true。

规则：

- `SessionValidated` true 需要同时具备已经验证过的 actor identity 和有效 persisted runtime session row。
- metadata-only identity 永远不能变成 session-validated。
- first-message bound identity 不能通过本 gate 满足普通 protected route policy。
- route policy 只有在独立 route-policy gate 定义哪些 route 需要 session-validated identity 以及失败如何映射后，才能使用 session-validated identity。
- Domain module 必须继续接收规范化 `RequestIdentity`，而不是 session repository record。

## 5. 错误与脱敏边界

未来对外 session validation 错误应该折叠为：

```yaml
public_error_class:
  - SESSION_INVALID
```

内部失败原因可以在测试和私有控制流中更具体：

```yaml
candidate_internal_failure_reasons:
  - missing_session_id
  - malformed_session_id
  - actor_identity_not_validated
  - session_not_found
  - session_inactive
  - session_expired
  - session_actor_mismatch
  - session_player_mismatch
  - session_repository_unavailable
```

规则：

- 对外错误不能泄漏失败来自 lookup、过期、撤销、actor mismatch、player mismatch、token linkage 或 repository failure。
- 错误、日志、事件和测试不能包含 raw access-token text、raw credential material、lookup digest、verifier digest、verifier key id、Authorization header、cookie、query-string token 或 WebSocket subprotocol token material。
- 除非后续 observability gate 把特定脱敏形式分类为 log-safe，否则 session id 和 player id 应视为运维敏感信息。

## 6. 与认证的关系

运行时会话验证不替代 token proof validation。

```yaml
access_token_validation_owner: runtime/internal/app/authentication
session_validation_owner: runtime/internal/app
session_validation_replaces_token_validation: false
session_validator_reads_token_digests: false
session_validator_reads_raw_tokens: false
```

规则：

- 当前 protected routes 的 proof-validation 路径仍然是 access-token validation。
- 未来 session validator 可以消费 authentication 产生的 validated actor identity，但不能计算 token digest、比较 token verifier、签发 token、刷新 token 或撤销 token。
- route policy 依赖 access-token validation 和 session validation 组合之前，必须显式定义组合方式。
- logout 和 revocation 的 active-connection 行为仍然放在独立 gate 后面。

## 7. 与 WebSocket 和协议的关系

本 gate 不改变 WebSocket 或 Protobuf 行为：

```yaml
websocket_transport_credential_neutral: true
websocket_handshake_authentication_added: false
transport_credential_carrier_added: false
protobuf_session_messages_added: false
existing_protobuf_envelope_change_added: false
connection_registry_added: false
```

规则：

- 本 gate 下 WebSocket transport 不能解析 Authorization header、bearer value、cookie、query-string token、session token 或 `Sec-WebSocket-Protocol` 认证材料。
- 现有 Protobuf envelope 保持不变。
- 这里不授权 session validation protocol message 或 generated output。
- 这里不授权 reconnect、resume、duplicate replacement、durable connection epoch、logout disconnect、presence、room、party、group 或 match attachment 行为。

## 8. 未来实现测试要求

后续运行时会话验证实现必须包含聚焦测试：

- repository lookup 前拒绝缺失或格式非法的 session id。
- 拒绝 metadata-only identity。
- 要求已经验证过的 actor identity。
- 通过 `session.Repository` 查询 active session。
- expired、revoked、not-found、actor-mismatch 都折叠为同一个对外 invalid-session 类。
- 成功验证只有在所有检查通过后才设置 `SessionValidated` true。
- repository error 保持脱敏，不泄漏 raw proof 或 digest material。
- 不解析 WebSocket transport credential。
- 不改变 Protobuf envelope shape。
- 除非后续 route-policy gate 授权，否则不让 route policy 使用 session-validated identity。

除非后续实现工作项要求，live PostgreSQL verification 可以保持 opt-in。

## 9. Nakama 与 Pitaya 参考映射

Nakama 参考映射：

```yaml
adopted_concepts:
  - sessions_have_lifetime_and_invalid_state
  - logout_and_revocation_affect_session_validity
  - realtime_socket_lifecycle_is_related_but_not_identical_to_session_validity
adapted_concepts:
  - vibit_uses_opaque_access_token_proof_plus_durable_runtime_session_records
  - public_invalid_session_failures_are_collapsed
  - session_validation_is_application_owned
deferred_concepts:
  - refresh_token_session_extension
  - logout_disconnect_active_socket
  - admin_session_management_api
  - single_socket_or_single_session_policy
rejected_for_now:
  - direct_nakama_session_api_compatibility
```

Pitaya 参考映射：

```yaml
adopted_concepts:
  - session_context_is_available_to_handlers_through_context
  - acceptor_transport_and_handler_logic_are_separate
  - session_lifecycle_callbacks_are_distinct_from_request_handler_logic
adapted_concepts:
  - durable_session_validation_builds_application_request_identity
  - transport_acceptor_remains_credential_neutral
  - route_handlers_receive_normalized_identity_not_storage_rows
deferred_concepts:
  - unique_session_enforcement
  - frontend_backend_cluster_session_routing
  - durable_connection_registry
  - group_or_room_session_membership
rejected_for_now:
  - direct_pitaya_session_api_compatibility
```

Nakama 和 Pitaya 指导能力规划。它们不覆盖 vibit 的 constitution、ADR、manifest、generated boundary 或 verification command。

## 10. 未来实现队列

本 gate 之后，未来工作仍应拆开：

```yaml
future_work_items:
  runtime_session_validation_implementation:
    may_add:
      - application-owned runtime session validator
      - tests with fake session repository
      - redacted error mapping
  session_creation_composition_gate:
    may_define:
      - whether login or BindConnection creates durable sessions
      - session id generation
      - token record linkage
      - expiration and last_seen semantics
  bound_identity_route_policy_gate:
    may_define:
      - which routes can use bound or session-validated identity
      - failure behavior
  logout_revocation_active_connection_gate:
    may_define:
      - whether revocation closes active WebSocket connections
  reconnect_connection_epoch_gate:
    may_define:
      - duplicate replacement
      - reconnect and resume behavior
```

不要在没有新 ADR 的情况下，把这些合并成一个宽泛的 session subsystem slice。

## 11. 验证

本边界的仓库验证是：

```bash
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change define-runtime-session-validation-gate --json
node tools/vibit check all --json
```

仓库检查规则是：

```yaml
runtime.runtime_session_validation_gate
```
