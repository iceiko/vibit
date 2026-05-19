# Conversation: Runtime Authentication Startup Composition Next Direction

Date: 2026-05-17
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-runtime-authentication-startup-composition/`

Related artifacts:

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/reference.yaml`

## Context

The work queue was blocked at `M-046/W-0118` after runtime authentication startup composition was implemented.

## Maintainer Narrative

The maintainer asked:

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。
```

## Agent Response Summary

The agent recommended selecting:

```text
add_authentication_command_protocol_messages_and_login_route_registration
```

The rationale was that the runtime can already validate access tokens for protected routes and startup composition now injects route protection, but clients still lack a formal in-band protocol route to perform device credential login and receive an opaque access token.

## Decisions

- Close `M-046/W-0118`.
- Select `add_authentication_command_protocol_messages_and_login_route_registration`.
- Create `M-047/W-0119` as the authentication command protocol and login route gate.

## Artifacts

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/reference.yaml`
- `changes/2026-05-17-confirm-next-direction-after-runtime-authentication-startup-composition/`

## Open Questions

- The implementation slice still needs to define exact Protobuf messages, generated output, bridge code, route handler registration, and startup composition.
- Session persistence, WebSocket handshake authentication, logout, refresh, cleanup, token rotation, operations posture, and direct compatibility posture remain separately gated.

## Follow-Up

- Define the authentication command protocol and login route gate.
- Keep Nakama and Pitaya as references for authenticate-before-gameplay sequencing and transport/session/route/handler separation without copying their public APIs.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, or GitHub tokens are recorded in this conversation log.
