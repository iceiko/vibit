# Change Request: Implement Pitaya-Aligned Serializer And Message Forwarding Source-First Map

Implement `W-0263 Implement Pitaya-aligned serializer and message forwarding source-first map`.

The maintainer asked to continue toward Pitaya. The narrow request is to make the serializer and message forwarding vocabulary from `ADR-0170` inspectable through a source-first repository inspection map.

## Scope

- Add `node tools/vibit inspect pitaya-serializer-forwarding --json`.
- Register `runtime.pitaya_aligned_serializer_message_forwarding_source_first_map`.
- Record `ADR-0171` and the W-0263 change artifacts.
- Open `W-0264 Select next Pitaya-aligned direction after serializer and message forwarding map` as next-ready.

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
