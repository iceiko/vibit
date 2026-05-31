# Impact

## Summary

This change implements the source-first Pitaya-aligned server-to-server RPC map. It moves the repository closer to Pitaya by making ADR-0158 `server_to_server_rpc` and `remote_call` vocabulary and current single-process dispatch/module-handler mappings inspectable through `tools/vibit`.

## Added

- `node tools/vibit inspect pitaya-rpc --json`.
- ADR-0159.
- Repository check rule registration.
- Change artifacts and conversation memory.
- Manifest updates completing W-0251 and opening W-0252.

## Not Added

- No runtime behavior.
- No server-to-server RPC implementation.
- No remote call behavior.
- No service discovery.
- No frontend/backend server role implementation.
- No distributed runtime implementation.
- No distributed groups, room broadcast fanout, or cluster-safe session routing.
- No protocol shape.
- No Protobuf source.
- No generated output.
- No repository or PostgreSQL adapter changes.
- No migrations.
- No dependencies.
- No startup wiring.
- No hosted deployment, SDK publication, release artifact, or direct Nakama/Pitaya API compatibility.

## Risk

The main risk is RPC vocabulary being mistaken for permission to add remoting, service discovery, or a distributed runtime. The inspection output and check rule reduce that risk by emitting explicit false deferral flags and opening only a future service discovery boundary gate.
