# Verification

Status: Verified

RED checks:

```text
node tools/vibit inspect rule runtime.pitaya_aligned_runtime_observability_boundary_gate
# Unknown rule_id: runtime.pitaya_aligned_runtime_observability_boundary_gate

node tools/vibit check change define-pitaya-aligned-runtime-observability-boundary-gate --json
# change directory does not exist
```

Required final checks:

```text
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_runtime_observability_boundary_gate
node tools/vibit check change define-pitaya-aligned-runtime-observability-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Observed result:

- 2026-06-01 fresh verification passed.
- `node -c tools/vibit`: passed.
- `node tools/vibit inspect next --json`: passed and reported `W-0272 Implement Pitaya-aligned runtime observability source-first map` as next-ready.
- `node tools/vibit inspect rule runtime.pitaya_aligned_runtime_observability_boundary_gate`: passed.
- `node tools/vibit check change define-pitaya-aligned-runtime-observability-boundary-gate --json`: passed with 13 passed, 0 warnings, 0 failures.
- `node tools/vibit check work --json`: passed with 1644 passed, 0 warnings, 0 failures.
- `node tools/vibit check runtime --json`: passed with 24338 passed, 1 warning, 0 failures.
- `node tools/vibit check memory --json`: passed with 4748 passed, 0 warnings, 0 failures.
- `node tools/vibit check schemas --json`: passed with 5118 passed, 0 warnings, 0 failures.
- `node tools/vibit check all --json`: passed with 326 subchecks passed, 1 warning, 0 failures.
- `git diff --check`: passed.
