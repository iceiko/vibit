# 注销撤销活跃连接 Gate

状态：Draft v0.1
最后更新：2026-05-18
范围：未来 logout、token revocation、runtime session revocation 和 active WebSocket connection invalidation 行为的 gate-only 边界
依赖：`docs/authentication-contract-error-permission-surfaces.md`、`docs/session-persistence-websocket-handshake-ratification.md`、`docs/first-message-connection-binding-implementation-gate.md`、`docs/runtime-session-validation-gate.md`、`docs/session-creation-composition-gate.md`、`docs/bound-identity-route-policy-gate.md`、`decisions/ADR-0070-bound-identity-route-policy-implementation.md`、`docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0071`

英文源文件是 `docs/logout-revocation-active-connection-gate.md`。英文文件是权威版本。

## 1. 目的

vibit 现在已经有足够多的 authentication 和 session 组件，因此 logout 和 revocation 不能再被当成纯 storage 细节：

- Request-level opaque access-token validation。
- Public device-credential login command。
- Durable access-token verifier records。
- 成功 login 时创建的 durable `runtime_sessions` rows。
- 可以把观测到的 WebSocket connection 和 validated player identity 关联起来的 first-message connection binding。
- 面向 request-token、bound-connection、session-validated 和 bound-session routes 的显式 route policy families。

缺失的边界是：未来 logout 或 revocation operation 应该如何影响活跃 WebSocket connections。成熟游戏服务器把这个压力暴露得很清楚：

- Nakama 把 authenticated session material 当成带有 logout、refresh、expiration 和 realtime socket 影响的 lifecycle object，同时仍然把 token/session validity 和 active socket behavior 作为不同 implementation concerns。
- Pitaya 分离 acceptors、session context、handlers 和 connection management；session binding 和类似 kick 的行为属于 connection/session lifecycle concern，不应该由 handler 自己解析 transport credential。

vibit 应该吸收这些经验，在任何代码可以 revoke tokens、revoke runtime sessions、close sockets 或添加 connection registry 之前先定义 policy boundary。本标准只是 gate-only。

