# Request

## Original Request

The maintainer said "继续推进" after the Termux PostgreSQL live verification succeeded.

## Clarified Requirement

Close the completed durable inventory runtime milestone and restore deterministic continuation by planning the next bounded milestone and first work item. The next milestone should prepare the player identity and session boundary, because inventory now has durable `player_id` references but vibit has not yet defined player account, authentication, or session ownership.

## User-Visible Outcome

`node tools/vibit inspect work` should show `M-002` completed, `M-003` active, and exactly one `next_ready` work item for defining the player identity and session boundary.

## Non-Goals

- Do not implement player accounts, authentication, login, tokens, or session validation.
- Do not change the WebSocket Protobuf envelope shape.
- Do not add database migrations for player identity.
- Do not introduce external auth providers or cryptography dependencies.
- Do not make PostgreSQL mandatory for default local tests.

## Unknowns

- The exact player account model, credential model, guest identity model, token model, and session lifecycle remain open for the next work item.
- Whether the first player/session implementation should include a persistent player module, a runtime-only session module, or both remains open.

## Acceptance Criteria

- [x] `M-002` is marked completed with a completion summary grounded in live PostgreSQL verification.
- [x] `M-003` is added and marked active.
- [x] The first `M-003` work item is marked `next_ready`.
- [x] The new work item has ask-first boundaries for auth, token, session, protocol, dependency, and migration decisions.
- [x] Relevant architecture manifests no longer describe durable inventory live verification as pending.
