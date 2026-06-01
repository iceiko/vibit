# ADR-0180: Pitaya-Aligned Runtime Observability Source-First Map

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-pitaya-aligned-runtime-observability-source-first-map/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-runtime-observability-source-first-map.md`

Related artifacts:

- `docs/pitaya-aligned-runtime-observability-boundary-gate.md`
- `docs/pitaya-aligned-runtime-observability-boundary-gate.zh-CN.md`
- `decisions/ADR-0179-pitaya-aligned-runtime-observability-boundary-gate.md`
- `decisions/ADR-0178-select-next-pitaya-aligned-direction-after-session-binding-kick-disconnect-session-data-map.md`
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

`ADR-0179` defined a gate-only Pitaya-aligned runtime observability vocabulary boundary. It allowed source-first planning vocabulary for runtime observability, operations snapshots, health/readiness signals, version/release posture, configuration posture, route inventory, verification posture, redaction posture, deferred operations surfaces, and node-local runtime surfaces.

The safe next step is a source-first repository inspection command that reports the boundary, current local operations surfaces, source surfaces, deferrals, and redaction posture without adding runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, dependencies, hosted surfaces, or direct compatibility.

## Decision

Implement `node tools/vibit inspect pitaya-observability --json` as the source-first Pitaya-aligned runtime observability map for `M-200/W-0272`.

The command reports:

- ADR-0179 as the source gate and ADR-0180 as the implementation decision.
- `runtime.pitaya_aligned_runtime_observability_source_first_map` as the check rule.
- Allowed runtime observability vocabulary from ADR-0179.
- Related vocabulary from current operations inspection, route, session, acceptor, and distributed-runtime planning surfaces.
- Current mappings for `node tools/vibit inspect operations --json`, `/healthz`, `/readyz`, `/version`, `/configz`, route inventory, verification posture, redaction posture, and deferred operations surfaces.
- Explicit false deferrals for runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, and direct compatibility.
- `M-201/W-0273 Select next Pitaya-aligned direction after runtime observability map` as the next-ready follow-up.

## Alternatives Considered

- Add metrics endpoints, tracing pipelines, observability pipelines, dashboards, or admin console behavior immediately.
- Expand `node tools/vibit inspect operations --json` instead of adding a Pitaya-focused map.
- Keep runtime observability mapping only in ADR-0179 without a tool inspection surface.
- Return directly to Nakama product module expansion after W-0271.

## Rationale

The operations inspection command already summarizes local alpha operational posture. A focused Pitaya-aligned observability map gives agents a narrow source-first place to inspect how those local facts map to future runtime observability vocabulary before any runtime behavior, telemetry dependency, dashboard, admin surface, or distributed runtime work is authorized.

## Agent Reasoning Summary

The active work item is an inspection-map implementation. The correct continuation is to add the `tools/vibit` command, repository check rule, change artifacts, ADR, and memory updates while preserving all runtime endpoint, metrics, tracing, observability pipeline, dashboard, admin console, protocol, generated output, persistence, dependency, hosted, SDK, release, distributed runtime, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  observability_mapping_clarity: high
  implementation_boundedness: high
  source_first_operations_reuse: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now inspect Pitaya-aligned runtime observability vocabulary without reading every architecture document.

This decision does not add runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, logging policy changes, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete metrics, tracing, dashboard, admin, or runtime observability implementation model;
- the runtime observability inspection output creates confusion with public API compatibility;
- current operations inspection, health/readiness/version/config, route inventory, verification, or redaction posture changes enough to require remapping;
- future Pitaya-aligned planning needs separate source-first maps for metrics, tracing, dashboards, admin console behavior, or distributed node telemetry.

## Follow-Up

- Complete `W-0273`: select the next Pitaya-aligned direction after the runtime observability source-first map.
- Keep runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, protocol changes, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.
