# Active Connection Registry Gate

状态：Draft v0.1
最后更新：2026-05-18
范围：未来 active WebSocket connection registry behavior 的 gate-only boundary
依赖：`docs/logout-revocation-active-connection-gate.md`、`docs/logout-access-token-behavior-gate.md`、`decisions/ADR-0073-logout-access-token-behavior-implementation.md`、`docs/first-message-connection-binding-implementation-gate.md`、`docs/bound-identity-route-policy-gate.md`、`docs/reference-game-server-alignment.md`
Canonical decision：`ADR-0074`

配套英文原文是 `docs/active-connection-registry-gate.md`。英文文件是权威版本。

## 1. 目的

`LogoutAccessToken` 现在已经可以 revoke 已验证的 presented opaque access-token record，但这个 revocation 不会关闭已经打开的 WebSocket connection，不会 revoke runtime session，也不会让 connection-bound identity 自动失效。这是有意为之。下一个缺失的边界，是 vibit 在把 active connections 当成 server-owned runtime state 处理前，必须先有 registry 边界。

Nakama 给出的产品压力很清楚：authenticated session material、realtime sockets、logout、expiration 和 revocation 是互相关联的 lifecycle concerns。Pitaya 给出的分层压力也很清楚：acceptors、sessions、route handlers 和 connection management 是不同 surface，handlers 不应解析 credentials，也不应拥有 transport lifecycle side effects。

vibit 吸收这些经验：在任何代码能按 player、runtime session、token record、connection id 或 epoch 定位 open sockets 前，必须先定义 application-owned active connection registry boundary。本标准只定义 gate，不添加实现。

