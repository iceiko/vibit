# Request

Define the next Pitaya-aligned distributed group and broadcast boundary gate after the service discovery source-first map.

## User Need

The maintainer requested continued progress toward Pitaya. The active queue points to `W-0254 Define Pitaya-aligned distributed group and broadcast boundary gate`.

## Scope

This change defines a gate-only distributed group and broadcast vocabulary boundary. It records how future `distributed_group`, `room_broadcast`, `broadcast_target`, `group_membership`, and `broadcast_fanout` vocabulary maps to current target-scope metadata and single-process server-push intent.

## Non-Goals

- No distributed group implementation.
- No room broadcast fanout.
- No delivery guarantees, retries, ordering, durable offsets, queueing, or backpressure behavior.
- No stream subscriptions, group membership registries, groups, parties, chat rooms, matchmaking, or match runtime behavior.
- No service discovery implementation, service registries, service selectors, node identity, server-to-server RPC implementation, remote call behavior, frontend/backend server role implementation, distributed runtime behavior, cluster-safe session routing, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf sources, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## Acceptance

- `ADR-0162` accepts the distributed group and broadcast boundary gate.
- `docs/pitaya-aligned-distributed-group-broadcast-boundary-gate.md` and `.zh-CN.md` define the vocabulary, mapping, ownership, stop conditions, and verification expectations.
- `runtime.pitaya_aligned_distributed_group_broadcast_boundary_gate` is registered and checked.
- `W-0254` is completed and `W-0255 Implement Pitaya-aligned distributed group and broadcast source-first map` is opened as next-ready.
