# Pitaya-Aligned Distributed Group And Broadcast Boundary Gate 中文版

状态：Accepted v0.1
最后更新：2026-06-01
范围：service discovery source-first map 之后，使用 Pitaya-aligned distributed group 和 room broadcast vocabulary 的 gate-only boundary
依赖：`decisions/ADR-0161-pitaya-aligned-service-discovery-source-first-map.md`、`docs/pitaya-aligned-service-discovery-boundary-gate.md`、`docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md`、`docs/first-server-push-realtime-messaging-gate.md`、`docs/realtime-protocol-websocket-outbound-delivery-gate.md`、`docs/game-protocol.md`、`docs/reference-game-server-alignment.md`、`.arch/reference.yaml`
权威决策：`ADR-0162`

英文文件 `docs/pitaya-aligned-distributed-group-broadcast-boundary-gate.md` 是权威版本。本文是简体中文译本。

本文只定义 distributed group and broadcast vocabulary gate。它不实现 distributed groups、room broadcast fanout、delivery guarantees、stream subscriptions、group membership registries、groups、parties、chat rooms、matchmaking、match runtime、service discovery implementation、service registries、service selectors、node identity、server-to-server RPC、remote calls、frontend/backend server roles、distributed runtime behavior、cluster-safe session routing、runtime endpoint behavior、metrics endpoints、observability pipelines、dashboards、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、authentication/session behavior changes、SDK publication、hosted deployments、release artifacts 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Pitaya-aligned distributed group and broadcast boundary gate 记录是：

