# Conversation: Login Method Token Format Ratification Standard

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-define-login-method-token-format-ratification-standard/`

Related artifacts:

- `docs/login-method-token-format-ratification.md`
- `docs/login-method-token-format-ratification.zh-CN.md`
- `decisions/ADR-0024-login-method-token-format-ratification-boundary.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`

## Context

After the authentication/token/session validation design milestone closed, the maintainer authorized the agent's professional recommendation to select `login_method_and_token_format_ratification` as the next milestone direction.

The next ready work item was `W-0064`, which required defining the standard for ratifying first login methods, token model, proof carrier boundaries, lifecycle semantics, storage implications, and implementation gates before code.

## Maintainer Narrative

The maintainer instructed the agent to keep advancing according to professional judgment:

> 按照你的建议和判断推进10步，除非有非常必要的，需要我决策的，再停下来问。

The maintainer had also previously emphasized that unnecessary confirmations should be avoided, but genuinely branching product or architecture decisions should be discussed.

## Agent Response Summary

The agent treated `W-0064` as a standards and queue-building step, not an implementation step.

The agent verified current Nakama and Pitaya reference material and used it as capability and vocabulary input only:

- Nakama remains the reference for authentication methods, account/session capability coverage, session token, refresh, expiration, logout, and realtime socket/session concepts.
- Pitaya remains the reference for handler context, session binding, frontend/backend session vocabulary, routing, push, and session lifecycle vocabulary.
- Neither reference governs vibit's public API shape.

The agent added a dedicated ratification standard and ADR before selecting concrete login methods or token formats.

## Decisions

- Create `docs/login-method-token-format-ratification.md` and the paired Simplified Chinese translation.
- Create `ADR-0024` for the login method and token format ratification boundary.
- Plan M-013 as bounded comparison, ratification, lifecycle, contract, schema, check, and milestone-close steps.
- Keep runtime authentication, token parsing, credential storage, external identity linking, session persistence, Protobuf envelope changes, WebSocket handshake authentication, runtime player handlers, and WebSocket routes unimplemented in this step.

## Artifacts

- `docs/login-method-token-format-ratification.md`
- `docs/login-method-token-format-ratification.zh-CN.md`
- `decisions/ADR-0024-login-method-token-format-ratification-boundary.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`

## Open Questions

- Which login methods should be selected as the first production set.
- Which token format should be selected first.
- Whether the first proof carrier should stay inside a login command payload, use a future protocol-level system message, or require a future envelope/handshake decision.
- Which storage and migration gates are required before token/session implementation.

## Follow-Up

- Advance `W-0065` to compare first login-method candidates.
- Continue M-013 without adding implementation code until the required ratification steps are complete.

## Redaction Notes

No secrets, tokens, account details, or private data are stored in this conversation log.
