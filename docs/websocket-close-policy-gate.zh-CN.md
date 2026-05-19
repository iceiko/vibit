# WebSocket 关闭策略 Gate

状态：Draft v0.1
最后更新：2026-05-18
范围：未来 WebSocket close、kick、disconnect 策略行为的 gate-only 边界
依赖：`docs/active-connection-registry-gate.md`、`decisions/ADR-0075-active-connection-registry-single-process-implementation.md`、`docs/logout-revocation-active-connection-gate.md`、`docs/bound-identity-route-policy-gate.md`、`docs/reference-game-server-alignment.md`
规范决策：`ADR-0076`

配对英文源文档是 `docs/websocket-close-policy-gate.md`。英文文件是权威版本。

## 1. 目的

Active connection registry 现在可以表达 server-observed connection state 和 validated identity linkage，但它刻意不能关闭 socket。vibit 需要先定义这个分离点，之后才能决定 logout、token revocation、session revocation、duplicate connection policy、reconnect policy、admin action 或 operational drain 命中 active WebSocket connection 时应该发生什么。

Nakama 给出的产品压力是：realtime sockets、authenticated sessions、logout、expiration、single-socket 这类策略都属于玩家可感知的生命周期问题，必须行为可预测。Pitaya 给出的分层压力是：acceptors、agents/sessions、handlers、connection management 是不同表面，handler 不应该以隐藏业务逻辑直接关闭网络连接。

vibit 吸收这些经验：在任何代码把 registry invalidation 转成 concrete socket close behavior 之前，必须先有 application-owned WebSocket close policy boundary。本标准只定义 gate。

