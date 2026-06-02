# Verification

Status: Verified

Required commands:

```text
node -c tools/vibit
node tools/vibit inspect pitaya-component-loading --json
node tools/vibit inspect rule runtime.pitaya_aligned_component_discovery_module_loading_source_first_map
node tools/vibit check change implement-pitaya-aligned-component-discovery-module-loading-source-first-map --json
node tools/vibit inspect next --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check schemas --json
node tools/vibit check memory --json
node tools/vibit check all --json
git diff --check
```

Notes:

- Existing unrelated warning remains expected: `runtime.identity_boundary` in `runtime/internal/platform/persistence/postgres/authentication_repository.go`.
- After verification succeeds, commit and push W-0287, then pause at W-0288 per maintainer instruction.

Result:

- `node -c tools/vibit`: passed.
- `node tools/vibit inspect pitaya-component-loading --json`: passed.
- `node tools/vibit inspect rule runtime.pitaya_aligned_component_discovery_module_loading_source_first_map`: passed.
- `node tools/vibit check change implement-pitaya-aligned-component-discovery-module-loading-source-first-map --json`: passed with 13 passed, 0 warnings, 0 failures.
- `node tools/vibit inspect next --json`: passed and reports W-0288 next-ready.
- `node tools/vibit check work --json`: passed with 1740 passed, 0 warnings, 0 failures.
- `node tools/vibit check runtime --json`: passed with 26926 passed, 1 existing warning, 0 failures.
