# Request

The maintainer asked to replan vibit as a Nakama/Pitaya-class product and continue development. The current roadmap prioritizes runtime lifecycle closure before broad feature expansion.

Implement the bounded `W-0172` gate-only slice:

```text
define_transport_close_handoff_gate
```

## Clarified Requirement

Define the future boundary between application-owned close policy and WebSocket transport-owned concrete socket close mechanics.

The gate must not implement concrete close behavior. It must specify how a future implementation can hand a redacted server-owned target from application policy to transport while preserving authentication, route policy, session revocation, reconnect, protocol carriers, operations, and direct compatibility deferrals.

## Reference Rationale

Nakama is the product reference for explicit lifecycle surfaces such as authentication sessions, logout, realtime socket disconnect, and server-directed disconnect behavior.

Pitaya is the architecture reference for separating acceptors, sessions, handlers, groups/RPC, and kick/disconnect connection management.

vibit should adapt those lessons into an application-to-transport handoff, not a transport-owned policy engine.

## Acceptance Criteria

- [x] `docs/transport-close-handoff-gate.md` exists.
- [x] `docs/transport-close-handoff-gate.zh-CN.md` exists.
- [x] `ADR-0080` records the decision.
- [x] The first future handoff target is `connection_id_and_epoch`.
- [x] Client-supplied metadata is explicitly rejected as transport close authority.
- [x] WebSocket transport remains credential-neutral and policy-neutral.
- [x] Concrete socket close implementation, close codes, close reason text, logout-triggered close, runtime session revocation, reconnect/epoch behavior, protocol session carriers, dependencies, and direct Nakama/Pitaya API compatibility remain deferred.
