# Request

## Original Request

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。
```

## Clarified Requirement

Recommend the next milestone direction after the public authentication login route and continue with the recommended sequence.

The selected direction is:

```text
ratify_session_persistence_and_websocket_handshake_authentication
```

## User-Visible Outcome

The work queue records that the next milestone after the public login route is session persistence and WebSocket handshake authentication ratification, with Nakama and Pitaya used as reference baselines.

## Recommendation

The recommended next ten-step sequence is:

1. Close `M-049/W-0121` by selecting the next direction.
2. Create a bounded `M-050/W-0122` gate before implementation.
3. Define the current production path as request-level validation through the existing authenticated payload wrapper.
4. Keep WebSocket transport credential-neutral.
5. Preserve the existing Protobuf envelope without adding token or session proof fields.
6. Ratify future first-message connection binding as the preferred next connection-level gate.
7. Defer WebSocket handshake header, cookie, query, and subprotocol credential carriers.
8. Defer session persistence schema, repository interfaces, migrations, and dependencies.
9. Record Nakama and Pitaya reference mapping without adopting direct API compatibility.
10. Add repository checks and verification for the new gate.

## Non-Goals

- Do not implement session persistence in the direction-confirmation step.
- Do not implement WebSocket handshake authentication.
- Do not add WebSocket `Authorization`, Bearer, cookie, query-string, or subprotocol proof carriers.
- Do not change the existing Protobuf envelope.
- Do not add session tables, migrations, repository interfaces, PostgreSQL adapter behavior, dependencies, logout, refresh, cleanup, or token rotation.
- Do not adopt direct Nakama or Pitaya public API compatibility.

## Acceptance Criteria

- [x] `W-0121` is marked completed.
- [x] The selected direction is recorded as `ratify_session_persistence_and_websocket_handshake_authentication`.
- [x] `M-050/W-0122` is created as the next gate.
- [x] Ask-first boundaries remain recorded for implementation of session persistence, WebSocket handshake authentication, repositories, migrations, dependencies, logout, refresh, cleanup, token rotation, operations posture, memory durable authentication behavior, and direct Nakama/Pitaya API compatibility.
- [x] Verification is recorded.
