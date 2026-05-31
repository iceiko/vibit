# Pitaya-Aligned Service Discovery Boundary Gate 中文版

状态：Accepted v0.1
最后更新：2026-05-31
范围：server-to-server RPC source-first map 之后，使用 Pitaya-aligned service discovery vocabulary 的 gate-only boundary
依赖：`decisions/ADR-0159-pitaya-aligned-server-to-server-rpc-source-first-map.md`、`docs/pitaya-aligned-server-to-server-rpc-boundary-gate.md`、`docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md`、`docs/reference-game-server-alignment.md`、`.arch/reference.yaml`
权威决策：`ADR-0160`

英文文件 `docs/pitaya-aligned-service-discovery-boundary-gate.md` 是权威版本。本文是简体中文译本。

本文只定义 service discovery vocabulary gate。它不实现 service discovery、service registries、service selectors、node registries、server identity、server-to-server RPC、remote calls、frontend/backend server roles、distributed runtime behavior、distributed groups、room broadcast fanout、cluster-safe session routing、runtime endpoint behavior、metrics endpoints、observability pipelines、dashboards、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、authentication/session behavior changes、SDK publication、hosted deployments、release artifacts 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Pitaya-aligned service discovery boundary gate 记录是：

```yaml
pitaya_aligned_service_discovery_boundary_gate: defined
completed_work_item: W-0252
decision: ADR-0160
check_rule: runtime.pitaya_aligned_service_discovery_boundary_gate
rpc_source_first_map_decision: ADR-0159
rpc_source_first_map_check_rule: runtime.pitaya_aligned_server_to_server_rpc_source_first_map
standard: docs/pitaya-aligned-service-discovery-boundary-gate.md
translation: docs/pitaya-aligned-service-discovery-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: service_discovery_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_service_discovery_vocabulary
future_implementation_work_item: W-0253
future_implementation_direction: implement_pitaya_aligned_service_discovery_source_first_map
allowed_service_discovery_vocabulary:
  - service_discovery
  - service_registry
  - service_instance
  - service_selector
related_vocabulary:
  - frontend_server
  - backend_server
  - server_to_server_rpc
  - remote_call
  - route_handler
  - module_handler
  - static_process_composition
current_single_process_service_discovery_mapping:
  service_discovery:
    current: no_service_discovery_current_static_single_process_composition
    future_vocabulary: service_discovery
    implementation_status: deferred_future_architecture_reference
  service_registry:
    current: no_registry_current_startup_composition
    future_vocabulary: service_registry
    implementation_status: deferred_future_architecture_reference
  service_instance:
    current: single_process_runtime_components_not_network_instances
    future_vocabulary: service_instance
    implementation_status: deferred_future_architecture_reference
  service_selector:
    current: no_selector_current_direct_in_process_dispatch
    future_vocabulary: service_selector
    implementation_status: deferred_future_architecture_reference
  route_handler:
    current: current_application_dispatch_and_protocol_bridge
    future_vocabulary: discoverable_backend_route_handler
    implementation_status: current_single_process_only
  module_handler:
    current: current_module_handler_in_process_function_call
    future_vocabulary: discoverable_backend_module_handler
    implementation_status: current_single_process_only
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

`ADR-0159` 已经通过 `node tools/vibit inspect pitaya-rpc --json` 让 server-to-server RPC 和 remote-call vocabulary 可检查。下一类容易被误解成实现许可的 Pitaya-aligned 概念是 service discovery。Distributed runtimes 需要找到 services 或 server instances，但 vibit 现在使用一个静态组合的 single-process。

本 gate 在任何实现之前记录 service discovery vocabulary 和 current mapping。它不添加 registry storage、registry clients、selection algorithms、node identity、heartbeat behavior、routing tables、process topology、network listeners、RPC transports 或 dependencies。

## 3. Service Discovery Vocabulary

允许的 service discovery vocabulary：

- `service_discovery`：用于 future architecture planning 的服务端能力或实例定位 family。
- `service_registry`：未来可用服务端能力或实例的记录。它在本 slice 中不是 database table、external registry 或 in-memory runtime structure。
- `service_instance`：未来 discoverable server-side instance 的表示。它不是当前 runtime process identity。
- `service_selector`：未来选择 service instance 的概念。它在本 slice 中不是 load balancer、routing algorithm 或 retry policy。

相关 vocabulary：

- `frontend_server` 和 `backend_server`：来自 prior Pitaya-aligned role boundary 的 role vocabulary。它们仍只是 future architecture vocabulary。
- `server_to_server_rpc` 和 `remote_call`：来自 prior boundary 的 RPC vocabulary。Service discovery 不授权 RPC implementation。
- `route_handler` 和 `module_handler`：当前 in-process application/module handlers，未来可映射到 discoverable backend ownership。
- `static_process_composition`：当前 vibit 的具体 implementation model。

禁止的 vocabulary 使用：

- 不要从 Pitaya 引入 concrete public API、package、route、method、wire、registry 或 configuration compatibility names。
- 不要把 service discovery vocabulary 当作添加 service registries、selectors、heartbeats、node identity、runtime topology、registry storage、external discovery dependencies、RPC transports、remote calls、endpoint behavior、protocol changes、generated output、persistence 或 dependencies 的许可。
- 不要用 future service discovery vocabulary 绕过 module contracts、application dispatch boundaries、authentication/session validation gates、permission checks、generated output rules、redaction rules 或 repository ownership。

## 4. Current Mapping

```yaml
current_single_process_service_discovery_mapping:
  service_discovery:
    current: none; composition is static and single-process
    future_vocabulary: service_discovery
    status: deferred_future_architecture_reference
  service_registry:
    current: no registry; startup wires concrete handlers directly
    future_vocabulary: service_registry
    status: deferred_future_architecture_reference
  service_instance:
    current: runtime components are not discoverable network instances
    future_vocabulary: service_instance
    status: deferred_future_architecture_reference
  service_selector:
    current: no selector; dispatch is direct and in-process
    future_vocabulary: service_selector
    status: deferred_future_architecture_reference
  route_handler:
    current: Protobuf bridge plus application dispatch and module handlers
    future_vocabulary: discoverable backend route handler
    status: current_application_dispatch_and_protocol_bridge
  module_handler:
    current: runtime/internal/modules handwritten behavior in one process
    future_vocabulary: discoverable backend module handler
    status: current_module_handler_in_process_function_call
