# Request

Implement the next Pitaya-aligned distributed group and broadcast source-first map after the distributed group and broadcast boundary gate.

## User Need

The maintainer requested continued progress toward Pitaya. The active queue points to `W-0255 Implement Pitaya-aligned distributed group and broadcast source-first map`.

## Scope

This change adds a source-first repository inspection map for the Pitaya-aligned distributed group and broadcast vocabulary defined by `ADR-0162`.

The map reports `distributed_group`, `room_broadcast`, `broadcast_target`, `group_membership`, and `broadcast_fanout` vocabulary, maps those terms to current target-scope metadata, server-push intent, and single-process delivery, and keeps group membership registries, room broadcast fanout, delivery guarantees, stream subscriptions, cluster-safe routing, service discovery, RPC, remote-call, and distributed runtime behavior deferred.

## Non-Goals

- No distributed group implementation.
- No group membership registry, subscription state, room state, fanout worker, queue, retry, ordering, or delivery guarantee behavior.
- No stream subscriptions, groups, parties, chat rooms, matchmaking, or match runtime behavior.
- No service discovery implementation, service registry behavior, service selector behavior, node identity, server-to-server RPC implementation, remote call behavior, frontend/backend server role implementation, distributed runtime behavior, cluster-safe session routing behavior, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf sources, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## Acceptance

- `node tools/vibit inspect pitaya-groups --json` emits the source-first distributed group and broadcast map.
- `ADR-0163` accepts the distributed group and broadcast source-first map.
- `runtime.pitaya_aligned_distributed_group_broadcast_source_first_map` is registered and checked.
- `W-0255` is completed and `W-0256 Define Pitaya-aligned cluster-safe session routing boundary gate` is opened as next-ready.
