# Impact

## Runtime

This change selects the next runtime direction only. It authorizes a follow-up bounded implementation work item for the PostgreSQL adapter that implements `runtime/internal/app/session.Repository`.

## Authentication

Authentication remains token-linkage-only for runtime sessions. This direction does not change access-token validation, login behavior, logout behavior, or token revocation behavior.

## Protocol And Transport

The WebSocket transport remains credential-neutral and the existing Protobuf envelope remains unchanged.

## References

Nakama is used as the durable session lifecycle reference. Pitaya is used as the separation-of-concerns reference for transport, session context, and handler routing.
