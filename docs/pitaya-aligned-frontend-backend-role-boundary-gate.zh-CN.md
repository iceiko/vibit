# Pitaya-Aligned Frontend Backend Role Boundary Gate 中文版

状态：Accepted v0.1
最后更新：2026-05-31
范围：source-first distributed runtime vocabulary map 之后，使用 Pitaya-aligned frontend/backend role vocabulary 的 gate-only boundary
依赖：`decisions/ADR-0155-pitaya-aligned-distributed-runtime-vocabulary-source-first-map.md`、`docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md`、`docs/reference-game-server-alignment.md`、`.arch/reference.yaml`
权威决策：`ADR-0156`

英文文件 `docs/pitaya-aligned-frontend-backend-role-boundary-gate.md` 是权威版本。本文是简体中文译本。

本文只定义 frontend/backend role vocabulary gate。它不实现 frontend/backend server roles、distributed runtime behavior、server-to-server RPC、remote calls、service discovery、distributed groups、room broadcast fanout、cluster-safe session routing、runtime endpoint behavior、metrics endpoints、observability pipelines、dashboards、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、authentication/session behavior changes、SDK publication、hosted deployments、release artifacts 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Pitaya-aligned frontend/backend role boundary gate 记录是：

```yaml
pitaya_aligned_frontend_backend_role_boundary_gate: defined
completed_work_item: W-0248
decision: ADR-0156
check_rule: runtime.pitaya_aligned_frontend_backend_role_boundary_gate
source_vocabulary_map_decision: ADR-0155
source_vocabulary_map_check_rule: runtime.pitaya_aligned_distributed_runtime_vocabulary_source_first_map
standard: docs/pitaya-aligned-frontend-backend-role-boundary-gate.md
translation: docs/pitaya-aligned-frontend-backend-role-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: frontend_backend_role_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_frontend_backend_role_vocabulary
future_implementation_work_item: W-0249
future_implementation_direction: implement_pitaya_aligned_frontend_backend_role_source_first_map
allowed_role_vocabulary:
  - frontend_server
  - backend_server
related_vocabulary:
  - acceptor
  - session_binding
  - route_handler
current_single_process_role_mapping:
  frontend_server:
    current: single_process_acceptor_and_dispatch
    future_vocabulary: frontend_server
    implementation_status: deferred_future_architecture_reference
  backend_server:
    current: application_dispatch_and_module_handlers_in_same_process
    future_vocabulary: backend_server
    implementation_status: deferred_future_architecture_reference
  acceptor:
    current: current_websocket_acceptor_single_process_only
    future_vocabulary: frontend_server_acceptor_boundary
    implementation_status: current_single_process_only
  session_binding:
    current: current_first_message_connection_binding_single_process_only
    future_vocabulary: frontend_server_session_ingress_boundary
    implementation_status: current_single_process_only
  route_handler:
    current: current_application_dispatch_and_protocol_bridge
    future_vocabulary: backend_server_handler_boundary
    implementation_status: current_single_process_only
frontend_server_implementation_added: false
backend_server_implementation_added: false
frontend_backend_server_roles_added: false
distributed_runtime_implementation_added: false
server_to_server_rpc_added: false
remote_call_behavior_added: false
service_discovery_added: false
distributed_groups_added: false
room_broadcast_fanout_added: false
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

`ADR-0155` 已经通过 `node tools/vibit inspect pitaya-vocabulary --json` 让更宽的 Pitaya-aligned distributed runtime vocabulary 可检查。下一步安全推进是为最容易导致误实现的两个 role 词汇定义更窄边界：`frontend_server` 和 `backend_server`。

本 gate 只允许 role vocabulary 用于 future architecture planning。它不拆分进程、不添加 topology、不添加 listeners、不添加 handler remoting、不添加 service discovery，也不改变 protocol routes。

## 3. Role Vocabulary

允许的 role vocabulary：

- `frontend_server`：未来可能拥有 acceptors、session ingress、connection lifecycle 和 routing handoff 的 client-facing role。当前 vibit 仍是 single-process WebSocket acceptor 和 application dispatch path。
- `backend_server`：未来可能在 frontend role 之后拥有 module handlers 的 service-facing role。当前 vibit 仍在同一进程里执行 application dispatch 和 module handlers。

相关 vocabulary：

- `acceptor`：当前是 single-process WebSocket acceptor；未来规划可把它关联到 frontend role。
- `session_binding`：当前是 first-message connection binding；未来规划可把它关联到 frontend session ingress。
- `route_handler`：当前是 application dispatch 加 protocol bridge；未来规划可把它关联到 backend handler ownership。

禁止的 vocabulary 使用：

- 不要从 Pitaya 引入 concrete public API、package、route、method 或 wire compatibility names。
- 不要把 `frontend_server` 或 `backend_server` 当作添加 process topology、runtime behavior、new listeners、RPC/remoting、service discovery、protocol changes、generated output、persistence 或 dependencies 的许可。
- 不要把 module ownership 转移到 role vocabulary。Module contracts 和 vibit ownership manifests 仍是权威来源。

## 4. Current Mapping

```yaml
current_single_process_role_mapping:
  frontend_server:
    current: single_process_acceptor_and_dispatch
    future_vocabulary: frontend_server
    status: deferred_future_architecture_reference
  backend_server:
    current: application_dispatch_and_module_handlers_in_same_process
    future_vocabulary: backend_server
    status: deferred_future_architecture_reference
  acceptor:
    current: runtime/cmd/vibit-server WebSocket acceptor
    future_role: frontend_server
    status: current_websocket_acceptor_single_process_only
  session_binding:
    current: first-message connection binding
    future_role: frontend_server
    status: current_first_message_connection_binding_single_process_only
  route_handler:
    current: Protobuf bridge plus application dispatch and module handlers
    future_role: backend_server
    status: current_application_dispatch_and_protocol_bridge
