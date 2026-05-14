# Request

## Original Request

它区域限制呗。按照你的建议和判断推进10步，除非有非常必要的，需要我决策的，再停下来问。

## Clarified Requirement

Advance `W-0078` by adding the SQL-first PostgreSQL migration source for the ratified `authentication_access_tokens` schema.

## User-Visible Outcome

Maintainers and agents can inspect the first token verifier persistence schema source under `runtime/migrations/postgres/`.

The migration creates:

- `authentication_access_tokens`: durable opaque access-token verifier records for the selected first token posture.

## Non-Goals

- Do not change the credential migration schema.
- Do not add authentication repository interfaces.
- Do not add PostgreSQL authentication adapters.
- Do not implement token issuance, token validation, logout, refresh, cleanup, credential lookup, handlers, routes, or production authentication behavior.
- Do not add JWT, signing, key-management, Redis-like, OAuth, OIDC, provider SDK, password-hashing, or other authentication dependencies.
- Do not add Protobuf messages, generated authentication shapes, WebSocket proof carriers, or WebSocket handshake authentication.
- Do not add generic metadata columns or store raw token material.
- Do not copy Nakama or Pitaya public API shapes.

## Unknowns

- The exact token verifier algorithm, pepper/secret configuration, and constant-time comparison behavior remain future implementation work.
- Live PostgreSQL apply/rollback is not required by this source-only work item and depends on an explicit disposable DSN.
- Authentication migration static checks remain a separate next work item.

## Acceptance Criteria

- [x] Add deterministic SQL-first goose migration source `runtime/migrations/postgres/000004_create_authentication_access_tokens.sql`.
- [x] Include `-- +goose Up`, `-- Module: runtime.authentication`, and `-- +goose Down`.
- [x] Create only `authentication_access_tokens`.
- [x] Include explicit constraints and indexes for ratified token verifier record semantics.
- [x] Preserve player account lifecycle table separation.
- [x] Preserve credential-token linkage without changing the credential migration.
- [x] Update runtime manifests, standards, guides, and checks to record token verifier migration source state while repositories, adapters, runtime authentication, Protobuf, WebSocket, and dependencies remain deferred.
