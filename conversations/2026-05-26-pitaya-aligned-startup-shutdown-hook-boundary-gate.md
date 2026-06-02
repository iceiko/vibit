# Conversation: Pitaya-Aligned Startup And Shutdown Hook Boundary Gate

Date: 2026-06-02
Work item: W-0289
Milestone: M-217
Decision: ADR-0197
Change: `changes/2026-05-26-define-pitaya-aligned-startup-shutdown-hook-boundary-gate`
Rule: `runtime.pitaya_aligned_startup_shutdown_hook_boundary_gate`

## Context

Continue advancing toward Pitaya alignment in bounded work-item order, with commits and pushes for completed increments.

`W-0288 Select next Pitaya-aligned direction after component discovery and module loading map` completed ADR-0196, registered `runtime.next_pitaya_aligned_direction_after_component_discovery_module_loading_map`, selected `define_pitaya_aligned_startup_shutdown_hook_boundary_gate`, and opened `W-0289` as the next-ready work item.

The RED checks confirmed the expected missing W-0289 surfaces:

```text
node tools/vibit inspect rule runtime.pitaya_aligned_startup_shutdown_hook_boundary_gate
# Unknown rule_id: runtime.pitaya_aligned_startup_shutdown_hook_boundary_gate

node tools/vibit check change define-pitaya-aligned-startup-shutdown-hook-boundary-gate --json
# change directory does not exist
```

## Maintainer Narrative

The maintainer asked to keep advancing toward Pitaya in bounded steps and to preserve commit/push discipline. The active repository queue selected `W-0289 Define Pitaya-aligned startup and shutdown hook boundary gate`.

## Agent Response Summary

The agent defined the startup and shutdown hook vocabulary gate as a documentation and repository-memory boundary only.

This completes `W-0289`, accepts `ADR-0197`, registers `runtime.pitaya_aligned_startup_shutdown_hook_boundary_gate`, defines `docs/pitaya-aligned-startup-shutdown-hook-boundary-gate.md`, and opens `M-218/W-0290 Implement Pitaya-aligned startup and shutdown hook source-first map` as next-ready.

## Decisions

- Accept `ADR-0197`.
- Complete `W-0289`.
- Open `W-0290` as next-ready.
- Define startup and shutdown hook vocabulary only.
- Keep startup hook behavior, shutdown hook behavior, lifecycle hook execution, dependency container behavior, component discovery behavior, component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, runtime endpoint behavior, dashboards, admin console behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, and direct compatibility deferred.

## Artifacts

- `docs/pitaya-aligned-startup-shutdown-hook-boundary-gate.md`
- `docs/pitaya-aligned-startup-shutdown-hook-boundary-gate.zh-CN.md`
- `decisions/ADR-0197-pitaya-aligned-startup-shutdown-hook-boundary-gate.md`
- `changes/2026-05-26-define-pitaya-aligned-startup-shutdown-hook-boundary-gate/`
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

No open questions for `W-0289`. The selected source-first follow-up is bounded and ready.

## Follow-Up

Proceed to `W-0290 Implement Pitaya-aligned startup and shutdown hook source-first map`.

## Redaction Notes

No ignored credential file contents were read or printed. No secrets, credentials, raw access tokens, verifier material, lookup digests, verifier digests, DSNs with credentials, database payloads, transport payloads, local secret values, route payloads, session data payloads, dashboard payloads, admin console payloads, component lifecycle payloads, handler registration payloads, component module payloads, startup hook payloads, shutdown hook payloads, or concrete operational payloads are recorded.
