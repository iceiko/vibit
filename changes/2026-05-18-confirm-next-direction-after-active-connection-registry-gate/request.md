# Request

The maintainer asked the agent to recommend the next ten steps after the active connection registry gate and continue, with Nakama and Pitaya as important references.

The selected direction is:

```text
implement_active_connection_registry_single_process
```

This direction adds the first application-owned, single-process, in-memory, non-durable active connection registry before any WebSocket close policy, reconnect behavior, protocol logout route, or broader game backend module expansion.
