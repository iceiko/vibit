# ADR-0177: Pitaya-Aligned Session Binding Kick Disconnect Session Data Source-First Map

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-pitaya-aligned-session-binding-kick-disconnect-session-data-source-first-map/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-session-binding-kick-disconnect-session-data-source-first-map.md`

Related artifacts:

- `docs/pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate.md`
- `docs/pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate.zh-CN.md`
- `decisions/ADR-0176-pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate.md`
- `decisions/ADR-0175-select-next-pitaya-aligned-direction-after-acceptor-connection-lifecycle-map.md`
- `decisions/ADR-0174-pitaya-aligned-acceptor-connection-lifecycle-source-first-map.md`
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

`ADR-0176` defined a gate-only Pitaya-aligned session binding, kick/disconnect, and session data vocabulary boundary. It allowed `session_binding_boundary`, `connection_bound_session`, `session_data`, `session_data_scope`, `server_initiated_disconnect`, `server_initiated_kick`, `session_unbind`, `session_close_reason`, `connection_session_handoff`, and `presence_session_handoff` as future architecture vocabulary while preserving vibit's current first-message binding, request-level session validation, active connection registry, logout/close handoff, and presence lifecycle behavior.

The safe next step is source-first inspection. Agents need a concrete repository command that summarizes the vocabulary, current single-process mappings, source surfaces, redaction posture, and deferrals before any session binding behavior, kick/disconnect behavior, session data behavior, protocol carrier, persistence, dependency, metrics, tracing, or distributed runtime behavior is authorized.

## Decision

Implement `node tools/vibit inspect pitaya-session-lifecycle --json` as the source-first Pitaya-aligned session binding, kick/disconnect, and session data map for `M-197/W-0269`.

The command reports:

- ADR-0176 as the source gate and ADR-0177 as the implementation decision.
- `runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_source_first_map` as the check rule.
- Allowed session lifecycle vocabulary: `session_binding_boundary`, `connection_bound_session`, `session_data`, `session_data_scope`, `server_initiated_disconnect`, `server_initiated_kick`, `session_unbind`, `session_close_reason`, `connection_session_handoff`, and `presence_session_handoff`.
- Related vocabulary: `first_message_binding`, `runtime_session`, `active_connection_registry`, `logout_service_behavior`, `websocket_close_handoff`, `presence_lifecycle`, `cluster_safe_session_routing`, `acceptor_boundary`, `connection_epoch`, and `route_request`.
- Current single-process mappings for first-message binding, runtime session validation, session metadata, session data scope, active connection registry, logout/close handoff, kick policy, session unbind, close reason, and presence lifecycle.
- Explicit false deferrals for session binding behavior, kick/disconnect behavior, server-initiated disconnect behavior, server-initiated kick behavior, session data behavior, session data persistence, session unbind behavior, close reason behavior changes, connection-session handoff behavior, presence-session handoff behavior, acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, route handler implementation, backend route targeting, cluster-safe session routing, distributed runtime, service discovery, RPC, remote calls, frontend/backend roles, protocol, generated output, persistence, dependencies, metrics, tracing, hosted surfaces, SDKs, and direct compatibility.
- Redaction flags for credential, token, digest, key, DSN, database payload, local secret file, node credential, transport metadata, connection metadata payload contents, session data payload contents, and route payload contents.
- `M-198/W-0270 Select next Pitaya-aligned direction after session binding, kick/disconnect, and session data map` as the next-ready follow-up.

## Alternatives Considered

- Implement session binding behavior, kick/disconnect behavior, session data behavior, or persistence immediately.
- Add WebSocket behavior changes, handshake authentication, reconnect routing, server-directed disconnect behavior, metrics endpoints, or tracing pipelines while adding the map.
- Keep session lifecycle mapping only in ADR-0176 without a tool inspection surface.
- Fold this mapping back into `node tools/vibit inspect pitaya-sessions --json` instead of adding a focused command.

## Rationale

Session binding, kick/disconnect, and session data vocabulary is useful for Pitaya-aligned planning, but it is easy to confuse with permission to add runtime behavior or direct API compatibility. A dedicated source-first inspection command gives agents a precise place to inspect vocabulary, current mappings, deferrals, and redaction posture before any implementation gate.

This preserves vibit's agent-native model: inspectable source surfaces, explicit check rules, bounded vocabulary, and continuation memory before runtime behavior.

## Agent Reasoning Summary

The active work item is an inspection-map implementation. The correct continuation is to add the `tools/vibit` command, repository check rule, change artifacts, ADR, and memory updates while preserving all session binding behavior, kick/disconnect behavior, session data behavior or persistence, acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, route handler implementation, backend targeting, service discovery, RPC, distributed runtime, protocol, persistence, dependency, metrics, tracing, hosted, SDK, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  session_lifecycle_mapping_clarity: high
  implementation_boundedness: high
  inspection_surface_value: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now inspect Pitaya-aligned session binding, kick/disconnect, and session data vocabulary without reading every architecture document.

This decision does not add session binding behavior, kick/disconnect behavior, server-initiated disconnect behavior, server-initiated kick behavior, session data behavior, session data persistence, session unbind behavior, close reason behavior changes, connection-session handoff behavior, presence-session handoff behavior, acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, backend route targeting, cluster-safe session routing behavior, session location registries, connection owner node registries, routing epoch behavior, session route targets, remote connection handoff, reconnect routing, distributed session routing, service discovery implementation, service registries, service selectors, node identity, server-to-server RPC implementation, remote call behavior, frontend/backend server role implementation, distributed runtime behavior, distributed group implementation, group membership registries, room broadcast fanout, delivery guarantees, stream subscriptions, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, or direct Nakama/Pitaya API compatibility.

Future Pitaya-aligned direction planning must start with W-0270 as a selection gate rather than implementation.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete session binding, kick/disconnect, session data, metrics, tracing, or distributed session routing model;
- the session lifecycle inspection output creates confusion with public API compatibility;
- first-message binding, runtime session validation, active connection registry, close handoff, presence lifecycle, outbound delivery, or route ownership changes enough to require remapping;
- future Pitaya-aligned planning needs separate inspection surfaces for session binding, server-directed disconnect, server-directed kick, session data, close reasons, or presence handoff.

## Follow-Up

- Complete `W-0270`: select the next Pitaya-aligned direction after the session binding, kick/disconnect, and session data map.
- Keep session binding behavior, kick/disconnect behavior, session data behavior or persistence, backend route targeting, service discovery implementation, server-to-server RPC, remote calls, frontend/backend role implementation, cluster-safe session routing behavior, distributed runtime behavior, distributed groups, room broadcast fanout, metrics, dashboards, SDKs, hosted deployment, and direct compatibility behind later bounded work items.
