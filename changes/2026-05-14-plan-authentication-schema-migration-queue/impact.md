# Impact

This change plans future work only.

Runtime impact:

- No Go runtime behavior changes.
- No authentication implementation.
- No credential lookup.
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
- Planned future migration order is credential first, then token verifier.

Protocol impact:

- No Protobuf source changes.
- No generated Protobuf output.
- No generated authentication contract shapes.
- No WebSocket handshake changes.
- No proof carrier implementation.

Architecture impact:

- M-014 now has a deterministic remaining work queue from migration source through adapter-boundary planning and milestone closeout.
- `W-0077` becomes the next ready work item.
- Runtime authentication remains outside M-014.

Reference impact:

- Nakama remains a capability reference for account authentication and token lifecycle pressure.
- Pitaya remains a vocabulary reference for sessions and handler context.
- Neither reference governs vibit's schema or public API shape.