```yaml
websocket_close_policy_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0164
decision: ADR-0076
check_rule: runtime.websocket_close_policy_gate
future_policy_owner: runtime/internal/app
future_registry_owner: runtime/internal/app/connection
future_transport_handoff_owner: runtime/internal/platform/transport/ws
future_protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
first_close_policy_posture: application_owned_policy_before_transport_handoff
transport_close_handoff_posture: deferred
registry_invalidation_to_close_default: not_selected_by_this_gate
logout_close_socket_default: not_selected_by_this_gate
kick_disconnect_behavior_added: false
websocket_close_implementation_added: false
close_code_mapping_added: false
close_reason_text_added: false
protocol_close_message_added: false
protobuf_logout_route_added: false
protobuf_session_carrier_added: false
reconnect_epoch_behavior_added: false
runtime_session_revocation_added: false
transport_credential_carrier_added: false
websocket_handshake_authentication_added: false
dependencies_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

本标准不添加 close behavior。

## 2. 所有权

未来 WebSocket close policy 必须由 application 拥有：

```yaml
future_close_policy_owner: runtime/internal/app
future_registry_owner: runtime/internal/app/connection
future_transport_handoff_owner: runtime/internal/platform/transport/ws
future_protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
authentication_service_owner: runtime/internal/app/authentication
session_repository_owner: runtime/internal/app/session
domain_handler_owner: runtime/internal/modules/*
```

规则：

- application layer 拥有连接是否应该 close、disconnect、kick、drain、invalidate 或保持 open 的策略决策。
- WebSocket transport 只能在未来 narrow concrete close handoff 中执行 application policy 产生的 redacted close intent。
- WebSocket transport 不得解析 credential，不得决定 identity、logout policy、session policy、route policy，也不得选择 player-facing close reason text。
- Authentication service 可以 revoke tokens，但不得直接关闭 socket，也不得拥有 transport close handoff。
- Active connection registry 可以把记录标记为 invalidated 或 closed，但 registry listing 和 invalidation 本身不是 close policy。
- Protobuf adapter 只有在后续 implementation gate 定义精确 route/payload 行为后，才可以调用 application close policy。
- Domain modules 不得直接 close、kick、disconnect 或 target WebSocket connections。

## 3. 未来 Close Intent 词汇

后续 implementation gate 必须选择精确 Go 类型，但第一版策略词汇应该把 intent、target、transport action 和 outcome 分开：

```yaml
candidate_close_intent:
  target:
    - connection_id_and_epoch
    - player_id
    - runtime_session_id
    - access_token_record_id
  reason_class:
    - token_revoked
    - logout_presented_token
    - session_revoked
    - duplicate_connection_policy
    - server_shutdown_or_drain
    - policy_violation
    - administrative_action
    - protocol_error
    - idle_timeout
    - unknown_internal
  transport_action:
    - close_socket
    - mark_invalidated_only
    - send_system_message_then_close
    - defer_to_reconnect_policy
  retryability:
    - retryable
    - not_retryable
    - unknown
  public_visibility:
    - silent
    - generic_disconnect
    - generic_reauth_required
```

规则：

- Close intent 必须来自 server-owned registry records 和 validated application identity，不能来自 client metadata alone。
- Close reason class 是 internal semantic class，不是 raw WebSocket close reason text。
- 未来 player-visible close text、system messages 或 protocol errors 必须单独 ratify。
- Registry invalidation 和 concrete socket close 是分离动作，直到 implementation gate 选择策略。
- Close policy 必须定义关闭 socket 失败时 caller 是否失败、是否带 redacted warning 成功、是否记录 retryable action，或是否忽略。
- Close policy 必须定义 target by connection id/epoch、player id、runtime session id、access-token record id，还是组合策略。

## 4. Close Code 和 Reason 边界

本 gate 不选择 WebSocket close codes 或 reason strings：

```yaml
close_code_mapping_added: false
close_reason_text_added: false
custom_close_codes_selected: false
websocket_status_code_dependency_added: false
player_visible_close_reason_added: false
```

规则：

- 本 gate 不授权任何 custom WebSocket close code。
- 本 gate 不授权 close reason text、kick reason、disconnect reason 或 player-facing system message。
- 未来实现必须把 internal reason classes 映射成 redacted close behavior，不能泄漏 raw token material、verifier material、session ids、internal repository ids、remote addresses、headers、cookies、query strings、subprotocol values 或 database errors。
- Transport-level abnormal close、peer disconnect、network failure 和 application-directed close 必须能在 internal records 里区分，但 public behavior 必须 redacted。

## 5. 与 Active Connection Registry 的关系

Registry 是 target model，不是 close policy：

```yaml
registry_lookup_allowed_in_future_policy: true
registry_invalidation_to_close_default: not_selected_by_this_gate
registry_mark_closed_from_transport_lifecycle: future_handoff_only
registry_mark_invalidated_from_policy: future_policy_only
```

规则：

- 未来 close policy 可以按 connection id/epoch、player id、runtime session id 或 access-token record id 查询 registry。
- Listing registry records 不得决定 socket 是否应该 close。
- Marking a registry record invalidated 不得关闭 concrete socket，除非后续 implementation gate 明确 wire 该动作。
- Transport-observed peer close 和 application-directed close 是不同 lifecycle facts，不能混为一谈。
- Duplicate connection replacement 需要单独 reconnect/epoch 或 duplicate-policy gate。

## 6. 与 Logout 和 Session Revocation 的关系

本 gate 不改变 logout 或 session behavior：

```yaml
logout_access_token_behavior_changed: false
token_revocation_behavior_changed: false
runtime_session_revocation_added: false
logout_close_socket_default: not_selected_by_this_gate
session_revocation_close_socket_default: not_selected_by_this_gate
admin_kick_default: not_selected_by_this_gate
```

规则：

- `LogoutAccessToken` 继续只 revoke verified presented access-token record。
- 现有 request-level access-token validation 在后续 protected request 运行时会拒绝 revoked material。
- 未来策略必须决定 token revocation 是按 access-token record id、runtime session id、player id 还是 connection id 查 active connections。
- 未来策略必须决定 socket close 失败时 logout 是否仍然成功。
- 未来策略必须决定 runtime session revocation 是关闭 active sockets、只 invalidate registry records，还是交给 reconnect/session validation policy。

## 7. 与 WebSocket 和 Protocol 的关系

本 gate 不改变 WebSocket 或 Protobuf behavior：

```yaml
websocket_transport_credential_neutral: true
websocket_close_implementation_added: false
transport_close_handoff_added: false
transport_credential_carrier_added: false
websocket_handshake_authentication_added: false
protobuf_logout_route_added: false
protobuf_session_carrier_added: false
protocol_close_message_added: false
existing_protobuf_envelope_change_added: false
generated_output_added: false
```

规则：

- WebSocket transport 不得为本 gate 解析 Authorization headers、bearer values、cookies、query-string tokens、session tokens 或 `Sec-WebSocket-Protocol` authentication material。
- 本 gate 不授权 close code、close reason、kick/disconnect message、logout route、session carrier、reconnect token、resume token 或 connection epoch protocol field。
- Generated Go Protobuf output 不得因本 gate 改变。
- 未来 close handoff 必须足够窄：transport 只关闭 concrete socket，不拥有 application policy。

## 8. 未来测试期望

后续 implementation gate 必须要求聚焦测试：

- 只从 server-owned registry records 产生 close intents。
- 拒绝 metadata-only player id、session id、token record id 或 client-supplied connection metadata 作为 close targets。
- Close decision policy 留在 application layer。
- Concrete socket close mechanics 留在 narrow transport handoff。
- Registry invalidation 和 concrete socket close 保持分离，除非明确选择。
- 把 internal reason classes 映射成 redacted public behavior。
- 除非未来策略另有选择，否则 socket close 失败不改变 logout/token revocation 行为。
- 保持 WebSocket transport credential neutrality。
- 除非单独授权，否则不添加 generated Protobuf changes 和 protocol close messages。

## 9. Nakama 和 Pitaya 参考映射

Nakama reference mapping：

- 采纳 realtime socket lifecycle、session expiration、logout 和 single-socket 风格策略需要显式 lifecycle semantics 的经验。
- 改造为 vibit-owned close policy vocabulary，未来可协调 registry records、token revocation、runtime sessions 和 reconnect behavior。
- 延后 direct Nakama session APIs、JWT/session token shape、realtime socket compatibility、dashboard operations、cluster session routing 和 exact close behavior。

Pitaya reference mapping：

- 采纳 acceptors、agents/sessions、route handlers、connection management 分离的经验。
- 改造为 application-owned close policy 加未来 narrow transport handoff。
- 延后 frontend/backend cluster routing、distributed kick/disconnect、groups/rooms integration 和 server-to-server RPC invalidation。

## 10. 非目标

本 gate 不授权：

- Go runtime WebSocket close implementation。
- Transport close handoff code。
- Close codes、close reason strings、kick messages、disconnect messages 或 protocol close messages。
- Logout-triggered socket close。
- Runtime session revocation-triggered socket close。
- Admin kick/disconnect behavior。
- Duplicate connection replacement。
- Reconnect、resume、durable epoch behavior 或 reconnect token behavior。
- Protocol logout routes、protocol session carriers、protocol close messages 或 existing envelope changes。
- WebSocket handshake authentication 或 transport credential carriers。
- PostgreSQL、Redis-like、distributed 或 durable active connection storage。
- 主要新依赖。
- Direct Nakama 或 Pitaya public API compatibility。

## 11. Next Gate

此 gate 完成后，work queue 必须停在新的 confirmation point。建议的下一个方向是：

```yaml
candidate_next_directions:
  - implement_websocket_close_policy_single_process
  - define_protocol_logout_route_gate
  - define_reconnect_connection_epoch_gate
  - define_protocol_session_carrier_gate
  - strengthen_operations_observability_and_admin_tooling
  - expand_core_game_backend_modules_after_nakama_pitaya_review
```

保守建议是 `implement_websocket_close_policy_single_process`，但必须由 maintainer 明确选择。
