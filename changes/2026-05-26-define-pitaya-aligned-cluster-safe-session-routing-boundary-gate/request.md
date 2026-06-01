# Request

Define the next Pitaya-aligned cluster-safe session routing boundary gate after the distributed group and broadcast source-first map.

## User Need

The maintainer requested continued progress toward Pitaya. The active queue points to `W-0256 Define Pitaya-aligned cluster-safe session routing boundary gate`.

## Scope

This change defines a gate-only cluster-safe session routing vocabulary boundary. It records how future `cluster_safe_session_routing`, `session_location`, `connection_owner_node`, `routing_epoch`, `session_route_target`, `connection_handoff`, and `reconnect_route` vocabulary maps to current single-process connection binding, active connection registry, runtime session validation, and request identity surfaces.

## Non-Goals

- No cluster-safe session routing behavior.
- No session location registry or connection owner node registry.
- No routing epoch behavior, session route target behavior, remote connection handoff, reconnect routing, or distributed session routing.
- No distributed runtime behavior.
- No distributed group implementation, group membership registries, room broadcast fanout, delivery guarantees, or stream subscriptions.
- No service discovery implementation, service registries, service selectors, node identity, server-to-server RPC implementation, remote call behavior, frontend/backend server role implementation, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf sources, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## Acceptance

- `ADR-0164` accepts the cluster-safe session routing boundary gate.
- `docs/pitaya-aligned-cluster-safe-session-routing-boundary-gate.md` and `.zh-CN.md` define the vocabulary, mapping, ownership, stop conditions, and verification expectations.
- `runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate` is registered and checked.
- `W-0256` is completed and `W-0257 Implement Pitaya-aligned cluster-safe session routing source-first map` is opened as next-ready.
