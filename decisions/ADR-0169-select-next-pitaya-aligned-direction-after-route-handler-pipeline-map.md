# ADR-0169: Select Next Pitaya-Aligned Direction After Route Handler Pipeline Map

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-route-handler-pipeline-map/`

Related conversations:

- `conversations/2026-05-26-select-next-pitaya-aligned-direction-after-route-handler-pipeline-map.md`

Related artifacts:

- `decisions/ADR-0168-pitaya-aligned-route-handler-pipeline-source-first-map.md`
- `decisions/ADR-0167-pitaya-aligned-route-handler-pipeline-boundary-gate.md`
- `docs/pitaya-aligned-route-handler-pipeline-boundary-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `modules/friends/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0168` implemented `node tools/vibit inspect pitaya-routes --json` as the source-first map for Pitaya-aligned route handler pipeline vocabulary. That map exposes `serializer_boundary` and `message_forwarding` as future vocabulary, while explicitly preserving current vibit protocol bridge ownership and single-process outbound message behavior.

After route handler pipeline vocabulary is inspectable, the next useful bounded Pitaya-aligned direction is a gate for serializer and message forwarding vocabulary. Those terms are close to implementation risk: serializer behavior can imply pluggable encoding or protocol shape changes, and message forwarding can imply cross-node delivery, backend targeting, service discovery, RPC, or distributed runtime behavior.

## Decision

Complete `M-189/W-0261` by selecting `define_pitaya_aligned_serializer_message_forwarding_boundary_gate` as the next bounded Pitaya-aligned direction.

Open `M-190/W-0262 Define Pitaya-aligned serializer and message forwarding boundary gate` as the next-ready work item.

This decision does not implement serializer behavior, message forwarding behavior, route handlers, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, backend route targeting, runtime behavior, protocol routes, Protobuf sources, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, service discovery, RPC, remote calls, frontend/backend roles, cluster-safe session routing, hosted surfaces, SDKs, release artifacts, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement serializer behavior or message forwarding immediately.
- Add backend route targeting, service discovery, RPC, or frontend/backend routing while selecting the next direction.
- Return to route handler implementation before bounding serializer and forwarding terminology.
- Select matchmaking, match runtime, chat, metrics, hosted deployment, or SDK publication before completing the immediate Pitaya route pipeline vocabulary chain.

## Rationale

Pitaya's route handler architecture connects handler routing, serialization, message forwarding, server roles, and remote dispatch. vibit's current source-first map already names serializer and forwarding surfaces, but it does not yet define a gate that keeps those terms separated from runtime behavior. Selecting a serializer and message forwarding boundary gate gives agents a narrow next step that can clarify vocabulary, current vibit mappings, and deferrals before implementation.

This preserves vibit's agent-native posture: choose the next boundary, record candidate directions and non-goals, register a repository rule, update continuation memory, and defer implementation until a later explicit work item.

## Agent Reasoning Summary

The active work item is direction selection only. The correct continuation is to record why serializer and message forwarding vocabulary should be bounded next, then open W-0262 as a gate-only follow-up. Runtime implementation remains out of scope.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  serializer_forwarding_clarity: high
  implementation_boundedness: high
  distributed_runtime_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents now have a precise next-ready Pitaya direction after the route handler pipeline source-first map: define a serializer and message forwarding boundary gate before any serializer, forwarding, backend targeting, route handler, pipeline, distributed runtime, protocol, generated output, persistence, or dependency work.

The repository checks can verify W-0261 completion, ADR-0169 presence, W-0262 next-ready state, and explicit no-implementation flags.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a different serializer or forwarding model;
- route handler pipeline vocabulary changes enough to require another map first;
- current protocol bridge, outbound delivery, or target-scope ownership changes enough to require remapping;
- the selected direction creates confusion with direct Nakama or Pitaya API compatibility.

## Follow-Up

- Complete `W-0262`: define the Pitaya-aligned serializer and message forwarding boundary gate.
- Keep serializer behavior, message forwarding behavior, route handler implementation, handler routing, handler pipelines, backend route targeting, protocol routes, generated output, persistence, dependencies, hosted surfaces, SDK publication, and direct compatibility behind later bounded work items.
