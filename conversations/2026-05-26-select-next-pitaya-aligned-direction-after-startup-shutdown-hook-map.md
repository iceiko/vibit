# Conversation: Select Next Pitaya-Aligned Direction After Startup And Shutdown Hook Map

Date: 2026-06-07
Work item: W-0291
Milestone: M-219
Decision: ADR-0199
Change: `changes/2026-05-26-select-next-pitaya-aligned-direction-after-startup-shutdown-hook-map`
Rule: `runtime.next_pitaya_aligned_direction_after_startup_shutdown_hook_map`

## Context

Continue repository work from the current `next_ready` item for up to 20 steps.

`W-0290 Implement Pitaya-aligned startup and shutdown hook source-first map` completed ADR-0198, registered `runtime.pitaya_aligned_startup_shutdown_hook_source_first_map`, implemented `node tools/vibit inspect pitaya-startup-shutdown --json`, and opened `W-0291` as the next direction selection work item.

The RED checks confirmed the expected missing W-0291 surfaces:

```text
node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_startup_shutdown_hook_map
# Unknown rule_id: runtime.next_pitaya_aligned_direction_after_startup_shutdown_hook_map

node tools/vibit check change select-next-pitaya-aligned-direction-after-startup-shutdown-hook-map --json
# change directory does not exist
```

`node tools/vibit inspect next --json` reported `W-0291` as next-ready.

## Maintainer Narrative

The maintainer asked to continue 20 steps. The repository queue selected `W-0291 Select next Pitaya-aligned direction after startup and shutdown hook map`.

## Agent Response Summary

The agent selected `define_currency_wallet_lifecycle_boundary_gate` as the next bounded direction after the Pitaya-aligned startup and shutdown hook source-first map.

This completes `W-0291`, accepts `ADR-0199`, registers `runtime.next_pitaya_aligned_direction_after_startup_shutdown_hook_map`, and opens `M-220/W-0292 Define currency wallet lifecycle boundary gate` as next-ready.

The selection ends the current Pitaya vocabulary pass at the startup/shutdown hook source-first map and returns the roadmap to Phase 3 core game backend modules.

## Decisions

- Accept `ADR-0199`.
- Complete `W-0291`.
- Open `W-0292` as next-ready.
- Select `define_currency_wallet_lifecycle_boundary_gate`.
- Keep currency wallet behavior deferred.
- Keep balance tables, wallet transaction behavior, reward integration, inventory integration, purchase behavior, grant/spend behavior, audit/event tables, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime, startup hook behavior, shutdown hook behavior, lifecycle execution, dependency containers, and direct compatibility deferred.

## Artifacts

- `decisions/ADR-0199-select-next-pitaya-aligned-direction-after-startup-shutdown-hook-map.md`
- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-startup-shutdown-hook-map/`
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

No open questions for `W-0291`. The selected next work item is bounded and ready.

## Follow-Up

Proceed to `W-0292 Define currency wallet lifecycle boundary gate`.

## Redaction Notes

No ignored credential file contents were read or printed. No secrets, credentials, raw access tokens, verifier material, lookup digests, verifier digests, DSNs with credentials, database payloads, transport payloads, local secret values, route payloads, session data payloads, dashboard payloads, admin console payloads, component lifecycle payloads, handler registration payloads, component module payloads, startup hook payloads, shutdown hook payloads, currency wallet payloads, wallet transaction payloads, or concrete operational payloads are recorded.
