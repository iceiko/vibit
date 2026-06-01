# Pitaya-Aligned Acceptor And Connection Lifecycle Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-01
Scope: Serializer and message forwarding 后续方向选择之后，使用 Pitaya-aligned acceptor and connection lifecycle vocabulary 的 gate-only boundary
Depends on: `decisions/ADR-0172-select-next-pitaya-aligned-direction-after-serializer-message-forwarding-map.md`, `decisions/ADR-0171-pitaya-aligned-serializer-message-forwarding-source-first-map.md`, `docs/pitaya-aligned-serializer-message-forwarding-boundary-gate.md`, `docs/runtime-protocol-adapter.md`, `docs/game-protocol.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0173`

配对英文原文是 `docs/pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.md`。英文文件是权威版本。

本文只定义 acceptor and connection lifecycle vocabulary gate。它不实现 acceptor behavior、TCP acceptors、WebSocket behavior changes、connection lifecycle behavior changes、session binding behavior、kick/disconnect behavior、serializer behavior、message forwarding behavior、route handler implementation、handler routing behavior、handler pipeline behavior、pipeline middleware behavior、backend route targeting、cluster-safe session routing behavior、distributed session routing、distributed runtime behavior、distributed groups、room broadcast fanout、delivery guarantees、stream subscriptions、service discovery implementation、service registries、service selectors、node identity、server-to-server RPC、remote calls、frontend/backend server roles、runtime endpoint behavior、metrics endpoints、observability pipelines、dashboards、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、authentication/session behavior changes、SDK publication、hosted deployments、release artifacts 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Pitaya-aligned acceptor and connection lifecycle boundary gate record 是：

```yaml
pitaya_aligned_acceptor_connection_lifecycle_boundary_gate: defined
completed_work_item: W-0265
decision: ADR-0173
check_rule: runtime.pitaya_aligned_acceptor_connection_lifecycle_boundary_gate
previous_direction_decision: ADR-0172
serializer_message_forwarding_source_first_map_decision: ADR-0171
serializer_message_forwarding_source_first_map_check_rule: runtime.pitaya_aligned_serializer_message_forwarding_source_first_map
standard: docs/pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.md
translation: docs/pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: acceptor_connection_lifecycle_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_acceptor_connection_lifecycle_vocabulary
future_implementation_work_item: W-0266
future_implementation_direction: implement_pitaya_aligned_acceptor_connection_lifecycle_source_first_map
allowed_acceptor_connection_lifecycle_vocabulary:
  - acceptor_boundary
  - websocket_acceptor
  - connection_id
  - connection_epoch
  - session_binding
  - active_connection_registry
  - close_handoff
  - presence_lifecycle_handoff
current_single_process_acceptor_connection_mapping:
  websocket_acceptor:
    current: single_process_websocket_server_accept_loop
    future_vocabulary: acceptor_boundary
    implementation_status: no_tcp_acceptor_or_distributed_acceptor
  connection_identity:
    current: server_observed_connection_id
    future_vocabulary: connection_id
    implementation_status: local_process_metadata
  connection_epoch:
    current: server_observed_connection_epoch
    future_vocabulary: connection_epoch
    implementation_status: no_distributed_routing_epoch
  first_message_binding:
    current: authentication_bind_connection_route
    future_vocabulary: session_binding
    implementation_status: no_handshake_authentication_or_reconnect_binding
  active_connection_registry:
    current: application_owned_connection_registry
    future_vocabulary: active_connection_registry
    implementation_status: no_connection_owner_node_registry
  close_handoff:
    current: transport_close_to_application_policy
    future_vocabulary: close_handoff
    implementation_status: no_remote_disconnect_handoff
  presence_lifecycle:
    current: server_owned_presence_snapshot
    future_vocabulary: presence_lifecycle_handoff
    implementation_status: no_distributed_presence_lifecycle
acceptor_behavior_added: false
tcp_acceptor_added: false
websocket_acceptor_behavior_changed: false
connection_lifecycle_behavior_changed: false
session_binding_behavior_added: false
kick_disconnect_behavior_added: false
concrete_socket_close_behavior_changed: false
serializer_behavior_added: false
message_forwarding_behavior_added: false
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
migration_added: false
dependency_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Purpose

`ADR-0172` 在 serializer and message forwarding source-first map 之后选择 acceptor and connection lifecycle vocabulary 作为下一项 Pitaya-aligned direction。

风险是 agent 可能把 acceptor、session binding、kick、disconnect 或 lifecycle 词汇当作添加 TCP acceptors、改变 WebSocket transport behavior、添加 handshake authentication、改变 connection close semantics、添加 remote disconnect behavior 或引入 distributed session routing 的许可。本 gate 只记录 vocabulary 和 mapping，保持 vibit 现有 single-process WebSocket accept loop、first-message binding、connection registry、close handoff 和 presence lifecycle behavior 不变。

## 3. Vocabulary

允许的 acceptor and connection lifecycle vocabulary：

- `acceptor_boundary`：未来规划 accept client connections 的 boundary。当前 owner 仍是现有 WebSocket transport。
- `websocket_acceptor`：未来规划第一种 accepted transport family。本文不授权改变 transport behavior。
- `connection_id`：未来规划 server-observed connection identity。当前 id 仍是 local metadata。
- `connection_epoch`：未来规划 connection generation metadata。当前 epoch 不是 distributed routing epoch。
- `session_binding`：未来规划 authenticated runtime session 与 connection 的绑定。当前绑定仍是既有 first-message route behavior。
- `active_connection_registry`：未来规划 active connections tracking。当前 state 仍是 application-owned single-process state。
- `close_handoff`：未来规划 transport 向 application lifecycle policy 交接 close facts。当前 close behavior 不变。
- `presence_lifecycle_handoff`：未来规划 connection lifecycle facts 与 presence lifecycle 的衔接。当前 presence 仍是 server-owned snapshot behavior。

Forbidden vocabulary use：

- 不要引入 Pitaya 或 Nakama 的 concrete public API、package、route、method、wire、handler、acceptor、session、disconnect、registry、selector 或 configuration compatibility names。
- 不要把 acceptor 或 connection lifecycle vocabulary 当作添加 TCP acceptors、WebSocket behavior changes、session binding behavior、kick/disconnect behavior、remote connection handoff、reconnect routing、protocol messages、generated output、persistence、dependencies、topology 或 distributed runtime behavior 的许可。
- 不要跨 transport、application、protocol、repository 或 startup boundaries 移动 session validation、first-message binding、close policy、presence lifecycle 或 delivery behavior。

## 4. Current Mapping

```yaml
current_single_process_acceptor_connection_mapping:
  websocket_acceptor:
    current: single-process WebSocket server accept loop
    future_vocabulary: acceptor_boundary
    status: no_tcp_acceptor_or_distributed_acceptor
  connection_identity:
    current: server-observed connection id
    future_vocabulary: connection_id
    status: local_process_metadata
  connection_epoch:
    current: server-observed connection epoch
    future_vocabulary: connection_epoch
    status: no_distributed_routing_epoch
  first_message_binding:
    current: authentication BindConnection route
    future_vocabulary: session_binding
    status: no_handshake_authentication_or_reconnect_binding
  active_connection_registry:
    current: application-owned connection registry
    future_vocabulary: active_connection_registry
    status: no_connection_owner_node_registry
  close_handoff:
    current: transport close to application policy handoff
    future_vocabulary: close_handoff
    status: no_remote_disconnect_handoff
  presence_lifecycle:
    current: server-owned presence snapshot
    future_vocabulary: presence_lifecycle_handoff
    status: no_distributed_presence_lifecycle
