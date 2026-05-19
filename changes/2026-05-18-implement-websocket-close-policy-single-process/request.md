# Request

Continue the next ten steps after the WebSocket close policy gate, using Nakama and Pitaya as important references and absorbing their design strengths.

The selected direction is:

```text
implement_websocket_close_policy_single_process
```

This change implements only the application-owned single-process close policy primitive. It does not close concrete WebSocket sockets, add transport close handoff code, choose close codes or close reason text, add logout protocol routes, add reconnect behavior, add protocol session carriers, or copy Nakama/Pitaya public APIs.
