# Impact

This change updates static checks only.

Runtime impact:

- No Go runtime behavior changes.
- No authentication adapter implementation yet.
- No token generation, token validation, verifier comparison, bearer parsing, login, logout execution, refresh, cleanup jobs, handlers, routes, Protobuf behavior, or WebSocket behavior.

Check impact:

- Authentication PostgreSQL adapter files under `runtime/internal/platform/persistence/postgres/` may use the ratified persistence vocabulary for `authentication_device_credentials` and `authentication_access_tokens`.
- Runtime, protocol, transport, generated output, and domain-module code remain protected by selected login/token boundaries.
- The authentication module repository interface remains storage-neutral and forbidden from SQL, pgx, transaction control, token generation, bearer parsing, and verifier comparison.

Workflow impact:

- `W-0083` becomes completed.
- `W-0084 Implement authentication PostgreSQL adapter` becomes the next ready work item.
