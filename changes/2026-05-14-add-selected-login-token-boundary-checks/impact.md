# Impact

This change adds a repository check standard, ADR, rule catalog entry, runtime check logic, and architecture metadata.

Runtime impact:

- No Go runtime behavior changes.
- No authentication implementation.
- No token generation, parsing, validation, revocation, rotation, cleanup, or storage.
- No credential lookup.
- No session persistence.
- No runtime player handlers or WebSocket routes.

Data impact:

- No migrations.
- No credential, token, external identity, runtime session, or audit tables.
- No repository interfaces.
- No PostgreSQL adapters.
- Player account lifecycle storage remains unchanged.

Protocol impact:

- No Protobuf source changes.
- No generated Protobuf changes.
- No generated authentication contract shapes.
- No WebSocket handshake changes.
- No request proof payload wire implementation.

Agent impact:

- Future agents get a specific machine-readable check rule for selected login/token boundary violations.
- The new rule keeps default verification static and local, with no live external-service requirement.
