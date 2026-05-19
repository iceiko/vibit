# Plan

1. Close `M-061/W-0133`.
2. Select `define_session_repository_boundary`.
3. Create `M-062/W-0134` as the next bounded gate-only slice.
4. Preserve repository implementation, adapter, runtime behavior, route policy, logout/revocation, reconnect, dependency, memory durable session, direct Nakama/Pitaya API compatibility, and broader game backend deferrals.
5. Verify change, work, runtime, and repository checks.

## Rollback

Reversal means reopening `M-061/W-0133`, removing the selected direction, and removing `M-062/W-0134` before any session repository boundary depends on it.
