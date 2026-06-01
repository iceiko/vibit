# Request

Implement `W-0257 Implement Pitaya-aligned cluster-safe session routing source-first map`.

The maintainer asked to continue toward Pitaya. The narrow request is to make the cluster-safe session routing vocabulary from `ADR-0164` inspectable through a source-first repository inspection map.

No cluster-safe session routing behavior, session location registries, connection owner node registries, routing epoch behavior, remote connection handoff, reconnect routing, distributed session routing, service discovery implementation, RPC implementation, remote calls, frontend/backend role implementation, distributed runtime behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, or direct Nakama/Pitaya API compatibility are requested.

## Acceptance

- Add `node tools/vibit inspect pitaya-sessions --json`.
- Register `runtime.pitaya_aligned_cluster_safe_session_routing_source_first_map`.
- Record ADR-0165, change artifacts, conversation memory, manifest updates, and next-ready W-0258.
- Preserve all implementation deferrals.
