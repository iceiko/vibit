# Request

Define the credential record schema boundary for the selected `device_credential_login` posture.

The boundary must ratify credential record ownership, lifecycle states, verifier semantics, uniqueness rules, player account relationship, disabled/revoked credential behavior, rotation/replacement semantics, redaction rules, and required future migration/repository/adapter gates.

Do not add migration source, tables, repository interfaces, PostgreSQL adapters, runtime lookup, handlers, routes, generated output, Protobuf messages, WebSocket behavior, authentication dependencies, or authentication implementation.

Preserve PostgreSQL as the default durable target and keep player account lifecycle tables credential-free and provider-subject-free.
