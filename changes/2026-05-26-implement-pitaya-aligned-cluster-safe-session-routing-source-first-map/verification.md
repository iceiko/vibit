# Verification

Status: Verified

## RED Checks

```text
node tools/vibit inspect pitaya-sessions --json
# Unknown command.

node tools/vibit inspect rule runtime.pitaya_aligned_cluster_safe_session_routing_source_first_map
# Unknown rule_id: runtime.pitaya_aligned_cluster_safe_session_routing_source_first_map

node tools/vibit check change implement-pitaya-aligned-cluster-safe-session-routing-source-first-map --json
# change directory does not exist
```

## Required Verification

- `node -c tools/vibit` - passed.
- `node tools/vibit inspect next --json` - passed; reports `W-0258 Select next Pitaya-aligned direction after cluster-safe session routing map`.
- `node tools/vibit inspect pitaya-sessions --json` - passed; reports the W-0257 source-first session-routing map and W-0258 follow-up.
- `node tools/vibit inspect rule runtime.pitaya_aligned_cluster_safe_session_routing_source_first_map` - passed.
- `node tools/vibit check change implement-pitaya-aligned-cluster-safe-session-routing-source-first-map --json` - passed with 13 passed, 0 warnings, 0 failures.
- `node tools/vibit check work --json` - passed with 1560 passed, 0 warnings, 0 failures.
- `node tools/vibit check runtime --json` - passed with 21831 passed, 1 known warning, 0 failures. The warning is the existing `runtime.identity_boundary` warning for `runtime/internal/platform/persistence/postgres/authentication_repository.go`.
- `node tools/vibit check memory --json` - passed with 4412 passed, 0 warnings, 0 failures.
- `node tools/vibit check schemas --json` - passed with 4810 passed, 0 warnings, 0 failures.
- `node tools/vibit check all --json` - passed with 312 subchecks passed, 1 known warning, 0 failures.
- `git diff --check` - passed.
