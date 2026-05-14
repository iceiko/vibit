# Conversation: Player Account PostgreSQL Persistence Direction

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-confirm-next-direction-after-tooling-hardening/`
- `changes/2026-05-14-define-player-account-postgresql-persistence-schema-boundary/`

Related artifacts:

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `docs/postgresql-persistence-boundary.md`
- `docs/player-account-session-contracts.md`
- `modules/player/module.yaml`

## Context

After `M-007 Agent Tooling And Generator Hardening` completed, the project stopped at `M-008 Next Direction Confirmation Gate`. The agent recommended `ratify_player_account_postgresql_persistence_schema` as the next conservative direction because player account semantic and wire contracts already exist, while runtime player handlers, authentication, token behavior, credentials, and session persistence remain deferred.

## Maintainer Narrative

The maintainer authorized the professional recommendation:

```text
按照你的专业建议进行推进。
```

## Agent Response Summary

The agent treated this as authorization to close the direction gate and open a bounded milestone for player account PostgreSQL persistence schema ratification.

The agent kept the work limited to schema boundary preparation first. This does not implement authentication, token behavior, credential storage, session persistence, WebSocket handshake authentication, runtime player handlers, WebSocket routes, direct Nakama/Pitaya API compatibility, or a major external dependency.

## Decisions

- Select `ratify_player_account_postgresql_persistence_schema` as the next major direction.
- Create `M-009 Player Account PostgreSQL Persistence Schema`.
- Start the milestone with a schema boundary standard before adding migration source or repository implementation.
- Continue using Nakama and Pitaya as reference baselines for capability coverage and framework vocabulary, not as public API shapes.
- Preserve the rule that metadata-only `player_id` and `session_id` are not authenticated proof.

## Artifacts

- Added this conversation log.
- Added the `confirm-next-direction-after-tooling-hardening` change spec.
- Planned player account PostgreSQL persistence schema boundary work under `M-009`.

## Open Questions

- Which later milestone should choose authentication and token/session validation behavior.
- Whether player account creation will initially be exposed through runtime routes or reserved for an internal/system actor after persistence exists.
- Whether future account identity linking needs a separate table before any login method is chosen.

## Follow-Up

- Execute `W-0050 Define player account PostgreSQL persistence schema boundary`.
- Keep migration source, repository adapter, runtime handlers, authentication, tokens, credentials, and session persistence as separate bounded work.

## Redaction Notes

No secrets, tokens, account identifiers, or private data were included in this conversation log.
