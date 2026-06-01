# Request

Implement the next Pitaya-aligned service discovery source-first map after the service discovery boundary gate.

## User Need

The maintainer requested continued progress toward Pitaya. The active queue points to `W-0253 Implement Pitaya-aligned service discovery source-first map`.

## Scope

This change adds a source-first repository inspection map for the Pitaya-aligned service discovery vocabulary defined by `ADR-0160`.

The map reports `service_discovery`, `service_registry`, `service_instance`, and `service_selector` vocabulary, maps those terms to current static single-process composition and direct dispatch, and keeps registry, selector, node identity, runtime topology, RPC, remote-call, distributed group, and room-broadcast behavior deferred.

## Non-Goals

- No service discovery implementation.
- No service registry, service selector, node registry, heartbeat, membership, or server identity behavior.
- No server-to-server RPC implementation or remote call behavior.
- No frontend/backend server role implementation.
- No distributed runtime behavior.
- No distributed groups, room broadcast fanout, delivery guarantees, cluster-safe session routing, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf sources, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## Acceptance

- `node tools/vibit inspect pitaya-discovery --json` emits the source-first service discovery map.
- `ADR-0161` accepts the service discovery source-first map.
- `runtime.pitaya_aligned_service_discovery_source_first_map` is registered and checked.
- `W-0253` is completed and `W-0254 Define Pitaya-aligned distributed group and broadcast boundary gate` is opened as next-ready.
