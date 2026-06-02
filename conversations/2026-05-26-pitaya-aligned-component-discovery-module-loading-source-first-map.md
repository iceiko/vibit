# Conversation: Pitaya-Aligned Component Discovery And Module Loading Source-First Map

Date: 2026-06-02
Work item: `W-0287`
Milestone: `M-215`
Decision: `ADR-0195`
Change: `changes/2026-05-26-implement-pitaya-aligned-component-discovery-module-loading-source-first-map/`

## Context

The maintainer asked to continue advancing toward Pitaya alignment and later clarified that the next phase should pause. This slice completes the current W-0287 map and stops before implementing W-0288.

## Maintainer Narrative

The requested direction was to continue moving toward Pitaya-class architecture. The newest instruction says to pause at the next phase, so the agent should finish, verify, commit, and push W-0287, then stop before implementing W-0288.

## Agent Response Summary

Before implementation, the expected checks failed:

- `node tools/vibit inspect rule runtime.pitaya_aligned_component_discovery_module_loading_source_first_map` returned `Unknown rule_id: runtime.pitaya_aligned_component_discovery_module_loading_source_first_map`.
- `node tools/vibit inspect pitaya-component-loading --json` returned `Unknown command`.
- `node tools/vibit check change implement-pitaya-aligned-component-discovery-module-loading-source-first-map --json` reported that the change directory does not exist.

Implemented `runtime.pitaya_aligned_component_discovery_module_loading_source_first_map` and `node tools/vibit inspect pitaya-component-loading --json`.

The map reports the component discovery and module loading vocabulary, current explicit bootstrap composition, application module ownership, handler module registration handoff, explicit source inventory, dynamic loading deferral, distributed discovery deferral, source surfaces, and redaction posture.

## Decisions

- Accepted `ADR-0195`.
- Registered `runtime.pitaya_aligned_component_discovery_module_loading_source_first_map`.
- Added `node tools/vibit inspect pitaya-component-loading --json`.
- Marked `W-0287` complete and opened `W-0288 Select next Pitaya-aligned direction after component discovery and module loading map` as next-ready.

## Artifacts

- `decisions/ADR-0195-pitaya-aligned-component-discovery-module-loading-source-first-map.md`
- `changes/2026-05-26-implement-pitaya-aligned-component-discovery-module-loading-source-first-map/`
- `tools/vibit`
- `rules/check-rules.json`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`

## Open Questions

- Which bounded Pitaya-aligned direction should W-0288 select after the component discovery and module loading map?
- Should future component loading planning split dependency container vocabulary from startup/shutdown hook vocabulary?

## Follow-Up

`W-0288 Select next Pitaya-aligned direction after component discovery and module loading map` is next-ready. Per the maintainer's latest instruction, pause at that next phase after W-0287 is committed and pushed.

## Redaction Notes

Keep component discovery behavior deferred.
Keep component loading behavior deferred.
Keep component module loading behavior deferred.
Keep handler module registration behavior deferred.
Keep handler registration behavior deferred.
Keep dynamic handler registration deferred.
Keep startup hooks and shutdown hooks deferred.
Keep runtime endpoint behavior, dashboards, and admin console behavior deferred.
Keep protocol, generated output, persistence, dependency, hosted, SDK, release, distributed runtime, and direct Nakama/Pitaya compatibility work deferred.
No raw credentials, tokens, verifier keys, DSNs, headers, route payloads, component inventory payloads, module loading payloads, or local secret file contents are recorded.
