# Change Request: Implement Pitaya-Aligned Acceptor And Connection Lifecycle Source-First Map

Implement `W-0266 Implement Pitaya-aligned acceptor and connection lifecycle source-first map`.

The maintainer asked to continue toward Pitaya. The narrow request is to make the acceptor and connection lifecycle vocabulary from `ADR-0173` inspectable through a source-first repository inspection map.

## Scope

- Add `node tools/vibit inspect pitaya-acceptor-connection --json`.
- Register `runtime.pitaya_aligned_acceptor_connection_lifecycle_source_first_map`.
- Record `ADR-0174` and the W-0266 change artifacts.
- Open `W-0267 Select next Pitaya-aligned direction after acceptor and connection lifecycle map` as next-ready.

## Non-Goals

- No acceptor behavior.
- No TCP acceptors.
- No WebSocket behavior changes.
- No connection lifecycle behavior changes.
- No session binding behavior.
- No kick/disconnect behavior.
- No concrete socket close behavior changes.
- No serializer behavior.
- No message forwarding behavior.
- No route handler implementation.
- No handler routing behavior.
- No handler pipeline behavior.
- No pipeline middleware behavior.
- No backend route targeting.
- No protocol messages or routes.
- No Protobuf source.
- No generated output.
- No repository interfaces, PostgreSQL adapters, migrations, or dependencies.
- No metrics endpoints, tracing pipelines, service discovery, RPC, remote calls, frontend/backend roles, cluster-safe session routing, or distributed runtime behavior.
- No hosted deployment, SDK publication, release artifact, or direct Nakama/Pitaya API compatibility.
