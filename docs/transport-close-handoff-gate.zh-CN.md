# Transport Close Handoff Gate

状态：Draft v0.1
最后更新：2026-05-18
范围：未来 application-to-WebSocket concrete close handoff 的 gate-only boundary
依赖：`docs/websocket-close-policy-gate.md`、`decisions/ADR-0077-websocket-close-policy-single-process-implementation.md`、`docs/protocol-logout-route-gate.md`、`docs/nakama-pitaya-product-parity-roadmap.md`、`docs/runtime-protocol-adapter.md`、`docs/reference-game-server-alignment.md`
Canonical decision：`ADR-0080`

英文源文件是 `docs/transport-close-handoff-gate.md`，英文文件是权威版本。本文件是简体中文翻译。

## 1. 目的

vibit 现在有三块已经分开的 lifecycle 基础：

- Active connection registry 记录服务端观察到的 connection state。
- Application close policy 可以解析 active bound records，并产出 redacted close intents。
- Protocol logout route 可以 revoke presented access-token record，但不关闭 socket。

缺失的部分，是从 application-owned close policy 到具体 WebSocket socket close mechanics 的未来窄 handoff。如果没有 gate，后续代码很容易让 WebSocket transport 拥有 authentication policy，让 logout 隐式关闭 socket，用客户端提供的 metadata 作为 close authority，或把 session/reconnect 决策藏进 protocol handlers。

Nakama 是 session、logout、realtime socket 和 server-directed disconnect 等显式 lifecycle behavior 的产品参考。Pitaya 是 acceptors、sessions、route handlers、groups、RPC 和 kick/disconnect 式 connection management 分层的架构参考。vibit 吸收这些经验：policy 留在 application layer，transport 只在收到 server-owned handoff 后执行窄的 concrete close action。

本标准只定义 gate。

