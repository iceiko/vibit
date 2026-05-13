# Request

## Original Request

Continue ten work items unless maintainer confirmation is required.

## Clarified Requirement

Advance `W-0010` by defining the PostgreSQL repository, migration, transaction, and ownership boundaries before adding persistent inventory behavior.

## User-Visible Outcome

Maintainers and agents have a durable standard for how PostgreSQL-backed inventory persistence must be introduced:

- Domain modules own repository interfaces and invariants.
- PostgreSQL adapters own `pgx`.
- Migration tooling owns `goose`.
- SQL migrations live under `runtime/migrations/postgres/`.
- State-changing commands use an application-owned unit of work.
- Inventory grants must lock the player inventory aggregate before applying capacity-sensitive mutations.

## Non-Goals

- Do not add a PostgreSQL repository implementation.
- Do not add SQL migration files yet.
- Do not add PostgreSQL runtime configuration.
- Do not add MinIO or an S3 SDK.
- Do not change the WebSocket Protobuf protocol.
- Do not add authentication or session validation.

## Unknowns

- The exact disposable PostgreSQL integration-test environment is not defined yet.
- The future event delivery or outbox standard remains deferred.

## Acceptance Criteria

- [x] Define repository, transaction, migration, and ownership rules before persistent implementation.
- [x] Record the decision as an Agent Decision Record.
- [x] Update runtime and inventory guidance.
- [x] Add follow-up work items for the durable inventory milestone.
- [x] Keep this change documentation-only and standards-only.
