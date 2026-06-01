# Pitaya-Aligned Route Handler Pipeline Boundary Gate 中文版

状态：Accepted v0.1
最后更新：2026-06-01
范围：cluster-safe session routing source-first map 之后，使用 Pitaya-aligned route handler、handler pipeline、serializer 和 message forwarding vocabulary 的 gate-only boundary
依赖：`decisions/ADR-0166-select-next-pitaya-aligned-direction-after-cluster-safe-session-routing-map.md`、`decisions/ADR-0165-pitaya-aligned-cluster-safe-session-routing-source-first-map.md`、`docs/runtime-protocol-adapter.md`、`docs/game-protocol.md`、`docs/pitaya-aligned-cluster-safe-session-routing-boundary-gate.md`、`docs/reference-game-server-alignment.md`、`.arch/reference.yaml`
Canonical decision：`ADR-0167`

英文文件 `docs/pitaya-aligned-route-handler-pipeline-boundary-gate.md` 是权威版本。本文是简体中文译本。

本文只定义 route handler pipeline vocabulary gate。它不实现 route handlers、handler routing behavior、pipeline middleware behavior、serializer behavior、message forwarding behavior、backend route targeting、cluster-safe session routing behavior、session location registries、connection owner node registries、routing epoch behavior、session route targets、remote connection handoff、reconnect routing、distributed session routing、distributed runtime behavior、distributed groups、group membership registries、room broadcast fanout、delivery guarantees、stream subscriptions、service discovery implementation、service registries、service selectors、node identity、server-to-server RPC、remote calls、frontend/backend server roles、runtime endpoint behavior、metrics endpoints、observability pipelines、dashboards、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、authentication/session behavior changes、SDK publication、hosted deployments、release artifacts 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Pitaya-aligned route handler pipeline boundary gate 记录是：

