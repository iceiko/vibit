# Request

## Original Request

```text
继续推进，注意提交和推送。提交和推送的key在Git忽略的文件里有。
```

## Clarified Requirement

Complete `M-172/W-0244 Define minimum operations inspection surface gate`.

The friends relationship protocol route family is now locally proven, and `ADR-0151` selected `admin_console_metrics_observability_and_operations` as the next Nakama-first prototype-ready capability family. This slice must define the minimum source-first operations inspection posture before implementation.

## User Story

```text
As a developer evaluating or extending vibit's prototype-ready foundation, I want a minimum source-first operations inspection surface that tells me what local alpha state can be inspected, how redaction is handled, and which operations surfaces remain future work before I add broader social, realtime, competitive, or match runtime behavior.
```

## Selected Posture

```text
source_first_local_operations_inspection
```

## Future Implementation Candidate

```text
tools/vibit inspect operations
```

## Accepted Existing Runtime Surfaces

- `/healthz`
- `/readyz`
- `/version`
- `/configz`

## Accepted Existing Source Surfaces

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `docs/runtime-runbook.md`
- `docs/alpha-developer-flow.md`
- `docs/alpha-acceptance-checklist.md`
- `examples/local-alpha-example-client.sh`
- `examples/local-alpha-request-loop.sh`
- `tools/vibit`

## Non-Goals

- Implement an operations/admin endpoint, admin console, metrics endpoint, observability pipeline, dashboard, hosted operations surface, or live state inspector.
- Add runtime behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, event/audit tables, SDK publication, hosted deployments, release artifacts, public announcements, paid promotion, distributed runtime, or direct Nakama/Pitaya API compatibility.
- Inspect database row payloads, arbitrary JSON object values, raw credentials, raw tokens, digests, verifier keys, DSNs with credentials, transport metadata, or full sensitive identifiers.

## Acceptance Criteria

- `docs/minimum-operations-inspection-surface-gate.md` and `.zh-CN.md` define the gate.
- `ADR-0152` records the decision.
- The gate records minimum inspectable state categories, ownership, redaction posture, future implementation candidates, tests, and stop conditions.
- `M-172/W-0244` is marked completed.
- `M-173/W-0245 Implement minimum operations inspection source-first surface` is opened as next-ready.
- Repository checks include `runtime.minimum_operations_inspection_surface_gate`.
- The gate adds no runtime behavior, endpoint implementation, protocol shape, generated output, persistence, dependency, hosted surface, SDK, distributed runtime, or direct compatibility scope.
