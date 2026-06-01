# Verification

Status: Verified

RED checks:

```text
node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_runtime_observability_map
# Unknown rule_id: runtime.next_pitaya_aligned_direction_after_runtime_observability_map

node tools/vibit check change select-next-pitaya-aligned-direction-after-runtime-observability-map --json
# change directory does not exist
```

Required final checks:

```text
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_runtime_observability_map
node tools/vibit check change select-next-pitaya-aligned-direction-after-runtime-observability-map --json
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
- `node tools/vibit inspect next --json`: passed and reported `W-0274 Define Pitaya-aligned metrics and tracing boundary gate` as next-ready.
- `node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_runtime_observability_map`: passed.
- `node tools/vibit check change select-next-pitaya-aligned-direction-after-runtime-observability-map --json`: passed with 13 passed, 0 warnings, 0 failures.
- `node tools/vibit check work --json`: passed with 1656 passed, 0 warnings, 0 failures.
- `node tools/vibit check runtime --json`: passed with 24641 passed, 1 existing warning, 0 failures.
- `node tools/vibit check memory --json`: passed with 4796 passed, 0 warnings, 0 failures.
- `node tools/vibit check schemas --json`: passed with 5162 passed, 0 warnings, 0 failures.
- `node tools/vibit check all --json`: passed with 328 subchecks passed, 1 existing warning, 0 failures.
- `git diff --check`: passed.
