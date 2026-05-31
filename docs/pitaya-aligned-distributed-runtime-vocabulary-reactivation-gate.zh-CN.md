# Pitaya-Aligned Distributed Runtime Vocabulary Reactivation Gate 中文版

状态：Accepted v0.1
最后更新：2026-05-31
范围：source-first operations inspection 之后，重新激活 Pitaya-aligned distributed runtime vocabulary 的 gate-only boundary
依赖：`decisions/ADR-0153-minimum-operations-inspection-source-first-surface-implementation.md`、`docs/reference-game-server-alignment.md`、`docs/nakama-pitaya-product-parity-roadmap.md`、`.arch/reference.yaml`
权威决策：`ADR-0154`

英文文件 `docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md` 是权威版本。本文是简体中文译本。

本文只定义 vocabulary gate。它不实现 distributed runtime behavior、frontend/backend server roles、server-to-server RPC、remote calls、service discovery、distributed groups、broadcast fanout、cluster-safe session routing、runtime endpoint behavior、metrics endpoints、observability pipelines、dashboards、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、authentication/session behavior changes、SDK publication、hosted deployments、release artifacts 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Pitaya-aligned distributed runtime vocabulary reactivation gate 记录是：

```yaml
pitaya_aligned_distributed_runtime_vocabulary_reactivation_gate: defined
completed_work_item: W-0246
decision: ADR-0154
check_rule: runtime.pitaya_aligned_distributed_runtime_vocabulary_reactivation_gate
source_operations_inspection_decision: ADR-0153
source_operations_inspection_check_rule: runtime.minimum_operations_inspection_source_first_surface_implementation
standard: docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md
translation: docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: vocabulary_reactivated_for_future_architecture_planning
implementation_scope: gate_only_architecture_vocabulary
future_implementation_work_item: W-0247
future_implementation_direction: implement_pitaya_aligned_distributed_runtime_vocabulary_source_first_map
allowed_vocabulary:
  - acceptor
  - frontend_server
  - backend_server
  - route_handler
  - session_binding
  - server_to_server_rpc
  - remote_call
  - service_discovery
  - distributed_group
  - room_broadcast
  - cluster_safe_session_routing
current_single_process_mapping:
  websocket_tcp_acceptors: current_websocket_acceptor_single_process_only
  session_binding: current_first_message_connection_binding_single_process_only
  route_handler_model: current_application_dispatch_and_protocol_bridge
  frontend_backend_server_roles: deferred_future_architecture_reference
  rpc_and_remote_calls: deferred_future_architecture_reference
  groups_rooms_broadcast: deferred_future_architecture_reference
  cluster_service_discovery: deferred_future_architecture_reference
distributed_runtime_implementation_added: false
frontend_backend_server_roles_added: false
server_to_server_rpc_added: false
service_discovery_added: false
distributed_groups_added: false
cluster_safe_session_routing_added: false
runtime_behavior_added: false
runtime_endpoint_behavior_added: false
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

`ADR-0153` 已经在 source-first operations inspection output 中显式记录 Pitaya architecture pressure。该信息有用，但仍嵌在 operations surface 里。下一步安全推进是定义独立 vocabulary gate，让未来 agents 可以讨论 distributed runtime concepts，同时不意外实现它们。

本 gate 只把 Pitaya 重新激活为 future architecture planning vocabulary。它不会让 Pitaya 成为当前 product capability driver。Nakama 仍是近期 capability breadth 的主要产品参考。只有讨论 transport acceptors、session binding、route handling、frontend/backend roles、RPC/remotes、service discovery、distributed groups、broadcast 和 cluster-safe routing 时，才允许使用 Pitaya vocabulary。

## 3. Vocabulary

允许的 vocabulary：

- `acceptor`：未来 client connection acceptors 的抽象，例如 WebSocket 或 TCP。当前 vibit 仍只有 single-process WebSocket acceptor。
- `frontend_server`：未来拥有 client-facing acceptors 和 session ingress 的角色。当前 vibit 没有 frontend/backend role split。
- `backend_server`：未来在 frontend server 之后拥有 backend services 或 domain handlers 的角色。当前 vibit 是 modular monolith。
- `route_handler`：当前 application dispatch 加 protocol bridge 最接近 Pitaya handler routing，但 vibit 仍保留 route contracts 和 module ownership。
- `session_binding`：当前 first-message connection binding 是未来 cluster-safe binding semantics 的 single-process 前身。
- `server_to_server_rpc`：未来 server-to-server call family。后续引入时不得绕过 module contracts。
- `remote_call`：未来 distributed call vocabulary item，与 public client protocol routes 分离。
- `service_discovery`：未来 cluster membership 和 service lookup concern。
- `distributed_group`：未来 cluster-aware group、room、party、stream 或 match broadcast target concern。
- `room_broadcast`：未来 broadcast vocabulary item。它不授权 delivery guarantees、fanout code 或 protocol messages。
- `cluster_safe_session_routing`：未来跨节点保持 player/session/connection ownership 一致的 routing semantics。

禁止的 vocabulary 使用：

- 不要按 Pitaya API 命名 concrete public APIs。
- 不要添加 Pitaya package、namespace、route、method 或 wire compatibility markers。
- 不要把 vocabulary reactivation 当作 implementation code 的授权。

## 4. Current Mapping

```yaml
current_single_process_mapping:
  websocket_tcp_acceptors:
    current: runtime/cmd/vibit-server WebSocket acceptor
    future_vocabulary: acceptor
    implementation_status: current_websocket_acceptor_single_process_only
  session_binding:
    current: first-message connection binding
    future_vocabulary: session_binding and cluster_safe_session_routing
    implementation_status: current_first_message_connection_binding_single_process_only
  route_handler_model:
    current: Protobuf bridge plus application dispatch and module handlers
    future_vocabulary: route_handler
    implementation_status: current_application_dispatch_and_protocol_bridge
  frontend_backend_server_roles:
    current: none
    future_vocabulary: frontend_server and backend_server
    implementation_status: deferred_future_architecture_reference
  rpc_and_remote_calls:
    current: none
    future_vocabulary: server_to_server_rpc and remote_call
    implementation_status: deferred_future_architecture_reference
  groups_rooms_broadcast:
    current: target scope vocabulary exists; distributed groups do not
    future_vocabulary: distributed_group and room_broadcast
    implementation_status: deferred_future_architecture_reference
  cluster_service_discovery:
    current: none
    future_vocabulary: service_discovery
    implementation_status: deferred_future_architecture_reference
