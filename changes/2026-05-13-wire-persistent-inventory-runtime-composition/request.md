# Request

## Original Request

The maintainer asked to continue development. The current next-ready work item is `W-0020: Wire persistent inventory runtime composition`.

## Clarified Requirement

Add a persistent inventory runtime composition path that binds the existing WebSocket Protobuf request loop to PostgreSQL inventory persistence through the application-owned unit-of-work boundary.

The persistent path must be explicit and opt-in. It must not make PostgreSQL mandatory for default local startup, remove the in-memory bootstrap path, introduce authentication or player account ownership, or replace the inventory account row-lock model.

## User-Visible Outcome

The runtime can still start with the existing in-memory inventory repository by default.

When explicitly configured for PostgreSQL, the runtime opens a PostgreSQL pool, creates a pgx-backed unit-of-work runner, wires inventory command handling through transaction-bound PostgreSQL repositories, and serves the same `/v1/ws` WebSocket Protobuf endpoint.

## Non-Goals

- Do not apply migrations automatically during normal server startup.
- Do not make PostgreSQL mandatory for default local development startup.
- Do not remove the in-memory bootstrap path.
- Do not introduce authentication, session validation, player account ownership, item catalog, outbox delivery, or event publication.
- Do not change the WebSocket Protobuf envelope shape.
- Do not change inventory command, query, event, permission, or error contracts.
- Do not replace the inventory account row-lock model.

## Unknowns

- Live PostgreSQL end-to-end verification remains a later work item.
- Migration apply/status behavior against a live disposable PostgreSQL environment remains a later work item.
- Production permission policy and authentication remain deferred.

## Acceptance Criteria

- Add an application-level composition point that can create inventory repositories from the current unit of work for command routes.
- Keep application and domain packages free of `pgx`, generated Protobuf, WebSocket, and migration dependencies.
- Add a PostgreSQL runtime store selection path that is opt-in through explicit process configuration.
- Keep the in-memory inventory runtime as the default path.
- Add focused tests proving command routes use the unit-of-work repository while query routes do not open a write transaction.
- Update runtime runbook, runtime agent guidance, architecture manifests, and inventory module metadata.
- Run Go and repository verification.
