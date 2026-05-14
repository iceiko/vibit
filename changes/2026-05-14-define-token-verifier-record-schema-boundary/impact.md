# Impact

This change ratifies a schema boundary only.

Runtime impact:

- No Go runtime behavior changes.
- No authentication implementation.
- No token issuance.
- No token validation.
- No logout behavior.
- No refresh behavior.
- No cleanup job.
- No runtime player handlers.
- No WebSocket routes.

Data impact:

- No migrations.
- No credential, token, external identity, runtime session, or audit tables.
- No repository interfaces.
- No PostgreSQL adapters.
- Future logical token verifier table target is ratified as `authentication_access_tokens`.

Protocol impact:

- No Protobuf source changes.
- No generated Protobuf output.
- No generated authentication contract shapes.
- No WebSocket handshake changes.
- No proof carrier implementation.

Architecture impact:

- `runtime.authentication` owns the future token verifier record boundary.
- Player account lifecycle tables remain credential-free, token-free, external-identity-free, session-free, WebSocket-state-free, and request-validation-free.
- Credential record and token verifier schema boundaries are now both ratified for M-014.
- W-0076 becomes the next ready work item for authentication schema migration queue planning.

Reference impact:

- Nakama remains a capability reference for token lifecycle, expiration, refresh, revocation, and logout pressure.
- Pitaya remains a vocabulary reference for sessions and handler context.
- Neither reference governs vibit's schema or public API shape.
