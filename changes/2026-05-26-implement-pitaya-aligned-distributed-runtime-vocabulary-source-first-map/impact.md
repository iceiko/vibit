# Impact

## Summary

This change implements the source-first Pitaya-aligned distributed runtime vocabulary map. It moves the repository closer to Pitaya by making ADR-0154 vocabulary and current single-process mappings inspectable through `tools/vibit`.

## Added

- `node tools/vibit inspect pitaya-vocabulary --json`.
- ADR-0155.
- Repository check rule registration.
- Change artifacts and conversation memory.
- Manifest updates completing W-0247 and opening W-0248.

## Not Added

- No runtime behavior.
- No frontend/backend server role implementation.
- No server-to-server RPC, remote calls, service discovery, distributed groups, room broadcast fanout, or cluster-safe session routing.
- No protocol shape.
- No Protobuf source.
- No generated output.
- No repository or PostgreSQL adapter changes.
- No migrations.
- No dependencies.
- No startup wiring.
- No direct Nakama/Pitaya API compatibility.

## Risk

The main risk is source-first vocabulary being mistaken for implementation permission. The inspection output and check rule reduce that risk by emitting explicit false deferral flags and opening only a future frontend/backend role boundary gate.
