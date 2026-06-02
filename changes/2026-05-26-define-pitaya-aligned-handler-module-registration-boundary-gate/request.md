# Request

Continue advancing toward Pitaya alignment from `W-0283 Define Pitaya-aligned handler module registration boundary gate`.

This slice is gate-only. It must define handler module registration vocabulary and current source-first mapping and must not add handler module registration behavior, handler registration behavior, dynamic handler registration, component discovery or loading, component module loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol changes, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## RED Checks

```text
node tools/vibit inspect rule runtime.pitaya_aligned_handler_module_registration_boundary_gate
# Unknown rule_id: runtime.pitaya_aligned_handler_module_registration_boundary_gate

node tools/vibit check change define-pitaya-aligned-handler-module-registration-boundary-gate --json
# change directory does not exist: changes/define-pitaya-aligned-handler-module-registration-boundary-gate
```
