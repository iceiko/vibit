# Conversation: Select Next Pitaya-Aligned Direction After Metrics And Tracing Map

Date: 2026-06-01
Work item: W-0276
Decision: ADR-0184
Rule: runtime.next_pitaya_aligned_direction_after_metrics_tracing_map

## Context

The maintainer asked to continue advancing toward Pitaya alignment in bounded steps and to commit and push completed increments. W-0275 completed the source-first metrics and tracing inspection map, accepted ADR-0183, and opened W-0276 as the next direction selection work item.

The RED checks confirmed the expected missing surfaces:

```text
node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_metrics_tracing_map
# Unknown rule_id: runtime.next_pitaya_aligned_direction_after_metrics_tracing_map

node tools/vibit check change select-next-pitaya-aligned-direction-after-metrics-tracing-map --json
# change directory does not exist
```

## Maintainer Narrative

Continue advancing in bounded steps toward Pitaya-class architecture coverage while preserving commit/push discipline and avoiding direct compatibility or runtime behavior expansion.

## Agent Response Summary

The agent selected `define_pitaya_aligned_dashboard_admin_operations_boundary_gate` as the next bounded Pitaya-aligned direction. This completes W-0276, accepts ADR-0184, registers `runtime.next_pitaya_aligned_direction_after_metrics_tracing_map`, and opens M-205/W-0277 as next-ready.

Keep dashboard and admin operations behavior deferred until a later bounded work item explicitly authorizes it. Keep dashboard behavior, admin console behavior, runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, player/session/token inspectors, event/audit tables, protocol shape, generated output, persistence, dependencies, distributed runtime behavior, hosted surfaces, SDKs, release artifacts, and direct compatibility deferred.

## Decisions

- Accept `ADR-0184`.
- Complete W-0276.
- Open W-0277 as next-ready.
- Select `define_pitaya_aligned_dashboard_admin_operations_boundary_gate`.
- Keep runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, protocol changes, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, and direct compatibility deferred.

## Artifacts

- `decisions/ADR-0184-select-next-pitaya-aligned-direction-after-metrics-tracing-map.md`
- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-metrics-tracing-map/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

No open questions for W-0276. The selected next work item is bounded and ready.

## Follow-Up

Proceed to `W-0277 Define Pitaya-aligned dashboard and admin operations boundary gate`.

## Redaction Notes

No ignored credential file contents were read or printed. No secrets, credentials, raw access tokens, verifier material, lookup digests, verifier digests, DSNs with credentials, database payloads, transport payloads, local secret values, route payloads, session data payloads, dashboard payloads, admin console payloads, or concrete operational payloads are recorded.