```yaml
pitaya_aligned_distributed_group_broadcast_boundary_gate: defined
completed_work_item: W-0254
decision: ADR-0162
check_rule: runtime.pitaya_aligned_distributed_group_broadcast_boundary_gate
service_discovery_source_first_map_decision: ADR-0161
service_discovery_source_first_map_check_rule: runtime.pitaya_aligned_service_discovery_source_first_map
standard: docs/pitaya-aligned-distributed-group-broadcast-boundary-gate.md
translation: docs/pitaya-aligned-distributed-group-broadcast-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: distributed_group_broadcast_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_distributed_group_broadcast_vocabulary
future_implementation_work_item: W-0255
future_implementation_direction: implement_pitaya_aligned_distributed_group_broadcast_source_first_map
allowed_group_broadcast_vocabulary:
  - distributed_group
  - room_broadcast
  - broadcast_target
  - group_membership
  - broadcast_fanout
related_vocabulary:
  - target_scope
  - server_push_intent
  - route_handler
  - module_handler
  - frontend_server
  - backend_server
  - service_discovery
  - server_to_server_rpc
  - remote_call
  - cluster_safe_session_routing
current_single_process_group_broadcast_mapping:
  target_scope:
    current: protocol_envelope_target_metadata_only
    future_vocabulary: broadcast_target
    implementation_status: current_single_process_intent_only
  server_push_intent:
    current: application_owned_realtime_outbound_intent_single_process_delivery
    future_vocabulary: room_broadcast
    implementation_status: current_single_process_delivery_only
  distributed_group:
    current: no_distributed_group_current_single_process_target_scope_only
    future_vocabulary: distributed_group
    implementation_status: deferred_future_architecture_reference
  group_membership:
    current: no_group_membership_registry_or_subscription_state
    future_vocabulary: group_membership
    implementation_status: deferred_future_architecture_reference
  room_broadcast:
    current: no_room_broadcast_fanout_current_server_push_intent_only
    future_vocabulary: room_broadcast
    implementation_status: deferred_future_architecture_reference
  broadcast_target:
    current: target_scope_values_are_metadata_not_distributed_targets
    future_vocabulary: broadcast_target
    implementation_status: current_metadata_only
  broadcast_fanout:
    current: no_cluster_fanout_no_delivery_guarantee
    future_vocabulary: broadcast_fanout
    implementation_status: deferred_future_architecture_reference
distributed_group_implementation_added: false
distributed_groups_added: false
group_membership_registry_added: false
room_broadcast_fanout_added: false
broadcast_delivery_guarantee_added: false
stream_subscription_added: false
groups_parties_chat_runtime_behavior_added: false
match_runtime_behavior_added: false
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

`ADR-0161` 已经通过 `node tools/vibit inspect pitaya-discovery --json` 让 service discovery vocabulary 可检查。该 map 有意把 distributed groups 和 room broadcast 留作后续 vocabulary。下一类高风险 Pitaya 概念是 broadcast：它容易被误解成可以添加 group membership state、room routing、fanout workers、delivery guarantees、stream subscriptions 或 cluster-safe session routing。

本 gate 在任何实现之前记录 vocabulary 和 current mapping。它把当前 single-process target-scope 和 server-push intent surfaces 映射到 future distributed group and broadcast concepts，但不改变 runtime behavior。

## 3. Group And Broadcast Vocabulary

允许的 group and broadcast vocabulary：

- `distributed_group`：用于 future architecture planning 的 cluster-aware grouping target。它在本 slice 中不是当前 group、party、chat room、match 或 stream implementation。
- `room_broadcast`：未来 broadcast vocabulary item。它在本 slice 中不是 delivery guarantee、wire route 或 fanout implementation。
- `broadcast_target`：未来 fanout planning 的 targeting vocabulary item。当前 `target_scope` values 仍是 protocol metadata，不标识 distributed targets。
- `group_membership`：未来 membership vocabulary item。它在本 slice 中不是 registry、subscription table、persistence model 或 runtime state。
- `broadcast_fanout`：未来 delivery topology vocabulary item。它在本 slice 中不是 worker、queue、retry policy、ordering policy 或 cluster mechanism。

相关 vocabulary：

- `target_scope`：game protocol 和 envelope posture 定义的当前 protocol target metadata。
- `server_push_intent`：当前 application-owned outbound realtime intent 和 single-process delivery posture。
- `route_handler` 和 `module_handler`：当前 in-process application/module behavior，未来可能产生 broadcast intent。
- `frontend_server`、`backend_server`、`service_discovery`、`server_to_server_rpc`、`remote_call` 和 `cluster_safe_session_routing`：既有 Pitaya-aligned vocabulary families，implementation 仍然 deferred。

禁止的 vocabulary 使用：

- 不要从 Pitaya 或 Nakama 引入 concrete public API、package、route、method、wire、group、room、channel、stream 或 configuration compatibility names。
- 不要把 group 或 broadcast vocabulary 当作添加 membership registries、room state、stream subscriptions、fanout workers、delivery guarantees、retries、ordering、durable offsets、queueing、cluster routing、service discovery、RPC、remote calls、protocol messages、generated output、persistence 或 dependencies 的许可。
- 不要用 future broadcast vocabulary 绕过 module contracts、application dispatch boundaries、authentication/session validation gates、permission checks、generated output rules、redaction rules 或 repository ownership。

## 4. Current Mapping

```yaml
current_single_process_group_broadcast_mapping:
  target_scope:
    current: Protobuf envelope target metadata and game-protocol target vocabulary
    future_vocabulary: broadcast_target
    status: current_single_process_intent_only
  server_push_intent:
    current: application-owned realtime outbound intent with single-process WebSocket delivery
    future_vocabulary: room_broadcast
    status: current_single_process_delivery_only
  distributed_group:
    current: no distributed group model; current target scope is metadata only
    future_vocabulary: distributed_group
    status: deferred_future_architecture_reference
  group_membership:
    current: no membership registry, subscription table, or distributed group state
    future_vocabulary: group_membership
    status: deferred_future_architecture_reference
  room_broadcast:
    current: no room broadcast fanout; current server push is narrow single-process delivery
    future_vocabulary: room_broadcast
    status: deferred_future_architecture_reference
  broadcast_target:
    current: target scope values are not distributed routing targets
    future_vocabulary: broadcast_target
    status: current_metadata_only
  broadcast_fanout:
    current: no fanout worker, queue, retry, ordering, or delivery guarantee
    future_vocabulary: broadcast_fanout
    status: deferred_future_architecture_reference
