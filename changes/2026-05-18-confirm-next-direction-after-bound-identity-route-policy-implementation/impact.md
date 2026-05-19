# Impact

This change records a queue decision only.

It completes `M-081/W-0153`, creates `M-082/W-0154` as the next bounded gate, and preserves ask-first boundaries for:

- Logout execution.
- Token revocation execution.
- Runtime session revocation execution.
- Active connection invalidation.
- Connection registry behavior.
- WebSocket close policy.
- Reconnect and epoch behavior.
- Protocol logout routes.
- Protocol session carriers.
- Operations and observability posture.
- Memory durable session behavior.
- Direct Nakama/Pitaya API compatibility.
- Broader game backend module expansion.

No Go runtime behavior, Protobuf source, generated output, migration, dependency, or WebSocket behavior is added by this direction-selection change.
