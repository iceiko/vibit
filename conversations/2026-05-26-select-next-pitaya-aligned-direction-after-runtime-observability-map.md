# Conversation: Select Next Pitaya-Aligned Direction After Runtime Observability Map

Date: 2026-06-01

## Context

The maintainer asked to continue advancing toward Pitaya alignment. The repository next-ready item was `W-0273 Select next Pitaya-aligned direction after runtime observability map`, opened by `ADR-0180` after `node tools/vibit inspect pitaya-observability --json`.

The RED checks confirmed the expected missing surfaces:

```text
node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_runtime_observability_map
# Unknown rule_id: runtime.next_pitaya_aligned_direction_after_runtime_observability_map

node tools/vibit check change select-next-pitaya-aligned-direction-after-runtime-observability-map --json
# change directory does not exist
```

## Maintainer Narrative

Continue advancing in bounded steps toward Pitaya-class architecture coverage while preserving commit/push discipline and avoiding direct compatibility or runtime behavior expansion.

## Agent Response Summary

The agent selected `define_pitaya_aligned_metrics_tracing_boundary_gate` as the next bounded Pitaya-aligned direction. This completes W-0273, accepts `ADR-0181`, registers `runtime.next_pitaya_aligned_direction_after_runtime_observability_map`, and opens `M-202/W-0274 Define Pitaya-aligned metrics and tracing boundary gate` as next-ready.

Keep metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, and runtime endpoint behavior deferred until a later bounded work item explicitly authorizes them.

## Decisions

- Accept `ADR-0181`.
- Complete W-0273.
- Open W-0274 as next-ready.
- Keep runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, protocol changes, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, and direct compatibility deferred.

## Artifacts

- `decisions/ADR-0181-select-next-pitaya-aligned-direction-after-runtime-observability-map.md`
- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-runtime-observability-map/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

No open questions for W-0273. The selected next work item is bounded and ready.

## Follow-Up

Proceed to `W-0274 Define Pitaya-aligned metrics and tracing boundary gate`.

## Redaction Notes

No ignored credential file contents were read or printed. No secrets, credentials, raw access tokens, verifier material, lookup digests, verifier digests, DSNs with credentials, database payloads, transport payloads, local secret values, or concrete operational payloads are recorded.
