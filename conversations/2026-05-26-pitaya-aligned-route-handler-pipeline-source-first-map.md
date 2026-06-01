# Conversation: Pitaya-Aligned Route Handler Pipeline Source-First Map

Date: 2026-06-01
Status: Accepted
Related change: `changes/2026-05-26-implement-pitaya-aligned-route-handler-pipeline-source-first-map/`
Related decision: `ADR-0168`

## Context

The maintainer asked to continue pushing toward Pitaya alignment with commit and push discipline. The active continuation queue after `W-0259` was `M-188/W-0260 Implement Pitaya-aligned route handler pipeline source-first map`.

`W-0259` had already defined the Pitaya-aligned route handler pipeline boundary gate, accepted `ADR-0167`, registered `runtime.pitaya_aligned_route_handler_pipeline_boundary_gate`, and opened the source-first route handler pipeline map as next-ready.

## Maintainer Narrative

Continue toward Pitaya, but keep the work bounded. Route handler pipeline vocabulary should become inspectable as future architecture vocabulary, not implemented as runtime behavior.

## Agent Response Summary

The agent treated W-0260 as a source-first inspection-map work item. It added `node tools/vibit inspect pitaya-routes --json`, accepted ADR-0168, registered the `runtime.pitaya_aligned_route_handler_pipeline_source_first_map` check rule, completed W-0260, and opened W-0261 as the next Pitaya-aligned direction selection follow-up.

RED checks confirmed the command, rule, and change artifacts were initially absent:

```text
node tools/vibit inspect pitaya-routes --json
# Unknown command.

node tools/vibit inspect rule runtime.pitaya_aligned_route_handler_pipeline_source_first_map
# Unknown rule_id: runtime.pitaya_aligned_route_handler_pipeline_source_first_map

node tools/vibit check change implement-pitaya-aligned-route-handler-pipeline-source-first-map --json
# change directory does not exist
```

## Decisions

- `ADR-0168` implements the Pitaya-aligned route handler pipeline source-first map.
- The inspection command is `node tools/vibit inspect pitaya-routes --json`.
- The allowed vocabulary remains `route_handler`, `route_key`, `handler_dispatch`, `handler_pipeline`, `pipeline_step`, `serializer_boundary`, `message_forwarding`, and `route_target`.
- Current vibit behavior remains structured `kind`/`module`/`name` protocol routing, route request handoff, explicit application dispatch, application unit-of-work transactions, generated Protobuf bridge functions, and single-process outbound message handling.
- W-0261 is the next-ready follow-up for selecting the next Pitaya-aligned direction.

## Artifacts

- `tools/vibit`
- `rules/check-rules.json`
- `decisions/ADR-0168-pitaya-aligned-route-handler-pipeline-source-first-map.md`
- `changes/2026-05-26-implement-pitaya-aligned-route-handler-pipeline-source-first-map/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `modules/friends/module.yaml`
- Repository navigation docs and module guide updates for the W-0261 next-ready state.

## Open Questions

No runtime implementation question is answered by this source-first map. A later bounded work item must separately choose any route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, serializer behavior, message forwarding behavior, backend route targeting, protocol carrier, persistence, dependency, service discovery, RPC, remote-call, frontend/backend role, cluster-safe session routing, or distributed runtime implementation.

## Follow-Up

- `M-189/W-0261 Select next Pitaya-aligned direction after route handler pipeline map`

## Redaction Notes

The inspection output exposes no raw credentials, raw access tokens, lookup digests, verifier digests, verifier keys, PostgreSQL DSNs, database payloads, local secret file contents, node credentials, transport metadata, route payloads, or local secret values.
