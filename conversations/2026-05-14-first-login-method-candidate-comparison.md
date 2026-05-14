# Conversation: First Login Method Candidate Comparison

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-compare-first-login-method-candidates/`

Related artifacts:

- `docs/first-login-method-candidates.md`
- `docs/first-login-method-candidates.zh-CN.md`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`

## Context

The maintainer asked the agent to continue through the work queue unless a truly necessary decision required human input.

After `W-0064` defined the login method and token format ratification standard, `W-0065` required comparing candidate first login methods and recommending a first login-method set for ratification.

## Maintainer Narrative

The maintainer authorized professional judgment for technical decisions and asked the agent to continue unless a decision genuinely needed maintainer input.

The maintainer also repeatedly asked the project to actively reference Nakama and Pitaya while preserving vibit's Agent-Native distinction.

## Agent Response Summary

The agent compared device credential login, guest/anonymous login, custom ID login, email/password login, external provider login, and service authentication.

The recommendation is to ratify `device_credential_login` as the first production player login method in W-0066, with all other candidate families deferred.

This recommendation adapts Nakama's low-friction device-auth capability without copying Nakama's API shape. It also preserves Pitaya-inspired separation between connection/session vocabulary and handler request identity, without putting authentication into WebSocket acceptors.

## Decisions

- Recommend `device_credential_login` as the first login-method set for W-0066.
- Treat the credential as high-entropy secret proof material, not raw OS device ID metadata.
- Recommend login command payload as the first proof exchange posture, leaving Protobuf envelope and WebSocket handshake behavior unchanged.
- Defer guest/anonymous, custom ID, email/password, external provider, and service-auth families.

## Artifacts

- `docs/first-login-method-candidates.md`
- `docs/first-login-method-candidates.zh-CN.md`
- `changes/2026-05-14-compare-first-login-method-candidates/`

## Open Questions

- Whether the initial credential is client-generated, server-issued during bootstrap, or both.
- Whether account creation is implicit or requires explicit create intent.
- Whether credential rotation and account recovery are included in the first implementation or deferred.
- Which token model will be selected in W-0067 and W-0068.

## Follow-Up

- Advance W-0066 to ratify the first login-method set.
- Do not implement authentication until later contract, schema, token, check, and test gates are complete.

## Redaction Notes

No secrets, tokens, account details, or private data are stored in this conversation log.
