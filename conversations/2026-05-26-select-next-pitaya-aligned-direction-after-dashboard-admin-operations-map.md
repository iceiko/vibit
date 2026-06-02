# Conversation: Select Next Pitaya-Aligned Direction After Dashboard/Admin Operations Map

Date: 2026-06-02
Work item: W-0279
Decision: ADR-0187
Rule: runtime.next_pitaya_aligned_direction_after_dashboard_admin_operations_map

## Context

The maintainer asked to continue advancing toward Pitaya alignment in bounded steps and to commit and push completed increments. W-0278 completed the source-first dashboard and admin operations inspection map, accepted ADR-0186, and opened W-0279 as the next direction selection work item.

The RED checks confirmed the expected missing surfaces:

```text
node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_dashboard_admin_operations_map
# Unknown rule_id: runtime.next_pitaya_aligned_direction_after_dashboard_admin_operations_map

node tools/vibit check change select-next-pitaya-aligned-direction-after-dashboard-admin-operations-map --json
# change directory does not exist
```

## Maintainer Narrative

Continue advancing in bounded steps toward Pitaya-class architecture coverage while preserving commit/push discipline and avoiding direct compatibility or runtime behavior expansion.

## Agent Response Summary

The agent selected `define_pitaya_aligned_runtime_component_lifecycle_boundary_gate` as the next bounded Pitaya-aligned direction. This completes W-0279, accepts ADR-0187, registers `runtime.next_pitaya_aligned_direction_after_dashboard_admin_operations_map`, and opens M-208/W-0280 as next-ready.

Keep runtime component lifecycle behavior deferred until a later bounded work item explicitly authorizes it. Keep handler registration behavior, component discovery or loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol shape, generated output, persistence, dependencies, distributed runtime behavior, hosted surfaces, SDKs, release artifacts, and direct compatibility deferred.

## Decisions

- Accept `ADR-0187`.
- Complete W-0279.
- Open W-0280 as next-ready.
- Select `define_pitaya_aligned_runtime_component_lifecycle_boundary_gate`.
- Keep runtime component lifecycle behavior, handler registration behavior, component discovery or loading, startup hooks, shutdown hooks, runtime endpoints, dashboards, admin console behavior, protocol changes, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, and direct compatibility deferred.

## Artifacts

- `decisions/ADR-0187-select-next-pitaya-aligned-direction-after-dashboard-admin-operations-map.md`
- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-dashboard-admin-operations-map/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

No open questions for W-0279. The selected next work item is bounded and ready.

## Follow-Up

Proceed to `W-0280 Define Pitaya-aligned runtime component lifecycle boundary gate`.

## Redaction Notes

No ignored credential file contents were read or printed. No secrets, credentials, raw access tokens, verifier material, lookup digests, verifier digests, DSNs with credentials, database payloads, transport payloads, local secret values, route payloads, session data payloads, dashboard payloads, admin console payloads, component lifecycle payloads, handler registration payloads, or concrete operational payloads are recorded.
