# Impact Analysis

## Affected Modules

- Runtime application session package: adds `runtime/internal/app/session`.
- Authentication module: manifest and guide references clarify that it remains token-record-linkage only.

## Module Ownership Impact

`runtime/internal/app/session` now owns the storage-neutral session repository interface vocabulary. PostgreSQL adapter ownership remains deferred to `runtime/internal/platform/persistence/postgres` behind a separate gate.

## Public Contract Impact

No public commands, queries, events, errors, permissions, Protobuf messages, or envelope fields are changed.

## Data And Migration Impact

No migration is added or changed. The interface references the existing `runtime_sessions` lifecycle concept only as application value vocabulary and does not execute SQL.

## Test Impact

Focused Go tests are added for:

- Repository interface storage-neutral conformance.
- Session status closed vocabulary.
- Create mutation normalization.
- Required field, actor-kind, actor/player mismatch, status, and lifetime rejection.
- Runtime session record normalization.
- Query and mutation normalization for lookup, active lookup, last-seen, expiration, revocation, and active-player listing.
- Absence of raw token, credential, digest, verifier key, WebSocket, or connection-state fields.

## Documentation Impact

ADR-0062, change specs, conversation memory, manifests, AGENTS guides, and repository checks are updated.

## Compatibility Risks

Runtime behavior is unchanged. The interface may need extension before the PostgreSQL adapter, but the chosen shape matches the already ratified `runtime_sessions` table and the prior boundary's candidate capability vocabulary.
