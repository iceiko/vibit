# Pitaya-Aligned Dashboard And Admin Operations Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-01
Scope: Gate-only boundary for using Pitaya-aligned dashboard and admin operations vocabulary after the metrics and tracing source-first map
Depends on: `decisions/ADR-0184-select-next-pitaya-aligned-direction-after-metrics-tracing-map.md`, `decisions/ADR-0183-pitaya-aligned-metrics-tracing-source-first-map.md`, `docs/pitaya-aligned-metrics-tracing-boundary-gate.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0185`

英文文件 `docs/pitaya-aligned-dashboard-admin-operations-boundary-gate.md` 是权威版本。

本文只定义 dashboard and admin operations vocabulary gate。它不实现 runtime endpoint behavior、metrics endpoints、tracing pipelines、observability pipelines、dashboards、admin console behavior、player/session/token inspectors、event/audit tables、transport behavior changes、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、hosted deployment、SDK publication、release artifacts、distributed runtime behavior 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

```yaml
pitaya_aligned_dashboard_admin_operations_boundary_gate: defined
completed_work_item: W-0277
decision: ADR-0185
check_rule: runtime.pitaya_aligned_dashboard_admin_operations_boundary_gate
selection_decision: ADR-0184
metrics_tracing_source_first_map_decision: ADR-0183
metrics_tracing_boundary_gate_decision: ADR-0182
standard: docs/pitaya-aligned-dashboard-admin-operations-boundary-gate.md
translation: docs/pitaya-aligned-dashboard-admin-operations-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: dashboard_admin_operations_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_dashboard_admin_operations_vocabulary
future_implementation_work_item: W-0278
future_implementation_direction: implement_pitaya_aligned_dashboard_admin_operations_source_first_map
allowed_dashboard_admin_operations_vocabulary:
  - dashboard_admin_boundary
  - admin_operation_surface
  - operator_action_boundary
  - dashboard_read_model
  - source_first_operations_snapshot
  - inspector_redaction_posture
  - event_audit_deferral
  - local_operations_diagnostic_surface
  - future_admin_authorization_boundary
  - hosted_operations_deferral
runtime_endpoint_behavior_added: false
metrics_endpoint_added: false
tracing_pipeline_added: false
observability_pipeline_added: false
dashboard_added: false
admin_console_added: false
player_session_token_inspector_added: false
event_audit_table_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
migration_added: false
dependency_added: false
hosted_deployment_added: false
sdk_added: false
release_artifact_added: false
distributed_runtime_implementation_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Purpose

`ADR-0184` 在 metrics and tracing source-first map 之后选择 dashboard and admin operations 作为下一项 Pitaya-aligned direction。

风险在于 agent 可能把 dashboard/admin 词汇当作添加 UI、live operations endpoints、admin authentication、sensitive inspectors、audit tables、deployment surfaces 或第三方 dependencies 的授权。此 gate 只记录 vocabulary 与 source-first mapping，保持当前 local alpha runtime behavior 不变，并为后续 source-first map 做准备。

## 3. Vocabulary

Allowed dashboard and admin operations vocabulary:

- `dashboard_admin_boundary`
- `admin_operation_surface`
- `operator_action_boundary`
- `dashboard_read_model`
- `source_first_operations_snapshot`
- `inspector_redaction_posture`
- `event_audit_deferral`
- `local_operations_diagnostic_surface`
- `future_admin_authorization_boundary`
- `hosted_operations_deferral`

Forbidden vocabulary use:

- 不要引入 Pitaya 或 Nakama 的 concrete public API、package、route、method、wire、handler、dashboard、metrics、trace、admin、console 或 inspector compatibility names。
- 不要用 dashboard/admin vocabulary 添加 runtime endpoints、metrics endpoints、tracing pipelines、observability pipelines、dashboards、admin console behavior、live player/session/token inspectors、event/audit tables、dependencies、protocol messages、generated output、persistence、hosted surfaces、SDKs、release artifacts 或 distributed runtime behavior。
- 不要把 raw tokens、credentials、lookup digests、verifier digests、verifier keys、DSNs、headers、cookies、query strings、subprotocol values、remote addresses、database payloads、local secret file contents、route payloads、session data payloads、dashboard payloads、admin console payloads 或 concrete transport metadata 标为 log-safe。

## 4. Current Mapping

```yaml
current_source_first_dashboard_admin_operations_mapping:
  operations_snapshot:
    current: node tools/vibit inspect operations --json
    future_vocabulary: source_first_operations_snapshot
    status: source_first_repository_inspection_only
  runtime_observability_map:
    current: node tools/vibit inspect pitaya-observability --json
    future_vocabulary: dashboard_read_model
    status: no_dashboard_or_admin_console_behavior
  metrics_tracing_map:
    current: node tools/vibit inspect pitaya-metrics-tracing --json
    future_vocabulary: dashboard_admin_boundary
    status: no_metrics_endpoint_or_tracing_pipeline_behavior
  verification_posture:
    current: tools/vibit repository checks and known warning posture
    future_vocabulary: admin_operation_surface
    status: no_admin_operation_execution_behavior
  redaction_policy:
    current: operations and metrics/tracing redaction requirements
    future_vocabulary: inspector_redaction_posture
    status: no_sensitive_inspector_behavior
  audit_event_posture:
    current: no event/audit tables or admin audit runtime behavior
    future_vocabulary: event_audit_deferral
    status: planning_vocabulary_only
```

## 5. Stop Conditions

在后续 bounded work item 明确授权前，不要添加 runtime endpoint behavior、metrics endpoints、tracing pipelines、observability pipelines、dashboards、admin console behavior、player/session/token inspectors、event/audit tables、admin users/roles/permissions/authentication behavior、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、hosted deployment、SDK publication、release artifacts、distributed runtime behavior 或 direct Nakama/Pitaya API compatibility。

## 6. Verification

Repository check rule: `runtime.pitaya_aligned_dashboard_admin_operations_boundary_gate`。
