# Impact

This change records a queue decision only.

It completes `M-089/W-0161`, creates `M-090/W-0162` as the next bounded implementation slice, and preserves ask-first boundaries for:

- WebSocket close policy.
- Kick/disconnect behavior.
- Runtime session revocation.
- Duplicate connection replacement.
- Reconnect and epoch behavior.
- Protocol logout routes.
- Protocol session carriers.
- Operations and observability posture.
- Memory durable session behavior.
- Direct Nakama/Pitaya API compatibility.
- Broader game backend module expansion.

No Protobuf source, generated output, migration, dependency, WebSocket close behavior, reconnect behavior, or direct Nakama/Pitaya compatibility is added by this direction-selection change.
