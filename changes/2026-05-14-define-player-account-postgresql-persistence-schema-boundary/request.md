# Request

## Original Request

Proceed according to professional recommendation after selecting player account PostgreSQL persistence schema ratification as the next direction.

## Clarified Requirement

Define the player account PostgreSQL persistence schema boundary before adding SQL migration source, repository adapter code, runtime player handlers, authentication, token behavior, credential storage, or session persistence.

## User-Visible Outcome

Future agents should be able to inspect the standards and understand exactly what the first player account migration may create and what it must not include.

## Non-Goals

- Do not add SQL migration source in this work item.
- Do not add player account repository implementation.
- Do not add runtime player handlers or WebSocket routes.
- Do not implement authentication.
- Do not choose login methods.
- Do not add token behavior.
- Do not add credential storage or password hashing.
- Do not add session persistence.
- Do not change Protobuf envelope or WebSocket handshake behavior.
- Do not copy Nakama or Pitaya public API shapes.
- Do not make metadata-only identity sufficient for production permissions.

## Acceptance Criteria

- [x] Player account lifecycle tables are named and owned.
- [x] Required columns, constraints, indexes, and audit/event records are documented.
- [x] Explicit exclusions prevent credentials, tokens, sessions, external identity links, WebSocket state, and inventory state from entering the account schema.
- [x] Module/runtime manifests mark schema as ratified while migration and implementation remain deferred.
- [x] Runtime checks allow future player account migrations only after schema ratification.
- [x] Verification is recorded.
