# Request

## Original Request

Continue advancing the project until no safe work item remains. The current next-ready work item is `W-0031 Define player account and session contract standard`.

## Clarified Requirement

Define a contract standard for player account lifecycle and runtime session validation before implementation. The standard must explicitly map Nakama and Pitaya reference patterns while preserving vibit's Agent-Native boundaries.

## User-Visible Outcome

Future agents can ratify player account and session contracts in small steps without silently choosing authentication scheme, token format, credential storage, database schema, session persistence, Protobuf envelope, or WebSocket handshake behavior.

## Non-Goals

- Do not implement production authentication.
- Do not choose login methods.
- Do not choose token format, signing, refresh, expiration, or revocation behavior.
- Do not add credential storage, password hashing, cryptography, OAuth, OIDC, or external identity provider dependencies.
- Do not add player account database schema or migrations.
- Do not add session persistence.
- Do not change Protobuf envelope shape.
- Do not change WebSocket handshake behavior.
- Do not copy Nakama or Pitaya public API shape.

## Acceptance Criteria

- Add an English player account/session contract standard.
- Add a Simplified Chinese translation.
- Record adopted, adapted, deferred, and rejected Nakama/Pitaya patterns.
- Update architecture manifests so agents can discover the standard.
- Complete `W-0031`.
- Add exactly one conservative next-ready work item if safe.
- Run repository verification.
