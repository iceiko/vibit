# Conversation: First Login Method Set Ratification

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-ratify-first-login-method-set/`

Related artifacts:

- `docs/first-login-method-set.md`
- `docs/first-login-method-set.zh-CN.md`
- `decisions/ADR-0025-first-login-method-set.md`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`

## Context

The maintainer asked the agent to continue through the work queue unless a truly necessary decision required human input.

After `W-0065` compared first login-method candidates, `W-0066` required ratifying the first login-method set with rationale, rejected alternatives, decision weights, known gaps, and implementation gates.

## Maintainer Narrative

The maintainer authorized professional judgment for technical decisions within the selected milestone and repeatedly asked the project to reference Nakama and Pitaya while preserving vibit's Agent-Native distinction.

## Agent Response Summary

The agent ratified `device_credential_login` as the first login-method set.

The selected method is defined as proof of possession of a high-entropy installation credential. It is explicitly not a raw device ID, player ID, session ID, connection ID, or metadata-only proof.

The ratification preserves implementation deferral. It does not add handlers, credential tables, token parsing, session persistence, Protobuf envelope changes, WebSocket handshake authentication, runtime player routes, or WebSocket routes.

## Decisions

- Ratify `device_credential_login` as the only first login method.
- Classify it as production-capable only after required gates.
- Allow account creation and existing-account authentication only after account creation and credential lookup policies are ratified.
- Defer account linking, account recovery, anonymous upgrade, guest login, custom ID login, email/password login, provider login, and service authentication.
- Keep Nakama as a capability reference and Pitaya as session/context vocabulary, without copying either public API shape.

## Artifacts

- `docs/first-login-method-set.md`
- `docs/first-login-method-set.zh-CN.md`
- `decisions/ADR-0025-first-login-method-set.md`
- `changes/2026-05-14-ratify-first-login-method-set/`

## Open Questions

- Whether installation credentials are client-generated, server-issued, or both.
- Whether first login creates a player account by default or requires explicit create intent.
- Which token format and carrier will be selected.
- Whether refresh tokens and persisted runtime sessions exist in the first implementation.

## Follow-Up

- Advance W-0067 to compare token format and carrier options.
- Do not implement authentication until token, lifecycle, contract, schema, check, and test gates are complete.

## Redaction Notes

No secrets, tokens, account details, or private data are stored in this conversation log.
