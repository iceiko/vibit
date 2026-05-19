# Conversation: Request-Level Route Protection Next Direction

Date: 2026-05-17
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-request-level-route-protection/`

Related artifacts:

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/reference.yaml`

## Context

The work queue was blocked at `M-043/W-0115` after the request-level access-token payload wrapper and route-protection slice was implemented.

## Maintainer Narrative

The maintainer asked:

```text
选择 wire_runtime_authentication_startup_composition，继续推进十步，注意，重点参考。Nakama 和 Pitaya。这些已有的游戏Server，他们对现在游戏的需求理解得非常透彻，我们要重点参考。
```

The maintainer also asked for Chinese communication.

## Agent Response Summary

The agent treated the message as explicit authorization to close the blocked direction gate and select:

```text
wire_runtime_authentication_startup_composition
```

The agent kept Nakama and Pitaya as active references:

- Nakama guides the expectation that clients authenticate and use session/token context before normal gameplay and realtime features.
- Pitaya guides the separation between acceptor/connection handling, session context, routing, and handlers.

The selection starts with a gate-only milestone before startup composition implementation.

## Decisions

- Close `M-043/W-0115`.
- Select `wire_runtime_authentication_startup_composition`.
- Create `M-044` with `W-0116` as the startup composition gate work item.

## Artifacts

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/reference.yaml`
- `changes/2026-05-17-confirm-next-direction-after-request-level-route-protection/`

## Open Questions

- The next direction after startup composition still needs maintainer selection because likely candidates cross session persistence, WebSocket handshake authentication, authentication command routes, logout/refresh/cleanup, token rotation, operations, or broader game backend scope.

## Follow-Up

- Define the runtime authentication startup composition gate.
- Keep Nakama and Pitaya as references for capability sequencing and runtime layering without copying their public APIs.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, or GitHub tokens are recorded in this conversation log.
