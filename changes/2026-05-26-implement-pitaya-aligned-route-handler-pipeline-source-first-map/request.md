# Change Request: Implement Pitaya-Aligned Route Handler Pipeline Source-First Map

Implement `W-0260 Implement Pitaya-aligned route handler pipeline source-first map`.

The maintainer asked to continue toward Pitaya. The narrow request is to make the route handler pipeline vocabulary from `ADR-0167` inspectable through a source-first repository inspection map.

## Scope

- Add `node tools/vibit inspect pitaya-routes --json`.
- Register `runtime.pitaya_aligned_route_handler_pipeline_source_first_map`.
- Record `ADR-0168` and the W-0260 change artifacts.
- Open `W-0261 Select next Pitaya-aligned direction after route handler pipeline map` as next-ready.

## Non-Goals

- No route handler implementation.
- No handler routing behavior.
- No handler pipeline behavior.
- No pipeline middleware behavior.
- No serializer behavior.
- No message forwarding behavior.
- No backend route targeting.
- No protocol messages or routes.
- No Protobuf source.
- No generated output.
- No repository interfaces, PostgreSQL adapters, migrations, or dependencies.
- No service discovery, RPC, remote calls, frontend/backend roles, cluster-safe session routing, or distributed runtime behavior.
- No hosted deployment, SDK publication, release artifact, or direct Nakama/Pitaya API compatibility.
