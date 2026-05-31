# Pitaya-Aligned Server To Server RPC Boundary Gate 中文版

状态：Accepted v0.1
最后更新：2026-05-31
范围：frontend/backend role source-first map 之后，使用 Pitaya-aligned server-to-server RPC 和 remote-call vocabulary 的 gate-only boundary
依赖：`decisions/ADR-0157-pitaya-aligned-frontend-backend-role-source-first-map.md`、`docs/pitaya-aligned-frontend-backend-role-boundary-gate.md`、`docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md`、`docs/reference-game-server-alignment.md`、`.arch/reference.yaml`
权威决策：`ADR-0158`

英文文件 `docs/pitaya-aligned-server-to-server-rpc-boundary-gate.md` 是权威版本。本文是简体中文译本。

本文只定义 server-to-server RPC vocabulary gate。它不实现 server-to-server RPC、remote calls、service discovery、frontend/backend server roles、distributed runtime behavior、distributed groups、room broadcast fanout、cluster-safe session routing、runtime endpoint behavior、metrics endpoints、observability pipelines、dashboards、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、authentication/session behavior changes、SDK publication、hosted deployments、release artifacts 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Pitaya-aligned server-to-server RPC boundary gate 记录是：

```yaml
pitaya_aligned_server_to_server_rpc_boundary_gate: defined
completed_work_item: W-0250
decision: ADR-0158
check_rule: runtime.pitaya_aligned_server_to_server_rpc_boundary_gate
role_source_first_map_decision: ADR-0157
role_source_first_map_check_rule: runtime.pitaya_aligned_frontend_backend_role_source_first_map
standard: docs/pitaya-aligned-server-to-server-rpc-boundary-gate.md
translation: docs/pitaya-aligned-server-to-server-rpc-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: server_to_server_rpc_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_server_to_server_rpc_vocabulary
future_implementation_work_item: W-0251
future_implementation_direction: implement_pitaya_aligned_server_to_server_rpc_source_first_map
allowed_rpc_vocabulary:
  - server_to_server_rpc
  - remote_call
related_vocabulary:
  - route_handler
  - module_handler
  - application_dispatch
  - service_discovery
current_single_process_rpc_mapping:
  server_to_server_rpc:
    current: no_rpc_current_single_process_application_dispatch
    future_vocabulary: server_to_server_rpc
    implementation_status: deferred_future_architecture_reference
  remote_call:
    current: no_remote_call_current_in_process_module_invocation
    future_vocabulary: remote_call
    implementation_status: deferred_future_architecture_reference
  route_handler:
    current: current_application_dispatch_and_protocol_bridge
    future_vocabulary: backend_server_route_handler_boundary
    implementation_status: current_single_process_only
  module_handler:
    current: current_module_handler_in_process_function_call
    future_vocabulary: backend_server_module_handler_boundary
    implementation_status: current_single_process_only
  service_discovery:
    current: no_service_discovery_current_static_single_process_composition
    future_vocabulary: service_discovery
    implementation_status: deferred_future_architecture_reference
server_to_server_rpc_implementation_added: false
remote_call_behavior_added: false
service_discovery_added: false
frontend_server_implementation_added: false
backend_server_implementation_added: false
frontend_backend_server_roles_added: false
distributed_runtime_implementation_added: false
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

`ADR-0157` 已经通过 `node tools/vibit inspect pitaya-roles --json` 让 frontend/backend role vocabulary 可检查。下一类最容易导致误实现的 Pitaya-aligned 概念是 server-to-server RPC。Pitaya 把 distributed server calls 作为架构的一部分；vibit 只能在边界明确后把这些词汇用于规划。

本 gate 在任何实现之前记录 vocabulary 和 current mapping。它不添加 RPC transports、remote handlers、service registries、node identity、frontend/backend process topology、protocol carriers 或 client-visible routes。

## 3. RPC Vocabulary

允许的 RPC vocabulary：

- `server_to_server_rpc`：用于 future architecture planning 的 server-internal call family。若后续实现，它不得绕过 vibit module contracts、route contracts、authorization boundaries、identity/session validation、generated boundaries 或 repository checks。
- `remote_call`：未来 distributed invocation vocabulary item。它不同于 client protocol commands、queries、events 和 WebSocket routes。

相关 vocabulary：

- `route_handler`：当前是 application dispatch 加 protocol bridge；未来规划可把它关联到 backend handler ownership。
- `module_handler`：当前是 in-process handwritten module behavior；未来规划可把它关联到 internal service handling，但 module ownership 仍是权威来源。
- `application_dispatch`：当前 vibit in-process command/query dispatch path。
- `service_discovery`：未来 dependency-sensitive architecture vocabulary item，不是当前实现。

禁止的 vocabulary 使用：

- 不要从 Pitaya 引入 concrete public API、package、route、method 或 wire compatibility names。
- 不要把 `server_to_server_rpc` 或 `remote_call` 当作添加 RPC transports、remote calls、service discovery、node registry、distributed process topology、new endpoint behavior、protocol changes、generated output、persistence 或 dependencies 的许可。
- 不要用 RPC vocabulary 绕过 module contracts、application dispatch boundaries、authentication/session validation gates、permission checks、generated output rules 或 repository ownership。

## 4. Current Mapping

```yaml
current_single_process_rpc_mapping:
  server_to_server_rpc:
    current: no server-to-server RPC; application dispatch is in-process
    future_vocabulary: server_to_server_rpc
    status: deferred_future_architecture_reference
  remote_call:
    current: no remote call; module handlers are invoked in-process
    future_vocabulary: remote_call
    status: deferred_future_architecture_reference
  route_handler:
    current: Protobuf bridge plus application dispatch and module handlers
    future_vocabulary: backend_server route handler boundary
    status: current_application_dispatch_and_protocol_bridge
  module_handler:
    current: runtime/internal/modules handwritten behavior in one process
    future_vocabulary: backend_server module handler boundary
    status: current_module_handler_in_process_function_call
  service_discovery:
    current: none; composition is static and single-process
    future_vocabulary: service_discovery
    status: deferred_future_architecture_reference
