# Request

Implement the bounded `LogoutAccessToken` behavior selected after the logout access-token behavior gate.

Scope:

- Application service only: `runtime/internal/app/authentication`.
- Revoke only the verified presented opaque access token.
- Reject missing or malformed token proof before opening a unit of work.
- Use existing lookup digest, verifier digest, and constant-time comparison helpers.
- Use the existing authentication repository `FindTokenByLookupDigest` and `RevokeToken` capabilities through the application unit-of-work boundary.
- Return revoked success only after unit-of-work commit.

Out of scope:

- Runtime session revocation.
- Active WebSocket close behavior.
- Connection registry behavior.
- Reconnect, resume, duplicate replacement, or connection epoch behavior.
- Protobuf logout route or protocol session carrier.
- WebSocket handshake authentication or transport credential carriers.
- Refresh, logout-all, admin revocation, cleanup jobs, dependencies, memory durable session behavior, broader game backend modules, or direct Nakama/Pitaya API compatibility.
