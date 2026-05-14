# Request

## Original Request

它区域限制呗。按照你的建议和判断推进10步，除非有非常必要的，需要我决策的，再停下来问。

## Clarified Requirement

Advance `W-0080` by defining the storage-neutral authentication repository interface boundary after credential and token verifier migration sources and static checks exist.

## User-Visible Outcome

Maintainers and agents can inspect the authentication repository interface at `runtime/internal/modules/authentication/repository.go`.

The boundary declares storage-neutral credential and token verifier record shapes, mutation/query shapes, repository methods, and normalization helpers for future persistence adapters.

## Non-Goals

- Do not add PostgreSQL authentication adapters.
- Do not implement runtime credential lookup, token issuance, token validation, logout behavior, refresh behavior, cleanup jobs, handlers, routes, or production authentication behavior.
- Do not add Protobuf messages, generated authentication shapes, WebSocket proof carriers, or WebSocket handshake authentication.
- Do not add JWT, signing, key-management, Redis-like, OAuth, OIDC, provider SDK, password-hashing, S3, MinIO, or other authentication dependencies.
- Do not change the ratified credential or token verifier migration schemas.
- Do not copy Nakama or Pitaya public API shapes.

## Unknowns

- PostgreSQL authentication adapter boundaries and implementation remain future work items.
- Runtime authentication behavior remains deferred until a later milestone explicitly authorizes it.
- Token verifier algorithms, pepper/key handling, and constant-time comparison remain future implementation details.

## Acceptance Criteria

- [x] Add storage-neutral authentication repository interface source under `runtime/internal/modules/authentication/`.
- [x] Add focused tests for repository interface shape, status closed sets, mutation normalization, digest copying, and UTC timestamp normalization.
- [x] Add `modules/authentication/` manifest and agent guides with English canonical text and Simplified Chinese translation.
- [x] Update runtime checks so this specific repository interface boundary is allowed while adapters and runtime authentication behavior remain blocked.
- [x] Update work queue state so `W-0080` is completed and `W-0081` becomes next ready.
