# Pitaya 对齐的运行时可观测性边界闸门

状态：Accepted v0.1
最后更新：2026-06-01
范围：在 session binding、kick/disconnect 与 session data 后续方向选择之后，仅定义 Pitaya 对齐的 runtime observability 词汇边界
依赖：`decisions/ADR-0178-select-next-pitaya-aligned-direction-after-session-binding-kick-disconnect-session-data-map.md`、`decisions/ADR-0177-pitaya-aligned-session-binding-kick-disconnect-session-data-source-first-map.md`、`docs/minimum-operations-inspection-surface-gate.md`、`docs/runtime-runbook.md`、`docs/reference-game-server-alignment.md`、`.arch/reference.yaml`
规范决策：`ADR-0179`

英文文件 `docs/pitaya-aligned-runtime-observability-boundary-gate.md` 是权威版本；本文是对应的简体中文翻译。

本文只定义 runtime observability 词汇闸门。它不实现 runtime endpoint behavior、metrics endpoints、tracing pipelines、observability pipelines、dashboards、admin console behavior、player/session/token inspectors、event/audit tables、transport behavior changes、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、hosted deployment、SDK publication、release artifacts、distributed runtime behavior 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Pitaya 对齐的 runtime observability 边界闸门记录如下：

```yaml
pitaya_aligned_runtime_observability_boundary_gate: defined
completed_work_item: W-0271
decision: ADR-0179
check_rule: runtime.pitaya_aligned_runtime_observability_boundary_gate
selection_decision: ADR-0178
session_lifecycle_source_first_map_decision: ADR-0177
minimum_operations_inspection_gate_decision: ADR-0152
minimum_operations_inspection_source_first_decision: ADR-0153
standard: docs/pitaya-aligned-runtime-observability-boundary-gate.md
translation: docs/pitaya-aligned-runtime-observability-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: runtime_observability_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_runtime_observability_vocabulary
future_implementation_work_item: W-0272
future_implementation_direction: implement_pitaya_aligned_runtime_observability_source_first_map
allowed_runtime_observability_vocabulary:
  - runtime_observability_boundary
  - operations_snapshot
  - health_readiness_signal
  - version_release_posture
  - configuration_posture
  - route_inventory_snapshot
  - verification_posture
  - redaction_posture
  - deferred_operations_surface
  - node_local_runtime_surface
current_source_first_runtime_observability_mapping:
  minimum_operations_inspection:
    current: node_tools_vibit_inspect_operations
    future_vocabulary: operations_snapshot
    implementation_status: source_first_repository_inspection_only
  health_and_readiness:
    current: existing_local_healthz_and_readyz_endpoint_summary
    future_vocabulary: health_readiness_signal
    implementation_status: no_new_runtime_endpoint_behavior
  version_surface:
    current: existing_local_version_endpoint_summary
    future_vocabulary: version_release_posture
    implementation_status: no_release_or_hosted_surface_change
  configuration_surface:
    current: existing_redacted_configz_endpoint_summary
    future_vocabulary: configuration_posture
    implementation_status: no_secret_or_dsn_exposure
  route_inventory:
    current: source_first_protocol_route_family_inventory
    future_vocabulary: route_inventory_snapshot
    implementation_status: no_protocol_route_or_wire_shape_change
  repository_verification:
    current: tools_vibit_check_work_runtime_memory_schemas_all
    future_vocabulary: verification_posture
    implementation_status: no_new_test_framework_or_dependency
  redaction_policy:
    current: existing_operations_redaction_requirements
    future_vocabulary: redaction_posture
    implementation_status: no_identifier_or_payload_log_safety_expansion
  deferred_operations_surfaces:
    current: admin_metrics_tracing_dashboard_inspectors_deferred
    future_vocabulary: deferred_operations_surface
    implementation_status: planning_vocabulary_only
runtime_endpoint_behavior_added: false
metrics_endpoint_added: false
tracing_pipeline_added: false
observability_pipeline_added: false
dashboard_added: false
admin_console_added: false
player_session_token_inspector_added: false
event_audit_table_added: false
session_binding_behavior_added: false
kick_disconnect_behavior_added: false
session_data_behavior_added: false
session_data_persistence_added: false
acceptor_behavior_added: false
tcp_acceptor_added: false
websocket_behavior_changed: false
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
remote_connection_handoff_added: false
distributed_session_routing_added: false
distributed_group_implementation_added: false
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
release_artifact_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Purpose

`ADR-0178` 在 session lifecycle source-first map 之后选择 runtime observability 作为下一项 Pitaya 对齐方向。主要风险是 agent 可能把 observability 误读为可以添加 metrics endpoint、tracing pipeline、dashboard behavior、admin operations 或 live state inspectors。

这个闸门只记录词汇和映射。它把既有 source-first operations inspection surface 连接到未来 Pitaya 对齐的 runtime observability 规划，同时保持当前 local alpha runtime behavior 不变。

## 3. Vocabulary

允许使用的 runtime observability 词汇：

- `runtime_observability_boundary`：未来规划中拥有 runtime inspection 与 operational visibility 的边界。
- `operations_snapshot`：未来规划中对 runtime posture、route families、verification 和 deferred operations state 的 source-first 快照。
- `health_readiness_signal`：未来规划中的 liveness/readiness 状态；当前仍是 `/healthz` 和 `/readyz`。
- `version_release_posture`：未来规划中的 source/runtime version posture；当前仍是 `/version`。
- `configuration_posture`：未来规划中的 redacted configuration state；当前仍是 `/configz`。
- `route_inventory_snapshot`：未来规划中的 committed route family inventory；不改变 protocol 或 dispatch。
- `verification_posture`：未来规划中的 repository checks 与 known warnings。
- `redaction_posture`：未来规划中 observability surfaces 可以暴露什么。
- `deferred_operations_surface`：未来规划中仍 deferred 的 admin、metrics、tracing、dashboard 和 inspector surfaces。
- `node_local_runtime_surface`：未来规划中的本地 single-process operational facts；不是 distributed node telemetry。

禁止用法：

- 不要引入来自 Pitaya 或 Nakama 的具体 public API、package、route、method、wire、handler、metrics、trace、dashboard 或 admin compatibility 名称。
- 不要把 runtime observability 词汇当作添加 runtime endpoints、metrics endpoints、tracing pipelines、observability pipelines、dashboards、admin console behavior、live player/session/token inspectors、event/audit tables、dependencies、protocol messages、generated output、persistence、hosted surfaces、SDKs 或 distributed runtime behavior 的许可。
- 此闸门不把 raw tokens、credentials、lookup digests、verifier digests、verifier keys、DSNs、headers、cookies、query strings、subprotocol values、remote addresses、database payloads、local secret file contents、route payloads、session data payloads 或 concrete transport metadata 分类为 log-safe。

## 4. Ownership

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-runtime-observability-boundary-gate.md
  - .arch/reference.yaml
  - .arch/runtime.yaml
source_first_map_candidate_owner:
  - tools/vibit
runtime_endpoint_owner: unchanged
metrics_owner: deferred
tracing_owner: deferred
dashboard_owner: deferred
admin_console_owner: deferred
protocol_owner: unchanged
persistence_owner: unchanged
dependency_owner: unchanged
```

文档和 manifest 可以定义 runtime observability 词汇和当前 source-first 映射。只有后续 bounded work item 明确授权时，`tools/vibit` 才可以输出 runtime observability source-first map。现有 `/healthz`、`/readyz`、`/version` 和 `/configz` 仍是此闸门接受的唯一 runtime troubleshooting endpoints。

## 5. Verification

此边界的仓库检查规则是：

```text
runtime.pitaya_aligned_runtime_observability_boundary_gate
```

验证命令：

```sh
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_runtime_observability_boundary_gate
node tools/vibit check change define-pitaya-aligned-runtime-observability-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

此闸门不添加 Go runtime behavior，因此不要求 Go tests。
