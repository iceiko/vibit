# ADR-0164: Pitaya-Aligned Cluster-Safe Session Routing Boundary Gate

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-pitaya-aligned-cluster-safe-session-routing-boundary-gate/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-cluster-safe-session-routing-boundary-gate.md`

Related artifacts:

- `docs/pitaya-aligned-cluster-safe-session-routing-boundary-gate.md`
- `docs/pitaya-aligned-cluster-safe-session-routing-boundary-gate.zh-CN.md`
- `docs/pitaya-aligned-distributed-group-broadcast-boundary-gate.md`
- `decisions/ADR-0163-pitaya-aligned-distributed-group-broadcast-source-first-map.md`
- `decisions/ADR-0162-pitaya-aligned-distributed-group-broadcast-boundary-gate.md`
- `docs/first-message-connection-binding-gate.md`
- `docs/active-connection-registry-gate.md`
- `docs/runtime-session-validation-gate.md`
- `docs/logout-revocation-active-connection-gate.md`
- `docs/bound-identity-route-policy-gate.md`
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

`ADR-0163` implemented `node tools/vibit inspect pitaya-groups --json` as a source-first map for Pitaya-aligned distributed group and broadcast vocabulary. It opened `M-184/W-0256` to define a cluster-safe session routing boundary gate before session routing vocabulary turns into cross-node lookup, remote handoff, reconnect routing, service discovery, RPC, or distributed runtime behavior.

Pitaya-style cluster routing is useful architecture vocabulary for future distributed runtime planning, but vibit's current runtime remains a single-process server with server-observed `connection_id`, `connection_epoch`, first-message connection binding, an application-owned active connection registry posture, request-level token identity, and runtime session validation vocabulary.

## Decision

Accept `docs/pitaya-aligned-cluster-safe-session-routing-boundary-gate.md` and its Simplified Chinese translation as the gate for Pitaya-aligned cluster-safe session routing vocabulary.

Register `runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate` as the repository check rule.

The gate defines:

- allowed session routing vocabulary: `cluster_safe_session_routing`, `session_location`, `connection_owner_node`, `routing_epoch`, `session_route_target`, `connection_handoff`, and `reconnect_route`;
- related vocabulary: `connection_id`, `connection_epoch`, `first_message_connection_binding`, `active_connection_registry`, `runtime_session`, `bound_connection_identity`, `request_token_identity`, `session_validated_identity`, `single_process_connection_binding`, `frontend_server`, `backend_server`, `service_discovery`, `server_to_server_rpc`, `remote_call`, `distributed_group`, and `room_broadcast`;
- current single-process mapping for connection id, connection epoch, first-message binding, active connection registry, runtime session, bound identity, request token identity, session validated identity, connection handoff, and reconnect route concepts;
- ownership for cluster-safe session routing vocabulary and a future source-first map;
- stop conditions for any implementation behavior.

Complete `M-184/W-0256` and open `M-185/W-0257 Implement Pitaya-aligned cluster-safe session routing source-first map` as next-ready.

This decision does not add cluster-safe session routing behavior, session location registries, connection owner node registries, routing epoch behavior, session route targets, remote connection handoff, reconnect routing, distributed session routing, distributed runtime behavior, distributed groups, group membership registries, room broadcast fanout, delivery guarantees, stream subscriptions, service discovery implementation, service registries, service selectors, node identity, server-to-server RPC implementation, remote call behavior, frontend/backend server role implementation, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement cluster-safe session routing immediately.
- Add a source-first session routing map without first defining a gate.
- Keep session routing vocabulary only inside the general Pitaya vocabulary or group/broadcast map.
- Defer session routing vocabulary entirely until distributed runtime implementation.

## Rationale

Session routing vocabulary is high-risk because it can be mistaken for permission to add session location registries, owner node registries, routing epochs, reconnect routing, remote handoff, service discovery, RPC, transport carriers, dependencies, or distributed runtime behavior. Defining a gate first lets vibit preserve Pitaya-aligned planning vocabulary while keeping the concrete runtime single-process.

This preserves vibit's agent-native model: vocabulary, ownership, deferrals, stop conditions, checks, and memory before implementation.

## Agent Reasoning Summary

The active work item is a gate. The correct continuation is to write the standard, ADR, change artifacts, repository checks, and manifest updates. The follow-up should implement a source-first cluster-safe session routing map, not session routing behavior, registries, handoff, reconnect routing, service discovery, RPC, remote calls, or distributed runtime behavior.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  session_routing_boundary_clarity: high
  implementation_boundedness: high
  distributed_runtime_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `cluster_safe_session_routing`, `session_location`, `connection_owner_node`, `routing_epoch`, `session_route_target`, `connection_handoff`, and `reconnect_route` are allowed as future architecture vocabulary only.
- `runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate` becomes the check rule for W-0256.
- `M-184/W-0256` is completed.
- `M-185/W-0257 Implement Pitaya-aligned cluster-safe session routing source-first map` becomes next-ready.
- Session location registries, connection owner node registries, routing epochs, remote handoff, reconnect routing, service discovery implementation, RPC, remote calls, frontend/backend role implementation, distributed runtime behavior, distributed groups, room broadcast fanout, delivery guarantees, protocol changes, generated output, persistence, dependencies, hosted deployment, SDKs, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete cluster-safe session routing or distributed runtime model;
- the session routing vocabulary creates confusion with public API compatibility;
- connection binding, active connection registry, runtime session validation, or route-policy ownership changes enough to require remapping;
- prototype feedback shows a stronger near-term need for a different product capability family.

## Follow-Up

- Complete `W-0257`: implement a source-first Pitaya-aligned cluster-safe session routing map.
- Keep cluster-safe session routing behavior, session location registries, connection owner node registries, routing epoch behavior, remote handoff, reconnect routing, service discovery implementation, server-to-server RPC, remote calls, frontend/backend role implementation, distributed runtime behavior, distributed groups, room broadcast fanout, metrics, dashboards, SDKs, hosted deployment, and direct compatibility behind later bounded work items.
