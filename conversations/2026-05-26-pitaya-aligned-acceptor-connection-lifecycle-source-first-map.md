# Conversation: Pitaya-Aligned Acceptor And Connection Lifecycle Source-First Map

Date: 2026-06-01
Status: Accepted
Related change: `changes/2026-05-26-implement-pitaya-aligned-acceptor-connection-lifecycle-source-first-map/`
Related decision: `ADR-0174`

## Context

The maintainer asked to continue pushing toward Pitaya alignment with commit and push discipline. The active continuation queue after `W-0265` was `M-194/W-0266 Implement Pitaya-aligned acceptor and connection lifecycle source-first map`.

`W-0265` had already defined the Pitaya-aligned acceptor and connection lifecycle boundary gate, accepted `ADR-0173`, registered `runtime.pitaya_aligned_acceptor_connection_lifecycle_boundary_gate`, and opened the source-first acceptor and connection lifecycle map as next-ready.

## Maintainer Narrative

Continue toward Pitaya, but keep the work bounded. Acceptor and connection lifecycle vocabulary should become inspectable as future architecture vocabulary, not implemented as runtime behavior.

## Agent Response Summary

The agent treated W-0266 as a source-first inspection-map work item. It added `node tools/vibit inspect pitaya-acceptor-connection --json`, accepted ADR-0174, registered the `runtime.pitaya_aligned_acceptor_connection_lifecycle_source_first_map` check rule, completed W-0266, and opened W-0267 as the next Pitaya-aligned direction selection follow-up.

RED checks confirmed the command, rule, and change artifacts were initially absent:

```text
node tools/vibit inspect pitaya-acceptor-connection --json
# Unknown command.

node tools/vibit inspect rule runtime.pitaya_aligned_acceptor_connection_lifecycle_source_first_map
# Unknown rule_id: runtime.pitaya_aligned_acceptor_connection_lifecycle_source_first_map

node tools/vibit check change implement-pitaya-aligned-acceptor-connection-lifecycle-source-first-map --json
# change directory does not exist
```

## Decisions

- `ADR-0174` implements the Pitaya-aligned acceptor and connection lifecycle source-first map.
- The inspection command is `node tools/vibit inspect pitaya-acceptor-connection --json`.
- The allowed vocabulary is `acceptor_boundary`, `websocket_acceptor`, `connection_id`, `connection_epoch`, `session_binding`, `active_connection_registry`, `close_handoff`, and `presence_lifecycle_handoff`.
- Current vibit behavior remains a single-process WebSocket accept loop, server-observed connection id and epoch metadata, first-message binding route, application-owned active connection registry, transport-to-application close handoff, and server-owned presence lifecycle snapshot.
- W-0267 is the next-ready follow-up for selecting the next Pitaya-aligned direction.

## Artifacts

- `tools/vibit`
- `rules/check-rules.json`
- `decisions/ADR-0174-pitaya-aligned-acceptor-connection-lifecycle-source-first-map.md`
- `changes/2026-05-26-implement-pitaya-aligned-acceptor-connection-lifecycle-source-first-map/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `modules/friends/module.yaml`
- Repository navigation docs and module guide updates for the W-0267 next-ready state.

## Open Questions

No runtime implementation question is answered by this source-first map. A later bounded work item must separately choose any acceptor behavior, TCP acceptor, WebSocket behavior change, session binding behavior, kick/disconnect behavior, connection lifecycle behavior, metrics endpoint, tracing pipeline, protocol carrier, persistence, dependency, service discovery, RPC, remote-call, frontend/backend role, cluster-safe session routing, or distributed runtime implementation.

## Follow-Up

- `M-195/W-0267 Select next Pitaya-aligned direction after acceptor and connection lifecycle map`

## Redaction Notes

The inspection output exposes no raw credentials, raw access tokens, lookup digests, verifier digests, verifier keys, PostgreSQL DSNs, database payloads, local secret file contents, node credentials, transport metadata, connection metadata payloads, route payloads, or local secret values. No ignored credential file contents were read or recorded.
