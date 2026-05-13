# Request

## Original Request

The maintainer asked to continue development after the previous PostgreSQL inventory repository adapter work was completed.

## Clarified Requirement

The work queue currently has no `next_ready` work item, so a continuation request cannot advance deterministically. Restore the continuation mechanism by adding the next bounded M-002 work items in dependency order and marking the next conservative step as `next_ready`.

## User-Visible Outcome

Future `continue` / `继续` requests have a concrete next work item again.

## Non-Goals

- Do not implement runtime persistence wiring in this change.
- Do not introduce a new persistence, transaction, event delivery, protocol, or deployment architecture.
- Do not make live PostgreSQL integration tests mandatory yet.
- Do not close M-002 until durable runtime wiring and verification have explicit work items.

## Unknowns

- The exact operational standard for disposable PostgreSQL integration tests remains deferred.
- The outbox or durable event delivery standard remains deferred.
- The final M-002 completion point may need adjustment after persistent runtime wiring is implemented and verified.

## Acceptance Criteria

- `.arch/work-items.yaml` includes the next M-002 work items after W-0015.
- Exactly one work item is `next_ready`.
- New work items have dependencies, completion criteria, and ask-first boundaries where needed.
- Repository work checks pass.
