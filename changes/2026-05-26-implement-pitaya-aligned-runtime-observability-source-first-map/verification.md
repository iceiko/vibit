# Verification

Status: Verified

RED checks:

```text
node tools/vibit inspect rule runtime.pitaya_aligned_runtime_observability_source_first_map
# Unknown rule_id: runtime.pitaya_aligned_runtime_observability_source_first_map

node tools/vibit inspect pitaya-observability --json
# Unknown command.

node tools/vibit check change implement-pitaya-aligned-runtime-observability-source-first-map --json
# change directory does not exist
```

Required final checks:

```text
node -c tools/vibit
node tools/vibit inspect pitaya-observability --json
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_runtime_observability_source_first_map
node tools/vibit check change implement-pitaya-aligned-runtime-observability-source-first-map --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Observed result:

- `node -c tools/vibit`: passed.
- `node tools/vibit inspect pitaya-observability --json`: passed and reported `W-0273 Select next Pitaya-aligned direction after runtime observability map`.
- `node tools/vibit inspect next --json`: passed and reported W-0273 as next-ready.
- `node tools/vibit inspect rule runtime.pitaya_aligned_runtime_observability_source_first_map`: passed.
- `node tools/vibit check change implement-pitaya-aligned-runtime-observability-source-first-map --json`: passed, 13 passed, 0 warnings, 0 failures.
- `node tools/vibit check work --json`: passed, 1650 passed, 0 warnings, 0 failures.
- `node tools/vibit check runtime --json`: passed, 24501 passed, 1 existing warning, 0 failures.
- `node tools/vibit check memory --json`: passed, 4772 passed, 0 warnings, 0 failures.
- `node tools/vibit check schemas --json`: passed, 5140 passed, 0 warnings, 0 failures.
- `node tools/vibit check all --json`: passed, 327/327 subchecks passed, 1 existing warning, 0 failures.
- `git diff --check`: passed.
