# ADR-0186: Pitaya-Aligned Dashboard And Admin Operations Source-First Map

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-pitaya-aligned-dashboard-admin-operations-source-first-map/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-dashboard-admin-operations-source-first-map.md`

Related artifacts:

- `docs/pitaya-aligned-dashboard-admin-operations-boundary-gate.md`
- `docs/pitaya-aligned-dashboard-admin-operations-boundary-gate.zh-CN.md`
- `decisions/ADR-0185-pitaya-aligned-dashboard-admin-operations-boundary-gate.md`
- `decisions/ADR-0184-select-next-pitaya-aligned-direction-after-metrics-tracing-map.md`
- `decisions/ADR-0183-pitaya-aligned-metrics-tracing-source-first-map.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0185` defined the gate-only Pitaya-aligned dashboard and admin operations vocabulary boundary and opened `M-206/W-0278`.

The repository already exposes source-first operations, runtime observability, and metrics/tracing inspection surfaces. The safe next step is a source-first dashboard/admin operations inspection command that reports vocabulary, current mappings, redaction posture, source surfaces, and implementation deferrals without adding dashboards, admin console behavior, sensitive inspectors, audit tables, runtime endpoints, dependencies, hosted surfaces, or direct compatibility.

## Decision

Implement `node tools/vibit inspect pitaya-dashboard-admin --json` as the source-first Pitaya-aligned dashboard and admin operations map for `M-206/W-0278`.

The command reports:

- ADR-0185 as the source gate and ADR-0186 as the implementation decision.
- `runtime.pitaya_aligned_dashboard_admin_operations_source_first_map` as the check rule.
- Allowed dashboard/admin operations vocabulary from ADR-0185.
- Related vocabulary from metrics/tracing, runtime observability, operations inspection, health/readiness, route inventory, verification, redaction, and deferred telemetry surfaces.
- Current mappings for `node tools/vibit inspect operations --json`, `node tools/vibit inspect pitaya-observability --json`, `node tools/vibit inspect pitaya-metrics-tracing --json`, health/readiness/version/config facts, route inventory, verification posture, operator action boundaries, inspector redaction policy, audit/event posture, admin authorization posture, and hosted operations posture.
- Explicit false deferrals for runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, and direct compatibility.
- `M-207/W-0279 Select next Pitaya-aligned direction after dashboard/admin operations map` as the next-ready follow-up.

This decision does not add runtime endpoint behavior. It also does not add metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, logging policy changes, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Add dashboard behavior immediately.
- Add admin console behavior immediately.
- Add player/session/token inspectors immediately.
- Add event/audit tables or admin audit runtime behavior immediately.
- Expand `node tools/vibit inspect pitaya-metrics-tracing --json` instead of adding a dashboard/admin-focused map.
- Return directly to Nakama product module expansion after W-0277.

## Rationale

Dashboard and admin operations are useful only after the repository can distinguish source-first operational metadata from live administrative behavior. A focused command gives agents a narrow, inspectable map before any UI, endpoint, inspector, audit storage, admin authorization, dependency, hosted deployment, or distributed runtime work is authorized.

## Agent Reasoning Summary

The active work item is a source-first inspection-map implementation. The correct continuation is to add the `tools/vibit` command, repository check rule, change artifacts, ADR, and memory updates while preserving all runtime endpoint, metrics endpoint, tracing pipeline, observability pipeline, dashboard, admin console, inspector, event/audit, protocol, generated output, persistence, dependency, hosted, SDK, release, distributed runtime, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  dashboard_admin_mapping_clarity: high
  implementation_boundedness: high
  source_first_operations_reuse: high
  redaction_boundary_value: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now inspect Pitaya-aligned dashboard and admin operations vocabulary without reading every architecture document.

This decision does not add runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, logging policy changes, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete dashboard, admin console, inspector, audit, metrics, tracing, or telemetry implementation model;
- the dashboard/admin inspection output creates confusion with public API compatibility or live operational introspection;
- current operations inspection, runtime observability mapping, metrics/tracing mapping, health/readiness/version/config, route inventory, verification, redaction, audit/event, admin authorization, or hosted operations posture changes enough to require remapping;
- future Pitaya-aligned planning needs separate source-first maps for dashboards, admin console behavior, sensitive inspectors, audit/event storage, admin authorization, or hosted operations.

## Follow-Up

- Complete `W-0279`: select the next Pitaya-aligned direction after the dashboard/admin operations source-first map.
- Keep runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, protocol changes, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.
