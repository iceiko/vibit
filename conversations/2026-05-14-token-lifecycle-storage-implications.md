# Conversation: Token Lifecycle Storage Implications

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-define-token-lifecycle-and-storage-implications/`

Related artifacts:

- `docs/token-lifecycle-storage-implications.md`
- `docs/token-lifecycle-storage-implications.zh-CN.md`
- `decisions/ADR-0027-token-lifecycle-and-storage-implications.md`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`

## Context

The maintainer asked the agent to continue through the work queue unless a truly necessary decision required human input.

After W-0068 ratified opaque high-entropy access tokens and explicit proof carriers, W-0069 required defining issuance, expiration, refresh, revocation, rotation, replay, logout, cleanup, redaction, audit, and storage implications before implementation.

## Maintainer Narrative

The maintainer has authorized professional technical judgment within the selected milestone and has asked the project to preserve long-term production seriousness rather than treating the first authentication path as a demo.

## Agent Response Summary

The agent defined a first lifecycle posture for opaque high-entropy access tokens.

The posture requires at least 256 bits of entropy, a one-hour access-token TTL, no refresh token in the first implementation, renewal by reauthenticating with the selected login method, revocation, rotation, replay controls, logout, cleanup, audit, raw-token redaction, and future non-plaintext verifier storage.

The posture keeps credential storage and token verifier storage behind future schema gates, defaults future token storage planning to PostgreSQL, and leaves session storage, external identity storage, Redis-like stores, Protobuf envelope changes, WebSocket handshake authentication, and runtime implementation deferred.

## Decisions

- Set the first access-token TTL to `1h`.
- Require at least 256 bits of token entropy.
- Exclude refresh tokens from the first implementation posture.
- Use reauthentication with `device_credential_login` as the first renewal method.
- Require revocation, rotation, logout, cleanup, audit, redaction, and replay controls before implementation.
- Require future non-plaintext token verifier storage.
- Use PostgreSQL as the default future durable target for token verifier storage, behind W-0071 schema gates.
- Keep session storage and external identity storage out of the first posture.
- Preserve player account lifecycle tables as credential-free, token-free, external-identity-free, and session-free.

## Artifacts

- `docs/token-lifecycle-storage-implications.md`
- `docs/token-lifecycle-storage-implications.zh-CN.md`
- `decisions/ADR-0027-token-lifecycle-and-storage-implications.md`
- `changes/2026-05-14-define-token-lifecycle-and-storage-implications/`

## Open Questions

- Exact future schema fields for credential and token verifier records.
- Exact cleanup mechanism and retention enforcement.
- Exact audit event catalog.
- Whether a later milestone introduces refresh tokens or persisted sessions.

## Follow-Up

- Advance W-0070 to define authentication contract, error, permission, and audit surfaces.
- Do not implement token lifecycle behavior until contract, schema, check, and test gates are complete.

## Redaction Notes

No secrets, tokens, account details, or private data are stored in this conversation log.
