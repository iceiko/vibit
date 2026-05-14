# Conversation: Login Method Token Format Direction

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-close-authentication-token-session-validation-design-milestone/`
- `changes/2026-05-14-confirm-next-direction-after-authentication-design/`

Related artifacts:

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `docs/authentication-token-session-validation.md`
- `docs/authentication-proof-token-session-contract-dimensions.md`
- `docs/credential-storage-external-identity-linking-boundaries.md`
- `docs/session-persistence-websocket-handshake-decision-gates.md`

## Context

After `M-011 Authentication And Token Session Validation Design` closed, the work queue stopped at `M-012 Next Direction Confirmation Gate`. The agent recommended `login_method_and_token_format_ratification` because player account persistence and authentication/session boundaries exist, but runtime player routes would still lack a ratified production identity basis.

## Maintainer Narrative

The maintainer authorized the agent's recommendation and asked the agent to continue up to ten steps unless a truly necessary decision point appears:

```text
按照你的建议和判断推进10步，除非有非常必要的，需要我决策的，再停下来问。
```

## Agent Response Summary

The agent interpreted this as authorization to choose `login_method_and_token_format_ratification` from the `W-0063` candidate directions.

The agent recorded the decision as a direction-gate closeout and created `M-013 Login Method And Token Format Ratification` with a bounded first work item. The next milestone should start with a ratification standard and reference mapping before implementation.

The agent preserved ask-first boundaries for production authentication implementation, concrete login handlers, token parsing, credential storage, external identity linking, session persistence, Protobuf envelope changes, WebSocket handshake authentication, runtime player handlers, WebSocket routes, major dependencies, direct Nakama/Pitaya API compatibility, and metadata-only identity trust.

## Decisions

- Select `login_method_and_token_format_ratification` as the next major direction.
- Create `M-013 Login Method And Token Format Ratification`.
- Start with a ratification standard before implementation.
- Continue using Nakama as the authentication/session-token capability reference and Pitaya as the session-binding vocabulary reference.
- Preserve the rule that metadata-only `player_id` and `session_id` are not authenticated proof.

## Artifacts

- Added this conversation log.
- Added the `confirm-next-direction-after-authentication-design` change spec.
- Planned login method and token format ratification work under `M-013`.

## Open Questions

- Which login methods should vibit ratify first.
- Whether the first token model should be opaque, signed, structured, split access/refresh, or another explicit model.
- Whether token/proof carrier behavior should remain request-level first or involve first-message or handshake gates.
- Whether credential storage, external identity linking, token storage, and session persistence should remain separate milestones or become a tightly staged implementation sequence.

## Follow-Up

- Execute `W-0064 Define login method and token format ratification standard`.
- Keep implementation, dependency adoption, migrations, Protobuf envelope changes, WebSocket handshake changes, runtime player handlers, and WebSocket routes separate until the ratification standard and follow-up work queue are explicit.

## Redaction Notes

No secrets, tokens, account identifiers, or private data were included in this conversation log.
