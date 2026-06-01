# Conversation: Pitaya-Aligned Service Discovery Source-First Map

Date: 2026-05-31
Status: Accepted
Related change: `changes/2026-05-26-implement-pitaya-aligned-service-discovery-source-first-map/`
Related decision: `ADR-0161`

## Context

The maintainer asked to continue pushing toward Pitaya alignment. The active continuation queue after `W-0252` was `M-181/W-0253 Implement Pitaya-aligned service discovery source-first map`.

`W-0252` had already defined the Pitaya-aligned service discovery boundary gate, accepted `ADR-0160`, registered `runtime.pitaya_aligned_service_discovery_boundary_gate`, and opened the source-first service discovery map as next-ready.

## Maintainer Narrative

Continue toward Pitaya, but keep the work bounded. Service discovery should become inspectable as future architecture vocabulary, not implemented as runtime behavior.

## Agent Response Summary

The agent treated W-0253 as a source-first inspection-map work item. It added `node tools/vibit inspect pitaya-discovery --json`, accepted ADR-0161, registered the `runtime.pitaya_aligned_service_discovery_source_first_map` check rule, completed W-0253, and opened W-0254 as the distributed group and broadcast boundary gate follow-up.

RED checks confirmed the command, rule, and change artifacts were initially absent:

```text
node tools/vibit inspect pitaya-discovery --json
# Unknown command.

node tools/vibit inspect rule runtime.pitaya_aligned_service_discovery_source_first_map
# Unknown rule_id: runtime.pitaya_aligned_service_discovery_source_first_map

node tools/vibit check change implement-pitaya-aligned-service-discovery-source-first-map --json
# change directory does not exist
```

## Decisions

- `ADR-0161` implements the Pitaya-aligned service discovery source-first map.
- The inspection command is `node tools/vibit inspect pitaya-discovery --json`.
- The allowed vocabulary remains `service_discovery`, `service_registry`, `service_instance`, and `service_selector`.
- Current vibit behavior remains static single-process startup composition and direct in-process dispatch.
- W-0254 is the next-ready follow-up for a distributed group and broadcast boundary gate.

## Artifacts

- `tools/vibit`
- `rules/check-rules.json`
- `decisions/ADR-0161-pitaya-aligned-service-discovery-source-first-map.md`
- `changes/2026-05-26-implement-pitaya-aligned-service-discovery-source-first-map/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`

## Open Questions

No runtime implementation question is answered by this source-first map. A later bounded work item must separately choose any service discovery model, registry storage, selector behavior, node identity, membership, heartbeat, topology, dependency, RPC, remote-call, distributed group, room broadcast, or cluster-safe session-routing implementation.

## Follow-Up

- `M-182/W-0254 Define Pitaya-aligned distributed group and broadcast boundary gate`

## Redaction Notes

The inspection output exposes no raw credentials, raw access tokens, lookup digests, verifier digests, verifier keys, PostgreSQL DSNs, database payloads, local secret file contents, node credentials, or transport metadata.
