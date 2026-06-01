# Conversation: Select Next Pitaya-Aligned Direction After Route Handler Pipeline Map

Date: 2026-06-01
Status: Accepted
Related change: `changes/2026-05-26-select-next-pitaya-aligned-direction-after-route-handler-pipeline-map/`
Related decision: `ADR-0169`

## Context

The maintainer asked to continue pushing toward Pitaya alignment. `W-0260` had completed the source-first route handler pipeline map, including `node tools/vibit inspect pitaya-routes --json`, and opened `W-0261` as the next-ready direction-selection step.

## Maintainer Narrative

Continue advancing the repository in bounded steps toward Pitaya-class architecture, commit and push completed slices, and avoid exposing ignored local credentials.

## Agent Response Summary

The agent treated W-0261 as selection-only. It selected `define_pitaya_aligned_serializer_message_forwarding_boundary_gate` as the next bounded Pitaya-aligned direction, accepted `ADR-0169`, registered `runtime.next_pitaya_aligned_direction_after_route_handler_pipeline_map`, completed W-0261, and opened W-0262 as next-ready.

## Decisions

- Select `define_pitaya_aligned_serializer_message_forwarding_boundary_gate` after the route handler pipeline map.
- Keep serializer behavior and message forwarding behavior deferred.
- Keep backend route targeting, service discovery, RPC, frontend/backend roles, distributed runtime behavior, protocol changes, generated output, persistence, dependencies, hosted surfaces, SDKs, and direct compatibility deferred.

## Artifacts

- `decisions/ADR-0169-select-next-pitaya-aligned-direction-after-route-handler-pipeline-map.md`
- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-route-handler-pipeline-map/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

None for W-0261. W-0262 must define the serializer and message forwarding boundary without implementing runtime behavior.

## Follow-Up

- Advance `W-0262 Define Pitaya-aligned serializer and message forwarding boundary gate`.
- Preserve all implementation deferrals until a later explicit bounded work item authorizes them.

## Redaction Notes

No ignored credential file contents, token values, DSNs with credentials, raw credentials, raw access tokens, verifier digests, lookup digests, node credentials, transport payloads, or route payload contents were recorded.
