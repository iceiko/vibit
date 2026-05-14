# Request

## Original Request

它区域限制呗。按照你的建议和判断推进10步，除非有非常必要的，需要我决策的，再停下来问。

## Clarified Requirement

Advance `W-0077` by adding the SQL-first PostgreSQL migration source for the ratified `authentication_device_credentials` schema.

## User-Visible Outcome

Maintainers and agents can inspect the first authentication credential persistence schema source under `runtime/migrations/postgres/`.

The migration creates:

- `authentication_device_credentials`: durable device credential verifier records for the selected first login posture.

## Non-Goals

- Do not add token verifier tables.
- Do not add authentication repository interfaces.
- Do not add PostgreSQL authentication adapters.
- Do not implement credential lookup, login, token issuance, token validation, logout, refresh, cleanup, handlers, routes, or production authentication behavior.
- Do not add password hashing, JWT, signing, key-management, Redis-like, OAuth, OIDC, provider SDK, or other authentication dependencies.
- Do not add Protobuf messages, generated authentication shapes, WebSocket proof carriers, or WebSocket handshake authentication.
- Do not add generic metadata columns or store raw credential material.
- Do not copy Nakama or Pitaya public API shapes.

## Unknowns

- The exact verifier algorithm and secret configuration remain future implementation work.
- Live PostgreSQL apply/rollback is not required by this source-only work item and depends on an explicit disposable DSN.
- Token verifier schema remains a separate next work item.

## Acceptance Criteria

- [x] Add deterministic SQL-first goose migration source `runtime/migrations/postgres/000003_create_authentication_device_credentials.sql`.
- [x] Include `-- +goose Up`, `-- Module: runtime.authentication`, and `-- +goose Down`.
- [x] Create only `authentication_device_credentials`.
- [x] Include explicit constraints and indexes for ratified credential record semantics.
- [x] Preserve player account lifecycle table separation.
- [x] Update runtime manifests, standards, guides, and checks to record credential migration source state while token migration, repositories, adapters, runtime authentication, Protobuf, WebSocket, and dependencies remain deferred.
