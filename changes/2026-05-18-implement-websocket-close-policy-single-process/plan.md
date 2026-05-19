# Plan

1. Confirm the current blocked work item after the WebSocket close policy gate.
2. Select `implement_websocket_close_policy_single_process`.
3. Add the close policy source under `runtime/internal/app/connection`.
4. Add focused tests for registry-backed targeting, redaction, invalidation, skipped records, and deferrals.
5. Add ADR-0077 and conversation memory.
6. Update work items and architecture manifests.
7. Update authentication module manifests and AGENTS guides.
8. Add the repository check rule `runtime.websocket_close_policy_single_process_implementation`.
9. Run focused Go tests.
10. Run full repository verification and stop at the next confirmation gate.
