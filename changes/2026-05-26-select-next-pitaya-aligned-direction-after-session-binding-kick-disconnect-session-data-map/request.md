# Request

Continue advancing toward Pitaya alignment from `W-0270 Select next Pitaya-aligned direction after session binding, kick/disconnect, and session data map`.

This slice is selection-only. It must choose the next bounded Pitaya-aligned direction after the source-first session lifecycle map and must not add runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, session binding behavior, kick/disconnect behavior, session data behavior or persistence, protocol changes, generated output, persistence, dependencies, hosted surfaces, SDKs, or direct Nakama/Pitaya API compatibility.

## RED Checks

```text
node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_session_binding_kick_disconnect_session_data_map
# Unknown rule_id: runtime.next_pitaya_aligned_direction_after_session_binding_kick_disconnect_session_data_map

node tools/vibit check change select-next-pitaya-aligned-direction-after-session-binding-kick-disconnect-session-data-map --json
# change directory does not exist
```
