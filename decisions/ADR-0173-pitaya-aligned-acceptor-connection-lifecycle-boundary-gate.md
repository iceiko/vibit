# ADR-0173: Pitaya-Aligned Acceptor And Connection Lifecycle Boundary Gate

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-pitaya-aligned-acceptor-connection-lifecycle-boundary-gate/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.md`

Related artifacts:

- `docs/pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.md`
- `docs/pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.zh-CN.md`
- `decisions/ADR-0172-select-next-pitaya-aligned-direction-after-serializer-message-forwarding-map.md`
- `decisions/ADR-0171-pitaya-aligned-serializer-message-forwarding-source-first-map.md`
- `docs/pitaya-aligned-serializer-message-forwarding-boundary-gate.md`
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

`ADR-0172` selected `define_pitaya_aligned_acceptor_connection_lifecycle_boundary_gate` as the next bounded Pitaya-aligned direction after the serializer and message forwarding source-first map.

Pitaya-style acceptor and connection lifecycle vocabulary is useful architecture pressure because vibit already has concrete single-process facts: a WebSocket accept loop, server-observed connection id and epoch metadata, first-message binding, active connection registry state, close handoff, and server-owned presence lifecycle snapshots.

## Decision

Accept `docs/pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.md` and its Simplified Chinese translation as the gate for Pitaya-aligned acceptor and connection lifecycle vocabulary.

Register `runtime.pitaya_aligned_acceptor_connection_lifecycle_boundary_gate` as the repository check rule.

Complete `M-193/W-0265` and open `M-194/W-0266 Implement Pitaya-aligned acceptor and connection lifecycle source-first map` as next-ready.

This decision does not add acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, session binding behavior, kick/disconnect behavior, serializer behavior, message forwarding behavior, route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, backend route targeting, cluster-safe session routing behavior, distributed session routing, service discovery implementation, RPC, remote calls, frontend/backend role behavior, distributed runtime behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, metrics endpoints, tracing pipelines, hosted deployment, SDK publication, release artifacts, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement acceptor behavior, session binding behavior, or kick/disconnect behavior immediately.
- Add TCP acceptors or WebSocket behavior changes immediately.
- Add a source-first acceptor and connection lifecycle map without first defining a gate.
- Fold acceptor and connection lifecycle vocabulary back into the cluster-safe session routing map.
- Return directly to monitoring, tracing, backend targeting, service discovery, RPC, or distributed runtime implementation.

## Rationale

Acceptor and connection lifecycle vocabulary is high-risk because it can be mistaken for permission to change the transport accept loop, connection lifetime, session binding, close handling, or presence lifecycle. Defining a gate first lets vibit preserve Pitaya-aligned planning vocabulary while keeping the concrete runtime single-process and source-first.

## Agent Reasoning Summary

The active work item is a gate. The correct continuation is to write the standard, ADR, change artifacts, repository checks, and manifest updates. The follow-up should implement a source-first acceptor and connection lifecycle map, not acceptor behavior, transport behavior changes, session binding changes, kick/disconnect behavior, protocol changes, generated output, persistence, dependencies, metrics/tracing, or distributed runtime behavior.

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

- `acceptor_boundary`, `websocket_acceptor`, `connection_id`, `connection_epoch`, `session_binding`, `active_connection_registry`, `close_handoff`, and `presence_lifecycle_handoff` are allowed as future architecture vocabulary only.
- `runtime.pitaya_aligned_acceptor_connection_lifecycle_boundary_gate` becomes the check rule for W-0265.
- `M-193/W-0265` is completed.
- `M-194/W-0266 Implement Pitaya-aligned acceptor and connection lifecycle source-first map` becomes next-ready.

## Reversal Conditions

Revisit this decision if a later architecture ADR selects a concrete acceptor or connection lifecycle model, if the vocabulary creates confusion with public API compatibility, or if WebSocket acceptor, connection registry, first-message binding, close handoff, or presence lifecycle ownership changes enough to require remapping.

## Follow-Up

- Complete `W-0266`: implement a source-first Pitaya-aligned acceptor and connection lifecycle map.
- Keep acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, session binding behavior, kick/disconnect behavior, backend route targeting, service discovery, RPC, remote calls, frontend/backend role behavior, distributed runtime behavior, protocol changes, generated output, persistence, dependencies, metrics endpoints, tracing pipelines, hosted deployment, SDKs, and direct compatibility behind later bounded work items.
