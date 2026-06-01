# ADR-0183: Pitaya-Aligned Metrics And Tracing Source-First Map

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-pitaya-aligned-metrics-tracing-source-first-map/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-metrics-tracing-source-first-map.md`

Related artifacts:

- `docs/pitaya-aligned-metrics-tracing-boundary-gate.md`
- `docs/pitaya-aligned-metrics-tracing-boundary-gate.zh-CN.md`
- `decisions/ADR-0182-pitaya-aligned-metrics-tracing-boundary-gate.md`
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

`ADR-0182` defined the gate-only Pitaya-aligned metrics and tracing vocabulary boundary and opened `M-203/W-0275`.

The repository already exposes source-first operations and runtime observability inspection surfaces. The safe next step is a source-first metrics/tracing inspection command that reports the boundary vocabulary, current mappings, source surfaces, redaction posture, and implementation deferrals without adding metrics endpoints, tracing pipelines, dashboards, telemetry dependencies, runtime behavior, hosted surfaces, or direct compatibility.

## Decision

Implement `node tools/vibit inspect pitaya-metrics-tracing --json` as the source-first Pitaya-aligned metrics and tracing map for `M-203/W-0275`.

The command reports:

- ADR-0182 as the source gate and ADR-0183 as the implementation decision.
- `runtime.pitaya_aligned_metrics_tracing_source_first_map` as the check rule.
- Allowed metrics and tracing vocabulary from ADR-0182.
- Related vocabulary from runtime observability, operations inspection, route inventory, verification, redaction, and deferred operations surfaces.
- Current mappings for `node tools/vibit inspect pitaya-observability --json`, `node tools/vibit inspect operations --json`, `/healthz`, `/readyz`, route inventory, verification posture, trace boundary, trace span boundary, trace context boundary, correlation posture, sampling posture, redaction posture, and deferred telemetry surfaces.
- Explicit false deferrals for runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, and direct compatibility.
- `M-204/W-0276 Select next Pitaya-aligned direction after metrics and tracing map` as the next-ready follow-up.

This decision does not add runtime endpoint behavior. It also does not add metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, logging policy changes, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Add metrics endpoints immediately.
- Add tracing pipelines or span behavior immediately.
- Expand `node tools/vibit inspect pitaya-observability --json` instead of adding a metrics/tracing-focused map.
- Keep metrics and tracing mapping only in ADR-0182 without a tool inspection surface.
- Return directly to Nakama product module expansion after W-0274.

## Rationale

Metrics and tracing vocabulary is useful only if agents can inspect the current source-first mapping without treating it as runtime authorization. A focused command gives agents a narrow source-first artifact before any exporter, dashboard, endpoint, sampling, correlation carrier, dependency, or distributed telemetry work is authorized.

## Agent Reasoning Summary

The active work item is a source-first inspection-map implementation. The correct continuation is to add the `tools/vibit` command, repository check rule, change artifacts, ADR, and memory updates while preserving all runtime endpoint, metrics endpoint, tracing pipeline, observability pipeline, dashboard, admin console, protocol, generated output, persistence, dependency, hosted, SDK, release, distributed runtime, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  metrics_tracing_mapping_clarity: high
  implementation_boundedness: high
  source_first_observability_reuse: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now inspect Pitaya-aligned metrics and tracing vocabulary without reading every architecture document.

This decision does not add runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, logging policy changes, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete metrics, tracing, dashboard, admin, or telemetry implementation model;
- the metrics/tracing inspection output creates confusion with public API compatibility or live operational introspection;
- current operations inspection, runtime observability mapping, health/readiness, route inventory, verification, redaction, trace, correlation, or sampling posture changes enough to require remapping;
- future Pitaya-aligned planning needs separate source-first maps for metrics, tracing, correlation, sampling, dashboards, admin console behavior, or distributed node telemetry.

## Follow-Up

- Complete `W-0276`: select the next Pitaya-aligned direction after the metrics and tracing source-first map.
- Keep runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, protocol changes, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.