```

## 5. Ownership

Service discovery vocabulary ownership：

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-service-discovery-boundary-gate.md
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

- Documentation 和 manifests 可以定义 service discovery vocabulary 与 mapping。
- 若后续 implementation work item 授权，`tools/vibit` 可以输出 source-first service discovery map。
- Runtime、protocol、repository、persistence、generated output、startup wiring、dependencies、RPC、remote calls 和 frontend/backend role behavior 不因本 gate 改变。
- Domain modules 默认不会获得 service discovery ownership。Module contracts 仍是 module behavior 和 data ownership 的来源。

## 6. Nakama And Pitaya Mapping

Nakama 仍是近期 capability breadth 的主要产品参考。Pitaya 仍是 frontend/backend roles、route handler placement、RPC/remotes、service discovery、groups、broadcast 和 cluster routing 的 architecture vocabulary reference。

采纳为 vocabulary：

- service discovery 作为 future architecture-planning vocabulary；
- service registry、service instance 和 service selector 作为 future discovery vocabulary；
- 把 current static startup composition 映射到 deferred service discovery concepts；
- 把 current route and module handlers 映射到可能的 future discoverable backend ownership。

适配到 vibit：

- 任何 future discovery model 都必须保留 vibit module ownership、application dispatch boundaries、server-authoritative validation、generated output rules、redaction 和 repository checks。
- 当前 single-process runtime 仍是具体实现。
- 任何 future service discovery implementation 都必须先经过单独 gate 和 verification。

当前拒绝：

- direct Pitaya API compatibility；
- Pitaya package、method、registry 或 route naming compatibility；
- service discovery implementation；
- service registry 或 selector behavior；
- node identity、heartbeats 或 runtime topology；
- server-to-server RPC 或 remote call behavior；
- frontend/backend process split；
- discovery 的 protocol messages 或 routes。

## 7. Future Implementation Work

打开：

```text
M-181/W-0253 Implement Pitaya-aligned service discovery source-first map
```

后续 work item 可以：

- 添加 service discovery vocabulary 的 source-first repository inspection map；
- 摘要 current static single-process composition 和 direct handler dispatch；
- 更新 runbooks 和 acceptance docs 指向 service discovery map；
- 添加 repository checks，确认 service discovery map 仍保持 gate-only 且 redacted。

后续 work item 禁止：

- 添加 service discovery implementation；
- 添加 service registry 或 selector behavior；
- 添加 node identity、heartbeat、membership 或 topology behavior；
- 添加 server-to-server RPC implementation；
- 添加 remote call behavior；
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
node tools/vibit inspect rule runtime.pitaya_aligned_service_discovery_boundary_gate
node tools/vibit check change define-pitaya-aligned-service-discovery-boundary-gate --json
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

- service discovery implementation；
- service registry、selector、membership、heartbeat 或 node identity behavior；
- server-to-server RPC 或 remote call behavior；
- process topology changes；
- new goroutines、listeners、network protocols 或 server roles；
- distributed group 或 broadcast behavior；
- cluster-safe session routing behavior；
- protocol 或 Protobuf changes；
- repository、adapter、migration、dependency 或 generated-output changes；
- 与 Pitaya 或 Nakama 的 public API compatibility；
- 任何 sensitive runtime state exposure。