```

## 5. Ownership

Acceptor and connection lifecycle vocabulary ownership：

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.md
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

Rules：

- 文档和 manifests 可以定义 acceptor and connection lifecycle vocabulary 与 current mapping。
- 后续 implementation work item 授权时，`tools/vibit` 可以输出 source-first acceptor and connection lifecycle map。
- Runtime、transport、protocol、repository、persistence、generated output、startup wiring、dependencies、service discovery、RPC、remote calls、frontend/backend role behavior、cluster-safe session routing、distributed group behavior 和 room broadcast behavior 不因本 gate 改变。
- Domain modules 不会默认获得 acceptor、transport、session binding、disconnect、close policy、presence lifecycle、service discovery、RPC 或 distributed runtime ownership。

## 6. Nakama And Pitaya Mapping

Nakama 仍是近期 capability breadth 的 primary product reference。Pitaya 仍作为 acceptors、sessions、route handlers、frontend/backend roles、RPC/remotes、service discovery、groups、broadcast、cluster routing、handler pipelines、serializers、forwarding 和 connection lifecycle 的 architecture vocabulary reference。

Adopted as vocabulary：

- 用于未来 architecture planning 的 acceptor and connection lifecycle vocabulary；
- 对当前 single-process WebSocket accept loop、connection metadata、binding、registry、close handoff 和 presence lifecycle 的 mapping；
- 后续 source-first inspection work 的明确 deferral language。

Not adopted：

- direct Nakama 或 Pitaya API compatibility；
- concrete TCP acceptors；
- concrete session binding behavior changes；
- concrete kick/disconnect behavior；
- distributed connection owner registries；
- distributed session routing、reconnect routing 或 remote handoff behavior；
- metrics、tracing、dashboards、hosted surfaces、SDKs 或 release artifacts。

## 7. Verification

Required repository checks：

```bash
node -c tools/vibit
node tools/vibit inspect rule runtime.pitaya_aligned_acceptor_connection_lifecycle_boundary_gate
node tools/vibit check change define-pitaya-aligned-acceptor-connection-lifecycle-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

## 8. Stop Conditions

以下情况必须先停止并询问：

- 添加或改变 acceptor behavior；
- 添加 TCP acceptors；
- 改变 WebSocket behavior；
- 改变 connection lifecycle behavior；
- 改变 session binding behavior；
- 添加 kick/disconnect behavior；
- 改变 protocol messages or routes；
- 添加 generated output；
- 改变 repositories、PostgreSQL adapters、migrations 或 dependencies；
- 添加 metrics endpoints、tracing pipelines、dashboards、hosted deployment、SDK publication、release artifacts 或 direct Nakama/Pitaya API compatibility。
