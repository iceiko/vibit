# Pitaya-Aligned Metrics And Tracing Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-01
Scope: runtime observability source-first map 之后的 metrics/tracing vocabulary gate-only boundary
Depends on: `decisions/ADR-0181-select-next-pitaya-aligned-direction-after-runtime-observability-map.md`, `decisions/ADR-0180-pitaya-aligned-runtime-observability-source-first-map.md`, `docs/pitaya-aligned-runtime-observability-boundary-gate.md`, `docs/minimum-operations-inspection-surface-gate.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0182`

英文文件 `docs/pitaya-aligned-metrics-tracing-boundary-gate.md` 是权威版本。本文是配套简体中文说明。

本文只定义 metrics 与 tracing 词汇边界。它不实现 runtime endpoint behavior、metrics endpoints、tracing pipelines、observability pipelines、dashboards、admin console behavior、player/session/token inspectors、event/audit tables、transport behavior changes、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、hosted deployment、SDK publication、release artifacts、distributed runtime behavior 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

```yaml
pitaya_aligned_metrics_tracing_boundary_gate: defined
completed_work_item: W-0274
decision: ADR-0182
check_rule: runtime.pitaya_aligned_metrics_tracing_boundary_gate
selection_decision: ADR-0181
runtime_observability_source_first_map_decision: ADR-0180
runtime_observability_boundary_gate_decision: ADR-0179
standard: docs/pitaya-aligned-metrics-tracing-boundary-gate.md
translation: docs/pitaya-aligned-metrics-tracing-boundary-gate.zh-CN.md
implementation_scope: gate_only_metrics_tracing_vocabulary
future_implementation_work_item: W-0275
future_implementation_direction: implement_pitaya_aligned_metrics_tracing_source_first_map
allowed_metrics_tracing_vocabulary:
  - metrics_tracing_boundary
  - metric_signal
  - metric_dimension
  - metric_source_surface
  - trace_signal
  - trace_span_boundary
  - trace_context_boundary
  - correlation_id_posture
  - sampling_posture
  - redaction_posture
  - deferred_telemetry_pipeline
  - node_local_telemetry_surface
runtime_endpoint_behavior_added: false
metrics_endpoint_added: false
tracing_pipeline_added: false
observability_pipeline_added: false
dashboard_added: false
admin_console_added: false
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

`ADR-0181` 在 runtime observability source-first map 之后选择了 metrics/tracing boundary gate。风险是 agent 可能把 metrics/tracing 词汇当成添加 endpoints、telemetry dependencies、trace pipelines、dashboards 或 live operational state 的授权。

本 gate 只记录 vocabulary 和 source-first mapping。它保持当前 local alpha runtime behavior 不变，并打开一个后续 source-first map work item。

## 3. Vocabulary

允许的 vocabulary 包括 `metrics_tracing_boundary`、`metric_signal`、`metric_dimension`、`metric_source_surface`、`trace_signal`、`trace_span_boundary`、`trace_context_boundary`、`correlation_id_posture`、`sampling_posture`、`redaction_posture`、`deferred_telemetry_pipeline` 和 `node_local_telemetry_surface`。

这些词汇不得被用来添加 runtime endpoints、metrics endpoints、tracing pipelines、observability pipelines、dashboards、admin console behavior、live inspectors、event/audit tables、dependencies、protocol messages、generated output、persistence、hosted surfaces、SDKs、release artifacts、distributed runtime behavior 或 direct compatibility。

## 4. Current Mapping

当前 mapping 只允许从 source-first facts 映射到 future vocabulary：

- `node tools/vibit inspect pitaya-observability --json` 映射到 `metrics_tracing_boundary`。
- `node tools/vibit inspect operations --json` 映射到 `metric_source_surface`。
- 现有 `/healthz` 和 `/readyz` 摘要映射到 `metric_signal`，但不新增 runtime endpoint behavior。
- source-first route family inventory 映射到 `metric_dimension`，但不改变 protocol route 或 wire shape。
- repository checks 映射到 `metric_signal`，但不新增 test framework 或 dependency。
- trace boundary、correlation posture 和 sampling posture 只保留 planning vocabulary。
- redaction policy 继续要求不扩大 identifier 或 payload log-safety。

## 5. Stop Conditions

添加以下任何内容都必须进入后续 bounded work item：

- runtime endpoint behavior、metrics endpoints、tracing pipelines、observability pipelines、dashboards、admin console behavior、player/session/token inspectors、event/audit tables、dependencies、hosted operations surfaces、SDK publication、release artifacts 或 direct Nakama/Pitaya API compatibility。
- protocol messages or routes、Protobuf sources、generated output、repository interfaces、PostgreSQL adapters、migrations、database inspection 或 persistence changes。
- trace context carriers、correlation id carriers、sampling logic、telemetry exporters、metric storage、alerting rules、dashboard panels 或 production operations integrations。

## 6. Verification

Repository check rule:

```text
runtime.pitaya_aligned_metrics_tracing_boundary_gate
```
