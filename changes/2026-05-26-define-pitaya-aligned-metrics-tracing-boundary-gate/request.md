# Request

Continue advancing toward Pitaya alignment from `W-0274 Define Pitaya-aligned metrics and tracing boundary gate`.

This slice is gate-only. It must define future metrics and tracing vocabulary after the runtime observability source-first map and must not add runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, protocol changes, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## RED Checks

```text
node tools/vibit inspect rule runtime.pitaya_aligned_metrics_tracing_boundary_gate
# Unknown rule_id: runtime.pitaya_aligned_metrics_tracing_boundary_gate

node tools/vibit check change define-pitaya-aligned-metrics-tracing-boundary-gate --json
# change directory does not exist
```