```yaml
pitaya_aligned_route_handler_pipeline_boundary_gate: defined
completed_work_item: W-0259
decision: ADR-0167
check_rule: runtime.pitaya_aligned_route_handler_pipeline_boundary_gate
previous_direction_decision: ADR-0166
cluster_safe_session_routing_source_first_map_decision: ADR-0165
cluster_safe_session_routing_source_first_map_check_rule: runtime.pitaya_aligned_cluster_safe_session_routing_source_first_map
standard: docs/pitaya-aligned-route-handler-pipeline-boundary-gate.md
translation: docs/pitaya-aligned-route-handler-pipeline-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: route_handler_pipeline_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_route_handler_pipeline_vocabulary
future_implementation_work_item: W-0260
future_implementation_direction: implement_pitaya_aligned_route_handler_pipeline_source_first_map
allowed_route_handler_pipeline_vocabulary:
  - route_handler
  - route_key
  - handler_dispatch
  - handler_pipeline
  - pipeline_step
  - serializer_boundary
  - message_forwarding
  - route_target
related_vocabulary:
  - protocol_envelope
  - route_request
  - application_dispatch
  - command_handler
  - query_handler
  - protocol_bridge
  - target_scope
  - frontend_server
  - backend_server
  - server_to_server_rpc
  - remote_call
  - service_discovery
  - cluster_safe_session_routing
current_single_process_route_handler_mapping:
  protocol_envelope:
    current: kind_module_name_structured_routing
    future_vocabulary: route_key
    implementation_status: current_protocol_adapter_owned_shape
  route_request:
    current: app_route_request_handoff
    future_vocabulary: route_handler
    implementation_status: current_application_dispatch
  application_dispatch:
    current: explicit_command_query_dispatch
    future_vocabulary: handler_dispatch
    implementation_status: no_pitaya_handler_pipeline
  transactional_dispatch:
    current: application_unit_of_work_wrapper
    future_vocabulary: pipeline_step
    implementation_status: current_vibit_transaction_boundary_only
  protocol_bridge:
    current: explicit_generated_payload_bridge
    future_vocabulary: serializer_boundary
    implementation_status: no_pluggable_serializer_behavior
  outbound_message:
    current: server_push_intent_to_protocol_envelope
    future_vocabulary: message_forwarding
    implementation_status: no_cross_node_forwarding
route_handler_implementation_added: false
handler_routing_behavior_added: false
handler_pipeline_behavior_added: false
pipeline_middleware_behavior_added: false
serializer_behavior_added: false
message_forwarding_behavior_added: false
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

`ADR-0166` 在 cluster-safe session routing source-first map 之后选择 route handler pipeline vocabulary 作为下一项 Pitaya-aligned direction。该 vocabulary 对未来分布式架构有价值，因为它把 client-facing route receipt、handler dispatch、handler pipeline、serialization 和 forwarding decisions 分开表达。

风险是 agent 可能把这些词当成替换 vibit 当前 protocol adapter、application dispatch、generated payload bridge、transaction wrapper 或 server-push path 的许可。这个 gate 只记录 vocabulary 和 mapping，当前 single-process WebSocket Protobuf route flow 不变。

## 3. Route Handler Pipeline Vocabulary

允许的 route handler pipeline vocabulary：

- `route_handler`：未来用于按 route key 选择 application-facing handler 的规划词汇。本 slice 不添加 handler code。
- `route_key`：未来用于 logical route identity 的规划词汇。当前 route identity 仍是结构化 `kind`、`module` 和 `name` fields。
- `handler_dispatch`：未来用于 handler selection 的规划词汇。当前 dispatch 仍是 explicit application dispatcher。
- `handler_pipeline`：未来用于 pre-handler 或 post-handler 处理顺序的规划词汇。本 slice 不添加 middleware behavior。
- `pipeline_step`：未来用于 bounded pipeline unit 的规划词汇。当前 transactional dispatch 不是通用 middleware。
- `serializer_boundary`：未来用于 encode/decode ownership 的规划词汇。当前 Protobuf bridge functions 仍是唯一 concrete serializer boundary。
- `message_forwarding`：未来用于把 message 转交给另一个 owner 或 node 的规划词汇。当前 runtime 没有 cross-node forwarding。
- `route_target`：未来用于 handler placement 或 target selection 的规划词汇。当前 target scope metadata 不是 backend route targeting。

相关 vocabulary：

- `protocol_envelope`、`route_request`、`application_dispatch`、`command_handler`、`query_handler`、`protocol_bridge` 和 `target_scope`：当前 vibit route-flow concepts。
- `frontend_server`、`backend_server`、`server_to_server_rpc`、`remote_call`、`service_discovery` 和 `cluster_safe_session_routing`：既有 Pitaya-aligned vocabulary families，implementation 仍然 deferred。

禁止用法：

- 不要从 Pitaya 或 Nakama 引入 concrete public API、package、route、method、wire、handler、pipeline、serializer、forwarding、registry、selector 或 configuration compatibility names。
- 不要把 route handler pipeline vocabulary 当成添加 handler routing behavior、middleware chains、serializer plugins、forwarding workers、backend route targeting、service discovery、RPC、remote calls、protocol messages、generated output、persistence、dependencies、topology 或 distributed runtime behavior 的许可。
- 不要把 domain behavior 移入 transport、Protobuf adapters、serializer boundaries 或 process startup。
- 不要绕过 application dispatch、request/session validation、bound identity route policy、generated output rules、redaction rules 或 module ownership。

## 4. Current Mapping

```yaml
current_single_process_route_handler_mapping:
  protocol_envelope:
    current: kind/module/name structured routing fields
    future_vocabulary: route_key
    status: current_protocol_adapter_owned_shape
  route_request:
    current: explicit application route request handoff
    future_vocabulary: route_handler
    status: current_application_dispatch
  application_dispatch:
    current: explicit command/query route registration and dispatch
    future_vocabulary: handler_dispatch
    status: no_pitaya_handler_pipeline
  transactional_dispatch:
    current: application-owned unit-of-work wrapper for commands
    future_vocabulary: pipeline_step
    status: current_vibit_transaction_boundary_only
  protocol_bridge:
    current: explicit generated Protobuf payload bridge functions
    future_vocabulary: serializer_boundary
    status: no_pluggable_serializer_behavior
  outbound_message:
    current: server-push intent converted to protocol envelope
    future_vocabulary: message_forwarding
    status: no_cross_node_forwarding