```

## 5. Ownership

Group and broadcast vocabulary ownership：

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-distributed-group-broadcast-boundary-gate.md
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

- Documentation 和 manifests 可以定义 distributed group and broadcast vocabulary 与 mapping。
- 若后续 implementation work item 授权，`tools/vibit` 可以输出 source-first distributed group and broadcast map。
- Runtime、protocol、repository、persistence、generated output、startup wiring、dependencies、RPC、remote calls、service discovery、frontend/backend role behavior 和 cluster-safe session routing 不因本 gate 改变。
- Domain modules 默认不会获得 group、party、chat、stream、match 或 broadcast ownership。Module contracts 仍是 module behavior 和 data ownership 的来源。

## 6. Nakama And Pitaya Mapping

Nakama 仍是近期 capability breadth 的主要产品参考。Pitaya 仍是 acceptors、sessions、route handlers、frontend/backend roles、RPC/remotes、service discovery、groups、broadcast 和 cluster routing 的 architecture vocabulary reference。

采纳为 vocabulary：

- distributed group 作为 future cluster-aware target vocabulary；
- room broadcast 作为 future fanout vocabulary；
- broadcast target、group membership 和 broadcast fanout 作为 planning vocabulary；
- 把 current target-scope metadata 和 server-push intent 映射到 deferred distributed group and broadcast concepts。

适配到 vibit：

- 任何 future group 或 broadcast model 都必须保留 vibit module ownership、application dispatch boundaries、server-authoritative validation、generated output rules、redaction 和 repository checks。
- 当前 single-process runtime 和窄口径 realtime outbound delivery 仍是具体实现。
- 任何 future distributed group 或 room broadcast implementation 都必须先经过单独 gate 和 verification。

当前拒绝：

- direct Pitaya 或 Nakama API compatibility；
- Pitaya 或 Nakama package、method、group、room、stream 或 route naming compatibility；
- distributed group implementation；
- group membership registry 或 subscription state；
- room broadcast fanout；
- delivery guarantees、retries、ordering、durable offsets 或 queueing；
- groups 或 broadcast 的 protocol messages 或 routes；
- service discovery、RPC、remote calls、frontend/backend process split 或 cluster-safe session routing。

## 7. Future Implementation Work

打开：

```text
M-183/W-0255 Implement Pitaya-aligned distributed group and broadcast source-first map
```

后续 work item 可以：

- 添加 distributed group and broadcast vocabulary 的 source-first repository inspection map；
- 摘要 current target-scope metadata、server-push intent 和 single-process delivery mapping；
- 更新 runbooks 和 acceptance docs 指向 distributed group and broadcast map；
- 添加 repository checks，确认该 map 仍保持 gate-only 且 redacted。

后续 work item 禁止：

- 添加 distributed group implementation；
- 添加 room broadcast fanout；
- 添加 delivery guarantees、retries、ordering、durable offsets、queueing 或 backpressure behavior；
- 添加 stream subscriptions、group membership registries、groups、parties、chat rooms、matchmaking 或 match runtime behavior；
- 添加 service discovery implementation、service registries、selectors、node identity 或 topology behavior；
- 添加 server-to-server RPC implementation；
- 添加 remote call behavior；
- 添加 frontend/backend server role implementation；
- 添加 distributed runtime implementation；
- 添加 cluster-safe session routing behavior；
- 添加 runtime endpoint behavior、metrics endpoints、observability pipelines、dashboards、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、authentication/session behavior changes、SDK publication、hosted deployments、release artifacts 或 direct Nakama/Pitaya API compatibility。

## 8. Verification Expectations

本 gate 应验证：

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_distributed_group_broadcast_boundary_gate
node tools/vibit check change define-pitaya-aligned-distributed-group-broadcast-boundary-gate --json
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

- distributed group implementation；
- room broadcast fanout behavior；
- group membership registry、subscription、stream、chat、party、group、room、matchmaking 或 match runtime behavior；
- delivery guarantees、retries、ordering、durable offsets、queueing 或 backpressure behavior；
- service discovery implementation、registry、selector、membership、heartbeat 或 node identity behavior；
- server-to-server RPC 或 remote call behavior；
- process topology changes；
- new goroutines、listeners、network protocols 或 server roles；
- cluster-safe session routing behavior；
- protocol 或 Protobuf changes；
- repository、adapter、migration、dependency 或 generated-output changes；
- 与 Pitaya 或 Nakama 的 public API compatibility；
- 任何 sensitive runtime state exposure。
