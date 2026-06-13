# Conversation: Pitaya-Aligned Startup And Shutdown Hook Source-First Map

Date: 2026-06-07
Work item: `W-0290`
Milestone: `M-218`
Decision: `ADR-0198`
Change: `changes/2026-05-26-implement-pitaya-aligned-startup-shutdown-hook-source-first-map/`

## Context

The maintainer asked to continue. The active next-ready work item was `W-0290 Implement Pitaya-aligned startup and shutdown hook source-first map`.

## Maintainer Narrative

The requested direction was to continue one next-ready work item according to `.arch/work-items.yaml` continuation semantics.

## Agent Response Summary

Before implementation, the expected checks failed:

- `node tools/vibit inspect rule runtime.pitaya_aligned_startup_shutdown_hook_source_first_map` returned `Unknown rule_id: runtime.pitaya_aligned_startup_shutdown_hook_source_first_map`.
- `node tools/vibit inspect pitaya-startup-shutdown --json` returned `Unknown command`.
- `node tools/vibit check change implement-pitaya-aligned-startup-shutdown-hook-source-first-map --json` reported that the change directory does not exist.

Implemented `runtime.pitaya_aligned_startup_shutdown_hook_source_first_map` and `node tools/vibit inspect pitaya-startup-shutdown --json`.

The map reports startup and shutdown hook vocabulary, current explicit bootstrap composition, startup hook deferral, shutdown hook deferral, lifecycle hook execution deferral, dependency handoff deferral, module loading handoff, distributed lifecycle deferral, source surfaces, and redaction posture.

## Decisions

- Accepted `ADR-0198`.
- Registered `runtime.pitaya_aligned_startup_shutdown_hook_source_first_map`.
- Added `node tools/vibit inspect pitaya-startup-shutdown --json`.
- Marked `W-0290` complete and opened `W-0291 Select next Pitaya-aligned direction after startup and shutdown hook map` as next-ready.

## Artifacts

- `decisions/ADR-0198-pitaya-aligned-startup-shutdown-hook-source-first-map.md`
- `changes/2026-05-26-implement-pitaya-aligned-startup-shutdown-hook-source-first-map/`
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

- Which bounded Pitaya-aligned direction should W-0291 select after the startup and shutdown hook map?
- Should future startup/shutdown planning split ordering vocabulary from lifecycle execution vocabulary?

## Follow-Up

`W-0291 Select next Pitaya-aligned direction after startup and shutdown hook map` is next-ready.

## Redaction Notes

Keep startup hook behavior deferred.
Keep shutdown hook behavior deferred.
Keep lifecycle hook execution deferred.
Keep dependency container behavior deferred.
Keep component discovery behavior deferred.
Keep component loading behavior deferred.
Keep component module loading behavior deferred.
Keep handler module registration behavior deferred.
Keep handler registration behavior deferred.
Keep dynamic handler registration deferred.
Keep runtime endpoint behavior, dashboards, and admin console behavior deferred.
Keep protocol, generated output, persistence, dependency, hosted, SDK, release, distributed runtime, and direct Nakama/Pitaya compatibility work deferred.
No raw credentials, tokens, verifier keys, DSNs, headers, route payloads, component inventory payloads, module loading payloads, startup hook payloads, shutdown hook payloads, or local secret file contents are recorded.
