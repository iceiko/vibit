# Request

The maintainer asked the agent to recommend the next ten steps after the active connection registry implementation and continue, with Nakama and Pitaya as important references.

The selected direction is:

```text
define_websocket_close_policy_gate
```

This direction adds a gate-only WebSocket close policy boundary before any registry invalidation, logout, token revocation, session revocation, reconnect, or admin operation can close concrete sockets.