```

## 5. Ownership

RPC vocabulary ownership：

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-server-to-server-rpc-boundary-gate.md
  - .arch/reference.yaml
  - .arch/runtime.yaml
source_first_map_candidate_owner:
  - tools/vibit
runtime_behavior_owner: unchanged
protocol_owner: unchanged
persistence_owner: unchanged
module_owner: unchanged
dependency_owner: unchanged
```

规则：

- Documentation 和 manifests 可以定义 RPC vocabulary 与 mapping。
- 若后续 implementation work item 授权，`tools/vibit` 可以输出 source-first RPC map。
- Runtime、protocol、repository、persistence、generated output、startup wiring、dependencies 和 service discovery 不因本 gate 改变。
- Domain modules 默认不会获得 RPC ownership。Module contracts 仍是 module behavior 和 data ownership 的来源。

## 6. Nakama And Pitaya Mapping

Nakama 仍是近期 capability breadth 的主要产品参考。Pitaya 仍是 frontend/backend roles、route handler placement、RPC/remotes、service discovery、groups、broadcast 和 cluster routing 的 architecture vocabulary reference。

采纳为 vocabulary：

- server-to-server RPC 作为 future architecture-planning vocabulary；
- remote calls 作为 future distributed invocation vocabulary；
- service discovery 作为 future dependency-sensitive vocabulary；
- 把 current in-process dispatch 和 module handlers 映射到 deferred RPC concepts。

适配到 vibit：

- 任何未来 RPC 都必须保留 vibit module ownership、application dispatch boundaries、server-authoritative validation、generated output rules、redaction 和 repository checks。
- 当前 single-process runtime 仍是具体实现。
- 任何未来 RPC implementation 都必须先经过单独 gate 和 verification。

当前拒绝：

- direct Pitaya API compatibility；
- Pitaya package、method 或 route naming compatibility；
- server-to-server RPC implementation；
- remote call behavior；
- service discovery；
- distributed runtime behavior；
- frontend/backend process split；
- RPC 的 protocol messages 或 routes。

## 7. Future Implementation Work

打开：

```text
M-179/W-0251 Implement Pitaya-aligned server-to-server RPC source-first map
```

后续 work item 可以：

- 添加 server-to-server RPC 和 remote-call vocabulary 的 source-first repository inspection map；
- 摘要 current single-process dispatch 和 module handler mapping；
- 更新 runbooks 和 acceptance docs 指向 RPC map；
- 添加 repository checks，确认 RPC map 仍保持 gate-only 且 redacted。

后续 work item 禁止：

- 添加 server-to-server RPC implementation；
- 添加 remote call behavior；
- 添加 service discovery；
- 添加 frontend/backend server role implementation；
- 添加 distributed runtime implementation；
- 添加 distributed groups、room broadcast fanout 或 delivery guarantees；
- 添加 cluster-safe session routing behavior；
- 添加 runtime endpoint behavior、metrics endpoints、observability pipelines、dashboards、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、authentication/session behavior changes、SDK publication、hosted deployments、release artifacts 或 direct Nakama/Pitaya API compatibility。

## 8. Verification Expectations

本 gate 应验证：

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_server_to_server_rpc_boundary_gate
node tools/vibit check change define-pitaya-aligned-server-to-server-rpc-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

本 gate 本身不需要 Go tests，因为它没有添加 Go runtime behavior。

## 9. Stop Conditions

如果 work 需要以下内容，停止并创建单独 gate：

- RPC transport behavior；
- remote call behavior；
- service discovery 或 node registry behavior；
- process topology changes；
- new goroutines、listeners、network protocols 或 server roles；
- distributed group 或 broadcast behavior；
- cluster-safe session routing behavior；
- protocol 或 Protobuf changes；
- repository、adapter、migration、dependency 或 generated-output changes；
- 与 Pitaya 或 Nakama 的 public API compatibility；
- 任何 sensitive runtime state exposure。
