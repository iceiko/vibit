# Impact

Runtime impact:

- No Go code is added by the gate-only step.
- The gate authorizes a future implementation to expose the existing `AuthenticateWithDeviceCredential` service method through the WebSocket Protobuf request loop.

Architecture impact:

- Public login route registration ownership is assigned to application composition.
- Protocol payload bridging is assigned to the Protobuf adapter.
- The first composed runtime store remains limited to `VIBIT_RUNTIME_STORE=postgres`.
- The memory runtime store remains a bootstrap path where durable login route behavior is unavailable.
- Nakama informs the authenticate-before-gameplay sequence.
- Pitaya informs transport/session/route/handler separation.

Data impact:

- No migrations are added or changed.
- No repository interfaces are changed.
- No PostgreSQL adapters are changed.

Compatibility impact:

- Existing Protobuf envelope remains unchanged.
- WebSocket handshake remains unchanged.
- Existing protected-route access-token wrapper remains unchanged.

Security impact:

- Credential proof and returned access-token text remain redacted.
- Public login failures must not reveal lookup hits, player existence, verifier mismatch, key ids, or token internals.
- The authentication service remains the owner of token generation and unit-of-work behavior.