```yaml
logout_revocation_active_connection_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0154
decision: ADR-0071
check_rule: runtime.logout_revocation_active_connection_gate
future_policy_owner: runtime/internal/app
future_authentication_service_owner: runtime/internal/app/authentication
future_session_repository_owner: runtime/internal/app/session
future_connection_registry_owner_candidate: runtime/internal/app
future_connection_registry_package_candidate: runtime/internal/app/connection
future_transport_invalidation_interface_owner: runtime/internal/platform/transport/ws
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
existing_logout_contract: contracts/runtime/authentication/commands/LogoutAccessToken.yaml
existing_logout_service_method: runtime/internal/app/authentication.Service.LogoutAccessToken
existing_logout_service_behavior: fail_closed_not_implemented
existing_refresh_service_behavior: refresh_not_supported
first_recommended_logout_scope: presented_access_token_only
token_revocation_and_session_revocation_separate: true
connection_registry_required_before_targeting_active_connections: true
logout_execution_added: false
token_revocation_execution_added: false
runtime_session_revocation_execution_added: false
active_connection_invalidation_added: false
connection_registry_added: false
websocket_close_policy_added: false
transport_credential_carrier_added: false
websocket_handshake_authentication_added: false
protobuf_logout_route_added: false
protobuf_session_carrier_added: false
existing_protobuf_envelope_change_added: false
reconnect_epoch_behavior_added: false
cleanup_jobs_added: false
dependencies_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

本标准不实现 logout 或 revocation behavior。

## 2. 所有权

未来 logout 和 revocation policy 必须保持 application-owned：

```yaml
future_policy_owner: runtime/internal/app
future_logout_execution_owner: runtime/internal/app/authentication
future_token_record_owner: runtime/internal/modules/authentication
future_session_repository_owner: runtime/internal/app/session
future_connection_registry_owner_candidate: runtime/internal/app
transport_owner: runtime/internal/platform/transport/ws
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
domain_handler_owner: runtime/internal/modules/*
postgres_adapter_owner: runtime/internal/platform/persistence/postgres
```

规则：

- Application layer 拥有 token、runtime session 或 active connection 必须 invalidated 的决策。
- Authentication service 只有在后续 implementation slice 授权后，才可以执行 token revocation。
- Session repository 只有在后续 implementation slice 授权后，才可以执行 runtime session revocation。
- 未来 connection registry 必须暴露窄的 application-owned invalidation target；它不能让 WebSocket transport 成为 authentication state 的 owner。
- WebSocket transport 只能通过由 application policy 请求的窄 transport-owned capability 关闭具体 socket。
- Protobuf adapters 只有在 protocol route gate 之后才可以映射 public logout 或 invalidation result；它们不得决定 revocation policy。
- Domain modules 永远不得解析 raw tokens、session ids、connection ids、WebSocket close metadata 或 transport credential carriers 来决定 logout 或 revocation behavior。

## 3. 未来策略问题

后续 implementation gate 必须在代码改动前回答这些问题：

```yaml
future_policy_questions:
  logout_scope:
    - Logout 是否只 revoke presented access token？
    - Logout 是否也 revoke linked runtime session？
    - Logout 是否支持 player 或 credential family 的 all sessions？
  revocation_source:
    - Revocation 来源是 user-requested logout、admin action、expiration cleanup、credential compromise，还是 account disablement？
  active_connection_targeting:
    - Server 能否按 token record id 找到 active connections？
    - Server 能否按 runtime session id 找到 active connections？
    - Server 能否按 player id 找到 active connections？
  invalidation_effect:
    - Active connection 是否应该立即关闭？
    - Connection 是否保持打开，但 protected routes 失败？
    - Server 是否先发送 system message 再 close？
  transport_boundary:
    - 哪一层选择 WebSocket close reason 和 code？
    - 哪一层执行具体 close？
  route_policy_boundary:
    - Revocation 后哪些 route families 受到影响？
    - Request-token validation alone 是否会在 dispatch 前捕获 revocation？
    - Bound-connection routes 是否需要 active registry check？
  reconnect_epoch_boundary:
    - Revocation 是否推进 epoch？
    - Revocation 是否阻止 resume？
  observability_boundary:
    - 后续需要哪些 counters、audit events 和 redacted logs？
```

规则：

- 这些问题不得在 WebSocket handlers 中被隐式回答。
- Metadata-only `player_id`、`session_id`、`connection_id` 或 `connection_epoch` 不足以 target 或 authorize revocation。
- Connection 只能从 server-observed binding/session/token state 被 target，绝不能只依赖 client-asserted metadata。
- Public errors 不得暴露失败原因是 token lookup、token revocation、session revocation、active connection lookup，还是 connection close。

## 4. 推荐的未来第一姿态

推荐的未来第一姿态是保守的：

```yaml
recommended_future_first_posture:
  logout_scope: presented_access_token_only
  runtime_session_revocation: linked_session_policy_deferred_until_implementation_gate
  active_connection_invalidation: policy_defined_before_implementation
  connection_registry_required_before_targeting_active_connections: true
  websocket_transport_auth_state_owner: false
  close_active_socket_on_logout_default: not_selected_by_this_gate
  bound_route_policy_reclassification: not_changed_by_this_gate
  reconnect_epoch_interaction: deferred
  public_logout_route_protocol: deferred
  direct_nakama_pitaya_api_compatibility: false
```

规则：

- 第一版未来 logout implementation 应只 revoke presented access token，除非后续 ADR 明确扩大 scope。
- 即使 session row 链接到 access-token record，runtime session revocation 也应该是独立 policy choice。
- Active connection invalidation 不能是隐藏的 best-effort 行为。后续 implementation 必须说明 invalidating active connections 失败时 logout 是失败、带 warning 成功，还是异步重试。
- 在 server 可以按 player、session、token、connection id 或 epoch 定位已打开 socket 之前，必须先有 future connection registry。
- Ordinary protected routes 默认仍然由 request-token 保护；本 gate 不重新分类 routes。

## 5. 未来事件和状态词汇

未来 implementation 只有在后续 implementation gate 后才可以使用这些词汇：

```yaml
candidate_internal_state_transitions:
  token:
    active_to_revoked: authentication_access_tokens.status
  runtime_session:
    active_to_revoked: runtime_sessions.session_status
  active_connection:
    bound_to_invalidated: future_connection_registry_state
    bound_to_closed: future_transport_close_result

candidate_internal_reasons:
  - logout_presented_access_token
  - token_revoked_by_policy
  - runtime_session_revoked_by_policy
  - player_account_disabled
  - credential_compromised
  - admin_session_revocation
  - duplicate_connection_replacement
```

规则：

- Token revocation、session revocation 和 active-connection invalidation 是不同 state transitions。
- 后续 implementation 必须明确选择 transaction boundaries。尤其要说明 token revocation 和 runtime session revocation 是否在同一个 unit of work 中发生。
- Active connection close 在 SQL transaction control 之外，必须被建模为 application/transport side effect，并带 redacted outcome handling。
- 未来 events 不得包含 raw token text、raw credential material、lookup digests、verifier digests、verifier key ids、Authorization headers、cookies、query strings、WebSocket subprotocol values 或 inner payload bytes。

## 6. 错误和脱敏边界

未来 public behavior 必须稳定且脱敏：

```yaml
candidate_public_errors:
  logout_token_missing: AUTHENTICATION_TOKEN_MISSING
  logout_token_malformed: AUTHENTICATION_TOKEN_MALFORMED
  logout_token_invalid: AUTHENTICATION_TOKEN_INVALID
  logout_unavailable: AUTHENTICATION_TOKEN_STORE_UNAVAILABLE
  session_invalid: SESSION_INVALID
  active_connection_invalidation_unavailable: CONNECTION_INVALIDATION_UNAVAILABLE
```

规则：

- Public logout 和 revocation failures 不得暴露 token 是 unknown、already revoked、expired、linked to missing session、linked to closed connection，还是 repository failure。
- Active connection invalidation errors 不得泄露 connection ids、session ids、token record ids、player ids、remote addresses、close codes 或 transport-specific internals。
- Logs 和 events 只能在 observability standard 定义什么是 log-safe 后，使用 redacted ids 或明确分类过的 internal ids。
- `verifier_key_id`、lookup digest、verifier digest、raw token text 和 raw credential material 仍然不是 log-safe。
- WebSocket close reason text 必须保持 generic，除非未来 protocol/transport close policy standard 授权更具体的 client-visible value。

## 7. 和现有 runtime 组件的关系

本 gate 不改变现有行为：

```yaml
authentication_service_logout_changed: false
authentication_service_refresh_changed: false
access_token_validation_changed: false
route_policy_changed: false
connection_binding_changed: false
runtime_session_validation_changed: false
session_creation_changed: false
session_repository_changed: false
postgres_adapter_changed: false
websocket_transport_changed: false
protobuf_protocol_changed: false
```

规则：

- `LogoutAccessToken` 可以继续 fail-closed，直到未来 implementation gate。
- `RefreshAccessToken` 可以继续 unsupported。
- Access-token validation 在被 per request 调用时，可以继续通过现有 record status checks 拒绝 revoked token。
- Bound connection identity 不等于 active connection registry。
- Runtime session repository revocation methods 在未来 behavior slice 使用之前仍然只是 storage-neutral capabilities。
- Route policy 继续把 ordinary protected routes 默认分类为 request-token required。
- 本 gate 不授权 session last-seen、cleanup、audit mutation 或 connection replacement behavior。

## 8. 和 WebSocket 及 Protocol 的关系

本 gate 不授权 WebSocket 或 Protobuf changes：

```yaml
websocket_transport_credential_neutral: true
websocket_close_policy_added: false
transport_credential_carrier_added: false
websocket_handshake_authentication_added: false
protobuf_logout_route_added: false
protobuf_session_carrier_added: false
existing_protobuf_envelope_change_added: false
generated_output_added: false
```

规则：

- WebSocket transport 不得为本 gate 解析 Authorization headers、bearer values、cookies、query-string tokens、session tokens 或 `Sec-WebSocket-Protocol` authentication material。
- 本 gate 不授权 logout command Protobuf source、generated output、authentication response session fields、revocation system message、close reason enum、reconnect message 或 envelope field。
- 未来 protocol logout route gate 必须在任何 client-visible logout command 暴露之前授权。
- 未来 transport close policy gate 必须授权任何 custom WebSocket close code 或 close reason。
- 未来 protocol session carrier gate 必须授权任何 client-visible session id、session proof、resume token 或 connection epoch carrier。

## 9. 延后事项

本 gate 不授权：

- 实现 `LogoutAccessToken`。
- 实现 token revocation execution。
- 实现 runtime session revocation execution。
- 实现 logout-all-sessions、account-wide revocation、credential-wide revocation、admin revocation 或 credential compromise workflows。
- 在 logout 或 revocation 时关闭 WebSocket connections。
- 添加 active connection registry。
- 把 connection registry 接入 frame handling。
- 添加 kick、disconnect、duplicate replacement、reconnect、resume 或 epoch behavior。
- 添加 WebSocket close codes、close reason policy 或 close system messages。
- 添加 Protobuf logout routes、session carriers、resume carriers、generated output 或 existing envelope changes。
- 添加 WebSocket handshake authentication 或 transport credential carriers。
- 添加 cleanup jobs、async invalidation workers、queues、metrics、dashboards、admin APIs 或 observability dependencies。
- 把 ordinary protected routes 从 request-level token proof 重新分类走。
- 添加 memory durable session behavior。
- 添加 direct Nakama 或 Pitaya public API compatibility。

## 10. 未来实现的测试要求

后续 implementation 必须包含聚焦测试：

- Missing、malformed、unknown、expired、already-revoked 和 active tokens collapse 到选定 public outcomes。
- Presented-token logout 只 revoke 所选 scope。
- Runtime session revocation behavior 要么明确 deferred，要么用测试展示 transaction boundary。
- Active connection invalidation 要么明确不选择，要么通过 narrow registry/transport interface 实现。
- Connection targeting 永远不依赖 metadata-only identity。
- Bound/session route policies 在 revocation 后继续 fail closed。
- 关闭 already-closed 或 missing connection 失败时遵循选定 policy。
- Logout side effects 不在 token proof validated 之前发生。
- Public errors、logs 和 events 保持 redacted。
- WebSocket transport 保持 credential-neutral。
- Protobuf envelope 和 authentication response shapes 保持不变，除非单独 protocol gate 授权 changes。

除非后续 implementation 改变 persistence behavior 且需要真实数据库检查，否则 live PostgreSQL verification 可以继续 opt-in。

## 11. Nakama 和 Pitaya 参考映射

Nakama reference mapping：

```yaml
adopted_concepts:
  - authenticated_session_material_has_logout_revocation_and_expiration_lifecycle_pressure
  - revoked_or_expired_credentials_must_not_authorize_gameplay_requests
  - realtime_socket_behavior_must_be_considered_when session_lifecycle_changes
adapted_concepts:
  - vibit_keeps_presented_token_logout_runtime_session_revocation_and_active_socket_invalidation_as_separate_policy_decisions
  - vibit_requires_connection_registry_policy_before_targeting_open_sockets
  - vibit_does_not_copy_nakama_session_api_jwt_shape_or_realtime_socket_contract
deferred_concepts:
  - refresh_token_flow
  - logout_all_sessions
  - session_management_api
  - admin_session_revocation_surface
  - dashboard_session_operations
rejected_for_now:
  - direct_nakama_api_compatibility
```

Pitaya reference mapping：

```yaml
adopted_concepts:
  - acceptors_sessions_and_handlers_are_separate_lifecycle_surfaces
  - session_binding_can_support_targeted_connection_management
  - handler_logic_should_receive_context_not_parse_transport_credentials
adapted_concepts:
  - vibit_future_connection_registry_is_application_runtime_owned_with_narrow_transport_close_handoff
  - vibit_keeps_socket_close_policy_out_of_authentication_repositories
  - vibit_defers_cluster_session_routing_until_single_process_boundaries_are_stable
deferred_concepts:
  - frontend_backend_cluster_session_invalidation
  - distributed_kick_or_disconnect_routing
  - groups_broadcast_integration
  - server_to_server_rpc_invalidation
rejected_for_now:
  - direct_pitaya_api_compatibility
```

Nakama 和 Pitaya 用于指导 capability planning。它们不覆盖 vibit 的 constitution、ADRs、manifests、generated boundaries 或 verification commands。

## 12. 未来实现队列

在这个 gate 之后，未来工作应该继续拆开：

```yaml
future_work_items:
  logout_access_token_behavior_gate:
    may_define:
      - presented-token logout execution
      - token revocation transaction boundary
      - whether linked runtime session revocation is included
  active_connection_registry_gate:
    may_define:
      - server-observed connection registry ownership
      - lookup keys and redaction
      - narrow transport invalidation handoff
  logout_revocation_active_connection_implementation:
    may_add:
      - token or session revocation execution
      - active connection invalidation only if registry and close policy are already selected
  reconnect_connection_epoch_gate:
    may_define:
      - reconnect, resume, duplicate replacement, and epoch mismatch behavior
  protocol_logout_route_gate:
    may_define:
      - client-visible logout command route and Protobuf messages
  protocol_session_carrier_gate:
    may_define:
      - whether and how clients receive or carry session ids, resume tokens, or connection epochs
```

不要在没有新 ADR 的情况下把这些合并成一个宽泛的 connection/session subsystem slice。

## 13. 验证

本 gate 的 repository check rule 是：

```text
runtime.logout_revocation_active_connection_gate
```
