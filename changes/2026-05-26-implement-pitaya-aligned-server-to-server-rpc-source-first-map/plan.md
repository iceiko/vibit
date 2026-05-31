# Plan

1. Confirm W-0251 is the active next-ready work item and that ADR-0158 authorizes only a source-first server-to-server RPC map.
2. Run RED checks for the missing `inspect pitaya-rpc` command, missing check rule, and missing change artifacts.
3. Add `node tools/vibit inspect pitaya-rpc --json` to `tools/vibit`.
4. Register `runtime.pitaya_aligned_server_to_server_rpc_source_first_map` in `tools/vibit` and `rules/check-rules.json`.
5. Add ADR-0159, this change directory, and conversation memory.
6. Update architecture manifests, module manifests, guides, README, alpha docs, maturity docs, and roadmap docs to complete W-0251 and open W-0252.
7. Run targeted and full repository verification.
8. Commit and push after verification.

## Boundaries

The implementation must stay source-first. It must not add RPC implementation, remote call behavior, service discovery, frontend/backend server role implementation, distributed runtime behavior, distributed groups, broadcast fanout, cluster-safe session routing, protocol messages or routes, Protobuf source, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, or direct Nakama/Pitaya API compatibility.
