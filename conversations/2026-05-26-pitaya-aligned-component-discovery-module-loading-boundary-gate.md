# Conversation: Pitaya-Aligned Component Discovery And Module Loading Boundary Gate

Date: 2026-06-02
Work item: W-0286
Decision: ADR-0194
Change: `changes/2026-05-26-define-pitaya-aligned-component-discovery-module-loading-boundary-gate`
Check rule: `runtime.pitaya_aligned_component_discovery_module_loading_boundary_gate`

## Context

Continue advancing the repository toward the Pitaya-aligned architecture reference in bounded work-item steps, with commit and push discipline.

`W-0285` selected `define_pitaya_aligned_component_discovery_module_loading_boundary_gate` after the handler module registration source-first map. `W-0286` was the next-ready gate-only boundary work item.

The RED checks confirmed the expected missing surfaces:

```text
node tools/vibit inspect rule runtime.pitaya_aligned_component_discovery_module_loading_boundary_gate
# Unknown rule_id: runtime.pitaya_aligned_component_discovery_module_loading_boundary_gate

node tools/vibit check change define-pitaya-aligned-component-discovery-module-loading-boundary-gate --json
# change directory does not exist
```

## Maintainer Narrative

The maintainer asked to continue advancing toward Pitaya alignment. The active repository queue selected `W-0286 Define Pitaya-aligned component discovery and module loading boundary gate`.

## Agent Response Summary

The agent defined the gate-only component discovery and module loading vocabulary boundary and kept runtime behavior unchanged.

This completes `W-0286`, accepts `ADR-0194`, registers `runtime.pitaya_aligned_component_discovery_module_loading_boundary_gate`, and opens `W-0287 Implement Pitaya-aligned component discovery and module loading source-first map` as next-ready.

## Decisions

- Accept `ADR-0194`.
- Complete `W-0286`.
- Open `W-0287` as next-ready.
- Select `implement_pitaya_aligned_component_discovery_module_loading_source_first_map` as the follow-up direction.
- Keep component discovery behavior deferred.
- Keep component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, and direct compatibility deferred.

## Artifacts

- `docs/pitaya-aligned-component-discovery-module-loading-boundary-gate.md`
- `docs/pitaya-aligned-component-discovery-module-loading-boundary-gate.zh-CN.md`
- `decisions/ADR-0194-pitaya-aligned-component-discovery-module-loading-boundary-gate.md`
- `changes/2026-05-26-define-pitaya-aligned-component-discovery-module-loading-boundary-gate/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

No open questions for `W-0286`. The selected next work item is bounded and ready.

## Follow-Up

Proceed to `W-0287 Implement Pitaya-aligned component discovery and module loading source-first map`.

## Redaction Notes

No ignored credential file contents were read or printed. No secrets, credentials, raw access tokens, verifier material, lookup digests, verifier digests, DSNs with credentials, database payloads, transport payloads, local secret values, route payloads, session data payloads, dashboard payloads, admin console payloads, component lifecycle payloads, handler registration payloads, component discovery payloads, component inventory payloads, module loading payloads, or concrete operational payloads are recorded.
