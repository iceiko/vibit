# Impact

Affected areas:

- Runtime application authentication helper package.
- Runtime authentication checks.
- Architecture manifests.
- Authentication module guidance and metadata.
- Work-item continuation state.
- Conversation memory.

No public command, query, event, error catalog, permission catalog, Protobuf message, WebSocket carrier, repository interface, SQL migration, persistence adapter, startup path, external dependency, or production authentication behavior changes.

The helper compares verifier digest bytes that future application-owned authentication behavior can use after repository record selection and lifecycle checks. This change does not select records, inspect accounts, validate proofs, map public failures, or expose protocol behavior.

## Boundary Notes

- Lookup digest equality remains record selection only.
- Mismatch is returned as a redacted non-match result, not as public authentication behavior.
- Authentication service orchestration remains deferred.
- Repository handoff remains digest-only and storage-neutral.
- Digest bytes, raw material, HMAC input, HMAC output, verifier keys, and full key set ids are not log-safe.
