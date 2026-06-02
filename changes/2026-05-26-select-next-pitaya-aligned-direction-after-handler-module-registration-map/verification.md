# Verification

Status: pending final command run during implementation.

RED evidence:

- `node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_handler_module_registration_map` failed with `Unknown rule_id: runtime.next_pitaya_aligned_direction_after_handler_module_registration_map`.
- `node tools/vibit check change select-next-pitaya-aligned-direction-after-handler-module-registration-map --json` failed because the change directory did not exist.

Expected verification:

- `node -c tools/vibit`
- `node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_handler_module_registration_map`
- `node tools/vibit check change select-next-pitaya-aligned-direction-after-handler-module-registration-map --json`
- `node tools/vibit check work --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check all --json`
- `git diff --check`