```

## 5. Ownership

Role vocabulary ownership：

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-frontend-backend-role-boundary-gate.md
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

- Documentation 和 manifests 可以定义 role vocabulary 与 mapping。
- 若后续 implementation work item 授权，`tools/vibit` 可以输出 source-first frontend/backend role map。
- Runtime、protocol、repository、persistence、generated output、startup wiring 和 dependencies 不因本 gate 改变。
- 默认情况下，没有任何 domain module 拥有 frontend/backend role vocabulary。

## 6. Nakama And Pitaya Mapping

Nakama 仍是近期 capability breadth 的主要产品参考。Pitaya 仍是 role topology、acceptor/session ingress、route-handler placement、RPC/remotes、service discovery 和 cluster routing 的 architecture vocabulary reference。

采纳为 vocabulary：

- frontend role 作为未来 client-facing ingress vocabulary；
- backend role 作为未来 service/handler ownership vocabulary；
- frontend/backend mapping 到当前 acceptor、session binding、protocol bridge、application dispatch 和 module handler responsibilities。

适配到 vibit：

- Role vocabulary 必须留在 contract-first、module-owned、repository-checkable boundaries 后面。
- 当前 single-process runtime 仍是具体实现。
- 任何未来 role split 都必须保留 vibit route contracts、module ownership、generated boundaries 和 verification commands。

当前拒绝：

- direct Pitaya API compatibility；
- Pitaya package 或 route naming compatibility；
- runtime topology changes；
- frontend/backend process split；
- handler remoting、service discovery、distributed groups、broadcast fanout 或 cluster-safe routing behavior。

## 7. Future Implementation Work

打开：

```text
M-177/W-0249 Implement Pitaya-aligned frontend/backend role source-first map
```

后续 work item 可以：

- 添加 frontend/backend role vocabulary 的 source-first repository inspection map；
- 摘要 current single-process role mapping；
- 更新 runbooks 和 acceptance docs 指向 role map；
- 添加 repository checks，确认 role map 仍保持 gate-only 且 redacted。

后续 work item 禁止：

- 添加 frontend/backend server role implementation；
- 添加 distributed runtime implementation；
- 添加 server-to-server RPC 或 remote call behavior；
- 添加 service discovery；
- 添加 distributed groups、room broadcast fanout 或 delivery guarantees；
- 添加 cluster-safe session routing behavior；
- 添加 runtime endpoint behavior、metrics endpoints、observability pipelines、dashboards、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、authentication/session behavior changes、SDK publication、hosted deployments、release artifacts 或 direct Nakama/Pitaya API compatibility。

## 8. Verification Expectations

本 gate 应验证：

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_frontend_backend_role_boundary_gate
node tools/vibit check change define-pitaya-aligned-frontend-backend-role-boundary-gate --json
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
