# Pitaya 对齐的会话绑定、踢出/断开与会话数据边界闸门

状态：Accepted v0.1
最后更新：2026-06-01
范围：在 acceptor 与 connection lifecycle 后续方向选择之后，仅定义 Pitaya 对齐的 session binding、kick/disconnect 与 session data 词汇边界
依赖：`decisions/ADR-0175-select-next-pitaya-aligned-direction-after-acceptor-connection-lifecycle-map.md`、`decisions/ADR-0174-pitaya-aligned-acceptor-connection-lifecycle-source-first-map.md`、`docs/pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.md`、`docs/runtime-protocol-adapter.md`、`docs/game-protocol.md`、`docs/reference-game-server-alignment.md`、`.arch/reference.yaml`
规范决策：`ADR-0176`

英文文件 `docs/pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate.md` 是权威版本；本文是对应的简体中文翻译。

本文只定义 session binding、kick/disconnect 与 session data 的词汇闸门。它不实现 session binding behavior、kick/disconnect behavior、session data behavior、session data persistence、acceptor behavior、TCP acceptors、WebSocket behavior changes、connection lifecycle behavior changes、route handler implementation、handler routing behavior、handler pipeline behavior、pipeline middleware behavior、backend route targeting、cluster-safe session routing behavior、distributed session routing、distributed runtime behavior、distributed groups、room broadcast fanout、delivery guarantees、stream subscriptions、service discovery implementation、service registries、service selectors、node identity、server-to-server RPC、remote calls、frontend/backend server roles、runtime endpoint behavior、metrics endpoints、observability pipelines、dashboards、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、authentication/session behavior changes、SDK publication、hosted deployments、release artifacts 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Pitaya 对齐的 session binding、kick/disconnect 与 session data 边界闸门记录如下：

