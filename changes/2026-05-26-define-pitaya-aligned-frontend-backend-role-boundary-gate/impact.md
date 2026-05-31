# Impact

## Summary

This change defines a gate-only Pitaya-aligned frontend/backend role boundary and opens a source-first role map follow-up.

## Added

- `docs/pitaya-aligned-frontend-backend-role-boundary-gate.md`
- `docs/pitaya-aligned-frontend-backend-role-boundary-gate.zh-CN.md`
- `decisions/ADR-0156-pitaya-aligned-frontend-backend-role-boundary-gate.md`
- `runtime.pitaya_aligned_frontend_backend_role_boundary_gate`
- Change artifacts and conversation memory for W-0248.

## Updated

- Work queue and architecture manifests completing W-0248 and opening W-0249.
- Repository guidance and roadmap references for the new next-ready work item.
- `tools/vibit` checks and expected continuation helpers.
- `rules/check-rules.json`.

## Not Changed

- No runtime behavior.
- No frontend/backend server role implementation.
- No distributed runtime behavior.
- No server-to-server RPC or remote call behavior.
- No service discovery.
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

W-0249 should implement a source-first Pitaya-aligned frontend/backend role map while preserving the same implementation deferrals.
