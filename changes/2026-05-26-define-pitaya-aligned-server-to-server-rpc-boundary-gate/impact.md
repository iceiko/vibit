# Impact

## Summary

This change defines a gate-only Pitaya-aligned server-to-server RPC boundary and opens a source-first RPC map follow-up.

## Added

- `docs/pitaya-aligned-server-to-server-rpc-boundary-gate.md`
- `docs/pitaya-aligned-server-to-server-rpc-boundary-gate.zh-CN.md`
- `decisions/ADR-0158-pitaya-aligned-server-to-server-rpc-boundary-gate.md`
- `runtime.pitaya_aligned_server_to_server_rpc_boundary_gate`
- Change artifacts and conversation memory for W-0250.

## Updated

- Work queue and architecture manifests completing W-0250 and opening W-0251.
- Repository guidance and roadmap references for the new next-ready work item.
- `tools/vibit` checks and expected continuation helpers.
- `rules/check-rules.json`.

## Not Changed

- No runtime behavior.
- No server-to-server RPC implementation.
- No remote call behavior.
- No service discovery.
- No frontend/backend server role implementation.
- No distributed runtime behavior.
- No distributed groups or room broadcast fanout.
- No cluster-safe session routing.
- No runtime endpoint behavior.
- No protocol shape.
- No Protobuf source.
- No generated output.
- No repository interfaces.
- No PostgreSQL adapters.
- No migrations.
- No dependencies.
- No SDKs, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## Follow-Up

W-0251 should implement a source-first Pitaya-aligned server-to-server RPC map while preserving the same implementation deferrals.
