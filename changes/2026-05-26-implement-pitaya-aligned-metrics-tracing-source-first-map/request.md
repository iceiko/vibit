# Request

Continue advancing toward Pitaya alignment by implementing the next-ready `W-0275 Implement Pitaya-aligned metrics and tracing source-first map`.

The change must stay source-first and must not add runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, protocol messages or routes, Protobuf sources, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## RED Checks

```text
node tools/vibit inspect rule runtime.pitaya_aligned_metrics_tracing_source_first_map
# Unknown rule_id: runtime.pitaya_aligned_metrics_tracing_source_first_map

node tools/vibit inspect pitaya-metrics-tracing --json
# Unknown command.

node tools/vibit check change implement-pitaya-aligned-metrics-tracing-source-first-map --json
# change directory does not exist
```
