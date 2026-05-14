# Request

## Original Request

Continue advancing the project. The maintainer clarified that agents should proceed on professional implementation and architecture sequencing details without stopping for unnecessary confirmation, and should only stop for real maintainer decisions.

## Clarified Requirement

Ratify the minimal runtime session validation semantic contracts for the existing application-owned session validation boundary.

This step should define machine-readable contract source files and registry entries for session validation without implementing production authentication, token validation, credential storage, session persistence, player account persistence, WebSocket handshake authentication, or player runtime handlers.

## User-Visible Outcome

Future agents can inspect runtime session validation contracts before changing the application session validator hook or adding real authentication/session behavior.

## Non-Goals

- Do not choose token format, signing, refresh, expiration, revocation, or token storage behavior.
- Do not choose concrete login methods.
- Do not add credential storage, password hashing, OAuth, OIDC, JWT, or external identity provider dependencies.
- Do not add session persistence, player account persistence, database schema, migrations, or indexes.
- Do not change Protobuf envelope shape.
- Do not change WebSocket handshake authentication behavior.
- Do not add runtime player account handlers or WebSocket player routes.
- Do not treat metadata-only `player_id` or `session_id` as authenticated proof.
- Do not copy Nakama or Pitaya public API shape.

## Acceptance Criteria

- Add runtime session validation semantic contract sources under a runtime-owned path.
- Register the runtime session contracts in `.arch/contracts.yaml`.
- Update runtime and account/session standards and Simplified Chinese translations.
- Update repository checks so runtime session contracts are machine-checkable without pretending they are domain module contracts.
- Complete `W-0034`.
- Preserve all ask-first boundaries for authentication, token, credential, persistence, protocol envelope, and WebSocket handshake decisions.
- Run repository verification.
