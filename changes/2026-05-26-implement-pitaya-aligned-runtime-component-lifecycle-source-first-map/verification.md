# Verification

Status: Final verification complete

## RED Checks

```bash
node tools/vibit inspect rule runtime.pitaya_aligned_runtime_component_lifecycle_source_first_map
# Unknown rule_id: runtime.pitaya_aligned_runtime_component_lifecycle_source_first_map

node tools/vibit inspect pitaya-component-lifecycle --json
# Unknown command.

node tools/vibit check change implement-pitaya-aligned-runtime-component-lifecycle-source-first-map --json
# change directory does not exist
```

## Final Checks

```bash
node -c tools/vibit
# passed

node tools/vibit inspect rule runtime.pitaya_aligned_runtime_component_lifecycle_source_first_map
# passed; rule_id runtime.pitaya_aligned_runtime_component_lifecycle_source_first_map is registered

node tools/vibit inspect pitaya-component-lifecycle --json
# passed; status source_first_pitaya_aligned_runtime_component_lifecycle_map
# next_ready: M-210 / W-0282 / select_next_pitaya_aligned_direction_after_runtime_component_lifecycle_map

node tools/vibit inspect next --json
# passed; status ready
# current_milestone: M-210 Select Next Pitaya-Aligned Direction After Runtime Component Lifecycle Map
# next_ready: W-0282 Select next Pitaya-aligned direction after runtime component lifecycle map

node tools/vibit check change implement-pitaya-aligned-runtime-component-lifecycle-source-first-map --json
# passed: 13 passed, 0 warnings, 0 failures

node tools/vibit check work --json
# passed: 1704 passed, 0 warnings, 0 failures

node tools/vibit check runtime --json
# passed: 25898 passed, 1 warning, 0 failures
# warning is the existing runtime.identity_boundary credential dependency warning

node tools/vibit check memory --json
# passed: 4988 passed, 0 warnings, 0 failures

node tools/vibit check schemas --json
# passed: 5337 passed, 0 warnings, 0 failures

git diff --check
# passed
```

## Final Aggregate Check

```bash
node tools/vibit check all --json
# passed: 336 subchecks passed, 1 warning, 0 failures
# warning is the existing runtime.identity_boundary credential dependency warning
```
