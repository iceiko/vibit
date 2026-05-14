# Conversation: Authentication Token Session Validation Direction

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-close-player-account-postgresql-persistence-milestone/`
- `changes/2026-05-14-confirm-next-direction-after-player-account-postgresql-persistence/`

Related artifacts:

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `docs/player-identity-session-boundary.md`
- `docs/player-account-session-contracts.md`

## Context

After `M-009 Player Account PostgreSQL Persistence Schema` closed, the work queue stopped at `M-010 Next Direction Confirmation Gate`. The agent recommended option 1, `ratify_authentication_and_token_session_validation_design`, because player account persistence exists but runtime player handlers would immediately need authenticated identity, token/session semantics, and production permission boundaries.

## Maintainer Narrative

The maintainer selected option 1:

```text
1
```

## Agent Response Summary

The agent interpreted this as a maintainer decision to choose `ratify_authentication_and_token_session_validation_design` as the next major direction.

The agent recorded the decision as a direction-gate closeout and planned a new bounded milestone. The next milestone should begin with authentication and token/session validation design standards, not production authentication implementation.

The agent preserved ask-first boundaries for concrete login methods, token format, token signing/refresh/expiration/revocation, credential storage, external identity linking, session persistence, Protobuf envelope behavior, WebSocket handshake authentication, runtime player handlers, WebSocket routes, direct Nakama/Pitaya API compatibility, and metadata-only identity trust.

## Decisions

- Select `ratify_authentication_and_token_session_validation_design` as the next major direction.
- Create `M-011 Authentication And Token Session Validation Design`.
- Start with a design standard and reference mapping before implementation.
- Continue using Nakama and Pitaya as references for account/auth/session capability coverage and session-binding vocabulary, not as public API shapes.
- Preserve the rule that metadata-only `player_id` and `session_id` are not authenticated proof.

## Artifacts

- Added this conversation log.
- Added the `confirm-next-direction-after-player-account-postgresql-persistence` change spec.
- Planned authentication and token/session validation design work under `M-011`.

## Open Questions

- Which login methods should vibit support first.
- Whether the first token model should be opaque server-issued tokens, signed tokens, or another explicit format.
- Whether credential storage, external identity linking, token storage, and session persistence should be separate milestones or a single tightly scoped milestone.
- Whether WebSocket handshake authentication should be introduced before or after request-level validation behavior.

## Follow-Up

- Execute `W-0057 Define authentication and token session validation design standard`.
- Keep implementation, dependency adoption, migrations, Protobuf envelope changes, WebSocket handshake changes, runtime player handlers, and WebSocket routes separate until design is ratified.

## Redaction Notes

No secrets, tokens, account identifiers, or private data were included in this conversation log.
