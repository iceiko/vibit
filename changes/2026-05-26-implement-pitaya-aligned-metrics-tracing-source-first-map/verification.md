# Verification

Status: Verified

RED checks:

```text
node tools/vibit inspect rule runtime.pitaya_aligned_metrics_tracing_source_first_map
# Unknown rule_id: runtime.pitaya_aligned_metrics_tracing_source_first_map

node tools/vibit inspect pitaya-metrics-tracing --json
# Unknown command.

node tools/vibit check change implement-pitaya-aligned-metrics-tracing-source-first-map --json
# change directory does not exist
```

Required final checks:

```text
node -c tools/vibit
node tools/vibit inspect pitaya-metrics-tracing --json
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_metrics_tracing_source_first_map
node tools/vibit check change implement-pitaya-aligned-metrics-tracing-source-first-map --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Observed result:

- 2026-06-01 focused verification passed.
- `node -c tools/vibit`: passed.
- `node tools/vibit inspect pitaya-metrics-tracing --json`: passed.
- `node tools/vibit inspect next --json`: passed and reported `W-0276 Select next Pitaya-aligned direction after metrics and tracing map` as next-ready.
- `node tools/vibit inspect rule runtime.pitaya_aligned_metrics_tracing_source_first_map`: passed.
- `node tools/vibit check change implement-pitaya-aligned-metrics-tracing-source-first-map --json`: passed with 13 passed, 0 warnings, 0 failures.
- `node tools/vibit check work --json`: passed with 1668 passed, 0 warnings, 0 failures.
- `node tools/vibit check runtime --json`: passed with 24960 passed, 1 existing warning, 0 failures.
- `node tools/vibit check memory --json`: passed with 4844 passed, 0 warnings, 0 failures.
- `node tools/vibit check schemas --json`: passed with 5206 passed, 0 warnings, 0 failures.
- `node tools/vibit check all --json`: passed with 330 subchecks passed, 1 existing warning, 0 failures.
- `git diff --check`: passed.
