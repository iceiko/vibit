# Conversation: Authentication Login Route Next Direction And Session Handshake Ratification

Date: 2026-05-17
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-authentication-login-route/`
- `changes/2026-05-17-ratify-session-persistence-websocket-handshake-authentication/`

Related artifacts:

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `docs/session-persistence-websocket-handshake-ratification.md`
- `docs/session-persistence-websocket-handshake-ratification.zh-CN.md`
- `decisions/ADR-0056-session-persistence-websocket-handshake-ratification.md`
- `tools/vibit`
- `rules/check-rules.json`

## Context

The public authentication login route was implemented in `M-048/W-0120`. The work queue was blocked at `M-049/W-0121`, asking for the next major direction after clients gained an in-band device credential login command route.

## Maintainer Narrative

The maintainer asked:

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。
```

The maintainer previously emphasized that Nakama and Pitaya should be treated as important references because mature game servers understand current game backend needs deeply.

## Agent Response Summary

The agent recommended selecting:

```text
ratify_session_persistence_and_websocket_handshake_authentication
```

The reasoning was that vibit now has public login and request-level route protection. The next mature game-server concern is the login-to-session-to-socket lifecycle: Nakama provides a broad authentication/session/socket capability reference, while Pitaya provides session binding and route/handler layering vocabulary.

The agent then defined a bounded standards gate rather than implementing session tables or WebSocket handshake authentication. The ratified posture keeps the current production path as request-level opaque access-token validation through the Protobuf authenticated payload wrapper, keeps WebSocket transport credential-neutral, keeps the existing envelope unchanged, and plans future connection binding through a protocol/application first-message binding gate.

## Decisions

- Close `M-049/W-0121`.
- Select `ratify_session_persistence_and_websocket_handshake_authentication`.
- Create `M-050/W-0122` as the ratification gate.
- Accept `ADR-0056`.
- Add `runtime.session_persistence_websocket_handshake_ratification` as a repository check rule.
- Defer session tables, migrations, repository interfaces, PostgreSQL adapters, dependencies, WebSocket handshake proof carriers, logout, refresh, cleanup, token rotation, and direct Nakama/Pitaya API compatibility.

## Artifacts

- `changes/2026-05-17-confirm-next-direction-after-authentication-login-route/`
- `changes/2026-05-17-ratify-session-persistence-websocket-handshake-authentication/`
- `docs/session-persistence-websocket-handshake-ratification.md`
- `docs/session-persistence-websocket-handshake-ratification.zh-CN.md`
- `decisions/ADR-0056-session-persistence-websocket-handshake-ratification.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Open Questions

- The exact first-message system/protocol message shape is deferred.
- Session persistence schema, repository interface, migration source, and cleanup behavior are deferred.
- Logout, refresh, cleanup, token rotation, and active-connection invalidation remain separate lifecycle gates.
- Whether WebSocket handshake authentication should ever parse header, cookie, query, or subprotocol carriers remains deferred and currently not selected.

## Follow-Up

- Define the first-message connection binding protocol gate if connection-bound session behavior becomes the next milestone.
- Define the PostgreSQL session persistence schema gate before adding any session table or repository behavior.
- Keep Nakama and Pitaya as references without copying their public APIs.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, or GitHub tokens are recorded in this conversation log.
