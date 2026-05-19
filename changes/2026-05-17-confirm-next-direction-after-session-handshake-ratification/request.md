# Request

Close `M-051/W-0123` by selecting the next milestone direction after session persistence and WebSocket handshake ratification.

The selected direction is:

```text
define_first_message_connection_binding_gate
```

The maintainer asked the agent to recommend the next ten steps and continue according to that recommendation, while continuing to treat Nakama and Pitaya as important reference baselines.

This change records the direction only. It does not implement first-message connection binding, add protocol messages, generate Protobuf output, add session persistence, change WebSocket handshake authentication, add migrations, change repositories, add dependencies, implement logout/revocation, implement reconnect/epoch behavior, or adopt direct Nakama/Pitaya API compatibility.
