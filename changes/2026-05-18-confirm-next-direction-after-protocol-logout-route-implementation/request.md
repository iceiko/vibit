# Request

The maintainer asked to replan vibit as a Nakama/Pitaya-class product that covers their common capability families, then continue development.

After implementing `W-0170`, the selected next direction is:

```text
define_transport_close_handoff_gate
```

This keeps the near-term roadmap on runtime lifecycle closure. The protocol logout route is now client-visible, and the application close policy can already produce redacted `CloseIntent` values with `mark_invalidated_only`, but WebSocket transport still has no narrow handoff for concrete socket close.

## Rationale

Nakama shows that logout, session invalidation, realtime socket disconnect, and operator-initiated disconnect are related product lifecycle capabilities but should have explicit semantics.

Pitaya shows that connection acceptors, sessions, handlers, groups, and kick/disconnect behavior should remain separate framework surfaces.

For vibit, the next correct move is not presence, chat, groups, matchmaking, or match runtime yet. Those systems need stable close, reconnect, and session-carrier behavior underneath them. The next work should define the transport close handoff gate so later implementation can close concrete WebSocket sockets without letting transport own authentication, route policy, session revocation, or business decisions.

## Non-Goals

- Do not implement concrete socket close handoff in this confirmation step.
- Do not choose WebSocket close codes or close reason text.
- Do not make logout automatically close sockets.
- Do not revoke runtime sessions.
- Do not invalidate active connection records beyond existing close policy behavior.
- Do not add reconnect, resume, duplicate replacement, or connection epoch behavior.
- Do not add protocol session carriers or change the existing Protobuf envelope.
- Do not add presence, chat, friends, groups, parties, matchmaking, match runtime, SDKs, cluster, or distributed runtime behavior.
- Do not add dependencies.
- Do not add direct Nakama/Pitaya API compatibility.

## Acceptance Criteria

- [x] `define_transport_close_handoff_gate` is selected.
- [x] The selection is recorded as the continuation after `W-0170`.
- [x] `M-100/W-0172` is created as the next bounded gate-only work item.
- [x] Runtime session revocation, reconnect/epoch, protocol session carriers, presence, operations, dependencies, direct compatibility, and broad product modules remain deferred.
