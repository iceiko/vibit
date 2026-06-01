# ADR-0176: Pitaya-Aligned Session Binding, Kick/Disconnect, And Session Data Boundary Gate

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate.md`

Related artifacts:

- `docs/pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate.md`
- `docs/pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate.zh-CN.md`
- `decisions/ADR-0175-select-next-pitaya-aligned-direction-after-acceptor-connection-lifecycle-map.md`
- `decisions/ADR-0174-pitaya-aligned-acceptor-connection-lifecycle-source-first-map.md`
- `docs/pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.md`
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

`ADR-0175` selected `define_pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate` as the next bounded Pitaya-aligned direction after the acceptor and connection lifecycle source-first map.

Pitaya-style session binding, kick/disconnect, and session data vocabulary is useful architecture pressure because vibit already has concrete single-process facts: first-message connection binding, request-level session validation, request identity handoff, active connection registry state, logout service behavior, close handoff, and server-owned presence lifecycle snapshots.

## Decision

Accept `docs/pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate.md` and its Simplified Chinese translation as the gate for Pitaya-aligned session binding, kick/disconnect, and session data vocabulary.

Register `runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate` as the repository check rule.

Complete `M-196/W-0268` and open `M-197/W-0269 Implement Pitaya-aligned session binding, kick/disconnect, and session data source-first map` as next-ready.

This decision does not add session binding behavior, kick/disconnect behavior, session data behavior, session data persistence, acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, backend route targeting, cluster-safe session routing behavior, distributed session routing, service discovery implementation, RPC, remote calls, frontend/backend role behavior, distributed runtime behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, metrics endpoints, tracing pipelines, hosted deployment, SDK publication, release artifacts, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement session binding behavior, kick/disconnect behavior, or session data behavior immediately.
- Add general session data persistence, kick routes, disconnect routes, or reconnect routing immediately.
- Add a source-first session binding, kick/disconnect, and session data map without first defining a gate.
- Fold session binding, kick/disconnect, and session data vocabulary back into the acceptor and connection lifecycle map.
- Return directly to monitoring, tracing, backend targeting, service discovery, RPC, or distributed runtime implementation.

## Rationale

Session binding, kick/disconnect, and session data vocabulary is high-risk because it can be mistaken for permission to alter authentication handoff, connection lifetime, logout behavior, close handling, session storage, or presence lifecycle. Defining a gate first lets vibit preserve Pitaya-aligned planning vocabulary while keeping the concrete runtime single-process and source-first.

## Agent Reasoning Summary

The active work item is a gate. The correct continuation is to write the standard, ADR, change artifacts, repository checks, and manifest updates. The follow-up should implement a source-first session binding, kick/disconnect, and session data map, not session binding behavior, kick/disconnect behavior, session data persistence, acceptor behavior, transport behavior changes, connection lifecycle changes, protocol changes, generated output, persistence, dependencies, metrics/tracing, or distributed runtime behavior.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  session_lifecycle_clarity: high
  implementation_boundedness: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `session_binding_boundary`, `connection_bound_session`, `session_data`, `session_data_scope`, `server_initiated_disconnect`, `server_initiated_kick`, `session_unbind`, `session_close_reason`, `connection_session_handoff`, and `presence_session_handoff` are allowed as future architecture vocabulary only.
- `runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate` becomes the check rule for W-0268.
- `M-196/W-0268` is completed.
- `M-197/W-0269 Implement Pitaya-aligned session binding, kick/disconnect, and session data source-first map` becomes next-ready.

## Reversal Conditions

Revisit this decision if a later architecture ADR selects a concrete session binding, kick/disconnect, session data, reconnect, or distributed session routing model, if the vocabulary creates confusion with public API compatibility, or if first-message binding, request identity handoff, logout behavior, close handoff, connection registry, or presence lifecycle ownership changes enough to require remapping.

## Follow-Up

- Complete `W-0269`: implement a source-first Pitaya-aligned session binding, kick/disconnect, and session data map.
- Keep session binding behavior, kick/disconnect behavior, session data behavior or persistence, acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, backend route targeting, service discovery, RPC, remote calls, frontend/backend role behavior, distributed runtime behavior, protocol changes, generated output, persistence, dependencies, metrics endpoints, tracing pipelines, hosted deployment, SDKs, and direct compatibility behind later bounded work items.
