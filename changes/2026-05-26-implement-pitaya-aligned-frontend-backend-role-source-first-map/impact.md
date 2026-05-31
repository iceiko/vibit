# Impact

## Summary

This change implements the source-first Pitaya-aligned frontend/backend role map. It moves the repository closer to Pitaya by making ADR-0156 `frontend_server` and `backend_server` vocabulary and current single-process role mappings inspectable through `tools/vibit`.

## Added

- `node tools/vibit inspect pitaya-roles --json`.
- ADR-0157.
- Repository check rule registration.
- Change artifacts and conversation memory.
- Manifest updates completing W-0249 and opening W-0250.

## Not Added

- No runtime behavior.
- No frontend/backend server role implementation.
- No distributed runtime implementation.
- No server-to-server RPC, remote calls, service discovery, distributed groups, room broadcast fanout, or cluster-safe session routing.
- No protocol shape.
- No Protobuf source.
- No generated output.
- No repository or PostgreSQL adapter changes.
- No migrations.
- No dependencies.
- No startup wiring.
- No hosted deployment, SDK publication, release artifact, or direct Nakama/Pitaya API compatibility.

## Risk

The main risk is source-first role vocabulary being mistaken for topology or process-split permission. The inspection output and check rule reduce that risk by emitting explicit false deferral flags and opening only a future server-to-server RPC boundary gate.
