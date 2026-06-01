# ADR-0172: Select Next Pitaya-Aligned Direction After Serializer And Message Forwarding Map

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-serializer-message-forwarding-map/`

Related conversations:

- `conversations/2026-05-26-select-next-pitaya-aligned-direction-after-serializer-message-forwarding-map.md`

Related artifacts:

- `decisions/ADR-0171-pitaya-aligned-serializer-message-forwarding-source-first-map.md`
- `decisions/ADR-0170-pitaya-aligned-serializer-message-forwarding-boundary-gate.md`
- `docs/pitaya-aligned-serializer-message-forwarding-boundary-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0171` implemented `node tools/vibit inspect pitaya-serializer-forwarding --json` as the source-first map for Pitaya-aligned serializer and message forwarding vocabulary. That completed the immediate route-handler-pipeline chain: route handler and pipeline vocabulary, serializer and forwarding vocabulary, and explicit deferrals for backend targeting, RPC, service discovery, distributed runtime behavior, and direct compatibility.

The remaining high-value Pitaya vocabulary area that is already reflected in vibit's single-process runtime is acceptor and connection lifecycle vocabulary. vibit has a WebSocket acceptor, server-observed connection id and epoch metadata, first-message binding, active connection registry state, close handoff, and presence lifecycle snapshots. Those surfaces should be mapped behind a gate before agents use Pitaya terms such as acceptor, session binding, kick, disconnect, or connection lifecycle as implementation permission.

## Decision

Complete `M-192/W-0264` by selecting `define_pitaya_aligned_acceptor_connection_lifecycle_boundary_gate` as the next bounded Pitaya-aligned direction.

Open `M-193/W-0265 Define Pitaya-aligned acceptor and connection lifecycle boundary gate` as the next-ready work item.

This decision does not implement acceptor behavior, add TCP acceptors, change WebSocket transport behavior, change connection lifecycle behavior, change session binding behavior, add kick/disconnect behavior, implement serializer behavior, implement message forwarding behavior, implement route handlers, implement handler routing behavior, implement handler pipeline behavior, add backend route targeting, change runtime behavior, add protocol routes, add Protobuf sources, change generated output, change repository interfaces, change PostgreSQL adapters, add migrations, add dependencies, add service discovery, add RPC, add remote calls, add frontend/backend roles, add cluster-safe session routing, add metrics endpoints, add tracing pipelines, add hosted surfaces, publish SDKs, create release artifacts, or add direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement serializer behavior or message forwarding immediately.
- Add acceptor, connection lifecycle, session binding, kick, or disconnect behavior immediately.
- Select a monitoring and tracing boundary gate first.
- Return to Nakama product module expansion before closing the remaining Pitaya acceptor and connection lifecycle vocabulary.
- Add backend route targeting, service discovery, RPC, or distributed runtime behavior while selecting the next direction.

## Rationale

Pitaya's architecture vocabulary starts at acceptors and session/connection binding before it reaches route handlers, serializers, forwarding, RPC, discovery, groups, and cluster routing. vibit already has concrete single-process WebSocket and connection lifecycle facts, but those facts are spread across transport, application connection registry, authentication binding, close handoff, and presence code. A boundary gate is the smallest next step that can map those facts to Pitaya vocabulary without authorizing transport or runtime behavior changes.

This keeps the repository aligned with the user's Pitaya direction while preserving the roadmap posture: Nakama remains the primary product reference, and Pitaya remains deferred architecture vocabulary until an explicit later work item authorizes implementation.

## Agent Reasoning Summary

The active work item is direction selection only. The correct continuation is to record why acceptor and connection lifecycle vocabulary should be bounded next, then open W-0265 as a gate-only follow-up. Runtime implementation remains out of scope.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  acceptor_connection_lifecycle_clarity: high
  implementation_boundedness: high
  transport_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents now have a precise next-ready Pitaya direction after the serializer and message forwarding source-first map: define an acceptor and connection lifecycle boundary gate before any acceptor, transport, session binding, kick/disconnect, protocol, generated output, persistence, dependency, metrics, tracing, or distributed runtime work.

The repository checks can verify W-0264 completion, ADR-0172 presence, W-0265 next-ready state, and explicit no-implementation flags.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a different transport or connection lifecycle model;
- the current WebSocket acceptor, connection registry, first-message binding, close handoff, or presence lifecycle ownership changes enough to require a different mapping;
- monitoring/tracing or operations visibility becomes a higher-risk near-term Pitaya vocabulary gap;
- the selected direction creates confusion with direct Nakama or Pitaya API compatibility.

## Follow-Up

- Complete `W-0265`: define the Pitaya-aligned acceptor and connection lifecycle boundary gate.
- Keep acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, session binding behavior, kick/disconnect behavior, serializer behavior, message forwarding behavior, backend route targeting, service discovery, RPC, remote calls, frontend/backend role behavior, distributed runtime behavior, protocol changes, generated output, persistence, dependencies, metrics endpoints, tracing pipelines, hosted deployment, SDKs, and direct compatibility behind later bounded work items.
