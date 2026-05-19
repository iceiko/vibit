# Impact

Runtime impact:

- `runtime/cmd/vibit-server` now wires the existing authentication service into the PostgreSQL runtime path.
- PostgreSQL startup loads verifier key configuration and fails closed when it is unavailable or invalid.
- The Protobuf frame handler receives an application route protector when PostgreSQL authentication startup composition succeeds.
- Memory startup remains a bootstrap path without authentication repository capability.

Architecture impact:

- Authentication validation remains application-owned.
- WebSocket transport remains credential-neutral.
- The Protobuf envelope remains unchanged.
- Startup composition stays in the process startup boundary.

Data impact:

- No migrations are added or changed.
- No repository interfaces are changed.
- No PostgreSQL adapters are changed.

Security impact:

- Protected routes can be enforced in the PostgreSQL runtime process path.
- Missing auth secret configuration fails startup closed for PostgreSQL.
- Startup errors avoid raw verifier key values.
- Raw access tokens, credential material, lookup digests, verifier digests, and verifier keys are not logged or stored by startup code.

Compatibility impact:

- Default memory startup remains available.
- PostgreSQL startup now requires the five `VIBIT_AUTH_*` verifier key environment variables.
- `VIBIT_AUTH_ACCESS_TOKEN_TTL` and `VIBIT_AUTH_TOKEN_AUDIENCE` are optional.
