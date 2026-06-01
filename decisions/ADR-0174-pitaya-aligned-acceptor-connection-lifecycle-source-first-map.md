# ADR-0174: Pitaya-Aligned Acceptor And Connection Lifecycle Source-First Map

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-pitaya-aligned-acceptor-connection-lifecycle-source-first-map/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-acceptor-connection-lifecycle-source-first-map.md`

Related artifacts:

- `docs/pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.md`
- `docs/pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.zh-CN.md`
- `decisions/ADR-0173-pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.md`
- `decisions/ADR-0172-select-next-pitaya-aligned-direction-after-serializer-message-forwarding-map.md`
- `decisions/ADR-0171-pitaya-aligned-serializer-message-forwarding-source-first-map.md`
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

`ADR-0173` defined a gate-only Pitaya-aligned acceptor and connection lifecycle vocabulary boundary. It allowed `acceptor_boundary`, `websocket_acceptor`, `connection_id`, `connection_epoch`, `session_binding`, `active_connection_registry`, `close_handoff`, and `presence_lifecycle_handoff` as future architecture vocabulary while preserving vibit's current single-process WebSocket, first-message binding, active connection registry, close handoff, and presence lifecycle behavior.

The safe next step is source-first inspection. Agents need a concrete repository command that summarizes acceptor and connection lifecycle vocabulary, current single-process mappings, source surfaces, redaction posture, and deferrals before any acceptor behavior, TCP acceptor, WebSocket behavior change, connection lifecycle behavior, session binding behavior, kick/disconnect behavior, protocol carrier, persistence, dependency, metrics, tracing, or distributed runtime behavior is authorized.

## Decision

Implement `node tools/vibit inspect pitaya-acceptor-connection --json` as the source-first Pitaya-aligned acceptor and connection lifecycle map for `M-194/W-0266`.

The command reports:

- ADR-0173 as the source gate and ADR-0174 as the implementation decision.
- `runtime.pitaya_aligned_acceptor_connection_lifecycle_source_first_map` as the check rule.
- Allowed acceptor and connection lifecycle vocabulary: `acceptor_boundary`, `websocket_acceptor`, `connection_id`, `connection_epoch`, `session_binding`, `active_connection_registry`, `close_handoff`, and `presence_lifecycle_handoff`.
- Related vocabulary: `first_message_binding`, `runtime_session`, `route_request`, `server_push_delivery`, `cluster_safe_session_routing`, `message_forwarding`, `route_handler`, `serializer_boundary`, `frontend_server`, `backend_server`, `server_to_server_rpc`, `remote_call`, and `service_discovery`.
- Current single-process mappings for WebSocket acceptor, server-observed connection id, server-observed connection epoch, first-message binding, active connection registry, close handoff, and presence lifecycle.
- Explicit false deferrals for acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, session binding behavior, kick/disconnect behavior, concrete socket close behavior changes, serializer behavior, message forwarding behavior, route handler implementation, backend route targeting, cluster-safe session routing, distributed runtime, service discovery, RPC, remote calls, frontend/backend roles, protocol, generated output, persistence, dependencies, metrics, tracing, hosted surfaces, SDKs, and direct compatibility.
- Redaction flags for credential, token, digest, key, DSN, database payload, local secret file, node credential, transport metadata, connection metadata payload contents, and route payload contents.
- `M-195/W-0267 Select next Pitaya-aligned direction after acceptor and connection lifecycle map` as the next-ready follow-up.

## Alternatives Considered

- Implement acceptor behavior, session binding behavior, kick/disconnect behavior, or lifecycle behavior immediately.
- Add TCP acceptors, WebSocket behavior changes, handshake authentication, reconnect routing, remote disconnect handoff, metrics endpoints, or tracing pipelines while adding the map.
- Keep acceptor and connection lifecycle mapping only in ADR-0173 without a tool inspection surface.
- Fold acceptor and connection lifecycle mapping back into `node tools/vibit inspect pitaya-sessions --json` instead of adding a focused command.

## Rationale

Acceptor and connection lifecycle vocabulary is useful for Pitaya-aligned planning, but it is easy to confuse with permission to change transport, session binding, close behavior, presence lifecycle, metrics, tracing, or multi-node routing. A dedicated source-first inspection command gives agents a precise place to inspect vocabulary, current mappings, deferrals, and redaction posture before any implementation gate.

This preserves vibit's agent-native model: inspectable source surfaces, explicit check rules, bounded vocabulary, and continuation memory before runtime behavior.

## Agent Reasoning Summary

The active work item is an inspection-map implementation. The correct continuation is to add the `tools/vibit` command, repository check rule, change artifacts, ADR, and memory updates while preserving all acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, session binding behavior, kick/disconnect behavior, serializer behavior, message forwarding behavior, route handler implementation, backend targeting, service discovery, RPC, distributed runtime, protocol, persistence, dependency, metrics, tracing, hosted, SDK, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  acceptor_connection_lifecycle_mapping_clarity: high
  implementation_boundedness: high
  inspection_surface_value: high
  transport_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now inspect Pitaya-aligned acceptor and connection lifecycle vocabulary without reading every architecture document.

This decision does not add acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, session binding behavior, kick/disconnect behavior, concrete socket close behavior changes, serializer behavior, message forwarding behavior, route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, backend route targeting, cluster-safe session routing behavior, session location registries, connection owner node registries, routing epoch behavior, session route targets, remote connection handoff, reconnect routing, distributed session routing, service discovery implementation, service registries, service selectors, node identity, server-to-server RPC implementation, remote call behavior, frontend/backend server role implementation, distributed runtime behavior, distributed group implementation, group membership registries, room broadcast fanout, delivery guarantees, stream subscriptions, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, or direct Nakama/Pitaya API compatibility.

Future Pitaya-aligned direction planning must start with W-0267 as a selection gate rather than implementation.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete acceptor, connection lifecycle, metrics, tracing, or distributed session routing model;
- the acceptor and connection lifecycle inspection output creates confusion with public API compatibility;
- WebSocket acceptor, connection registry, first-message binding, close handoff, presence lifecycle, outbound delivery, or route ownership changes enough to require remapping;
- future Pitaya-aligned planning needs separate inspection surfaces for acceptors, connection identity, session binding, close handoff, or presence lifecycle.

## Follow-Up

- Complete `W-0267`: select the next Pitaya-aligned direction after the acceptor and connection lifecycle map.
- Keep acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, session binding behavior, kick/disconnect behavior, backend route targeting, service discovery implementation, server-to-server RPC, remote calls, frontend/backend role implementation, cluster-safe session routing behavior, distributed runtime behavior, distributed groups, room broadcast fanout, metrics, dashboards, SDKs, hosted deployment, and direct compatibility behind later bounded work items.
