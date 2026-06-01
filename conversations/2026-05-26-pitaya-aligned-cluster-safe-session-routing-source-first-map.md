# Conversation: Pitaya-Aligned Cluster-Safe Session Routing Source-First Map

Date: 2026-06-01
Status: Accepted
Related change: `changes/2026-05-26-implement-pitaya-aligned-cluster-safe-session-routing-source-first-map/`
Related decision: `ADR-0165`

## Context

The maintainer asked to continue pushing toward Pitaya alignment. The active continuation queue after `W-0256` was `M-185/W-0257 Implement Pitaya-aligned cluster-safe session routing source-first map`.

`W-0256` had already defined the Pitaya-aligned cluster-safe session routing boundary gate, accepted `ADR-0164`, registered `runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate`, and opened the source-first session routing map as next-ready.

## Maintainer Narrative

Continue toward Pitaya, but keep the work bounded. Cluster-safe session routing vocabulary should become inspectable as future architecture vocabulary, not implemented as runtime behavior.

## Agent Response Summary

The agent treated W-0257 as a source-first inspection-map work item. It added `node tools/vibit inspect pitaya-sessions --json`, accepted ADR-0165, registered the `runtime.pitaya_aligned_cluster_safe_session_routing_source_first_map` check rule, completed W-0257, and opened W-0258 as the next Pitaya-aligned direction selection follow-up.

RED checks confirmed the command, rule, and change artifacts were initially absent:

```text
node tools/vibit inspect pitaya-sessions --json
# Unknown command.

node tools/vibit inspect rule runtime.pitaya_aligned_cluster_safe_session_routing_source_first_map
# Unknown rule_id: runtime.pitaya_aligned_cluster_safe_session_routing_source_first_map

node tools/vibit check change implement-pitaya-aligned-cluster-safe-session-routing-source-first-map --json
# change directory does not exist
```

## Decisions

- `ADR-0165` implements the Pitaya-aligned cluster-safe session routing source-first map.
- The inspection command is `node tools/vibit inspect pitaya-sessions --json`.
- The allowed vocabulary remains `cluster_safe_session_routing`, `session_location`, `connection_owner_node`, `routing_epoch`, `session_route_target`, `connection_handoff`, and `reconnect_route`.
- Current vibit behavior remains single-process connection binding, active connection registry state, runtime session validation vocabulary, and request-level token identity.
- W-0258 is the next-ready follow-up for selecting the next Pitaya-aligned direction.

## Artifacts

- `tools/vibit`
- `rules/check-rules.json`
- `decisions/ADR-0165-pitaya-aligned-cluster-safe-session-routing-source-first-map.md`
- `changes/2026-05-26-implement-pitaya-aligned-cluster-safe-session-routing-source-first-map/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`

## Open Questions

No runtime implementation question is answered by this source-first map. A later bounded work item must separately choose any session location registry, connection owner node registry, routing epoch behavior, session route target model, connection handoff, reconnect route, distributed session routing, protocol carrier, persistence, dependency, service discovery, RPC, remote-call, or distributed runtime implementation.

## Follow-Up

- `M-186/W-0258 Select next Pitaya-aligned direction after cluster-safe session routing map`

## Redaction Notes

The inspection output exposes no raw credentials, raw access tokens, lookup digests, verifier digests, verifier keys, PostgreSQL DSNs, database payloads, local secret file contents, node credentials, or transport metadata.
