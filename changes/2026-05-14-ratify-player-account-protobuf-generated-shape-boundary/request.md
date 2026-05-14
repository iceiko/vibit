# Request

## Original Request

Continue advancing the project. The maintainer clarified that professional implementation steps should proceed without separate confirmation unless they require a product, security, persistence, or compatibility decision.

## Clarified Requirement

Ratify the Protobuf source and generated Go Protobuf shape for the already-ratified player account semantic contracts.

This change may add player account Protobuf source and generated Protobuf output because those are the next required shape artifacts for the existing semantic contracts.

## User-Visible Outcome

The player account contracts will have an explicit wire schema and generated Go Protobuf package that future protocol adapter work can inspect.

## Non-Goals

- Do not implement player runtime handlers.
- Do not add WebSocket routes for player account commands or queries.
- Do not change the Protobuf envelope.
- Do not change WebSocket handshake behavior.
- Do not implement authentication.
- Do not choose login methods.
- Do not choose token format, signing, refresh, expiration, or revocation behavior.
- Do not add credential storage, password hashing, cryptography, OAuth, OIDC, or external identity provider dependencies.
- Do not add player account database schema, migrations, indexes, repositories, or session persistence.
- Do not copy Nakama or Pitaya public API shape.

## Acceptance Criteria

- Add `proto/vibit/player/v1/player.proto`.
- Generate `runtime/internal/generated/proto/vibit/player/v1/player.pb.go` from the Protobuf source using the accepted Buf generation path.
- Update player module and architecture manifests to record the ratified wire shape.
- Update relevant standards and translations.
- Update `runtime.identity_boundary` checks so the newly ratified player wire source is allowed while runtime handlers, authentication, and persistence remain blocked.
- Complete `W-0033`.
- Add the next conservative work item if safe.
- Run protocol, generated-output, runtime, work, change, and repository verification.
