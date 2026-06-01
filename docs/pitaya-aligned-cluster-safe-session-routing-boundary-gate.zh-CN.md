# Pitaya-Aligned Cluster-Safe Session Routing Boundary Gate 中文版

状态：Accepted v0.1
最后更新：2026-06-01
范围：distributed group and broadcast source-first map 之后，使用 Pitaya-aligned cluster-safe session routing vocabulary 的 gate-only boundary
依赖：`decisions/ADR-0163-pitaya-aligned-distributed-group-broadcast-source-first-map.md`、`docs/pitaya-aligned-distributed-group-broadcast-boundary-gate.md`、`docs/first-message-connection-binding-gate.md`、`docs/active-connection-registry-gate.md`、`docs/runtime-session-validation-gate.md`、`docs/logout-revocation-active-connection-gate.md`、`docs/bound-identity-route-policy-gate.md`、`docs/reference-game-server-alignment.md`、`.arch/reference.yaml`
权威决策：`ADR-0164`

英文文件 `docs/pitaya-aligned-cluster-safe-session-routing-boundary-gate.md` 是权威版本。本文是简体中文译本。

本文只定义 cluster-safe session routing vocabulary gate。它不实现 cluster-safe session routing、session location registries、connection owner node registries、routing epoch behavior、session route target behavior、remote connection handoff、reconnect routing、distributed session routing、distributed runtime behavior、distributed groups、group membership registries、room broadcast fanout、delivery guarantees、stream subscriptions、service discovery implementation、service registries、service selectors、node identity、server-to-server RPC、remote calls、frontend/backend server roles、runtime endpoint behavior、metrics endpoints、observability pipelines、dashboards、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、authentication/session behavior changes、SDK publication、hosted deployments、release artifacts 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Pitaya-aligned cluster-safe session routing boundary gate 记录是：