```yaml
pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate: defined
completed_work_item: W-0268
decision: ADR-0176
check_rule: runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate
previous_direction_decision: ADR-0175
acceptor_connection_lifecycle_source_first_map_decision: ADR-0174
acceptor_connection_lifecycle_source_first_map_check_rule: runtime.pitaya_aligned_acceptor_connection_lifecycle_source_first_map
standard: docs/pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate.md
translation: docs/pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: session_binding_kick_disconnect_session_data_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_session_binding_kick_disconnect_session_data_vocabulary
future_implementation_work_item: W-0269
future_implementation_direction: implement_pitaya_aligned_session_binding_kick_disconnect_session_data_source_first_map
allowed_session_binding_kick_disconnect_session_data_vocabulary:
  - session_binding_boundary
  - connection_bound_session
  - session_data
  - session_data_scope
  - server_initiated_disconnect
  - server_initiated_kick
  - session_unbind
  - session_close_reason
  - connection_session_handoff
  - presence_session_handoff
current_single_process_session_binding_kick_disconnect_session_data_mapping:
  first_message_binding:
    current: authentication_bind_connection_route
    future_vocabulary: session_binding_boundary
    implementation_status: no_handshake_authentication_or_reconnect_binding_change
  runtime_session_validation:
    current: request_level_access_token_validation_and_request_identity_handoff
    future_vocabulary: connection_bound_session
    implementation_status: no_session_persistence_or_every_request_policy_change
  session_metadata:
    current: request_identity_and_connection_metadata
    future_vocabulary: session_data
    implementation_status: no_session_data_store_or_public_api
  session_data_scope:
    current: no_general_session_data_scope
    future_vocabulary: session_data_scope
    implementation_status: planning_vocabulary_only
  active_connection_registry:
    current: application_owned_connection_registry
    future_vocabulary: connection_session_handoff
    implementation_status: no_cluster_safe_session_location_registry
  logout_disconnect_handoff:
    current: logout_service_revokes_token_and_transport_close_policy_remains_unchanged
    future_vocabulary: server_initiated_disconnect
    implementation_status: no_server_initiated_disconnect_behavior
  kick_policy:
    current: no_kick_policy_or_route
    future_vocabulary: server_initiated_kick
    implementation_status: planning_vocabulary_only
  session_unbind:
    current: close_handoff_and_connection_registry_cleanup
    future_vocabulary: session_unbind
    implementation_status: no_remote_unbind_or_reconnect_routing
  close_reason:
    current: existing_transport_close_reason_mapping
    future_vocabulary: session_close_reason
    implementation_status: no_close_policy_change
  presence_lifecycle:
    current: server_owned_presence_snapshot
    future_vocabulary: presence_session_handoff
    implementation_status: no_distributed_presence_session_handoff
session_binding_behavior_added: false
kick_disconnect_behavior_added: false
session_data_behavior_added: false
session_data_persistence_added: false
acceptor_behavior_added: false
tcp_acceptor_added: false
websocket_acceptor_behavior_changed: false
connection_lifecycle_behavior_changed: false
route_handler_implementation_added: false
handler_routing_behavior_added: false
handler_pipeline_behavior_added: false
pipeline_middleware_behavior_added: false
backend_route_targeting_added: false
cluster_safe_session_routing_added: false
session_location_registry_added: false
connection_owner_node_registry_added: false
routing_epoch_behavior_added: false
session_route_target_added: false
remote_connection_handoff_added: false
reconnect_route_added: false
distributed_session_routing_added: false
distributed_group_implementation_added: false
distributed_groups_added: false
group_membership_registry_added: false
room_broadcast_fanout_added: false
broadcast_delivery_guarantee_added: false
stream_subscription_added: false
service_discovery_implementation_added: false
service_registry_added: false
service_selector_added: false
node_registry_added: false
server_identity_added: false
server_to_server_rpc_implementation_added: false
remote_call_behavior_added: false
frontend_server_implementation_added: false
backend_server_implementation_added: false
frontend_backend_server_roles_added: false
distributed_runtime_implementation_added: false
runtime_behavior_added: false
runtime_endpoint_behavior_added: false
metrics_endpoint_added: false
tracing_pipeline_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
migration_added: false
dependency_added: false
authentication_session_behavior_changed: false
hosted_deployment_added: false
sdk_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Purpose

`ADR-0175` 在 acceptor 与 connection lifecycle source-first map 之后，选择了 session binding、kick/disconnect 与 session data 作为下一项 Pitaya 对齐方向。

风险在于，agent 可能把 session、kick、disconnect 或 session data 词汇误读为可以修改 first-message binding、添加 socket kick routes、添加通用 session storage、改变 logout 或 close behavior、添加 reconnect routing，或开始 distributed session routing。这个闸门只记录词汇和当前映射，保持 vibit 现有 single-process WebSocket runtime、request-level authentication、active connection registry、logout behavior、close handoff 和 presence lifecycle behavior 不变。

## 3. Vocabulary

允许使用的词汇：

- `session_binding_boundary`：规划中用于把已验证 runtime identity 绑定到 connection 的边界词汇。当前仍是 first-message route。
- `connection_bound_session`：规划中用于描述与 active connection 关联的 validated session。当前仍是 request-level validation 与 metadata handoff。
- `session_data`：规划中用于 server-owned session metadata。此 slice 不添加通用持久化数据存储。
- `session_data_scope`：规划中用于未来 session data ownership 边界。当前不添加具体 scope。
- `server_initiated_disconnect`：规划中用于 server-originated close intent。当前 close behavior 不变。
- `server_initiated_kick`：规划中用于 policy-driven removal。当前不添加 kick policy 或 route。
- `session_unbind`：规划中用于从 connection 解绑 session。当前 cleanup 仍是现有 close handoff 与 registry cleanup。
- `session_close_reason`：规划中用于分类 server-visible close reason。当前 close reason mapping 不变。
- `connection_session_handoff`：规划中用于在 application boundary 内传递 connection/session association facts。当前不添加 distributed owner registry。
- `presence_session_handoff`：规划中用于连接 session lifecycle facts 与 presence lifecycle。当前 presence 仍是 server-owned snapshot behavior。

禁止用法：

- 不要引入来自 Pitaya 或 Nakama 的具体 public API、package、route、method、wire、handler、session、disconnect、kick、registry、selector 或 configuration compatibility 名称。
- 不要把 session binding、kick/disconnect 或 session data 词汇当作添加 handshake authentication changes、reconnect routing、kick routes、disconnect routes、general session data persistence、protocol messages、generated output、persistence、dependencies、topology 或 distributed runtime behavior 的许可。
- 不要把 authentication validation、request identity、first-message binding、logout、close policy、connection registry、presence lifecycle 或 delivery behavior 移到 transport、application、protocol、repository 或 startup 边界之外。

## 4. Ownership

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate.md
  - .arch/reference.yaml
  - .arch/runtime.yaml
source_first_map_candidate_owner:
  - tools/vibit
runtime_behavior_owner: unchanged
transport_owner: unchanged
protocol_owner: unchanged
persistence_owner: unchanged
module_owner: unchanged
dependency_owner: unchanged
```

文档和 manifest 可以定义 session binding、kick/disconnect 与 session data 词汇及当前映射。只有后续 bounded work item 明确授权时，`tools/vibit` 才可以暴露对应 source-first map。此闸门不改变 runtime、transport、protocol、repository、persistence、generated output、startup wiring、dependencies、service discovery、RPC、remote calls、frontend/backend role behavior、cluster-safe session routing、distributed group behavior、room broadcast behavior 或 session data persistence。

## 5. Verification

此边界的仓库检查规则是：

```text
runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate
```

验证命令：

```sh
node -c tools/vibit
node tools/vibit inspect rule runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate
node tools/vibit inspect next --json
node tools/vibit check change define-pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```
