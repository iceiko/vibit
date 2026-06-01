# Plan

1. Confirm W-0255 is the active next-ready work item and that ADR-0162 authorizes only a source-first distributed group and broadcast map.
2. Run RED checks for the missing `inspect pitaya-groups` command, missing check rule, and missing change artifacts.
3. Add `node tools/vibit inspect pitaya-groups --json` to `tools/vibit`.
4. Register `runtime.pitaya_aligned_distributed_group_broadcast_source_first_map` in `tools/vibit` and `rules/check-rules.json`.
5. Add ADR-0163, this change directory, and conversation memory.
6. Update architecture manifests, module manifests, guides, README, alpha docs, maturity docs, and roadmap docs to complete W-0255 and open W-0256.
7. Run targeted and full repository verification.
8. Commit and push after verification.

## Boundaries

The implementation must stay source-first. It must not add distributed group implementation, group membership registry behavior, room broadcast fanout, delivery guarantees, stream subscriptions, groups, parties, chat rooms, matchmaking, match runtime, service discovery implementation, service registries, service selectors, node identity, RPC implementation, remote call behavior, frontend/backend server role implementation, distributed runtime behavior, cluster-safe session routing, protocol messages or routes, Protobuf source, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, or direct Nakama/Pitaya API compatibility.
