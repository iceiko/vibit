# Conversation: Pitaya-Aligned Cluster-Safe Session Routing Boundary Gate

Date: 2026-06-01
Status: Accepted
Related change: `changes/2026-05-26-define-pitaya-aligned-cluster-safe-session-routing-boundary-gate/`
Related decision: `ADR-0164`

## Context

The maintainer asked to continue pushing toward Pitaya alignment. The active continuation queue after `W-0255` was `M-184/W-0256 Define Pitaya-aligned cluster-safe session routing boundary gate`.

`W-0255` had already implemented `node tools/vibit inspect pitaya-groups --json`, accepted `ADR-0163`, registered `runtime.pitaya_aligned_distributed_group_broadcast_source_first_map`, and opened the cluster-safe session routing boundary gate as next-ready.

## Maintainer Narrative

Continue toward Pitaya, but keep the work bounded. Cluster-safe session routing is useful future architecture vocabulary, not permission to add cross-node session lookup, connection owner registries, remote handoff, reconnect routing, service discovery, RPC, or distributed runtime behavior.

## Agent Response Summary

The agent treated W-0256 as a gate-only work item. It added a cluster-safe session routing boundary standard and translation, accepted ADR-0164, registered the `runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate` check rule, completed W-0256, and opened W-0257 as the source-first cluster-safe session routing map follow-up.

RED checks confirmed the rule and change artifacts were initially absent:

```text
node tools/vibit inspect rule runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate
# Unknown rule_id: runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate

node tools/vibit check change define-pitaya-aligned-cluster-safe-session-routing-boundary-gate --json
# change directory does not exist
```

## Decisions

- `ADR-0164` defines the Pitaya-aligned cluster-safe session routing boundary gate.
- The allowed vocabulary is `cluster_safe_session_routing`, `session_location`, `connection_owner_node`, `routing_epoch`, `session_route_target`, `connection_handoff`, and `reconnect_route`.
- Current vibit behavior remains single-process connection binding, active connection registry vocabulary, request-level token identity, and runtime session validation vocabulary.
- W-0257 is the next-ready follow-up for a source-first cluster-safe session routing map.

## Artifacts

- `docs/pitaya-aligned-cluster-safe-session-routing-boundary-gate.md`
- `docs/pitaya-aligned-cluster-safe-session-routing-boundary-gate.zh-CN.md`
- `decisions/ADR-0164-pitaya-aligned-cluster-safe-session-routing-boundary-gate.md`
- `changes/2026-05-26-define-pitaya-aligned-cluster-safe-session-routing-boundary-gate/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

No runtime implementation question is answered by this gate. A later bounded work item must separately choose any session location model, connection owner node registry, routing epoch semantics, reconnect routing, remote handoff, service discovery, RPC, remote-call, distributed runtime, distributed group, or broadcast implementation.

## Follow-Up

- `M-185/W-0257 Implement Pitaya-aligned cluster-safe session routing source-first map`

## Redaction Notes

The gate exposes no raw credentials, raw access tokens, lookup digests, verifier digests, verifier keys, PostgreSQL DSNs, database payloads, local secret file contents, node credentials, transport metadata, connection payloads, session payloads, or routing metadata contents.
