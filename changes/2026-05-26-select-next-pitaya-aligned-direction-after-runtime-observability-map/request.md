# Request

Continue advancing toward Pitaya alignment from `W-0273 Select next Pitaya-aligned direction after runtime observability map`.

This slice is selection-only. It must choose the next bounded Pitaya-aligned direction after the source-first runtime observability map and must not add runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, protocol changes, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## RED Checks

```text
node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_runtime_observability_map
# Unknown rule_id: runtime.next_pitaya_aligned_direction_after_runtime_observability_map

node tools/vibit check change select-next-pitaya-aligned-direction-after-runtime-observability-map --json
# change directory does not exist
```
