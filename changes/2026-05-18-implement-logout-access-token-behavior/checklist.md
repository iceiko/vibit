# Checklist

- [x] Implemented `Service.LogoutAccessToken`.
- [x] Rejected missing and malformed proof before unit of work.
- [x] Used lookup digest before repository lookup.
- [x] Required active access-token posture before revocation.
- [x] Compared verifier digest before revocation.
- [x] Called `RevokeToken` with `logout_presented_access_token`.
- [x] Returned success only after commit.
- [x] Collapsed public invalid-token failures.
- [x] Kept raw token text out of result and error strings.
- [x] Added focused tests for success, rejection, invalid token posture, verifier mismatch, dependency failures, and no session/player side effects.
- [x] Preserved WebSocket, Protobuf, session revocation, connection registry, reconnect, refresh, logout-all, admin revocation, cleanup, dependency, memory durable session behavior, broader module, and direct Nakama/Pitaya compatibility deferrals.
