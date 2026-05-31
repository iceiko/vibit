# Request

Define the next Pitaya-aligned service discovery boundary gate after the server-to-server RPC source-first map.

## User Need

The maintainer requested continued progress toward Pitaya. The active queue points to `W-0252 Define Pitaya-aligned service discovery boundary gate`.

## Scope

This change defines a gate-only service discovery vocabulary boundary. It records how future `service_discovery`, `service_registry`, `service_instance`, and `service_selector` vocabulary maps to the current static single-process runtime.

## Non-Goals

- No service discovery implementation.
- No service registry, service selector, node registry, heartbeat, membership, or server identity behavior.
- No server-to-server RPC implementation or remote call behavior.
- No frontend/backend server role implementation.
- No distributed runtime behavior.
- No distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf sources, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## Acceptance

- `ADR-0160` accepts the service discovery boundary gate.
- `docs/pitaya-aligned-service-discovery-boundary-gate.md` and `.zh-CN.md` define the vocabulary, mapping, ownership, stop conditions, and verification expectations.
- `runtime.pitaya_aligned_service_discovery_boundary_gate` is registered and checked.
- `W-0252` is completed and `W-0253 Implement Pitaya-aligned service discovery source-first map` is opened as next-ready.
