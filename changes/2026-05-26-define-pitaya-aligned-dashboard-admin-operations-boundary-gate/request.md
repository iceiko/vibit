# Request

Continue advancing toward Pitaya alignment from `W-0277 Define Pitaya-aligned dashboard and admin operations boundary gate`.

This slice is gate-only. It must define dashboard/admin operations vocabulary and source-first mapping after the metrics and tracing source-first map and must not add runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, protocol changes, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## RED Checks

```text
node tools/vibit inspect rule runtime.pitaya_aligned_dashboard_admin_operations_boundary_gate
# Unknown rule_id: runtime.pitaya_aligned_dashboard_admin_operations_boundary_gate

node tools/vibit check change define-pitaya-aligned-dashboard-admin-operations-boundary-gate --json
# change directory does not exist: changes/define-pitaya-aligned-dashboard-admin-operations-boundary-gate
```
