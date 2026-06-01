# Pitaya-Aligned Metrics And Tracing Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-01
Scope: Gate-only boundary for using Pitaya-aligned metrics and tracing vocabulary after the runtime observability source-first map
Depends on: `decisions/ADR-0181-select-next-pitaya-aligned-direction-after-runtime-observability-map.md`, `decisions/ADR-0180-pitaya-aligned-runtime-observability-source-first-map.md`, `docs/pitaya-aligned-runtime-observability-boundary-gate.md`, `docs/minimum-operations-inspection-surface-gate.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0182`

The paired Simplified Chinese translation is `docs/pitaya-aligned-metrics-tracing-boundary-gate.zh-CN.md`. The English file is authoritative.

This document defines a metrics and tracing vocabulary gate only. It does not implement runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, transport behavior changes, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The Pitaya-aligned metrics and tracing boundary gate record is:

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
primary_product_reference: Nakama
pitaya_reference_status: metrics_tracing_vocabulary_boundary_defined_for_future_architecture_planning
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
current_source_first_metrics_tracing_mapping:
  runtime_observability_map:
    current: node_tools_vibit_inspect_pitaya_observability
    future_vocabulary: metrics_tracing_boundary
    implementation_status: source_first_repository_inspection_only
  operations_snapshot:
    current: node_tools_vibit_inspect_operations
    future_vocabulary: metric_source_surface
    implementation_status: no_metrics_endpoint_behavior
  health_readiness_signal:
    current: existing_local_healthz_and_readyz_endpoint_summary
    future_vocabulary: metric_signal
    implementation_status: no_new_runtime_endpoint_behavior
  route_inventory_snapshot:
    current: source_first_protocol_route_family_inventory
    future_vocabulary: metric_dimension
    implementation_status: no_protocol_route_or_wire_shape_change
  verification_posture:
    current: tools_vibit_check_work_runtime_memory_schemas_all
    future_vocabulary: metric_signal
    implementation_status: no_new_test_framework_or_dependency
  trace_boundary:
    current: no_trace_pipeline_or_span_runtime_behavior
    future_vocabulary: trace_signal
    implementation_status: planning_vocabulary_only
  correlation_posture:
    current: no_cross_request_trace_context_contract
    future_vocabulary: correlation_id_posture
    implementation_status: no_protocol_or_header_carrier_change
  sampling_posture:
    current: no_sampling_runtime_or_dependency
    future_vocabulary: sampling_posture
    implementation_status: planning_vocabulary_only
  redaction_policy:
    current: existing_operations_redaction_requirements
    future_vocabulary: redaction_posture
    implementation_status: no_identifier_or_payload_log_safety_expansion
  deferred_telemetry:
    current: metrics_endpoints_tracing_pipelines_dashboards_admin_and_hosted_surfaces_deferred
    future_vocabulary: deferred_telemetry_pipeline
    implementation_status: planning_vocabulary_only
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

`ADR-0181` selected metrics and tracing as the next Pitaya-aligned direction after the runtime observability source-first map.

The risk is that agents may treat metrics and tracing vocabulary as permission to add endpoints, telemetry dependencies, trace pipelines, dashboards, or live operational state. This gate records vocabulary and source-first mapping only. It keeps the current local alpha runtime behavior unchanged and prepares a narrow source-first map follow-up.

## 3. Vocabulary

Allowed metrics and tracing vocabulary:

- `metrics_tracing_boundary`: future planning vocabulary for the boundary that owns metrics and tracing semantics.
- `metric_signal`: future planning vocabulary for a named measurement or status signal.
- `metric_dimension`: future planning vocabulary for bounded classification fields on future metric signals.
- `metric_source_surface`: future planning vocabulary for source-first facts that could later inform metrics.
- `trace_signal`: future planning vocabulary for traceable activity facts.
- `trace_span_boundary`: future planning vocabulary for where trace spans may begin and end.
- `trace_context_boundary`: future planning vocabulary for how trace context might cross runtime boundaries later.
- `correlation_id_posture`: future planning vocabulary for request or connection correlation without selecting a carrier.
- `sampling_posture`: future planning vocabulary for later trace or metric sampling semantics.
- `redaction_posture`: future planning vocabulary for what metrics and tracing surfaces may expose.
- `deferred_telemetry_pipeline`: future planning vocabulary for metrics/tracing pipeline work that is still deferred.
- `node_local_telemetry_surface`: future planning vocabulary for local single-process telemetry facts. It is not distributed node telemetry.

Forbidden vocabulary use:

