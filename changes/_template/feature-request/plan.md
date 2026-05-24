# Plan

## Intake

- [ ] Preserve the original request in `request.md`.
- [ ] Clarify user-visible outcome, non-goals, unknowns, and acceptance criteria.
- [ ] Map the request to a Nakama capability family or record `no_mapping_applies`.
- [ ] Confirm Pitaya remains deferred unless a later ADR explicitly reactivates it.

## Files To Create

- None yet.

## Files To Edit

- None yet.

## Contracts And Schemas

- Record commands, queries, events, errors, permissions, routes, protocol payloads, generated output, and migrations before implementation when applicable.

## Implementation Boundary

- Allowed: Replace with bounded file, package, module, or artifact areas.
- Forbidden: Runtime, protocol, generated output, migration, dependency, SDK, hosted, distributed runtime, or direct compatibility scope not explicitly authorized by this change.

## Tests

- [ ] Positive behavior tests or not-applicable rationale.
- [ ] Negative behavior tests or not-applicable rationale.
- [ ] Permission/authentication tests or not-applicable rationale.
- [ ] Persistence/protocol/integration tests or not-applicable rationale.
- [ ] Repository checks.

## Verification Commands

- `node tools/vibit check change {{CHANGE_ID}} --json`
- `node tools/vibit check work --json`

## Rollback Or Migration Notes

Describe rollback, migration, or cleanup notes when relevant.
