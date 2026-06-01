# ADR-0166: Select Next Pitaya-Aligned Direction After Cluster-Safe Session Routing Map

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-cluster-safe-session-routing-map/`

Related conversations:

- `conversations/2026-05-26-select-next-pitaya-aligned-direction-after-cluster-safe-session-routing-map.md`

Related artifacts:

- `decisions/ADR-0165-pitaya-aligned-cluster-safe-session-routing-source-first-map.md`
- `decisions/ADR-0164-pitaya-aligned-cluster-safe-session-routing-boundary-gate.md`
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

`ADR-0165` implemented the source-first Pitaya-aligned cluster-safe session routing map and opened `M-186/W-0258` to select the next bounded Pitaya-aligned direction.

The completed Pitaya alignment slices now cover distributed runtime vocabulary, frontend/backend role vocabulary, server-to-server RPC vocabulary, service discovery vocabulary, distributed group and broadcast vocabulary, and cluster-safe session routing vocabulary. The next useful source-first planning area is the route handler model that connects route naming, handler dispatch, handler pipelines, serialization boundaries, and future forwarding decisions.

## Decision

Complete `M-186/W-0258` by selecting `define_pitaya_aligned_route_handler_pipeline_boundary_gate` as the next bounded Pitaya-aligned direction.

Open `M-187/W-0259 Define Pitaya-aligned route handler pipeline boundary gate` as the next-ready work item.

This decision does not implement route handlers, handler routing behavior, pipeline middleware behavior, serializer behavior, message forwarding behavior, runtime behavior, protocol routes, Protobuf sources, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted surfaces, SDKs, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Jump directly to route handler or pipeline implementation.
- Add serializer or message forwarding behavior as part of direction selection.
- Return to service discovery, RPC, distributed group, or cluster-safe session routing implementation immediately.
- Select chat, matchmaking, match runtime, SDK publication, hosted deployment, metrics, or dashboards as the next direction before route handler pipeline vocabulary is bounded.

## Rationale

Pitaya's distributed architecture depends on clear route and handler vocabulary before broader distributed behavior is safe to implement. Selecting a route handler pipeline boundary gate gives agents a narrow next step that can define vocabulary and deferrals before any runtime or protocol behavior changes.

This keeps vibit's agent-native posture intact: choose the next boundary, record the decision, register a repository rule, update continuation memory, and defer implementation until a later explicit work item.

## Agent Reasoning Summary

The active work item is direction selection only. The correct continuation is to record why the next Pitaya-aligned slice should bound route handler, pipeline, serializer, and forwarding vocabulary, then open W-0259 as a gate-only follow-up. Runtime implementation remains out of scope.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  route_handler_pipeline_clarity: high
  implementation_boundedness: high
  dependency_risk: none
  protocol_risk: none_in_this_step
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents now have a precise next-ready Pitaya direction after cluster-safe session routing: define a route handler pipeline boundary gate before any handler, pipeline, serializer, forwarding, distributed routing, protocol, generated output, persistence, or dependency work.

The repository checks can verify W-0258 completion, ADR-0166 presence, W-0259 next-ready state, and explicit no-implementation flags.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a different Pitaya-aligned route or handler model;
- route handler pipeline vocabulary becomes less important than another bounded Pitaya planning surface;
- existing vibit route dispatch, protocol adapter, or module boundary standards change enough to require another direction first;
- the selected direction creates confusion with direct Nakama or Pitaya API compatibility.

## Follow-Up

- Complete `W-0259`: define the Pitaya-aligned route handler pipeline boundary gate.
- Keep route handler implementation, pipeline middleware behavior, serializer behavior, message forwarding behavior, runtime behavior, protocol routes, generated output, persistence, dependencies, hosted surfaces, SDK publication, and direct compatibility behind later bounded work items.
