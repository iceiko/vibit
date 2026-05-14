# Request

Advance `W-0083 Refine authentication adapter implementation checks`.

The checks must allow the separately ratified PostgreSQL adapter boundary to use authentication credential and token verifier persistence vocabulary inside `runtime/internal/platform/persistence/postgres/`, while continuing to block runtime authentication, token generation, token validation, verifier comparison, bearer parsing, WebSocket behavior, Protobuf behavior, generated authentication shapes, transaction ownership, and major authentication dependencies.
