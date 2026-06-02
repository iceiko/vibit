# Verification

Status: Verified

RED checks:

```text
node tools/vibit inspect rule runtime.pitaya_aligned_dashboard_admin_operations_source_first_map
# Unknown rule_id: runtime.pitaya_aligned_dashboard_admin_operations_source_first_map

node tools/vibit inspect pitaya-dashboard-admin --json
# Unknown command.

node tools/vibit check change implement-pitaya-aligned-dashboard-admin-operations-source-first-map --json
# change directory does not exist
```

Required final checks:

```text
node -c tools/vibit
node tools/vibit inspect pitaya-dashboard-admin --json
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_dashboard_admin_operations_source_first_map
node tools/vibit check change implement-pitaya-aligned-dashboard-admin-operations-source-first-map --json
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
- `node tools/vibit inspect pitaya-dashboard-admin --json`: passed.
- `node tools/vibit inspect next --json`: passed and reported `W-0279 Select next Pitaya-aligned direction after dashboard/admin operations map` as next-ready.
- `node tools/vibit inspect rule runtime.pitaya_aligned_dashboard_admin_operations_source_first_map`: passed.
- `node tools/vibit check change implement-pitaya-aligned-dashboard-admin-operations-source-first-map --json`: passed with 13 passed, 0 warnings, 0 failures.
- `node tools/vibit check work --json`: passed with 1686 passed, 0 warnings, 0 failures.
- `node tools/vibit check runtime --json`: passed with 25423 passed, 1 existing warning, 0 failures.
- `node tools/vibit check memory --json`: passed with 4916 passed, 0 warnings, 0 failures.
- `node tools/vibit check schemas --json`: passed with 5272 passed, 0 warnings, 0 failures.
- `node tools/vibit check all --json`: passed with 333 subchecks passed, 1 existing warning, 0 failures.
- `git diff --check`: passed.
