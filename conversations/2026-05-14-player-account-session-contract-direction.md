# Conversation: Player Account And Session Contract Direction

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-confirm-next-milestone-direction/`

Related artifacts:

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `docs/reference-game-server-alignment.md`
- `docs/player-identity-session-boundary.md`

## Context

After `M-003 Player Identity And Session Boundary` was completed, the work queue intentionally stopped at `M-004 Next Direction Confirmation Gate`. The agent recommended choosing player account and session contract ratification as the next direction instead of immediately implementing production authentication.

## Maintainer Narrative

The maintainer selected the first candidate direction and emphasized that Nakama and Pitaya should remain active references:

```text
选1 注意参考nakama pitaya
```

## Agent Response Summary

The agent treated this as confirmation to close the direction gate and create the next milestone around player account and session contracts.

The agent preserved the existing ask-first boundaries: this direction does not yet choose authentication scheme, token format, credential storage, player account schema, session persistence, Protobuf envelope changes, WebSocket handshake changes, major external dependencies, or production permission semantics.

## Decisions

- Select `ratify_player_account_and_session_contracts` as the next major direction.
- Create `M-005 Player Account And Session Contracts`.
- Start the milestone with a contract standard work item, not production authentication implementation.
- Continue using Nakama as a reference for account, user, authentication, and session capability coverage.
- Continue using Pitaya as a reference for session binding, route handler, frontend/backend, and realtime server vocabulary.
- Do not copy Nakama or Pitaya public APIs without a separate compatibility ADR.

## Artifacts

- Planned update to `.arch/work-items.yaml`.
- Planned update to `.arch/runtime.yaml`.
- Planned update to `.arch/reference.yaml`.
- Added this conversation log.
- Added the `confirm-next-milestone-direction` change spec.

## Open Questions

- Which account commands and queries should be ratified first.
- Whether `session` should remain an application/runtime contract family or eventually receive its own module manifest.
- When to ratify token format and session persistence.
- Whether player account creation should support only internal/system creation first or expose client-facing authentication routes later.

## Follow-Up

- Execute `W-0031 Define player account and session contract standard`.
- Use Nakama/Pitaya references during contract vocabulary design.
- Keep authentication implementation blocked until a later explicit decision.

## Redaction Notes

No secrets, tokens, account identifiers, or private data were included in this conversation log.
