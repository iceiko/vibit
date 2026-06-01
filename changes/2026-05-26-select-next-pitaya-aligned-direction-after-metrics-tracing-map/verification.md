# Verification

Status: Verified

RED checks:

```text
node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_metrics_tracing_map
# Unknown rule_id: runtime.next_pitaya_aligned_direction_after_metrics_tracing_map

node tools/vibit check change select-next-pitaya-aligned-direction-after-metrics-tracing-map --json
# change directory does not exist
```

Required final checks:

```text
node -c tools/vibit
node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_metrics_tracing_map
node tools/vibit inspect next --json
node tools/vibit check change select-next-pitaya-aligned-direction-after-metrics-tracing-map --json
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
- `node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_metrics_tracing_map`: passed.
- `node tools/vibit inspect next --json`: passed and reported `W-0277 Define Pitaya-aligned dashboard and admin operations boundary gate` as next-ready.
- `node tools/vibit check change select-next-pitaya-aligned-direction-after-metrics-tracing-map --json`: passed with 13 passed, 0 warnings, 0 failures.
- `node tools/vibit check work --json`: passed with 1674 passed, 0 warnings, 0 failures.
- `node tools/vibit check runtime --json`: passed with 25097 passed, 1 existing warning, 0 failures.
- Existing runtime warning: `runtime.identity_boundary` on `runtime/internal/platform/persistence/postgres/authentication_repository.go` for credential dependency boundary posture.
- `node tools/vibit check memory --json`: passed with 4868 passed, 0 warnings, 0 failures.
- `node tools/vibit check schemas --json`: passed with 5228 passed, 0 warnings, 0 failures.
- `node tools/vibit check all --json`: passed with 331 subchecks passed, 1 existing warning, 0 failures.
- `git diff --check`: passed.
