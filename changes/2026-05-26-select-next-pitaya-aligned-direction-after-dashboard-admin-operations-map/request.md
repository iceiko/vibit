# Request

Continue advancing toward Pitaya alignment from `W-0279 Select next Pitaya-aligned direction after dashboard/admin operations map`.

This slice is selection-only. It must choose the next bounded Pitaya-aligned direction after the source-first dashboard/admin operations map and must not add runtime component lifecycle behavior, handler registration behavior, component discovery or loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol changes, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## RED Checks

```text
node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_dashboard_admin_operations_map
# Unknown rule_id: runtime.next_pitaya_aligned_direction_after_dashboard_admin_operations_map

node tools/vibit check change select-next-pitaya-aligned-direction-after-dashboard-admin-operations-map --json
# change directory does not exist: changes/select-next-pitaya-aligned-direction-after-dashboard-admin-operations-map
```
