# Request

Define the token verifier record schema boundary for the selected opaque access-token posture.

The boundary must ratify token verifier record ownership, non-plaintext verifier storage, statuses, expiration, revocation, credential-token linkage, retention, cleanup, replay-sensitive failure classes, redaction rules, and required future migration/repository/adapter gates.

Do not add migration source, tables, repository interfaces, PostgreSQL adapters, runtime token validation, token issuance, logout, refresh, cleanup jobs, handlers, routes, generated output, Protobuf messages, WebSocket behavior, authentication dependencies, or authentication implementation.

Preserve PostgreSQL as the default durable target, keep refresh-token storage and runtime session persistence deferred, and keep player account lifecycle tables token-free.
