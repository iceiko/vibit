# ADR-0175: Select Next Pitaya-Aligned Direction After Acceptor And Connection Lifecycle Map

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-acceptor-connection-lifecycle-map/`

Related conversations:

- `conversations/2026-05-26-select-next-pitaya-aligned-direction-after-acceptor-connection-lifecycle-map.md`

Related artifacts:

- `decisions/ADR-0174-pitaya-aligned-acceptor-connection-lifecycle-source-first-map.md`
- `decisions/ADR-0173-pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.md`
- `docs/pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0174` implemented `node tools/vibit inspect pitaya-acceptor-connection --json` as the source-first map for Pitaya-aligned acceptor and connection lifecycle vocabulary. That completed the immediate chain for distributed runtime vocabulary, frontend/backend roles, server-to-server RPC, service discovery, distributed groups and broadcast, cluster-safe session routing, route handler pipeline vocabulary, serializer and message forwarding vocabulary, and acceptor/connection lifecycle vocabulary.

The remaining Pitaya reference bullet with high planning value is "Sessions, binding, kick/disconnect, and session data." vibit already has authenticated runtime sessions, first-message connection binding, logout and close handoff, active connection registry state, and presence lifecycle snapshots, but it does not yet have a bounded Pitaya-aligned gate for session binding, kick/disconnect, and session data vocabulary. The repository needs that gate before agents treat those terms as implementation permission.

## Decision

Complete `M-195/W-0267` by selecting `define_pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate` as the next bounded Pitaya-aligned direction.

Open `M-196/W-0268 Define Pitaya-aligned session binding, kick/disconnect, and session data boundary gate` as the next-ready work item.

This decision does not implement session binding behavior, add kick/disconnect behavior, add session data behavior, add session data persistence, implement acceptor behavior, add TCP acceptors, change WebSocket behavior, change connection lifecycle behavior, implement route handlers, implement handler routing behavior, implement handler pipeline behavior, add backend route targeting, change runtime behavior, add protocol routes, add Protobuf sources, change generated output, change repository interfaces, change PostgreSQL adapters, add migrations, add dependencies, add metrics endpoints, add tracing pipelines, add hosted surfaces, publish SDKs, create release artifacts, or add direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement session binding, kick/disconnect, or session data behavior immediately.
- Return to Nakama product module expansion before closing this Pitaya vocabulary gap.
- Add route handler implementation or backend route targeting before session binding vocabulary is bounded.
- Add acceptor, TCP acceptor, WebSocket transport, or connection lifecycle behavior while selecting the next direction.
- Add operations metrics or tracing behavior before defining the session binding vocabulary gate.

## Rationale

The acceptor and connection lifecycle map made the current WebSocket accept loop, connection ids, epochs, first-message binding, active connection registry, close handoff, and presence lifecycle inspectable. The next planning risk is terminology drift around session binding, kick/disconnect, and session data. Those concepts sit between transport connection state and future distributed session routing, so a gate is the smallest next step that can record vocabulary, ownership, deferrals, and check rules without adding behavior.

This keeps the repository aligned with the user's Pitaya direction while preserving the roadmap posture: Nakama remains the primary product reference, and Pitaya remains deferred architecture vocabulary until an explicit later work item authorizes implementation.

## Agent Reasoning Summary

The active work item is direction selection only. The correct continuation is to select the next gate for session binding, kick/disconnect, and session data, then open W-0268 as a gate-only follow-up. Runtime implementation remains out of scope.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  session_binding_vocabulary_clarity: high
  implementation_boundedness: high
  transport_behavior_risk: none_in_this_step
  persistence_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents now have a precise next-ready Pitaya direction after the acceptor and connection lifecycle source-first map: define a session binding, kick/disconnect, and session data boundary gate before any related runtime, transport, protocol, generated output, persistence, dependency, metrics, tracing, or distributed runtime work.

The repository checks can verify W-0267 completion, ADR-0175 presence, W-0268 next-ready state, and explicit no-implementation flags.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a different session binding or session data model;
- the current session validation, connection binding, logout, close handoff, or active connection registry ownership changes enough to require a different mapping;
- another Pitaya-aligned vocabulary gap becomes a higher-risk blocker before W-0268 starts;
- the selected direction creates confusion with direct Nakama or Pitaya API compatibility.

## Follow-Up

- Complete `W-0268`: define the Pitaya-aligned session binding, kick/disconnect, and session data boundary gate.
- Keep session binding behavior, kick/disconnect behavior, session data behavior and persistence, acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, route handler implementation, backend route targeting, service discovery, RPC, remote calls, frontend/backend role behavior, distributed runtime behavior, protocol changes, generated output, persistence, dependencies, metrics endpoints, tracing pipelines, hosted deployment, SDKs, and direct compatibility behind later bounded work items.
