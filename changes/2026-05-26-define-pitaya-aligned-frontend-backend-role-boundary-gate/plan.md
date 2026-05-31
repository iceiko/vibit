# Plan

## Goal

Define W-0248 as a gate-only Pitaya-aligned frontend/backend role boundary and open W-0249 as the source-first role map follow-up.

## Steps

1. Run RED checks for the missing W-0248 rule and change artifacts.
2. Add the English and Simplified Chinese boundary standard.
3. Accept `ADR-0156`.
4. Add conversation memory.
5. Add change artifacts.
6. Register `runtime.pitaya_aligned_frontend_backend_role_boundary_gate`.
7. Add repository check coverage.
8. Update work queue and architecture manifests to complete W-0248 and open W-0249.
9. Update continuation docs and module manifests.
10. Run verification.

## Verification Commands

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_frontend_backend_role_boundary_gate
node tools/vibit check change define-pitaya-aligned-frontend-backend-role-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

## Boundaries

Do not add frontend/backend server role implementation, distributed runtime behavior, server-to-server RPC, remote calls, service discovery, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, or direct Nakama/Pitaya API compatibility.
