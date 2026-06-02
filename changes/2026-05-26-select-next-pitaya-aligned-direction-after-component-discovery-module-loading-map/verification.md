# Verification

Status: verified during implementation.

RED evidence:

- `node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_component_discovery_module_loading_map` failed with `Unknown rule_id: runtime.next_pitaya_aligned_direction_after_component_discovery_module_loading_map`.
- `node tools/vibit check change select-next-pitaya-aligned-direction-after-component-discovery-module-loading-map --json` failed because the change directory did not exist.

Expected verification:

- `node -c tools/vibit`
- `node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_component_discovery_module_loading_map`
- `node tools/vibit check change select-next-pitaya-aligned-direction-after-component-discovery-module-loading-map --json`
- `node tools/vibit inspect next --json`
- `node tools/vibit check work --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check all --json`
- `git diff --check`

Observed verification:

- `node -c tools/vibit` passed.
- `node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_component_discovery_module_loading_map` passed.
- `node tools/vibit inspect next --json` reported `M-217/W-0289` as next-ready.
- `node tools/vibit check change select-next-pitaya-aligned-direction-after-component-discovery-module-loading-map --json` passed with 13 passed, 0 warnings, 0 failures.
- `node tools/vibit check work --json` passed with `W-0289` as the only next-ready work item.
- `node tools/vibit check runtime --json` passed with 27071 passed, 1 existing warning, 0 failures.
- `node tools/vibit check memory --json` passed with 5156 passed, 0 warnings, 0 failures.
- `node tools/vibit check schemas --json` passed with 5490 passed, 0 warnings, 0 failures.
- `node tools/vibit check all --json` passed with 343 subchecks passed, 1 existing warning, 0 failures.
- `git diff --check` passed.
