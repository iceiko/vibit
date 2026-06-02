# Request

Continue advancing toward Pitaya alignment from `W-0280 Define Pitaya-aligned runtime component lifecycle boundary gate`.

This slice is gate-only. It must define runtime component lifecycle vocabulary and current source-first mapping and must not add runtime component lifecycle behavior, handler registration behavior, component discovery or loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol changes, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## RED Checks

```text
node tools/vibit inspect rule runtime.pitaya_aligned_runtime_component_lifecycle_boundary_gate
# Unknown rule_id: runtime.pitaya_aligned_runtime_component_lifecycle_boundary_gate

node tools/vibit check change define-pitaya-aligned-runtime-component-lifecycle-boundary-gate --json
# change directory does not exist: changes/define-pitaya-aligned-runtime-component-lifecycle-boundary-gate
```
