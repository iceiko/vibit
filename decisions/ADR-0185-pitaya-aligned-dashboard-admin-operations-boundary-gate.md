# ADR-0185: Pitaya-Aligned Dashboard And Admin Operations Boundary Gate

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-pitaya-aligned-dashboard-admin-operations-boundary-gate/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-dashboard-admin-operations-boundary-gate.md`

Related artifacts:

- `docs/pitaya-aligned-dashboard-admin-operations-boundary-gate.md`
- `docs/pitaya-aligned-dashboard-admin-operations-boundary-gate.zh-CN.md`
- `docs/pitaya-aligned-metrics-tracing-boundary-gate.md`
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

`ADR-0184` selected `define_pitaya_aligned_dashboard_admin_operations_boundary_gate` as the next bounded Pitaya-aligned direction after the metrics and tracing source-first map.

The repository already has source-first operations, runtime observability, and metrics/tracing inspection surfaces. Dashboard and admin operations remain deliberately undefined because they can imply UI, runtime endpoints, sensitive inspectors, audit storage, admin authentication, hosted deployments, and third-party dependencies.

## Decision

Define the Pitaya-aligned dashboard and admin operations boundary gate for `M-205/W-0277`.

Register `runtime.pitaya_aligned_dashboard_admin_operations_boundary_gate` as the repository check rule.

The gate allows future planning vocabulary for `dashboard_admin_boundary`, `admin_operation_surface`, `operator_action_boundary`, `dashboard_read_model`, `source_first_operations_snapshot`, `inspector_redaction_posture`, `event_audit_deferral`, `local_operations_diagnostic_surface`, `future_admin_authorization_boundary`, and `hosted_operations_deferral`.

The gate maps current source-first operations, runtime observability, metrics/tracing, health/readiness/version/config, route inventory, verification, redaction, audit/event deferral, and admin authorization deferral posture to that vocabulary.

Complete `M-205/W-0277` and open `M-206/W-0278 Implement Pitaya-aligned dashboard and admin operations source-first map` as next-ready.

The selected follow-up direction is `implement_pitaya_aligned_dashboard_admin_operations_source_first_map`.

This decision does not add dashboard behavior. It also does not add admin console behavior, runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, player/session/token inspectors, event/audit tables, logging policy changes, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement dashboard behavior immediately.
- Implement admin console behavior immediately.
- Add player/session/token inspectors immediately.
- Add event/audit tables immediately.
- Fold dashboard/admin vocabulary into the metrics and tracing gate.
- Return directly to Nakama product module expansion after the metrics and tracing map.

## Rationale

Dashboard and admin operations are operationally useful but high-risk because they easily imply public endpoints, sensitive data exposure, admin authorization behavior, audit/event persistence, UI scope, hosted operations surfaces, and deployment posture.

A gate keeps this step vocabulary-only and source-first. It also lets a later source-first map report committed repository facts before any runtime endpoint, dashboard, admin console, inspector, audit table, dependency, hosted deployment, or distributed runtime work is authorized.

## Agent Reasoning Summary

The active work item is a gate-only boundary. The correct continuation is to add the standard, translation, ADR, change artifacts, rule catalog entry, `tools/vibit` check coverage, and repository memory updates while preserving all runtime endpoint, metrics endpoint, tracing pipeline, observability pipeline, dashboard, admin console, inspector, event/audit table, protocol, generated output, persistence, dependency, hosted, SDK, release, distributed runtime, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  dashboard_admin_boundary_clarity: high
  implementation_boundedness: high
  redaction_boundary_value: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now use dashboard/admin operations vocabulary without treating it as permission to add endpoints, UI, live inspectors, audit/event persistence, dependencies, hosted operations, or runtime behavior.

This decision does not add runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, logging policy changes, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Reversal Conditions

Revisit this decision if:

- a later ADR selects a concrete dashboard, admin console, inspector, audit, metrics, tracing, or telemetry implementation model;
- the boundary creates confusion with public API compatibility or live operational introspection;
- existing operations inspection, runtime observability mapping, metrics/tracing mapping, health/readiness/version/config, route inventory, verification, or redaction posture changes enough to require remapping;
- future Pitaya-aligned planning needs separate gates for dashboards, admin console behavior, sensitive inspectors, audit/event storage, hosted operations, or admin authorization.

## Follow-Up

- Complete `W-0278`: implement a source-first Pitaya-aligned dashboard and admin operations map.
- Keep runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, protocol changes, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.
