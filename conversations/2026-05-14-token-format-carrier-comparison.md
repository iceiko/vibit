# Conversation: Token Format Carrier Comparison

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-compare-token-format-and-carrier-options/`

Related artifacts:

- `docs/token-format-carrier-options.md`
- `docs/token-format-carrier-options.zh-CN.md`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`

## Context

The maintainer asked the agent to continue through the work queue unless a truly necessary decision required human input.

After `W-0066` ratified `device_credential_login` as the first login-method set, `W-0067` required comparing token format and proof carrier options before ratifying token posture in W-0068.

## Maintainer Narrative

The maintainer authorized professional judgment within the selected M-013 milestone and asked the project to continue actively referencing Nakama and Pitaya.

## Agent Response Summary

The agent compared opaque high-entropy tokens, signed structured tokens, external provider tokens, and plain session-ID-as-secret posture.

The recommendation is to use `opaque_high_entropy_token` as the first token format, issue it through a future login command response, and carry subsequent proof through explicit request payload fields until a later protocol or session-binding decision is ratified.

The recommendation rejects using current `Session.session_id` metadata as proof and defers Protobuf envelope changes, WebSocket handshake authentication, first system-message binding, refresh tokens, and direct Nakama/Pitaya compatibility.

## Decisions

- Recommend `opaque_high_entropy_token` for W-0068 ratification.
- Recommend `login_command_response_token` as the issuance carrier.
- Recommend `explicit_request_proof_payload` as the first authenticated-route proof carrier.
- Reject current Protobuf `Session` metadata as proof.
- Defer signed structured tokens, external provider tokens, plain session-ID-as-secret posture, envelope extension, first system-message binding, and WebSocket handshake carriers.

## Artifacts

- `docs/token-format-carrier-options.md`
- `docs/token-format-carrier-options.zh-CN.md`
- `changes/2026-05-14-compare-token-format-and-carrier-options/`

## Open Questions

- Whether the first opaque token is named access token, session token, or another vibit term.
- Whether token storage starts in PostgreSQL or waits for another session store decision.
- Whether refresh tokens exist in the first implementation.
- Whether explicit proof payload fields later become generated wrapper conventions.

## Follow-Up

- Advance W-0068 to ratify first token format and proof carrier posture.
- Do not implement token parsing, token storage, Protobuf changes, handshake authentication, runtime player handlers, or WebSocket routes.

## Redaction Notes

No secrets, tokens, account details, or private data are stored in this conversation log.
