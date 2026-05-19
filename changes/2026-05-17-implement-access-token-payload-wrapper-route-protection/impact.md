# Impact

Runtime impact:

- Adds `vibit.authentication.v1.AuthenticatedRequest` as an explicit request-level Protobuf payload wrapper.
- Adds generated Go Protobuf output under `runtime/internal/generated/proto/`.
- Adds application-owned route policy and route access-token validation handoff.
- Adds Protobuf adapter wrapper handling so protected routes validate proof before domain dispatch.
- Adds focused app, protocol adapter, and WebSocket transport tests.

Architecture impact:

- Keeps the existing Protobuf envelope unchanged.
- Keeps WebSocket transport credential-neutral.
- Keeps route protection in the application/protocol boundary instead of transport.
- Preserves session persistence and WebSocket handshake authentication as separate future decisions.
- Moves the work queue to a blocked next-direction confirmation gate.

Data impact:

- No migrations are added or changed.
- No repository interfaces are changed.
- No PostgreSQL adapters are changed.

Compatibility impact:

- Existing envelope route fields remain the domain route.
- Protected routes require the wrapper payload posture when route protection is configured.
- Public device credential login route remains explicit.

Security impact:

- Access-token proof must be explicit request payload material, not metadata-only identity.
- Protected route dispatch requires validated player identity.
- Public invalid/malformed/unavailable errors remain collapsed and redacted.
- Raw token values are not stored in repositories by this change.
