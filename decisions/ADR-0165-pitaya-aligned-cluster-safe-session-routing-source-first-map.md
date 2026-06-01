# ADR-0165: Pitaya-Aligned Cluster-Safe Session Routing Source-First Map

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-pitaya-aligned-cluster-safe-session-routing-source-first-map/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-cluster-safe-session-routing-source-first-map.md`

Related artifacts:

- `docs/pitaya-aligned-cluster-safe-session-routing-boundary-gate.md`
- `docs/pitaya-aligned-cluster-safe-session-routing-boundary-gate.zh-CN.md`
- `decisions/ADR-0164-pitaya-aligned-cluster-safe-session-routing-boundary-gate.md`
- `decisions/ADR-0163-pitaya-aligned-distributed-group-broadcast-source-first-map.md`
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

`ADR-0164` defined a gate-only Pitaya-aligned cluster-safe session routing boundary. It allowed `cluster_safe_session_routing`, `session_location`, `connection_owner_node`, `routing_epoch`, `session_route_target`, `connection_handoff`, and `reconnect_route` as future architecture vocabulary while preserving vibit's current single-process runtime posture.

The safe next step is source-first inspection. Agents need a concrete repository command that summarizes session-routing vocabulary, current connection/session identity mapping, redaction posture, and deferrals before any session location registry, owner-node registry, routing epoch behavior, remote handoff, reconnect route, service discovery, RPC, distributed runtime, protocol carrier, persistence, or dependency is authorized.

## Decision

Implement `node tools/vibit inspect pitaya-sessions --json` as the source-first Pitaya-aligned cluster-safe session routing map for `M-185/W-0257`.

The command reports:

- ADR-0164 as the source gate and ADR-0165 as the implementation decision.
- `runtime.pitaya_aligned_cluster_safe_session_routing_source_first_map` as the check rule.
- Allowed session-routing vocabulary: `cluster_safe_session_routing`, `session_location`, `connection_owner_node`, `routing_epoch`, `session_route_target`, `connection_handoff`, and `reconnect_route`.
- Related vocabulary: `connection_id`, `connection_epoch`, `first_message_connection_binding`, `active_connection_registry`, `runtime_session`, `bound_connection_identity`, `request_token_identity`, `session_validated_identity`, `single_process_connection_binding`, `frontend_server`, `backend_server`, `service_discovery`, `server_to_server_rpc`, `remote_call`, `distributed_group`, and `room_broadcast`.
- Current single-process mappings for connection id, connection epoch, first-message binding, active connection registry, runtime session, bound connection identity, request token identity, session validated identity, connection handoff, reconnect route, and distributed session routing planning.
- Explicit false deferrals for cluster-safe session routing behavior, session location registries, connection owner node registries, routing epoch behavior, session route targets, remote connection handoff, reconnect routing, distributed session routing, service discovery implementation, RPC implementation, remote calls, frontend/backend role implementation, distributed runtime implementation, distributed groups, room broadcast fanout, protocol, generated output, persistence, dependency, hosted, SDK, and direct compatibility surfaces.
- Redaction flags for credential, token, digest, key, DSN, database payload, local secret file, node credential, and transport metadata contents.
- `M-186/W-0258 Select next Pitaya-aligned direction after cluster-safe session routing map` as the next-ready follow-up.

## Alternatives Considered

- Implement cluster-safe session routing immediately.
- Add session location registries, owner-node registries, routing epochs, remote handoff, or reconnect routing while adding the map.
- Keep session-routing mapping only in ADR-0164 without a tool inspection surface.
- Fold session-routing vocabulary into `node tools/vibit inspect pitaya-groups --json` instead of adding a focused command.

## Rationale

Cluster-safe session routing vocabulary is useful for Pitaya-aligned planning, but it is easy to confuse with permission to add distributed session lookup, connection-owner placement, reconnect routing, remote handoff, service discovery, RPC, transport carriers, dependencies, or multi-node behavior. A dedicated source-first inspection command gives agents a precise place to inspect vocabulary, current mappings, deferrals, and redaction posture before any implementation gate.

This preserves vibit's agent-native model: inspectable source surfaces, explicit check rules, bounded vocabulary, and continuation memory before runtime behavior.

## Agent Reasoning Summary

The active work item is an inspection-map implementation. The correct continuation is to add the `tools/vibit` command, repository check rule, change artifacts, ADR, and memory updates while preserving all session-routing implementation, registry, handoff, reconnect, service discovery, RPC, distributed runtime, protocol, persistence, dependency, hosted, SDK, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  session_routing_mapping_clarity: high
  implementation_boundedness: high
  inspection_surface_value: high
  distributed_runtime_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now inspect Pitaya-aligned cluster-safe session routing vocabulary without reading every architecture document.

This decision does not add cluster-safe session routing behavior, session location registries, connection owner node registries, routing epoch behavior, session route targets, remote connection handoff, reconnect routing, distributed session routing, service discovery implementation, service registries, service selectors, node identity, server-to-server RPC implementation, remote call behavior, frontend/backend server role implementation, distributed runtime behavior, distributed group implementation, group membership registries, room broadcast fanout, delivery guarantees, stream subscriptions, runtime endpoint behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, or direct Nakama/Pitaya API compatibility.

Future Pitaya-aligned direction planning must start with W-0258 as a selection gate rather than implementation.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete distributed session routing or distributed runtime model;
- the session-routing inspection output creates confusion with public API compatibility;
- connection binding, active connection registry, runtime session validation, or route-policy ownership changes enough to require remapping;
- future Pitaya-aligned planning needs separate inspection surfaces for handoff, reconnect, session locations, or node ownership.

## Follow-Up

- Complete `W-0258`: select the next Pitaya-aligned direction after the cluster-safe session routing map.
- Keep cluster-safe session routing behavior, session location registries, connection owner node registries, routing epoch behavior, remote handoff, reconnect routing, service discovery implementation, server-to-server RPC, remote calls, frontend/backend role implementation, distributed runtime behavior, distributed groups, room broadcast fanout, metrics, dashboards, SDKs, hosted deployment, and direct compatibility behind later bounded work items.
