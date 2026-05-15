# Impact

Affected areas:

- Runtime application authentication helper package.
- Runtime authentication checks.
- Architecture manifests.
- Authentication module guidance and metadata.
- Work-item continuation state.
- Conversation memory.

No public command, query, event, error catalog, permission catalog, Protobuf message, WebSocket carrier, repository interface, SQL migration, persistence adapter, startup path, external dependency, or production authentication behavior changes.

The helper computes digest bytes that future application-owned authentication behavior can pass to repositories or verifier comparison helpers. This change does not store those digests, compare those digests, select accounts, validate proofs, or expose protocol behavior.

## Boundary Notes

- Lookup digest equality remains record selection only.
- Verifier digest comparison remains deferred.
- Authentication service orchestration remains deferred.
- Repository handoff remains digest-only and storage-neutral.
- Digest bytes, raw material, HMAC input, HMAC output, verifier keys, and full key set ids are not log-safe.
