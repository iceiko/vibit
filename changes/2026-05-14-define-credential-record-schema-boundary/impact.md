# Impact

This change ratifies a schema boundary only.

Runtime impact:

- No Go runtime behavior changes.
- No authentication implementation.
- No credential lookup.
- No token behavior.
- No runtime player handlers.
- No WebSocket routes.

Data impact:

- No migrations.
- No credential, token, external identity, runtime session, or audit tables.
- No repository interfaces.
- No PostgreSQL adapters.
- Future logical credential table target is ratified as `authentication_device_credentials`.

Protocol impact:

- No Protobuf source changes.
- No generated Protobuf output.
- No generated authentication contract shapes.
- No WebSocket handshake changes.
- No proof carrier implementation.

Architecture impact:

- `runtime.authentication` owns the future credential record boundary.
- Player account lifecycle tables remain credential-free, token-free, external-identity-free, session-free, WebSocket-state-free, and request-validation-free.
- W-0075 becomes the next ready work item for token verifier record schema boundary.

Reference impact:

- Nakama remains a capability reference for device-style login and session vocabulary.
- Pitaya remains a vocabulary reference for sessions and handler context.
- Neither reference governs vibit's schema or public API shape.
