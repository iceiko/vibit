# Request

## Original Request

Continue ten work items unless maintainer confirmation is required.

## Clarified Requirement

Advance `W-0012` by adding vibit-owned unit-of-work interfaces and application orchestration points without binding domain modules or application handlers to PostgreSQL driver types.

## User-Visible Outcome

Maintainers and agents can see the first transaction boundary skeleton in code:

- `runtime/internal/platform/tx/` owns unit-of-work interfaces.
- `runtime/internal/app/` can wrap command dispatch in a unit of work.
- Query dispatch remains outside a write unit of work by default.
- Existing in-memory runtime dispatch still works without persistence wiring.

## Non-Goals

- Do not add PostgreSQL adapter implementation.
- Do not add SQL migration files.
- Do not introduce `pgx` imports into application or domain packages.
- Do not introduce an outbox or durable event delivery standard.
- Do not move the transaction boundary owner away from `runtime/internal/platform/tx` and `runtime/internal/app`.
- Do not require persistent runtime configuration yet.

## Unknowns

- Concrete PostgreSQL transaction begin/commit/rollback behavior is deferred until the PostgreSQL platform adapter exists.
- Durable event recording remains deferred until an explicit event storage or outbox decision exists.

## Acceptance Criteria

- [x] Add `platform/tx` unit-of-work interfaces with no third-party persistence dependency.
- [x] Add application-level transaction dispatch orchestration for command routes.
- [x] Leave query routes outside write transactions by default.
- [x] Keep existing dispatcher behavior and request-loop tests working.
- [x] Update runtime guidance and manifests.
- [x] Update runtime checks so application code may import only the transaction boundary package, not arbitrary platform adapters.
