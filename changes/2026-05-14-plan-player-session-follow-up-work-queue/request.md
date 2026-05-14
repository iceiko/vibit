# Request

## Original Request

Continue advancing the project work queue. The active next-ready work item is `W-0025 Plan player/session follow-up work queue`.

## Clarified Requirement

After the application request identity handoff exists, restore deterministic continuation for the rest of `M-003 Player Identity And Session Boundary` by adding bounded follow-up work items in dependency order.

The planned queue must keep identity/session work moving without crossing any deferred authentication, token, credential, persistence, Protobuf envelope, or WebSocket handshake decision.

## User-Visible Outcome

Future `continue` / `继续` requests can advance through player/session boundary work one work item at a time, with clear stop conditions and ask-first boundaries.

## Non-Goals

- Do not choose an authentication scheme.
- Do not choose a token format.
- Do not add credential storage.
- Do not add player account database schema or migrations.
- Do not change the Protobuf envelope.
- Do not change the WebSocket handshake.
- Do not make inventory depend directly on the player module.
- Do not claim metadata-only identity is authenticated proof.

## Unknowns

- The future production authentication model remains undecided.
- The future session persistence and expiration model remains undecided.
- The future player account schema remains undecided.
- The future API surface for real login and session renewal remains undecided.

## Acceptance Criteria

- `.arch/work-items.yaml` contains the remaining conservative `M-003` work items.
- Exactly one work item is `next_ready`.
- Each planned work item has dependencies, completion criteria, and ask-first boundaries.
- The queue separates vocabulary, runtime hook, inventory permission migration, repository checks, and milestone closure.
- Repository work checks pass.
