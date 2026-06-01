# ADR-0182: Pitaya-Aligned Metrics And Tracing Boundary Gate

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-pitaya-aligned-metrics-tracing-boundary-gate/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-metrics-tracing-boundary-gate.md`

Related artifacts:

- `docs/pitaya-aligned-metrics-tracing-boundary-gate.md`
- `docs/pitaya-aligned-metrics-tracing-boundary-gate.zh-CN.md`
- `docs/pitaya-aligned-runtime-observability-boundary-gate.md`
- `decisions/ADR-0181-select-next-pitaya-aligned-direction-after-runtime-observability-map.md`
- `decisions/ADR-0180-pitaya-aligned-runtime-observability-source-first-map.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0181` selected `define_pitaya_aligned_metrics_tracing_boundary_gate` as the next bounded Pitaya-aligned direction after the runtime observability source-first map.

The repository already has source-first operations and runtime observability inspection surfaces, but metrics and tracing remain deliberately undefined. The safe next step is a gate-only vocabulary boundary that defines what future metrics and tracing terms mean without adding endpoints, telemetry dependencies, tracing pipelines, dashboards, hosted surfaces, or runtime behavior.

## Decision

Define the Pitaya-aligned metrics and tracing boundary gate for `M-202/W-0274`.

Register `runtime.pitaya_aligned_metrics_tracing_boundary_gate` as the repository check rule.

The gate allows future planning vocabulary for `metrics_tracing_boundary`, `metric_signal`, `metric_dimension`, `metric_source_surface`, `trace_signal`, `trace_span_boundary`, `trace_context_boundary`, `correlation_id_posture`, `sampling_posture`, `redaction_posture`, `deferred_telemetry_pipeline`, and `node_local_telemetry_surface`.

The gate maps current source-first runtime observability, operations, health/readiness, route inventory, verification, redaction, and deferred telemetry posture to that vocabulary.

Complete `M-202/W-0274` and open `M-203/W-0275 Implement Pitaya-aligned metrics and tracing source-first map` as next-ready.

The selected follow-up direction is `implement_pitaya_aligned_metrics_tracing_source_first_map`.

This decision does not add metrics endpoints. It also does not add runtime endpoint behavior, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, logging policy changes, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement metrics endpoints immediately.
- Implement tracing pipelines immediately.
- Add dashboard or admin console behavior before vocabulary boundaries are ratified.
- Fold metrics and tracing into the existing runtime observability boundary without a dedicated gate.
- Return directly to Nakama product module expansion after the runtime observability map.

## Rationale

Metrics and tracing are operationally useful but high-risk because they easily imply public endpoints, telemetry dependencies, live state exposure, sampling behavior, correlation carriers, dashboards, and hosted operations surfaces.

A gate keeps this step vocabulary-only and source-first. It also lets a later source-first map report committed repository facts before any runtime endpoint, exporter, trace pipeline, dashboard, or dependency work is authorized.

## Agent Reasoning Summary

The active work item is a gate-only boundary. The correct continuation is to add the standard, translation, ADR, change artifacts, rule catalog entry, `tools/vibit` check coverage, and repository memory updates while preserving all runtime endpoint, metrics endpoint, tracing pipeline, observability pipeline, dashboard, admin console, protocol, generated output, persistence, dependency, hosted, SDK, release, distributed runtime, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  metrics_tracing_boundary_clarity: high
  implementation_boundedness: high
  redaction_boundary_value: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now use metrics and tracing vocabulary without treating it as permission to add endpoints, pipelines, dashboards, admin operations, dependencies, or live inspectors.

This decision does not add runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, logging policy changes, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Reversal Conditions

Revisit this decision if:

- a later ADR selects a concrete metrics, tracing, dashboard, admin, or telemetry implementation model;
- the boundary creates confusion with public API compatibility or live operational introspection;
- existing operations inspection, runtime observability mapping, health/readiness/version/config, route inventory, verification, or redaction posture changes enough to require remapping;
- future Pitaya-aligned planning needs separate gates for metrics, tracing, correlation, sampling, dashboards, or distributed node telemetry.

## Follow-Up

- Complete `W-0275`: implement a source-first Pitaya-aligned metrics and tracing map.
- Keep runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, protocol changes, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.