```yaml
transport_close_handoff_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0172
decision: ADR-0080
check_rule: runtime.transport_close_handoff_gate
parity_phase: phase_2r_runtime_lifecycle_closure
application_policy_owner: runtime/internal/app/connection
active_connection_registry_owner: runtime/internal/app/connection
future_transport_handoff_owner: runtime/internal/platform/transport/ws
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
authentication_service_owner: runtime/internal/app/authentication
first_handoff_target: connection_id_and_epoch
server_observed_target_required: true
client_metadata_authority_allowed: false
transport_policy_ownership_allowed: false
transport_credential_parsing_allowed: false
transport_session_revocation_allowed: false
first_transport_action_candidate: close_socket
close_code_mapping_added: false
close_reason_text_added: false
protocol_close_message_added: false
logout_triggered_socket_close_added: false
runtime_session_revocation_added: false
reconnect_epoch_behavior_added: false
protocol_session_carrier_added: false
dependencies_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. 所有权

未来 handoff 必须保持这些 ownership boundaries：

```yaml
application_policy_owner: runtime/internal/app/connection
active_connection_registry_owner: runtime/internal/app/connection
future_transport_handoff_owner: runtime/internal/platform/transport/ws
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
authentication_service_owner: runtime/internal/app/authentication
domain_module_owner: runtime/internal/modules/*
```

规则：

- Application close policy 拥有“某个 connection 应该被关闭”的决策。
- Active connection registry 拥有 server-observed connection records 以及 invalidation/closed markers。
- WebSocket transport 只可以在 application policy 交给它 server-owned target 后，拥有 concrete socket close mechanics。
- WebSocket transport 不得解析 credentials、验证 tokens、选择 players、选择 runtime sessions、评估 logout state、决定 reconnect behavior，或选择玩家可见文本。
- Protocol adapters 不得直接关闭 socket，也不得从 client payload metadata 创建 close targets。
- Authentication service 可以 revoke token records，但不得调用 WebSocket transport 或拥有 close handoff。
- Domain modules 不得导入 WebSocket transport 或直接关闭 concrete sockets。

## 3. Handoff Target

第一版未来 handoff target 必须是：

```yaml
first_handoff_target: connection_id_and_epoch
connection_id_source: server_observed_websocket_accept_metadata
connection_epoch_source: server_observed_websocket_accept_metadata
requires_active_registry_record: true
client_supplied_connection_id_authority: false
client_supplied_epoch_authority: false
player_id_transport_authority: false
runtime_session_id_transport_authority: false
access_token_record_id_transport_authority: false
```

规则：

- Future transport handoff 必须用 server-observed `connection_id` 和 `connection_epoch` 定位 concrete sockets。
- Application policy 可以先把 player、runtime session 或 access-token record targets 解析成 concrete connection/epoch targets，再交给 transport。
- Transport 不得按 player id、runtime session id、access-token record id、route identity、request identity、envelope session metadata、headers、cookies、query strings、subprotocol values 或 remote address 关闭连接。
- Connection epoch 必须防止 stale close intent 关闭后来复用同一 connection id 的新 socket。
- 如果 concrete socket 不存在或 epoch 不匹配，失败必须 redacted 且 policy-neutral。

## 4. 未来 Handoff Shape

后续 implementation gate 可以选择精确的 Go types。第一版 vocabulary 应保持很窄：

```yaml
candidate_transport_close_request:
  connection_id: server_observed_connection_id
  connection_epoch: server_observed_connection_epoch
  reason_class: internal_redacted_close_reason_class
  public_visibility: silent_or_generic_disconnect_or_generic_reauth_required
  retryability: retryable_or_not_retryable_or_unknown
  requested_at: server_time

candidate_transport_close_result:
  connection_id: server_observed_connection_id
  connection_epoch: server_observed_connection_epoch
  outcome:
    - close_requested
    - socket_not_found
    - epoch_mismatch
    - already_closed
    - close_failed
  closed_at: server_time_optional
```

规则：

- Request 不得携带 raw access-token material、raw credential material、lookup digests、verifier digests、verifier key ids、headers、cookies、query strings、subprotocol values、remote addresses、database errors 或完整 repository errors。
- `reason_class` 是 internal semantic class，不是 WebSocket close reason string。
- 本 gate 不选择 WebSocket close codes、close reason text、protocol close messages 或 player-facing system messages。
- 未来实现必须定义 transport close errors 是返回给 application policy、记录为 redacted outcomes，还是两者都做。
- 这个 handoff 应尽量可在没有 live network dependency 的情况下测试。

## 5. 与现有 Close Policy 的关系

当前 close policy 产出的 `CloseIntent` 使用：

```yaml
current_transport_action: mark_invalidated_only
concrete_socket_close_added: false
```

未来 handoff work 不得悄悄把现有 `mark_invalidated_only` 行为重新解释为 concrete close。后续 implementation 必须显式选择 application policy 如何产出 close requests：

```yaml
future_options:
  - keep_mark_invalidated_only_and_add_separate_close_socket_action
  - allow_close_policy_to_emit_close_socket_after_gate
  - compose_invalidated_marker_then_transport_close_request
```

规则：

- Registry invalidation 和 concrete socket close 仍然是不同的 lifecycle facts。
- 除非 application policy 选择了 close action，否则不得生成 transport close request。
- Transport-observed peer close 和 application-directed close 必须能在 registry updates 中区分。
- 除非后续 policy 明确定义耦合关系，close concrete socket 失败不得让 token revocation、logout 或 session mutation 表现成成功或失败。

## 6. 与 Logout、Sessions、Reconnect 的关系

本 gate 不改变 logout、session 或 reconnect behavior：

```yaml
logout_triggered_socket_close_added: false
runtime_session_revocation_added: false
session_revocation_close_added: false
duplicate_connection_replacement_added: false
reconnect_epoch_behavior_added: false
protocol_session_carrier_added: false
```

规则：

- `LogoutAccessToken` 仍然只 revoke verified presented access-token record。
- 本 gate 不决定 logout 后是否关闭当前 socket，或关闭与 token record 关联的其它 sockets。
- Runtime session revocation 仍然是单独的未来 behavior。
- Duplicate connection replacement 和 reconnect/resume behavior 仍然是单独的 future gates。
- Protocol session carriers 与 transport close handoff 分离。

## 7. WebSocket 和 Protocol 边界

本 gate 保持 WebSocket transport credential-neutral 和 protocol-neutral：

```yaml
websocket_transport_credential_neutral: true
transport_credential_parsing_added: false
websocket_handshake_authentication_added: false
protocol_close_message_added: false
protobuf_source_added: false
generated_output_added: false
existing_protobuf_envelope_change_added: false
```

规则：

- 本 gate 不授权任何 `.proto` source 或 generated output。
- 本 gate 不得改变现有 Protobuf envelope fields。
- 本 gate 不授权 close route、kick route、disconnect route、system close message 或 admin disconnect protocol surface。
- WebSocket transport 不得读取 Authorization headers、bearer values、cookies、query-string tokens、session tokens 或 subprotocol authentication material。
- Future close reason text 必须保持 redacted，并由单独决策 ratify。

## 8. 后续测试要求

后续 implementation gate 必须要求 focused tests：

```yaml
required_tests:
  - application_policy_remains_owner_of_close_decision
  - transport_handoff_targets_connection_id_and_epoch_only
  - stale_epoch_does_not_close_new_socket
  - missing_socket_returns_redacted_not_found_outcome
  - transport_does_not_parse_credentials_or_tokens
  - protocol_adapter_does_not_close_sockets_directly
  - authentication_service_does_not_call_transport
  - domain_modules_do_not_import_websocket_transport
  - close_reason_class_does_not_leak_secrets_or_internal_ids
  - existing_protobuf_envelope_remains_unchanged
  - logout_route_behavior_remains_token_record_scoped
```

Live end-to-end socket tests 以后可能有价值，但本 gate-only standard 不要求。

## 9. Nakama 和 Pitaya Reference Mapping

Nakama reference mapping：

- 采用其 product lesson：logout、session invalidation、realtime disconnect 和 server-directed disconnect 都是显式 lifecycle concerns。
- 将其转化为 vibit 的 application-owned close policy 和 narrow WebSocket close handoff。
- 推迟 Nakama session API compatibility、realtime socket compatibility、dashboard disconnect behavior、session logout-all behavior 和 cluster session routing。

Pitaya reference mapping：

- 采用其 layering lesson：acceptors、sessions、handlers、groups、RPC 和 kick/disconnect behavior 应保持分离。
- 将其转化为 transport handoff：只有 application policy 选择了 server-observed connection target 后，transport 才关闭 concrete sockets。
- 推迟 Pitaya route naming compatibility、frontend/backend topology、distributed kick/disconnect、group broadcast integration 和 RPC/session propagation。

## 10. Non-Goals

本 gate 不授权：

- Concrete WebSocket socket close implementation。
- Close code mapping。
- Close reason text。
- Player-facing kick/disconnect messages。
- Protocol close messages。
- Protobuf source 或 generated output changes。
- Logout-triggered socket close。
- Runtime session revocation。
- Duplicate connection replacement。
- Reconnect/resume/epoch behavior。
- Protocol session carriers。
- Presence、chat、groups、parties、matchmaking、match runtime、SDK、cluster 或 distributed runtime behavior。
- New dependencies。
- Direct Nakama/Pitaya API compatibility。

## 11. Agent Rules

Agents 必须：

- 在添加 concrete WebSocket close handoff code 前阅读本标准。
- 把 close policy decisions 留在 `runtime/internal/app/connection`。
- 把 concrete socket mechanics 留在 `runtime/internal/platform/transport/ws`。
- 第一版 transport target 只使用 server-owned connection id 和 epoch。
- 保持 tokens、credentials、digests、verifier keys、headers、cookies、query strings、subprotocol values、remote addresses、database errors 和 internal repository errors 的 redaction boundaries。
- 除非后续 protocol gate 授权，否则保持 generated Protobuf output 不变。

Agents 不得：

- 从 authentication service behavior 关闭 socket。
- 从 protocol bridge code 直接关闭 socket。
- 从 domain modules 关闭 socket。
- 把 client-supplied player id、session id、token record id、connection id、epoch 或 envelope metadata 当成 close authority。
- 在本 gate 中添加 close codes、reason text、reconnect behavior、session carriers、operations/admin disconnect、dependencies 或 direct Nakama/Pitaya compatibility。

## 12. Verification

本 gate 需要的 repository checks：

```bash
node -c tools/vibit
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change define-transport-close-handoff-gate --json
node tools/vibit inspect next --json
```

由于本 gate-only standard 不添加 Go runtime behavior，因此不要求 Go runtime tests。
