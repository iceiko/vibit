# Request

Define the next gate after bound identity route policy implementation: logout/revocation active-connection behavior.

The gate must focus on Nakama and Pitaya lessons:

- Nakama shows that authenticated session material has logout, revocation, expiration, refresh, and realtime socket lifecycle pressure.
- Pitaya shows that acceptors, session context, handlers, and connection management should stay separated.

This change must not implement logout, token revocation, runtime session revocation, socket close behavior, connection registries, Protobuf routes, reconnect behavior, or direct external API compatibility.
