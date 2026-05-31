# ADR-0152: Minimum Operations Inspection Surface Gate

Status: Accepted
Date: 2026-05-31
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-minimum-operations-inspection-surface-gate/`

Related conversations:

- `conversations/2026-05-26-minimum-operations-inspection-surface-gate.md`

Related artifacts:

- `docs/minimum-operations-inspection-surface-gate.md`
- `docs/minimum-operations-inspection-surface-gate.zh-CN.md`
- `decisions/ADR-0151-select-next-nakama-prototype-ready-capability-after-friends-route-proof.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0151` selected `admin_console_metrics_observability_and_operations` as the next Nakama-first prototype-ready capability family after the protected friends relationship route family was proven through the local alpha request flow.

The local alpha already has small troubleshooting endpoints for liveness, readiness, version, and redacted configuration posture. It also has source-first state in architecture manifests, runbooks, example scripts, repository checks, generated-output standards, and focused tests. What is missing is a bounded operations inspection posture that tells prototype authors which current state can be inspected, what must remain redacted, and which operations surfaces remain future work.

Building an admin console, metrics endpoint, telemetry pipeline, dashboard, or live data inspector in this slice would be premature. It would create sensitive-state, compatibility, and operational semantics before the project has agreed on the minimum source-first surface.

## Decision

Accept `docs/minimum-operations-inspection-surface-gate.md` as the gate for the first minimum operations inspection surface.

The selected first posture is:

```text
source_first_local_operations_inspection
```

The future implementation candidate is:

```text
tools/vibit inspect operations
```

The future implementation should summarize existing committed source, manifests, runbooks, local alpha scripts, route families, status endpoints, generated-output posture, migration posture, and verification posture. It must keep secrets and sensitive identifiers redacted.

Open:

```text
M-173/W-0245 Implement minimum operations inspection source-first surface
```

as the next-ready work item.

This decision completes `M-172/W-0244`. It does not implement operations/admin endpoints, metrics endpoints, observability pipelines, dashboards, runtime behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, event/audit tables, SDK publication, hosted deployments, release artifacts, public announcements, paid promotion, distributed runtime, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement an HTTP operations endpoint immediately.
- Add a Prometheus-style metrics endpoint immediately.
- Add a hosted or local admin dashboard immediately.
- Add live player/session/token/friends/storage inspection endpoints.
- Add event/audit tables before defining inspection semantics.
- Continue directly to groups, parties, chat, leaderboards, matchmaking, or match runtime.
- Reactivate Pitaya distributed architecture review before source-first operations inspection.

## Rationale

A source-first inspection surface fits vibit's current maturity. The project is still pre-alpha, but it already has enough moving parts that prototype authors need a single way to understand local runtime posture without reading every ADR. A `tools/vibit` inspection command can summarize existing state while preserving redaction and avoiding new runtime behavior.

Nakama motivates the capability family because production-useful backend frameworks need operations visibility. vibit adapts that pressure to its current stage: source-first inspection before admin console, metrics, or hosted operations. Pitaya remains deferred because distributed runtime and cluster operations would expand the architecture before the single-process foundation is clear.

## Agent Reasoning Summary

The current `next_ready` work item is a gate. The correct continuation is to define inspectable categories, ownership, redaction, future allowed files, verification expectations, and stop conditions before implementation. The narrow next implementation should be a source-first repository inspection surface rather than a new runtime endpoint.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  operations_visibility_value: high
  redaction_safety: high
  implementation_boundedness: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none_in_this_step
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `docs/minimum-operations-inspection-surface-gate.md` and its Simplified Chinese translation are accepted.
- `runtime.minimum_operations_inspection_surface_gate` becomes the repository check rule for this gate.
- `M-172/W-0244` is completed.
- `M-173/W-0245 Implement minimum operations inspection source-first surface` becomes next-ready.
- The first implementation should prefer `tools/vibit inspect operations` over runtime endpoint expansion.
- Admin console, metrics endpoint, observability pipeline, dashboard, live state inspectors, event/audit tables, hosted operations, SDK publication, distributed runtime, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- prototype feedback shows that source-first operations inspection is insufficient without a runtime endpoint;
- a maintainer explicitly authorizes metrics or admin console work before source-first inspection;
- redaction policy changes to classify previously sensitive identifiers as inspectable in a specific surface;
- the project chooses hosted operations or production-candidate observability as a near-term milestone;
- distributed runtime becomes active before single-process operations posture is implemented.

## Follow-Up

- Complete `W-0245`: implement the minimum operations inspection source-first surface within this gate.
- Keep HTTP operations endpoints, metrics, dashboards, live inspectors, event/audit tables, generated clients, SDKs, hosted deployment, distributed runtime, and direct compatibility behind later bounded work items.
