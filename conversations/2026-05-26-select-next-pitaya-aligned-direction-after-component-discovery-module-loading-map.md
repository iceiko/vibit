# Conversation: Select Next Pitaya-Aligned Direction After Component Discovery And Module Loading Map

Date: 2026-06-02
Work item: W-0288
Milestone: M-216
Decision: ADR-0196
Change: `changes/2026-05-26-select-next-pitaya-aligned-direction-after-component-discovery-module-loading-map`
Rule: `runtime.next_pitaya_aligned_direction_after_component_discovery_module_loading_map`

## Context

Continue advancing toward Pitaya alignment in bounded work-item order, with commits and pushes for completed increments.

`W-0287 Implement Pitaya-aligned component discovery and module loading source-first map` completed ADR-0195, registered `runtime.pitaya_aligned_component_discovery_module_loading_source_first_map`, implemented `node tools/vibit inspect pitaya-component-loading --json`, and opened `W-0288` as the next direction selection work item.

The RED checks confirmed the expected missing W-0288 surfaces:

```text
node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_component_discovery_module_loading_map
# Unknown rule_id: runtime.next_pitaya_aligned_direction_after_component_discovery_module_loading_map

node tools/vibit check change select-next-pitaya-aligned-direction-after-component-discovery-module-loading-map --json
# change directory does not exist
```

## Maintainer Narrative

The maintainer asked to keep advancing toward Pitaya in bounded steps and to preserve commit/push discipline. The active repository queue selected `W-0288 Select next Pitaya-aligned direction after component discovery and module loading map`.

## Agent Response Summary

The agent selected `define_pitaya_aligned_startup_shutdown_hook_boundary_gate` as the next bounded Pitaya-aligned direction after the component discovery and module loading source-first map.

This completes `W-0288`, accepts `ADR-0196`, registers `runtime.next_pitaya_aligned_direction_after_component_discovery_module_loading_map`, and opens `M-217/W-0289 Define Pitaya-aligned startup and shutdown hook boundary gate` as next-ready.

## Decisions

- Accept `ADR-0196`.
- Complete `W-0288`.
- Open `W-0289` as next-ready.
- Select `define_pitaya_aligned_startup_shutdown_hook_boundary_gate`.
- Keep component discovery behavior deferred.
- Keep component loading behavior, component module loading behavior, startup hook behavior, shutdown hook behavior, lifecycle hook execution, dependency container behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, runtime endpoint behavior, dashboards, admin console behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, and direct compatibility deferred.

## Artifacts

- `decisions/ADR-0196-select-next-pitaya-aligned-direction-after-component-discovery-module-loading-map.md`
- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-component-discovery-module-loading-map/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`

## Open Questions

No open questions for `W-0288`. The selected next work item is bounded and ready.

## Follow-Up

Proceed to `W-0289 Define Pitaya-aligned startup and shutdown hook boundary gate`.

## Redaction Notes

No ignored credential file contents were read or printed. No secrets, credentials, raw access tokens, verifier material, lookup digests, verifier digests, DSNs with credentials, database payloads, transport payloads, local secret values, route payloads, session data payloads, dashboard payloads, admin console payloads, component lifecycle payloads, handler registration payloads, component module payloads, startup hook payloads, shutdown hook payloads, or concrete operational payloads are recorded.
