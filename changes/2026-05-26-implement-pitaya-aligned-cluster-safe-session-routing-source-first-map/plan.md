# Plan

1. Confirm W-0257 next-ready context and RED checks.
2. Add `node tools/vibit inspect pitaya-sessions --json`.
3. Register `runtime.pitaya_aligned_cluster_safe_session_routing_source_first_map`.
4. Add ADR-0165 and conversation memory.
5. Add W-0257 change artifacts.
6. Complete W-0257 in work-item manifests.
7. Open W-0258 as the next Pitaya-aligned direction selection work item.
8. Update repository docs, module manifests, and agent guides.
9. Run targeted verification.
10. Run full repository verification.
11. Commit and push.

## Boundaries

Do not add cluster-safe session routing behavior, session location registries, connection owner node registries, routing epoch behavior, remote handoff, reconnect routing, distributed session routing, distributed runtime behavior, service discovery implementation, RPC implementation, remote calls, frontend/backend role implementation, distributed groups, room broadcast fanout, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, or direct Nakama/Pitaya API compatibility.