- Do not introduce concrete public API, package, route, method, wire, handler, metrics, trace, dashboard, or admin compatibility names from Pitaya or Nakama.
- Do not use metrics or tracing vocabulary as permission to add runtime endpoints, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, live player/session/token inspectors, event/audit tables, dependencies, protocol messages, generated output, persistence, hosted surfaces, SDKs, release artifacts, or distributed runtime behavior.
- Do not classify raw tokens, credentials, lookup digests, verifier digests, verifier keys, DSNs, headers, cookies, query strings, subprotocol values, remote addresses, database payloads, local secret file contents, route payloads, session data payloads, or concrete transport metadata as log-safe in this gate.

## 4. Current Mapping

```yaml
current_source_first_metrics_tracing_mapping:
  runtime_observability_map:
    current: node tools/vibit inspect pitaya-observability --json
    future_vocabulary: metrics_tracing_boundary
    status: source_first_repository_inspection_only
  operations_snapshot:
    current: node tools/vibit inspect operations --json
    future_vocabulary: metric_source_surface
    status: no_metrics_endpoint_behavior
  health_readiness_signal:
    current: existing /healthz and /readyz local troubleshooting endpoint summaries
    future_vocabulary: metric_signal
    status: no_new_runtime_endpoint_behavior
  route_inventory_snapshot:
    current: source-first route family inventory from committed source
    future_vocabulary: metric_dimension
    status: no_protocol_route_or_wire_shape_change
  verification_posture:
    current: tools/vibit repository checks and known warning posture
    future_vocabulary: metric_signal
    status: no_new_test_framework_or_dependency
  trace_boundary:
    current: no trace pipeline or span runtime behavior
    future_vocabulary: trace_signal
    status: planning_vocabulary_only
  correlation_posture:
    current: no cross-request trace context contract
    future_vocabulary: correlation_id_posture
    status: no_protocol_or_header_carrier_change
  sampling_posture:
    current: no sampling runtime or dependency
    future_vocabulary: sampling_posture
    status: planning_vocabulary_only
  redaction_policy:
    current: minimum operations inspection redaction requirements
    future_vocabulary: redaction_posture
    status: no_identifier_or_payload_log_safety_expansion
  deferred_telemetry:
    current: metrics endpoints, tracing pipelines, dashboards, admin, hosted, SDK, and distributed telemetry surfaces deferred
    future_vocabulary: deferred_telemetry_pipeline
    status: planning_vocabulary_only
```

## 5. Ownership

Metrics and tracing vocabulary ownership:

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-metrics-tracing-boundary-gate.md
  - .arch/reference.yaml
  - .arch/runtime.yaml
source_first_map_candidate_owner:
  - tools/vibit
runtime_endpoint_owner: unchanged
metrics_endpoint_owner: deferred
tracing_pipeline_owner: deferred
dashboard_owner: deferred
admin_console_owner: deferred
protocol_owner: unchanged
persistence_owner: unchanged
dependency_owner: unchanged
```

Rules:

- Documentation and manifests may define metrics and tracing vocabulary and current source-first mapping.
- `tools/vibit` may later emit a source-first metrics/tracing map if a follow-up implementation work item authorizes it.
- Existing `/healthz`, `/readyz`, `/version`, and `/configz` endpoints remain unchanged by this gate.
- Runtime behavior, transport behavior, protocol payloads, repository interfaces, migrations, generated output, startup wiring, metrics endpoints, tracing pipelines, dashboards, admin console behavior, event/audit tables, dependencies, hosted surfaces, SDKs, release artifacts, and distributed runtime behavior remain unchanged by this gate.

## 6. Nakama And Pitaya Mapping

Nakama remains the primary product reference for admin console, metrics, observability, and operations capability pressure. Pitaya remains an architecture vocabulary reference for runtime, node, session, route, RPC, service discovery, and group concerns.

This gate maps those references into vibit-owned vocabulary only. It does not create direct compatibility, public API parity, metrics parity, tracing parity, dashboard parity, or runtime behavior.

## 7. Stop Conditions

Stop and require a later bounded work item before adding:

- Runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, logging policy changes, dependency changes, hosted operations surfaces, SDK publication, release artifacts, or direct Nakama/Pitaya API compatibility.
- Protocol messages or routes, Protobuf sources, generated output, repository interfaces, PostgreSQL adapters, migrations, database inspection, arbitrary state dumps, or persistence changes.
- Trace context carriers, correlation id carriers, sampling logic, telemetry exporters, metric storage, alerting rules, dashboard panels, or production operations integrations.
- Session binding behavior, kick/disconnect behavior, session data behavior or persistence, acceptor behavior, WebSocket behavior changes, connection lifecycle behavior changes, route handler implementation, backend route targeting, cluster-safe session routing, service discovery, RPC, frontend/backend roles, or distributed runtime behavior.

## 8. Verification

The repository check rule for this boundary is:

```text
runtime.pitaya_aligned_metrics_tracing_boundary_gate
```

Verification commands:

```sh
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_metrics_tracing_boundary_gate
node tools/vibit check change define-pitaya-aligned-metrics-tracing-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```
