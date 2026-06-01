# Conversation: Pitaya-Aligned Distributed Group And Broadcast Source-First Map

Date: 2026-06-01
Status: Accepted
Related change: `changes/2026-05-26-implement-pitaya-aligned-distributed-group-broadcast-source-first-map/`
Related decision: `ADR-0163`

## Context

The maintainer asked to continue pushing toward Pitaya alignment. The active continuation queue after `W-0254` was `M-183/W-0255 Implement Pitaya-aligned distributed group and broadcast source-first map`.

`W-0254` had already defined the Pitaya-aligned distributed group and broadcast boundary gate, accepted `ADR-0162`, registered `runtime.pitaya_aligned_distributed_group_broadcast_boundary_gate`, and opened the source-first group/broadcast map as next-ready.

## Maintainer Narrative

Continue toward Pitaya, but keep the work bounded. Distributed group and broadcast vocabulary should become inspectable as future architecture vocabulary, not implemented as runtime behavior.

## Agent Response Summary

The agent treated W-0255 as a source-first inspection-map work item. It added `node tools/vibit inspect pitaya-groups --json`, accepted ADR-0163, registered the `runtime.pitaya_aligned_distributed_group_broadcast_source_first_map` check rule, completed W-0255, and opened W-0256 as the cluster-safe session routing boundary gate follow-up.

RED checks confirmed the command, rule, and change artifacts were initially absent:

```text
node tools/vibit inspect pitaya-groups --json
# Unknown command.

node tools/vibit inspect rule runtime.pitaya_aligned_distributed_group_broadcast_source_first_map
# Unknown rule_id: runtime.pitaya_aligned_distributed_group_broadcast_source_first_map

node tools/vibit check change implement-pitaya-aligned-distributed-group-broadcast-source-first-map --json
# change directory does not exist
```

## Decisions

- `ADR-0163` implements the Pitaya-aligned distributed group and broadcast source-first map.
- The inspection command is `node tools/vibit inspect pitaya-groups --json`.
- The allowed vocabulary remains `distributed_group`, `room_broadcast`, `broadcast_target`, `group_membership`, and `broadcast_fanout`.
- Current vibit behavior remains target-scope metadata, application-owned server-push intent, and single-process WebSocket delivery.
- W-0256 is the next-ready follow-up for a cluster-safe session routing boundary gate.

## Artifacts

- `tools/vibit`
- `rules/check-rules.json`
- `decisions/ADR-0163-pitaya-aligned-distributed-group-broadcast-source-first-map.md`
- `changes/2026-05-26-implement-pitaya-aligned-distributed-group-broadcast-source-first-map/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`

## Open Questions

No runtime implementation question is answered by this source-first map. A later bounded work item must separately choose any distributed group model, group membership registry, room broadcast fanout, delivery guarantee, stream subscription, cluster-safe routing, protocol carrier, persistence, dependency, RPC, remote-call, or distributed runtime implementation.

## Follow-Up

- `M-184/W-0256 Define Pitaya-aligned cluster-safe session routing boundary gate`

## Redaction Notes

The inspection output exposes no raw credentials, raw access tokens, lookup digests, verifier digests, verifier keys, PostgreSQL DSNs, database payloads, local secret file contents, node credentials, or transport metadata.
