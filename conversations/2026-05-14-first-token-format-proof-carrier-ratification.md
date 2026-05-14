# Conversation: First Token Format Proof Carrier Ratification

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-ratify-first-token-format-and-proof-carrier-posture/`

Related artifacts:

- `docs/first-token-format-proof-carrier-posture.md`
- `docs/first-token-format-proof-carrier-posture.zh-CN.md`
- `decisions/ADR-0026-first-token-format-and-proof-carrier-posture.md`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`

## Context

The maintainer asked the agent to advance through the work queue unless a truly necessary decision required human input.

After `W-0067` compared token format and carrier options, `W-0068` required ratifying the first token format and proof carrier posture with issuer, verifier, subject, audience, expiration, refresh, revocation, rotation, replay, redaction, and storage implications.

## Maintainer Narrative

The maintainer authorized professional technical judgment within the selected milestone, asked the project to continue referencing Nakama and Pitaya, and emphasized that routine technical details should not stop progress unless they require a real product or architecture decision.

## Agent Response Summary

The agent ratified `opaque_high_entropy_token` as the first access-token format.

The selected issuance carrier is a future login command response. The selected request proof carrier is a future explicit request proof payload field.

The ratification keeps refresh tokens, session token vocabulary, Protobuf envelope changes, WebSocket handshake authentication, first system-message binding, token lifecycle details, storage schema, and implementation deferred.

Current Protobuf session metadata remains metadata-only and must not become authenticated proof.

## Decisions

- Ratify `opaque_high_entropy_token` as the first access-token format.
- Ratify `login_command_response_token` as the token issuance carrier.
- Ratify `explicit_request_proof_payload` as the first request proof carrier.
- Require finite expiration, redaction, and non-plaintext verifier storage before implementation.
- Defer exact expiration, refresh, revocation, rotation, replay, logout, cleanup, audit, and storage policy to W-0069.
- Defer refresh tokens and session token vocabulary.
- Preserve Protobuf envelope and WebSocket handshake behavior unchanged.
- Reject current `Session` metadata as proof.

## Artifacts

- `docs/first-token-format-proof-carrier-posture.md`
- `docs/first-token-format-proof-carrier-posture.zh-CN.md`
- `decisions/ADR-0026-first-token-format-and-proof-carrier-posture.md`
- `changes/2026-05-14-ratify-first-token-format-and-proof-carrier-posture/`

## Open Questions

- Exact access-token expiration duration.
- Whether refresh tokens are part of the first implementation.
- Token revocation, rotation, replay, logout, cleanup, and audit policy.
- Token verifier storage schema.
- Future contract shape for login responses and authenticated request proof fields.

## Follow-Up

- Advance W-0069 to define token lifecycle and storage implications.
- Do not implement authentication until lifecycle, contract, schema, check, and test gates are complete.

## Redaction Notes

No secrets, tokens, account details, or private data are stored in this conversation log.
