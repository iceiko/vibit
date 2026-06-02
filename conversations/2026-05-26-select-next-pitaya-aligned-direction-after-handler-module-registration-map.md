# Conversation: Select Next Pitaya-Aligned Direction After Handler Module Registration Map

Date: 2026-06-02
Work item: W-0285
Milestone: M-213
Decision: ADR-0193
Change: `changes/2026-05-26-select-next-pitaya-aligned-direction-after-handler-module-registration-map`
Rule: `runtime.next_pitaya_aligned_direction_after_handler_module_registration_map`

## Context

Continue advancing toward Pitaya alignment in bounded work-item order, with commits and pushes for completed increments.

`W-0284 Implement Pitaya-aligned handler module registration source-first map` completed ADR-0192, registered `runtime.pitaya_aligned_handler_module_registration_source_first_map`, implemented `node tools/vibit inspect pitaya-handler-modules --json`, and opened `W-0285` as the next direction selection work item.

The RED checks confirmed the expected missing W-0285 surfaces:

```text
node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_handler_module_registration_map
# Unknown rule_id: runtime.next_pitaya_aligned_direction_after_handler_module_registration_map

node tools/vibit check change select-next-pitaya-aligned-direction-after-handler-module-registration-map --json
# change directory does not exist
```

## Maintainer Narrative

The maintainer asked to keep advancing toward Pitaya in bounded steps and to preserve commit/push discipline. The active repository queue selected `W-0285 Select next Pitaya-aligned direction after handler module registration map`.

## Agent Response Summary

The agent selected `define_pitaya_aligned_component_discovery_module_loading_boundary_gate` as the next bounded Pitaya-aligned direction after the handler module registration source-first map.

This completes `W-0285`, accepts `ADR-0193`, registers `runtime.next_pitaya_aligned_direction_after_handler_module_registration_map`, and opens `M-214/W-0286 Define Pitaya-aligned component discovery and module loading boundary gate` as next-ready.

## Decisions

- Accept `ADR-0193`.
- Complete `W-0285`.
- Open `W-0286` as next-ready.
- Select `define_pitaya_aligned_component_discovery_module_loading_boundary_gate`.
- Keep handler module registration behavior deferred.
- Keep handler module registration behavior, handler registration behavior, dynamic handler registration, component discovery behavior, component loading behavior, component module loading behavior, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, and direct compatibility deferred.

## Artifacts

- `decisions/ADR-0193-select-next-pitaya-aligned-direction-after-handler-module-registration-map.md`
- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-handler-module-registration-map/`
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

No open questions for `W-0285`. The selected next work item is bounded and ready.

## Follow-Up

Proceed to `W-0286 Define Pitaya-aligned component discovery and module loading boundary gate`.

## Redaction Notes

No ignored credential file contents were read or printed. No secrets, credentials, raw access tokens, verifier material, lookup digests, verifier digests, DSNs with credentials, database payloads, transport payloads, local secret values, route payloads, session data payloads, dashboard payloads, admin console payloads, component lifecycle payloads, handler registration payloads, component module payloads, or concrete operational payloads are recorded.
