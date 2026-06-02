# Conversation: Pitaya-Aligned Runtime Component Lifecycle Source-First Map

Date: 2026-06-02
Work item: W-0281
Decision: ADR-0189
Change: `changes/2026-05-26-implement-pitaya-aligned-runtime-component-lifecycle-source-first-map`
Check rule: `runtime.pitaya_aligned_runtime_component_lifecycle_source_first_map`

## Context

Continue advancing the repository toward the Pitaya-aligned architecture reference in bounded work-item steps, with commit and push discipline.

`W-0280` completed the runtime component lifecycle boundary gate. `W-0281` was the next-ready source-first map implementation.

## Maintainer Narrative

The maintainer asked to continue advancing toward Pitaya alignment. The active repository queue selected `W-0281 Implement Pitaya-aligned runtime component lifecycle source-first map`.

## RED Evidence

Before implementation, the repository did not yet know the W-0281 source-first map rule or command:

```text
Unknown rule_id: runtime.pitaya_aligned_runtime_component_lifecycle_source_first_map
Unknown command
change directory does not exist
```

The missing command was expected to be `node tools/vibit inspect pitaya-component-lifecycle --json`.

## Agent Response Summary

The agent implemented the source-first map and kept runtime behavior unchanged.

## Decisions

Implement a source-first runtime component lifecycle inspection map only.

The work registers `runtime.pitaya_aligned_runtime_component_lifecycle_source_first_map`, implements `node tools/vibit inspect pitaya-component-lifecycle --json`, records ADR-0189, updates repository memory, and opens `W-0282 Select next Pitaya-aligned direction after runtime component lifecycle map` as next-ready.

## Artifacts

- `tools/vibit`
- `rules/check-rules.json`
- `decisions/ADR-0189-pitaya-aligned-runtime-component-lifecycle-source-first-map.md`
- `changes/2026-05-26-implement-pitaya-aligned-runtime-component-lifecycle-source-first-map/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`

## Open Questions

- Which Pitaya-aligned direction should follow the runtime component lifecycle source-first map remains for `W-0282`.

## Follow-Up

- Complete `W-0282 Select next Pitaya-aligned direction after runtime component lifecycle map`.

## Redaction Notes

No secret values, local ignored configuration contents, raw credentials, access tokens, lookup digests, verifier digests, DSNs, headers, cookies, query strings, route payloads, session data payloads, component lifecycle payloads, or concrete transport metadata were recorded.

## Boundaries

Keep runtime component lifecycle behavior, handler registration behavior, component discovery or loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, and direct compatibility deferred.
