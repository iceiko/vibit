# ADR-0179: Pitaya-Aligned Runtime Observability Boundary Gate

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-pitaya-aligned-runtime-observability-boundary-gate/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-runtime-observability-boundary-gate.md`

Related artifacts:

- `docs/pitaya-aligned-runtime-observability-boundary-gate.md`
- `docs/pitaya-aligned-runtime-observability-boundary-gate.zh-CN.md`
- `decisions/ADR-0178-select-next-pitaya-aligned-direction-after-session-binding-kick-disconnect-session-data-map.md`
- `decisions/ADR-0177-pitaya-aligned-session-binding-kick-disconnect-session-data-source-first-map.md`
- `decisions/ADR-0153-minimum-operations-inspection-source-first-surface-implementation.md`
- `decisions/ADR-0152-minimum-operations-inspection-surface-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0178` selected `define_pitaya_aligned_runtime_observability_boundary_gate` as the next bounded Pitaya-aligned direction after the session binding, kick/disconnect, and session data source-first map.

The repository already has a source-first local operations inspection surface from `ADR-0152` and `ADR-0153`, plus minimal local troubleshooting endpoints `/healthz`, `/readyz`, `/version`, and `/configz`. The next safe step is to define runtime observability vocabulary that maps those existing surfaces without adding runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, dependencies, or direct API compatibility.

## Decision

Define the Pitaya-aligned runtime observability boundary gate for `M-199/W-0271`.

Register `runtime.pitaya_aligned_runtime_observability_boundary_gate` as the repository check rule.

The gate allows future planning vocabulary for `runtime_observability_boundary`, `operations_snapshot`, `health_readiness_signal`, `version_release_posture`, `configuration_posture`, `route_inventory_snapshot`, `verification_posture`, `redaction_posture`, `deferred_operations_surface`, and `node_local_runtime_surface`.

The gate maps current source-first operations inspection, health/readiness/version/config endpoint summaries, route inventory, repository verification, redaction posture, and deferred operations surfaces to that vocabulary.

Complete `M-199/W-0271` and open `M-200/W-0272 Implement Pitaya-aligned runtime observability source-first map` as next-ready.

This decision does not add runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, logging policy changes, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, or admin console behavior immediately.
- Expand `node tools/vibit inspect operations --json` directly without a Pitaya-aligned boundary gate.
- Define separate metrics, tracing, dashboard, and admin gates before a shared runtime observability vocabulary exists.
- Return directly to Nakama product module expansion after the session lifecycle map.

## Rationale

Observability is useful for Pitaya-aligned runtime planning, but it easily crosses into runtime behavior, telemetry dependencies, dashboards, and live state inspection. A gate keeps this step source-first and vocabulary-only. It also reuses the existing local operations inspection posture instead of inventing a second operations model.

This preserves vibit's agent-native model: explicit boundaries, inspectable current source surfaces, redaction-first operational posture, and small next-ready implementation slices.

## Agent Reasoning Summary

The active work item is a gate-only boundary. The correct continuation is to add the standard, translation, ADR, change artifacts, rule catalog entry, `tools/vibit` check coverage, and repository memory updates while preserving all runtime endpoint, metrics, tracing, observability pipeline, dashboard, admin console, protocol, generated output, persistence, dependency, hosted, SDK, distributed runtime, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  operations_mapping_clarity: high
  implementation_boundedness: high
  redaction_boundary_value: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now use Pitaya-aligned runtime observability vocabulary without treating it as permission to add metrics, tracing, dashboards, admin operations, or live inspectors.

This decision does not add runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, logging policy changes, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Reversal Conditions

Revisit this decision if:

- a later ADR selects a concrete metrics, tracing, logging, dashboard, admin, or runtime observability implementation model;
- the boundary creates confusion with public API compatibility or live operational introspection;
- existing operations inspection, health/readiness/version/config, route inventory, verification, or redaction posture changes enough to require remapping;
- future Pitaya-aligned planning needs separate gates for metrics, tracing, dashboards, admin console behavior, or distributed node telemetry.

## Follow-Up

- Complete `W-0272`: implement a source-first Pitaya-aligned runtime observability map.
- Keep runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, protocol changes, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.
