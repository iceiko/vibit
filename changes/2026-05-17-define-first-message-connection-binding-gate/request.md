# Request

Define the first-message connection binding gate after session persistence and WebSocket handshake ratification.

This change selects the future connection binding shape at the standards level only:

```text
runtime.authentication.BindConnection
```

The future binding message is an application/protocol message carried inside the existing Protobuf WebSocket frame loop. It is not a WebSocket handshake credential carrier and it does not change the existing Protobuf envelope.

This change must not implement connection binding, add `.proto` messages, generate Protobuf output, add connection registries, change route protection behavior, add session persistence, add migrations, change repositories, add dependencies, implement logout/revocation, implement reconnect/epoch behavior, or adopt direct Nakama/Pitaya public API compatibility.
