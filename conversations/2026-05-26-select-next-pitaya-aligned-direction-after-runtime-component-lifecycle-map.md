# Conversation: Select Next Pitaya-Aligned Direction After Runtime Component Lifecycle Map

Date: 2026-06-02
Work item: W-0282
Decision: ADR-0190
Rule: runtime.next_pitaya_aligned_direction_after_runtime_component_lifecycle_map

## Context

The maintainer asked to continue advancing toward Pitaya alignment in bounded steps and to commit and push completed increments. W-0281 completed the source-first runtime component lifecycle inspection map, accepted ADR-0189, and opened W-0282 as the next direction selection work item.

The RED checks confirmed the expected missing surfaces:

```text
node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_runtime_component_lifecycle_map
# Unknown rule_id: runtime.next_pitaya_aligned_direction_after_runtime_component_lifecycle_map

node tools/vibit check change select-next-pitaya-aligned-direction-after-runtime-component-lifecycle-map --json
# change directory does not exist
```

## Maintainer Narrative

Continue advancing in bounded steps toward Pitaya-class architecture coverage while preserving commit/push discipline and avoiding direct compatibility or runtime behavior expansion.

## Agent Response Summary

The agent selected `define_pitaya_aligned_handler_module_registration_boundary_gate` as the next bounded Pitaya-aligned direction. This completes W-0282, accepts ADR-0190, registers `runtime.next_pitaya_aligned_direction_after_runtime_component_lifecycle_map`, and opens M-211/W-0283 Define Pitaya-aligned handler module registration boundary gate as next-ready.

Keep runtime component lifecycle behavior deferred until a later bounded work item explicitly authorizes it. Keep handler registration behavior, dynamic handler registration, component discovery or loading, component module loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol shape, generated output, persistence, dependencies, distributed runtime behavior, hosted surfaces, SDKs, release artifacts, and direct compatibility deferred.

## Decisions

- Accept `ADR-0190`.
- Complete W-0282.
- Open W-0283 as next-ready.
- Select `define_pitaya_aligned_handler_module_registration_boundary_gate`.
- Keep runtime component lifecycle behavior, handler registration behavior, dynamic handler registration, component discovery or loading, component module loading, startup hooks, shutdown hooks, runtime endpoints, dashboards, admin console behavior, protocol changes, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, and direct compatibility deferred.

## Artifacts

- `decisions/ADR-0190-select-next-pitaya-aligned-direction-after-runtime-component-lifecycle-map.md`
- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-runtime-component-lifecycle-map/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

No open questions for W-0282. The selected next work item is bounded and ready.

## Follow-Up

Proceed to `W-0283 Define Pitaya-aligned handler module registration boundary gate`.

## Redaction Notes

No ignored credential file contents were read or printed. No secrets, credentials, raw access tokens, verifier material, lookup digests, verifier digests, DSNs with credentials, database payloads, transport payloads, local secret values, route payloads, session data payloads, dashboard payloads, admin console payloads, component lifecycle payloads, handler registration payloads, component module payloads, or concrete operational payloads are recorded.
