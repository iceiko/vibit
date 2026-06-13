# Verification

Status: Verified.

## RED

- `node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_startup_shutdown_hook_map`: failed with unknown rule before implementation.
- `node tools/vibit check change select-next-pitaya-aligned-direction-after-startup-shutdown-hook-map --json`: failed because the change directory did not exist before implementation.
- `node tools/vibit inspect next --json`: reported `W-0291` as next-ready before implementation.

## GREEN

- `node -c tools/vibit`: passed.
- `node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_startup_shutdown_hook_map`: passed; the rule is registered in `rules/check-rules.json`.
- `node tools/vibit inspect next --json`: passed; current milestone is `M-220` and next-ready is `W-0292`.
- `node tools/vibit check change select-next-pitaya-aligned-direction-after-startup-shutdown-hook-map --json`: passed; 13 passed, 0 warnings, 0 failures.
- `node tools/vibit check work --json`: passed; 1764 passed, 0 warnings, 0 failures.
- `node tools/vibit check runtime --json`: passed; 27602 passed, 1 warning, 0 failures. The warning is the existing `runtime.identity_boundary` credential-dependency warning for `runtime/internal/platform/persistence/postgres/authentication_repository.go`.
- `node tools/vibit check schemas --json`: passed; 5558 passed, 0 warnings, 0 failures.
- `node tools/vibit check memory --json`: passed; 5228 passed, 0 warnings, 0 failures.
- `node tools/vibit check all --json`: passed; 346 subchecks passed, 1 warning, 0 failures.
- `git diff --check`: passed.
