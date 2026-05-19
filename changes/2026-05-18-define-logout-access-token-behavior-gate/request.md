# Request

Define the gate-only boundary for future `LogoutAccessToken` behavior.

The selected posture is `presented_access_token_only`: a later implementation may revoke only the verified presented opaque access token and must not revoke runtime sessions, close active WebSocket connections, add protocol logout routes, add refresh behavior, add logout-all behavior, or copy Nakama/Pitaya APIs.
