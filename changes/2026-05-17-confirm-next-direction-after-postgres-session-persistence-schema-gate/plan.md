# Plan

1. Close `M-059/W-0131`.
2. Select `implement_runtime_sessions_migration_source`.
3. Create `M-060/W-0132` as the next bounded implementation slice.
4. Preserve repository, adapter, runtime behavior, route policy, logout/revocation, reconnect, dependency, memory durable session, direct Nakama/Pitaya API compatibility, and broader game backend deferrals.
5. Verify change, work, runtime, and repository checks.

## Rollback

Reversal means reopening `M-059/W-0131`, removing the selected direction, and removing `M-060/W-0132` before any migration-source implementation depends on it.