```

## 5. Ownership

Route handler pipeline vocabulary ownership：

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-route-handler-pipeline-boundary-gate.md
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

规则：

- Documentation 和 manifests 可以定义 route handler pipeline vocabulary 与当前 mapping。
- 后续 implementation work item 明确授权时，`tools/vibit` 可以输出 source-first route handler pipeline map。
- Runtime、transport、protocol、repository、persistence、generated output、startup wiring、dependencies、service discovery、RPC、remote calls、frontend/backend role behavior、cluster-safe session routing、distributed group behavior 和 room broadcast behavior 不因本 gate 改变。
- Domain modules 默认不会获得 handler pipeline、serializer、forwarding、backend route targeting、service discovery、RPC 或 transport ownership。

## 6. Nakama And Pitaya Mapping

Nakama 仍是近期 capability breadth 的 primary product reference。Pitaya 仍是 acceptors、sessions、route handlers、frontend/backend roles、RPC/remotes、service discovery、groups、broadcast、cluster routing、handler pipelines、serializers 和 forwarding 的 architecture vocabulary reference。

采纳为 vocabulary：

- route handler pipeline vocabulary，作为未来 architecture planning vocabulary；
- route key、handler dispatch、handler pipeline、pipeline step、serializer boundary、message forwarding 和 route target，作为 planning vocabulary；
- 把 vibit 现有 protocol adapter、route request、application dispatcher、transactional dispatch、protocol bridges 和 server-push intent 映射到 deferred Pitaya-aligned concepts。

按 vibit 方式适配：

- 当前 route identity 仍是结构化 `kind`、`module` 和 `name` fields。
- 当前 dispatch 仍由 application 层显式拥有。
- 当前 serialization 仍由 Protobuf adapter 通过 explicit bridge functions 拥有。
- 当前 server push 仍是 single-process，不暗示 cross-node forwarding。
- 任何未来 route handler pipeline implementation 都必须先经过单独 gate 和 verification。

当前拒绝：

- direct Pitaya 或 Nakama API compatibility；
- Pitaya 或 Nakama package、method、route、handler、serializer、pipeline、forwarding、registry、selector 或 configuration naming compatibility；
- route handler implementation、handler routing behavior、pipeline middleware behavior、serializer behavior、message forwarding behavior 或 backend route targeting；
- service discovery、RPC、remote calls、frontend/backend process split、cluster-safe session routing、distributed groups、room broadcast fanout 或 distributed runtime behavior；
- 为 route handler pipelines 添加 protocol messages or routes、generated output、persistence、migrations、dependencies、hosted deployment、SDK publication 或 release artifacts。

## 7. Future Implementation Work

打开：

```text
M-188/W-0260 Implement Pitaya-aligned route handler pipeline source-first map
```

未来 work item 可以：

- 添加 route handler pipeline vocabulary 的 source-first repository inspection map；
- 总结当前 protocol envelope、route request、application dispatch、transactional dispatch、protocol bridge 和 outbound message mapping；
- 更新 runbooks 和 acceptance docs 指向 route handler pipeline map；
- 添加 repository checks，验证该 map 仍保持 gate-only 和 redacted。

未来 work item 不能：

- 添加 route handler implementation；
- 添加 handler routing behavior、handler pipeline behavior、pipeline middleware behavior、serializer behavior、message forwarding behavior 或 backend route targeting；
- 添加 cluster-safe session routing behavior、session location registries、connection owner node registries、routing epochs、session route targets、remote connection handoff、reconnect routing 或 distributed session routing；
- 添加 service discovery implementation、service registries、selectors、node identity 或 topology behavior；
- 添加 server-to-server RPC implementation 或 remote call behavior；
- 添加 frontend/backend server role implementation；
- 添加 distributed runtime implementation；
- 添加 distributed group implementation、group membership registries、stream subscriptions、room broadcast fanout 或 delivery guarantees；
- 添加 runtime endpoint behavior、metrics endpoints、observability pipelines、dashboards、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、authentication/session behavior changes、SDK publication、hosted deployments、release artifacts 或 direct Nakama/Pitaya API compatibility。

## 8. Verification Expectations

该 gate 应验证：

- `runtime.pitaya_aligned_route_handler_pipeline_boundary_gate` 已注册。
- `ADR-0167`、本 standard、简体中文译本、change artifacts 和 conversation memory 存在。
- W-0259 已完成，W-0260 是 next-ready。
- 当前 route dispatch 和 protocol adapter mapping 已记录。
- route handler implementation、handler routing behavior、handler pipeline behavior、pipeline middleware behavior、serializer behavior、message forwarding behavior、backend route targeting、cluster-safe session routing、distributed runtime、service discovery、RPC、remote calls、protocol、generated output、persistence、dependencies、hosted surfaces、SDKs 和 direct compatibility 的 deferrals 保持显式。

Required commands：

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_route_handler_pipeline_boundary_gate
node tools/vibit check change define-pitaya-aligned-route-handler-pipeline-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

## 9. Stop Conditions

在以下情况前必须停止并要求新的 bounded work item：

- 实现 route handlers、handler routing behavior、handler pipelines、pipeline middleware、serializer behavior、message forwarding 或 backend route targeting；
- 改变 route identity、protocol envelope shape、Protobuf sources、generated output、application dispatch semantics、transaction behavior、protocol bridge behavior 或 outbound delivery behavior；
- 添加 service discovery、RPC、remote calls、frontend/backend server role behavior、cluster-safe session routing、distributed groups、room broadcast fanout 或 distributed runtime behavior；
- 添加 dependencies、migrations、repository interfaces、PostgreSQL adapters、hosted deployment surfaces、SDK publication、release artifacts 或 direct Nakama/Pitaya API compatibility。