```yaml
pitaya_aligned_cluster_safe_session_routing_boundary_gate: defined
completed_work_item: W-0256
decision: ADR-0164
check_rule: runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate
distributed_group_broadcast_source_first_map_decision: ADR-0163
distributed_group_broadcast_source_first_map_check_rule: runtime.pitaya_aligned_distributed_group_broadcast_source_first_map
standard: docs/pitaya-aligned-cluster-safe-session-routing-boundary-gate.md
translation: docs/pitaya-aligned-cluster-safe-session-routing-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: cluster_safe_session_routing_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_cluster_safe_session_routing_vocabulary
future_implementation_work_item: W-0257
future_implementation_direction: implement_pitaya_aligned_cluster_safe_session_routing_source_first_map
allowed_cluster_safe_session_routing_vocabulary:
  - cluster_safe_session_routing
  - session_location
  - connection_owner_node
  - routing_epoch
  - session_route_target
  - connection_handoff
  - reconnect_route
related_vocabulary:
  - connection_id
  - connection_epoch
  - first_message_connection_binding
  - active_connection_registry
  - runtime_session
  - bound_connection_identity
  - request_token_identity
  - session_validated_identity
  - single_process_connection_binding
  - frontend_server
  - backend_server
  - service_discovery
  - server_to_server_rpc
  - remote_call
  - distributed_group
  - room_broadcast
current_single_process_session_routing_mapping:
  connection_id:
    current: server_observed_connection_id_epoch
    future_vocabulary: connection_owner_node
    implementation_status: current_single_process_connection_binding
  connection_epoch:
    current: server_observed_connection_id_epoch
    future_vocabulary: routing_epoch
    implementation_status: current_single_process_connection_binding
  first_message_connection_binding:
    current: current_single_process_connection_binding
    future_vocabulary: session_location
    implementation_status: active_connection_registry_single_process
  active_connection_registry:
    current: active_connection_registry_single_process
    future_vocabulary: session_location
    implementation_status: no_cross_node_session_location
  runtime_session:
    current: metadata_only_session_id_not_routing_proof
    future_vocabulary: session_route_target
    implementation_status: no_cluster_route_target
  bound_connection_identity:
    current: current_single_process_connection_binding
    future_vocabulary: session_location
    implementation_status: no_distributed_session_routing
  request_token_identity:
    current: request_level_identity_not_cluster_route
    future_vocabulary: session_route_target
    implementation_status: no_cluster_route_target
  session_validated_identity:
    current: request_validation_status_not_cluster_route
    future_vocabulary: session_route_target
    implementation_status: no_cross_node_session_location
  connection_handoff:
    current: no_remote_connection_handoff
    future_vocabulary: connection_handoff
    implementation_status: deferred_future_architecture_reference
  reconnect_route:
    current: reconnect_epoch_local_only_not_cluster_routing
    future_vocabulary: reconnect_route
    implementation_status: deferred_future_architecture_reference
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

`ADR-0163` 已经通过 `node tools/vibit inspect pitaya-groups --json` 让 distributed group and broadcast vocabulary 可检查。该 map 有意把 cluster-safe session routing 留作下一组 deferred Pitaya-aligned planning vocabulary。Session routing 风险很高，因为它容易被误解成可以添加 cross-node connection lookup、remote connection handoff、reconnect routing、service discovery、RPC 或 transport-level authentication behavior。

本 gate 在任何实现之前记录 vocabulary 和 current mapping。它把当前 single-process connection binding、active connection registry、runtime session validation 和 request-level identity surfaces 映射到 future session routing concepts，但不改变 runtime behavior。

## 3. Session Routing Vocabulary

允许的 cluster-safe session routing vocabulary：

- `cluster_safe_session_routing`：用于 future planning 的 across-node validated session 或 connection target routing vocabulary。本 slice 不添加 behavior。
- `session_location`：用于 future planning 的 validated runtime session 或 bound connection 当前所在位置 vocabulary。本 slice 不添加 registry、cache、table 或 service discovery record。
- `connection_owner_node`：用于 future planning 的 open connection owner node vocabulary。当前 runtime 没有 node ownership registry。
- `routing_epoch`：用于 future planning、防止 stale routing decisions 的 vocabulary。当前 `connection_epoch` 仍是 single-process connection lifecycle metadata。
- `session_route_target`：用于 future planning 的 application-selected route target vocabulary。当前 session metadata 不是 cluster route target。
- `connection_handoff`：用于 future planning 的 connection handling handoff vocabulary。本 slice 不添加 remote call、socket migration 或 close policy。
- `reconnect_route`：用于 future planning 的 reconnect direction vocabulary。当前 reconnect 和 epoch behavior 仍是 local 且 non-cluster。

相关 vocabulary：

- `connection_id` 和 `connection_epoch`：当前 single-process runtime 中 server-observed connection lifecycle metadata。
- `first_message_connection_binding`：当前把 authenticated identity 绑定到 WebSocket connection 的 application/protocol posture。
- `active_connection_registry`：当前 server-observed open connections 的 single-process runtime state vocabulary。
- `runtime_session`、`bound_connection_identity`、`request_token_identity` 和 `session_validated_identity`：当前 authentication/session validation concepts；本 slice 中都不是 cluster route target。
- `frontend_server`、`backend_server`、`service_discovery`、`server_to_server_rpc`、`remote_call`、`distributed_group` 和 `room_broadcast`：既有 Pitaya-aligned vocabulary families，implementation 仍然 deferred。

禁止的 vocabulary 使用：

- 不要从 Pitaya 或 Nakama 引入 concrete public API、package、route、method、wire、session、channel、registry、selector、handoff 或 configuration compatibility names。
- 不要把 session routing vocabulary 当作添加 cross-node registries、session-location tables、node identity、service discovery、RPC、remote calls、transport carriers、protocol messages、generated output、persistence、dependencies、routing caches、handoff workers、reconnect routers 或 cluster runtime behavior 的许可。
- 不要把 metadata-only `session_id`、client-supplied connection metadata、target-scope metadata 或 transport metadata 当作 routing proof。
- 不要绕过 application dispatch、authenticated request validation、bound identity route policy、generated output rules、redaction rules 或 module ownership。

## 4. Current Mapping

```yaml
current_single_process_session_routing_mapping:
  connection_id:
    current: server-observed connection id and epoch metadata
    future_vocabulary: connection_owner_node
    status: current_single_process_connection_binding
  connection_epoch:
    current: server-observed connection epoch metadata
    future_vocabulary: routing_epoch
    status: current_single_process_connection_binding
  first_message_connection_binding:
    current: application-owned first-message bind posture for one process
    future_vocabulary: session_location
    status: active_connection_registry_single_process
  active_connection_registry:
    current: single-process active connection registry vocabulary
    future_vocabulary: session_location
    status: no_cross_node_session_location
  runtime_session:
    current: session validation metadata; metadata-only session id is not proof
    future_vocabulary: session_route_target
    status: no_cluster_route_target
  bound_connection_identity:
    current: one-process bound identity vocabulary
    future_vocabulary: session_location
    status: no_distributed_session_routing
  request_token_identity:
    current: request-level token identity
    future_vocabulary: session_route_target
    status: no_cluster_route_target
  session_validated_identity:
    current: validation status vocabulary, not a route record
    future_vocabulary: session_route_target
    status: no_cross_node_session_location
  connection_handoff:
    current: no remote connection handoff
    future_vocabulary: connection_handoff
    status: deferred_future_architecture_reference
  reconnect_route:
    current: reconnect epoch behavior is local, not cluster routing
    future_vocabulary: reconnect_route
    status: deferred_future_architecture_reference
