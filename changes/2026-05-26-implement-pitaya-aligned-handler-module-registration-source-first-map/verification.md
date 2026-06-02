# Verification

Status: Verified.

RED evidence:

- `node tools/vibit inspect rule runtime.pitaya_aligned_handler_module_registration_source_first_map` exited 1 with `Unknown rule_id: runtime.pitaya_aligned_handler_module_registration_source_first_map`.
- `node tools/vibit inspect pitaya-handler-modules --json` exited 2 with `Unknown command`.
- `node tools/vibit check change implement-pitaya-aligned-handler-module-registration-source-first-map --json` exited 1 because the change directory did not exist.

Implementation verification:

- `node -c tools/vibit` passed.
- `node tools/vibit inspect pitaya-handler-modules --json` passed and reported `kind: pitaya_aligned_handler_module_registration_inspection`, `check_rule: runtime.pitaya_aligned_handler_module_registration_source_first_map`, and `next_ready_work_item: W-0285`.
- `node tools/vibit inspect next --json` passed and reported `W-0285 Select next Pitaya-aligned direction after handler module registration map` as next-ready.
- `node tools/vibit inspect rule runtime.pitaya_aligned_handler_module_registration_source_first_map` passed.
- `node tools/vibit check change implement-pitaya-aligned-handler-module-registration-source-first-map --json` passed with 13 passed, 0 warnings, and 0 failures.
- `node tools/vibit check work --json` passed with 1722 passed, 0 warnings, and 0 failures.
- `node tools/vibit check schemas --json` passed with 5404 passed, 0 warnings, and 0 failures.
- `node tools/vibit check memory --json` passed with 5060 passed, 0 warnings, and 0 failures.
- `node tools/vibit check runtime --json` passed with 26410 passed, 1 existing warning, and 0 failures.
- `node tools/vibit check all --json` passed with 339 subchecks passed, 1 existing warning, and 0 failures.
- `git diff --check` passed.

Existing warning:

- `runtime.identity_boundary` on `runtime/internal/platform/persistence/postgres/authentication_repository.go`: credential dependency boundary posture.
