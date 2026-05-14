# Request

Plan the authentication schema migration queue now that the credential record and token verifier record schema boundaries are both ratified.

The plan must define the credential and token verifier migration sequence, repository-interface gates, PostgreSQL adapter gates, redaction checks, live PostgreSQL verification expectations, and future work-item ordering.

Do not add migration source, tables, repository interfaces, PostgreSQL adapters, runtime credential lookup, token issuance, token validation, logout, refresh, cleanup jobs, handlers, routes, generated output, Protobuf messages, WebSocket behavior, authentication dependencies, or authentication implementation.

Preserve PostgreSQL as the default durable target and keep player account lifecycle tables credential-free, token-free, external-identity-free, session-free, WebSocket-state-free, and request-validation-free.