```

## 5. Ownership

Cluster-safe session routing vocabulary ownership：

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-cluster-safe-session-routing-boundary-gate.md
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

- Documentation 和 manifests 可以定义 cluster-safe session routing vocabulary 与 mapping。
- 若后续 implementation work item 授权，`tools/vibit` 可以输出 source-first cluster-safe session routing map。
- Runtime、transport、protocol、repository、persistence、generated output、startup wiring、dependencies、service discovery、RPC、remote calls、frontend/backend role behavior、distributed group behavior 和 room broadcast behavior 不因本 gate 改变。
- Domain modules 默认不会获得 session routing、connection targeting、reconnect routing、handoff、service discovery、RPC 或 transport ownership。

## 6. Nakama And Pitaya Mapping

Nakama 仍是近期 capability breadth 的主要产品参考。Pitaya 仍是 acceptors、sessions、route handlers、frontend/backend roles、RPC/remotes、service discovery、groups、broadcast 和 cluster routing 的 architecture vocabulary reference。

采纳为 vocabulary：

- cluster-safe session routing 作为 future architecture-planning vocabulary；
- session location、connection owner node、routing epoch、session route target、connection handoff 和 reconnect route 作为 planning vocabulary；
- 把 current single-process connection binding、active connection registry 和 request/session validation vocabulary 映射到 deferred cluster-routing concepts。

适配到 vibit：

- 任何 future routing model 都必须保留 vibit application-owned identity validation、bound identity route policy、module boundaries、generated output rules、redaction 和 repository checks。
- 当前 single-process runtime、connection binding、active connection registry 和 request-level validation 仍是具体实现。
- 任何 future cluster-safe session routing implementation 都必须先经过单独 gate 和 verification。

当前拒绝：

- direct Pitaya 或 Nakama API compatibility；
- Pitaya 或 Nakama package、method、session、route、registry、selector 或 handoff naming compatibility；
- cluster-safe session routing behavior；
- session location registry、connection owner node registry、routing epoch behavior、session route targets、remote connection handoff、reconnect routing 或 distributed session routing；
- service discovery、RPC、remote calls、frontend/backend process split、distributed groups、room broadcast fanout 或 distributed runtime behavior；
- routing 的 protocol messages or routes、generated output、persistence、migrations、dependencies、hosted deployment、SDK publication 或 release artifacts。

## 7. Future Implementation Work

打开：

```text
M-185/W-0257 Implement Pitaya-aligned cluster-safe session routing source-first map
```

后续 work item 可以：

- 添加 cluster-safe session routing vocabulary 的 source-first repository inspection map；
- 摘要 current connection id、connection epoch、first-message binding、active connection registry、runtime session validation、bound identity 和 request identity mappings；
- 更新 runbooks 和 acceptance docs 指向 cluster-safe session routing map；
- 添加 repository checks，确认该 map 仍保持 gate-only 且 redacted。

后续 work item 禁止：

- 添加 cluster-safe session routing behavior；
- 添加 session location registries 或 connection owner node registries；
- 添加 routing epoch behavior、session route targets、remote connection handoff 或 reconnect routing；
- 添加 service discovery implementation、service registries、selectors、node identity 或 topology behavior；
- 添加 server-to-server RPC implementation 或 remote call behavior；
- 添加 frontend/backend server role implementation；
- 添加 distributed runtime implementation；
- 添加 distributed group implementation、group membership registries、stream subscriptions、room broadcast fanout 或 delivery guarantees；
- 添加 runtime endpoint behavior、metrics endpoints、observability pipelines、dashboards、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、authentication/session behavior changes、SDK publication、hosted deployments、release artifacts 或 direct Nakama/Pitaya API compatibility。

## 8. Verification Expectations

本 gate 应验证：

- `runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate` 已注册。
- `ADR-0164`、本 standard、简体中文译本、change artifacts 和 conversation memory 存在。
- W-0256 已完成，W-0257 是 next-ready。
- Current single-process connection/session mapping 已记录。
- 对 cluster-safe routing behavior、registries、handoff、distributed runtime、service discovery、RPC、remote calls、protocol、generated output、persistence、dependencies、hosted surfaces、SDKs 和 direct compatibility 的 deferrals 仍然明确。

Required commands：

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate
node tools/vibit check change define-pitaya-aligned-cluster-safe-session-routing-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

## 9. Stop Conditions

在添加以下内容前必须停止并等待后续 bounded work item：

- cluster-safe session routing behavior；
- session location registry 或 connection owner node registry behavior；
- routing epoch behavior、route target resolution、remote connection handoff 或 reconnect route behavior；
- service discovery、service registry、service selector、node identity、server identity、RPC、remote calls 或 frontend/backend roles；
- distributed runtime behavior、distributed groups、group membership registries、room broadcast fanout、delivery guarantees 或 stream subscriptions；
- runtime endpoint behavior、protocol messages、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、hosted surfaces、SDKs 或 direct compatibility。

## 10. Non-Authorization

This is a boundary-only standard。它只授权 vocabulary、mapping、manifests、checks、ADRs 和 memory。它不授权 runtime behavior、protocol behavior、generated output、persistence、service discovery、RPC、remote calls、distributed runtime behavior、hosted deployment、SDK publication、release execution 或 direct Nakama/Pitaya API compatibility。
