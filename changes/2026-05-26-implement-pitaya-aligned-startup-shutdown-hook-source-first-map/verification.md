# Verification

Status: Verified

Required commands:

```text
node -c tools/vibit
node tools/vibit inspect pitaya-startup-shutdown --json
node tools/vibit inspect rule runtime.pitaya_aligned_startup_shutdown_hook_source_first_map
node tools/vibit inspect next --json
node tools/vibit check change implement-pitaya-aligned-startup-shutdown-hook-source-first-map --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check schemas --json
node tools/vibit check memory --json
node tools/vibit check all --json
git diff --check
```

RED evidence:

- `node tools/vibit inspect pitaya-startup-shutdown --json` exited 2 with `Unknown command`.
- `node tools/vibit inspect rule runtime.pitaya_aligned_startup_shutdown_hook_source_first_map` exited 1 with `Unknown rule_id: runtime.pitaya_aligned_startup_shutdown_hook_source_first_map`.
- `node tools/vibit check change implement-pitaya-aligned-startup-shutdown-hook-source-first-map --json` exited 1 because the change directory did not exist.

Result:

- `node -c tools/vibit`: passed.
- `node tools/vibit inspect pitaya-startup-shutdown --json`: passed; output kind `pitaya_aligned_startup_shutdown_hook_inspection`, status `source_first_pitaya_aligned_startup_shutdown_hook_map`, and next-ready `W-0291`.
- `node tools/vibit inspect rule runtime.pitaya_aligned_startup_shutdown_hook_source_first_map`: passed; rule is registered in `rules/check-rules.json`.
- `node tools/vibit inspect next --json`: passed; current milestone `M-219`, next-ready `W-0291`.
- `node tools/vibit check change implement-pitaya-aligned-startup-shutdown-hook-source-first-map --json`: passed with 13 passed, 0 warnings, 0 failures.
- `node tools/vibit check work --json`: passed with 1758 passed, 0 warnings, 0 failures.
- `node tools/vibit check runtime --json`: passed with 27454 passed, 1 warning, 0 failures.
- `node tools/vibit check schemas --json`: passed with 5534 passed, 0 warnings, 0 failures.
- `node tools/vibit check memory --json`: passed with 5204 passed, 0 warnings, 0 failures.
- `node tools/vibit check all --json`: passed with 345 subchecks passed, 1 warning, 0 failures.
- `git diff --check`: passed.