```

## 5. Ownership

Vocabulary ownership：

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md
  - .arch/reference.yaml
  - .arch/runtime.yaml
source_first_map_candidate_owner:
  - tools/vibit
runtime_behavior_owner: unchanged
protocol_owner: unchanged
persistence_owner: unchanged
module_owner: unchanged
```

规则：

- Documentation 和 manifests 可以定义 vocabulary 与 mapping。
- 若后续 implementation work item 授权，`tools/vibit` 可以输出 source-first vocabulary map。
- Runtime、protocol、repository、persistence、generated output、startup wiring 和 dependencies 不因本 gate 改变。
- 默认情况下，没有任何 game module 拥有 distributed runtime vocabulary。

## 6. Nakama And Pitaya Mapping

Nakama 仍是主要产品参考：

- 当前 product path 继续优先 common backend capability coverage、local alpha clarity、route proof 和 prototype usefulness。
- 本 gate 不选择新的 social、realtime、matchmaking、match runtime、SDK 或 operations implementation slice。

Pitaya 只作为 future architecture vocabulary 重新激活：

- 采纳为 vocabulary：acceptors、session binding、route handlers、frontend/backend roles、server-to-server RPC、remotes、service discovery、groups、broadcast 和 cluster routing。
- 适配到 vibit：vocabulary 必须留在 contract-first、module-owned、repository-checkable boundaries 后面。
- 当前拒绝：cluster runtime implementation、direct Pitaya API compatibility、public route naming compatibility、package namespace compatibility，以及任何绕过 vibit module contracts 的行为。

## 7. Future Implementation Work

打开：

```text
M-175/W-0247 Implement Pitaya-aligned distributed runtime vocabulary source-first map
```

后续 work item 可以：

- 添加 source-first `tools/vibit inspect pitaya-vocabulary` 命令或等价 repository inspection；
- 摘要本 gate 定义的 vocabulary 和 current single-process mapping；
- 更新 runbooks 和 acceptance docs 指向 vocabulary map；
- 添加 repository checks，确认 map 仍保持 gate-only 且 redacted。

后续 work item 禁止：

- 添加 distributed runtime implementation；
- 添加 frontend/backend server role implementation；
- 添加 server-to-server RPC 或 remote call behavior；
- 添加 service discovery；
- 添加 distributed groups、broadcast fanout 或 delivery guarantees；
- 添加 cluster-safe session routing behavior；
- 添加 runtime endpoint behavior、metrics endpoints、observability pipelines、dashboards、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、authentication/session behavior changes、SDK publication、hosted deployments、release artifacts 或 direct Nakama/Pitaya API compatibility。

## 8. Verification Expectations

本 gate 应验证：

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_distributed_runtime_vocabulary_reactivation_gate
node tools/vibit check change define-pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

本 gate 本身不要求 Go tests，因为它不添加 Go runtime behavior。

## 9. Stop Conditions

如果 work 需要以下内容，应停止并创建 separate gate：

- process topology changes；
- new goroutines、listeners、network protocols 或 server roles；
- RPC/remoting behavior；
- service discovery dependencies；
- distributed group 或 broadcast behavior；
- cluster-safe session routing behavior；
- protocol 或 Protobuf changes；
- repository、adapter、migration 或 generated-output changes；
- public API compatibility with Pitaya or Nakama；
- 任何 sensitive runtime state exposure。
