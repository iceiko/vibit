# ADR-0184: Select Next Pitaya-Aligned Direction After Metrics And Tracing Map

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-metrics-tracing-map/`

Related conversations:

- `conversations/2026-05-26-select-next-pitaya-aligned-direction-after-metrics-tracing-map.md`

Related artifacts:

- `decisions/ADR-0183-pitaya-aligned-metrics-tracing-source-first-map.md`
- `decisions/ADR-0182-pitaya-aligned-metrics-tracing-boundary-gate.md`
- `decisions/ADR-0181-select-next-pitaya-aligned-direction-after-runtime-observability-map.md`
- `docs/pitaya-aligned-metrics-tracing-boundary-gate.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0183` completed the source-first Pitaya-aligned metrics and tracing map through `node tools/vibit inspect pitaya-metrics-tracing --json`.

The map consolidated metrics and tracing vocabulary, runtime observability mapping, operations posture, trace/correlation/sampling deferrals, redaction posture, source surfaces, and implementation deferrals. The next bounded continuation should keep operations capability planning moving without turning the source-first map into live dashboard, admin, metrics endpoint, tracing pipeline, or runtime behavior authorization.

The safe next step is a boundary gate for dashboard and admin operations vocabulary. That gate can define how future dashboard and admin operations surfaces relate to existing source-first observability, metrics, tracing, health, readiness, version, config, route inventory, verification, and redaction posture.

## Decision

Select `define_pitaya_aligned_dashboard_admin_operations_boundary_gate` as the next bounded Pitaya-aligned direction after the metrics and tracing source-first map.

Register `runtime.next_pitaya_aligned_direction_after_metrics_tracing_map` as the repository check rule.

Complete `M-204/W-0276` and open `M-205/W-0277 Define Pitaya-aligned dashboard and admin operations boundary gate` as next-ready.

This decision does not add dashboard behavior. It also does not add admin console behavior, runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, player/session/token inspectors, event/audit tables, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement dashboard or admin console behavior immediately.
- Add player, session, token, event, or audit inspectors immediately.
- Add runtime operations endpoints immediately.
- Extend metrics and tracing into live endpoints or pipelines immediately.
- Return directly to Nakama product module expansion after the metrics and tracing map.
- Start distributed runtime implementation after the metrics and tracing vocabulary sequence.

## Rationale

Dashboard and admin operations are high-leverage for Pitaya-class operational alignment, but they carry broad risk: runtime endpoints, sensitive data inspection, authentication and authorization surfaces, event/audit storage, observability backends, UI scope, deployment posture, and direct compatibility confusion.

A gate keeps the next step bounded. It can define vocabulary, source boundaries, redaction expectations, sensitive-inspection deferrals, and future implementation gates before any dashboard, admin console, endpoint, protocol, persistence, dependency, hosted, SDK, release, or distributed behavior is authorized.

## Agent Reasoning Summary

The active work item is a selection-only continuation. The correct continuation is to record the selected follow-up direction, update repository memory, register the check rule, and open W-0277. Dashboard and admin operations behavior remain deferred until a later gate and implementation slice explicitly authorize them.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  operations_surface_continuity: high
  implementation_boundedness: high
  redaction_boundary_value: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

W-0277 becomes the next-ready work item and must define only a dashboard and admin operations boundary gate.

This selection does not authorize runtime endpoints, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, protocol changes, generated output, persistence changes, dependencies, hosted surfaces, SDK publication, distributed runtime behavior, or direct compatibility.

## Reversal Conditions

Revisit this decision if a later ADR chooses to return to Nakama product module expansion before any further Pitaya-aligned operations vocabulary, or if dashboard and admin operations must be split into separate gates before any shared boundary is useful.

## Follow-Up

- Complete `W-0277`: define the Pitaya-aligned dashboard and admin operations boundary gate.
- Keep runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, protocol changes, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.
