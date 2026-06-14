# ADR-0228: Pitaya-Aligned Heartbeat And Liveness Source-First Map

Status: Accepted
Date: 2026-06-14
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-06-14-implement-pitaya-aligned-heartbeat-liveness-source-first-map/`

Related conversations:

- `conversations/2026-06-14-implement-pitaya-aligned-heartbeat-liveness-source-first-map.md`

Related artifacts:

- `docs/pitaya-aligned-heartbeat-liveness-boundary-gate.md`
- `docs/pitaya-aligned-heartbeat-liveness-boundary-gate.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

The maintainer confirmed a twenty-step continuation to move vibit closer to Pitaya-class distributed operations planning after `ADR-0220` completed the Pitaya service dispatch source-first maps. The sequence is source-first and must not implement distributed runtime behavior, RPC behavior, frontend/backend behavior, service discovery behavior, protocol changes, migrations, dependencies, or direct Nakama/Pitaya API compatibility.

## Decision

Implement `node tools/vibit inspect pitaya-heartbeat-liveness --json` and register `runtime.pitaya_aligned_heartbeat_liveness_source_first_map` as the source-first map check rule for `M-248/W-0320`.

This decision records `W-0320` as completed, with `ADR-0228` as the decision and `runtime.pitaya_aligned_heartbeat_liveness_source_first_map` as the check rule.

This decision does not add heartbeat and liveness behavior, service registry behavior, service selector behavior, heartbeat/liveness behavior, route targeting behavior, remote timeout/retry behavior, distributed session ownership behavior, distributed presence fanout behavior, cross-node error mapping behavior, cluster observability behavior, runtime behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement heartbeat and liveness behavior immediately.
- Add concrete distributed runtime APIs or direct Pitaya-compatible names.
- Return to broad product module expansion without source-first Pitaya distributed operations planning.
- Defer the confirmed twenty-step sequence until after another product slice.

## Rationale

Heartbeat And Liveness is useful Pitaya-class vocabulary, but it can easily imply distributed runtime or operations behavior. A narrow source-first inspection map lets agents reason about future architecture without changing the running server.

## Agent Reasoning Summary

The correct continuation is to add bounded docs, ADRs, change artifacts, manifests, check rules, and inspection output while preserving all runtime, protocol, generated, persistence, dependency, hosted, SDK, release, distributed runtime, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  source_first_clarity: high
  implementation_boundedness: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now inspect or reason about Pitaya-aligned heartbeat and liveness vocabulary without treating that vocabulary as live runtime behavior.

The repository remains single-process and source-first for this concern.

## Reversal Conditions

Revisit this decision if a later ADR selects a concrete heartbeat and liveness, distributed runtime, service discovery, RPC, frontend/backend, event bus, protocol, persistence, hosted, or direct compatibility implementation model.

## Follow-Up

- Continue through the confirmed twenty-step Pitaya distributed operations sequence with `W-0321`.
- Keep runtime behavior, protocol changes, generated output, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.

## Verification Markers

```yaml
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
dependency_added: false
distributed_runtime_implementation_added: false
direct_nakama_pitaya_api_compatibility_added: false
```
