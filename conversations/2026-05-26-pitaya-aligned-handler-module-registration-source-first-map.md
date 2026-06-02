# Conversation: Pitaya-Aligned Handler Module Registration Source-First Map

Date: 2026-06-02
Work item: W-0284
Decision: ADR-0192
Change: `changes/2026-05-26-implement-pitaya-aligned-handler-module-registration-source-first-map`
Check rule: `runtime.pitaya_aligned_handler_module_registration_source_first_map`

## Context

Continue advancing the repository toward the Pitaya-aligned architecture reference in bounded work-item steps, with commit and push discipline.

`W-0283` completed the handler module registration boundary gate. `W-0284` was the next-ready source-first map implementation.

## Maintainer Narrative

The maintainer asked to continue advancing toward Pitaya alignment. The active repository queue selected `W-0284 Implement Pitaya-aligned handler module registration source-first map`.

## RED Evidence

Before implementation, the repository did not yet know the W-0284 source-first map rule or command:

```text
Unknown rule_id: runtime.pitaya_aligned_handler_module_registration_source_first_map
Unknown command
change directory does not exist
```

The missing command was expected to be `node tools/vibit inspect pitaya-handler-modules --json`.

## Agent Response Summary

The agent implemented the source-first map and kept runtime behavior unchanged.

## Decisions

Implement a source-first handler module registration inspection map only.

The work registers `runtime.pitaya_aligned_handler_module_registration_source_first_map`, implements `node tools/vibit inspect pitaya-handler-modules --json`, records ADR-0192, updates repository memory, and opens `W-0285 Select next Pitaya-aligned direction after handler module registration map` as next-ready.

## Artifacts

- `tools/vibit`
- `rules/check-rules.json`
- `decisions/ADR-0192-pitaya-aligned-handler-module-registration-source-first-map.md`
- `changes/2026-05-26-implement-pitaya-aligned-handler-module-registration-source-first-map/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`

## Open Questions

- Which Pitaya-aligned direction should follow the handler module registration source-first map remains for `W-0285`.

## Follow-Up

- Complete `W-0285 Select next Pitaya-aligned direction after handler module registration map`.

## Redaction Notes

No secret values, local ignored configuration contents, raw credentials, access tokens, lookup digests, verifier digests, verifier keys, DSNs, headers, cookies, query strings, route payloads, session data payloads, component lifecycle payloads, handler registration payloads, handler module registration payloads, or concrete transport metadata were recorded.

## Boundaries

Keep handler module registration behavior, handler registration behavior, dynamic handler registration, component discovery or loading, component module loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, and direct compatibility deferred.
