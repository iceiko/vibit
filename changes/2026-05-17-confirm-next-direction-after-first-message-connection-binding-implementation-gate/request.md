# Request

The maintainer asked:

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya
```

## Clarified Requirement

Recommend the next direction after the first-message connection binding implementation gate, then continue according to that recommendation.

## Recommendation

Select:

```text
implement_first_message_connection_binding
```

## Basis

- `ADR-0057` selected the future `runtime.authentication.BindConnection` system route.
- `ADR-0058` defined the implementation boundary for this exact slice.
- Nakama's session/socket model supports the idea that authenticated session material precedes authenticated realtime socket use.
- Pitaya's acceptor/session/handler separation supports keeping WebSocket transport credential-neutral and binding identity in an application-owned layer.

## Non-Goals

- Do not select PostgreSQL session persistence schema in this direction confirmation.
- Do not add WebSocket handshake authentication.
- Do not add logout/revocation active-connection invalidation.
- Do not add reconnect or duplicate connection replacement behavior.
- Do not adopt direct Nakama or Pitaya public API compatibility.

## Acceptance Criteria

- [x] The selected direction is recorded.
- [x] `M-055/W-0127` is completed.
- [x] `M-056/W-0128` is created for the bounded implementation slice.
- [x] Ask-first boundaries remain recorded for session persistence, logout/revocation, reconnect/epoch, operations, direct external API compatibility, and broader game backend expansion.
