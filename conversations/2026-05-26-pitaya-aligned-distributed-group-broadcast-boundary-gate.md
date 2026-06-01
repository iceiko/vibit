# Conversation: Pitaya-Aligned Distributed Group And Broadcast Boundary Gate

Date: 2026-06-01
Status: Accepted
Related change: `changes/2026-05-26-define-pitaya-aligned-distributed-group-broadcast-boundary-gate/`
Related decision: `ADR-0162`

## Context

The maintainer asked to continue pushing toward Pitaya alignment. The active continuation queue after `W-0253` was `M-182/W-0254 Define Pitaya-aligned distributed group and broadcast boundary gate`.

`W-0253` had already implemented `node tools/vibit inspect pitaya-discovery --json`, accepted `ADR-0161`, registered `runtime.pitaya_aligned_service_discovery_source_first_map`, and opened the distributed group and broadcast boundary gate as next-ready.

## Maintainer Narrative

Continue toward Pitaya, but keep the work bounded. Distributed groups and room broadcast are useful future architecture vocabulary, not permission to add fanout, membership, groups, chat, streams, or cluster routing.

## Agent Response Summary

The agent treated W-0254 as a gate-only work item. It added a distributed group and broadcast boundary standard and translation, accepted ADR-0162, registered the `runtime.pitaya_aligned_distributed_group_broadcast_boundary_gate` check rule, completed W-0254, and opened W-0255 as the source-first distributed group and broadcast map follow-up.

RED checks confirmed the rule and change artifacts were initially absent:

```text
node tools/vibit inspect rule runtime.pitaya_aligned_distributed_group_broadcast_boundary_gate
# Unknown rule_id: runtime.pitaya_aligned_distributed_group_broadcast_boundary_gate

node tools/vibit check change define-pitaya-aligned-distributed-group-broadcast-boundary-gate --json
# change directory does not exist
```

## Decisions

- `ADR-0162` defines the Pitaya-aligned distributed group and broadcast boundary gate.
- The allowed vocabulary is `distributed_group`, `room_broadcast`, `broadcast_target`, `group_membership`, and `broadcast_fanout`.
- Current vibit behavior remains target-scope metadata, application-owned server-push intent, and single-process outbound delivery.
- W-0255 is the next-ready follow-up for a source-first distributed group and broadcast map.

## Artifacts

- `docs/pitaya-aligned-distributed-group-broadcast-boundary-gate.md`
- `docs/pitaya-aligned-distributed-group-broadcast-boundary-gate.zh-CN.md`
- `decisions/ADR-0162-pitaya-aligned-distributed-group-broadcast-boundary-gate.md`
- `changes/2026-05-26-define-pitaya-aligned-distributed-group-broadcast-boundary-gate/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

No runtime implementation question is answered by this gate. A later bounded work item must separately choose any distributed group model, group membership storage, stream subscription semantics, broadcast fanout, delivery guarantee, room routing, service discovery, RPC, remote-call, or cluster-safe session-routing implementation.

## Follow-Up

- `M-183/W-0255 Implement Pitaya-aligned distributed group and broadcast source-first map`

## Redaction Notes

The gate exposes no raw credentials, raw access tokens, lookup digests, verifier digests, verifier keys, PostgreSQL DSNs, database payloads, local secret file contents, node credentials, transport metadata, group membership payloads, or broadcast payloads.
