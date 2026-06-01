# Verification

Status: Verified

RED checks:

```text
node tools/vibit inspect pitaya-session-lifecycle --json
# Unknown command.

node tools/vibit inspect rule runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_source_first_map
# Unknown rule_id: runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_source_first_map

node tools/vibit check change implement-pitaya-aligned-session-binding-kick-disconnect-session-data-source-first-map --json
# change directory does not exist
```

Required final checks:

```text
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect pitaya-session-lifecycle --json
node tools/vibit inspect rule runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_source_first_map
node tools/vibit check change implement-pitaya-aligned-session-binding-kick-disconnect-session-data-source-first-map --json
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
- `node tools/vibit inspect pitaya-session-lifecycle --json`: passed and reported `kind: pitaya_session_lifecycle_inspection`, `status: source_first_pitaya_aligned_session_binding_kick_disconnect_session_data_map`, and `next_ready_work_item: W-0270`.
- `node tools/vibit inspect rule runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_source_first_map`: passed.
- `node tools/vibit inspect next --json`: passed and reported `W-0270` as the only next-ready work item.
- `node tools/vibit check change implement-pitaya-aligned-session-binding-kick-disconnect-session-data-source-first-map --json`: passed with 13 passed, 0 warnings, 0 failures.
- `node tools/vibit check work --json`: passed with 1632 passed, 0 warnings, 0 failures.
- `node tools/vibit check runtime --json`: passed with 23997 passed, 1 warning, 0 failures.
- `node tools/vibit check memory --json`: passed with 4700 passed, 0 warnings, 0 failures.
- `node tools/vibit check schemas --json`: passed with 5074 passed, 0 warnings, 0 failures.
- `node tools/vibit check all --json`: passed with 324 subchecks passed, 1 warning, 0 failures.
- `git diff --check`: passed.
