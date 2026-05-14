# Conversation: Authentication Token Session Validation Design Standard

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-confirm-next-direction-after-player-account-postgresql-persistence/`
- `changes/2026-05-14-define-authentication-token-session-validation-design-standard/`

Related artifacts:

- `docs/authentication-token-session-validation.md`
- `docs/authentication-token-session-validation.zh-CN.md`
- `decisions/ADR-0023-authentication-token-session-validation-design-boundary.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`

## Context

After player account PostgreSQL persistence was completed, the maintainer selected option `1` from the next direction gate. That direction is:

```text
ratify_authentication_and_token_session_validation_design
```

The active work item became `W-0057 Define authentication and token session validation design standard`.

## Maintainer Narrative

The maintainer selected:

```text
1
```

This selected the authentication and token/session validation design direction before runtime player handlers, WebSocket routes, credential storage, external identity linking, or session persistence implementation.

## Agent Response Summary

The agent treated this as authorization to complete the next bounded design-standard step, not as permission to implement authentication.

The agent created the authentication/token/session validation design standard, mapped Nakama and Pitaya concepts into vibit positions, preserved metadata-only identity as non-authenticated, and planned the remaining M-011 work queue.

The agent did not choose a concrete login method, token format, credential store, external identity model, session persistence model, Protobuf envelope change, WebSocket handshake authentication behavior, runtime player handler, or WebSocket route.

## Decisions

- Ratify `ADR-0023` as the authentication/token/session validation design boundary.
- Keep metadata-only `player_id`, `session_id`, and `connection_id` non-authenticated.
- Treat Nakama session/token/realtime socket concepts as reference vocabulary and capability coverage.
- Treat Pitaya session binding, route handler, frontend/backend, and realtime session vocabulary as architecture reference vocabulary.
- Continue M-011 through bounded follow-up work instead of implementing authentication immediately.

## Artifacts

- `docs/authentication-token-session-validation.md`
- `docs/authentication-token-session-validation.zh-CN.md`
- `decisions/ADR-0023-authentication-token-session-validation-design-boundary.md`
- `changes/2026-05-14-define-authentication-token-session-validation-design-standard/`

## Open Questions

- Which login method should be implemented first?
- Which token model should be ratified first?
- Should first production validation be request-level, handshake-level, first-message, every-request, or hybrid?
- Should runtime sessions be persisted, and if so in which store?
- When should runtime player account handlers and WebSocket routes be added?

## Follow-Up

- Add architecture checks for the authentication/token/session boundary.
- Ratify semantic dimensions for authentication proof and token/session validation.
- Define credential and external identity storage boundaries.
- Define session persistence and WebSocket handshake decision gates.

## Redaction Notes

No secrets, tokens, account credentials, or private data were recorded in this conversation log.