```yaml
active_connection_registry_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0160
decision: ADR-0074
check_rule: runtime.active_connection_registry_gate
future_registry_owner: runtime/internal/app/connection
future_policy_owner: runtime/internal/app
future_transport_handoff_owner: runtime/internal/platform/transport/ws
future_protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
first_registry_posture: single_process_in_memory
registry_persistence_posture: non_durable_runtime_state
cluster_registry_posture: deferred
server_observed_connection_id_required: true
connection_epoch_required: true
metadata_only_targeting_allowed: false
bind_connection_integration_candidate: runtime.authentication.BindConnection
logout_revocation_integration_candidate: future_application_policy_only
active_connection_registry_added: false
active_connection_invalidation_added: false
websocket_close_policy_added: false
kick_disconnect_behavior_added: false
duplicate_connection_replacement_added: false
reconnect_epoch_behavior_added: false
runtime_session_revocation_added: false
protobuf_logout_route_added: false
protobuf_session_carrier_added: false
existing_protobuf_envelope_change_added: false
transport_credential_carrier_added: false
websocket_handshake_authentication_added: false
dependencies_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

本标准不添加 registry implementation。

## 2. 所有权

未来 active connection registry behavior 必须由 application 拥有：

```yaml
future_registry_owner: runtime/internal/app/connection
future_policy_owner: runtime/internal/app
future_connection_binding_caller: runtime/internal/platform/protocol/protobuf
future_transport_handoff_owner: runtime/internal/platform/transport/ws
authentication_service_owner: runtime/internal/app/authentication
session_repository_owner: runtime/internal/app/session
domain_handler_owner: runtime/internal/modules/*
```

规则：

- Application layer 拥有 active connection registry state 和 lifecycle policy。
- WebSocket transport 可以报告 server-observed connection open/close facts；只有后续 narrow transport handoff 授权后，才可以关闭 concrete sockets。
- WebSocket transport 不得拥有 player identity、session identity、token validity、logout policy、route policy 或 credential parsing。
- Protobuf adapters 只有在后续 implementation gate 定义 exact handoff 后，才可以调用 application-owned binding 或 registry methods。
- Authentication service 可以 revoke tokens，但不得从 token repository code 直接 close sockets 或 mutate registry state。
- Domain modules 不得直接 read、write、close、kick、disconnect 或 target WebSocket connections。

## 3. 未来 Registry 形状

后续 implementation gate 必须选择具体 Go types，但第一版姿态应使用以下 vocabulary 建模 registry record：

```yaml
candidate_registry_record:
  connection_id: server_observed_connection_id
  connection_epoch: server_observed_epoch
  state:
    - open_unbound
    - bound
    - closing
    - closed
  bound_actor:
    actor_kind: player
    player_id: validated_player_id
  runtime_session_id: optional_validated_runtime_session_id
  access_token_record_id: optional_server_token_record_id
  opened_at: server_clock_time
  bound_at: optional_server_clock_time
  last_seen_at: optional_server_clock_time
  closed_at: optional_server_clock_time
  close_reason_class: optional_redacted_internal_reason
```

规则：

- Registry record 代表 server-observed connection state，不是 client proof。
- `connection_id` 和 `connection_epoch` 是 server-owned identifiers。它们可以帮助关联 transport events，但不是 authentication proof。
- Bound player、runtime session 或 token linkage 只能来自 validated application identity，不能只来自 client-supplied metadata。
- 第一版 registry posture 是 single-process、in-memory、non-durable。Durable connection state、cross-node lookup、Redis-like stores、service discovery 和 server-to-server RPC 继续 deferred。
- Registry state 不得存储 raw access-token text、raw credential material、lookup digests、verifier digests、verifier key ids、Authorization headers、cookies、query strings、WebSocket subprotocol values、remote addresses 或 inner payload bytes。

## 4. 未来 Capabilities

后续 implementation gate 可以定义窄 capabilities，例如：

```yaml
candidate_registry_capabilities:
  - RegisterOpenConnection
  - BindConnectionIdentity
  - MarkConnectionClosed
  - FindConnectionByID
  - ListConnectionsByPlayerID
  - ListConnectionsByRuntimeSessionID
  - ListConnectionsByAccessTokenRecordID
  - MarkConnectionInvalidated
```

规则：

- Registration 必须由 server-observed transport lifecycle events 驱动。
- Binding 必须要求来自现有 access-token 和可选 runtime-session validation path 的 validated application identity。
- 按 player、session 或 token record id listing 只是 targeting primitive。它本身不得决定 close policy。
- Invalidating registry record 和 closing WebSocket 是分离的 future actions。
- 未来任何 target connections 的 public route 或 admin surface，都必须经过 application policy 和 redacted authorization checks。

## 5. 与 Logout 和 Revocation 的关系

这个 gate 存在是因为 logout/revocation 需要安全的 target model，但它不实现 active invalidation：

```yaml
logout_access_token_behavior_changed: false
token_revocation_behavior_changed: false
runtime_session_revocation_added: false
active_connection_invalidation_added: false
websocket_close_policy_added: false
logout_close_socket_default: not_selected_by_this_gate
logout_registry_lookup_default: not_selected_by_this_gate
```

规则：

- `LogoutAccessToken` 继续只 revoke verified presented access-token record。
- 现有 request-level access-token validation 在后续 protected requests 运行时，会拒绝 revoked material。
- 后续 policy 必须决定 token revocation 应按 access-token record id、runtime session id、player id 还是 connection id 查找 active connections。
- 后续 policy 必须决定 active connection invalidation failure 是让 logout 失败、带 redacted warning 成功，还是异步重试。
- 后续 policy 必须决定 bound connection 是保持 open 但 protected routes 失败、先收到 system message、立即关闭，还是等待 reconnect/epoch handling。

## 6. 与 WebSocket 和 Protocol 的关系

本 gate 不改变 WebSocket 或 Protobuf behavior：

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

- WebSocket transport 不得为了本 gate 解析 Authorization headers、bearer values、cookies、query-string tokens、session tokens 或 `Sec-WebSocket-Protocol` authentication material。
- 本 gate 不授权 close code、close reason、kick/disconnect message、logout route、session carrier、reconnect token、resume token 或 connection epoch protocol field。
- Generated Go Protobuf output 不得因本 gate 改变。
- 现有 first-message `BindConnection` route 仍是唯一已选择的 connection identity route。本 gate 不添加 registry-backed route-policy use of bound identity。

## 7. 未来测试预期

后续 implementation gate 必须要求 focused tests 覆盖：

- 从 server-observed connection id 和 epoch 注册 open connection。
- 在 replacement policy 被明确选择前，拒绝同一 connection id 和 epoch 的 duplicate active records。
- 只有在 validated player identity 可用后才允许 binding。
- 拒绝 metadata-only player id、session id、token record id 或 client-supplied connection metadata 作为 targeting proof。
- 从 server-observed transport lifecycle events 标记 close state。
- 按 player id、runtime session id 和 access-token record id listing active connections，同时不泄漏 raw proof material。
- 保持 WebSocket transport credential neutrality。
- 在未单独授权前，把 invalidation policy 和 concrete socket close behavior 留在 registry 外部。
- 从 errors 和 logs 中 redacts raw token text、lookup digests、verifier digests、verifier key ids、remote addresses、headers、cookies、query strings 和 subprotocol values。

## 8. Nakama 和 Pitaya 参考映射

Nakama reference mapping：

- 采用 realtime sockets 与 authenticated session material 需要 coordinated lifecycle semantics 的经验。
- 将该经验适配为 vibit-owned registry，用于跟踪 server-observed active connections 和 validated identity linkage。
- 推迟 direct Nakama session APIs、JWT/session token shape、realtime socket compatibility、dashboard operations 和 cluster session routing。

Pitaya reference mapping：

- 采用 acceptors、sessions、handlers 和 connection management 分离的经验。
- 将 Pitaya-style session/context separation 适配为 application-owned registry state 加 narrow transport lifecycle handoff。
- 推迟 frontend/backend cluster routing、distributed kick/disconnect、groups/rooms integration 和 server-to-server RPC invalidation。

## 9. 非目标

本 gate 不授权：

- Go runtime registry implementation。
- WebSocket socket close、kick、disconnect 或 close reason behavior。
- Runtime session revocation。
- Logout-all、admin revocation、refresh、cleanup jobs 或 token rotation。
- Reconnect、resume、duplicate replacement 或 durable epoch behavior。
- Protocol logout routes、protocol session carriers、protocol close messages 或 existing envelope changes。
- WebSocket handshake authentication 或 transport credential carriers。
- PostgreSQL、Redis-like、distributed 或 durable active connection storage。
- Major new dependencies。
- Direct Nakama 或 Pitaya public API compatibility。

## 10. 下一个 Gate

完成本 gate 后，work queue 必须停在新的 confirmation point。推荐候选方向：

```yaml
candidate_next_directions:
  - implement_active_connection_registry_single_process
  - define_websocket_close_policy_gate
  - define_protocol_logout_route_gate
  - define_reconnect_connection_epoch_gate
  - define_protocol_session_carrier_gate
  - strengthen_operations_observability_and_admin_tooling
  - expand_core_game_backend_modules_after_nakama_pitaya_review
```

保守建议是 `implement_active_connection_registry_single_process`，但必须显式选择。
